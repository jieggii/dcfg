SOURCES := ./cmd/ ./internal/
MAX_LINE_LENGTH := 120

.PHONY: fmt
fmt:
	golines --max-len 120 --write-output $(SOURCES)
	gofmt -w -s $(SOURCES)
