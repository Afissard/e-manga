package download

import (
	"io"
	"net/http"
	"os"
)

func Image(url, filename string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

/*
for i, page := range chapter.Pages {
    download.Image(page, fmt.Sprintf("%03d.jpg", i+1))
}
*/
