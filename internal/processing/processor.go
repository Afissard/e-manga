package processing

import (
	"e-manga/internal/config"
	"e-manga/internal/library"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sync"
)

type pageResult struct {
	index int
	img   image.Image
	err   error
}

func ProcessToCBZ(mangaName string, target config.Target) error {
	config.LogSrv.LogMessage(fmt.Sprintf("Processing manga %s to CBZ with options: %+v", mangaName, target), config.LogLevelInfo)

	// Loading manga
	manga, err := library.LoadManga(mangaName)
	manga.Target = target.Name
	if err != nil {
		config.LogSrv.LogMessage(fmt.Sprintf("Failed to load manga: %v", err), config.LogLevelError)
		return err
	}
	manga.LoadMetadata()

	// Compare source and metadata, update if necessary
	if len(manga.Chapters) != len(manga.Metadata.Chapters) {
		config.LogSrv.LogMessage(fmt.Sprintf("Metadata are outdated. Updating metadata for manga: %s", mangaName), config.LogLevelWarning)
		manga.UpdateMetadata("", "", "", "", false)
	} else {
		config.LogSrv.LogMessage(fmt.Sprintf("Metadata are up-to-date for manga: %s", mangaName), config.LogLevelInfo)
	}

	// Creating output file
	err = CreateOutputForTarget(manga, target)
	if err != nil {
		config.LogSrv.LogMessage(fmt.Sprintf("failed to create output directory: %v", err), config.LogLevelError)
		return err
	}
	out, err := os.Create(filepath.Join(manga.OutputDir(), manga.Target, manga.Title+".cbz"))
	if err != nil {
		config.LogSrv.LogMessage(fmt.Sprintf("failed to create output file: %v", err), config.LogLevelError)
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
			config.LogSrv.LogMessage(fmt.Sprintf("Failed to load cover image %s for manga %s. Skipping cover.", manga.Cover, manga.Title), config.LogLevelWarning)
		} else {
			coverImg = imageTraitment(coverImg, target)
			err := cbz.AddImage("0000.png", coverImg)
			if err != nil {
				config.LogSrv.LogMessage(fmt.Sprintf("failed to add cover image to CBZ: %v", err), config.LogLevelError)
				return err
			} else {
				config.LogSrv.LogMessage(fmt.Sprintf("Added cover image %s to CBZ for manga %s.", manga.Cover, manga.Title), config.LogLevelInfo)
			}
		}
	} else {
		config.LogSrv.LogMessage(fmt.Sprintf("No cover image specified for manga %s. Skipping cover.", manga.Title), config.LogLevelWarning)
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
				img, err := GetPageImageToCBZ(manga, &chapter, filename, target)
				results <- pageResult{index: pageIndex, img: img, err: err}
			}()

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
			config.LogSrv.LogMessage(fmt.Sprintf("failed to process image n°%d: %v (%d/%d)", res.index, res.err, res.index, index-1), config.LogLevelError)
		}
		pages[res.index] = res.img
	}

	for i := 1; i < index; i++ {
		name := fmt.Sprintf("%06d.png", i)
		if err := cbz.AddImage(name, pages[i]); err != nil {
			config.LogSrv.LogMessage(fmt.Sprintf("failed to add image %s to CBZ: %v (%d/%d)", name, err, i, index-1), config.LogLevelError)
			return err
		} else {
			config.LogSrv.LogMessage(fmt.Sprintf("Added image %s to CBZ for manga %s. (%d/%d)", name, manga.Title, i, index-1), config.LogLevelInfo)
		}
	}

	// Generate ComicInfo.xml and add to CBZ
	if err := cbz.GenerateComicInfoXML(manga); err != nil {
		config.LogSrv.LogMessage(fmt.Sprintf("failed to generate ComicInfo.xml: %v", err), config.LogLevelError)
		return err
	} else {
		config.LogSrv.LogMessage(fmt.Sprintf("Successfully generated ComicInfo.xml for manga %s.", manga.Title), config.LogLevelInfo)
	}

	// Save updated metadata
	if err := manga.Save(); err != nil {
		config.LogSrv.LogMessage(fmt.Sprintf("failed to save metadata: %v", err), config.LogLevelError)
		return err
	}

	config.LogSrv.LogMessage(fmt.Sprintf("Successfully created CBZ file: %s", manga.Title+".cbz"), config.LogLevelInfo)
	return nil
}

func GetPageImageToCBZ(manga *library.Manga, chapter *library.Chapter, filename string, target config.Target) (image.Image, error) {
	img, err := manga.LoadImageFromCache(chapter.Name, filename)
	if err != nil {
		config.LogSrv.LogMessage(fmt.Sprintf("Failed to load image %s from cache for manga %s, chapter %s.", filename, manga.Title, chapter.Name), config.LogLevelError)
		config.LogSrv.LogMessage(fmt.Sprintf("Processing image %s from manga %s, chapter %s with options: %+v", filename, manga.Title, chapter.Name, target), config.LogLevelInfo)

		path := filepath.Join(manga.SourceDir(), chapter.Name, filename)

		img, err = LoadSourceImage(path)
		if err != nil {
			return nil, err
		}

		img = imageTraitment(img, target)

		filename = filename[:len(filename)-len(filepath.Ext(filename))] + ".png"
		manga.SaveImageToCache(chapter.Name, filename, img)
	} else {
		config.LogSrv.LogMessage(fmt.Sprintf("Loaded image %s from cache for manga %s, chapter %s.", filename, manga.Title, chapter.Name), config.LogLevelInfo)
	}

	return img, nil
}

func imageTraitment(img image.Image, target config.Target) image.Image {
	if !target.Color {
		img = Grayscale(img)
	}

	if target.AutoRotate && img.Bounds().Dx() > img.Bounds().Dy() {
		img = Rotate90CW(img)
	}

	if target.Width > 0 {
		img = Resize(img, target.Width, target.Height)
	}

	return img
}
