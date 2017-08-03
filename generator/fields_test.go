package generator

import (
	"testing"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
)

func TestGenerateCountForEmptyFieldCountArray(t *testing.T) {
	actual := generateCount(Range{})
	AssertEqual(t,1, actual)
}

func TestGenerateCountWhenRangeHasNoDifference(t *testing.T) {
	actual := generateCount(Range{2,2})
	AssertEqual(t,2, actual)
}

func TestGenerateCountWhenRangeDifferenceIsOne(t *testing.T) {
	expecteds := []struct {
		seed int64
		count int
	}{
		{ int64(5), 3},
		{ int64(1), 4},
	}
	for _, expected := range expecteds {
		actual := generateCount(Range{3,4}, expected.seed)
		AssertEqual(t, expected.count, actual)
	}
}

func TestGenerateCountWhenRangeHasLargeDifference(t *testing.T) {
	actual := generateCount(Range{5,11})
	AssertWithinRange(t, 5, 11, actual)
}

func TestCanGenerateSingleValueForReferenceFields(t *testing.T) {
	generator := NewGenerator("person", GetLogger(t))
	generator.WithField("field", "integer", [2]int{2,2}, Range{})
	reference := &ReferenceField{referred: generator, fieldName: "field"}
	value := reference.GenerateValue()
	referredValue:= generator.GetField("field").GenerateValue()
	AssertEqual(t, referredValue, value)
	AssertEqualTypes(t, referredValue, value)
}

func TestCanGenerateMultipleValuesForReferenceFields(t *testing.T) {
	generator := NewGenerator("person", GetLogger(t))
	generator.WithField("field", "integer", [2]int{2,2}, Range{2,2})
	reference := &ReferenceField{referred: generator, fieldName: "field"}
	value := reference.GenerateValue()
	AssertEqual(t, 2, len(value.([]interface{})))
}

func TestGenerateValueOnIntegerField(t *testing.T) {
	field := &IntegerField{2,2, Range{min: 2, max: 2}}
	actual := field.GenerateValue()
	AssertEqual(t, 2, actual)
}
