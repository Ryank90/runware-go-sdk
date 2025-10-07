package runware

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// DefaultWSURL is the default WebSocket URL for the Runware API
	DefaultWSURL = "wss://ws-api.runware.ai/v1"

	// Default timeouts and intervals
	defaultConnectTimeout    = 30 * time.Second
	defaultPingInterval      = 30 * time.Second
	defaultPongTimeout       = 10 * time.Second
	defaultReconnectDelay    = 5 * time.Second
	defaultMaxReconnectDelay = 60 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultReadBufferSize    = 4096
	defaultWriteBufferSize   = 4096
)

// WSConfig contains WebSocket configuration options
type WSConfig struct {
	// URL is the WebSocket endpoint URL
	URL string

	// ConnectTimeout is the timeout for establishing a connection
	ConnectTimeout time.Duration

	// PingInterval is the interval between ping messages
	PingInterval time.Duration

	// PongTimeout is the timeout for receiving a pong response
	PongTimeout time.Duration

	// ReconnectDelay is the initial delay before reconnecting
	ReconnectDelay time.Duration

	// MaxReconnectDelay is the maximum delay between reconnect attempts
	MaxReconnectDelay time.Duration

	// WriteTimeout is the timeout for write operations
	WriteTimeout time.Duration

	// EnableAutoReconnect enables automatic reconnection on disconnect
	EnableAutoReconnect bool

	// ReadBufferSize specifies the size of the read buffer
	ReadBufferSize int

	// WriteBufferSize specifies the size of the write buffer
	WriteBufferSize int
}

// DefaultWSConfig returns a default WebSocket configuration
func DefaultWSConfig() *WSConfig {
	return &WSConfig{
		URL:                 DefaultWSURL,
		ConnectTimeout:      defaultConnectTimeout,
		PingInterval:        defaultPingInterval,
		PongTimeout:         defaultPongTimeout,
		ReconnectDelay:      defaultReconnectDelay,
		MaxReconnectDelay:   defaultMaxReconnectDelay,
		WriteTimeout:        defaultWriteTimeout,
		EnableAutoReconnect: true,
		ReadBufferSize:      defaultReadBufferSize,
		WriteBufferSize:     defaultWriteBufferSize,
	}
}

// wsClient manages the WebSocket connection
type wsClient struct {
	config        *WSConfig
	apiKey        string
	conn          *websocket.Conn
	mu            sync.RWMutex
	connected     bool
	stopChan      chan struct{}
	reconnectChan chan struct{}
	messageChan   chan []byte
	errorChan     chan error
	handlers      map[string]ResponseHandler
	handlersMu    sync.RWMutex
	wg            sync.WaitGroup
}

// ResponseHandler is a function that handles responses for a specific task
type ResponseHandler func(data interface{}, err error)

// newWSClient creates a new WebSocket client
func newWSClient(apiKey string, config *WSConfig) *wsClient {
	if config == nil {
		config = DefaultWSConfig()
	}

	return &wsClient{
		config:        config,
		apiKey:        apiKey,
		stopChan:      make(chan struct{}),
		reconnectChan: make(chan struct{}, 1),
		messageChan:   make(chan []byte, 100),
		errorChan:     make(chan error, 10),
		handlers:      make(map[string]ResponseHandler),
	}
}

// Connect establishes a WebSocket connection
func (c *wsClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return ErrAlreadyConnected
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: c.config.ConnectTimeout,
		ReadBufferSize:   c.config.ReadBufferSize,
		WriteBufferSize:  c.config.WriteBufferSize,
	}

	conn, _, err := dialer.DialContext(ctx, c.config.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.conn = conn
	c.connected = true

	// Send authentication
	if err := c.authenticate(); err != nil {
		c.conn.Close()
		c.connected = false
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Start background goroutines
	c.wg.Add(3)
	go c.readLoop()
	go c.processMessages()
	go c.pingLoop()

	if c.config.EnableAutoReconnect {
		c.wg.Add(1)
		go c.reconnectLoop(ctx)
	}

	return nil
}

// Disconnect closes the WebSocket connection
func (c *wsClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	close(c.stopChan)
	c.wg.Wait()

	if c.conn != nil {
		err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
		c.conn = nil
		c.connected = false
		return err
	}

	c.connected = false
	return nil
}

// IsConnected returns whether the client is connected
func (c *wsClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// Send sends a request and registers a handler for the response
func (c *wsClient) Send(ctx context.Context, request interface{}, handler ResponseHandler) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}

	data, err := json.Marshal([]interface{}{request})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Extract taskUUID directly from the request object
	reqData, _ := json.Marshal(request)
	var taskMap map[string]interface{}
	var taskUUID string
	if err := json.Unmarshal(reqData, &taskMap); err == nil {
		if uuid, ok := taskMap["taskUUID"].(string); ok {
			taskUUID = uuid
		}
	}

	if taskUUID == "" {
		return fmt.Errorf("request missing taskUUID")
	}

	// Register handler
	c.handlersMu.Lock()
	c.handlers[taskUUID] = handler
	c.handlersMu.Unlock()

	// Send the request
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return ErrNotConnected
	}

	if err := conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// authenticate sends authentication credentials
