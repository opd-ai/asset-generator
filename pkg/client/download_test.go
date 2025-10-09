package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadImages(t *testing.T) {
	// Create a test HTTP server that serves fake image data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is for an image
		if r.URL.Path == "/View/local/raw/2024-05-19/test-image.png" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("fake-image-data"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		}
	}))
	defer server.Close()

	// Create temporary directory for test output
	tmpDir := t.TempDir()

	// Create test client
	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
		Verbose: false,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name        string
		imagePaths  []string
		outputDir   string
		expectError bool
		expectFiles int
	}{
		{
			name:        "Single image download",
			imagePaths:  []string{"View/local/raw/2024-05-19/test-image.png"},
			outputDir:   tmpDir,
			expectError: false,
			expectFiles: 1,
		},
		{
			name:        "Empty image paths",
			imagePaths:  []string{},
			outputDir:   tmpDir,
			expectError: true,
			expectFiles: 0,
		},
		{
			name:        "Invalid image path",
			imagePaths:  []string{"View/local/raw/2024-05-19/nonexistent.png"},
			outputDir:   tmpDir,
			expectError: true,
			expectFiles: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			savedPaths, err := client.DownloadImages(ctx, tt.imagePaths, tt.outputDir)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(savedPaths) != tt.expectFiles {
				t.Errorf("Expected %d files, got %d", tt.expectFiles, len(savedPaths))
			}

			// Verify files exist if expected
			if !tt.expectError {
				for _, path := range savedPaths {
					if _, err := os.Stat(path); os.IsNotExist(err) {
						t.Errorf("Expected file does not exist: %s", path)
					}
				}
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test-content"))
	}))
	defer server.Close()

	// Create temporary directory
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-file.txt")

	// Create test client
	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test download
	ctx := context.Background()
	err = client.downloadFile(ctx, server.URL+"/test", outputPath)
	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	// Verify file exists and has correct content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != "test-content" {
		t.Errorf("Expected content 'test-content', got '%s'", string(content))
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		expectError bool
		setup       func(string) error
	}{
		{
			name:        "Create new directory",
			path:        filepath.Join(tmpDir, "new-dir"),
			expectError: false,
		},
		{
			name:        "Existing directory",
			path:        tmpDir,
			expectError: false,
		},
		{
			name: "Path is a file",
			path: filepath.Join(tmpDir, "file-not-dir"),
			setup: func(path string) error {
				return os.WriteFile(path, []byte("test"), 0644)
			},
			expectError: true,
		},
		{
			name:        "Nested directory creation",
			path:        filepath.Join(tmpDir, "nested", "dir", "structure"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(tt.path); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			err := ensureDir(tt.path)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify directory exists if no error expected
			if !tt.expectError {
				info, err := os.Stat(tt.path)
				if err != nil {
					t.Errorf("Directory does not exist: %v", err)
				} else if !info.IsDir() {
					t.Errorf("Path exists but is not a directory")
				}
			}
		})
	}
}

func TestDownloadImagesWithAuth(t *testing.T) {
	// Create a test HTTP server that requires authentication
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated-content"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()

	// Create test client with API key
	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test download with authentication
	ctx := context.Background()
	imagePaths := []string{"test-image.png"}
	savedPaths, err := client.DownloadImages(ctx, imagePaths, tmpDir)

	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	if len(savedPaths) != 1 {
		t.Errorf("Expected 1 saved path, got %d", len(savedPaths))
	}

	// Verify file content
	if len(savedPaths) > 0 {
		content, err := os.ReadFile(savedPaths[0])
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		if string(content) != "authenticated-content" {
			t.Errorf("Expected 'authenticated-content', got '%s'", string(content))
		}
	}
}

func TestDownloadImagesPartialFailure(t *testing.T) {
	// Create a test HTTP server that returns 404 for some images
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount%2 == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("image-%d", callCount)))
		}
	}))
	defer server.Close()

	tmpDir := t.TempDir()

	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	imagePaths := []string{"image1.png", "image2.png", "image3.png", "image4.png"}
	savedPaths, err := client.DownloadImages(ctx, imagePaths, tmpDir)

	// Should get partial success error
	if err == nil {
		t.Error("Expected error for partial failure")
	}

	// Should have some successful downloads
	if len(savedPaths) == 0 {
		t.Error("Expected some successful downloads")
	}

	// Verify successful downloads
	for _, path := range savedPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file does not exist: %s", path)
		}
	}
}

