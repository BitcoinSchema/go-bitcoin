# Common makefile commands & variables between projects
include .make/Makefile.common

# Common Golang makefile commands & variables between projects
include .make/Makefile.go

## Not defined? Use default repo name which is the application
ifeq ($(REPO_NAME),)
	REPO_NAME="go-bitcoin"
endif

## Not defined? Use default repo owner
ifeq ($(REPO_OWNER),)
	REPO_OWNER="bitcoinschema"
endif

.PHONY: clean

all: ## Runs multiple commands
	@$(MAKE) test

clean: ## Remove previous builds and any test cache data
	@go clean -cache -testcache -i -r
	@test $(DISTRIBUTIONS_DIR)
	@if [ -d $(DISTRIBUTIONS_DIR) ]; then rm -r $(DISTRIBUTIONS_DIR); fi

lint:: ## Installs the linter and runs
	@if [ "$(shell command -v golangci-lint)" = "" ]; then curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $GOPATH/bin v1.23.8; fi
	@golangci-lint
	#@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $GOPATH/bin v1.31.0

release:: ## Runs common.release then runs godocs
	@$(MAKE) godocs

update-linter: ## Update the golangci-lint package
	@brew upgrade golangci-lint