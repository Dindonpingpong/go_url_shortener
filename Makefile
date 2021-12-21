.PHONY: build
build:
		go build -v ./cmd/shortener

.DEFAULT_GOAL: build