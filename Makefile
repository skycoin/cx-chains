.DEFAULT_GOAL := help
.PHONY: run run-help test test-386 test-amd64 check check-newcoin
.PHONY: integration-tests-stable
.PHONY: integration-test-stable
.PHONY: integration-test-stable-disable-csrf
.PHONY: integration-test-stable-disable-wallet-api
.PHONY: integration-test-stable-disable-seed-api
.PHONY: integration-test-stable-enable-seed-api
.PHONY: integration-test-stable-enable-seed-api
.PHONY: integration-test-stable-disable-gui
.PHONY: integration-test-stable-db-no-unconfirmed
.PHONY: integration-test-stable-auth
.PHONY: integration-test-live integration-test-live-wallet
.PHONY: install-linters format release clean-release clean-coverage
.PHONY: install-deps-ui build-ui build-ui-travis help newcoin merge-coverage
.PHONY: generate update-golden-files
.PHONY: fuzz-base58 fuzz-encoder

COIN ?= skycoin

# Platform specific checks
OSNAME = $(TRAVIS_OS_NAME)

# Tooling versions
GOLANGCI_LINT_VERSION ?= v1.32.0

install: ## Installs cxchain and cxchain-cli
	go install ./cmd/...

run-client:  ## Run skycoin with desktop client configuration. To add arguments, do 'make ARGS="--foo" run'.
	./run-client.sh ${ARGS}

run-daemon:  ## Run skycoin with server daemon configuration. To add arguments, do 'make ARGS="--foo" run'.
	./run-daemon.sh ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/$(COIN)/$(COIN).go --help

run-integration-test-live: ## Run the skycoin node configured for live integration tests
	./ci-scripts/run-live-integration-test-node.sh

run-integration-test-live-cover: ## Run the skycoin node configured for live integration tests with coverage
	./ci-scripts/run-live-integration-test-node-cover.sh

test: clean-vendor ## Run tests for Skycoin
	@mkdir -p coverage/
	COIN=$(COIN) go test -coverpkg="github.com/$(COIN)/$(COIN)/..." -coverprofile=coverage/go-test-cmd.coverage.out -timeout=5m ./cmd/...
	COIN=$(COIN) go test -coverpkg="github.com/$(COIN)/$(COIN)/..." -coverprofile=coverage/go-test-src.coverage.out -timeout=5m ./src/...

test-386: clean-vendor ## Run tests for Skycoin with GOARCH=386
	GOARCH=386 COIN=$(COIN) go test ./cmd/... -timeout=5m
	GOARCH=386 COIN=$(COIN) go test ./src/... -timeout=5m

test-amd64: clean-vendor ## Run tests for Skycoin with GOARCH=amd64
	GOARCH=amd64 COIN=$(COIN) go test ./cmd/... -timeout=5m
	GOARCH=amd64 COIN=$(COIN) go test ./src/... -timeout=5m

lint: ## Run linters. Use make install-linters first.
	golangci-lint run -c .golangci.yml ./...
	@# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

check: lint clean-coverage test test-386 integration-tests-stable check-newcoin ## Run tests and linters

integration-tests-stable: integration-test-stable \
	integration-test-stable-disable-csrf \
	integration-test-stable-disable-wallet-api \
	integration-test-stable-disable-seed-api \
	integration-test-stable-enable-seed-api \
	integration-test-stable-disable-gui \
	integration-test-stable-auth \
	integration-test-stable-db-no-unconfirmed ## Run all stable integration tests

integration-test-stable: ## Run stable integration tests
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -c -x -n enable-csrf-header-check

integration-test-stable-disable-header-check: ## Run stable integration tests with header check disabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -n disable-header-check

integration-test-stable-disable-csrf: ## Run stable integration tests with CSRF disabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -n disable-csrf

integration-test-stable-disable-wallet-api: ## Run disable wallet api integration tests
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-disable-wallet-api.sh

integration-test-stable-enable-seed-api: ## Run enable seed api integration test
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-enable-seed-api.sh

integration-test-stable-disable-gui: ## Run tests with the GUI disabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-disable-gui.sh

integration-test-stable-db-no-unconfirmed: ## Run stable tests against the stable database that has no unconfirmed transactions
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-stable.sh -d -n no-unconfirmed

integration-test-stable-auth: ## Run stable tests with HTTP Basic auth enabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-auth.sh

integration-test-live: ## Run live integration tests
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh -c

integration-test-live-wallet: ## Run live integration tests with wallet
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh -w

integration-test-live-enable-header-check: ## Run live integration tests against a node with header check enabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh

integration-test-live-disable-csrf: ## Run live integration tests against a node with CSRF disabled
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh

integration-test-live-disable-networking: ## Run live integration tests against a node with networking disabled (requires wallet)
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-live.sh -c -k

install-linters: ## Install linters
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOBIN) $(GOLANGCI_LINT_VERSION)

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/cx-chains ./cmd/
	goimports -w -local github.com/skycoin/cx-chains ./src/

clean-coverage: ## Remove coverage output files
	rm -rf ./coverage/

clean-vendor: ## Clean vendor directory.
	go mod tidy && go mod vendor

newcoin: ## Rebuild cmd/$COIN/$COIN.go file from the template. Call like "make newcoin COIN=foo".
	go run cmd/newcoin/newcoin.go createcoin --coin $(COIN)

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

update-golden-files: ## Run integration tests in update mode
	./ci-scripts/integration-test-stable.sh -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -c -x -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -d -u >/dev/null 2>&1 || true
	./ci-scripts/integration-test-stable.sh -c -x -d -u >/dev/null 2>&1 || true

merge-coverage: ## Merge coverage files and create HTML coverage output. gocovmerge is required, install with `go get github.com/wadey/gocovmerge`
	@echo "To install gocovmerge do:"
	@echo "go get github.com/wadey/gocovmerge"
	gocovmerge coverage/*.coverage.out > coverage/all-coverage.merged.out
	go tool cover -html coverage/all-coverage.merged.out -o coverage/all-coverage.html
	@echo "Total coverage HTML file generated at coverage/all-coverage.html"
	@echo "Open coverage/all-coverage.html in your browser to view"

fuzz-base58: ## Fuzz the base58 package. Requires https://github.com/dvyukov/go-fuzz
	go-fuzz-build github.com/skycoin/skycoin/src/cipher/base58/internal
	go-fuzz -bin=base58fuzz-fuzz.zip -workdir=src/cipher/base58/internal

fuzz-encoder: ## Fuzz the encoder package. Requires https://github.com/dvyukov/go-fuzz
	go-fuzz-build github.com/skycoin/skycoin/src/cipher/encoder/internal
	go-fuzz -bin=encoderfuzz-fuzz.zip -workdir=src/cipher/encoder/internal

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
