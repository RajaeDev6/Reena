package main

import (
	"fmt"
	"os"
)

func main() {
	image_path := "/home/metrix/Pictures/screenshot-2025-12-12-132414.png"

	image_content, err := os.ReadFile(image_path)

	if err != nil {
		panic("an error occured")
	}

	fmt.Println(image_content)

}
