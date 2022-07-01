# sentry-exporter

Export event counts for organization based on [an API](https://docs.sentry.io/api/organizations/retrieve-event-counts-for-an-organization-v2/).

## Configuration

| Environment variable name  | Required | Description |
|---|---|---|
| SENTRY_API_KEY| Yes | Scopes: [org:admin, org:read, org:write] |
| SENTRY_ORGANIZATION_SLUG | Yes | |
| SENTRY_API_ENDPOINT | No | Default: "https://sentry.io/api/0/" |
| WEB_LISTEN_ADDRESS| No | Default: ":9115"
