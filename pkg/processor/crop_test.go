package processor

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createCropTestImage creates a test image with whitespace borders
func createCropTestImage(width, height, contentX, contentY, contentW, contentH int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with white background
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, white)
		}
	}

	// Add content area (dark pixels)
	black := color.RGBA{0, 0, 0, 255}
	for y := contentY; y < contentY+contentH && y < height; y++ {
		for x := contentX; x < contentX+contentW && x < width; x++ {
			img.Set(x, y, black)
		}
	}

	return img
}

// saveTestImage saves an image to a file
func saveTestImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func TestAutoCropImage(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		imgWidth       int
		imgHeight      int
		contentX       int
		contentY       int
		contentW       int
		contentH       int
		opts           CropOptions
		expectWidth    int
		expectHeight   int
		expectError    bool
		expectNoChange bool
	}{
		{
			name:         "Crop whitespace borders",
			imgWidth:     100,
			imgHeight:    100,
			contentX:     20,
			contentY:     30,
			contentW:     40,
			contentH:     30,
			opts:         CropOptions{Threshold: 250, Tolerance: 10},
			expectWidth:  40,
			expectHeight: 30,
			expectError:  false,
		},
		{
			name:           "No whitespace - no change",
			imgWidth:       50,
			imgHeight:      50,
			contentX:       0,
			contentY:       0,
			contentW:       50,
			contentH:       50,
			opts:           CropOptions{Threshold: 250, Tolerance: 10},
			expectWidth:    50,
			expectHeight:   50,
			expectError:    false,
			expectNoChange: true,
		},
		{
			name:         "Preserve aspect ratio",
			imgWidth:     200,
			imgHeight:    100,
			contentX:     50,
			contentY:     25,
			contentW:     50,
			contentH:     50,
			opts:         CropOptions{Threshold: 250, Tolerance: 10, PreserveAspectRatio: true},
			expectWidth:  100, // Will be expanded to maintain 2:1 aspect ratio
			expectHeight: 50,
			expectError:  false,
		},
		{
			name:         "Top-left corner content",
			imgWidth:     100,
			imgHeight:    100,
			contentX:     0,
			contentY:     0,
			contentW:     30,
			contentH:     40,
			opts:         CropOptions{Threshold: 250, Tolerance: 10},
			expectWidth:  30,
			expectHeight: 40,
			expectError:  false,
		},
		{
			name:         "Bottom-right corner content",
			imgWidth:     100,
			imgHeight:    100,
			contentX:     70,
			contentY:     60,
			contentW:     30,
			contentH:     40,
			opts:         CropOptions{Threshold: 250, Tolerance: 10},
			expectWidth:  30,
			expectHeight: 40,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test image
			img := createCropTestImage(tt.imgWidth, tt.imgHeight, tt.contentX, tt.contentY, tt.contentW, tt.contentH)

			// Save to temp file
			inputPath := filepath.Join(tmpDir, tt.name+"_input.png")
			if err := saveTestImage(img, inputPath); err != nil {
				t.Fatalf("Failed to save test image: %v", err)
			}

			// Perform crop
			outputPath := filepath.Join(tmpDir, tt.name+"_output.png")
			err := AutoCropImage(inputPath, outputPath, tt.opts)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.expectError {
				return
			}

			// Load and verify output
			outputFile, err := os.Open(outputPath)
			if err != nil {
				t.Fatalf("Failed to open output file: %v", err)
			}
			defer outputFile.Close()

			outputImg, _, err := image.Decode(outputFile)
			if err != nil {
				t.Fatalf("Failed to decode output image: %v", err)
			}

			bounds := outputImg.Bounds()
			width := bounds.Dx()
			height := bounds.Dy()

			// Verify dimensions
			if width != tt.expectWidth {
				t.Errorf("Expected width %d, got %d", tt.expectWidth, width)
			}
			if height != tt.expectHeight {
				t.Errorf("Expected height %d, got %d", tt.expectHeight, height)
			}
		})
	}
}

