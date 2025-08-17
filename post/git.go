package post

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"time"
)

func getLastCommit(filepath string) (hash string, date string, err error) {
	r, err := git.PlainOpen("./.git")
	if err != nil {
		return hash, date, fmt.Errorf("error opening git repo: %w", err)
	}
	ref, err := r.Head()
	if err != nil {
		return hash, date, fmt.Errorf("error opening head: %w", err)
	}

	cIter, err := r.Log(&git.LogOptions{
		From:     ref.Hash(),
		FileName: &filepath,
	})
	if err != nil {
		return hash, date, fmt.Errorf("iter error: %w", err)
	}

	c, err := cIter.Next()
	if err != nil {
		return hash, date, fmt.Errorf("iter.next for %q, error: %w", filepath, err)
	}

	ftime := c.Author.When.Format(time.DateOnly)
	return c.Hash.String(), ftime, nil
}
