package processing

/*
Creates the final ebook.
- CBZ
later add other formats:
- EPUB
*/

import (
	"archive/zip"
	"image"
	"io"
)

type CBZWriter struct {
	zip *zip.Writer
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
