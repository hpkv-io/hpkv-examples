# HPKV Examples

This repository contains example code with HPKV in multiple programming languages.

## Structure

```
examples/
├── atomic-inc/      # Atomic increment operations
├── basic-crud/      # Basic CRUD operations
├── range-queries/   # Range query operations
└── web-sockets/     # WebSocket-based operations
```

Each example category is available in multiple programming languages:
- C# (.NET)
- Go
- Java
- JavaScript (Node.js)
- Python

## Example Categories

### 1. Basic CRUD Operations
Demonstrates fundamental Create, Read, Update, and Delete operations with HPKV.
- Create new records
- Read existing records
- Update record values
- Delete records
- Verify operations

### 2. Atomic Increment Operations
Shows how to perform atomic increment/decrement operations on numeric values.
- Create a counter
- Increment values atomically
- Decrement values atomically
- Handle concurrent operations safely

### 3. Range Queries
Demonstrates how to efficiently query multiple records within a key range.
- Create sample records with sequential keys
- Perform basic range queries
- Use range queries with limits
- Filter results based on record values

### 4. WebSocket Operations
Shows how to use HPKV's WebSocket API for real-time operations.
- WebSocket-based communication
- Full CRUD operations support
- Automatic reconnection handling
- Message correlation
- SSL/TLS support
- API key authentication

## Prerequisites

Before running any of the examples, you'll need:

1. An HPKV API key
2. The HPKV base URL for your instance
3. The appropriate runtime/SDK for your chosen language:
   - C#: .NET 9.0 or later
   - Go: Go 1.16 or later
   - Java: Java 11 or later
   - JavaScript: Node.js 16 or later
   - Python: Python 3.7 or later

## Setup

1. Clone this repository
2. Navigate to your chosen example and language
3. Follow the setup instructions in the example's README.md
4. Configure your HPKV credentials (API key and base URL)

## Additional Resources

- [HPKV Documentation](https://hpkv.io/docs)
- [HPKV REST API Reference](https://hpkv.io/docs/rest-api)
- [HPKV Best Practices](https://hpkv.io/docs/best-practices)
- [HPKV Dashboard](https://hpkv.io/dashboard) 

## License

This project is licensed under the MIT License - see the individual example directories for details. 