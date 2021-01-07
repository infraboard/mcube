MAIN_FILE := "main.go"
PROJECT_NAME := "mcube"
PKG := "github.com/infraboard/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v redis | grep -v broker | grep -v etcd)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep lint vet test test-coverage build clean

all: build

push: lint vet test build## git push
	@git push
	@rm -f build/*

dep: ## Get the dependencies
	@go mod download

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt

build: dep ## Build the binary file
	@go build -i -o build/$(PROJECT_NAME) $(MAIN_FILE)

clean: ## Remove previous build
	@rm -f build/*

codegen: # Init Service
	@protoc -I=.  -I${GOPATH}/src --go_out=module=github/infraboard/mcube:. pb/page/*.proto

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'