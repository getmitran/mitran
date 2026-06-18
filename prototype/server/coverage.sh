#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

echo "Running Go tests with coverage..."
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
echo ""
echo "HTML report: coverage.html"
go tool cover -html=coverage.out -o coverage.html
