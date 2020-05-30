.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ./cmd/bin/gocr ./cmd/main.go
