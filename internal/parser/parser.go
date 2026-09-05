package parser

import (
	"fmt"

	"github.com/lucasch37/nsspi/internal/errors"
	"github.com/lucasch37/nsspi/internal/ir"
	"github.com/lucasch37/nsspi/internal/lexer"
	"github.com/lucasch37/nsspi/internal/tokens"
)

/*
program : PROGRAM variable SEMI block DOT

block : declarationList compound_statement

declarationList
    : declaration declarationList
    | empty

declaration
    : VAR variableDeclaration SEMI
    | procedureDeclaration
    | functionDeclaration
    | empty

variableDeclaration : ID (COMMA ID)* COLON type_spec

procedureDeclaration :
     PROCEDURE ID (LPAREN formalParameterList RPAREN)? SEMI block SEMI

functionDeclaration : FUNCTION ID (LPAREN formalParameterList RPAREN)? COLON type_spec SEMI block SEMI

formalParameterList : formalParameters
                    | formalParameters SEMI formalParameterList

formalParameters : ID (COMMA ID)* COLON type_spec

typeSpec : INTEGER | REAL | STRING | BOOLEAN

compoundStatement : BEGIN statementList END

statementList : statement
              | statement SEMI statementList

statement : compoundStatement
          | callStatement
          | assignmentStatement
          | writeStatement
          | empty

callStatement : call

call : ID LPAREN (expr (COMMA expr)*)? RPAREN

writeStatement : (WRITE | WRITELN) (LPAREN (expr (COMMA expr)*)? RPAREN)?

assignmentStatement : identifier ASSIGN expr

ifStatement : IF expr THEN statement (ELSE statement)?

empty :

expr: artihmeticExpr (relationalOperator arithmeticExpr)?

relOp : EQUAL | NOT_EQUAL | LESS_THAN | LESS_THAN_EQUAL | GREATER_THAN | GREATER_THAN_EQUAL

arithmeticExpr : term ((PLUS | MINUS) term)*

term : factor ((MUL | INTEGER_DIV | FLOAT_DIV | MOD) factor)*

factor : PLUS factor
       | MINUS factor
       | INTEGER_LIT
       | REAL_LIT
       | STRING_LIT
       | TRUE
       | FALSE
       | LPAREN expr RPAREN
       | call
       | identifier

identifier: ID
*/

type Parser struct {
	lexer        *lexer.Lexer
	currentToken tokens.Token
	nextToken    tokens.Token
}

func NewParser(lexer *lexer.Lexer) (*Parser, error) {
	token, err := lexer.GetNextToken()
	if err != nil {
		return nil, err
	}

	nextToken, err := lexer.GetNextToken()
	if err != nil {
		return nil, err
	}

	return &Parser{
		lexer:        lexer,
		currentToken: token,
		nextToken:    nextToken,
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

	p.currentToken = p.nextToken

	token, err := p.lexer.GetNextToken()
	if err != nil {
		return err
	}

	p.nextToken = token

	return nil
}

func (p *Parser) program() (ir.Node, error) {
	if err := p.eat(tokens.PROGRAM); err != nil {
		return nil, err
	}

	varNode, err := p.identifier()
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
	declarationNodes, err := p.declarationList()
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

func (p *Parser) declarationList() ([]ir.Node, error) {
	var declarations []ir.Node

	for p.currentToken.Type == tokens.VAR || p.currentToken.Type == tokens.PROCEDURE || p.currentToken.Type == tokens.FUNCTION {
		newDeclarations, err := p.declaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, newDeclarations...)
	}

	return declarations, nil
}

func (p *Parser) declaration() ([]ir.Node, error) {
	var declarations []ir.Node
	if p.currentToken.Type == tokens.VAR {
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

		return declarations, nil
	}

	if p.currentToken.Type == tokens.PROCEDURE {
		procDecl, err := p.procedureDeclaraton()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, procDecl)
		return declarations, nil
	}

	if p.currentToken.Type == tokens.FUNCTION {
		funcDecl, err := p.functionDeclaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, funcDecl)
		return declarations, nil
	}

	return nil, p.error(errors.UnexpectedToken)
}

