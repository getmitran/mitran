FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /mitran-engine .

FROM alpine:3.19
RUN apk add --no-cache wget
COPY --from=builder /mitran-engine /usr/local/bin/mitran-engine
EXPOSE 7780
CMD ["mitran-engine"]
