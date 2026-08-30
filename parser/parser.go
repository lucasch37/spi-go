package parser

import (
	"fmt"

	"github.com/lucasch37/spi-go/ast"
	"github.com/lucasch37/spi-go/lexer"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
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

func (p *Parser) error() error {
	return fmt.Errorf("Invalid syntax")
}

func (p *Parser) program() (ast.Node, error) {
	if err := p.eat(lexer.PROGRAM); err != nil {
		return nil, err
	}

	varNode, err := p.variable()
	if err != nil {
		return nil, err
	}

	progName := varNode.Value

	if err := p.eat(lexer.SEMI); err != nil {
		return nil, err
	}

	blockNode, err := p.block()
	if err != nil {
		return nil, err
	}

	programNode := ast.NewProgram(progName, blockNode)

	if err := p.eat(lexer.DOT); err != nil {
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
	if p.currentToken.Type == lexer.VAR {
		if err := p.eat(lexer.VAR); err != nil {
			return nil, err
		}

		for p.currentToken.Type == lexer.ID {
			varDecls, err := p.variableDeclaration()
			if err != nil {
				return nil, err
			}

			for _, vd := range varDecls {
				declarations = append(declarations, vd)
			}

			if err := p.eat(lexer.SEMI); err != nil {
				return nil, err
			}
		}
	}

	return declarations, nil
}

func (p *Parser) variableDeclaration() ([]*ast.VarDecl, error) {
	varNodes := []*ast.Var{{
		Token: p.currentToken,
		Value: p.currentToken.Value.(string),
	}}
	if err := p.eat(lexer.ID); err != nil {
		return nil, err
	}

	for p.currentToken.Type == lexer.COMMA {
		if err := p.eat(lexer.COMMA); err != nil {
			return nil, err
		}
		varNodes = append(varNodes, ast.NewVar(p.currentToken))
		if err := p.eat(lexer.ID); err != nil {
			return nil, err
		}
	}

	if err := p.eat(lexer.COLON); err != nil {
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
	if p.currentToken.Type == lexer.INTEGER {
		if err := p.eat(lexer.INTEGER); err != nil {
			return nil, err
		}
	} else {
		if err := p.eat(lexer.REAL); err != nil {
			return nil, err
		}
	}
	node := ast.NewType(token)
	return node, nil
}

func (p *Parser) compoundStatement() (*ast.Compound, error) {
	if err := p.eat(lexer.BEGIN); err != nil {
		return nil, err
	}

	nodes, err := p.statementList()
	if err != nil {
		return nil, err
	}

	if err := p.eat(lexer.END); err != nil {
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

	for p.currentToken.Type == lexer.SEMI {
		if err := p.eat(lexer.SEMI); err != nil {
			return nil, err
		}

		node, err := p.statement()
		if err != nil {
			return nil, err
		}
		results = append(results, node)
	}

	if p.currentToken.Type == lexer.ID {
		p.error()
	}

	return results, nil
}

func (p *Parser) statement() (ast.Node, error) {
	switch p.currentToken.Type {
	case lexer.BEGIN:
		return p.compoundStatement()
	case lexer.ID:
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

	if err := p.eat(lexer.ASSIGN); err != nil {
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
	if p.currentToken.Type != lexer.ID {
		return nil, p.error()
	}

	node := ast.NewVar(p.currentToken)

	if err := p.eat(lexer.ID); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) empty() (*ast.NoOp, error) {
	return ast.NewNoOp(), nil
}

func (p *Parser) eat(tokenType lexer.TokenType) error {
	if p.currentToken.Type != tokenType {
		return p.error()
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
	case lexer.PLUS:
		if err := p.eat(lexer.PLUS); err != nil {
			return nil, err
		}

		expr, err := p.factor()
		if err != nil {
			return nil, err
		}

		node := ast.NewUnaryOp(token, expr)

		return node, nil

	case lexer.MINUS:
		if err := p.eat(lexer.MINUS); err != nil {
			return nil, err
		}

		expr, err := p.factor()
		if err != nil {
			return nil, err
		}

		node := ast.NewUnaryOp(token, expr)
		return node, nil

	case lexer.INTEGER_CONST:
		if err := p.eat(lexer.INTEGER_CONST); err != nil {
			return nil, err
		}

		return ast.NewIntegerLit(token), nil

	case lexer.REAL_CONST:
		if err := p.eat(lexer.REAL_CONST); err != nil {
			return nil, err
		}

		return ast.NewRealLit(token), nil

	case lexer.LPAREN:
		if err := p.eat(lexer.LPAREN); err != nil {
			return nil, err
		}

		node, err := p.expr()
		if err != nil {
			return nil, err
		}

		if err := p.eat(lexer.RPAREN); err != nil {
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

	for p.currentToken.Type == lexer.MUL ||
		p.currentToken.Type == lexer.INTEGER_DIV || p.currentToken.Type == lexer.FLOAT_DIV {

		token := p.currentToken

		switch token.Type {
		case lexer.MUL:
			if err := p.eat(lexer.MUL); err != nil {
				return nil, err
			}
		case lexer.INTEGER_DIV:
			if err := p.eat(lexer.INTEGER_DIV); err != nil {
				return nil, err
			}
		case lexer.FLOAT_DIV:
			if err := p.eat(lexer.FLOAT_DIV); err != nil {
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

	for p.currentToken.Type == lexer.PLUS ||
		p.currentToken.Type == lexer.MINUS {

		token := p.currentToken

		if token.Type == lexer.PLUS {
			if err := p.eat(lexer.PLUS); err != nil {
				return nil, err
			}
		} else {
			if err := p.eat(lexer.MINUS); err != nil {
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

	if p.currentToken.Type != lexer.EOF {
		return nil, p.error()
	}

	return tree, nil
}
