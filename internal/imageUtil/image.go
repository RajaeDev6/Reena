package imageutil

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

func LoadAndEncodeImage(path string) (string, string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", ""
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	ext := getExtension(path)

	imageURL := "data:image/" + ext + ";base64," + encoded
	return ext, imageURL
}

func getExtension(file string) string {
	return strings.TrimPrefix(filepath.Ext(file), ".")
}

func Rename(oldName string, newName string, ext string) {
	err := os.Rename(oldName, newName+"."+ext)

	if err != nil {
		panic(err)
	}
}
