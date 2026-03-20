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
