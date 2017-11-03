package main

import (
	"os"
	"strconv"
	"text/scanner"
	"unicode"
)

func ParseCSS(inputFileName string) (*Stylesheet, error) {
	f, err := os.Open(inputFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s Scanner
	s.Init(f)

	sheet := parseStylesheet(&s)

	return sheet, nil
}

func parseStylesheet(s *Scanner) *Stylesheet {
	return &Stylesheet{
		Rules: parseRules(s),
	}
}

func parseRules(s *Scanner) []Rule {
	var rules []Rule

	for s.Peek() != scanner.EOF {
		if s.NextChar() == rune('}') {
			s.Scan()
			continue
		}

		rule := parseRule(s)
		if len(rule.Declarations) == 0 || len(rule.Selectors) == 0 {
			continue
		}
		rules = append(rules, rule)
	}

	return rules
}

func parseRule(s *Scanner) Rule {
	return Rule{
		Selectors:    parseSelectors(s),
		Declarations: parseDeclarations(s),
	}
}

func parseSelectors(s *Scanner) []Selector {
	var selectors []Selector

	for s.Peek() != scanner.EOF {
		selector := parseSelector(s)
		selectors = append(selectors, selector)
		if s.NextChar() == rune(',') {
			s.Scan()
			continue
		} else if s.NextChar() == rune('{') {
			s.Scan()
			break
		}
	}

	return selectors
}

func parseSelector(s *Scanner) Selector {
	selector := Selector{}

	for s.Peek() != scanner.EOF {
		if s.NextChar() == rune(',') || s.NextChar() == rune('{') {
			break
		}

		s.Scan()
		value := s.TokenText()

		if value == "#" {
			s.Scan()
			selector.ID = s.TokenText()
		} else if value == "." {
			s.Scan()
			selector.Classes = []string{s.TokenText()}
		} else {
			selector.TagName = value
		}
	}

	return selector
}

func parseDeclarations(s *Scanner) []Declaration {
	var declarations []Declaration

	for s.Peek() != scanner.EOF {
		if s.NextChar() == rune(';') {
			s.Scan()
			continue
		}
		if s.NextChar() == rune('}') {
			break
		}
		declarations = append(declarations, parseDeclaration(s))
	}

	return declarations
}

func parseDeclaration(s *Scanner) Declaration {
	identifier := parseIdentifier(s)
	value := parseValue(s)

	d := Declaration{
		Name:  identifier,
		Value: value,
	}

	return d
}

func parseIdentifier(s *Scanner) string {
	name := ""
	for s.Scan() != scanner.EOF {
		if s.TokenText() == ":" || s.TokenText() == ";" {
			break
		}
		name += s.TokenText()
	}
	return name
}

func parseValue(s *Scanner) Value {
	v := Value{}

	next := s.NextChar()

	if unicode.IsDigit(next) {
		v.Length = parseLength(s)
	} else if next == rune('#') {
		v.Color = parseColor(s)
	} else {
		keyword := parseIdentifier(s)
		switch keyword {
		case "red":
			v.Color = Color{A: 255, R: 255}
		case "blue":
			v.Color = Color{A: 255, B: 255}
		case "green":
			v.Color = Color{A: 255, G: 255}
		case "white":
			v.Color = Color{A: 255, R: 255, G: 255, B: 255}
		case "black":
			v.Color = Color{A: 255}
		default:
			v.Keyword = keyword
		}
	}

	return v
}

func parseLength(s *Scanner) Length {
	// TODO Remove this and implement the correct logic
	s.Scan()
	if s.NextChar() != rune(';') {
		s.Scan()
	}

	return Length{Quantity: 66, Unit: Px}
}

func parseColor(s *Scanner) Color {
	if s.NextChar() == rune('#') {
		s.Scan()
		s.Scan()
	}
	text := s.TokenText()

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
	} else if len(text) == 8 {
		color.A = hexToUint8(text[0:2])
		color.R = hexToUint8(text[2:4])
		color.G = hexToUint8(text[4:6])
		color.B = hexToUint8(text[6:8])
	}
	return color
}

func hexToUint8(hex string) uint8 {
	val, err := strconv.ParseUint(hex, 16, 8)
	if err != nil {
		return 0
	}
	return uint8(val)
}
