package console

import (
	"strings"
)

type TextParser struct {
}

func NewTextParser() *TextParser {
	return &TextParser{}
}

func (textParser *TextParser) Parse(message string) UserMessage {
	if strings.HasPrefix(message, "auth") {
		return textParser.parseAuthMessage(message)
	} else if strings.HasPrefix(message, "join") {
		return textParser.parseJoinToRoomMessage(message)
	}
	return nil
}

func (textParser *TextParser) parseAuthMessage(message string) *UserAuthMessage {
	res := strings.Split(message, ":")
	if len(res) != 2 {
		return &UserAuthMessage{
			name:     "",
			password: "",
		}
	}
	command, payload := res[0], res[1]
	if command != "auth" || payload == "" {
		return &UserAuthMessage{
			name:     "",
			password: "",
		}
	}
	payloadParts := strings.Split(payload, "|")
	if len(payloadParts) != 2 {
		return &UserAuthMessage{
			name:     "",
			password: "",
		}
	}
	for i, part := range payloadParts {
		payloadParts[i] = strings.Trim(part, "{}")
	}
	userName, password := payloadParts[0], payloadParts[1]
	if userName == "" || password == "" {
		return &UserAuthMessage{
			name:     "",
			password: "",
		}
	}
	return &UserAuthMessage{
		name:     userName,
		password: password,
	}
}

func (textParser *TextParser) parseJoinToRoomMessage(message string) *UserJoinToRoomMessage {
	return nil
}
