package main

import (
	"bufio"
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
	// Meta tags
	Title       string
	Description string
	Date        string

	// These we add ourselves
	MarkdownFile string
	Rellink      string
	Permalink    string
	Onionlink    string

	// these get generated
	Content template.HTML
	Footer  template.HTML
}

func (p *Post) OutPath() string {
	return strings.TrimSuffix(p.MarkdownFile, ".md") + ".html"
}

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

	// walk our directory tree and add all posts
	rfs := os.DirFS(basedir)
	fs.WalkDir(rfs, ".", func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		// TODO check if meta has md assoc. with it

		if filepath.Ext(fpath) == ".meta" {
			fullpath := path.Join(basedir, fpath)
			perma := baseUrl + "/" + strings.TrimSuffix(fullpath, ".meta") + ".html"
			onion := strings.Replace(perma, baseUrl, onionUrl, -1)

			// associated markdown file
			mdPath := strings.TrimSuffix(fullpath, ".meta") + ".md"
			if _, err := os.Stat(mdPath); errors.Is(err, os.ErrNotExist) {
				log.Printf("%s: has no markdown associated with it - ignoring")
				return nil
			}

			// parse meta file into new post
			p := parseMeta(fullpath)
			p.MarkdownFile = mdPath
			p.Permalink = perma
			p.Onionlink = onion
			p.Rellink = strings.TrimSuffix(fpath, ".meta") + ".html"
			ps = append(ps, p)
		}
		return nil
	})

	return ps
}
