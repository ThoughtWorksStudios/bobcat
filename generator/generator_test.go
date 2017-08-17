package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"github.com/rs/xid"
	"reflect"
	"testing"
	"time"
)

func isBetween(actual, lower, upper float64) bool {
	return actual >= lower && actual <= upper
}

/*
 * is this a cheap hack? you bet it is.
 */
func equiv(expected, actual *Field) bool {
	return fmt.Sprintf("%v", expected.fieldType) == fmt.Sprintf("%v", actual.fieldType)
}

func AssertEquiv(t *testing.T, expected, actual *Field) {
	if !equiv(expected, actual) {
		t.Errorf("Expected: \n [%v] \n\n but got: \n [%v]", expected.fieldType, actual.fieldType)
	}
}

func TestExtendGenerator(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)

	g.WithField("name", "string", int64(10), nil)
	g.WithField("age", "decimal", [2]float64{2, 4}, nil)
	g.WithStaticField("species", "human")

	m := ExtendGenerator("thang", g)
	m.WithStaticField("species", "h00man")
	m.WithStaticField("name", "kyle")

	data := GeneratedEntities{}

	data = g.Generate(1)

	base := data[0]

	AssertEqual(t, "human", base["species"])
	AssertEqual(t, 10, len(base["name"].(string)))
	Assert(t, isBetween(base["age"].(float64), 2, 4), "base entity failed to generate the correct age")

	data = m.Generate(1)

	extended := data[0]
	AssertEqual(t, "h00man", extended["species"])
	AssertEqual(t, "kyle", extended["name"].(string))
	Assert(t, isBetween(extended["age"].(float64), 2, 4), "extended entity failed to generate the correct age")
}

func TestSubentityHasParentReference(t *testing.T) {
	logger := GetLogger(t)

	subentityGenerator := NewGenerator("Cat", logger)
	subentityGenerator.WithField("name", "string", 5, nil)

	g := NewGenerator("Person", logger)
	g.WithField("name", "string", int64(10), nil)
	g.WithEntityField("pet", subentityGenerator, 1, nil)

	entities := g.Generate(3)
	person := entities[0]
	cat := entities[0]["pet"].(EntityResult)

	if person["$id"] != cat["$parent"] {
		t.Errorf("Parent id (%v) on subentity does not match the parent entity's id (%v)", cat["$parent"], person["$id"])
	}

	nextCat := subentityGenerator.Generate(1)[0]

	if val, ok := nextCat["$parent"]; ok {
		t.Errorf("Cat should not have a parent (%v) when generated on it's own", val)
	}
}

func TestWithFieldCreatesCorrectFields(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("login", "string", int64(2), nil)
	g.WithField("age", "integer", [2]int64{2, 4}, nil)
	g.WithField("stars", "decimal", [2]float64{2.85, 4.50}, nil)
	g.WithField("dob", "date", [2]time.Time{timeMin, timeMax}, nil)

	expectedFields := []struct {
		fieldName string
		field     *Field
	}{
		{"login", NewField(&StringType{2}, nil)},
		{"age", NewField(&IntegerType{2, 4}, nil)},
		{"stars", NewField(&FloatType{2.85, 4.50}, nil)},
		{"dob", NewField(&DateType{timeMin, timeMax}, nil)},
		{"$id", NewField(&MongoIDType{}, nil)},
	}

	for _, expectedField := range expectedFields {
		AssertEquiv(t, expectedField.field, g.fields[expectedField.fieldName])
	}
}

func TestIntegerRangeIsCorrect(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	ExpectsError(t, fmt.Sprintf("max %d cannot be less than min %d", 2, 4), g.WithField("age", "integer", [2]int64{4, 2}, nil))
}

