package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestFetchOrCreate(t *testing.T) {
	em, p := make(EmitterMap), &TestProvider{}
	AssertEqual(t, 0, len(em), "Baseline - EmitterMap should not have any entries")
	emitter, _ := em.FetchOrCreate("foo", p)
	AssertEqual(t, emitter, em["foo"], "Should create new emitter and associate with key")
	fetched, _ := em.FetchOrCreate("foo", p)
	AssertEqual(t, emitter, fetched, "Should return the same emitter in subsequent FetchOrCreate() using same key")
}

func TestFinalize(t *testing.T) {
	em, p := make(EmitterMap), &TestProvider{}
	em.FetchOrCreate("foo", p)
	em.FetchOrCreate("bar", p)

	AssertEqual(t, 2, len(em), "Should have 2 emitters")

	for _, e := range em {
		te := e.(*TestEmitter)
		Assert(t, !te.Closed(), "Emitters should not be finalized yet")
	}

	em.Finalize()

	for _, e := range em {
		te := e.(*TestEmitter)
		Assert(t, te.Closed(), "All emitters should be finalized")
	}
}

type TestProvider struct{}

func (p *TestProvider) Get(key string) (Emitter, error) {
	return NewTestEmitter(), nil
}
