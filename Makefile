.PHONY: build
build:
		go build -o ./cmd/shortener ./cmd/shortener
test:
		go test ./...

run:	build
		./cmd/shortener/shortener
.DEFAULT_GOAL: build