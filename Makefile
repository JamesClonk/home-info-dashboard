.PHONY: run gin binary setup glide test
SHELL := /bin/bash

all: run

run:
	source .env
	go run main.go

gin:
	source .env
	gin --all run main.go

binary:
	GOARCH=amd64 GOOS=linux go build -o weather_app

setup:
	go get -v -u github.com/codegangsta/gin
	go get -v -u github.com/Masterminds/glide

glide:
	glide install --force

test:
	GOARCH=amd64 GOOS=linux go test $$(go list ./... | grep -v /vendor/)
