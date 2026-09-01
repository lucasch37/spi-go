package semantic

import (
	"fmt"

	"github.com/lucasch37/spi-go/ast"
)

type SemanticAnalyzer struct {
	CurrentScope *SymbolTable
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	stb := &SemanticAnalyzer{
		CurrentScope: nil,
	}

	return stb
}

func (stb *SemanticAnalyzer) Visit(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		return stb.VisitProgram(node)

	case *ast.ProcedureDecl:
		return stb.VisitProcedureDecl(node)

	case *ast.Block:
		for _, decl := range node.Declarations {
			if err := stb.Visit(decl); err != nil {
				return err
			}
		}
		if err := stb.Visit(node.CompoundStatement); err != nil {
			return err
		}

	case *ast.BinOp:
		if err := stb.Visit(node.Left); err != nil {
			return err
		}

		if err := stb.Visit(node.Right); err != nil {
			return err
		}

	case *ast.IntegerLit:
	case *ast.RealLit:

	case *ast.UnaryOp:
		if err := stb.Visit(node.Expr); err != nil {
			return err
		}

	case *ast.Compound:
		for _, node := range node.Children {
			if err := stb.Visit(node); err != nil {
				return err
			}
		}

	case *ast.NoOp:

	case *ast.VarDecl:
		return stb.VisitVarDecl(node)

	case *ast.Assign:
		return stb.VisitAssign(node)

	case *ast.Var:
		return stb.VisitVar(node)

	default:
		return fmt.Errorf("no visit method for %T", node)
	}

	return nil
}

func (stb *SemanticAnalyzer) VisitProgram(node *ast.Program) error {
	fmt.Println("ENTER scope: global")

	globalScope := NewSymbolTable("global", 1, stb.CurrentScope)
	globalScope.InitBuiltins()
	stb.CurrentScope = globalScope

	if err := stb.Visit(node.Block); err != nil {
		return err
	}

	fmt.Println(globalScope.String())

	stb.CurrentScope = stb.CurrentScope.EnclosingScope
	fmt.Println("LEAVE scope: global")

	return nil
}

func (stb *SemanticAnalyzer) VisitProcedureDecl(node *ast.ProcedureDecl) error {
	procName := node.ProcName
	procSymbol := NewProcedureSymbol(procName, make([]*VarSymbol, 0))
	stb.CurrentScope.Insert(procSymbol)

	fmt.Println("ENTER scope: ", procName)
	procScope := NewSymbolTable(procName, stb.CurrentScope.ScopeLevel+1, stb.CurrentScope)
	stb.CurrentScope = procScope

	for _, param := range node.Params {
		paramType := stb.CurrentScope.Lookup(param.TypeNode.Value, false)
		paramName := param.VarNode.Value
		varSymbol := NewVarSymbol(paramName, paramType)

		stb.CurrentScope.Insert(varSymbol)
		procSymbol.params = append(procSymbol.params, varSymbol)
	}

	if err := stb.Visit(node.Block); err != nil {
		return err
	}

	fmt.Println(procScope.String())

	stb.CurrentScope = stb.CurrentScope.EnclosingScope
	fmt.Println("LEAVE scope: ", procName)

	return nil
}

func (stb *SemanticAnalyzer) VisitVarDecl(node *ast.VarDecl) error {
	typeName := node.TypeNode.Value
	typeSymbol := stb.CurrentScope.Lookup(typeName, false)
	varName := node.VarNode.Value
	varSymbol := NewVarSymbol(varName, typeSymbol)

	if stb.CurrentScope.Lookup(varName, true) != nil {
		return fmt.Errorf("Duplicate identifier: %q", varName)
	}

	stb.CurrentScope.Insert(varSymbol)

	return nil
}

func (stb *SemanticAnalyzer) VisitAssign(node *ast.Assign) error {
	if err := stb.Visit(node.Right); err != nil {
		return err
	}

	if err := stb.Visit(node.Left); err != nil {
		return err
	}

	return nil
}

func (stb *SemanticAnalyzer) VisitVar(node *ast.Var) error {
	varName := node.Value
	varSymbol := stb.CurrentScope.Lookup(varName, false)
	if varSymbol == nil {
		return fmt.Errorf("NameError: %q", varName)
	}

	return nil
}
