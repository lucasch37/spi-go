package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

var reservedKeywords = map[string]Token{
	"BEGIN":     {Type: BEGIN, Value: "BEGIN"},
	"END":       {Type: END, Value: "END"},
	"PROGRAM":   {Type: PROGRAM, Value: "PROGRAM"},
	"DIV":       {Type: INTEGER_DIV, Value: "DIV"},
	"INTEGER":   {Type: INTEGER, Value: "INTEGER"},
	"VAR":       {Type: VAR, Value: "VAR"},
	"REAL":      {Type: REAL, Value: "REAL"},
	"PROCEDURE": {Type: PROCEDURE, Value: "PROCEDURE"},
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

func (l *Lexer) error() error {
	return fmt.Errorf("Invalid Character")
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

func (l *Lexer) skipComment() error {
	for l.currentChar != 0 && l.currentChar != '}' {
		l.advance()
	}

	if l.currentChar == 0 {
		return fmt.Errorf("Unterminated comment")
	}

	l.advance()
	return nil
}

func (l *Lexer) peek() byte {
	peekPos := l.pos + 1

	if peekPos > len(l.text)-1 {
		return 0
	} else {
		return l.text[peekPos]
	}
}

func (l *Lexer) number() (Token, error) {
	var result strings.Builder

	for l.currentChar != 0 && l.currentChar >= '0' && l.currentChar <= '9' {
		result.WriteString(string(l.currentChar))
		l.advance()
	}

	if l.currentChar == '.' {
		result.WriteString(string(l.currentChar))
		l.advance()

		for l.currentChar != 0 && l.currentChar >= '0' && l.currentChar <= '9' {
			result.WriteString(string(l.currentChar))
			l.advance()
		}

		value, err := strconv.ParseFloat(result.String(), 64)

		if err != nil {
			return Token{}, err
		}

		return Token{
			Type:  REAL_CONST,
			Value: value,
		}, nil
	}

	value, err := strconv.Atoi(result.String())

	if err != nil {
		return Token{}, err
	}

	return Token{
		Type:  INTEGER_CONST,
		Value: value,
	}, nil
}

func (l *Lexer) id() Token {
	result := ""

	for l.currentChar != 0 && isAlphaNumeric(l.currentChar) {
		result += string(l.currentChar)
		l.advance()
	}

	if token, exists := reservedKeywords[strings.ToUpper(result)]; exists {
		return token
	}

	return Token{
		Type:  ID,
		Value: result,
	}
}

func (l *Lexer) GetNextToken() (Token, error) {
	for l.currentChar != 0 {

		if isWhitespace(l.currentChar) {
			l.skipWhitespace()
			continue
		}

		if l.currentChar >= '0' && l.currentChar <= '9' {
			return l.number()
		}

		if isAlpha(l.currentChar) {
			return l.id(), nil
		}

		if l.currentChar == ':' && l.peek() == '=' {
			l.advance()
			l.advance()
			return Token{
				Type:  ASSIGN,
				Value: ":=",
			}, nil
		}

		switch l.currentChar {
		case '{':
			l.advance()
			if err := l.skipComment(); err != nil {
				return Token{}, err
			}
			continue

		case ';':
			l.advance()
			return Token{SEMI, ";"}, nil

		case '.':
			l.advance()
			return Token{DOT, "."}, nil

		case '+':
			l.advance()
			return Token{PLUS, "+"}, nil

		case '-':
			l.advance()
			return Token{MINUS, "-"}, nil

		case '*':
			l.advance()
			return Token{MUL, "*"}, nil

		case ':':
			l.advance()
			return Token{COLON, ":"}, nil

		case ',':
			l.advance()
			return Token{COMMA, ","}, nil

		case '/':
			l.advance()
			return Token{FLOAT_DIV, "/"}, nil

		case '(':
			l.advance()
			return Token{LPAREN, "("}, nil

		case ')':
			l.advance()
			return Token{RPAREN, ")"}, nil
		}

		l.error()
	}

	return Token{
		Type:  EOF,
		Value: nil,
	}, nil
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
