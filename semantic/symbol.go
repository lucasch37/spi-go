package semantic

import (
	"fmt"

	"github.com/lucasch37/spi-go/ast"
)

type Symbol interface {
	Name() string
	Type() string
}

type BuiltinTypeSymbol struct {
	name string
}

func NewBuiltinTypeSymbol(name string) *BuiltinTypeSymbol {
	return &BuiltinTypeSymbol{
		name: name,
	}
}

func (b *BuiltinTypeSymbol) Name() string {
	return b.name
}

func (b *BuiltinTypeSymbol) Type() string {
	return ""
}

type VarSymbol struct {
	name string
	typ  Symbol
}

func NewVarSymbol(name string, typ Symbol) *VarSymbol {
	return &VarSymbol{
		name: name,
		typ:  typ,
	}
}

func (v *VarSymbol) Name() string {
	return v.name
}

func (v *VarSymbol) Type() string {
	return v.typ.Name()
}

type SymbolTable struct {
	Symbols map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	s := &SymbolTable{
		Symbols: make(map[string]Symbol),
	}

	s.Symbols["INTEGER"] = NewBuiltinTypeSymbol("INTEGER")
	s.Symbols["REAL"] = NewBuiltinTypeSymbol("REAL")

	return s
}

func (st *SymbolTable) String() string {
	return fmt.Sprintf("Symbols: %v", st.Symbols)
}

func (st *SymbolTable) Define(symbol Symbol) {
	// fmt.Printf("Define: %s\n", symbol.Name())
	st.Symbols[symbol.Name()] = symbol
}

func (st *SymbolTable) Lookup(name string) Symbol {
	// fmt.Printf("Lookup: %s\n", name)
	symbol, exists := st.Symbols[name]
	if !exists {
		return nil
	}

	return symbol
}

type SymbolTableBuilder struct {
	SymbolTable *SymbolTable
}

func NewSymbolTableBuilder() *SymbolTableBuilder {
	st := NewSymbolTable()
	stb := &SymbolTableBuilder{
		SymbolTable: st,
	}

	return stb
}

func (stb *SymbolTableBuilder) Visit(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		if err := stb.Visit(node.Block); err != nil {
			return err
		}

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

	case *ast.ProcedureDecl:

	default:
		return fmt.Errorf("no visit method for %T", node)
	}

	return nil
}

func (stb *SymbolTableBuilder) VisitVarDecl(node *ast.VarDecl) error {
	typeName := node.TypeNode.Value
	typeSymbol := stb.SymbolTable.Lookup(typeName)
	varName := node.VarNode.Value
	varSymbol := NewVarSymbol(varName, typeSymbol)
	stb.SymbolTable.Define(varSymbol)

	return nil
}

func (stb *SymbolTableBuilder) VisitAssign(node *ast.Assign) error {
	varName := node.Left.Value
	varSymbol := stb.SymbolTable.Lookup(varName)
	if varSymbol == nil {
		return fmt.Errorf("NameError: %q", varName)
	}

	return stb.Visit(node.Right)
}

func (stb *SymbolTableBuilder) VisitVar(node *ast.Var) error {
	varName := node.Value
	varSymbol := stb.SymbolTable.Lookup(varName)
	if varSymbol == nil {
		return fmt.Errorf("NameError: %q", varName)
	}

	return nil
}
