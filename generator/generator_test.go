package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
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

func TestExtendGenerator(t *testing.T) {
	logger := GetLogger(t)
	buffer := new(bytes.Buffer)
	g := NewGenerator("thing", logger)

	g.WithField("name", "string", 10)
	g.WithField("age", "integer", [2]int{2, 4})
	g.WithStaticField("species", "human")

	m := ExtendGenerator("thang", g)
	m.WithStaticField("species", "h00man")
	m.WithStaticField("name", "kyle")

	var data []map[string]interface{}

	g.Generate(1, buffer)
	json.Unmarshal(buffer.Bytes(), &data)

	base := data[0]

	AssertEqual(t, "human", base["species"])
	AssertEqual(t, 10, len(base["name"].(string)))
	Assert(t, isBetween(base["age"].(float64), 2, 4), "base entity failed to generate the correct age")

	buffer = new(bytes.Buffer)
	m.Generate(1, buffer)
	json.Unmarshal(buffer.Bytes(), &data)

	extended := data[0]
	AssertEqual(t, "h00man", extended["species"])
	AssertEqual(t, "kyle", extended["name"].(string))
	Assert(t, isBetween(extended["age"].(float64), 2, 4), "extended entity failed to generate the correct age")
}

func TestWithFieldCreatesCorrectFields(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("login", "string", 2)
	g.WithField("age", "integer", [2]int{2, 4})
	g.WithField("stars", "decimal", [2]float64{2.85, 4.50})
	g.WithField("dob", "date", [2]time.Time{timeMin, timeMax})
	g.WithField("boo", "dict", "silly_name")

	expectedFields := []struct {
		fieldName string
		field     Field
	}{
		{"login", &StringField{2}},
		{"age", &IntegerField{2, 4}},
		{"stars", &FloatField{2.85, 4.50}},
		{"dob", &DateField{timeMin, timeMax}},
		{"boo", &DictField{"silly_name"}},
	}

	for _, expectedField := range expectedFields {
		if !equiv(expectedField.field, g.fields[expectedField.fieldName]) {
			t.Errorf("Field '%s' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
				expectedField.fieldName, expectedField.field, g.fields[expectedField.fieldName])
		}
	}
}

func TestNameChangesWhenOverridingFieldType(t *testing.T) {
	logger := GetLogger(t)
	parent := NewGenerator("thing", logger)
	parent.WithField("age", "integer", [2]int{2, 4})

	child := ExtendGenerator("blah", parent)
	child.WithField("age", "decimal", [2]float64{2.0, 4.0})
	if child.GetName() == "blah" {
		t.Errorf("Expected generator name to change when field type changes")
	}

}

func TestIntegerRangeIsCorrect(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	err := g.WithField("age", "integer", [2]int{4, 2})
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
	err := g.WithField("dob", "date", [2]time.Time{timeMax, timeMin})
	expected := fmt.Sprintf("max %s cannot be before min %s", timeMin, timeMax)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDecimalRangeIsCorrect(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	err := g.WithField("stars", "decimal", [2]float64{4.4, 2.0})
	expected := fmt.Sprintf("max %v cannot be less than min %v", 2.0, 4.4)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDuplicateFieldIsLogged(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)

	AssertNil(t, g.WithField("login", "string", 2), "Should not return an error")
	AssertNil(t, g.WithField("login", "string", 5), "Should not return an error")

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
	expectedField := &LiteralField{"something"}
	if !equiv(expectedField, g.fields["login"]) {
		t.Errorf("Field 'login' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
			expectedField, g.fields["login"])
	}
}

func TestInvalidFieldType(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	ExpectsError(t, fmt.Sprintf("Invalid field type '%s'", "foo"), g.WithField("login", "foo", 2))
}

func TestFieldOptsCantBeNil(t *testing.T) {
	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	ExpectsError(t, "FieldOpts are nil for field 'login', this should never happen!", g.WithField("login", "foo", nil))
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
		AssertNotNil(t, g.WithField("fieldName", field.fieldType, field.fieldOpts), "Mismatched field opts type for field type '%s' should be logged", field.fieldType)
	}
}

func TestGenerateProducesCorrectJSON(t *testing.T) {
	buffer := new(bytes.Buffer)

	logger := GetLogger(t)
	g := NewGenerator("thing", logger)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("a", "string", 2)
	g.WithField("b", "integer", [2]int{2, 4})
	g.WithField("c", "decimal", [2]float64{2.85, 4.50})
	g.WithField("d", "date", [2]time.Time{timeMin, timeMax})
	g.WithField("e", "dict", "last_name")
	g.WithField("f", "dict", "first_name")
	g.WithField("g", "dict", "city")
	g.WithField("h", "dict", "country")
	g.WithField("i", "dict", "state")
	g.WithField("j", "dict", "street")
	g.WithField("k", "dict", "address")
	g.WithField("l", "dict", "email")
	g.WithField("m", "dict", "zip_code")
	g.WithField("n", "dict", "full_name")
	g.WithField("o", "dict", "random_string")
	g.WithField("p", "dict", "invalid_type")
	g.Generate(3, buffer)

	var data []map[string]interface{}
	json.Unmarshal(buffer.Bytes(), &data)

	AssertEqual(t, 3, len(data))

	var testFields = []struct {
		fieldName string
		fieldType interface{}
	}{
		{"a", "string"},
		{"b", 1.2},
		{"c", 2.1},
		{"d", "string"},
		{"e", "string"},
		{"f", "string"},
		{"g", "string"},
		{"h", "string"},
		{"i", "string"},
		{"j", "string"},
		{"k", "string"},
		{"l", "string"},
		{"m", "string"},
		{"n", "string"},
		{"o", "string"},
		{"p", nil},
	}

	entity := data[0]
	for _, field := range testFields {
		actual := reflect.TypeOf(entity[field.fieldName])
		expected := reflect.TypeOf(field.fieldType)
		AssertEqual(t, expected, actual)
	}
}
