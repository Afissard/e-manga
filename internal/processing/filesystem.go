package processing

import (
	"e-manga/internal/config"
	"os"
	"path/filepath"
)

/*
Wrapper around filesystem operations.
- create temp directories
- copy files
- remove temporary data
*/

func CreateLibraryDir() error {
	if _, err := os.Stat(config.AppConfig.LibraryPath); os.IsNotExist(err) {
		return os.MkdirAll(config.AppConfig.LibraryPath, 0755)
	}
	return nil
}

// in the library, create a new directory for the book, with a source folder, a cache folder and a metadata.json file and a output folder
func CreateNewBookDir(root string) (string, error) {
	_, folder := filepath.Split(root)
	bookDir := filepath.Join(config.AppConfig.LibraryPath, folder)

	if _, err := os.Stat(bookDir); os.IsNotExist(err) {
		err := os.MkdirAll(bookDir, 0755)
		if err != nil {
			return "", err
		}
	}

	sourceDir := filepath.Join(bookDir, "source")
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		err := os.MkdirAll(sourceDir, 0755)
		if err != nil {
			return "", err
		}
	}

	cacheDir := filepath.Join(bookDir, "cache")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return "", err
		}
	}

	outputDir := filepath.Join(bookDir, "output")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return "", err
		}
	}

	return bookDir, nil
}

func OutputPath(root string) string {
	_, folder := filepath.Split(root)
	return filepath.Join(config.AppConfig.LibraryPath, folder, "output/", folder)
}

func GetBookDir(root string) string {
	_, folder := filepath.Split(root)
	return filepath.Join(config.AppConfig.LibraryPath, folder)
}

func GetSourceDir(root string) string {
	_, folder := filepath.Split(root)
	return filepath.Join(config.AppConfig.LibraryPath, folder, "source")
}

func GetCacheDir(root string) string {
	_, folder := filepath.Split(root)
	return filepath.Join(config.AppConfig.LibraryPath, folder, "cache")
}

func GetOutputDir(root string) string {
	_, folder := filepath.Split(root)
	return filepath.Join(config.AppConfig.LibraryPath, folder, "output")
}