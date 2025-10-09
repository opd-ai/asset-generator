package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

// AssetClient is the main client for interacting with asset generation APIs
type AssetClient struct {
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

// Model represents an asset generation model
type Model struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Loaded      bool   `json:"loaded"`
}

// NewAssetClient creates a new asset generation API client
func NewAssetClient(config *Config) (*AssetClient, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	return &AssetClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 10 * time.Minute, // Extended timeout for Flux generation (can take 5-10 minutes)
		},
		sessions: make(map[string]*GenerationSession),
	}, nil
}

// GetNewSession gets a new session ID from the asset generation API
func (c *AssetClient) GetNewSession(ctx context.Context) (string, error) {
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

// GenerateImage generates an image using the asset generation API
func (c *AssetClient) GenerateImage(ctx context.Context, req *GenerationRequest) (*GenerationResult, error) {
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
	}

	// Handle images count (batch size)
	if images, ok := req.Parameters["images"]; ok && images != nil {
		if img, isInt := images.(int); isInt && img > 0 {
			body["images"] = img
		} else {
			body["images"] = 1 // Default to 1 if invalid
		}
	} else {
		body["images"] = 1 // Default to 1 image
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

	// Handle negative prompt explicitly (only include if non-empty)
	if negPrompt, ok := req.Parameters["negative_prompt"]; ok {
		if negPromptStr, isString := negPrompt.(string); isString && negPromptStr != "" {
			body["negative_prompt"] = negPromptStr
		}
	}

	// Add any other parameters from the request
	for k, v := range req.Parameters {
		// Skip parameters we've already handled explicitly
		if k != "batch_size" && k != "width" && k != "height" && k != "cfgscale" && k != "steps" && k != "seed" && k != "negative_prompt" {
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
		// Handle session expiration - retry with new session
		if apiResp.ErrorID == "invalid_session_id" {
			// Clear expired session
			c.mu.Lock()
			oldSessionID := c.sessionID
			c.sessionID = ""
			c.mu.Unlock()

			// Only retry if we actually had a cached session
			// This prevents infinite recursion if the session was already cleared
			if oldSessionID != "" {
				return c.GenerateImage(ctx, req)
			}
		}

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

// GenerateImageWS generates an image using WebSocket for real-time progress updates.
// This connects to the GenerateText2ImageWS endpoint for live progress information.
// Falls back to HTTP GenerateImage() if WebSocket connection fails.
//
// WebSocket provides authentic progress updates from SwarmUI instead of simulated progress.
// This is particularly beneficial for long-running generations (e.g., Flux models: 5-10 minutes).
func (c *AssetClient) GenerateImageWS(ctx context.Context, req *GenerationRequest) (*GenerationResult, error) {
	// Get or create session ID
	sessionID, err := c.ensureSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Build WebSocket URL (convert http:// to ws:// or https:// to wss://)
	// SwarmUI WebSocket endpoint: ws://host/API/GenerateText2ImageWS
	wsURL := c.config.BaseURL
	if len(wsURL) > 7 && wsURL[:7] == "http://" {
		wsURL = "ws://" + wsURL[7:]
	} else if len(wsURL) > 8 && wsURL[:8] == "https://" {
		wsURL = "wss://" + wsURL[8:]
	}
	wsURL += "/API/GenerateText2ImageWS"

	// Build request body with correct SwarmUI parameter names
	body := map[string]interface{}{
		"session_id": sessionID,
		"prompt":     req.Prompt,
		"images":     1, // Default to 1 image
	}

	// Override images count if specified in parameters
	if images, ok := req.Parameters["images"]; ok && images != nil {
		if img, isInt := images.(int); isInt && img > 0 {
			body["images"] = img
		}
	}

	// Add model if specified
	if req.Model != "" {
		body["model"] = req.Model
	}

	// Add all other parameters
	for key, value := range req.Parameters {
		if key != "images" { // Already handled above
			body[key] = value
		}
	}

	// Connect to WebSocket
	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		// Fallback to HTTP if WebSocket fails (e.g., server doesn't support WS, network issues)
		// This ensures backward compatibility and graceful degradation
		if c.config.Verbose {
			fmt.Printf("WebSocket connection failed, falling back to HTTP: %v\n", err)
		}
		return c.GenerateImage(ctx, req)
	}
	defer conn.Close()

	// Send initial request - SwarmUI expects same JSON format as HTTP endpoint
	if err := conn.WriteJSON(body); err != nil {
		return nil, fmt.Errorf("failed to send WebSocket request: %w", err)
	}

	// Create session tracking
	session := &GenerationSession{
		ID:        sessionID,
		Status:    "generating",
		Progress:  0.0,
		StartTime: time.Now(),
	}

	c.mu.Lock()
	c.sessions[sessionID] = session
	c.mu.Unlock()

	// Listen for progress updates
	// SwarmUI sends multiple JSON messages over the WebSocket connection:
	// 1. Progress updates: {"progress": 0.45, "status": "generating"}
	// 2. Final result: {"images": ["path/to/image.png"], "info": {...}}
	var finalResult *GenerationResult
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Read message from WebSocket
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				// Check if it's a normal close (generation complete)
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					break
				}
				return nil, fmt.Errorf("WebSocket read error: %w", err)
			}

			// Handle error messages from SwarmUI
			if errMsg, ok := msg["error"].(string); ok && errMsg != "" {
				// Handle session expiration - retry with new session (same as HTTP behavior)
				if errID, hasErrID := msg["error_id"].(string); hasErrID && errID == "invalid_session_id" {
					c.mu.Lock()
					oldSessionID := c.sessionID
					c.sessionID = ""
					c.mu.Unlock()

					// Only retry if we had a cached session (prevents infinite recursion)
					if oldSessionID != "" {
						return c.GenerateImageWS(ctx, req)
					}
				}
				return nil, fmt.Errorf("SwarmUI error: %s", errMsg)
			}

			// Parse real-time progress updates from SwarmUI
			// Unlike simulated progress, these reflect actual generation state
			if progress, ok := msg["progress"].(float64); ok {
				c.mu.Lock()
				session.Progress = progress
				c.mu.Unlock()

				if req.ProgressCallback != nil {
					status := "Generating..."
					if statusStr, ok := msg["status"].(string); ok {
						status = statusStr
					}
					req.ProgressCallback(progress, status)
				}
			}

			// Check for completion (images field indicates generation is complete)
			// This is the final message from SwarmUI containing the result
			if images, ok := msg["images"].([]interface{}); ok && len(images) > 0 {
				// Convert images to string slice
				imagePaths := make([]string, len(images))
				for i, img := range images {
					if imgStr, ok := img.(string); ok {
						imagePaths[i] = imgStr
					}
				}

				// Extract metadata if present
				metadata := make(map[string]interface{})
				if info, ok := msg["info"].(map[string]interface{}); ok {
					metadata = info
				}

				finalResult = &GenerationResult{
					SessionID:  sessionID,
					ImagePaths: imagePaths,
					Metadata:   metadata,
					Status:     "completed",
					CreatedAt:  time.Now(),
				}

				// Update session
				c.mu.Lock()
				session.Status = "completed"
				session.Progress = 1.0
				session.Result = finalResult
				c.mu.Unlock()

				// Report completion
				if req.ProgressCallback != nil {
					req.ProgressCallback(1.0, "Generation completed")
				}

				break
			}
		}

		// Break if we got final result
		if finalResult != nil {
			break
		}
	}

	if finalResult == nil {
		return nil, fmt.Errorf("WebSocket closed without returning images")
	}

	return finalResult, nil
}

