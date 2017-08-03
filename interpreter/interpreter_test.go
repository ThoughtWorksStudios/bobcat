package interpreter

import (
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/generator"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"testing"
	"time"
)

func AssertShouldHaveField(t *testing.T, entity *generator.Generator, field dsl.Node) {
	AssertNotNil(t, entity.GetField(field.Name), "Expected entity to have field %s, but it did not", field.Name)
}

func AssertFieldShouldBeOverriden(t *testing.T, entity *generator.Generator, field dsl.Node) {
	AssertEqual(t, field.Value.(dsl.Node).Value, entity.GetField(field.Name).GenerateValue())
}

var validFields = []dsl.Node{
	FieldNode("name", BuiltinNode("string"), IntArgs(10), IntArgs(1, 3)),
	FieldNode("age", BuiltinNode("integer"), IntArgs(1, 10), []dsl.Node{}),
	FieldNode("weight", BuiltinNode("decimal"), FloatArgs(1.0, 200.0), []dsl.Node{}),
	FieldNode("dob", BuiltinNode("date"), DateArgs("2015-01-01", "2017-01-01"), []dsl.Node{}),
	FieldNode("last_name", BuiltinNode("dict"), StringArgs("last_name"), []dsl.Node{}),
	FieldNode("catch_phrase", StringNode("Grass.... Tastes bad"), []dsl.Node{}, []dsl.Node{}),
}

var nestedFields = []dsl.Node{
	FieldNode("name", BuiltinNode("string"), IntArgs(10), []dsl.Node{}),
	FieldNode("pet", IdNode("Goat"), IntArgs(2), []dsl.Node{}),
	FieldNode("friend", EntityNode("Horse", validFields), IntArgs(1), []dsl.Node{}),
}

var overridenFields = []dsl.Node{
	FieldNode("catch_phrase", StringNode("Grass.... Tastes good"), []dsl.Node{}, []dsl.Node{}),
}

func interp() *Interpreter {
	return New()
}

func TestScopingResolvesOtherEntities(t *testing.T) {
	scope := NewRootScope()
	i := interp()
	node := RootNode(EntityNode("person", dsl.NodeSet{
		FieldNode("pet", EntityNode("kitteh", overridenFields)),
		FieldNode("pets_can_have_pets_too", EntityNode("lolcat", dsl.NodeSet{
			FieldNode("cheezburgrz", StringNode("can has")),
			FieldNode("protoype", IdNode("kitteh")),
		})),
	}))
	err := i.Visit(node, scope)
	AssertNil(t, err, "`lolcat` should be able to resolve `kitteh` because it lives within the scope hierarchy. error was %v", err)

	// using same root scope to simulate previously defined symbols
	err = i.Visit(RootNode(GenerationNode(IdNode("person"), 2)), scope)
	AssertNil(t, err, "Should be able to resolve `person` because it is defined in the root scope. error was %v", err)

	// using same root scope to simulate previously defined symbols; here, `kitteh` was defined in a child scope of `person`,
	// but not at the root scope, so we should not be able to resolve it.
	ExpectsError(t, "Cannot resolve symbol \"kitteh\"", i.Visit(RootNode(GenerationNode(IdNode("kitteh"), 1)), scope))
}

func TestValidVisit(t *testing.T) {
	node := RootNode(EntityNode("person", validFields), GenerationNode(IdNode("person"), 2))
	i := interp()
	scope := NewRootScope()
	err := i.Visit(node, scope)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	for _, entry := range scope.symbols {
		entity := entry.Value.(*generator.Generator)
		for _, field := range validFields {
			AssertShouldHaveField(t, entity, field)
		}
	}
}

func TestValidVisitWithNesting(t *testing.T) {
	node := RootNode(EntityNode("Goat", validFields), EntityNode("person", nestedFields),
		GenerationNode(IdNode("person"), 2))
	i := interp()

	scope := NewRootScope()
	err := i.Visit(node, scope)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	person, _ := i.ResolveEntity(IdNode("person"), scope)
	for _, field := range nestedFields {
		AssertShouldHaveField(t, person, field)
	}
}

