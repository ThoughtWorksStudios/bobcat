package common

import (
	"fmt"
	"math/rand"
)

type CountRange struct {
	Min int64
	Max int64
}

func (r *CountRange) Multiple() bool {
	return r != nil
}

func (r *CountRange) Count() int64 {
	return determineCount(r.Min, r.Max)
}

func (r *CountRange) Validate() error {
	if r.Min < 0 || r.Max < 0 {
		return fmt.Errorf("Count range bounds must not be negative")
	}

	if r.Max < r.Min {
		return fmt.Errorf("Count range max cannot be less than min")
	}

	return nil
}

func determineCount(min int64, max int64) int64 {
	if max == min {
		return max
	}

	return rand.Int63n(max-min+1) + min
}
