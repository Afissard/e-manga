package processing

/*
Creates the final ebook.
- CBZ
later add other formats:
- EPUB
*/

import (
	"archive/zip"
	"e-manga/internal/config"
	"e-manga/internal/library"
	"image"
	"io"
	"os"
	"path/filepath"
)

type CBZWriter struct {
	zip *zip.Writer
}

func CreateOutputForTarget(manga *library.Manga, target config.Target) error {
	// Create output directory for the target if it doesn't exist
	outputDir := filepath.Join(manga.OutputDir(), target.Name)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	return nil
}

func NewCBZ(w io.Writer) *CBZWriter {
	return &CBZWriter{
		zip: zip.NewWriter(w),
	}
}

func (c *CBZWriter) AddImage(name string, img image.Image) error {
	entry, err := c.zip.Create(name)
	if err != nil {
		return err
	}

	return EncodePNG(entry, img)
}

func (c *CBZWriter) Close() error {
	return c.zip.Close()
}
