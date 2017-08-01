package generator

import (
	"github.com/ThoughtWorksStudios/fake"
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
)

type Field interface {
	Type() string
	GenerateValue() interface{}
}

type FieldSet map[string]Field

type ReferenceField struct {
	referred  *Generator
	fieldName string
}

func (field *ReferenceField) Type() string {
	return "reference"
}

func (field *ReferenceField) GenerateValue() interface{} {
	referredField := field.referred.fields[field.fieldName]
	return referredField.GenerateValue()
}

func (field *ReferenceField) referencedField() Field {
	f := field.referred.fields[field.fieldName]
	if f.Type() == field.Type() {
		return f.(*ReferenceField).referencedField()
	} else {
		return f
	}
}

type EntityField struct {
	entityGenerator *Generator
	count           int
}

func (field *EntityField) Type() string {
	return "entity"
}

func (field *EntityField) GenerateValue() interface{} {
	entities := make(map[string]GeneratedEntities)
	entities[field.entityGenerator.Name] = field.entityGenerator.Generate(int64(field.count))
	return entities
}

type UuidField struct{}

func (field *UuidField) Type() string {
	return "uuid"
}

func (field *UuidField) GenerateValue() interface{} {
	return uuid.NewV4()
}

type LiteralField struct {
	value interface{}
}

func (field *LiteralField) Type() string {
	return "literal"
}

func (field *LiteralField) GenerateValue() interface{} {
	return field.value
}

type StringField struct {
	length int
}

func (field *StringField) Type() string {
	return "string"
}

func (field *StringField) GenerateValue() interface{} {
	return fake.Password(field.length, field.length, true, true, true)
}

type IntegerField struct {
	min int
	max int
}

func (field *IntegerField) Type() string {
	return "integer"
}

func (field *IntegerField) GenerateValue() interface{} {
	result := float64(rand.Intn(int(field.max - field.min)))
	result += float64(field.min)
	return int(result)
}

type FloatField struct {
	min float64
	max float64
}

func (field *FloatField) Type() string {
	return "float"
}

func (field *FloatField) GenerateValue() interface{} {
	return float64(rand.Intn(int(field.max-field.min))) + field.min + rand.Float64()
}

type DateField struct {
	min time.Time
	max time.Time
}

func (field *DateField) Type() string {
	return "date"
}

func (field *DateField) ValidBounds() bool {
	return field.min.Before(field.max)
}

func (field *DateField) GenerateValue() interface{} {
	min, max := field.min.Unix(), field.max.Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min

	return time.Unix(sec, 0)
}

type DictField struct {
	category string
}

var CustomDictPath = ""

func (field *DictField) Type() string {
	return "dict"
}

func (field *DictField) GenerateValue() interface{} {
	//TODO: Figure out how to get rid of the switch statement bellow
	switch field.category {
	case "addresses":
		return fake.StreetAddress()
	case "emails":
		return fake.EmailAddress()
	default:
		fake.SetCustomDataLocation(CustomDictPath)
		return fake.ValueFromDictionary(field.category)
	}
}
