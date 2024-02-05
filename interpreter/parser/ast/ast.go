package ast

import (
	"bytes"
	"go-learning/interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, str := range p.Statements {
		out.WriteString(str.String())
	}

	return out.String()
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

func (letStatement *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(letStatement.TokenLiteral() + " ")
	out.WriteString(letStatement.Name.String())
	out.WriteString(" = ")

	if letStatement.Value != nil {
		out.WriteString(letStatement.Value.String())
	}

	out.WriteString(";")

	return out.String()
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

func (identifier *Identifier) String() string { return identifier.Value }

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

func (returnStmt *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(returnStmt.TokenLiteral() + " ")

	if returnStmt.ReturnValue != nil {
		out.WriteString(returnStmt.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

func (returnStmt *ReturnStatement) statementNode() {

}

func (returnStmt *ReturnStatement) TokenLiteral() string {
	return returnStmt.Token.Literal
}

// ExpressionStatement data
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expression *ExpressionStatement) String() string {
	if expression.Expression != nil {
		return expression.Expression.String()
	}
	return ""
}

func (expression *ExpressionStatement) statementNode() {

}

func (expression *ExpressionStatement) TokenLiteral() string {
	return expression.Token.Literal
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (intLiteral *IntegerLiteral) expressionNode()      {}
func (intLiteral *IntegerLiteral) TokenLiteral() string { return intLiteral.Token.Literal }
func (intLiteral *IntegerLiteral) String() string       { return intLiteral.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (prefixExpr *PrefixExpression) expressionNode() {}

func (prefixExpr *PrefixExpression) TokenLiteral() string {
	return prefixExpr.Token.Literal
}

func (prefixExpr *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(prefixExpr.Operator)
	out.WriteString(prefixExpr.Right.String())
	out.WriteString(")")

	return out.String()
}
