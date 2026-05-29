# EventBridge Schedule Expression Guide

## Overview

AWS EventBridge uses **cron expressions** to schedule the cleanup worker. This guide explains the syntax and provides examples.

## Cron Expression Syntax

```
cron(minute hour day month? day-of-week? year?)
```

### Fields

| Field | Values | Wildcards | Notes |
|-------|--------|-----------|-------|
| **minute** | 0-59 | , - * / | Required |
| **hour** | 0-23 | , - * / | UTC timezone. 0=midnight, 12=noon |
| **day** | 1-31 | , - * ? / | 1-based. Use `?` if `day-of-week` is specified |
| **month** | 1-12 | , - * / | Jan=1, Dec=12 |
| **day-of-week** | 1-7 | , - * ? / | 1=Sunday, 7=Saturday. Use `?` if `day` is specified |
| **year** | 1970-2199 | , - * / | Optional. Defaults to all years |

### Wildcard Meanings

- **`*`** — Any value
- **`?`** — No specific value (used for `day` OR `day-of-week`, never both)
- **`,`** — List separator: `1,3,5` = days 1, 3, and 5
- **`-`** — Range: `9-17` = 9 through 17 inclusive
- **`/`** — Increment: `0 */4` = every 4 hours starting at 0

## Common Schedules

### Daily

```cron
# Every day at 2:00 AM UTC
cron(0 2 * * ? *)

# Every day at 2:30 AM UTC
cron(30 2 * * ? *)

# Every day at noon UTC
cron(0 12 * * ? *)
```

### Multiple times per day

```cron
# At 2 AM, 8 AM, 2 PM, 8 PM UTC
cron(0 2,8,14,20 * * ? *)

# Every 6 hours (2 AM, 8 AM, 2 PM, 8 PM UTC)
cron(0 2 * * ? *)

# Every 4 hours starting at midnight
cron(0 0,4,8,12,16,20 * * ? *)
```

### Weekly

```cron
# Every Monday at 2:00 AM UTC
cron(0 2 ? * MON *)

# Every Monday and Friday at 2:00 AM UTC
cron(0 2 ? * MON,FRI *)

# Every Sunday at 3:00 AM UTC
cron(0 3 ? * SUN *)

# Monday-Friday (weekdays) at 2:00 AM UTC
cron(0 2 ? * MON-FRI *)

# Saturday-Sunday (weekends) at 2:00 AM UTC
cron(0 2 ? * SAT,SUN *)
```

### Monthly

```cron
# 1st of each month at 2:00 AM UTC
cron(0 2 1 * ? *)

# 15th of each month at 2:00 AM UTC
cron(0 2 15 * ? *)

# Last day of each month at 2:00 AM UTC
cron(0 2 L * ? *)

# Last Friday of each month at 2:00 AM UTC
cron(0 2 ? * 6L *)
```

### Specific dates

```cron
# January 1st at 2:00 AM UTC
cron(0 2 1 1 ? *)

# Christmas Day at 2:00 AM UTC
cron(0 2 25 12 ? *)

# Every October 15th at 3:00 AM UTC
cron(0 3 15 10 ? *)
```

## Rate Expressions (Alternative)

For simpler recurring schedules, use `rate()` expressions:

```
rate(value unit)
```

### Units

- **minute** or **minutes**
- **hour** or **hours**
- **day** or **days**

### Examples

```cron
# Every 30 minutes
rate(30 minutes)

# Every 1 hour
rate(1 hour)

# Every 2 hours
rate(2 hours)

# Every 1 day
rate(1 day)

# Every 7 days
rate(7 days)
```

## Cleanup Worker Recommendations

### Development

Run frequently for testing:

```cron
# Every hour
rate(1 hour)

# Or specific time daily
cron(0 2 * * ? *)
```

### Production

Choose based on database size and retention period:

#### Small database (< 100K posts/month)
```cron
# Once daily at 2 AM UTC
cron(0 2 * * ? *)
```

#### Medium database (100K-1M posts/month)
```cron
# Twice daily at 2 AM and 2 PM UTC
cron(0 2,14 * * ? *)
```

#### Large database (> 1M posts/month)
```cron
# Every 6 hours (midnight, 6 AM, noon, 6 PM UTC)
cron(0 0,6,12,18 * * ? *)
```

#### Very large database or high load
```cron
# Every 4 hours but only off-peak times
cron(0 1,5,9,13,17,21 * * ? *)

# Or use rate expression
rate(6 hours)
```

## Time Zones

**Important**: EventBridge always uses **UTC** for cron expressions.

### Converting from local time

| Your Timezone | UTC Offset | Convert To UTC |
|---------------|-----------|---|
| EST | UTC-5 | Add 5 hours |
| CST | UTC-6 | Add 6 hours |
| MST | UTC-7 | Add 7 hours |
| PST | UTC-8 | Add 8 hours |
| GMT | UTC±0 | No change |
| CET | UTC+1 | Subtract 1 hour |
| IST | UTC+5:30 | Subtract 5.5 hours |

