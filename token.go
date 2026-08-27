package main

import "fmt"

type TokenType string

const (
	INTEGER TokenType = "INTEGER"
	PLUS    TokenType = "PLUS"
	MINUS   TokenType = "MINUS"
	MUL     TokenType = "MUL"
	DIV     TokenType = "DIV"
	LPAREN  TokenType = "("
	RPAREN  TokenType = ")"
	EOF     TokenType = "EOF"
	BEGIN   TokenType = "BEGIN"
	END     TokenType = "END"
	DOT     TokenType = "DOT"
	ID      TokenType = "ID"
	ASSIGN  TokenType = "ASSIGN"
	SEMI    TokenType = "SEMI"
)

type Token struct {
	Type  TokenType
	Value any
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %v)", t.Type, t.Value)
}
