# Auth Module - Wiring Instructions

Add to `main.go`:

```go
import "mitran/auth"

// In your route setup, register auth routes:
auth.RegisterRoutes(mux)

// Wrap your main handler with auth middleware:
handler := auth.AuthMiddleware(mux)
http.ListenAndServe(":7780", handler)
```

## Environment Variables

- `MITRAN_JWT_SECRET` - HMAC signing key. If empty, auth is disabled (dev mode).

## Default Credentials

- Username: `admin` / Password: `admin`

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /api/v1/auth/login | No | Get JWT token |
| POST | /api/v1/auth/token | Admin | Generate API key |
| GET | /api/v1/auth/me | Yes | Current user info |
| POST | /api/v1/auth/register | Admin | Create new user |
