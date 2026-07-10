package command

import (
	"e-manga/internal/config"
	"e-manga/internal/processing"
	"log"
)

type ProcessMangaOptions struct {
	Manga  string
	Target string
}

func ProcessManga(opts ProcessMangaOptions) error {
	// check if the target exists in the config file
	for _, target := range config.Configuration.Targets {
		if target.Name == opts.Target {
			return processing.ProcessToCBZ(opts.Manga, target)
		}
	}
	/*
		target, ok := config.Configuration.Targets[opts.Target]
		if !ok {
			log.Fatalf("unknown target: %s", opts.Target)
		}
		options := processing.Options{
			Target: target,
		}

		return processing.ProcessToCBZ(opts.Manga, options)
	*/
	log.Fatalf("unknown target: %s", opts.Target)
	return nil
}
