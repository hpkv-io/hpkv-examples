# HPKV Atomic Increment Example

This example demonstrates how to use the Atomic Increment operation with HPKV Key Value Store. The example shows how to increment and decrement a counter value atomically.

## Prerequisites

- Python 3.6 or higher
- HPKV account and API credentials

## Setup

1. Install the required dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. Create a `.env` file in the same directory with your HPKV credentials:
   ```
   HPKV_BASE_URL=your_base_url_here
   HPKV_API_KEY=your_api_key_here
   ```

   You can get these credentials from your HPKV dashboard at https://hpkv.io/dashboard/api-keys

## Running the Example

Run the example script:
```bash
python atomic_increment.py
```

The script will:
1. Increment a counter by 1
2. Increment the same counter by 5
3. Decrement the counter by 2

Each operation will print the result showing the new value after the operation.

## Code Explanation

The example demonstrates:
- How to set up HPKV client with API credentials
- How to perform atomic increment operations
- How to handle the API response
- Error handling for missing credentials

## API Documentation

For more information about the HPKV API, visit:
https://hpkv.io/docs/rest-api 