# HPKV Atomic Increment Example (C#)

This example demonstrates how to use HPKV to perform atomic increment operations on a key-value pair using C#.

## Prerequisites

- .NET 9.0 or later
- HPKV API key and base URL

## Setup

1. Copy `appsettings.json` and fill in your HPKV API credentials:
   ```json
   {
     "HPKV_BASE_URL": "your_base_url_here",
     "HPKV_API_KEY": "your_api_key_here"
   }
   ```

2. Edit the `appsettings.json` file with your actual HPKV API key and base URL.

3. Restore the NuGet packages:
   ```bash
   dotnet restore
   ```

## Running the Example

To run the example:

```bash
dotnet run
```

The example will:
1. Create a new key with an initial value of 0
2. Increment the value by 1
3. Increment the value by 5
4. Decrement the value by 2

## Code Structure

- `AtomicIncrement.cs`: Main program file containing the implementation
- `AtomicIncrement.csproj`: Project file with dependencies
- `appsettings.json`: Configuration file for HPKV credentials
- `README.md`: This file

## Dependencies

- Microsoft.Extensions.Configuration: For configuration management
- Microsoft.Extensions.Configuration.Json: For JSON configuration support
- Microsoft.Extensions.Configuration.EnvironmentVariables: For environment variable support 