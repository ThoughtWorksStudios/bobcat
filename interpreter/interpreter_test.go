package interpreter

import "testing"
import "time"
import "fmt"
import "github.com/ThoughtWorksStudios/datagen/dsl"
import "github.com/ThoughtWorksStudios/datagen/generator"
import . "github.com/ThoughtWorksStudios/datagen/test_helpers"

var validFields = []dsl.Node{
	dsl.Node{Kind: "field", Name: "name", Value: Builtin("string"), Args: IntArgs(10)},
	dsl.Node{Kind: "field", Name: "age", Value: Builtin("integer"), Args: IntArgs(1, 10)},
	dsl.Node{Kind: "field", Name: "weight", Value: Builtin("decimal"), Args: FloatArgs(1.0, 200.0)},
	dsl.Node{Kind: "field", Name: "dob", Value: Builtin("date"), Args: DateArgs("2015-01-01", "2017-01-01")},
	dsl.Node{Kind: "field", Name: "last_name", Value: Builtin("dict"), Args: StringArgs("last_name")},
	dsl.Node{Kind: "field", Name: "catch_phrase", Value: StaticNode("Grass.... Tastes bad")},
}

func TestTranslate(t *testing.T) {
	tree := RootNode(EntityNode("person", validFields), GenerationNode("person", 2))
	err := Translate(tree)
	if err != nil {
		t.Errorf("Failed to translate tree because: %v", err.Error())
	}
}

func TestTranslateEntity(t *testing.T) {
	entity := translateEntity(EntityNode("person", validFields))
	for _, field := range validFields {
		AssertShouldHaveField(t, entity, field)
	}
}

func TestTranslateEntities(t *testing.T) {
	entity1 := EntityNode("cat", validFields)
	entity2 := EntityNode("dog", validFields)
	for _, entity := range translateEntities(RootNode(entity1, entity2)) {
		for _, field := range validFields {
			AssertShouldHaveField(t, entity, field)
		}
	}
}

func TestValidGenerateEntities(t *testing.T) {
	entities := make(map[string]*generator.Generator)
	entities["person"] = translateEntity(EntityNode("person", validFields))
	node := RootNode(GenerationNode("person", 2))
	err := generateEntities(node, entities)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}
}

func TestGenerateEntitiesOnlyAcceptIntCounts(t *testing.T) {
	entities := make(map[string]*generator.Generator)
	entities["burp"] = translateEntity(EntityNode("burp", validFields))
	generationNode := dsl.Node{Kind: "generation", Name: "burp", Args: StringArgs("blah")}
	node := RootNode(generationNode)
	err := generateEntities(node, entities)
	if err == nil || err.Error() != "ERROR: generate burp takes an integer count" {
		t.Errorf("There was a problem generating entities: %v", err)
	}
}

func TestGenerateEntisiesRequiresCountTobeGreaterThatZero(t *testing.T) {
	entities := make(map[string]*generator.Generator)
	entities["person"] = translateEntity(EntityNode("person", validFields))
	node := RootNode(GenerationNode("person", 0))
	err := generateEntities(node, entities)
	if err == nil || err.Error() != "ERROR: Must generate at least 1 `person` entity" {
		t.Errorf("There was a problem generating entities: %v", err)
	}
}

func TestGenerateEntitiesReturnsErrorIfEntityDoesNotExist(t *testing.T) {
	entities := make(map[string]*generator.Generator)
	node := RootNode(GenerationNode("tree", 2))
	err := generateEntities(node, entities)
	if err == nil || err.Error() != "ERROR: tree is undefined; expected entity" {
		t.Errorf("There was a problem generating entities: %v", err)
	}
}

func TestDefaultArgument(t *testing.T) {
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "2017-01-01")
	defaults := map[string]interface{}{
		"string":  5,
		"integer": [2]int{1, 10},
		"decimal": [2]float64{1, 10},
		"date":    [2]time.Time{timeMin, timeMax},
	}

	for kind, expected_value := range defaults {
		actual := defaultArgumentFor(kind)
		if actual != expected_value {
			t.Errorf("default value for argument type '%s' was expected to be %v but was %v", kind, expected_value, actual)
		}
	}
}

func TestDefaultArgumentDiesForTypesWithDefaults(t *testing.T) {
	var died bool = false
	var deathMessage string
	die = func(msg string, args ...interface{}) {
		died = true
		deathMessage = fmt.Sprintf(msg, args)

	}
	defaultArgumentFor("dict")
	if died != true || deathMessage != "Field of type `[dict]` requires arguments" {
		t.Errorf("should throw error for fieldtypes that're not string, integer, decimal, or date")
	}
}

func TestTranslateFieldsForEntity(t *testing.T) {
	testEntity := generator.NewGenerator("person")
	translateFieldsForEntity(testEntity, validFields)
	for _, field := range validFields {
		AssertShouldHaveField(t, testEntity, field)
	}
}

func TestConfiguringFieldForEntity(t *testing.T) {
	testEntity := generator.NewGenerator("person")
	for _, field := range validFields {
		configureFieldOn(testEntity, field)
		AssertShouldHaveField(t, testEntity, field)
	}

	if testEntity.GetField("wubba lubba dub dub") != nil {
		t.Error("should not get field for non existent field")
	}
}

func TestConfiguringFieldWithoutArguments(t *testing.T) {
	testEntity := generator.NewGenerator("person")
	fieldNoArgs := dsl.Node{Kind: "field", Name: "name", Value: Builtin("string")}
	configureFieldOn(testEntity, fieldNoArgs)
	AssertShouldHaveField(t, testEntity, fieldNoArgs)
}

func TestConfiguringFieldsForEntityErrors(t *testing.T) {
	testEntity := generator.NewGenerator("person")
	badField := dsl.Node{Kind: "field", Name: "last_name", Value: Builtin("dict"), Args: IntArgs(1, 10)}

	var died bool = false
	var deathMessage string
	die = func(msg string, args ...interface{}) {
		died = true
		deathMessage = msg

	}
	withDynamicField(testEntity, badField)
	if died != true || deathMessage != "field type `dict` requires exactly 1 argument" {
		t.Errorf("should have died because dict requires exactly 1 argument")
	}
}

func TestDynamicFieldThrowsErrorWhenGivenAStiticField(t *testing.T) {
	testEntity := generator.NewGenerator("person")
	badField := dsl.Node{Kind: "field", Name: "last_name", Value: StaticNode(2), Args: IntArgs(1, 10)}

	var died bool = false
	var deathMessage string
	die = func(msg string, args ...interface{}) {
		died = true
		deathMessage = fmt.Sprintf(msg, args...)

	}
	withDynamicField(testEntity, badField)
	if died != true || deathMessage != "Could not parse field-type for field `last_name`. Expected one of the builtin generator types, but instead got: 2" {
		t.Errorf("Should have died because field-type was not a string")
	}
}

func TestValInt(t *testing.T) {
	expected := 666
	actual := valInt(IntArgs(666)[0])
	AssertExpectedEqsActual(t, expected, actual)
}

func TestValStr(t *testing.T) {
	expected := "blah"
	actual := valStr(StringArgs("blah")[0])
	AssertExpectedEqsActual(t, expected, actual)
}

func TestValFloat(t *testing.T) {
	expected := 4.2
	actual := valFloat(FloatArgs(4.2)[0])
	AssertExpectedEqsActual(t, expected, actual)
}

func TestValTime(t *testing.T) {
	expected, _ := time.Parse("2006-01-02", "1945-01-01")
	actual := valTime(DateArgs("1945-01-01")[0])
	AssertExpectedEqsActual(t, expected, actual)
}
