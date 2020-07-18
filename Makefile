.PHONY:run
.PHONY: build
.PHONY: test


run:
	go run cmd/server/main.go

build:
	go build -o server -ldflags="-s -w" cmd/main.go

test:
	go test -v