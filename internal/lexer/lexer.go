package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lucasch37/spi-go/internal/errors"
	"github.com/lucasch37/spi-go/internal/tokens"
)

type Lexer struct {
	text        string
	pos         int
	CurrentChar byte
	lineNo      int
	column      int
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		text:        text,
		pos:         0,
		CurrentChar: text[0],
		lineNo:      1,
		column:      1,
	}
}

func (l *Lexer) error(code errors.ErrorCode) error {
	s := fmt.Sprintf("%s -> line: %d column: %d", code.String(), l.lineNo, l.column)
	return errors.NewLexicalError(-1, tokens.Token{}, s)
}

func (l *Lexer) advance() {
	if l.CurrentChar == '\n' {
		l.lineNo++
		l.column = 0
	}

	l.pos++

	if l.pos >= len(l.text) {
		// 0 represents EOF
		l.CurrentChar = 0
	} else {
		l.CurrentChar = l.text[l.pos]
		l.column++
	}
}

func (l *Lexer) skipWhitespace() {
	for l.CurrentChar != 0 && isWhitespace(l.CurrentChar) {
		l.advance()
	}
}

func (l *Lexer) skipComment() error {
	for l.CurrentChar != 0 && l.CurrentChar != '}' {
		l.advance()
	}

	if l.CurrentChar == 0 {
		return l.error(errors.UnterminatedComment)
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

func (l *Lexer) number() (tokens.Token, error) {
	var result strings.Builder

	for l.CurrentChar != 0 && l.CurrentChar >= '0' && l.CurrentChar <= '9' {
		result.WriteString(string(l.CurrentChar))
		l.advance()
	}

	if l.CurrentChar == '.' {
		result.WriteString(string(l.CurrentChar))
		l.advance()

		for l.CurrentChar != 0 && l.CurrentChar >= '0' && l.CurrentChar <= '9' {
			result.WriteString(string(l.CurrentChar))
			l.advance()
		}

		value, err := strconv.ParseFloat(result.String(), 64)

		if err != nil {
			return tokens.Token{}, l.error(errors.InvalidNumberLiteral)
		}

		return tokens.Token{
			Type:   tokens.REAL_CONST,
			Value:  value,
			LineNo: l.lineNo,
			Column: l.column,
		}, nil
	}

	value, err := strconv.Atoi(result.String())

	if err != nil {
		return tokens.Token{}, l.error(errors.InvalidNumberLiteral)
	}

	return tokens.Token{
		Type:   tokens.INTEGER_CONST,
		Value:  value,
		LineNo: l.lineNo,
		Column: l.column,
	}, nil
}

func (l *Lexer) id() tokens.Token {
	token := tokens.Token{
		Type:   -1,
		Value:  nil,
		LineNo: l.lineNo,
		Column: l.column,
	}

	result := ""

	for l.CurrentChar != 0 && isAlphaNumeric(l.CurrentChar) {
		result += string(l.CurrentChar)
		l.advance()
	}

	if resToken, exists := tokens.ReservedKeywords[strings.ToUpper(result)]; exists {
		token.Type = resToken.Type
		token.Value = resToken.Value
		return token
	}

	token.Type = tokens.ID
	token.Value = result

	return token
}

func (l *Lexer) GetNextToken() (tokens.Token, error) {
	for l.CurrentChar != 0 {

		if isWhitespace(l.CurrentChar) {
			l.skipWhitespace()
			continue
		}

		if l.CurrentChar >= '0' && l.CurrentChar <= '9' {
			return l.number()
		}

		if isAlpha(l.CurrentChar) {
			return l.id(), nil
		}

		if l.CurrentChar == ':' && l.peek() == '=' {
			l.advance()
			l.advance()
			return tokens.Token{
				Type:   tokens.ASSIGN,
				Value:  ":=",
				LineNo: l.lineNo,
				Column: l.column,
			}, nil
		}

		if l.CurrentChar == '{' {
			l.advance()
			if err := l.skipComment(); err != nil {
				return tokens.Token{}, err
			}
			continue
		}

		tokenType, ok := tokens.TokenTypeFromName(l.CurrentChar)
		if !ok {
			l.error(errors.InvalidChar)
		}

		token := tokens.Token{
			Type:   tokenType,
			Value:  string(l.CurrentChar),
			LineNo: l.lineNo,
			Column: l.column,
		}

		l.advance()

		return token, nil
	}

	return tokens.Token{
		Type:   tokens.EOF,
		Value:  nil,
		LineNo: l.lineNo,
		Column: l.column,
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
