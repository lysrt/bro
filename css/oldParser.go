package css

// import (
// 	"io"
// 	"io/ioutil"
// 	"strconv"
// 	"text/scanner"
// 	"unicode"
// )

// func ParseCSS(r io.Reader) *Stylesheet {
// 	b, err := ioutil.ReadAll(r)
// 	if err != nil {
// 		return nil
// 	}

// 	var s scanner.Scanner
// 	s.Init(r)
// 	sheet := parseStylesheet(&s)

// 	return sheet
// }

// func parseStylesheet(s *scanner.Scanner) *Stylesheet {
// 	return &Stylesheet{
// 		Rules: parseRules(s),
// 	}
// }

// func parseRules(s *scanner.Scanner) []Rule {
// 	var rules []Rule

// 	for s.Peek() != scanner.EOF {
// 		if s.NextChar() == rune('}') {
// 			s.Scan()
// 			continue
// 		}

// 		rule := parseRule(s)
// 		if len(rule.Declarations) == 0 || len(rule.Selectors) == 0 {
// 			continue
// 		}
// 		rules = append(rules, rule)
// 	}

// 	return rules
// }

// func parseRule(s *scanner.Scanner) Rule {
// 	return Rule{
// 		Selectors:    parseSelectors(s),
// 		Declarations: parseDeclarations(s),
// 	}
// }

// func parseSelectors(s *scanner.Scanner) []Selector {
// 	var selectors []Selector

// 	for s.Peek() != scanner.EOF {
// 		selector := parseSelector(s)
// 		selectors = append(selectors, selector)
// 		if s.NextChar() == rune(',') {
// 			s.Scan()
// 			continue
// 		} else if s.NextChar() == rune('{') {
// 			s.Scan()
// 			break
// 		}
// 	}

// 	return selectors
// }

// func parseSelector(s *scanner.Scanner) Selector {
// 	selector := Selector{}

// 	for s.Peek() != scanner.EOF {
// 		if s.NextChar() == rune(',') || s.NextChar() == rune('{') {
// 			break
// 		}

// 		s.Scan()
// 		value := s.TokenText()

// 		if value == "#" {
// 			s.Scan()
// 			selector.ID = s.TokenText()
// 		} else if value == "." {
// 			s.Scan()
// 			selector.Classes = []string{s.TokenText()}
// 		} else {
// 			selector.TagName = value
// 		}
// 	}

// 	return selector
// }

// func parseDeclarations(s *scanner.Scanner) []Declaration {
// 	var declarations []Declaration

// 	for s.Peek() != scanner.EOF {
// 		if s.NextChar() == rune(';') {
// 			s.Scan()
// 			continue
// 		}
// 		if s.NextChar() == rune('}') {
// 			break
// 		}
// 		declarations = append(declarations, parseDeclaration(s))
// 	}

// 	return declarations
// }

// func parseDeclaration(s *scanner.Scanner) Declaration {
// 	identifier := parseIdentifier(s)
// 	value := parseValue(s)

// 	d := Declaration{
// 		Name:  identifier,
// 		Value: value,
// 	}

// 	return d
// }

// func parseIdentifier(s *scanner.Scanner) string {
// 	name := ""
// 	for s.Scan() != scanner.EOF {
// 		if s.TokenText() == ":" || s.TokenText() == ";" {
// 			break
// 		}
// 		name += s.TokenText()
// 	}
// 	return name
// }

// func parseValue(s *scanner.Scanner) Value {
// 	v := Value{}

// 	next := s.NextChar()

// 	if unicode.IsDigit(next) {
// 		v.Length = parseLength(s)
// 	} else if next == rune('#') {
// 		v.Color = parseColor(s)
// 	} else {
// 		keyword := parseIdentifier(s)
// 		switch keyword {
// 		case "red":
// 			v.Color = Color{A: 255, R: 255}
// 		case "blue":
// 			v.Color = Color{A: 255, B: 255}
// 		case "green":
// 			v.Color = Color{A: 255, G: 255}
// 		case "white":
// 			v.Color = Color{A: 255, R: 255, G: 255, B: 255}
// 		case "black":
// 			v.Color = Color{A: 255}
// 		default:
// 			v.Keyword = keyword
// 		}
// 	}

// 	return v
// }

// func parseLength(s *scanner.Scanner) Length {
// 	s.Scan()

// 	value := s.TokenText()
// 	f, err := strconv.ParseFloat(value, 32)
// 	if err != nil {
// 		f = 0
// 	}

// 	if s.NextChar() != rune(';') {
// 		s.Scan()
// 	}

// 	t := s.TokenText()
// 	var unit Unit
// 	switch t {
// 	case "px":
// 		unit = Px
// 	case "em":
// 		unit = Em
// 	case "%":
// 		unit = Percent
// 	}

// 	return Length{Quantity: float32(f), Unit: unit}
// }

// func parseColor(s *scanner.Scanner) Color {
// 	if s.NextChar() == rune('#') {
// 		s.Scan()
// 		s.Scan()
// 	}
// 	text := s.TokenText()

// 	color := Color{}
// 	if len(text) == 3 {
// 		color.A = 255
// 		color.R = hexToUint8(string(text[0]) + string(text[0]))
// 		color.G = hexToUint8(string(text[1]) + string(text[1]))
// 		color.B = hexToUint8(string(text[2]) + string(text[2]))
// 	} else if len(text) == 6 {
// 		color.A = 255
// 		color.R = hexToUint8(text[0:2])
// 		color.G = hexToUint8(text[2:4])
// 		color.B = hexToUint8(text[4:6])
// 	} else if len(text) == 8 {
// 		color.R = hexToUint8(text[0:2])
// 		color.G = hexToUint8(text[2:4])
// 		color.B = hexToUint8(text[4:6])
// 		color.A = hexToUint8(text[6:8])
// 	}
// 	return color
// }

// func hexToUint8(hex string) uint8 {
// 	val, err := strconv.ParseUint(hex, 16, 8)
// 	if err != nil {
// 		return 0
// 	}
// 	return uint8(val)
// }
