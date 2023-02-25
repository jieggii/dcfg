BIN_DEST := /usr/bin/dcfg

.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p ./dist/
	go build -o ./dist/dcfg dcfg.go

.PHONY: install
install:
	install -m 755 ./dist/dcfg $(BIN_DEST)

.PHONY: uninstall
uninstall:
	rm -f $(BIN_DEST)