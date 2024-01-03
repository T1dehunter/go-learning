package ast

import "go-learning/interpreter/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program data
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement data
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) statementNode() {

}

func (letStatement *LetStatement) TokenLiteral() string {
	return letStatement.Token.Literal
}

// Identifier data
type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode() {

}

func (identifier *Identifier) TokenLiteral() string {
	return identifier.Token.Literal
}

// ReturnStatement data
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStmt *ReturnStatement) statementNode() {

}

func (returnStmt *ReturnStatement) TokenLiteral() string {
	return returnStmt.Token.Literal
}
