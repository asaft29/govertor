APP_NAME = govertor
MAIN = cmd/govertor/main.go

.PHONY: help build clean

help:
	@echo "Makefile commands:"
	@echo "  build               Build the Go binary"
	@echo "  clean               Remove built binary"

build:
	go build -o $(APP_NAME) $(MAIN)

clean:
	rm -f $(APP_NAME)

