// ascii/image_creator.go
package ascii

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"golang.org/x/image/draw"
)

const asciiChars = "@%#*+=-:."

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
}

type ImageCreator struct {
	input *string
}

func (ic *ImageCreator) GetInput() *string { return ic.input }

func (ic *ImageCreator) Prepare(filePath string, targetWidth, targetHeight int) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return toGrayscale(dst), nil
}

func toGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			grayImg.SetGray(x, y, color.Gray{Y: gray})
		}
	}
	return grayImg
}

func (ic *ImageCreator) PrintToASCII(img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var line string
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
			idx := int(float64(gray) / 256 * float64(len(asciiChars)))
			if idx >= len(asciiChars) {
				idx = len(asciiChars) - 1
			}
			line += string(asciiChars[idx])
		}
		fmt.Println(line)
	}
}
