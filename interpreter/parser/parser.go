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
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	//parser.registerPrefix(token.ELSE, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)

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

func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: parser.curToken, Value: parser.curToken.Literal}
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

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
		Left:     left,
	}

	precedence := parser.curPrecedence()

	parser.parseNextToken()

	expression.Right = parser.parseExpression(precedence)

	return expression
}

func (parser *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: parser.curToken, Value: parser.curTokenIs(token.TRUE)}
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.parseNextToken()

	exp := parser.parseExpression(LOWEST)

	if !parser.expectNext(token.RPAREN) {
		return nil
	}

	return exp
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.curToken}

	if !parser.expectNext(token.LPAREN) {
		return nil
	}

	parser.parseNextToken()

	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectNext(token.RPAREN) {
		return nil
	}

	if !parser.expectNext(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.nextTokenIs(token.ELSE) {
		parser.parseNextToken()

		if !parser.expectNext(token.LBRACE) {
			return nil
		}

		expression.Alternative = parser.parseBlockStatement()
	}

	return expression
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.curToken}
	block.Statements = []ast.Statement{}

	parser.parseNextToken()

	for !parser.curTokenIs(token.RBRACE) && !parser.curTokenIs(token.EOF) {
		stmt := parser.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		parser.parseNextToken()
	}

	return block
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	funcLietral := &ast.FunctionLiteral{Token: parser.curToken}

	if !parser.expectNext(token.LPAREN) {
		return nil
	}

	funcLietral.Parameters = parser.parseFunctionParameters()

	if !parser.expectNext(token.LBRACE) {
		return nil
	}

	funcLietral.Body = parser.parseBlockStatement()

	return funcLietral
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if parser.nextTokenIs(token.RPAREN) {
		parser.parseNextToken()
		return identifiers
	}

	parser.parseNextToken()

	ident := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	identifiers = append(identifiers, ident)

	for parser.nextTokenIs(token.COMMA) {
		// skip comma
		parser.parseNextToken()
		// parse next identifier
		parser.parseNextToken()
		ident := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !parser.expectNext(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: parser.curToken, Function: function}
	exp.Arguments = parser.parseCallArguments()
	return exp
}

func (parser *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if parser.nextTokenIs(token.RPAREN) {
		parser.parseNextToken()
		return args
	}

	parser.parseNextToken()
	args = append(args, parser.parseExpression(LOWEST))

	for parser.nextTokenIs(token.COMMA) {
		parser.parseNextToken()
		parser.parseNextToken()
		args = append(args, parser.parseExpression(LOWEST))
	}

	if !parser.expectNext(token.RPAREN) {
		return nil
	}

	return args
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
	stmt := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectNext(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.expectNext(token.ASSIGN) {
		return nil
	}

	parser.parseNextToken()

	stmt.Value = parser.parseExpression(LOWEST)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.parseNextToken()
	}

	return stmt
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.curToken}

	parser.parseNextToken()

	stmt.ReturnValue = parser.parseExpression(LOWEST)

	if parser.nextTokenIs(token.SEMICOLON) {
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
	isEndExpression := parser.nextTokenIs(token.SEMICOLON)
	isNextTokenHasHigherPrecedence := precedence < parser.nextTokenPrecedence()
	for !isEndExpression && isNextTokenHasHigherPrecedence {
		infix := parser.infixParseFns[parser.nextToken.Type]
		if infix == nil {
			return leftExp
		}

		parser.parseNextToken()

		leftExp = infix(leftExp)
	}

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

var precedences = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTERISK:  PRODUCT,
	token.LPAREN:    CALL,
}

func (parser *Parser) nextTokenPrecedence() int {
	if p, ok := precedences[parser.nextToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (parser *Parser) curPrecedence() int {
	if p, ok := precedences[parser.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
