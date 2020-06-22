.DEFAULT_GOAL := help
SHELL := /bin/bash
APP ?= $(shell basename $$(pwd))
COMMIT_SHA = $(shell git rev-parse HEAD)

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: run
## run: runs main.go with local sqlite database
run:
	@cp _fixtures/test.db _fixtures/temp.db
	source .env_sqlite; go run -race main.go

.PHONY: dev
## dev: runs main.go on a local postgres
dev:
	source .env; source .env_dev; go run -race main.go

.PHONY: gin
## gin: runs main.go via gin on a local postgres
gin:
	source .env; source .env_dev; gin --all --immediate --appPort 9099 --port 9090 run main.go

.PHONY: build
## build: builds the application
build: clean
	@echo "Building binary ..."
	go build -o ${APP}

.PHONY: clean
## clean: cleans up binary files
clean:
	@echo "Cleaning up ..."
	@go clean

.PHONY: test
## test: runs go test with the race detector
test:
	@cp _fixtures/test.db _fixtures/temp.db
	@source .env_sqlite; GOARCH=amd64 GOOS=linux go test -v -race ./...

.PHONY: init
## init: sets up go modules
init:
	@echo "Setting up modules ..."
	@go mod init 2>/dev/null; go mod tidy && go mod vendor

.PHONY: push
## push: pushes the application onto CF
push: test build
	cf push

.PHONY: db-admin
## db-admin: opens up database administration gui
db-admin:
	dbeaver &

.PHONY: postgres
## postgres: runs postgres backend on docker
postgres: postgres-network postgres-stop postgres-start
	docker logs postgres -f

.PHONY: postgres-network
postgres-network:
	docker network create postgres-network --driver bridge || true

.PHONY: postgres-cleanup
## postgres-cleanup: cleans up postgres backend
postgres-cleanup: postgres-stop
.PHONY: postgres-stop
postgres-stop:
	docker rm -f postgres || true

.PHONY: postgres-start
postgres-start:
	docker run --name postgres \
		--network postgres-network \
		-e POSTGRES_USER='dev-user' \
		-e POSTGRES_PASSWORD='dev-secret' \
		-e POSTGRES_DB='home_info_db' \
		-p 5432:5432 \
		-d postgres:9-alpine
	scripts/db_setup.sh

.PHONY: postgres-client
## postgres-client: connects to postgres backend with CLI
postgres-client:
	docker exec -it \
		-e PGPASSWORD='dev-secret' \
		postgres psql -U 'dev-user' -d 'home_info_db'

.PHONY: cleanup
cleanup: docker-cleanup
.PHONY: docker-cleanup
## docker-cleanup: cleans up local docker images and volumes
docker-cleanup:
	docker system prune --volumes -a

########################################################################################################################
####### docker/kubernetes related stuff ################################################################################
########################################################################################################################
.PHONY: image-login
## image-login: login to docker hub
image-login:
	@export PATH="$$HOME/bin:$$PATH"
	@echo $$DOCKER_PASS | docker login -u $$DOCKER_USER --password-stdin

.PHONY: build image-build
## image-build: build docker image
image-build:
	@export PATH="$$HOME/bin:$$PATH"
	docker build -t jamesclonk/${APP}:${COMMIT_SHA} .

.PHONY: image-publish
## image-publish: build and publish docker image
image-publish:
	@export PATH="$$HOME/bin:$$PATH"
	docker push jamesclonk/${APP}:${COMMIT_SHA}
	docker tag jamesclonk/${APP}:${COMMIT_SHA} jamesclonk/${APP}:latest
	docker push jamesclonk/${APP}:latest

.PHONY: image-run
## image-run: run docker image
image-run:
	@export PATH="$$HOME/bin:$$PATH"
	docker run --rm --env-file .dockerenv -p 9090:9090 jamesclonk/${APP}:${COMMIT_SHA}
