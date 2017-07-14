package generator

import (
	"fmt"
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

	if err == nil {
		t.Error("Field 'age' had invalid range, error was returnedgged")
	}
}

func TestDecimalRangeIsCorrect(t *testing.T) {
	g := NewGenerator("thing")
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	err := g.WithField("dob", "date", [2]time.Time{timeMax, timeMin})

	if err == nil {
		t.Error("Field 'dob' had invalid range, but error was returned")
	}
}

func TestDateRangeIsCorrect(t *testing.T) {
	g := NewGenerator("thing")
	err := g.WithField("stars", "decimal", [2]float64{4.4, 2.0})

	if err == nil {
		t.Error("Field 'stars' had invalid range, but error was returned")
	}
}

func TestDuplicatedFieldIsLogged(t *testing.T) {
	g := NewGenerator("thing")
	g.WithField("login", "string", 2)
	err := g.WithField("login", "string", 2)

	if err == nil {
		t.Error("Field 'login' duplicated, but not logged")
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
		t.Error("Static field 'login' duplicated, but no error was returned")
	}
}

func TestInvalidFieldTypeIsLogged(t *testing.T) {
	g := NewGenerator("thing")
	err := g.WithField("login", "foo", 2)

	if err == nil {
		t.Error("Invalid field type 'foo'")
	}
}

func TestFieldOptsCantBeNil(t *testing.T) {
	g := NewGenerator("thing")
	err := g.WithField("login", "foo", nil)

	if err == nil {
		t.Error("Expected an error when fieldOpts are nil, but did not receive it")
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
	var fileCreated string

	saved := writeToFile
	defer func() { writeToFile = saved }()

	writeToFile = func(payload []byte, filename string) {
		fileCreated = filename
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
	g.Generate(1)

	if fileCreated != "thing.json" {
		t.Errorf("Did not write JSON to file (with correct file name)")
	}
}
