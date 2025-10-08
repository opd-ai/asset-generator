package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Config holds the client configuration
type Config struct {
	BaseURL string
	APIKey  string
	Verbose bool
}

// SwarmClient is the main client for interacting with SwarmUI API
type SwarmClient struct {
	config     *Config
	httpClient *http.Client
	wsConn     *websocket.Conn
	mu         sync.RWMutex
	sessions   map[string]*GenerationSession
}

// GenerationRequest represents a request to generate an asset
type GenerationRequest struct {
	Prompt     string                 `json:"prompt"`
	Model      string                 `json:"model,omitempty"`
	Parameters map[string]interface{} `json:"parameters"`
	SessionID  string                 `json:"session_id,omitempty"`
}

// GenerationResult represents the result of a generation
type GenerationResult struct {
	SessionID  string                 `json:"session_id"`
	ImagePaths []string               `json:"image_paths"`
	Metadata   map[string]interface{} `json:"metadata"`
	Status     string                 `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
}

// GenerationSession tracks an ongoing generation
type GenerationSession struct {
	ID        string
	Status    string
	Progress  float64
	StartTime time.Time
	Result    *GenerationResult
}

// Model represents a SwarmUI model
type Model struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Loaded      bool   `json:"loaded"`
}

// NewSwarmClient creates a new SwarmUI client
func NewSwarmClient(config *Config) (*SwarmClient, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	return &SwarmClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		sessions: make(map[string]*GenerationSession),
	}, nil
}

// GetNewSession gets a new session ID from SwarmUI API
func (c *SwarmClient) GetNewSession(ctx context.Context) (string, error) {
	endpoint := fmt.Sprintf("%s/API/GetNewSession", c.config.BaseURL)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return "", fmt.Errorf("failed to create session request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Request: POST %s\n", endpoint)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("session request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("session API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResp struct {
		SessionID string `json:"session_id"`
		Error     string `json:"error,omitempty"`
		ErrorID   string `json:"error_id,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("failed to decode session response: %w", err)
	}

	if apiResp.Error != "" {
		return "", fmt.Errorf("SwarmUI session error: %s", apiResp.Error)
	}

	if apiResp.SessionID == "" {
		return "", fmt.Errorf("SwarmUI did not return a session ID")
	}

	return apiResp.SessionID, nil
}

// GenerateImage generates an image using the SwarmUI API
func (c *SwarmClient) GenerateImage(ctx context.Context, req *GenerationRequest) (*GenerationResult, error) {
	// Get a new session ID from SwarmUI API
	sessionID, err := c.GetNewSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Create local session tracking
	session := &GenerationSession{
		ID:        sessionID,
		Status:    "pending",
		Progress:  0,
		StartTime: time.Now(),
	}

	c.mu.Lock()
	c.sessions[sessionID] = session
	c.mu.Unlock()

	// Ensure session cleanup on function exit (success or error)
	defer c.cleanupSession(sessionID)

	// Make HTTP request to generate endpoint
	endpoint := fmt.Sprintf("%s/API/GenerateText2Image", c.config.BaseURL)

	// Build request body with correct SwarmUI parameter names
	body := map[string]interface{}{
		"session_id": sessionID, // Required by SwarmUI API
		"prompt":     req.Prompt,
		"images":     1, // Default to 1 image, SwarmUI uses "images" not "batch_size"
	}

	// Add batch size parameter if specified (SwarmUI expects "images" field)
	if batchSize, ok := req.Parameters["batch_size"]; ok && batchSize != nil {
		if bs, isInt := batchSize.(int); isInt && bs > 0 {
			body["images"] = bs
		}
	}

	// Add model if specified
	if req.Model != "" {
		body["model"] = req.Model
	}

	// Add standard SwarmUI parameters with defaults
	if width, ok := req.Parameters["width"]; ok {
		body["width"] = width
	} else {
		body["width"] = 1024 // Default width
	}

	if height, ok := req.Parameters["height"]; ok {
		body["height"] = height
	} else {
		body["height"] = 1024 // Default height
	}

	if cfgScale, ok := req.Parameters["cfgscale"]; ok {
		body["cfgscale"] = cfgScale
	} else {
		body["cfgscale"] = 7.5 // Default CFG scale
	}

	if steps, ok := req.Parameters["steps"]; ok {
		body["steps"] = steps
	} else {
		body["steps"] = 20 // Default steps
	}

	if seed, ok := req.Parameters["seed"]; ok {
		body["seed"] = seed
	} else {
		body["seed"] = -1 // Random seed
	}

	// Add any other parameters from the request
	for k, v := range req.Parameters {
		// Skip parameters we've already handled
		if k != "batch_size" && k != "width" && k != "height" && k != "cfgscale" && k != "steps" && k != "seed" {
			body[k] = v
		}
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Request: POST %s\n", endpoint)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response with SwarmUI error handling
	var apiResp struct {
		Images  []string               `json:"images"`
		Info    map[string]interface{} `json:"info"`
		Error   string                 `json:"error,omitempty"`
		ErrorID string                 `json:"error_id,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Handle SwarmUI-specific errors
	if apiResp.Error != "" {
		return nil, fmt.Errorf("SwarmUI error: %s", apiResp.Error)
	}

	if apiResp.ErrorID != "" {
		return nil, fmt.Errorf("SwarmUI error (ID: %s)", apiResp.ErrorID)
	}

	// Build result
	result := &GenerationResult{
		SessionID:  sessionID,
		ImagePaths: apiResp.Images,
		Metadata:   apiResp.Info,
		Status:     "completed",
		CreatedAt:  time.Now(),
	}

	// Update session
	c.mu.Lock()
	session.Status = "completed"
	session.Progress = 1.0
	session.Result = result
	c.mu.Unlock()

	return result, nil
}

// cleanupSession removes a session from memory to prevent memory leaks
func (c *SwarmClient) cleanupSession(sessionID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.sessions, sessionID)
}

// cleanupOldSessions removes sessions older than the specified duration
func (c *SwarmClient) cleanupOldSessions(maxAge time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	cutoff := time.Now().Add(-maxAge)
	for sessionID, session := range c.sessions {
		if session.StartTime.Before(cutoff) {
			delete(c.sessions, sessionID)
		}
	}
}

// ListModels lists all available models
func (c *SwarmClient) ListModels() ([]Model, error) {
	endpoint := fmt.Sprintf("%s/API/ListModels", c.config.BaseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Request: GET %s\n", endpoint)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResp struct {
		Models []Model `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return apiResp.Models, nil
}

// GetModel gets details about a specific model
func (c *SwarmClient) GetModel(name string) (*Model, error) {
	// Using correct SwarmUI API endpoint pattern
	endpoint := fmt.Sprintf("%s/API/ListModels", c.config.BaseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Request: GET %s\n", endpoint)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var model Model
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &model, nil
}

// Close closes any open connections
func (c *SwarmClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.wsConn != nil {
		return c.wsConn.Close()
	}

	return nil
}
