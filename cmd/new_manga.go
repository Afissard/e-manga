package main

import (
	"e-manga/internal/config"
	"e-manga/internal/library"
	"os"
)

func NewManga(name, author, summary, cover, target, url string, leftToRight bool) error {
	if _, err := os.Stat(config.AppConfig.LibraryPath); os.IsNotExist(err) {
		return os.MkdirAll(config.AppConfig.LibraryPath, 0755)
	}

	manga := &library.Manga{}
	err := manga.Init(name, author, summary, cover, target, url, leftToRight)
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
