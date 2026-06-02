# admin-operator

Scales the **admin-auth-api** Deployment from an `AdminAppProfile` resource.

Implementation: [#216](https://github.com/ua-academy-projects/share-bite/issues/216). Conventions: [operators-overview.md](./operators-overview.md). Stack: [operator-framework-go.md](./operator-framework-go.md).

## Quick reference

| Item | Value |
|------|--------|
| Binary | `cmd/admin-operator` |
| CRD | `AdminAppProfile` (`admin.sharebite.dev/v1alpha1`) |
| Default Deployment | `admin-auth-api` |
| Manifests | `deploy/k8s/operators/admin/` |
| Run locally | `make run-admin-operator` |

Install and verification steps will land here once the operator is in the repo.
