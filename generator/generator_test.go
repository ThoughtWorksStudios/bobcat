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
