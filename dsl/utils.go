package dsl

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"regexp"
	"strings"
	"time"
)

func assembleTime(date, localTime interface{}) (time.Time, error) {
	iso8601Date := date.(string)
	var ts []string

	if localTime != nil {
		ts = localTime.([]string)
	}

	str := strings.Join(append([]string{iso8601Date}, ts...), "")
	return parseDateLikeJS(str)
}

/**
 * Parses date and date + timestamp in ISO-8601 variations just like
 * JavaScript. Specifically:
 *
 * YYYY-MM-DD
 * YYYY-mm-ddTHH:MM:SS
 * YYYY-mm-ddTHH:MM:SSZ
 * YYYY-mm-ddTHH:MM:SS-0000
 * YYYY-mm-ddTHH:MM:SS-00:00
 */
func parseDateLikeJS(tstamp string) (time.Time, error) {
	// you'll just have to take my word on this
	re := regexp.MustCompile("^([0-9]{4}-[0-9]{2}-[0-9]{2})(?:(T[0-9]{2}:[0-9]{2}:[0-9]{2})(Z|(?:[+-][0-9]{2}:?[0-9]{2}))?)?$")

	format := "2006-01-02" // default to parsing only the date

	m := re.FindStringSubmatch(tstamp)

	if m == nil {
		return time.Time{}, fmt.Errorf("Not a parsable timestamp: %s", tstamp)
	}

	parts := []string{m[1], m[2], strings.Replace(strings.Replace(m[3], ":", "", -1), "Z", "", -1)}

	if m[2] != "" {
		format = format + "T15:04:05"
	}

	if m[3] != "" && m[3] != "Z" {
		format = format + "-0700"
	}

	return time.Parse(format, strings.Join(parts, ""))
}

/** Extracts a *Location value from *current */
func ref(c *current) *Location {
	filename, _ := c.globalStore["filename"].(string)
	return NewLocation(
		filename,
		c.pos.line,
		c.pos.col,
		c.pos.offset,
	)
}
