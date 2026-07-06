package library

import (
	"e-manga/internal/config"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func (m *Manga) Init(title string) error {
	m.Title = title
	m.Chapters = []Chapter{}

	dirs := []string{
		m.SourceDir(),
		m.CacheDir(),
		m.OutputDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manga) SourceDir() string {
	return filepath.Join(config.AppConfig.LibraryPath, m.Title, "source")
}

func (m *Manga) CacheDir() string {
	return filepath.Join(config.AppConfig.LibraryPath, m.Title, "cache")
}

func (m *Manga) OutputDir() string {
	return filepath.Join(config.AppConfig.LibraryPath, m.Title, "output")
}

func (m *Manga) MetadataPath() string {
	return filepath.Join(config.AppConfig.LibraryPath, m.Title, "metadata.json")
}

func LoadManga(title string) (*Manga, error) {
	manga := &Manga{
		Title: title,
	}

	entries, err := os.ReadDir(manga.SourceDir())
	if err != nil {
		log.Printf("Error reading input directory: \"%s\", error: %v\n", manga.SourceDir(), err)
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		chapterPath := filepath.Join(manga.SourceDir(), entry.Name())

		files, err := os.ReadDir(chapterPath)
		if err != nil {
			log.Println("Error reading chapter directory:", err)
			return nil, err
		}

		var images []string

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			images = append(images, file.Name())
		}

		sort.Strings(images)

		manga.Chapters = append(manga.Chapters, Chapter{
			Name:   entry.Name(),
			Images: images,
		})
	}

	sort.Slice(manga.Chapters, func(i, j int) bool {
		return manga.Chapters[i].Name < manga.Chapters[j].Name
	})

	return manga, nil
}
