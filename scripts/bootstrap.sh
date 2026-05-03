#!/usr/bin/env bash
# =============================================================================
# scripts/bootstrap.sh
# One-shot bootstrap for the local Garage dev node.
# Run ONCE after the first `docker compose -f docker/compose.yaml up -d garage`.
# Safe to re-run: bucket/key creation will just print "already exists".
# =============================================================================

set -euo pipefail

CONTAINER="s3"
BUCKET="app-dev-bucket"
KEY_NAME="app-dev-key"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

# Shortcut: run any garage CLI command inside the container
garage() { docker compose -f docker/compose.yaml exec "$CONTAINER" /garage "$@"; }

echo "==> Waiting for Garage to be healthy..."
until docker compose -f docker/compose.yaml ps "$CONTAINER" | grep -q "healthy"; do
  sleep 1
done
echo -e "${GREEN}    Garage is healthy${RESET}"

# --------------------------------------------------------------------------
# 1. Cluster layout
# --------------------------------------------------------------------------
echo ""
echo "==> Step 1: Reading node ID..."
NODE_ID=$(garage status 2>&1 | awk '/NO ROLE ASSIGNED/{print $1; exit}')

if [[ -z "$NODE_ID" ]]; then
  echo "    Node already has a role. Skipping layout step."
else
  echo "    Node ID: $NODE_ID"
  echo "==> Step 1a: Assigning layout..."
  garage layout assign -z dc1 -c 1G "$NODE_ID"

  echo "==> Step 1b: Applying layout (version 1)..."
  garage layout apply --version 1
  echo -e "${GREEN}    Layout applied.${RESET}"
fi

# --------------------------------------------------------------------------
# 2. Bucket
# --------------------------------------------------------------------------
echo ""
echo "==> Step 2: Creating bucket '$BUCKET'..."
garage bucket create "$BUCKET" 2>&1 | grep -v "already exists" || true
echo -e "${GREEN}    Bucket ready.${RESET}"

# --------------------------------------------------------------------------
# 3. API key
# --------------------------------------------------------------------------
echo ""
echo "==> Step 3: Creating API key '$KEY_NAME'..."
EXISTING_KEY=$(garage key list 2>&1 | grep "$KEY_NAME" || true)
if [[ -z "$EXISTING_KEY" ]]; then
  KEY_CREATE_OUTPUT=$(garage key create "$KEY_NAME" 2>&1)
  SECRET_KEY=$(echo "$KEY_CREATE_OUTPUT" | awk '/Secret key:/{print $3}')
else
  echo "    Key already exists. Skipping."
  SECRET_KEY="(run: garage key info to get credentials)"
fi

echo "==> Step 3a: Granting read/write/owner access..."
garage bucket allow --read --write --owner "$BUCKET" --key "$KEY_NAME" 2>&1 || true
echo -e "${GREEN}    Key ready.${RESET}"

# --------------------------------------------------------------------------
# 4. Extract and print credentials
# --------------------------------------------------------------------------
KEY_INFO=$(garage key info "$KEY_NAME" 2>&1)

ACCESS_KEY=$(echo "$KEY_INFO" | awk '/Key ID:/{print $3}')

echo ""
echo -e "${BOLD}================================================================${RESET}"
echo -e "${YELLOW}  Copy S3_ACCESS_KEY and S3_SECRET_KEY values into your .env file:${RESET}"
echo -e "${BOLD}================================================================${RESET}"
echo ""
echo -e "  ${CYAN}S3_ACCESS_KEY${RESET}     = ${GREEN}${ACCESS_KEY}${RESET}"
echo -e "  ${CYAN}S3_SECRET_KEY${RESET}     = ${GREEN}${SECRET_KEY}${RESET}"

echo ""
echo -e "${BOLD}================================================================${RESET}"