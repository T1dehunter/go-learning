package parser

import (
	"fmt"
	"go-learning/interpreter/lexer"
	"go-learning/interpreter/parser/ast"
	"go-learning/interpreter/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	nextToken token.Token
	errors    []string
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	parser.parseNextToken()
	parser.parseNextToken()

	return parser
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) addError(token token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", token, parser.nextToken.Type)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) parseNextToken() {
	parser.curToken = parser.nextToken
	parser.nextToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.curToken.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.parseNextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectNext(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.expectNext(token.ASSIGN) {
		return nil
	}

	//TODO: skip until semicolon will be done
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.parseNextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.curToken}

	parser.parseNextToken()

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.parseNextToken()
	}

	return stmt
}

func (parser *Parser) curTokenIs(token token.TokenType) bool {
	return parser.curToken.Type == token
}

// old name in book peekTokenIs
func (parser *Parser) nextTokenIs(token token.TokenType) bool {
	return parser.nextToken.Type == token
}

// old name in book expectPeek
func (parser *Parser) expectNext(token token.TokenType) bool {
	if parser.nextTokenIs(token) {
		parser.parseNextToken()
		return true
	} else {
		parser.addError(token)
		return false
	}
}
