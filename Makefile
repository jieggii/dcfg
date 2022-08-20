.PHONY: all build install

all: build

build:
	mkdir -p ./dist/
	go build -o ./dist/dcfg ./cmd/main.go

install: build
	install -m 755 ./dist/dcfg /usr/bin/dcfg
