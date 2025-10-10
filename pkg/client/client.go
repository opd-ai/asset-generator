package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opd-ai/asset-generator/pkg/processor"
)

const (
	// stateFileName is the name of the file used to persist generation state
	stateFileName = ".asset-generator-state.json"
	// stateMaxAge is the maximum age of sessions to keep in the state file
	stateMaxAge = 24 * time.Hour
)

// persistedState represents the structure of the state file
type persistedState struct {
	Sessions  map[string]*PersistedSession `json:"sessions"`
	UpdatedAt time.Time                    `json:"updated_at"`
}

// PersistedSession is the serializable version of GenerationSession
type PersistedSession struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Progress  float64   `json:"progress"`
	StartTime time.Time `json:"start_time"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Config holds the client configuration
type Config struct {
	BaseURL string
	APIKey  string
	Verbose bool
}

// AssetClient is the main client for interacting with asset generation APIs
type AssetClient struct {
	config        *Config
	httpClient    *http.Client
	wsConn        *websocket.Conn // Reserved for future WebSocket implementation
	mu            sync.RWMutex
	sessions      map[string]*GenerationSession
	sessionID     string // Current session ID for API calls
	stateFilePath string // Path to the persistent state file
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

	// Determine state file path (current working directory)
	cwd, err := os.Getwd()
	if err != nil {
		// If we can't get cwd, use temp directory as fallback
		cwd = os.TempDir()
	}
	stateFilePath := filepath.Join(cwd, stateFileName)

	client := &AssetClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 40 * time.Minute, // Extended timeout for Flux generation (can take up to 40 minutes for complex generations)
		},
		sessions:      make(map[string]*GenerationSession),
		stateFilePath: stateFilePath,
	}

	// Load existing state from file
	client.loadStateFromFile()

	// Clean up old sessions
	client.cleanupOldPersistedSessions()

	return client, nil
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

	// Save initial state to file
	c.saveStateToFile()

	// Ensure session cleanup on function exit (success or error)
	defer c.removeSessionState(sessionID)

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
		c.mu.Lock()
		session.Progress = 0.0
		session.Status = "starting"
		c.mu.Unlock()
		c.saveStateToFile()
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
		HandshakeTimeout: 10 * time.Minute, // 10 minutes for WebSocket handshake
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

	// Save initial state to file
	c.saveStateToFile()

	// Ensure session cleanup on function exit
	defer c.removeSessionState(sessionID)

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

				// Persist progress to file
				c.saveStateToFile()

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

// loadStateFromFile loads persisted generation sessions from the state file
func (c *AssetClient) loadStateFromFile() {
	data, err := os.ReadFile(c.stateFilePath)
	if err != nil {
		// File doesn't exist or can't be read - this is fine for first run
		if c.config.Verbose && !os.IsNotExist(err) {
			fmt.Printf("Could not read state file: %v\n", err)
		}
		return
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		if c.config.Verbose {
			fmt.Printf("Could not parse state file: %v\n", err)
		}
		return
	}

	// Load persisted sessions into memory
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, ps := range state.Sessions {
		// Only load sessions that are still active
		if ps.Status == "generating" || ps.Status == "starting" || ps.Status == "pending" {
			c.sessions[id] = &GenerationSession{
				ID:        ps.ID,
				Status:    ps.Status,
				Progress:  ps.Progress,
				StartTime: ps.StartTime,
			}
		}
	}

	if c.config.Verbose && len(state.Sessions) > 0 {
		fmt.Printf("Loaded %d session(s) from state file\n", len(state.Sessions))
	}
}

// saveStateToFile persists current generation sessions to the state file
func (c *AssetClient) saveStateToFile() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := persistedState{
		Sessions:  make(map[string]*PersistedSession),
		UpdatedAt: time.Now(),
	}

	// Convert in-memory sessions to persistable format
	for id, session := range c.sessions {
		// Only persist active sessions
		if session.Status == "generating" || session.Status == "starting" || session.Status == "pending" {
			state.Sessions[id] = &PersistedSession{
				ID:        session.ID,
				Status:    session.Status,
				Progress:  session.Progress,
				StartTime: session.StartTime,
				UpdatedAt: time.Now(),
			}
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write to file (use temp file + rename for atomicity)
	tempPath := c.stateFilePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	if err := os.Rename(tempPath, c.stateFilePath); err != nil {
		os.Remove(tempPath) // Clean up temp file on error
		return fmt.Errorf("failed to rename state file: %w", err)
	}

	if c.config.Verbose && len(state.Sessions) > 0 {
		fmt.Printf("Saved %d session(s) to state file\n", len(state.Sessions))
	}

	return nil
}

// cleanupOldPersistedSessions removes old sessions from the state file
func (c *AssetClient) cleanupOldPersistedSessions() {
	data, err := os.ReadFile(c.stateFilePath)
	if err != nil {
		return // File doesn't exist or can't be read - nothing to clean
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		return // Can't parse - leave it alone
	}

	// Remove old sessions
	cutoff := time.Now().Add(-stateMaxAge)
	needsUpdate := false

	for id, session := range state.Sessions {
		// Remove if too old or if already completed/failed
		if session.StartTime.Before(cutoff) ||
			(session.Status != "generating" && session.Status != "starting" && session.Status != "pending") {
			delete(state.Sessions, id)
			needsUpdate = true
		}
	}

	// Save updated state if we removed anything
	if needsUpdate {
		state.UpdatedAt = time.Now()
		data, err := json.MarshalIndent(state, "", "  ")
		if err == nil {
			os.WriteFile(c.stateFilePath, data, 0644)
		}

		if c.config.Verbose {
			fmt.Printf("Cleaned up old sessions from state file\n")
		}
	}
}

// updateSessionState updates a session in both memory and persistent state
func (c *AssetClient) updateSessionState(sessionID, status string, progress float64) error {
	c.mu.Lock()
	if session, exists := c.sessions[sessionID]; exists {
		session.Status = status
		session.Progress = progress
	}
	c.mu.Unlock()

	// Persist to file
	return c.saveStateToFile()
}

// removeSessionState removes a session from both memory and persistent state
func (c *AssetClient) removeSessionState(sessionID string) error {
	c.mu.Lock()
	delete(c.sessions, sessionID)
	c.mu.Unlock()

	// Persist to file
	return c.saveStateToFile()
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

// Interrupt cancels the current generation in progress
func (c *AssetClient) Interrupt(ctx context.Context) error {
	// Get session ID
	sessionID, err := c.ensureSession()
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	endpoint := fmt.Sprintf("%s/API/InterruptGeneration", c.config.BaseURL)

	payload := map[string]interface{}{
		"session_id": sessionID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Request: POST %s\n", endpoint)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var apiResp struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
		ErrorID string `json:"error_id,omitempty"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		// If we can't parse the response, assume success if status was OK
		if c.config.Verbose {
			fmt.Printf("Warning: Failed to parse interrupt response: %v\n", err)
		}
		return nil
	}

	if apiResp.Error != "" {
		// Handle session expiration
		if apiResp.ErrorID == "invalid_session_id" {
			// Clear session and retry once
			c.mu.Lock()
			c.sessionID = ""
			c.mu.Unlock()
			return c.Interrupt(ctx)
		}
		if apiResp.ErrorID != "" {
			return fmt.Errorf("SwarmUI error (%s): %s", apiResp.ErrorID, apiResp.Error)
		}
		return fmt.Errorf("SwarmUI error: %s", apiResp.Error)
	}

	return nil
}

