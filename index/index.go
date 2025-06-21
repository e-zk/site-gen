package index

import (
	"html/template"
	"log"
	"os"
	"sort"
	"time"

	"site-gen/post"
)

type Index struct {
	Title     string
	Templates []string
	List      []*post.Post
}

func New(posts []*post.Post) *Index {
	sort.Slice(posts, func(i, j int) bool {
		ti, err := time.Parse("2006-01-02", posts[i].Metadata.Date)
		if err != nil {
			log.Fatal(err)
		}
		tj, err := time.Parse("2006-01-02", posts[j].Metadata.Date)
		if err != nil {
			log.Fatal(err)
		}
		return ti.After(tj)
	})
	return &Index{List: posts}
}

func (data *Index) Execute(indexFile string) error {
	os.Create(indexFile)
	err := os.Truncate(indexFile, 0)
	if err != nil {
		return err
	}

	outputFile, err := os.OpenFile(indexFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	//t, err := template.ParseFiles("html/base.html", "html/posts.html")
	t, err := template.ParseFiles(data.Templates...)
	if err != nil {
		return err
	}

	err = t.Execute(outputFile, data)
	if err != nil {
		return err
	}
	return nil
}
