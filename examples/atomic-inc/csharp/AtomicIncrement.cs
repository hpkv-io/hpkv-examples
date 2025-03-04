using System;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using Microsoft.Extensions.Configuration;

class Program
{
    private static readonly string? HpkvBaseUrl;
    private static readonly string? HpkvApiKey;

    static Program()
    {
        var configuration = new ConfigurationBuilder()
            .SetBasePath(Directory.GetCurrentDirectory())
            .AddJsonFile("appsettings.json", optional: false, reloadOnChange: true)
            .AddEnvironmentVariables()
            .Build();

        HpkvBaseUrl = configuration["HPKV_BASE_URL"];
        HpkvApiKey = configuration["HPKV_API_KEY"];
    }

    private static async Task<bool> CreateKey(string key, int initialValue = 0)
    {
        try
        {
            var payload = new
            {
                key = key,
                value = initialValue.ToString()
            };

            Console.WriteLine($"Creating key with payload: {JsonSerializer.Serialize(payload)}");

            using var client = new HttpClient();
            client.DefaultRequestHeaders.Add("x-api-key", HpkvApiKey);

            var content = new StringContent(
                JsonSerializer.Serialize(payload),
                Encoding.UTF8,
                "application/json"
            );

            var response = await client.PostAsync($"{HpkvBaseUrl}/record", content);
            var responseContent = await response.Content.ReadAsStringAsync();

            Console.WriteLine($"Create response status: {(int)response.StatusCode}");
            Console.WriteLine($"Create response body: {responseContent}");

            if (!response.IsSuccessStatusCode)
            {
                Console.WriteLine($"Create failed with status {(int)response.StatusCode} - {responseContent}");
                return false;
            }

            return true;
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error creating record: {ex.Message}");
            return false;
        }
    }

    private static async Task<JsonDocument> AtomicIncrement(string key, int increment)
    {
        if (string.IsNullOrEmpty(HpkvBaseUrl) || string.IsNullOrEmpty(HpkvApiKey))
        {
            throw new ArgumentException("HPKV_BASE_URL and HPKV_API_KEY must be set in environment variables");
        }

        try
        {
            var payload = new
            {
                key = key,
                increment = increment
            };

            Console.WriteLine($"Attempting atomic increment with payload: {JsonSerializer.Serialize(payload)}");

            using var client = new HttpClient();
            client.DefaultRequestHeaders.Add("x-api-key", HpkvApiKey);

            var content = new StringContent(
                JsonSerializer.Serialize(payload),
                Encoding.UTF8,
                "application/json"
            );

            var response = await client.PostAsync($"{HpkvBaseUrl}/record/atomic", content);
            var responseContent = await response.Content.ReadAsStringAsync();

            Console.WriteLine($"Atomic increment response status: {(int)response.StatusCode}");
            Console.WriteLine($"Atomic increment response body: {responseContent}");

            if (response.StatusCode == System.Net.HttpStatusCode.NotFound)
            {
                Console.WriteLine($"Key '{key}' doesn't exist. Creating it with initial value 0...");
                if (!await CreateKey(key, 0))
                {
                    throw new Exception("Failed to create key with initial value");
                }

                Console.WriteLine("Retrying atomic increment after key creation...");
                response = await client.PostAsync($"{HpkvBaseUrl}/record/atomic", content);
                responseContent = await response.Content.ReadAsStringAsync();

                Console.WriteLine($"Retry response status: {(int)response.StatusCode}");
                Console.WriteLine($"Retry response body: {responseContent}");
            }

            response.EnsureSuccessStatusCode();
            var result = JsonDocument.Parse(responseContent);

            if (!result.RootElement.GetProperty("success").GetBoolean())
            {
                var errorMsg = result.RootElement.GetProperty("message").GetString() ?? "Unknown error";
                throw new Exception($"HPKV API error: {errorMsg}");
            }

            return result;
        }
        catch (HttpRequestException ex)
        {
            throw new Exception($"Failed to connect to HPKV API: {ex.Message}");
        }
        catch (Exception ex)
        {
            throw new Exception($"Invalid response from HPKV API: {ex.Message}");
        }
    }

    static async Task Main(string[] args)
    {
        var key = "counter:example";

        try
        {
            // Clean up any existing key
            Console.WriteLine("Cleaning up any existing key...");
            using var client = new HttpClient();
            client.DefaultRequestHeaders.Add("x-api-key", HpkvApiKey);

            var deleteResponse = await client.DeleteAsync($"{HpkvBaseUrl}/record/{key}");
            var deleteContent = await deleteResponse.Content.ReadAsStringAsync();
            Console.WriteLine($"Delete response status: {(int)deleteResponse.StatusCode}");
            Console.WriteLine($"Delete response body: {deleteContent}");

            await CreateKey(key, 0);

            var getResponse = await client.GetAsync($"{HpkvBaseUrl}/record/{key}");
            Console.WriteLine($"Get response status: {(int)getResponse.StatusCode}");

            // Increment by 1
            Console.WriteLine("\nIncrementing counter by 1...");
            var result = await AtomicIncrement(key, 1);
            Console.WriteLine($"Result: {result.RootElement.GetRawText()}");

            // Increment by 5
            Console.WriteLine("\nIncrementing counter by 5...");
            result = await AtomicIncrement(key, 5);
            Console.WriteLine($"Result: {result.RootElement.GetRawText()}");

            // Decrement by 2
            Console.WriteLine("\nDecrementing counter by 2...");
            result = await AtomicIncrement(key, -2);
            Console.WriteLine($"Result: {result.RootElement.GetRawText()}");
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error: {ex.Message}");
            Environment.Exit(1);
        }
    }
} 