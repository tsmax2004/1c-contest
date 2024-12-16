package utils

import (
	"os"
	"path/filepath"
)

func GetRelPath(url string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(pwd, url)
	if err != nil {
		return "", err
	}

	return rel, nil
}
