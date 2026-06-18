# Security Guide

## CORS Configuration

Production: Set `MITRAN_CORS_ORIGINS` to your dashboard domain only.
```bash
export MITRAN_CORS_ORIGINS=https://app.yourdomain.com
```

Development: defaults to `*` (all origins) when `MITRAN_JWT_SECRET` is empty.

## Secret Rotation

1. Generate new JWT secret: `openssl rand -hex 32`
2. Set `MITRAN_JWT_SECRET=<new-secret>`
3. Restart engine: existing tokens invalidate, users re-login
4. For zero-downtime: support both old+new secrets during rotation window

## API Key Management

- Keys stored hashed (SHA-256) in the auth store
- Rotate via: POST /api/v1/auth/token (generates new), revoke old
- Set expiry on all keys: default 90 days

## Environment Variables (secrets)

| Variable | Purpose | Rotation |
|----------|---------|----------|
| MITRAN_JWT_SECRET | JWT signing | 90 days |
| AWS_ACCESS_KEY_ID | Bedrock LLM | Per AWS policy |
| GITHUB_TOKEN | GitHub API | 90 days |
| SLACK_BOT_TOKEN | Slack integration | On compromise |

## Best Practices

- Never commit .env files
- Use `.env.example` as template
- In production: use AWS Secrets Manager or Vault
- Enable rate limiting (100 req/s default)
- Enable request ID tracing for audit logs
