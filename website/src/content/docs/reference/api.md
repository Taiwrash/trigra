---
title: API Reference
description: TRIGRA HTTP API endpoints
---

# API Reference

TRIGRA exposes a simple HTTP API for webhook processing and health checks.

## Endpoints

### POST /webhook

Receives GitHub webhook events.

| Status | Description |
|--------|-------------|
| `200 OK` | Success |
| `401 Unauthorized` | Invalid signature |
| `400 Bad Request` | Invalid payload |

### GET /health

Health check endpoint.

```bash
curl http://trigra-service/health
```

```json
{"status": "ok"}
```

### GET /ready

Readiness check endpoint.

### GET /

Service info.

```json
{
  "service": "trigra",
  "version": "1.0.0",
  "status": "running"
}
```

## Webhook Payload

TRIGRA processes GitHub push events:

```json
{
  "ref": "refs/heads/main",
  "repository": {
    "full_name": "user/repo",
    "clone_url": "https://github.com/user/repo.git"
  },
  "commits": [
    {
      "added": ["new-service.yaml"],
      "modified": ["deployment.yaml"]
    }
  ]
}
```

TRIGRA fetches and applies `.yaml` and `.yml` files from the commit.
