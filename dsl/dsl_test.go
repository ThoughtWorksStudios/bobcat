package dsl

import (
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"testing"
)

func RequiresDefOrGenerateStatements(t *testing.T) {
	_, err := Parse("", []byte("eek"))
	expectedErrorMsg := "1:1 (0): no match found, expected: \"def\", \"generate\", [ \t\r\n] or EOF"
	ExpectsError(t, expectedErrorMsg, err)
}
