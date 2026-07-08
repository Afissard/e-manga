package processing

import (
	"e-manga/internal/library"
	"encoding/xml"
	"os"
	"path/filepath"
)

type comicInfoXML struct {
	XMLName xml.Name `xml:"ComicInfo"`

	Title         string `xml:"Title,omitempty"`
	Series        string `xml:"Series,omitempty"`
	Writer        string `xml:"Writer,omitempty"`
	Artist        string `xml:"Artist,omitempty"`
	Summary       string `xml:"Summary,omitempty"`
	LanguageISO   string `xml:"LanguageISO,omitempty"`
	PageCount     int    `xml:"PageCount,omitempty"`
	BlackAndWhite string `xml:"BlackAndWhite,omitempty"`
	Manga         string `xml:"Manga,omitempty"`

	Pages comicPages `xml:"Pages"`
}

type comicPages struct {
	Pages []comicPage `xml:"Page"`
}

type comicPage struct {
	Image    int    `xml:"Image,attr"`
	Type     string `xml:"Type,attr,omitempty"`
	Bookmark string `xml:"Bookmark,attr,omitempty"`
}

func (c *CBZWriter) GenerateComicInfoXML(manga *library.Manga) error {

	xmlInfo := comicInfoXML{
		Title:         manga.Title,
		Series:        manga.Title,
		Writer:        manga.Author,
		Artist:        manga.Author,
		Summary:       manga.Summary,
		LanguageISO:   manga.Language,
		BlackAndWhite: "Yes",
	}

	if manga.LeftToRight {
		xmlInfo.Manga = "YesAndRightToLeft"
	} else {
		xmlInfo.Manga = "No"
	}

	pageIndex := 0

	for i, chapter := range manga.Chapters {
		if i == 0 {
			// Cover page
			xmlInfo.Pages.Pages = append(xmlInfo.Pages.Pages, comicPage{
				Image: pageIndex,
				Type:  "FrontCover",
			})
			// Chapter 1
			xmlInfo.Pages.Pages = append(xmlInfo.Pages.Pages, comicPage{
				Image:    pageIndex + 1,
				Bookmark: chapter.Name,
				Type:     "Chapter",
			})
		} else {
			xmlInfo.Pages.Pages = append(xmlInfo.Pages.Pages, comicPage{
				Image:    pageIndex,
				Bookmark: chapter.Name,
				Type:     "Chapter",
			})
		}

		xmlInfo.PageCount += len(chapter.Images)
		pageIndex += len(chapter.Images)
	}

	data, err := xml.MarshalIndent(xmlInfo, "", "    ")
	if err != nil {
		return err
	}

	header := []byte(xml.Header)
	data = append(header, data...)

	entry, err := c.zip.Create("ComicInfo.xml")
	if err != nil {
		return err
	}

	_, err = entry.Write(data)

	// save to cache the ComicInfo.xml for future use
	err = os.WriteFile(filepath.Join(manga.CacheDir(), manga.Metadata.Target, "ComicInfo.xml"), data, 0644)
	if err != nil {
		return err
	}

	return err
}
