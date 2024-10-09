package userInputParser

import (
	"testing"
)

func TestNewTextParser(test *testing.T) {
	textParser := NewUserInputParser()
	if textParser == nil {
		test.Errorf("NewUserInputParser() returns nil")
	}
	testParseCredentials(test)
}

func testParseCredentials(test *testing.T) {
	textParser := NewUserInputParser()

	testData := []struct {
		input        string
		expectedName string
		expectedPass string
	}{
		{"[Sandor Clegane]:[Test1234%]", "Sandor Clegane", "Test1234%"},
		{"[Sandor Clegane]|[Test1234%]", "", ""},
		{"[Sandor Clegane]", "", ""},
		{"[]:[]", "", ""},
		{"asdad", "", ""},
		{":", "", ""},
		{"", "", ""},
	}

	for _, data := range testData {
		name, pass := textParser.ParseCredentials(data.input)
		if name != data.expectedName || pass != data.expectedPass {
			test.Errorf("ParseCredentials() returns incorrect values, expected %s, %s, got: %s, %s", data.expectedName, data.expectedPass, name, pass)
		}
	}
}
