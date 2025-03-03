# HPKV Go CRUD Example

This example demonstrates how to perform CRUD (Create, Read, Update, Delete) operations with HPKV using Go.

## Prerequisites

- Go 1.19 or later
- HPKV API Key
- RHPKV Endpoint Url

## Setup

1. Clone the repository and navigate to the Go example directory:
   ```bash
   cd examples/basic-crud/go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file from the example:
   ```bash
   cp .env.example .env
   ```

4. Edit the `.env` file and set your HPKV credentials:
   ```
   # Get this from https://hpkv.io/dashboard/api-keys
   HPKV_API_KEY=your_api_key_here
   HPKV_BASE_URL=your_hpkv_base_url
   ```

   You can obtain your API key and base URL from the HPKV dashboard.

## Running the Example

Run the example using:
```bash
go run .
```

The example will:
1. Create a new user record
2. Read the record
3. Update the user's age
4. Read the updated record
5. Delete the record
6. Verify the deletion

## Error Handling

The example includes error handling for:
- Missing or invalid credentials
- Network errors
- API errors
- Invalid responses

All errors are logged to stderr with descriptive messages.

## Code Structure

- `client.go`: The main client implementation for interacting with HPKV
- `main.go`: Example usage of the client
- `go.mod`: Go module file with dependencies

## Additional Resources

- [HPKV Documentation](https://hpkv.io/docs)
- [HPKV REST API Reference](https://hpkv.io/docs/rest-api)
- [HPKV Best Practices](https://hpkv.io/docs/best-practices)
- [HPKV Dashboard](https://hpkv.io/dashboard) 