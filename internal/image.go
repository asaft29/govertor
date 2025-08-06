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

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
}

type ImageCreator struct {
	input *string
	w     *uint
	h     *uint
}

func (img ImageCreator) GetInput() *string {
	return img.input
}

func (img ImageCreator) GetWidth() *uint {
	return img.w
}

func (img ImageCreator) GetHeight() *uint {
	return img.h
}

const asciiChars = "@%#*+=-:."

func PrepareImage(filePath string, targetWidth, targetHeight int) (*image.Gray, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ERROR : Opening the file : %s", err)
	}

	img, _, err := image.Decode(file)

	if err != nil {
		return nil, fmt.Errorf("ERROR : Decoding : %s", err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	gray := toGrayscale(dst)

	return gray, nil
}

func toGrayscale(img image.Image) *image.Gray {

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			gray := uint8(0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8))

			grayImg.SetGray(x, y, color.Gray{Y: gray})
		}
	}
	return grayImg
}

func PrintImageToASCII(img image.Gray) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		var line string
		for x := 0; x < img.Bounds().Dx(); x++ {
			gray := img.GrayAt(x, y).Y
			index := int(float64(gray) / 256.0 * float64(len(asciiChars)))
			char := asciiChars[index]
			line += string(char)
		}
		fmt.Println(line)
	}

}
