package lexer

import (
	"go-learning/interpreter/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current currentChar)
	readPosition int  // current reading position in input (after current currentChar)
	currentChar  byte // current currentChar under examination
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readCurrentChar()
	return lexer
}

func (lexer *Lexer) readCurrentChar() {
	lexer.currentChar = lexer.input[lexer.readPosition]
}

func (lexer *Lexer) readNextChar() {
	lexer.readPosition += 1
	if lexer.readPosition >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		lexer.currentChar = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
}

func (lexer *Lexer) NextToken() token.Token {
	var currentToken token.Token
	switch lexer.currentChar {
	case '=':
		currentToken = createToken(token.ASSIGN, lexer.currentChar)
	case ';':
		currentToken = createToken(token.SEMICOLON, lexer.currentChar)
	case '(':
		currentToken = createToken(token.LPAREN, lexer.currentChar)
	case ')':
		currentToken = createToken(token.RPAREN, lexer.currentChar)
	case ',':
		currentToken = createToken(token.COMMA, lexer.currentChar)
	case '+':
		currentToken = createToken(token.PLUS, lexer.currentChar)
	case '{':
		currentToken = createToken(token.LBRACE, lexer.currentChar)
	case '}':
		currentToken = createToken(token.RBRACE, lexer.currentChar)
	case 0:
		currentToken.Literal = ""
		currentToken.Type = token.EOF
	}
	lexer.readNextChar()
	return currentToken
}

func createToken(tokenType token.TokenType, currentChar byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(currentChar)}
}
