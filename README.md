# sentry-exporter

A Prometheus exporter of sentry event counts for organization.

## Configuration

| Environment variable name  | Required | Description |
|---|---|---|
| SENTRY_API_KEY| Yes | Scopes: [org:admin, org:read, org:write] |
| SENTRY_ORGANIZATION_SLUG | Yes | |
| SENTRY_API_ENDPOINT | No | Default: "https://sentry.io/api/0/" |
| WEB_LISTEN_ADDRESS| No | Default: ":9115"

## Metrics
```
curl localhost:9115/metrics
# HELP sentry_errors_total Total errors
# TYPE sentry_errors_total counter
sentry_errors_total{project="XXXXXXX",project_slug="yyyy"} 1
# HELP sentry_transactions_total Total transactions
# TYPE sentry_transactions_total counter
sentry_transactions_total{project="XXXXXXX",project_slug="zzzz"} 193
```
