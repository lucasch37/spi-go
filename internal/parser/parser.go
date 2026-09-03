package parser

import (
	"fmt"

	"github.com/lucasch37/spi-go/internal/errors"
	"github.com/lucasch37/spi-go/internal/ir"
	"github.com/lucasch37/spi-go/internal/lexer"
	"github.com/lucasch37/spi-go/internal/tokens"
)

/*
program : PROGRAM variable SEMI block DOT

block : declarations compound_statement

declarations : (VAR (variableDeclaration SEMI)+)? procedure_declaration*

variableDeclaration : ID (COMMA ID)* COLON type_spec

procedureDeclaration :
     PROCEDURE ID (LPAREN formalParameterList RPAREN)? SEMI block SEMI

formalParamsList : formalParameters
                   | formalParameters SEMI formalParameterList

formalParameters : ID (COMMA ID)* COLON type_spec

typeSpec : INTEGER | REAL | STRING

compoundStatement : BEGIN statementList END

statementList : statement
               | statement SEMI statementList

statement : compoundStatement
          | proccallStatement
          | assignmentStatement
          | writeStatement
          | empty

writeStatement : (WRITE | WRITELN) LPAREN (expr (COMMA expr)*)? RPAREN

procCallStatement : ID LPAREN (expr (COMMA expr)*)? RPAREN

assignmentStatement : variable ASSIGN expr

empty :

expr : term ((PLUS | MINUS) term)*

term : factor ((MUL | INTEGER_DIV | FLOAT_DIV) factor)*

factor : PLUS factor
       | MINUS factor
       | INTEGER_LIT
       | REAL_LIT
			 | STRING_LIT
       | LPAREN expr RPAREN
       | variable

variable: ID
*/

type Parser struct {
	lexer        *lexer.Lexer
	currentToken tokens.Token
}

func NewParser(lexer *lexer.Lexer) (*Parser, error) {
	token, err := lexer.GetNextToken()
	if err != nil {
		return nil, err
	}

	return &Parser{
		lexer:        lexer,
		currentToken: token,
	}, nil
}

func (p *Parser) Parse() (ir.Node, error) {
	tree, err := p.program()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != tokens.EOF {
		return nil, p.error(errors.UnexpectedToken)
	}

	return tree, nil
}

func (p *Parser) error(code errors.ErrorCode) error {
	return errors.NewSyntaxError(code, p.currentToken, fmt.Sprintf("%s -> %s", code.String(), p.currentToken.String()))
}

func (p *Parser) eat(tokenType tokens.TokenType) error {
	if p.currentToken.Type != tokenType {
		return p.error(errors.UnexpectedToken)
	}

	token, err := p.lexer.GetNextToken()
	if err != nil {
		return err
	}

	p.currentToken = token
	return nil
}

func (p *Parser) program() (ir.Node, error) {
	if err := p.eat(tokens.PROGRAM); err != nil {
		return nil, err
	}

	varNode, err := p.variable()
	if err != nil {
		return nil, err
	}

	progName := varNode.Value

	if err := p.eat(tokens.SEMI); err != nil {
		return nil, err
	}

	blockNode, err := p.block()
	if err != nil {
		return nil, err
	}

	programNode := ir.NewProgram(progName, blockNode)

	if err := p.eat(tokens.DOT); err != nil {
		return nil, err
	}

	return programNode, nil
}

func (p *Parser) block() (*ir.Block, error) {
	declarationNodes, err := p.declarations()
	if err != nil {
		return nil, err
	}

	compoundStatementNode, err := p.compoundStatement()
	if err != nil {
		return nil, err
	}

	node := ir.NewBlock(declarationNodes, compoundStatementNode)
	return node, nil
}

func (p *Parser) declarations() ([]ir.Node, error) {
	var declarations []ir.Node
	for p.currentToken.Type == tokens.VAR {
		if err := p.eat(tokens.VAR); err != nil {
			return nil, err
		}

		for p.currentToken.Type == tokens.ID {
			varDecls, err := p.variableDeclaration()
			if err != nil {
				return nil, err
			}

			for _, vd := range varDecls {
				declarations = append(declarations, vd)
			}

			if err := p.eat(tokens.SEMI); err != nil {
				return nil, err
			}
		}
	}

	for p.currentToken.Type == tokens.PROCEDURE {
		procDecl, err := p.ProcedureDeclaraton()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, procDecl)
	}

	return declarations, nil
}

func (p *Parser) ProcedureDeclaraton() (*ir.ProcedureDecl, error) {
	if err := p.eat(tokens.PROCEDURE); err != nil {
		return nil, err
	}

	procName := p.currentToken.Value.(string)

	if err := p.eat(tokens.ID); err != nil {
		return nil, err
	}

	var params []*ir.Param

	if p.currentToken.Type == tokens.LPAREN {
		if err := p.eat(tokens.LPAREN); err != nil {
			return nil, err
		}

		paramNodes, err := p.formalParameterList()
		if err != nil {
			return nil, err
		}

		params = append(params, paramNodes...)

		if err := p.eat(tokens.RPAREN); err != nil {
			return nil, err
		}
	}

	if err := p.eat(tokens.SEMI); err != nil {
		return nil, err
	}

	blockNode, err := p.block()
	if err != nil {
		return nil, err
	}

	procDecl := ir.NewProcedureDecl(procName, blockNode, params)

	if err := p.eat(tokens.SEMI); err != nil {
		return nil, err
	}

	return procDecl, nil
}

