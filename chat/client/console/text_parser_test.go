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
		input         string
		parsedMessage UserAuthMessage
	}{
		{"auth:{Sandor Clegane}|{Test1234%}", UserAuthMessage{"Sandor Clegane", "Test1234%"}},
		{"auth:{Sandor Clegane}", UserAuthMessage{"", ""}},
		{"auth:{Test1234%}", UserAuthMessage{"", ""}},
		{"asdad", UserAuthMessage{"", ""}},
		{"", UserAuthMessage{"", ""}},
	}
	for _, data := range testData {
		userAuthMessage := textParser.parseAuthMessage(data.input)
		if userAuthMessage.name != data.parsedMessage.name {
			test.Errorf("parseAuthMessage() returns incorrect name, got: %s", userAuthMessage.name)
		}
		if userAuthMessage.password != data.parsedMessage.password {
			test.Errorf("parseAuthMessage() returns incorrect password got: %s", userAuthMessage.password)
		}
	}
}
