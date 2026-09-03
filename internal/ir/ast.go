package ir

import (
	"github.com/lucasch37/spi-go/internal/tokens"
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
	ProcedureDeclNode
	ParamNode
	ProcedureCallNode
)

var nodeTypeNames = map[NodeType]string{
	BinOpNode:         "BinOp",
	UnaryOpNode:       "UnaryOp",
	IntegerLitNode:    "IntegerLit",
	RealLitNode:       "RealLit",
	CompoundNode:      "Compound",
	AssignNode:        "Assign",
	VarNode:           "Var",
	NoOpNode:          "NoOp",
	ProgramNode:       "Program",
	BlockNode:         "Block",
	VarDeclNode:       "VarDecl",
	TypeNode:          "Type",
	ProcedureDeclNode: "ProcedureDecl",
	ParamNode:         "Param",
	ProcedureCallNode: "ProcedureCall",
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
	Token tokens.Token
	Op    tokens.Token
	Right Node
}

func NewBinOp(left Node, token tokens.Token, right Node) *BinOp {
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
	Token tokens.Token
	Op    tokens.Token
	Expr  Node
}

func NewUnaryOp(token tokens.Token, expr Node) *UnaryOp {
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
	Token tokens.Token
	Op    tokens.Token
	Right Node
}

func NewAssign(left *Var, token tokens.Token, right Node) *Assign {
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
	Token tokens.Token
	Value string
}

func NewVar(token tokens.Token) *Var {
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
	Token tokens.Token
	Value string
}

func NewType(token tokens.Token) *Type {
	return &Type{
		NodeType: TypeNode,
		Token:    token,
		Value:    token.Value.(string),
	}
}

type IntegerLit struct {
	NodeType
	Token tokens.Token
	Value int
}

func NewIntegerLit(token tokens.Token) *IntegerLit {
	return &IntegerLit{
		NodeType: IntegerLitNode,
		Token:    token,
		Value:    token.Value.(int),
	}
}

type RealLit struct {
	NodeType
	Token tokens.Token
	Value float64
}

func NewRealLit(token tokens.Token) *RealLit {
	return &RealLit{
		NodeType: RealLitNode,
		Token:    token,
		Value:    token.Value.(float64),
	}
}

type ProcedureDecl struct {
	NodeType
	ProcName string
	Block    *Block
	Params   []*Param
}

func NewProcedureDecl(procName string, block *Block, params []*Param) *ProcedureDecl {
	return &ProcedureDecl{
		NodeType: ProcedureDeclNode,
		ProcName: procName,
		Block:    block,
		Params:   params,
	}
}

type Param struct {
	NodeType
	VarNode  *Var
	TypeNode *Type
}

func NewParam(varNode *Var, typeNode *Type) *Param {
	return &Param{
		NodeType: ParamNode,
		VarNode:  varNode,
		TypeNode: typeNode,
	}
}

type ProcedureCall struct {
	NodeType
	ProcName     string
	ActualParams []Node
	Token        tokens.Token
	ProcSymbol   *ProcedureSymbol
}

func NewProcedureCall(procName string, actualParams []Node, token tokens.Token) *ProcedureCall {
	return &ProcedureCall{
		NodeType:     ProcedureCallNode,
		ProcName:     procName,
		ActualParams: actualParams,
		Token:        token,
	}
}
