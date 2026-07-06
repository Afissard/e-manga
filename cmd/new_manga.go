package main

import "e-manga/internal/processing"

func NewManga(name string) error {
	err := processing.CreateLibraryDir()
	if err != nil {
		return err
	}

	_, err = processing.CreateNewBookDir(name)
	if err != nil {
		return err
	}

	return nil
}
