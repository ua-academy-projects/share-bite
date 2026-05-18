# Multi-EC2 and AWS RDS PostgreSQL Deployment Guide

This document presents the verified distributed deployment architecture for the Share Bite application within AWS. Each service runs on an isolated EC2 instance, utilizing a shared managed Amazon RDS PostgreSQL database and an Nginx API Gateway.

---

## 1. Target Architecture & VPC Layout

The entire infrastructure is deployed inside **VPC:** `vpc-0de75c5ee1159a965`

### Compute Layer (EC2 Instances)

- **EC2 Instance 1 (Guest API):** Go / Gin (`10.0.136.8`, Internal Port: `3800`, Public Port: `8082`)
- **EC2 Instance 2 (Business API):** Go / Gin (`10.0.145.25`, Internal Port: `3900`, Public Port: `8081`)
- **EC2 Instance 3 (Admin & Auth API):** Go / Gin (`10.0.152.61`, Internal Port: `3850`, Public Port: `8080`)
- **EC2 Instance 4 (API Gateway / Reverse Proxy):** Nginx (`10.0.3.74`, Ports: `80`, `443`)

### Data Layer (Managed Service)

- **Amazon RDS PostgreSQL:** Shared Managed Database Core
    - **Endpoint:** `share-bite-db.cxmyqis8a0d9.us-east-2.rds.amazonaws.com`
    - **Port:** `5432`

*(Note: Formally decoupled from compute instances as a managed service within private DB subnets).*

---

## 2. Security Group Matrix (Least Privilege Verification)

The network configuration enforces a strict security perimeter based on three dedicated Security Groups:

### 1. `share-bite-proxy-sg` (`sg-0935400affebd0feb`)

Public group for Nginx reverse proxy. Allows inbound web traffic from the internet.

- **Inbound Rules:**
    - `80/TCP` (HTTP) from `0.0.0.0/0`
    - `443/TCP` (HTTPS) from `0.0.0.0/0`
    - `22/TCP` (SSH) from `0.0.0.0/0` (Developer access)

### 2. `share-bite-services-sg` (`sg-0228543305fc6ddb3`)

Private group for API microservices. Isolates applications from direct internet exposure.

- **Inbound Rules:**
    - `8080/TCP` (Admin API) from `sg-0935400affebd0feb` (Proxy Only)
    - `8081/TCP` (Business API) from `sg-0935400affebd0feb` (Proxy Only)
    - `8082/TCP` (Guest API) from `sg-0935400affebd0feb` (Proxy Only)
    - `8080/TCP`, `8081/TCP`, `8082/TCP`, `443/TCP` from `sg-0228543305fc6ddb3` (Allows secure Service-to-Service communication)
    - `22/TCP` (SSH) restricted **only** from Gateway IP `10.0.3.74/32`

### 3. `share-bite-rds-sg` (`sg-0610fb6359e951e60`)

Private group protecting the managed PostgreSQL layer.

- **Inbound Rules:**
    - `5432/TCP` (PostgreSQL) restricted **only** from `sg-0228543305fc6ddb3` (`share-bite-services-sg`)

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

All application EC2 instances run under a unified IAM Instance Profile:

- **Role Name:** `Training-GolangShareBiteEC2Role`
- **Role ARN:** `arn:aws:iam::897201144750:role/Training-GolangShareBiteEC2Role`

### Attached Permissions Policies

1. **AmazonEC2ContainerRegistryReadOnly** (AWS Managed): Grants permissions to pull compiled application Docker images from private AWS ECR repositories.
2. **AmazonSSMManagedInstanceCore** (AWS Managed): Enables secure infrastructure management and potential parameters lookup via Systems Manager.
3. **ec2-notifications-sqs-access** (Customer Inline Policy): Custom permission policy granting backend components access to AWS SQS queues for notification workloads.

---

## 5. Deployment Verification Checklist (Smoke Tests)

### ✓ Test 1: Public Gateway Health

```bash
curl -i http://localhost/gateway-health
# Expected Output: HTTP/1.1 200 OK
```

---

## 6. Restart and Rollback Procedures

### Service Restart

To safely restart a service container on a specific EC2 instance without altering the environment configuration:

```bash
docker compose -f build/aws/compose.<service-name>.yaml restart
```

### Rollback Process

If a newly deployed image causes failures or anomalies during smoke tests:

1. Revert the `IMAGE_TAG` variable inside the `.env` file to the last verified stable version tag.
2. Force the recreate lifecycle using Docker Compose to pull the previous state from ECR:

```bash
docker compose -f build/aws/compose.<service-name>.yaml up -d --force-recreate
```

---

## 7. Cost Optimization & Cleanup Instructions

To eliminate unnecessary AWS expenditures when the active testing or grading cycle is concluded:

- **EC2 Instances (Compute Layer):** Terminate all 4 deployed EC2 compute instances (`guest-api`, `business-api`, `admin-auth-api`, and the Nginx proxy). This stops instant per-second compute billing.
- **Amazon RDS PostgreSQL (Data Layer):** Delete the managed database instance via the AWS Console or CLI. Uncheck the "Create final snapshot" option to avoid ongoing storage costs for obsolete test schemas.
- **IAM Roles & Security Groups:** Keep or safely archive the `Training-GolangShareBiteEC2Role` and custom security groups since inactive network policies do not incur standalone infrastructure charges.