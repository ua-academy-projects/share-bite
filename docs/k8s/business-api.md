# business-api Deployment Runbook

This guide deploys the **business-api** microservice to a local Kubernetes cluster. The service provides core business logic and business post management functionality.

## Service Overview

| Item | Details |
|------|---------|
| **Service Name** | `business-api` |
| **Port** | `3900` |
| **Namespace** | `share-bite-local` |
| **Image** | `business-api:latest` |
| **Manifests** | `deploy/k8s/business/` |

## Prerequisites

Ensure the following steps are completed **before** deploying business-api:

- Local Kubernetes cluster is running (Docker Desktop, Podman Desktop, or k3s)
- `kubectl` is configured and points to your local cluster
- Infrastructure (PostgreSQL, Redis, ConfigMap, Secrets) is deployed
- Migrator Job has run successfully and completed

Refer to [local-kubernetes.md](./local-kubernetes.md) for initial cluster setup and bootstrapping.

## Deployment Structure

The `deploy/k8s/business/` directory contains:

```
deploy/k8s/business/
├── business-deploy.yaml      # Deployment specification
├── business-service.yaml      # Service definition (port 3900)
├── business-configmap.yaml    # Configuration overrides
└── kustomization.yaml         # Kustomize configuration
```

## Configuration

Environment variables are sourced from:

1. **Infrastructure ConfigMap** (`share-bite-infra-config`): Postgres, Redis, and shared settings from `.env.example`
2. **Business ConfigMap** (`share-bite-business-config`): Business-specific settings:
   - `H3_RESOLUTION=7`
   - `H3_RECOMMEND_RADIUS=2`
   - CORS, headers, and HTTP server settings

3. **Secrets** (`share-bite-secrets`): Sensitive values (API keys, credentials)

### S3/Garage Configuration

The business-api integrates with object storage (Garage or AWS S3) for image uploads and file management.

**For Local Development (Garage on Host Machine):**

When Garage runs outside the cluster (e.g., on your host via `compose.s3.yaml`), the business-api must access it through the host network:

- **Docker Desktop / Podman Desktop on macOS or Windows:**
  ```
  S3_ENDPOINT=http://host.docker.internal:4300
  S3_REGION=garage
  S3_ACCESS_KEY=<your-key>
  S3_SECRET_KEY=<your-secret>
  S3_BUCKET=app-dev-bucket
  S3_USE_PATH_STYLE=true
  ```

- **Linux with k3s or Host Bridge:**
  Replace `host.docker.internal` with your host's bridge IP (commonly `172.17.0.1` for Docker or check with `ip route show | grep docker`):
  ```
  S3_ENDPOINT=http://172.17.0.1:4300
  ```

- **Garage Running Inside Cluster:**
  If Garage is deployed as a Kubernetes Service in the same namespace:
  ```
  S3_ENDPOINT=http://garage:4300
  ```

**For AWS S3:**
  Leave `S3_ENDPOINT` empty and set `S3_REGION` to a real AWS region (e.g., `us-east-1`), with `S3_USE_PATH_STYLE=false`.

> **Note:** Update `docs/k8s/secrets.local.yaml` with your actual S3/Garage credentials before deployment.

## Deployment Order

> [!IMPORTANT]
> The **shared migrator Job must be applied and completed successfully** before deploying the business-api Deployment. Database migrations must run first to ensure the schema is ready.

### Step 1: Deploy business-api

From the project root, build and apply the business-api manifests using Kustomize:

```bash
cd deploy/k8s
kustomize build business | kubectl apply -f -
```

This single command:
- Builds all manifests from `deploy/k8s/business/` (ConfigMap, Service, Deployment)
- Applies them to the cluster in the correct order
- Uses the pre-built `business-api:latest` image (ensure it's built and available in your container runtime)

### Step 2: Verify Infrastructure and Migrator

Check that infrastructure is running and migrations have completed:

```bash
# Verify migrations have completed
kubectl wait --for=condition=complete --timeout=300s job/share-bite-migrator -n share-bite-local
```

### Step 3: Wait for Rollout

Monitor the deployment rollout:

```bash
kubectl rollout status deployment/business-api -n share-bite-local --timeout=180s
```

## Verification and Acceptance Criteria

### 1. Check Pod Status

Verify that the business-api pod is running:

```bash
kubectl get pods -n share-bite-local -l app=business-api
```

Expected output:
```
NAME                            READY   STATUS    RESTARTS   AGE
business-api-7d8f9c6b4a-k2mxj   1/1     Running   0          30s
```

Detailed pod information:
```bash
kubectl describe pod -n share-bite-local -l app=business-api
```

### 2. Port-Forward to Local Machine

Forward the service port to your localhost:

```bash
kubectl port-forward svc/business-api -n share-bite-local 3900:3900
```

You should see:
```
Forwarding from 127.0.0.1:3900 -> 3900
```

### 3. Health Check

In a separate terminal, verify the service is healthy:

```bash
curl http://localhost:3900/business/healthz -i
```

Expected response: `200 OK` (typically a JSON response with health status).

### 4. Readiness Check

```bash
curl http://localhost:3900/business/ready -i
```

Expected response: `200 OK` indicating the service is ready to handle traffic.

### 5. API Verification with Authentication

Test an authenticated endpoint. Replace `<AUTH_TOKEN>` with a valid JWT token (obtained from the auth flow) or generate a test token:

```bash
curl -X GET http://localhost:3900/business/posts \
  -H "Authorization: Bearer <AUTH_TOKEN>" \
  -H "Content-Type: application/json"
```

**Example with a placeholder token (for testing):**
```bash
curl -X GET http://localhost:3900/business/posts \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...." \
  -H "Content-Type: application/json" \
  -H "Accept: application/json"
```

Expected response: List of business posts or a 401 error if token is invalid (which confirms the endpoint exists and auth is enforced).

### 6. Log Inspection

View application logs:

```bash
kubectl logs -n share-bite-local -l app=business-api -f
```

Check for errors or warnings. Look for startup messages confirming successful connection to PostgreSQL, Redis, and S3/Garage.

## Troubleshooting

### Pod Fails to Start

```bash
kubectl describe pod -n share-bite-local -l app=business-api
kubectl logs -n share-bite-local -l app=business-api
```

Common causes:
- **CrashLoopBackOff**: Check logs for configuration errors (missing env vars, S3 connectivity)
- **ImagePullBackOff**: Verify `business-api:latest` image is available in your container runtime
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

### Connection to Garage/S3 Failed

Verify the S3 endpoint is reachable:

```bash
# From inside the business-api pod
kubectl exec -it -n share-bite-local -l app=business-api -- \
  curl -v http://host.docker.internal:4300  # or your S3 endpoint
```

Check S3 credentials in `docs/k8s/secrets.local.yaml`.

### Migrations Not Run Before Deployment

If you accidentally deployed business-api before migrations:

1. Delete the failed deployment:
   ```bash
   kubectl delete deployment business-api -n share-bite-local
   ```

2. Ensure migrations complete:
   ```bash
   kubectl wait --for=condition=complete --timeout=300s job/share-bite-migrator -n share-bite-local
   ```

3. Redeploy from `deploy/k8s/`:
   ```bash
   cd deploy/k8s
   kustomize build business | kubectl apply -f -
   ```

## Teardown

To remove the business-api deployment:

```bash
kubectl delete deployment business-api -n share-bite-local
kubectl delete service business-api -n share-bite-local
kubectl delete configmap share-bite-business-config -n share-bite-local
```

To remove the entire namespace (careful—this removes all services):

```bash
kubectl delete namespace share-bite-local
```