package processor

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestImage creates a simple test image file
func createTestImage(t *testing.T, path string, width, height int) {
	t.Helper()

	// Create a simple RGBA image with a gradient pattern
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a simple gradient
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Save as PNG
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
}

func TestDownscaleImage(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		srcWidth      int
		srcHeight     int
		targetWidth   int
		targetHeight  int
		filter        ResizeFilter
		expectError   bool
		errorContains string
	}{
		{
			name:         "Downscale both dimensions with Lanczos",
			srcWidth:     1024,
			srcHeight:    768,
			targetWidth:  512,
			targetHeight: 384,
			filter:       FilterLanczos,
			expectError:  false,
		},
		{
			name:         "Downscale with bilinear filter",
			srcWidth:     800,
			srcHeight:    600,
			targetWidth:  400,
			targetHeight: 300,
			filter:       FilterBiLinear,
			expectError:  false,
		},
		{
			name:         "Downscale with nearest neighbor",
			srcWidth:     640,
			srcHeight:    480,
			targetWidth:  320,
			targetHeight: 240,
			filter:       FilterNearestNeighbor,
			expectError:  false,
		},
		{
			name:         "Auto-calculate width from height",
			srcWidth:     1920,
			srcHeight:    1080,
			targetWidth:  0, // Auto-calculate
			targetHeight: 540,
			filter:       FilterLanczos,
			expectError:  false,
		},
		{
			name:         "Auto-calculate height from width",
			srcWidth:     1920,
			srcHeight:    1080,
			targetWidth:  960,
			targetHeight: 0, // Auto-calculate
			filter:       FilterLanczos,
			expectError:  false,
		},
		{
			name:          "Error: Both dimensions zero",
			srcWidth:      800,
			srcHeight:     600,
			targetWidth:   0,
			targetHeight:  0,
			filter:        FilterLanczos,
			expectError:   true,
			errorContains: "at least one dimension",
		},
		{
			name:          "Error: Upscaling not allowed",
			srcWidth:      512,
			srcHeight:     512,
			targetWidth:   1024,
			targetHeight:  1024,
			filter:        FilterLanczos,
			expectError:   true,
			errorContains: "not smaller than source",
		},
		{
			name:          "Error: Negative dimensions",
			srcWidth:      800,
			srcHeight:     600,
			targetWidth:   -100,
			targetHeight:  300,
			filter:        FilterLanczos,
			expectError:   true,
			errorContains: "cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create source image
			srcPath := filepath.Join(tmpDir, tt.name+"_source.png")
			createTestImage(t, srcPath, tt.srcWidth, tt.srcHeight)

			// Create output path
			dstPath := filepath.Join(tmpDir, tt.name+"_output.png")

			// Run downscale
			opts := DownscaleOptions{
				Width:       tt.targetWidth,
				Height:      tt.targetHeight,
				Filter:      tt.filter,
				JPEGQuality: 90,
			}

			err := DownscaleImage(srcPath, dstPath, opts)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got no error", tt.errorContains)
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Fatalf("Expected error containing '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify output file exists
			if _, err := os.Stat(dstPath); os.IsNotExist(err) {
				t.Fatalf("Output file was not created")
			}

			// Verify output dimensions
			width, height, err := GetImageDimensions(dstPath)
			if err != nil {
				t.Fatalf("Failed to read output image dimensions: %v", err)
			}

			// Calculate expected dimensions
			expectedWidth := tt.targetWidth
			expectedHeight := tt.targetHeight
			if expectedWidth == 0 {
				expectedWidth = int(float64(tt.targetHeight) * float64(tt.srcWidth) / float64(tt.srcHeight))
			}
			if expectedHeight == 0 {
				expectedHeight = int(float64(tt.targetWidth) * float64(tt.srcHeight) / float64(tt.srcWidth))
			}

			if width != expectedWidth || height != expectedHeight {
				t.Fatalf("Expected dimensions %dx%d, got %dx%d", expectedWidth, expectedHeight, width, height)
			}
		})
	}
}

func TestDownscaleInPlace(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test image
	srcPath := filepath.Join(tmpDir, "test_inplace.png")
	createTestImage(t, srcPath, 1024, 768)

	// Get original file info
	originalInfo, err := os.Stat(srcPath)
	if err != nil {
		t.Fatalf("Failed to stat original file: %v", err)
	}

	// Downscale in place
	opts := DownscaleOptions{
		Width:       512,
		Height:      384,
		Filter:      FilterLanczos,
		JPEGQuality: 90,
	}

	err = DownscaleInPlace(srcPath, opts)
	if err != nil {
		t.Fatalf("DownscaleInPlace failed: %v", err)
	}

	// Verify file still exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		t.Fatalf("File was deleted instead of replaced")
	}

	// Verify dimensions changed
	width, height, err := GetImageDimensions(srcPath)
	if err != nil {
		t.Fatalf("Failed to read dimensions: %v", err)
	}

	if width != 512 || height != 384 {
		t.Fatalf("Expected dimensions 512x384, got %dx%d", width, height)
	}

	// Verify file was modified (modification time or size changed)
	newInfo, err := os.Stat(srcPath)
	if err != nil {
		t.Fatalf("Failed to stat modified file: %v", err)
	}

	// Size should be different after downscaling
	if newInfo.Size() == originalInfo.Size() {
		t.Fatalf("File size did not change after downscaling")
	}
}