// InterruptAll cancels all queued generations
func (c *AssetClient) InterruptAll(ctx context.Context) error {
	// Get session ID
	sessionID, err := c.ensureSession()
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	endpoint := fmt.Sprintf("%s/API/InterruptAll", c.config.BaseURL)

	payload := map[string]interface{}{
		"session_id": sessionID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Request: POST %s\n", endpoint)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var apiResp struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
		ErrorID string `json:"error_id,omitempty"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		// If we can't parse the response, assume success if status was OK
		if c.config.Verbose {
			fmt.Printf("Warning: Failed to parse interrupt response: %v\n", err)
		}
		return nil
	}

	if apiResp.Error != "" {
		// Handle session expiration
		if apiResp.ErrorID == "invalid_session_id" {
			// Clear session and retry once
			c.mu.Lock()
			c.sessionID = ""
			c.mu.Unlock()
			return c.InterruptAll(ctx)
		}
		if apiResp.ErrorID != "" {
			return fmt.Errorf("SwarmUI error (%s): %s", apiResp.ErrorID, apiResp.Error)
		}
		return fmt.Errorf("SwarmUI error: %s", apiResp.Error)
	}

	return nil
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

// DownloadOptions holds options for downloading images
type DownloadOptions struct {
	OutputDir        string                 // Directory to save images
	FilenameTemplate string                 // Template for generating filenames (e.g., "image-{index}.png")
	Metadata         map[string]interface{} // Metadata for template variables

	// Postprocessing options - applied locally after download
	// Auto-crop (runs first, before downscaling)
	AutoCrop               bool  // Enable automatic cropping of whitespace borders
	AutoCropThreshold      uint8 // Whitespace detection threshold (0-255, default: 250)
	AutoCropTolerance      uint8 // Tolerance for near-white colors (0-255, default: 10)
	AutoCropPreserveAspect bool  // Preserve original aspect ratio when cropping

	// Downscaling (runs after auto-crop if enabled)
	DownscaleWidth      int     // Target width for downscaling (0 means auto-calculate from height)
	DownscaleHeight     int     // Target height for downscaling (0 means auto-calculate from width)
	DownscalePercentage float64 // Scale by percentage (1-100, takes precedence over Width/Height if > 0)
	DownscaleFilter     string  // Downscaling algorithm: "lanczos" (default), "bilinear", "nearest"
	JPEGQuality         int     // JPEG quality for downscaled images (1-100, default: 90)
}

// DownloadImages downloads generated images from the server and saves them to the specified directory.
// imagePaths should be the paths returned by the generation API (e.g., "View/local/raw/2024-05-19/file.png")
// outputDir is the local directory where images will be saved.
// Returns a slice of local file paths where images were saved.
func (c *AssetClient) DownloadImages(ctx context.Context, imagePaths []string, outputDir string) ([]string, error) {
	return c.DownloadImagesWithOptions(ctx, imagePaths, &DownloadOptions{
		OutputDir: outputDir,
	})
}

// DownloadImagesWithOptions downloads generated images with custom filename options.
func (c *AssetClient) DownloadImagesWithOptions(ctx context.Context, imagePaths []string, opts *DownloadOptions) ([]string, error) {
	if len(imagePaths) == 0 {
		return nil, fmt.Errorf("no images to download")
	}

	if opts == nil {
		opts = &DownloadOptions{OutputDir: "."}
	}

	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = "."
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

		// Extract original filename and extension from path
		parts := strings.Split(imagePath, "/")
		if len(parts) == 0 {
			downloadErrors = append(downloadErrors, fmt.Errorf("invalid image path: %s", imagePath))
			continue
		}
		originalFilename := parts[len(parts)-1]

		// Determine the filename to use
		var filename string
		if opts.FilenameTemplate != "" {
			// Generate filename from template
			filename = generateFilename(opts.FilenameTemplate, i, originalFilename, opts.Metadata)
		} else {
			// Use original filename
			filename = originalFilename
		}

		// Create output file path
		outputPath := fmt.Sprintf("%s/%s", outputDir, filename)

		// Download the image
		if err := c.downloadFile(ctx, imageURL, outputPath); err != nil {
			downloadErrors = append(downloadErrors, fmt.Errorf("failed to download image %d (%s): %w", i+1, filename, err))
			continue
		}

		// Apply postprocessing pipeline (order matters: crop first, then downscale)

		// Step 1: Auto-crop if enabled (removes whitespace before downscaling)
		if opts.AutoCrop {
			if err := c.applyAutoCrop(outputPath, opts); err != nil {
				downloadErrors = append(downloadErrors, fmt.Errorf("failed to auto-crop image %d (%s): %w", i+1, filename, err))
				continue
			}
		}

		// Step 2: Downscale if options are set
		if opts.DownscaleWidth > 0 || opts.DownscaleHeight > 0 {
			if err := c.applyDownscale(outputPath, opts); err != nil {
				downloadErrors = append(downloadErrors, fmt.Errorf("failed to downscale image %d (%s): %w", i+1, filename, err))
				continue
			}
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

	// Mandatory: Strip all PNG metadata from downloaded images
	// This removes any generation parameters, timestamps, or other sensitive data
	// that may have been embedded by the generation API
	if err := processor.StripPNGMetadata(filepath); err != nil {
		return fmt.Errorf("failed to strip PNG metadata: %w", err)
	}

	return nil
}

// generateFilename creates a filename from a template with variable substitution.
// Supported placeholders:
// - {index} or {i}: Zero-padded index (e.g., 001, 002)
// - {index1} or {i1}: One-based index (e.g., 1, 2, 3)
// - {timestamp} or {ts}: Unix timestamp
// - {datetime} or {dt}: Formatted datetime (YYYY-MM-DD_HH-MM-SS)
// - {date}: Date only (YYYY-MM-DD)
// - {time}: Time only (HH-MM-SS)
// - {original}: Original filename from server
// - {ext}: Original file extension (including dot)
// - {seed}: Seed value from metadata
// - {model}: Model name from metadata
// - {width}: Image width from metadata
// - {height}: Image height from metadata
// - {prompt}: Truncated prompt (first 50 chars, sanitized)
func generateFilename(template string, index int, originalFilename string, metadata map[string]interface{}) string {
	result := template
	now := time.Now()

	// Extract extension from original filename
	ext := ""
	if dotIdx := strings.LastIndex(originalFilename, "."); dotIdx != -1 {
		ext = originalFilename[dotIdx:]
	}

	// Zero-padded index (assuming max 999 images)
	result = strings.ReplaceAll(result, "{index}", fmt.Sprintf("%03d", index))
	result = strings.ReplaceAll(result, "{i}", fmt.Sprintf("%03d", index))

	// One-based index
	result = strings.ReplaceAll(result, "{index1}", fmt.Sprintf("%d", index+1))
	result = strings.ReplaceAll(result, "{i1}", fmt.Sprintf("%d", index+1))

	// Timestamps
	result = strings.ReplaceAll(result, "{timestamp}", fmt.Sprintf("%d", now.Unix()))
	result = strings.ReplaceAll(result, "{ts}", fmt.Sprintf("%d", now.Unix()))
	result = strings.ReplaceAll(result, "{datetime}", now.Format("2006-01-02_15-04-05"))
	result = strings.ReplaceAll(result, "{dt}", now.Format("2006-01-02_15-04-05"))
	result = strings.ReplaceAll(result, "{date}", now.Format("2006-01-02"))
	result = strings.ReplaceAll(result, "{time}", now.Format("15-04-05"))

	// Original filename and extension
	result = strings.ReplaceAll(result, "{original}", originalFilename)
	result = strings.ReplaceAll(result, "{ext}", ext)

	// Metadata-based replacements
	if metadata != nil {
		if seed, ok := metadata["seed"]; ok {
			result = strings.ReplaceAll(result, "{seed}", fmt.Sprintf("%v", seed))
		}
		if model, ok := metadata["model"]; ok {
			result = strings.ReplaceAll(result, "{model}", fmt.Sprintf("%v", model))
		}
		if width, ok := metadata["width"]; ok {
			result = strings.ReplaceAll(result, "{width}", fmt.Sprintf("%v", width))
		}
		if height, ok := metadata["height"]; ok {
			result = strings.ReplaceAll(result, "{height}", fmt.Sprintf("%v", height))
		}
		if prompt, ok := metadata["prompt"]; ok {
			// Sanitize and truncate prompt for filename
			promptStr := fmt.Sprintf("%v", prompt)
			promptStr = sanitizeForFilename(promptStr)
			if len(promptStr) > 50 {
				promptStr = promptStr[:50]
			}
			result = strings.ReplaceAll(result, "{prompt}", promptStr)
		}
	}

	// If no extension in template but we have one, append it
	if !strings.Contains(result, ".") && ext != "" {
		result += ext
	}

	return result
}

// sanitizeForFilename removes or replaces characters that are invalid in filenames
func sanitizeForFilename(s string) string {
	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")

	// Remove or replace invalid filename characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t"}
	for _, char := range invalidChars {
		s = strings.ReplaceAll(s, char, "")
	}

	// Remove multiple consecutive underscores
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}

	// Trim leading/trailing underscores
	s = strings.Trim(s, "_")

	return s
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

// applyAutoCrop applies local postprocessing auto-crop to an image file.
// This is called after downloading to remove whitespace borders before downscaling.
func (c *AssetClient) applyAutoCrop(imagePath string, opts *DownloadOptions) error {
	threshold := opts.AutoCropThreshold
	if threshold == 0 {
		threshold = 250 // default: very light colors are whitespace
	}

	tolerance := opts.AutoCropTolerance
	if tolerance == 0 {
		tolerance = 10 // default tolerance
	}

	quality := opts.JPEGQuality
	if quality == 0 {
		quality = 90 // default JPEG quality
	}

	if c.config.Verbose {
		fmt.Printf("Auto-cropping image: %s (threshold: %d, tolerance: %d, preserve aspect: %v)\n",
			imagePath, threshold, tolerance, opts.AutoCropPreserveAspect)
	}

	// Build crop options
	cropOpts := processor.CropOptions{
		Threshold:           threshold,
		Tolerance:           tolerance,
		JPEGQuality:         quality,
		PreserveAspectRatio: opts.AutoCropPreserveAspect,
	}

	// Apply auto-crop in-place (replaces the original file)
	if err := processor.AutoCropInPlace(imagePath, cropOpts); err != nil {
		return fmt.Errorf("auto-crop operation failed: %w", err)
	}

	return nil
}

// applyDownscale applies local postprocessing downscaling to an image file.
// This is called after downloading to reduce image size using high-quality Lanczos filtering.
func (c *AssetClient) applyDownscale(imagePath string, opts *DownloadOptions) error {
	// Map filter name to processor filter type
	filter := "lanczos" // default
	if opts.DownscaleFilter != "" {
		filter = strings.ToLower(opts.DownscaleFilter)
	}

	// Validate and map filter to processor.ResizeFilter
	var filterType processor.ResizeFilter
	switch filter {
	case "lanczos":
		filterType = processor.FilterLanczos
	case "bilinear":
		filterType = processor.FilterBiLinear
	case "nearest":
		filterType = processor.FilterNearestNeighbor
	default:
		return fmt.Errorf("invalid downscale filter: %s (valid options: lanczos, bilinear, nearest)", filter)
	}

	quality := opts.JPEGQuality
	if quality == 0 {
		quality = 90 // default JPEG quality
	}

	if c.config.Verbose {
		if opts.DownscalePercentage > 0 {
			fmt.Printf("Downscaling image: %s (scale: %.0f%%, filter: %s)\n",
				imagePath, opts.DownscalePercentage, filter)
		} else {
			fmt.Printf("Downscaling image: %s (target: %dx%d, filter: %s)\n",
				imagePath, opts.DownscaleWidth, opts.DownscaleHeight, filter)
		}
	}

	// Build downscale options
	downscaleOpts := processor.DownscaleOptions{
		Width:       opts.DownscaleWidth,
		Height:      opts.DownscaleHeight,
		Percentage:  opts.DownscalePercentage,
		Filter:      filterType,
		JPEGQuality: quality,
	}

	// Apply downscaling in-place (replaces the original file)
	if err := processor.DownscaleInPlace(imagePath, downscaleOpts); err != nil {
		return fmt.Errorf("downscale operation failed: %w", err)
	}

	return nil
}

// ServerStatus represents the status of the SwarmUI server
type ServerStatus struct {
	ServerURL          string                 `json:"server_url"`
	Status             string                 `json:"status"`
	ResponseTime       string                 `json:"response_time"`
	Version            string                 `json:"version,omitempty"`
	SessionID          string                 `json:"session_id,omitempty"`
	Backends           []BackendStatus        `json:"backends,omitempty"`
	ModelsCount        int                    `json:"models_count"`
	ModelsLoaded       int                    `json:"models_loaded"`
	SystemInfo         map[string]interface{} `json:"system_info,omitempty"`
	ActiveGenerations  []ActiveGeneration     `json:"active_generations,omitempty"`
	GenerationsRunning int                    `json:"generations_running"`
}

// ActiveGeneration represents a currently running generation
type ActiveGeneration struct {
	SessionID string    `json:"session_id"`
	Status    string    `json:"status"`
	Progress  float64   `json:"progress"`
	StartTime time.Time `json:"start_time"`
	Duration  string    `json:"duration"`
}

// BackendStatus represents the status of a single backend
type BackendStatus struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	ModelLoaded string `json:"model_loaded,omitempty"`
	GPU         string `json:"gpu,omitempty"`
}

// GetActiveGenerations returns a list of currently running generation sessions
func (c *AssetClient) GetActiveGenerations() []ActiveGeneration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	activeGens := make([]ActiveGeneration, 0)

	for _, session := range c.sessions {
		// Only include sessions that are actively generating
		if session.Status == "generating" || session.Status == "starting" || session.Status == "pending" {
			duration := time.Since(session.StartTime)
			activeGens = append(activeGens, ActiveGeneration{
				SessionID: session.ID,
				Status:    session.Status,
				Progress:  session.Progress,
				StartTime: session.StartTime,
				Duration:  formatDuration(duration),
			})
		}
	}

	return activeGens
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// GetServerStatus queries the SwarmUI server for its current status
func (c *AssetClient) GetServerStatus(ctx context.Context) (*ServerStatus, error) {
	status := &ServerStatus{
		ServerURL: c.config.BaseURL,
		Status:    "unknown",
	}

	// Try to get a session to verify connectivity
	sessionID, sessionErr := c.GetNewSession(ctx)
	if sessionErr == nil {
		status.SessionID = sessionID
		status.Status = "online"
	} else {
		// Server might be offline or have issues
		status.Status = "offline"
		return status, fmt.Errorf("server unreachable: %w", sessionErr)
	}

	// Try to get backend status information
	// SwarmUI exposes backend info through the ListBackends API
	backendInfo, backendErr := c.getBackendStatus(ctx, sessionID)
	if backendErr == nil && backendInfo != nil {
		status.Backends = backendInfo.Backends
		status.SystemInfo = backendInfo.SystemInfo
		status.Version = backendInfo.Version
	}

	// Get model information
	models, modelsErr := c.ListModels()
	if modelsErr == nil {
		status.ModelsCount = len(models)
		// Count loaded models
		loadedCount := 0
		for _, model := range models {
			if model.Loaded {
				loadedCount++
			}
		}
		status.ModelsLoaded = loadedCount
	}

	// Get active generation information from local tracking
	activeGens := c.GetActiveGenerations()
	status.ActiveGenerations = activeGens
	status.GenerationsRunning = len(activeGens)

	// If we don't have local session tracking, infer from backend status
	// A backend with status "running" likely indicates an active generation
	if len(activeGens) == 0 && len(status.Backends) > 0 {
		for _, backend := range status.Backends {
			if backend.Status == "running" || backend.Status == "generating" {
				status.GenerationsRunning++
			}
		}
	}

	return status, nil
}

// backendInfo holds backend status information from SwarmUI
type backendInfo struct {
	Backends   []BackendStatus        `json:"backends"`
	Version    string                 `json:"version"`
	SystemInfo map[string]interface{} `json:"system_info"`
}

// getBackendStatus queries SwarmUI for backend status information
func (c *AssetClient) getBackendStatus(ctx context.Context, sessionID string) (*backendInfo, error) {
	endpoint := fmt.Sprintf("%s/API/ListBackends", c.config.BaseURL)

	payload := map[string]interface{}{
		"session_id": sessionID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	if c.config.Verbose {
		fmt.Printf("Request: POST %s\n", endpoint)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Backend listing might not be available in all SwarmUI versions
		// This is not a critical error, so we'll return nil instead
		if c.config.Verbose {
			fmt.Printf("Backend status unavailable (status %d): %s\n", resp.StatusCode, string(bodyBytes))
		}
		return nil, nil
	}

	// Parse the response
	var apiResp struct {
		Backends   []map[string]interface{} `json:"backends"`
		Version    string                   `json:"version"`
		SystemInfo map[string]interface{}   `json:"system_info"`
		Error      string                   `json:"error,omitempty"`
		ErrorID    string                   `json:"error_id,omitempty"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		// If we can't parse the response, it's not critical
		if c.config.Verbose {
			fmt.Printf("Failed to parse backend response: %v\n", err)
		}
		return nil, nil
	}

	if apiResp.Error != "" {
		// Backend errors are not critical for status command
		if c.config.Verbose {
			fmt.Printf("Backend API error: %s\n", apiResp.Error)
		}
		return nil, nil
	}

	// Convert backend data to structured format
	info := &backendInfo{
		Version:    apiResp.Version,
		SystemInfo: apiResp.SystemInfo,
		Backends:   make([]BackendStatus, 0, len(apiResp.Backends)),
	}

	for _, b := range apiResp.Backends {
		backend := BackendStatus{}

		if id, ok := b["backend_id"].(string); ok {
			backend.ID = id
		} else if id, ok := b["id"].(string); ok {
			backend.ID = id
		}

		if bType, ok := b["type"].(string); ok {
			backend.Type = bType
		}

		if status, ok := b["status"].(string); ok {
			backend.Status = status
		}

		if model, ok := b["model_loaded"].(string); ok {
			backend.ModelLoaded = model
		} else if model, ok := b["current_model"].(string); ok {
			backend.ModelLoaded = model
		}

		if gpu, ok := b["gpu"].(string); ok {
			backend.GPU = gpu
		} else if gpu, ok := b["gpu_id"].(string); ok {
			backend.GPU = gpu
		}

		info.Backends = append(info.Backends, backend)
	}

	return info, nil
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
