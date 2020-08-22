PROJECTNAME=jira2trello
VERSION ?= $(shell git describe --tags --always --match=v* || echo v0)
LDFLAGS=-ldflags "-X=main.version=$(VERSION)"
LINTERFLAGS=--enable-all --disable gochecknoinits --disable gochecknoglobals --disable goimports --disable wsl --out-format=tab --tests=false
BUILDFLAGS=$(LDFLAGS)
GOEXE := $(shell go env GOEXE)
GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
BIN=bin/$(PROJECTNAME)$(GOEXE)
MODULE=github.com/Brialius/$(PROJECTNAME)
LINT_PATH := ./bin/golangci-lint
LINT_PATH_WIN := golangci-lint
LINT_SETUP := curl -sfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s latest

# -race doesn't work in Windows
ifneq ($(GOOS), windows)
	RACE = -race
endif

ifeq ($(GOOS), windows)
	LINT_PATH := $(LINT_PATH_WIN)
	LINT_SETUP := go install github.com/golangci/golangci-lint/cmd/golangci-lint
endif

export

.PHONY: setup
setup: ## Install all the build and lint dependencies
	$(LINT_SETUP)
	go get ./...

.PHONY: test
test: ## Run all the tests
	go test -cover $(RACE) -v $(BUILDFLAGS) ./...

.PHONY: lint
lint: ## Run all the linters
	$(LINT_PATH) run $(LINTERFLAGS) ./...

.PHONY: ci
ci: setup lint build test ## Run all the tests and code checks

.PHONY: build
build: clean mod-refresh

.PHONY: build
build: mod-refresh ## Build a version
	go build $(BUILDFLAGS) -o $(BIN) $(MODULE)

.PHONY: install
install: mod-refresh ## Install a binary
	go install $(BUILDFLAGS)

.PHONY: clean
clean: ## Remove temporary files
	go clean $(BUILDFLAGS) $(MODULE)

.PHONY: mod-refresh
mod-refresh: ## Refresh modules
	go mod tidy -v

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: release
release:
	git tag $(ver)

.DEFAULT_GOAL := build
