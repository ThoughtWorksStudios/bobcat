package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"strings"
	"time"
)

var DEFAULT_PK_CONFIG *PrimaryKey = &PrimaryKey{name: "$id", kind: UID_TYPE}

type Generator struct {
	name            string
	extends         string
	declaredType    string
	fields          *FieldSet
	disableMetadata bool
	pkey            *PrimaryKey
}

func ExtendGenerator(name string, parent *Generator, pkey *PrimaryKey, disableMetadata bool) *Generator {
	var gen *Generator

	if pkey == nil {
		gen = NewGenerator(name, parent.pkey, disableMetadata)
		parent.pkey.Inherit(gen, parent) // must reference parent's primary key; this is important for serial fields to continue the sequence
	} else {
		gen = NewGenerator(name, pkey, disableMetadata)
	}

	gen.extends = parent.Type()

	gen.recalculateType()

	if !disableMetadata {
		gen.fields.AddField("$extends", NewField(&LiteralType{value: gen.extends}, nil, false))
		gen.fields.AddField("$type", NewField(&LiteralType{value: gen.Type()}, nil, false))
	}

	for _, fieldEntry := range parent.fields.AllFields() {
		key := fieldEntry.Name
		f := fieldEntry.Field
		if !gen.fields.HasField(key) && !strings.HasPrefix(key, "$") && key != parent.PrimaryKeyName() {
			gen.fields.AddField(key, NewField(&ReferenceType{referred: parent.fields, fieldName: key}, f.count, false))
		}
	}

	return gen
}

func NewGenerator(name string, pkey *PrimaryKey, disableMetadata bool) *Generator {
	if name == "" {
		name = "$"
	}

	g := &Generator{name: name, fields: NewFieldSet(), disableMetadata: disableMetadata}

	g.recalculateType()

	if pkey == nil {
		pkey = DEFAULT_PK_CONFIG
	}

	pkey.Attach(g)

	if !disableMetadata {
		g.fields.AddField("$type", NewField(&LiteralType{value: g.Type()}, nil, false))
	}

	return g
}

func (g *Generator) HasField(name string) bool {
	return g.fields.HasField(name)
}

func (g *Generator) GetField(name string) *Field {
	return g.fields.GetField(name)
}

func (g *Generator) PrimaryKeyName() string {
	return g.pkey.name
}

func (g *Generator) WithDeferredField(fieldName string, fieldValue DeferredResolver) error {
	g.fields.AddField(fieldName, NewDeferredField(fieldValue))
	return nil
}

func (g *Generator) WithLiteralField(fieldName string, fieldValue interface{}) error {
	g.fields.AddField(fieldName, NewLiteralField(fieldValue))
	return nil
}

func (g *Generator) WithEntityField(fieldName string, entityGenerator *Generator, countRange *CountRange) error {
	g.fields.AddField(fieldName, NewField(&EntityType{entityGenerator: entityGenerator}, countRange, false))
	return nil
}

