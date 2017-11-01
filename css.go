package main

import (
	"fmt"
	"os"
	"strings"
	"text/scanner"
	"unicode"
)

// Stylesheet represents a whole CSS file
type Stylesheet struct {
	Rules []Rule
}

func (s Stylesheet) String() string {
	r := "Spreadsheet\n"
	for _, rule := range s.Rules {
		r += fmt.Sprint(rule)
	}
	return r
}

// Rule represents a CSS block
type Rule struct {
	Selectors    []Selector
	Declarations []Declaration
}

func (s Rule) String() string {
	r := " Rule\n"
	for _, selector := range s.Selectors {
		r += fmt.Sprintf("%v, ", selector)
	}
	r += fmt.Sprintln()
	for _, declaration := range s.Declarations {
		r += fmt.Sprintf("%v\n", declaration)
	}
	return r
}

// Selector represents a CSS selector, present before each CSS block
type Selector struct {
	TagName string
	ID      string
	Classes []string
}

func (s Selector) String() string {
	r := "  Selector ("
	if s.TagName != "" {
		r += fmt.Sprintf("TAG: %s)", s.TagName)
	} else if s.ID != "" {
		r += fmt.Sprintf("ID: %s)", s.ID)
	} else {
		r += fmt.Sprintf("CLASSES: [%v])", strings.Join(s.Classes, ", "))
	}
	return r
}

// Specificity computes and returns the specificity of a selector.
func (s *Selector) Specificity() Specificity {
	return Specificity{
		A: len(s.ID),
		B: len(s.Class),
		C: len(s.TagName),
	}
}

// Declaration represents a single CSS property
type Declaration struct {
	Name  string
	Value Value
}

func (s Declaration) String() string {
	r := "  Declaration: "
	if s.Name != "" {
		r += fmt.Sprint(s.Name)
		r += fmt.Sprint(s.Value)
	} else {
		r += "No Name..."
	}
	return r
}

// Value represents the possible value of a CSS declaration
type Value struct {
	Keyword string
	Length  Length
	Color   Color
}

// Length describes a unit of length in CSS
type Length struct {
	Quantity float32
	Unit     Unit
}

type Unit int

const (
	Px Unit = iota + 1
	Em
	Percent
)

type Color struct {
	Name       string
	A, R, G, B uint8
}

// Specificity represents the specificity of a Rule.
// It is use to compute rule precedence.
type Specificity struct {
	A, B, C int
}

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
	// TODO Implement #ABCDEF
	s.Scan()
	return Color{A: 255, R: 255, G: 0, B: 0}
}
