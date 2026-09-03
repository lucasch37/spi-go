package semantic

import (
	"fmt"

	"github.com/lucasch37/spi-go/internal/errors"
	"github.com/lucasch37/spi-go/internal/ir"
	"github.com/lucasch37/spi-go/internal/tokens"
)

type SemanticAnalyzer struct {
	CurrentScope   *ir.SymbolTable
	ShouldLogScope bool
}

func NewSemanticAnalyzer(shouldLogScope bool) *SemanticAnalyzer {
	sa := &SemanticAnalyzer{
		CurrentScope:   nil,
		ShouldLogScope: shouldLogScope,
	}

	return sa
}

func (sa *SemanticAnalyzer) error(code errors.ErrorCode, token tokens.Token) error {
	return errors.NewSemanticError(code, token, fmt.Sprintf("%s -> %s", code.String(), token.String()))
}

func (sa *SemanticAnalyzer) log(msg string) {
	if sa.ShouldLogScope {
		fmt.Println(msg)
	}
}

func (sa *SemanticAnalyzer) Visit(node ir.Node) error {
	switch node := node.(type) {
	case *ir.Program:
		return sa.visitProgram(node)

	case *ir.ProcedureDecl:
		return sa.visitProcedureDecl(node)

	case *ir.Block:
		for _, decl := range node.Declarations {
			if err := sa.Visit(decl); err != nil {
				return err
			}
		}
		if err := sa.Visit(node.CompoundStatement); err != nil {
			return err
		}

	case *ir.BinOp:
		if err := sa.Visit(node.Left); err != nil {
			return err
		}

		if err := sa.Visit(node.Right); err != nil {
			return err
		}

	case *ir.IntegerLit:
	case *ir.RealLit:
	case *ir.StringLit:

	case *ir.UnaryOp:
		if node.Expr.Type() == ir.StringLitNode {
			return sa.error(errors.InvalidOperand, node.Expr.(*ir.StringLit).Token)
		}

		if err := sa.Visit(node.Expr); err != nil {
			return err
		}

	case *ir.Compound:
		for _, node := range node.Children {
			if err := sa.Visit(node); err != nil {
				return err
			}
		}

	case *ir.NoOp:

	case *ir.VarDecl:
		return sa.visitVarDecl(node)

	case *ir.Assign:
		return sa.visitAssign(node)

	case *ir.Var:
		return sa.visitVar(node)

	case *ir.ProcedureCall:
		return sa.visitProcedureCall(node)

	case *ir.WriteStatement:

	default:
		return fmt.Errorf("no visit method for %T", node)
	}

	return nil
}

func (sa *SemanticAnalyzer) visitProgram(node *ir.Program) error {
	sa.log("ENTER scope: global")

	globalScope := ir.NewSymbolTable("global", 1, sa.CurrentScope, sa.ShouldLogScope)
	globalScope.InitBuiltins()
	sa.CurrentScope = globalScope

	if err := sa.Visit(node.Block); err != nil {
		return err
	}

	sa.log(globalScope.String())

	sa.CurrentScope = sa.CurrentScope.EnclosingScope
	sa.log("LEAVE scope: global")

	return nil
}

func (sa *SemanticAnalyzer) visitProcedureDecl(node *ir.ProcedureDecl) error {
	procName := node.ProcName
	procSymbol := ir.NewProcedureSymbol(procName, make([]*ir.VarSymbol, 0))
	sa.CurrentScope.Insert(procSymbol)

	sa.log(fmt.Sprintf("ENTER scope: %s", procName))
	procScope := ir.NewSymbolTable(procName, sa.CurrentScope.ScopeLevel+1, sa.CurrentScope, sa.ShouldLogScope)
	sa.CurrentScope = procScope

	for _, param := range node.Params {
		paramType := sa.CurrentScope.Lookup(param.TypeNode.Value, false)
		paramName := param.VarNode.Value
		varSymbol := ir.NewVarSymbol(paramName, paramType)

		sa.CurrentScope.Insert(varSymbol)
		procSymbol.Params = append(procSymbol.Params, varSymbol)
	}

	if err := sa.Visit(node.Block); err != nil {
		return err
	}

	sa.log(procScope.String())

	sa.CurrentScope = sa.CurrentScope.EnclosingScope
	sa.log(fmt.Sprintf("LEAVE scope: %s", procName))

	procSymbol.BlockNode = node.Block

	return nil
}

func (sa *SemanticAnalyzer) visitVarDecl(node *ir.VarDecl) error {
	typeName := node.TypeNode.Value
	typeSymbol := sa.CurrentScope.Lookup(typeName, false)
	varName := node.VarNode.Value
	varSymbol := ir.NewVarSymbol(varName, typeSymbol)

	if sa.CurrentScope.Lookup(varName, true) != nil {
		return sa.error(errors.DuplicateID, node.VarNode.Token)
	}

	sa.CurrentScope.Insert(varSymbol)

	return nil
}

func (sa *SemanticAnalyzer) visitAssign(node *ir.Assign) error {
	if err := sa.Visit(node.Right); err != nil {
		return err
	}

	if err := sa.Visit(node.Left); err != nil {
		return err
	}

	return nil
}

func (sa *SemanticAnalyzer) visitVar(node *ir.Var) error {
	varName := node.Value
	varSymbol := sa.CurrentScope.Lookup(varName, false)
	if varSymbol == nil {
		return sa.error(errors.IDNotFound, node.Token)
	}

	return nil
}

func (sa *SemanticAnalyzer) visitProcedureCall(node *ir.ProcedureCall) error {
	actualParamCount := len(node.ActualParams)

	symbol := sa.CurrentScope.Lookup(node.ProcName, false)
	if symbol == nil {
		return sa.error(errors.IDNotFound, node.Token)
	}

	procSymbol, ok := symbol.(*ir.ProcedureSymbol)

	if !ok {
		return sa.error(errors.NotCallable, node.Token)
	}

	paramCount := len(procSymbol.Params)

	if paramCount != actualParamCount {
		return sa.error(errors.WrongParamCount, node.Token)
	}

	for _, param := range node.ActualParams {
		if err := sa.Visit(param); err != nil {
			return err
		}
	}

	node.ProcSymbol = procSymbol

	return nil
}
