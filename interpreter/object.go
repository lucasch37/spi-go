package interpreter

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	REAL_OBJ    ObjectType = "REAL"
)

type Object interface {
	Type() ObjectType
	String() string
}

type Integer struct {
	Value int
}

func (i Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type Real struct {
	Value float64
}

func (r Real) Type() ObjectType {
	return REAL_OBJ
}

func (r Real) String() string {
	return fmt.Sprintf("%g", r.Value)
}
