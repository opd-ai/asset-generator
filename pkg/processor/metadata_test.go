package processor

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestPNG creates a simple test PNG file
func createTestPNGForMetadata(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with a simple pattern
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func TestStripPNGMetadata(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		filename  string
		shouldErr bool
	}{
		{
			name:      "Valid PNG file",
			filename:  "test.png",
			shouldErr: false,
		},
		{
			name:      "Non-PNG file (should skip)",
			filename:  "test.jpg",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imagePath := filepath.Join(tmpDir, tt.filename)

			if filepath.Ext(tt.filename) == ".png" {
				// Create a test PNG
				if err := createTestPNGForMetadata(imagePath); err != nil {
					t.Fatalf("Failed to create test image: %v", err)
				}
			} else {
				// Create a dummy file for non-PNG test
				if err := os.WriteFile(imagePath, []byte("dummy"), 0644); err != nil {
					t.Fatalf("Failed to create dummy file: %v", err)
				}
			}

			// Get original file size
			origInfo, err := os.Stat(imagePath)
			if err != nil {
				t.Fatalf("Failed to stat original file: %v", err)
			}

			// Strip metadata
			err = StripPNGMetadata(imagePath)

			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify file still exists
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				t.Error("File was deleted instead of being processed")
			}

			// For PNG files, verify the image is still valid
			if filepath.Ext(tt.filename) == ".png" && !tt.shouldErr {
				file, err := os.Open(imagePath)
				if err != nil {
					t.Fatalf("Failed to open processed file: %v", err)
				}
				defer file.Close()

				img, err := png.Decode(file)
				if err != nil {
					t.Fatalf("Failed to decode processed PNG: %v", err)
				}

				// Check dimensions are preserved
				bounds := img.Bounds()
				if bounds.Dx() != 100 || bounds.Dy() != 100 {
					t.Errorf("Image dimensions changed: got %dx%d, want 100x100",
						bounds.Dx(), bounds.Dy())
				}

				// Verify file was actually rewritten (file modification or size might change)
				newInfo, err := os.Stat(imagePath)
				if err != nil {
					t.Fatalf("Failed to stat processed file: %v", err)
				}

				// File should still exist and be a valid size
				if newInfo.Size() == 0 {
					t.Error("Processed file is empty")
				}

				// For this simple test image, size should be reasonable
				if newInfo.Size() > origInfo.Size()*2 {
					t.Logf("Warning: Processed file is significantly larger than original (%d vs %d bytes)",
						newInfo.Size(), origInfo.Size())
				}
			}
		})
	}
}

func TestStripPNGMetadata_NonexistentFile(t *testing.T) {
	err := StripPNGMetadata("/nonexistent/file.png")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestEnsureCleanPNG(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test_clean.png")

	// Create a test PNG
	if err := createTestPNGForMetadata(imagePath); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Ensure it's clean
	err := EnsureCleanPNG(imagePath)
	if err != nil {
		t.Errorf("EnsureCleanPNG failed: %v", err)
	}

	// Verify file is still valid
	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf("Failed to open cleaned file: %v", err)
	}
	defer file.Close()

	_, err = png.Decode(file)
	if err != nil {
		t.Errorf("Failed to decode cleaned PNG: %v", err)
	}
}

func TestStripPNGMetadata_PreservesImageData(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test_preserve.png")

	// Create a test PNG with specific colors
	if err := createTestPNGForMetadata(imagePath); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Read a sample pixel before stripping
	fileBefore, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf("Failed to open file before: %v", err)
	}
	imgBefore, err := png.Decode(fileBefore)
	fileBefore.Close()
	if err != nil {
		t.Fatalf("Failed to decode before: %v", err)
	}
	colorBefore := imgBefore.At(50, 50)

	// Strip metadata
	if err := StripPNGMetadata(imagePath); err != nil {
		t.Fatalf("StripPNGMetadata failed: %v", err)
	}

	// Read the same pixel after stripping
	fileAfter, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf("Failed to open file after: %v", err)
	}
	defer fileAfter.Close()
	imgAfter, err := png.Decode(fileAfter)
	if err != nil {
		t.Fatalf("Failed to decode after: %v", err)
	}
	colorAfter := imgAfter.At(50, 50)

	// Compare colors
	rBefore, gBefore, bBefore, aBefore := colorBefore.RGBA()
	rAfter, gAfter, bAfter, aAfter := colorAfter.RGBA()

	if rBefore != rAfter || gBefore != gAfter || bBefore != bAfter || aBefore != aAfter {
		t.Errorf("Image data changed: before RGBA(%d,%d,%d,%d) != after RGBA(%d,%d,%d,%d)",
			rBefore, gBefore, bBefore, aBefore,
			rAfter, gAfter, bAfter, aAfter)
	}
}
