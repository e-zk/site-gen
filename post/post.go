package post

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

// converter function type def (see md.go + djot.go for implementations)
type ConvFunc func(data []byte) []byte

type Post struct {
	Metadata PostMetadata

	BlogFQDN    string
	BlogBaseURL string

	// generated
	HTMLPath     string
	GitHash      string
	LastModified string
	PermanentURL string // We don't even need this? It can be calculated inside the templates
	RelativeURL  string

	Content template.HTML
}

// func New(path, baseURL string) (p *Post, err error) {
func New(metafile string) (p *Post, err error) {
	p = new(Post)
	err = p.loadMetadata(metafile)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(p.Metadata.Path); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%q: plaintext path does not exist!", metafile)
	}

	p.HTMLPath = strings.TrimSuffix(p.Metadata.Path, ".md") + ".html"
	p.RelativeURL = strings.TrimLeft(p.HTMLPath, "./") // cuz sometimes its prefixed w/ .

	// if we have a publish date, we probably want a last modified date.
	// we use the last git commit on the markdown file for that.
	// fails silently.
	if p.Metadata.Date != "" {
		// ignore error just leave them blank
		hash, modified, err := getLastCommit(p.Metadata.Path)
		//debugging: log.Printf("got hash: %v", hash)
		if err != nil {
			log.Printf("%w (ignored)", err)
		} else {
			p.GitHash = hash
			p.LastModified = modified
		}
	}

	return p, nil
}

func (p *Post) Parse(xhtml bool) error {
	fc, err := os.ReadFile(p.Metadata.Path)
	if err != nil {
		return err
	}

	var htmlContent string
	//var converter ConvFunc

	if strings.HasSuffix(p.Metadata.Path, ".djot") {
		if xhtml {
			htmlContent = string(djotToXHTML(fc[:]))
			htmlContent = template.HTMLEscapeString(htmlContent)
		} else {
			htmlContent = string(djotToHTML(fc[:]))
		}
	} else {
		if xhtml {
			htmlContent = string(mdToXHTML(fc[:]))
			htmlContent = template.HTMLEscapeString(htmlContent)
		} else {
			htmlContent = string(mdToHTML(fc[:]))
		}
	}

	p.Content = template.HTML(htmlContent)

	return nil
}

// we should compile the markdown here too
func (p *Post) Execute(filenames ...string) error {

	// create + truncate (flush)
	os.Create(p.HTMLPath)
	err := os.Truncate(p.HTMLPath, 0)
	if err != nil {
		return err
	}

	// open
	htmlFile, err := os.OpenFile(p.HTMLPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	t, err := template.ParseFiles(filenames...)
	if err != nil {
		return err
	}

	// TODO write to html file
	if err = t.Execute(htmlFile, p); err != nil {
		return err
	}

	return nil
}
