# HPKV Atomic Increment Example (Go)

This example demonstrates how to use HPKV to perform atomic increment operations on a key-value pair using Go.

## Prerequisites

- Go 1.21 or later
- HPKV API key and base URL

## Setup

1. Copy `.env.example` to `.env` and fill in your HPKV API credentials:
   ```bash
   cp .env.example .env
   ```

2. Edit the `.env` file with your actual HPKV API key and base URL.

3. Download the dependencies:
   ```bash
   go mod tidy
   ```

## Running the Example

To run the example:

```bash
go run atomic_increment.go
```

The example will:
1. Create a new key with an initial value of 0
2. Increment the value by 1
3. Increment the value by 5
4. Decrement the value by 2

## Code Structure

- `atomic_increment.go`: Main program file containing the implementation
- `go.mod`: Go module file with dependencies
- `.env`: Configuration file for HPKV credentials (not included in repository)
- `.env.example`: Example configuration file

## Dependencies

- github.com/joho/godotenv: For loading environment variables from .env file 