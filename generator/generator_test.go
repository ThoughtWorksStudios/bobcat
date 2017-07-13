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

	expectedFields := make(map[string]Field)
	expectedFields["login"] = &StringField{2}
	expectedFields["age"] = &IntegerField{2, 4}
	expectedFields["stars"] = &FloatField{2.85, 4.50}
	expectedFields["dob"] = &DateField{timeMin, timeMax}
	expectedFields["boo"] = &DictField{"silly_name"}

	for field, value := range expectedFields {
		if !equiv(g.fields[field], value) {
			t.Errorf("Field '%s' does have appropriate value. \n Expected: \n [%v] \n\n but generated: \n [%v]", field, value, g.fields[field])
		}
	}
}

func TestDuplicatedFieldIsLogged(t *testing.T) {
	saved := inform
	defer func() { inform = saved }()

	messageLogged := false

	inform = func(message string, values ...interface{}) {
		messageLogged = true
	}

	g := NewGenerator("thing")
	g.WithField("login", "string", 2)
	g.WithField("login", "string", 2)

	if !messageLogged {
		t.Error("Field 'login' duplicated, but not logged")
	}
}

func TestDuplicatedStaticFieldIsLogged(t *testing.T) {
	saved := inform
	defer func() { inform = saved }()

	messageLogged := false

	inform = func(message string, values ...interface{}) {
		messageLogged = true
	}

	g := NewGenerator("thing")
	g.WithStaticField("login", "something")
	g.WithStaticField("login", "other")

	if !messageLogged {
		t.Error("Static field 'login' duplicated, but not logged")
	}
}

func TestInvalidFieldTypeIsLogged(t *testing.T) {
	saved := inform
	defer func() { inform = saved }()

	messageLogged := false

	inform = func(message string, values ...interface{}) {
		messageLogged = true
	}

	g := NewGenerator("thing")
	g.WithField("login", "foo", 2)

	if !messageLogged {
		t.Error("Invalid field type 'foo'")
	}
}

func TestFieldOptsCantBeNil(t *testing.T) {
	g := NewGenerator("thing")
	_, error := g.WithField("login", "foo", nil)

	if error == nil {
		t.Error("Expected an error when fieldOpts are nil, but did not receive it")
	}
}

func TestFieldOptsMatchesFieldType(t *testing.T) {
	saved := inform
	defer func() { inform = saved }()

	messageLogged := false

	inform = func(message string, values ...interface{}) {
		messageLogged = true
	}

	var testFields = []struct {
		fieldType string
		fieldOpts interface{}
	} {
		{"string", "string" },
		{"integer", "string" },
		{"decimal", "string" },
		{"date", "string" },
		{"dict", 0 },
	}

	g := NewGenerator("thing")

	for _, field := range testFields {
		messageLogged = false
		g.WithField("fieldName", field.fieldType, field.fieldOpts)

		if !messageLogged {
			t.Errorf("Mismatched field opts type for field type '%s' should be logged", field.fieldType)
		}
	}
}



