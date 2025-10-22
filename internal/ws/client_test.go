package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Ryank90/runware-go-sdk/models"
	"github.com/gorilla/websocket"
)

// mockLogger for testing
type mockLogger struct {
	logs []string
}

func (m *mockLogger) Printf(format string, v ...interface{}) {
	// Store logs for assertion
	m.logs = append(m.logs, format)
}

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	config := DefaultWSConfig()
	logger := &mockLogger{}

	client := NewClient(apiKey, config, logger)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.apiKey != apiKey {
		t.Errorf("apiKey = %v, want %v", client.apiKey, apiKey)
	}

	if client.config != config {
		t.Errorf("config mismatch")
	}
}

func TestDefaultWSConfig(t *testing.T) {
	config := DefaultWSConfig()

	if config == nil {
		t.Fatal("DefaultWSConfig returned nil")
	}

	if config.URL == "" {
		t.Error("URL is empty")
	}

	if config.ConnectTimeout == 0 {
		t.Error("ConnectTimeout is zero")
	}

	if config.PingInterval == 0 {
		t.Error("PingInterval is zero")
	}

	if config.PongTimeout == 0 {
		t.Error("PongTimeout is zero")
	}

	if config.WriteTimeout == 0 {
		t.Error("WriteTimeout is zero")
	}

	if config.ReadBufferSize == 0 {
		t.Error("ReadBufferSize is zero")
	}

	if config.WriteBufferSize == 0 {
		t.Error("WriteBufferSize is zero")
	}
}

func TestIsConnected(t *testing.T) {
	client := NewClient("test-key", DefaultWSConfig(), &mockLogger{})

	if client.IsConnected() {
		t.Error("IsConnected() = true for new client, want false")
	}

	// Simulate connection
	client.mu.Lock()
	client.connected = true
	client.mu.Unlock()

	if !client.IsConnected() {
		t.Error("IsConnected() = false after setting connected, want true")
	}
}

// mockWebSocketServer creates a test WebSocket server
func mockWebSocketServer(t *testing.T, handler func(*websocket.Conn)) *httptest.Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Upgrade error: %v", err)
			return
		}
		defer conn.Close()

		handler(conn)
	}))

	return server
}

