.PHONY: all build run test clean docker

all: build

build:
	cd prototype/server && go build -o ../../bin/mitran-engine .
	cd prototype/cli && go build -o ../../bin/mitran .

run: build
	./bin/mitran-engine &
	cd prototype/agent-worker && python -m flask --app worker.app run --port 8888 &
	cd prototype/dashboard && npm run dev

test:
	cd prototype/server && go test ./... -v
	cd prototype/agent-worker && pytest tests/ -v || true
	cd prototype/dashboard && npx tsc --noEmit

clean:
	rm -rf bin/ prototype/server/data/

docker:
	docker compose up --build

docker-down:
	docker compose down -v

install-deps:
	cd prototype/dashboard && npm ci
	cd prototype/agent-worker && pip install -e .

lint:
	cd prototype/server && golangci-lint run ./...
	cd prototype/dashboard && npx tsc --noEmit