func (c *wsClient) authenticate() error {
	// Per Runware docs, auth must be an array with taskType and apiKey
	authMsg := []map[string]interface{}{
		{
			"taskType": "authentication",
			"apiKey":   c.apiKey,
		},
	}

	data, err := json.Marshal(authMsg)
	if err != nil {
		return err
	}

	if err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout)); err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// readLoop reads messages from the WebSocket
func (c *wsClient) readLoop() {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopChan:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					c.errorChan <- fmt.Errorf("unexpected close: %w", err)
					c.triggerReconnect()
				}
				return
			}

			c.messageChan <- message
		}
	}
}

// processMessages processes incoming messages
func (c *wsClient) processMessages() {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopChan:
			return
		case message := <-c.messageChan:
			c.handleMessage(message)
		}
	}
}

// handleMessage handles a single message
func (c *wsClient) handleMessage(message []byte) {
	var response struct {
		Data []json.RawMessage `json:"data,omitempty"`
	}

	if err := json.Unmarshal(message, &response); err != nil {
		// Try to parse as error response
		var errResp ErrorResponse
		if err := json.Unmarshal(message, &errResp); err == nil {
			c.handlersMu.RLock()
			handler, ok := c.handlers[errResp.TaskUUID]
			c.handlersMu.RUnlock()

			if ok {
				handler(nil, NewAPIError(&errResp))
				c.handlersMu.Lock()
				delete(c.handlers, errResp.TaskUUID)
				c.handlersMu.Unlock()
			}
		}
		return
	}

	// Process each data item
	for _, item := range response.Data {
		var baseResp struct {
			TaskUUID string `json:"taskUUID"`
			TaskType string `json:"taskType"`
		}

		if err := json.Unmarshal(item, &baseResp); err != nil {
			continue
		}

		c.handlersMu.RLock()
		handler, ok := c.handlers[baseResp.TaskUUID]
		c.handlersMu.RUnlock()

		if !ok {
			continue
		}

		// Parse based on task type
		var result interface{}
		switch baseResp.TaskType {
		case TaskTypeImageInference:
			var resp ImageInferenceResponse
			if err := json.Unmarshal(item, &resp); err == nil {
				result = &resp
			}
		case TaskTypeImageUpload:
			var resp UploadImageResponse
			if err := json.Unmarshal(item, &resp); err == nil {
				result = &resp
			}
		case TaskTypeUpscaleGan:
			var resp UpscaleGanResponse
			if err := json.Unmarshal(item, &resp); err == nil {
				result = &resp
			}
		case TaskTypeImageBackgroundRemoval:
			var resp RemoveImageBackgroundResponse
			if err := json.Unmarshal(item, &resp); err == nil {
				result = &resp
			}
		case TaskTypePromptEnhance:
			var resp EnhancePromptResponse
			if err := json.Unmarshal(item, &resp); err == nil {
				result = &resp
			}
		case TaskTypeImageCaption:
			var resp ImageCaptionResponse
			if err := json.Unmarshal(item, &resp); err == nil {
				result = &resp
			}
		case "videoInference", "getResponse":
			var resp VideoInferenceResponse
			if err := json.Unmarshal(item, &resp); err == nil {
				result = &resp
			}
		}

		handler(result, nil)

		// Remove handler after processing (single response per task)
		c.handlersMu.Lock()
		delete(c.handlers, baseResp.TaskUUID)
		c.handlersMu.Unlock()
	}
}

// pingLoop sends periodic ping messages
func (c *wsClient) pingLoop() {
	defer c.wg.Done()

	ticker := time.NewTicker(c.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			c.mu.RLock()
			conn := c.conn
			c.mu.RUnlock()

			if conn == nil {
				return
			}

			if err := conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout)); err != nil {
				c.triggerReconnect()
				return
			}
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.triggerReconnect()
				return
			}
		}
	}
}

// reconnectLoop handles automatic reconnection
func (c *wsClient) reconnectLoop(ctx context.Context) {
	defer c.wg.Done()

	delay := c.config.ReconnectDelay

	for {
		select {
		case <-c.stopChan:
			return
		case <-c.reconnectChan:
			time.Sleep(delay)

			c.mu.Lock()
			if c.connected {
				c.mu.Unlock()
				continue
			}

			// Close old connection
			if c.conn != nil {
				c.conn.Close()
				c.conn = nil
			}
			c.mu.Unlock()

			// Attempt reconnection
			if err := c.Connect(ctx); err != nil {
				delay *= 2
				if delay > c.config.MaxReconnectDelay {
					delay = c.config.MaxReconnectDelay
				}
				c.triggerReconnect()
			} else {
				delay = c.config.ReconnectDelay
			}
		}
	}
}

// triggerReconnect triggers a reconnection attempt
func (c *wsClient) triggerReconnect() {
	c.mu.Lock()
	c.connected = false
	c.mu.Unlock()

	select {
	case c.reconnectChan <- struct{}{}:
	default:
	}
}
