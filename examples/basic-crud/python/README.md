# HPKV Python CRUD Example

This example demonstrates basic CRUD (Create, Read, Update, Delete) operations using HPKV with Python. The example includes a simple client class that interacts with HPKV's REST API and a demonstration of common operations.

## Prerequisites

- Python 3.7 or higher
- pip (Python package installer)
- HPKV API Key (get it from [HPKV Dashboard](https://hpkv.io/dashboard/api-keys))
- HPKV Endpoint Urls

## Setup

1. Create a virtual environment (recommended):
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows, use: venv\Scripts\activate
   ```

2. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

3. Set up your HPKV credentials:
   - Copy the `.env.example` file to `.env`:
     ```bash
     cp .env.example .env
     ```
   - Edit `.env` and add your HPKV credentials:
     ```
     # Get this from https://hpkv.io/dashboard/api-keys
     HPKV_API_KEY=your_api_key_here  
     HPKV_BASE_URL=your_hpkv_base_url
     ```

## Running the Example

1. Make sure you have:
   - Created your `.env` file with valid credentials
   - Activated your virtual environment
   - Installed all dependencies

2. Run the example script:
   ```bash
   python hpkv_crud_example.py
   ```

The script will perform the following operations:
1. Create a new user record
2. Read the user record
3. Update the user's age
4. Read the updated record
5. Delete the record
6. Verify the deletion

## Configuration

The example can be configured in two ways:

1. Environment Variables (recommended):
   - Set `HPKV_API_KEY` and `HPKV_BASE_URL` in your `.env` file
   - Or set them in your environment:
     ```bash
     export HPKV_API_KEY=your_api_key_here
     export HPKV_BASE_URL=https://api-eu-2.hpkv.io
     ```

2. Direct Configuration:
   - Pass the values when creating the client:
     ```python
     client = HPKVClient(
         api_key="your_api_key_here",
         base_url="https://api-eu-2.hpkv.io"
     )
     ```

## Getting HPKV Credentials

1. Visit the [HPKV Dashboard](https://hpkv.io/dashboard/api-keys)
2. Log in to your account
3. Navigate to the API Keys section
4. Generate a new API key
5. Copy your API key and base URL
6. Add them to your `.env` file

## Example Output

When running the script, you should see output similar to this:

```
HPKV CRUD Operations Example
===========================

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

The example includes comprehensive error handling:
- Failed operations return `false` or `null`
- HTTP status codes are checked for success
- Try/catch blocks for network errors
- Detailed error logging
- Configuration validation

## Additional Resources

- [HPKV Documentation](https://hpkv.io/docs)
- [HPKV REST API Reference](https://hpkv.io/docs/rest-api)
- [HPKV Best Practices](https://hpkv.io/docs/best-practices)
- [HPKV Dashboard](https://hpkv.io/dashboard) 