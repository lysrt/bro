package main

import (
	"bufio"
	"fmt"
	"os"
)

// Stylesheet represents a whole CSS file
type Stylesheet struct {
	Rules []Rule
}

// Rule represents a CSS block
type Rule struct {
	Selectors    []Selector
	Declatations []Declaration
}

// Selector represents a CSS selector, present before each CSS block
type Selector struct {
	TagName string
	ID      string
	Class   []string
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
	name  string
	value Value
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

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	sheet := parseStylesheet(scanner)

	return sheet, nil
}

func parseStylesheet(s *bufio.Scanner) *Stylesheet {
	fmt.Println("Parsing Stylesheet")

	return &Stylesheet{
		Rules: parseRules(s),
	}
}

func parseRules(s *bufio.Scanner) []Rule {
	fmt.Println("Parsing Rules")

	var rules []Rule
	for s.Err() == nil {
		rules = append(rules, parseRule(s))
	}
	if err := s.Err(); err != nil {
		//return nil, fmt.Errorf("error reading css: %q", err)
		fmt.Println("Error: ", err)
		return nil
	}
	return rules
}

func parseRule(s *bufio.Scanner) Rule {
	fmt.Println("\tParsing Rule")
	return Rule{
		Selectors:    parseSelectors(s),
		Declatations: parseDeclarations(s),
	}
}

func parseSelectors(s *bufio.Scanner) []Selector {
	fmt.Println("\t\tParsing Selectors")

	// Return when reading {
	for s.Scan() {
		text := s.Text()
		fmt.Println("Scanning: ", text)
		if text == "{" {
			break
		}
	}
	return nil
}

func parseDeclarations(s *bufio.Scanner) []Declaration {
	fmt.Println("\t\tParsing Declarations")

	// Return when reading }
	for s.Scan() {
		text := s.Text()
		fmt.Println("Scanning: ", text)
		if text == "}" {
			break
		}
	}
	return nil
}
