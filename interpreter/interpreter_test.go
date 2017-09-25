package interpreter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/dsl"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/ThoughtWorksStudios/bobcat/generator"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func AssertShouldHaveField(t *testing.T, entity *generator.Generator, fieldName string) {
	emitter := NewDummyEmitter()
	result := entity.One(nil, emitter, NewRootScope())
	AssertNotNil(t, result[fieldName], "Expected entity to have field %s, but it did not", fieldName)
}

func AssertFieldYieldsValue(t *testing.T, entity *generator.Generator, field *Node) {
	emitter := NewDummyEmitter()
	result := entity.One(nil, emitter, NewRootScope())
	AssertEqual(t, field.ValNode().Value, result[field.Name])
}

var validFields = NodeSet{
	Field("name", Builtin(STRING_TYPE), IntArgs(10)...),
	Field("age", Builtin(INT_TYPE), IntArgs(1, 10)...),
	Field("weight", Builtin(FLOAT_TYPE), FloatArgs(1.0, 200.0)...),
	Field("dob", Builtin(DATE_TYPE), DateArgs("2015-01-01", "2017-01-01")...),
	Field("last_name", Builtin(DICT_TYPE), StringArgs("last_name")...),
	Field("status", Builtin(ENUM_TYPE), NodeSet{StringCollection("enabled", "disabled")}...),
	Field("status", Builtin(SERIAL_TYPE)),
	Field("catch_phrase", StringVal("Grass.... Tastes bad")),
}

var nestedFields = NodeSet{
	Field("name", Builtin(STRING_TYPE), IntArgs(10)...),
	Field("pet", Id("Goat")),
	Field("friend", Entity("Horse", validFields)),
}

var overridenFields = NodeSet{
	Field("catch_phrase", StringVal("Grass.... Tastes good")),
}

func interp() *Interpreter {
	emitter := NewDummyEmitter()
	return New(emitter, false)
}

func TestScopingResolvesOtherEntities(t *testing.T) {
	scope := NewRootScope()
	i := interp()
	node := Root(Entity("person", NodeSet{
		Field("pet", Entity("kitteh", overridenFields)),
		Field("pets_can_have_pets_too", Entity("lolcat", NodeSet{
			Field("cheezburgrz", StringVal("can has")),
			Field("protoype", Id("kitteh")),
		})),
	}))
	_, err := i.Visit(node, scope, false)
	AssertNil(t, err, "`lolcat` should be able to resolve `kitteh` because it lives within the scope hierarchy. error was %v", err)

	// using same root scope to simulate previously defined symbols
	_, err = i.Visit(Root(Generation(2, Id("person"))), scope, false)
	AssertNil(t, err, "Should be able to resolve `person` because it is defined in the root scope. error was %v", err)

	// using same root scope to simulate previously defined symbols; here, `kitteh` was defined in a child scope of `person`,
	// but not at the root scope, so we should not be able to resolve it.
	_, err = i.Visit(Root(Generation(1, Id("kitteh"))), scope, false)
	ExpectsError(t, "Cannot resolve symbol \"kitteh\"", err)
}

func TestValidVisit(t *testing.T) {
	node := Root(Entity("person", validFields), Generation(2, Id("person")))
	i := interp()
	scope := NewRootScope()
	_, err := i.Visit(node, scope, false)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	for _, entry := range scope.Symbols {
		entity := entry.(*generator.Generator)
		for _, field := range validFields {
			AssertShouldHaveField(t, entity, field.Name)
		}
	}
}

func TestValidVisitWithNesting(t *testing.T) {
	node := Root(Entity("Goat", validFields), Entity("person", nestedFields),
		Generation(2, Id("person")))
	i := interp()

	scope := NewRootScope()
	_, err := i.Visit(node, scope, false)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	person, _ := i.ResolveEntityFromNode(Id("person"), scope)
	for _, field := range nestedFields {
		AssertShouldHaveField(t, person, field.Name)
	}
}

