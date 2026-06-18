FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o mitran-engine .

FROM alpine:3.19
RUN apk add --no-cache curl
COPY --from=builder /app/mitran-engine /usr/local/bin/
EXPOSE 7780
CMD ["mitran-engine"]
