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
	wsConn     *websocket.Conn // Reserved for future WebSocket implementation
	mu         sync.RWMutex
	sessions   map[string]*GenerationSession
	sessionID  string // Current session ID for API calls
}

// ProgressCallback is called with progress updates during generation
type ProgressCallback func(progress float64, status string)

// GenerationRequest represents a request to generate an asset
type GenerationRequest struct {
	Prompt           string                 `json:"prompt"`
	Model            string                 `json:"model,omitempty"`
	Parameters       map[string]interface{} `json:"parameters"`
	SessionID        string                 `json:"session_id,omitempty"`
	ProgressCallback ProgressCallback       `json:"-"` // Not serialized, used for progress updates
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
			Timeout: 10 * time.Minute, // Extended timeout for Flux generation (can take 5-10 minutes)
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
	// Get or reuse session ID
	sessionID, err := c.ensureSession()
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
		body["width"] = 512 // Default width (matches CLI and documentation)
	}

	if height, ok := req.Parameters["height"]; ok {
		body["height"] = height
	} else {
		body["height"] = 512 // Default height (matches CLI and documentation)
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

	// Report initial progress
	if req.ProgressCallback != nil {
		req.ProgressCallback(0.0, "Starting generation...")
		session.Progress = 0.0
		session.Status = "starting"
	}

	// Start progress simulation in background for HTTP requests
	// Since we're using HTTP (not WebSocket), we simulate progress
	var progressDone chan bool
	if req.ProgressCallback != nil {
		progressDone = make(chan bool, 1)
		go c.simulateProgress(sessionID, req.ProgressCallback, progressDone)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		if progressDone != nil {
			progressDone <- true // Stop progress simulation
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Stop progress simulation
	if progressDone != nil {
		progressDone <- true
	}

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

	// Report completion
	if req.ProgressCallback != nil {
		req.ProgressCallback(1.0, "Generation completed")
	}

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

// parseSwarmUIError attempts to parse a SwarmUI error response from raw body bytes
// Returns nil if no SwarmUI error format is detected
func parseSwarmUIError(body []byte) error {
	var errResp struct {
		Error   string `json:"error,omitempty"`
		ErrorID string `json:"error_id,omitempty"`
	}

	// Try to parse as JSON
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil // Not a JSON error response
	}

	// Check if SwarmUI error fields are present
	if errResp.Error != "" {
		if errResp.ErrorID != "" {
			return fmt.Errorf("SwarmUI error (%s): %s", errResp.ErrorID, errResp.Error)
		}
		return fmt.Errorf("SwarmUI error: %s", errResp.Error)
	}

	return nil // No SwarmUI error detected
}

// ensureSession ensures we have a valid session ID, getting a new one if needed
func (c *SwarmClient) ensureSession() (string, error) {
	c.mu.RLock()
	sessionID := c.sessionID
	c.mu.RUnlock()

	if sessionID != "" {
		return sessionID, nil
	}

	// Need to get a new session
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check again in case another goroutine got a session while we waited
	if c.sessionID != "" {
		return c.sessionID, nil
	}

	newSessionID, err := c.getNewSession()
	if err != nil {
		return "", err
	}

	c.sessionID = newSessionID
	return newSessionID, nil
}

// getNewSession gets a new session ID from SwarmUI
func (c *SwarmClient) getNewSession() (string, error) {
	endpoint := fmt.Sprintf("%s/API/GetNewSession", c.config.BaseURL)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return "", fmt.Errorf("failed to create session request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Getting new session: POST %s\n", endpoint)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("session request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read session response: %w", err)
	}

	if c.config.Verbose {
		fmt.Printf("Session Response Status: %d\nSession Response Body: %s\n", resp.StatusCode, string(bodyBytes))
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("session API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var sessionResp struct {
		SessionID string `json:"session_id"`
		Error     string `json:"error,omitempty"`
		ErrorID   string `json:"error_id,omitempty"`
	}

	if err := json.Unmarshal(bodyBytes, &sessionResp); err != nil {
		return "", fmt.Errorf("failed to decode session response: %w", err)
	}

	if sessionResp.Error != "" {
		if sessionResp.ErrorID != "" {
			return "", fmt.Errorf("SwarmUI session error (%s): %s", sessionResp.ErrorID, sessionResp.Error)
		}
		return "", fmt.Errorf("SwarmUI session error: %s", sessionResp.Error)
	}

	if sessionResp.SessionID == "" {
		return "", fmt.Errorf("session response did not contain session_id")
	}

	return sessionResp.SessionID, nil
}

// ListModels lists all available models
func (c *SwarmClient) ListModels() ([]Model, error) {
	return c.ListModelsWithOptions(ListModelsOptions{})
}

// ListModelsOptions configures the ListModels API call
type ListModelsOptions struct {
	Path        string // What folder path to search within. Use empty string for root.
	Depth       int    // Maximum depth (number of recursive folders) to search.
	Subtype     string // Model sub-type - LoRA, Wildcards, etc. Default: "Stable-Diffusion"
	SortBy      string // What to sort the list by - Name, DateCreated, or DateModified. Default: "Name"
	AllowRemote bool   // If true, allow remote models. If false, only local models. Default: true
	SortReverse bool   // If true, the sorting should be done in reverse. Default: false
	DataImages  bool   // If true, provide model images in raw data format. If false, use URLs. Default: false
}

// ListModelsWithOptions lists available models with specific options
func (c *SwarmClient) ListModelsWithOptions(options ListModelsOptions) ([]Model, error) {
	// Get session ID if we don't have one
	sessionID, err := c.ensureSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	endpoint := fmt.Sprintf("%s/API/ListModels", c.config.BaseURL)

	// Set defaults for required parameters
	if options.Subtype == "" {
		options.Subtype = "Stable-Diffusion"
	}
	if options.SortBy == "" {
		options.SortBy = "Name"
	}
	if options.Depth == 0 {
		options.Depth = 5 // Reasonable default depth
	}

	// Build request payload
	payload := map[string]interface{}{
		"session_id":  sessionID,
		"path":        options.Path,
		"depth":       options.Depth,
		"subtype":     options.Subtype,
		"sortBy":      options.SortBy,
		"allowRemote": options.AllowRemote,
		"sortReverse": options.SortReverse,
		"dataImages":  options.DataImages,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Request: POST %s\nPayload: %s\n", endpoint, string(payloadBytes))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if c.config.Verbose {
		fmt.Printf("Response Status: %d\nResponse Body: %s\n", resp.StatusCode, string(bodyBytes))
	}

	// Check for non-OK status and parse SwarmUI error format
	if resp.StatusCode != http.StatusOK {
		if swarmErr := parseSwarmUIError(bodyBytes); swarmErr != nil {
			return nil, swarmErr
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResp struct {
		Folders []string `json:"folders"`
		Files   []Model  `json:"files"`
		Error   string   `json:"error,omitempty"`
		ErrorID string   `json:"error_id,omitempty"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Handle SwarmUI-specific errors in successful responses
	if apiResp.Error != "" {
		// Handle session expiration
		if apiResp.ErrorID == "invalid_session_id" {
			// Clear session and retry once
			c.mu.Lock()
			c.sessionID = ""
			c.mu.Unlock()
			return c.ListModelsWithOptions(options)
		}
		if apiResp.ErrorID != "" {
			return nil, fmt.Errorf("SwarmUI error (%s): %s", apiResp.ErrorID, apiResp.Error)
		}
		return nil, fmt.Errorf("SwarmUI error: %s", apiResp.Error)
	}

	return apiResp.Files, nil
}

// GetModel gets details about a specific model
func (c *SwarmClient) GetModel(name string) (*Model, error) {
	// Get all models and find the specific one
	models, err := c.ListModels()
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	// Find the specific model by name
	for _, model := range models {
		if model.Name == name {
			return &model, nil
		}
	}

	return nil, fmt.Errorf("model '%s' not found", name)
}

// simulateProgress provides progress updates for HTTP-based generation
// This is a temporary solution until WebSocket support is implemented
// TODO: Replace with actual WebSocket implementation using GenerateText2ImageWS endpoint
func (c *SwarmClient) simulateProgress(sessionID string, callback ProgressCallback, done chan bool) {
	ticker := time.NewTicker(500 * time.Millisecond) // Update every 500ms
	defer ticker.Stop()

	progress := 0.1   // Start at 10%
	increment := 0.05 // Increase by 5% each tick

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			progress += increment
			if progress > 0.9 { // Cap at 90% until completion
				progress = 0.9
				increment = 0.01 // Slow down near completion
			}

			// Update session progress
			c.mu.Lock()
			if session, exists := c.sessions[sessionID]; exists {
				session.Progress = progress
				session.Status = "generating"
			}
			c.mu.Unlock()

			// Call progress callback
			callback(progress, "Generating...")
		}
	}
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
