GOPATH ?= $(HOME)/go
BIN_DIR = $(GOPATH)/bin
APP_NAME = battleship

TMPDIR ?= $(shell dirname $$(mktemp -u))
COVER_FILE ?= $(TMPDIR)/$(APP_NAME)-coverage.out

SWAGDIR = $(GOPATH)/src/github.com/swaggo

# Tools
GOLINT = $(BIN_DIR)/golint
SWAG = $(BIN_DIR)/swag

$(GOLINT):
	GO111MODULE=off go get golang.org/x/lint/golint

.PHONY: $(SWAG)
$(SWAG):
	GO111MODULE=on go install github.com/swaggo/swag/cmd/swag

.PHONY: tools
tools: $(GOLINT) ## Install all needed tools, e.g. for static checks

# Main targets

.PHONY: build
build: ## Build the project binary
	go build ./cmd/$(APP_NAME)/

.PHONY: test
test: ## Run unit tests
	go test ./... -coverprofile=$(COVER_FILE)
	go tool cover -func=$(COVER_FILE) | grep ^total

$(COVER_FILE):
	$(MAKE) test

.PHONY: cover
cover: $(COVER_FILE) ## Output coverage in human readable form in html
	go tool cover -html=$(COVER_FILE)
	rm -f $(COVER_FILE)

.PHONY: lint
lint: $(GOLINT) ## Check the project with lint
	golint -set_exit_status ./...

.PHONY: vet
vet: ## Check the project with vet
	go vet ./...

.PHONY: fmt
fmt: ## Run go fmt for the whole project
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do gofmt -e -l -w $$d/*.go; done)

.PHONY: static_check
static_check: fmt vet lint  ## Run static checks (fmt, lint, imports, vet, ...) all over the project

.PHONY: check
check: static_check test ## Check project with static checks and unit tests

.PHONY: dependencies
dependencies: ## Manage go mod dependencies, beautify go.mod and go.sum files
	go mod tidy

.PHONY: run
run: build ## Start the project
	./$(APP_NAME)

.PHONY: docs
docs: $(SWAG) ## Generate docs with go-swagl
	swag init -g cmd/$(APP_NAME)/main.go

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: dependencies-download
dependencies-download:
	go mod download
