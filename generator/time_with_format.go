package generator

import (
	"github.com/jehiah/go-strftime"
	"strconv"
	"time"
)

type TimeWithFormat struct {
	Format string
	Time   time.Time
}

func (t *TimeWithFormat) Formatted() string {
	return strftime.Format(t.Format, t.Time)
}

func (t *TimeWithFormat) MarshalJSON() ([]byte, error) {
	if "" == t.Format {
		return t.Time.MarshalJSON()
	}
	return []byte(strconv.Quote(t.Formatted())), nil
}

func NewTimeWithFormat(t time.Time, format string) *TimeWithFormat {
	return &TimeWithFormat{Time: t, Format: format}
}