func TestValidVisitWithOverrides(t *testing.T) {
	node := Root(
		Entity("person", validFields),
		Generation(
			2,
			EntityExtension("lazyPerson", "person", overridenFields),
		),
	)

	i := interp()
	scope := NewRootScope()

	if _, err := i.Visit(node, scope, false); err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	AssertEqual(t, 2, len(scope.Symbols), "Should have 2 entities defined")

	for _, key := range []string{"person", "lazyPerson"} {
		_, isPresent := scope.Symbols[key]
		// don't try to use AssertNotNil here; it won't work because it is unable to detect
		// whether a nil pointer passed as an interface{} param to AssertNotEqual is nil.
		// see this crazy shit: https://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
		Assert(t, isPresent, "`%v` should be defined in scope", key)

		if isPresent {
			entity, isGeneratorType := scope.Symbols[key].(*generator.Generator)
			Assert(t, isGeneratorType, "`key` should be defined")

			if key != "person" {
				for _, field := range overridenFields {
					AssertFieldYieldsValue(t, entity, field)
				}
			}
		}
	}
}

func TestDeferredEvaluation(t *testing.T) {
	scope := NewRootScope()
	scope.SetSymbol("foo", int64(10))

	ast, err := dsl.Parse("testScript", []byte(`1 + 2 + 4 * foo * (foo + 18) - foo`))
	AssertNil(t, err, "Should not receive error while parsing")

	i := interp()
	actual, err := i.Visit(ast.(*Node), scope, true)
	AssertNil(t, err, "Should not receive error while interpreting")

	result, ok := actual.(DeferredResolver)
	Assert(t, ok, "Should return a DeferredResolver")

	val, err := result(scope)
	AssertNil(t, err, "Should not receive error while evaluating resolver")

	AssertEqual(t, int64(1113), val)
}

type EvalSpec map[string]interface{}

func TestBinaryExpressionComposition(t *testing.T) {
	i := interp()
	scope := NewRootScope()

	for expr, expected := range (EvalSpec{
		"1 + 2 * 3":                        int64(7),
		"1 * 2 + 3":                        int64(5),
		"(1 + 2) * 3":                      int64(9),
		"5 * 2":                            int64(10),
		"5 / 2":                            float64(2.5),
		"5.0 / 2":                          float64(2.5),
		"5 / 2.0":                          float64(2.5),
		"(\"hi \" + \"thar\" + 5) + false": "hi thar5false",
		"3 * \"hi\"":                       "hihihi",
		"\"hi\" * 3":                       "hihihi",
		"5 * 3.0":                          float64(15),
		"3.0 * 5":                          float64(15),
		"true + \" that\"":                 "true that",
		"1 + 2 + 4 * 10 * (10 + 18) - 10":  int64(1113),
		"(-2 * (6 - 7) / 2) * 88 / 4":      float64(22),
	}) {
		ast, err := dsl.Parse("testScript", []byte(expr))
		AssertNil(t, err, "Should not receive error while parsing %q", expr)

		actual, err := i.Visit(ast.(*Node), scope, false)
		AssertNil(t, err, "Should not receive error while interpreting %q", expr)

		AssertEqual(t, expected, actual, "Incorrect result for %q", expr)
	}
}

func TestBinaryExpressionAsEntityField(t *testing.T) {
	i := interp()
	scope := NewRootScope()
	expr := "entity foo { field: 1 + 1 }"

	ast, err := dsl.Parse("testScript", []byte(expr))
	AssertNil(t, err, "Should not receive error while parsing %q", expr)

	actual, err := i.Visit(ast.(*Node), scope, false)
	AssertNil(t, err, "Should not receive error while interpreting %q", expr)

	entity := actual.(*generator.Generator)
	expectedFieldName := "field"

	Assert(t, entity.HasField(expectedFieldName), "Field %q does not exist", expectedFieldName)
}

func TestLambdaExpression(t *testing.T) {
	scope := NewRootScope()

	script := `
  lambda Square(x) {
    x * x
  }

  let foo = 2, bar = 4

  # demonstrate nested call within inlined call with closure
  (lambda () {
    foo = bar * foo
    Square(foo)
  })()
  `
	ast, err := dsl.Parse("testScript", []byte(script))
	AssertNil(t, err, "Should not receive error while parsing")

	i := interp()
	actual, err := i.Visit(ast.(*Node), scope, false)
	AssertNil(t, err, "Should not receive error while interpreting")

	AssertEqual(t, int64(64), actual, "Unexpected result %T %v", actual, actual)
}