### Examples

If you want 2 AM EST:
- EST = UTC-5
- 2 AM EST = 7 AM UTC
- **Expression**: `cron(0 7 * * ? *)`

If you want 3 PM PST:
- PST = UTC-8
- 3 PM PST = 11 PM UTC
- **Expression**: `cron(0 23 * * ? *)`

## Performance Considerations

### Batch size and schedule

| Dataset Size | Frequency | Config |
|---|---|---|
| Small | Daily | `CLEANUP_BATCH_SIZE=100` |
| Medium | Twice daily | `CLEANUP_BATCH_SIZE=500` |
| Large | Every 4-6 hours | `CLEANUP_BATCH_SIZE=1000` |

### Off-peak scheduling

Schedule cleanup during low activity:

```cron
# Off-peak in North America (early morning UTC)
cron(0 2 * * ? *)

# Off-peak globally (try UTC 2-4 AM)
cron(0 2,3 * * ? *)

# Or weekends only
cron(0 2 ? * SAT *)
```

### Avoiding overlap

Ensure cleanup takes less time than the schedule:

- **Schedule**: Every 4 hours
- **Timeout**: 5 minutes max (set Lambda timeout to 300s)
- **Safe**: ✓ Won't overlap

## Testing Schedules

Before deploying, test your cron expression:

```bash
# Use AWS cron expression validator
# https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-cron-expressions.html

# Or manually test with rate expressions first
rate(1 hour)

# Then switch to cron after testing
cron(0 2 * * ? *)
```

## Terraform / Infrastructure as Code

If using IaC, here are common configurations:

### Terraform

```hcl
resource "aws_cloudwatch_event_rule" "cleanup_schedule" {
  name                = "share-bite-cleanup-schedule"
  schedule_expression = "cron(0 2 * * ? *)"  # Daily at 2 AM UTC
  is_enabled          = true
}

resource "aws_cloudwatch_event_target" "cleanup_lambda" {
  rule      = aws_cloudwatch_event_rule.cleanup_schedule.name
  target_id = "CleanupLambda"
  arn       = aws_lambda_function.cleanup.arn
  role_arn  = aws_iam_role.eventbridge_lambda_role.arn
}
```

### AWS SAM

```yaml
Cleanup:
  Type: AWS::Lambda::Function
  Properties:
    Handler: main
    Runtime: go1.x
    CodeUri: ./lambda/cleanup-worker/
    Timeout: 300
    Events:
      DailyCleanup:
        Type: Schedule
        Properties:
          Schedule: 'cron(0 2 * * ? *)'  # Daily at 2 AM UTC
```

### AWS CloudFormation

```yaml
CleanupScheduleRule:
  Type: AWS::Events::Rule
  Properties:
    Name: share-bite-cleanup-schedule
    ScheduleExpression: 'cron(0 2 * * ? *)'
    State: ENABLED
    Targets:
      - Arn: !GetAtt CleanupLambda.Arn
        RoleArn: !GetAtt EventBridgeLambdaRole.Arn
```

## Troubleshooting

### Rule not triggering

1. **Check timezone**: Verify time is in UTC
2. **Check enabled state**: `aws events describe-rule --name share-bite-cleanup-schedule`
3. **Check target**: `aws events list-targets-by-rule --rule share-bite-cleanup-schedule`
4. **Check permissions**: Lambda must have permission from EventBridge

### Lambda not invoking

```bash
# Verify permission exists
aws lambda get-policy --function-name share-bite-cleanup

# Add if missing
aws lambda add-permission \
  --function-name share-bite-cleanup \
  --statement-id AllowEventBridgeInvoke \
  --action lambda:InvokeFunction \
  --principal events.amazonaws.com
```

### Schedule takes too long

If cleanup takes longer than the schedule interval:

1. **Reduce batch size**: `CLEANUP_BATCH_SIZE=100`
2. **Reduce retention period**: `CLEANUP_RETENTION_DAYS=14`
3. **Increase schedule interval**: `cron(0 0,6,12,18 * * ? *)` (every 6 hours)

## Examples

### Business hours cleanup (UTC)

```cron
# Every day at 1 AM UTC (before business hours)
cron(0 1 * * ? *)
```

### Weekend maintenance

```cron
# Saturday and Sunday at 2 AM UTC
cron(0 2 ? * SAT,SUN *)
```

### Conservative production

```cron
# Once per week on Sunday at 3 AM UTC
cron(0 3 ? * SUN *)
```

### Aggressive cleanup

```cron
# Every 4 hours on off-peak times
cron(0 1,5,9,13,17,21 * * ? *)
```

## Resources

- [AWS EventBridge Cron Expressions](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-cron-expressions.html)
- [AWS EventBridge User Guide](https://docs.aws.amazon.com/eventbridge/latest/userguide/)
- [Cron Expression Format](https://en.wikipedia.org/wiki/Cron)
