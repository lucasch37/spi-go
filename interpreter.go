package main

type Interpreter struct {
	parser      *Parser
	GLOBALSCOPE map[string]int
}

func NewInterpreter(parser *Parser) *Interpreter {
	return &Interpreter{
		parser:      parser,
		GLOBALSCOPE: make(map[string]int),
	}
}

func (i *Interpreter) visit(node AST) int {
	switch node := node.(type) {

	case BinOp:
		return i.visitBinOp(node)

	case Num:
		return i.visitNum(node)

	case UnaryOp:
		return i.visitUnaryOp(node)

	case Compound:
		return i.visitCompound(node)

	case NoOp:
		return i.visitNoOp(node)

	case Assign:
		return i.visitAssign(node)

	case Var:
		return i.visitVar(node)

	default:
		panic("No visit method for node")
	}
}

func (i *Interpreter) visitBinOp(node BinOp) int {
	switch node.Op.Type {

	case PLUS:
		return i.visit(node.Left) + i.visit(node.Right)

	case MINUS:
		return i.visit(node.Left) - i.visit(node.Right)

	case MUL:
		return i.visit(node.Left) * i.visit(node.Right)

	case DIV:
		return i.visit(node.Left) / i.visit(node.Right)

	default:
		panic("Unknown binary operator")
	}
}

func (i *Interpreter) visitNum(node Num) int {
	return node.Value
}

func (i *Interpreter) visitUnaryOp(node UnaryOp) int {
	switch node.Op.Type {

	case PLUS:
		return i.visit(node.Expr)

	case MINUS:
		return i.visit(node.Expr) * -1

	default:
		panic("Unknown unary operator")
	}
}

func (i *Interpreter) visitCompound(node Compound) int {
	for _, child := range node.Children {
		i.visit(child)
	}
	return 0
}

func (i *Interpreter) visitNoOp(node NoOp) int {
	return 0
}

func (i *Interpreter) visitAssign(node Assign) int {
	varName := node.Left.Value
	i.GLOBALSCOPE[varName] = int(i.visit(node.Right))
	return 0
}

func (i *Interpreter) visitVar(node Var) int {
	varName := node.Value
	if val, exists := i.GLOBALSCOPE[varName]; exists {
		return val
	} else {
		panic("Variable not found: " + varName)
	}
}

func (i *Interpreter) Interpret() int {
	tree := i.parser.Parse()
	return i.visit(tree)
}
