package main

import (
	"log"
	"os"
	"strconv"

	ascii "github.com/asaft29/govertor/internal"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("Usage: %s <image-path> <width> <height>", os.Args[0])
	}

	filePath := os.Args[1]

	width, err := strconv.Atoi(os.Args[2])
	if err != nil || width <= 0 {
		log.Fatalf("Invalid width: %s", os.Args[2])
	}

	height, err := strconv.Atoi(os.Args[3])
	if err != nil || height <= 0 {
		log.Fatalf("Invalid height: %s", os.Args[3])
	}

	if ascii.CheckFile(filePath) {
		grayscale, err := ascii.PrepareImage(filePath, width, height)
		if err != nil {
			log.Fatalf("ERROR : %s", err)
		}
		ascii.PrintImageToASCII(*grayscale)
	} else {
		log.Fatalf("File does not exist or is not accessible: %s", filePath)
	}
}
