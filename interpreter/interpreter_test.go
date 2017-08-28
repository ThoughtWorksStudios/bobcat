package interpreter

import (
	ast "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/generator"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
	"time"
)

func AssertShouldHaveField(t *testing.T, entity *generator.Generator, field *ast.Node) {
	result := entity.One("")
	AssertNotNil(t, result[field.Name], "Expected entity to have field %s, but it did not", field.Name)
}

func AssertFieldYieldsValue(t *testing.T, entity *generator.Generator, field *ast.Node) {
	result := entity.One("")
	AssertEqual(t, field.ValNode().Value, result[field.Name])
}

var validFields = ast.NodeSet{
	Field("name", Builtin("string"), IntArgs(10)...),
	Field("age", Builtin("integer"), IntArgs(1, 10)...),
	Field("weight", Builtin("decimal"), FloatArgs(1.0, 200.0)...),
	Field("dob", Builtin("date"), DateArgs("2015-01-01", "2017-01-01")...),
	Field("last_name", Builtin("dict"), StringArgs("last_name")...),
	Field("status", Builtin("enum"), ast.NodeSet{StringCollection("enabled", "disabled")}...),
	Field("catch_phrase", StringVal("Grass.... Tastes bad")),
}

var nestedFields = ast.NodeSet{
	Field("name", Builtin("string"), IntArgs(10)...),
	Field("pet", Id("Goat")),
	Field("friend", Entity("Horse", validFields)),
}

var overridenFields = ast.NodeSet{
	Field("catch_phrase", StringVal("Grass.... Tastes good")),
}

func interp() *Interpreter {
	return New(false, false)
}

func TestScopingResolvesOtherEntities(t *testing.T) {
	scope := NewRootScope()
	i := interp()
	node := Root(Entity("person", ast.NodeSet{
		Field("pet", Entity("kitteh", overridenFields)),
		Field("pets_can_have_pets_too", Entity("lolcat", ast.NodeSet{
			Field("cheezburgrz", StringVal("can has")),
			Field("protoype", Id("kitteh")),
		})),
	}))
	_, err := i.Visit(node, scope)
	AssertNil(t, err, "`lolcat` should be able to resolve `kitteh` because it lives within the scope hierarchy. error was %v", err)

	// using same root scope to simulate previously defined symbols
	_, err = i.Visit(Root(Generation(2, Id("person"))), scope)
	AssertNil(t, err, "Should be able to resolve `person` because it is defined in the root scope. error was %v", err)

	// using same root scope to simulate previously defined symbols; here, `kitteh` was defined in a child scope of `person`,
	// but not at the root scope, so we should not be able to resolve it.
	_, err = i.Visit(Root(Generation(1, Id("kitteh"))), scope)
	ExpectsError(t, "Cannot resolve symbol \"kitteh\"", err)
}

func TestValidVisit(t *testing.T) {
	node := Root(Entity("person", validFields), Generation(2, Id("person")))
	i := interp()
	scope := NewRootScope()
	_, err := i.Visit(node, scope)
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
	_, err := i.Visit(node, scope)
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

	if _, err := i.Visit(node, scope); err != nil {
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

func TestInvalidGenerationNodeBadCountArg(t *testing.T) {
	i := interp()
	scope := NewRootScope()
	i.EntityFromNode(Entity("person", validFields), scope)
	node := Generation(0, Id("person"))
	_, err := i.GenerateFromNode(node, scope)
	ExpectsError(t, "Must generate at least 1 person{} entity", err)
}

func TestEntityWithUndefinedParent(t *testing.T) {
	ent := Entity("person", validFields)
	unresolvable := Id("nope")
	ent.Related = unresolvable
	_, err := interp().EntityFromNode(ent, NewRootScope())
	ExpectsError(t, `Cannot resolve parent entity "nope" for entity "person"`, err)
}

func TestGenerateEntitiesCannotResolveEntity(t *testing.T) {
	node := Generation(2, Id("tree"))
	_, err := interp().GenerateFromNode(node, NewRootScope())
	ExpectsError(t, `Cannot resolve symbol "tree"`, err)
}

func TestDefaultArguments(t *testing.T) {
	i := interp()
	defaults := map[string]interface{}{
		"string":  int64(5),
		"integer": [2]int64{1, 10},
		"decimal": [2]float64{1, 10},
		"date":    [2]time.Time{UNIX_EPOCH, NOW},
		"bool":    nil,
	}

	for kind, expected_value := range defaults {
		actual, _ := i.defaultArgumentFor(kind)
		if actual != expected_value {
			t.Errorf("default value for argument type '%s' was expected to be %v but was %v", kind, expected_value, actual)
		}
	}
}

func TestDisallowNondeclaredEntityAsFieldIdentifier(t *testing.T) {
	i := interp()
	_, e := i.EntityFromNode(Entity("hiccup", nestedFields), NewRootScope())
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
	entity := generator.NewGenerator("cat", false)
	ExpectsError(t, "Field of type `dict` requires arguments", i.withDynamicField(entity, badNode, NewRootScope()))
}

func TestConfiguringFieldWithoutArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", false)
	fieldNoArgs := Field("last_name", Builtin("string"))
	i.withDynamicField(testEntity, fieldNoArgs, NewRootScope())
	AssertShouldHaveField(t, testEntity, fieldNoArgs)
}

func TestConfiguringFieldsForEntityErrors(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", false)
	badNode := Field("last_name", Builtin("dict"), IntArgs(1, 10)...)
	ExpectsError(t, "Field type `dict` expected 1 args, but 2 found.", i.withDynamicField(testEntity, badNode, NewRootScope()))
}