func TestConnect(t *testing.T) {
	server := mockWebSocketServer(t, func(conn *websocket.Conn) {
		// Read auth message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Verify it's an auth message
		var authMsg map[string]interface{}
		if err := json.Unmarshal(msg, &authMsg); err != nil {
			return
		}

		if authMsg["newConnectionToken"] == nil {
			t.Error("Expected auth message with newConnectionToken")
		}

		// Send auth success
		authResp := map[string]interface{}{
			"connectionSessionUUID": "test-session-uuid",
		}
		conn.WriteJSON(authResp)

		// Keep connection alive
		time.Sleep(100 * time.Millisecond)
	})
	defer server.Close()

	// Convert http://... to ws://...
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := DefaultWSConfig()
	config.URL = wsURL
	config.EnableAutoReconnect = false

	client := NewClient("test-api-key", config, &mockLogger{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if !client.IsConnected() {
		t.Error("Client not connected after Connect()")
	}

	// Cleanup
	client.Disconnect()
}

func TestConnectTimeout(t *testing.T) {
	// Use a context with timeout instead of configuring client timeout
	// since the WebSocket connection happens quickly but auth might not complete
	config := DefaultWSConfig()
	config.URL = "ws://localhost:65535" // Non-existent server
	config.EnableAutoReconnect = false

	client := NewClient("test-api-key", config, &mockLogger{})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := client.Connect(ctx)

	if err == nil {
		t.Error("Connect() should fail but succeeded")
		client.Disconnect()
	}
}

func TestDisconnect(t *testing.T) {
	server := mockWebSocketServer(t, func(conn *websocket.Conn) {
		// Read auth
		conn.ReadMessage()
		// Send auth success
		authResp := map[string]interface{}{
			"connectionSessionUUID": "test-session",
		}
		conn.WriteJSON(authResp)

		// Wait for close
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := DefaultWSConfig()
	config.URL = wsURL
	config.EnableAutoReconnect = false

	client := NewClient("test-api-key", config, &mockLogger{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if err := client.Disconnect(); err != nil {
		t.Errorf("Disconnect() error = %v", err)
	}

	if client.IsConnected() {
		t.Error("Client still connected after Disconnect()")
	}

	// Disconnect again should be safe
	if err := client.Disconnect(); err != nil {
		t.Errorf("Second Disconnect() error = %v", err)
	}
}

func TestSend(t *testing.T) {
	receivedMessage := make(chan []byte, 1)

	server := mockWebSocketServer(t, func(conn *websocket.Conn) {
		// Read and respond to auth
		conn.ReadMessage()
		conn.WriteJSON(map[string]interface{}{
			"connectionSessionUUID": "test-session",
		})

		// Read the actual message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Read error: %v", err)
			return
		}
		receivedMessage <- msg

		// Keep connection open
		time.Sleep(100 * time.Millisecond)
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := DefaultWSConfig()
	config.URL = wsURL
	config.EnableAutoReconnect = false

	client := NewClient("test-api-key", config, &mockLogger{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer client.Disconnect()

	// Create a test request
	req := models.NewImageInferenceRequest("test prompt", "test-model", 512, 512)

	handler := func(data interface{}, err error) {
		// Handler for the response
	}

	sendCtx, sendCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer sendCancel()

	if err := client.Send(sendCtx, req, handler); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	// Verify message was received by server
	select {
	case msg := <-receivedMessage:
		// The message is sent as an array wrapping the request
		var received []map[string]interface{}
		if err := json.Unmarshal(msg, &received); err != nil {
			t.Fatalf("Failed to unmarshal received message: %v", err)
		}

		if len(received) == 0 {
			t.Fatal("Received empty array")
		}

		if received[0]["taskType"] != models.TaskTypeImageInference {
			t.Errorf("taskType = %v, want %v", received[0]["taskType"], models.TaskTypeImageInference)
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestHandlerRegistrationAndRemoval(t *testing.T) {
	client := NewClient("test-key", DefaultWSConfig(), &mockLogger{})

	taskUUID := "test-task-uuid"
	handlerCalled := false

	handler := func(data interface{}, err error) {
		if err != nil {
			t.Errorf("Handler received error: %v", err)
		}
		handlerCalled = true
	}

	client.handlersMu.Lock()
	client.handlers[taskUUID] = handler
	client.handlersMu.Unlock()

	// Verify handler is registered
	client.handlersMu.RLock()
	_, exists := client.handlers[taskUUID]
	client.handlersMu.RUnlock()

	if !exists {
		t.Error("Handler not registered")
	}

	// Test calling handler
	handler(nil, nil)

	time.Sleep(50 * time.Millisecond)

	if !handlerCalled {
		t.Error("Handler was not called")
	}

	// Test handler removal
	client.RemoveHandler(taskUUID)

	client.handlersMu.RLock()
	_, exists = client.handlers[taskUUID]
	client.handlersMu.RUnlock()

	if exists {
		t.Error("Handler still exists after removal")
	}
}

func TestErrorHandling(t *testing.T) {
	client := NewClient("test-key", DefaultWSConfig(), &mockLogger{})

	// Start error logger
	go client.logErrorsLoop()

	// Send an error
	client.errorChan <- fmt.Errorf("test error")

	// Give it time to process
	time.Sleep(50 * time.Millisecond)

	// Should not panic or block
}

func TestConcurrentSend(t *testing.T) {
	server := mockWebSocketServer(t, func(conn *websocket.Conn) {
		// Auth
		conn.ReadMessage()
		conn.WriteJSON(map[string]interface{}{
			"connectionSessionUUID": "test-session",
		})

		// Read multiple messages
		for i := 0; i < 10; i++ {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}

		time.Sleep(100 * time.Millisecond)
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := DefaultWSConfig()
	config.URL = wsURL
	config.EnableAutoReconnect = false

	client := NewClient("test-api-key", config, &mockLogger{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer client.Disconnect()

	// Send multiple messages concurrently
	errChan := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			req := models.NewImageInferenceRequest("test", "model", 512, 512)
			handler := func(data interface{}, err error) {}
			sendCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			errChan <- client.Send(sendCtx, req, handler)
		}(i)
	}

	// Check all sends succeeded
	for i := 0; i < 10; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("Concurrent Send() #%d error = %v", i, err)
		}
	}
}

func TestSendWithoutConnection(t *testing.T) {
	client := NewClient("test-key", DefaultWSConfig(), &mockLogger{})

	req := models.NewImageInferenceRequest("test", "model", 512, 512)
	handler := func(data interface{}, err error) {}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := client.Send(ctx, req, handler)
	if err == nil {
		t.Error("Send() should fail when not connected")
	}
}
