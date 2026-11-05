VERSION := 0.1.0

.PHONY: build

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.VERSION=$(VERSION)"
