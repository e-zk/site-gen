package main

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Post struct {
	// Taken from meta file
	Title       string
	Description string
	Date        string

	// These we add ourselves
	MarkdownFile string
	HtmlFile     string
	Rellink      string
	Permalink    string
	Onionlink    string
	Archived     bool

	// these get generated
	Content template.HTML
	Footer  template.HTML
}

// TODO remove: this should not be needed anymore
func (p *Post) OutPath() string {
	return strings.TrimSuffix(p.MarkdownFile, ".md") + ".html"
}

// TODO remove: this should not be needed anymore
func (p *Post) IsArchived() bool {
	archive := false
	for _, link := range archivePerma {
		if p.Permalink == link {
			archive = true
		}
	}
	return archive
}

// execute post template
func (p *Post) Execute() {
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

// convert markdown => html for post
// and generate footer as well
func (p *Post) ConvPost() {
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

// parse a .meta file into a new post
func parseMeta(fpath string) *Post {
	parseKV := func(pair string) (k, v string) {
		k, v, _ = strings.Cut(pair, ":")
		return strings.Trim(k, " "), strings.Trim(v, " ")
	}

	p := new(Post)

	fp, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	sc := bufio.NewScanner(fp)
	for sc.Scan() {
		t := sc.Text()
		if strings.HasPrefix(t, "#") || !strings.Contains(t, ":") {
			continue
		}
		switch k, v := parseKV(t); k {
		case "title":
			p.Title = v
		case "description":
			p.Description = v
		case "date":
			p.Date = v
		default:
			// nothing
		}

	}

	return p
}

// output list of posts from $postsDir
func getAllPosts(basedir string) []*Post {
	ps := make([]*Post, 0)

	// walk our directory tree and add all posts (.meta) we can find
	rfs := os.DirFS(basedir)
	fs.WalkDir(rfs, ".", func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if filepath.Ext(fpath) == ".meta" {

			// construct links
			fullpath := path.Join(basedir, fpath)
			rel := strings.TrimSuffix(fpath, ".meta") + ".html"
			perma := baseUrl + "/" + strings.TrimSuffix(fullpath, ".meta") + ".html"
			onion := strings.Replace(perma, baseUrl, onionUrl, -1)

			// associated markdown file
			mdPath := strings.TrimSuffix(fullpath, ".meta") + ".md"
			if _, err := os.Stat(mdPath); errors.Is(err, os.ErrNotExist) {
				log.Printf("%s: has no markdown associated with it - ignoring", fullpath)
				return nil
			}

			// html output will be same as mdPath but with .html extension
			htmlPath := strings.TrimSuffix(mdPath, ".md") + ".html"

			// is it in the list of archived posts?
			archive := false
			for _, link := range archivePerma {
				if perma == link {
					archive = true
				}
			}

			// parse meta file into new post
			p := parseMeta(fullpath)

			p.MarkdownFile = mdPath
			p.HtmlFile = htmlPath
			p.Permalink = perma
			p.Onionlink = onion
			p.Rellink = rel
			p.Archived = archive

			ps = append(ps, p)
		}
		return nil
	})

	return ps
}
