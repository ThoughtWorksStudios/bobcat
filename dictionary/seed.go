package dictionary

import (
	"math/rand"
	"sync"
	"time"
)

var r = rand.New(&rndSrc{src: rand.NewSource(time.Now().UnixNano())})

// Seed uses the provided seed value to initialize the internal PRNG to a
// deterministic state.
func Seed(seed int64) {
	r.Seed(seed)
}

type rndSrc struct {
	mtx sync.Mutex
	src rand.Source
}

func (s *rndSrc) Int63() int64 {
	s.mtx.Lock()
	n := s.src.Int63()
	s.mtx.Unlock()
	return n
}

func (s *rndSrc) Seed(n int64) {
	s.mtx.Lock()
	s.src.Seed(n)
	s.mtx.Unlock()
}
