package main

import (
	"e-manga/internal/config"
	"e-manga/internal/processing"
	"log"
)

func ProcessManga(manga string, targetName string) error {
	target, ok := config.Targets[targetName]
	if !ok {
		log.Fatalf("unknown target: %s", targetName)
	}

	opts := processing.Options{
		Target: target,
	}

	return processing.ProcessToCBZ(manga, opts)
}
