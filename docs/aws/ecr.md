# Share Bite - Amazon ECR Setup & Usage Guide

This document outlines the setup and manual usage of the Elastic Container Registry (ECR) for Share Bite services.

## Available Repositories
The following private repositories have been provisioned:
- `share-bite/guest-api`
- `share-bite/business-api`
- `share-bite/admin-auth-api`
- `share-bite/migrator`

### Security & Lifecycle Configurations
- **Immutability:** All repositories have `IMMUTABLE` image tags enabled to prevent accidental overwriting of release artifacts.
- **Lifecycle Policy:** A rule is configured to automatically expire `Untagged` images when the count exceeds `5`. This prevents storage bloat.
- **Access Control:** The repositories use a least-privilege resource-based policy, restricting push access to authorized CI/CD roles/admins and pull access to AWS compute services (EC2/ECS).

---

## How to Authenticate and Push Locally

To push a Docker image from your local machine to ECR, ensure you have the [AWS CLI](https://aws.amazon.com/cli/) installed and configured (`aws configure`).

### 1. Authenticate Docker with your ECR registry
Retrieve an authentication token and authenticate your local Docker client to your registry.
*(Replace `<AWS_ACCOUNT_ID>` and `<REGION>` with your actual AWS details, e.g., `us-east-2`)*:
```bash
aws ecr get-login-password --region <REGION> | docker login --username AWS --password-stdin <AWS_ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com
```

### 2. Build your Docker image
Build your local service (e.g., Guest API).
```bash
docker build -t share-bite/guest-api:latest .
```

### 3. Tag the image for ECR
Tag the local image with the remote ECR repository URI and add a specific version (since we use immutable tags, avoid just using `:latest` for production).
```bash
docker tag share-bite/guest-api:latest <AWS_ACCOUNT_ID>.dkr.ecr.<REGION>[.amazonaws.com/share-bite/guest-api:v1.0.0](https://.amazonaws.com/share-bite/guest-api:v1.0.0)
```

### 4. Push the image
Push the tagged image to the AWS ECR repository.
```bash
docker push <AWS_ACCOUNT_ID>.dkr.ecr.<REGION>[.amazonaws.com/share-bite/guest-api:v1.0.0](https://.amazonaws.com/share-bite/guest-api:v1.0.0)
```