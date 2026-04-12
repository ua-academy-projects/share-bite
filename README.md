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

_(To run them individually, you can use `make run-guest`, `make run-business`, or `make run-auth`)._

### 5. Run local S3 storage

### Start

```bash
make s3-up
```

After bootstrap, copy credentials printed in the terminal into `.env`:

```env
S3_ENDPOINT=http://localhost:3900
S3_REGION=garage
S3_ACCESS_KEY=<printed by bootstrap>
S3_SECRET_KEY=<printed by bootstrap>
S3_BUCKET=app-dev-bucket
S3_USE_PATH_STYLE=true
```

> Bootstrap runs once. On subsequent starts credentials stay the same.

### Web UI (optional)

```bash
make s3-ui
```

Open http://localhost:3909

### 6. Code Generation & API Clients

To prevent excessive git churn, **generated Swagger documentation and API clients are not committed to the repository**. You must generate them locally before building or running tests.

No external binaries (`swag` or `go-swagger`) need to be installed — both are managed as Go tool dependencies.

Run the full generation suite:

```bash
make generate
```

This orchestrates the correct execution order:

1. **Generates Swagger specs** for all microservices into `docs/api/`.
2. **Generates type-safe Go clients** required for inter-service communication.

> Run `make clean` to wipe all generated files and start fresh.

---

### 7. Testing

> Ensure API clients are generated (step 6) before running tests.

Run the full test suite:

```bash
make test
```

By default tests run once (`COUNT=1`). To catch flaky behaviour, increase the count:

```bash
make test COUNT=5
```

Generate a coverage report and open it in your browser:

```bash
make test-cover
```