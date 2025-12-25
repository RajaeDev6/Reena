package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/RajaeDev6/reena/internal/dir"
	imageutil "github.com/RajaeDev6/reena/internal/imageUtil"
	"github.com/RajaeDev6/reena/internal/namer"
	provider "github.com/RajaeDev6/reena/internal/openai"
)

func main() {

	api := os.Getenv("OPENAI_API_KEY")
	var dirMode bool
	var generator namer.ImageNamer

	if api == "" {
		panic("OPENAI_API_KEY not set")
	}

	generator = provider.NewOpenAIGenerator(api)

	flag.BoolVar(&dirMode, "d", false, "treat argument as directory")
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage:")
		fmt.Fprintln(os.Stderr, "  reena <file>")
		fmt.Fprintln(os.Stderr, "  reena -d <directory>")
		os.Exit(1)
	}

	path := args[0]

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if dirMode {
		if !info.IsDir() {
			fmt.Fprintln(os.Stderr, "error: -d flag requires a directory")
			os.Exit(1)
		}

		dir.Scan(path, generator)

		os.Exit(0)
	}

	filename, err := generator.GenerateFilename(path)

	if err != nil {
		log.Fatal(err)
	}

	ext := imageutil.GetExtension(path)

	err = imageutil.Rename(path, filename, ext)
	if err != nil {
		log.Fatal(err)
	}
}
