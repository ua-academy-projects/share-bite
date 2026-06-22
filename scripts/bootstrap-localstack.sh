#!/usr/bin/env bash
# =============================================================================
# scripts/bootstrap-localstack.sh
# Provision the local SNS topic + SQS queue + filtered subscription in LocalStack,
# mirroring terraform/main.tf (the `to_sse` subscription) for offline dev.
#
# Run ONCE after `docker compose -f build/compose.infra.yaml up -d localstack`.
# Safe to re-run: topic/queue/subscription creation is idempotent.
# =============================================================================

set -euo pipefail

CONTAINER="share-bite-localstack"
REGION="us-east-2"
ACCOUNT_ID="000000000000"
TOPIC_NAME="notifications"
QUEUE_NAME="notifications-sse"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

# Run awslocal inside the LocalStack container (addressed by name, so this works
# regardless of which compose project started it).
awslocal() {
  MSYS_NO_PATHCONV=1 docker exec -i "$CONTAINER" awslocal "$@"
}

echo "==> Waiting for LocalStack to be healthy..."
until [ "$(docker inspect -f '{{.State.Health.Status}}' "$CONTAINER" 2>/dev/null)" = "healthy" ]; do
  sleep 1
done
echo -e "${GREEN}    LocalStack is healthy${RESET}"

# --------------------------------------------------------------------------
# 1. SNS topic
# --------------------------------------------------------------------------
echo "==> Creating SNS topic '${TOPIC_NAME}'..."
TOPIC_ARN=$(awslocal sns create-topic --name "$TOPIC_NAME" --query 'TopicArn' --output text | tr -d '\r')
echo -e "${GREEN}    Topic: ${TOPIC_ARN}${RESET}"

# --------------------------------------------------------------------------
# 2. SQS queue (consumed by notifications-service)
# --------------------------------------------------------------------------
echo "==> Creating SQS queue '${QUEUE_NAME}'..."
awslocal sqs create-queue --queue-name "$QUEUE_NAME" >/dev/null
QUEUE_URL=$(awslocal sqs get-queue-url --queue-name "$QUEUE_NAME" --query 'QueueUrl' --output text | tr -d '\r')
QUEUE_ARN=$(awslocal sqs get-queue-attributes --queue-url "$QUEUE_URL" \
  --attribute-names QueueArn --query 'Attributes.QueueArn' --output text | tr -d '\r')
echo -e "${GREEN}    Queue: ${QUEUE_URL}${RESET}"

# --------------------------------------------------------------------------
# 3. Subscription with eventType filter policy (mirrors terraform `to_sse`)
# --------------------------------------------------------------------------
echo "==> Subscribing queue to topic with eventType filter policy..."
FILTER_POLICY='{"eventType":["post_liked","post_commented","post_mentioned","post_invitation_received","post_published","follow_added","business_verified","business_rejected"]}'

SUB_ARN=$(awslocal sns subscribe \
  --topic-arn "$TOPIC_ARN" \
  --protocol sqs \
  --notification-endpoint "$QUEUE_ARN" \
  --return-subscription-arn --output text | tr -d '\r')

# Set attributes separately — the --attributes shorthand can't carry JSON values.
awslocal sns set-subscription-attributes \
  --subscription-arn "$SUB_ARN" \
  --attribute-name RawMessageDelivery --attribute-value true
awslocal sns set-subscription-attributes \
  --subscription-arn "$SUB_ARN" \
  --attribute-name FilterPolicy --attribute-value "$FILTER_POLICY"
echo -e "${GREEN}    Subscription ready (raw delivery, eventType filtered).${RESET}"

# --------------------------------------------------------------------------
# 4. Print the values to put in .env
# --------------------------------------------------------------------------
echo ""
echo -e "${BOLD}================================================================${RESET}"
echo -e "${YELLOW}  Put these in your .env (host services reach LocalStack on :4566):${RESET}"
echo -e "${BOLD}================================================================${RESET}"
echo ""
echo -e "  ${CYAN}OUTBOX_SNS_TOPIC_ARN${RESET}        = ${GREEN}arn:aws:sns:${REGION}:${ACCOUNT_ID}:${TOPIC_NAME}${RESET}"
echo -e "  ${CYAN}OUTBOX_SNS_ENDPOINT_URL${RESET}     = ${GREEN}http://localhost:4566${RESET}"
echo -e "  ${CYAN}NOTIFICATION_SQS_QUEUE_URL${RESET}  = ${GREEN}http://localhost:4566/${ACCOUNT_ID}/${QUEUE_NAME}${RESET}"
echo -e "  ${CYAN}NOTIFICATION_SQS_ENDPOINT_URL${RESET} = ${GREEN}http://localhost:4566${RESET}"
echo -e "  ${CYAN}NOTIFICATION_AWS_REGION${RESET}     = ${GREEN}${REGION}${RESET}"
echo -e "  ${CYAN}AWS_ACCESS_KEY_ID${RESET}           = ${GREEN}test${RESET}"
echo -e "  ${CYAN}AWS_SECRET_ACCESS_KEY${RESET}       = ${GREEN}test${RESET}"
echo ""