func (p *Parser) procedureDeclaraton() (*ir.ProcedureDecl, error) {
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

func (p *Parser) functionDeclaration() (*ir.FunctionDecl, error) {
	if err := p.eat(tokens.FUNCTION); err != nil {
		return nil, err
	}

	funcName := p.currentToken.Value.(string)

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

	if err := p.eat(tokens.COLON); err != nil {
		return nil, err
	}

	typeNode, err := p.typeSpec()
	if err != nil {
		return nil, err
	}

	if err := p.eat(tokens.SEMI); err != nil {
		return nil, err
	}

	blockNode, err := p.block()
	if err != nil {
		return nil, err
	}

	funcDecl := ir.NewFunctionDecl(funcName, blockNode, params, typeNode)

	if err := p.eat(tokens.SEMI); err != nil {
		return nil, err
	}

	return funcDecl, nil
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
		paramNode := ir.NewParam(ir.NewIdentifier(paramToken), typeNode)
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
	varNodes := []*ir.Identifier{{
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

		varNodes = append(varNodes, ir.NewIdentifier(p.currentToken))
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

func (p *Parser) typeSpec() (*ir.TypeN, error) {
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

	case tokens.BOOLEAN:
		if err := p.eat(tokens.BOOLEAN); err != nil {
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
		if p.nextToken.Type == tokens.ASSIGN {
			return p.assignmentStatement()
		} else {
			return p.call()
		}
	case tokens.WRITELN:
		return p.writeStatement()

	case tokens.WRITE:
		return p.writeStatement()

	case tokens.IF:
		return p.ifStatement()

	default:
		return p.empty()
	}
}

func (p *Parser) ifStatement() (*ir.IfStatement, error) {

	if err := p.eat(tokens.IF); err != nil {
		return nil, err
	}

	condition, err := p.expr()
	if err != nil {
		return nil, err
	}

	if err := p.eat(tokens.THEN); err != nil {
		return nil, err
	}

	statementNode, err := p.statement()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type == tokens.ELSE {
		if err := p.eat(tokens.ELSE); err != nil {
			return nil, err
		}

		alternativeNode, err := p.statement()
		if err != nil {
			return nil, err
		}

		return ir.NewIfStatement(condition, statementNode, alternativeNode, p.currentToken), nil
	}

	return ir.NewIfStatement(condition, statementNode, nil, p.currentToken), nil
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

	var expressions []ir.Node

	if p.currentToken.Type == tokens.LPAREN {
		if err := p.eat(tokens.LPAREN); err != nil {
			return nil, err
		}

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
	}

	if token.Type == tokens.WRITE {
		return ir.NewWriteStatement(false, expressions, token), nil
	} else {
		return ir.NewWriteStatement(true, expressions, token), nil
	}
}

func (p *Parser) call() (ir.Node, error) {
	token := p.currentToken
	callName := token.Value.(string)

	if err := p.eat(tokens.ID); err != nil {
		return nil, err
	}

	if p.currentToken.Type == tokens.LPAREN {
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

		return ir.NewCall(callName, actualParams, token), nil
	}

	return ir.NewCall(callName, nil, token), nil
}

func (p *Parser) assignmentStatement() (*ir.Assign, error) {
	left, err := p.identifier()
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

func (p *Parser) identifier() (*ir.Identifier, error) {
	node := ir.NewIdentifier(p.currentToken)

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

	case tokens.TRUE:
		if err := p.eat(tokens.TRUE); err != nil {
			return nil, err
		}

		return ir.NewBooleanLit(token), nil

	default:
		if p.currentToken.Type == tokens.ID && p.nextToken.Type == tokens.LPAREN {
			return p.call()
		}
		return p.identifier()
	}
}

func (p *Parser) term() (ir.Node, error) {
	node, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type == tokens.MUL ||
		p.currentToken.Type == tokens.INTEGER_DIV ||
		p.currentToken.Type == tokens.FLOAT_DIV ||
		p.currentToken.Type == tokens.MOD {

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

		case tokens.MOD:
			if err := p.eat(tokens.MOD); err != nil {
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

func (p *Parser) arithmeticExpr() (ir.Node, error) {
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

func (p *Parser) expr() (ir.Node, error) {
	node, err := p.arithmeticExpr()
	if err != nil {
		return nil, err
	}

	if isRelOp(p.currentToken.Type) {
		token := p.currentToken
		if err := p.eat(token.Type); err != nil {
			return nil, err
		}

		right, err := p.arithmeticExpr()
		if err != nil {
			return nil, err
		}

		node = ir.NewBinOp(node, token, right)
	}

	return node, nil
}

func isRelOp(tokenType tokens.TokenType) bool {
	return tokenType == tokens.EQUAL ||
		tokenType == tokens.NOT_EQUAL ||
		tokenType == tokens.LESS_THAN ||
		tokenType == tokens.LESS_THAN_EQUAL ||
		tokenType == tokens.GREATER_THAN ||
		tokenType == tokens.GREATER_THAN_EQUAL
}
