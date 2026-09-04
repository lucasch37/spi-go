package interpreter

import (
	"fmt"

	"github.com/lucasch37/spi-go/internal/errors"
	"github.com/lucasch37/spi-go/internal/ir"
	"github.com/lucasch37/spi-go/internal/tokens"
)

type Interpreter struct {
	Tree           ir.Node
	CallStack      *CallStack
	ShouldLogStack bool
}

func NewInterpreter(tree ir.Node, shouldLogStack bool) *Interpreter {
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

func (i *Interpreter) visit(node ir.Node) (Object, error) {
	switch node := node.(type) {
	case *ir.BinOp:
		return i.visitBinOp(node)

	case *ir.IntegerLit:
		return i.visitIntegerLit(node)

	case *ir.RealLit:
		return i.visitRealLit(node)

	case *ir.StringLit:
		return i.visitStringLit(node)

	case *ir.BooleanLit:
		return i.visitBooleanLit(node)

	case *ir.UnaryOp:
		return i.visitUnaryOp(node)

	case *ir.Compound:
		return i.visitCompound(node)

	case *ir.NoOp:
		return i.visitNoOp(node)

	case *ir.Assign:
		return i.visitAssign(node)

	case *ir.Var:
		return i.visitVar(node)

	case *ir.Program:
		return i.visitProgram(node)

	case *ir.Block:
		return i.visitBlock(node)

	case *ir.VarDecl:
		return i.visitVarDecl(node)

	case *ir.TypeN:
		return i.visitType(node)

	case *ir.ProcedureDecl:
		return i.visitProcedureDecl(node)

	case *ir.ProcedureCall:
		return i.visitProcedureCall(node)

	case *ir.WriteStatement:
		return i.visitWriteStatement(node)

	case *ir.IfStatement:
		return i.visitIfStatement(node)

	default:
		return nil, fmt.Errorf("no visit method for %T", node)
	}
}

func (i *Interpreter) visitProgram(node *ir.Program) (Object, error) {
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

func toFloat(obj Object) (float64, error) {
	switch value := obj.(type) {
	case IntegerObject:
		return float64(value.Value), nil

	case RealObject:
		return value.Value, nil

	default:
		return 0, fmt.Errorf("Cannot convert %T to float", obj)
	}
}

func compareObject(op tokens.TokenType, left, right Object) (BooleanObject, bool) {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		leftValue := left.(IntegerObject).Value
		rightValue := right.(IntegerObject).Value

		switch op {
		case tokens.EQUAL:
			return BooleanObject{Value: leftValue == rightValue}, true
		case tokens.NOT_EQUAL:
			return BooleanObject{Value: leftValue != rightValue}, true
		case tokens.LESS_THAN:
			return BooleanObject{Value: leftValue < rightValue}, true
		case tokens.LESS_THAN_EQUAL:
			return BooleanObject{Value: leftValue <= rightValue}, true
		case tokens.GREATER_THAN:
			return BooleanObject{Value: leftValue > rightValue}, true
		case tokens.GREATER_THAN_EQUAL:
			return BooleanObject{Value: leftValue >= rightValue}, true
		}

	case left.Type() == REAL_OBJ || right.Type() == REAL_OBJ:
		leftValue, err := toFloat(left)
		if err != nil {
			return BooleanObject{}, false
		}

		rightValue, err := toFloat(right)
		if err != nil {
			return BooleanObject{}, false
		}

		switch op {
		case tokens.EQUAL:
			return BooleanObject{Value: leftValue == rightValue}, true
		case tokens.NOT_EQUAL:
			return BooleanObject{Value: leftValue != rightValue}, true
		case tokens.LESS_THAN:
			return BooleanObject{Value: leftValue < rightValue}, true
		case tokens.LESS_THAN_EQUAL:
			return BooleanObject{Value: leftValue <= rightValue}, true
		case tokens.GREATER_THAN:
			return BooleanObject{Value: leftValue > rightValue}, true
		case tokens.GREATER_THAN_EQUAL:
			return BooleanObject{Value: leftValue >= rightValue}, true
		}

	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		leftValue := left.(StringObject).Value
		rightValue := right.(StringObject).Value

		switch op {
		case tokens.EQUAL:
			return BooleanObject{Value: leftValue == rightValue}, true
		case tokens.NOT_EQUAL:
			return BooleanObject{Value: leftValue != rightValue}, true
		case tokens.LESS_THAN:
			return BooleanObject{Value: leftValue < rightValue}, true
		case tokens.LESS_THAN_EQUAL:
			return BooleanObject{Value: leftValue <= rightValue}, true
		case tokens.GREATER_THAN:
			return BooleanObject{Value: leftValue > rightValue}, true
		case tokens.GREATER_THAN_EQUAL:
			return BooleanObject{Value: leftValue >= rightValue}, true
		}
	}

	return BooleanObject{}, false
}

func (i *Interpreter) visitBinOp(node *ir.BinOp) (Object, error) {
	left, err := i.visit(node.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.visit(node.Right)
	if err != nil {
		return nil, err
	}

	if result, ok := compareObject(node.Op.Type, left, right); ok {
		return result, nil
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
			return RealObject{Value: leftValue + rightValue}, nil

		case tokens.MINUS:
			return RealObject{Value: leftValue - rightValue}, nil

		case tokens.MUL:
			return RealObject{Value: leftValue * rightValue}, nil

		case tokens.FLOAT_DIV:
			if rightValue == 0 {
				return nil, i.error(errors.DivideByZero, node.Token)
			}

			return RealObject{Value: leftValue / rightValue}, nil
		}
	}

	if left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ {
		leftValue := left.(IntegerObject).Value
		rightValue := right.(IntegerObject).Value

		switch node.Op.Type {
		case tokens.PLUS:
			return IntegerObject{Value: leftValue + rightValue}, nil

		case tokens.MINUS:
			return IntegerObject{Value: leftValue - rightValue}, nil

		case tokens.MUL:
			return IntegerObject{Value: leftValue * rightValue}, nil

		case tokens.INTEGER_DIV:
			if rightValue == 0 {
				return nil, i.error(errors.DivideByZero, node.Token)
			}

			return IntegerObject{Value: leftValue / rightValue}, nil

		case tokens.FLOAT_DIV:
			if rightValue == 0 {
				return nil, i.error(errors.DivideByZero, node.Token)
			}

			return RealObject{
				Value: float64(leftValue) / float64(rightValue),
			}, nil
		}
	}

	if left.Type() == STRING_OBJ && right.Type() == STRING_OBJ {
		leftValue := left.(StringObject).Value
		rightValue := right.(StringObject).Value

		if node.Op.Type == tokens.PLUS {
			return StringObject{Value: leftValue + rightValue}, nil
		}
	}

	return nil, fmt.Errorf(
		"Invalid binary operation: %s %s %s",
		left.Type(),
		node.Op.Type.String(),
		right.Type(),
	)
}

func (i *Interpreter) visitIntegerLit(node *ir.IntegerLit) (Object, error) {
	return IntegerObject{Value: node.Value}, nil
}

func (i *Interpreter) visitRealLit(node *ir.RealLit) (Object, error) {
	return RealObject{Value: node.Value}, nil
}

func (i *Interpreter) visitStringLit(node *ir.StringLit) (Object, error) {
	return StringObject{Value: node.Value}, nil
}

func (i *Interpreter) visitBooleanLit(node *ir.BooleanLit) (Object, error) {
	return BooleanObject{Value: node.Value}, nil
}

func (i *Interpreter) visitUnaryOp(node *ir.UnaryOp) (Object, error) {
	value, err := i.visit(node.Expr)
	if err != nil {
		return nil, err
	}

	switch node.Op.Type {
	case tokens.PLUS:
		return value, nil

	case tokens.MINUS:
		switch value := value.(type) {
		case IntegerObject:
			return IntegerObject{Value: -value.Value}, nil

		case RealObject:
			return RealObject{Value: -value.Value}, nil

		default:
			return nil, fmt.Errorf("Invalid unary operand: %s", value.Type())
		}

	default:
		return nil, fmt.Errorf("unknown unary operator: %s", node.Op.Type.String())
	}
}

func (i *Interpreter) visitCompound(node *ir.Compound) (Object, error) {
	for _, child := range node.Children {
		if _, err := i.visit(child); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) visitNoOp(node *ir.NoOp) (Object, error) {
	return nil, nil
}

func (i *Interpreter) visitAssign(node *ir.Assign) (Object, error) {
	varName := node.Left.Value
	value, err := i.visit(node.Right)
	if err != nil {
		return nil, err
	}

	ar := i.CallStack.Peek()
	ar.Set(varName, value)

	return nil, nil
}

func (i *Interpreter) visitVar(node *ir.Var) (Object, error) {
	varName := node.Value

	value := i.CallStack.Lookup(varName)

	return value, nil
}

func (i *Interpreter) visitBlock(node *ir.Block) (Object, error) {
	for _, declaration := range node.Declarations {
		if _, err := i.visit(declaration); err != nil {
			return nil, err
		}
	}

	return i.visit(node.CompoundStatement)
}

func (i *Interpreter) visitVarDecl(node *ir.VarDecl) (Object, error) {
	return nil, nil
}

func (i *Interpreter) visitType(node *ir.TypeN) (Object, error) {
	return nil, nil
}

func (i *Interpreter) visitProcedureDecl(node *ir.ProcedureDecl) (Object, error) {
	return nil, nil
}

func (i *Interpreter) visitProcedureCall(node *ir.ProcedureCall) (Object, error) {
	procName := node.ProcName

	ar := NewActivationRecord(procName, PROCEDURE, node.ProcSymbol.ScopeLevel+1)

	formalParams := node.ProcSymbol.Params
	actualParams := node.ActualParams

	for index, paramSymbol := range formalParams {
		paramValue, err := i.visit(actualParams[index])
		if err != nil {
			return nil, err
		}

		ar.Set(paramSymbol.Name(), paramValue)
	}

	i.CallStack.Push(ar)

	i.log(fmt.Sprintf("ENTER: PROCEDURE %s", procName))
	i.log(i.CallStack.String())

	if _, err := i.visit(node.ProcSymbol.BlockNode); err != nil {
		return nil, err
	}

	i.log(fmt.Sprintf("LEAVE: PROCEDURE %s", procName))
	i.log(i.CallStack.String())

	i.CallStack.Pop()

	return nil, nil
}

func (i *Interpreter) visitWriteStatement(node *ir.WriteStatement) (Object, error) {
	if len(node.Exprs) == 0 {
		if node.NewLine {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}

		return nil, nil
	}

	for _, expr := range node.Exprs {
		val, err := i.visit(expr)
		if err != nil {
			return nil, err
		}

		fmt.Print(val.String())
	}

	if node.NewLine {
		fmt.Println()
	}

	return nil, nil
}

func (i *Interpreter) visitIfStatement(node *ir.IfStatement) (Object, error) {
	condition, err := i.visit(node.Condition)
	if err != nil {
		return nil, err
	}

	boolCondition, ok := condition.(BooleanObject)
	if !ok {
		return nil, fmt.Errorf("Condition must be a boolean expression, got %T", condition)
	}

	if boolCondition.Value {
		i.visit(node.Statement)
	} else if node.Alternative != nil {
		i.visit(node.Alternative)
	}

	return nil, nil
}
