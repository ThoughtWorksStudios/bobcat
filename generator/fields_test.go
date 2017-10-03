package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestGenerateEntity(t *testing.T) {
	g := NewGenerator("testEntity", nil, false)
	fieldType := &EntityType{g}
	emitter := NewTestEmitter()
	subId, err := fieldType.One(nil, emitter, nil)
	AssertNil(t, err, "Should not receive error")

	e := emitter.Shift()

	if nil == e {
		t.Errorf("Expected to generate an entity but got %T %v", e, e)
	}

	AssertEqual(t, "testEntity", e["$type"], "Should have generated an entity of type \"testEntity\"")
	AssertEqual(t, subId, e[g.PrimaryKeyName()])
}

func TestMultiValueGenerate(t *testing.T) {
	field := NewField(NewLiteralType("foo"), &CountRange{3, 3})

	v, err := field.GenerateValue(nil, NewDummyEmitter(), nil)
	AssertNil(t, err, "Should not receive error")

	actual := len(v.([]interface{}))
	AssertEqual(t, 3, actual)
}

func TestDeferredType(t *testing.T) {
	closure := func(scope *Scope) (interface{}, error) {
		return scope.ResolveSymbol("bar"), nil
	}
	scope := NewRootScope()
	scope.SetSymbol("bar", "foo")

	generatedType := NewDeferredType(closure)
	actual, err := generatedType.One(nil, nil, scope)
	AssertNil(t, err, "Should not receive error")
	AssertEqual(t, "foo", actual)
}

func TestLiteralType(t *testing.T) {
	generatedType := NewLiteralType("foo")
	actual, err := generatedType.One(nil, nil, nil)
	AssertNil(t, err, "Should not receive error")
	AssertEqual(t, "foo", actual)
}
