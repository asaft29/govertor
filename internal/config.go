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
	Prepare(filePath string, w, h int) (image.Image, error)
	PrintToASCII(img image.Image) error
	IsVideo() bool
}

type Flags struct {
	Input  *string
	Save   *bool
	Frames *int
	Output *string
}

func CreateConfig() (ASCIIConvertor, error) {
	input := flag.String("i", "", "Path to input file (image or video)")
	save := flag.Bool("s", false, "Flag for saving")
	frames := flag.Int("f", 50, "Frames for video")
	output := flag.String("o", "", "Output of saved file")
	flag.Parse()

	if *input == "" {
		return nil, errors.New("input file required")
	}

	if *frames > maxFrames || *frames < 0 {
		return nil, errors.New("max frames should be [0, 600]")
	}

	ext := filepath.Ext(*input)
	flags := Flags{input, save, frames, output}

	switch {
	case imageExts[ext]:
		return &ImageCreator{input: input, ConfFlags: flags}, nil
	case videoExts[ext]:

		return &VideoCreator{input: input, ConfFlags: flags}, nil
	default:
		return nil, errors.New("unsupported file type: " + ext)
	}
}
