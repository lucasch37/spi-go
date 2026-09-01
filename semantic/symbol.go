package semantic

import (
	"fmt"
	"strings"
)

type Symbol interface {
	Name() string
	Type() string
	String() string
}

type BuiltinTypeSymbol struct {
	name string
}

func NewBuiltinTypeSymbol(name string) *BuiltinTypeSymbol {
	return &BuiltinTypeSymbol{
		name: name,
	}
}

func (b *BuiltinTypeSymbol) String() string {
	return fmt.Sprintf("<BuiltInTypeSymbol(name='%s')>", b.name)
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

func (v *VarSymbol) String() string {
	return fmt.Sprintf("<VarSymbol(name='%s', type='%s')>", v.name, v.typ.Name())
}

func (v *VarSymbol) Name() string {
	return v.name
}

func (v *VarSymbol) Type() string {
	return v.typ.Name()
}

type ProcedureSymbol struct {
	name   string
	params []*VarSymbol
}

func NewProcedureSymbol(name string, params []*VarSymbol) *ProcedureSymbol {
	return &ProcedureSymbol{
		name:   name,
		params: params,
	}
}

func (p *ProcedureSymbol) String() string {
	return fmt.Sprintf("<ProcedureSymbol(name='%s', params='%v')>", p.name, p.params)
}

func (p *ProcedureSymbol) Name() string {
	return p.name
}

func (p *ProcedureSymbol) Type() string {
	return ""
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
	st.Insert(NewBuiltinTypeSymbol("INTEGER"))
	st.Insert(NewBuiltinTypeSymbol("REAL"))
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
