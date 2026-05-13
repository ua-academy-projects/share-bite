# Share Bite

## How to Run Locally

### 1. Configuration

Create a local `.env` file based on the provided example:

```bash
cp .env.example .env
```

> **Note:** Please use distinct, strong credentials per service in non-local environments (e.g. for `POSTGRES_PASSWORD` and `REDIS_PASSWORD`).

### 2. Database Infrastructure

Start the local PostgreSQL database using Docker Compose:

```bash
docker compose -f build/compose.yaml up -d
```

### Optional pgAdmin:

```bash
docker compose -f build/compose.yaml --profile tools up -d
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

### 5. Object Storage (S3 / Garage)

The application uses S3-compatible storage for storing media files. It supports both local development via Garage and real AWS S3 for production.

#### Local Development (Garage)

Start the local Garage instance:
```bash
make s3-up
```

After bootstrap, copy the credentials printed in the terminal into your `.env`:
```env
S3_ENDPOINT=http://localhost:4300
S3_REGION=garage
S3_ACCESS_KEY=<printed by bootstrap>
S3_SECRET_KEY=<printed by bootstrap>
S3_BUCKET=app-dev-bucket
S3_USE_PATH_STYLE=true
```

> Bootstrap runs once. On subsequent starts credentials stay the same.

**Web UI (optional)**
```bash
make s3-ui
```

> Open http://localhost:4309

#### Production (AWS S3)

For real AWS deployment, configure your production `.env` like this:

```env
S3_REGION=eu-central-1
S3_BUCKET=your-production-bucket-name
```

**Recommended: IAM Roles (no keys needed)**
When running on AWS infrastructure (EC2, ECS, Lambda), leave `S3_ACCESS_KEY` and `S3_SECRET_KEY` empty.
The client will automatically use the instance's IAM role.

**Alternative: Static credentials**
```env
S3_ACCESS_KEY=your-aws-access-key
S3_SECRET_KEY=your-aws-secret-key
```

> `S3_ENDPOINT` and `S3_USE_PATH_STYLE` are not required for native AWS S3.

### 6. Redis

For local development, add Redis connection values to `.env` based on `.env.example`:

```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis_example_password
```

Start the Redis service using the provided docker-compose configuration:
```bash
docker compose -f build/compose.infra.yaml up -d
```

### 7. Notifications service

The `notifications-service` consumes UI notifications from the SQS queue, enriches them with actor metadata, stores them in `notifications_history`, and forwards them to SSE clients.

Add the notifications-specific env vars to `.env`:

```env
NOTIFICATION_HTTP_SERVER_HOST=0.0.0.0
NOTIFICATION_HTTP_SERVER_PORT=4005
NOTIFICATION_SQS_QUEUE_URL=<terraform output notifications_sse_queue_url>
```

Run it locally:

```bash
go run cmd/notifications-service/main.go
```

For local SQS testing, point the service to a local emulator or your AWS queue with:

```env
NOTIFICATION_AWS_REGION=us-east-2
NOTIFICATION_SQS_ENDPOINT_URL=http://localhost:4566
```

If you are using localstack, create the queue in the `notifications-service` consumer path and keep the SSE/UI queue separated as in Terraform.

If you deploy infrastructure with Terraform, export the queue URL from the outputs and point the service at the SSE queue.

### 8. Code Generation & API Clients

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

### 9. Testing

> Ensure API clients are generated (step 8) before running tests.

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
