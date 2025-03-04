# HPKV WebSocket Example (C#)

This example demonstrates how to use the HPKV WebSocket API to perform CRUD operations using C#.

## Prerequisites

- .NET 6.0 or higher
- HPKV account and API key

## Project Setup

1. Create a new .NET project:
```bash
dotnet new console -n HPKVWebSocketExample
cd HPKVWebSocketExample
```

2. Add required NuGet packages:
```bash
dotnet add package System.Net.WebSockets.Client
dotnet add package System.Text.Json
```

## Configuration

Set the following environment variables:
```bash
export HPKV_BASE_URL=your_hpkv_server_url
export HPKV_API_KEY=your_api_key
```

Or on Windows:
```cmd
set HPKV_BASE_URL=your_hpkv_server_url
set HPKV_API_KEY=your_api_key
```

## Building

```bash
dotnet build
```

## Running

```bash
dotnet run
```

## Features

- WebSocket-based communication
- Full CRUD operations support
- Automatic reconnection handling
- Message correlation using message IDs
- SSL/TLS support
- API key authentication

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

```csharp
using (var client = new HPKVWebSocketClient(baseUrl, apiKey))
{
    await client.ConnectAsync();

    // Create a record
    await client.CreateAsync("key", value);

    // Read a record
    var value = await client.ReadAsync("key");

    // Update a record
    await client.UpdateAsync("key", newValue, partialUpdate: true);

    // Delete a record
    await client.DeleteAsync("key");
}
```

## Dependencies

- `System.Net.WebSockets.Client`: For WebSocket communication
- `System.Text.Json`: For JSON serialization/deserialization 