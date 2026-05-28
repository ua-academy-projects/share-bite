# ShareBite Cleanup Worker

A scheduled job system for cleaning up stale data in the ShareBite database.

## Quick Start

### Local CLI

```bash
# Build
go build -o cleanup-worker ./cmd/cleanup-worker

# Run (dry-run by default from env or use -dry-run flag)
./cleanup-worker

# Run with dry-run flag
./cleanup-worker -dry-run

# Run with custom config
CLEANUP_RETENTION_DAYS=60 CLEANUP_BATCH_SIZE=200 ./cleanup-worker
```

### Lambda

```bash
# Build
GOOS=linux GOARCH=amd64 go build -o main ./lambda/cleanup-worker
zip function.zip main

# Deploy
aws lambda create-function \
  --function-name share-bite-cleanup \
  --runtime go1.x \
  --handler main \
  --zip-file fileb://function.zip \
  --timeout 300 \
  --environment Variables='{"CLEANUP_RETENTION_DAYS=30","CLEANUP_BATCH_SIZE=500",...}'
```

## Configuration

```bash
CLEANUP_RETENTION_DAYS=30          # Records older than this are cleaned
CLEANUP_BATCH_SIZE=100             # Records deleted per batch
CLEANUP_DRY_RUN=true               # Log what would be deleted (no actual deletes)
CLEANUP_SCHEDULE_ENABLED=true      # Enable scheduling (used by Lambda)
```

## What Gets Cleaned

1. **Old guest posts** - Posts older than retention period (but not archived)
2. **Expired password reset tokens** - Tokens that never got used and have expired

## Features

- ✅ Dry-run mode (test without deleting)
- ✅ Batch processing (safe, configurable batch sizes)
- ✅ Idempotent (safe to retry)
- ✅ CLI support (local testing)
- ✅ Lambda ready (AWS integration)
- ✅ Comprehensive logging
- ✅ Cascade deletes handled properly

## Testing

```bash
# Unit tests
go test -v ./internal/cleanup/...

# Dry-run test
CLEANUP_DRY_RUN=true ./cleanup-worker

# See results without modifying database
CLEANUP_DRY_RUN=true CLEANUP_RETENTION_DAYS=1 ./cleanup-worker
```

## Documentation

- [Cleanup Worker Guide](../cleanup-worker.md) - Full documentation
- [EventBridge Schedule Guide](../eventbridge-schedule-guide.md) - Schedule expressions

## Example Output

```
Starting cleanup with retention period: 720h0m0s, batch size: 100, dry-run: false

=== Cleanup Results ===

Job: expire-old-posts
  Records Found:   250
  Records Deleted: 250
  Duration:        2.345s

Job: cleanup-expired-tokens
  Records Found:   75
  Records Deleted: 75
  Duration:        0.567s

=== Summary ===
Total Records Found:   325
Total Records Deleted: 325
```

## Architecture

```
cmd/cleanup-worker/main.go      → CLI entry point
lambda/cleanup-worker/main.go   → Lambda handler
internal/cleanup/
  ├── entity/                   → Data types
  ├── repository/               → Database layer
  └── service/                  → Business logic
internal/config/env/cleanup.go  → Configuration
```

## See Also

- [Database Schema Documentation](../architecture/database.md)
- [AWS Lambda Setup](../aws/lambda.md)
- [Makefile targets](#makefile)

## Status

- [x] Cleanup logic implemented
- [x] Job runs locally from CLI
- [x] Lambda handler invokes business logic
- [x] EventBridge schedule documented
- [x] Dry-run mode logs what would be deleted
- [x] Job is idempotent and safe to retry
- [x] Unit tests included
- [x] Documentation complete
