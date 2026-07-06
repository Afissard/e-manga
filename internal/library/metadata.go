package library

type Metadata struct {
	Name     string                     `json:"name"`
	Target   string                     `json:"target"`
	Chapters map[string]ChapterMetadata `json:"chapters"`
}

type ChapterMetadata struct {
	ID        string `json:"id"`   // e.g. "001", "057", "100.5"
	Name      string `json:"name"` // Display name or folder name
	Checksum  string `json:"checksum"`
	PageCount int    `json:"pageCount"`
	Processed bool   `json:"processed"`
}
