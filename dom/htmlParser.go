package dom

import (
	"io"

	"golang.org/x/net/html"
)

func ParseHTML(r io.Reader) (*html.Node, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return node, nil
}