func TestLambdaExpressionNoOp(t *testing.T) {
	scope := NewRootScope()

	script := `
  lambda noop(x) {}
  noop(5)
  `
	ast, err := dsl.Parse("testScript", []byte(script))
	AssertNil(t, err, "Should not receive error while parsing")

	i := interp()
	actual, err := i.Visit(ast.(*Node), scope, false)
	AssertNil(t, err, "Should not receive error while interpreting")
	AssertNil(t, actual, "noop() should do nothing and return nil")
}

func TestLambdaExpressionVariableShadowing(t *testing.T) {
	scope := NewRootScope()

	script := `
  let foo = 1

  lambda BoundParamShadows(foo) {
    "bound lambda arg 'foo' is " + foo
  }

  lambda VarDeclShadows(x) {
    let foo = x
    "declared variable 'foo' within lambda is " + foo
  }

  let shadowed = BoundParamShadows(5) + ", " + VarDeclShadows(10)

  shadowed + ", " + "but outer scoped 'foo' is still " + foo
  `
	ast, err := dsl.Parse("testScript", []byte(script))
	AssertNil(t, err, "Should not receive error while parsing")

	i := interp()

	expected := "bound lambda arg 'foo' is 5, declared variable 'foo' within lambda is 10, but outer scoped 'foo' is still 1"
	actual, err := i.Visit(ast.(*Node), scope, false)

	AssertNil(t, err, "Should not receive error while interpreting")
	AssertEqual(t, expected, actual)
}

func TestLambdaExpressionsAllowComments(t *testing.T) {
	scope := NewRootScope()
	script := `
  let baz = 10

  lambda test1() {
    let foo = 1, bar = 2 # multiple declarations with comment

    baz = 0 # don't break


    # this comment shouldn't break anything
    foo + bar # nor should a terminal comment
  }

  lambda test2() {
    # shouldn't break when first token of lambda body is a comment
    test1()
    # sequential expressions should still work too; this one should return 6
    4, 5, 6
  }

  test1() + test2() # comments outside of lambdas should still be ok
  `

	ast, err := dsl.Parse("testScript", []byte(script))
	AssertNil(t, err, "Should not receive error while parsing; comments may be interfering with parsing.")

	i := interp()

	actual, err := i.Visit(ast.(*Node), scope, false)

	AssertNil(t, err, "Should not receive error while interpreting")
	AssertEqual(t, int64(9), actual, "Comments should not affect interpretation of lambdas and calls")
}

func TestLambdaExpressionUsesStaticScoping(t *testing.T) {
	scope := NewRootScope()

	script := `
  let b = "static"

  lambda fn() {
    lambda foo() {
      b + " scoping"
    }

    lambda bar() { # test that foo has static scope
      let b = "dynamic"
      foo()
    }
  }

  lambda baz() { # test that fn has static scope
    let b = "dynamic"
    (fn())()
  }

  baz() # should invoke bar()
`
	ast, err := dsl.Parse("testScript", []byte(script))
	AssertNil(t, err, "Should not receive error while parsing")

	i := interp()

	actual, err := i.Visit(ast.(*Node), scope, false)

	AssertNil(t, err, "Should not receive error while interpreting")
	AssertEqual(t, "static scoping", actual, "Should be using a lexical/static scope for lambda declarations")
}

func TestLambdaExpressionsWithClosuresContinueToWorkAfterFirstInvocation(t *testing.T) {
	script := `
  let foo

  lambda outer() {
    foo = 1

    lambda inner() {
      foo = foo * 2
    }
  }

  let pow2 = outer()

  pow2() # foo => 2
  pow2() # foo => 4; there was a bug that prevented inner() from executing its body when invoked more than once
  pow2() # foo => 8; there was a bug that prevented inner() from executing its body when invoked more than once

  foo
  `
	scope := NewRootScope()
	ast, err := dsl.Parse("testScript", []byte(script))
	AssertNil(t, err, "Should not receive error while parsing")

	i := interp()

	actual, err := i.Visit(ast.(*Node), scope, false)

	AssertNil(t, err, "Should not receive error while interpreting")
	AssertEqual(t, int64(8), actual, "Lambda closures should continue to work after first invocation when symbols are involved")
}

