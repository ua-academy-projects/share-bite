# Guest API Local Kubernetes Guide

This document outlines how to deploy the `guest-api` to a local Kubernetes cluster (k3s, Docker Desktop, or Podman Desktop) and verify its functionality.

## Prerequisites

Ensure the base infrastructure (PostgreSQL, Redis, Migrations) is already running in your local cluster. See [local-kubernetes.md](./local-kubernetes.md) for details.

You also need to build the `guest-api` image locally so the cluster can use it:
```bash
make docker-build
```

## Apply Order

Deploy the application using `kustomize` (which applies the Deployment, Service, and ConfigMap):
```bash
kubectl apply -k deploy/k8s/guest
```

Wait for the Pod to reach the `Running` state:
```bash
kubectl get pods -n share-bite-local -l app=guest-api -w
```
*(Check the logs to ensure a successful database connection using `kubectl logs -l app=guest-api -n share-bite-local`).*

## Accessing the API Locally

To access the `guest-api` from your local machine, use port-forwarding:
```bash
kubectl port-forward svc/guest-api -n share-bite-local 3800:80
```

## Verification / Sample Curl

In a new terminal window, verify that the application is running and the health probe is working correctly:
```bash
curl -v http://localhost:3800/health
```

**Expected output:**
```json
{"status":"ok"}
```