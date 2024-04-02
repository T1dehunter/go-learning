package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
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
