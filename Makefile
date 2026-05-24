.PHONY: run-guest run-business run-auth migrate-up run-all build tidy s3-up s3-ui
.PHONY: test test-cover docs docs-guest docs-business docs-admin-auth
.PHONY: generate generate-guest-business-client clean
.PHONY: goose-up goose-down goose-status goose-create
.PHONY: docker-build

COUNT ?= 1
MIGRATIONS_DIR := migrations

-include .env
DB_DSN := host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password='$(POSTGRES_PASSWORD)' dbname=$(POSTGRES_DB) sslmode=$(POSTGRES_SSL)

run-guest:
	go run ./cmd/guest-api

run-business:
	go run ./cmd/business-api

run-auth:
	go run ./cmd/admin-auth-api

migrate-up:
	go run ./cmd/migrator

run-all:
	$(MAKE) -j 3 run-guest run-business run-auth

build:
	go build -o bin/migrator ./cmd/migrator
	go build -o bin/guest-api ./cmd/guest-api
	go build -o bin/business-api ./cmd/business-api
	go build -o bin/admin-auth-api ./cmd/admin-auth-api
	go build -o bin/notifications-service ./cmd/notifications-service
	go build -o bin/outbox-worker ./cmd/outbox-worker
	go build -o bin/notifications-lambda ./cmd/notifications-lambda

tidy:
	go mod tidy

test:
	go test -v ./... -count=$(COUNT)

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

s3-up:
	docker compose -f build/compose.yaml up -d s3
	bash scripts/bootstrap.sh

s3-ui:
	docker compose -f build/compose.yaml up -d garage_webui
	@echo "web_ui: http://localhost:4309"

goose-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DB_DSN)" goose -dir $(MIGRATIONS_DIR) up

goose-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DB_DSN)" goose -dir $(MIGRATIONS_DIR) down

goose-status:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DB_DSN)" goose -dir $(MIGRATIONS_DIR) status

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
	rm -rf internal/guest/gateway/business/client
	mkdir -p internal/guest/gateway/business/client
	go tool swagger generate client \
		-f docs/api/business/swagger.yaml \
		-t internal/guest/gateway/business/client \
		-c business_client \
		-m dto

generate: docs generate-guest-business-client

clean:
	@echo "cleaning generated files..."
	rm -rf docs/api
	rm -rf internal/guest/gateway/business/client

docker-build:
	docker build -t guest-api -f build/Dockerfile.guest .
	docker build -t business-api -f build/Dockerfile.business .
	docker build -t admin-auth-api -f build/Dockerfile.admin .
	docker build -t migrator -f build/Dockerfile.migrator .
