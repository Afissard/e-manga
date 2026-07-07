package main

import (
	"flag"
	"fmt"
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
		url := fs.String("url", "", "manga url")

		if err := fs.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		if err := NewManga(*name, *url); err != nil {
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

	default:
		log.Fatal("expected 'new-manga' or 'process' subcommands")
		// TODO: later, load a tui for the app
	}
}

func logErrorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
