package lexer

import (
	"go-learning/interpreter/token"
	"testing"
)

func TestNextToken(test *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for index, testToken := range tests {
		nextToken := lexer.NextToken()

		if nextToken.Type != testToken.expectedType {
			test.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", index, testToken.expectedType, nextToken.Type)
		}

		if nextToken.Literal != testToken.expectedLiteral {
			test.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", index, testToken.expectedLiteral, nextToken.Literal)
		}

	}

}
