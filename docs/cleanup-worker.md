# Cleanup Worker Documentation

## Overview

The ShareBite Cleanup Worker is a scheduled job system that performs maintenance tasks on the database and removes stale data. It can run:

- **Locally** via CLI (`cmd/cleanup-worker`)
- **On AWS Lambda** (`lambda/cleanup-worker`)
- **On a schedule** via AWS EventBridge

## Features

- **Dry-run mode**: Test what would be deleted without making changes
- **Batch processing**: Delete records in configurable batch sizes for safety
- **Idempotent**: Safe to retry without side effects
- **Comprehensive logging**: Detailed output of cleanup operations
- **Multiple cleanup strategies**: Extensible design for adding more cleanup tasks

## Cleanup Operations

### 1. Expire Old Posts

Removes guest posts that are older than the retention period (default: 30 days).

- **Query**: Finds posts with `created_at < (now - retention_period)`
- **Filter**: Skips archived posts (preserves audit trail)
- **Batch processing**: Deletes in batches to avoid locking tables

**Cascade effects**: 
- Deletes associated `guest.post_images` (via FK cascade)
- Deletes associated `guest.comments` (via FK cascade)
- Deletes associated `guest.post_likes` (via FK cascade)

### 2. Clean Expired Password Reset Tokens

Removes password reset tokens that have expired and were never used.

- **Query**: Finds tokens with `expires_at < (now - retention_period)` AND `used_at IS NULL`
- **Batch processing**: Deletes in batches for safety

**Why this matters**: Prevents accumulation of expired tokens that could impact query performance.

## Configuration

All configuration is via environment variables:

```bash
# Number of days records must be old to be cleaned (default: 30)
CLEANUP_RETENTION_DAYS=30

# Number of records to delete per batch (default: 100)
# Larger batch sizes are faster but lock tables longer
CLEANUP_BATCH_SIZE=100

# Enable dry-run mode (don't actually delete, just log)
CLEANUP_DRY_RUN=true

# Enable schedule (used by Lambda/EventBridge)
CLEANUP_SCHEDULE_ENABLED=true
```

### Recommended Configuration

**Development**:
```bash
CLEANUP_RETENTION_DAYS=7
CLEANUP_BATCH_SIZE=50
CLEANUP_DRY_RUN=true
```

**Production**:
```bash
CLEANUP_RETENTION_DAYS=30
CLEANUP_BATCH_SIZE=500
CLEANUP_DRY_RUN=false
CLEANUP_SCHEDULE_ENABLED=true
```

## Usage

### CLI - Local Execution

#### Basic Usage

```bash
# Run with default config (dry-run off, 30 day retention, batch size 100)
./cmd/cleanup-worker/main.go

# Run in dry-run mode
./cmd/cleanup-worker/main.go -dry-run

# Run with custom config via environment
CLEANUP_RETENTION_DAYS=60 CLEANUP_BATCH_SIZE=200 ./cmd/cleanup-worker/main.go
```

#### Build and Run

```bash
# Build the CLI
go build -o cleanup-worker ./cmd/cleanup-worker/main.go

# Run
./cleanup-worker

# Run with dry-run
./cleanup-worker -dry-run

# Run with verbose output
CLEANUP_DRY_RUN=true ./cleanup-worker
```

#### Output Example

```
Starting cleanup with retention period: 720h0m0s, batch size: 100, dry-run: true
Lambda: Successfully connected to database
Starting cleanup of posts older than 2026-04-10 14:30:00
Found 250 posts older than 2026-04-10 14:30:00
[DRY-RUN] Would delete 250 posts

=== Cleanup Results ===

Job: expire-old-posts
  Records Found:  250
  Would Delete:   250 (DRY-RUN)
  Duration:       1.234s
  
Job: cleanup-expired-tokens
  Records Found:  75
  Would Delete:   75 (DRY-RUN)
  Duration:       0.456s

=== Summary ===
Total Records Found:     325
Would Delete (DRY-RUN):  325
Mode: DRY-RUN (no changes made)
```

