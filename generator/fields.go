package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/dictionary"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/rs/xid"
	"math"
	"math/rand"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

type Field struct {
	fieldType      FieldType
	count          *CountRange
	UniqueValue    bool
	previousValues []interface{}
}

func (f *Field) MultiValue() bool {
	return f.count.Multiple()
}

func (f *Field) Type() string {
	return f.fieldType.Type()
}

func (f *Field) underlyingType() string {
	ft := f.fieldType

	for ft.Type() == "reference" {
		rt := ft.(*ReferenceType)
		ft = rt.referred.GetField(rt.fieldName).fieldType
	}

	return ft.Type()
}

func (f *Field) numberOfPossibilities() int64 {
	if int64(-1) == f.fieldType.numberOfPossibilities() {
		return int64(-1)
	}
	return f.fieldType.numberOfPossibilities() - int64(len(f.previousValues))
}

func (f *Field) Uniquable() bool {
	if f.count.Multiple() {
		return false
	}

	switch f.underlyingType() {
	case "dict", "enum", "string", "date", "integer", "float":
		return true
	default:
		return false
	}
}

func (f Field) String() string {
	return fmt.Sprintf(`{ type: %q, underlying: %q, multiVal: %t, unique: %t`, f.Type(), f.underlyingType(), f.MultiValue(), f.Uniquable() && f.UniqueValue)
}

func (f *Field) GenerateValue(parentId interface{}, emitter Emitter, scope *Scope) (result interface{}, err error) {
	result, err = f.value(parentId, emitter, scope)

	if err != nil {
		return nil, err
	}

	if f.Uniquable() && f.UniqueValue {
		tries := 1000
		for contains(f.previousValues, result) {
			if tries < 0 {
				return nil, fmt.Errorf("Failed to generate unique value for %q", f.underlyingType())
			}

			result, err = f.value(parentId, emitter, scope)

			if err != nil {
				return nil, err
			}

			tries--
		}
		f.previousValues = append(f.previousValues, result)
	}

	return result, nil
}

func (f *Field) value(parentId interface{}, emitter Emitter, scope *Scope) (result interface{}, err error) {
	if !f.count.Multiple() {
		result, err = f.fieldType.One(parentId, emitter, scope)
	} else {
		count := f.count.Count()
		values := make([]interface{}, count)

		for i := int64(0); i < count; i++ {
			if values[i], err = f.fieldType.One(parentId, emitter, scope); err != nil {
				break
			}
		}

		result = values
	}
	return result, err
}

func contains(sl []interface{}, value interface{}) bool {
	for _, a := range sl {
		if a == value {
			return true
		}
	}
	return false

}

type FieldType interface {
	Type() string
	One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error)
	// If numberOfPossibilities returns -1, then there are infinite possibilities
	numberOfPossibilities() int64
}

func NewField(fieldType FieldType, count *CountRange, unique bool) *Field {
	return &Field{fieldType: fieldType, count: count, UniqueValue: unique, previousValues: make([]interface{}, 0)}
}

type SerialType struct {
	current uint64
}

func (field *SerialType) Type() string {
	return "serial"
}

func (field *SerialType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	field.current++
	return field.current, nil
}

func (field *SerialType) numberOfPossibilities() int64 {
	return int64(-1)
}

type DeferredType struct {
	closure DeferredResolver
}

func (f *DeferredType) Type() string {
	return "deferred"
}

func (f *DeferredType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return f.closure(scope)
}

func (f *DeferredType) numberOfPossibilities() int64 {
	return int64(-1)
}

type ReferenceType struct {
	referred  *FieldSet
	fieldName string
}

func (f *ReferenceType) Type() string {
	return "reference"
}

func (f *ReferenceType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	ref := f.referred.GetField(f.fieldName).fieldType
	return ref.One(parentId, emitter, scope)
}

func (f *ReferenceType) numberOfPossibilities() int64 {
	ref := f.referred.GetField(f.fieldName).fieldType
	return ref.numberOfPossibilities()
}

type EntityType struct {
	entityGenerator *Generator
}

func (field *EntityType) Type() string {
	return "entity"
}

func (field *EntityType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	g := field.entityGenerator
	if result, err := g.One(parentId, emitter, scope); err == nil {
		return result[g.PrimaryKeyName()], nil
	} else {
		return nil, err
	}
}

func (field *EntityType) numberOfPossibilities() int64 {
	return int64(-1)
}

type BoolType struct {
}

func (field *BoolType) Type() string {
	return "boolean"
}

