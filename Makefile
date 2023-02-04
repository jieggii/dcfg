SOURCES := ./cmd/ ./internal/


.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p ./dist/
	go build -o ./dist/dcfg ./cmd/dcfg/dcfg.go

.PHONY: install
install:
	install -m 755 ./dist/dcfg /usr/bin/dcfg

.PHONY: fmt
fmt:
	gofmt -w -s $(SOURCES)
