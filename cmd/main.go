package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("expected 'new-manga' or 'process' subcommands")
	}

	switch os.Args[1] {
	case "new-manga":
		fs := flag.NewFlagSet("new-manga", flag.ExitOnError)
		name := fs.String("name", "", "manga name")
		author := fs.String("author", "", "manga author")
		summary := fs.String("summary", "", "manga summary")
		cover := fs.String("cover", "", "manga cover image path")
		target := fs.String("target", "", "Target device for resizing images")
		url := fs.String("url", "", "manga url")
		leftToRight := fs.Bool("left-to-right", false, "Set reading direction to left-to-right")

		if err := fs.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		if err := NewManga(*name, *author, *summary, *cover, *target, *url, *leftToRight); err != nil {
			log.Fatal(err)
		}

	case "process":
		fs := flag.NewFlagSet("process", flag.ExitOnError)
		manga := fs.String("manga", "", "manga name")
		target := fs.String("target", "none", "Target device for resizing images")

		if err := fs.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		if err := ProcessManga(*manga, *target); err != nil {
			log.Fatal(err)
		}

	case "update-metadata":
		fs := flag.NewFlagSet("update-metadata", flag.ExitOnError)
		manga := fs.String("manga", "", "manga name")
		author := fs.String("author", "", "manga author")
		summary := fs.String("summary", "", "manga summary")
		cover := fs.String("cover", "", "manga cover image path")
		target := fs.String("target", "", "Target device for resizing images")
		url := fs.String("url", "", "manga url")
		leftToRight := fs.Bool("left-to-right", false, "Set reading direction to left-to-right")

		if err := fs.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		if err := UpdateMangaMetadata(*manga, *author, *summary, *cover, *target, *url, *leftToRight); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("expected 'new-manga' or 'process' subcommands")
		// TODO: later, load a tui for the app
	}
}
