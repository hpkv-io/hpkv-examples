# HPKV Java CRUD Example

This example demonstrates how to perform CRUD (Create, Read, Update, Delete) operations with HPKV using Java.

## Prerequisites

- Java 11 or later
- Maven 3.6 or later
- HPKV API Key
- RHPKV Endpoint Url

## Setup

1. Clone the repository and navigate to the Java example directory:
   ```bash
   cd examples/basic-crud/java
   ```

2. Create a `.env` file from the example:
   ```bash
   cp .env.example .env
   ```

3. Edit the `.env` file and set your HPKV credentials:
   ```
   # Get this from https://hpkv.io/dashboard/api-keys
   HPKV_API_KEY=your_api_key_here
   HPKV_BASE_URL=your_hpkv_base_url
   ```

   You can obtain your API key and base URL from the HPKV dashboard.

4. Build the project:
   ```bash
   mvn clean package
   ```

## Running the Example

Run the example using:
```bash
java -jar target/hpkv-example-1.0-SNAPSHOT.jar
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

- `HPKVClient.java`: The main client class for interacting with HPKV
- `User.java`: Data model class for the example
- `Main.java`: Example usage of the client
- `pom.xml`: Maven project configuration and dependencies

## Additional Resources

- [HPKV Documentation](https://hpkv.io/docs)
- [HPKV REST API Reference](https://hpkv.io/docs/rest-api)
- [HPKV Best Practices](https://hpkv.io/docs/best-practices)
- [HPKV Dashboard](https://hpkv.io/dashboard) 