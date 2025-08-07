package main

import (
	"log"
	"os"

	ascii "github.com/asaft29/govertor/internal"
	"golang.org/x/term"
)

func main() {
	conf, err := ascii.CreateConfig()

	if err != nil {
		log.Fatalf("ERROR : %v", err)
	}

	termW, termH, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		log.Fatalf("ERROR : %v", err)
	}

	img, err := conf.Prepare(*conf.GetInput(), termW, termH)

	if err != nil {
		log.Fatalf("ERROR : %v", err)
	}
	conf.PrintToASCII(img)
}
