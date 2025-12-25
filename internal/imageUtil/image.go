package imageutil

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadAndEncodeImage(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	return encoded, err
}

func GetExtension(file string) string {
	return strings.TrimPrefix(filepath.Ext(file), ".")
}

func IsImage(file string) bool {
	ext := GetExtension(file)

	switch ext {
	case ".png", ".jpg":
		return false
	default:
		return true
	}
}

func Rename(oldName string, newName string, ext string) error {
	fmt.Println("Renaming File...")
	err := os.Rename(oldName, newName+"."+ext)

	if err != nil {
		return err
	}

	return nil
}
