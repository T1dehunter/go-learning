package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
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
