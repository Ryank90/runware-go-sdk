package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	models "github.com/Ryank90/runware-go-sdk/models"
	"github.com/gorilla/websocket"
)

const (
	// DefaultWSURL is the default WebSocket URL for the Runware API
	DefaultWSURL = "wss://ws-api.runware.ai/v1"

	// Default timeouts and intervals
	defaultConnectTimeout = 30 * time.Second
	defaultPingInterval   = 30 * time.Second
	// Ensure pong timeout comfortably exceeds ping interval to avoid premature timeouts
	defaultPongTimeout       = 90 * time.Second
	defaultReconnectDelay    = 5 * time.Second
	defaultMaxReconnectDelay = 60 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultReadBufferSize    = 4096
	defaultWriteBufferSize   = 4096
)

// DebugLogger matches the parent package interface for debug logging
type DebugLogger interface {
	Printf(format string, v ...interface{})
}

// WSConfig contains WebSocket configuration options
type WSConfig struct {
	URL                 string
	ConnectTimeout      time.Duration
	PingInterval        time.Duration
	PongTimeout         time.Duration
	ReconnectDelay      time.Duration
	MaxReconnectDelay   time.Duration
	WriteTimeout        time.Duration
	EnableAutoReconnect bool
	ReadBufferSize      int
	WriteBufferSize     int
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

// ResponseHandler handles responses for a specific task
type ResponseHandler func(data interface{}, err error)

// Client manages the WebSocket connection
type Client struct {
	config        *WSConfig
	apiKey        string
	conn          *websocket.Conn
	mu            sync.RWMutex
	writeMu       sync.Mutex
	connected     bool
	stopChan      chan struct{}
	reconnectChan chan struct{}
	messageChan   chan []byte
	errorChan     chan error
	handlers      map[string]ResponseHandler
	handlersMu    sync.RWMutex
	wg            sync.WaitGroup
	debugLogger   DebugLogger
}

// NewClient creates a new WebSocket client
func NewClient(apiKey string, config *WSConfig, debugLogger DebugLogger) *Client {
	if config == nil {
		config = DefaultWSConfig()
	}
	if debugLogger == nil {
		debugLogger = &noopLogger{}
	}
	return &Client{
		config:        config,
		apiKey:        apiKey,
		stopChan:      make(chan struct{}),
		reconnectChan: make(chan struct{}, 1),
		messageChan:   make(chan []byte, 100),
		errorChan:     make(chan error, 10),
		handlers:      make(map[string]ResponseHandler),
		debugLogger:   debugLogger,
	}
}

type noopLogger struct{}

func (n *noopLogger) Printf(string, ...interface{}) {}

// Connect establishes a WebSocket connection
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connected {
		return fmt.Errorf("already connected")
	}

	c.debugLogger.Printf("Connecting to %s", c.config.URL)

	dialer := websocket.Dialer{
		HandshakeTimeout: c.config.ConnectTimeout,
		ReadBufferSize:   c.config.ReadBufferSize,
		WriteBufferSize:  c.config.WriteBufferSize,
	}

	conn, resp, err := dialer.DialContext(ctx, c.config.URL, nil)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.conn = conn
	c.connected = true

	c.debugLogger.Printf("WebSocket connected, authenticating...")

	// Set initial read deadline and pong handler
	if err := c.conn.SetReadDeadline(time.Now().Add(c.config.PongTimeout)); err == nil {
		c.conn.SetPongHandler(func(string) error {
			return c.conn.SetReadDeadline(time.Now().Add(c.config.PongTimeout))
		})
	}

	if err := c.authenticate(); err != nil {
		_ = c.conn.Close()
		c.connected = false
		return fmt.Errorf("authentication failed: %w", err)
	}

	c.debugLogger.Printf("Authenticated successfully")

	c.wg.Add(4)
	go c.readLoop()
	go c.processMessages()
	go c.pingLoop()
	go c.logErrorsLoop()

	if c.config.EnableAutoReconnect {
		c.wg.Add(1)
		go c.reconnectLoop(ctx)
	}
	return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() error {
	c.mu.Lock()
	if !c.connected {
		c.mu.Unlock()
		return nil
	}
	c.connected = false
	select {
	case <-c.stopChan:
	default:
		close(c.stopChan)
	}
	c.mu.Unlock()

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		_ = c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// Send sends a request and registers a handler for the response
func (c *Client) Send(ctx context.Context, request interface{}, handler ResponseHandler) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	data, err := json.Marshal([]interface{}{request})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Extract task fields via optional interface to avoid extra JSON work
	var taskUUID, taskType string
	if ti, ok := request.(models.TaskIdentifiable); ok {
		taskUUID = ti.GetTaskUUID()
		taskType = ti.GetTaskType()
	}
	if taskUUID == "" {
		return fmt.Errorf("request missing taskUUID")
	}

	c.debugLogger.Printf("Sending request: %s (TaskUUID: %s)", taskType, taskUUID)

	c.handlersMu.Lock()
	c.handlers[taskUUID] = handler
	c.handlersMu.Unlock()

	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return fmt.Errorf("not connected")
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	c.debugLogger.Printf("Request sent successfully: %s (TaskUUID: %s)", taskType, taskUUID)
	return nil
}

func (c *Client) authenticate() error {
	authMsg := []map[string]interface{}{
		{"taskType": "authentication", "apiKey": c.apiKey},
	}
	data, err := json.Marshal(authMsg)
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout)); err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) readLoop() {
	defer c.wg.Done()
	for {
		select {
		case <-c.stopChan:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				c.errorChan <- fmt.Errorf("read error: %w", err)
				c.triggerReconnect()
				return
			}
			_ = c.conn.SetReadDeadline(time.Now().Add(c.config.PongTimeout))
			c.messageChan <- message
		}
	}
}

