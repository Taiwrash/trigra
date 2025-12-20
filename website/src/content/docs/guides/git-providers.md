---
title: Git Providers
description: How to configure Trigra with GitHub, GitLab, Gitea, and more.
---

Trigra is designed to be provider-agnostic. While it defaults to GitHub, it supports any Git provider that can send a webhook or is accessible via standard Git protocols.

## 🚀 Supported Providers

Trigra currently supports the following providers natively:

| Provider | `GIT_PROVIDER` | Type | Auto-Webhook |
|----------|----------------|------|:------------:|
| GitHub | `github` | API-driven | ✅ |
| GitLab | `gitlab` | API-driven | ✅ |
| Gitea / Forgejo | `gitea` | API-driven | ✅ |
| Bitbucket Cloud | `bitbucket` | API-driven | ✅ |
| Generic Git | `git` | Local Clone | ❌ |

---

## 🦊 GitLab

Trigra supports both **GitLab.com** and **Self-Managed** instances.

### Configuration
- `GIT_PROVIDER`: `gitlab`
- `GIT_TOKEN`: A Personal Access Token with `api` permissions.
- `GIT_BASE_URL`: (Optional) URL of your self-managed instance (e.g., `https://gitlab.mycompany.com`).

### Manual Webhook Setup
If you don't use the `PUBLIC_URL` auto-setup:
1. Go to your project -> **Settings** -> **Webhooks**.
2. URL: `http://<your-trigra-ip>/webhook`
3. Secret Token: Your `WEBHOOK_SECRET`.
4. Trigger: **Push events**.

---

## 🍵 Gitea / Forgejo

Gitea is a popular lightweight choice for homelabs and private servers.

### Configuration
- `GIT_PROVIDER`: `gitea`
- `GIT_TOKEN`: An Access Token generated in your profile settings.
- `GIT_BASE_URL`: The URL of your Gitea instance (e.g., `https://gitea.local`).

### Manual Webhook Setup
1. Go to Repository -> **Settings** -> **Webhooks**.
2. Add Webhook -> **Gitea**.
3. Target URL: `http://<your-trigra-ip>/webhook`
4. HTTP Method: `POST`
5. Secret: Your `WEBHOOK_SECRET`.

---

## 🪣 Bitbucket Cloud

Trigra integrates with Bitbucket Cloud using Basic Authentication.

### Configuration
- `GIT_PROVIDER`: `bitbucket`
- `BITBUCKET_USER`: Your Bitbucket username.
- `GIT_TOKEN`: A Bitbucket **App Password** with `Repository: Read` permissions.

### Manual Webhook Setup
1. Repository settings -> **Webhooks**.
2. URL: `http://<your-trigra-ip>/webhook`
3. Triggers: **Repository -> Push**.

---

## 🛠 Generic Git (SSH & HTTPS)

Use this for any Git server not listed above, or for local testing. This provider clones the repository locally into the controller.

### Configuration
- `GIT_PROVIDER`: `git`
- `GIT_REPO_URL`: The full Git URL (e.g., `git@github.com:user/repo.git` or `https://server.com/repo.git`).
- `GIT_SSH_KEY_FILE`: (Optional) Path to your SSH private key for private repositories.

### Triggering Syncs
Since generic Git doesn't have a standard API for webhook registration, you must manually trigger a sync by sending an empty POST request to Trigra:

```bash
curl -X POST http://<your-trigra-ip>/webhook
```

---

## 🔄 Automated Webhook Registration

Trigra can automatically register webhooks for **GitHub**, **GitLab**, and **Gitea**.

### Setup
Ensure these three variables are set:
1. `PUBLIC_URL`: The reachable address of your Trigra instance (e.g., `https://trigra.example.com`).
2. `GIT_TOKEN`: A token with permission to manage hooks.
3. `WEBHOOK_SECRET`: The secret Trigra will use to validate incoming events.

On startup, Trigra will check if a webhook exists for its `PUBLIC_URL` and create it if missing.
