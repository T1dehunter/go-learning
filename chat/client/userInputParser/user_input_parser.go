package userInputParser

import (
	"strings"
)

type UserInputParser struct {
}

func NewUserInputParser() *UserInputParser {
	return &UserInputParser{}
}

func (inputParser *UserInputParser) ParseCredentials(message string) (string, string) {
	res := strings.Split(message, ":")
	if len(res) != 2 {
		return "", ""
	}

	name, pass := res[0], res[1]
	if name == "" || pass == "" {
		return "", ""
	}

	name = strings.Trim(name, "[]")
	pass = strings.Trim(pass, "[]")
	if name == "" || pass == "" {
		return "", ""
	}

	return name, pass
}
