SOURCES := ./cmd/ ./internal/

.PHONY: fmt
fmt:
	gofmt -w -s $(SOURCES)
