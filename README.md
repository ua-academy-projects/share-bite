# Share Bite

## How to Run Locally

### 1. Configuration
Create a local `.env` file based on the provided example:
```bash
cp .env.example .env
```

### 2. Database Infrastructure
Start the local PostgreSQL database using Docker Compose:
```bash
docker compose -f docker/compose.yaml up -d
```

### Optional pgAdmin:
```bash
docker compose -f docker/compose.yaml --profile tools up -d
```

### 3. Migrations
Apply the latest database schema using the built-in migrator:
```bash
go run cmd/migrator/main.go
```

### 4. Start the Microservices
Use the `Makefile` to start **all three services concurrently** (Guest, Business, Admin-Auth):
```bash
make run-all
```

*(To run them individually, you can use `make run-guest`, `make run-business`, or `make run-auth`).*