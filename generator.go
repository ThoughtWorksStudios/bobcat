package main

import "encoding/json"

import "fmt"

type Generator struct {
	name   string
	fields map[string]Field
}

func NewGenerator(name string) *Generator {
	return &Generator{name: name, fields: make(map[string]Field)}
}

func (g *Generator) withField(fieldName, fieldType string, fieldOpts interface{}) *Generator {
	if _, ok := g.fields[fieldName]; ok {
		fmt.Printf("already defined field %s\n", fieldName)
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
			fmt.Printf("max %d cannot be less than min %d\n", max, min)
		}

		if ok {
			g.fields[fieldName] = IntegerField{min: min, max: max}
		} else {
			expectsType("(min:int, max:int)", fieldName, fieldType, fieldOpts)
		}
	case "date":
		bounds, ok := fieldOpts.([2]string)
		min, max := bounds[0], bounds[1]

		if ok {
			field := DateField{min: min, max: max}
			if !field.ValidBounds() {
				fmt.Printf("max %s cannot be before %s\n", max, min)
			}
			g.fields[fieldName] = field
		} else {
			expectsType("string", fieldName, fieldType, fieldOpts)
		}

	}

	return g
}

func expectsType(expectedType, fieldName, fieldType string, fieldOpts interface{}) {
	fmt.Println("expected options to be ", expectedType, " for field ", fieldName, " (", fieldType, ")")
}

func (g *Generator) generate(count int) string {

	result := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {

		obj := make(map[string]interface{})
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
		withField("DOB", "date", [2]string{"2010-01-04", "2015-01-04"}).
		withField("weight", "integer", [2]int{100, 200})
	fmt.Println(person.generate(10))

}
