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

// PrefixExpression data
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

// PrefixExpression data
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (infixExpr *InfixExpression) expressionNode() {

}

func (infixExpr *InfixExpression) TokenLiteral() string {
	return infixExpr.Token.Literal
}

func (infixExpr *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(infixExpr.Left.String())
	out.WriteString(" " + infixExpr.Operator + " ")
	out.WriteString(infixExpr.Right.String())
	out.WriteString(")")

	return out.String()
}

// Boolean data
type Boolean struct {
	Token token.Token
	Value bool
}

func (boolean *Boolean) expressionNode() {}

func (boolean *Boolean) TokenLiteral() string {
	return boolean.Token.Literal
}

func (boolean *Boolean) String() string {
	return boolean.Token.Literal
}

// IfExpression data
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpr *IfExpression) expressionNode()      {}
func (ifExpr *IfExpression) TokenLiteral() string { return ifExpr.Token.Literal }
func (ifExpr *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ifExpr.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifExpr.Consequence.String())

	if ifExpr.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ifExpr.Alternative.String())
	}

	return out.String()
}

// BlockStatement data
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (blockStmt *BlockStatement) statementNode()       {}
func (blockStmt *BlockStatement) TokenLiteral() string { return blockStmt.Token.Literal }
func (blockStmt *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range blockStmt.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
