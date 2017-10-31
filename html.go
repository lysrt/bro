package main

import (
	"os"

	"golang.org/x/net/html"
)

func ParseHTML(name string) (*html.Node, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	node, err := html.Parse(f)
	if err != nil {
		return nil, err
	}

	return node, nil
}
