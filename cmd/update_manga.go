package main

import "e-manga/internal/library"

func UpdateMangaMetadata(mangaName string, author string, summary string, cover string, target string, url string, leftToRight bool) error {
	manga, err := library.LoadManga(mangaName)
	if err != nil {
		return err
	}
	return manga.UpdateMetadata(author, summary, cover, target, url, leftToRight)
}
