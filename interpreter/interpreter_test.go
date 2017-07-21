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
	FieldNode("name", BuiltinNode("string"), IntArgs(10)...),
	FieldNode("age", BuiltinNode("integer"), IntArgs(1, 10)...),
	FieldNode("weight", BuiltinNode("decimal"), FloatArgs(1.0, 200.0)...),
	FieldNode("dob", BuiltinNode("date"), DateArgs("2015-01-01", "2017-01-01")...),
	FieldNode("last_name", BuiltinNode("dict"), StringArgs("last_name")...),
	FieldNode("catch_phrase", StringNode("Grass.... Tastes bad")),
}

var overridenFields = []dsl.Node{
	FieldNode("catch_phrase", StringNode("Grass.... Tastes good")),
}

func interp() *Interpreter {
	return New()
}

func TestValidVisit(t *testing.T) {
	node := RootNode(EntityNode("person", validFields), GenerationNode("person", 2))
	i := interp()
	err := i.Visit(node)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	for _, entity := range i.entities {
		for _, field := range validFields {
			AssertShouldHaveField(t, entity, field)
		}
	}
}

func TestValidVisitWithOverrides(t *testing.T) {
	node := RootNode(EntityNode("person", validFields),
		GenerationNodeWithOverrides("person", overridenFields, 2))
	i := interp()
	err := i.Visit(node)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	for _, entity := range i.entities {
		if entity.GetName() != "person" { // want entity personX where X is random int
			for _, field := range overridenFields {
				AssertFieldShouldBeOverriden(t, entity, field)
			}
		}
	}
}

func TestInvalidGenerationNodeBadArgType(t *testing.T) {
	i := interp()
	i.EntityFromNode(EntityNode("burp", validFields))
	node := dsl.Node{Kind: "generation", Name: "burp", Args: StringArgs("blah")}
	ExpectsError(t, "generate burp takes an integer count", i.GenerateFromNode(node))
}

func TestInvalidGenerationNodeBadCountArg(t *testing.T) {
	i := interp()
	i.EntityFromNode(EntityNode("person", validFields))
	node := GenerationNode("person", 0)
	ExpectsError(t, "Must generate at least 1 `person` entity", i.GenerateFromNode(node))
}

func TestGenerateEntitiesCannotResolveEntity(t *testing.T) {
	node := GenerationNode("tree", 2)
	ExpectsError(t, "Unknown symbol `tree` -- expected an entity. Did you mean to define an entity named `tree`?", interp().GenerateFromNode(node))
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

	badNode := FieldNode("name", BuiltinNode("dict"))
	entity := generator.NewGenerator("cat", GetLogger(t))
	ExpectsError(t, "Field of type `dict` requires arguments", i.withDynamicField(entity, badNode))
}

func TestConfiguringFieldWithoutArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", GetLogger(t))
	fieldNoArgs := FieldNode("last_name", BuiltinNode("string"))
	i.withDynamicField(testEntity, fieldNoArgs)
	AssertShouldHaveField(t, testEntity, fieldNoArgs)
}

func TestConfiguringFieldsForEntityErrors(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", GetLogger(t))
	badNode := FieldNode("last_name", BuiltinNode("dict"), IntArgs(1, 10)...)
	ExpectsError(t, "Field type `dict` expected 1 args, but 2 found.", i.withDynamicField(testEntity, badNode))
}

func TestDynamicFieldRejectsStaticFieldDecl(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", GetLogger(t))
	badField := FieldNode("last_name", IntNode(2), IntArgs(1, 10)...)
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
