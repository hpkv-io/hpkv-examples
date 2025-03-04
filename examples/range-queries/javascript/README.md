# HPKV Range Queries Example (JavaScript)

This example demonstrates how to use HPKV's range queries functionality in JavaScript. It shows how to:
1. Create sample records with sequential keys
2. Perform basic range queries
3. Use range queries with limits
4. Filter results based on record values

## Prerequisites

- Node.js (v14 or higher)
- npm (Node Package Manager)
- HPKV API key and base URL

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create a `.env` file in the project root with your HPKV credentials:
   ```
   HPKV_API_KEY=your_api_key_here
   HPKV_API_BASE_URL=your_api_base_url
   ```

## Running the Example

Run the example script:
```bash
npm start
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
  - `createSampleRecords()`: Creates sample user records
  - `performRangeQuery(startKey, endKey, limit)`: Performs range queries
  - `cleanupRecords()`: Removes sample records

## Error Handling

The example includes comprehensive error handling:
- Validates environment variables
- Handles API request errors
- Provides detailed error messages
- Uses try-catch blocks for error recovery

## Dependencies

- `axios`: For making HTTP requests
- `dotenv`: For loading environment variables 