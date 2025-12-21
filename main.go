package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	image_path := "image.png"

	api := os.Getenv("OPENAI_API_KEY")

	if api == "" {
		panic("OPENAI_API_KEY not set")
	}

	image_content, err := os.ReadFile(image_path)
	if err != nil {
		panic(err)
	}

	encoded := base64.StdEncoding.EncodeToString(image_content)

	ext := strings.TrimPrefix(filepath.Ext(image_path), ".")
	imgURL := "data:image/" + ext + ";base64," + encoded

	client := &http.Client{}

	reqBody := fmt.Sprintf(`{
	  "model": "gpt-4.1-mini",
	  "input": [{
	    "role": "user",
	    "content": [
	      {
	        "type": "input_text",
	        "text": "Look at this image and return a short descriptive filename that matches whats inside the image use exactly 3 words. Rules: lowercase only, underscores instead of spaces, no file extension, 3 words, no punctuation, return only the name."
	      },
	      {
	        "type": "input_image",
	        "image_url": "%s"
	      }
	    ]
	  }]
	}`, imgURL)

	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/responses",
		strings.NewReader(reqBody),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]any
	json.Unmarshal(body, &data)

	filename := data["output"].([]any)[0].(map[string]any)["content"].([]any)[0].(map[string]any)["text"].(string)

	err = os.Rename(image_path, filename+"."+ext)

	if err != nil {
		log.Fatal(err)
	}
}
