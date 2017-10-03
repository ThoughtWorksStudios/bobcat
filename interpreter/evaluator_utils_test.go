package interpreter

import (
	"testing"
	"time"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
)

func TestAddSubtractFromDateWorks(t *testing.T) {
	i := interp()
	startTime, _ := time.Parse(time.RFC3339, "1977-10-11T12:59:51-07:00")
	futureTime, _ := time.Parse(time.RFC3339, "1977-10-12T12:59:51-07:00")
	pastTime, _ := time.Parse(time.RFC3339, "1977-10-10T12:59:51-07:00")
	timeDelta := int64(86400)
	timeDeltaFloat := float64(86400.11)

	newTime, err := i.addToDate("+", startTime, timeDelta, NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when adding time to a date")
	AssertEqual(t, futureTime, newTime, "expected end date to be one day ahead of start date, which was: '%v'", startTime)

	newTime, err = i.addToDate("-", startTime, timeDelta, NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when removing time from a date")
	AssertEqual(t, pastTime, newTime, "expected end date to be one day behind start date, which was: '%v'", startTime)

	newTime, err = i.addToDate("+", startTime, timeDeltaFloat, NewRootScope(), false)
	AssertNil(t, err, "should not receive an error when adding time to a date")
	AssertEqual(t, futureTime, newTime, "expected start and end date to be one day apart, start date was: '%v'", startTime)
}