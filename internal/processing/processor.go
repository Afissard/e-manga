package processing

import (
	"e-manga/internal/config"
	"e-manga/internal/library"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Options struct {
	Target     config.Target
	AutoRotate bool
}

func ProcessToCBZ(mangaName string, opts Options) error {
	log.Printf("Processing manga from %s to CBZ with options: %+v", mangaName, opts)

	// Loading manga
	manga, err := library.LoadManga(mangaName)
	if err != nil {
		log.Fatalf("failed to load manga: %v", err)
		return err
	}
	manga.LoadMetadata()

	// Compare source and metadata, update if necessary
	if len(manga.Chapters) != len(manga.Metadata.Chapters) || manga.Metadata.Target != opts.Target.Name {
		log.Printf("Metadata are outdated. Updating metadata for manga: %s", mangaName)
		manga.UpdateMetadata(opts.Target.Name)
	} else {
		log.Printf("Metadata are up-to-date for manga: %s", mangaName)
	}

	// Creating output
	out, err := os.Create(filepath.Join(manga.OutputDir(), manga.Title+".cbz"))
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
		return err
	}
	defer out.Close()

	cbz := NewCBZ(out)
	defer cbz.Close()

	index := 1

	for _, chapter := range manga.Chapters {
		for _, filename := range chapter.Images {
			path := filepath.Join(manga.SourceDir(), chapter.Name, filename)

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
		// Update chapter metadata
		chapterMetadata := manga.Metadata.Chapters[chapter.Name]
		chapterMetadata.PageCount = len(chapter.Images)
		manga.Metadata.Chapters[chapter.Name] = chapterMetadata
	}

	// Save updated metadata
	if err := manga.Save(); err != nil {
		log.Fatalf("failed to save metadata: %v", err)
		return err
	}

	log.Printf("Successfully created CBZ file: %s", manga.Title+".cbz")
	return nil
}
