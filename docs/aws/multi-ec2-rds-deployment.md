# Multi-EC2 and AWS RDS PostgreSQL Deployment Guide

This document presents the verified distributed deployment architecture for the Share Bite application within AWS. Each service runs on an isolated EC2 instance, utilizing a shared managed Amazon RDS PostgreSQL database, asynchronous Redis event queues, and an Nginx API Gateway.

---

## 1. Target Architecture & VPC Layout

The entire infrastructure is deployed inside **VPC:** <AWS_VPC_ID>

### Compute Layer (EC2 Instances)

- **EC2 Instance 1 (Guest API):** Go / Gin (<GUEST_API_PRIVATE_IP>, Internal Port: `3800`, Public Port: `8082`)
- **EC2 Instance 2 (Business API):** Go / Gin (<BUSINESS_API_PRIVATE_IP>, Internal Port: `3900`, Public Port: `8081`)
- **EC2 Instance 3 (Admin & Auth API):** Go / Gin (<ADMIN_API_PRIVATE_IP>, Internal Port: `3850`, Public Port: `8080`)
- **EC2 Instance 4 (API Gateway / Reverse Proxy):** Nginx (<PROXY_PRIVATE_IP>, Ports: `80`, `443`)
- **EC2 Instance 5 (Workers / Background Services):** Go (<WORKERS_PRIVATE_IP>, Decoupled Asynchronous Processing)
    - *Services:* `share-bite-notifications` (Port: `4005`), `share-bite-outbox`
    - *Local Queue:* `share-bite-redis-workers` (Port: `6379`)

> 💡 **Architectural Note:** EC2 Instance 5 (Workers) operates purely as an event consumer. It handles asynchronous outbox processing and notification deliveries by polling the RDS database and pulling tasks from Redis queues. It does not accept direct inbound public HTTP traffic and therefore does not require an Nginx reverse proxy layer.

### Data Layer (Managed Services & Queues)

- **Amazon RDS PostgreSQL:** Shared Managed Database Core
    - **Endpoint:** <RDS_POSTGRESQL_ENDPOINT>
    - **Port:** `5432`
- **AWS ECR Private Container Registry:** Custom image storage
    - **Application Repositories:** share-bite/guest-api, share-bite/business-api, share-bite/admin-auth-api, share-bite/notifications-service, share-bite/outbox-worker
    - **Infrastructure Repository:** <AWS_ACCOUNT_ID>.dkr.ecr.<AWS_REGION>.amazonaws.com/

---

## 2. Security Group Matrix (Least Privilege Verification)

The network configuration enforces a strict security perimeter based on three dedicated Security Groups, optimized with internal loopback rules for secure cross-instance messaging:

### 1. `share-bite-proxy-sg` (<PROXY_SG_ID>)
Public group for Nginx reverse proxy. Allows inbound web traffic from the internet.
- **Inbound Rules:**
    - `80/TCP` (HTTP) from `0.0.0.0/0`
    - `443/TCP` (HTTPS) from `0.0.0.0/0`
    - `22/TCP` (SSH) from `0.0.0.0/0` (<YOUR_DEVELOPER_IP_CIDR>)

### 2. `share-bite-services-sg` (<SERVICES_SG_ID>)
Private group protecting all API microservices and background workers. Isolates application runtimes from direct internet exposure.
- **Inbound Rules:**
    - `8080/TCP` (Admin API) from <PROXY_SG_ID> (Proxy Only)
    - `8081/TCP` (Business API) from <PROXY_SG_ID> (Proxy Only)
    - `8082/TCP` (Guest API) from <PROXY_SG_ID> (Proxy Only)
    - `8080/TCP`, `8081/TCP`, `8082/TCP`, `443/TCP` from <SERVICES_SG_ID> (Secure Service-to-Service communication)
    - `6379/TCP` (Redis Mesh Rule):** Allowed from <SERVICES_SG_ID> (Self-referencing rule allowing the Worker node to interact with Redis brokers on Admin, Business, and Guest nodes).
    - `22/TCP` (SSH) restricted **only** from Gateway IP <PROXY_PRIVATE_IP>/32

### 3. `share-bite-rds-sg` (<RDS_SG_ID>)
Private group protecting the managed PostgreSQL database layer.
- **Inbound Rules:**
    - `5432/TCP` (PostgreSQL) restricted **only** from <SERVICES_SG_ID> (`share-bite-services-sg`), automatically granting database access to both the API nodes and the Worker node.

---

## 3. RDS Setup Instructions

