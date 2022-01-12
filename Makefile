.PHONY: build
build:
		go build -o ./cmd/shortener ./cmd/shortener
test:
		go test ./...
.DEFAULT_GOAL: build