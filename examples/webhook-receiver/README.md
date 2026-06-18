# Webhook Receiver Example

Minimal Go server that receives and validates Mitran webhook events.

## Usage

```bash
export WEBHOOK_SECRET="your-shared-secret"
go run main.go
```

The server listens on `:9090` for `POST /webhook`.

## Register in Mitran

In your Mitran project settings, add a webhook:

- **URL:** `http://localhost:9090/webhook`
- **Secret:** same value as `WEBHOOK_SECRET`

Mitran signs each payload with HMAC-SHA256 and sends the signature in the `X-Signature` header. This receiver validates it before processing.