func (c *Client) processMessages() {
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

// The message handler and routing are kept here for simplicity
func (c *Client) handleMessage(message []byte) {
	var response struct {
		Data []json.RawMessage `json:"data,omitempty"`
	}
	if err := json.Unmarshal(message, &response); err != nil {
		c.handleErrorResponse(message)
		return
	}

	for _, item := range response.Data {
		c.processResponseItem(item)
	}
}

func (c *Client) handleErrorResponse(message []byte) {
	var errResp struct {
		Error    string `json:"error"`
		ErrorID  string `json:"errorId,omitempty"`
		Code     string `json:"code,omitempty"`
		Message  string `json:"message,omitempty"`
		TaskUUID string `json:"taskUUID,omitempty"`
		TaskType string `json:"taskType,omitempty"`
	}
	if json.Unmarshal(message, &errResp) == nil && errResp.TaskUUID != "" {
		c.handlersMu.RLock()
		h, ok := c.handlers[errResp.TaskUUID]
		c.handlersMu.RUnlock()
		if ok {
			h(nil, fmt.Errorf("api error: %s", errResp.Error))
			c.removeHandler(errResp.TaskUUID)
		}
	}
}

func (c *Client) processResponseItem(item json.RawMessage) {
	var baseResp struct {
		TaskUUID string `json:"taskUUID"`
		TaskType string `json:"taskType"`
		Status   string `json:"status,omitempty"`
	}
	if err := json.Unmarshal(item, &baseResp); err != nil {
		return
	}

	c.handlersMu.RLock()
	h, ok := c.handlers[baseResp.TaskUUID]
	c.handlersMu.RUnlock()
	if !ok {
		return
	}

	result := c.parseResponseByType(baseResp.TaskType, item)
	h(result, nil)
}

func (c *Client) parseResponseByType(taskType string, item json.RawMessage) interface{} {
	switch taskType {
	case models.TaskTypeImageInference:
		return c.parseImageInferenceResponse(item)
	case models.TaskTypeImageUpload:
		return c.parseUploadImageResponse(item)
	case models.TaskTypeUpscaleGan:
		return c.parseUpscaleGanResponse(item)
	case models.TaskTypeImageBackgroundRemoval:
		return c.parseRemoveBackgroundResponse(item)
	case models.TaskTypePromptEnhance:
		return c.parseEnhancePromptResponse(item)
	case models.TaskTypeImageCaption:
		return c.parseImageCaptionResponse(item)
	case models.TaskTypeVideoInference:
		return c.parseVideoInferenceResponse(item)
	case models.TaskTypeAudioInference:
		return c.parseAudioInferenceResponse(item)
	case models.TaskTypeGetResponse:
		return c.parseGetResponse(item)
	}
	return nil
}

func (c *Client) parseImageInferenceResponse(item json.RawMessage) interface{} {
	var resp models.ImageInferenceResponse
	if err := json.Unmarshal(item, &resp); err == nil {
		return &resp
	}
	return nil
}

func (c *Client) parseUploadImageResponse(item json.RawMessage) interface{} {
	var resp models.UploadImageResponse
	if err := json.Unmarshal(item, &resp); err == nil {
		return &resp
	}
	return nil
}

func (c *Client) parseUpscaleGanResponse(item json.RawMessage) interface{} {
	var resp models.UpscaleGanResponse
	if err := json.Unmarshal(item, &resp); err == nil {
		return &resp
	}
	return nil
}

func (c *Client) parseRemoveBackgroundResponse(item json.RawMessage) interface{} {
	var resp models.RemoveImageBackgroundResponse
	if err := json.Unmarshal(item, &resp); err == nil {
		return &resp
	}
	return nil
}

func (c *Client) parseEnhancePromptResponse(item json.RawMessage) interface{} {
	var resp models.EnhancePromptResponse
	if err := json.Unmarshal(item, &resp); err == nil {
		return &resp
	}
	return nil
}

func (c *Client) parseImageCaptionResponse(item json.RawMessage) interface{} {
	var resp models.ImageCaptionResponse
	if err := json.Unmarshal(item, &resp); err == nil {
		return &resp
	}
	return nil
}

func (c *Client) parseVideoInferenceResponse(item json.RawMessage) interface{} {
	var resp models.VideoInferenceResponse
	if err := json.Unmarshal(item, &resp); err == nil {
		return &resp
	}
	return nil
}

func (c *Client) parseAudioInferenceResponse(item json.RawMessage) interface{} {
	var resp models.AudioInferenceResponse
	if err := json.Unmarshal(item, &resp); err == nil {
		return &resp
	}
	return nil
}

func (c *Client) parseGetResponse(item json.RawMessage) interface{} {
	var videoResp models.VideoInferenceResponse
	if err := json.Unmarshal(item, &videoResp); err == nil {
		if videoResp.Status != "" || videoResp.VideoUUID != "" || videoResp.VideoURL != nil || videoResp.ThumbnailURL != nil {
			return &videoResp
		}
	}
	var audioResp models.AudioInferenceResponse
	if err := json.Unmarshal(item, &audioResp); err == nil {
		if audioResp.Status != "" || audioResp.AudioUUID != "" || audioResp.AudioURL != nil || audioResp.AudioBase64Data != nil || audioResp.AudioDataURI != nil {
			return &audioResp
		}
	}
	return nil
}

func (c *Client) pingLoop() {
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
			c.writeMu.Lock()
			if err := conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout)); err != nil {
				c.writeMu.Unlock()
				c.triggerReconnect()
				return
			}
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.writeMu.Unlock()
				c.triggerReconnect()
				return
			}
			c.writeMu.Unlock()
		}
	}
}

func (c *Client) reconnectLoop(ctx context.Context) {
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
			if c.conn != nil {
				_ = c.conn.Close()
				c.conn = nil
			}
			c.mu.Unlock()
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

func (c *Client) triggerReconnect() {
	select {
	case <-c.stopChan:
		return
	default:
	}
	if c.mu.TryLock() {
		c.connected = false
		c.mu.Unlock()
	}
	select {
	case c.reconnectChan <- struct{}{}:
	default:
	}
}

func (c *Client) removeHandler(taskUUID string) {
	c.handlersMu.Lock()
	delete(c.handlers, taskUUID)
	c.handlersMu.Unlock()
}

// RemoveHandler exposes handler removal for external coordination (e.g., after final response)
func (c *Client) RemoveHandler(taskUUID string) { c.removeHandler(taskUUID) }

func (c *Client) logErrorsLoop() {
	defer c.wg.Done()
	for {
		select {
		case <-c.stopChan:
			return
		case err := <-c.errorChan:
			if err != nil {
				c.debugLogger.Printf("websocket error: %v", err)
			}
		}
	}
}
