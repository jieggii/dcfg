.PHONY: all build

all: build

build:
	mkdir -p ./dist/
	go build -o ./dist/dcfg ./cmd/main.go

