package ascii

import (
	"errors"
	"flag"
	"image"
	"path/filepath"
)

type Config interface {
	GetInput() *string
	Prepare(filePath string, w, h int) (image.Image, error)
	PrintToASCII(img image.Image)
}

func CreateConfig() (Config, error) {
	input := flag.String("input", "", "Path to input file (image or video)")

	flag.Parse()

	if *input == "" {
		return nil, errors.New("input file required")
	}

	ext := filepath.Ext(*input)

	switch {
	case imageExts[ext]:
		return &ImageCreator{input: input}, nil
	case videoExts[ext]:
		return &VideoCreator{input: input}, nil
	default:
		return nil, errors.New("unsupported file type: " + ext)
	}
}
