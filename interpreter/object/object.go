package object

import (
	"bytes"
	"fmt"
	"go-learning/interpreter/parser/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer data
type Integer struct {
	Value int64
}

func (int *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (int *Integer) Inspect() string {
	return fmt.Sprintf("%d", int.Value)
}

// String data
type String struct {
	Value string
}

func (str *String) Type() ObjectType { return STRING_OBJ }

func (str *String) Inspect() string { return str.Value }

// Boolean data
type Boolean struct {
	Value bool
}

func (bool *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (bool *Boolean) Inspect() string {
	return fmt.Sprintf("%t", bool.Value)
}

// Null data
type Null struct{}

func (null *Null) Type() ObjectType {
	return NULL_OBJ
}

func (null *Null) Inspect() string {
	return "null"
}

// ReturnValue data
type ReturnValue struct {
	Value Object
}

func (returnVal *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (returnVal *ReturnValue) Inspect() string  { return returnVal.Value.Inspect() }

// Error data
type Error struct {
	Message string
}

func (error *Error) Type() ObjectType { return ERROR_OBJ }
func (error *Error) Inspect() string  { return "ERROR: " + error.Message }

// Function data
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fun *Function) Type() ObjectType { return FUNCTION_OBJ }
func (fun *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fun.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fun.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (builtin *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (builtin *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (arrObj *Array) Type() ObjectType { return ARRAY_OBJ }
func (arrObj *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range arrObj.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
