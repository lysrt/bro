package parser

import (
	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() *dom.Document {
	doc := &dom.Document{}
	doc.Children = []dom.Node{}

	for p.curToken.Type != lexer.TokenEOF {
		n := p.parseNode()
		if n != nil {
			doc.Children = append(doc.Children, n)
		}
		p.nextToken()
	}
	return nil
}

func (p *Parser) parseNode() dom.Node {
	return nil
}
