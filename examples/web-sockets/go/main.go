package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	// Load environment variables
	baseURL := os.Getenv("HPKV_BASE_URL")
	apiKey := os.Getenv("HPKV_API_KEY")

	if baseURL == "" || apiKey == "" {
		log.Fatal("Please set HPKV_BASE_URL and HPKV_API_KEY environment variables")
	}

	// Initialize client
	client, err := NewHPKVWebSocketClient(baseURL, apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("HPKV WebSocket CRUD Operations Example")
	fmt.Println("=====================================")
	fmt.Printf("\nUsing HPKV WebSocket server: %s/ws?apiKey=%s\n", baseURL, apiKey)

	// Create operation
	userData := UserData{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	fmt.Println("\n1. Creating a new user record...")
	err = client.Create("user:1", userData)
	if err != nil {
		log.Fatalf("Failed to create record: %v", err)
	}
	fmt.Println("Create operation succeeded")

	// Read operation
	fmt.Println("\n2. Reading the user record...")
	retrievedData, err := client.Read("user:1")
	if err != nil {
		log.Printf("Failed to retrieve data: %v", err)
	} else {
		jsonData, _ := json.MarshalIndent(retrievedData, "", "  ")
		fmt.Printf("Retrieved data: %s\n", jsonData)
	}

	// Update operation
	fmt.Println("\n3. Updating the user's age...")
	userData.Age = 31
	err = client.Update("user:1", userData, true)
	if err != nil {
		log.Printf("Update operation failed: %v", err)
	} else {
		fmt.Println("Update operation succeeded")
	}

	// Read after update
	fmt.Println("\n4. Reading the updated user record...")
	updatedData, err := client.Read("user:1")
	if err != nil {
		log.Printf("Failed to retrieve data: %v", err)
	} else {
		jsonData, _ := json.MarshalIndent(updatedData, "", "  ")
		fmt.Printf("Retrieved data: %s\n", jsonData)
	}

	// Delete operation
	fmt.Println("\n5. Deleting the user record...")
	err = client.Delete("user:1")
	if err != nil {
		log.Printf("Delete operation failed: %v", err)
	} else {
		fmt.Println("Delete operation succeeded")
	}

	// Verify deletion
	fmt.Println("\n6. Attempting to read deleted record...")
	_, err = client.Read("user:1")
	if err != nil {
		fmt.Println("Record was successfully deleted")
	} else {
		fmt.Println("Record still exists")
	}
}
