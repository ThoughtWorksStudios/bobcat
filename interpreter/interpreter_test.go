package interpreter

import "testing"
import "time"
import "github.com/ThoughtWorksStudios/datagen/dsl"
import "github.com/ThoughtWorksStudios/datagen/generator"

var validFields = []dsl.Node{
	dsl.Node{Kind: "field", Name: "name", Value: dsl.Node{Kind: "builtin", Value: "string"}, Args: stringArgs(10)},
	dsl.Node{Kind: "field", Name: "age", Value: dsl.Node{Kind: "builtin", Value: "integer"}, Args: intArgs(1, 10)},
	dsl.Node{Kind: "field", Name: "weight", Value: dsl.Node{Kind: "builtin", Value: "decimal"}, Args: floatArgs(1, 200)},
	dsl.Node{Kind: "field", Name: "dob", Value: dsl.Node{Kind: "builtin", Value: "date"}, Args: timeArgs("2015-01-01", "2017-01-01")},
	dsl.Node{Kind: "field", Name: "last_name", Value: dsl.Node{Kind: "builtin", Value: "dict"}, Args: dictArgs("last_name")},
}

func TestTranslateEntity(t *testing.T) {
	entity := translateEntity(newEntity("person", validFields))
	for _, field := range validFields {
		assertShouldHaveField(t, entity, field)
	}

}

func TestTranslateEntities(t *testing.T) {
	entity1 := newEntity("cat", validFields)
	entity2 := newEntity("dog", validFields)
	for _, entity := range translateEntities(rootNode(entity1, entity2)) {
		for _, field := range validFields {
			assertShouldHaveField(t, entity, field)
		}
	}
}

func TestValidGenerateEntities(t *testing.T) {
	entities := make(map[string]*generator.Generator)
	entities["person"] = translateEntity(newEntity("person", validFields))
	node := rootNode(generationNode("person", 2))
	err := generateEntities(node, entities)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}
}

func TestGenerateEntisiesRequiresCountTobeGreaterThatZero(t *testing.T) {
	entities := make(map[string]*generator.Generator)
	entities["person"] = translateEntity(newEntity("person", validFields))
	node := rootNode(generationNode("person", 0))
	err := generateEntities(node, entities)
	if err == nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}
}

func TestGenerateEntitiesReturnsErrorIfEntityDoesNotExist(t *testing.T) {
	entities := make(map[string]*generator.Generator)
	node := rootNode(generationNode("person", 0))
	err := generateEntities(node, entities)
	if err == nil {
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

func TestTranslateFieldsForEntity(t *testing.T) {
	testEntity := generator.NewGenerator("person")
	translateFieldsForEntity(testEntity, validFields)
	for _, field := range validFields {
		assertShouldHaveField(t, testEntity, field)
	}
}

func TestConfiguringFieldForEntity(t *testing.T) {
	testEntity := generator.NewGenerator("person")
	for _, field := range validFields {
		configureFieldOn(testEntity, field)
		assertShouldHaveField(t, testEntity, field)
	}

	if testEntity.GetField("wubba lubba dub dub") != nil {
		t.Error("should not get field for non existent field")
	}
}

func TestValInt(t *testing.T) {
	expected := 666
	actual := valInt(stringArg(666))
	assertExpectedEqsActual(t, expected, actual)
}

func TestValStr(t *testing.T) {
	expected := "blah"
	actual := valStr(dictArg("blah"))
	assertExpectedEqsActual(t, expected, actual)
}

func TestValFloat(t *testing.T) {
	expected := 4.2
	actual := valFloat(floatArg(4.2))
	assertExpectedEqsActual(t, expected, actual)
}

func TestValTime(t *testing.T) {
	expected, _ := time.Parse("2006-01-02", "1945-01-01")
	actual := valTime(timeArg(expected))
	assertExpectedEqsActual(t, expected, actual)
}