func TestDateRangeIsCorrect(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	err := g.WithField("dob", "date", [2]time.Time{timeMax, timeMin}, nil)
	expected := fmt.Sprintf("max %s cannot be before min %s", timeMin, timeMax)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDecimalRangeIsCorrect(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	err := g.WithField("stars", "decimal", [2]float64{4.4, 2.0}, nil)
	expected := fmt.Sprintf("max %v cannot be less than min %v", 2.0, 4.4)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestWithStaticFieldCreatesCorrectField(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	g.WithStaticField("login", "something")
	expectedField := NewField(&LiteralType{"something"}, nil)
	AssertEquiv(t, expectedField, g.fields["login"])
}

func TestWithEntityFieldCreatesCorrectField(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	countRange := &CountRange{3, 3}
	g.WithEntityField("food", g, 3, countRange)
	expectedField := NewField(&EntityType{g}, countRange)
	AssertEquiv(t, expectedField, g.fields["food"])
}

func TestInvalidFieldType(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	ExpectsError(t, fmt.Sprintf("Invalid field type '%s'", "foo"),
		g.WithField("login", "foo", 2, nil))
}

func TestFieldArgsMatchesFieldType(t *testing.T) {
	var testFields = []struct {
		fieldType string
		fieldArgs interface{}
	}{
		{"string", "string"},
		{"integer", "string"},
		{"decimal", "string"},
		{"date", "string"},
		{"enum", "string"},
		{"dict", 0},
	}

	logger := GetLogger(t)
	g := NewGenerator("thing", logger)

	for _, field := range testFields {
		AssertNotNil(t, g.WithField("fieldName", field.fieldType, field.fieldArgs, nil),
			"Mismatched field args type for field type '%s' should be logged", field.fieldType)
	}
}

func TestGenerateProducesGeneratedContent(t *testing.T) {
	data := GeneratedEntities{}
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("a", "string", int64(2), nil)
	g.WithField("b", "integer", [2]int64{2, 4}, nil)
	g.WithField("c", "decimal", [2]float64{2.85, 4.50}, nil)
	g.WithField("d", "date", [2]time.Time{timeMin, timeMax}, nil)
	g.WithField("e", "dict", "last_name", nil)
	g.WithField("f", "mongoid", "", nil)
	g.WithField("g", "enum", []interface{}{"eek", "two"}, nil)

	data = g.Generate(3)

	AssertEqual(t, 3, len(data))

	var testFields = []struct {
		fieldName string
		fieldType interface{}
	}{
		{"a", "string"},
		{"b", int64(1)},
		{"c", 2.1},
		{"d", time.Time{}},
		{"e", "string"},
		{"f", xid.New().String()},
		{"g", "string"},
	}

	entity := data[0]
	for _, field := range testFields {
		actual := reflect.TypeOf(entity[field.fieldName])
		expected := reflect.TypeOf(field.fieldType)
		AssertEqual(t, expected, actual)
	}
}

func TestGenerateWithBoundsArgumentProducesCorrectCountOfValues(t *testing.T) {
	data := GeneratedEntities{}
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("a", "string", int64(2), &CountRange{2, 2})
	g.WithField("b", "integer", [2]int64{2, 4}, &CountRange{3, 3})
	g.WithField("c", "decimal", [2]float64{2.85, 4.50}, &CountRange{4, 4})
	g.WithField("d", "date", [2]time.Time{timeMin, timeMax}, &CountRange{5, 5})
	g.WithField("e", "dict", "last_name", &CountRange{6, 6})
	g.WithEntityField("f", NewGenerator("subthing", logger), 1, &CountRange{7, 7})
	g.WithField("g", "enum", []interface{}{"1"}, &CountRange{8, 8})

	data = g.Generate(1)

	var testFields = []struct {
		fieldName string
		count     int
	}{
		{"a", 2},
		{"b", 3},
		{"c", 4},
		{"d", 5},
		{"e", 6},
		{"f", 7},
		{"g", 8},
	}

	entity := data[0]
	for _, field := range testFields {
		actual := len(entity[field.fieldName].([]interface{}))
		AssertEqual(t, field.count, actual)
	}
}
