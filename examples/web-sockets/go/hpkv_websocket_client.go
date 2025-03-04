package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type OperationCode int

const (
	Get OperationCode = iota + 1
	Insert
	Update
	Delete
)

type HPKVWebSocketClient struct {
	conn            *websocket.Conn
	url             string
	messageID       uint64
	responseFutures sync.Map
}

type Message struct {
	Op        OperationCode `json:"op"`
	Key       string        `json:"key"`
	Value     interface{}   `json:"value,omitempty"`
	MessageID uint64        `json:"messageId"`
}

type Response struct {
	MessageID uint64      `json:"messageId"`
	Value     interface{} `json:"value,omitempty"`
	Error     string      `json:"error,omitempty"`
}

type ResponseFuture struct {
	Response chan Response
	Error    chan error
}

func NewHPKVWebSocketClient(baseURL, apiKey string) (*HPKVWebSocketClient, error) {
	// Convert HTTP/HTTPS URL to WS/WSS URL
	wsURL := baseURL
	if strings.HasPrefix(baseURL, "https://") {
		wsURL = "wss://" + strings.TrimPrefix(baseURL, "https://")
	} else if strings.HasPrefix(baseURL, "http://") {
		wsURL = "ws://" + strings.TrimPrefix(baseURL, "http://")
	} else if !strings.HasPrefix(baseURL, "ws://") && !strings.HasPrefix(baseURL, "wss://") {
		wsURL = "wss://" + baseURL
	}

	wsURL = fmt.Sprintf("%s/ws?apiKey=%s", wsURL, apiKey)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	client := &HPKVWebSocketClient{
		conn:      conn,
		url:       wsURL,
		messageID: 0,
	}

	go client.handleMessages()

	return client, nil
}

func (c *HPKVWebSocketClient) handleMessages() {
	for {
		var response Response
		err := c.conn.ReadJSON(&response)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		if future, ok := c.responseFutures.LoadAndDelete(response.MessageID); ok {
			f := future.(*ResponseFuture)
			if response.Error != "" {
				f.Error <- fmt.Errorf(response.Error)
			} else {
				f.Response <- response
			}
		}
	}
}

func (c *HPKVWebSocketClient) sendMessage(message Message) (*Response, error) {
	messageID := atomic.AddUint64(&c.messageID, 1)
	message.MessageID = messageID

	future := &ResponseFuture{
		Response: make(chan Response, 1),
		Error:    make(chan error, 1),
	}
	c.responseFutures.Store(messageID, future)

	err := c.conn.WriteJSON(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %v", err)
	}

	select {
	case response := <-future.Response:
		return &response, nil
	case err := <-future.Error:
		return nil, err
	}
}

func (c *HPKVWebSocketClient) Create(key string, value interface{}) error {
	// Convert value to JSON string if it's not already a string
	var jsonValue interface{}
	switch v := value.(type) {
	case string:
		jsonValue = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
		jsonValue = string(jsonBytes)
	}

	message := Message{
		Op:    Insert,
		Key:   key,
		Value: jsonValue,
	}

	response, err := c.sendMessage(message)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	return nil
}

func (c *HPKVWebSocketClient) Read(key string) (interface{}, error) {
	message := Message{
		Op:  Get,
		Key: key,
	}

	response, err := c.sendMessage(message)
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}

	// Handle string values that might be JSON
	if strValue, ok := response.Value.(string); ok {
		var jsonValue interface{}
		if err := json.Unmarshal([]byte(strValue), &jsonValue); err == nil {
			return jsonValue, nil
		}
	}

	return response.Value, nil
}

func (c *HPKVWebSocketClient) Update(key string, value interface{}, partialUpdate bool) error {
	// For partial updates, first read the existing value
	var existingValue interface{}
	if partialUpdate {
		var err error
		existingValue, err = c.Read(key)
		if err != nil {
			return fmt.Errorf("failed to read existing value: %v", err)
		}

		// Merge the existing value with the new value
		if existingMap, ok := existingValue.(map[string]interface{}); ok {
			if newMap, ok := value.(map[string]interface{}); ok {
				for k, v := range newMap {
					existingMap[k] = v
				}
				value = existingMap
			}
		}
	}

	// Convert value to JSON string if it's not already a string
	var jsonValue interface{}
	switch v := value.(type) {
	case string:
		jsonValue = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
		jsonValue = string(jsonBytes)
	}

	message := Message{
		Op:    Insert, // Always use full update
		Key:   key,
		Value: jsonValue,
	}

	response, err := c.sendMessage(message)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	return nil
}

func (c *HPKVWebSocketClient) Delete(key string) error {
	message := Message{
		Op:  Delete,
		Key: key,
	}

	response, err := c.sendMessage(message)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	return nil
}

func (c *HPKVWebSocketClient) Close() error {
	return c.conn.Close()
}
