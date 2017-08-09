package common

import (
  "math/rand"
)

type Bound struct {
  Min int
  Max int
}

func(b *Bound) Multiple() bool {
  return b != nil
}

func (b *Bound) Amount() int {
  return determineAmount(b.Min, b.Max)
}

func determineAmount(min int, max int) int {
  if max == 0 && min == 0 {
    return 1
  } else if max - min == 0 {
    return min
  }

  return rand.Intn(max - min + 1) + min
}
