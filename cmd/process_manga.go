package main

import (
	"e-manga/internal/config"
	"e-manga/internal/processing"
)

func ProcessManga(manga string, targetName string) error {
	target, ok := config.Targets[targetName]
	if !ok {
		return logErrorf("unknown target: %s", targetName)
	}

	opts := processing.Options{
		Target: target,
	}

	return processing.Process(manga, processing.OutputPath(manga), opts)
}
