package interpreter

import (
	g "github.com/ThoughtWorksStudios/datagen/generator"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"reflect"
	"testing"
)

func TestAppendingToGenerationOutput(t *testing.T) {
	actual := GenerationOutput{}

	beast := g.GeneratedEntities{g.GeneratedFields{"of the beast": 666}}

	expected := GenerationOutput{"sign": beast}
	actual.addAndAppend("sign", beast)

	Assert(t, reflect.DeepEqual(expected, actual), "expected \n%v\n to be equal to \n%v\n but wasn't", expected, actual)

	rick := g.GeneratedFields{
		"of Rick": "wubba lubba dub dub!!!!",
	}

	expected = GenerationOutput{
		"sign": append(beast, rick),
	}

	actual.addAndAppend("sign", g.GeneratedEntities{rick})
	Assert(t, reflect.DeepEqual(expected, actual), "expected \n%v\n to be equal to \n%v\n but wasn't", expected, actual)
}
