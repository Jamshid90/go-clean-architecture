.PHONY: run, build

.DEFAULT_GOAL = build

build:
	go build -o server -ldflags="-s -w" cmd/http/main.go

run:
	go run cmd/http/main.go