# Variables
PACKAGES := $(shell go list ./...)

.PHONY: all
all: help

## help: Display this help message
.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo

## test: Run tests with race detection and coverage
.PHONY: test
test:
	go test -race -cover $(PACKAGES)

## lint: Run golangci-lint code quality checks
.PHONY: lint
lint:
	golangci-lint run ./...

## lint-fix: Run golangci-lint with auto-fix for common issues
.PHONY: lint-fix
lint-fix:
	golangci-lint fmt
	golangci-lint run --fix ./...

## clean: Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/