### Lambda - AWS Function

#### Build the Lambda

```bash
# Install AWS Lambda Go runtime
go get github.com/aws/aws-lambda-go/cmd/build-lambda-zip

# Build Lambda function
GOOS=linux GOARCH=amd64 go build -o main ./lambda/cleanup-worker
zip function.zip main

# Deploy to AWS (assuming aws CLI is configured)
aws lambda create-function \
  --function-name share-bite-cleanup \
  --runtime go1.x \
  --role arn:aws:iam::ACCOUNT_ID:role/lambda-role \
  --handler main \
  --zip-file fileb://function.zip \
  --timeout 300 \
  --environment Variables='{
    CLEANUP_RETENTION_DAYS=30,
    CLEANUP_BATCH_SIZE=500,
    CLEANUP_DRY_RUN=false,
    CLEANUP_SCHEDULE_ENABLED=true,
    DB_DSN=postgresql://user:pass@host:5432/db,
    ...
  }'
```

#### Invoke Lambda Manually

```bash
# Invoke and capture response
aws lambda invoke \
  --function-name share-bite-cleanup \
  --log-type Tail \
  output.json

# View response
cat output.json

# View logs
aws lambda get-function-url-config \
  --function-name share-bite-cleanup
```

#### Lambda Response Format

```json
{
  "success": true,
  "total_found": 325,
  "total_deleted": 325,
  "dry_run": false,
  "duration": "15s",
  "results": [
    {
      "name": "expire-old-posts",
      "records_found": 250,
      "records_deleted": 250,
      "dry_run": false,
      "duration_ms": 12340
    },
    {
      "name": "cleanup-expired-tokens",
      "records_found": 75,
      "records_deleted": 75,
      "dry_run": false,
      "duration_ms": 4560
    }
  ]
}
```

### EventBridge Schedule

Cleanup jobs can be scheduled using AWS EventBridge to run on a regular schedule.

#### EventBridge Schedule Expression

EventBridge uses cron expressions for scheduling. Format: `cron(minute hour day month? day-of-week? year?)`

**Common Expressions**:

```
# Every day at 2:00 AM UTC
cron(0 2 * * ? *)

# Every Sunday at 3:00 AM UTC
cron(0 3 ? * SUN *)

# Every 6 hours (at 12 AM, 6 AM, 12 PM, 6 PM UTC)
cron(0 0,6,12,18 * * ? *)

# Every Monday, Wednesday, Friday at 1:00 AM UTC
cron(0 1 ? * MON,WED,FRI *)

# 1st of each month at 4:00 AM UTC
cron(0 4 1 * ? *)

# Every 30 minutes (special case using rate expression)
rate(30 minutes)

# Every 1 hour
rate(1 hour)

# Every 2 days
rate(2 days)
```

#### Setup EventBridge Rule

```bash
# Create the rule
aws events put-rule \
  --name share-bite-cleanup-schedule \
  --schedule-expression "cron(0 2 * * ? *)" \
  --state ENABLED

# Add Lambda target
aws events put-targets \
  --rule share-bite-cleanup-schedule \
  --targets "Id"="1","Arn"="arn:aws:lambda:REGION:ACCOUNT_ID:function:share-bite-cleanup","RoleArn"="arn:aws:iam::ACCOUNT_ID:role/eventbridge-lambda-role"

# Give EventBridge permission to invoke Lambda
aws lambda add-permission \
  --function-name share-bite-cleanup \
  --statement-id AllowEventBridgeInvoke \
  --action lambda:InvokeFunction \
  --principal events.amazonaws.com \
  --source-arn arn:aws:events:REGION:ACCOUNT_ID:rule/share-bite-cleanup-schedule
```

#### Monitor EventBridge Executions

```bash
# List recent invocations
aws events list-rule-names-by-target \
  --target-arn arn:aws:lambda:REGION:ACCOUNT_ID:function:share-bite-cleanup

# Check CloudWatch logs
aws logs tail /aws/lambda/share-bite-cleanup --follow
```

