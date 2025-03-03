package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize HPKV client using environment variables
	client := NewHPKVClient(
		os.Getenv("HPKV_BASE_URL"),
		os.Getenv("HPKV_API_KEY"),
	)

	fmt.Println("HPKV CRUD Operations Example")
	fmt.Println("===========================")

	// Create operation
	userData := User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	fmt.Println("\n1. Creating a new user record...")
	if err := client.Create("user:1", userData); err != nil {
		fmt.Println("Failed to create record. Exiting...")
		os.Exit(1)
	}
	fmt.Println("Create operation succeeded")

	// Read operation
	fmt.Println("\n2. Reading the user record...")
	var retrievedData User
	if err := client.Read("user:1", &retrievedData); err != nil {
		fmt.Println("Failed to retrieve data:", err)
	} else {
		data, _ := json.MarshalIndent(retrievedData, "", "  ")
		fmt.Printf("Retrieved data: %s\n", string(data))
	}

	// Update operation
	fmt.Println("\n3. Updating the user's age...")
	userData.Age = 31
	if err := client.Update("user:1", userData, false); err != nil {
		fmt.Printf("Update operation failed: %v\n", err)
	} else {
		fmt.Println("Update operation succeeded")
	}

	// Read after update
	fmt.Println("\n4. Reading the updated user record...")
	if err := client.Read("user:1", &retrievedData); err != nil {
		fmt.Println("Failed to retrieve data:", err)
	} else {
		data, _ := json.MarshalIndent(retrievedData, "", "  ")
		fmt.Printf("Retrieved data: %s\n", string(data))
	}

	// Delete operation
	fmt.Println("\n5. Deleting the user record...")
	if err := client.Delete("user:1"); err != nil {
		fmt.Printf("Delete operation failed: %v\n", err)
	} else {
		fmt.Println("Delete operation succeeded")
	}

	// Verify deletion
	fmt.Println("\n6. Attempting to read deleted record...")
	if err := client.Read("user:1", &retrievedData); err != nil {
		fmt.Println("Record was successfully deleted")
	} else {
		fmt.Println("Record still exists")
	}
} 