package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"github.com/satori/go.uuid"
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
func equiv(expected, actual Field) bool {
	return fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", actual)
}

func defaultFieldCount() Range { return Range{1,1} }

func TestExtendGenerator(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)

	g.WithField("name", "string", 10, defaultFieldCount())
	g.WithField("age", "decimal", [2]float64{2, 4}, defaultFieldCount())
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
	subentityGenerator.WithField("name", "string", 5, defaultFieldCount())

	generator := NewGenerator("Person", logger)
	generator.WithField("name", "string", 10, defaultFieldCount())
	generator.WithEntityField("pet", "entity", subentityGenerator, defaultFieldCount())

	entities := generator.Generate(3)
	person_id := entities[0]["$id"]
	cat_parent := entities[0]["pet"].(map[string]interface{})["Cat"].(GeneratedEntities)[0]["$parent"]

	if person_id != cat_parent {
		t.Errorf("Parent id (%v) on subentity does not match the parent entity's id (%v)", cat_parent, person_id)
	}
}

func TestWithFieldCreatesCorrectFields(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("login", "string", 2, Range{2,2})
	g.WithField("age", "integer", [2]int{2, 4}, defaultFieldCount())
	g.WithField("stars", "decimal", [2]float64{2.85, 4.50}, defaultFieldCount())
	g.WithField("dob", "date", [2]time.Time{timeMin, timeMax}, defaultFieldCount())

	expectedFields := []struct {
		fieldName string
		field     Field
	}{
		{"login", &StringField{2, Range{2,2}}},
		{"age", &IntegerField{2, 4, defaultFieldCount()}},
		{"stars", &FloatField{2.85, 4.50, defaultFieldCount()}},
		{"dob", &DateField{timeMin, timeMax, defaultFieldCount()}},
		{"$id", &UuidField{defaultFieldCount()}},
	}

	for _, expectedField := range expectedFields {
		if !equiv(expectedField.field, g.fields[expectedField.fieldName]) {
			t.Errorf("Field '%s' does not have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
				expectedField.fieldName, expectedField.field, g.fields[expectedField.fieldName])
		}
	}
}

func TestIntegerRangeIsCorrect(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	err := g.WithField("age", "integer", [2]int{4, 2}, defaultFieldCount())
	expected := fmt.Sprintf("max %d cannot be less than min %d", 2, 4)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDateRangeIsCorrect(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	err := g.WithField("dob", "date", [2]time.Time{timeMax, timeMin}, defaultFieldCount())
	expected := fmt.Sprintf("max %s cannot be before min %s", timeMin, timeMax)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDecimalRangeIsCorrect(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	err := g.WithField("stars", "decimal", [2]float64{4.4, 2.0}, defaultFieldCount())
	expected := fmt.Sprintf("max %v cannot be less than min %v", 2.0, 4.4)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDuplicateFieldIsLogged(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)

	AssertNil(t, g.WithField("login", "string", 2, defaultFieldCount()), "Should not return an error")
	AssertNil(t, g.WithField("login", "string", 5, defaultFieldCount()), "Should not return an error")

	logger.AssertWarning("Field thing.login is already defined; overriding to string(5)")
}

func TestDuplicateStaticFieldIsLogged(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)

	AssertNil(t, g.WithStaticField("login", "something"), "Should not return an error")
	AssertNil(t, g.WithStaticField("login", "other"), "Should not return an error")

	logger.AssertWarning("Field thing.login is already defined; overriding to other")
}

func TestWithStaticFieldCreatesCorrectField(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	g.WithStaticField("login", "something")
	expectedField := &LiteralField{"something", defaultFieldCount()}
	if !equiv(expectedField, g.fields["login"]) {
		t.Errorf("Field 'login' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
			expectedField, g.fields["login"])
	}
}

func TestWithEntityFieldCreatesCorrectField(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	g.WithEntityField("food", "Peanut", g, Range{3,3})
	expectedField := &EntityField{g, Range{3,3}}
	if !equiv(expectedField, g.fields["food"]) {
		t.Errorf("Field 'food' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
			expectedField, g.fields["food"])
	}
}

func TestInvalidFieldType(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	ExpectsError(t, fmt.Sprintf("Invalid field type '%s'", "foo"),
		g.WithField("login", "foo", 2, defaultFieldCount()))
}

func TestFieldOptsCantBeNil(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	ExpectsError(t, "FieldOpts are nil for field 'login', this should never happen!",
		g.WithField("login", "foo", nil, defaultFieldCount()))
}

func TestFieldOptsMatchesFieldType(t *testing.T) {
	var testFields = []struct {
		fieldType string
		fieldOpts interface{}
	}{
		{"string", "string"},
		{"integer", "string"},
		{"decimal", "string"},
		{"date", "string"},
		{"dict", 0},
	}

	logger := GetLogger(t)
	g := NewGenerator("thing", logger)

	for _, field := range testFields {
		AssertNotNil(t, g.WithField("fieldName", field.fieldType, field.fieldOpts, defaultFieldCount()),
			"Mismatched field opts type for field type '%s' should be logged", field.fieldType)
	}
}

func TestGenerateProducesGeneratedContent(t *testing.T) {
	data := GeneratedEntities{}
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("a", "string", 2, defaultFieldCount())
	g.WithField("b", "integer", [2]int{2, 4}, defaultFieldCount())
	g.WithField("c", "decimal", [2]float64{2.85, 4.50}, defaultFieldCount())
	g.WithField("d", "date", [2]time.Time{timeMin, timeMax}, defaultFieldCount())
	g.WithField("e", "dict", "last_name", defaultFieldCount())
	g.WithField("f", "uuid", "", defaultFieldCount())

	data = g.Generate(3)

	AssertEqual(t, 3, len(data))

	var testFields = []struct {
		fieldName string
		fieldType interface{}
	}{
		{"a", "string"},
		{"b", 1},
		{"c", 2.1},
		{"d", time.Time{}},
		{"e", "string"},
		{"f", uuid.NewV4()},
	}

	entity := data[0]
	for _, field := range testFields {
		actual := reflect.TypeOf(entity[field.fieldName])
		expected := reflect.TypeOf(field.fieldType)
		AssertEqual(t, expected, actual)
	}
}
