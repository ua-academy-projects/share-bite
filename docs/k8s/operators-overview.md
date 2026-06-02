# Domain operator conventions

Guest, business, and admin each ship a separate Kubernetes operator. Treat this file as the shared contract so the three stay consistent without merging into one binary.

Each operator is a Go controller using [Operator SDK](https://sdk.operatorframework.io/) and [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime). Read [operator-framework-go.md](./operator-framework-go.md) for Manager, Reconcile, and CRD basics before you start coding.

## Principles

* **One operator per API team** — separate `main`, CRD, RBAC, and docs.
* **Same reconcile semantics** — scale a single named Deployment via CR spec.
* **Parallel development** — teams do not block each other’s PRs.

## CRD shape (all domains)

| Field | Type | Description |
|-------|------|-------------|
| `spec.replicas` | int32 | Desired replicas when `enabled: true` |
| `spec.enabled` | bool | When `false`, operator sets Deployment replicas to `0` |
| `spec.deploymentName` | string, optional | Target Deployment; default is the team’s standard name |

**Status**

| Field | Description |
|-------|-------------|
| `status.conditions[]` | At least `type: Ready` — `True` when observed replicas match desired |

## API groups

| Domain | Group | Kind | Default Deployment |
|--------|-------|------|-------------------|
| Guest | `guest.sharebite.dev/v1alpha1` | `GuestAppProfile` | `guest-api` |
| Business | `business.sharebite.dev/v1alpha1` | `BusinessAppProfile` | `business-api` |
| Admin | `admin.sharebite.dev/v1alpha1` | `AdminAppProfile` | `admin-auth-api` |

## Reconciler behavior

1. Load CR; compute desired replicas (`0` if `enabled: false`, else `spec.replicas`).
2. Get Deployment in the same namespace (name from spec or default).
3. If Deployment missing → `Ready=False`, record message, requeue.
4. Patch `deployment.spec.replicas` if drift.
5. Set `Ready=True` when the Deployment’s `status.readyReplicas` matches desired.

No admission webhooks or cross-namespace logic in the first version.

## Repository layout (per team)

```text
operators/<domain>-operator/
  api/v1alpha1/..._types.go
  internal/controller/..._controller.go
cmd/<domain>-operator/main.go
deploy/k8s/operators/<domain>/
docs/k8s/<domain>-operator.md
```

## Example CR (guest)

```yaml
apiVersion: guest.sharebite.dev/v1alpha1
kind: GuestAppProfile
metadata:
  name: guest-api-local
  namespace: share-bite
spec:
  replicas: 2
  enabled: true
```

## Verification checklist (each team)

- [ ] CRD applies cleanly
- [ ] Operator pod or `make run-*-operator` starts
- [ ] Example CR changes Deployment replica count
- [ ] `enabled: false` scales to zero
- [ ] Unit tests in CI

## Implementation notes

* Use controller-runtime `Reconciler` + `Manager` (`operator-sdk create api … --controller` or equivalent).
* Keep reconcilers idempotent; patch the Deployment only when replica count drifts.
* Write **status** from the controller, not spec; clients set spec, the operator sets `Ready`.
* Unit-test pure helpers (desired replica count); integration tests optional via envtest—see Operator SDK [testing guide](https://sdk.operatorframework.io/docs/building-operators/golang/testing/).

## Related

- [Operator Framework Go SDK](./operator-framework-go.md)
- [Operators index](./operator.md)
