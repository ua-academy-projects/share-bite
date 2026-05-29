# AWS Lambda Cleanup Worker

Serverless cleanup job for ShareBite database maintenance.

## Quick Start

### Build

```bash
# Set Go to build for Linux
GOOS=linux GOARCH=amd64 go build -o main ./lambda/cleanup-worker

# Create deployment package
zip function.zip main
```

### Deploy

```bash
# Using AWS CLI
aws lambda create-function \
  --function-name share-bite-cleanup \
  --runtime go1.x \
  --role arn:aws:iam::ACCOUNT_ID:role/lambda-role \
  --handler main \
  --zip-file fileb://function.zip \
  --timeout 300 \
  --memory-size 256 \
  --environment Variables='{
    "CLEANUP_RETENTION_DAYS":"30",
    "CLEANUP_BATCH_SIZE":"500",
    "CLEANUP_DRY_RUN":"false",
    "DB_DSN":"postgresql://user:pass@host:5432/db",
    "CLEANUP_SCHEDULE_ENABLED":"true"
  }'
```

### Invoke

```bash
# Manual invocation
aws lambda invoke \
  --function-name share-bite-cleanup \
  --log-type Tail \
  output.json

# View response
cat output.json | jq
```

## Response Format

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

## EventBridge Integration

Schedule cleanup to run automatically:

```bash
# Create EventBridge rule
aws events put-rule \
  --name share-bite-cleanup-schedule \
  --schedule-expression "cron(0 2 * * ? *)" \
  --state ENABLED

# Add Lambda target
aws events put-targets \
  --rule share-bite-cleanup-schedule \
  --targets "Id"="1","Arn"="arn:aws:lambda:REGION:ACCOUNT_ID:function:share-bite-cleanup","RoleArn"="arn:aws:iam::ACCOUNT_ID:role/eventbridge-lambda-role"

# Grant EventBridge permission to invoke Lambda
aws lambda add-permission \
  --function-name share-bite-cleanup \
  --statement-id AllowEventBridgeInvoke \
  --action lambda:InvokeFunction \
  --principal events.amazonaws.com
```

## Configuration

Set via Lambda environment variables:

```
CLEANUP_RETENTION_DAYS=30         # Records older than this are cleaned
CLEANUP_BATCH_SIZE=500            # Records per batch (increased for Lambda)
CLEANUP_DRY_RUN=false             # Enable dry-run mode
CLEANUP_SCHEDULE_ENABLED=true     # Enable EventBridge scheduling
```

## Dry-Run Testing

Test before enabling in production:

```bash
# Invoke with dry-run
aws lambda invoke \
  --function-name share-bite-cleanup \
  --environment Variables='{"CLEANUP_DRY_RUN":"true"}' \
  output.json

# Check results (won't delete anything)
cat output.json | jq .total_deleted  # Should be 0
cat output.json | jq .total_found    # Shows what would be deleted
```

## Monitoring

### View logs

```bash
# Follow logs in real-time
aws logs tail /aws/lambda/share-bite-cleanup --follow

# View last 10 lines
aws logs tail /aws/lambda/share-bite-cleanup --max-items 10
```

### Metrics

Check CloudWatch for:
- **Duration**: Lambda execution time
- **Errors**: Invocation failures
- **Throttles**: Concurrent execution limits

## Troubleshooting

### Lambda timeout

Increase timeout and reduce batch size:

```bash
aws lambda update-function-configuration \
  --function-name share-bite-cleanup \
  --timeout 900 \
  --environment Variables='{"CLEANUP_BATCH_SIZE":"100"}'
```

### Database connection issues

Check VPC and security group:

```bash
# Lambda must be in same VPC as database
aws lambda get-function-vpc-config --function-name share-bite-cleanup

# Check database security group allows access from Lambda security group
aws ec2 describe-security-groups --group-ids sg-xxxxx
```

### Not executing on schedule

Verify EventBridge rule:

```bash
# Check rule is enabled
aws events describe-rule --name share-bite-cleanup-schedule

# Check targets
aws events list-targets-by-rule --rule share-bite-cleanup-schedule

# Check CloudWatch logs for rule executions
aws logs tail /aws/events/share-bite-cleanup-schedule
```

## Update Function

```bash
# Rebuild
GOOS=linux GOARCH=amd64 go build -o main ./lambda/cleanup-worker
zip function.zip main

# Update
aws lambda update-function-code \
  --function-name share-bite-cleanup \
  --zip-file fileb://function.zip
```

## IAM Policy

Lambda needs permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ec2:CreateNetworkInterface",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DeleteNetworkInterface"
      ],
      "Resource": "*"
    }
  ]
}
```

## Estimated Costs

- **Execution time**: 10-30 seconds typical
- **Memory**: 256 MB
- **Monthly cost**: ~$1-5 depending on frequency

Estimate: `(number of invocations × execution time × memory) / 1,600,000`

## See Also

- [Cleanup Worker Documentation](../cleanup-worker.md)
- [EventBridge Schedule Guide](../eventbridge-schedule-guide.md)
- [AWS Lambda Documentation](https://docs.aws.amazon.com/lambda/)
