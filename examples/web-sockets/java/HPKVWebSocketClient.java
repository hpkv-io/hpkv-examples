import com.fasterxml.jackson.databind.ObjectMapper;
import org.java_websocket.client.WebSocketClient;
import org.java_websocket.handshake.ServerHandshake;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.URI;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class HPKVWebSocketClient extends WebSocketClient {
    private static final Logger logger = LoggerFactory.getLogger(HPKVWebSocketClient.class);
    private final ObjectMapper objectMapper;
    private final AtomicInteger messageId;
    private final Map<Integer, CompletableFuture<Map<String, Object>>> responseFutures;

    public enum OperationCode {
        GET(1),
        INSERT(2),
        UPDATE(3),
        DELETE(4);

        private final int value;

        OperationCode(int value) {
            this.value = value;
        }

        public int getValue() {
            return value;
        }
    }

    public HPKVWebSocketClient(String serverUri, String apiKey) throws Exception {
        super(new URI(serverUri + "/ws?apiKey=" + apiKey));
        this.objectMapper = new ObjectMapper();
        this.messageId = new AtomicInteger(0);
        this.responseFutures = new ConcurrentHashMap<>();
    }

    @Override
    public void onOpen(ServerHandshake handshakedata) {
        logger.info("WebSocket connection established");
    }

    @Override
    public void onMessage(String message) {
        try {
            Map<String, Object> response = objectMapper.readValue(message, Map.class);
            Integer msgId = (Integer) response.get("messageId");
            if (msgId != null && responseFutures.containsKey(msgId)) {
                CompletableFuture<Map<String, Object>> future = responseFutures.remove(msgId);
                if (response.containsKey("error")) {
                    future.completeExceptionally(new Exception((String) response.get("error")));
                } else {
                    future.complete(response);
                }
            }
        } catch (Exception e) {
            logger.error("Error handling message: {}", e.getMessage());
        }
    }

    @Override
    public void onClose(int code, String reason, boolean remote) {
        logger.info("WebSocket connection closed: {} - {}", code, reason);
    }

    @Override
    public void onError(Exception ex) {
        logger.error("WebSocket error: {}", ex.getMessage());
    }

    private CompletableFuture<Map<String, Object>> sendMessage(Map<String, Object> message) {
        try {
            int msgId = messageId.incrementAndGet();
            message.put("messageId", msgId);
            
            CompletableFuture<Map<String, Object>> future = new CompletableFuture<>();
            responseFutures.put(msgId, future);
            
            send(objectMapper.writeValueAsString(message));
            return future;
        } catch (Exception e) {
            CompletableFuture<Map<String, Object>> future = new CompletableFuture<>();
            future.completeExceptionally(e);
            return future;
        }
    }

    public CompletableFuture<Boolean> create(String key, Object value) {
        try {
            Map<String, Object> message = Map.of(
                "op", OperationCode.INSERT.getValue(),
                "key", key,
                "value", value instanceof String ? value : objectMapper.writeValueAsString(value)
            );
            
            return sendMessage(message)
                .thenApply(response -> !response.containsKey("error"));
        } catch (Exception e) {
            CompletableFuture<Boolean> future = new CompletableFuture<>();
            future.completeExceptionally(e);
            return future;
        }
    }

    public CompletableFuture<Object> read(String key) {
        Map<String, Object> message = Map.of(
            "op", OperationCode.GET.getValue(),
            "key", key
        );
        
        return sendMessage(message)
            .thenApply(response -> {
                if (response.containsKey("error")) {
                    return null;
                }
                try {
                    return objectMapper.readValue((String) response.get("value"), Object.class);
                } catch (Exception e) {
                    return response.get("value");
                }
            });
    }

    public CompletableFuture<Boolean> update(String key, Object value, boolean partialUpdate) {
        try {
            Map<String, Object> message = Map.of(
                "op", partialUpdate ? OperationCode.UPDATE.getValue() : OperationCode.INSERT.getValue(),
                "key", key,
                "value", value instanceof String ? value : objectMapper.writeValueAsString(value)
            );
            
            return sendMessage(message)
                .thenApply(response -> !response.containsKey("error"));
        } catch (Exception e) {
            CompletableFuture<Boolean> future = new CompletableFuture<>();
            future.completeExceptionally(e);
            return future;
        }
    }

    public CompletableFuture<Boolean> delete(String key) {
        Map<String, Object> message = Map.of(
            "op", OperationCode.DELETE.getValue(),
            "key", key
        );
        
        return sendMessage(message)
            .thenApply(response -> !response.containsKey("error"));
    }
} 