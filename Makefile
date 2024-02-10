APP_NAME        := web-stream-recorder

SHELL           = /usr/bin/env bash -eu

TAG             ?= $(shell git describe --abbrev=1 --tags --always)

DOCKER_REGISTRY 		?= noredine69/public
DOCKER_USERNAME         ?= noredine69
APP_BUILDER_IMAGE		?= golang:1.17.5-buster


TAG             ?= $(shell git describe --abbrev=1 --tags --always)
IMAGE_CURRENT   := $(DOCKER_REGISTRY):$(APP_NAME)-$(TAG)
IMAGE_LATEST    := $(DOCKER_REGISTRY):$(APP_NAME)-latest

APP_DIR    		?= $(shell pwd -P)

.PHONY: help all compile lint test build push push-latest

default: help

help: ## Print available Makefile commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: | compile lint test build push push-latest

compile: ## Build Go application
	@docker run --rm \
		-v $(APP_DIR):/$(APP_NAME) \
		$(APP_BUILDER_IMAGE) /$(APP_NAME)/.scripts/compile.sh $(APP_NAME) $(APP_NAME)

lint: ## Build Go application
	@docker run --rm \
		-v $(APP_DIR):/$(APP_NAME) \
		$(APP_BUILDER_IMAGE) /$(APP_NAME)/.scripts/lint.sh $(APP_NAME)

test: ## Build Go application
	@docker run --rm \
		-v $(APP_DIR):/$(APP_NAME) \
		$(APP_BUILDER_IMAGE) /$(APP_NAME)/.scripts/test.sh $(APP_NAME)

build: ## Build the docker image
	@docker image build \
			--build-arg BINARY_NAME_ARG=$(APP_NAME) \
			--pull --no-cache -t $(IMAGE_CURRENT) .
	@docker image tag $(IMAGE_CURRENT) $(IMAGE_LATEST)

push: | repo-login ## Push the docker image to repository
	@docker image push $(IMAGE_CURRENT)

push-latest: | repo-login ## Push the docker image latest tag to repository
	@docker image push $(IMAGE_LATEST)

start:
	@docker run -p 8080:8080 $(IMAGE_CURRENT)

start-latest:
	@docker run -p 8080:8080 $(IMAGE_LATEST)

pull:
	@docker pull $(IMAGE_CURRENT)
	@docker pull $(IMAGE_LATEST)

repo-login:
	@docker login -u $(DOCKER_USERNAME) --password-stdin < .docker_hub_token