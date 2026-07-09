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
	target, ok := config.Targets[opts.Target]
	if !ok {
		log.Fatalf("unknown target: %s", opts.Target)
	}

	options := processing.Options{
		Target: target,
	}

	return processing.ProcessToCBZ(opts.Manga, options)
}
