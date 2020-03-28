# -----------------------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------------------
GIT_HASH := $(shell git rev-parse HEAD)
GIT_SUMMARY := $(shell git describe --tags --dirty --always)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG := $(shell git describe --abbrev=0 --tags || echo "v0.0.0")
BUILD_STAMP := $(shell date -u '+%Y-%m-%dT%H:%M:%S%z')
VERSION := $(shell echo $(GIT_TAG) | sed -e "s/^v//")
LDFLAGS = -ldflags '-w -s \
	-X "main.ApplicationName=$(1)" \
	-X "main.ApplicationDescription=$(2)" \
	-X "main.ApplicationVersion=$(VERSION)" \
	-X "main.GitSummary=$(GIT_SUMMARY)" \
	-X "main.GitBranch=$(GIT_BRANCH)" \
	-X "main.BuildStamp=$(BUILD_STAMP)"'
.info:
	@echo GIT_SUMMARY: $(GIT_SUMMARY)
	@echo GIT_BRANCH: $(GIT_BRANCH)
	@echo BUILD_STAMP: $(BUILD_STAMP)
	@echo LDFLAGS: $(LDFLAGS)
	@echo VERSION: $(VERSION)

NAMESPACE=flaggio
DOCKER_CONTAINER_NAME=flaggio
DOCKER_REPOSITORY=$(NAMESPACE)/$(DOCKER_CONTAINER_NAME)

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

docker-build:
	docker build -t $(DOCKER_REPOSITORY):latest .

release: docker-build
	docker tag $(DOCKER_REPOSITORY):latest $(DOCKER_REPOSITORY):$(VERSION)
	docker push $(DOCKER_REPOSITORY):latest
	docker push $(DOCKER_REPOSITORY):$(VERSION)
