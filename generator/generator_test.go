package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/builtins"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"strings"
	"testing"
	"time"
)

func collection(vals ...interface{}) []interface{} {
	return vals
}

func isBetween(actual, lower, upper float64) bool {
	return actual >= lower && actual <= upper
}

func AssertEquivField(t *testing.T, expected, actual *Field) {
	// is this a cheap hack? you bet it is.
	if expected.String() != actual.String() {
		t.Errorf("Expected: \n [%v] \n\n but got: \n [%v]", expected.fieldType, actual.fieldType)
	}
}

func TestExtendGenerator(t *testing.T) {
	g := NewGenerator("thing", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	addBuiltin(g, "age", FLOAT_TYPE, float64(2), float64(4))
	addLiteral(g, "species", "human")

	m := ExtendGenerator("thang", g, nil, false)
	addLiteral(m, "species", "h00man")
	addLiteral(m, "name", "kyle")

	emitter := NewTestEmitter()
	scope := NewRootScope()

	g.Generate(1, emitter, scope)

	base := emitter.Shift()
	AssertNotNil(t, base, "Should have generated an entity result")

	AssertEqual(t, "human", base["species"])
	AssertEqual(t, 10, len(base["name"].(string)))
	Assert(t, isBetween(base["age"].(float64), 2, 4), "base entity failed to generate the correct age")

	m.Generate(1, emitter, scope)

	extended := emitter.Shift()
	AssertNotNil(t, extended, "Should have generated an entity result")
	AssertEqual(t, "h00man", extended["species"])
	AssertEqual(t, "kyle", extended["name"].(string))
	Assert(t, isBetween(extended["age"].(float64), 2, 4), "extended entity failed to generate the correct age")
}

func TestNoMetadataGeneratedWhenDisabled(t *testing.T) {
	g := NewGenerator("Cat", nil, true)
	addBuiltin(g, "name", STRING_TYPE, int64(5))
	scope := NewRootScope()
	emitter := NewTestEmitter()

	g.One("foo", emitter, scope)
	entity := emitter.Shift()

	for name, _ := range entity {
		if strings.HasPrefix(name, "$") && name != g.PrimaryKeyName() && name != "$parent" {
			t.Errorf("Found metadata in entity when there should be none, '%v'", name)
		}
	}
}

func TestSubentityHasParentReference(t *testing.T) {
	subentityGenerator := NewGenerator("Cat", nil, false)
	addBuiltin(subentityGenerator, "name", STRING_TYPE, int64(5))

	g := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	g.WithEntityField("pet", subentityGenerator, nil)
	scope := NewRootScope()
	emitter := NewTestEmitter()

	g.Generate(1, emitter, scope)
	cat := emitter.Shift()
	person := emitter.Shift()

	if person[g.PrimaryKeyName()] != cat["$parent"] {
		t.Errorf("Parent id (%v) on subentity does not match the parent entity's id (%v)", cat["$parent"], person[g.PrimaryKeyName()])
	}

	subentityGenerator.Generate(1, emitter, scope)
	nextCat := emitter.Shift()

	if val, ok := nextCat["$parent"]; ok {
		t.Errorf("Cat should not have a parent (%v) when generated on it's own", val)
	}
}

func TestWithStaticFieldCreatesCorrectField(t *testing.T) {
	g := NewGenerator("thing", nil, false)
	g.WithLiteralField("login", "something")
	expectedField := NewField(&LiteralType{"something"}, nil)
	AssertEquivField(t, expectedField, g.fields.GetField("login"))
}

func TestWithEntityFieldCreatesCorrectField(t *testing.T) {
	g := NewGenerator("thing", nil, false)
	countRange := &CountRange{3, 3}
	g.WithEntityField("food", g, countRange)
	expectedField := NewField(&EntityType{g}, countRange)
	AssertEquivField(t, expectedField, g.fields.GetField("food"))
}

func TestGenerateProducesGeneratedContent(t *testing.T) {
	g := NewGenerator("thing", nil, false)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	addBuiltin(g, "a", STRING_TYPE, int64(2))
	addBuiltin(g, "b", INT_TYPE, int64(2), int64(4))
	addBuiltin(g, "c", FLOAT_TYPE, float64(2.85), float64(4.50))
	addBuiltin(g, "d", DATE_TYPE, timeMin, timeMax, "")
	addBuiltin(g, "e", DICT_TYPE, "last_name")
	addBuiltin(g, "f", UID_TYPE)
	addBuiltin(g, "g", ENUM_TYPE, collection("a", "b"))
	g.WithEntityField("h", NewGenerator("thang", nil, false), nil)
	addBuiltin(g, "i", SERIAL_TYPE)
	scope := NewRootScope()
	emitter := NewTestEmitter()

	data, err := g.Generate(3, emitter, scope)
	AssertNil(t, err, "Should not receive error")

	emitter.Shift()
	entity := emitter.Shift()

	AssertEqual(t, 3, len(data))

	testFields := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}

	for _, fieldName := range testFields {
		fieldValue, ok := entity[fieldName]
		Assert(t, ok, "entity should have field %q", fieldName)

		switch fieldType := fieldValue.(type) {
		case int64:
			Assert(t, fieldName == "b", "field %q should have yielded a int64", fieldName)
		case float64:
			Assert(t, fieldName == "c", "field %q should have yielded a float64", fieldName)
		case string:
			Assert(t, strings.Contains("a, e, f, g, h", fieldName), "field %q should have yielded a string", fieldName)
		case *TimeWithFormat:
			Assert(t, fieldName == "d", "field %q should have yielded a Time", fieldName)
		case uint64:
			Assert(t, fieldName == "i", "field %q should have yielded a int", fieldName)
		default:
			Assert(t, false, "Don't know what to do with the field type for %q! The type is %v", fieldName, fieldType)
		}
	}
}

