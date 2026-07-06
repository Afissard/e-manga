package main

import (
	"flag"
	"log"

	"e-manga/internal/config"
	"e-manga/internal/processing"
)

func main() {
	input := flag.String("input", "", "Path to the manga folder")
	targetFlag := flag.String("target", "none", "Target device for resizing images (e.g., 'kindle-paperwhite-7')")

	flag.Parse()

	target, ok := config.Targets[*targetFlag]
	if !ok {
		log.Fatalf("unknown target: %s", *targetFlag)
	}

	opts := processing.Options{
		Target: target,
	}

	processing.CreateOutputDir()
	err := processing.Process(*input, processing.OutputPath(*input), opts)

	if err != nil {
		log.Fatal(err)
	}
}
