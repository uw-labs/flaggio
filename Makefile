# -----------------------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------------------
GIT_SUMMARY := $(shell git describe --tags --dirty --always)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_STAMP := $(shell date -u '+%Y-%m-%dT%H:%M:%S%z')
LDFLAGS = -ldflags '-s \
	-X "github.com/victorkohl/flaggio/main.ApplicationName=$(1)" \
	-X "github.com/victorkohl/flaggio/main.ApplicationDescription=$(2)" \
	-X "github.com/victorkohl/flaggio/main.GitSummary=$(GIT_SUMMARY)" \
	-X "github.com/victorkohl/flaggio/main.GitBranch=$(GIT_BRANCH)" \
	-X "github.com/victorkohl/flaggio/main.BuildStamp=$(BUILD_STAMP)"'
.info:
	@echo GIT_SUMMARY: $(GIT_SUMMARY)
	@echo GIT_BRANCH: $(GIT_BRANCH)
	@echo BUILD_STAMP: $(BUILD_STAMP)
	@echo LDFLAGS: $(LDFLAGS)

NAMESPACE=victorkohl
DOCKER_REGISTRY=hub.docker.com
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

build: CGO_ENABLED=0 ## Build project
build:
	rm -rf bin
	mkdir bin
	go build -a $(call LDFLAGS,admin,Manage flaggio application) -o bin/admin ./cmd/admin/
	go build -a $(call LDFLAGS,flaggio,Evaluate user context and return flag values) -o bin/flaggio ./cmd/flaggio/

gqlgen:
	go run github.com/99designs/gqlgen -v -c gqlgen.admin.yml
	go run github.com/99designs/gqlgen -v -c gqlgen.api.yml

gen:
	go generate ./...
