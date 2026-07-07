package model

type Manga struct {
    Title       string
    Description string
    Chapters    []Chapter
}

type Chapter struct {
    Title  string
    Number int
    URL    string
    Pages  []string
}