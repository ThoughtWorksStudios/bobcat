package generator

import (
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"reflect"
	"testing"
)

func TestAppendingToGeneratedContent(t *testing.T) {
	actual := NewGeneratedContent()

	beast := GeneratedEntities{GeneratedFields{"of the beast": 666}}

	expected := GeneratedContent{"sign": beast}
	actual.Append(expected)

	Assert(t, reflect.DeepEqual(expected, actual), "expected \n%v\n to be equal to \n%v\n but wasn't", expected, actual)

	rick := GeneratedFields{
		"of Rick": "wubba lubba dub dub!!!!",
	}

	expected = GeneratedContent{
		"sign": append(beast, rick),
	}

	actual.Append(GeneratedContent{"sign": GeneratedEntities{rick}})
	Assert(t, reflect.DeepEqual(expected, actual), "expected \n%v\n to be equal to \n%v\n but wasn't", expected, actual)
}
