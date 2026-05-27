# business-operator

Scales the **business-api** Deployment from a `BusinessAppProfile` resource.

Implementation: [#215](https://github.com/ua-academy-projects/share-bite/issues/215). Conventions: [operators-overview.md](./operators-overview.md). Stack: [operator-framework-go.md](./operator-framework-go.md).

## Quick reference

| Item | Value |
|------|--------|
| Binary | `cmd/business-operator` |
| CRD | `BusinessAppProfile` (`business.sharebite.dev/v1alpha1`) |
| Default Deployment | `business-api` |
| Manifests | `deploy/k8s/operators/business/` |
| Run locally | `make run-business-operator` |

---

## Overview

The **Business Operator** is a custom Kubernetes controller designed for the **Share-Bite** platform. Its primary responsibility is to orchestrate and manage the lifecycle of business-related backend deployments dynamically. By extending the Kubernetes API with a Custom Resource Definition (CRD), the operator automates the scaling and status tracking of applications based on custom business configurations, ensuring alignment with the project's domain operator conventions.

---

## Custom Resource Definition (CRD)

The operator manages a custom resource named `BusinessAppProfile`. This resource acts as a declarative contract defining the desired operational state of a target business deployment.

### Specification (`spec`)

| Field | Type | Description |
|:------|:-----|:------------|
| `deploymentName` | `string` | Optional. The exact name of the target `Deployment` to be managed. Defaults to `business-api`. |
| `enabled` | `boolean` | Global power switch. `true` scales to the desired replicas; `false` scales the deployment to `0`. |
| `replicas` | `int32` | The desired number of active pods to maintain when `enabled` is `true`. |

### Status (`status`)

| Field | Type | Description |
|:------|:-----|:------------|
| `conditions` | `array` | A list of structured operational states (e.g., `Ready=True` when observed replicas match the desired count). |

---

## Architecture & Reconciliation Loop

The core engine of the operator is its **Reconciliation Loop**, which continuously enforces synchronization between the declared `BusinessAppProfile` and the actual state of the cluster. The operator uses `controller-runtime` to watch both the CRD and the underlying `Deployment` for changes.

```text
+---------------------------------------------------------+
|                  Kubernetes API Server                  |
+---------------------------+-----------------------------+
                            |
                            | Watches Events (CRD & Deployment)
                            v
+---------------------------------------------------------+
|                Business Operator (Go)                   |
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

1. **Event Trigger:** The controller catches an event (Create/Update/Delete) for `BusinessAppProfile` or its owned `Deployment`.

2. **Resource Fetching:** Retrieves the CRD and determines the target deployment name.

3. **Drift Check:** Compares the actual `deployment.spec.replicas` against the desired state. If there is a drift, it issues a Patch request.

4. **Status Update:** Evaluates `deployment.status.readyReplicas`.
   - If they match desired replicas → sets `Ready=True` (Reason: `Scaled`).
   - If they do not match → sets `Ready=False` (Reason: `Scaling`).

5. **Error Resilience:** If the deployment is missing, sets `Ready=False` (Reason: `DeploymentNotFound`) and safely requeues.

---

## RBAC Configuration

To safely interact with the cluster, the operator requires the following RBAC permissions:

| API Group | Resources | Verbs | Purpose |
|:----------|:----------|:------|:--------|
| `apps` | `deployments` | `get, list, watch, update, patch` | Allows observing and scaling target deployments. |
| `business.sharebite.dev` | `businessappprofiles` | `get, list, watch, create, update, patch, delete` | Full lifecycle control over the CRD. |
| `business.sharebite.dev` | `businessappprofiles/status` | `get, update, patch` | Restricted permission to report status conditions. |
| `coordination.k8s.io` | `leases` | `get, list, watch, create, update, patch, delete` | Manages Leader Election locks. |
| `` (core) | `events` | `create, patch` | Writes controller events to the cluster. |

---

## Developer Guide

### Prerequisites

- A running Kubernetes cluster (e.g., local k3s or Docker Desktop).
- A configured `~/.kube/config`.

### Running Locally

You can use the project's Makefile to automatically apply the CRD, setup RBAC, and run the controller locally outside the cluster:

```bash
make run-business-operator
```

### Example Usage

To test the operator, apply a sample CR:

```yaml
apiVersion: business.sharebite.dev/v1alpha1
kind: BusinessAppProfile
metadata:
  name: business-api-local
  namespace: default
spec:
  replicas: 3
  enabled: true
  deploymentName: business-api
```

Observe the status changes using:

```bash
kubectl describe businessappprofile business-api-local
```