## Testing

### Unit Tests

Run all cleanup tests:

```bash
# Run all cleanup tests
go test ./internal/cleanup/...

# Run with coverage
go test -cover ./internal/cleanup/...

# Run specific test
go test -run TestExpireOldPostsDryRun ./internal/cleanup/service/

# Verbose output
go test -v ./internal/cleanup/...
```

### Integration Testing

Test with actual database:

```bash
# With PostgreSQL running locally
DB_DSN="postgresql://user:password@localhost:5432/share_bite_test" \
  CLEANUP_RETENTION_DAYS=1 \
  CLEANUP_BATCH_SIZE=10 \
  CLEANUP_DRY_RUN=false \
  go test -v ./internal/cleanup/repository/
```

### Dry-Run Testing

**Recommended for production testing**:

```bash
# Test cleanup logic without modifying database
CLEANUP_DRY_RUN=true \
CLEANUP_RETENTION_DAYS=30 \
./cleanup-worker

# Or with Lambda
aws lambda invoke \
  --function-name share-bite-cleanup \
  --payload '{}' \
  --env-var CLEANUP_DRY_RUN=true \
  output.json
```

## Idempotency Guarantee

The cleanup worker is fully idempotent and safe to retry:

1. **State-driven**: Deletes records based on timestamp, not sequences
2. **Batch-safe**: If a batch fails, only that batch needs to be retried
3. **No cursors**: Doesn't use pagination cursors that could timeout
4. **No side effects**: Each run is independent

**Safe to retry scenarios**:
- Network timeout during execution
- Lambda timeout (will retry from start of next batch)
- Database connection loss
- Partial batch failure

## Monitoring

### CloudWatch Metrics

Cleanup jobs should emit metrics:

```go
// Example metric emission (to be added to service)
cloudwatch.PutMetricData(&cloudwatch.PutMetricDataInput{
  Namespace: "ShareBite/Cleanup",
  MetricData: []*cloudwatch.MetricDatum{
    {
      MetricName: "RecordsDeleted",
      Value:      aws.Float64(float64(result.RecordsDeleted)),
      Unit:       cloudwatch.StandardUnitCount,
    },
  },
})
```

### Logs

Check logs for:

- **INFO**: Regular cleanup progress
- **WARN**: Records that couldn't be deleted (investigate)
- **ERROR**: Job failures (check database connectivity)

### Alerts

Recommended alerts:

- Job fails 2 times in a row
- Duration exceeds 5 minutes (timeout risk)
- No cleanup jobs run for 2 days (scheduling issue)

## Performance Considerations

### Batch Size Tuning

- **Small batches (50)**: Safer, less table locking, slower
- **Large batches (1000)**: Faster, but increases lock duration
- **Recommended**: 100-500 records per batch

### Index Usage

Ensure these indexes exist for performance:

```sql
-- On guest.posts
CREATE INDEX idx_posts_created_at ON guest.posts(created_at);
CREATE INDEX idx_posts_status_created ON guest.posts(status, created_at);

-- On auth.password_reset_tokens
CREATE INDEX idx_tokens_expires_at ON auth.password_reset_tokens(expires_at);
CREATE INDEX idx_tokens_used_at ON auth.password_reset_tokens(used_at);
```

### Execution Time

Typical execution times:
- **Dry-run**: < 1 second
- **Delete 100 records**: 1-2 seconds
- **Delete 1000 records**: 5-10 seconds
- **Delete 10000 records**: 30-60 seconds

## Troubleshooting

### Issue: Lambda timeout

**Symptoms**: Job doesn't complete, logs cut off mid-execution

**Solutions**:
1. Increase Lambda timeout (max 15 min)
2. Reduce `CLEANUP_BATCH_SIZE`
3. Reduce `CLEANUP_RETENTION_DAYS`
4. Check database performance

### Issue: Connection pool exhausted

**Symptoms**: "connection refused" errors

**Solutions**:
1. Check database connection limits
2. Verify database is accessible from Lambda VPC
3. Reduce concurrent Lambda invocations