func TestValidVisitWithOverrides(t *testing.T) {
	node := RootNode(
		EntityNode("person", validFields),
		GenerationNode(
			EntityExtensionNode("lazyPerson", "person", overridenFields),
			2,
		),
	)
	i := interp()
	err := i.Visit(node, NewRootScope())
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	for _, entity := range i.entities {
		if entity.Name != "person" { // want entity personX where X is random int
			for _, field := range overridenFields {
				AssertFieldShouldBeOverriden(t, entity, field)
			}
		}
	}
}

func TestInvalidGenerationNodeBadArgType(t *testing.T) {
	i := interp()
	scope := NewRootScope()
	i.EntityFromNode(EntityNode("burp", validFields), scope)
	node := dsl.Node{Kind: "generation", Value: IdNode("burp"), Args: StringArgs("blah")}
	ExpectsError(t, `generate "burp" takes an integer count`, i.GenerateFromNode(node, scope))
}

func TestInvalidGenerationNodeBadCountArg(t *testing.T) {
	i := interp()
	scope := NewRootScope()
	i.EntityFromNode(EntityNode("person", validFields), scope)
	node := GenerationNode(IdNode("person"), 0)
	ExpectsError(t, "Must generate at least 1 `person` entity", i.GenerateFromNode(node, scope))
}

func TestEntityWithUndefinedParent(t *testing.T) {
	ent := EntityNode("person", validFields)
	unresolvable := IdNode("nope")
	ent.Related = &unresolvable
	_, err := interp().EntityFromNode(ent, NewRootScope())
	ExpectsError(t, `Cannot resolve parent entity "nope" for entity "person"`, err)
}

func TestGenerateEntitiesCannotResolveEntity(t *testing.T) {
	node := GenerationNode(IdNode("tree"), 2)
	ExpectsError(t, `Cannot resolve symbol "tree"`, interp().GenerateFromNode(node, NewRootScope()))
}

func TestDefaultArguments(t *testing.T) {
	i := interp()
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "2017-01-01")
	defaults := map[string]interface{}{
		"string":  5,
		"integer": [2]int{1, 10},
		"decimal": [2]float64{1, 10},
		"date":    [2]time.Time{timeMin, timeMax},
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
	_, e := i.EntityFromNode(EntityNode("hiccup", nestedFields), NewRootScope())
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

	badNode := FieldNode("name", BuiltinNode("dict"), []dsl.Node{}, []dsl.Node{})
	entity := generator.NewGenerator("cat", GetLogger(t))
	ExpectsError(t, "Field of type `dict` requires arguments", i.withDynamicField(entity, badNode, NewRootScope()))
}

func TestConfiguringFieldWithoutArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", GetLogger(t))
	fieldNoArgs := FieldNode("last_name", BuiltinNode("string"), []dsl.Node{}, []dsl.Node{})
	i.withDynamicField(testEntity, fieldNoArgs, NewRootScope())
	AssertShouldHaveField(t, testEntity, fieldNoArgs)
}

func TestConfiguringFieldsForEntityErrors(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", GetLogger(t))
	badNode := FieldNode("last_name", BuiltinNode("dict"), IntArgs(1, 10), []dsl.Node{})
	ExpectsError(t, "Field type `dict` expected 1 args, but 2 found.", i.withDynamicField(testEntity, badNode), NewRootScope())
}

func TestDynamicFieldRejectsStaticFieldDecl(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", GetLogger(t))
	badField := FieldNode("last_name", IntNode(2), IntArgs(1, 10), []dsl.Node{})
	ExpectsError(t, "Could not parse field-type for field `last_name`. Expected one of the builtin generator types, but instead got: 2", i.withDynamicField(testEntity, badField))
}

func TestValInt(t *testing.T) {
	expected := 666
	actual := valInt(IntArgs(666)[0])
	AssertEqual(t, expected, actual)
}

func TestValStr(t *testing.T) {
	expected := "blah"
	actual := valStr(StringArgs("blah")[0])
	AssertEqual(t, expected, actual)
}

func TestValFloat(t *testing.T) {
	expected := 4.2
	actual := valFloat(FloatArgs(4.2)[0])
	AssertEqual(t, expected, actual)
}

func TestValTime(t *testing.T) {
	expected, _ := time.Parse("2006-01-02", "1945-01-01")
	actual := valTime(DateArgs("1945-01-01")[0])
	AssertEqual(t, expected, actual)
}
