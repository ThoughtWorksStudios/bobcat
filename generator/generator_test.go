package generator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

/*
 * is this a cheap hack? you bet it is.
 */
func equiv(expected, actual Field) bool {
	return fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", actual)
}

func TestWithFieldCreatesCorrectFields(t *testing.T) {
	g := NewGenerator("thing")
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

func TestIntegerRangeIsCorrect(t *testing.T) {
	g := NewGenerator("thing")
	err := g.WithField("age", "integer", [2]int{4, 2})
	expected := fmt.Sprintf("max %d cannot be less than min %d", 2, 4)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDateRangeIsCorrect(t *testing.T) {
	g := NewGenerator("thing")
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	err := g.WithField("dob", "date", [2]time.Time{timeMax, timeMin})
	expected := fmt.Sprintf("max %s cannot be before min %s", timeMin, timeMax)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDecimalRangeIsCorrect(t *testing.T) {
	g := NewGenerator("thing")
	err := g.WithField("stars", "decimal", [2]float64{4.4, 2.0})
	expected := fmt.Sprintf("max %v cannot be less than min %v", 2.0, 4.4)
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestDuplicatedFieldIsLogged(t *testing.T) {
	g := NewGenerator("thing")
	g.WithField("login", "string", 2)
	err := g.WithField("login", "string", 2)
	if err == nil {
		t.Errorf("Expected warning, but got none")
	} else if _, ok := err.(WarningError); ok {
		t.Errorf("Expected warning, but got none")
	} else if err.Error() != "already defined field: login" {
		t.Errorf("Expected already defined field error, but received %v", err.Error())
	}
}

func TestWithStaticFieldCreatesCorrectField(t *testing.T) {
	g := NewGenerator("thing")
	g.WithStaticField("login", "something")
	expectedField := &LiteralField{"something"}
	if !equiv(expectedField, g.fields["login"]) {
		t.Errorf("Field 'login' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
			expectedField, g.fields["login"])
	}
}

func TestDuplicatedStaticFieldIsLogged(t *testing.T) {
	g := NewGenerator("thing")
	g.WithStaticField("login", "something")
	err := g.WithStaticField("login", "other")
	if err == nil {
		t.Errorf("Expected warning, but got none")
	} else if _, ok := err.(WarningError); ok {
		t.Errorf("Expected warning, but got none")
	} else if err.Error() != "already defined field: login" {
		t.Errorf("Expected already defined field error, but received %v", err.Error())
	}
}

func TestInvalidFieldTypeIsLogged(t *testing.T) {
	g := NewGenerator("thing")
	err := g.WithField("login", "foo", 2)
	expected := fmt.Sprintf("Invalid field type '%s'", "foo")
	if err == nil || err.Error() != expected {
		t.Errorf("expected error: %v\n but got %v", expected, err)
	}
}

func TestFieldOptsCantBeNil(t *testing.T) {
	g := NewGenerator("thing")
	err := g.WithField("login", "foo", nil)
	if err == nil || err.Error() != "FieldOpts are nil for field 'login', this should never happen!" {
		t.Errorf("expected error")
	}

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

	g := NewGenerator("thing")

	for _, field := range testFields {
		err := g.WithField("fieldName", field.fieldType, field.fieldOpts)
		if err == nil {
			t.Errorf("Mismatched field opts type for field type '%s' should be logged", field.fieldType)
		}
	}
}

func TestGenerateProducesCorrectJSON(t *testing.T) {
	var fileOutput []byte

	saved := writeToFile
	defer func() { writeToFile = saved }()

	writeToFile = func(payload []byte, filename string) {
		fileOutput = payload
	}

	g := NewGenerator("thing")
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
	g.Generate(3)

	var data []map[string]interface{}
	json.Unmarshal(fileOutput, &data)

	if len(data) != 3 {
		t.Errorf("Did not generate the appropriate number of entities")
	}

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
		if expected != actual {
			t.Errorf("Field type of '%v' is not correct: '%v' not '%v'", field.fieldName, actual, expected)
		}
	}
}
