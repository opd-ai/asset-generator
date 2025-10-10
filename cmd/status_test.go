package cmd

import (
	"testing"

	"github.com/opd-ai/asset-generator/pkg/client"
)

func TestFormatStatusTable(t *testing.T) {
	status := &client.ServerStatus{
		ServerURL:    "http://localhost:7801",
		Status:       "online",
		ResponseTime: "123ms",
		Version:      "1.0.0",
		SessionID:    "test-session-123",
		Backends: []client.BackendStatus{
			{
				ID:          "backend-1",
				Type:        "ComfyUI",
				Status:      "running",
				ModelLoaded: "stable-diffusion-xl",
				GPU:         "NVIDIA RTX 3090",
			},
		},
		ModelsCount:        10,
		ModelsLoaded:       2,
		GenerationsRunning: 0,
		SystemInfo: map[string]interface{}{
			"gpu_memory": "24GB",
			"cpu_count":  16,
		},
	}

	result := formatStatusTable(status)

	// Basic checks to ensure the output contains expected information
	if result == "" {
		t.Error("formatStatusTable returned empty string")
	}

	// Check for key elements in the output
	expectedStrings := []string{
		"SwarmUI Server Status",
		"http://localhost:7801",
		"online",
		"123ms",
		"1.0.0",
		"test-session-123",
		"backend-1",
		"ComfyUI",
		"running",
		"stable-diffusion-xl",
		"NVIDIA RTX 3090",
		"No generations currently running",
	}

	for _, expected := range expectedStrings {
		if !contains(result, expected) {
			t.Errorf("formatStatusTable output missing expected string: %s", expected)
		}
	}
}

func TestFormatStatusTableWithActiveGenerations(t *testing.T) {
	status := &client.ServerStatus{
		ServerURL:    "http://localhost:7801",
		Status:       "online",
		ResponseTime: "123ms",
		Version:      "1.0.0",
		SessionID:    "test-session-123",
		ActiveGenerations: []client.ActiveGeneration{
			{
				SessionID: "gen-session-1",
				Status:    "generating",
				Progress:  0.45,
				Duration:  "2.5m",
			},
		},
		GenerationsRunning: 1,
		ModelsCount:        10,
		ModelsLoaded:       2,
	}

	result := formatStatusTable(status)

	// Check for active generation information
	expectedStrings := []string{
		"Active Generations",
		"gen-session-1",
		"generating",
		"45.0%",
		"2.5m",
	}

	for _, expected := range expectedStrings {
		if !contains(result, expected) {
			t.Errorf("formatStatusTable output missing expected string for active generation: %s", expected)
		}
	}
}

func TestColorizeStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"running status", "running"},
		{"loaded status", "loaded"},
		{"active status", "active"},
		{"online status", "online"},
		{"ready status", "ready"},
		{"generating status", "generating"},
		{"idle status", "idle"},
		{"unloaded status", "unloaded"},
		{"pending status", "pending"},
		{"starting status", "starting"},
		{"error status", "error"},
		{"failed status", "failed"},
		{"offline status", "offline"},
		{"unknown status", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorizeStatus(tt.status)
			if result == "" {
				t.Errorf("colorizeStatus(%s) returned empty string", tt.status)
			}
			// The result should contain the original status text
			if !contains(result, tt.status) {
				t.Errorf("colorizeStatus(%s) does not contain original status", tt.status)
			}
		})
	}
}

func TestValueOrNA(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", "N/A"},
		{"non-empty string", "test", "test"},
		{"whitespace", "   ", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueOrNA(tt.input)
			if result != tt.expected {
				t.Errorf("valueOrNA(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
