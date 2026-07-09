package main

import (
	"e-manga/internal/command"
	"e-manga/internal/tui"
	"flag"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		if err := tui.RunTUI(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if len(os.Args) < 2 {
		log.Fatal("expected 'new-manga' or 'process' subcommands")
	}

	switch os.Args[1] {
	case "new-manga":
		fs := flag.NewFlagSet("new-manga", flag.ExitOnError)
		name := *fs.String("name", "", "manga name")
		author := *fs.String("author", "", "manga author")
		summary := *fs.String("summary", "", "manga summary")
		cover := *fs.String("cover", "", "manga cover image path")
		target := *fs.String("target", "", "Target device for resizing images")
		url := *fs.String("url", "", "manga url")
		leftToRight := *fs.Bool("left-to-right", false, "Set reading direction to left-to-right")

		if err := fs.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		opts := command.NewMangaOptions{
			Name:        name,
			Author:      author,
			Summary:     summary,
			Cover:       cover,
			Target:      target,
			URL:         url,
			LeftToRight: leftToRight,
		}

		if err := command.NewManga(opts); err != nil {
			log.Fatal(err)
		}

	case "process":
		fs := flag.NewFlagSet("process", flag.ExitOnError)

		manga := fs.String("manga", "", "manga name")
		target := fs.String("target", "none", "Target device for resizing images")

		if err := fs.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		opts := command.ProcessMangaOptions{
			Manga:  *manga,
			Target: *target,
		}

		if err := command.ProcessManga(opts); err != nil {
			log.Fatal(err)
		}

	case "update-metadata":
		fs := flag.NewFlagSet("update-metadata", flag.ExitOnError)

		manga := *fs.String("manga", "", "manga name")
		author := *fs.String("author", "", "manga author")
		summary := *fs.String("summary", "", "manga summary")
		cover := *fs.String("cover", "", "manga cover image path")
		URL := *fs.String("url", "", "manga url")
		leftToRight := *fs.Bool("left-to-right", false, "Set reading direction to left-to-right")

		if err := fs.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		opts := command.UpdateMangaMetadataOptions{
			Manga:       manga,
			Author:      author,
			Summary:     summary,
			Cover:       cover,
			URL:         URL,
			LeftToRight: leftToRight,
		}

		if err := command.UpdateMangaMetadata(opts); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("expected 'new-manga' or 'process' subcommands")
		// TODO: later, load a tui for the app
	}
}
