package ascii

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const savedDirVids = "converted/videos"

var videoExts = map[string]bool{
	".mp4": true, ".avi": true, ".mov": true, ".gif": true,
	".svg": true, ".mkv": true, "flv": true, "webm": true,
}

type VideoCreator struct {
	input       *string
	cmd         *exec.Cmd
	stdout      io.ReadCloser
	reader      *bufio.Reader
	width       int
	height      int
	frameBuf    []byte
	closed      bool
	save        bool
	asciiFrames []string
}

func (vid *VideoCreator) GetInput() *string {
	return vid.input
}

func (vid *VideoCreator) GetExtension() string {
	return filepath.Ext(*vid.input)
}

func (vid *VideoCreator) IsVideo() bool {
	return true
}

func (vid *VideoCreator) Prepare(filePath string, targetWidth, targetHeight int) (image.Image, error) {
	if vid.cmd == nil {
		vid.width = targetWidth
		vid.height = targetHeight

		vid.cmd = exec.Command(
			"ffmpeg",
			"-stream_loop", "-1",
			"-i", filePath,
			"-vf", fmt.Sprintf("scale=%d:%d", targetWidth, targetHeight),
			"-f", "rawvideo",
			"-pix_fmt", "gray",
			"-",
		)
		stdout, err := vid.cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
		}
		vid.stdout = stdout
		vid.reader = bufio.NewReader(stdout)
		err = vid.cmd.Start()
		if err != nil {
			return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
		}
		vid.frameBuf = make([]byte, targetWidth*targetHeight)
		vid.closed = false
		vid.asciiFrames = []string{}
	}

	if vid.closed {
		return nil, io.EOF
	}

	n, err := io.ReadFull(vid.reader, vid.frameBuf)
	if err != nil {
		vid.closed = true
		vid.cmd.Wait()
		return nil, err
	}

	if n != vid.width*vid.height {
		return nil, fmt.Errorf("incomplete frame read")
	}

	img := image.NewGray(image.Rect(0, 0, vid.width, vid.height))
	copy(img.Pix, vid.frameBuf)
	return img, nil
}

func (vid *VideoCreator) PrintToASCII(img image.Image) error {
	bounds := img.Bounds()
	var b strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
			idx := int(float64(gray) / 256 * float64(len(asciiChars)))
			if idx >= len(asciiChars) {
				idx = len(asciiChars) - 1
			}
			b.WriteByte(asciiChars[idx])
		}
		b.WriteByte('\n')
	}

	asciiFrame := b.String()

	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	fmt.Print("\033[3J")
	fmt.Print(asciiFrame)
	os.Stdout.Sync()

	if vid.save && len(vid.asciiFrames) < 50 {
		vid.asciiFrames = append(vid.asciiFrames, asciiFrame)
	}

	time.Sleep(time.Millisecond * 50)
	return nil
}

func (vid *VideoCreator) SaveGIF() error {
	if len(vid.asciiFrames) == 0 {
		return fmt.Errorf("no frames to save")
	}

	var images []*image.Paletted
	var delays []int

	for _, frame := range vid.asciiFrames {
		img := vid.textToImage(frame)
		images = append(images, img)
		delays = append(delays, 10)
	}

	baseName := strings.TrimSuffix(filepath.Base(*vid.input), filepath.Ext(*vid.input))
	outputPath := filepath.Join(savedDirVids, baseName+"_ascii.gif")

	os.MkdirAll(savedDirVids, 0755)
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gifAnim := &gif.GIF{Image: images, Delay: delays}
	err = gif.EncodeAll(file, gifAnim)
	if err != nil {
		return err
	}

	fmt.Printf("GIF saved: %s (%d frames)\n", outputPath, len(images))
	return nil
}

func (vid *VideoCreator) textToImage(text string) *image.Paletted {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil
	}

	face := basicfont.Face7x13
	charWidth := 7
	charHeight := 13

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	width := charWidth * maxWidth
	height := charHeight * len(lines)

	palette := color.Palette{
		color.White,
		color.Black,
	}

	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

	draw.Draw(img, img.Bounds(), &image.Uniform{palette[0]}, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{palette[1]},
		Face: face,
		Dot:  fixed.P(0, charHeight-2),
	}

	for y, line := range lines {
		d.Dot = fixed.P(0, (y+1)*charHeight-2)
		d.DrawString(line)
	}

	return img
}
