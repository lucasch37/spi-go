package ir

import (
	"github.com/lucasch37/spi-go/internal/tokens"
)

type Node interface {
	String() string
	Type() NodeType
	SourceToken() tokens.Token
}

type NodeType int

const (
	BinOpNode NodeType = iota
	UnaryOpNode
	IntegerLitNode
	RealLitNode
	StringLitNode
	BooleanLitNode
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
	WriteStatementNode
	IfStatementNode
)

var nodeTypeNames = map[NodeType]string{
	BinOpNode:          "BinOp",
	UnaryOpNode:        "UnaryOp",
	IntegerLitNode:     "IntegerLit",
	BooleanLitNode:     "BooleanLit",
	RealLitNode:        "RealLit",
	StringLitNode:      "StringLit",
	CompoundNode:       "Compound",
	AssignNode:         "Assign",
	VarNode:            "Var",
	NoOpNode:           "NoOp",
	ProgramNode:        "Program",
	BlockNode:          "Block",
	VarDeclNode:        "VarDecl",
	TypeNode:           "Type",
	ProcedureDeclNode:  "ProcedureDecl",
	ParamNode:          "Param",
	ProcedureCallNode:  "ProcedureCall",
	WriteStatementNode: "WriteStatement",
	IfStatementNode:    "IfStatement",
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

func (n *BinOp) SourceToken() tokens.Token {
	return n.Token
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

func (n *UnaryOp) SourceToken() tokens.Token {
	return n.Token
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

func (n *Compound) SourceToken() tokens.Token {
	return tokens.Token{}
}

func NewCompound(children []Node) *Compound {
	return &Compound{
		NodeType: NodeType(CompoundNode),
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

func (n *Assign) SourceToken() tokens.Token {
	return n.Token
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

func (n *Var) SourceToken() tokens.Token {
	return n.Token
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

func (n *NoOp) SourceToken() tokens.Token {
	return tokens.Token{}
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

func (n *Program) SourceToken() tokens.Token {
	return tokens.Token{}
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

func (n *Block) SourceToken() tokens.Token {
	return tokens.Token{}
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
	TypeNode *TypeN
}

func (n *VarDecl) SourceToken() tokens.Token {
	return n.VarNode.SourceToken()
}

func NewVarDecl(varNode *Var, typeNode *TypeN) *VarDecl {
	return &VarDecl{
		NodeType: NodeType(VarDeclNode),
		VarNode:  varNode,
		TypeNode: typeNode,
	}
}

type TypeN struct {
	NodeType
	Token tokens.Token
	Value string
}

func (n *TypeN) SourceToken() tokens.Token {
	return n.Token
}

func NewType(token tokens.Token) *TypeN {
	return &TypeN{
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

func (n *IntegerLit) SourceToken() tokens.Token {
	return n.Token
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

func (n *RealLit) SourceToken() tokens.Token {
	return n.Token
}

func NewRealLit(token tokens.Token) *RealLit {
	return &RealLit{
		NodeType: RealLitNode,
		Token:    token,
		Value:    token.Value.(float64),
	}
}

type StringLit struct {
	NodeType
	Token tokens.Token
	Value string
}

func (n *StringLit) SourceToken() tokens.Token {
	return n.Token
}

func NewStringLit(token tokens.Token) *StringLit {
	return &StringLit{
		NodeType: StringLitNode,
		Token:    token,
		Value:    token.Value.(string),
	}
}

type BooleanLit struct {
	NodeType
	Token tokens.Token
	Value bool
}

func (n *BooleanLit) SourceToken() tokens.Token {
	return n.Token
}

func NewBooleanLit(token tokens.Token) *BooleanLit {
	return &BooleanLit{
		NodeType: BooleanLitNode,
		Token:    token,
		Value:    token.Value.(bool),
	}
}

type ProcedureDecl struct {
	NodeType
	ProcName string
	Block    *Block
	Params   []*Param
}

func (n *ProcedureDecl) SourceToken() tokens.Token {
	return tokens.Token{}
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
	TypeNode *TypeN
}

func (n *Param) SourceToken() tokens.Token {
	return n.VarNode.SourceToken()
}

func NewParam(varNode *Var, typeNode *TypeN) *Param {
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

func (n *ProcedureCall) SourceToken() tokens.Token {
	return n.Token
}

func NewProcedureCall(procName string, actualParams []Node, token tokens.Token) *ProcedureCall {
	return &ProcedureCall{
		NodeType:     ProcedureCallNode,
		ProcName:     procName,
		ActualParams: actualParams,
		Token:        token,
	}
}

type WriteStatement struct {
	NodeType
	NewLine bool
	Exprs   []Node
	Token   tokens.Token
}

func (n *WriteStatement) SourceToken() tokens.Token {
	return n.Token
}

func NewWriteStatement(newLine bool, exprs []Node, token tokens.Token) *WriteStatement {
	return &WriteStatement{
		NodeType: WriteStatementNode,
		NewLine:  newLine,
		Exprs:    exprs,
		Token:    token,
	}
}

type IfStatement struct {
	NodeType
	Condition   Node
	Statement   Node
	Alternative Node
	Token       tokens.Token
}

func (n *IfStatement) SourceToken() tokens.Token {
	return n.Token
}

func NewIfStatement(condition Node, statement Node, alternative Node, token tokens.Token) *IfStatement {
	return &IfStatement{
		Condition:   condition,
		NodeType:    IfStatementNode,
		Statement:   statement,
		Alternative: alternative,
		Token:       token,
	}
}
