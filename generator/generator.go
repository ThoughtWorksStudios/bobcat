package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/logging"
	"strings"
	"time"
)

type Generator struct {
	name   string
	base   string
	fields FieldSet
	log    logging.ILogger
}

func ExtendGenerator(name string, parent *Generator) *Generator {
	gen := NewGenerator(name, parent.log)
	gen.base = parent.Type()
	gen.fields["$extends"] = NewField(&LiteralType{value: gen.base}, nil)
	gen.fields["$type"] = NewField(&LiteralType{value: gen.Type()}, nil)

	for key, f := range parent.fields {
		if _, hasField := gen.fields[key]; !hasField || !strings.HasPrefix(key, "$") {
			gen.fields[key] = NewField(&ReferenceType{referred: parent, fieldName: key}, f.count)
		}
	}

	return gen
}

func NewGenerator(name string, logger logging.ILogger) *Generator {
	if logger == nil {
		logger = &logging.DefaultLogger{}
	}

	if name == "" {
		name = "$"
	}

	g := &Generator{name: name, fields: make(FieldSet), log: logger}

	g.fields["$id"] = NewField(&MongoIDType{}, nil)

	g.fields["$type"] = NewField(&LiteralType{value: g.name}, nil)
	g.fields["$species"] = NewField(&LiteralType{value: g.name}, nil)

	return g
}

func (g *Generator) WithStaticField(fieldName string, fieldValue interface{}) error {
	g.fields[fieldName] = NewField(&LiteralType{value: fieldValue}, nil)
	return nil
}

func (g *Generator) WithEntityField(fieldName string, entityGenerator *Generator, fieldArgs interface{}, countRange *CountRange) error {
	g.fields[fieldName] = NewField(&EntityType{entityGenerator: entityGenerator}, countRange)
	return nil
}

func (g *Generator) WithField(fieldName, fieldType string, fieldArgs interface{}, countRange *CountRange) error {
	switch fieldType {
	case "string":
		if ln, ok := fieldArgs.(int); ok {
			g.fields[fieldName] = NewField(&StringType{length: ln}, countRange)
		} else {
			return fmt.Errorf("expected field args to be of type 'int' for field %s (%s), but got %v",
				fieldName, fieldType, fieldArgs)
		}
	case "integer":
		if bounds, ok := fieldArgs.([2]int); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return fmt.Errorf("max %v cannot be less than min %v", max, min)
			}

			g.fields[fieldName] = NewField(&IntegerType{min: min, max: max}, countRange)
		} else {
			return fmt.Errorf("expected field args to be of type '(min:int, max:int)' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case "decimal":
		if bounds, ok := fieldArgs.([2]float64); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return fmt.Errorf("max %v cannot be less than min %v", max, min)
			}
			g.fields[fieldName] = NewField(&FloatType{min: min, max: max}, countRange)
		} else {
			return fmt.Errorf("expected field args to be of type '(min:float64, max:float64)' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case "date":
		if bounds, ok := fieldArgs.([2]time.Time); ok {
			min, max := bounds[0], bounds[1]
			dateType := &DateType{min: min, max: max}
			if !dateType.ValidBounds() {
				return fmt.Errorf("max %v cannot be before min %v", max, min)
			}
			g.fields[fieldName] = NewField(dateType, countRange)
		} else {
			return fmt.Errorf("expected field args to be of type 'time.Time' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case "mongoid":
		g.fields[fieldName] = NewField(&MongoIDType{}, nil)
	case "bool":
		g.fields[fieldName] = NewField(&BoolType{}, countRange)
	case "dict":
		if dict, ok := fieldArgs.(string); ok {
			g.fields[fieldName] = NewField(&DictType{category: dict}, countRange)
		} else {
			return fmt.Errorf("expected field args to be of type 'string' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case "enum":
		if values, ok := fieldArgs.([]interface{}); ok {
			g.fields[fieldName] = NewField(&EnumType{values: values}, countRange)
		} else {
			return fmt.Errorf("expected field args to be a list of values for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	default:
		return fmt.Errorf("Invalid field type '%v'", fieldType)
	}

	return nil
}

func (g *Generator) Type() string {
	if (strings.HasPrefix(g.name, "$") || g.name == "") && g.base != "" {
		return g.base
	}
	return g.name
}

func (g *Generator) Generate(count int64) GeneratedEntities {
	entities := NewGeneratedEntities(count)
	for i := int64(0); i < count; i++ {
		entities[i] = g.One("")
	}
	return entities
}

func (g *Generator) One(parentId string) EntityResult {
	entity := EntityResult{}
	id, _ := g.fields["$id"].GenerateValue("").(string)
	entity["$id"] = id // create this first because we may use it as reference in $parent
	if parentId != "" {
		entity["$parent"] = parentId
	}

	for name, field := range g.fields {
		if _, hasVal := entity[name]; !hasVal {
			entity[name] = field.GenerateValue(id)
		}
	}
	return entity
}

func (g *Generator) String() string {
	return fmt.Sprintf("%s{}", g.name)
}
