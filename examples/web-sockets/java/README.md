# HPKV WebSocket Example (Java)

This example demonstrates how to use the HPKV WebSocket API to perform CRUD operations using Java.

## Prerequisites

- Java 11 or higher
- Maven
- HPKV account and API key

## Dependencies

Add the following dependencies to your `pom.xml`:

```xml
<dependencies>
    <dependency>
        <groupId>org.java-websocket</groupId>
        <artifactId>Java-WebSocket</artifactId>
        <version>1.5.3</version>
    </dependency>
    <dependency>
        <groupId>com.fasterxml.jackson.core</groupId>
        <artifactId>jackson-databind</artifactId>
        <version>2.15.2</version>
    </dependency>
    <dependency>
        <groupId>org.slf4j</groupId>
        <artifactId>slf4j-api</artifactId>
        <version>2.0.7</version>
    </dependency>
    <dependency>
        <groupId>ch.qos.logback</groupId>
        <artifactId>logback-classic</artifactId>
        <version>1.4.11</version>
    </dependency>
</dependencies>
```

## Configuration

Set the following environment variables:
```bash
export HPKV_BASE_URL=your_hpkv_server_url
export HPKV_API_KEY=your_api_key
```

## Building

```bash
mvn clean package
```

## Running

```bash
java -jar target/hpkv-websocket-example.jar
```

## Features

- WebSocket-based communication
- Full CRUD operations support
- Automatic reconnection handling
- Message correlation using message IDs
- SSL/TLS support
- API key authentication

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