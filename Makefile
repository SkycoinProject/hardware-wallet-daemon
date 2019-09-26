.DEFAULT_GOAL := help
.PHONY: run run-usb run-emulator test test-race
.PHONY: test-integration-emulator test-integration-wallet test-integration-emulator-enable-csrf test-integration-wallet-enable-csrf
.PHONY: check mocks lint
.PHONY: clean-coverage update-golden-files merge-coverage
.PHONY: install-linters format generate-client
.PHONY: release

run: ## Run hardware wallet daemon
	./run.sh ${ARGS}

run-usb: ## Run daemon in usb mode
	./run.sh -daemon-mode USB

run-emulator: ## Run daemon in emulator mode
	./run.sh -daemon-mode EMULATOR

run-help: ## Show daemon help
	./run.sh -help

test: ## Run tests for hardware wallet daemon
	@mkdir -p coverage/
	go test -coverpkg="github.com/skycoin/hardware-wallet-daemon/..." -coverprofile=coverage/go-test-cmd.coverage.out -timeout=5m ./src/...

test-race: ## Run tests for hardware wallet daemon with race flag
	@mkdir -p coverage/
	go test -race -coverpkg="github.com/skycoin/hardware-wallet-daemon/..." -coverprofile=coverage/go-test-cmd.coverage.out -timeout=5m ./src/...

test-integration-emulator: ## Run emulator integration tests
	./ci-scripts/integration-test.sh -a -m EMULATOR -n emulator-integration

test-integration-wallet: ## Run wallet integration tests
	./ci-scripts/integration-test.sh -m USB -n wallet-integration

test-integration-emulator-enable-csrf: ## Run wallet integration tests with CSRF enabled
	./ci-scripts/integration-test.sh -a -m EMULATOR -c -n emulator-integration-enable-csrf

test-integration-wallet-enable-csrf: ## Run emulator integration tests with CSRF enabled
	./ci-scripts/integration-test.sh -m USB -c -n wallet-integration-enable-csrf

check: test \
    test-integration-emulator \
    test-integration-wallet ## run unit and integration tests

mocks: ## Create all mock files for unit tests
	echo "Generating mock files"
	go generate ./src/...

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	@# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

clean-coverage: ## Remove coverage output files
	rm -rf ./coverage/

update-golden-files: ## Run integration tests in update mode
	./ci-scripts/integration-test.sh -u -m ${DEVICE_TYPE} -n update-golden-files

merge-coverage: ## Merge coverage files and create HTML coverage output. gocovmerge is required, install with `go get github.com/wadey/gocovmerge`
	@echo "To install gocovmerge do:"
	@echo "go get github.com/wadey/gocovmerge"
	gocovmerge coverage/*.coverage.out > coverage/all-coverage.merged.out
	go tool cover -html coverage/all-coverage.merged.out -o coverage/all-coverage.html
	@echo "Total coverage HTML file generated at coverage/all-coverage.html"
	@echo "Open coverage/all-coverage.html in your browser to view"

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.16.0

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/hardware-wallet-daemon ./cmd
	goimports -w -local github.com/skycoin/hardware-wallet-daemon ./src

generate-client: ## Generate go client using swagger
	swagger generate client swagger.yml --template-dir templates -t ./src

release: ## Build daemon binaries
	./ci-scripts/build-daemon.sh

clean-release: ## Remove release files
	rm -rf ./build/release
	rm -rf ./build/.xgo_output
	rm -rf ./build/.daemon_output

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
