package post

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// could make a custom marshall'er or use existing library but maybe later lol

var (
	ErrNoTitle error = errors.New("no title specified")
	ErrNoPath  error = errors.New("no path specified")
)

type PostMetadata struct {
	Title       string // title (mandatory)
	Path        string // markdown path
	Description string
	Date        string
	ImageURL    string
	Archived    bool
	Tags        []string // unused
	Type        string   // unused
}

func isStringBool(s string) bool {
	s = strings.ToLower(s)
	if s == "" {
		return false
	}
	if s == "true" {
		return true
	}
	return false
}

// parse a .meta file into a new post
func (p *Post) loadMetadata(path string) error {
	// give us the key:value pair
	parseKV := func(pair string) (k, v string) {
		k, v, _ = strings.Cut(pair, ":")
		return strings.TrimSpace(k), strings.TrimSpace(v)
	}

	m := PostMetadata{}

	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	sc := bufio.NewScanner(fp)
	for sc.Scan() {
		t := sc.Text()
		if strings.HasPrefix(t, "#") || !strings.Contains(t, ":") {
			continue
		}
		switch k, v := parseKV(t); k {
		case "path":
			m.Path = v
		case "title":
			m.Title = v
		case "description":
			m.Description = v
		case "date":
			m.Date = v
		case "image_url":
			m.ImageURL = v
		case "archived":
			m.Archived = isStringBool(v)
		default:
			// nothing
			// TODO specify tags + type (?)
		}
	}

	// if no path given assume its the same as .meta, but w/ .md ext
	if len(m.Path) == 0 {
		m.Path = strings.TrimSuffix(path, ".meta") + ".md"
	}
	//
	if len(m.Title) == 0 {
		return ErrNoTitle
	}

	p.Metadata = m

	return nil
}
