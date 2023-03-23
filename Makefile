.PHONY: up
up:
	docker compose up --remove-orphans -d --build

.PHONY: down
down:
	docker compose -p hmm down

PROJECT_NAME := hmm
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
CONTAINER_DIR := /go/src/github.com/dmartzol/$(PROJECT_NAME)
.PHONY: lint-ci
lint-ci:
	echo $(ROOT) && \
	docker run \
	-v $(ROOT):$(CONTAINER_DIR) \
	-w $(CONTAINER_DIR)/ \
	--rm \
	-t golangci/golangci-lint:v1.50 \
	golangci-lint run -v --timeout 5m0s ./...
