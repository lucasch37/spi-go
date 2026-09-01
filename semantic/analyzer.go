package semantic

import (
	"fmt"

	"github.com/lucasch37/spi-go/ast"
	"github.com/lucasch37/spi-go/errors"
	"github.com/lucasch37/spi-go/tokens"
)

type SemanticAnalyzer struct {
	CurrentScope   *SymbolTable
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

func (sa *SemanticAnalyzer) Visit(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		return sa.VisitProgram(node)

	case *ast.ProcedureDecl:
		return sa.VisitProcedureDecl(node)

	case *ast.Block:
		for _, decl := range node.Declarations {
			if err := sa.Visit(decl); err != nil {
				return err
			}
		}
		if err := sa.Visit(node.CompoundStatement); err != nil {
			return err
		}

	case *ast.BinOp:
		if err := sa.Visit(node.Left); err != nil {
			return err
		}

		if err := sa.Visit(node.Right); err != nil {
			return err
		}

	case *ast.IntegerLit:
	case *ast.RealLit:

	case *ast.UnaryOp:
		if err := sa.Visit(node.Expr); err != nil {
			return err
		}

	case *ast.Compound:
		for _, node := range node.Children {
			if err := sa.Visit(node); err != nil {
				return err
			}
		}

	case *ast.NoOp:

	case *ast.VarDecl:
		return sa.VisitVarDecl(node)

	case *ast.Assign:
		return sa.VisitAssign(node)

	case *ast.Var:
		return sa.VisitVar(node)

	default:
		return fmt.Errorf("no visit method for %T", node)
	}

	return nil
}

func (sa *SemanticAnalyzer) VisitProgram(node *ast.Program) error {
	sa.log("ENTER scope: global")

	globalScope := NewSymbolTable("global", 1, sa.CurrentScope, sa.ShouldLogScope)
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

func (sa *SemanticAnalyzer) VisitProcedureDecl(node *ast.ProcedureDecl) error {
	procName := node.ProcName
	procSymbol := NewProcedureSymbol(procName, make([]*VarSymbol, 0))
	sa.CurrentScope.Insert(procSymbol)

	sa.log(fmt.Sprintf("ENTER scope: %s", procName))
	procScope := NewSymbolTable(procName, sa.CurrentScope.ScopeLevel+1, sa.CurrentScope, sa.ShouldLogScope)
	sa.CurrentScope = procScope

	for _, param := range node.Params {
		paramType := sa.CurrentScope.Lookup(param.TypeNode.Value, false)
		paramName := param.VarNode.Value
		varSymbol := NewVarSymbol(paramName, paramType)

		sa.CurrentScope.Insert(varSymbol)
		procSymbol.params = append(procSymbol.params, varSymbol)
	}

	if err := sa.Visit(node.Block); err != nil {
		return err
	}

	sa.log(procScope.String())

	sa.CurrentScope = sa.CurrentScope.EnclosingScope
	sa.log(fmt.Sprintf("LEAVE scope: %s", procName))

	return nil
}

func (sa *SemanticAnalyzer) VisitVarDecl(node *ast.VarDecl) error {
	typeName := node.TypeNode.Value
	typeSymbol := sa.CurrentScope.Lookup(typeName, false)
	varName := node.VarNode.Value
	varSymbol := NewVarSymbol(varName, typeSymbol)

	if sa.CurrentScope.Lookup(varName, true) != nil {
		return sa.error(errors.DuplicateID, node.VarNode.Token)
	}

	sa.CurrentScope.Insert(varSymbol)

	return nil
}

func (sa *SemanticAnalyzer) VisitAssign(node *ast.Assign) error {
	if err := sa.Visit(node.Right); err != nil {
		return err
	}

	if err := sa.Visit(node.Left); err != nil {
		return err
	}

	return nil
}

func (sa *SemanticAnalyzer) VisitVar(node *ast.Var) error {
	varName := node.Value
	varSymbol := sa.CurrentScope.Lookup(varName, false)
	if varSymbol == nil {
		return sa.error(errors.IDNotFound, node.Token)
	}

	return nil
}
