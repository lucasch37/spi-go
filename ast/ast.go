package ast

import (
	"github.com/lucasch37/spi-go/lexer"
)

type Node interface {
	String() string
	Type() NodeType
}

type NodeType int

const (
	BinOpNode NodeType = iota
	UnaryOpNode
	IntegerLitNode
	RealLitNode
	CompoundNode
	AssignNode
	VarNode
	NoOpNode
	ProgramNode
	BlockNode
	VarDeclNode
	TypeNode
)

var nodeTypeNames = map[NodeType]string{
	BinOpNode:      "BinOp",
	UnaryOpNode:    "UnaryOp",
	IntegerLitNode: "IntegerLit",
	RealLitNode:    "RealLit",
	CompoundNode:   "Compound",
	AssignNode:     "Assign",
	VarNode:        "Var",
	NoOpNode:       "NoOp",
	ProgramNode:    "Program",
	BlockNode:      "Block",
	VarDeclNode:    "VarDecl",
	TypeNode:       "Type",
}

func (nt NodeType) Type() NodeType {
	return nt
}

func (nt NodeType) String() string {
	if name, ok := nodeTypeNames[nt]; ok {
		return name
	}

	return "Unknown"
}

type BinOp struct {
	NodeType
	Left  Node
	Token lexer.Token
	Op    lexer.Token
	Right Node
}

func NewBinOp(left Node, token lexer.Token, right Node) *BinOp {
	return &BinOp{
		NodeType: BinOpNode,
		Left:     left,
		Token:    token,
		Op:       token,
		Right:    right,
	}
}

type UnaryOp struct {
	NodeType
	Token lexer.Token
	Op    lexer.Token
	Expr  Node
}

func NewUnaryOp(token lexer.Token, expr Node) *UnaryOp {
	return &UnaryOp{
		NodeType: UnaryOpNode,
		Token:    token,
		Op:       token,
		Expr:     expr,
	}
}

type Compound struct {
	NodeType
	Children []Node
}

func NewCompound(children []Node) *Compound {
	return &Compound{
		NodeType: CompoundNode,
		Children: children,
	}
}

type Assign struct {
	NodeType
	Left  *Var
	Token lexer.Token
	Op    lexer.Token
	Right Node
}

func NewAssign(left *Var, token lexer.Token, right Node) *Assign {
	return &Assign{
		NodeType: AssignNode,
		Left:     left,
		Token:    token,
		Op:       token,
		Right:    right,
	}
}

type Var struct {
	NodeType
	Token lexer.Token
	Value string
}

func NewVar(token lexer.Token) *Var {
	return &Var{
		NodeType: VarNode,
		Token:    token,
		Value:    token.Value.(string),
	}
}

type NoOp struct {
	NodeType
}

func NewNoOp() *NoOp {
	return &NoOp{
		NodeType: NoOpNode,
	}
}

type Program struct {
	NodeType
	Name  string
	Block *Block
}

func NewProgram(name string, block *Block) *Program {
	return &Program{
		NodeType: ProgramNode,
		Name:     name,
		Block:    block,
	}
}

type Block struct {
	NodeType
	Declarations      []Node
	CompoundStatement *Compound
}

func NewBlock(declarations []Node, compoundStatement *Compound) *Block {
	return &Block{
		NodeType:          BlockNode,
		Declarations:      declarations,
		CompoundStatement: compoundStatement,
	}
}

type VarDecl struct {
	NodeType
	VarNode  *Var
	TypeNode *Type
}

func NewVarDecl(varNode *Var, typeNode *Type) *VarDecl {
	return &VarDecl{
		NodeType: NodeType(VarDeclNode),
		VarNode:  varNode,
		TypeNode: typeNode,
	}
}

type Type struct {
	NodeType
	Token lexer.Token
	Value string
}

func NewType(token lexer.Token) *Type {
	return &Type{
		NodeType: TypeNode,
		Token:    token,
		Value:    token.Value.(string),
	}
}

type IntegerLit struct {
	NodeType
	Token lexer.Token
	Value int
}

func NewIntegerLit(token lexer.Token) *IntegerLit {
	return &IntegerLit{
		NodeType: IntegerLitNode,
		Token:    token,
		Value:    token.Value.(int),
	}
}

type RealLit struct {
	NodeType
	Token lexer.Token
	Value float64
}

func NewRealLit(token lexer.Token) *RealLit {
	return &RealLit{
		NodeType: RealLitNode,
		Token:    token,
		Value:    token.Value.(float64),
	}
}
