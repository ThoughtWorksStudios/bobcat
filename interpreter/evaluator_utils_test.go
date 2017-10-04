package interpreter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/generator"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
	"time"
)

const (
	DAY = int64(86400000) // 1 day in ms

	PRESENT = "1977-10-11T12:59:51-07:00"
	PAST    = "1977-10-10T12:59:51-07:00"
	FUTURE  = "1977-10-12T12:59:51-07:00"
)

func TestAddSubtractFromDateWorks(t *testing.T) {
	i := interp()

	start := mktime(PRESENT)
	timeDeltaFloat := float64(DAY) + 0.11

	newTime, err := i.addToTime("+", start, DAY, NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when adding time to a date")
	AssertEqual(t, mktime(FUTURE).Formatted(), newTime.(*generator.TimeWithFormat).Formatted(), "expected end date to be one day ahead of start date, which was: %q", start.Formatted())

	newTime, err = i.addToTime("-", start, DAY, NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when removing time from a date")
	AssertEqual(t, mktime(PAST).Formatted(), newTime.(*generator.TimeWithFormat).Formatted(), "expected end date to be one day behind start date, which was: %q", start.Formatted())

	newTime, err = i.addToTime("+", start, timeDeltaFloat, NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when adding time to a date")
	AssertEqual(t, mktime(FUTURE).Formatted(), newTime.(*generator.TimeWithFormat).Formatted(), "expected start and end date to be one day apart, start date was: %q", start.Formatted())

	duration, err := i.addToTime("-", mktime(FUTURE), start, NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when calculating duration")
	AssertEqual(t, DAY, duration, "expected (end_date - start_date) to yield 1 day (%d) ms", DAY)

	str, err := i.ApplyOperator("+", mktime(PRESENT), " ok", NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when concatenating")
	AssertEqual(t, "1977-10-11 ok", str)

	str, err = i.ApplyOperator("+", "date is ", mktime(PRESENT), NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when concatenating")
	AssertEqual(t, "date is 1977-10-11", str)
}

func mktime(t string) *generator.TimeWithFormat {
	ts, _ := time.Parse(time.RFC3339, t)
	return generator.NewTimeWithFormat(ts, "%Y-%m-%d")
}
