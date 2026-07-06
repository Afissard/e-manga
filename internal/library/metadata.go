package library

import (
	"encoding/json"
	"os"
)

func (m *Manga) LoadMetadata() error {
	path := m.MetadataPath()

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		m.Metadata = MangaMetadata{
			Name:     m.Title,
			Target:   "",
			Chapters: make(map[string]ChapterMetadata),
		}

		return m.Save()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.Metadata)
}

func (m *Manga) Save() error {
	data, err := json.MarshalIndent(m.Metadata, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.MetadataPath(), data, 0644)
}

func (m *Manga) UpdateMetadata(target string) error {
	m.Metadata.Target = target
	m.Metadata.Chapters = make(map[string]ChapterMetadata)
	for _, chapter := range m.Chapters {
		m.Metadata.Chapters[chapter.Name] = ChapterMetadata{
			Name:      chapter.Name,
			PageCount: len(chapter.Images),
		}
	}

	return m.Save()
}
