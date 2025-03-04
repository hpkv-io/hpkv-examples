import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

public class HPKVWebSocketExample {
    public static void main(String[] args) {
        try {
            // Load environment variables (you might want to use a proper configuration library)
            String baseUrl = System.getenv("HPKV_BASE_URL");
            String apiKey = System.getenv("HPKV_API_KEY");
            
            if (baseUrl == null || apiKey == null) {
                System.err.println("Please set HPKV_BASE_URL and HPKV_API_KEY environment variables");
                System.exit(1);
            }
            
            // Initialize client
            HPKVWebSocketClient client = new HPKVWebSocketClient(baseUrl, apiKey);
            client.connect();
            
            // Wait for connection to be established
            Thread.sleep(1000);
            
            System.out.println("HPKV WebSocket CRUD Operations Example");
            System.out.println("=====================================");
            System.out.println("\nUsing HPKV WebSocket server: " + client.getURI());
            
            // Create operation
            Map<String, Object> userData = Map.of(
                "name", "John Doe",
                "email", "john@example.com",
                "age", 30
            );
            
            System.out.println("\n1. Creating a new user record...");
            CompletableFuture<Boolean> createFuture = client.create("user:1", userData);
            if (!createFuture.get()) {
                System.err.println("Failed to create record. Exiting...");
                System.exit(1);
            }
            System.out.println("Create operation succeeded");
            
            // Read operation
            System.out.println("\n2. Reading the user record...");
            CompletableFuture<Object> readFuture = client.read("user:1");
            Object retrievedData = readFuture.get();
            if (retrievedData != null) {
                System.out.println("Retrieved data: " + new ObjectMapper().writeValueAsString(retrievedData));
            } else {
                System.out.println("Failed to retrieve data");
            }
            
            // Update operation
            System.out.println("\n3. Updating the user's age...");
            userData = Map.of(
                "name", "John Doe",
                "email", "john@example.com",
                "age", 31
            );
            CompletableFuture<Boolean> updateFuture = client.update("user:1", userData, true);
            System.out.println("Update operation " + (updateFuture.get() ? "succeeded" : "failed"));
            
            // Read after update
            System.out.println("\n4. Reading the updated user record...");
            readFuture = client.read("user:1");
            retrievedData = readFuture.get();
            if (retrievedData != null) {
                System.out.println("Retrieved data: " + new ObjectMapper().writeValueAsString(retrievedData));
            } else {
                System.out.println("Failed to retrieve data");
            }
            
            // Delete operation
            System.out.println("\n5. Deleting the user record...");
            CompletableFuture<Boolean> deleteFuture = client.delete("user:1");
            System.out.println("Delete operation " + (deleteFuture.get() ? "succeeded" : "failed"));
            
            // Verify deletion
            System.out.println("\n6. Attempting to read deleted record...");
            readFuture = client.read("user:1");
            retrievedData = readFuture.get();
            if (retrievedData == null) {
                System.out.println("Record was successfully deleted");
            } else {
                System.out.println("Record still exists");
            }
            
            // Clean up
            client.close();
            
        } catch (Exception e) {
            System.err.println("Error running example: " + e.getMessage());
            System.exit(1);
        }
    }
} 