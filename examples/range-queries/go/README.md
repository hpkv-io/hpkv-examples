# HPKV Range Queries Example (Go)

This example demonstrates how to use HPKV's range queries functionality in Go. It shows how to:
1. Create sample records with sequential keys
2. Perform basic range queries
3. Use range queries with limits
4. Filter results based on record values

## Prerequisites

- Go 1.21 or higher
- HPKV API key and base URL

## Setup

1. Create a `.env` file in the project root with your HPKV credentials:
   ```
   HPKV_API_KEY=your_api_key_here
   HPKV_API_BASE_URL=your_api_base_url
   ```

2. Download dependencies:
   ```bash
   go mod tidy
   ```

## Running the Example

Run the example using the Go CLI:
```bash
go run hpkv_range_queries_example.go
```

## Example Output

The script will:
1. Create 10 sample user records with sequential IDs
2. Demonstrate three different range query scenarios:
   - Basic range query (users 1-5)
   - Range query with limit (users 1-10, limit 3)
   - Range query for users in New York (even IDs)
3. Clean up the sample records (commented out by default)

## Code Structure

- `HPKVRangeQueriesExample` struct handles all HPKV operations
- Methods:
  - `createSampleRecords()`: Creates sample user records
  - `performRangeQuery(startKey, endKey, limit)`: Performs range queries
  - `cleanupRecords()`: Removes sample records

## Data Structures

- `UserData`: Represents user information
  ```go
  type UserData struct {
      Name  string `json:"name"`
      Email string `json:"email"`
      Age   int    `json:"age"`
      City  string `json:"city"`
  }
  ```
- `Record`: Represents a key-value record
  ```go
  type Record struct {
      Key   string `json:"key"`
      Value string `json:"value"`
  }
  ```

## Error Handling

The example includes comprehensive error handling:
- Validates environment variables
- Handles API request errors
- Provides detailed error messages
- Uses Go's error return pattern
- Properly closes HTTP response bodies

## Dependencies

- `github.com/joho/godotenv`: For loading environment variables
- Standard library packages:
  - `net/http`: For HTTP requests
  - `encoding/json`: For JSON serialization
  - `bytes`: For request body handling
  - `fmt`: For string formatting

## Project Structure

- `hpkv_range_queries_example.go`: Main example implementation
- `go.mod`: Go module definition
- `.env`: Environment variables (not included in source control) 