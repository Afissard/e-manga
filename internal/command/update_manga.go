package command

import "e-manga/internal/library"

type UpdateMangaMetadataOptions struct {
	Manga       string
	Author      string
	Summary     string
	Cover       string
	URL         string
	LeftToRight bool
}

func UpdateMangaMetadata(opts UpdateMangaMetadataOptions) error {
	manga, err := library.LoadManga(opts.Manga)
	if err != nil {
		return err
	}
	return manga.UpdateMetadata(opts.Author, opts.Summary, opts.Cover, opts.URL, opts.LeftToRight)
}
