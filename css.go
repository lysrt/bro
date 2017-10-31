package main

import (
	"fmt"
	"io/ioutil"
)

type CSSParser struct {
	input string
}

func newCSSParser(inputFileName string) (*CSSParser, error) {
	css, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		return nil, err
	}

	parser := &CSSParser{string(css)}

	return parser, nil
}

func (p *CSSParser) Parse() (*CSS, error) {
	fmt.Println("Parse Css")
	fmt.Print(p.input)

	return nil, nil
}
