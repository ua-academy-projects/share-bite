# guest-operator

Scales the **guest-api** Deployment from a `GuestAppProfile` resource.

Implementation: [#214](https://github.com/ua-academy-projects/share-bite/issues/214). Conventions: [operators-overview.md](./operators-overview.md). Stack: [operator-framework-go.md](./operator-framework-go.md).

## Quick reference

| Item | Value |
|------|--------|
| Binary | `cmd/guest-operator` |
| CRD | `GuestAppProfile` (`guest.sharebite.dev/v1alpha1`) |
| Default Deployment | `guest-api` |
| Manifests | `deploy/k8s/operators/guest/` |
| Run locally | `make run-guest-operator` |

Install and verification steps will land here once the operator is in the repo.
