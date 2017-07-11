package generator

import "encoding/json"
import "fmt"
import "os"

type Generator struct {
	name   string
	fields map[string]Field
}

func NewGenerator(name string) *Generator {
	return &Generator{name: name, fields: make(map[string]Field)}
}

func (g *Generator) WithField(fieldName, fieldType string, fieldOpts interface{}) *Generator {
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
	case "float":
		bounds, ok := fieldOpts.([2]float64)
		min, max := bounds[0], bounds[1]
		if max < min {
			fmt.Printf("max %d cannot be less than min %d\n", max, min)
		}

		if ok {
			g.fields[fieldName] = FloatField{min: min, max: max}
		} else {
			expectsType("(min:float64, max:float64)", fieldName, fieldType, fieldOpts)
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
	case "dict":
		dict, ok := fieldOpts.(string)
		if ok {
			g.fields[fieldName] = DictField{category: dict}
		} else {
			expectsType("string", fieldName, fieldType, fieldOpts)
		}
	}
	return g
}

func expectsType(expectedType, fieldName, fieldType string, fieldOpts interface{}) {
	fmt.Println("expected options to be ", expectedType, " for field ", fieldName, " (", fieldType, ")")
}

func (g *Generator) Generate(count int) {

	result := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {

		obj := make(map[string]interface{})
		for name, field := range g.fields {
			obj[name] = field.GenerateValue()
		}
		result[i] = obj
	}
	marsh, _ := json.MarshalIndent(result, "", "\t")
	g.writeToFile(marsh)
}

func (g *Generator) filename() string {
	return fmt.Sprintf("%s.json", g.name)
}

func (g *Generator) writeToFile(json []byte) {
	dest, _ := os.Create(g.filename())
	defer dest.Close()
	dest.Write(json)
}

func TestThis() {
	person := NewGenerator("Person").
		WithField("first_name", "dict", "first_name").
		WithField("age", "integer", [2]int{5, 70}).
		WithField("DOB", "date", [2]string{"1945-01-04", "2010-01-04"}).
		WithField("email", "dict", "email").
		WithField("last_name", "dict", "last_name").
		WithField("zip_code", "dict", "zip_code").
		WithField("address", "dict", "address").
		WithField("weight", "float", [2]float64{100.2, 200.66})
	person.Generate(10)
}
