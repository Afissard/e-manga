package main

import (
	"e-manga/internal/command"
	"e-manga/internal/config"
	"e-manga/internal/tui"
	"flag"
	"log"
	"os"
)

func main() {
	// logger initialization
	var err error = nil
	config.LogSrv, err = config.InitLogger(config.LogLevelInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer config.LogSrv.Close()
	go config.LogSrv.Run()

	config.LogSrv.LogMessage("Starting e-manga application...", config.LogLevelInfo)

	// Config initialization
	err = config.LoadConfigFile()
	if err != nil {
		config.LogSrv.LogMessage("Error loading config file: "+err.Error(), config.LogLevelError)
		panic(err)
	}

	// If no arguments are provided, run the TUI
	if len(os.Args) == 1 {
		config.Configuration.TuiMode = true
		if err := tui.RunTUI(); err != nil {
			config.LogSrv.LogMessage("Error running TUI: "+err.Error(), config.LogLevelError)
			panic(err)
		}
		return
	}

	// Handle command-line arguments
	if len(os.Args) < 2 {
		config.LogSrv.LogMessage("expected 'new-manga' or 'process' subcommands", config.LogLevelError)
		panic("expected 'new-manga' or 'process' subcommands")
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
			config.LogSrv.LogMessage("Error parsing flags: "+err.Error(), config.LogLevelError)
			panic(err)
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
			config.LogSrv.LogMessage("Error creating new manga: "+err.Error(), config.LogLevelError)
			panic(err)
		}

	case "process":
		fs := flag.NewFlagSet("process", flag.ExitOnError)

		manga := fs.String("manga", "", "manga name")
		target := fs.String("target", "none", "Target device for resizing images")

		if err := fs.Parse(os.Args[2:]); err != nil {
			config.LogSrv.LogMessage("Error parsing flags: "+err.Error(), config.LogLevelError)
			panic(err)
		}

		opts := command.ProcessMangaOptions{
			Manga:  *manga,
			Target: *target,
		}

		if err := command.ProcessManga(opts); err != nil {
			config.LogSrv.LogMessage("Error processing manga: "+err.Error(), config.LogLevelError)
			panic(err)
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
			config.LogSrv.LogMessage("Error parsing flags: "+err.Error(), config.LogLevelError)
			panic(err)
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
			config.LogSrv.LogMessage("Error updating manga metadata: "+err.Error(), config.LogLevelError)
			panic(err)
		}

	default:
		config.LogSrv.LogMessage("expected 'new-manga' or 'process' subcommands", config.LogLevelError)
		panic("expected 'new-manga' or 'process' subcommands")
	}
}
