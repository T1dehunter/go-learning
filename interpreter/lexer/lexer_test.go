package lexer

import (
	"go-learning/interpreter/token"
	"testing"
)

func TestLexer(test *testing.T) {
	TestParseChars(test)
	TestParseCodeBlock(test)
}

func TestParseChars(test *testing.T) {
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
		//{token.EOF, ""},
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

func TestParseCodeBlock(test *testing.T) {
	input := `
		let five = 5;
		let ten = 10;
		let add = fn(x, y) {
	       x + y;
	    };
		let result = add(five, ten);
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		//{token.EOF, ""},
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
