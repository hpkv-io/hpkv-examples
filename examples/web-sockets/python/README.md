# HPKV WebSocket Example

This example demonstrates how to use the HPKV WebSocket API to perform CRUD (Create, Read, Update, Delete) operations on key-value pairs.

## Features

- WebSocket-based communication for low latency
- Full CRUD operations support
- Automatic reconnection handling
- Message correlation using message IDs
- SSL/TLS support
- API key authentication

## Prerequisites

- Python 3.7 or higher
- HPKV account and API key
- Required Python packages (see requirements.txt)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/hpkv-io/hpkv-examples.git
cd hpkv-examples
```

2. Install the required packages:
```bash
pip install -r examples/web-sockets/python/requirements.txt
```

3. Create a `.env` file in the project root with your HPKV credentials:
```bash
HPKV_BASE_URL=your_hpkv_server_url
HPKV_API_KEY=your_api_key
```

## Usage

Run the example script:
```bash
python3 examples/web-sockets/python/hpkv_websocket_example.py
```

The example will demonstrate:
1. Creating a new user record
2. Reading the created record
3. Updating the user's age
4. Reading the updated record
5. Deleting the record
6. Verifying the deletion

## WebSocket Operations

The example implements the following WebSocket operations:

- **Create/Insert** (op: 2)
  ```python
  await client.create(key="user:1", value=user_data)
  ```

- **Read/Get** (op: 1)
  ```python
  await client.read(key="user:1")
  ```

- **Update** (op: 3)
  ```python
  await client.update(key="user:1", value=updated_data, partial_update=True)
  ```

- **Delete** (op: 4)
  ```python
  await client.delete(key="user:1")
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

## Troubleshooting

1. **Connection Issues**
   - Verify your HPKV server URL is correct
   - Check your internet connection
   - Ensure your API key is valid

2. **SSL/TLS Errors**
   - The example disables certificate verification by default
   - For production use, configure proper certificate verification

3. **Authentication Failures**
   - Double-check your API key in the `.env` file
   - Ensure the API key has the necessary permissions

## Contributing

Feel free to submit issues and enhancement requests! 