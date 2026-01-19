package post

import (
	"bytes"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var MyExtensions = parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes

var MyRenderer = &html.Renderer{}
var MyOptions html.RendererOptions

var MyXMLRenderer = &html.Renderer{}
var MyXMLOptions html.RendererOptions

type Box struct {
	ast.Leaf
	Content []byte
}

var boxMatch = []byte(":box\n")
var boxMatchEnd = []byte("::\n")

func init() {
	htmlFlags := html.CommonFlags | html.FootnoteReturnLinks
	MyOptions = html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: boxRenderHook,
	}
	MyRenderer = html.NewRenderer(MyOptions)

	MyXMLOptions = html.RendererOptions{
		Flags:          html.UseXHTML,
		RenderNodeHook: boxRenderHookXML,
	}
	MyXMLRenderer = html.NewRenderer(MyXMLOptions)
}

func parseBox(data []byte) (ast.Node, []byte, int) {
	if !bytes.HasPrefix(data, boxMatch) {
		return nil, nil, 0
	}
	i := len(boxMatch)
	end := bytes.Index(data[i:], boxMatchEnd) + len(boxMatchEnd)
	if end < 0 {
		return nil, data, 0
	}
	end = end + i
	lines := data[i:end]
	res := &Box{
		// we want to chop off the trailing ending token
		Content: lines[:len(lines)-len(boxMatchEnd)],
	}
	return res, nil, end
}

func parserHook(data []byte) (ast.Node, []byte, int) {
	if node, d, n := parseBox(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

func boxRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if box, ok := node.(*Box); ok {
		if entering {
			io.WriteString(w, "<section class=\"box\">\n")
			// render markdown
			p := parser.NewWithExtensions(MyExtensions)
			p.Opts.ParserHook = parserHook
			renderer := *MyRenderer
			boxDoc := p.Parse(box.Content)
			boxHTML := markdown.Render(boxDoc, &renderer)
			io.WriteString(w, string(boxHTML))
			io.WriteString(w, "</section>\n\n")
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func boxRenderHookXML(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if box, ok := node.(*Box); ok {
		if entering {
			io.WriteString(w, "<section>\n")
			// render markdown
			p := parser.NewWithExtensions(MyExtensions)
			p.Opts.ParserHook = parserHook
			renderer := *MyXMLRenderer
			boxDoc := p.Parse(box.Content)
			boxHTML := markdown.Render(boxDoc, &renderer)
			io.WriteString(w, string(boxHTML))
			io.WriteString(w, "</section>\n\n")
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func mdToHTML(md []byte) []byte {
	p := parser.NewWithExtensions(MyExtensions)
	p.Opts.ParserHook = parserHook
	doc := p.Parse(md)

	renderer := *MyRenderer
	return markdown.Render(doc, &renderer)
}

func mdToXHTML(md []byte) []byte {
	p := parser.NewWithExtensions(MyExtensions)
	p.Opts.ParserHook = parserHook
	doc := p.Parse(md)

	renderer := *MyXMLRenderer
	return markdown.Render(doc, &renderer)
}
