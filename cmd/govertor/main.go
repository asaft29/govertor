package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"

	ascii "github.com/asaft29/govertor/internal"
	"golang.org/x/term"
)

func main() {
	input := flag.String("i", "", "Path to input file (image or video)")
	save := flag.Bool("s", false, "Flag for saving")
	frames := flag.Int("f", 50, "Frames for video")
	output := flag.String("o", "", "Output of saved file")
	flag.Parse()

	inputs := []string{*input}
	inputs = append(inputs, flag.Args()...)

	var validInputs []string
	for _, in := range inputs {
		if in != "" {
			validInputs = append(validInputs, in)
		}
	}

	if len(validInputs) == 0 {
		log.Fatalf("ERROR : input file required")
	}
	termW, termH, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalf("ERROR : %v", err)
	}

	for _, filePath := range validInputs {
		conf, err := ascii.CreateConfig(filePath, save, frames, output)
		if err != nil {
			log.Printf("ERROR : %v", err)
			continue
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
					err = conf.PrintToASCII(img)
					if err != nil {
						log.Fatalf("ERROR : %v", err)
					}
				}
			}()

			wg.Wait()

			if v, ok := conf.(*ascii.VideoCreator); ok {
				if *v.ConfFlags.Save {
					err := v.SaveGIF()
					if err != nil {
						log.Fatalf("Failed to save GIF: %v", err)
					}
				}
			}
		}

		} else {
			img, err := conf.Prepare(*conf.GetInput(), termW, termH)
			if err != nil {
				log.Fatalf("ERROR : %v", err)
			}
			conf.PrintToASCII(img)
		}
	}
}
