package main

import (
	"bufio"
	"bytes"
	"errors"
	"html"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

type Post struct {
	// Taken from meta file
	Title       string
	Description string
	Date        string
	Modified    string

	// These we add ourselves
	MarkdownFile     string
	MarkdownFileHash string
	HtmlFile         string
	Rellink          string
	Permalink        string
	Onionlink        string
	Archived         bool

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
	os.Create(p.OutPath())
	err := os.Truncate(p.OutPath(), 0)
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.OpenFile(p.OutPath(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

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
func (p *Post) ConvPost(xhtml bool) {
	fc, err := os.ReadFile(p.MarkdownFile)
	if err != nil {
		log.Fatal(err)
	}

	// TODO this relative path may need changing
	md := "./" + path.Base(p.MarkdownFile)

	// convert markdown to html
	var content []byte
	if xhtml {
		content = mdToXHTML(fc)
	} else {
		content = mdToHTML(fc)
	}
	contentStr := string(content[:])
	if xhtml {
		contentStr = html.EscapeString(contentStr)
	}
	p.Content = template.HTML(contentStr)

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

	outStr := out.String()
	//if xhtml {
	//	outStr = html.EscapeString(outStr)
	//}
	p.Footer = template.HTML(outStr)
}

// parse a .meta file into a new post
func newPostFromMeta(fpath string) *Post {
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
			// stat
			_, err := os.Stat(mdPath)
			if errors.Is(err, os.ErrNotExist) {
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
			p := newPostFromMeta(fullpath)

			if p.Date != "" {
				r, err := git.PlainOpen("./.git")
				if err != nil {
					log.Println("err opening .git:")
					log.Fatal(err)
				}
				ref, err := r.Head()
				if err != nil {
					log.Println("err getting head:")
					log.Fatal(err)
				}

				cIter, err := r.Log(&git.LogOptions{
					From:     ref.Hash(),
					FileName: &mdPath,
				})
				if err != nil {
					log.Println("iter err:")
					log.Fatal(err)
				}

				//err = cIter.ForEach(func(c *object.Commit) error {
				//	fmt.Println(c)
				//	return nil
				//})
				c, err := cIter.Next()
				//if errors.Is(err, errors.New("EOF")) {
				//
				//} else if err != nil {
				if err != nil {
					log.Fatalf("iter.next err for %q: %v", mdPath, err)
				}

				ctime := c.Author.When.Format(time.DateOnly)
				log.Printf("got lmod for %q: %s (%s)", mdPath, ctime, c.Hash.String())
				p.Modified = ctime
				p.MarkdownFileHash = c.Hash.String()
			}
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
