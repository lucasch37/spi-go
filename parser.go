package main

type Parser struct {
	lexer        *Lexer
	currentToken Token
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer:        lexer,
		currentToken: lexer.GetNextToken(),
	}
}

func (p *Parser) error() {
	panic("Invalid syntax")
}

func (p *Parser) program() AST {
	node := p.compoundStatement()
	p.eat(DOT)
	return node
}

func (p *Parser) compoundStatement() AST {
	p.eat(BEGIN)
	nodes := p.statementList()
	p.eat(END)

	root := Compound{
		Children: nodes,
	}

	return root
}

func (p *Parser) statementList() []AST {
	node := p.statement()
	results := []AST{
		node,
	}

	for p.currentToken.Type == SEMI {
		p.eat(SEMI)
		results = append(results, p.statement())
	}

	if p.currentToken.Type == ID {
		p.error()
	}

	return results
}

func (p *Parser) statement() AST {
	var node AST

	if p.currentToken.Type == BEGIN {
		node = p.compoundStatement()
	} else if p.currentToken.Type == ID {
		node = p.assignmentStatement()
	} else {
		node = p.empty()
	}

	return node
}

func (p *Parser) assignmentStatement() AST {
	left := p.variable()
	token := p.currentToken
	p.eat(ASSIGN)
	right := p.expr()

	node := Assign{
		Left:  left,
		Token: token,
		Op:    token,
		Right: right,
	}

	return node
}

func (p *Parser) variable() Var {
	node := Var{
		Token: p.currentToken,
		Value: p.currentToken.Value.(string),
	}
	p.eat(ID)
	return node
}

func (p *Parser) empty() AST {
	return NoOp{}
}

func (p *Parser) eat(tokenType TokenType) {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.GetNextToken()
	} else {
		p.error()
	}
}

func (p *Parser) factor() AST {
	token := p.currentToken

	switch token.Type {

	case PLUS:
		p.eat(PLUS)
		node := UnaryOp{
			Token: token,
			Op:    token,
			Expr:  p.factor(),
		}
		return node

	case MINUS:
		p.eat(MINUS)
		node := UnaryOp{
			Token: token,
			Op:    token,
			Expr:  p.factor(),
		}
		return node

	case INTEGER:
		p.eat(INTEGER)

		return Num{
			Token: token,
			Value: token.Value.(int),
		}

	case LPAREN:
		p.eat(LPAREN)

		node := p.expr()

		p.eat(RPAREN)

		return node

	default:
		node := p.variable()
		return node
	}
}

func (p *Parser) term() AST {
	node := p.factor()

	for p.currentToken.Type == MUL ||
		p.currentToken.Type == DIV {

		token := p.currentToken

		if token.Type == MUL {
			p.eat(MUL)
		} else {
			p.eat(DIV)
		}

		node = BinOp{
			Left:  node,
			Token: token,
			Op:    token,
			Right: p.factor(),
		}
	}

	return node
}

func (p *Parser) expr() AST {
	node := p.term()

	for p.currentToken.Type == PLUS ||
		p.currentToken.Type == MINUS {

		token := p.currentToken

		if token.Type == PLUS {
			p.eat(PLUS)
		} else {
			p.eat(MINUS)
		}

		node = BinOp{
			Left:  node,
			Token: token,
			Op:    token,
			Right: p.term(),
		}
	}

	return node
}

func (p *Parser) Parse() AST {
	node := p.program()
	if p.currentToken.Type != EOF {
		p.error()
	}

	return node
}
