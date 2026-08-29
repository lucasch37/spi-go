package interpreter

import (
	"fmt"

	"github.com/lucasch37/spi-go/ast"
	"github.com/lucasch37/spi-go/lexer"
	"github.com/lucasch37/spi-go/parser"
)

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	REAL_OBJ    ObjectType = "REAL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int
}

func (i Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Real struct {
	Value float64
}

func (r Real) Type() ObjectType {
	return REAL_OBJ
}

func (r Real) Inspect() string {
	return fmt.Sprintf("%g", r.Value)
}

type Interpreter struct {
	parser       *parser.Parser
	GLOBAL_SCOPE map[string]Object
}

func NewInterpreter(p *parser.Parser) *Interpreter {
	return &Interpreter{
		parser:       p,
		GLOBAL_SCOPE: make(map[string]Object),
	}
}

func (i *Interpreter) visit(node ast.Node) (Object, error) {
	switch node := node.(type) {
	case *ast.BinOp:
		return i.visitBinOp(node)

	case *ast.IntegerLit:
		return i.visitIntegerLit(node)

	case *ast.RealLit:
		return i.visitRealLit(node)

	case *ast.UnaryOp:
		return i.visitUnaryOp(node)

	case *ast.Compound:
		return i.visitCompound(node)

	case *ast.NoOp:
		return i.visitNoOp(node)

	case *ast.Assign:
		return i.visitAssign(node)

	case *ast.Var:
		return i.visitVar(node)

	case *ast.Program:
		return i.visitProgram(node)

	case *ast.Block:
		return i.visitBlock(node)

	case *ast.VarDecl:
		return i.visitVarDecl(node)

	case *ast.Type:
		return i.visitType(node)

	default:
		return nil, fmt.Errorf("no visit method for %T", node)
	}
}

func (i *Interpreter) visitBinOp(node *ast.BinOp) (Object, error) {
	left, err := i.visit(node.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.visit(node.Right)
	if err != nil {
		return nil, err
	}

	// If either operand is REAL, perform floating-point arithmetic.
	if left.Type() == REAL_OBJ || right.Type() == REAL_OBJ {
		leftValue := toFloat(left)
		rightValue := toFloat(right)

		switch node.Op.Type {
		case lexer.PLUS:
			return Real{Value: leftValue + rightValue}, nil

		case lexer.MINUS:
			return Real{Value: leftValue - rightValue}, nil

		case lexer.MUL:
			return Real{Value: leftValue * rightValue}, nil

		case lexer.FLOAT_DIV:
			if rightValue == 0 {
				panic("division by zero")
			}
			return Real{Value: leftValue / rightValue}, nil

		default:
			panic("Unknown binary operator")
		}
	}

	// Both operands are INTEGER.
	leftValue := left.(Integer).Value
	rightValue := right.(Integer).Value

	switch node.Op.Type {
	case lexer.PLUS:
		return Integer{Value: leftValue + rightValue}, nil

	case lexer.MINUS:
		return Integer{Value: leftValue - rightValue}, nil

	case lexer.MUL:
		return Integer{Value: leftValue * rightValue}, nil

	case lexer.INTEGER_DIV:
		if rightValue == 0 {
			panic("division by zero")
		}
		return Integer{Value: leftValue / rightValue}, nil

	case lexer.FLOAT_DIV:
		if rightValue == 0 {
			panic("division by zero")
		}

		// Use floating-point division for DIV.
		return Real{Value: float64(leftValue) / float64(rightValue)}, nil

	default:
		panic("Unknown binary operator")
	}
}

func toFloat(obj Object) float64 {
	switch value := obj.(type) {
	case Integer:
		return float64(value.Value)

	case Real:
		return value.Value

	default:
		panic(fmt.Sprintf("Cannot convert %T to float", obj))
	}
}

func (i *Interpreter) visitIntegerLit(node *ast.IntegerLit) (Object, error) {
	return Integer{Value: node.Value}, nil
}

func (i *Interpreter) visitRealLit(node *ast.RealLit) (Object, error) {
	return Real{Value: node.Value}, nil
}

func (i *Interpreter) visitUnaryOp(node *ast.UnaryOp) (Object, error) {
	value, err := i.visit(node.Expr)
	if err != nil {
		return nil, err
	}

	switch node.Op.Type {
	case lexer.PLUS:
		return value, nil

	case lexer.MINUS:
		switch value := value.(type) {
		case Integer:
			return Integer{Value: -value.Value}, nil

		case Real:
			return Real{Value: -value.Value}, nil

		default:
			panic(fmt.Sprintf("Invalid unary operand: %s", value.Type()))
		}

	default:
		return nil, fmt.Errorf("unknown unary operator: %s", node.Op.Type)
	}
}

func (i *Interpreter) visitCompound(node *ast.Compound) (Object, error) {
	for _, child := range node.Children {
		if _, err := i.visit(child); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) visitNoOp(node *ast.NoOp) (Object, error) {
	return nil, nil
}

func (i *Interpreter) visitAssign(node *ast.Assign) (Object, error) {
	varName := node.Left.Value
	value, err := i.visit(node.Right)
	if err != nil {
		return nil, err
	}

	i.GLOBAL_SCOPE[varName] = value

	return nil, nil
}

func (i *Interpreter) visitVar(node *ast.Var) (Object, error) {
	varName := node.Value

	value, exists := i.GLOBAL_SCOPE[varName]
	if !exists {
		panic("Variable not found: " + varName)
	}

	return value, nil
}

func (i *Interpreter) visitProgram(node *ast.Program) (Object, error) {
	return i.visit(node.Block)
}

func (i *Interpreter) visitBlock(node *ast.Block) (Object, error) {
	for _, declaration := range node.Declarations {
		if _, err := i.visit(declaration); err != nil {
			return nil, err
		}
	}

	return i.visit(node.CompoundStatement)
}

func (i *Interpreter) visitVarDecl(node *ast.VarDecl) (Object, error) {
	return nil, nil
}

func (i *Interpreter) visitType(node *ast.Type) (Object, error) {
	return nil, nil
}

func (i *Interpreter) Interpret() error {
	tree, err := i.parser.Parse()
	if err != nil {
		return err
	}

	if tree == nil {
		return nil
	}

	if _, err := i.visit(tree); err != nil {
		return err
	}

	return nil
}
