package lexer

import "fmt"

type TokenType int

const (
	INTEGER TokenType = iota
	REAL
	INTEGER_CONST
	REAL_CONST
	PLUS
	MINUS
	MUL
	INTEGER_DIV
	FLOAT_DIV
	LPAREN
	RPAREN
	ID
	ASSIGN
	BEGIN
	END
	SEMI
	DOT
	PROGRAM
	VAR
	COLON
	COMMA
	EOF
	PROCEDURE
)

var tokenTypeNames = [...]string{
	"INTEGER",
	"REAL",
	"INTEGER_CONST",
	"REAL_CONST",
	"PLUS",
	"MINUS",
	"MUL",
	"INTEGER_DIV",
	"FLOAT_DIV",
	"LPAREN",
	"RPAREN",
	"ID",
	"ASSIGN",
	"BEGIN",
	"END",
	"SEMI",
	"DOT",
	"PROGRAM",
	"VAR",
	"COLON",
	"COMMA",
	"EOF",
	"PROCEDURE",
}

func (t TokenType) String() string {
	if int(t) >= 0 && int(t) < len(tokenTypeNames) {
		return tokenTypeNames[t]
	}
	return "UNKNOWN"
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %v)", t.Type.String(), t.Value)
}

type Token struct {
	Type  TokenType
	Value any
}
