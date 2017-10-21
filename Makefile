PHONY: run gin binary setup

all: run

run:
	go run main.go

gin:
	gin --all run main.go

binary:
	go build -o weatherapp

setup:
	go get -v -u github.com/codegangsta/gin
	go get -v -u github.com/Masterminds/glide
