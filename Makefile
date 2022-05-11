.PHONY: install clean help 
.DEFAULT: help

ENV_LIST=$(shell basename -a -s .env env/*.env)
ENV=local
BUILD_NAME=backend/sad
BUILD_TAG=latest

help:
	@echo "make install: compile packages and dependencies"
	@echo "make build: build docker image"
	@echo "make clean: remove object files and cached files"
	@echo "make tool: run specified go tool"
	@echo "make test: run all unit tests"

install:
ifeq (, $(wildcard $(shell which swag)))
	@go get -u github.com/swaggo/swag/cmd/swag
endif
	swag init
	@go build -v .
ifneq ($(findstring $(ENV),$(ENV_LIST)),)
	cp ./env/$(ENV).env .env
endif

build:
	docker build -t "$(BUILD_NAME):$(BUILD_TAG)" .

clean:
	rm -f SADBackend
	rm -rf docs
	rm -f .env
	go clean -i .

tool:
	go vet -composites=false ./...; true
	gofmt -w .
test: install
	go test -v -covermode=count -coverprofile=test.out ./...
	go tool cover -html=test.out
