using System;
using System.Net.Http;
using System.Net.Http.Json;
using System.Text.Json;
using System.Threading.Tasks;

namespace HPKVExample
{
    public class HPKVClient
    {
        private readonly HttpClient _httpClient;
        private readonly string _baseUrl;
        private readonly string _apiKey;

        public HPKVClient(string? baseUrl = null, string? apiKey = null)
        {
            // Load from environment variables if not provided
            var envBaseUrl = baseUrl ?? Environment.GetEnvironmentVariable("HPKV_BASE_URL");
            var envApiKey = apiKey ?? Environment.GetEnvironmentVariable("HPKV_API_KEY");

            if (string.IsNullOrEmpty(envBaseUrl))
                throw new ArgumentException("HPKV base URL not provided. Set HPKV_BASE_URL environment variable or pass baseUrl parameter.");
            if (string.IsNullOrEmpty(envApiKey))
                throw new ArgumentException("HPKV API key not provided. Set HPKV_API_KEY environment variable or pass apiKey parameter.");

            _baseUrl = envBaseUrl.TrimEnd('/');
            _apiKey = envApiKey;
            _httpClient = new HttpClient();
            _httpClient.DefaultRequestHeaders.Add("x-api-key", _apiKey);
            _httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        }

        private string SerializeValue<T>(T value)
        {
            if (value is string strValue)
                return strValue;
            return JsonSerializer.Serialize(value);
        }

        public async Task<bool> CreateAsync<T>(string key, T value)
        {
            try
            {
                var payload = new
                {
                    key = key,
                    value = SerializeValue(value)
                };

                Console.WriteLine($"Sending request to {_baseUrl}/record");
                Console.WriteLine($"Payload: {JsonSerializer.Serialize(payload, new JsonSerializerOptions { WriteIndented = true })}");

                var response = await _httpClient.PostAsJsonAsync($"{_baseUrl}/record", payload);
                
                if (!response.IsSuccessStatusCode)
                {
                    var error = await response.Content.ReadAsStringAsync();
                    Console.Error.WriteLine($"Create failed with status {response.StatusCode} - {error}");
                    return false;
                }

                Console.WriteLine($"Create succeeded with status {response.StatusCode}");
                return true;
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error creating record: {ex.Message}");
                return false;
            }
        }

        public async Task<T?> ReadAsync<T>(string key)
        {
            try
            {
                var response = await _httpClient.GetAsync($"{_baseUrl}/record/{key}");
                if (!response.IsSuccessStatusCode)
                {
                    var error = await response.Content.ReadAsStringAsync();
                    Console.Error.WriteLine($"Error in read: Status {response.StatusCode} - {error}");
                    return default;
                }

                var data = await response.Content.ReadFromJsonAsync<RecordResponse>();
                if (data?.Value == null)
                    return default;

                try
                {
                    return JsonSerializer.Deserialize<T>(data.Value);
                }
                catch
                {
                    return (T)Convert.ChangeType(data.Value, typeof(T));
                }
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error reading record: {ex.Message}");
                return default;
            }
        }

        public async Task<bool> UpdateAsync<T>(string key, T value, bool partialUpdate = false)
        {
            try
            {
                var payload = new
                {
                    key = key,
                    value = SerializeValue(value),
                    partialUpdate = partialUpdate
                };

                var response = await _httpClient.PostAsJsonAsync($"{_baseUrl}/record", payload);
                if (!response.IsSuccessStatusCode)
                {
                    var error = await response.Content.ReadAsStringAsync();
                    Console.Error.WriteLine($"Update failed with status {response.StatusCode} - {error}");
                    return false;
                }

                return true;
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error updating record: {ex.Message}");
                return false;
            }
        }

        public async Task<bool> DeleteAsync(string key)
        {
            try
            {
                var response = await _httpClient.DeleteAsync($"{_baseUrl}/record/{key}");
                if (!response.IsSuccessStatusCode)
                {
                    var error = await response.Content.ReadAsStringAsync();
                    Console.Error.WriteLine($"Delete failed with status {response.StatusCode} - {error}");
                    return false;
                }

                return true;
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error deleting record: {ex.Message}");
                return false;
            }
        }

        private class RecordResponse
        {
            public string Value { get; set; } = string.Empty;
        }
    }
} 