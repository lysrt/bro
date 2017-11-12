package css

import (
	"fmt"
	"strings"
)

// Stylesheet represents a whole CSS file
type Stylesheet struct {
	Rules []Rule
}

func (s Stylesheet) String() string {
	r := "Stylesheet\n"
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
		B: len(s.Classes),
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

// ToPx is a Helper method needed by layout.go to get the actual pixel value
func (v Value) ToPx() float32 {
	if v.Length.Unit == Px {
		return v.Length.Quantity
	}
	return 0.0
}

// Length describes a unit of length in CSS
type Length struct {
	Quantity float32
	Unit     Unit
}

type Unit string

const (
	Px      Unit = "px"
	Em           = "em"
	Percent      = "pc"
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
