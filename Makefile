.PHONY: run dev binary setup glide start-mysql stop-mysql test update
SHELL := /bin/bash

all: run

run: binary
	scripts/run.sh

dev: start-mysql
	scripts/dev.sh

binary:
	GOARCH=amd64 GOOS=linux go build -i -o weather_app

setup:
	go get -v -u github.com/codegangsta/gin
	go get -v -u github.com/Masterminds/glide

glide:
	glide install --force

start-mysql:
	docker run --name weatherdb \
		-e MYSQL_ROOT_PASSWORD=blibb \
		-e MYSQL_DATABASE=weather_db \
		-e MYSQL_USER=blubb \
		-e MYSQL_PASSWORD=blabb \
		-p "3306:3306" \
		-d mariadb:10

stop-mysql:
	docker kill weatherdb || true
	docker rm -f weatherdb

test:
	GOARCH=amd64 GOOS=linux go test $$(go list ./... | grep -v /vendor/)

update:
	git fetch --all
	git merge upstream/master
