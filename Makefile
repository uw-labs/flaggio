# -----------------------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------------------
GIT_SUMMARY := $(shell git describe --tags --dirty --always)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_STAMP := $(shell date -u '+%Y-%m-%dT%H:%M:%S%z')
LDFLAGS = -ldflags '-w -s \
	-X "main.ApplicationName=$(1)" \
	-X "main.ApplicationDescription=$(2)" \
	-X "main.GitSummary=$(GIT_SUMMARY)" \
	-X "main.GitBranch=$(GIT_BRANCH)" \
	-X "main.BuildStamp=$(BUILD_STAMP)"'
.info:
	@echo GIT_SUMMARY: $(GIT_SUMMARY)
	@echo GIT_BRANCH: $(GIT_BRANCH)
	@echo BUILD_STAMP: $(BUILD_STAMP)
	@echo LDFLAGS: $(LDFLAGS)

NAMESPACE=flaggio
DOCKER_REGISTRY=docker.io
DOCKER_CONTAINER_NAME=flaggio
DOCKER_REPOSITORY=$(DOCKER_REGISTRY)/$(NAMESPACE)/$(DOCKER_CONTAINER_NAME)

# -----------------------------------------------------------------------------------------
# Application Tasks
# -----------------------------------------------------------------------------------------

install: ## Install dependencies
	go get -t -d -v ./...

lint: ## Run linting
	go vet ./...
	[[ -z "$(shell gofmt -e -l .)" ]] || exit 1 # gofmt

test: ## Run tests
	go test -race -v ./...

build:
	rm -rf bin && \
	mkdir bin && \
	CGO_ENABLED=0 go build -a $(call LDFLAGS,flaggio,Self hosted feature flag solution) -o bin/flaggio ./cmd/flaggio/

gen:
	go generate ./...
