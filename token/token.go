package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	LITERAL  = "LITERAL"
	VARIABLE = "VARIABLE"

	NOT     = "~"
	AND     = "&"
	OR      = "|"
	XOR     = "^"
	IMPLIES = "=>"
	IFF     = "<>"

	LPAREN = "("
	RPAREN = ")"
)
