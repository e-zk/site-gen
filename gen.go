package main

// template functions
// todo: add opengraph tag functionality

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Index struct {
	List   []*Post
	After  template.HTML
	Before template.HTML
	Footer template.HTML
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.FootnoteReturnLinks
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// func genIndex(articles []*Post, indexFile string) {
func genIndex(data Index, indexFile string) {
	//err := os.Truncate(indexFile, 0)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//outputFile, err := os.OpenFile(indexFile, os.O_WRONLY|os.O_CREATE, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer outputFile.Close()
	outputFile, err := os.Create(indexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	t, err := template.ParseFiles("html/base.html", "html/posts.html")
	if err != nil {
		log.Fatal(err)
	}

	//data := Index{
	//	List:   articles,
	//	Before: template.HTML(`<p>Writings...</p>`),
	//	After:  template.HTML(`<p><a href="archive.html">&laquo; archive</a></p>`),
	//	Footer: template.HTML(`<p><a href="https://creativecommons.org/licenses/by-sa/4.0/">&copy; CC BY-SA 4.0</a></p>`),
	//}

	err = t.Execute(outputFile, data)
	if err != nil {
		log.Fatal(err)
	}
}

// convert markdown => html for post
// and generate footer as well
func convPost(p *Post) {
	fc, err := os.ReadFile(p.MarkdownFile)
	if err != nil {
		log.Fatal(err)
	}

	// TODO this relative path may need changing
	md := "./" + path.Base(p.MarkdownFile)

	// convert markdown to html
	p.Content = template.HTML(mdToHTML(fc))

	// struct for footer template
	footerData := struct {
		Plaintext string
		Onion     string
	}{
		Plaintext: md,
		Onion:     p.Onionlink,
	}

	// generate footer
	var out bytes.Buffer
	t, _ := template.New("footer").Parse(footerTemplate)
	err = t.Execute(&out, footerData)

	p.Footer = template.HTML(out.String())

}

func genPost(p *Post) {
	err := os.Truncate(p.OutPath(), 0)
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.OpenFile(p.OutPath(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles("html/base.html", "html/post.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(outputFile, p)
	if err != nil {
		log.Fatal(err)
	}
}