func (g *Generator) newFieldType(fieldName, fieldType string, fieldArgs interface{}, countRange *CountRange, uniqueValue bool) (*Field, error) {
	switch fieldType {
	case STRING_TYPE:
		if ln, ok := fieldArgs.(int64); ok {
			return NewField(&StringType{length: ln}, countRange, uniqueValue), nil
		} else {
			return nil, fmt.Errorf("expected field args to be of type 'int' for field %s (%s), but got %v",
				fieldName, fieldType, fieldArgs)
		}
	case INT_TYPE:
		if bounds, ok := fieldArgs.([2]int64); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return nil, fmt.Errorf("max %v cannot be less than min %v", max, min)
			}

			return NewField(&IntegerType{min: min, max: max}, countRange, uniqueValue), nil
		} else {
			return nil, fmt.Errorf("expected field args to be of type '(min:int, max:int)' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case FLOAT_TYPE:
		if bounds, ok := fieldArgs.([2]float64); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return nil, fmt.Errorf("max %v cannot be less than min %v", max, min)
			}
			return NewField(&FloatType{min: min, max: max}, countRange, uniqueValue), nil
		} else {
			return nil, fmt.Errorf("expected field args to be of type '(min:float64, max:float64)' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case DATE_TYPE:
		if bounds, ok := fieldArgs.([]interface{}); ok {
			min, max, format := bounds[0].(time.Time), bounds[1].(time.Time), bounds[2].(string)
			dateType := &DateType{min: min, max: max, format: format}
			if !dateType.ValidBounds() {
				return nil, fmt.Errorf("max %v cannot be before min %v", max, min)
			}
			return NewField(dateType, countRange, uniqueValue), nil
		} else {
			return nil, fmt.Errorf("expected field args to be of type 'time.Time' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case UID_TYPE:
		return NewField(&MongoIDType{}, nil, false), nil
	case BOOL_TYPE:
		if uniqueValue {
			return nil, fmt.Errorf("boolean fields cannot be unique")
		}
		return NewField(&BoolType{}, countRange, false), nil
	case DICT_TYPE:
		if dict, ok := fieldArgs.(string); ok {
			return NewField(&DictType{category: dict}, countRange, uniqueValue), nil
		} else {
			return nil, fmt.Errorf("expected field args to be of type 'string' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case ENUM_TYPE:
		if values, ok := fieldArgs.([]interface{}); ok {
			return NewField(&EnumType{values: values, size: int64(len(values))}, countRange, uniqueValue), nil
		} else {
			return nil, fmt.Errorf("expected field args to be of type 'collection' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case SERIAL_TYPE:
		if countRange != nil {
			return nil, fmt.Errorf("serial fields can only have a single value")
		}
		return NewField(&SerialType{}, nil, false), nil
	default:
		return nil, fmt.Errorf("Invalid field type '%v'", fieldType)
	}

	return nil, nil
}

func (g *Generator) WithField(fieldName, fieldType string, fieldArgs interface{}, countRange *CountRange, uniqueValue bool) error {
	if field, err := g.newFieldType(fieldName, fieldType, fieldArgs, countRange, uniqueValue); err == nil {
		g.fields.AddField(fieldName, field)
	} else {
		return err
	}
	return nil
}

func (g *Generator) WithDistribution(fieldName, distType string, fieldTypes []FieldType, weights []float64) error {
	distribution, err := g.newDistribution(distType, weights)

	if err != nil {
		return err
	}

	if !distribution.supportsMultipleIntervals() && len(fieldTypes) > 1 {
		return fmt.Errorf("%v distributions do not support multiple domains", distribution.Type())
	}

	for _, field := range fieldTypes {
		if "literal" == field.Type() {
			v, _ := field.One(nil, nil, nil)
			var valueType string = "anything"
			switch v.(type) {
			case int64:
				valueType = INT_TYPE
			case float64:
				valueType = FLOAT_TYPE
			}

			if !distribution.isCompatibleDomain(valueType) {
				return fmt.Errorf("Invalid Distribution Domain: %v is not a valid domain for %v distributions", valueType, distribution.Type())
			}
		}
	}

	g.fields.AddField(fieldName, NewField(&DistributionType{domain: Domain{intervals: fieldTypes}, dist: distribution}, nil, false))

	return nil
}

func (g *Generator) newDistribution(distType string, weights []float64) (Distribution, error) {
	switch distType {
	case NORMAL_DIST:
		return &NormalDistribution{}, nil
	case WEIGHT_DIST:
		for _, w := range weights {
			if w < 0 {
				return nil, fmt.Errorf("weights cannot be negative: %f", w)
			}
		}
		return &WeightDistribution{weights: weights}, nil
	case PERCENT_DIST:
		total := float64(0)

		for _, w := range weights {
			if w < 0 {
				return nil, fmt.Errorf("weights cannot be negative: %f", w)
			}
			total += w
		}

		if total != float64(1) {
			return nil, fmt.Errorf("percentage weights do not add to 100%% (i.e. 1.0). total = %f", total)
		}

		return &WeightDistribution{weights: weights}, nil
	default:
		return nil, fmt.Errorf("Unsupported distribution %q", distType)
	}
}

func (g *Generator) Type() string {
	return g.declaredType
}

func (g *Generator) recalculateType() {
	if (strings.HasPrefix(g.name, "$") || g.name == "") && g.extends != "" {
		g.declaredType = g.extends
	} else {
		g.declaredType = g.name
	}
}

func (g *Generator) EnsureGeneratable(count int64) error {
	for _, fieldEntry := range g.fields.AllFields() {
		field := fieldEntry.Field
		name := fieldEntry.Name
		if field.Uniquable() && field.UniqueValue {
			numberOfPossibilities := field.numberOfPossibilities()
			if numberOfPossibilities != int64(-1) && numberOfPossibilities < count {
				return fmt.Errorf("Not enough unique values for field '%v': There are only %v unique values available for the '%v' field, and you're trying to generate %v entities", name, numberOfPossibilities, name, count)
			}
		}
	}
	return nil
}

func (g *Generator) Generate(count int64, emitter Emitter, scope *Scope) ([]interface{}, error) {
	ids := make([]interface{}, count)
	idKey := g.PrimaryKeyName()

	for i := int64(0); i < count; i++ {
		if r, err := g.One(nil, emitter, scope); err == nil {
			ids[i] = r[idKey]
		} else {
			return nil, err
		}
	}
	return ids, nil
}

func (g *Generator) One(parentId interface{}, emitter Emitter, scope *Scope) (EntityResult, error) {
	entity := EntityResult{}
	// Need to extend TransientScope once more as a protective layer so that any symbols declared
	// by expressions are NOT set as fields in the final entity result
	childScope := ExtendScope(TransientScope(scope, SymbolTable(entity)))

	idKey := g.PrimaryKeyName()
	id, err := g.fields.GetField(idKey).GenerateValue(nil, emitter, childScope)

	if err != nil {
		return nil, err
	}

	entity[idKey] = id // create this first because we may use it as the parentId when generating child entities

	if parentId != nil {
		entity["$parent"] = parentId
	}

	for _, entry := range g.fields.AllFields() {
		name, field := entry.Name, entry.Field
		if name != idKey { // don't GenerateValue() more than once for id -- it throws off the sequence for serial types
			// GenerateValue MAY populate entity[name] for entity fields
			value, err := field.GenerateValue(id, emitter.NextEmitter(entity, name, field.MultiValue()), childScope)
			if err != nil {
				return nil, err
			}
			if _, isAlreadySet := entity[name]; !isAlreadySet {
				entity[name] = value
			}
		}
	}

	if err = emitter.Emit(entity, g.Type()); err != nil {
		return nil, err
	}

	return entity, nil
}

func (g *Generator) String() string {
	return fmt.Sprintf("%s{}", g.name)
}
