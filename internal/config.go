package ascii

import (
	"errors"
	"flag"
	"image"
	"path/filepath"
)

const asciiChars = ".%#*+=-:^"

type Config interface {
	GetInput() *string
	GetExtension() string
	Prepare(filePath string, w, h int) (image.Image, error)
	PrintToASCII(img image.Image) error
	IsVideo() bool
}

func CreateConfig() (Config, error) {
	input := flag.String("input", "", "Path to input file (image or video)")
	save := flag.Bool("save", false, "Flag for saving")

	flag.Parse()

	if *input == "" {
		return nil, errors.New("input file required")
	}

	ext := filepath.Ext(*input)

	switch {
	case imageExts[ext]:
		return &ImageCreator{input: input, save: *save}, nil
	case videoExts[ext]:
		return &VideoCreator{input: input, save: *save}, nil
	default:
		return nil, errors.New("unsupported file type: " + ext)
	}
}
