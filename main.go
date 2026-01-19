package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"site-gen/bookshelf"
	"site-gen/index"
	"site-gen/post"
)

const (
	FQDN    = "zakaria.org"
	BaseURL = "https://" + FQDN + "/"
)

func getAllPosts(basedir string) []*post.Post {
	ps := make([]*post.Post, 0)

	// walk our directory tree and add all posts (.meta) we can find
	rfs := os.DirFS(basedir)
	fs.WalkDir(rfs, ".", func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if filepath.Ext(fpath) == ".meta" {
			fullPath := filepath.Join(basedir, fpath)
			if strings.HasPrefix(filepath.Base(fpath), ".") {
				// skip hidden files
				log.Printf(">> skipping %q (hidden)", fullPath)
				return nil
			}
			log.Printf(">> found %q", fullPath)
			p, err := post.New(fullPath)
			if err != nil {
				log.Fatal(err)
			}

			p.BlogFQDN = FQDN
			p.BlogBaseURL = BaseURL

			ps = append(ps, p)
		}
		return nil
	})
	return ps
}

func compile() {
	log.Printf("> compiling...")
	posts := getAllPosts(".")

	// compile all posts
	for _, p := range posts {
		log.Printf(">> loading post %q", p.Metadata.Title)

		err := p.Parse(false)
		if err != nil {
			log.Fatal(err)
		}

		err = p.Execute("./html/base.html", "./html/post.html")
		if err != nil {
			log.Fatal(err)
		}
	}

	// generate index (only from ./posts)
	// todo make this callable via arguments:
	// site-gen index \
	// 	-path "./posts" -title "Web log" \
	// 	-template "./html/posts.html" -output "./pots/index.html"
}

func indexAll() {
	log.Printf("> indexing...")
	blogPosts := getAllPosts("./posts")
	blogIndex := index.New(blogPosts)
	blogIndex.Title = "Web log"
	blogIndex.Templates = []string{"html/base.html", "html/posts.html"}
	err := blogIndex.Execute("./posts/index.html")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("> generating rss...")
	rssIndex := index.New(blogPosts)

	// compile+change date format afterwards (idk if this works)
	for _, p := range rssIndex.List {
		log.Printf(">> compiling post as xhtml %q", p.Metadata.Title)
		d, err := time.Parse("2006-01-02", p.Metadata.Date)
		if err != nil {
			log.Fatal(err)
		}
		p.Metadata.Date = d.Format(time.RFC1123Z)
		err = p.Parse(true)
		if err != nil {
			log.Fatal(err)
		}
	}
	rssIndex.Title = "Web log" // don't even need this
	rssIndex.Templates = []string{"html/rss.xml"}
	err = rssIndex.Execute("./rss.xml")
	if err != nil {
		log.Fatal(err)
	}
}

func books() {
	err := bookshelf.New(
		"./bookshelf/bookshelf.toml",
		"./bookshelf/index.html",
		[]string{"html/base.html", "html/bookshelf.html"},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) == 0 {
		compile()
		indexAll()
		books()
	}
	switch os.Args[1] {
	case "compile":
		compile()
	case "index":
		indexAll()
	case "bookshelf":
		books()
	default:
		compile()
		indexAll()
		books()
	}
	// would be cool not worth implementing since i only have one index page:
	// site-gen index <-dir path> [-output file] [-title "index"]
	// where <dir> dir containing posts to index
	//       [output] is the generated index file (defaults to <dir>/index.html)
}
