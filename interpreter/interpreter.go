package interpreter

import (
	"fmt"

	"github.com/lucasch37/spi-go/ast"
	"github.com/lucasch37/spi-go/errors"
	"github.com/lucasch37/spi-go/tokens"
)

type Interpreter struct {
	Tree           ast.Node
	CallStack      *CallStack
	ShouldLogStack bool
}

func NewInterpreter(tree ast.Node, shouldLogStack bool) *Interpreter {
	return &Interpreter{
		Tree:           tree,
		CallStack:      NewCallStack(),
		ShouldLogStack: shouldLogStack,
	}
}

func (i *Interpreter) Interpret() error {
	if i.Tree == nil {
		return nil
	}

	if _, err := i.visit(i.Tree); err != nil {
		return err
	}

	return nil
}

func (i *Interpreter) log(msg string) {
	if i.ShouldLogStack {
		fmt.Println(msg)
	}
}

func (i *Interpreter) error(code errors.ErrorCode, token tokens.Token) error {
	return errors.NewRuntimeError(code, token, fmt.Sprintf("%s -> %s", code.String(), token.String()))
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

	case *ast.ProcedureDecl:
		return i.visitProcedureDecl(node)

	case *ast.ProcedureCall:
		return i.visitProcedureCall(node)

	default:
		return nil, fmt.Errorf("no visit method for %T", node)
	}
}

func (i *Interpreter) visitProgram(node *ast.Program) (Object, error) {
	programName := node.Name
	i.log(fmt.Sprintf("\nENTER: PROGRAM %s", programName))

	ar := NewActivationRecord(programName, PROGRAM, 1)
	i.CallStack.Push(ar)
	i.log(i.CallStack.String())

	i.visit(node.Block)

	i.log(fmt.Sprintf("LEAVE: PROGRAM %s", programName))
	i.log(i.CallStack.String())
	i.CallStack.Pop()

	return nil, nil
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

	if left.Type() == REAL_OBJ || right.Type() == REAL_OBJ {
		leftValue, err := toFloat(left)
		if err != nil {
			return nil, err
		}

		rightValue, err := toFloat(right)
		if err != nil {
			return nil, err
		}

		switch node.Op.Type {
		case tokens.PLUS:
			return Real{Value: leftValue + rightValue}, nil

		case tokens.MINUS:
			return Real{Value: leftValue - rightValue}, nil

		case tokens.MUL:
			return Real{Value: leftValue * rightValue}, nil

		case tokens.FLOAT_DIV:
			if rightValue == 0 {
				return nil, i.error(errors.DivideByZero, node.Token)
			}
			return Real{Value: leftValue / rightValue}, nil

		default:
			return nil, fmt.Errorf("Unknown binary operator")
		}
	}

	leftValue := left.(Integer).Value
	rightValue := right.(Integer).Value

	switch node.Op.Type {
	case tokens.PLUS:
		return Integer{Value: leftValue + rightValue}, nil

	case tokens.MINUS:
		return Integer{Value: leftValue - rightValue}, nil

	case tokens.MUL:
		return Integer{Value: leftValue * rightValue}, nil

	case tokens.INTEGER_DIV:
		if rightValue == 0 {
			return nil, i.error(errors.DivideByZero, node.Token)
		}
		return Integer{Value: leftValue / rightValue}, nil

	case tokens.FLOAT_DIV:
		if rightValue == 0 {
			return nil, i.error(errors.DivideByZero, node.Token)
		}

		return Real{Value: float64(leftValue) / float64(rightValue)}, nil

	default:
		return nil, fmt.Errorf("Unknown binary operator: %s", node.Op.Type.String())
	}
}

func toFloat(obj Object) (float64, error) {
	switch value := obj.(type) {
	case Integer:
		return float64(value.Value), nil

	case Real:
		return value.Value, nil

	default:
		return 0, fmt.Errorf("Cannot convert %T to float", obj)
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
	case tokens.PLUS:
		return value, nil

	case tokens.MINUS:
		switch value := value.(type) {
		case Integer:
			return Integer{Value: -value.Value}, nil

		case Real:
			return Real{Value: -value.Value}, nil

		default:
			return nil, fmt.Errorf("Invalid unary operand: %s", value.Type())
		}

	default:
		return nil, fmt.Errorf("unknown unary operator: %s", node.Op.Type.String())
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

	ar := i.CallStack.Peek()
	ar.Set(varName, value)

	return nil, nil
}

func (i *Interpreter) visitVar(node *ast.Var) (Object, error) {
	varName := node.Value

	ar := i.CallStack.Peek()
	value := ar.Get(varName)

	return value, nil
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

func (i *Interpreter) visitProcedureDecl(node *ast.ProcedureDecl) (Object, error) {
	return nil, nil
}

func (i *Interpreter) visitProcedureCall(node *ast.ProcedureCall) (Object, error) {
	return nil, nil
}
