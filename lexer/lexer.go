package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lucasch37/spi-go/errors"
	"github.com/lucasch37/spi-go/tokens"
)

type Lexer struct {
	text        string
	pos         int
	currentChar byte
	lineNo      int
	column      int
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		text:        text,
		pos:         0,
		currentChar: text[0],
		lineNo:      1,
		column:      1,
	}
}

func (l *Lexer) error(code errors.ErrorCode) error {
	s := fmt.Sprintf("%s -> line: %d column: %d", code.String(), l.lineNo, l.column)
	return errors.NewLexicalError(-1, tokens.Token{}, s)
}

func (l *Lexer) advance() {
	if l.currentChar == '\n' {
		l.lineNo++
		l.column = 0
	}

	l.pos++

	if l.pos >= len(l.text) {
		// 0 represents EOF
		l.currentChar = 0
	} else {
		l.currentChar = l.text[l.pos]
		l.column++
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

	for l.currentChar != 0 && isAlphaNumeric(l.currentChar) {
		result += string(l.currentChar)
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
			return tokens.Token{
				Type:   tokens.ASSIGN,
				Value:  ":=",
				LineNo: l.lineNo,
				Column: l.column,
			}, nil
		}

		if l.currentChar == '{' {
			l.advance()
			if err := l.skipComment(); err != nil {
				return tokens.Token{}, err
			}
			continue
		}

		tokenType, ok := tokens.TokenTypeFromName(l.currentChar)
		if !ok {
			l.error(errors.InvalidChar)
		}

		token := tokens.Token{
			Type:   tokenType,
			Value:  string(l.currentChar),
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
