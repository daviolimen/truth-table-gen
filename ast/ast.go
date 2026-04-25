package ast

import (
	"bytes"
	"truth-table/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Variable struct {
	Token token.Token
	Value string
}

func (v *Variable) TokenLiteral() string { return v.Token.Literal }
func (v *Variable) String() string       { return v.Value }

type Literal struct {
	Token token.Token
	Value bool
}

func (l *Literal) TokenLiteral() string { return l.Token.Literal }
func (l *Literal) String() string       { return l.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Node
}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Node
	Operator string
	Right    Node
}

func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
