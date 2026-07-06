package processing

import (
	"e-manga/internal/config"
)

type Chapter struct {
	Name   string
	Images []string
}

type Manga struct {
	Title    string
	Chapters []Chapter
	// TODO: add metadata fields like author, genre, etc.
}

type Options struct {
	Target config.Target
	AutoRotate bool
}
