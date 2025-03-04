using System;
using System.Collections.Concurrent;
using System.Net.WebSockets;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;

namespace HPKV.Examples.WebSocket
{
    public enum OperationCode
    {
        Get = 1,
        Insert = 2,
        Update = 3,
        Delete = 4
    }

    public class HPKVWebSocketClient : IDisposable
    {
        private readonly ClientWebSocket _webSocket;
        private readonly string _url;
        private ulong _messageId;
        private readonly Dictionary<ulong, TaskCompletionSource<JsonElement>> _responseFutures;
        private readonly CancellationTokenSource _cts;
        private Task? _receiveTask;

        public HPKVWebSocketClient(string baseUrl, string apiKey)
        {
            baseUrl = baseUrl.Replace("http://", "ws://").Replace("https://", "wss://");
            if (!baseUrl.StartsWith("ws://") && !baseUrl.StartsWith("wss://"))
            {
                baseUrl = "wss://" + baseUrl;
            }
            _url = $"{baseUrl}/ws?apiKey={apiKey}";
            _webSocket = new ClientWebSocket();
            _messageId = 0;
            _responseFutures = new Dictionary<ulong, TaskCompletionSource<JsonElement>>();
            _cts = new CancellationTokenSource();
        }

        public async Task ConnectAsync()
        {
            await _webSocket.ConnectAsync(new Uri(_url), _cts.Token);
            _receiveTask = ReceiveLoop();
            Console.WriteLine("WebSocket connection established");
        }

        private async Task ReceiveLoop()
        {
            var buffer = new byte[4096];
            try
            {
                while (!_cts.Token.IsCancellationRequested)
                {
                    var result = await _webSocket.ReceiveAsync(new ArraySegment<byte>(buffer), _cts.Token);
                    if (result.MessageType == WebSocketMessageType.Close)
                    {
                        await _webSocket.CloseAsync(WebSocketCloseStatus.NormalClosure, string.Empty, _cts.Token);
                        break;
                    }

                    var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
                    var response = JsonSerializer.Deserialize<JsonElement>(message);
                    
                    if (response.TryGetProperty("messageId", out var messageIdElement))
                    {
                        var messageId = messageIdElement.GetUInt64();
                        if (_responseFutures.TryGetValue(messageId, out var tcs))
                        {
                            _responseFutures.Remove(messageId);
                            tcs.SetResult(response);
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error in receive loop: {ex.Message}");
            }
        }

        private async Task<JsonElement> SendMessageAsync(object message)
        {
            var messageId = ++_messageId;
            var messageWithId = JsonSerializer.SerializeToElement(new
            {
                messageId,
                op = message.GetType().GetProperty("op")?.GetValue(message),
                key = message.GetType().GetProperty("key")?.GetValue(message),
                value = message.GetType().GetProperty("value")?.GetValue(message)
            });

            var tcs = new TaskCompletionSource<JsonElement>();
            _responseFutures[messageId] = tcs;

            var messageBytes = Encoding.UTF8.GetBytes(messageWithId.ToString());
            await _webSocket.SendAsync(new ArraySegment<byte>(messageBytes), WebSocketMessageType.Text, true, _cts.Token);

            var response = await tcs.Task;
            if (response.TryGetProperty("error", out var error) && !error.ValueEquals(""))
            {
                throw new Exception(error.GetString());
            }

            return response;
        }

        public async Task<bool> CreateAsync(string key, object value)
        {
            try
            {
                var message = new { op = OperationCode.Insert, key, value = JsonSerializer.Serialize(value) };
                await SendMessageAsync(message);
                return true;
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error creating record: {ex.Message}");
                return false;
            }
        }

        public async Task<JsonElement?> ReadAsync(string key)
        {
            try
            {
                var message = new { op = OperationCode.Get, key };
                var response = await SendMessageAsync(message);
                
                if (response.TryGetProperty("value", out var value))
                {
                    if (value.ValueKind == JsonValueKind.String)
                    {
                        return JsonSerializer.Deserialize<JsonElement>(value.GetString()!);
                    }
                    return value;
                }
                return null;
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error reading record: {ex.Message}");
                return null;
            }
        }

        public async Task<bool> UpdateAsync(string key, object value, bool partialUpdate)
        {
            try
            {
                // For partial updates, first read the existing value
                if (partialUpdate)
                {
                    var existingData = await ReadAsync(key);
                    if (existingData != null)
                    {
                        var newData = JsonSerializer.Deserialize<JsonElement>(JsonSerializer.Serialize(value));
                        var existingDict = JsonSerializer.Deserialize<Dictionary<string, JsonElement>>(existingData.Value.GetRawText());
                        var newDict = JsonSerializer.Deserialize<Dictionary<string, JsonElement>>(newData.GetRawText());

                        foreach (var (k, v) in newDict)
                        {
                            existingDict[k] = v;
                        }

                        value = existingDict;
                    }
                }

                var message = new { op = OperationCode.Insert, key, value = JsonSerializer.Serialize(value) };
                await SendMessageAsync(message);
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
                var message = new { op = OperationCode.Delete, key };
                await SendMessageAsync(message);
                return true;
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error deleting record: {ex.Message}");
                return false;
            }
        }

        public void Dispose()
        {
            _cts.Cancel();
            if (_webSocket.State == WebSocketState.Open)
            {
                _webSocket.CloseAsync(WebSocketCloseStatus.NormalClosure, string.Empty, CancellationToken.None).Wait();
            }
            _webSocket.Dispose();
            _cts.Dispose();
        }
    }
} 