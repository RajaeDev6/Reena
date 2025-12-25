package dir

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	imageutil "github.com/RajaeDev6/reena/internal/imageUtil"
	"github.com/RajaeDev6/reena/internal/namer"
)

func Scan(path string, generator namer.ImageNamer) {

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	if len(entries) == 0 {
		log.Fatal("Dir is empty")
	}

	for _, de := range entries {

		if de.IsDir() {
			continue
		}

		if imageutil.IsImage(de.Name()) {

			fp := filepath.Join(path, de.Name())
			fmt.Println("FILE: " + fp)
			ext := imageutil.GetExtension(de.Name())

			filename, err := generator.GenerateFilename(fp)

			fmt.Println("FILENAME: " + filename)
			if err != nil {
				log.Fatal(err)
			}

			err = imageutil.Rename(fp, filepath.Join(path, filename), ext)
			if err != nil {
				log.Fatal()
			}
			fmt.Println("New Name: " + filename)

		}

	}
}
