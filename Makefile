VERSION?= v$(shell git rev-list HEAD --count)
ARCH?= $(shell uname -m)
SYSTEM?= $(shell uname)
PROJECT_NAME?= feature-toggle
export CGO_ENABLED=0
export IMAGE_TAG=${PROJECT_NAME}_${ARCH}:${VERSION}

.PHONY: mod build docker docker-compose-up docker-compose-down clean vet fmt

mod:
	@echo "+ $@"
	go mod verify
	go mod vendor
	go mod tidy

build:
	@echo "+ $@"
	@go build -o $(BUILD)/$(PROJECT_NAME) $(CURDIR)/cmd/$(PROJECT_NAME)

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v sandbox)

fmt:
	@echo "+ $@"
	@test -z "$$(gofmt -s -l . 2>&1 | grep -v ^vendor/ | tee /dev/stderr)" || \
		(echo >&2 "+ please format Go code with 'gofmt -s'" && false)

docker:
	@echo "+ $@"
	@docker build -t ${PROJECT_NAME}_${ARCH}:${VERSION} $(CURDIR)/.

docker-compose-up:
	@echo "+ $@"
	@docker-compose -f $(CURDIR)/docker-compose.yml up -d

docker-compose-down:
	@echo "+ $@"
	@docker-compose -f $(CURDIR)/docker-compose.yml down -v

clean:
	@echo "+ $@"
	@go clean -cache -testcache

.DEFAULT_GOAL := build
