package ascii

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var videoExts = map[string]bool{
	".mp4": true, ".avi": true, ".mov": true, ".gif": true,
}

type VideoCreator struct {
	input    *string
	cmd      *exec.Cmd
	stdout   io.ReadCloser
	reader   *bufio.Reader
	width    int
	height   int
	frameBuf []byte
	closed   bool
}

func (vid *VideoCreator) GetInput() *string {
	return vid.input
}

func (vid *VideoCreator) IsVideo() bool { return true }

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

func (vid *VideoCreator) PrintToASCII(img image.Image) {

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

	// I TRIED SO MANY OTHER THINGS
	// BUT THIS IS THE ONLY ONE THAT DOESN'T BREAK THE TERMINAL

	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	fmt.Print("\033[3J")

	fmt.Print(b.String())

	os.Stdout.Sync()
	time.Sleep(time.Millisecond * 50)
}
