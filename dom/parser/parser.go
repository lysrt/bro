package parser

import (
	"fmt"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

// Error represents a parser error.
type Error struct {
	Token lexer.Token
	Msg   string
}

func (e Error) Error() string {
	return e.Msg
}

// Parser represents an HTML parser.
type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	elements []*dom.Node

	errors []Error
}

// New instanciates a new Parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []Error{},
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []Error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parse parses the document and return a dom tree.
// Before using node check for parser error with Errors.
func (p *Parser) Parse() *dom.Node {
	root := &dom.Node{}

	var previous *dom.Node
	for p.curToken.Type != lexer.TokenEOF {
		n := p.parseNode(root)
		if n == nil {
			p.nextToken()
			continue
		}
		if root.FirstChild == nil {
			root.FirstChild = n
		}
		n.PrevSibling = previous
		if previous != nil {
			previous.NextSibling = n
		}
		previous = n
		p.nextToken()
	}
	root.LastChild = previous
	return root
}

func (p *Parser) parseNode(parent *dom.Node) *dom.Node {
	switch p.curToken.Type {
	case lexer.TokenLBracket:
		if p.peekTokenIs(lexer.TokenBang) {
			//TODO: handle doctype
		} else if p.peekTokenIs(lexer.TokenSlash) {
			//TODO: handle closing element
		} else if p.peekTokenIs(lexer.TokenIdent) {
			return p.parseElement(parent)
		} else {
			p.unexpectedPeekError()
			p.nextToken()
			return nil
		}
	}
	return nil
}

func (p *Parser) parseElement(parent *dom.Node) *dom.Node {
	n := &dom.Node{Type: dom.ElementNode}
	n.Tag = p.peekToken.Literal
	p.nextToken()

	//TODO: parse attributes

	if !p.expectPeek(lexer.TokenRBracket) {
		return nil
	}
	return n
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t lexer.TokenType) {
	err := Error{Token: p.peekToken}
	err.Msg = fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) unexpectedPeekError() {
	err := Error{Token: p.peekToken}
	err.Msg = fmt.Sprintf("unexpected token %s", err.Token.Type)
	p.errors = append(p.errors, err)
}
