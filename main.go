package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	htmlIn := flag.String("html", "input.html", "-html input.html")
	cssIn := flag.String("css", "", "-css input.css")
	output := flag.String("o", "out.png", "-o out.png")
	flag.Parse()

	// 1. Constructing the DOM tree
	dom, err := ParseHTML(*htmlIn)
	if err != nil {
		log.Fatalf("cannot parse HTML file: %q", err)
	}

	// 2. Parsing the CSS to a *Stylesheet
	var css *Stylesheet
	if *cssIn != "" {
		css, err = ParseCSS(*cssIn)
		if err != nil {
			log.Fatalf("cannot parse CSS file: %q", err)
		}
	}

	Parcour(dom)
	fmt.Println(css)

	// 3. Decorating the DOM to generate the Style Tree
	styleTree, err := GenerateStyleTree(dom, css)
	if err != nil {
		log.Fatalf("cannot build style tree: %q", err)
	}

	// 4. Build the Layout Tree
	layoutTree, err := GenerateLayoutTree(styleTree)
	if err != nil {
		log.Fatalf("cannot build layout tree: %q", err)
	}

	// 5. Paint the output from the Layout Tree
	pixels, err := Paint(layoutTree)
	if err != nil {
		log.Fatalf("cannot paint from layout tree: %q", err)
	}

	writeOutput(*output, pixels)
}

func writeOutput(outputFileName string, pixels image.Image) {
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("cannot open output file: %q", err)
	}

	if err := png.Encode(f, pixels); err != nil {
		f.Close()
		log.Fatalf("cannot encode PNG: %q", err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("cannot close output file: %q", err)
	}
}
