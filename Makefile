GO ?= go

.PHONY: help setup migrate run test vet build fmt

help:
	@echo "Available commands:"
	@echo "  make setup    Run migrations against local PostgreSQL"
	@echo "  make migrate  Run database migrations"
	@echo "  make run      Run the API server"
	@echo "  make test     Run tests"
	@echo "  make vet      Run go vet"
	@echo "  make build    Build the API binary"
	@echo "  make fmt      Format Go files"

setup: migrate

migrate:
	$(GO) run ./cmd/migrate

run:
	$(GO) run ./cmd/api

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

build:
	$(GO) build -o bin/gotask ./cmd/api

fmt:
	$(GO)fmt -w $$(find . -name '*.go' -not -path './vendor/*')
