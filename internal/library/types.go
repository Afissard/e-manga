package library

type Manga struct {
	Title       string
	Author      string
	Summary     string
	Cover       string
	URL         string
	LeftToRight bool
	Chapters    []Chapter
	Metadata    MangaMetadata
	ComicInfo   ComicInfo
}

type MangaMetadata struct {
	Title       string                     `json:"title"`
	Author      string                     `json:"author"`
	Summary     string                     `json:"summary"`
	Cover       string                     `json:"cover"`
	URL         string                     `json:"url"`
	LeftToRight bool                       `json:"leftToRight"`
	Target      string                     `json:"target"`
	Chapters    map[string]ChapterMetadata `json:"chapters"`
}

type Chapter struct {
	Name     string
	Number   float64
	URL      string
	Images   []string
	Metadata ChapterMetadata
}

type ChapterMetadata struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Checksum  string `json:"checksum"`
	PageCount int    `json:"pageCount"`
}

type ComicInfo struct {
	Title    string
	Author   string
	Summary  string
	Language string
	Manga    bool
	Chapters []Bookmark
}

type Bookmark struct {
	Name      string
	PageCount int
}
