package lexer

import (
	"truth-table/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '>' {
			l.readChar()
			tok = token.Token{Type: token.IMPLIES, Literal: "=>"}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '<':
		if l.peekChar() == '>' {
			l.readChar()
			tok = token.Token{Type: token.IFF, Literal: "<>"}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '~':
		tok = newToken(token.NOT, l.ch)
	case '&':
		tok = newToken(token.AND, l.ch)
	case '|':
		tok = newToken(token.OR, l.ch)
	case '^':
		tok = newToken(token.XOR, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 'T':
		tok = newToken(token.LITERAL, l.ch)
	case 'F':
		tok = newToken(token.LITERAL, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isVariable(l.ch) {
			tok = newToken(token.VARIABLE, l.ch)
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isVariable(ch byte) bool {
	return ('a' <= ch && ch <= 'z')
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
