package processing

import (
	"e-manga/internal/library"
	"encoding/xml"
)

type comicInfoXML struct {
	XMLName xml.Name `xml:"ComicInfo"`

	Title         string `xml:"Title,omitempty"`
	Series        string `xml:"Series,omitempty"`
	Writer        string `xml:"Writer,omitempty"`
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

func (c *CBZWriter) GenerateComicInfoXML(info *library.ComicInfo) error {
	xmlInfo := comicInfoXML{
		Title:         info.Title,
		Series:        info.Title,
		Writer:        info.Author,
		Summary:       info.Summary,
		LanguageISO:   info.Language,
		BlackAndWhite: "Yes",
	}

	if info.Manga {
		xmlInfo.Manga = "YesAndRightToLeft"
	} else {
		xmlInfo.Manga = "No"
	}

	pageIndex := 0

	for i, chapter := range info.Chapters {

		if i == 0 {
			xmlInfo.Pages.Pages = append(xmlInfo.Pages.Pages, comicPage{
				Image: pageIndex,
				Type:  "FrontCover",
			})
		} else {
			xmlInfo.Pages.Pages = append(xmlInfo.Pages.Pages, comicPage{
				Image:    pageIndex,
				Bookmark: chapter.Name,
			})
		}

		xmlInfo.PageCount += chapter.PageCount
		pageIndex += chapter.PageCount
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

	return err
}
