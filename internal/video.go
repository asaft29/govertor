package ascii

import "image"

var videoExts = map[string]bool{
	".mp4": true, ".avi": true, ".mov": true,
}

type VideoCreator struct {
	input *string
}

func (vid VideoCreator) GetInput() *string {
	return vid.input
}

func (vid VideoCreator) Prepare(filePath string, w, h int) (image.Image, error) {

	return nil, nil
}

func (vid VideoCreator) PrintToASCII(img image.Image) {

}
