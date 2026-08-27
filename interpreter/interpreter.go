package interpreter

import (
	"github.com/lucasch37/spi-go/ast"
	"github.com/lucasch37/spi-go/lexer"
	"github.com/lucasch37/spi-go/parser"
)

type Interpreter struct {
	parser      *parser.Parser
	GLOBALSCOPE map[string]int
}

func NewInterpreter(parser *parser.Parser) *Interpreter {
	return &Interpreter{
		parser:      parser,
		GLOBALSCOPE: make(map[string]int),
	}
}

func (i *Interpreter) visit(node ast.Node) int {
	switch node := node.(type) {

	case ast.BinOp:
		return i.visitBinOp(node)

	case ast.Num:
		return i.visitNum(node)

	case ast.UnaryOp:
		return i.visitUnaryOp(node)

	case ast.Compound:
		return i.visitCompound(node)

	case ast.NoOp:
		return i.visitNoOp(node)

	case ast.Assign:
		return i.visitAssign(node)

	case ast.Var:
		return i.visitVar(node)

	default:
		panic("No visit method for node")
	}
}

func (i *Interpreter) visitBinOp(node ast.BinOp) int {
	switch node.Op.Type {

	case lexer.PLUS:
		return i.visit(node.Left) + i.visit(node.Right)

	case lexer.MINUS:
		return i.visit(node.Left) - i.visit(node.Right)

	case lexer.MUL:
		return i.visit(node.Left) * i.visit(node.Right)

	case lexer.DIV:
		return i.visit(node.Left) / i.visit(node.Right)

	default:
		panic("Unknown binary operator")
	}
}

func (i *Interpreter) visitNum(node ast.Num) int {
	return node.Value
}

func (i *Interpreter) visitUnaryOp(node ast.UnaryOp) int {
	switch node.Op.Type {

	case lexer.PLUS:
		return i.visit(node.Expr)

	case lexer.MINUS:
		return i.visit(node.Expr) * -1

	default:
		panic("Unknown unary operator")
	}
}

func (i *Interpreter) visitCompound(node ast.Compound) int {
	for _, child := range node.Children {
		i.visit(child)
	}
	return 0
}

func (i *Interpreter) visitNoOp(node ast.NoOp) int {
	return 0
}

func (i *Interpreter) visitAssign(node ast.Assign) int {
	varName := node.Left.Value
	i.GLOBALSCOPE[varName] = int(i.visit(node.Right))
	return 0
}

func (i *Interpreter) visitVar(node ast.Var) int {
	varName := node.Value
	if val, exists := i.GLOBALSCOPE[varName]; exists {
		return val
	} else {
		panic("Variable not found: " + varName)
	}
}

func (i *Interpreter) Interpret() int {
	tree := i.parser.Parse()
	return i.visit(tree)
}
