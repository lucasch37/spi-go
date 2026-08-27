package ast

import (
	"github.com/lucasch37/spi-go/lexer"
)

type Node interface {
	astNode()
}

type BinOp struct {
	Left  Node
	Token lexer.Token
	Op    lexer.Token
	Right Node
}

type UnaryOp struct {
	Token lexer.Token
	Op    lexer.Token
	Expr  Node
}

type Num struct {
	Token lexer.Token
	Value int
}

type Compound struct {
	Children []Node
}

type Assign struct {
	Left  Var
	Token lexer.Token
	Op    lexer.Token
	Right Node
}

type Var struct {
	Token lexer.Token
	Value string
}

type NoOp struct{}

func (BinOp) astNode()    {}
func (Num) astNode()      {}
func (UnaryOp) astNode()  {}
func (Compound) astNode() {}
func (Assign) astNode()   {}
func (Var) astNode()      {}
func (NoOp) astNode()     {}
