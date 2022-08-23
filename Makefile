.PHONY: all build install fmt

all: build

build:
	mkdir -p ./dist/
	go build -o ./dist/dcfg ./cmd/main.go

install:
	install -m 755 ./dist/dcfg /usr/bin/dcfg

fmt:
	# https://github.com/segmentio/golines
	golines --max-len 120 ./cmd/
	go fmt -w ./cmd/
