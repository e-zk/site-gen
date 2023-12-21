package main

import (
	"html/template"
	"log"
	"sort"
	"time"
)

func inArchive(p *Post) bool {
	archive := false
	for _, link := range archivePerma {
		if p.Permalink == link {
			archive = true
		}
	}
	return archive
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
		if !inArchive(p) {
			psa = append(psa, p)
		}
	}

	data := Index{
		List:   psa,
		Before: template.HTML("<p>Writings...</p>"),
		After:  template.HTML(`<p><a href="./archive.html">&laquo; archive</a></p>`),
		Footer: template.HTML(`<p><a href="https://creativecommons.org/licenses/by-sa/4.0/">&copy; CC BY-SA 4.0</a></p>`),
	}

	genIndex(data, indexPath)
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
		if inArchive(p) {
			psa = append(psa, p)
		}
	}

	data := Index{
		List:   psa,
		Before: template.HTML("<p>Old/archived posts.</p>"),
		After:  template.HTML(`<p><a href="/posts">&laquo; main posts</a></p>`),
		Footer: template.HTML(`<p><a href="https://creativecommons.org/licenses/by-sa/4.0/">&copy; CC BY-SA 4.0</a></p>`),
	}

	genIndex(data, archivePath)
}
