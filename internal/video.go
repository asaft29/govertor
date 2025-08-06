package ascii

var videoExts = map[string]bool{
	".mp4": true, ".mov": true, ".avi": true, ".mkv": true, ".wmv": true,
}

type VideoCreator struct {
	input *string
	w     *uint
	h     *uint
}

func (img VideoCreator) GetInput() *string {
	return img.input
}

func (img VideoCreator) GetWidth() *uint {
	return img.w
}

func (img VideoCreator) GetHeight() *uint {
	return img.h
}
