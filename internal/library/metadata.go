package library

import (
	"encoding/json"
	"log"
	"os"
)

func (m *Manga) LoadMetadata() error {
	path := m.MetadataPath()

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		m.Metadata = MangaMetadata{
			Title:    m.Title,
			Target:   "",
			Chapters: make(map[string]ChapterMetadata),
		}

		return m.Save()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &m.Metadata)
	if err != nil {
		return err
	}

	m.Author = m.Metadata.Author
	m.Summary = m.Metadata.Summary
	m.Cover = m.Metadata.Cover
	m.URL = m.Metadata.URL
	m.LeftToRight = m.Metadata.LeftToRight

	// Log all Manga info
	log.Printf("Manga loaded:\nTitle: %s\nAuthor: %s\nSummary: %s\nCover: %s\nURL: %s", m.Title, m.Author, m.Summary, m.Cover, m.URL)

	return nil
}

func (m *Manga) Save() error {
	data, err := json.MarshalIndent(m.Metadata, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.MetadataPath(), data, 0644)
}

func (m *Manga) UpdateMetadata(author, summary, cover, target, url string, leftToRight bool) error {
	if author == "" {
		author = m.Metadata.Author
	}
	if summary == "" {
		summary = m.Metadata.Summary
	}
	if cover == "" {
		cover = m.Metadata.Cover
	}
	if target == "" {
		target = m.Metadata.Target
	}
	if url == "" {
		url = m.Metadata.URL
	}
	if !leftToRight {
		leftToRight = m.Metadata.LeftToRight
	}

	m.Metadata.Title = m.Title
	m.Metadata.Author = author
	m.Metadata.Summary = summary
	m.Metadata.Cover = cover
	m.Metadata.Target = target
	m.Metadata.URL = url
	m.Metadata.LeftToRight = leftToRight
	m.Metadata.Chapters = make(map[string]ChapterMetadata)
	for _, chapter := range m.Chapters {
		m.Metadata.Chapters[chapter.Name] = ChapterMetadata{
			Name:      chapter.Name,
			PageCount: len(chapter.Images),
		}
	}

	return m.Save()
}