func TestValidGenerationNodeIdentifierAsCountArg(t *testing.T) {
	i := interp()
	scope := NewRootScope()
	i.EntityFromNode(Entity("person", validFields), scope, false)
	scope.SetSymbol("count", int64(1))
	node := GenNode(nil, NodeSet{Id("count"), Id("person")})
	_, err := i.GenerateFromNode(node, scope, false)
	AssertNil(t, err, "Should be able to use identifiers as count argument")
}

func TestInvalidGenerationNodeBadCountArg(t *testing.T) {
	i := interp()
	scope := NewRootScope()
	i.EntityFromNode(Entity("person", validFields), scope, false)
	node := Generation(0, Id("person"))
	_, err := i.GenerateFromNode(node, scope, false)
	ExpectsError(t, "Must generate at least 1 person{} entity", err)

	scope.SetSymbol("count", "ten")
	node = GenNode(nil, NodeSet{Id("count"), Id("person")})
	_, err = i.GenerateFromNode(node, scope, false)
	ExpectsError(t, "Expected an integer, but got ten", err)
}

func TestEntityWithUndefinedParent(t *testing.T) {
	ent := Entity("person", validFields)
	unresolvable := Id("nope")
	ent.Related = unresolvable
	_, err := interp().EntityFromNode(ent, NewRootScope(), false)
	ExpectsError(t, `Cannot resolve parent entity "nope" for entity "person"`, err)
}

func TestGenerateEntitiesCannotResolveEntityFromNode(t *testing.T) {
	node := Generation(2, Id("tree"))
	_, err := interp().GenerateFromNode(node, NewRootScope(), false)
	ExpectsError(t, `Cannot resolve symbol "tree"`, err)
}

func TestDefaultArguments(t *testing.T) {
	i := interp()
	defaults := map[string]interface{}{
		STRING_TYPE: int64(5),
		INT_TYPE:    [2]int64{1, 10},
		FLOAT_TYPE:  [2]float64{1, 10},
		DATE_TYPE:   []interface{}{UNIX_EPOCH, NOW, ""},
		BOOL_TYPE:   nil,
	}

	for kind, expected := range defaults {
		actual, _ := i.defaultArgumentFor(kind)
		AssertDeepEqual(t, expected, actual)
	}
}

func TestDisallowNondeclaredEntityAsFieldIdentifier(t *testing.T) {
	i := interp()
	_, e := i.EntityFromNode(Entity("hiccup", nestedFields), NewRootScope(), false)
	ExpectsError(t, `Cannot resolve symbol "Goat"`, e)

}

func TestDefaultArgumentsReturnsErrorOnUnsupportedFieldType(t *testing.T) {
	i := interp()
	arg, err := i.defaultArgumentFor("dict")
	if err == nil || err.Error() != "Field of type `dict` requires arguments" {
		t.Errorf("expected an error when getting a default Argument for an unsupported field Type")
	}
	AssertNil(t, arg, "defaultArgumentFor(\"dict\") Should not have returned anything")
}

func TestConfiguringFieldDiesWhenFieldWithoutArgsHasNoDefaults(t *testing.T) {
	i := interp()

	entity := generator.NewGenerator("cat", nil, false)
	ExpectsError(t, "Field of type `dict` requires arguments", i.AddBuiltinField(entity, "name", "dict", []interface{}{}, nil, false))
}

func TestConfiguringFieldWithoutArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	i.AddBuiltinField(testEntity, "last_name", STRING_TYPE, []interface{}{}, nil, false)
	AssertShouldHaveField(t, testEntity, "last_name")
}

func TestConfiguringFieldsForEntityErrors(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	ExpectsError(t, "Field type `$dict` expected 1 args, but 2 found.", i.AddBuiltinField(testEntity, "last_name", DICT_TYPE, []interface{}{1, 10}, nil, false))
}

