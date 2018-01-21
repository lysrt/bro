package css

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

type Parser struct {
	lexer *Lexer

	curToken  CSSToken
	peekToken CSSToken

	errors []string
}

func NewParser(r io.Reader) *Parser {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	lexer := NewLexer(string(b))
	p := Parser{
		lexer:  lexer,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return &p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) tokenError(expected CSSTokenType) {
	msg := fmt.Sprintf("expected %s, got %s", expected, p.curToken.Litteral)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()

	// Skip comments
	if p.curToken.Type == COMMENT {
		p.nextToken()
	}
}

func (p *Parser) ParseStylesheet() *Stylesheet {
	stylesheet := &Stylesheet{}
	stylesheet.Rules = []Rule{}

	for p.curToken.Type != EOF {
		rule := p.parseRule()
		if len(rule.Declarations) != 0 && len(rule.Selectors) != 0 {
			stylesheet.Rules = append(stylesheet.Rules, rule)
		}
	}

	return stylesheet
}

func (p *Parser) parseRule() Rule {
	rule := Rule{}
	rule.Selectors = p.parseSelectors()

	if p.curToken.Type != LBRACE {
		p.tokenError(LBRACE)
	}
	p.nextToken()

	rule.Declarations = []Declaration{}
	for p.curToken.Type != RBRACE && p.peekToken.Type != EOF {
		declaration := p.parseDeclaration()
		if declaration.Name != "" {
			rule.Declarations = append(rule.Declarations, declaration)
		}
	}

	if p.curToken.Type != RBRACE {
		p.tokenError(RBRACE)
	}
	p.nextToken()

	return rule
}

func (p *Parser) parseSelectors() []Selector {
	var selectors []Selector

	for p.curToken.Type != LBRACE && p.curToken.Type != EOF {
		selector := p.parseSelector()
		selectors = append(selectors, selector)

		if p.curToken.Type == COMMA {
			p.nextToken()
			continue
		}
	}

	return selectors
}

func (p *Parser) parseSelector() Selector {
	selector := Selector{
		Classes: []string{},
	}

	for p.curToken.Type != COMMA && p.curToken.Type != LBRACE && p.curToken.Type != EOF {
		switch p.curToken.Type {
		case STAR:
			selector.TagName = "*"
			p.nextToken()
			continue
		case IDENTIFIER:
			selector.TagName = p.curToken.Litteral
			p.nextToken()
			continue
		case HASH:
			if p.peekToken.Type != IDENTIFIER {
				p.tokenError(IDENTIFIER)
				p.nextToken()
				return selector
			}
			selector.ID = p.peekToken.Litteral
			p.nextToken()
			p.nextToken()
			continue
		case DOT:
			if p.peekToken.Type != IDENTIFIER {
				p.tokenError(IDENTIFIER)
				p.nextToken()
				return selector
			}
			selector.Classes = append(selector.Classes, p.peekToken.Litteral)
			p.nextToken()
			p.nextToken()
			continue
		default:
			p.tokenError(IDENTIFIER)
			p.tokenError(STAR)
			p.tokenError(DOT)
			p.tokenError(HASH)
			p.nextToken()
			return selector
		}
	}

	return selector
}

func (p *Parser) parseDeclaration() Declaration {
	d := Declaration{}

	if p.curToken.Type != IDENTIFIER {
		p.tokenError(IDENTIFIER)
		p.nextToken()
		return d
	}

	d.Name = p.curToken.Litteral
	p.nextToken()

	if p.curToken.Type != COLON {
		p.tokenError(COLON)
		return d
	}
	p.nextToken()

	d.Value = p.parseValue()

	if p.curToken.Type != SEMICOLON {
		p.tokenError(SEMICOLON)
		return d
	}
	p.nextToken()

	return d
}

func (p *Parser) parseValue() Value {
	v := Value{}

	if p.peekToken.Type == SEMICOLON || p.peekToken.Type == EOF {
		// Value made of one token
		if p.curToken.Type == IDENTIFIER {
			v.Keyword = p.curToken.Litteral
			p.nextToken()
		} else if p.curToken.Type == NUMBER {
			v.Length = p.parseLength()
		}
	} else {
		// Value made of two or more tokens
		if p.curToken.Type == HASH {
			v.Color = p.parseColor()
		}

		if p.curToken.Type == NUMBER {
			v.Length = p.parseLength()
		}
	}

	return v
}

func (p *Parser) parseLength() Length {
	length := Length{}

	if p.curToken.Type != NUMBER {
		panic("wrong length, expecting NUMBER")
	}

	value := p.curToken.Litteral
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		f = 0
	}
	length.Quantity = f

	if p.peekToken.Type == SEMICOLON {
		return length
	}
	if p.peekToken.Type != IDENTIFIER {
		panic("bad number, expected IDENTIFIER or SEMICOLON(;) after NUMBER")
	}

	p.nextToken()
	t := p.curToken.Litteral

	var unit Unit
	switch t {
	case "px":
		unit = Px
	case "em":
		unit = Em
	case "%":
		unit = Percent
	}

	length.Unit = unit

	p.nextToken()

	return length
}

func (p *Parser) parseColor() Color {
	if p.curToken.Type != HASH || (p.peekToken.Type != IDENTIFIER && p.peekToken.Type != NUMBER) {
		panic("bad color")
	}

	text := p.peekToken.Litteral

	color := Color{}
	if len(text) == 3 {
		color.A = 255
		color.R = hexToUint8(string(text[0]) + string(text[0]))
		color.G = hexToUint8(string(text[1]) + string(text[1]))
		color.B = hexToUint8(string(text[2]) + string(text[2]))
	} else if len(text) == 6 {
		color.A = 255
		color.R = hexToUint8(text[0:2])
		color.G = hexToUint8(text[2:4])
		color.B = hexToUint8(text[4:6])
	}

	p.nextToken()
	p.nextToken()

	return color
}

func hexToUint8(hex string) uint8 {
	val, err := strconv.ParseUint(hex, 16, 8)
	if err != nil {
		return 0
	}
	return uint8(val)
}
