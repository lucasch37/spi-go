package interpreter

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	REAL_OBJ    ObjectType = "REAL"
	STRING_OBJ  ObjectType = "STRING"
)

type Object interface {
	Type() ObjectType
	String() string
}

type IntegerObject struct {
	Value int
}

func (i IntegerObject) Type() ObjectType {
	return INTEGER_OBJ
}

func (i IntegerObject) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type RealObject struct {
	Value float64
}

func (r RealObject) Type() ObjectType {
	return REAL_OBJ
}

func (r RealObject) String() string {
	return fmt.Sprintf("%g", r.Value)
}

type StringObject struct {
	Value string
}

func (s StringObject) Type() ObjectType {
	return STRING_OBJ
}

func (s StringObject) String() string {
	return fmt.Sprintf("%s", s.Value)
}
