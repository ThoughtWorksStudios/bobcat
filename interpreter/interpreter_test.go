package interpreter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/dsl"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/ThoughtWorksStudios/bobcat/generator"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func AssertShouldHaveField(t *testing.T, entity *generator.Generator, field *Node) {
	emitter := NewDummyEmitter()
	result := entity.One(nil, emitter)
	AssertNotNil(t, result[field.Name], "Expected entity to have field %s, but it did not", field.Name)
}

func AssertFieldYieldsValue(t *testing.T, entity *generator.Generator, field *Node) {
	emitter := NewDummyEmitter()
	result := entity.One(nil, emitter)
	AssertEqual(t, field.ValNode().Value, result[field.Name])
}

var validFields = NodeSet{
	Field("name", Builtin("string"), IntArgs(10)...),
	Field("age", Builtin("integer"), IntArgs(1, 10)...),
	Field("weight", Builtin("decimal"), FloatArgs(1.0, 200.0)...),
	Field("dob", Builtin("date"), DateArgs("2015-01-01", "2017-01-01")...),
	Field("last_name", Builtin("dict"), StringArgs("last_name")...),
	Field("status", Builtin("enum"), NodeSet{StringCollection("enabled", "disabled")}...),
	Field("status", Builtin("serial")),
	Field("catch_phrase", StringVal("Grass.... Tastes bad")),
}

var nestedFields = NodeSet{
	Field("name", Builtin("string"), IntArgs(10)...),
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

	for _, entry := range scope.symbols {
		entity := entry.(*generator.Generator)
		for _, field := range validFields {
			AssertShouldHaveField(t, entity, field)
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

	person, _ := i.ResolveEntity(Id("person"), scope)
	for _, field := range nestedFields {
		AssertShouldHaveField(t, person, field)
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

	AssertEqual(t, 2, len(scope.symbols), "Should have 2 entities defined")

	for _, key := range []string{"person", "lazyPerson"} {
		_, isPresent := scope.symbols[key]
		// don't try to use AssertNotNil here; it won't work because it is unable to detect
		// whether a nil pointer passed as an interface{} param to AssertNotEqual is nil.
		// see this crazy shit: https://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
		Assert(t, isPresent, "`%v` should be defined in scope", key)

		if isPresent {
			entity, isGeneratorType := scope.symbols[key].(*generator.Generator)
			Assert(t, isGeneratorType, "`key` should be defined")

			if key != "person" {
				for _, field := range overridenFields {
					AssertFieldYieldsValue(t, entity, field)
				}
			}
		}
	}
}

type EvalSpec map[string]interface{}

func TestBinaryExpressionComposition(t *testing.T) {
	i := interp()
	scope := NewRootScope()

	for expr, expected := range (EvalSpec{
		"1 + 2 * 3":                      int64(7),
		"1 * 2 + 3":                      int64(5),
		"(1 + 2) * 3":                    int64(9),
		"5 * 2":                          int64(10),
		"5 / 2":                          float64(2.5),
		"5.0 / 2":                        float64(2.5),
		"5 / 2.0":                        float64(2.5),
		"(\"hi \" + \"thar\" + 5) + false": "hi thar5false",
		"3 * \"hi\"":                     "hihihi",
		"\"hi\" * 3":                     "hihihi",
		"5 * 3.0":                        float64(15),
		"3.0 * 5":                        float64(15),
		"true + \" that\"":               "true that",
		"1 + 2 + 4 * 10 * (10 + 18) - 10": int64(1113),
		"(-2 * (6 - 7) / 2) * 88 / 4":     float64(22),
	}) {
		ast, err := dsl.Parse("testScript", []byte(expr))
		AssertNil(t, err, "Should not receive error while parsing %q", expr)

		actual, err := i.Visit(ast.(*Node), scope, false)
		AssertNil(t, err, "Should not receive error while interpreting %q", expr)

		AssertEqual(t, expected, actual, "Incorrect result for %q", expr)
	}
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

func TestGenerateEntitiesCannotResolveEntity(t *testing.T) {
	node := Generation(2, Id("tree"))
	_, err := interp().GenerateFromNode(node, NewRootScope(), false)
	ExpectsError(t, `Cannot resolve symbol "tree"`, err)
}

func TestDefaultArguments(t *testing.T) {
	i := interp()
	defaults := map[string]interface{}{
		"string":  int64(5),
		"integer": [2]int64{1, 10},
		"decimal": [2]float64{1, 10},
		"date":    []interface{}{UNIX_EPOCH, NOW, ""},
		"bool":    nil,
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

	badNode := Field("name", Builtin("dict"))
	entity := generator.NewGenerator("cat", nil, false)
	ExpectsError(t, "Field of type `dict` requires arguments", i.withDynamicField(entity, badNode, NewRootScope(), false))
}

func TestConfiguringFieldWithoutArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	fieldNoArgs := Field("last_name", Builtin("string"))
	i.withDynamicField(testEntity, fieldNoArgs, NewRootScope(), false)
	AssertShouldHaveField(t, testEntity, fieldNoArgs)
}

func TestConfiguringFieldsForEntityErrors(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	badNode := Field("last_name", Builtin("dict"), IntArgs(1, 10)...)
	ExpectsError(t, "Field type `dict` expected 1 args, but 2 found.", i.withDynamicField(testEntity, badNode, NewRootScope(), false))
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
	resolvedEntity := entity.(*generator.Generator).One(nil, NewTestEmitter())
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
	fieldNoArgs := Field("age", Builtin("integer"))
	field := Field("age", Distribution("uniform"), fieldNoArgs)
	i.withDistributionField(testEntity, field, NewRootScope(), false)
	AssertShouldHaveField(t, testEntity, field)
}

func TestConfiguringDistributionWithArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", nil, false)
	fieldArgs := Field("age", Builtin("integer"), IntArgs(1, 10)...)
	field := Field("age", Distribution("uniform"), fieldArgs)
	i.withDistributionField(testEntity, field, NewRootScope(), false)
	AssertShouldHaveField(t, testEntity, field)
}
