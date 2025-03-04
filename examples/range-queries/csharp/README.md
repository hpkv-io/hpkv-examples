# HPKV Range Queries Example (C#)

This example demonstrates how to use HPKV's range queries functionality in C#. It shows how to:
1. Create sample records with sequential keys
2. Perform basic range queries
3. Use range queries with limits
4. Filter results based on record values

## Prerequisites

- .NET 7.0 SDK or higher
- HPKV API key and base URL

## Setup

1. Create a `.env` file in the project root with your HPKV credentials:
   ```
   HPKV_API_KEY=your_api_key_here
   HPKV_API_BASE_URL=your_api_base_url
   ```

2. Restore NuGet packages:
   ```bash
   dotnet restore
   ```

## Running the Example

Run the example using the .NET CLI:
```bash
dotnet run
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

- `HPKVRangeQueriesExample` class handles all HPKV operations
- Methods:
  - `CreateSampleRecordsAsync()`: Creates sample user records
  - `PerformRangeQueryAsync(startKey, endKey, limit)`: Performs range queries
  - `CleanupRecordsAsync()`: Removes sample records

## Error Handling

The example includes comprehensive error handling:
- Validates environment variables
- Handles API request errors
- Provides detailed error messages
- Uses try-catch blocks for error recovery

## Dependencies

- `DotEnv.Net`: For loading environment variables
- Built-in `System.Net.Http` for HTTP requests
- Built-in `System.Text.Json` for JSON serialization

## Project Structure

- `HPKVRangeQueriesExample.cs`: Main example implementation
- `HPKVRangeQueriesExample.csproj`: Project configuration file
- `.env`: Environment variables (not included in source control) 