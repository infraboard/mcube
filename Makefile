MCUBE_MAIN := "cmd/mcube/main.go"
PROTOC_GEN_GO_HTTP_MAIN = "cmd/protoc-gen-go-http/main.go"
PROJECT_NAME := "mcube"
PKG := "github.com/infraboard/$(PROJECT_NAME)"
MOD_DIR := $(shell go env GOMODCACHE)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v redis | grep -v broker | grep -v etcd | grep -v examples)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep lint vet test test-coverage build clean

all: build

push: lint test build## git push
	@git push -u gitee
	@git push -u origin
	@rm -f build/*

dep: ## Get the dependencies
	@go mod download

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

install: ## install mcube cli
	@go install ${PKG}/cmd/mcube

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt

test-hg:  ## test http gen
	@protoc -I=. -I=${GOPATH}/src --go-http_out=. examples/http/hello.proto --go-http_opt=module="github.com/infraboard/mcube"

build: dep ## Build the binary file
	@go build -o build/$(PROJECT_NAME) $(MCUBE_MAIN)

clean: ## Remove previous build
	@rm -f build/*

gen: # Generate code
	@protoc -I=.. -I=/usr/local/include --go_out=. --go_opt=module=${PKG} ../mcube/pb/*/*.proto
	@protoc-go-inject-tag -input=pb/*/*.pb.go
	@mcube generate enum -p -m pb/*/*.pb.go
	
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'