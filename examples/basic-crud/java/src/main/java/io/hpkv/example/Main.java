package io.hpkv.example;

import io.github.cdimascio.dotenv.Dotenv;
import com.fasterxml.jackson.databind.ObjectMapper;

public class Main {
    public static void main(String[] args) {
        try {
            // Load environment variables from .env file
            Dotenv dotenv = Dotenv.load();
            
            // Initialize HPKV client using environment variables
            HPKVClient client = new HPKVClient(
                dotenv.get("HPKV_BASE_URL"),
                dotenv.get("HPKV_API_KEY")
            );

            System.out.println("HPKV CRUD Operations Example");
            System.out.println("===========================");

            // Create operation
            User userData = new User("John Doe", "john@example.com", 30);

            System.out.println("\n1. Creating a new user record...");
            boolean success = client.create("user:1", userData);
            if (!success) {
                System.out.println("Failed to create record. Exiting...");
                System.exit(1);
            }
            System.out.println("Create operation succeeded");

            // Read operation
            System.out.println("\n2. Reading the user record...");
            User retrievedData = client.read("user:1", User.class);
            if (retrievedData != null) {
                System.out.println("Retrieved data: " + new ObjectMapper()
                    .writerWithDefaultPrettyPrinter()
                    .writeValueAsString(retrievedData));
            } else {
                System.out.println("Failed to retrieve data");
            }

            // Update operation
            System.out.println("\n3. Updating the user's age...");
            userData.setAge(31);
            success = client.update("user:1", userData, false);
            System.out.println("Update operation " + (success ? "succeeded" : "failed"));

            // Read after update
            System.out.println("\n4. Reading the updated user record...");
            retrievedData = client.read("user:1", User.class);
            if (retrievedData != null) {
                System.out.println("Retrieved data: " + new ObjectMapper()
                    .writerWithDefaultPrettyPrinter()
                    .writeValueAsString(retrievedData));
            } else {
                System.out.println("Failed to retrieve data");
            }

            // Delete operation
            System.out.println("\n5. Deleting the user record...");
            success = client.delete("user:1");
            System.out.println("Delete operation " + (success ? "succeeded" : "failed"));

            // Verify deletion
            System.out.println("\n6. Attempting to read deleted record...");
            retrievedData = client.read("user:1", User.class);
            if (retrievedData == null) {
                System.out.println("Record was successfully deleted");
            } else {
                System.out.println("Record still exists");
            }

        } catch (Exception e) {
            System.err.println("Error running example: " + e.getMessage());
            System.exit(1);
        }
    }
} 