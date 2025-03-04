using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using System.IO;

namespace HPKVExamples
{
    public class HPKVRangeQueriesExample
    {
        private readonly string _apiKey;
        private readonly string _baseUrl;
        private readonly HttpClient _httpClient;

        public HPKVRangeQueriesExample()
        {
            LoadEnvFile();
            _apiKey = Environment.GetEnvironmentVariable("HPKV_API_KEY") ?? throw new ArgumentException("HPKV_API_KEY not found in environment variables");
            _baseUrl = Environment.GetEnvironmentVariable("HPKV_API_BASE_URL") ?? throw new ArgumentException("HPKV_API_BASE_URL not found in environment variables");

            _httpClient = new HttpClient();
            _httpClient.DefaultRequestHeaders.Add("x-api-key", _apiKey);
        }

        private void LoadEnvFile()
        {
            var envPath = Path.Combine(Directory.GetCurrentDirectory(), ".env");
            if (File.Exists(envPath))
            {
                foreach (var line in File.ReadAllLines(envPath))
                {
                    var trimmedLine = line.Trim();
                    if (!string.IsNullOrEmpty(trimmedLine) && !trimmedLine.StartsWith("#"))
                    {
                        var parts = trimmedLine.Split('=', 2);
                        if (parts.Length == 2)
                        {
                            Environment.SetEnvironmentVariable(parts[0].Trim(), parts[1].Trim());
                        }
                    }
                }
            }
        }

        public async Task CreateSampleRecordsAsync()
        {
            for (int i = 1; i <= 10; i++)
            {
                var userData = new
                {
                    name = $"User {i}",
                    email = $"user{i}@example.com",
                    age = 20 + i,
                    city = i % 2 == 0 ? "New York" : "San Francisco"
                };

                var payload = new
                {
                    key = $"user:{i}",
                    value = JsonSerializer.Serialize(userData)
                };

                var content = new StringContent(
                    JsonSerializer.Serialize(payload),
                    Encoding.UTF8,
                    "application/json"
                );

                try
                {
                    var response = await _httpClient.PostAsync($"{_baseUrl}/record", content);
                    if (!response.IsSuccessStatusCode)
                    {
                        var errorBody = await response.Content.ReadAsStringAsync();
                        throw new HttpRequestException($"Error creating record for user:{i}. Status: {response.StatusCode}, Body: {errorBody}");
                    }
                    Console.WriteLine($"Created record for user:{i}");
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error creating record for user:{i}: {ex.Message}");
                    throw;
                }
            }
        }

        public async Task<Dictionary<string, object>> PerformRangeQueryAsync(string startKey, string endKey, int? limit = null)
        {
            var queryParams = new List<string>
            {
                $"startKey={Uri.EscapeDataString(startKey)}",
                $"endKey={Uri.EscapeDataString(endKey)}"
            };

            if (limit.HasValue)
            {
                queryParams.Add($"limit={limit.Value}");
            }

            try
            {
                var response = await _httpClient.GetAsync($"{_baseUrl}/records?{string.Join("&", queryParams)}");
                if (!response.IsSuccessStatusCode)
                {
                    var errorBody = await response.Content.ReadAsStringAsync();
                    throw new HttpRequestException($"Error performing range query. Status: {response.StatusCode}, Body: {errorBody}");
                }
                var content = await response.Content.ReadAsStringAsync();
                return JsonSerializer.Deserialize<Dictionary<string, object>>(content) ?? throw new JsonException("Failed to deserialize response");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error performing range query: {ex.Message}");
                throw;
            }
        }

        public async Task CleanupRecordsAsync()
        {
            for (int i = 1; i <= 10; i++)
            {
                try
                {
                    var response = await _httpClient.DeleteAsync($"{_baseUrl}/record/user:{i}");
                    if (!response.IsSuccessStatusCode)
                    {
                        var errorBody = await response.Content.ReadAsStringAsync();
                        throw new HttpRequestException($"Error deleting record for user:{i}. Status: {response.StatusCode}, Body: {errorBody}");
                    }
                    Console.WriteLine($"Deleted record for user:{i}");
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error deleting record for user:{i}: {ex.Message}");
                    throw;
                }
            }
        }
    }

    class Program
    {
        static async Task Main(string[] args)
        {
            try
            {
                var example = new HPKVRangeQueriesExample();

                // Create sample records
                Console.WriteLine("\nCreating sample records...");
                await example.CreateSampleRecordsAsync();

                // Example 1: Basic range query
                Console.WriteLine("\nExample 1: Basic range query (users 1-5)");
                var result1 = await example.PerformRangeQueryAsync("user:1", "user:5");
                Console.WriteLine(JsonSerializer.Serialize(result1, new JsonSerializerOptions { WriteIndented = true }));

                // Example 2: Range query with limit
                Console.WriteLine("\nExample 2: Range query with limit (users 1-10, limit 3)");
                var result2 = await example.PerformRangeQueryAsync("user:1", "user:9", 3);
                Console.WriteLine(JsonSerializer.Serialize(result2, new JsonSerializerOptions { WriteIndented = true }));

                // Example 3: Range query for specific city
                Console.WriteLine("\nExample 3: Range query for users in New York (even IDs)");
                var result3 = await example.PerformRangeQueryAsync("user:2", "user:9");
                var recordsJson = result3["records"]?.ToString() ?? throw new KeyNotFoundException("records not found in response");
                var records = JsonSerializer.Deserialize<List<Dictionary<string, object>>>(recordsJson) ?? new List<Dictionary<string, object>>();
                var newYorkUsers = records.Where(record =>
                {
                    var value = record["value"]?.ToString() ?? throw new KeyNotFoundException("value not found in record");
                    var userData = JsonSerializer.Deserialize<Dictionary<string, object>>(value) ?? throw new JsonException("Failed to deserialize user data");
                    return userData["city"]?.ToString() == "New York";
                }).ToList();
                Console.WriteLine(JsonSerializer.Serialize(newYorkUsers, new JsonSerializerOptions { WriteIndented = true }));
            }
            catch (Exception ex)
            {
                Console.WriteLine($"\nError running example: {ex.Message}");
                if (ex.InnerException != null)
                {
                    Console.WriteLine($"Inner exception: {ex.InnerException.Message}");
                }
                return;
            }
            finally
            {
                // Clean up the sample records
                Console.WriteLine("\nCleaning up sample records...");
                // await example.CleanupRecordsAsync();
            }
        }
    }
} 