func TestGenerateWithBoundsArgumentProducesCorrectCountOfValues(t *testing.T) {
	g := NewGenerator("thing", nil, false)
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	g.WithEntityField("a", NewGenerator("subthing", nil, false), &CountRange{1, 1})
	addBuiltinWithCount(g, "b", STRING_TYPE, &CountRange{2, 2}, int64(2))
	addBuiltinWithCount(g, "c", INT_TYPE, &CountRange{3, 3}, int64(2), int64(4))
	addBuiltinWithCount(g, "d", FLOAT_TYPE, &CountRange{4, 4}, float64(2.85), float64(4.50))
	addBuiltinWithCount(g, "e", DATE_TYPE, &CountRange{5, 5}, timeMin, timeMax, "")
	addBuiltinWithCount(g, "f", DICT_TYPE, &CountRange{6, 6}, "last_name")
	addBuiltinWithCount(g, "g", ENUM_TYPE, &CountRange{7, 7}, collection("a", "b"))
	scope := NewRootScope()
	emitter := NewTestEmitter()
	g.Generate(1, emitter, scope)
	emitter.Shift()
	entity := emitter.Shift()

	var testFields = []struct {
		fieldName string
		count     int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
		{"e", 5},
		{"f", 6},
		{"g", 7},
	}

	for _, field := range testFields {
		actual := len(entity[field.fieldName].([]interface{}))
		AssertEqual(t, field.count, actual)
	}
}

func TestHasField(t *testing.T) {
	g := NewGenerator("thing", nil, false)
	addBuiltin(g, "eek", FLOAT_TYPE, float64(2.0), float64(4.0), nil, true)
	Assert(t, g.HasField("eek"), "Expected field 'eek' to exist, but it does not!")
}

func TestGeneratedFieldsUsesExistingFieldValuesWhenAvailable(t *testing.T) {
	g := NewGenerator("generator", nil, false)
	addBuiltin(g, "price", FLOAT_TYPE, float64(2.0), float64(4.0))
	closure := func(scope *Scope) (interface{}, error) { return scope.ResolveSymbol("price"), nil }
	g.WithDeferredField("price_clone", closure)
	scope := NewRootScope()

	result, err := g.One(nil, NewTestEmitter(), scope)
	AssertNil(t, err, "Should not receive error")

	AssertEqual(t, result["price"], result["price_clone"],
		"Expected 'price' and 'price_clone' fields to match, but got: '%v', '%v'",
		result["price"], result["price_clone"])
}

func TestGeneratedFieldsDoesNotUseExistingFieldValuesWhenNotAvailable(t *testing.T) {
	g := NewGenerator("generator", nil, false)
	closure := func(scope *Scope) (interface{}, error) { return scope.ResolveSymbol("foo"), nil }
	g.WithDeferredField("price_clone", closure)

	result, err := g.One(nil, NewTestEmitter(), NewRootScope())
	AssertNil(t, err, "Should not receive error")

	AssertEqual(t, result["price_clone"], nil,
		"Expected 'price_clone' to not exist, but got: '%v'", result["price_clone"])
}

func TestPercentDistributionValidatesPercentages(t *testing.T) {
	_, err := NewDistribution(PERCENT_DIST, []float64{0.10, 0.20, 0.10}, []FieldType{})
	ExpectsError(t, "percentage weights do not add to 100%", err)

	_, err = NewDistribution(PERCENT_DIST, []float64{0.10, 0.20, 0.10, 0.60}, []FieldType{})
	AssertNil(t, err, "Should not receive error when percentages add to 100%")
}

func TestWithDistributionFieldCanParseMultipleFieldTypes(t *testing.T) {
	g := NewGenerator("thing", nil, false)
	fields := []FieldType{
		NewDeferredType(func(scope *Scope) (interface{}, error) {
			return 3.14159, nil
		}),

		NewDeferredType(func(scope *Scope) (interface{}, error) {
			return int64(22), nil
		}),

		NewLiteralType("just a string"),
	}

	weights := []float64{0.10, 0.80, 0.10}

	dist, err := NewDistribution(PERCENT_DIST, weights, fields)
	AssertNil(t, err, "Should not receive error while constructing distribution")

	g.WithField("eek", dist, nil)

	_, err = g.One(nil, NewDummyEmitter(), nil)
	AssertNil(t, err, "Should not generate error")
}
