.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p ./dist/
	go build -o ./dist/dcfg dcfg.go

.PHONY: install
install:
	install -m 755 ./dist/dcfg /usr/bin/dcfg
