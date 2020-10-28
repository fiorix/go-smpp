UNAME := $(shell uname)
PWD = $(shell pwd)

PROFILE ?= go-smpp
DOCKER_COMPOSE_PATH = $(PWD)/docker/docker-compose.test.yml
SRC = `go list -f {{.Dir}} ./... | grep -v /vendor/`

ndef = $(if $(value $(1)),,$(error $(1) not set))

install:
	@echo "==> Installing tools..."
	@go install golang.org/x/tools/...
	@go install golang.org/x/lint/golint
	@GO111MODULE=off go get mvdan.cc/gofumpt/gofumports
	@GO111MODULE=off go get github.com/daixiang0/gci
	@brew install golangci/tap/golangci-lint
	@brew upgrade golangci/tap/golangci-lint

fmt:
	@echo "==> Formatting source code..."
	@go fmt $(SRC)
	@goimports -w $(SRC)
	@gofumports -w $(SRC)
	@-gci -w $(SRC)

lint:
	@echo "==> Running lint check..."
	@golangci-lint --config .golangci.yml run
	@golint $(SRC)
	@go vet $(SRC)

test:
	@echo "==> Running tests..."
	@go clean -testcache ./...
	@go test `go list ./... | grep -v cmd` -race --cover

test-ci-up:
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose \
		-f $(DOCKER_COMPOSE_PATH) \
		-p $(PROFILE) up \
		--force-recreate \
		--abort-on-container-exit \
		--exit-code-from app \
		--build

test-ci-down:
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose \
		-f $(DOCKER_COMPOSE_PATH) \
		-p $(PROFILE) down \
		-v --rmi local

.PHONY: fmt lint test ci-test-up ci-test-down