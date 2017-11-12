package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/lysrt/bro/css"
	"github.com/lysrt/bro/dom"
)

func main() {
	htmlIn := flag.String("html", "input.html", "-html input.html")
	cssIn := flag.String("css", "", "-css input.css")
	output := flag.String("o", "out.png", "-o out.png")
	flag.Parse()

	// 1. Constructing the DOM tree
	d, err := dom.ParseHTML(*htmlIn)
	if err != nil {
		log.Fatalf("cannot parse HTML file: %q", err)
	}

	// 2. Parsing the CSS to a *Stylesheet
	var s *css.Stylesheet
	if *cssIn != "" {
		f, err := os.Open(*cssIn)
		if err != nil {
			log.Fatal("fail to open stylesheet:", err)
		}
		s = css.ParseCSS(f)
	}

	//dom.Parcour(d)
	//fmt.Println(s)

	// 3. Decorating the DOM to generate the Style Tree
	styleTree := GenerateStyleTree(d, s)
	// if err != nil {
	// log.Fatalf("cannot build style tree: %q", err)
	// }

	//fmt.Println(styleTree)

	// 4.1 Build the Layout Tree
	layoutTree := GenerateLayoutTree(styleTree)
	// if err != nil {
	// 	log.Fatalf("cannot build layout tree: %q", err)
	// }

	fmt.Println(layoutTree)

	// 4.2 Parcour the layout tree to compute boxes dimensions
	layoutTree.Layout(Dimensions{content: Rect{0, 0, 200, 100}})
	fmt.Println(layoutTree)

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