func TestDownloadImagesWithTemplate(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test-image-data"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()

	client, err := NewAssetClient(&Config{
		BaseURL: server.URL,
		Verbose: false,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name             string
		imagePaths       []string
		filenameTemplate string
		metadata         map[string]interface{}
		expectedPattern  string // Regex pattern or exact match
	}{
		{
			name:             "Index placeholder",
			imagePaths:       []string{"View/local/raw/2024-05-19/original.png"},
			filenameTemplate: "image-{index}.png",
			expectedPattern:  "image-000.png",
		},
		{
			name:             "One-based index",
			imagePaths:       []string{"View/local/raw/2024-05-19/original.png"},
			filenameTemplate: "image-{i1}.png",
			expectedPattern:  "image-1.png",
		},
		{
			name:             "Seed placeholder",
			imagePaths:       []string{"View/local/raw/2024-05-19/original.png"},
			filenameTemplate: "image-{seed}.png",
			metadata:         map[string]interface{}{"seed": 12345},
			expectedPattern:  "image-12345.png",
		},
		{
			name:             "Multiple placeholders",
			imagePaths:       []string{"View/local/raw/2024-05-19/original.png"},
			filenameTemplate: "{model}-{width}x{height}-{index}.png",
			metadata: map[string]interface{}{
				"model":  "sdxl",
				"width":  1024,
				"height": 768,
			},
			expectedPattern: "sdxl-1024x768-000.png",
		},
		{
			name:             "Extension from original",
			imagePaths:       []string{"View/local/raw/2024-05-19/image.jpg"},
			filenameTemplate: "custom-{index}",
			expectedPattern:  "custom-000.jpg",
		},
		{
			name:             "Prompt placeholder",
			imagePaths:       []string{"View/local/raw/2024-05-19/original.png"},
			filenameTemplate: "{prompt}-{index}.png",
			metadata:         map[string]interface{}{"prompt": "a beautiful sunset over mountains"},
			expectedPattern:  "a_beautiful_sunset_over_mountains-000.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			opts := &DownloadOptions{
				OutputDir:        tmpDir,
				FilenameTemplate: tt.filenameTemplate,
				Metadata:         tt.metadata,
			}

			savedPaths, err := client.DownloadImagesWithOptions(ctx, tt.imagePaths, opts)
			if err != nil {
				t.Fatalf("Download failed: %v", err)
			}

			if len(savedPaths) != 1 {
				t.Fatalf("Expected 1 file, got %d", len(savedPaths))
			}

			// Extract filename from path
			filename := filepath.Base(savedPaths[0])

			if filename != tt.expectedPattern {
				t.Errorf("Expected filename '%s', got '%s'", tt.expectedPattern, filename)
			}

			// Verify file exists
			if _, err := os.Stat(savedPaths[0]); os.IsNotExist(err) {
				t.Errorf("Expected file does not exist: %s", savedPaths[0])
			}
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		name             string
		template         string
		index            int
		originalFilename string
		metadata         map[string]interface{}
		expectedContains []string // Strings that should be in the result
		expectedExact    string   // Exact expected result (if applicable)
	}{
		{
			name:             "Zero-padded index",
			template:         "img-{index}.png",
			index:            5,
			originalFilename: "test.png",
			expectedExact:    "img-005.png",
		},
		{
			name:             "One-based index",
			template:         "img-{i1}.png",
			index:            0,
			originalFilename: "test.png",
			expectedExact:    "img-1.png",
		},
		{
			name:             "Original filename",
			template:         "copy-of-{original}",
			index:            0,
			originalFilename: "myimage.jpg",
			expectedExact:    "copy-of-myimage.jpg",
		},
		{
			name:             "Extension extraction",
			template:         "image-{index}{ext}",
			index:            1,
			originalFilename: "photo.jpeg",
			expectedExact:    "image-001.jpeg",
		},
		{
			name:             "Date and time",
			template:         "{date}-{time}-{index}.png",
			index:            0,
			originalFilename: "test.png",
			expectedContains: []string{"-000.png"},
		},
		{
			name:             "Metadata seed",
			template:         "seed-{seed}-img.png",
			index:            0,
			originalFilename: "test.png",
			metadata:         map[string]interface{}{"seed": 42},
			expectedExact:    "seed-42-img.png",
		},
		{
			name:             "Multiple metadata",
			template:         "{model}-{width}x{height}.png",
			index:            0,
			originalFilename: "test.png",
			metadata: map[string]interface{}{
				"model":  "flux",
				"width":  512,
				"height": 512,
			},
			expectedExact: "flux-512x512.png",
		},
		{
			name:             "Prompt sanitization",
			template:         "{prompt}.png",
			index:            0,
			originalFilename: "test.png",
			metadata:         map[string]interface{}{"prompt": "a cat/dog in <the> rain?"},
			expectedExact:    "a_catdog_in_the_rain.png",
		},
		{
			name:             "Auto extension append",
			template:         "image-{index}",
			index:            0,
			originalFilename: "photo.png",
			expectedExact:    "image-000.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateFilename(tt.template, tt.index, tt.originalFilename, tt.metadata)

			if tt.expectedExact != "" {
				if result != tt.expectedExact {
					t.Errorf("Expected '%s', got '%s'", tt.expectedExact, result)
				}
			}

			for _, substr := range tt.expectedContains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected result to contain '%s', got '%s'", substr, result)
				}
			}
		})
	}
}

func TestSanitizeForFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello world",
			expected: "hello_world",
		},
		{
			input:    "file/with\\slashes",
			expected: "filewithslashes",
		},
		{
			input:    "question?mark*asterisk",
			expected: "questionmarkasterisk",
		},
		{
			input:    "multiple   spaces",
			expected: "multiple_spaces",
		},
		{
			input:    "__leading_trailing__",
			expected: "leading_trailing",
		},
		{
			input:    "tabs\tand\nnewlines\r",
			expected: "tabsandnewlines",
		},
		{
			input:    "valid-filename_123.txt",
			expected: "valid-filename_123.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeForFilename(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeForFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
