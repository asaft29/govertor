package utils

import (
	"errors"
	"os"
)

func CheckFile(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}
