package post

import (
	"os"
	"strings"
)

func trimAllSuffix(s string, suffixes []string) string {
	for _, trim := range suffixes {
		s = strings.TrimSuffix(s, trim)
	}
	return s
}
func hasAnySuffix(s string, suffixes []string) bool {
	for _, suf := range suffixes {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}
	log.Printf("file %q may or may not exist; this path should not be reached", path)
	return false
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
