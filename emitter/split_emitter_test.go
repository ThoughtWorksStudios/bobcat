package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestEmitCreatesEmitterPerType(t *testing.T) {
	se := &SplitEmitter{provider: &TestProvider{}, emitters: make(EmitterMap)}

	AssertEqual(t, 0, len(se.emitters), "Should start with empty map")
	se.Emit(EntityResult{"foo": 0}, "foo")
	se.Emit(EntityResult{"bar": 1}, "bar")

	AssertNotNil(t, se.emitters["foo"], "Should create emitter for foo type")
	AssertNotNil(t, se.emitters["bar"], "Should create emitter for bar type")

	foomitter := se.emitters["foo"].(*TestEmitter)
	barmitter := se.emitters["bar"].(*TestEmitter)
	AssertNotEqual(t, foomitter, barmitter, "Emitters should be different instances per type")
}
