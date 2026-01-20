package post

import (
	"bytes"

	"ser1.net/godjot/v2/djot_html"
	"ser1.net/godjot/v2/djot_parser"
)

func djotToHTML(data []byte) []byte {
	ast := djot_parser.BuildDjotAst(data)
	content := djot_html.New().ConvertDjot(&djot_html.HtmlWriter{}, ast...)

	return bytes.NewBufferString(content.String()).Bytes()
}

// TODO verify this works ? dunno
func djotToXHTML(data []byte) []byte {
	ast := djot_parser.BuildDjotAst(data)
	content := djot_html.New().ConvertDjot(&djot_html.HtmlWriter{}, ast...)

	return bytes.NewBufferString(content.String()).Bytes()
}
