package main

import (
//	"regexp"
//	"log"
//	"bytes"
//	"fmt"
//	"io"

	"github.com/gomarkdown/markdown"
//	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var myExtensions = parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes

// if not already linked-to, link to openbsd man page
/*
const manr = "([\\w-]+[^]`])\\(([0-9])\\)"
func preproc(data []byte) []byte {
	re := regexp.MustCompile(manr)
	return re.ReplaceAll(
		data,
		[]byte(`[$1($2)](https://man.openbsd.org/$1 "$1 man page")`),
	)
}
*/


func mdToHTML(md []byte) []byte {
	//md = preproc(md)
	p := parser.NewWithExtensions(myExtensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.FootnoteReturnLinks
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func mdToXHTML(md []byte) []byte {
	p := parser.NewWithExtensions(myExtensions)
	doc := p.Parse(md)

	htmlFlags := html.UseXHTML
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer)
}

