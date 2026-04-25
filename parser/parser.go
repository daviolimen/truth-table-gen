package parser

import (
	"fmt"
	"truth-table/ast"
	"truth-table/lexer"
	"truth-table/token"
)

const (
	_ int = iota
	LOWEST
	IMPLICATION
	OR
	XOR
	AND
	NOT
)

type Parser struct {
	l         *lexer.Lexer
	errorList []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errorList: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.LITERAL, p.parseLiteral)
	p.registerPrefix(token.VARIABLE, p.parseVariable)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.XOR, p.parseInfixExpression)
	p.registerInfix(token.IMPLIES, p.parseInfixExpression)
	p.registerInfix(token.IFF, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) errors() []string {
	return p.errorList
}

func (p *Parser) Parse() (ast.Node, []string) {
	expr := p.ParseExpression(LOWEST)

	if !p.peekTokenIs(token.EOF) {
		msg := fmt.Sprintf("unexpected token: %s", p.peekToken.Literal)
		p.errorList = append(p.errorList, msg)
	}

	return expr, p.errorList
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errorList = append(p.errorList, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errorList = append(p.errorList, msg)
}

func (p *Parser) ParseExpression(precedence int) ast.Node {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseVariable() ast.Node {
	return &ast.Variable{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseLiteral() ast.Node {
	return &ast.Literal{Token: p.curToken, Value: p.curToken.Literal == "T"}
}

func (p *Parser) parsePrefixExpression() ast.Node {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.ParseExpression(NOT)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Node) ast.Node {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.ParseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Node {
	p.nextToken()

	exp := p.ParseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

type (
	prefixParseFn func() ast.Node
	infixParseFn  func(ast.Node) ast.Node
)

var precedences = map[token.TokenType]int{
	token.NOT:     NOT,
	token.AND:     AND,
	token.OR:      OR,
	token.XOR:     XOR,
	token.IMPLIES: IMPLICATION,
	token.IFF:     IMPLICATION,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
