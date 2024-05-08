package lexer

import (
	"go-learning/interpreter/token"
	"testing"
)

func TestLexer(testFramework *testing.T) {
	TestParseChars(testFramework)
	TestParseCodeBlock(testFramework)
	TestParseIfStatement(testFramework)
	TestNextToken(testFramework)
}

func TestParseChars(testFramework *testing.T) {
	input := "=+-!*/<>(){},;==!="

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.BANG, "!"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EQUAL, "=="},
		{token.NOT_EQUAL, "!="},

		//{token.EOF, ""},
	}

	lexer := New(input)

	runTests(lexer, tests, testFramework)
}

func TestParseCodeBlock(testFramework *testing.T) {
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

		{token.EOF, ""},
	}

	lexer := New(input)

	runTests(lexer, tests, testFramework)
}

func TestParseIfStatement(testFramework *testing.T) {
	input := `
		if (5 < 10) {
			return true
		} else {
			return false
		}
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
	}

	lexer := New(input)

	runTests(lexer, tests, testFramework)
}

func TestNextToken(testFramework *testing.T) {
	input := `
		let five = 5;
		let ten = 10;
		let add = fn(x, y) {
	       x + y;
	    };
		let result = add(five, ten);
		"foobar"
		"foo bar"
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
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},

		{token.EOF, ""},
	}

	lexer := New(input)

	runTests(lexer, tests, testFramework)
}

func runTests(lexer *Lexer, tests []struct {
	expectedType    token.TokenType
	expectedLiteral string
}, testFramework *testing.T) {
	for index, testToken := range tests {
		nextToken := lexer.NextToken()

		if nextToken.Type != testToken.expectedType {
			testFramework.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", index, testToken.expectedType, nextToken.Type)
		}

		if nextToken.Literal != testToken.expectedLiteral {
			testFramework.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", index, testToken.expectedLiteral, nextToken.Literal)
		}

	}
}
