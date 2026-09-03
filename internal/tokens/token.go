package tokens

import "fmt"

type TokenType int

const (
	// single-character types
	PLUS TokenType = iota
	MINUS
	MUL
	FLOAT_DIV
	LPAREN
	RPAREN
	SEMI
	DOT
	COLON
	COMMA

	// keywords block
	PROGRAM
	INTEGER
	REAL
	INTEGER_DIV
	VAR
	PROCEDURE
	BEGIN
	END

	// misc
	ID
	INTEGER_CONST
	REAL_CONST
	ASSIGN
	EOF
)

var tokenTypeNames = [...]string{
	"+",
	"-",
	"*",
	"/",
	"(",
	")",
	";",
	".",
	":",
	",",

	"PROGRAM",
	"INTEGER",
	"REAL",
	"DIV",
	"VAR",
	"PROCEDURE",
	"BEGIN",
	"END",

	"ID",
	"INTEGER_CONST",
	"REAL_CONST",
	"ASSIGN",
	"EOF",
}

func (t TokenType) String() string {
	if int(t) >= 0 && int(t) < len(tokenTypeNames) {
		return tokenTypeNames[t]
	}
	return "UNKNOWN"
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %v, position=%d:%d)", t.Type.String(), t.Value, t.LineNo, t.Column)
}

type Token struct {
	Type   TokenType
	Value  any
	LineNo int
	Column int
}

var ReservedKeywords = makeReservedKeywords()

func makeReservedKeywords() map[string]Token {
	keywords := make(map[string]Token)

	for t := PROGRAM; t <= END; t++ {
		name := t.String()
		keywords[name] = Token{
			Type:  t,
			Value: name,
		}
	}

	return keywords
}

func TokenTypeFromName(char byte) (TokenType, bool) {
	name := string(char)

	for i, tokenName := range tokenTypeNames {
		if tokenName == name {
			return TokenType(i), true
		}
	}

	return 0, false
}
