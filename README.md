<picture>
  <source media="(prefers-color-scheme: dark)" srcset="redis-go-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="redis-go-light.svg">
  <img alt="redis-go" src="redis-go-light.svg" width="720">
</picture>

[![CI](https://github.com/zulerne/redis-go/actions/workflows/ci.yml/badge.svg)](https://github.com/zulerne/redis-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/zulerne/redis-go/graph/badge.svg?token=7PDUZ56Q6B)](https://codecov.io/gh/zulerne/redis-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/zulerne/redis-go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A Redis server implementation in Go, featuring RESP protocol support, in-memory data structures, and concurrent client handling.

## Features

- RESP (Redis Serialization Protocol) for client-server communication
- String, List, and Stream data types
- Thread-safe operations with goroutine-based connection handling
- Blocking list operations (BLPOP)
- TTL support with lazy deletion
- Graceful shutdown with structured logging (slog)

## Getting Started

### Prerequisites

- Go 1.26+
- [Task](https://taskfile.dev/) (optional, for task automation)

### Run

```bash
task run
# or
go run ./cmd/server
```

The server starts on `0.0.0.0:6379` by default.

### Build

```bash
task build
# or
go build -o redis-server ./cmd/server
```

### Connect

```bash
redis-cli -h localhost -p 6379
```

## Supported Commands

### Connection
- `PING` — test connectivity
- `ECHO <message>` — echo the given string

### String
- `SET <key> <value> [EX seconds] [PX milliseconds]` — set key with optional TTL
- `GET <key>` — get value by key

### List
- `RPUSH <key> <value> [value ...]` — append to the end
- `LPUSH <key> <value> [value ...]` — prepend to the start
- `LRANGE <key> <start> <stop>` — get range (supports negative indices)
- `LLEN <key>` — get length
- `RPOP <key> [count]` — remove and return from the end
- `LPOP <key> [count]` — remove and return from the start
- `BLPOP <key> <timeout>` — blocking pop from head

### Stream
- `XADD <key> <ID> <field> <value> [field value ...]` — add entry

### Generic
- `TYPE <key>` — get key type (string, list, stream, none)

## Architecture

```
cmd/server/          Entry point
internal/
  commands/          Command handlers
  config/            Server configuration
  server/            TCP server and connection handling
  store/             In-memory data store
pkg/
  resp/              RESP protocol parser and encoder
```

## Development

```bash
task test       # run tests with -race
task lint       # run golangci-lint
task check      # lint + test
task coverage   # test with coverage report
task fmt        # format code
```

## License

This project is open source and available for educational purposes.

**Note**: This is an educational Redis implementation. For production use, see [Redis](https://redis.io/).
