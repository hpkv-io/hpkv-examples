# HPKV WebSocket Example (JavaScript)

This example demonstrates how to use the HPKV WebSocket API to perform CRUD operations using JavaScript.

## Prerequisites

- Node.js 14 or higher
- npm or yarn
- HPKV account and API key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/hpkv-io/hpkv-examples.git
cd hpkv-examples
```

2. Install dependencies:
```bash
cd examples/web-sockets/javascript
npm install
```

## Configuration

Create a `.env` file in the project root with your HPKV credentials:
```bash
HPKV_BASE_URL=your_hpkv_server_url
HPKV_API_KEY=your_api_key
```

## Running

```bash
node hpkv_websocket_example.js
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

## Dependencies

- `dotenv`: For loading environment variables
- `ws`: For WebSocket communication

## Browser Usage

To use this example in a browser, you'll need to:

1. Bundle the code using a tool like webpack or rollup
2. Include the bundled file in your HTML
3. Use the client in your browser code:

```javascript
const client = new HPKVWebSocketClient(baseUrl, apiKey);
await client.connect();

// Create a record
await client.create('key', value);

// Read a record
const value = await client.read('key');

// Update a record
await client.update('key', newValue);

// Delete a record
await client.delete('key');
``` 