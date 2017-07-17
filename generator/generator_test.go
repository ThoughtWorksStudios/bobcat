package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
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
	g := NewGenerator("thing", nil)
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
	log := GetLogger()
	g := NewGenerator("thing", log)
	g.WithField("age", "integer", [2]int{4, 2})
	log.AssertMessage(t, "max %d cannot be less than min %d\n", 2, 4)
}

func TestDecimalRangeIsCorrect(t *testing.T) {
	log := GetLogger()
	g := NewGenerator("thing", log)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithField("dob", "date", [2]time.Time{timeMax, timeMin})
	log.AssertMessage(t, "max %s cannot be before min %s\n", timeMin, timeMax)
}

func TestDateRangeIsCorrect(t *testing.T) {
	log := GetLogger()
	g := NewGenerator("thing", log)
	g.WithField("stars", "decimal", [2]float64{4.4, 2.0})
	log.AssertMessage(t, "max %d cannot be less than min %d\n", 2.0, 4.4)
}

func TestDuplicatedFieldIsLogged(t *testing.T) {
	log := GetLogger()
	g := NewGenerator("thing", log)
	g.WithField("login", "string", 2)
	g.WithField("login", "string", 2)
	log.AssertWarning(t, "already defined field: %s", "login")
}

func TestWithStaticFieldCreatesCorrectField(t *testing.T) {
	g := NewGenerator("thing", nil)
	g.WithStaticField("login", "something")
	expectedField := &LiteralField{"something"}
	if !equiv(expectedField, g.fields["login"]) {
		t.Errorf("Field 'login' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]",
			expectedField, g.fields["login"])
	}
}

func TestDuplicatedStaticFieldIsLogged(t *testing.T) {
	log := GetLogger()
	g := NewGenerator("thing", log)
	g.WithStaticField("login", "something")
	g.WithStaticField("login", "other")
	log.AssertWarning(t, "already defined field: %s", "login")
}

func TestInvalidFieldTypeIsLogged(t *testing.T) {
	log := GetLogger()
	g := NewGenerator("thing", log)
	g.WithField("login", "foo", 2)
	log.AssertMessage(t, "Invalid field type '%s'", "foo")
}

func TestFieldOptsCantBeNil(t *testing.T) {
	log := GetLogger()
	g := NewGenerator("thing", log)
	g.WithField("login", "foo", nil)
	log.AssertMessage(t, "FieldOpts are nil for field '%s', this should never happen!", "login")

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

	log := GetLogger()
	g := NewGenerator("thing", log)

	for _, field := range testFields {
		g.WithField("fieldName", field.fieldType, field.fieldOpts)
		switch field.fieldType {
		case "string":
			log.AssertMessage(t, "expected field options to be of type 'int' for field %s (%s), but got %v", "fieldName", field.fieldType, field.fieldOpts)
		case "integer":
			log.AssertMessage(t, "expected field options to be of type '(min:int, max:int)' for field %s (%s), but got %v", "fieldName", field.fieldType, field.fieldOpts)
		case "decimal":
			log.AssertMessage(t, "expected field options to be of type '(min:float64, max:float64)' for field %s (%s), but got %v", "fieldName", field.fieldType, field.fieldOpts)
		case "date":
			log.AssertMessage(t, "expected field options to be of type 'time.Time' for field %s (%s), but got %v", "fieldName", field.fieldType, field.fieldOpts)
		case "dict":
			log.AssertMessage(t, "expected field options to be of type 'string' for field %s (%s), but got %v", "fieldName", field.fieldType, field.fieldOpts)
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

	g := NewGenerator("thing", nil)
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