func (p *Parser) formalParamaters() ([]*ir.Param, error) {
	var paramNodes []*ir.Param

	paramTokens := []tokens.Token{p.currentToken}

	if err := p.eat(tokens.ID); err != nil {
		return nil, err
	}

	for p.currentToken.Type == tokens.COMMA {
		if err := p.eat(tokens.COMMA); err != nil {
			return nil, err
		}

		paramTokens = append(paramTokens, p.currentToken)

		if err := p.eat(tokens.ID); err != nil {
			return nil, err
		}
	}

	if err := p.eat(tokens.COLON); err != nil {
		return nil, err
	}

	typeNode, err := p.typeSpec()
	if err != nil {
		return nil, err
	}

	for _, paramToken := range paramTokens {
		paramNode := ir.NewParam(ir.NewVar(paramToken), typeNode)
		paramNodes = append(paramNodes, paramNode)
	}

	return paramNodes, nil
}

func (p *Parser) formalParameterList() ([]*ir.Param, error) {
	if p.currentToken.Type != tokens.ID {
		return make([]*ir.Param, 0), nil
	}

	paramNodes, err := p.formalParamaters()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == tokens.SEMI {
		if err := p.eat(tokens.SEMI); err != nil {
			return nil, err
		}

		newParamNodes, err := p.formalParamaters()
		if err != nil {
			return nil, err
		}

		paramNodes = append(paramNodes, newParamNodes...)
	}

	return paramNodes, nil
}

func (p *Parser) variableDeclaration() ([]*ir.VarDecl, error) {
	varNodes := []*ir.Var{{
		Token: p.currentToken,
		Value: p.currentToken.Value.(string),
	}}
	if err := p.eat(tokens.ID); err != nil {
		return nil, err
	}

	for p.currentToken.Type == tokens.COMMA {
		if err := p.eat(tokens.COMMA); err != nil {
			return nil, err
		}

		varNodes = append(varNodes, ir.NewVar(p.currentToken))
		if err := p.eat(tokens.ID); err != nil {
			return nil, err
		}
	}

	if err := p.eat(tokens.COLON); err != nil {
		return nil, err
	}

	typeNode, err := p.typeSpec()
	if err != nil {
		return nil, err
	}

	var varDeclarations []*ir.VarDecl
	for _, varNode := range varNodes {
		varDeclarations = append(varDeclarations, ir.NewVarDecl(varNode, typeNode))
	}

	return varDeclarations, nil
}

func (p *Parser) typeSpec() (*ir.Type, error) {
	token := p.currentToken

	switch p.currentToken.Type {
	case tokens.INTEGER:
		if err := p.eat(tokens.INTEGER); err != nil {
			return nil, err
		}
	case tokens.REAL:
		if err := p.eat(tokens.REAL); err != nil {
			return nil, err
		}
	case tokens.STRING:
		if err := p.eat(tokens.STRING); err != nil {
			return nil, err
		}
	default:
		return nil, p.error(errors.UnexpectedToken)
	}

	node := ir.NewType(token)

	return node, nil
}

func (p *Parser) compoundStatement() (*ir.Compound, error) {
	if err := p.eat(tokens.BEGIN); err != nil {
		return nil, err
	}

	nodes, err := p.statementList()
	if err != nil {
		return nil, err
	}

	if err := p.eat(tokens.END); err != nil {
		return nil, err
	}

	root := ir.NewCompound(nodes)
	return root, nil
}

func (p *Parser) statementList() ([]ir.Node, error) {
	node, err := p.statement()
	if err != nil {
		return nil, err
	}

	results := []ir.Node{
		node,
	}

	for p.currentToken.Type == tokens.SEMI {
		if err := p.eat(tokens.SEMI); err != nil {
			return nil, err
		}

		node, err := p.statement()
		if err != nil {
			return nil, err
		}
		results = append(results, node)
	}

	if p.currentToken.Type == tokens.ID {
		p.error(errors.ExpectedAssign)
	}

	return results, nil
}

func (p *Parser) statement() (ir.Node, error) {
	switch p.currentToken.Type {
	case tokens.BEGIN:
		return p.compoundStatement()
	case tokens.ID:
		if p.lexer.CurrentChar == '(' {
			return p.ProcCallStatement()
		} else {
			return p.assignmentStatement()
		}
	case tokens.WRITELN:
		return p.writeStatement()

	case tokens.WRITE:
		return p.writeStatement()

	default:
		return p.empty()
	}
}

