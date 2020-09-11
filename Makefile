.PHONY: build run
build:
	mkdir -p  ./build
	go build -v -o build/service ./cmd/service

run:
	go run ./cmd/service/main.go

.DEFAULT_GOAL := build