package main

import (
	"log"
)

const (
	hostname       = "zakaria.org"
	baseUrl        = "https://" + hostname
	onionUrl       = "http://64wv2uqwjacqer7z5d6ooqgrvjwlioizmo7hgmxm7zxerbvgnoqhafid.onion"
	postsDir       = "./posts"
	postIndex      = postsDir + "/index.html"
	postArchive    = postsDir + "/archive.html"
	footerTemplate = `<p><a href="https://creativecommons.org/licenses/by-sa/4.0/">&copy; CC BY-SA 4.0</a> <a href="{{.Plaintext}}">plaintext</a> <a href="{{.Onion}}">onion</a></p>`
)

var (
	archivePerma = []string{
		"https://zakaria.org/posts/2020-08-01-shblog.html",
		"https://zakaria.org/posts/2020-08-03-m4.html",
		"https://zakaria.org/posts/2020-09-09-tmux.html",
		"https://zakaria.org/posts/2020-11-07-malthusian-belt.html",
		"https://zakaria.org/posts/2020-12-05-fonts.html",
	}
)

func genAllPosts() {
	ps := getAllPosts("./")
	for _, p := range ps {
		log.Printf("%s => %s", p.MarkdownFile, p.OutPath())

		// prepare for template execution by conveting mardown => html
		p.ConvPost()

		// execute template to generate full .html
		p.Execute()
	}
}

func main() {
	//ps := getAllPosts(".")
	//for _, post := range ps {
	//	log.Printf("%s / %s / %s / %s ==> %s", post.MarkdownFile, post.Date, post.Title, post.Description, post.Permalink)
	//}

	log.Println("starting...")
	genAllPosts()
	log.Println("generating index...")
	genIndexSorted(postIndex)
	log.Println("generating archive index...")
	genArchiveSorted(postArchive)
}
