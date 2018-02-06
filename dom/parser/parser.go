package parser

import (
	"fmt"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

const NodeClosingElement dom.NodeType = "closing element"

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
func (p *Parser) Parse() (nodes []*dom.Node) {
	for {
		n := p.parseNode()
		if n == nil {
			return
		}
		nodes = append(nodes, n)
	}
}

func (p *Parser) parseNode() *dom.Node {
	switch p.curToken.Type {
	case lexer.TokenLBracket:
		if p.peekTokenIs(lexer.TokenBang) {
			//TODO: handle doctype
		} else if p.peekTokenIs(lexer.TokenSlash) {
			return p.parseClosingElement()
		} else if p.peekTokenIs(lexer.TokenIdent) {
			return p.parseElement()
		} else {
			p.unexpectedPeekError()
			p.nextToken()
			return nil
		}
	}
	return nil
}

// parseElement parses nodes like: `<a><b></b></a>`
func (p *Parser) parseElement() *dom.Node {
	startToken := p.curToken

	// skip LBracket
	p.nextToken()

	root := &dom.Node{
		Type: dom.NodeElement,
		Tag:  p.curToken.Literal,
	}

	//TODO: parse attributes & handle autoclose

	if !p.expectsPeek(lexer.TokenRBracket) {
		return nil
	}
	// skip RBracket
	p.nextToken()

	p.elements = append(p.elements, root)
	for {
		n := p.parseNode()
		if n == nil {
			p.addError(startToken, "missing closing element")
			break
		}
		if n.Type == NodeClosingElement {
			last := p.elements[len(p.elements)-1]
			if last.Tag != n.Tag {
				p.addError(startToken, "unexpected closing element. expected=%q got=%q", last.Tag, n.Tag)
			}
			p.elements = p.elements[:len(p.elements)-1]
			break
		}
		if root.FirstChild == nil {
			root.FirstChild = n
		}
		n.PrevSibling = root.LastChild
		root.LastChild.NextSibling = n
		root.LastChild = n
	}

	return root
}

// parseClosingElement parses elements like: `</node>`
func (p *Parser) parseClosingElement() *dom.Node {
	// skip LBracket
	p.nextToken()
	if !p.expectsPeek(lexer.TokenIdent) {
		return nil
	}
	n := &dom.Node{
		Type: NodeClosingElement,
		Tag:  p.curToken.Literal,
	}

	if !p.expectsPeek(lexer.TokenRBracket) {
		return nil
	}
	// skip RBracket
	p.nextToken()
	return n
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectsPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t lexer.TokenType) {
	p.addError(p.peekToken, "expected next token to be %s, got %s instead", t, p.peekToken.Type)
}

func (p *Parser) unexpectedPeekError() {
	p.addError(p.peekToken, "unexpected token %s", p.peekToken.Type)
}

func (p *Parser) addError(tok lexer.Token, format string, a ...interface{}) {
	p.errors = append(p.errors, Error{
		Token: tok,
		Msg:   fmt.Sprintf(format, a),
	})
}
