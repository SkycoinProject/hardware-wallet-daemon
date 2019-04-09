.DEFAULT_GOAL := help
.PHONY: run lint format generate-client
.PHONY: test integration-test-emulator integration-test-wallet
.PHONY: clean-coverage update-golden-files merge-coverage
.PHONY: mocks

build: ## daemon build release
	@mkdir -p release
	go build -o ./release/skyd ./cmd/daemon

run: ## Run hardware wallet daemon
	./run.sh ${ARGS}

run-help: ## Show daemon help
	./run.sh -help

test: ## Run tests for hardware wallet daemon
	@mkdir -p coverage/
	go test -coverpkg="github.com/skycoin/hardware-wallet-daemon/..." -coverprofile=coverage/go-test-cmd.coverage.out -timeout=5m ./src/...

integration-test-emulator: ## Run emulator integration tests
	./ci-scripts/integration-test.sh -m emulator -n emulator-integration

integration-test-wallet: ## Run wallet integration tests
	./ci-scripts/integration-test.sh -m wallet -n wallet-integration

integration-test-wallet-enable-csrf: ## Run wallet integration tests with CSRF enabled
	./ci-scripts/integration-test.sh -m emulator -c -n emulator-integration-enable-csrf

integration-test-emulator-enable-csrf: ## Run emulator integration tests with CSRF enabled
	./ci-scripts/integration-test.sh -m wallet -c -n wallet-integration-enable-csrf

check: test \
    integration-test-emulator \
    integration-test-wallet ## run unit and integration tests

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
	./ci-scripts/integration-test.sh -u >/dev/null 2>&1 || true

merge-coverage: ## Merge coverage files and create HTML coverage output. gocovmerge is required, install with `go get github.com/wadey/gocovmerge`
	@echo "To install gocovmerge do:"
	@echo "go get github.com/wadey/gocovmerge"
	gocovmerge coverage/*.coverage.out > coverage/all-coverage.merged.out
	go tool cover -html coverage/all-coverage.merged.out -o coverage/all-coverage.html
	@echo "Total coverage HTML file generated at coverage/all-coverage.html"
	@echo "Open coverage/all-coverage.html in your browser to view"

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/hardware-wallet-daemon ./cmd
	goimports -w -local github.com/skycoin/hardware-wallet-daemon ./src

generate-client: ## Generate go client using swagger
	swagger generate client swagger.yml --template-dir templates -t ./src

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
