# Proto Definitions

gRPC service contracts between Go engine and Python workers.

## Generate Go code

```bash
# Option 1: buf (recommended)
buf generate

# Option 2: protoc
protoc --go_out=../server/grpc_gen --go-grpc_out=../server/grpc_gen \
  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
  mitran.proto
```

## Generate Python code

```bash
python -m grpc_tools.protoc -I. \
  --python_out=../agent-worker/worker/grpc_gen \
  --grpc_python_out=../agent-worker/worker/grpc_gen \
  mitran.proto
```

## Services

- **AgentWorker** — Engine dispatches tasks, worker streams progress back
- **ToolService** — Worker requests tool calls, engine handles approval + execution
