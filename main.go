package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/net/html"
)

func main() {
	htmlIn := flag.String("h", "in.html", "-h in.html")
	cssIn := flag.String("c", "in.css", "-c in.css")
	// output := flag.String("o", "out", "-o out")
	flag.Parse()

	dom := parseHTML(*htmlIn)

	css := parseCSS(*cssIn)

	/*layout := */
	generateLayout(dom, css)

	// writeOutput(*output, layout)
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

	fmt.Printf("CSS: %v\n", css)

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

	fmt.Printf("Successfully wrote to %s\n", outputFileName)
}
