# notifications-lambda — deployment & operational guide

Short, practical reference for deploying and operating the `notifications-lambda` Lambda (container image).

---

## Overview
- Purpose: consume SQS notification messages, validate and process them (current `Processor` republishes to Redis and can be extended to email/push).
- Package type: `Image` (ECR image)

---

## ECR image (publish & reference)

- Repository: `share-bite/notifications-lambda`
- Image URI pattern:

```text
<ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/<REPO>:<TAG>
```

Examples (CLI):

```powershell
# build
docker build --platform linux/arm64 --provenance=false -t notifications-lambda -f build/Dockerfile.notifications-lambda .

# tag (fill placeholders)
docker tag notifications-lambda:latest <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/<REPO>:<TAG>

# push (login via `aws ecr get-login-password` beforehand)
docker push <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/<REPO>:<TAG>
```

Use the image URI when updating Lambda code (see below).

---

## Lambda configuration (recommended minimal example)

Replace placeholders before running commands.

```json
{
  "FunctionName": "<FUNCTION_NAME>",
  "PackageType": "Image",
  "Role": "<LAMBDA_ROLE_ARN>",
  "MemorySize": 256,
  "Timeout": 20,
  "Architectures": ["arm64"],
  "Environment": {"Variables": {"REDIS_HOST":"<REDIS_HOST>","REDIS_PORT":"<REDIS_PORT>","APP_STAGE":"prod"}}
}
```

CLI: update image & config

```powershell
aws lambda update-function-code \
  --function-name <FUNCTION_NAME> \
  --image-uri <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/<REPO>:<TAG> \
  --region <REGION>

aws lambda update-function-configuration \
  --function-name <FUNCTION_NAME> \
  --memory-size 256 \
  --timeout 20 \
  --environment Variables={REDIS_HOST=<REDIS_HOST>,REDIS_PORT=<REDIS_PORT>,APP_STAGE=prod} \
  --region <REGION>
```

---

## SQS configuration (reliable delivery)

Example attributes to check/adjust (replace placeholders):

```json
{
  "QueueArn": "<SQS_QUEUE_ARN>",
  "VisibilityTimeout": "30",
  "MessageRetentionPeriod": "345600",
  "ReceiveMessageWaitTimeSeconds": "10",
  "RedrivePolicy": "{\"deadLetterTargetArn\":\"<DLQ_ARN>\",\"maxReceiveCount\":3}",
  "SqsManagedSseEnabled": "true"
}
```

Create / ensure event-source mapping (SQS → Lambda):

```powershell
aws lambda create-event-source-mapping \
  --function-name <FUNCTION_NAME> \
  --event-source-arn <SQS_QUEUE_ARN> \
  --batch-size 10 \
  --enabled \
  --function-response-types ReportBatchItemFailures \
  --region <REGION>
```

If mapping exists, list/update via `aws lambda list-event-source-mappings`.

---

## IAM role (Lambda) — example policy statements

Lambda role must allow CloudWatch logging and, if in VPC, ENI actions. Also ensure SQS permissions are included so the function can poll and acknowledge messages. Minimal statements:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    { "Effect": "Allow", "Action": ["logs:CreateLogGroup","logs:CreateLogStream","logs:PutLogEvents"], "Resource": "arn:aws:logs:*:*:*" },
    { "Effect": "Allow", "Action": ["ec2:CreateNetworkInterface","ec2:DescribeNetworkInterfaces","ec2:DeleteNetworkInterface"], "Resource": "*" },
    { "Effect": "Allow", "Action": ["sqs:ReceiveMessage","sqs:DeleteMessage","sqs:ChangeMessageVisibility","sqs:GetQueueAttributes","sqs:GetQueueUrl"], "Resource": "<SQS_QUEUE_ARN>" }
  ]
}
```

## Operational commands & checks

- Check deployed image tags:
  - `aws ecr describe-images --repository-name <REPO> --region <REGION>`
- Check Lambda config:
  - `aws lambda get-function-configuration --function-name <FUNCTION_NAME> --region <REGION>`
- Check SQS attributes:
  - `aws sqs get-queue-attributes --queue-url <QUEUE_URL> --attribute-names All`

---
