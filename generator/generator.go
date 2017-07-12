package generator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Generator struct {
	name   string
	fields map[string]Field
}

func NewGenerator(name string) *Generator {
	return &Generator{name: name, fields: make(map[string]Field)}
}

// For testing purposes
func (g *Generator) GetField(name string) Field {
	return g.fields[name]
}

func (g *Generator) WithStaticField(fieldName string, fieldValue interface{}) *Generator {
	if _, ok := g.fields[fieldName]; ok {
		log.Fatalln("already defined field: ", fieldName)
	}

	g.fields[fieldName] = &LiteralField{value: fieldValue}

	return g
}

func (g *Generator) WithField(fieldName, fieldType string, fieldOpts interface{}) *Generator {
	if _, ok := g.fields[fieldName]; ok {
		log.Fatalln("already defined field: ", fieldName)
	}

	switch fieldType {
	case "string":
		if ln, ok := fieldOpts.(int); ok {
			g.fields[fieldName] = &StringField{length: ln}
		} else {
			expectsType("int", fieldName, fieldType, fieldOpts)
		}
	case "integer":
		if bounds, ok := fieldOpts.([2]int); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				fmt.Printf("max %d cannot be less than min %d\n", max, min)
			}

			g.fields[fieldName] = &IntegerField{min: min, max: max}
		} else {
			expectsType("(min:int, max:int)", fieldName, fieldType, fieldOpts)
		}
	case "decimal":
		if bounds, ok := fieldOpts.([2]float64); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				fmt.Printf("max %d cannot be less than min %d\n", max, min)
			}
			g.fields[fieldName] = &FloatField{min: min, max: max}
		} else {
			expectsType("(min:float64, max:float64)", fieldName, fieldType, fieldOpts)
		}
	case "date":
		if bounds, ok := fieldOpts.([2]time.Time); ok {
			min, max := bounds[0], bounds[1]
			field := &DateField{min: min, max: max}
			if !field.ValidBounds() {
				fmt.Printf("max %s cannot be before min %s\n", max, min)
			}
			g.fields[fieldName] = field
		} else {
			expectsType("time.Time", fieldName, fieldType, fieldOpts)
		}
	case "dict":
		if dict, ok := fieldOpts.(string); ok {
			g.fields[fieldName] = &DictField{category: dict}
		} else {
			expectsType("string", fieldName, fieldType, fieldOpts)
		}
	}
	return g
}

func expectsType(expectedType, fieldName, fieldType string, fieldOpts interface{}) {
	fmt.Println("expected options to be ", expectedType, " for field ", fieldName, " (", fieldType, ")")
}

func (g *Generator) Generate(count int64) {

	result := make([]map[string]interface{}, count)
	for i := int64(0); i < count; i++ {

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
