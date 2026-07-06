package processing

/*
Responsible for:
- finding chapter folders
- sorting them
- listing images
- validating filenames
*/

import (
	"log"
	"os"
	"path/filepath"
	"sort"
)

func LoadChapters(root string) ([]Chapter, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Printf("Error reading input directory: \"%s\", error: %v\n", root, err)
		return nil, err
	}

	var chapters []Chapter

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		chapterPath := filepath.Join(root, entry.Name())

		files, err := os.ReadDir(chapterPath)
		if err != nil {
			log.Println("Error reading chapter directory:", err)
			return nil, err
		}

		var images []string

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			images = append(images, file.Name())
		}

		sort.Strings(images)

		chapters = append(chapters, Chapter{
			Name:   entry.Name(),
			Images: images,
		})
	}

	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].Name < chapters[j].Name
	})

	return chapters, nil
}
