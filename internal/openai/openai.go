package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	imageutil "github.com/RajaeDev6/reena/internal/imageUtil"
)

type OpenAIClient struct {
	APIKey string
}

func NewOpenAIGenerator(APIKey string) *OpenAIClient {
	return &OpenAIClient{
		APIKey: APIKey,
	}
}

type Response struct {
	Output []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func buildRequestBody(imageURL string) string {
	return fmt.Sprintf(`{
	  "model": "gpt-4.1-mini",
	  "input": [
	    {
	      "role": "user",
	      "content": [
	        {
	          "type": "input_text",
	          "text": "Look at this image and return a short descriptive filename that matches whats inside the image. Use exactly 3 words. Rules: lowercase only, underscores instead of spaces, no file extension, no punctuation. Return only the name."
	        },
	        {
	          "type": "input_image",
	          "image_url": "%s"
	        }
	      ]
	    }
	  ]
	}`, imageURL)
}

func (c *OpenAIClient) GenerateFilename(path string) (string, error) {
	client := &http.Client{}

	encoded, err := imageutil.LoadAndEncodeImage(path)
	if err != nil {
		return "", err
	}

	ext := imageutil.GetExtension(path)

	imageURL := "data:image/" + ext + ";base64," + encoded

	reqBody := buildRequestBody(imageURL)

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

	if resp.StatusCode != http.StatusOK {
		var er ErrorResponse
		if err := json.Unmarshal(body, &er); err != nil {
			return "", fmt.Errorf("openai error: %s", string(body))
		}
		return "", fmt.Errorf("openai error: %s", er.Error.Message)
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		return "", err
	}

	if len(r.Output) == 0 ||
		len(r.Output[0].Content) == 0 ||
		r.Output[0].Content[0].Text == "" {
		return "", fmt.Errorf("empty model output")
	}

	return r.Output[0].Content[0].Text, nil
}
