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

## Installation and Local Verification

1. Apply the CRD to your cluster:
   `kubectl apply -f deploy/k8s/operators/guest/crd.yaml`
2. Start the operator locally:
   `make run-guest-operator`
3. In a new terminal, apply the example profile:
   `kubectl apply -f deploy/k8s/operators/guest/guest-app-profile.example.yaml`
4. Verify scaling by editing the example file (`enabled: false`) and reapplying it. You should see the `guest-api` Deployment scale down to 0.