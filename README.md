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
