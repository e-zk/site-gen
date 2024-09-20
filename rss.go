package main

import (
	"html/template"
	"log"
	"os"
	"sort"
	"time"
)

type rssData struct {
	Posts []*Post
}

func genRssFile(rssFile string) {
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
			p.ConvPost()
			d, _ := time.Parse("2006-01-02", p.Date)
			p.Date = d.Format(time.RFC1123Z)
			psa = append(psa, p)
		}
	}

	data := rssData{
		Posts: psa,
	}

	os.Create(rssFile)
	err := os.Truncate(rssFile, 0)
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.OpenFile(rssFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	t, err := template.ParseFiles("html/rss.xml")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(outputFile, data)
	if err != nil {
		log.Fatal(err)
	}
}
