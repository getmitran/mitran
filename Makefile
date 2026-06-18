.PHONY: dev build test test-go test-python test-ts lint lint-go lint-python lint-ts docker clean

dev:
	cd prototype/server && go run . &
	cd prototype/agent-worker && python -m flask --app worker.app run --port 8888 &
	cd prototype/dashboard && npm run dev

build:
	cd prototype/server && go build -o ../../bin/mitran-engine .
	cd prototype/dashboard && npm run build

test: test-go test-python test-ts

test-go:
	cd prototype/server && go test ./...

test-python:
	cd prototype/agent-worker && pytest tests/ -v

test-ts:
	cd prototype/dashboard && npx tsc --noEmit

lint: lint-go lint-python lint-ts

lint-go:
	cd prototype/server && gofmt -l .

lint-python:
	cd prototype/agent-worker && ruff check .

lint-ts:
	cd prototype/dashboard && npx prettier --check "src/**/*.{ts,tsx}"

docker:
	docker compose build

clean:
	rm -rf bin/ prototype/server/data/ prototype/dashboard/dist/
