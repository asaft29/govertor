package ascii

import (
	"bufio"
	"os"
	"strings"
)

type QuitHandler struct {
	quit chan bool
}

func NewQuitHandler() *QuitHandler {
	qh := &QuitHandler{
		quit: make(chan bool, 1),
	}

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(input)) == "q" {
				qh.quit <- true
				return
			}
		}
	}()

	return qh
}

func (qh *QuitHandler) ShouldQuit() bool {
	select {
	case <-qh.quit:
		return true
	default:
		return false
	}
}
