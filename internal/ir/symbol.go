package ir

import (
	"fmt"
	"strings"
)

type DataType int

const (
	NoType DataType = iota
	IntegerType
	RealType
	StringType
	BooleanType
)

func (d DataType) Type() DataType {
	return d
}

type Symbol interface {
	Name() string
	String() string
	SetScopeLevel(level int)
}

type BuiltinTypeSymbol struct {
	DataType
	name       string
	ScopeLevel int
}

func NewBuiltinTypeSymbol(dt DataType, name string) *BuiltinTypeSymbol {
	return &BuiltinTypeSymbol{
		DataType: dt,
		name:     name,
	}
}

func (b *BuiltinTypeSymbol) String() string {
	return fmt.Sprintf("<BuiltInTypeSymbol(name='%s')>", b.name)
}

func (b *BuiltinTypeSymbol) Name() string {
	return b.name
}

func (b *BuiltinTypeSymbol) Type() string {
	return b.name
}

func (b *BuiltinTypeSymbol) SetScopeLevel(level int) {
	b.ScopeLevel = level
}

type VarSymbol struct {
	name       string
	TypeSymbol *BuiltinTypeSymbol
	ScopeLevel int
}

func NewVarSymbol(name string, typ *BuiltinTypeSymbol) *VarSymbol {
	return &VarSymbol{
		name:       name,
		TypeSymbol: typ,
	}
}

func (v *VarSymbol) String() string {
	return fmt.Sprintf("<VarSymbol(name='%s', type='%s')>", v.name, v.TypeSymbol.Name())
}

func (v *VarSymbol) Name() string {
	return v.name
}

func (v *VarSymbol) SetScopeLevel(level int) {
	v.ScopeLevel = level
}

type ProcedureSymbol struct {
	name       string
	Params     []*VarSymbol
	BlockNode  *Block
	ScopeLevel int
}

func NewProcedureSymbol(name string, params []*VarSymbol) *ProcedureSymbol {
	return &ProcedureSymbol{
		name:   name,
		Params: params,
	}
}

func (p *ProcedureSymbol) String() string {
	return fmt.Sprintf("<ProcedureSymbol(name='%s', params='%v')>", p.name, p.Params)
}

func (p *ProcedureSymbol) Name() string {
	return p.name
}

func (p *ProcedureSymbol) SetScopeLevel(level int) {
	p.ScopeLevel = level
}

type FunctionSymbol struct {
	name           string
	Params         []*VarSymbol
	BlockNode      *Block
	ReturnTypeNode *TypeN
	ScopeLevel     int
}

func NewFunctionSymbol(name string, params []*VarSymbol) *FunctionSymbol {
	return &FunctionSymbol{
		name:   name,
		Params: params,
	}
}

func (f *FunctionSymbol) String() string {
	return fmt.Sprintf("<FunctionSymbol(name='%s', params='%v')>", f.name, f.Params)
}

func (f *FunctionSymbol) Name() string {
	return f.name
}

func (f *FunctionSymbol) SetScopeLevel(level int) {
	f.ScopeLevel = level
}

func (f *FunctionSymbol) ReturnType() DataType {
	switch f.ReturnTypeNode.Value {
	case "INTEGER":
		return IntegerType
	case "REAL":
		return RealType
	case "STRING":
		return StringType
	case "BOOLEAN":
		return BooleanType
	default:
		return NoType
	}
}

type SymbolTable struct {
	Symbols        map[string]Symbol
	ScopeName      string
	ScopeLevel     int
	EnclosingScope *SymbolTable
	ShouldLogScope bool
}

func NewSymbolTable(scopeName string, scopeLevel int, enclosingScope *SymbolTable, shouldLogScope bool) *SymbolTable {
	st := &SymbolTable{
		Symbols:        make(map[string]Symbol),
		ScopeName:      scopeName,
		ScopeLevel:     scopeLevel,
		EnclosingScope: enclosingScope,
		ShouldLogScope: shouldLogScope,
	}

	return st
}

func (st *SymbolTable) InitBuiltins() {
	st.Insert(NewBuiltinTypeSymbol(IntegerType, "INTEGER"))
	st.Insert(NewBuiltinTypeSymbol(RealType, "REAL"))
	st.Insert(NewBuiltinTypeSymbol(StringType, "STRING"))
	st.Insert(NewBuiltinTypeSymbol(BooleanType, "BOOLEAN"))
}

func (st *SymbolTable) String() string {
	h1 := "SYMBOL TABLE"

	lines := []string{
		"",
		h1,
		strings.Repeat("=", 20),
	}

	lines = append(lines,
		fmt.Sprintf("%-15s: %s", "Scope name", st.ScopeName),
		fmt.Sprintf("%-15s: %d", "Scope level", st.ScopeLevel),
	)

	if st.EnclosingScope != nil {
		lines = append(lines, fmt.Sprintf("%-15s: %v", "Enclosing scope", st.EnclosingScope.ScopeName))
	}

	h2 := "\nSymbol Table Contents"
	lines = append(lines,
		h2,
		strings.Repeat("-", len(h2)),
	)

	for key, value := range st.Symbols {
		lines = append(lines,
			fmt.Sprintf("%7s: %v", key, value.String()),
		)
	}

	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

func (st *SymbolTable) log(msg string) {
	if st.ShouldLogScope {
		fmt.Println(msg)
	}
}

func (st *SymbolTable) Insert(symbol Symbol) {
	st.log(fmt.Sprintf("Insert: %s", symbol.Name()))

	symbol.SetScopeLevel(st.ScopeLevel)

	st.Symbols[symbol.Name()] = symbol
}

func (st *SymbolTable) Lookup(name string, currentScopeOnly bool) Symbol {
	st.log(fmt.Sprintf("Lookup: %s (scope: %s)", name, st.ScopeName))
	symbol, exists := st.Symbols[name]
	if exists {
		return symbol
	}

	if currentScopeOnly {
		return nil
	}

	if st.EnclosingScope != nil {
		return st.EnclosingScope.Lookup(name, false)
	}

	return nil
}
