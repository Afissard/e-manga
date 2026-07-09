package processing

import (
	"e-manga/internal/config"
	"e-manga/internal/library"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Options struct {
	Target     config.Target
	AutoRotate bool
}

type pageResult struct {
	index int
	img   image.Image
	err   error
}

func ProcessToCBZ(mangaName string, opts Options) error {
	log.Printf("Processing manga %s to CBZ with options: %+v", mangaName, opts)

	// Loading manga
	manga, err := library.LoadManga(mangaName)
	if err != nil {
		log.Fatalf("failed to load manga: %v", err)
		return err
	}
	manga.LoadMetadata()

	// Compare source and metadata, update if necessary
	if len(manga.Chapters) != len(manga.Metadata.Chapters) {
		log.Printf("Metadata are outdated. Updating metadata for manga: %s", mangaName)
		manga.UpdateMetadata("", "", "", "", false)
	} else {
		log.Printf("Metadata are up-to-date for manga: %s", mangaName)
	}

	// Creating output file
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
			err := cbz.AddImage("0000.png", coverImg)
			if err != nil {
				log.Fatalf("failed to add cover image to CBZ: %v", err)
				return err
			} else {
				log.Printf("Added cover image %s to CBZ for manga %s.", manga.Cover, manga.Title)
			}
		}
	} else {
		log.Printf("No cover image specified for manga %s. Skipping cover.", manga.Title)
	}

	// Add chapters and images to CBZ
	results := make(chan pageResult)
	var wg sync.WaitGroup
	index := 1
	for _, chapter := range manga.Chapters {
		chapter := chapter
		for _, filename := range chapter.Images {
			filename := filename
            pageIndex := index

			wg.Add(1)
			go func() {
				defer wg.Done()
				img, err := GetPageImageToCBZ(manga, &chapter, filename, opts)
				results <- pageResult{index: pageIndex, img: img, err: err}
			}()

			//GetPageImageToCBZ(manga, &chapter, filename, opts, index, cbz)
			index++
		}
		// Update chapter metadata
		chapterMetadata := manga.Metadata.Chapters[chapter.Name]
		chapterMetadata.PageCount = len(chapter.Images)
		manga.Metadata.Chapters[chapter.Name] = chapterMetadata
	}
	// Wait for all images to be processed
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Collect results and add images to CBZ
	pages := make(map[int]image.Image, index-1)
	for res := range results {
		if res.err != nil {
			log.Fatalf("failed to process image n°%d: %v", res.index, res.err)
		}
		pages[res.index] = res.img
	}

	for i := 1; i < index; i++ {
		name := fmt.Sprintf("%06d.png", i)
		if err := cbz.AddImage(name, pages[i]); err != nil {
			log.Fatalf("failed to add image %s to CBZ: %v", name, err)
			return err
		} else {
			log.Printf("Added image %s to CBZ for manga %s.", name, manga.Title)
		}
	}

	// Generate ComicInfo.xml and add to CBZ
	if err := cbz.GenerateComicInfoXML(manga); err != nil {
		log.Fatalf("failed to generate ComicInfo.xml: %v", err)
		return err
	} else {
		log.Printf("Successfully generated ComicInfo.xml for manga %s.", manga.Title)
	}

	// Save updated metadata
	if err := manga.Save(); err != nil {
		log.Fatalf("failed to save metadata: %v", err)
		return err
	}

	log.Printf("Successfully created CBZ file: %s", manga.Title+".cbz")
	return nil
}

func GetPageImageToCBZ(manga *library.Manga, chapter *library.Chapter, filename string, opts Options) (image.Image, error) {
	img, err := manga.LoadImageFromCache(chapter.Name, filename)
	if err != nil {
		log.Printf("Failed to load image %s from cache for manga %s, chapter %s.", filename, manga.Title, chapter.Name)
		log.Printf("Processing image %s from manga %s, chapter %s with options: %+v", filename, manga.Title, chapter.Name, opts)

		path := filepath.Join(manga.SourceDir(), chapter.Name, filename)

		img, err = LoadSourceImage(path)
		if err != nil {
			return nil, err
		}

		img = imageTraitment(img, opts)

		filename = filename[:len(filename)-len(filepath.Ext(filename))] + ".png"
		manga.SaveImageToCache(chapter.Name, filename, img)
	} else {
		log.Printf("Loaded image %s from cache for manga %s, chapter %s.", filename, manga.Title, chapter.Name)
	}

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
