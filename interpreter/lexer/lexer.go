package lexer

import (
	"go-learning/interpreter/token"
)

type CharactersMap map[byte]func() token.Token

type Lexer struct {
	input       string
	charPointer int  // points to the current char
	currentChar byte // current char under examination
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readFirstChar()
	return lexer
}

func (lexer *Lexer) readFirstChar() {
	lexer.currentChar = lexer.input[lexer.charPointer]
}

func (lexer *Lexer) readNextChar() {
	nexCharPointer := lexer.charPointer + 1
	if nexCharPointer >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		lexer.currentChar = lexer.input[nexCharPointer]
	}
	lexer.charPointer = nexCharPointer
}

func (lexer *Lexer) peekNextChar() byte {
	if lexer.charPointer >= len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.charPointer+1]
}

func (lexer *Lexer) NextToken() token.Token {
	var currentToken token.Token

	var charactersMap CharactersMap
	charactersMap = map[byte]func() token.Token{}

	lexer.skipWhitespace()

	charactersMap['='] = func() token.Token {
		if lexer.peekNextChar() == '=' {
			lexer.readNextChar()
			return token.Token{Type: token.EQUAL, Literal: token.EQUAL}
		}
		return createToken(token.ASSIGN, lexer.currentChar)
	}
	charactersMap['+'] = func() token.Token {
		return createToken(token.PLUS, lexer.currentChar)
	}
	charactersMap['-'] = func() token.Token {
		return createToken(token.MINUS, lexer.currentChar)
	}
	charactersMap['!'] = func() token.Token {
		if lexer.peekNextChar() == '=' {
			currentChar := lexer.currentChar
			lexer.readNextChar()
			nextChar := lexer.currentChar
			literal := string(currentChar) + string(nextChar)
			return token.Token{Type: token.NOT_EQUAL, Literal: literal}
		}
		return createToken(token.BANG, lexer.currentChar)
	}
	charactersMap['*'] = func() token.Token {
		return createToken(token.ASTERISK, lexer.currentChar)
	}
	charactersMap['/'] = func() token.Token {
		return createToken(token.SLASH, lexer.currentChar)
	}
	charactersMap['<'] = func() token.Token {
		return createToken(token.LT, lexer.currentChar)
	}
	charactersMap['>'] = func() token.Token {
		return createToken(token.GT, lexer.currentChar)
	}
	charactersMap[';'] = func() token.Token {
		return createToken(token.SEMICOLON, lexer.currentChar)
	}
	charactersMap[':'] = func() token.Token {
		return createToken(token.COLON, lexer.currentChar)
	}
	charactersMap['('] = func() token.Token {
		return createToken(token.LPAREN, lexer.currentChar)
	}
	charactersMap[')'] = func() token.Token {
		return createToken(token.RPAREN, lexer.currentChar)
	}
	charactersMap[','] = func() token.Token {
		return createToken(token.COMMA, lexer.currentChar)
	}
	charactersMap['{'] = func() token.Token {
		return createToken(token.LBRACE, lexer.currentChar)
	}
	charactersMap['}'] = func() token.Token {
		return createToken(token.RBRACE, lexer.currentChar)
	}
	charactersMap['['] = func() token.Token {
		return createToken(token.LBRACKET, lexer.currentChar)
	}
	charactersMap[']'] = func() token.Token {
		return createToken(token.RBRACKET, lexer.currentChar)
	}
	charactersMap['"'] = func() token.Token {
		currentToken.Type = token.STRING
		currentToken.Literal = lexer.readString()
		return currentToken
	}
	charactersMap[0] = func() token.Token {
		currentToken.Literal = ""
		currentToken.Type = token.EOF
		return currentToken
	}

	if charHandler, ok := charactersMap[lexer.currentChar]; ok {
		currentToken = charHandler()
		lexer.readNextChar()
		return currentToken
	}

	if isLetter(lexer.currentChar) {
		currentToken.Literal = lexer.readIdentifier()
		currentToken.Type = token.LookupIdent(currentToken.Literal)
		return currentToken
	}

	if isDigit(lexer.currentChar) {
		currentToken.Literal = lexer.readNumber()
		currentToken.Type = token.INT
		return currentToken
	}

	currentToken = createToken(token.ILLEGAL, lexer.currentChar)

	lexer.readNextChar()

	return currentToken
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\t' || lexer.currentChar == '\n' || lexer.currentChar == '\r' {
		lexer.readNextChar()
	}
}

func (lexer *Lexer) readIdentifier() string {
	identifierStartChar := lexer.charPointer
	for isLetter(lexer.currentChar) {
		lexer.readNextChar()
	}
	identifierEndChar := lexer.charPointer
	return lexer.input[identifierStartChar:identifierEndChar]
}

func (lexer *Lexer) readNumber() string {
	numberStartChar := lexer.charPointer
	for isDigit(lexer.currentChar) {
		lexer.readNextChar()
	}
	numberEndChar := lexer.charPointer
	return lexer.input[numberStartChar:numberEndChar]
}

func (lexer *Lexer) readString() string {
	stringStart := lexer.charPointer + 1
	for {
		lexer.readNextChar()
		if lexer.currentChar == '"' || lexer.currentChar == 0 {
			break
		}
	}
	stringEnd := lexer.charPointer
	return lexer.input[stringStart:stringEnd]
}

func createToken(tokenType token.TokenType, currentChar byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(currentChar)}
}
