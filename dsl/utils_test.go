package dsl

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"log"
	"testing"
	"time"
)

func TestAssembleTimeReturnsError(t *testing.T) {
	_, err := assembleTime("2017-07-19", []string{"13:00:00-0700"})

	ExpectsError(t, "Not a parsable timestamp", err)
}

func TestParseDateLikeJS(t *testing.T) {
	specs := map[string]time.Time{
		"2017-07-11":                parse("2017-07-11 00:00:00 +0000"),
		"2017-07-11T00:14:56":       parse("2017-07-11 00:14:56 +0000"),
		"2017-07-11T00:14:56Z":      parse("2017-07-11 00:14:56 +0000"),
		"2017-07-11T00:14:56-0730":  parse("2017-07-11 00:14:56 -0730"),
		"2017-07-11T00:14:56-08:30": parse("2017-07-11 00:14:56 -0830"),
	}

	for ts, expected := range specs {
		actual, err := parseDateLikeJS(ts)
		AssertNil(t, err, "Got an error while parsing date: %v", err)
		AssertTimeEqual(t, expected, actual)
	}
}

func TestParseDateLikeJSReturnsError(t *testing.T) {
	input := "2017-07-19T13:00:00Z-700"
	expected := "Not a parsable timestamp: 2017-07-19T13:00:00Z-700"
	_, err := parseDateLikeJS(input)
	ExpectsError(t, expected, err)
}

func TestRefReturnsLocationFromCurrent(t *testing.T) {
	c := &current{
		pos:         position{line: 4, col: 3, offset: 2},
		globalStore: map[string]interface{}{"filename": "whatever.spec"},
	}

	expected := NewLocation("whatever.spec", 4, 3, 2).String()
	actual := ref(c).String()
	AssertEqual(t, expected, actual)
}

func parse(stamp string) time.Time {
	t, e := time.Parse("2006-01-02 15:04:05 -0700", stamp)
	if e != nil {
		log.Fatalf("error? %v", e)
	}
	return t
}
