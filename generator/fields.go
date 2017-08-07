package generator

import (
	"github.com/ThoughtWorksStudios/bobcat/dictionary"
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
)

type Field interface {
	Type() string
	GenerateValue() interface{}
	Amount() int
}

type FieldSet map[string]Field

type ReferenceField struct {
	referred  *Generator
	fieldName string
	minBound int
	maxBound int
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

func (field *ReferenceField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

type EntityField struct {
	entityGenerator *Generator
	minBound        int
	maxBound        int
}

func (field *EntityField) Type() string {
	return "entity"
}

func (field *EntityField) GenerateValue() interface{} {
	entities := make(map[string]GeneratedEntities)
	entities[field.entityGenerator.Type()] = field.entityGenerator.Generate(1)
	return entities
}

func (field *EntityField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

type UuidField struct{
	minBound int
	maxBound int
}

func (field *UuidField) Type() string {
	return "uuid"
}

func (field *UuidField) GenerateValue() interface{} {
	return uuid.NewV4()
}

func (field *UuidField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

type LiteralField struct {
	value interface{}
	minBound int
	maxBound int
}

func (field *LiteralField) Type() string {
	return "literal"
}

func (field *LiteralField) GenerateValue() interface{} {
	return field.value
}

func (field *LiteralField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

type StringField struct {
	length int
	minBound int
	maxBound int
}

func (field *StringField) Type() string {
	return "string"
}

func (field *StringField) GenerateValue() interface{} {
	allowedChars := []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!'@#$%^&*()_+-=[]{};:",./?`)
	result := []rune{}
	nTimes := rand.Intn(field.length-field.length+1) + field.length
	for i := 0; i < nTimes; i++ {
		result = append(result, allowedChars[rand.Intn(len(allowedChars))])
	}
	return string(result)
}

func (field *StringField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

type IntegerField struct {
	min int
	max int
	minBound int
	maxBound int
}

func (field *IntegerField) Type() string {
	return "integer"
}

func (field *IntegerField) GenerateValue() interface{} {
	result := float64(rand.Intn(int(field.max - field.min + 1)))
	result += float64(field.min)
	return int(result)
}

func (field *IntegerField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

type FloatField struct {
	min float64
	max float64
	minBound int
	maxBound int
}

func (field *FloatField) Type() string {
	return "float"
}

func (field *FloatField) GenerateValue() interface{} {
	return float64(rand.Intn(int(field.max-field.min))) + field.min + rand.Float64()
}

func (field *FloatField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

type DateField struct {
	min time.Time
	max time.Time
	minBound int
	maxBound int
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

func (field *DateField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

type DictField struct {
	category string
	minBound int
	maxBound int
}

var CustomDictPath = ""

func (field *DictField) Type() string {
	return "dict"
}

func (field *DictField) GenerateValue() interface{} {
	dictionary.SetCustomDataLocation(CustomDictPath)
	return dictionary.ValueFromDictionary(field.category)
}

func (field *DictField) Amount() int {
	return determineAmount(field.minBound, field.maxBound)
}

func determineAmount(min int, max int) int {
	if max == 0 && min == 0 {
		return 1
	} else if max - min == 0 {
		return min
	}

	return rand.Intn(max - min + 1) + min
}