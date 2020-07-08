.PHONY:run
run:
	go run cmd/server/main.go

.PHONY: build
build:
	go build -o server -ldflags="-s -w" cmd/main.go