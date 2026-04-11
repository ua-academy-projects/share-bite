run-guest:
	go run cmd/guest-api/main.go

run-business:
	go run cmd/business-api/main.go

run-auth:
	go run cmd/admin-auth-api/main.go

migrate-up:
	go run cmd/migrator/main.go

run-all:
	make -j 3 run-guest run-business run-auth

build:
	go build -o bin/migrator cmd/migrator/main.go
	go build -o bin/guest-api cmd/guest-api/main.go
	go build -o bin/business-api cmd/business-api/main.go
	go build -o bin/admin-auth-api cmd/admin-auth-api/main.go

tidy:
	go mod tidy

s3-up:
	docker compose -f docker/compose.yaml up -d s3
	bash scripts/bootstrap.sh

s3-ui:
	docker compose -f docker/compose.yaml up -d garage_webui
	@echo "web_ui: http://localhost:3909"

-include .env
DB_DSN="host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password='$(POSTGRES_PASSWORD)' dbname=$(POSTGRES_DB) sslmode=$(POSTGRES_SSL)"
MIGRATIONS_DIR=migrations
.PHONY: goose-up goose-down goose-status goose-create

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