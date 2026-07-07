package bigsolo

import (
	"e-manga/internal/scraper/model"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	client *http.Client
}

func New() *Scraper {
	return &Scraper{
		client: &http.Client{},
	}
}

func (s *Scraper) GetManga(url string) (*model.Manga, error) {

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	manga := &model.Manga{}

	manga.Title = doc.Find("h1").Text()
	manga.Description = doc.Find(".description").Text()

	doc.Find(".chapter-list a").Each(func(i int, sel *goquery.Selection) {

		href, _ := sel.Attr("href")

		manga.Chapters = append(manga.Chapters, model.Chapter{
			Title: sel.Text(),
			URL:   href,
		})
	})

	return manga, nil
}

func (s *Scraper) GetChapter(url string) (*model.Chapter, error) {

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	chapter := &model.Chapter{}

	doc.Find("img.page").Each(func(i int, sel *goquery.Selection) {

		src, _ := sel.Attr("src")
		chapter.Pages = append(chapter.Pages, src)

	})

	return chapter, nil
}
