package console

import (
	"testing"
)

func TestNewTextParser(test *testing.T) {
	textParser := NewTextParser()
	if textParser == nil {
		test.Errorf("NewTextParser() returns nil")
	}
	testParseAuthMessage(test)
}

func testParseAuthMessage(test *testing.T) {
	input := "auth:{Sandor Clegane}|{Test1234%}"
	textParser := NewTextParser()
	userAuthMessage := textParser.parseAuthMessage(input)
	if userAuthMessage == nil {
		test.Errorf("parseAuthMessage() returns nil")
	}
	if userAuthMessage.name != "Sandor Clegane" {
		test.Errorf("parseAuthMessage() returns incorrect name, got: %s", userAuthMessage.name)
	}
	if userAuthMessage.password != "Test1234%" {
		test.Errorf("parseAuthMessage() returns incorrect password got: %s", userAuthMessage.password)
	}

	testData := []struct {
		input string
	}{
		{"auth:{Sandor Clegane}"},
		{"auth:{Test1234%}"},
		{"asdad"},
		{""},
	}
	for _, data := range testData {
		parsedMessage := textParser.parseAuthMessage(data.input)
		if parsedMessage != nil {
			test.Errorf("parseAuthMessage() returns incorrect value, expected nill, got: %v", parsedMessage)
		}
	}
}
