package library

type Manga struct {
	Title    string
	Chapters []Chapter
	Metadata MangaMetadata
	// TODO: add metadata fields like author, genre, etc.
}

type MangaMetadata struct {
	Name     string                     `json:"name"`
	Target   string                     `json:"target"`
	Chapters map[string]ChapterMetadata `json:"chapters"`
}

type Chapter struct {
	Name     string
	Images   []string
	Metadata ChapterMetadata
}

type ChapterMetadata struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Checksum  string `json:"checksum"`
	PageCount int    `json:"pageCount"`
}
