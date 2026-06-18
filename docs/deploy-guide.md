# Deployment Guide

## Local Development

```bash
cp .env.example .env  # Fill in values
make install-deps
make run
```

Dashboard: http://localhost:5173 | Engine: http://localhost:7780 | Worker: http://localhost:8888

## Docker Compose

```bash
docker compose up --build
```

Dashboard: http://localhost:3000 | Engine: http://localhost:7780

## Kubernetes

```bash
kubectl create secret generic mitran-secrets --from-literal=jwt-secret=$(openssl rand -hex 32)
kubectl apply -f deploy/k8s/
```

## Environment Variables

See `.env.example` for full list. Required for production:
- `MITRAN_JWT_SECRET` -- JWT signing key
- `AWS_ACCESS_KEY_ID` + `AWS_SECRET_ACCESS_KEY` -- Bedrock access
- `MITRAN_CORS_ORIGINS` -- Dashboard domain

## Health Check

```bash
curl http://localhost:7780/health
```

Returns: status, uptime, version, goroutines, memory usage.

## Scaling

- Engine: single instance (stateful, SQLite)
- Worker: scale horizontally (stateless, 2+ replicas)
- Dashboard: scale horizontally (static assets)

For production with multiple engine instances, migrate to PostgreSQL.
