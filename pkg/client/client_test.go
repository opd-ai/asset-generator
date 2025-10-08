package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewAssetClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				BaseURL: "http://localhost:7801",
				APIKey:  "test-key",
			},
			wantErr: false,
		},
		{
			name: "missing base URL",
			config: &Config{
				APIKey: "test-key",
			},
			wantErr: true,
		},
		{
			name: "valid without API key",
			config: &Config{
				BaseURL: "http://localhost:7801",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewAssetClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAssetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("Expected client to be non-nil")
			}
		})
	}
}

func TestListModels(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/API/ListModels" {
			t.Errorf("Expected path /API/ListModels, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"models": [
				{
					"name": "stable-diffusion-xl",
					"type": "text-to-image",
					"description": "SDXL model",
					"version": "1.0",
					"loaded": true
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	models, err := client.ListModels()
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}

	if len(models) != 1 {
		t.Errorf("Expected 1 model, got %d", len(models))
	}

	if models[0].Name != "stable-diffusion-xl" {
		t.Errorf("Expected model name 'stable-diffusion-xl', got '%s'", models[0].Name)
	}
}

func TestGenerateImage(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle GetNewSession call first
		if r.URL.Path == "/API/GetNewSession" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"session_id": "test-session-123"}`))
			return
		}

		if r.URL.Path != "/API/GenerateText2Image" {
			t.Errorf("Expected path /API/GenerateText2Image, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"images": ["/output/image1.png"],
			"info": {
				"seed": 12345,
				"steps": 20
			}
		}`))
	}))
	defer server.Close()

	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &GenerationRequest{
		Prompt: "test prompt",
		Parameters: map[string]interface{}{
			"batch_size": 1,
		},
	}

	result, err := client.GenerateImage(context.Background(), req)
	if err != nil {
		t.Fatalf("GenerateImage() error = %v", err)
	}

	if len(result.ImagePaths) != 1 {
		t.Errorf("Expected 1 image path, got %d", len(result.ImagePaths))
	}

	if result.ImagePaths[0] != "/output/image1.png" {
		t.Errorf("Expected image path '/output/image1.png', got '%s'", result.ImagePaths[0])
	}

	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", result.Status)
	}
}

func TestGenerateImageBatch(t *testing.T) {
	// Create a test server that returns multiple images
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/API/GetNewSession" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"session_id": "test-session-123"}`))
			return
		}

		if r.URL.Path != "/API/GenerateText2Image" {
			t.Errorf("Expected path /API/GenerateText2Image, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Return multiple images for batch generation
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"images": [
				"/output/image1.png",
				"/output/image2.png",
				"/output/image3.png"
			],
			"info": {
				"seed": 12345,
				"steps": 20,
				"batch_count": 3
			}
		}`))
	}))
	defer server.Close()

	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &GenerationRequest{
		Prompt: "test prompt",
		Parameters: map[string]interface{}{
			"batch_size": 3,
		},
	}

	result, err := client.GenerateImage(context.Background(), req)
	if err != nil {
		t.Fatalf("GenerateImage() error = %v", err)
	}

	// Verify we got 3 images
	if len(result.ImagePaths) != 3 {
		t.Errorf("Expected 3 image paths for batch_size=3, got %d", len(result.ImagePaths))
	}

	// Verify the image paths are correct
	expectedPaths := []string{
		"/output/image1.png",
		"/output/image2.png",
		"/output/image3.png",
	}
	for i, expectedPath := range expectedPaths {
		if result.ImagePaths[i] != expectedPath {
			t.Errorf("Expected image path %d to be '%s', got '%s'", i, expectedPath, result.ImagePaths[i])
		}
	}

	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", result.Status)
	}
}

func TestGenerateImageWithContext(t *testing.T) {
	// Create a server that handles both GetNewSession and generation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/API/GetNewSession" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"session_id": "test-session-123"}`))
			return
		}
		// Wait a bit before responding to generation
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"images": [], "info": {}}`))
	}))
	defer server.Close()

	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := &GenerationRequest{
		Prompt: "test prompt",
		Parameters: map[string]interface{}{
			"batch_size": 1,
		},
	}

	_, err = client.GenerateImage(ctx, req)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}
}
