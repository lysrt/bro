package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/net/html"
)

func main() {
	htmlIn := flag.String("html", "input.html", "-html input.html")
	cssIn := flag.String("css", "input.css", "-css input.css")
	output := flag.String("o", "out", "-o out")
	flag.Parse()

	// 1. Constructing the DOM tree
	dom := parseHTML(*htmlIn)

	// 2. Parsing the CSS to a *Stylesheet
	css := parseCSS(*cssIn)

	// 3. Decorating the DOM to generate the Render Tree
	layout := generateLayout(dom, css)

	// 4. Building the output from the Render Tree
	writeOutput(*output, layout)
}

func parseHTML(inputFileName string) *html.Node {
	n, err := ParseHTML(inputFileName)
	if err != nil {
		log.Fatalf("cannot parse HTML file: %q", err)
	}

	Parcour(n)

	return n
}

func parseCSS(inputFileName string) *Stylesheet {
	css, err := ParseCSS(inputFileName)
	if err != nil {
		log.Fatalf("cannot parse CSS file: %q", err)
	}

	fmt.Println(css)

	return css
}

func generateLayout(dom *html.Node, css *Stylesheet) []byte {
	return []byte("")
}

func writeOutput(outputFileName string, layout []byte) {
	err := ioutil.WriteFile(outputFileName, layout, 0644)
	if err != nil {
		log.Fatalf("cannot write to output file: %q", err)
	}
}