func (field *BoolType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return 49 < rand.Intn(100), nil
}

func (field *BoolType) numberOfPossibilities() int64 {
	return 2
}

type MongoIDType struct {
}

func (field *MongoIDType) Type() string {
	return "uid"
}

func (field *MongoIDType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return xid.New().String(), nil
}

func (field *MongoIDType) numberOfPossibilities() int64 {
	return 0
}

type LiteralType struct {
	value interface{}
}

func (field *LiteralType) Type() string {
	return "literal"
}

func (field *LiteralType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return field.value, nil
}

func (field *LiteralType) numberOfPossibilities() int64 {
	return 1
}

type StringType struct {
	length int64
}

func (field *StringType) Type() string {
	return "string"
}

const ALLOWED_CHARACTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@"

var LETTER_INDEX_BITS uint = uint(math.Ceil(math.Log2(float64(len(ALLOWED_CHARACTERS))))) // number of bits to represent ALLOWED_CHARACTERS
var LETTER_BIT_MASK int64 = 1<<LETTER_INDEX_BITS - 1                                      // All 1-bits, as many as LETTER_INDEX_BITS
var LETTERS_PER_INT63 uint = 63 / LETTER_INDEX_BITS                                       // # of letter indices fitting in 63 bits as generated by src.Int63

func (field *StringType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	n := field.length
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), LETTERS_PER_INT63; i >= int64(0); {
		if remain == 0 {
			cache, remain = src.Int63(), LETTERS_PER_INT63
		}
		if idx := int(cache & LETTER_BIT_MASK); idx < len(ALLOWED_CHARACTERS) {
			b[i] = ALLOWED_CHARACTERS[idx]
			i--
		}
		cache >>= LETTER_INDEX_BITS
		remain--
	}

	return string(b), nil
}

func (field *StringType) numberOfPossibilities() int64 {
	if field.length > 10 {
		return -1
	}
	return int64(math.Pow(float64(len(ALLOWED_CHARACTERS)), float64(field.length)))
}

type IntegerType struct {
	min int64
	max int64
}

func (field *IntegerType) Type() string {
	return "integer"
}

func (field *IntegerType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return field.min + rand.Int63n(field.max-field.min+1), nil
}

func (field *IntegerType) numberOfPossibilities() int64 {
	return field.max - field.min + 1
}

type FloatType struct {
	min float64
	max float64
}

func (field *FloatType) Type() string {
	return "float"
}

func (field *FloatType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return rand.Float64()*(field.max-field.min) + field.min, nil
}

func (field *FloatType) numberOfPossibilities() int64 {
	if field.min == field.max {
		return int64(1)
	} else {
		return int64(-1)
	}
}

type DateType struct {
	min    time.Time
	max    time.Time
	format string
}

func (field *DateType) Type() string {
	return "date"
}

func (field *DateType) ValidBounds() bool {
	return field.min.Before(field.max)
}

func (field *DateType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	min, max := field.min.Unix(), field.max.Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min

	return &TimeWithFormat{Time: time.Unix(sec, 0), Format: field.format}, nil
}

func (field *DateType) numberOfPossibilities() int64 {
	//Number of seconds between the min and max
	return int64(field.max.Sub(field.min).Seconds())
}

type DictType struct {
	category string
}

var CustomDictPath = ""

func (field *DictType) Type() string {
	return "dict"
}

func (field *DictType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	dictionary.SetCustomDataLocation(CustomDictPath)
	return dictionary.ValueFromDictionary(field.category), nil
}

func (field *DictType) numberOfPossibilities() int64 {
	dictionary.SetCustomDataLocation(CustomDictPath)
	return dictionary.CalculatePossibilities(field.category)
}

type EnumType struct {
	size   int64
	values []interface{}
}

func (field *EnumType) Type() string {
	return "enum"
}

func (field *EnumType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return field.values[rand.Int63n(field.size)], nil
}

func (field *EnumType) numberOfPossibilities() int64 {
	return field.size
}

type DistributionType struct {
	bins []*Field
	dist Distribution
}

func (field *DistributionType) Type() string {
	return "distribution"
}

func (field *DistributionType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return field.dist.One(field.domain(), parentId, emitter, scope)
}

func (field *DistributionType) domain() Domain {
	intervals := make([]FieldType, len(field.bins))
	for i := 0; i < len(field.bins); i++ {
		intervals[i] = field.bins[i].fieldType
	}
	return Domain{intervals: intervals}

}

func (field *DistributionType) numberOfPossibilities() int64 {
	return 1
}
