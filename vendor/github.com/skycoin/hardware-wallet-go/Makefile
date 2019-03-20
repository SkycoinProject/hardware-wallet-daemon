.DEFAULT_GOAL := help
.PHONY: all build
.PHONY: test_unit test_integration test
.PHONY: dep vendor_proto proto mocks
.PHONY: clean lint check format

all: build

build: ## install cli
	cd cmd/cli && ./install.sh

dep: vendor_proto
	dep ensure
	# Ensure sources for protoc-gen-go and protobuf/proto are in sync
	dep ensure -add github.com/gogo/protobuf/protoc-gen-gofast ## setup dependencies

vendor_proto: proto
	mkdir -p vendor/github.com/google/protobuf
	cp -r -p src/device-wallet/messages/go/google/protobuf/descriptor.pb.go vendor/github.com/google/protobuf ## init proto messages package

mocks: ## Create all mock files for unit tests
	echo "Generating mock files"
	mockery -name Devicer -dir ./src/device-wallet -case underscore -inpkg -testonly
	mockery -name DeviceDriver -dir ./src/device-wallet -case underscore -inpkg -testonly

test_unit: ## Run unit tests
	go test -v github.com/skycoin/hardware-wallet-go/src/device-wallet

test_integration: ## Run integration tests
	go test -v github.com/skycoin/hardware-wallet-go/src/device-wallet/integration

test: test_unit test_integration ## Run all tests

proto: ## build proto files
	make -C src/device-wallet/messages build-go GO_IMPORT=github.com/skycoin/hardware-wallet-go/src/device-wallet/messages

clean: ## clean proto files
	make -C src/device-wallet/messages clean-go
	rm -rf vendor/github.com/google

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...

check: lint \
	test ## run checks

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/hardware-wallet-go ./cmd
	goimports -w -local github.com/skycoin/hardware-wallet-go ./src

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
