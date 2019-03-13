.DEFAULT_GOAL := help

run: ## Run hardware wallet daemon
	go run cmd/daemon/daemon.go

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	@# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/hardware-wallet-daemon ./cmd
	goimports -w -local github.com/skycoin/hardware-wallet-daemon ./src

generate-client: ## Generate client using swagger
	swagger generate client swagger.yml -t ./src

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
