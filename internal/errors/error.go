package errors

import (
	"fmt"

	"github.com/lucasch37/spi-go/internal/tokens"
)

type ErrorCode int

const (
	// syntax
	UnexpectedToken ErrorCode = iota
	ExpectedAssign

	// semantic
	IDNotFound
	DuplicateID
	WrongParamCount
	NotCallable
	InvalidOperand

	// lexical
	InvalidChar
	InvalidNumberLiteral
	UnterminatedComment

	// runtime
	DivideByZero
)

var errorCodeNames = [...]string{
	"Unexpected token",
	"Expected assignment operator",

	"Identifier not found",
	"Duplicate identifier",
	"Wrong number of parameters",
	"Not callable",
	"Invalid operand",

	"Invalid character",
	"Invalid number literal",
	"Unterminated comment",

	"Divide by zero",
}

func (e ErrorCode) String() string {
	if int(e) >= 0 && int(e) < len(errorCodeNames) {
		return errorCodeNames[e]
	}

	return "UNKNOWN"
}

type Diagnostic struct {
	Code    ErrorCode
	Token   tokens.Token
	Message string
}

func (d Diagnostic) Error() string {
	return d.Message
}

type LexicalError struct {
	Diagnostic
}

func NewLexicalError(code ErrorCode, token tokens.Token, message string) *LexicalError {
	return &LexicalError{
		Diagnostic{
			Code:    code,
			Token:   token,
			Message: fmt.Sprintf("LexerError: %s", message),
		},
	}
}

type SyntaxError struct {
	Diagnostic
}

func NewSyntaxError(code ErrorCode, token tokens.Token, message string) *SyntaxError {
	return &SyntaxError{
		Diagnostic{
			Code:    code,
			Token:   token,
			Message: fmt.Sprintf("SyntaxError: %s", message),
		},
	}
}

type SemanticError struct {
	Diagnostic
}

func NewSemanticError(code ErrorCode, token tokens.Token, message string) *SemanticError {
	return &SemanticError{
		Diagnostic{
			Code:    code,
			Token:   token,
			Message: fmt.Sprintf("SemanticError: %s", message),
		},
	}
}

type RuntimeError struct {
	Diagnostic
}

func NewRuntimeError(code ErrorCode, token tokens.Token, message string) *RuntimeError {
	return &RuntimeError{
		Diagnostic{
			Code:    code,
			Token:   token,
			Message: fmt.Sprintf("RuntimeError: %s", message),
		},
	}
}
