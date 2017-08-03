package generator

import (
	"github.com/ThoughtWorksStudios/datagen/dictionary"
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
)

type Field interface {
	Type() string
	GenerateValue() interface{}
	Count() int
}

type Range struct {
	min int
	max int
}

func NewRange(min, max int) Range {
	return Range{ min: min, max: max}
}

type FieldSet map[string]Field

type ReferenceField struct {
	referred  *Generator
	fieldName string
	count Range
}

func (field *ReferenceField) Type() string {
	return "reference"
}

func (field *ReferenceField) GenerateValue() interface{} {
	referredField := field.referred.fields[field.fieldName]
	fieldCount := referredField.Count()
	if fieldCount > 1 {
		fieldValue := []interface{}{}
		for i := 0; i < fieldCount; i++ {
			fieldValue = append(fieldValue, referredField.GenerateValue())
		}
		return fieldValue
 	}
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

func (field *ReferenceField) Count() int {
	return generateCount(field.count)
}

func generateCount(count Range, seed... int64) int {
	if len(seed) > 0 {
		rand.Seed(seed[0])
	}

	min, max := count.min, count.max
	if min == 0 && max == 0 {
		return 1
	}
	return rand.Intn(max - min + 1) + min
}

type EntityField struct {
	entityGenerator *Generator
	count Range
}

func (field *EntityField) Type() string {
	return "entity"
}

func (field *EntityField) GenerateValue() interface{} {
	return field.entityGenerator.Generate(1)
}

func (field *EntityField) Count() int {
	return generateCount(field.count)
}

type UuidField struct{
	count Range
}

func (field *UuidField) Type() string {
	return "uuid"
}

func (field *UuidField) GenerateValue() interface{} {
	return uuid.NewV4()
}

func (field *UuidField) Count() int {
	return generateCount(field.count)
}


type LiteralField struct {
	value interface{}
	count Range
}

func (field *LiteralField) Type() string {
	return "literal"
}

func (field *LiteralField) GenerateValue() interface{} {
	return field.value
}


func (field *LiteralField) Count() int {
	return generateCount(field.count)
}

type StringField struct {
	length int
	count Range
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

func (field *StringField) Count() int {
	return generateCount(field.count)
}

type IntegerField struct {
	min int
	max int
	count Range
}

func (field *IntegerField) Type() string {
	return "integer"
}

func (field *IntegerField) GenerateValue() interface{} {
	result := float64(rand.Intn(int(field.max - field.min + 1)))
	result += float64(field.min)
	return int(result)
}

func (field *IntegerField) Count() int {
	return generateCount(field.count)
}

type FloatField struct {
	min float64
	max float64
	count Range
}

func (field *FloatField) Type() string {
	return "float"
}

func (field *FloatField) GenerateValue() interface{} {
	return float64(rand.Intn(int(field.max-field.min))) + field.min + rand.Float64()
}

func (field *FloatField) Count() int {
	return generateCount(field.count)
}

type DateField struct {
	min time.Time
	max time.Time
	count Range
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

func (field *DateField) Count() int {
	return generateCount(field.count)
}

type DictField struct {
	category string
	count Range
}

var CustomDictPath = ""

func (field *DictField) Type() string {
	return "dict"
}

func (field *DictField) GenerateValue() interface{} {
	dictionary.SetCustomDataLocation(CustomDictPath)
	return dictionary.ValueFromDictionary(field.category)
}

func (field *DictField) Count() int {
	return generateCount(field.count)
}
