package main

import (
	"os"

	imageutil "github.com/RajaeDev6/reena/internal/imageUtil"
	"github.com/RajaeDev6/reena/internal/namer"
	"github.com/RajaeDev6/reena/internal/openai"
)

func main() {
	imagePath := "image.png"

	api := os.Getenv("OPENAI_API_KEY")

	if api == "" {
		panic("OPENAI_API_KEY not set")
	}

	ext, imageURL := imageutil.LoadAndEncodeImage(imagePath)

	var generator namer.ImageNamer

	generator = &openai.OpenAIClient{APIKey: api}

	filename, err := generator.GenerateFilename(imageURL)
	if err != nil {
		panic(err)
	}

	imageutil.Rename(imagePath, filename, ext)

}
