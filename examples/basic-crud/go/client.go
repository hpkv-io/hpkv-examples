package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HPKVClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type Record struct {
	Key          string      `json:"key"`
	Value        string      `json:"value"`
	PartialUpdate *bool      `json:"partialUpdate,omitempty"`
}

type Response struct {
	Value string `json:"value"`
}

func NewHPKVClient(baseURL, apiKey string) *HPKVClient {
	if baseURL == "" {
		panic("HPKV base URL not provided. Set HPKV_BASE_URL environment variable")
	}
	if apiKey == "" {
		panic("HPKV API key not provided. Set HPKV_API_KEY environment variable")
	}

	return &HPKVClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (c *HPKVClient) serializeValue(value interface{}) (string, error) {
	if str, ok := value.(string); ok {
		return str, nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to serialize value: %w", err)
	}
	return string(data), nil
}

func (c *HPKVClient) Create(key string, value interface{}) error {
	serializedValue, err := c.serializeValue(value)
	if err != nil {
		return err
	}

	record := Record{
		Key:   key,
		Value: serializedValue,
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	fmt.Printf("Sending request to %s/record\n", c.baseURL)
	fmt.Printf("Payload: %s\n", prettyJSON(payload))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/record", c.baseURL), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Create succeeded with status %d\n", resp.StatusCode)
	return nil
}

func (c *HPKVClient) Read(key string, value interface{}) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/record/%s", c.baseURL, key), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("read failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if err := json.Unmarshal([]byte(response.Value), value); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

func (c *HPKVClient) Update(key string, value interface{}, partialUpdate bool) error {
	serializedValue, err := c.serializeValue(value)
	if err != nil {
		return err
	}

	record := Record{
		Key:          key,
		Value:        serializedValue,
		PartialUpdate: &partialUpdate,
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/record", c.baseURL), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *HPKVClient) Delete(key string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/record/%s", c.baseURL, key), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func prettyJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return string(data)
	}
	return prettyJSON.String()
} 