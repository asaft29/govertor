package ascii

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const savedDirImgs = "converted/images"

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".bmp": true, ".webp": true, ".ico": true,
}

type ImageCreator struct {
	input     *string
	ConfFlags Flags
}

func (ic *ImageCreator) GetInput() *string {
	return ic.input
}

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

func (ic *ImageCreator) IsVideo() bool {
	return false
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

func (ic *ImageCreator) PrintToASCII(img image.Image) error {
	bounds := img.Bounds()
	var asciiBuilder strings.Builder
	var asciiLines []string

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var line strings.Builder
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
			idx := int(float64(gray) / 256 * float64(len(asciiChars)))
			if idx >= len(asciiChars) {
				idx = len(asciiChars) - 1
			}
			char := asciiChars[idx]
			asciiBuilder.WriteByte(char)
			line.WriteByte(char)
		}
		asciiBuilder.WriteByte('\n')
		asciiLines = append(asciiLines, line.String())
	}

	asciiArt := asciiBuilder.String()
	fmt.Print(asciiArt)

	if *ic.ConfFlags.Save {
		res := ic.saveAsPNG(asciiLines)
		if res != nil {
			log.Fatalf("Failed to save IMAGE: image will not be saved : %v", res)
		}
		return res
	}
	return nil
}

func (ic *ImageCreator) saveAsPNG(asciiLines []string) error {

	imgWidth := len(asciiLines[0]) * charWidth
	imgHeight := len(asciiLines) * charHeight

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for y := range imgHeight {
		for x := range imgWidth {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: basicfont.Face7x13,
	}

	for i, line := range asciiLines {
		d.Dot = fixed.P(0, (i+1)*charHeight)
		d.DrawString(line)
	}

	baseName := strings.TrimSuffix(filepath.Base(*ic.input), filepath.Ext(*ic.input))
	var outputPath string

	if ic.ConfFlags.Output == nil || *ic.ConfFlags.Output == "" {
		outputPath = filepath.Join(savedDirImgs, baseName+"_ascii.png")
	} else {
		out := strings.TrimSpace(*ic.ConfFlags.Output)

		if fi, err := os.Stat(out); err == nil {
			if fi.IsDir() {
				outputPath = filepath.Join(out, baseName+"_ascii.png")
			} else {
				outputPath = out
			}
		} else if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("output path does not exist: %s", out)
		} else {
			return fmt.Errorf("error accessing output path: %v", err)
		}
	}

	os.MkdirAll(savedDirImgs, 0755)
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	fmt.Printf("ASCII art saved as PNG: %s\n", outputPath)
	return nil
}
