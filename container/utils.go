package container

import (
	"errors"
	"os"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func volumeUrlExtract(volume string) []string {
	volumeURLs := strings.Split(volume, ":")
	return volumeURLs
}
