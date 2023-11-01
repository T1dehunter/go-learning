package lexer

import (
	"fmt"
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

	lexer.skipWhitespace()

	fmt.Printf("currentChar=%q", lexer.currentChar)

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
	default:
		if isLetter(lexer.currentChar) {
			currentToken.Literal = lexer.readIdentifier()
			currentToken.Type = token.LookupIdent(currentToken.Literal)
			return currentToken
		} else if isDigit(lexer.currentChar) {
			currentToken.Literal = lexer.readNumber()
			currentToken.Type = token.INT
			return currentToken
		} else {
			fmt.Printf("readPosition=%d", lexer.readPosition)
			fmt.Printf("currentChar=%b", lexer.currentChar)
			currentToken = createToken(token.ILLEGAL, lexer.currentChar)
		}
	}
	lexer.readNextChar()
	return currentToken
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\t' || lexer.currentChar == '\n' || lexer.currentChar == '\r' {
		lexer.readNextChar()
	}
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isLetter(lexer.currentChar) {
		lexer.readNextChar()
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.currentChar) {
		lexer.readNextChar()
	}
	return lexer.input[position:lexer.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func createToken(tokenType token.TokenType, currentChar byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(currentChar)}
}
