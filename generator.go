package main

import "encoding/json"

import "github.com/Pallinder/go-randomdata"
import "fmt"

type Generator struct {
	name   string
	fields map[string]Field
}

func NewGenerator(name string) *Generator {
	return &Generator{name: name, fields: make(map[string]Field)}
}

type Field interface {
	Type() string
	GenerateValue() interface{}
}

type StringField struct {
	length int
}

func (field StringField) Type() string {
	return "string"
}

func (field StringField) GenerateValue() interface{} {
	return randomdata.RandStringRunes(field.length)
}

type IntegerField struct {
	min int
	max int
}

func (field IntegerField) Type() string {
	return "integer"
}

func (field IntegerField) GenerateValue() interface{} {
	return randomdata.Number(field.min, field.max)
}

type FloatField struct {
}

func (g *Generator) withField(fieldName string, fieldType string, fieldOpts interface{}) *Generator {
	if _, ok := g.fields[fieldName]; ok {
		fmt.Printf("already defined field %s", fieldName)
	}

	switch fieldType {
	case "string":
		len, ok := fieldOpts.(int)
		if ok {
			g.fields[fieldName] = StringField{length: len}
		} else {
			expectsType("int", fieldName, fieldType, fieldOpts)
		}
	case "integer":
		bounds, ok := fieldOpts.([2]int)
		min, max := bounds[0], bounds[1]
		if max < min {
			fmt.Printf("max %d cannot be less than min %d", max, min)
		}

		if ok {
			g.fields[fieldName] = IntegerField{min: min, max: max}
		} else {
			expectsType("(min:int, max:int)", fieldName, fieldType, fieldOpts)
		}
	}

	return g
}

func expectsType(expectedType string, fieldName string, fieldType string, fieldOpts interface{}) {
	fmt.Println("expected options to be ", expectedType, " for field ", fieldName, " (", fieldType, ")")
}

func NewObj() map[string]interface{} {
	return make(map[string]interface{})
}

func (g *Generator) generate(count int) string {

	result := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {

		obj := NewObj()
		for name, field := range g.fields {
			obj[name] = field.GenerateValue()
		}
		result[i] = obj
	}
	marsh, _ := json.Marshal(result)
	return string(marsh)
}

func TestThis() {
	person := NewGenerator("Person").
		withField("name", "string", 10).
		withField("age", "integer", [2]int{5, 70}).
		withField("weight", "integer", [2]int{100, 200})
	fmt.Println(person.generate(10))

}
