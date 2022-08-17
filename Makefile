.PHONY: all build test install

all: build

test:
	@echo "No tests written yet..."

build:
	mkdir -p ./dist/
	go build -o ./dist/dcfg ./cmd/main.go

install: build
	install -m 755 ./dist/dcfg /usr/bin/dcfg
