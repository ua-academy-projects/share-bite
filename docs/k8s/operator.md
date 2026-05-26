# Share Bite Kubernetes operators

Share Bite runs three small domain operators, one per API team. Each one scales that team’s Deployment from a custom resource. Behavior matches across teams; binaries, CRDs, and manifests do not.

The code follows [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) with [Operator SDK](https://sdk.operatorframework.io/) layout. [operator-framework-go.md](./operator-framework-go.md) explains Manager, the reconcile loop, and CRDs. [operators-overview.md](./operators-overview.md) is the implementation contract ([#209](https://github.com/ua-academy-projects/share-bite/issues/209)).

## Operators

| Operator | Team | CRD | Default Deployment |
|----------|------|-----|-------------------|
| [guest-operator](./guest-operator.md) | Guest | `GuestAppProfile` | `guest-api` |
| [business-operator](./business-operator.md) | Business | `BusinessAppProfile` | `business-api` |
| [admin-operator](./admin-operator.md) | Admin | `AdminAppProfile` | `admin-auth-api` |

## Layout

```text
cmd/
├── guest-operator/main.go
├── business-operator/main.go
└── admin-operator/main.go

operators/
├── guest-operator/
├── business-operator/
└── admin-operator/

deploy/k8s/operators/
├── guest/       # CRD, RBAC, operator Deployment, example CR
├── business/
└── admin/
```

## Makefile (target names)

| Target | Action |
|--------|--------|
| `make run-guest-operator` | Run guest controller locally |
| `make run-business-operator` | Run business controller locally |
| `make run-admin-operator` | Run admin controller locally |
| `make test-operators` | `go test` all operator packages |

## Running multiple operators

All three can run in namespace `share-bite` at once. Each watches only its own CRD and Deployment; there is no shared reconciler process.

## Related

- [Operator Framework Go SDK](./operator-framework-go.md)
- [Local Kubernetes](./local-kubernetes.md)
- [Amazon EKS](./amazon-eks.md)
