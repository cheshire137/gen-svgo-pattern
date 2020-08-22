all: build

default: build

build:
	go build -o bin/gen-svgo-pattern gen-svgo-pattern.go