func (p *Parser) writeStatement() (*ir.WriteStatement, error) {
	token := p.currentToken

	if token.Type == tokens.WRITE {
		if err := p.eat(tokens.WRITE); err != nil {
			return nil, err
		}
	} else {
		if err := p.eat(tokens.WRITELN); err != nil {
			return nil, err
		}
	}

	if err := p.eat(tokens.LPAREN); err != nil {
		return nil, err
	}

	var expressions []ir.Node

	if p.currentToken.Type != tokens.RPAREN {
		node, err := p.expr()
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, node)
	}

	for p.currentToken.Type == tokens.COMMA {
		if err := p.eat(tokens.COMMA); err != nil {
			return nil, err
		}

		node, err := p.expr()
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, node)
	}

	if err := p.eat(tokens.RPAREN); err != nil {
		return nil, err
	}

	if token.Type == tokens.WRITE {
		return ir.NewWriteStatement(false, expressions, token), nil
	} else {
		return ir.NewWriteStatement(true, expressions, token), nil
	}
}

func (p *Parser) ProcCallStatement() (*ir.ProcedureCall, error) {
	token := p.currentToken
	procName := token.Value.(string)

	if err := p.eat(tokens.ID); err != nil {
		return nil, err
	}

	if err := p.eat(tokens.LPAREN); err != nil {
		return nil, err
	}

	var actualParams []ir.Node

	if p.currentToken.Type != tokens.RPAREN {
		node, err := p.expr()
		if err != nil {
			return nil, err
		}

		actualParams = append(actualParams, node)
	}

	for p.currentToken.Type == tokens.COMMA {
		if err := p.eat(tokens.COMMA); err != nil {
			return nil, err
		}

		node, err := p.expr()
		if err != nil {
			return nil, err
		}

		actualParams = append(actualParams, node)
	}

	if err := p.eat(tokens.RPAREN); err != nil {
		return nil, err
	}

	return ir.NewProcedureCall(procName, actualParams, token), nil
}

func (p *Parser) assignmentStatement() (*ir.Assign, error) {
	left, err := p.variable()
	if err != nil {
		return nil, err
	}

	token := p.currentToken

	if err := p.eat(tokens.ASSIGN); err != nil {
		return nil, err
	}

	right, err := p.expr()
	if err != nil {
		return nil, err
	}

	node := ir.NewAssign(left, token, right)

	return node, nil
}

func (p *Parser) variable() (*ir.Var, error) {
	node := ir.NewVar(p.currentToken)

	if err := p.eat(tokens.ID); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) empty() (*ir.NoOp, error) {
	return ir.NewNoOp(), nil
}

func (p *Parser) factor() (ir.Node, error) {
	token := p.currentToken

	switch token.Type {
	case tokens.PLUS:
		if err := p.eat(tokens.PLUS); err != nil {
			return nil, err
		}

		expr, err := p.factor()
		if err != nil {
			return nil, err
		}

		node := ir.NewUnaryOp(token, expr)

		return node, nil

	case tokens.MINUS:
		if err := p.eat(tokens.MINUS); err != nil {
			return nil, err
		}

		expr, err := p.factor()
		fmt.Println(expr.Type().String())
		if err != nil {
			return nil, err
		}

		node := ir.NewUnaryOp(token, expr)
		return node, nil

	case tokens.INTEGER_LIT:
		if err := p.eat(tokens.INTEGER_LIT); err != nil {
			return nil, err
		}

		return ir.NewIntegerLit(token), nil

	case tokens.REAL_LIT:
		if err := p.eat(tokens.REAL_LIT); err != nil {
			return nil, err
		}

		return ir.NewRealLit(token), nil

	case tokens.LPAREN:
		if err := p.eat(tokens.LPAREN); err != nil {
			return nil, err
		}

		node, err := p.expr()
		if err != nil {
			return nil, err
		}

		if err := p.eat(tokens.RPAREN); err != nil {
			return nil, err
		}
		return node, nil

	case tokens.STRING_LIT:
		if err := p.eat(tokens.STRING_LIT); err != nil {
			return nil, err
		}

		return ir.NewStringLit(token), nil

	default:
		return p.variable()
	}
}

func (p *Parser) term() (ir.Node, error) {
	node, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == tokens.MUL ||
		p.currentToken.Type == tokens.INTEGER_DIV || p.currentToken.Type == tokens.FLOAT_DIV {

		token := p.currentToken

		switch token.Type {
		case tokens.MUL:
			if err := p.eat(tokens.MUL); err != nil {
				return nil, err
			}
		case tokens.INTEGER_DIV:
			if err := p.eat(tokens.INTEGER_DIV); err != nil {
				return nil, err
			}
		case tokens.FLOAT_DIV:
			if err := p.eat(tokens.FLOAT_DIV); err != nil {
				return nil, err
			}
		}

		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		node = ir.NewBinOp(node, token, right)
	}

	return node, nil
}

func (p *Parser) expr() (ir.Node, error) {
	node, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == tokens.PLUS ||
		p.currentToken.Type == tokens.MINUS {

		token := p.currentToken

		if token.Type == tokens.PLUS {
			if err := p.eat(tokens.PLUS); err != nil {
				return nil, err
			}
		} else {
			if err := p.eat(tokens.MINUS); err != nil {
				return nil, err
			}
		}

		right, err := p.term()
		if err != nil {
			return nil, err
		}

		node = ir.NewBinOp(node, token, right)
	}

	return node, nil
}
