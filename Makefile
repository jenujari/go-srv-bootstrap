tidy:
	go mod tidy

build:tidy
	go build -o server.exe main.go

dev:
	air

run:build
	.\server.exe

hello:
	echo "Hello Hi"