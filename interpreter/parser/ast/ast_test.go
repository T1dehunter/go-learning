package ast

import (
	"fmt"
	"go-learning/interpreter/token"
	"testing"
)

func TestString(test *testing.T) {
	program := Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		test.Errorf("program.String() wrong. got=%q", program.String())
	}
	testFFF()
}

type User interface {
	Name() string
}

type MyUser struct {
	name string
}

func (receiver *MyUser) Name() string {
	return receiver.name
}

type Data struct {
	Users []User
}

func testFFF() {
	data := Data{
		Users: []User{&MyUser{name: "TEST"}},
	}
	fmt.Println("TEST: ", data)
}
