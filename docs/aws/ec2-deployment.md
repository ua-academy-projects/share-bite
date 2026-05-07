# Share Bite - Infrastructure & Deployment Guide

This guide outlines the infrastructure setup, automated CI/CD deployment process, and server management for the Share Bite microservices on AWS EC2.

## 1. Automated CI/CD Pipeline
Deployment is fully automated via **GitHub Actions**.
When code is merged or pushed to the `main` branch, the pipeline automatically:
1. Builds the Docker images.
2. Pushes them to AWS ECR.
3. Securely connects to the EC2 instance via AWS Systems Manager (SSM).
4. Updates the `compose.aws.yaml` file, pulls the latest images, and restarts the containers.

**No manual deployment steps are required for daily operations.**

---

## 2. AWS EC2 Setup Requirements (One-Time Setup)

### Instance Type Recommendation
- **Type:** `t3.small` (Minimum 2GB RAM is recommended to run multiple Go services + Postgres 18 concurrently).
- **OS:** Ubuntu Server 24.04 LTS.

### IAM Permissions (Crucial)
The EC2 instance **must** have an IAM Role (Instance Profile) attached with the following policies:
1. `AmazonEC2ContainerRegistryReadOnly` - To allow pulling Docker images from private ECR.
2. `AmazonSSMManagedInstanceCore` - To allow secure, passwordless terminal access via AWS SSM (Session Manager).

### Security Group (Inbound Rules)
Configure the EC2 Security Group. **For security, Port 22 (SSH) is strictly CLOSED.**
- **Port 80/443 (TCP):** HTTP/HTTPS (If using a reverse proxy/load balancer).
- **Ports 8080, 8081, 8082 (TCP):** Direct API access for testing (Restrict to your IP).

---

## 3. Initial Server Provisioning

Since SSH is disabled, all server management is done via AWS Systems Manager (SSM).

### Connect to the EC2 Instance
Use the AWS CLI on your local machine to start a session:
```bash
aws ssm start-session --target i-<YOUR_INSTANCE_ID>