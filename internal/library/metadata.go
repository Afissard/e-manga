package library

import (
	"e-manga/internal/config"
	"encoding/json"
	"os"
)

func (m *Manga) LoadMetadata() error {
	path := m.MetadataPath()

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		m.Metadata = MangaMetadata{
			Title:    m.Title,
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
	config.LogSrv.LogMessage("Manga loaded:\nTitle: "+m.Title+"\nAuthor: "+m.Author+"\nSummary: "+m.Summary+"\nCover: "+m.Cover+"\nURL: "+m.URL, config.LogLevelInfo)

	return nil
}

func (m *Manga) Save() error {
	data, err := json.MarshalIndent(m.Metadata, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.MetadataPath(), data, 0644)
}

func (m *Manga) UpdateMetadata(author, summary, cover, url string, leftToRight bool) error {
	if author == "" {
		author = m.Metadata.Author
	}
	if summary == "" {
		summary = m.Metadata.Summary
	}
	if cover == "" {
		cover = m.Metadata.Cover
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
