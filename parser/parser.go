package parser

import (
	"fmt"

	"github.com/lucasch37/spi-go/ast"
	"github.com/lucasch37/spi-go/errors"
	"github.com/lucasch37/spi-go/lexer"
	"github.com/lucasch37/spi-go/tokens"
)

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

func (p *Parser) error(code errors.ErrorCode) error {
	return errors.NewSyntaxError(code, p.currentToken, fmt.Sprintf("%s -> %s", code.String(), p.currentToken.String()))
}

func (p *Parser) program() (ast.Node, error) {
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

	programNode := ast.NewProgram(progName, blockNode)

	if err := p.eat(tokens.DOT); err != nil {
		return nil, err
	}

	return programNode, nil
}

func (p *Parser) block() (*ast.Block, error) {
	declarationNodes, err := p.declarations()
	if err != nil {
		return nil, err
	}

	compoundStatementNode, err := p.compoundStatement()
	if err != nil {
		return nil, err
	}

	node := ast.NewBlock(declarationNodes, compoundStatementNode)
	return node, nil
}

func (p *Parser) declarations() ([]ast.Node, error) {
	var declarations []ast.Node
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
		if err := p.eat(tokens.PROCEDURE); err != nil {
			return nil, err
		}

		procName := p.currentToken.Value.(string)

		if err := p.eat(tokens.ID); err != nil {
			return nil, err
		}

		var params []*ast.Param

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

		procDecl := ast.NewProcedureDecl(procName, blockNode, params)
		declarations = append(declarations, procDecl)

		if err := p.eat(tokens.SEMI); err != nil {
			return nil, err
		}
	}

	return declarations, nil
}

func (p *Parser) formalParamaters() ([]*ast.Param, error) {
	var paramNodes []*ast.Param

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
		paramNode := ast.NewParam(ast.NewVar(paramToken), typeNode)
		paramNodes = append(paramNodes, paramNode)
	}

	return paramNodes, nil
}

func (p *Parser) formalParameterList() ([]*ast.Param, error) {
	if p.currentToken.Type != tokens.ID {
		return make([]*ast.Param, 0), nil
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

func (p *Parser) variableDeclaration() ([]*ast.VarDecl, error) {
	varNodes := []*ast.Var{{
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

		varNodes = append(varNodes, ast.NewVar(p.currentToken))
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

	var varDeclarations []*ast.VarDecl
	for _, varNode := range varNodes {
		varDeclarations = append(varDeclarations, ast.NewVarDecl(varNode, typeNode))
	}

	return varDeclarations, nil
}

func (p *Parser) typeSpec() (*ast.Type, error) {
	token := p.currentToken
	if p.currentToken.Type == tokens.INTEGER {
		if err := p.eat(tokens.INTEGER); err != nil {
			return nil, err
		}
	} else {
		if err := p.eat(tokens.REAL); err != nil {
			return nil, err
		}
	}
	node := ast.NewType(token)
	return node, nil
}

func (p *Parser) compoundStatement() (*ast.Compound, error) {
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

	root := ast.NewCompound(nodes)
	return root, nil
}

func (p *Parser) statementList() ([]ast.Node, error) {
	node, err := p.statement()
	if err != nil {
		return nil, err
	}

	results := []ast.Node{
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

func (p *Parser) statement() (ast.Node, error) {
	switch p.currentToken.Type {
	case tokens.BEGIN:
		return p.compoundStatement()
	case tokens.ID:
		return p.assignmentStatement()
	default:
		return p.empty()
	}
}

func (p *Parser) assignmentStatement() (*ast.Assign, error) {
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

	node := ast.NewAssign(left, token, right)

	return node, nil
}

func (p *Parser) variable() (*ast.Var, error) {
	node := ast.NewVar(p.currentToken)

	if err := p.eat(tokens.ID); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) empty() (*ast.NoOp, error) {
	return ast.NewNoOp(), nil
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

func (p *Parser) factor() (ast.Node, error) {
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

		node := ast.NewUnaryOp(token, expr)

		return node, nil

	case tokens.MINUS:
		if err := p.eat(tokens.MINUS); err != nil {
			return nil, err
		}

		expr, err := p.factor()
		if err != nil {
			return nil, err
		}

		node := ast.NewUnaryOp(token, expr)
		return node, nil

	case tokens.INTEGER_CONST:
		if err := p.eat(tokens.INTEGER_CONST); err != nil {
			return nil, err
		}

		return ast.NewIntegerLit(token), nil

	case tokens.REAL_CONST:
		if err := p.eat(tokens.REAL_CONST); err != nil {
			return nil, err
		}

		return ast.NewRealLit(token), nil

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

	default:
		return p.variable()
	}
}

func (p *Parser) term() (ast.Node, error) {
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

		node = ast.NewBinOp(node, token, right)
	}

	return node, nil
}

func (p *Parser) expr() (ast.Node, error) {
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

		node = ast.NewBinOp(node, token, right)
	}

	return node, nil
}

func (p *Parser) Parse() (ast.Node, error) {
	tree, err := p.program()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != tokens.EOF {
		return nil, p.error(errors.UnexpectedToken)
	}

	return tree, nil
}
