package dictionary

import (
	"strconv"
	"strings"
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
	SetLang("en")
	EnFallback(true)
	brand := lookup(lang, "companies", true)
	if brand == "" {
		t.Error("Fake call for name with no samples with callback should not return blank string")
	}
}

func TestCompositeFormat(t *testing.T) {
	result := compositeFormat("first_names| |last_names")
	if result == "" {
		t.Error("Expected to get results, but got nothing :(")
	}
	components := strings.Split(result, " ")

	if len(components) != 2 {
		t.Errorf("Expected to get only 2 components (first name and last name), but got %v of length %v", result, len(components))
	}

	for _, v := range components {
		if v == "" {
			t.Error("Expected to get results, but got nothing :(")
		}
	}
}

func TestCompositeFormatShouldProcessSubFormats(t *testing.T) {
	result := compositeFormat("first_names| |full_names_format")
	if result == "" {
		t.Error("Expected to get results, but got nothing :(")
	}

	components := strings.Split(result, " ")

	if len(components) != 3 {
		t.Errorf("Expected to get only 3 components (first name and last name), but got %v of length %v", result, len(components))
	}

	for _, v := range components {
		if v == "" {
			t.Error("Expected to get results, but got nothing :(")
		}
	}
}

func TestCompositeFormatWithSubFormatCompositeComponents(t *testing.T) {
	result := compositeFormat("email_address_format| |phone_numbers_format")
	if result == "" {
		t.Error("Expected to get results, but got nothing :(")
	}

	components := strings.Split(result, " ")

	if len(components) != 2 {
		t.Errorf("Expected to get only 2 components, but got %v of length %v", result, len(components))
	}

	if components[0] == "email_address_format" {
		t.Errorf("The dictionary component was not processed!")
	}

	if components[1] == "phone_numbers_format" {
		t.Errorf("The dictionary component was not processed!")
	}

	for _, part := range strings.Split(components[1], "-") {
		if num, err := strconv.Atoi(part); err != nil {
			t.Errorf("Expected to get an integer back, but got: %v", num)
		}
	}

}

func TestvalueFromFormat(t *testing.T) {
	result := valueFromFormat("###")
	if result == "" {
		t.Error("Expected to get results, but got nothing :(")
	}

	if num, err := strconv.Atoi(result); err != nil {
		t.Errorf("Expected to get an integer back, but got: %v", num)
	}
}

func TestvalueFromFormatWithCompositeComponents(t *testing.T) {
	result := valueFromFormat("first_names| |###")
	if result == "" {
		t.Error("Expected to get results, but got nothing :(")
	}

	components := strings.Split(result, " ")

	if len(components) != 2 {
		t.Errorf("Expected to get only 2 components, but got %v of length %v", result, len(components))
	}

	if components[0] == "first_names" {
		t.Errorf("The dictionary component was not processed!")
	}

	if num, err := strconv.Atoi(components[1]); err != nil {
		t.Errorf("Expected to get an integer back, but got: %v", num)
	}
}

func TestValueFromDictionaryShouldTakeFormatWithoutFormatSuffix(t *testing.T) {
	result := ValueFromDictionary("full_names")
	if result == "" {
		t.Error("Expected to get results, but got nothing :(")
	}

	components := strings.Split(result, " ")

	if len(components) != 2 {
		t.Errorf("Expected to get only 2 components, but got %v of length %v", result, len(components))
	}

	if components[0] == "first_names" || components[1] == "last_names" {
		t.Errorf("The dictionary component was not processed!")
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
				ValueFromDictionary("email_address")
				ValueFromDictionary("email_address")
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
