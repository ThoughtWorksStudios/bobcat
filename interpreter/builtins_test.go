package interpreter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/generator"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"regexp"
	"strings"
	"testing"
	"time"
)

type Args []interface{}

func TestIntBuiltin(t *testing.T) {
	min, max := int64(5), int64(7)

	CallBuiltin(t, INT_TYPE, Args{min, max}, func(result interface{}) {
		actual, ok := result.(int64)
		Assert(t, ok, "Should have generated a int64, but was %T", result)

		if actual < min || actual > max {
			t.Errorf("Generated value '%v' is outside of expected range min: '%v', max: '%v'", actual, min, max)
		}
	})
}

func TestFloatBuiltin(t *testing.T) {
	min, max := 4.25, 4.3

	CallBuiltin(t, FLOAT_TYPE, Args{min, max}, func(result interface{}) {
		actual, ok := result.(float64)
		Assert(t, ok, "Should have generated a float64, but was %T", result)

		if actual < min || actual > max {
			t.Errorf("Generated value '%v' is outside of expected range min: '%v', max: '%v'", actual, min, max)
		}
	})
}

func TestStringBuiltin(t *testing.T) {
	l := int64(8)

	CallBuiltin(t, STRING_TYPE, Args{l}, func(result interface{}) {
		actual, ok := result.(string)
		Assert(t, ok, "Should have generated a string, but was %T", result)

		if len(actual) != int(l) {
			t.Errorf("Generated value '%v' should be of length", actual, l)
		}
	})
}

func TestBoolBuiltin(t *testing.T) {
	CallBuiltin(t, BOOL_TYPE, Args{}, func(result interface{}) {
		_, ok := result.(bool)
		Assert(t, ok, "Should have generated a bool, but was %T", result)
	})
}

func TestSerialBuiltin(t *testing.T) {
	builtin, err := NewBuiltin(SERIAL_TYPE)
	AssertNil(t, err, "Should not receive error instantiating builtin %q", SERIAL_TYPE)

	result, err := builtin.Call()
	AssertNil(t, err, "Should not receive error calling %q()", SERIAL_TYPE)
	_, ok := result.(uint64)
	Assert(t, ok, "Should have generated a bool, but was %T", result)
	AssertEqual(t, uint64(0), result, "Starts at zero")

	result, err = builtin.Call()
	AssertNil(t, err, "Should not receive error calling %q()", SERIAL_TYPE)
	AssertEqual(t, uint64(1), result, "Increments to 1")

	result, err = builtin.Call()
	AssertNil(t, err, "Should not receive error calling %q()", SERIAL_TYPE)
	AssertEqual(t, uint64(2), result, "Increments to 2")

	result, err = builtin.Call(int64(100))
	AssertNil(t, err, "Should not receive error calling %q()", SERIAL_TYPE)
	AssertEqual(t, uint64(103), result, "Can add an offset to generated number to effectively start sequence elsewhere")
}

func TestUniqueIntBuiltin(t *testing.T) {
	CallBuiltin(t, UNIQUE_INT_TYPE, Args{}, func(result interface{}) {
		_, ok := result.(uint64)
		Assert(t, ok, "Should have generated a uint64, but was %T", result)
	})
}

func TestUidBuiltin(t *testing.T) {
	CallBuiltin(t, UID_TYPE, Args{}, func(result interface{}) {
		actual, ok := result.(string)
		Assert(t, ok, "Should have generated a string, but was %T", result)

		AssertEqual(t, 20, len(actual), "UIDs are 20 characters")
	})
}

func TestDateBuiltin(t *testing.T) {
	min, _ := time.Parse("2006-01-02", "1945-01-01")
	max, _ := time.Parse("2006-01-02", "1945-01-02")

	CallBuiltin(t, DATE_TYPE, Args{min, max}, func(result interface{}) {
		actual, ok := result.(*generator.TimeWithFormat)
		Assert(t, ok, "Should have generated a *generator.TimeWithFormat, but was %T", result)

		if actual.Time.Before(min) || actual.Time.After(max) {
			t.Errorf("Generated value '%v' is outside of expected range min: '%v', max: '%v'", actual.Time, min, max)
		}
	})
}

func TestDictBuiltin(t *testing.T) {
	expected := regexp.MustCompile("^[\\w]+@[\\w]+\\.(?:[a-z]{3,4})$")
	CallBuiltin(t, DICT_TYPE, Args{"email_address"}, func(result interface{}) {
		actual, ok := result.(string)
		Assert(t, ok, "Should have generated a string, but was %T", result)
		Assert(t, expected.Match([]byte(actual)), "Should have returned an email address, but was %q", actual)
	})
}

func TestEnumBuiltin(t *testing.T) {
	collection := []interface{}{"a", "B", "z"}
	CallBuiltin(t, ENUM_TYPE, Args{collection}, func(result interface{}) {
		actual, ok := result.(string)
		Assert(t, ok, "Should have generated a string because all elements in collection are strings, but was %T", result)
		AssertEqual(t, 1, len(actual), "Should pick an element from collection, which in this case are all single chars")
		Assert(t, strings.Contains("aBz", actual), "Should have picked an element from the collection, but was %q", actual)
	})
}

func CallBuiltin(t *testing.T, name string, args Args, assertion func(result interface{})) {
	builtin, err := NewBuiltin(name)
	AssertNil(t, err, "Should not receive error instantiating builtin %q", name)

	result, err := builtin.Call(args...)
	AssertNil(t, err, "Should not receive error calling %q()", name)

	assertion(result)
}
