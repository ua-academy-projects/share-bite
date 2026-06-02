# Operator Framework — Go SDK overview

Share Bite’s guest, business, and admin operators are Go controllers on the [Operator SDK](https://sdk.operatorframework.io/docs/overview/) stack. The notes below cover the pieces you need when writing reconcilers here.

## What is the Operator Framework?

The [Operator Framework](https://github.com/operator-framework) is an open source toolkit for building **Operators**: controllers that encode operational knowledge (deploy, scale, recover) for applications running on Kubernetes.

For Go, the main pieces are:

| Component | Role |
|-----------|------|
| **[Operator SDK](https://sdk.operatorframework.io/)** | CLI, project layout, scaffolding, testing helpers, optional OLM bundles |
| **[controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)** | Libraries your reconciler code calls (`Manager`, `Client`, `Reconciler`) |
| **[Kubebuilder](https://book.kubebuilder.io/)** | CRD/API codegen and project conventions; Operator SDK Go projects use the **Kubebuilder layout** (`go/v4`) |

Operator SDK sits on top of controller-runtime; reconciler code calls controller-runtime whether you scaffold with `operator-sdk init` or wire `main.go` yourself.

Docs worth bookmarking:

- [Operator SDK — Overview](https://sdk.operatorframework.io/docs/overview/)
- [Operator SDK — Go operators](https://sdk.operatorframework.io/docs/building-operators/golang/)
- [Operator SDK — Go quickstart](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/)
- [controller-runtime — Getting started](https://pkg.go.dev/sigs.k8s.io/controller-runtime)

## Typical Go operator workflow

From the [Operator SDK overview](https://sdk.operatorframework.io/docs/overview/):

1. **Initialize** a project (`operator-sdk init`, or equivalent layout under `operators/<name>/`).
2. **Define APIs** — CRD types in Go (`api/v1alpha1/…_types.go`) and YAML under `config/crd/` or `deploy/k8s/crd/`.
3. **Implement a controller** — struct that satisfies `reconcile.Reconciler`.
4. **Register** the controller with a `Manager` in `main.go`.
5. **Build & deploy** — container image + RBAC + Deployment (or `make run` locally against a dev cluster).

Per team we do steps 2–5 from that list; scaffold with Operator SDK/Kubebuilder or hand-copy the layout in [operators-overview.md](./operators-overview.md).

## Core concepts (controller-runtime)

Terms you will see in SDK docs and in our operator code.

### Manager

The **Manager** runs the operator process: shared client cache, leader election (when configured), metrics, health probes, and registered controllers.

```go
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
    Scheme: scheme,
})
```

Local dev (`make run-guest-operator`) typically starts the manager outside the cluster with kubeconfig pointing at k3s, Docker Desktop, or Podman Kubernetes.

### Controller and Reconcile loop

A **controller** watches one or more Kubernetes types (for example `GuestAppProfile`) and calls your **`Reconcile`** function when objects change.

```go
func (r *GuestAppProfileReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // 1. Load GuestAppProfile
    // 2. Compute desired state (replicas)
    // 3. Patch Deployment
    // 4. Update status.conditions (Ready)
    return ctrl.Result{}, nil
}
```

`Reconcile` should be **idempotent**: calling it repeatedly with the same spec should converge to the same cluster state. Return `ctrl.Result{RequeueAfter: …}` only when you need a delayed retry.

Register the controller in `SetupWithManager` with `ctrl.NewControllerManagedBy(mgr).For(&YourCRD{}).Complete(r)`. Share Bite operators watch only their AppProfile CR; they patch an existing Deployment rather than creating one, so `Owns()` is not required for v1.

### Client and API reader

Use the manager’s **client** (`mgr.GetClient()`) to `Get`/`List`/`Patch` resources (Deployments, your CRDs). Reads come from the informer cache; writes go to the API server. Update CR **status** via `r.Status().Update()` or `Patch()`, not `r.Update()` on the whole object. See the [controller-runtime Client API](https://sdk.operatorframework.io/docs/building-operators/golang/references/client/) reference.

### Custom Resource Definitions (CRDs)

A **CRD** registers a new API kind with the Kubernetes API server (for example `GuestAppProfile` in group `guest.sharebite.dev`). Users (or CI) apply CR instances; the controller reacts to them.

Share Bite uses a small spec (`replicas`, `enabled`, optional `deploymentName`) and a **status** subresource with `conditions`—see [operators-overview.md](./operators-overview.md).

### RBAC

Controllers need permission to read/write their CRDs and managed resources. Kubebuilder/Operator SDK generate `+kubebuilder:rbac` markers and `config/rbac/role.yaml`. Our manifests live under `deploy/k8s/operators/<domain>/`.

## Operator SDK CLI (optional for Share Bite)

Scaffold with the SDK or copy the layout from [operators-overview.md](./operators-overview.md). Handy commands:

| Command | Purpose |
|---------|---------|
| `operator-sdk init --domain sharebite.dev --repo github.com/ua-academy-projects/share-bite` | New operator module (use `go/v4` plugin on Apple Silicon) |
| `operator-sdk create api --group guest --version v1alpha1 --kind GuestAppProfile --resource --controller` | CRD types + controller stub |
| `make install run` | Install CRDs and run controller locally ([quickstart](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/)) |
| `make test` | Unit tests (envtest optional for integration) |

Production clusters here use **direct deploy** (manifests in `deploy/k8s/operators/`), not OLM bundles, unless you extend scope later.

## How Share Bite uses this stack

| Framework piece | Share Bite usage |
|-----------------|------------------|
| CRD + API types | Per team: `GuestAppProfile`, `BusinessAppProfile`, `AdminAppProfile` |
| Reconciler | Scale one Deployment; update `Ready` condition |
| Manager | One process per domain operator (`cmd/guest-operator`, etc.) |
| Operator SDK | Optional scaffolding; **controller-runtime** is required |
| OLM / bundles | Out of scope for first version |

```text
GuestAppProfile (CR)
        │
        ▼
guest-operator Reconcile()
        │
        ├── GET Deployment guest-api
        ├── PATCH spec.replicas
        └── UPDATE status.conditions
```

Business and admin operators follow the same pattern for `business-api` and `admin-auth-api`.

## Local development tips

- Run against a [local cluster](./local-kubernetes.md) with the target API Deployment already applied (#211 / #212 / #213).
- Prefer **leader election off** for single-replica local runs unless testing HA.
- Use `kubectl describe guestappprofile <name>` (actual resource name depends on CRD plural) to inspect status conditions.
- [Operator SDK FAQ](https://sdk.operatorframework.io/docs/faqs/) clarifies how controller-runtime, Kubebuilder, and the SDK relate.

## Further reading

| Topic | Link |
|-------|------|
| Go quickstart (memcached example) | https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/ |
| controller-runtime Client API | https://sdk.operatorframework.io/docs/building-operators/golang/references/client/ |
| Go tutorial (in-depth) | https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/ |
| Project layout (`go/v4`) | https://sdk.operatorframework.io/docs/building-operators/golang/project-layout/ |
| Testing operators | https://sdk.operatorframework.io/docs/building-operators/golang/testing/ |
| Kubebuilder book | https://book.kubebuilder.io/ |
| Kubernetes operator pattern | https://kubernetes.io/docs/concepts/extend-kubernetes/operator/ |

## Related (Share Bite)

- [Operators index](./operator.md)
- [Domain operator conventions](./operators-overview.md)
- [guest-operator](./guest-operator.md) · [business-operator](./business-operator.md) · [admin-operator](./admin-operator.md)
