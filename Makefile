GO ?= go
AIR_BIN ?= $(shell command -v air 2>/dev/null || printf '%s/bin/air' "$$($(GO) env GOPATH)")

.PHONY: help setup migrate run dev start test vet build fmt

help:
	@echo "Available commands:"
	@echo "  make setup    Run migrations against local PostgreSQL"
	@echo "  make migrate  Run database migrations"
	@echo "  make run      Run the API server with hot reload"
	@echo "  make start    Run the API server once without hot reload"
	@echo "  make test     Run tests"
	@echo "  make vet      Run go vet"
	@echo "  make build    Build the API binary"
	@echo "  make fmt      Format Go files"

setup: migrate

migrate:
	$(GO) run ./cmd/migrate

run:
	$(MAKE) dev

dev:
	@test -x "$(AIR_BIN)" || { \
		echo "air is required for hot reload. Install it with: go install github.com/air-verse/air@v1.61.0"; \
		exit 1; \
	}
	$(AIR_BIN) -c .air.toml

start:
	$(GO) run ./cmd/api

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

build:
	$(GO) build -o bin/gotask ./cmd/api

fmt:
	$(GO)fmt -w $$(find . -name '*.go' -not -path './vendor/*')
