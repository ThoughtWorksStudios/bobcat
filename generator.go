package main

// import "encoding/json"
// import "github.com/Pallinder/go-randomdata"
import "fmt"

type Generator struct {
  name string
  fields map[string]interface{}
}

func NewGenerator(name string) Generator {
  return Generator{name: name, fields: make(map[string]interface{})}
}

type StringField struct {
  length int
}

type IntegerField struct {
  min int
  max int
}

type FloatField struct {
}

func (g Generator) withField(fieldName string, fieldType string, fieldOpts interface{}) Generator {
  if _, ok := g.fields[fieldName]; ok {
    fmt.Printf("already defined field %s", fieldName)
  }

  switch fieldType {
    case "string":
      len, ok := fieldOpts.(int)
      if ok {
        g.fields[fieldName] = StringField {length: len}
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
        g.fields[fieldName] = IntegerField {min: min, max: max}
      } else {
        expectsType("(min:int, max:int)", fieldName, fieldType, fieldOpts)
      }
  }

  return g
}

func expectsType(expectedType string, fieldName string, fieldType string, fieldOpts interface{}) {
  fmt.Println("expected options to be ", expectedType, " for field ", fieldName, " (", fieldType, ")")
}
