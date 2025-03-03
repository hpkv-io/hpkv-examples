using System;
using System.Text.Json;
using System.Threading.Tasks;
using DotNetEnv;

namespace HPKVExample
{
    public class User
    {
        public string Name { get; set; } = string.Empty;
        public string Email { get; set; } = string.Empty;
        public int Age { get; set; }
    }

    class Program
    {
        static async Task Main(string[] args)
        {
            try
            {
                // Load environment variables from .env file
                Env.Load();

                // Initialize HPKV client using environment variables
                var client = new HPKVClient();

                Console.WriteLine("HPKV CRUD Operations Example");
                Console.WriteLine("===========================");

                // Create operation
                var userData = new User
                {
                    Name = "John Doe",
                    Email = "john@example.com",
                    Age = 30
                };

                Console.WriteLine("\n1. Creating a new user record...");
                var success = await client.CreateAsync("user:1", userData);
                if (!success)
                {
                    Console.WriteLine("Failed to create record. Exiting...");
                    Environment.Exit(1);
                }
                Console.WriteLine("Create operation succeeded");

                // Read operation
                Console.WriteLine("\n2. Reading the user record...");
                var retrievedData = await client.ReadAsync<User>("user:1");
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
                var updatedData = new User
                {
                    Name = "John Doe",
                    Email = "john@example.com",
                    Age = 31
                };
                success = await client.UpdateAsync("user:1", updatedData);
                Console.WriteLine($"Update operation {(success ? "succeeded" : "failed")}");

                // Read after update
                Console.WriteLine("\n4. Reading the updated user record...");
                retrievedData = await client.ReadAsync<User>("user:1");
                if (retrievedData != null)
                {
                    Console.WriteLine($"Retrieved data: {JsonSerializer.Serialize(retrievedData, new JsonSerializerOptions { WriteIndented = true })}");
                }
                else
                {
                    Console.WriteLine("Failed to retrieve data");
                }

                // Delete operation
                Console.WriteLine("\n5. Deleting the user record...");
                success = await client.DeleteAsync("user:1");
                Console.WriteLine($"Delete operation {(success ? "succeeded" : "failed")}");

                // Verify deletion
                Console.WriteLine("\n6. Attempting to read deleted record...");
                retrievedData = await client.ReadAsync<User>("user:1");
                if (retrievedData == null)
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
} 