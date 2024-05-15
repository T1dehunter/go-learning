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

func (lexer *Lexer) peekNextChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.readPosition+1]
}

func (lexer *Lexer) NextToken() token.Token {
	var currentToken token.Token

	lexer.skipWhitespace()

	switch lexer.currentChar {
	case '=':
		if lexer.peekNextChar() == '=' {
			lexer.readNextChar()
			currentToken = token.Token{Type: token.EQUAL, Literal: token.EQUAL}
		} else {
			currentToken = createToken(token.ASSIGN, lexer.currentChar)
		}
	case '+':
		currentToken = createToken(token.PLUS, lexer.currentChar)
	case '-':
		currentToken = createToken(token.MINUS, lexer.currentChar)
	case '!':
		if lexer.peekNextChar() == '=' {
			currentChar := lexer.currentChar
			lexer.readNextChar()
			nextChar := lexer.currentChar
			literal := string(currentChar) + string(nextChar)
			currentToken = token.Token{Type: token.NOT_EQUAL, Literal: literal}
		} else {
			currentToken = createToken(token.BANG, lexer.currentChar)
		}
	case '*':
		currentToken = createToken(token.ASTERISK, lexer.currentChar)
	case '/':
		currentToken = createToken(token.SLASH, lexer.currentChar)
	case '<':
		currentToken = createToken(token.LT, lexer.currentChar)
	case '>':
		currentToken = createToken(token.GT, lexer.currentChar)
	case ';':
		currentToken = createToken(token.SEMICOLON, lexer.currentChar)
	case '(':
		currentToken = createToken(token.LPAREN, lexer.currentChar)
	case ')':
		currentToken = createToken(token.RPAREN, lexer.currentChar)
	case ',':
		currentToken = createToken(token.COMMA, lexer.currentChar)
	case '{':
		currentToken = createToken(token.LBRACE, lexer.currentChar)
	case '}':
		currentToken = createToken(token.RBRACE, lexer.currentChar)
	case '[':
		currentToken = createToken(token.LBRACKET, lexer.currentChar)
	case ']':
		currentToken = createToken(token.RBRACKET, lexer.currentChar)
	case '"':
		currentToken.Type = token.STRING
		currentToken.Literal = lexer.readString()
	case 0:
		currentToken.Literal = ""
		currentToken.Type = token.EOF
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
	identifierStartChar := lexer.position
	for isLetter(lexer.currentChar) {
		lexer.readNextChar()
	}
	identifierEndChar := lexer.position
	return lexer.input[identifierStartChar:identifierEndChar]
}

func (lexer *Lexer) readNumber() string {
	numberStartChar := lexer.position
	for isDigit(lexer.currentChar) {
		lexer.readNextChar()
	}
	numberEndChar := lexer.position
	return lexer.input[numberStartChar:numberEndChar]
}

func (lexer *Lexer) readString() string {
	stringStart := lexer.position + 1
	for {
		lexer.readNextChar()
		if lexer.currentChar == '"' || lexer.currentChar == 0 {
			break
		}
	}
	stringEnd := lexer.position
	return lexer.input[stringStart:stringEnd]
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
