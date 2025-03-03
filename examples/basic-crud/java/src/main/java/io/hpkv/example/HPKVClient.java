package io.hpkv.example;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import okhttp3.*;
import java.io.IOException;

public class HPKVClient {
    private final OkHttpClient httpClient;
    private final ObjectMapper objectMapper;
    private final String baseUrl;
    private final String apiKey;

    public HPKVClient(String baseUrl, String apiKey) {
        if (baseUrl == null || baseUrl.trim().isEmpty()) {
            throw new IllegalArgumentException("HPKV base URL not provided. Set HPKV_BASE_URL environment variable or pass baseUrl parameter.");
        }
        if (apiKey == null || apiKey.trim().isEmpty()) {
            throw new IllegalArgumentException("HPKV API key not provided. Set HPKV_API_KEY environment variable or pass apiKey parameter.");
        }

        this.baseUrl = baseUrl.replaceAll("/$", "");
        this.apiKey = apiKey;
        this.httpClient = new OkHttpClient();
        this.objectMapper = new ObjectMapper();
    }

    private String serializeValue(Object value) {
        if (value instanceof String) {
            return (String) value;
        }
        try {
            return objectMapper.writeValueAsString(value);
        } catch (Exception e) {
            throw new RuntimeException("Failed to serialize value", e);
        }
    }

    public boolean create(String key, Object value) {
        try {
            String payload = objectMapper.writeValueAsString(new Record(key, serializeValue(value)));
            
            System.out.println("Sending request to " + baseUrl + "/record");
            System.out.println("Payload: " + objectMapper.writerWithDefaultPrettyPrinter().writeValueAsString(
                objectMapper.readTree(payload)));

            Request request = new Request.Builder()
                .url(baseUrl + "/record")
                .addHeader("Content-Type", "application/json")
                .addHeader("x-api-key", apiKey)
                .post(RequestBody.create(payload, MediaType.parse("application/json")))
                .build();

            try (Response response = httpClient.newCall(request).execute()) {
                if (!response.isSuccessful()) {
                    String error = response.body() != null ? response.body().string() : "";
                    System.err.println("Create failed with status " + response.code() + " - " + error);
                    return false;
                }
                
                System.out.println("Create succeeded with status " + response.code());
                return true;
            }
        } catch (Exception e) {
            System.err.println("Error creating record: " + e.getMessage());
            return false;
        }
    }

    public <T> T read(String key, Class<T> valueType) {
        try {
            Request request = new Request.Builder()
                .url(baseUrl + "/record/" + key)
                .addHeader("Content-Type", "application/json")
                .addHeader("x-api-key", apiKey)
                .build();

            try (Response response = httpClient.newCall(request).execute()) {
                if (!response.isSuccessful()) {
                    String error = response.body() != null ? response.body().string() : "";
                    System.err.println("Error in read: Status " + response.code() + " - " + error);
                    return null;
                }

                String responseBody = response.body().string();
                JsonNode data = objectMapper.readTree(responseBody);
                if (!data.has("value")) {
                    return null;
                }

                String value = data.get("value").asText();
                try {
                    return objectMapper.readValue(value, valueType);
                } catch (Exception e) {
                    return (T) value;
                }
            }
        } catch (Exception e) {
            System.err.println("Error reading record: " + e.getMessage());
            return null;
        }
    }

    public boolean update(String key, Object value, boolean partialUpdate) {
        try {
            String payload = objectMapper.writeValueAsString(
                new Record(key, serializeValue(value), partialUpdate));

            Request request = new Request.Builder()
                .url(baseUrl + "/record")
                .addHeader("Content-Type", "application/json")
                .addHeader("x-api-key", apiKey)
                .post(RequestBody.create(payload, MediaType.parse("application/json")))
                .build();

            try (Response response = httpClient.newCall(request).execute()) {
                if (!response.isSuccessful()) {
                    String error = response.body() != null ? response.body().string() : "";
                    System.err.println("Update failed with status " + response.code() + " - " + error);
                    return false;
                }
                return true;
            }
        } catch (Exception e) {
            System.err.println("Error updating record: " + e.getMessage());
            return false;
        }
    }

    public boolean delete(String key) {
        try {
            Request request = new Request.Builder()
                .url(baseUrl + "/record/" + key)
                .addHeader("Content-Type", "application/json")
                .addHeader("x-api-key", apiKey)
                .delete()
                .build();

            try (Response response = httpClient.newCall(request).execute()) {
                if (!response.isSuccessful()) {
                    String error = response.body() != null ? response.body().string() : "";
                    System.err.println("Delete failed with status " + response.code() + " - " + error);
                    return false;
                }
                return true;
            }
        } catch (Exception e) {
            System.err.println("Error deleting record: " + e.getMessage());
            return false;
        }
    }

    private static class Record {
        public String key;
        public String value;
        public Boolean partialUpdate;

        public Record(String key, String value) {
            this.key = key;
            this.value = value;
        }

        public Record(String key, String value, Boolean partialUpdate) {
            this.key = key;
            this.value = value;
            this.partialUpdate = partialUpdate;
        }
    }
} 