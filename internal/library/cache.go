package library

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
)

func getCachedFilePath(manga *Manga, chapterName, filename string) string {
	return filepath.Join(manga.CacheDir(), manga.Metadata.Target, chapterName, filename)
}

func (m *Manga) SaveImageToCache(chapterName, filename string, img image.Image) error {
	cacheFilePath := getCachedFilePath(m, chapterName, filename)

	if err := os.MkdirAll(filepath.Dir(cacheFilePath), 0755); err != nil {
		return err
	}

	f, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func (m *Manga) LoadImageFromCache(chapterName, filename string) (image.Image, error) {
	// make sure the filename end with .png and if not, add it instead of the current extension
	if filepath.Ext(filename) != ".png" {
		filename = filename[:len(filename)-len(filepath.Ext(filename))] + ".png"
	}

	cacheFilePath := getCachedFilePath(m, chapterName, filename)
	f, err := os.Open(cacheFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, err
}

func (m *Manga) IsImageInCache(chapterName, filename string) bool {
	if filepath.Ext(filename) != ".png" {
		filename = filename[:len(filename)-len(filepath.Ext(filename))] + ".png"
	}

	cacheFilePath := getCachedFilePath(m, chapterName, filename)
	_, err := os.Stat(cacheFilePath)
	return !os.IsNotExist(err)
}

func (m *Manga) ClearCache() error {
	cacheDir := m.CacheDir()
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return nil // Cache directory does not exist, nothing to clear
	}
	return os.RemoveAll(cacheDir)
}
