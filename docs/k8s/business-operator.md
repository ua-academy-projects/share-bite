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

Install and verification steps will land here once the operator is in the repo.
