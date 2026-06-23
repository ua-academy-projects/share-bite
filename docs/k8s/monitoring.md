# Monitoring

We use the `kube-prometheus-stack` Helm chart. It installs Prometheus, Grafana, and built-in collectors for node and Kubernetes metrics.

## What the chart gives you

- Node metrics: CPU, memory, disk, network per machine
- Kubernetes metrics: pod CPU/memory, container restarts, PVC usage
- API server, etcd, scheduler metrics

We add application metrics on top: HTTP requests, latency, Go runtime stats.

## Install

```bash
make monitoring-up
```

You can override the chart version:

```bash
CHART_VERSION=86.3.0 make monitoring-up
```

Or run manually:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f - || true
kubectl apply -f deploy/k8s/monitoring/grafana-secret.yaml

helm upgrade --install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
    --version 86.3.0 \
    --namespace monitoring \
    --create-namespace \
    --values ./deploy/k8s/monitoring/metrics-values.yaml
```

The `metrics-values.yaml` tells Prometheus to discover our ServiceMonitors in the `share-bite-local` namespace. Without it, Prometheus only looks inside `monitoring`.

## Access

```bash
make monitoring-forward-grafana
```

- Grafana: http://localhost:3000

```bash
make monitoring-forward-prometheus
```

- Prometheus: http://localhost:9090

## Dashboard

The dashboard is loaded _automatically_ via Kustomize. The source JSON is at `deploy/k8s/monitoring/share-bite-services.json`.

The dashboard has a service dropdown at the top. New services appear automatically once they start scraping.

## Adding a new service

1. Add a `/metrics` route to your application server.
2. Create a `ServiceMonitor` in your app's k8s config folder. It must have the label `release: kube-prometheus-stack`.
3. Your `Service` must have `metadata.labels.app` matching the selector in `ServiceMonitor`.
