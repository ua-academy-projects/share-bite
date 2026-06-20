# admin-operator

Scales the **business-api** Deployment from a `AdminAppProfile` resource.

Implementation: [#216](https://github.com/ua-academy-projects/share-bite/issues/216).
Conventions: [operators-overview.md](./operators-overview.md).
Stack: [operator-framework-go.md](./operator-framework-go.md).

## Quick reference

| Item               | Value                                              |
|--------------------|----------------------------------------------------|
| Binary             | `cmd/admin-operator`                               |
| CRD                | `AdminAppProfile` (`admin.sharebite.dev/v1alpha1`) |
| Default Deployment | `admin-auth-api`                                   |
| Manifests          | `deploy/k8s/operators/admin-auth/`                 |
| Run locally        | `make run-admin-operator`                          |

---

## Overview

The **Admin Operator** is a custom Kubernetes controller designed for the **Share-Bite** platform. Its primary
responsibility is
to orchestrate, protect, and manage the lifecycle of the admin authentication and management backend deployments (
admin-auth-api) dynamically. By extending the Kubernetes API with a Custom Resource Definition (CRD), the operator
automates scaling, healing, and status tracking based on administrative profiles while adhering strictly to the
project's hardening and isolation standards.

---

## Custom Resource Definition (CRD)

The operator manages a custom resource named `AdminAppProfile`. This resource acts as a declarative contract defining
the desired operational state of a target business deployment.

### Specification (`spec`)

| Field            | Type      | Description                                                                                       |
|:-----------------|:----------|:--------------------------------------------------------------------------------------------------|
| `deploymentName` | `string`  | Optional. The exact name of the target `Deployment` to be managed. Defaults to `admin-auth-api`.  |
| `enabled`        | `boolean` | Global power switch. `true` scales to the desired replicas; `false` scales the deployment to `0`. |
| `replicas`       | `int32`   | The desired number of active pods to maintain when `enabled` is `true`.                           |

### Status (`status`)

| Field        | Type    | Description                                                                                                  |
|:-------------|:--------|:-------------------------------------------------------------------------------------------------------------|
| `conditions` | `array` | A list of structured operational states (e.g., `Ready=True` when observed replicas match the desired count). |

---

## Architecture & Reconciliation Loop

The core engine of the operator is its **Reconciliation Loop**, which continuously enforces synchronization between the
declared `AdminAppProfile` and the actual state of the cluster. The operator uses `controller-runtime` to watch both
the CRD and the underlying `Deployment` for changes.

```text
+---------------------------------------------------------+
|                  Kubernetes API Server                  |
+---------------------------+-----------------------------+
                            |
                            | Watches Events (CRD & Deployment)
                            v
+---------------------------------------------------------+
|                 Admin Operator (Go)                     |
|  Maintains loop to align actual state with desired spec |
+---------------------------+-----------------------------+
                            |
              +-------------+-------------+
              |                           |
      [enabled == true]           [enabled == false]
              |                           |
              v                           v
+--------------------------+  +--------------------------+
| Fetch Target Deployment  |  | Fetch Target Deployment  |
| Patch spec.replicas      |  | Patch replicas to 0      |
| Status: Ready / Scaled   |  | Status: Ready / Scaled   |
+--------------------------+  +--------------------------+
```

### Reconciliation Workflow

1. **Event Trigger:** The controller catches an event (Create/Update/Delete) for `AdminAppProfile` or its owned
   `Deployment`.

2. **Resource Fetching:** Retrieves the CRD and determines the target deployment name.

3. **Drift Check:** Compares the actual `deployment.spec.replicas` against the desired state. If there is a drift, it
   issues a Patch request.

4. **Status Update:** Evaluates `deployment.status.readyReplicas`.
    - If they match desired replicas → sets `Ready=True` (Reason: `Scaled`).
    - If they do not match → sets `Ready=False` (Reason: `Scaling`).

5. **Error Resilience:** If the deployment is missing, sets `Ready=False` (Reason: `DeploymentNotFound`) and safely
   requeues.

---

## RBAC Configuration

To safely interact with the cluster, the operator requires the following RBAC permissions:

| API Group             | Resources                    | Verbs                                             | Purpose                                            |
|:----------------------|:-----------------------------|:--------------------------------------------------|:---------------------------------------------------|
| `apps`                | `deployments`                | `get, list, watch, update, patch`                 | Allows observing and scaling target deployments.   |
| `admin.sharebite.dev` | `adminappprofiles`           | `get, list, watch, create, update, patch, delete` | Full lifecycle control over the CRD.               |
| `admin.sharebite.dev` | `businessappprofiles/status` | `get, update, patch`                              | Restricted permission to report status conditions. |
| `coordination.k8s.io` | `leases`                     | `get, list, watch, create, update, patch, delete` | Manages Leader Election locks.                     |
| `` (core)             | `events`                     | `create, patch`                                   | Writes controller events to the cluster.           |

---

## Security Context & Hardening

The admin-operator deployment is fully hardened for Production environments using advanced security features to
guarantee an immutable runtime footprint:

- Privilege Restriction: allowPrivilegeEscalation: false prevents the binary from gaining more permissions than its
  parent process.

- Rootless Execution: Runs strictly as a non-root user (runAsNonRoot: true, runAsUser: 1000 or 65532) to mitigate
  container breakout risks.

- Read-Only Root File System: readOnlyRootFilesystem: true enforces an immutable file system. No malware or arbitrary
  code can be written to the runtime container.

- Linux Capabilities: Linux kernel capabilities are entirely stripped using capabilities.drop: ["ALL"].

- Syscall Filtering: Uses seccompProfile: { type: RuntimeDefault } to strictly limit the allowed system calls to
  standard safe behaviors.

## Developer Guide

### Prerequisites

- A running Kubernetes cluster (e.g., local k3s or Docker Desktop).
- Local Kubernetes context set up `kubectl cluster-info`
- The project's central **Makefile** located in the root workspace directory.

### Automation Workflow (Using the Makefile)

The entire runtime environment can be stood up, managed, and torn down cleanly using predefined short targets:

1. Spin up the Core Admin Infrastructure
   Deploy the supporting stack (Postgres Database, Redis Cache, Admin Auth API API, and Database Migrations) via
   Kustomize:
   ```bash
   make run-auth-service
   ```
2. Build and Deploy the Admin Operator Container
   Automatically compile the Go codebase, package it into a secure Docker image, register the CRDs, apply the
   ClusterRBAC rules, and launch the operator Pod:
   ```bash
   make run-admin-operator
   ```

3. Apply the Custom Configuration Contract
   Trigger your desired operational state using a local sample custom resource profile:
   ```bash
   make apply-cr
   ```
4. Clean up and Teardown
   To cleanly wipe out the running infrastructure and controllers from your local cluster namespace without breaking
   configurations, execute:
   ```bash
   make stop-admin-operator
   make stop-auth-service
   ```

### Example Usage

To test the operator, apply a sample CR:

```yaml
apiVersion: admin.sharebite.dev/v1alpha1
kind: AdminAppProfile
metadata:
  name: admin-api-profile
  namespace: share-bite-local
spec:
  deploymentName: admin-auth-api
  replicas: 2
  enabled: true
```
