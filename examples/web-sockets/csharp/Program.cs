using System;
using System.Text.Json;
using System.Threading.Tasks;
using HPKV.Examples.WebSocket;

class Program
{
    static async Task Main(string[] args)
    {
        try
        {
            // Load environment variables
            var baseUrl = Environment.GetEnvironmentVariable("HPKV_BASE_URL");
            var apiKey = Environment.GetEnvironmentVariable("HPKV_API_KEY");
            
            if (string.IsNullOrEmpty(baseUrl) || string.IsNullOrEmpty(apiKey))
            {
                Console.Error.WriteLine("Please set HPKV_BASE_URL and HPKV_API_KEY environment variables");
                Environment.Exit(1);
            }
            
            // Initialize client
            using var client = new HPKVWebSocketClient(baseUrl, apiKey);
            await client.ConnectAsync();
            
            Console.WriteLine("HPKV WebSocket CRUD Operations Example");
            Console.WriteLine("=====================================");
            Console.WriteLine($"\nUsing HPKV WebSocket server: {baseUrl}/ws?apiKey={apiKey}");
            
            // Create operation
            var userData = new
            {
                name = "John Doe",
                email = "john@example.com",
                age = 30
            };
            
            Console.WriteLine("\n1. Creating a new user record...");
            var createSuccess = await client.CreateAsync("user:1", userData);
            if (!createSuccess)
            {
                Console.Error.WriteLine("Failed to create record. Exiting...");
                Environment.Exit(1);
            }
            Console.WriteLine("Create operation succeeded");
            
            // Read operation
            Console.WriteLine("\n2. Reading the user record...");
            var retrievedData = await client.ReadAsync("user:1");
            if (retrievedData != null)
            {
                Console.WriteLine($"Retrieved data: {JsonSerializer.Serialize(retrievedData, new JsonSerializerOptions { WriteIndented = true })}");
            }
            else
            {
                Console.WriteLine("Failed to retrieve data");
            }
            
            // Update operation
            Console.WriteLine("\n3. Updating the user's age...");
            var updatedUserData = new
            {
                name = "John Doe",
                email = "john@example.com",
                age = 31
            };
            var updateSuccess = await client.UpdateAsync("user:1", updatedUserData, true);
            Console.WriteLine($"Update operation {(updateSuccess ? "succeeded" : "failed")}");
            
            // Read after update
            Console.WriteLine("\n4. Reading the updated user record...");
            var updatedData = await client.ReadAsync("user:1");
            if (updatedData != null)
            {
                Console.WriteLine($"Retrieved data: {JsonSerializer.Serialize(updatedData, new JsonSerializerOptions { WriteIndented = true })}");
            }
            else
            {
                Console.WriteLine("Failed to retrieve data");
            }
            
            // Delete operation
            Console.WriteLine("\n5. Deleting the user record...");
            var deleteSuccess = await client.DeleteAsync("user:1");
            Console.WriteLine($"Delete operation {(deleteSuccess ? "succeeded" : "failed")}");
            
            // Verify deletion
            Console.WriteLine("\n6. Attempting to read deleted record...");
            var deletedData = await client.ReadAsync("user:1");
            if (deletedData == null)
            {
                Console.WriteLine("Record was successfully deleted");
            }
            else
            {
                Console.WriteLine("Record still exists");
            }
            
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"Error running example: {ex.Message}");
            Environment.Exit(1);
        }
    }
} 