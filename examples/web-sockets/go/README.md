# HPKV WebSocket Example (Go)

This example demonstrates how to use the HPKV WebSocket API to perform CRUD operations using Go.

## Prerequisites

- Go 1.16 or higher
- HPKV account and API key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/hpkv-io/hpkv-examples.git
cd hpkv-examples
```

2. Install dependencies:
```bash
cd examples/web-sockets/go
go mod tidy
```

## Configuration

Set the following environment variables:
```bash
export HPKV_BASE_URL=your_hpkv_server_url
export HPKV_API_KEY=your_api_key
```

## Running

```bash
go run .
```

## Features

- WebSocket-based communication
- Full CRUD operations support
- Automatic reconnection handling
- Message correlation using message IDs
- SSL/TLS support
- API key authentication

## Project Structure

```
go/
├── hpkv_websocket_client.go  # WebSocket client implementation
├── main.go                   # Example usage
└── go.mod                    # Go module file
```

## Example Output

```
HPKV WebSocket CRUD Operations Example
=====================================

Using HPKV WebSocket server: wss://api-eu-2.hpkv.io/ws?apiKey=your_api_key

1. Creating a new user record...
Create operation succeeded

2. Reading the user record...
Retrieved data: {
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30
}

3. Updating the user's age...
Update operation succeeded

4. Reading the updated user record...
Retrieved data: {
  "name": "John Doe",
  "email": "john@example.com",
  "age": 31
}

5. Deleting the user record...
Delete operation succeeded

6. Attempting to read deleted record...
Record was successfully deleted
```

## Error Handling

The example includes comprehensive error handling for:
- Connection issues
- Authentication failures
- Invalid operations
- Network timeouts
- SSL/TLS errors

## Security

- API key is passed as a query parameter in the WebSocket URL
- SSL/TLS is enabled by default
- Certificate verification can be configured through the SSL context

## Using the Client

```go
client, err := NewHPKVWebSocketClient(baseURL, apiKey)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Create a record
err = client.Create("key", value)
if err != nil {
    log.Printf("Failed to create record: %v", err)
}

// Read a record
value, err := client.Read("key")
if err != nil {
    log.Printf("Failed to read record: %v", err)
}

// Update a record
err = client.Update("key", newValue, true)
if err != nil {
    log.Printf("Failed to update record: %v", err)
}

// Delete a record
err = client.Delete("key")
if err != nil {
    log.Printf("Failed to delete record: %v", err)
}
```

## Dependencies

- `github.com/gorilla/websocket`: For WebSocket communication
- `encoding/json`: For JSON serialization/deserialization

