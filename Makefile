GO ?= go
AIR_BIN ?= $(shell command -v air 2>/dev/null || printf '%s/bin/air' "$$($(GO) env GOPATH)")

.PHONY: help setup migration migrate migrate-down migrate-reset migrate-version run dev start test vet build fmt swagger

help:
	@echo "Available commands:"
	@echo "  make setup    Run migrations against local PostgreSQL"
	@echo "  make migration  Interactively create a date-prefixed migration pair"
	@echo "  make migrate  Run database migrations"
	@echo "  make migrate-down  Roll back the most recent database migration"
	@echo "  make migrate-reset  Clear all public tables and migration records"
	@echo "  make migrate-version  Show the current database migration version"
	@echo "  make run      Run the API server with hot reload"
	@echo "  make start    Run the API server once without hot reload"
	@echo "  make test     Run tests"
	@echo "  make vet      Run go vet"
	@echo "  make build    Build the API binary"
	@echo "  make fmt      Format Go files"
	@echo "  make swagger  Regenerate Swagger documentation"

setup: migrate

migration:
	$(GO) run ./cmd/migration

migrate:
	$(GO) run -mod=readonly ./cmd/migrate

migrate-down:
	$(GO) run -mod=readonly ./cmd/migrate -direction down

migrate-reset:
	$(GO) run -mod=readonly ./cmd/migrate -direction reset

migrate-version:
	$(GO) run -mod=readonly ./cmd/migrate -direction version

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

swagger:
	$(GO) run github.com/swaggo/swag/cmd/swag@v1.16.6 init --generalInfo cmd/api/main.go --parseInternal --output docs
