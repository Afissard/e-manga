package main

import (
	"e-manga/internal/config"
	"e-manga/internal/library"
	"os"
)

func NewManga(name string) error {
	if _, err := os.Stat(config.AppConfig.LibraryPath); os.IsNotExist(err) {
		return os.MkdirAll(config.AppConfig.LibraryPath, 0755)
	}

	manga := &library.Manga{}
	err := manga.Init(name)
	if err != nil {
		return err
	}

	return nil
}
