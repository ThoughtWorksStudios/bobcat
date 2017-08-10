package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/dictionary"
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
)

type Field struct {
	fieldType FieldType
	count     *CountRange
}

func (f *Field) Type() string {
	return f.fieldType.Type()
}

func (f *Field) GenerateValue() interface{} {
	if !f.count.Multiple() {
		return f.fieldType.GenerateSingle()
	} else {
		count := f.count.Count()
		values := make([]interface{}, count)

		for i := 0; i < count; i++ {
			values[i] = f.fieldType.GenerateSingle()
		}

		return values
	}
}

type FieldSet map[string]*Field

type FieldType interface {
	Type() string
	GenerateSingle() interface{}
}

func NewField(fieldType FieldType, count *CountRange) *Field {
	return &Field{fieldType: fieldType, count: count}
}

type ReferenceType struct {
	referred  *Generator
	fieldName string
}

func (field *ReferenceType) Type() string {
	return "reference"
}

func (field *ReferenceType) GenerateSingle() interface{} {
	ref := field.referred.fields[field.fieldName].fieldType
	return ref.GenerateSingle()
}

type EntityType struct {
	entityGenerator *Generator
}

func (field *EntityType) Type() string {
	return "entity"
}

func (field *EntityType) GenerateSingle() interface{} {
	return field.entityGenerator.Generate(1)[0]
}

type BoolType struct {
}

func (field *BoolType) Type() string {
	return "boolean"
}

func (field *BoolType) GenerateSingle() interface{} {
	return 49 < rand.Intn(100)
}

type UuidType struct {
}

func (field *UuidType) Type() string {
	return "uuid"
}

func (field *UuidType) GenerateSingle() interface{} {
	return uuid.NewV4()
}

type LiteralType struct {
	value interface{}
}

func (field *LiteralType) Type() string {
	return "literal"
}

func (field *LiteralType) GenerateSingle() interface{} {
	return field.value
}

type StringType struct {
	length int
}

func (field *StringType) Type() string {
	return "string"
}

func (field *StringType) GenerateSingle() interface{} {
	allowedChars := []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!'@#$%^&*()_+-=[]{};:",./?`)
	result := []rune{}
	nTimes := rand.Intn(field.length-field.length+1) + field.length
	for i := 0; i < nTimes; i++ {
		result = append(result, allowedChars[rand.Intn(len(allowedChars))])
	}
	return string(result)
}

type IntegerType struct {
	min int
	max int
}

func (field *IntegerType) Type() string {
	return "integer"
}

func (field *IntegerType) GenerateSingle() interface{} {
	result := float64(rand.Intn(int(field.max - field.min + 1)))
	result += float64(field.min)
	return int(result)
}

type FloatType struct {
	min float64
	max float64
}

func (field *FloatType) Type() string {
	return "float"
}

func (field *FloatType) GenerateSingle() interface{} {
	return rand.Float64()*(field.max-field.min) + field.min
}

type DateType struct {
	min time.Time
	max time.Time
}

func (field *DateType) Type() string {
	return "date"
}

func (field *DateType) ValidBounds() bool {
	return field.min.Before(field.max)
}

func (field *DateType) GenerateSingle() interface{} {
	min, max := field.min.Unix(), field.max.Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min

	return time.Unix(sec, 0)
}

type DictType struct {
	category string
}

var CustomDictPath = ""

func (field *DictType) Type() string {
	return "dict"
}

func (field *DictType) GenerateSingle() interface{} {
	dictionary.SetCustomDataLocation(CustomDictPath)
	return dictionary.ValueFromDictionary(field.category)
}
