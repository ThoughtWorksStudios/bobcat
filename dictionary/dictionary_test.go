package dictionary

import (
	"testing"
)

func TestSetLang(t *testing.T) {
	err := SetLang("en")
	if err != nil {
		t.Error("SetLang should successfully set lang")
	}
}

func TestFakerFullPath(t *testing.T) {
	expected := "/data/en/cat"
	actual := fullPath("en", "cat")
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
	UseExternalData(true)
	expected = "cat"
	actual = fullPath("en", "cat")
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	SetCustomDataLocation("/custom/path")
	expected = "/custom/path/kitty"
	actual = fullPath("en", "kitty")
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
	SetCustomDataLocation("")
	UseExternalData(false)
}

func TestFakerRuWithCallback(t *testing.T) {
	SetLang("ru")
	EnFallback(true)
	brand := lookup(lang, "companies", true)
	if brand == "" {
		t.Error("Fake call for name with no samples with callback should not return blank string")
	}
}

// TestConcurrentSafety runs fake methods in multiple go routines concurrently.
// This test should be run with the race detector enabled.
func TestConcurrentSafety(t *testing.T) {
	workerCount := 10
	doneChan := make(chan struct{})

	for i := 0; i < workerCount; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				lookup(lang, "first_names", true)
				lookup(lang, "last_names", true)
				lookup(lang, "genders", true)
				ValueFromDictionary("full_names")
				lookup(lang, "companies", true)
				lookup(lang, "companies", true)
			}
			doneChan <- struct{}{}
		}()
	}

	for i := 0; i < workerCount; i++ {
		<-doneChan
	}
}
