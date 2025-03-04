package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	City  string `json:"city"`
}

type Record struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HPKVRangeQueriesExample struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewHPKVRangeQueriesExample() (*HPKVRangeQueriesExample, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("HPKV_API_KEY")
	baseURL := os.Getenv("HPKV_API_BASE_URL")

	if apiKey == "" || baseURL == "" {
		return nil, fmt.Errorf("please set HPKV_API_KEY and HPKV_API_BASE_URL in your .env file")
	}

	return &HPKVRangeQueriesExample{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

func (e *HPKVRangeQueriesExample) getHeaders() map[string]string {
	return map[string]string{
		"x-api-key": e.apiKey,
	}
}

func (e *HPKVRangeQueriesExample) createSampleRecords() error {
	for i := 1; i <= 10; i++ {
		userData := UserData{
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   20 + i,
			City: func() string {
				if i%2 == 0 {
					return "New York"
				}
				return "San Francisco"
			}(),
		}

		value, err := json.Marshal(userData)
		if err != nil {
			return fmt.Errorf("error marshaling user data: %v", err)
		}

		record := Record{
			Key:   fmt.Sprintf("user:%d", i),
			Value: string(value),
		}

		payload, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("error marshaling record: %v", err)
		}

		req, err := http.NewRequest("POST", e.baseURL+"/record", bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}

		for key, value := range e.getHeaders() {
			req.Header.Set(key, value)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := e.client.Do(req)
		if err != nil {
			return fmt.Errorf("error creating record for user:%d: %v", i, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("error creating record for user:%d: status code %d, body: %s", i, resp.StatusCode, string(body))
		}

		fmt.Printf("Created record for user:%d\n", i)
	}
	return nil
}

func (e *HPKVRangeQueriesExample) performRangeQuery(startKey, endKey string, limit *int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("startKey", startKey)
	params.Set("endKey", endKey)
	if limit != nil {
		params.Set("limit", fmt.Sprintf("%d", *limit))
	}

	req, err := http.NewRequest("GET", e.baseURL+"/records?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range e.getHeaders() {
		req.Header.Set(key, value)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing range query: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error performing range query: status code %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

func (e *HPKVRangeQueriesExample) cleanupRecords() error {
	for i := 1; i <= 10; i++ {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/record/user:%d", e.baseURL, i), nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}

		for key, value := range e.getHeaders() {
			req.Header.Set(key, value)
		}

		resp, err := e.client.Do(req)
		if err != nil {
			return fmt.Errorf("error deleting record for user:%d: %v", i, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("error deleting record for user:%d: status code %d, body: %s", i, resp.StatusCode, string(body))
		}

		fmt.Printf("Deleted record for user:%d\n", i)
	}
	return nil
}

func main() {
	example, err := NewHPKVRangeQueriesExample()
	if err != nil {
		fmt.Printf("Error initializing example: %v\n", err)
		return
	}

	defer func() {
		fmt.Println("\nCleaning up sample records...")
		// if err := example.cleanupRecords(); err != nil {
		//     fmt.Printf("Error cleaning up records: %v\n", err)
		// }
	}()

	// Create sample records
	fmt.Println("\nCreating sample records...")
	if err := example.createSampleRecords(); err != nil {
		fmt.Printf("Error creating sample records: %v\n", err)
		return
	}

	// Example 1: Basic range query
	fmt.Println("\nExample 1: Basic range query (users 1-5)")
	result1, err := example.performRangeQuery("user:1", "user:5", nil)
	if err != nil {
		fmt.Printf("Error performing range query: %v\n", err)
		return
	}
	prettyJSON, _ := json.MarshalIndent(result1, "", "  ")
	fmt.Println(string(prettyJSON))

	// Example 2: Range query with limit
	fmt.Println("\nExample 2: Range query with limit (users 1-10, limit 3)")
	limit := 3
	result2, err := example.performRangeQuery("user:1", "user:9", &limit)
	if err != nil {
		fmt.Printf("Error performing range query: %v\n", err)
		return
	}
	prettyJSON, _ = json.MarshalIndent(result2, "", "  ")
	fmt.Println(string(prettyJSON))

	// Example 3: Range query for specific city
	fmt.Println("\nExample 3: Range query for users in New York (even IDs)")
	result3, err := example.performRangeQuery("user:2", "user:9", nil)
	if err != nil {
		fmt.Printf("Error performing range query: %v\n", err)
		return
	}

	records, ok := result3["records"].([]interface{})
	if !ok {
		fmt.Println("Error: records not found in response")
		return
	}

	var newYorkUsers []interface{}
	for _, record := range records {
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		value, ok := recordMap["value"].(string)
		if !ok {
			continue
		}

		var userData UserData
		if err := json.Unmarshal([]byte(value), &userData); err != nil {
			continue
		}

		if userData.City == "New York" {
			newYorkUsers = append(newYorkUsers, record)
		}
	}

	prettyJSON, _ = json.MarshalIndent(newYorkUsers, "", "  ")
	fmt.Println(string(prettyJSON))
}
