package main

import (
	"html/template"
	"log"
	"os"
	"sort"
	"time"
)

type Index struct {
	List   []*Post
	Title  string
	After  template.HTML
	Before template.HTML
	Footer template.HTML
}

func (data *Index) Execute(indexFile string) {
	outputFile, err := os.Create(indexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	t, err := template.ParseFiles("html/base.html", "html/posts.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(outputFile, data)
	if err != nil {
		log.Fatal(err)
	}
}

func genIndexSorted(indexPath string) {
	ps := getAllPosts(postsDir)
	sort.Slice(ps, func(i, j int) bool {
		ti, err := time.Parse("2006-01-02", ps[i].Date)
		if err != nil {
			log.Fatal(err)
		}
		tj, err := time.Parse("2006-01-02", ps[j].Date)
		if err != nil {
			log.Fatal(err)
		}
		return ti.After(tj)
	})

	psa := make([]*Post, 0)
	for _, p := range ps {
		if !p.Archived {
			psa = append(psa, p)
		}
	}

	// TODO: make this configurable somewhere / not hard-coded
	data := Index{
		List:   psa,
		Title:  "Web log",
		Before: template.HTML(`<p>Writings... Subscribe to the <a href="/rss.xml">rss feed.</a></p>`),
		After:  template.HTML(`<p><a href="./archive.html">&laquo; archive</a></p>`),
		Footer: template.HTML(`<p><a href="https://creativecommons.org/licenses/by-sa/4.0/">&copy; CC BY-SA 4.0</a></p>`),
	}

	data.Execute(indexPath)
}

func genArchiveSorted(archivePath string) {
	ps := getAllPosts(postsDir)
	sort.Slice(ps, func(i, j int) bool {
		ti, err := time.Parse("2006-01-02", ps[i].Date)
		if err != nil {
			log.Fatal(err)
		}
		tj, err := time.Parse("2006-01-02", ps[j].Date)
		if err != nil {
			log.Fatal(err)
		}
		return ti.After(tj)
	})

	psa := make([]*Post, 0)
	for _, p := range ps {
		if p.Archived {
			psa = append(psa, p)
		}
	}

	data := Index{
		List:   psa,
		Title:  "Archived posts",
		Before: template.HTML(`<p>Archived posts.</p>`),
		After:  template.HTML(`<p><a href="/posts">main posts &raquo;</a></p>`),
		Footer: template.HTML(`<p><a href="https://creativecommons.org/licenses/by-sa/4.0/">&copy; CC BY-SA 4.0</a></p>`),
	}

	data.Execute(archivePath)
}
