APP_NAME = govertor
MAIN = cmd/govertor/main.go

.PHONY: help build run clean

help:
	@echo "Makefile commands:"
	@echo "  build               Build the Go binary"
	@echo "  run INPUT W H       Run the program with input file and width & height"
	@echo "  clean               Remove built binary"

build:
	go build -o $(APP_NAME) $(MAIN)

run:
ifdef INPUT
	@echo "Running with input=$(INPUT), width=$(W), height=$(H)"
	go run $(MAIN) $(INPUT) $(W) $(H)
else
	@echo "Error: Please provide INPUT, W, and H variables"
	@echo "Example: make run INPUT=input.txt W=640 H=480"
	exit 1
endif

clean:
	rm -f $(APP_NAME)

