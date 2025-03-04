package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	hpkvBaseURL string
	hpkvAPIKey  string
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	hpkvBaseURL = os.Getenv("HPKV_BASE_URL")
	hpkvAPIKey = os.Getenv("HPKV_API_KEY")
}

func createKey(key string, initialValue int) (bool, error) {
	payload := map[string]string{
		"key":   key,
		"value": fmt.Sprintf("%d", initialValue),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("error marshaling payload: %v", err)
	}

	fmt.Printf("Creating key with payload: %s\n", string(jsonData))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/record", hpkvBaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", hpkvAPIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Printf("Create response status: %d\n", resp.StatusCode)
	fmt.Printf("Create response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("create failed with status %d - %s", resp.StatusCode, string(body))
	}

	return true, nil
}

func atomicIncrement(key string, increment int) (map[string]interface{}, error) {
	if hpkvBaseURL == "" || hpkvAPIKey == "" {
		return nil, fmt.Errorf("HPKV_BASE_URL and HPKV_API_KEY must be set in environment variables")
	}

	payload := map[string]interface{}{
		"key":       key,
		"increment": increment,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	fmt.Printf("Attempting atomic increment with payload: %s\n", string(jsonData))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/record/atomic", hpkvBaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", hpkvAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Printf("Atomic increment response status: %d\n", resp.StatusCode)
	fmt.Printf("Atomic increment response body: %s\n", string(body))

	// If the key doesn't exist (404), create it first and try again
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Key '%s' doesn't exist. Creating it with initial value 0...\n", key)
		if success, err := createKey(key, 0); !success {
			return nil, fmt.Errorf("failed to create key with initial value: %v", err)
		}

		// Retry the increment operation
		fmt.Println("Retrying atomic increment after key creation...")
		req, err = http.NewRequest("POST", fmt.Sprintf("%s/record/atomic", hpkvBaseURL), bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("error creating retry request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", hpkvAPIKey)

		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making retry request: %v", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading retry response body: %v", err)
		}

		fmt.Printf("Retry response status: %d\n", resp.StatusCode)
		fmt.Printf("Retry response body: %s\n", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		message := "Unknown error"
		if msg, ok := result["message"].(string); ok {
			message = msg
		}
		return nil, fmt.Errorf("HPKV API error: %s", message)
	}

	return result, nil
}

func main() {
	key := "counter:example"

	// Clean up any existing key
	fmt.Println("Cleaning up any existing key...")
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/record/%s", hpkvBaseURL, key), nil)
	if err != nil {
		fmt.Printf("Error creating delete request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", hpkvAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making delete request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading delete response body: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Delete response status: %d\n", resp.StatusCode)
	fmt.Printf("Delete response body: %s\n", string(body))

	if _, err := createKey(key, 0); err != nil {
		fmt.Printf("Error creating key: %v\n", err)
		os.Exit(1)
	}

	// Get the initial value
	getReq, err := http.NewRequest("GET", fmt.Sprintf("%s/record/%s", hpkvBaseURL, key), nil)
	if err != nil {
		fmt.Printf("Error creating get request: %v\n", err)
		os.Exit(1)
	}

	getReq.Header.Set("Content-Type", "application/json")
	getReq.Header.Set("x-api-key", hpkvAPIKey)

	getResp, err := client.Do(getReq)
	if err != nil {
		fmt.Printf("Error making get request: %v\n", err)
		os.Exit(1)
	}
	defer getResp.Body.Close()

	fmt.Printf("Get response status: %d\n", getResp.StatusCode)

	// Increment by 1
	fmt.Println("\nIncrementing counter by 1...")
	result, err := atomicIncrement(key, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result: %+v\n", result)

	// Increment by 5
	fmt.Println("\nIncrementing counter by 5...")
	result, err = atomicIncrement(key, 5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result: %+v\n", result)

	// Decrement by 2
	fmt.Println("\nDecrementing counter by 2...")
	result, err = atomicIncrement(key, -2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result: %+v\n", result)
}
