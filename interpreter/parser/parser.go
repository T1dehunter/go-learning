package parser

import (
	"fmt"
	"go-learning/interpreter/lexer"
	"go-learning/interpreter/parser/ast"
	"go-learning/interpreter/token"
	"strconv"
)

type prefixParseFn func() ast.Expression

type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	curToken  token.Token
	nextToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	parser.parseNextToken()
	parser.parseNextToken()

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)

	return parser
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
	}

	parser.parseNextToken()

	expression.Right = parser.parseExpression(PREFIX)

	return expression
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
		return parser.parseExpressionStatement()
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

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFns[parser.curToken.Type]
	if prefix == nil {
		parser.noPrefixParseFnError(parser.curToken.Type)
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser.curToken}
	statement.Expression = parser.parseExpression(LOWEST)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.parseNextToken()
	}

	return statement
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

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
}

func (parser *Parser) noPrefixParseFnError(tokenType token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	parser.errors = append(parser.errors, message)
}
