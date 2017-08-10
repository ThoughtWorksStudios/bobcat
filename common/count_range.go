package common

import (
	"math/rand"
)

type CountRange struct {
	Min int
	Max int
}

func (b *CountRange) Multiple() bool {
	return b != nil
}

func (b *CountRange) Count() int {
	return determineCount(b.Min, b.Max)
}

func determineCount(min int, max int) int {
	if max == min {
		return max
	}

	return rand.Intn(max-min+1) + min
}
