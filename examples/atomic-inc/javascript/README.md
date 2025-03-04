# HPKV Atomic Increment Example (JavaScript)

This example demonstrates how to use HPKV to perform atomic increment operations on a key-value pair using JavaScript (Node.js).

## Prerequisites

- Node.js 14 or later
- HPKV API key and base URL

## Setup

1. Copy `.env.example` to `.env` and fill in your HPKV API credentials:
   ```bash
   cp .env.example .env
   ```

2. Edit the `.env` file with your actual HPKV API key and base URL.

3. Install the dependencies:
   ```bash
   npm install
   ```

## Running the Example

To run the example:

```bash
npm start
```

The example will:
1. Create a new key with an initial value of 0
2. Increment the value by 1
3. Increment the value by 5
4. Decrement the value by 2

## Code Structure

- `atomic_increment.js`: Main program file containing the implementation
- `package.json`: Project file with dependencies
- `.env`: Configuration file for HPKV credentials (not included in repository)
- `.env.example`: Example configuration file

## Dependencies

- axios: For making HTTP requests
- dotenv: For loading environment variables from .env file 