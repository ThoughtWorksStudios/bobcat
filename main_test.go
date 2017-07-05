package main

import "testing"

func TestValidSpec(t *testing.T) {
	_, err := parseSpec("testdata/valid_person.lang")
	if err != nil {
		t.Error("should not have thrown error", err)
	}
}

func TestInvalidSpec(t *testing.T) {
	_, err := parseSpec("testdata/invalid_person.lang")
	if err == nil {
		t.Error("should have thrown error")
	}

}