1. **Engine:** PostgreSQL 15+ (or 18-alpine compatible features).
2. **Instance Class:** `db.t3.micro` (eligible for AWS Free Tier).
3. **Connectivity:**
    - Public accessibility: **No**.
    - Assigned to `share-bite-rds-sg`.
4. **Database Configuration:**
    - Initial Database Name: `sharebite_db`
    - Master Username: `sharebite`
5. **Backup & Maintenance:** Backup retention period set to 7 days; storage autoscaling enabled (min 20 GB).

---

## 4. Required IAM Roles for EC2 Access

All application and worker EC2 instances run under a unified IAM Instance Profile:

- **Role Name:** <IAM_ROLE_NAME>
- **Role ARN:** arn:aws:iam::<AWS_ACCOUNT_ID>:role/<IAM_ROLE_NAME>

### Attached Permissions Policies

1. **AmazonEC2ContainerRegistryReadOnly** (AWS Managed): Grants permissions to pull compiled application Docker images and the containerized custom Redis image from private AWS ECR repositories.
2. **AmazonSSMManagedInstanceCore** (AWS Managed): Enables secure infrastructure management and potential parameters lookup via Systems Manager.
3. **ec2-notifications-sqs-access** (Customer Inline Policy): Custom permission policy granting backend components access to AWS SQS queues for notification workloads.

---

## 5. Database Migrations Workflow (Fail-Fast Strategy)

To safely update the database schema without risking data inconsistency, migrations run as an ephemeral (one-off) job sequentially **before** application deployment.

### ⚠️ Concurrent Migration Risk & Network Isolation
- **Network Isolation:** Because AWS RDS blocks all direct public traffic, migrations must be executed from **EC2 Instance 3 (Admin & Auth API)**. This instance resides within `share-bite-services-sg`, which holds exclusive access to the database.
- **Concurrency Warning:** Never trigger migrations from multiple instances or parallel deployment pipelines simultaneously. Concurrent schema execution can trigger table locks, race conditions, or corrupt the `goose_db_version` schema execution history.
- **Fail-Fast Enforcement:** The migrator runs with specific Docker Compose boundaries. If a migration script contains syntactic/logical errors or connectivity drops, the process aborts immediately, blocking the deployment chain and protecting active production systems.

### Migration Execution Command:
```bash
docker-compose -f compose.migrator.yaml up --abort-on-container-exit --exit-code-from migrator
```

## 6. Deployment Verification Checklist (Smoke Tests)

### ✓ Test 1: Public Gateway Health

```bash
curl -i http://localhost/gateway-health
# Expected Output: HTTP/1.1 200 OK
```
# Verify connections to the remote application Redis queues
nc -zv <ADMIN_API_PRIVATE_IP> 6379   # Admin Redis -> Should return 'succeeded!'
nc -zv <BUSINESS_API_PRIVATE_IP> 6379   # Business Redis -> Should return 'succeeded!'
nc -zv <GUEST_API_PRIVATE_IP> 6379    # Guest Redis -> Should return 'succeeded!'

# Verify connection to the centralized database
nc -zv <RDS_POSTGRESQL_ENDPOINT> 5432 # PostgreSQL -> Should return 'succeeded!'
---

## 7. Restart and Rollback Procedures

### Service Restart

To safely restart a service container on a specific EC2 instance without altering the environment configuration:
```bash

# On Admin EC2 Node
docker compose -f compose.admin-auth-api.yaml restart

# On Business EC2 Node
docker compose -f compose.business-api.yaml restart

# On Workers EC2 Node
docker compose -f compose.workers.yaml restart
```

### Rollback Process
If a newly deployed image causes failures or anomalies during smoke tests:

Revert the IMAGE_TAG variable inside the localized .env file to the last verified stable version tag.

Force the recreate lifecycle using Docker Compose to pull the previous state from ECR:
```bash
docker compose -f compose.<service-name>.yaml up -d --force-recreate
```

---

## 8. Cost Optimization & Cleanup Instructions

To eliminate unnecessary AWS expenditures when the active testing or grading cycle is concluded:

**EC2 Instances (Compute Layer)**: Terminate all 5 deployed EC2 compute instances (guest-api, business-api, admin-auth-api, workers-background, and the Nginx proxy). This stops instant per-second compute billing.

**Amazon RDS PostgreSQL (Data Layer)**: Delete the managed database instance via the AWS Console or CLI. Uncheck the "Create final snapshot" option to avoid ongoing storage costs for obsolete test schemas.

**IAM Roles & Security Groups**: Keep or safely archive the customized roles and security groups since inactive network policies do not incur standalone infrastructure charges.