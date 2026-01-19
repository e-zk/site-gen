package bookshelf

import (
	"html/template"
	"os"
	"sort"
	"strconv"

	"github.com/BurntSushi/toml"
)

type (
	Bookshelf struct {
		InfoMarkdown string `toml:"infomd"`
		Reading      []Book
		Books        map[string][]Book
		SortedKeys   []string `toml:"-"`
	}
	Book struct {
		Author     string
		Title      string
		Translator string
		Rating     int
	}
)

func New(config, output string, templates []string) error {
	var b Bookshelf
	_, err := toml.DecodeFile(config, &b)
	if err != nil {
		return err
	}
	//log.Printf("%+v", b.Reading)
	//log.Printf("%+v", b.InfoMarkdown)
	///for y, b := range b.Books {
	///	log.Printf("%v", y)
	///	log.Printf("%+v", b)
	///}
	keys := make([]int, 0, len(b.Books))
	for k, _ := range b.Books {
		ii, _ := strconv.Atoi(k)
		keys = append(keys, ii)
	}

	// we want zero to be last
	sort.Slice(keys, func(i, j int) bool {
		if keys[i] == 0 {
			return false
		} else if keys[j] == 0 {
			return true
		}
		if keys[i] > keys[j] {
			return true
		}
		return false
	})
	sortedKeys := make([]string, 0, len(b.Books))
	for _, k := range keys {
		is := strconv.Itoa(k)
		sortedKeys = append(sortedKeys, is)
	}
	b.SortedKeys = sortedKeys
	//log.Printf("%+v", sortedKeys)

	os.Create(output)
	err = os.Truncate(output, 0)
	if err != nil {
		return err
	}

	outputFile, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	t, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}
	err = t.Execute(outputFile, b)
	if err != nil {
		return err
	}
	return nil
}