### Issue: Some records not deleted

**Symptoms**: `RecordsFound > RecordsDeleted` in production (not dry-run)

**Solutions**:
1. Check for foreign key constraints
2. Verify `status != 'archived'` filter is correct
3. Check for row-level security (RLS) policies

### Issue: Performance degradation

**Symptoms**: Cleanup takes longer than expected, production queries slow down during cleanup

**Solutions**:
1. Run cleanup during off-peak hours
2. Reduce batch size to minimize lock duration
3. Run multiple smaller cleanup jobs instead of one large one
4. Add appropriate indexes

## Architecture

### Directory Structure

```
cmd/cleanup-worker/
  main.go              # CLI entry point

lambda/cleanup-worker/
  main.go              # Lambda handler

internal/cleanup/
  entity/
    cleanup.go         # Data types (CleanupResult, CleanupJob, etc.)
    cleanup_test.go    # Unit tests
  
  repository/
    cleanup_repository.go       # Interface
    postgres_cleanup.go         # PostgreSQL implementation
  
  service/
    cleanup_service.go          # Business logic
    cleanup_service_test.go     # Unit tests

internal/config/
  env/
    cleanup.go         # Config parsing
```

### Data Flow

```
EventBridge (schedule)
    ↓
Lambda Handler (lambda/cleanup-worker/main.go)
    ↓
Cleanup Service (internal/cleanup/service/)
    ↓
Cleanup Repository (internal/cleanup/repository/)
    ↓
PostgreSQL Database
```

Or locally:

```
CLI (cmd/cleanup-worker/main.go)
    ↓
Cleanup Service (internal/cleanup/service/)
    ↓
Cleanup Repository (internal/cleanup/repository/)
    ↓
PostgreSQL Database
```

## Future Enhancements

1. **Custom cleanup strategies**: Allow plugins for custom cleanup logic
2. **Metrics emission**: CloudWatch metrics for monitoring
3. **Parallel cleanup**: Run multiple cleanup jobs concurrently
4. **Cleanup scheduling**: Built-in cron scheduler (separate from EventBridge)
5. **Notification**: Send alerts on cleanup completion/failure
6. **Soft delete support**: Archive instead of hard delete for audit trails

## Database Schema Notes

### Posts table structure

```sql
CREATE TABLE guest.posts (
  id BIGSERIAL PRIMARY KEY,
  customer_id UUID NOT NULL,
  venue_id BIGINT NOT NULL,
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  ...
);

-- Related tables (cascade deleted)
CREATE TABLE guest.post_images (...);
CREATE TABLE guest.comments (...);
CREATE TABLE guest.post_likes (...);
```

### Password reset tokens structure

```sql
CREATE TABLE auth.password_reset_tokens (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  token_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ
);
```

## Emergency Procedures

### Manual cleanup via SQL (if worker fails)

```sql
-- WARNING: Only run if worker is broken

-- Dry run: count records
SELECT COUNT(*) FROM guest.posts WHERE created_at < NOW() - INTERVAL '30 days';

-- Actually delete (after verification)
DELETE FROM guest.posts 
WHERE created_at < NOW() - INTERVAL '30 days'
  AND status != 'archived'
LIMIT 1000;
-- Repeat in batches until done

-- Check tokens
SELECT COUNT(*) FROM auth.password_reset_tokens 
WHERE expires_at < NOW() AND used_at IS NULL;

-- Delete tokens
DELETE FROM auth.password_reset_tokens
WHERE expires_at < NOW() AND used_at IS NULL;
```

### Disable cleanup (if causing issues)

```bash
# Disable EventBridge rule
aws events disable-rule --name share-bite-cleanup-schedule

# Or set environment variable
export CLEANUP_SCHEDULE_ENABLED=false
```

## Questions?

See related documentation:
- [Database Schema](../architecture/database.md)
- [AWS Lambda Setup](../aws/lambda.md)
- [EventBridge Configuration](../aws/eventbridge.md)
