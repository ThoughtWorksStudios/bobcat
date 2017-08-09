package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"github.com/satori/go.uuid"
	"reflect"
	"testing"
	"time"
	. "github.com/ThoughtWorksStudios/bobcat/common"
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

func TestExtendGenerator(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)

	g.WithField("name", "string", 10, nil)
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
	g.WithField("name", "string", 10, nil)
	g.WithEntityField("pet", subentityGenerator, 1, nil)

	entities := g.Generate(3)
	person_id := entities[0]["$id"]
	cat_parent := entities[0]["pet"].(map[string]GeneratedEntities)["Cat"][0]["$parent"]

	if person_id != cat_parent {
		t.Errorf("Parent id (%v) on subentity does not match the parent entity's id (%v)", cat_parent, person_id)
	}
}

func TestWithFieldCreatesCorrectFields(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("login", "string", 2, nil)
	g.WithField("age", "integer", [2]int{2, 4}, nil)
	g.WithField("stars", "decimal", [2]float64{2.85, 4.50}, nil)
	g.WithField("dob", "date", [2]time.Time{timeMin, timeMax}, nil)

	expectedFields := []struct {
		fieldName string
		field     Field
	}{
		{"login", &StringField{2, nil}},
		{"age", &IntegerField{2, 4, nil}},
		{"stars", &FloatField{2.85, 4.50, nil}},
		{"dob", &DateField{timeMin, timeMax, nil}},
		{"$id", &UuidField{}},
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
	err := g.WithField("age", "integer", [2]int{4, 2}, nil)
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
	expectedField := &LiteralField{"something", nil}
	if !equiv(expectedField, g.fields["login"]) {
		t.Errorf("Field 'login' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
			expectedField, g.fields["login"])
	}
}

func TestWithEntityFieldCreatesCorrectField(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	bound := &Bound{3, 3}
	g.WithEntityField("food", g, 3, bound)
	expectedField := &EntityField{g, bound}
	if !equiv(expectedField, g.fields["food"]) {
		t.Errorf("Field 'food' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
			expectedField, g.fields["food"])
	}
}

func TestInvalidFieldType(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	ExpectsError(t, fmt.Sprintf("Invalid field type '%s'", "foo"),
		g.WithField("login", "foo", 2, nil))
}

func TestFieldArgsCantBeNil(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	ExpectsError(t, "FieldArgs are nil for field 'login', this should never happen!",
		g.WithField("login", "foo", nil, nil))
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
	g.WithField("a", "string", 2, nil)
	g.WithField("b", "integer", [2]int{2, 4}, nil)
	g.WithField("c", "decimal", [2]float64{2.85, 4.50}, nil)
	g.WithField("d", "date", [2]time.Time{timeMin, timeMax}, nil)
	g.WithField("e", "dict", "last_name", nil)
	g.WithField("f", "uuid", "", nil)

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

func TestGenerateWithBoundsArgumentProducesCorrectAmountOfValues(t *testing.T) {
	data := GeneratedEntities{}
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("a", "string", 2, &Bound{2,2})
	g.WithField("b", "integer", [2]int{2, 4}, &Bound{3,3})
	g.WithField("c", "decimal", [2]float64{2.85, 4.50}, &Bound{4,4})
	g.WithField("d", "date", [2]time.Time{timeMin, timeMax}, &Bound{5,5})
	g.WithField("e", "dict", "last_name", &Bound{6,6})
	g.WithEntityField("f", NewGenerator("subthing", logger), 1, &Bound{7,7})

	data = g.Generate(1)

	var testFields = []struct {
		fieldName string
		amount int
	}{
		{"a", 2},
		{"b", 3},
		{"c", 4},
		{"d", 5},
		{"e", 6},
		{"f", 7},
	}

	entity := data[0]
	for _, field := range testFields {
		actual := len(entity[field.fieldName].([]interface{}))
		AssertEqual(t, field.amount, actual)
	}
}
