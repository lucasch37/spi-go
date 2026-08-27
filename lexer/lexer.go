package lexer

import "strconv"

var reservedKeywords = map[string]Token{
	"BEGIN": {Type: BEGIN, Value: "BEGIN"},
	"END":   {Type: END, Value: "END"},
}

type Lexer struct {
	text        string
	pos         int
	currentChar byte
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		text:        text,
		pos:         0,
		currentChar: text[0],
	}
}

func (l *Lexer) error() {
	panic("Invalid character")
}

func (l *Lexer) advance() {
	l.pos++

	if l.pos >= len(l.text) {
		// 0 represents EOF
		l.currentChar = 0
	} else {
		l.currentChar = l.text[l.pos]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar != 0 && isWhitespace(l.currentChar) {
		l.advance()
	}
}

func (l *Lexer) peek() byte {
	peekPos := l.pos + 1

	if peekPos > len(l.text)-1 {
		return 0
	} else {
		return l.text[peekPos]
	}
}

func (l *Lexer) integer() int {
	result := ""

	for l.currentChar != 0 && l.currentChar >= '0' && l.currentChar <= '9' {
		result += string(l.currentChar)
		l.advance()
	}

	value, err := strconv.Atoi(result)
	if err != nil {
		panic(err)
	}

	return value
}

func (l *Lexer) id() Token {
	result := ""

	for l.currentChar != 0 && isAlphaNumeric(l.currentChar) {
		result += string(l.currentChar)
		l.advance()
	}

	if token, exists := reservedKeywords[result]; exists {
		return token
	}

	return Token{
		Type:  ID,
		Value: result,
	}
}

func (l *Lexer) GetNextToken() Token {
	for l.currentChar != 0 {

		if isWhitespace(l.currentChar) {
			l.skipWhitespace()
			continue
		}

		if l.currentChar >= '0' && l.currentChar <= '9' {
			return Token{
				Type:  INTEGER,
				Value: l.integer(),
			}
		}

		if isAlpha(l.currentChar) {
			return l.id()
		}

		if l.currentChar == ':' && l.peek() == '=' {
			l.advance()
			l.advance()
			return Token{
				Type:  ASSIGN,
				Value: ":=",
			}
		}

		switch l.currentChar {
		case ';':
			l.advance()
			return Token{SEMI, ";"}

		case '.':
			l.advance()
			return Token{DOT, "."}

		case '+':
			l.advance()
			return Token{PLUS, "+"}

		case '-':
			l.advance()
			return Token{MINUS, "-"}

		case '*':
			l.advance()
			return Token{MUL, "*"}

		case '/':
			l.advance()
			return Token{DIV, "/"}

		case '(':
			l.advance()
			return Token{LPAREN, "("}

		case ')':
			l.advance()
			return Token{RPAREN, ")"}
		}

		l.error()
	}

	return Token{
		Type:  EOF,
		Value: nil,
	}
}

func isWhitespace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	default:
		return false
	}
}

func isAlphaNumeric(b byte) bool {
	switch {
	case '0' <= b && b <= '9':
		return true
	case 'a' <= b && b <= 'z':
		return true
	case 'A' <= b && b <= 'Z':
		return true
	}
	return false
}

func isAlpha(b byte) bool {
	switch {
	case 'a' <= b && b <= 'z':
		return true
	case 'A' <= b && b <= 'Z':
		return true
	}
	return false
}
