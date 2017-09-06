package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"strings"
	"testing"
	"time"
)

func collection(vals ...interface{}) []interface{} {
	return vals
}

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
	g := NewGenerator("thing", false)

	g.WithField("name", "string", int64(10), nil, false)
	g.WithField("age", "decimal", [2]float64{2, 4}, nil, false)
	g.WithStaticField("species", "human")

	m := ExtendGenerator("thang", false, g)
	m.WithStaticField("species", "h00man")
	m.WithStaticField("name", "kyle")

	emitter := NewTestEmitter()

	g.Generate(1, emitter)

	base := emitter.Shift()
	AssertNotNil(t, base, "Should have generated an entity result")

	AssertEqual(t, "human", base["species"])
	AssertEqual(t, 10, len(base["name"].(string)))
	Assert(t, isBetween(base["age"].(float64), 2, 4), "base entity failed to generate the correct age")

	m.Generate(1, emitter)

	extended := emitter.Shift()
	AssertNotNil(t, extended, "Should have generated an entity result")
	AssertEqual(t, "h00man", extended["species"])
	AssertEqual(t, "kyle", extended["name"].(string))
	Assert(t, isBetween(extended["age"].(float64), 2, 4), "extended entity failed to generate the correct age")
}

func TestNoMetadataGeneratedWhenDisabled(t *testing.T) {
	g := NewGenerator("Cat", true)
	g.WithField("name", "string", 5, nil, false)

	emitter := NewTestEmitter()
	g.One("foo", emitter)
	entity := emitter.Shift()

	for name, _ := range entity {
		if strings.HasPrefix(name, "_") && name != "_id" && name != "_parent" {
			t.Errorf("Found metadata in entity when there should be none, '%v'", name)
		}
	}
}

func TestSubentityHasParentReference(t *testing.T) {
	subentityGenerator := NewGenerator("Cat", false)
	subentityGenerator.WithField("name", "string", 5, nil, false)

	g := NewGenerator("Person", false)
	g.WithField("name", "string", int64(10), nil, false)
	g.WithEntityField("pet", subentityGenerator, 1, nil)

	emitter := NewTestEmitter()

	g.Generate(1, emitter)
	cat := emitter.Shift()
	person := emitter.Shift()

	if person["_id"] != cat["_parent"] {
		t.Errorf("Parent id (%v) on subentity does not match the parent entity's id (%v)", cat["_parent"], person["_id"])
	}

	subentityGenerator.Generate(1, emitter)
	nextCat := emitter.Shift()

	if val, ok := nextCat["_parent"]; ok {
		t.Errorf("Cat should not have a parent (%v) when generated on it's own", val)
	}
}

func TestWithFieldCreatesCorrectFields(t *testing.T) {
	g := NewGenerator("thing", false)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("login", "string", int64(2), nil, false)
	g.WithField("age", "integer", [2]int64{2, 4}, nil, false)
	g.WithField("stars", "decimal", [2]float64{2.85, 4.50}, nil, false)
	g.WithField("dob", "date", [2]time.Time{timeMin, timeMax}, nil, false)
	g.WithField("counter", "serial", nil, nil, false)

	expectedFields := []struct {
		fieldName string
		field     *Field
	}{
		{"login", NewField(&StringType{2}, nil, false)},
		{"age", NewField(&IntegerType{2, 4}, nil, false)},
		{"stars", NewField(&FloatType{2.85, 4.50}, nil, false)},
		{"dob", NewField(&DateType{timeMin, timeMax}, nil, false)},
		{"_id", NewField(&MongoIDType{}, nil, false)},
		{"counter", NewField(&SerialType{}, nil, false)},
	}

	for _, expectedField := range expectedFields {
		AssertEquiv(t, expectedField.field, g.fields[expectedField.fieldName])
	}
}

func TestIntegerRangeIsCorrect(t *testing.T) {
	g := NewGenerator("thing", false)
	ExpectsError(t, fmt.Sprintf("max %d cannot be less than min %d", 2, 4), g.WithField("age", "integer", [2]int64{4, 2}, nil, false))
}

