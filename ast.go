package main

type AST interface {
	astNode()
}

type BinOp struct {
	Left  AST
	Token Token
	Op    Token
	Right AST
}

type UnaryOp struct {
	Token Token
	Op    Token
	Expr  AST
}

type Num struct {
	Token Token
	Value int
}

type Compound struct {
	Children []AST
}

type Assign struct {
	Left  Var
	Token Token
	Op    Token
	Right AST
}

type Var struct {
	Token Token
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
