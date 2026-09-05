package semantic

import (
	"fmt"

	"github.com/lucasch37/nsspi/internal/errors"
	"github.com/lucasch37/nsspi/internal/ir"
	"github.com/lucasch37/nsspi/internal/tokens"
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

func (sa *SemanticAnalyzer) Analyze(tree ir.Node) error {
	if _, err := sa.visit(tree); err != nil {
		return err
	}

	return nil
}

func (sa *SemanticAnalyzer) error(code errors.ErrorCode, token tokens.Token) error {
	return errors.NewSemanticError(code, token, fmt.Sprintf("%s -> %s", code.String(), token.String()))
}

func (sa *SemanticAnalyzer) log(msg string) {
	if sa.ShouldLogScope {
		fmt.Println(msg)
	}
}

func (sa *SemanticAnalyzer) visit(node ir.Node) (ir.DataType, error) {
	switch node := node.(type) {
	case *ir.Program:
		return sa.visitProgram(node)

	case *ir.ProcedureDecl:
		return sa.visitProcedureDecl(node)

	case *ir.FunctionDecl:
		return sa.visitFunctionDecl(node)

	case *ir.Block:
		for _, decl := range node.Declarations {
			if _, err := sa.visit(decl); err != nil {
				return ir.NoType, err
			}
		}
		if _, err := sa.visit(node.CompoundStatement); err != nil {
			return ir.NoType, err
		}

	case *ir.BinOp:
		return sa.visitBinOp(node)

	case *ir.IntegerLit:
		return ir.IntegerType, nil

	case *ir.RealLit:
		return ir.RealType, nil

	case *ir.StringLit:
		return ir.StringType, nil

	case *ir.BooleanLit:
		return ir.BooleanType, nil

	case *ir.UnaryOp:
		exprType, err := sa.visit(node.Expr)
		if err != nil {
			return ir.NoType, err
		}

		if exprType != ir.IntegerType && exprType != ir.RealType {
			return ir.NoType, sa.error(errors.InvalidOperand, node.Token)
		}

	case *ir.Compound:
		for _, node := range node.Children {
			if _, err := sa.visit(node); err != nil {
				return ir.NoType, err
			}
		}

	case *ir.NoOp:

	case *ir.VarDecl:
		return sa.visitVarDecl(node)

	case *ir.Assign:
		return sa.visitAssign(node)

	case *ir.Identifier:
		return sa.visitIdentifier(node)

	case *ir.Call:
		return sa.visitCall(node)

	case *ir.WriteStatement:
		for _, node := range node.Exprs {
			if _, err := sa.visit(node); err != nil {
				return ir.NoType, err
			}
		}

	case *ir.IfStatement:
		if _, err := sa.visit(node.Condition); err != nil {
			return ir.NoType, err
		}

		if _, err := sa.visit(node.Statement); err != nil {
			return ir.NoType, err
		}

		if node.Alternative != nil {
			if _, err := sa.visit(node.Alternative); err != nil {
				return ir.NoType, err
			}
		}

	default:
		return ir.NoType, fmt.Errorf("no visit method for %T", node)
	}

	return ir.NoType, nil
}

func isNumeric(t ir.DataType) bool {
	return t == ir.IntegerType || t == ir.RealType
}

func (sa *SemanticAnalyzer) visitBinOp(node *ir.BinOp) (ir.DataType, error) {
	leftType, err := sa.visit(node.Left)
	if err != nil {
		return ir.NoType, err
	}

	rightType, err := sa.visit(node.Right)
	if err != nil {
		return ir.NoType, err
	}

	switch node.Op.Type {
	case tokens.MINUS, tokens.MUL, tokens.FLOAT_DIV:
		if !isNumeric(leftType) || !isNumeric(rightType) {
			return ir.NoType, sa.error(errors.TypeMismatch, node.Token)
		}

		if leftType == ir.RealType || rightType == ir.RealType {
			return ir.RealType, nil
		}

		return ir.IntegerType, nil

	case tokens.PLUS:
		if leftType == ir.BooleanType || rightType == ir.BooleanType {
			return ir.NoType, sa.error(errors.TypeMismatch, node.Token)
		}

		if leftType == ir.StringType || rightType == ir.StringType {
			if leftType != rightType {
				return ir.NoType, sa.error(errors.TypeMismatch, node.Token)
			}

			return ir.StringType, nil
		}

		if leftType == ir.RealType || rightType == ir.RealType {
			return ir.RealType, nil
		}

		return ir.IntegerType, nil

	case tokens.MOD:
		if leftType != ir.IntegerType || rightType != ir.IntegerType {
			return ir.NoType, sa.error(errors.TypeMismatch, node.Token)
		}

		return ir.IntegerType, nil

	case tokens.INTEGER_DIV:
		if leftType != ir.IntegerType || rightType != ir.IntegerType {
			return ir.NoType, sa.error(errors.TypeMismatch, node.Token)
		}

		return ir.IntegerType, nil

	case tokens.EQUAL, tokens.NOT_EQUAL, tokens.LESS_THAN, tokens.LESS_THAN_EQUAL, tokens.GREATER_THAN, tokens.GREATER_THAN_EQUAL:
		if leftType != rightType {
			return ir.NoType, sa.error(errors.TypeMismatch, node.Token)
		}

		return ir.BooleanType, nil

	default:
		return ir.NoType, sa.error(errors.InvalidOperand, node.Token)
	}
}

func (sa *SemanticAnalyzer) visitProgram(node *ir.Program) (ir.DataType, error) {
	sa.log("ENTER scope: global")

	globalScope := ir.NewSymbolTable("global", 1, sa.CurrentScope, sa.ShouldLogScope)
	globalScope.InitBuiltins()
	sa.CurrentScope = globalScope

	if _, err := sa.visit(node.Block); err != nil {
		return ir.NoType, err
	}

	sa.log(globalScope.String())

	sa.CurrentScope = sa.CurrentScope.EnclosingScope
	sa.log("LEAVE scope: global")

	return ir.NoType, nil
}

func (sa *SemanticAnalyzer) visitProcedureDecl(node *ir.ProcedureDecl) (ir.DataType, error) {
	procName := node.ProcName
	procSymbol := ir.NewProcedureSymbol(procName, make([]*ir.VarSymbol, 0))
	sa.CurrentScope.Insert(procSymbol)

	sa.log(fmt.Sprintf("ENTER scope: %s", procName))
	procScope := ir.NewSymbolTable(procName, sa.CurrentScope.ScopeLevel+1, sa.CurrentScope, sa.ShouldLogScope)
	sa.CurrentScope = procScope

	for _, param := range node.Params {
		symbol := sa.CurrentScope.Lookup(param.TypeNode.Value, false)

		paramType, ok := symbol.(*ir.BuiltinTypeSymbol)
		if !ok {
			return ir.NoType, sa.error(errors.InvalidType, param.TypeNode.Token)
		}

		paramName := param.IdNode.Value
		varSymbol := ir.NewVarSymbol(paramName, paramType)

		sa.CurrentScope.Insert(varSymbol)
		procSymbol.Params = append(procSymbol.Params, varSymbol)
	}

	if _, err := sa.visit(node.Block); err != nil {
		return ir.NoType, err
	}

	sa.log(procScope.String())

	sa.CurrentScope = sa.CurrentScope.EnclosingScope
	sa.log(fmt.Sprintf("LEAVE scope: %s", procName))

	procSymbol.BlockNode = node.Block

	return ir.NoType, nil
}

func (sa *SemanticAnalyzer) visitFunctionDecl(node *ir.FunctionDecl) (ir.DataType, error) {
	funcName := node.FuncName
	funcSymbol := ir.NewFunctionSymbol(funcName, make([]*ir.VarSymbol, 0))
	sa.CurrentScope.Insert(funcSymbol)

	funcSymbol.ReturnTypeNode = node.ReturnType

	sa.log(fmt.Sprintf("ENTER scope: %s", funcName))
	funcScope := ir.NewSymbolTable(funcName, sa.CurrentScope.ScopeLevel+1, sa.CurrentScope, sa.ShouldLogScope)
	sa.CurrentScope = funcScope

	for _, param := range node.Params {
		symbol := sa.CurrentScope.Lookup(param.TypeNode.Value, false)

		paramType, ok := symbol.(*ir.BuiltinTypeSymbol)
		if !ok {
			return ir.NoType, sa.error(errors.InvalidType, param.TypeNode.Token)
		}

		paramName := param.IdNode.Value
		varSymbol := ir.NewVarSymbol(paramName, paramType)

		sa.CurrentScope.Insert(varSymbol)
		funcSymbol.Params = append(funcSymbol.Params, varSymbol)
	}

	if _, err := sa.visit(node.Block); err != nil {
		return ir.NoType, err
	}

	sa.log(funcScope.String())

	sa.CurrentScope = sa.CurrentScope.EnclosingScope
	sa.log(fmt.Sprintf("LEAVE scope: %s", funcName))

	funcSymbol.BlockNode = node.Block

	return ir.NoType, nil
}

func (sa *SemanticAnalyzer) visitVarDecl(node *ir.VarDecl) (ir.DataType, error) {
	typeName := node.TypeNode.Value
	symbol := sa.CurrentScope.Lookup(typeName, false)

	typeSymbol, ok := symbol.(*ir.BuiltinTypeSymbol)
	if !ok {
		return ir.NoType, sa.error(errors.InvalidType, node.TypeNode.Token)
	}

	varName := node.IdNode.Value
	varSymbol := ir.NewVarSymbol(varName, typeSymbol)

	if sa.CurrentScope.Lookup(varName, true) != nil {
		return ir.NoType, sa.error(errors.DuplicateID, node.IdNode.Token)
	}

	sa.CurrentScope.Insert(varSymbol)

	node.IdNode.Symbol = varSymbol

	return ir.NoType, nil
}

func (sa *SemanticAnalyzer) visitAssign(node *ir.Assign) (ir.DataType, error) {
	leftType, err := sa.visit(node.Right)
	if err != nil {
		return ir.NoType, err
	}

	rightType, err := sa.visit(node.Left)
	if err != nil {
		return ir.NoType, err
	}

	if leftType != rightType {
		return ir.NoType, sa.error(errors.TypeMismatch, node.Right.SourceToken())
	}

	return ir.NoType, nil
}

func (sa *SemanticAnalyzer) visitIdentifier(node *ir.Identifier) (ir.DataType, error) {
	varName := node.Value
	symbol := sa.CurrentScope.Lookup(varName, false)
	if symbol == nil {
		return ir.NoType, sa.error(errors.IDNotFound, node.Token)
	}

	switch symbol := symbol.(type) {
	case *ir.VarSymbol:
		node.Symbol = symbol
		return symbol.TypeSymbol.DataType, nil

	case *ir.FunctionSymbol:
		node.Symbol = symbol
		return symbol.ReturnType(), nil

	}
	return ir.NoType, sa.error(errors.IDNotFound, node.Token)
}

func (sa *SemanticAnalyzer) visitCall(node *ir.Call) (ir.DataType, error) {
	symbol := sa.CurrentScope.Lookup(node.CallName, false)
	if symbol == nil {
		return ir.NoType, sa.error(errors.IDNotFound, node.Token)
	}

	var params []*ir.VarSymbol

	switch symbol := symbol.(type) {
	case *ir.ProcedureSymbol:
		params = symbol.Params
		node.Symbol = symbol

	case *ir.FunctionSymbol:
		params = symbol.Params
		node.Symbol = symbol

	default:
		return ir.NoType, nil
	}

	if len(params) != len(node.ActualParams) {
		return ir.NoType, sa.error(errors.WrongParamCount, node.Token)
	}

	for i, param := range node.ActualParams {
		paramType, err := sa.visit(param)
		if err != nil {
			return ir.NoType, err
		}

		if paramType != params[i].TypeSymbol.DataType {
			return ir.NoType, sa.error(errors.TypeMismatch, param.SourceToken())
		}
	}

	funcSymbol, ok := symbol.(*ir.FunctionSymbol)
	if ok {
		return funcSymbol.ReturnType(), nil
	}

	return ir.NoType, nil
}
