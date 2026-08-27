package parser

import (
	"github.com/lucasch37/spi-go/ast"
	"github.com/lucasch37/spi-go/lexer"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
}

func NewParser(lexer *lexer.Lexer) *Parser {
	return &Parser{
		lexer:        lexer,
		currentToken: lexer.GetNextToken(),
	}
}

func (p *Parser) error() {
	panic("Invalid syntax")
}

func (p *Parser) program() ast.Node {
	node := p.compoundStatement()
	p.eat(lexer.DOT)
	return node
}

func (p *Parser) compoundStatement() ast.Node {
	p.eat(lexer.BEGIN)
	nodes := p.statementList()
	p.eat(lexer.END)

	root := ast.Compound{
		Children: nodes,
	}

	return root
}

func (p *Parser) statementList() []ast.Node {
	node := p.statement()
	results := []ast.Node{
		node,
	}

	for p.currentToken.Type == lexer.SEMI {
		p.eat(lexer.SEMI)
		results = append(results, p.statement())
	}

	if p.currentToken.Type == lexer.ID {
		p.error()
	}

	return results
}

func (p *Parser) statement() ast.Node {
	var node ast.Node

	switch p.currentToken.Type {
	case lexer.BEGIN:
		node = p.compoundStatement()
	case lexer.ID:
		node = p.assignmentStatement()
	default:
		node = p.empty()
	}

	return node
}

func (p *Parser) assignmentStatement() ast.Node {
	left := p.variable()
	token := p.currentToken
	p.eat(lexer.ASSIGN)
	right := p.expr()

	node := ast.Assign{
		Left:  left,
		Token: token,
		Op:    token,
		Right: right,
	}

	return node
}

func (p *Parser) variable() ast.Var {
	node := ast.Var{
		Token: p.currentToken,
		Value: p.currentToken.Value.(string),
	}
	p.eat(lexer.ID)
	return node
}

func (p *Parser) empty() ast.Node {
	return ast.NoOp{}
}

func (p *Parser) eat(tokenType lexer.TokenType) {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.GetNextToken()
	} else {
		p.error()
	}
}

func (p *Parser) factor() ast.Node {
	token := p.currentToken

	switch token.Type {

	case lexer.PLUS:
		p.eat(lexer.PLUS)
		node := ast.UnaryOp{
			Token: token,
			Op:    token,
			Expr:  p.factor(),
		}
		return node

	case lexer.MINUS:
		p.eat(lexer.MINUS)
		node := ast.UnaryOp{
			Token: token,
			Op:    token,
			Expr:  p.factor(),
		}
		return node

	case lexer.INTEGER:
		p.eat(lexer.INTEGER)

		return ast.Num{
			Token: token,
			Value: token.Value.(int),
		}

	case lexer.LPAREN:
		p.eat(lexer.LPAREN)

		node := p.expr()

		p.eat(lexer.RPAREN)

		return node

	default:
		node := p.variable()
		return node
	}
}

func (p *Parser) term() ast.Node {
	node := p.factor()

	for p.currentToken.Type == lexer.MUL ||
		p.currentToken.Type == lexer.DIV {

		token := p.currentToken

		if token.Type == lexer.MUL {
			p.eat(lexer.MUL)
		} else {
			p.eat(lexer.DIV)
		}

		node = ast.BinOp{
			Left:  node,
			Token: token,
			Op:    token,
			Right: p.factor(),
		}
	}

	return node
}

func (p *Parser) expr() ast.Node {
	node := p.term()

	for p.currentToken.Type == lexer.PLUS ||
		p.currentToken.Type == lexer.MINUS {

		token := p.currentToken

		if token.Type == lexer.PLUS {
			p.eat(lexer.PLUS)
		} else {
			p.eat(lexer.MINUS)
		}

		node = ast.BinOp{
			Left:  node,
			Token: token,
			Op:    token,
			Right: p.term(),
		}
	}

	return node
}

func (p *Parser) Parse() ast.Node {
	node := p.program()
	if p.currentToken.Type != lexer.EOF {
		p.error()
	}

	return node
}
