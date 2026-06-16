## admin-auth-api Deployment Runbook

This guide outlines how to deploy the admin-auth-api microservice to a local Kubernetes cluster and verify its
functionality. The service manages authentication, tokens, session storage, and role-based permissions for the admin
platform.

## Service Overview

| Item             | Details                  |
|------------------|--------------------------|
| **Service Name** | `admin-auth-api`         |
| **Port**         | `3850`                   |
| **Namespace**    | `share-bite-local`       |
| **Image**        | `admin-auth-api:latest`  |
| **Manifests**    | `deploy/k8s/admin-auth/` |

## Prerequisites

Ensure the following steps are completed before deploying admin-auth-api:

- Local Kubernetes cluster is running (Docker Desktop, Podman Desktop, or k3s).
- kubectl is configured and points to your local cluster namespace.
- Shared infrastructure (postgres-0, redis-864c5d97f9-...) is deployed and healthy.
- The share-bite-migrator Job has completed successfully.

Refer to [local-kubernetes.md](./local-kubernetes.md) for initial cluster setup and bootstrapping.

## Deployment Structure

The `deploy/k8s/admin-auth/` directory contains:

```
deploy/k8s/admin-auth/
├── auth-deployment.yaml      # Deployment specification (with Always pull policy)
├── auth-service.yaml         # Service definition (port 3850)
├── auth-configmap.yaml       # Admin Auth local configuration overrides
└── kustomization.yaml        # Kustomize configuration file
```

## Configuration

Environment variables are sourced from:

1. **Shared Infra ConfigMap** (share-bite-infra-config): Shared global variables (Database DSNs, base system variables).

2. **Auth ConfigMap** (admin-auth-config): Service-specific settings (token TTLs, provider mappings, rates limits).

3. **Secrets** (share-bite-secrets): Holds critical cryptographic keys and credentials.

# Secure Handling of JWT Secrets

[!WARNING]
CRITICAL SECURITY RULE: Never commit actual JWT signing secrets, plain string passwords, or production .env credentials
to Git. All private keys must be handled via local Kubernetes Secrets.

# Local Generation Workflow:

To secure local authentication tokens, generate secure 32-byte cryptographic keys on your local machine and register
them within the cluster.

1. Generate random hex keys for Access and Refresh tokens:

# Generate Access Secret

```bash
openssl rand -hex 32
```

# Generate Refresh Secret

```bash
openssl rand -hex 32
```

2. Create or inject these values directly into your local cluster instance under the share-bite-secrets map:

``` bash
kubectl create secret generic share-bite-secrets \
--from-literal=JWT_ACCESS_SECRET_KEY="your_generated_access_key" \
--from-literal=JWT_REFRESH_SECRET_KEY="your_generated_refresh_key" \
--from-literal=POSTGRES_DB_PASSWORD="your_db_password" \
-n share-bite-local \
--dry-run=client -o yaml | kubectl apply -f -
```

## Deployment Order

> [!IMPORTANT]
> The **shared migrator Job must be applied and completed successfully** before deploying the business-api Deployment.
> Database migrations must run first to ensure the schema is ready.

### Step 1: Deploy admin-auth-api

From the project root, build and apply the admin-auth-api manifests using Kustomize:

```bash
kubectl apply -k deploy/k8s/admin-auth
```

This single command:

- Builds all manifests from `deploy/k8s/admin-auth/` (ConfigMap, Service, Deployment)
- Applies them to the cluster in the correct order
- Uses the pre-built `admin-auth-api:latest` image (ensure it's built and available in your container runtime)

### Step 2: Live Updates During Development

Because the deployment configuration utilizes imagePullPolicy: Always, you don't need to rebuild or wipe out the cluster
cache when updating Go files. Simply build the image and kick the deployment:

```bash
# Rebuild Go executable into Docker layer
docker build -t admin-auth-api:latest -f build/Dockerfile.admin .

# Force Kubernetes to fetch the fresh Docker SHA digest
kubectl rollout restart deployment admin-auth-api -n share-bite-local
```

### Step 3: Wait for Stability

Monitor the rolling update cycle to ensure zero downtime transitions:

```bash
kubectl rollout status deployment/admin-auth-api -n share-bite-local --timeout=90s
```

## Verification and Acceptance Criteria

### 1. Check Pod Status

Query the local namespace to ensure the pod is stabilized and working smoothly.

```bash
kubectl get pods -n share-bite-local -l app=admin-auth-api
```

Expected output:

```
NAME                              READY   STATUS    RESTARTS   AGE
admin-auth-api-6799986fd9-mdxwn   1/1     Running   0          45s
```

Detailed pod information:

```bash
kubectl describe pod -n share-bite-local -l app=admin-auth-api
```

### 2. Establish Port-Forwarding Tunnel

Forward the cluster network interface to your local host space:

```bash
kubectl port-forward deployment/admin-auth-api -n share-bite-local 3850:3850
```

You should see:

```
Forwarding from 127.0.0.1:3850 -> 3850
```

## 3. Verification Curl Flows

**A. Global Infrastructure Health Check**
Verify the root-level health endpoint exposed for Kubernetes liveness probes:

```bash
curl -i http://localhost:3850/health
```

Expected Response (200 OK):

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 16 Jun 2026 12:00:00 GMT
Content-Length: 15

{"status":"ok"}

**B. Authenticated Login/Registration Flow Discovery**
Verify that routing groups (/auth) are working properly by making a sample payload request to registration:
``` bash
curl -i -X POST http://localhost:3850/auth/register \
-H "Content-Type: application/json" \
-d '{"email":"test-user@sharebite.com", "slug":"user", password":"securepassword123"}'
```
## Troubleshooting

### Pod Fails to Start

```bash
kubectl describe pod -n share-bite-local -l app=admin-auth-api
kubectl logs -n share-bite-local -l app=admin-auth-api
```

Common causes:

- **CrashLoopBackOff**: Check logs for configuration errors (missing env vars)
- **ImagePullBackOff**: Verify `admin-auth-api:latest` image is available in your container runtime
- **Pending**: Cluster may be out of resources or secrets not applied

### Connection to PostgreSQL Failed

Verify infra is running and accessible:

```bash
kubectl get pods -n share-bite-local -l app=postgres
kubectl port-forward svc/postgres -n share-bite-local 5432:5432
```

In another terminal:

```bash
PGPASSWORD=bite psql -h 127.0.0.1 -U share -d share-bite -c "SELECT 1"
```

### Migrations Not Run Before Deployment

If you accidentally deployed business-api before migrations:

1. Delete the failed deployment:
   ```bash
   kubectl delete deployment admin-auth-api -n share-bite-local
   ```

2. Ensure migrations complete:
   ```bash
   kubectl wait --for=condition=complete --timeout=300s job/share-bite-migrator -n share-bite-local
   ```

3. Redeploy from root of the project:
   ```bash
   kubectl apply -k deploy/k8s/admin-auth
   ```

## Teardown

To remove the business-api deployment:

```bash
kubectl delete deployment admin-auth-api -n share-bite-local
kubectl delete service admin-auth-api -n share-bite-local
kubectl delete configmap share-bite-business-config -n share-bite-local
```

To remove the entire namespace (careful—this removes all services):

```bash
kubectl delete namespace share-bite-local
```