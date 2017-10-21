PHONY: run gin binary

all: run

run:
	go run main.go

gin:
	gin --all run main.go

binary:
	go build -o weatherapp
