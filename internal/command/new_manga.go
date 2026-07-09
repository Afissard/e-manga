package command

import (
	"e-manga/internal/config"
	"e-manga/internal/library"
	"os"
)

type NewMangaOptions struct {
	Name        string
	Author      string
	Summary     string
	Cover       string
	Target      string
	URL         string
	LeftToRight bool
}

func NewManga(opts NewMangaOptions) error {
	if _, err := os.Stat(config.AppConfig.LibraryPath); os.IsNotExist(err) {
		return os.MkdirAll(config.AppConfig.LibraryPath, 0755)
	}

	manga := &library.Manga{}
	err := manga.Init(opts.Name, opts.Author, opts.Summary, opts.Cover, opts.Target, opts.URL, opts.LeftToRight)
	if err != nil {
		return err
	}

	/*
		err = scraper.ScrapDemo(url)
		if err != nil {
			return err
		}
	*/

	return nil
}
