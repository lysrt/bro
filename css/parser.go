package css

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

type Parser struct {
	Lexer *Lexer

	curToken  CSSToken
	peekToken CSSToken
}

func NewParser(r io.Reader) *Parser {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	lexer := New(string(b))
	p := Parser{Lexer: lexer}

	p.nextToken()
	p.nextToken()

	return &p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.Lexer.NextToken()
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

	if p.curToken.Type == LBRACE {
		p.nextToken() // TODO Check this is always the case
	}

	rule.Declarations = []Declaration{}
	for p.curToken.Type != RBRACE && p.peekToken.Type != EOF {
		rule.Declarations = append(rule.Declarations, p.parseDeclaration())
	}

	if p.curToken.Type == RBRACE {
		p.nextToken() // TODO Check this is always the case
	}

	return rule
}

func (p *Parser) parseSelectors() []Selector {
	var selectors []Selector

	for p.curToken.Type != LBRACE && p.curToken.Type != EOF {
		selector, err := p.parseSelector()
		if err != nil {
			panic(err)
		}

		selectors = append(selectors, selector)

		if p.curToken.Type == COMMA {
			p.nextToken()
			continue
		}
	}

	return selectors
}

func (p *Parser) parseSelector() (Selector, error) {
	selector := Selector{}

	switch p.curToken.Type {
	case HASH:
		if p.peekToken.Type != IDENTIFIER {
			return selector, fmt.Errorf("bad id selector: %s", p.peekToken.Litteral)
		}
		selector.ID = p.peekToken.Litteral
		p.nextToken()
		p.nextToken()
		return selector, nil
	case DOT:
		if p.peekToken.Type != IDENTIFIER {
			return selector, fmt.Errorf("bad class selector: %s", p.peekToken.Litteral)
		}
		selector.Classes = []string{p.peekToken.Litteral}
		p.nextToken()
		p.nextToken()
		return selector, nil
	case IDENTIFIER:
		selector.TagName = p.curToken.Litteral
		p.nextToken()
		return selector, nil
	default:
		return selector, fmt.Errorf("bad selector: %s %s", p.curToken.Litteral, p.peekToken.Litteral)
	}
}

func (p *Parser) parseDeclaration() Declaration {
	if p.curToken.Type != IDENTIFIER {
		panic("bad declaration, should start witn an identifier")
	}

	identifier := p.curToken.Litteral
	p.nextToken()

	if p.curToken.Type != COLON {
		panic("bad declaration, expected COLON(:)")
	}
	p.nextToken()

	value := p.parseValue()

	if p.curToken.Type != SEMICOLON {
		panic("bad declaration, expected SEMICOLON(;), got " + p.curToken.Litteral)
	}
	p.nextToken()

	d := Declaration{
		Name:  identifier,
		Value: value,
	}

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
