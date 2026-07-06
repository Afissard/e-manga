package processing

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Process(inputDir, outputCBZ string, opts Options) error {
	log.Printf("Processing manga from %s to %s with options: %+v", inputDir, outputCBZ, opts)
	
	chapters, err := LoadChapters(GetSourceDir(inputDir))
	if err != nil {
		log.Fatalf("failed to load chapters: %v", err)
		return err
	}
	log.Printf("Loaded %d chapters from %s", len(chapters), inputDir)

	out, err := os.Create(outputCBZ + ".cbz")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
		return err
	}
	defer out.Close()

	cbz := NewCBZ(out)
	defer cbz.Close()

	index := 1

	for _, chapter := range chapters {
		for _, filename := range chapter.Images {
			path := filepath.Join(GetSourceDir(inputDir), chapter.Name, filename)

			img, err := LoadImage(path)
			if err != nil {
				log.Fatalf("failed to load image %s: %v", path, err)
				return err
			}

			if !opts.Target.Color {
				img = Grayscale(img)
			}

			if opts.AutoRotate && img.Bounds().Dx() > img.Bounds().Dy() {
				img = Rotate90CW(img)
			}

			if opts.Target.Width > 0 {
				img = Resize(img, opts.Target.Width, opts.Target.Height)
			}

			name := fmt.Sprintf("%06d.png", index)
			index++

			if err := cbz.AddImage(name, img); err != nil {
				log.Fatalf("failed to add image %s to CBZ: %v", name, err)
				return err
			}
		}
	}

	log.Printf("Successfully created CBZ file: %s", outputCBZ+".cbz")
	return nil
}
