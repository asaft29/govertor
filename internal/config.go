// Package ascii provides utilities for converting images and videos
// into ASCII art representations, including saving ASCII frames as GIFs.
package ascii

import (
	"errors"
	"flag"
	"image"
	"path/filepath"
)

const asciiChars = ".%#*+=-:^"
const charWidth = 7
const charHeight = 13

const maxFrames = 600

type ASCIIConvertor interface {
	GetInput() *string
	GetExtension() string
	Prepare(filePath string, w, h int) (image.Image, error)
	PrintToASCII(img image.Image) error
	IsVideo() bool
}

type Flags struct {
	Input  *string
	Save   *bool
	Frames *int
}

func CreateConfig() (ASCIIConvertor, error) {
	input := flag.String("i", "", "Path to input file (image or video)")
	save := flag.Bool("s", false, "Flag for saving")
	frames := flag.Int("f", 50, "Frames for video")
	flag.Parse()

	if *input == "" {
		return nil, errors.New("input file required")
	}

	if *frames > maxFrames || *frames < 0 {
		return nil, errors.New("max frames should be [0, 600]")
	}

	ext := filepath.Ext(*input)

	switch {
	case imageExts[ext]:
		return &ImageCreator{input: input, save: *save}, nil
	case videoExts[ext]:
		flags := Flags{input, save, frames}
		return &VideoCreator{input: input, ConfFlags: flags}, nil
	default:
		return nil, errors.New("unsupported file type: " + ext)
	}
}
