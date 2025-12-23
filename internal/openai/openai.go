package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OpenAIClient struct {
	APIKey string
}

func (c *OpenAIClient) GenerateFilename(imageURL string) (string, error) {
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
	}`, imageURL)

	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/responses",
		strings.NewReader(reqBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return extractFilename(body), nil
}

func extractFilename(body []byte) string {
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	return data["output"].([]any)[0].(map[string]any)["content"].([]any)[0].(map[string]any)["text"].(string)
}
