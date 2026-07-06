package processing

import (
	"fmt"
	"os"
	"path/filepath"
)

func Process(inputDir, outputCBZ string, opts Options) error {
	chapters, err := LoadChapters(inputDir)
	if err != nil {
		return err
	}

	out, err := os.Create(outputCBZ + ".cbz")
	if err != nil {
		return err
	}
	defer out.Close()

	cbz := NewCBZ(out)
	defer cbz.Close()

	index := 1

	for _, chapter := range chapters {
		for _, filename := range chapter.Images {
			path := filepath.Join(inputDir, chapter.Name, filename)

			img, err := LoadImage(path)
			if err != nil {
				return err
			}

			img = Grayscale(img)

			if opts.AutoRotate && img.Bounds().Dx() > img.Bounds().Dy() {
				img = Rotate90CW(img)
			}

			if opts.Target.Width > 0 {
				img = Resize(img, opts.Target.Width, opts.Target.Height)
			}

			name := fmt.Sprintf("%06d.png", index)
			index++

			if err := cbz.AddImage(name, img); err != nil {
				return err
			}
		}
	}

	return nil
}
