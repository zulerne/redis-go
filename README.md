# redis-go

A Redis server implementation in Go, featuring RESP (Redis Serialization Protocol) support, in-memory data structures, and concurrent client handling.

## 🎯 Features

- **RESP Protocol**: Full implementation of Redis Serialization Protocol for client-server communication
- **Multiple Data Types**: String, List, and Stream data structures
- **Concurrency**: Thread-safe operations with goroutine-based connection handling
- **Blocking Operations**: Support for blocking list operations (BLPOP)
- **Key Expiration**: TTL support for string keys with automatic cleanup
- **Production Ready**: Graceful shutdown, structured logging, and comprehensive test coverage

## 🏗️ Architecture

```
redis-go/
├── cmd/server/          # Application entry point
├── internal/
│   ├── commands/        # Command handlers (connection, string, list, stream, generic)
│   ├── config/          # Server configuration
│   ├── server/          # TCP server and connection handling
│   └── store/           # In-memory data store with thread-safe operations
└── pkg/
    └── resp/            # RESP protocol parser and encoder
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25.4 or higher

### Installation

```bash
# Build the server
go build -o redis-server ./cmd/server

# Run the server
./redis-server
```

The server starts on `0.0.0.0:6379` by default.

### Connect with redis-cli

```bash
redis-cli -h localhost -p 6379
```

## 📝 Supported Commands

### Connection Commands
- `PING` - Test server connectivity
- `ECHO <message>` - Echo the given string

### String Commands
- `SET <key> <value> [EX seconds] [PX milliseconds]` - Set key with optional TTL
- `GET <key>` - Get value by key

### List Commands
- `RPUSH <key> <value> [value ...]` - Append values to the end of list
- `LPUSH <key> <value> [value ...]` - Prepend values to the start of list
- `LRANGE <key> <start> <stop>` - Get range of elements (supports negative indices)
- `LLEN <key>` - Get list length
- `RPOP <key> [count]` - Remove and return elements from the end
- `LPOP <key> [count]` - Remove and return elements from the start
- `BLPOP <key> <timeout>` - Blocking pop from list head

### Stream Commands
- `XADD <key> <ID> <field> <value> [field value ...]` - Add entry to stream

### Generic Commands
- `TYPE <key>` - Get the type of key (string, list, stream, none)

## 💡 Usage Examples

```bash
# String operations with TTL
SET mykey "Hello World" EX 60
GET mykey

# List operations
RPUSH mylist "item1" "item2" "item3"
LRANGE mylist 0 -1
LPOP mylist
LLEN mylist

# Blocking operations
BLPOP queue 5  # Wait up to 5 seconds

# Stream operations
XADD mystream 1-1 temperature 20 humidity 65

# Type checking
TYPE mykey
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/store
go test ./pkg/resp
```

## 🔧 Implementation Details

### RESP Protocol
The `pkg/resp` package provides a complete RESP protocol implementation:
- Parser for reading RESP-encoded data
- Encoders for all RESP data types (strings, errors, integers, bulk strings, arrays)
- Support for both simple and bulk strings
- Nil value handling

### Store
The in-memory store (`internal/store`) features:
- Thread-safe operations using `sync.RWMutex`
- Multiple data type support (string, list, stream)
- Key expiration with lazy deletion
- Blocking operations with timeout support
- Efficient memory management

### Server
The TCP server (`internal/server`) implements:
- Concurrent client handling with goroutines
- Graceful shutdown with connection draining
- Context-based cancellation
- Structured logging with slog

### Concurrency Model
- One goroutine per client connection
- Shared store protected by read-write locks
- Channel-based blocking operations for BLPOP
- Wait group for graceful shutdown

## 🎯 Design Principles

- **Clean Architecture**: Separation of concerns with clear boundaries between layers
- **Idiomatic Go**: Following Go best practices and standard library patterns
- **Performance**: Efficient data structures and minimal allocations
- **Testability**: Comprehensive unit tests with high coverage
- **Reliability**: Graceful error handling and safe concurrent access

## 📊 Performance Considerations

- **Lock Granularity**: RWMutex allows concurrent reads
- **Memory Efficiency**: Slice pre-allocation and string builders
- **Lazy Deletion**: Expired keys removed on access
- **Channel Buffering**: Optimized for blocking operations

## 🛣️ Roadmap

Potential future enhancements:
- Additional Redis commands (HSET, SADD, ZADD, etc.)
- Persistence (AOF, RDB snapshots)
- Pub/Sub messaging
- Transactions (MULTI/EXEC)
- Clustering support
- Metrics and monitoring

## 🤝 Contributing

Contributions are welcome! Areas for improvement:
- Add more Redis commands
- Optimize performance
- Improve test coverage
- Add benchmarks
- Documentation enhancements

## 📄 License

This project is open source and available for educational purposes.

---

**Note**: This is an educational implementation of Redis. For production use, please use the official [Redis](https://redis.io/) server.