func TestDateRangeIsCorrect(t *testing.T) {
	g := NewGenerator("thing", false)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	err := g.WithField("dob", "date", [2]time.Time{timeMax, timeMin}, nil, false)
	expected := fmt.Sprintf("max %s cannot be before min %s", timeMin, timeMax)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDecimalRangeIsCorrect(t *testing.T) {
	g := NewGenerator("thing", false)
	err := g.WithField("stars", "decimal", [2]float64{4.4, 2.0}, nil, false)
	expected := fmt.Sprintf("max %v cannot be less than min %v", 2.0, 4.4)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestWithStaticFieldCreatesCorrectField(t *testing.T) {
	g := NewGenerator("thing", false)
	g.WithStaticField("login", "something")
	expectedField := NewField(&LiteralType{"something"}, nil, false)
	AssertEquiv(t, expectedField, g.fields["login"])
}

func TestWithEntityFieldCreatesCorrectField(t *testing.T) {
	g := NewGenerator("thing", false)
	countRange := &CountRange{3, 3}
	g.WithEntityField("food", g, 3, countRange)
	expectedField := NewField(&EntityType{g}, countRange, false)
	AssertEquiv(t, expectedField, g.fields["food"])
}

func TestInvalidFieldType(t *testing.T) {
	g := NewGenerator("thing", false)
	ExpectsError(t, fmt.Sprintf("Invalid field type '%s'", "foo"),
		g.WithField("login", "foo", 2, nil, false))
}

func TestWithFieldThrowsErrorOnBadFieldArgs(t *testing.T) {
	var testFields = []struct {
		fieldType   string
		badArgsType interface{}
	}{
		{"string", "string"},
		{"integer", "string"},
		{"decimal", "string"},
		{"date", "string"},
		{"enum", "string"},
		{"dict", 0},
	}

	g := NewGenerator("thing", false)

	for _, field := range testFields {
		ExpectsError(t, "expected field args to be of type", g.WithField("fieldName", field.fieldType, field.badArgsType, nil, false))
	}
}

func TestGenerateProducesGeneratedContent(t *testing.T) {
	g := NewGenerator("thing", false)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("a", "string", int64(2), nil, false)
	g.WithField("b", "integer", [2]int64{2, 4}, nil, false)
	g.WithField("c", "decimal", [2]float64{2.85, 4.50}, nil, false)
	g.WithField("d", "date", [2]time.Time{timeMin, timeMax}, nil, false)
	g.WithField("e", "dict", "last_name", nil, false)
	g.WithField("f", "mongoid", "", nil, false)
	g.WithField("g", "enum", collection("a", "b"), nil, false)
	g.WithEntityField("h", NewGenerator("thang", false), false, nil)
	g.WithField("i", "serial", nil, nil, false)

	emitter := NewTestEmitter()
	data := g.Generate(3, emitter)
	emitter.Shift()
	entity := emitter.Shift()

	AssertEqual(t, 3, len(data))

	testFields := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}

	for _, fieldName := range testFields {
		fieldValue, ok := entity[fieldName]
		Assert(t, ok, "entity should have field %q", fieldName)

		switch fieldType := fieldValue.(type) {
		case int64:
			Assert(t, fieldName == "b", "field %q should have yielded a int64", fieldName)
		case float64:
			Assert(t, fieldName == "c", "field %q should have yielded a float64", fieldName)
		case string:
			Assert(t, strings.Contains("a, e, f, g, h", fieldName), "field %q should have yielded a string", fieldName)
		case time.Time:
			Assert(t, fieldName == "d", "field %q should have yielded a Time", fieldName)
		case uint64:
			Assert(t, fieldName == "i", "field %q should have yielded a int", fieldName)
		default:
			Assert(t, false, "Don't know what to do with the field type for %q! The type is %v", fieldName, fieldType)
		}
	}
}

func TestGenerateWithBoundsArgumentProducesCorrectCountOfValues(t *testing.T) {
	g := NewGenerator("thing", false)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithEntityField("a", NewGenerator("subthing", false), 1, &CountRange{1, 1})
	g.WithField("b", "string", int64(2), &CountRange{2, 2}, false)
	g.WithField("c", "integer", [2]int64{2, 4}, &CountRange{3, 3}, false)
	g.WithField("d", "decimal", [2]float64{2.85, 4.50}, &CountRange{4, 4}, false)
	g.WithField("e", "date", [2]time.Time{timeMin, timeMax}, &CountRange{5, 5}, false)
	g.WithField("f", "dict", "last_name", &CountRange{6, 6}, false)
	g.WithField("g", "enum", collection("a", "b"), &CountRange{7, 7}, false)

	emitter := NewTestEmitter()
	g.Generate(1, emitter)
	emitter.Shift()
	entity := emitter.Shift()

	var testFields = []struct {
		fieldName string
		count     int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
		{"e", 5},
		{"f", 6},
		{"g", 7},
	}

	for _, field := range testFields {
		actual := len(entity[field.fieldName].([]interface{}))
		AssertEqual(t, field.count, actual)
	}
}

func TestEnsureGeneratable(t *testing.T) {
	g := NewGenerator("thing", false)
	g.WithField("eek", "integer", [2]int64{2, 4}, nil, true)
	ExpectsError(t, "Not enough unique values for field 'eek': There are only 3 unique values available for the 'eek' field, and you're trying to generate 5 entities", g.EnsureGeneratable(5))
}

func TestEnsureGeneratableInfinitePossibilitiesFieldType(t *testing.T) {
	g := NewGenerator("thing", false)
	g.WithField("eek", "float", [2]int64{2.0, 4.0}, nil, true)
	AssertNil(t, g.EnsureGeneratable(55), "There should be infinite number of possible float values")
}