func TestGeneratedFieldAddedToInterpreterIfPreviousValueExists(t *testing.T) {
	scope := NewRootScope()
	i := interp()
	price := 2.0
	node := Root(Entity("cart", NodeSet{
		Field("price", FloatVal(price)),
		Field("price_clone", Id("price")),
	}))

	entity, _ := i.Visit(node, scope, false)
	resolvedEntity := entity.(*generator.Generator).One(nil, NewTestEmitter(), scope)
	AssertEqual(t, price, resolvedEntity["price_clone"])
}

func TestGeneratedFieldNotAddedToInterpreterIfPreviousValueDoesNotExist(t *testing.T) {
	scope := NewRootScope()
	i := interp()
	node := Root(Entity("cart", NodeSet{
		Field("price_clone", Id("price")),
	}))

	_, err := i.Visit(node, scope, false)

	ExpectsError(t, "Cannot resolve symbol \"price\"", err)
}

func TestConfiguringDistributionWithoutArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	fieldNoArgs := Field("age", Builtin(INT_TYPE))
	field := Field("age", Distribution("uniform"), fieldNoArgs)
	i.withDistributionField(testEntity, field, NewRootScope(), false)
	AssertShouldHaveField(t, testEntity, field.Name)
}

func TestConfiguringDistributionWithArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	fieldArgs := Field("age", Builtin(INT_TYPE), IntArgs(1, 10)...)
	field := Field("age", Distribution("uniform"), fieldArgs)
	i.withDistributionField(testEntity, field, NewRootScope(), false)
	AssertShouldHaveField(t, testEntity, field.Name)
}

func TestConfiguringDistributionWithStaticFields(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	fieldArgs := Field("age", StringVal("blah"))
	field := Field("age", Distribution("percent"), fieldArgs)
	i.withDistributionField(testEntity, field, NewRootScope(), false)
	AssertShouldHaveField(t, testEntity, field.Name)
}

func TestConfiguringDistributionWithMixedFieldTypesShouldBeOkay(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	fieldArgs1 := Field("name", StringVal("disabeled"))
	fieldArgs2 := Field("name", Builtin(ENUM_TYPE), NodeSet{StringCollection("enabled", "pending")}...)
	field := Field("age", Distribution("percent"), fieldArgs1, fieldArgs2)
	i.withDistributionField(testEntity, field, NewRootScope(), false)
	AssertShouldHaveField(t, testEntity, field.Name)
}

func TestConfiguringDistributionWithEntityField(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	scope := NewRootScope()
	i.Visit(Entity("Goat", validFields), scope, false)

	fieldArg1 := Field("friend", Entity("Horse", validFields))
	fieldArg2 := Field("pet", Id("Goat"))
	field := Field("friend", Distribution("percent"), fieldArg1, fieldArg2)
	i.withDistributionField(testEntity, field, scope, false)
	AssertShouldHaveField(t, testEntity, field.Name)
}

func TestConfiguringDistributionWithDeferredFields(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	scope := NewRootScope()

	AssertNil(t, i.AddBuiltinField(testEntity, "age", INT_TYPE, []interface{}{int64(1), int64(10)}, nil, false), "Should not receive error for age field")
	AssertNil(t, i.AddBuiltinField(testEntity, "weight", INT_TYPE, []interface{}{int64(20), int64(30)}, nil, false), "Should not receive error for weight field")
	AssertNil(t, i.withExpressionField(testEntity, "err", "disabled"), "Should not receive error for lit field")

	fieldArg1 := Field("a", Id("age"))
	fieldArg2 := Field("w", Id("weight"))
	fieldArg3 := Field("e", Id("weight"))
	field := Field("ageOrWeightOrErr", Distribution("percent"), fieldArg1, fieldArg2, fieldArg3)
	AssertNil(t, i.withDistributionField(testEntity, field, scope, false), "Should not receive error for ageOrWeightOrErr field")
	AssertShouldHaveField(t, testEntity, field.Name)
}

func TestConfiguringDistributionShouldNotAllowSubDistributions(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	fieldArgs1 := Field("name", StringVal("disabled"))
	fieldArgs2 := Field("age", Distribution("percent"), fieldArgs1)
	field := Field("age", Distribution("percent"), fieldArgs2)
	err := i.withDistributionField(testEntity, field, NewRootScope(), false)
	Assert(t, err != nil, "sub distributions are not allowed!")
}
