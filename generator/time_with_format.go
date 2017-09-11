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

func (t *TimeWithFormat) MarshalJSON() ([]byte, error) {
	if "" == t.Format {
		return t.Time.MarshalJSON()
	}
	return []byte(strconv.Quote(strftime.Format(t.Format, t.Time))), nil
}