// cleanupSession removes a session from memory to prevent memory leaks
func (c *AssetClient) cleanupSession(sessionID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.sessions, sessionID)
}

// cleanupOldSessions removes sessions older than the specified duration
func (c *AssetClient) cleanupOldSessions(maxAge time.Duration) {
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
func (c *AssetClient) ensureSession() (string, error) {
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

// getNewSession gets a new session ID from the asset generation service
func (c *AssetClient) getNewSession() (string, error) {
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
func (c *AssetClient) ListModels() ([]Model, error) {
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
func (c *AssetClient) ListModelsWithOptions(options ListModelsOptions) ([]Model, error) {
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
func (c *AssetClient) GetModel(name string) (*Model, error) {
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
func (c *AssetClient) simulateProgress(sessionID string, callback ProgressCallback, done chan bool) {
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

// DownloadImages downloads generated images from the server and saves them to the specified directory.
// imagePaths should be the paths returned by the generation API (e.g., "View/local/raw/2024-05-19/file.png")
// outputDir is the local directory where images will be saved.
// Returns a slice of local file paths where images were saved.
func (c *AssetClient) DownloadImages(ctx context.Context, imagePaths []string, outputDir string) ([]string, error) {
	if len(imagePaths) == 0 {
		return nil, fmt.Errorf("no images to download")
	}

	// Create output directory if it doesn't exist
	if err := ensureDir(outputDir); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	savedPaths := make([]string, 0, len(imagePaths))
	var downloadErrors []error

	for i, imagePath := range imagePaths {
		// Build full image URL
		// SwarmUI returns paths like "View/local/raw/2024-05-19/filename.png"
		imageURL := fmt.Sprintf("%s/%s", c.config.BaseURL, imagePath)

		// Extract filename from path
		parts := strings.Split(imagePath, "/")
		if len(parts) == 0 {
			downloadErrors = append(downloadErrors, fmt.Errorf("invalid image path: %s", imagePath))
			continue
		}
		filename := parts[len(parts)-1]

		// Create output file path
		outputPath := fmt.Sprintf("%s/%s", outputDir, filename)

		// Download the image
		if err := c.downloadFile(ctx, imageURL, outputPath); err != nil {
			downloadErrors = append(downloadErrors, fmt.Errorf("failed to download image %d (%s): %w", i+1, filename, err))
			continue
		}

		savedPaths = append(savedPaths, outputPath)

		if c.config.Verbose {
			fmt.Printf("Downloaded: %s -> %s\n", imageURL, outputPath)
		}
	}

	// Return error if any downloads failed
	if len(downloadErrors) > 0 {
		// If some succeeded and some failed, return partial success with error
		if len(savedPaths) > 0 {
			return savedPaths, fmt.Errorf("partial download failure: %d/%d images downloaded successfully; errors: %v",
				len(savedPaths), len(imagePaths), downloadErrors)
		}
		return nil, fmt.Errorf("all downloads failed: %v", downloadErrors)
	}

	return savedPaths, nil
}

// downloadFile downloads a file from the given URL and saves it to the specified path
func (c *AssetClient) downloadFile(ctx context.Context, url, filepath string) error {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header if API key is set
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Create output file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Copy response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ensureDir creates a directory if it doesn't exist
func ensureDir(dir string) error {
	// Check if directory exists
	info, err := os.Stat(dir)
	if err == nil {
		// Directory exists, check if it's actually a directory
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", dir)
		}
		return nil
	}

	// If error is not "not exists", return it
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory: %w", err)
	}

	// Create directory with permissions 0755
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// Close closes any open connections
func (c *AssetClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.wsConn != nil {
		return c.wsConn.Close()
	}

	return nil
}
