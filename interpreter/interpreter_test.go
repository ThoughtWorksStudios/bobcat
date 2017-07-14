package interpreter

import (
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/generator"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"testing"
	"time"
)

var validFields = []dsl.Node{
	FieldNode("name", BuiltinNode("string"), IntArgs(10)...),
	FieldNode("age", BuiltinNode("integer"), IntArgs(1, 10)...),
	FieldNode("weight", BuiltinNode("decimal"), FloatArgs(1.0, 200.0)...),
	FieldNode("dob", BuiltinNode("date"), DateArgs("2015-01-01", "2017-01-01")...),
	FieldNode("last_name", BuiltinNode("dict"), StringArgs("last_name")...),
	FieldNode("catch_phrase", StringNode("Grass.... Tastes bad")),
}

func interp() *Interpreter {
	return New(GetLogger())
}

func TestTranslateEntities(t *testing.T) {
	entity1 := EntityNode("cat", validFields)
	entity2 := EntityNode("dog", validFields)
	for _, entity := range interp().translateEntities(RootNode(entity1, entity2)) {
		for _, field := range validFields {
			AssertShouldHaveField(t, entity, field)
		}
	}
}

func TestValidConsume(t *testing.T) {
	node := RootNode(EntityNode("person", validFields), GenerationNode("person", 2))
	i := interp()
	err := i.Consume(node)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}
}

func TestInvalidGenerationNodeBadArgType(t *testing.T) {
	i := interp()
	i.EntityFromNode(EntityNode("burp", validFields))
	node := RootNode(dsl.Node{Kind: "generation", Name: "burp", Args: StringArgs("blah")})
	ExpectsError(t, "ERROR: generate burp takes an integer count", i.generateEntities(node))
}

func TestInvalidGenerationNodeBadCountArg(t *testing.T) {
	i := interp()
	i.EntityFromNode(EntityNode("person", validFields))
	node := RootNode(GenerationNode("person", 0))
	ExpectsError(t, "ERROR: Must generate at least 1 `person` entity", i.generateEntities(node))
}

func TestGenerateEntitiesCannotResolveEntity(t *testing.T) {
	node := RootNode(GenerationNode("tree", 2))
	ExpectsError(t, "ERROR: Unknown symbol `tree` -- expected an entity. Did you mean to define an entity named `tree`?", interp().generateEntities(node))
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
		actual := i.defaultArgumentFor(kind)
		if actual != expected_value {
			t.Errorf("default value for argument type '%s' was expected to be %v but was %v", kind, expected_value, actual)
		}
	}
}

func TestDefaultArgumentsDiesOnUnsupportedFieldType(t *testing.T) {
	i := interp()
	AssertNil(t, i.defaultArgumentFor("dict"), "defaultArgumentFor(\"dict\") Should not have returned anything")
	i.l.(*TestLogger).AssertMessage(t, "Field of type `dict` requires arguments")
}

func TestConfiguringFieldWithoutArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person")
	fieldNoArgs := FieldNode("last_name", BuiltinNode("string"))
	i.withDynamicField(testEntity, fieldNoArgs)
	AssertShouldHaveField(t, testEntity, fieldNoArgs)
}

func TestConfiguringFieldsForEntityErrors(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person")
	badNode := FieldNode("last_name", BuiltinNode("dict"), IntArgs(1, 10)...)
	i.withDynamicField(testEntity, badNode)
	i.l.(*TestLogger).AssertMessage(t, "Field type `dict` requires exactly 1 argument")
}

func TestDynamicFieldRejectsStaticFieldDecl(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person")
	badField := FieldNode("last_name", IntNode(2), IntArgs(1, 10)...)
	i.withDynamicField(testEntity, badField)
	i.l.(*TestLogger).AssertMessage(t, "Could not parse field-type for field `last_name`. Expected one of the builtin generator types, but instead got: 2")
}

func TestValInt(t *testing.T) {
	expected := 666
	actual := valInt(IntArgs(666)[0])
	AsserEqual(t, expected, actual)
}

func TestValStr(t *testing.T) {
	expected := "blah"
	actual := valStr(StringArgs("blah")[0])
	AsserEqual(t, expected, actual)
}

func TestValFloat(t *testing.T) {
	expected := 4.2
	actual := valFloat(FloatArgs(4.2)[0])
	AsserEqual(t, expected, actual)
}

func TestValTime(t *testing.T) {
	expected, _ := time.Parse("2006-01-02", "1945-01-01")
	actual := valTime(DateArgs("1945-01-01")[0])
	AsserEqual(t, expected, actual)
}
