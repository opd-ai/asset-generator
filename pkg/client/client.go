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

// GenerateImage generates an image using the SwarmUI API
func (c *SwarmClient) GenerateImage(ctx context.Context, req *GenerationRequest) (*GenerationResult, error) {
	// Create a unique session ID
	sessionID := fmt.Sprintf("gen-%d", time.Now().UnixNano())
	req.SessionID = sessionID

	// Create session
	session := &GenerationSession{
		ID:        sessionID,
		Status:    "pending",
		Progress:  0,
		StartTime: time.Now(),
	}

	c.mu.Lock()
	c.sessions[sessionID] = session
	c.mu.Unlock()

	// Make HTTP request to generate endpoint
	endpoint := fmt.Sprintf("%s/API/GenerateText2Image", c.config.BaseURL)

	// Build request body
	body := map[string]interface{}{
		"prompt": req.Prompt,
		"images": req.Parameters["batch_size"],
	}

	// Add model if specified
	if req.Model != "" {
		body["model"] = req.Model
	}

	// Add parameters
	for k, v := range req.Parameters {
		body[k] = v
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

	// Parse response
	var apiResp struct {
		Images []string               `json:"images"`
		Info   map[string]interface{} `json:"info"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
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
	endpoint := fmt.Sprintf("%s/API/GetModel?name=%s", c.config.BaseURL, name)

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
