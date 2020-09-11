.PHONY: build
build:
	mkdir -p  ./build
	go build -v -o build/service ./cmd/service

.DEFAULT_GOAL := build