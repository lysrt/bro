package main

import (
	"io/ioutil"
)

// Stylesheet represents a whole CSS file
type Stylesheet struct {
	Rules []Rule
}

// Rule represents a CSS block
type Rule struct {
	Selector     Selector
	Declatations []Declaration
}

// Selector represents a CSS selector, present before each CSS block
type Selector struct {
	TagName string
	ID      string
	Class   string
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
	A, R, G, B int
}

func ParseCSS(inputFileName string) (*Stylesheet, error) {
	css, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		return nil, err
	}

	sheet := &CSSParser{string(css)}

	return sheet, nil
}

func parse()

func parseSimpleSelector(buffer) Selector {

}
