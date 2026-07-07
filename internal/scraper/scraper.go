package scraper

import (
	"e-manga/internal/scraper/model"
	"e-manga/internal/scraper/site/bigsolo"
	"fmt"
)

type Scraper interface {
	GetManga(url string) (*model.Manga, error)
	GetChapter(url string) (*model.Chapter, error)
}


func ScrapDemo(url string) error {
	scraper := bigsolo.New()
	manga, err := scraper.GetManga(url)
	if err != nil {
		return err
	}
	
	fmt.Println(manga.Title)
	for _, chapter := range manga.Chapters {
		fmt.Println(chapter.Title)
	}
	return nil
}