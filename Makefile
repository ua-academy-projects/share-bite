.PHONY: build build-guest build-business build-auth build-migrator build-notifications build-outbox build-lambda
.PHONY: run-all run-guest run-business run-auth migrate-up
.PHONY: test test-cover tidy clean
.PHONY: docs docs-guest docs-business docs-admin-auth
.PHONY: generate generate-guest-business-client
.PHONY: s3-up s3-ui docker-build
.PHONY: goose-up goose-down goose-status goose-create
.PHONY: docker-build docker-build-business-operator docker-push-business-operator
.PHONY: k8s-secrets k8s-up k8s-down k8s-migrate

REGISTRY ?= mykolashevchenko
TAG ?= latest
OPERATOR_IMAGE := $(REGISTRY)/business-operator:$(TAG)

COUNT ?= 1
MIGRATIONS_DIR := migrations
K8S_NAMESPACE ?= share-bite-local
K8S_SECRETS_FILE ?= docs/k8s/secrets.local.yaml
K8S_SECRET_NAME ?= share-bite-secrets
K8S_READY_TIMEOUT ?= 180s

VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "development")
COMMIT     ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
PKG        := github.com/ua-academy-projects/share-bite/pkg/version
ifndef BUILD_TIME
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
endif

LDFLAGS := -ldflags "\
  -X '$(PKG).Version=$(VERSION)' \
  -X '$(PKG).CommitHash=$(COMMIT)' \
  -X '$(PKG).BuildTime=$(BUILD_TIME)'"

-include .env
DB_DSN := host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password='$(POSTGRES_PASSWORD)' dbname=$(POSTGRES_DB) sslmode=$(POSTGRES_SSL)

run-guest: build-guest
	./bin/guest-api

run-business: build-business
	./bin/business-api

run-auth: build-auth
	./bin/admin-auth-api

migrate-up: build-migrator
	./bin/migrator

run-all:
	$(MAKE) -j 3 run-guest run-business run-auth

build-guest:
	go build $(LDFLAGS) -o bin/guest-api ./cmd/guest-api

build-business:
	go build $(LDFLAGS) -o bin/business-api ./cmd/business-api

build-auth:
	go build $(LDFLAGS) -o bin/admin-auth-api ./cmd/admin-auth-api

build-migrator:
	go build $(LDFLAGS) -o bin/migrator ./cmd/migrator

build-notifications:
	go build $(LDFLAGS) -o bin/notifications-service ./cmd/notifications-service

build-outbox:
	go build $(LDFLAGS) -o bin/outbox-worker ./cmd/outbox-worker

build-lambda:
	go build $(LDFLAGS) -o bin/notifications-lambda ./cmd/notifications-lambda

build: build-guest build-business build-auth build-migrator build-notifications build-outbox build-lambda

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

docker-build-business-operator:
	docker build -t $(OPERATOR_IMAGE) -f build/Dockerfile.business-operator .

docker-push-business-operator:
	docker push $(OPERATOR_IMAGE)

run-business-operator:
	kubectl apply -f deploy/k8s/operators/business/crd.yaml
	kubectl apply -f deploy/k8s/operators/business/rbac.yaml
	go run ./cmd/business-operator

deploy-operator:
	kubectl apply -f deploy/k8s/operators/business/operator-deployment.yaml
	docker build \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg COMMIT=$(COMMIT) \
	  --build-arg BUILD_TIME=$(BUILD_TIME) \
	  -t guest-api -f build/Dockerfile.guest .
	docker build \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg COMMIT=$(COMMIT) \
	  --build-arg BUILD_TIME=$(BUILD_TIME) \
	  -t business-api -f build/Dockerfile.business .
	docker build \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg COMMIT=$(COMMIT) \
	  --build-arg BUILD_TIME=$(BUILD_TIME) \
	  -t admin-auth-api -f build/Dockerfile.admin .
	docker build \
	  --build-arg VERSION=$(VERSION) \
	  --build-arg COMMIT=$(COMMIT) \
	  --build-arg BUILD_TIME=$(BUILD_TIME) \
	  -t migrator -f build/Dockerfile.migrator .

k8s-secrets:
	kubectl apply -f deploy/k8s/infra/namespace.yaml
	kubectl apply -f $(K8S_SECRETS_FILE)

k8s-up: k8s-secrets
	kubectl apply -k deploy/k8s/infra
	kubectl wait --for=create secret/$(K8S_SECRET_NAME) -n $(K8S_NAMESPACE) --timeout=$(K8S_READY_TIMEOUT)
	kubectl rollout status statefulset/postgres -n $(K8S_NAMESPACE) --timeout=$(K8S_READY_TIMEOUT)
	kubectl rollout status deployment/redis -n $(K8S_NAMESPACE) --timeout=$(K8S_READY_TIMEOUT)
	@echo "Infrastructure ready in namespace $(K8S_NAMESPACE)."
	@echo "Next step: run 'make k8s-migrate'."

# Job pod templates are immutable; set image in deploy/k8s/infra/migrator-job.yaml before apply.
k8s-migrate:
	kubectl delete job share-bite-migrator -n $(K8S_NAMESPACE) --ignore-not-found=true
	kubectl apply -f deploy/k8s/infra/migrator-job.yaml
	kubectl wait --for=condition=complete --timeout=180s job/share-bite-migrator -n $(K8S_NAMESPACE)

k8s-down:
	kubectl delete namespace $(K8S_NAMESPACE) --ignore-not-found=true
