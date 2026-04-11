.PHONY: run-guest run-business run-auth migrate-up run-all build tidy s3-up s3-ui install-tools docs docs-guest docs-business generate-guest-business-client test test-cover
.PHONY: goose-up goose-down goose-status goose-create

COUNT ?= 1

run-guest: docs-guest
	go run ./cmd/guest-api

run-business: docs-business
	go run ./cmd/business-api

run-auth: docs-admin-auth
	go run ./cmd/admin-auth-api

migrate-up:
	go run ./cmd/migrator

run-all:
	$(MAKE) -j 3 run-guest run-business run-auth

build: docs
	go build -o bin/migrator ./cmd/migrator
	go build -o bin/guest-api ./cmd/guest-api
	go build -o bin/business-api ./cmd/business-api
	go build -o bin/admin-auth-api ./cmd/admin-auth-api

tidy:
	go mod tidy

test:
	go test -v ./... -count=$(COUNT)

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

s3-up:
	docker compose -f docker/compose.yaml up -d s3
	bash scripts/bootstrap.sh

s3-ui:
	docker compose -f docker/compose.yaml up -d garage_webui
	@echo "web_ui: http://localhost:3909"

-include .env
DB_DSN="host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password='$(POSTGRES_PASSWORD)' dbname=$(POSTGRES_DB) sslmode=$(POSTGRES_SSL)"
MIGRATIONS_DIR=migrations

goose-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) up

goose-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) down

goose-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) status

goose-create:
	@if [ -z "$(name)" ]; then \
		echo "Помилка: вкажіть назву міграції. Приклад: make goose-create name=add_users_table"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

docs: docs-guest docs-business docs-admin-auth

docs-guest:
	@echo "generating swagger for guest service api..."
	go tool swag init -g main.go -d cmd/guest-api,internal/guest -o docs/api/guest --parseInternal --parseDependency

docs-business:
	@echo "generating swagger for business service api..."
	go tool swag init -g main.go -d cmd/business-api,internal/business -o docs/api/business --parseInternal --parseDependency

docs-admin-auth:
	@echo "generating swagger for admin-auth service api..."
	go tool swag init -g main.go -d cmd/admin-auth-api,internal/admin-auth -o docs/api/admin-auth --parseInternal --parseDependency

generate-guest-business-client: docs-business
	@echo "generating business client for guest service..."
	mkdir -p internal/guest/gateway/business/client
	go tool swagger generate client \
		-f docs/api/business/swagger.yaml \
		-t internal/guest/gateway/business/client \
		-c business_client \
		-m dto

