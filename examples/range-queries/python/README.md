# HPKV Range Queries Example

This example demonstrates how to use HPKV's range query functionality in Python. Range queries allow you to retrieve multiple records within a specified key range efficiently.

## Prerequisites

- Python 3.7 or higher
- HPKV API key and base URL
- pip (Python package installer)

## Setup

1. Clone this repository or download the example files.

2. Create a virtual environment and activate it:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows, use: venv\Scripts\activate
   ```

3. Install the required dependencies:
   ```bash
   pip install -r requirements.txt
   ```

4. Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```

5. Edit the `.env` file and add your HPKV API key and base URL:
   ```
   HPKV_API_KEY=your_api_key_here
   HPKV_API_BASE_URL=your_api_base_url
   ```

## Running the Example

Run the example script:
```bash
python hpkv_range_queries_example.py
```

The example demonstrates three different use cases:

1. Basic range query: Retrieves all records between user:1 and user:5
2. Range query with limit: Retrieves up to 3 records between user:1 and user:10
3. Filtered range query: Retrieves users in New York from the range user:2 to user:10

## Example Output

The script will:
1. Create sample user records
2. Perform various range queries
3. Display the results
4. Clean up the sample records

## Understanding the Code

The example is structured into a class `HPKVRangeQueriesExample` that handles:
- Initialization with API credentials
- Creating sample records
- Performing range queries
- Cleaning up records

The `perform_range_query` method demonstrates how to use HPKV's range query endpoint with:
- Required parameters: `startKey` and `endKey`
- Optional parameter: `limit`

## Error Handling

The example includes basic error handling:
- Validates environment variables
- Uses `raise_for_status()` to check for HTTP errors
- Cleans up sample records in a `finally` block


## Additional Resources

- [HPKV Documentation](https://hpkv.io/docs)
- [HPKV REST API Reference](https://hpkv.io/docs/rest-api)
- [HPKV Best Practices](https://hpkv.io/docs/best-practices)
- [HPKV Dashboard](https://hpkv.io/dashboard) 