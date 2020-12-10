.DEFAULT_GOAL := help
.PHONY: run run-help test test-386 test-amd64 check check-newcoin
.PHONY: install-linters format release clean-release clean-coverage
.PHONY: install-deps-ui build-ui build-ui-travis help newcoin merge-coverage
.PHONY: generate update-golden-files
.PHONY: fuzz-base58 fuzz-encoder

# Platform specific checks
OSNAME = $(TRAVIS_OS_NAME)

# Tooling versions
GOLANGCI_LINT_VERSION ?= v1.32.0

install: ## Installs cxchain and cxchain-cli
	go install ./cmd/...

run-cxchain:  ## Run cxchain with default configuration. To add arguments, do 'make ARGS="--foo" run'.
	./run-cxchain.sh ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/cxchain/cxchain.go --help

test: clean-vendor ## Run tests for Skycoin
	@mkdir -p coverage/
	go test -coverpkg="github.com/skycoin/cx-chains/..." -coverprofile=coverage/go-test-cmd.coverage.out -timeout=5m ./cmd/...
	go test -coverpkg="github.com/skycoin/cx-chains/..." -coverprofile=coverage/go-test-src.coverage.out -timeout=5m ./src/...

test-386: clean-vendor ## Run tests for Skycoin with GOARCH=386
	GOARCH=386 go test ./cmd/... -timeout=5m
	GOARCH=386 go test ./src/... -timeout=5m

test-amd64: clean-vendor ## Run tests for Skycoin with GOARCH=amd64
	GOARCH=amd64 go test ./cmd/... -timeout=5m
	GOARCH=amd64 go test ./src/... -timeout=5m

lint: ## Run linters. Use make install-linters first.
	golangci-lint run -c .golangci.yml ./...
	@# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

check: lint clean-coverage test test-386 ## Run tests and linters

install-linters: ## Install linters
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOBIN) $(GOLANGCI_LINT_VERSION)

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/cx-chains ./cmd/
	goimports -w -local github.com/skycoin/cx-chains ./src/

clean-coverage: ## Remove coverage output files
	rm -rf ./coverage/

clean-vendor: ## Clean vendor directory.
	go mod tidy && go mod vendor

generate: ## Generate test interface mocks and struct encoders
	go generate ./src/...
	# mockery can't generate the UnspentPooler mock in package visor, patch it
	mv ./src/visor/blockdb/mock_unspent_pooler_test.go ./src/visor/mock_unspent_pooler_test.go
	sed -i "" -e 's/package blockdb/package visor/g' ./src/visor/mock_unspent_pooler_test.go
	sed -i "" -e 's/AddressHashes/blockdb.AddressHashes/g' ./src/visor/mock_unspent_pooler_test.go
	goimports -w -local github.com/skycoin/skycoin ./src/visor/mock_unspent_pooler_test.go

install-generators: ## Install tools used by go generate
	go get github.com/vektra/mockery/.../
	go get github.com/skycoin/skyencoder/cmd/skyencoder

merge-coverage: ## Merge coverage files and create HTML coverage output. gocovmerge is required, install with `go get github.com/wadey/gocovmerge`
	@echo "To install gocovmerge do:"
	@echo "go get github.com/wadey/gocovmerge"
	gocovmerge coverage/*.coverage.out > coverage/all-coverage.merged.out
	go tool cover -html coverage/all-coverage.merged.out -o coverage/all-coverage.html
	@echo "Total coverage HTML file generated at coverage/all-coverage.html"
	@echo "Open coverage/all-coverage.html in your browser to view"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