func TestAutoCropInPlace(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test image
	img := createCropTestImage(100, 100, 20, 20, 40, 40)
	imagePath := filepath.Join(tmpDir, "test_inplace.png")
	if err := saveTestImage(img, imagePath); err != nil {
		t.Fatalf("Failed to save test image: %v", err)
	}

	// Get original file info
	origInfo, err := os.Stat(imagePath)
	if err != nil {
		t.Fatalf("Failed to stat original file: %v", err)
	}

	// Perform in-place crop
	opts := CropOptions{Threshold: 250, Tolerance: 10}
	err = AutoCropInPlace(imagePath, opts)
	if err != nil {
		t.Fatalf("AutoCropInPlace failed: %v", err)
	}

	// Verify file was modified
	newInfo, err := os.Stat(imagePath)
	if err != nil {
		t.Fatalf("Failed to stat cropped file: %v", err)
	}

	if newInfo.Size() >= origInfo.Size() {
		t.Errorf("Expected cropped file to be smaller, orig: %d, new: %d", origInfo.Size(), newInfo.Size())
	}

	// Verify dimensions
	f, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf("Failed to open cropped file: %v", err)
	}
	defer f.Close()

	croppedImg, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("Failed to decode cropped image: %v", err)
	}

	bounds := croppedImg.Bounds()
	if bounds.Dx() != 40 || bounds.Dy() != 40 {
		t.Errorf("Expected 40x40 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestIsWhitespace(t *testing.T) {
	tests := []struct {
		name      string
		color     color.Color
		threshold uint8
		tolerance uint8
		expect    bool
	}{
		{
			name:      "Pure white",
			color:     color.RGBA{255, 255, 255, 255},
			threshold: 250,
			tolerance: 10,
			expect:    true,
		},
		{
			name:      "Near white within tolerance",
			color:     color.RGBA{245, 248, 250, 255},
			threshold: 250,
			tolerance: 10,
			expect:    true,
		},
		{
			name:      "Light gray outside tolerance",
			color:     color.RGBA{200, 200, 200, 255},
			threshold: 250,
			tolerance: 10,
			expect:    false,
		},
		{
			name:      "Pure black",
			color:     color.RGBA{0, 0, 0, 255},
			threshold: 250,
			tolerance: 10,
			expect:    false,
		},
		{
			name:      "Transparent pixel",
			color:     color.RGBA{255, 255, 255, 0},
			threshold: 250,
			tolerance: 10,
			expect:    true,
		},
		{
			name:      "One dark channel",
			color:     color.RGBA{100, 255, 255, 255},
			threshold: 250,
			tolerance: 10,
			expect:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWhitespace(tt.color, tt.threshold, tt.tolerance)
			if result != tt.expect {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestDetectContentBounds(t *testing.T) {
	tests := []struct {
		name       string
		imgWidth   int
		imgHeight  int
		contentX   int
		contentY   int
		contentW   int
		contentH   int
		expectMinX int
		expectMinY int
		expectMaxX int
		expectMaxY int
	}{
		{
			name:       "Centered content",
			imgWidth:   100,
			imgHeight:  100,
			contentX:   30,
			contentY:   40,
			contentW:   20,
			contentH:   10,
			expectMinX: 30,
			expectMinY: 40,
			expectMaxX: 50,
			expectMaxY: 50,
		},
		{
			name:       "Full image",
			imgWidth:   50,
			imgHeight:  50,
			contentX:   0,
			contentY:   0,
			contentW:   50,
			contentH:   50,
			expectMinX: 0,
			expectMinY: 0,
			expectMaxX: 50,
			expectMaxY: 50,
		},
		{
			name:       "Top-left corner",
			imgWidth:   100,
			imgHeight:  100,
			contentX:   0,
			contentY:   0,
			contentW:   25,
			contentH:   25,
			expectMinX: 0,
			expectMinY: 0,
			expectMaxX: 25,
			expectMaxY: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createCropTestImage(tt.imgWidth, tt.imgHeight, tt.contentX, tt.contentY, tt.contentW, tt.contentH)
			bounds := detectContentBounds(img, 250, 10)

			if bounds.Min.X != tt.expectMinX {
				t.Errorf("Expected MinX=%d, got %d", tt.expectMinX, bounds.Min.X)
			}
			if bounds.Min.Y != tt.expectMinY {
				t.Errorf("Expected MinY=%d, got %d", tt.expectMinY, bounds.Min.Y)
			}
			if bounds.Max.X != tt.expectMaxX {
				t.Errorf("Expected MaxX=%d, got %d", tt.expectMaxX, bounds.Max.X)
			}
			if bounds.Max.Y != tt.expectMaxY {
				t.Errorf("Expected MaxY=%d, got %d", tt.expectMaxY, bounds.Max.Y)
			}
		})
	}
}

func TestPreserveAspectRatio(t *testing.T) {
	tests := []struct {
		name         string
		srcBounds    image.Rectangle
		cropBounds   image.Rectangle
		expectWidth  int
		expectHeight int
	}{
		{
			name:         "Wider crop needs height expansion",
			srcBounds:    image.Rect(0, 0, 200, 100),  // 2:1 ratio
			cropBounds:   image.Rect(50, 40, 150, 60), // 100x20, 5:1 ratio
			expectWidth:  100,
			expectHeight: 50, // Expanded to match 2:1 ratio
		},
		{
			name:         "Taller crop needs width expansion",
			srcBounds:    image.Rect(0, 0, 100, 100), // 1:1 ratio
			cropBounds:   image.Rect(40, 20, 60, 80), // 20x60, 1:3 ratio
			expectWidth:  60,                         // Expanded to match 1:1 ratio
			expectHeight: 60,
		},
		{
			name:         "Already matching aspect ratio",
			srcBounds:    image.Rect(0, 0, 200, 100),  // 2:1 ratio
			cropBounds:   image.Rect(50, 25, 150, 75), // 100x50, 2:1 ratio
			expectWidth:  100,
			expectHeight: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preserveAspectRatio(tt.srcBounds, tt.cropBounds)
			width := result.Dx()
			height := result.Dy()

			if width != tt.expectWidth {
				t.Errorf("Expected width %d, got %d", tt.expectWidth, width)
			}
			if height != tt.expectHeight {
				t.Errorf("Expected height %d, got %d", tt.expectHeight, height)
			}

			// Verify aspect ratio matches source
			srcAspect := float64(tt.srcBounds.Dx()) / float64(tt.srcBounds.Dy())
			resultAspect := float64(width) / float64(height)
			aspectDiff := resultAspect / srcAspect

			// Allow 1% tolerance
			if aspectDiff < 0.99 || aspectDiff > 1.01 {
				t.Errorf("Aspect ratio mismatch: source=%.2f, result=%.2f", srcAspect, resultAspect)
			}
		})
	}
}

func TestCropImage(t *testing.T) {
	// Create test image
	img := createCropTestImage(100, 100, 20, 30, 40, 30)

	// Crop to content area
	cropBounds := image.Rect(20, 30, 60, 60)
	cropped := cropImage(img, cropBounds)

	bounds := cropped.Bounds()
	if bounds.Dx() != 40 {
		t.Errorf("Expected width 40, got %d", bounds.Dx())
	}
	if bounds.Dy() != 30 {
		t.Errorf("Expected height 30, got %d", bounds.Dy())
	}

	// Verify that cropped image contains the expected content
	// Check a pixel that should be black (from content area)
	c := cropped.At(bounds.Min.X, bounds.Min.Y)
	r, g, b, _ := c.RGBA()
	if r > 100 || g > 100 || b > 100 {
		t.Errorf("Expected dark pixel in content area, got RGB(%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}
