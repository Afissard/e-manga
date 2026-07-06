package processing

import (
	"os"
	"path/filepath"
)

/*
Wrapper around filesystem operations.
- create temp directories
- copy files
- remove temporary data
*/

const TempDir = "temp"
const OutputDir = "output"

func TempPath(root string) string {
	_, folder := filepath.Split(root)
	return filepath.Join(TempDir, folder)
}

func CreateOutputDir() error {
	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		return os.MkdirAll(OutputDir, 0755)
	}
	return nil
}

func OutputPath(root string) string {
	_, folder := filepath.Split(root)
	return filepath.Join(OutputDir, folder)
}

/*
func outputPath(outputRoot, chapter, filename string) string {
	return filepath.Join(outputRoot, chapter, filename)
}
*/