func TestGetImageDimensions(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		width         int
		height        int
		expectError   bool
		errorContains string
	}{
		{
			name:        "Standard dimensions",
			width:       1920,
			height:      1080,
			expectError: false,
		},
		{
			name:        "Square image",
			width:       512,
			height:      512,
			expectError: false,
		},
		{
			name:        "Portrait orientation",
			width:       768,
			height:      1024,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgPath := filepath.Join(tmpDir, tt.name+".png")
			createTestImage(t, imgPath, tt.width, tt.height)

			width, height, err := GetImageDimensions(imgPath)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error, got none")
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Fatalf("Expected error containing '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if width != tt.width || height != tt.height {
				t.Fatalf("Expected dimensions %dx%d, got %dx%d", tt.width, tt.height, width, height)
			}
		})
	}
}

func TestGetImageDimensions_NonexistentFile(t *testing.T) {
	_, _, err := GetImageDimensions("/nonexistent/file.png")
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got none")
	}
}

func TestDownscaleImage_NonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	dstPath := filepath.Join(tmpDir, "output.png")

	opts := DownscaleOptions{
		Width:  512,
		Height: 512,
		Filter: FilterLanczos,
	}

	err := DownscaleImage("/nonexistent/file.png", dstPath, opts)
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got none")
	}
}

func TestDownscaleImage_AspectRatioPreservation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a 16:9 aspect ratio image
	srcPath := filepath.Join(tmpDir, "aspect_test.png")
	createTestImage(t, srcPath, 1920, 1080)

	dstPath := filepath.Join(tmpDir, "aspect_output.png")

	// Downscale by width only
	opts := DownscaleOptions{
		Width:       960,
		Height:      0, // Auto-calculate
		Filter:      FilterLanczos,
		JPEGQuality: 90,
	}

	err := DownscaleImage(srcPath, dstPath, opts)
	if err != nil {
		t.Fatalf("Downscale failed: %v", err)
	}

	// Verify aspect ratio preserved
	width, height, err := GetImageDimensions(dstPath)
	if err != nil {
		t.Fatalf("Failed to read dimensions: %v", err)
	}

	// Check aspect ratio (960:540 = 16:9 = 1920:1080)
	if width != 960 || height != 540 {
		t.Fatalf("Aspect ratio not preserved: expected 960x540, got %dx%d", width, height)
	}
}

// TestDownscaleByPercentage tests percentage-based downscaling
func TestDownscaleByPercentage(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		srcWidth      int
		srcHeight     int
		percentage    float64
		expectError   bool
		errorContains string
	}{
		{
			name:        "Downscale to 50%",
			srcWidth:    1000,
			srcHeight:   800,
			percentage:  50,
			expectError: false,
		},
		{
			name:        "Downscale to 75%",
			srcWidth:    1920,
			srcHeight:   1080,
			percentage:  75,
			expectError: false,
		},
		{
			name:        "Downscale to 25%",
			srcWidth:    800,
			srcHeight:   600,
			percentage:  25,
			expectError: false,
		},
		{
			name:          "Error: Percentage over 100",
			srcWidth:      640,
			srcHeight:     480,
			percentage:    150,
			expectError:   true,
			errorContains: "percentage must be between 0 and 100",
		},
		{
			name:          "Error: Negative percentage",
			srcWidth:      640,
			srcHeight:     480,
			percentage:    -10,
			expectError:   true,
			errorContains: "percentage must be between 0 and 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create source image
			srcPath := filepath.Join(tmpDir, tt.name+"_source.png")
			createTestImage(t, srcPath, tt.srcWidth, tt.srcHeight)

			// Create output path
			dstPath := filepath.Join(tmpDir, tt.name+"_output.png")

			// Run downscale with percentage
			opts := DownscaleOptions{
				Percentage:  tt.percentage,
				Filter:      FilterLanczos,
				JPEGQuality: 90,
			}

			err := DownscaleImage(srcPath, dstPath, opts)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got no error", tt.errorContains)
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Fatalf("Expected error containing '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify output file exists
			if _, err := os.Stat(dstPath); os.IsNotExist(err) {
				t.Fatalf("Output file was not created")
			}

			// Verify output dimensions match percentage
			width, height, err := GetImageDimensions(dstPath)
			if err != nil {
				t.Fatalf("Failed to read output image dimensions: %v", err)
			}

			// Calculate expected dimensions
			scale := tt.percentage / 100.0
			expectedWidth := int(float64(tt.srcWidth) * scale)
			expectedHeight := int(float64(tt.srcHeight) * scale)

			if width != expectedWidth || height != expectedHeight {
				t.Fatalf("Expected dimensions %dx%d (%.0f%% of %dx%d), got %dx%d",
					expectedWidth, expectedHeight, tt.percentage, tt.srcWidth, tt.srcHeight, width, height)
			}

			// Verify aspect ratio preserved
			srcAspect := float64(tt.srcWidth) / float64(tt.srcHeight)
			dstAspect := float64(width) / float64(height)
			aspectDiff := srcAspect - dstAspect
			if aspectDiff < 0 {
				aspectDiff = -aspectDiff
			}
			if aspectDiff > 0.01 { // Allow small rounding error
				t.Fatalf("Aspect ratio not preserved: source %.3f, output %.3f", srcAspect, dstAspect)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
