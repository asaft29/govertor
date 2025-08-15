APP_NAME = govertor
MAIN = cmd/govertor/main.go

ifeq ($(OS),Windows_NT)
    RM = del /Q /F
    EXE = .exe
else
    RM = rm -f
    EXE =
endif

.PHONY: help build clean

help:
	@echo "Makefile commands:"
	@echo "  build               Build the Go binary"
	@echo "  clean               Remove built binary"

build:
	go build -o $(APP_NAME)$(EXE) $(MAIN)

clean:
	$(RM) $(APP_NAME)$(EXE)
