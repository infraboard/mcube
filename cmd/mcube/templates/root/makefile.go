package root

// MakefileTemplate todo
const MakefileTemplate = `BINARY_NAME := "{{.Name}}"
MAIN_FILE_PAHT := "main.go"
PROJECT_NAME := "{{.Name}}"
PKG := "{{.PKG}}"
IMAGE_PREFIX := "{{.PKG}}"

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v redis)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep lint vet test test-coverage build clean

all: build

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
	@sh ./script/build.sh local ${BINARY_NAME} ${MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}

linux: ## Linux build
	@sh ./script/build.sh linux ${BINARY_NAME} ${MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	
run: dep ## Run Server
	@go build -o ${BINARY_NAME} ${MAIN_FILE_PAHT} ${IMAGE_PREFIX} ${PKG}
	@./${BINARY_NAME} service start

clean: ## Remove previous build
	@go clean .
	@rm -f ${BINARY_NAME}

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'`
