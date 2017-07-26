package generator

import (
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"reflect"
	"testing"
)

func TestAppendingToGeneratedContent(t *testing.T) {
	beast := []map[string]interface{}{{
		"of the beast": 666,
	}}

	expected := GeneratedContent{
		"sign": beast,
	}
	gc := NewGeneratedContent()
	gc.Append(expected)

	Assert(t, reflect.DeepEqual(gc, expected), "expected \n%v\n to be equal to \n%v\n but wasn't", gc, expected)

	rick := map[string]interface{}{
		"of Rick": "wubba lubba dub dub!!!!",
	}

	expected = GeneratedContent{
		"sign": append(beast, rick),
	}

	gc.Append(GeneratedContent{"sign": []map[string]interface{}{rick}})
	Assert(t, reflect.DeepEqual(gc, expected), "expected \n%v\n to be equal to \n%v\n but wasn't", gc, expected)
}
