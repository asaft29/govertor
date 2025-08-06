package ascii

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
)

func CheckFile(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

type Config interface {
	GetInput() *string
	GetWidth() *uint
	GetHeight() *uint
}

func CreateConfig() (Config, error) {
	input := flag.String("input", "", "Path to input file (image or video)")
	w := flag.Uint("width", 0, "Width")
	h := flag.Uint("height", 0, "Height")

	flag.Parse()

	if *input == "" {
		return nil, errors.New("input file required")
	}

	if *w == 0 || *h == 0 {
		return nil, errors.New("width and height must be both > 0")
	}

	ext := filepath.Ext(*input)

	switch {
	case imageExts[ext]:
		return &ImageCreator{input: input, w: w, h: h}, nil
	case videoExts[ext]:
		return &VideoCreator{input: input, w: w, h: h}, nil
	default:
		return nil, errors.New("unsupported file type: " + ext)
	}
}
