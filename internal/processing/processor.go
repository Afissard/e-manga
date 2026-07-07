package processing

import (
	"e-manga/internal/config"
	"e-manga/internal/library"
	"fmt"
	"image"
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
	manga.UpdateComicInfo()

	// Compare source and metadata, update if necessary
	if len(manga.Chapters) != len(manga.Metadata.Chapters) || manga.Metadata.Target != opts.Target.Name {
		log.Printf("Metadata are outdated. Updating metadata for manga: %s", mangaName)
		manga.UpdateMetadata("", "", "", "", opts.Target.Name, false)
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

	// Add cover image to CBZ
	if manga.Cover != "" {
		coverPath := filepath.Join(manga.SourceDir(), manga.Cover)
		coverImg, err := LoadSourceImage(coverPath)
		if err != nil {
			log.Printf("Failed to load cover image %s for manga %s. Skipping cover.", manga.Cover, manga.Title)
		} else {
			coverImg = imageTraitment(coverImg, opts)
			if err := cbz.AddImage("0000", coverImg); err != nil {
				log.Fatalf("failed to add cover image to CBZ: %v", err)
				return err
			} else {
				log.Printf("Added cover image %s to CBZ for manga %s.", manga.Cover, manga.Title)
			}
		}
	} else {
		log.Printf("No cover image specified for manga %s. Skipping cover.", manga.Title)
	}

	// Generate ComicInfo.xml and add to CBZ
	if err := cbz.GenerateComicInfoXML(&manga.ComicInfo); err != nil {
		log.Fatalf("failed to generate ComicInfo.xml: %v", err)
		return err
	} else {
		log.Printf("Successfully generated ComicInfo.xml for manga %s.", manga.Title)
	}

	// Add chapters and images to CBZ
	index := 1
	for _, chapter := range manga.Chapters {
		for _, filename := range chapter.Images {

			img, err := manga.LoadImageFromCache(chapter.Name, filename)
			if err != nil {
				log.Printf("Failed to load image %s from cache for manga %s, chapter %s. Reprocessing.", filename, manga.Title, chapter.Name)
				img, err = ProcessImage(manga, &chapter, filename, opts)
				if err != nil {
					return err
				}
			} else {
				log.Printf("Loaded image %s from cache for manga %s, chapter %s.", filename, manga.Title, chapter.Name)
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

func ProcessImage(manga *library.Manga, chapter *library.Chapter, filename string, opts Options) (img image.Image, err error) {
	log.Printf("Processing image %s from manga %s, chapter %s with options: %+v", filename, manga.Title, chapter.Name, opts)

	path := filepath.Join(manga.SourceDir(), chapter.Name, filename)

	img, err = LoadSourceImage(path)
	if err != nil {
		log.Fatalf("failed to load image %s: %v", path, err)
		return nil, err
	}

	/*
		if !opts.Target.Color {
			img = Grayscale(img)
		}

		if opts.AutoRotate && img.Bounds().Dx() > img.Bounds().Dy() {
			img = Rotate90CW(img)
		}

		if opts.Target.Width > 0 {
			img = Resize(img, opts.Target.Width, opts.Target.Height)
		}
	*/

	img = imageTraitment(img, opts)

	filename = filename[:len(filename)-len(filepath.Ext(filename))] + ".png"
	manga.SaveImageToCache(chapter.Name, filename, img)

	return img, nil
}

func imageTraitment(img image.Image, opts Options) image.Image {
	if !opts.Target.Color {
		img = Grayscale(img)
	}

	if opts.AutoRotate && img.Bounds().Dx() > img.Bounds().Dy() {
		img = Rotate90CW(img)
	}

	if opts.Target.Width > 0 {
		img = Resize(img, opts.Target.Width, opts.Target.Height)
	}

	return img
}
