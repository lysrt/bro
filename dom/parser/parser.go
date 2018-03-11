package parser

import (
	"fmt"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

// nodeClosingElement is a special node for HTML closing element.
// It should not escapes outside of the parser.
const nodeClosingElement dom.NodeType = "closing element"

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

// Errors returns a list of error encounter during parsing.
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
		n := p.parseNode(nil)
		if n == nil {
			return
		}
		nodes = append(nodes, n)
	}
}

func (p *Parser) parseNode(parent *dom.Node) (n *dom.Node) {
	switch p.curToken.Type {
	case lexer.TokenLBracket:
		if p.peekTokenIs(lexer.TokenBang) {
			//TODO: handle doctype
			p.addError(p.peekToken, "doctype not implemented")
			return
		} else if p.peekTokenIs(lexer.TokenSlash) {
			n = p.parseClosingElement()
		} else if p.peekTokenIs(lexer.TokenIdent) {
			n = p.parseElement()
		} else {
			p.addError(p.peekToken, "unexpected peek token: %q", p.peekToken.Type)
			p.nextToken()
			return
		}
	case lexer.TokenEOF:
		if len(p.elements) != 0 {
			p.addError(p.curToken, "unexpected end of file")
		}
	default:
		p.addError(p.curToken, "unexpected token: %q", p.curToken.Type)
	}
	if n != nil {
		n.Parent = parent
	}
	return
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
		n := p.parseNode(root)
		if n == nil {
			p.addError(startToken, "missing closing element")
			break
		}
		if n.Type == nodeClosingElement {
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
		if root.LastChild != nil {
			root.LastChild.NextSibling = n
		}
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
		Type: nodeClosingElement,
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
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t lexer.TokenType) {
	p.addError(p.peekToken, "expected next token to be %q, got %q instead", t, p.peekToken.Type)
}

func (p *Parser) addError(tok lexer.Token, format string, a ...interface{}) {
	p.errors = append(p.errors, Error{
		Token: tok,
		Msg:   fmt.Sprintf(format, a...),
	})
}
