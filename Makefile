VERSION=$(shell date -u '+%y.%m.%d-%H.%M')

GO := go
GO111MODULE := on
CGO_ENABLED := 0
GOLANGCI_LINT := golangci-lint
BIN := bin/redins

.PHONY: all
all: fmt lint test build

.PHONY: build
build:
	$(GO) build -o $(BIN) ./

.PHONY: install
install:
	$(GO) install ./

.PHONY: fmt
fmt:
	@echo "Running gofmt..."
	@gofmt -s -w .
	@echo "Running goimports..."
	@goimports -w .

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run --timeout 4m --config .golangci.yaml

.PHONY: test
test:
	$(GO) test ./... -race

.PHONY: clean
clean:
	rm -f $(BIN)
	rm -f coverage.out
	rm -rf .bin/
	rm -rf bin/

.PHONY: image
image:
	docker build -t ealebed/redins:${VERSION} .
	docker push ealebed/redins:${VERSION}

.PHONY: update
update:
	$(GO) get -u -v ./
	$(GO) mod verify
	$(GO) mod tidy

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all      - build, format, lint, and test"
	@echo "  build    - build the binary"
	@echo "  install  - install the binary"
	@echo "  fmt      - format code with gofmt and goimports"
	@echo "  lint     - run golangci-lint"
	@echo "  test     - run tests with race detector"
	@echo "  clean    - clean build artifacts"
	@echo "  image    - build and push Docker image"
	@echo "  update   - update dependencies"
	@echo "  help     - show this help message"
