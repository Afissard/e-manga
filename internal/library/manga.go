package library

import (
	"e-manga/internal/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

func (m *Manga) Init(title, author, summary, cover, target, url string, leftToRight bool) error {
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

	m.LoadMetadata()
	m.UpdateMetadata(author, summary, cover, target, url, leftToRight)

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

		// get the chapter number from the directory name with a regex, if it fails, set it to 0.0
		// the regex get the last number in the string, with optional decimal point, and optional leading zeros
		re := regexp.MustCompile(`(\d+(?:\.\d+)?)$`)
		var chapterNumber float64 = 0.0
		match := re.FindStringSubmatch(entry.Name())
		if match == nil {
			fmt.Printf("%q -> no match\n", entry.Name())
			continue
		}

		chapterNumber, err = strconv.ParseFloat(match[1], 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		manga.Chapters = append(manga.Chapters, Chapter{
			Name:   entry.Name(),
			Number: chapterNumber,
			Images: images,
		})
	}

	/*
		sort.Slice(manga.Chapters, func(i, j int) bool {
			return manga.Chapters[i].Name < manga.Chapters[j].Name
		})
	*/

	// sort chapters by number, if number is 0, sort by name
	sort.Slice(manga.Chapters, func(i, j int) bool {
		if manga.Chapters[i].Number == 0 && manga.Chapters[j].Number == 0 {
			return manga.Chapters[i].Name < manga.Chapters[j].Name
		}
		if manga.Chapters[i].Number == 0 {
			return false
		}
		if manga.Chapters[j].Number == 0 {
			return true
		}
		return manga.Chapters[i].Number < manga.Chapters[j].Number
	})

	/*
		// print the chapters found in order to verify the sorting
		log.Printf("Chapters found for manga \"%s\":\n", manga.Title)
		for _, chapter := range manga.Chapters {
			log.Printf("  - %s (Number: %v)\n", chapter.Name, chapter.Number)
		}
	*/

	return manga, nil
}
