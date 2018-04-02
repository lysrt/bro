package main

import (
	"flag"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/lysrt/bro/css"
	"github.com/lysrt/bro/html"
	"github.com/lysrt/bro/html/lexer"
	"github.com/lysrt/bro/html/parser"
	"github.com/lysrt/bro/layout"
	"github.com/lysrt/bro/paint"
	"github.com/lysrt/bro/style"
)

func main() {
	var (
		htmlInput string
		cssInput  string
		pngOutput string
	)

	flag.StringVar(&htmlInput, "html", "input.html", "-html input.html")
	flag.StringVar(&cssInput, "css", "input.css", "-css input.css")
	flag.StringVar(&pngOutput, "o", "out.png", "-o out.png")
	flag.Parse()

	var (
		domNodes   *html.Node
		styleSheet *css.Stylesheet
	)

	//
	// 1. Construct the DOM tree
	//
	htmlFile, err := os.Open(htmlInput)
	if err != nil {
		log.Fatalf("cannot open HTML file: %q", err)
	}

	// domNodes, err = html.Parse(htmlFile)
	b, err := ioutil.ReadAll(htmlFile)
	if err != nil {
		log.Fatalf("cannot read HTML gile: %q", err)
	}

	l := lexer.New(string(b))
	p := parser.New(l)
	domNodes = p.Parse()
	if len(p.Errors()) > 0 {
		log.Println("cannot parse HTML file")
		for _, e := range p.Errors() {
			log.Printf("%q (l: %d, c: %d)\n", e, e.Token.Line, e.Token.LinePosition)
		}
		log.Fatal()
	}
	htmlFile.Close()
	// dom.Parcour(domNodes)

	//
	// 2. Parse the CSS to a *Stylesheet
	//
	cssFile, err := os.Open(cssInput)
	if err != nil {
		log.Fatalf("cannot open stylesheet: %q", err)
	}

	parser := css.NewParser(cssFile)
	cssFile.Close()
	styleSheet = parser.ParseStylesheet()
	if len(parser.Errors()) > 0 {
		for _, e := range parser.Errors() {
			log.Printf("parsing error: %q\n", e)
		}
	}
	// fmt.Println(styleSheet)

	//
	// 3. Decorate the DOM to generate the Style Tree
	//
	styleTree := style.GenerateStyleTree(domNodes, styleSheet)
	//fmt.Println(styleTree)

	//
	// 4.1 Build the Layout Tree
	//
	layoutTree := layout.GenerateLayoutTree(styleTree)
	// fmt.Println(layoutTree)

	//
	// 4.2 Parcour the layout tree to compute boxes dimensions
	//
	// Height must be zero here!
	layoutTree.Layout(layout.Dimensions{Content: layout.Rect{X: 0, Y: 0, Width: 300, Height: 0}})
	// fmt.Println(layoutTree)

	//
	// 5. Paint the output from the Layout Tree
	//
	pixels, err := paint.Paint(layoutTree)
	if err != nil {
		log.Fatalf("cannot paint from layout tree: %q", err)
	}
	//fmt.Println(pixels)

	writeOutput(pngOutput, pixels)
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
