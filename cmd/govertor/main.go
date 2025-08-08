package main

import (
	"io"
	"log"
	"os"
	"sync"

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

	if conf.IsVideo() {
		quitHandler := ascii.NewQuitHandler()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			for {
				if quitHandler.ShouldQuit() {
					return
				}

				img, err := conf.Prepare(*conf.GetInput(), termW, termH)
				if err != nil {
					if err == io.EOF {
						return
					}
					log.Printf("ERROR : %v", err)
					return
				}
				conf.PrintToASCII(img)
			}
		}()

		wg.Wait()

	} else {
		img, err := conf.Prepare(*conf.GetInput(), termW, termH)
		if err != nil {
			log.Fatalf("ERROR : %v", err)
		}
		conf.PrintToASCII(img)
	}
}
