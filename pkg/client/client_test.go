package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSwarmClient(t *testing.T) {
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
			client, err := NewSwarmClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSwarmClient() error = %v, wantErr %v", err, tt.wantErr)
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
	client, err := NewSwarmClient(&Config{
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

	client, err := NewSwarmClient(&Config{
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

func TestGenerateImageWithContext(t *testing.T) {
	// Create a server that takes a while to respond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wait a bit before responding
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"images": [], "info": {}}`))
	}))
	defer server.Close()

	client, err := NewSwarmClient(&Config{
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
