package interpreter

import (
	g "github.com/ThoughtWorksStudios/bobcat/generator"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"reflect"
	"testing"
)

func TestAppendingToNestedOutput(t *testing.T) {
	actual := NestedOutput{}

	beast := g.GeneratedEntities{g.EntityResult{"of the beast": 666}}

	expected := NestedOutput{"sign": beast}
	actual.addAndAppend("sign", beast)

	Assert(t, reflect.DeepEqual(expected, actual), "expected \n%v\n to be equal to \n%v\n but wasn't", expected, actual)

	rick := g.EntityResult{
		"of Rick": "wubba lubba dub dub!!!!",
	}

	expected = NestedOutput{
		"sign": append(beast, rick),
	}

	actual.addAndAppend("sign", g.GeneratedEntities{rick})
	Assert(t, reflect.DeepEqual(expected, actual), "expected \n%v\n to be equal to \n%v\n but wasn't", expected, actual)
}
