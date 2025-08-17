package post

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var myExtensions = parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes

func mdToHTML(md []byte) []byte {
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
