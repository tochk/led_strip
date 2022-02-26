generate:
	go generate ./...

build: fmt
	GOOS=linux GOARCH=arm64 go build -o bin/led_strip cmd/led_strip/main.go

fmt:
	go fmt ./...