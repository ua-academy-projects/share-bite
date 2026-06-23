# Local Kubernetes Setup

This guide bootstraps Share Bite infrastructure on a local Kubernetes cluster and runs migrations in-cluster.
Items included: Infra: PostgreSQL, Redis, ConfigMap/Secret, and migrator Job.

It aligns with:
- `.env.example` defaults (`POSTGRES_PORT=5432`, `REDIS_PORT=6379`)
- `build/compose.infra.yaml`

## Prerequisites

- `kubectl`
- `make`
- One runtime: k3s, Podman Desktop, Docker Desktop, or kind
- A container builder available in your environment (`docker` or `podman`)

## 1) Create/enable a local cluster

### Docker Desktop

1. Open Docker Desktop.
2. Go to **Settings -> Kubernetes**.
3. Enable Kubernetes and wait until status is `Running`.
4. Verify context:
   ```bash
   kubectl config current-context
   kubectl cluster-info
   ```

### Podman Desktop

1. Open Podman Desktop and create/start a Kubernetes-enabled local cluster from the UI.
2. Ensure your kubeconfig context points to that cluster.
3. Verify:
   ```bash
   kubectl config current-context
   kubectl cluster-info
   ```

### k3s

1. Install k3s:
   ```bash
   curl -sfL https://get.k3s.io | sh -
   ```
2. Use k3s kubeconfig:
   ```bash
   export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
   kubectl config current-context
   kubectl cluster-info
   ```

### kind

1. Install kind:
    ```bash
    brew install kind # macOS
    # or see https://kind.sigs.k8s.io/docs/user/quick-start/
    ```
2. Create cluster:
    ```bash
    kind create cluster --name share-bite
    ```

## 2) Build and load local migrator image

The migrator Job uses `migrator:latest` (see `deploy/k8s/infra/migrator-job.yaml`).

Build image:
```bash
docker build -t migrator:latest -f build/Dockerfile.migrator .
```

If you use Podman instead of Docker:
```bash
podman build -t migrator:latest -f build/Dockerfile.migrator .
```

Runtime-specific image loading:
- Docker Desktop Kubernetes: no extra step (same local image store is used).
- Podman Desktop Kubernetes: if your cluster cannot see local images, export/import the image into the cluster runtime.
- k3s:
  ```bash
  docker save migrator:latest -o /tmp/migrator.tar
  sudo k3s ctr images import /tmp/migrator.tar
  ```
- kind:
    ```bash
    kind load docker-image migrator:latest
    ```

If you use `kind` and want to load all service images at once after building:
    ```bash
    make kind-load
    ```
This builds guest-api, business-api, admin-auth-api, and migrator, then loads them into the kind cluster.

## 3) Create local secrets (placeholders are committed only)

Do not edit and commit the example file directly.

```bash
cp docs/k8s/secrets.example.yaml docs/k8s/secrets.local.yaml
```

Edit `docs/k8s/secrets.local.yaml` and replace placeholder values.

Use `key: value` YAML syntax (not `key=value`).

## 4) Bootstrap order

1. Create namespace and apply secrets (required before Postgres starts):
   ```bash
   make k8s-secrets
   ```
2. Create infra (Postgres + Redis + ConfigMap):
   ```bash
   make k8s-up
   ```
3. Confirm infra pods/services:
   ```bash
   kubectl get pods -n share-bite-local
   kubectl get svc -n share-bite-local
   ```
4. Run migrations as a one-shot Job:
   ```bash
   make k8s-migrate
   ```
5. Confirm migration completed:
   ```bash
   kubectl get jobs -n share-bite-local
   kubectl logs job/share-bite-migrator -n share-bite-local
   ```

Run this order every time you bootstrap a new cluster: infra first, migration second.

## 5) Verification commands

Check cluster state:
```bash
kubectl get pods -n share-bite-local
kubectl get jobs -n share-bite-local
```

Port-forward Postgres and validate TCP reachability:
```bash
kubectl port-forward svc/postgres -n share-bite-local 15432:5432
```

In a second terminal:

**PowerShell (Windows):**
```powershell
Test-NetConnection -ComputerName 127.0.0.1 -Port 15432
```

**Bash / Git Bash (or use real curl on Windows):**
```bash
curl.exe -v telnet://127.0.0.1:15432
```

Port-forward Redis and validate TCP reachability:
```bash
kubectl port-forward svc/redis -n share-bite-local 16379:6379
```

In a second terminal:

**PowerShell (Windows):**
```powershell
Test-NetConnection -ComputerName 127.0.0.1 -Port 16379
```

**Bash / Git Bash:**
```bash
curl.exe -v telnet://127.0.0.1:16379
```

> On Windows, `curl` is often an alias for `Invoke-WebRequest`, which does not support `telnet://`. Use `Test-NetConnection` or `curl.exe` instead.

## 6) Tear down

```bash
make k8s-down
```

## Troubleshooting

### `ImagePullBackOff` for migrator Job

- Cause: cluster runtime cannot find `migrator:latest`.
- Fix:
  1. Rebuild image: `docker build -t migrator:latest -f build/Dockerfile.migrator .`
  2. Re-import image for your runtime (especially k3s/isolated runtimes).
  3. Re-run migration: `make k8s-migrate`.
  4. Inspect pod events:
     ```bash
     kubectl describe pod -n share-bite-local -l job-name=share-bite-migrator
     ```

### Postgres `CrashLoopBackOff` (postgres:18 volume mount)

- Symptom: logs mention data in `/var/lib/postgresql/data` and pg_ctlcluster compatibility.
- Cause: `postgres:18` requires the PVC mounted at `/var/lib/postgresql`, not `/var/lib/postgresql/data`.
- Fix after updating manifests:
  ```bash
  kubectl delete statefulset postgres -n share-bite-local
  kubectl delete pvc postgres-data-postgres-0 -n share-bite-local
  kubectl apply -k deploy/k8s/infra
  ```

### Database host misconfiguration

- Symptom: migrator fails with Postgres connection errors.
- Expected in-cluster host is `postgres` (Kubernetes Service), not `localhost` or `pg`.
- Verify effective config:
  ```bash
  kubectl get configmap share-bite-infra-config -n share-bite-local -o yaml
  kubectl get svc postgres -n share-bite-local
  ```

### Migration order issues

- Symptom: migrations fail, or DB schema is missing after bootstrap.
- Required order is strict:
  1. `make k8s-secrets`
  2. `make k8s-up`
  3. `make k8s-migrate`
- If migration already exists and must be rerun:
  ```bash
  kubectl delete job share-bite-migrator -n share-bite-local --ignore-not-found=true
  make k8s-migrate
  ```
