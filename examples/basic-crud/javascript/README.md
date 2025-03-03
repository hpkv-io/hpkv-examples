# HPKV JavaScript CRUD Example

This example demonstrates how to perform CRUD (Create, Read, Update, Delete) operations with HPKV using JavaScript/Node.js.

## Prerequisites

- Node.js 16.x or later
- npm or yarn
- HPKV API Key
- RHPKV Endpoint Url

## Setup

1. Clone the repository and navigate to the JavaScript example directory:
   ```bash
   cd examples/basic-crud/javascript
   ```

2. Install dependencies:
   ```bash
   npm install
   # or
   yarn install
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
npm start
# or
yarn start
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

- `hpkv-client.js`: The main client class for interacting with HPKV
- `index.js`: Example usage of the client
- `package.json`: Project configuration and dependencies

## Additional Resources

- [HPKV Documentation](https://hpkv.io/docs)
- [HPKV REST API Reference](https://hpkv.io/docs/rest-api)
- [HPKV Best Practices](https://hpkv.io/docs/best-practices)
- [HPKV Dashboard](https://hpkv.io/dashboard) 