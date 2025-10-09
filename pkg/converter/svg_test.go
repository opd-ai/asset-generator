package converter

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createTestImage creates a simple test image
func createTestImage(t *testing.T, path string) {
	// Create a simple 100x100 image with some patterns
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Draw a simple pattern
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if (x+y)%20 < 10 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red
			} else {
				img.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue
			}
		}
	}

	// Save the image
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
}

func TestNewSVGConverter(t *testing.T) {
	converter := NewSVGConverter()

	if converter == nil {
		t.Fatal("Expected converter to be created")
	}

	if converter.PrimitiveShapes != 100 {
		t.Errorf("Expected PrimitiveShapes to be 100, got %d", converter.PrimitiveShapes)
	}

	if converter.PrimitiveMode != 1 {
		t.Errorf("Expected PrimitiveMode to be 1, got %d", converter.PrimitiveMode)
	}
}

func TestConvertWithPrimitive(t *testing.T) {
	// Create temp directory for test files
	tempDir := t.TempDir()

	// Create test image
	inputPath := filepath.Join(tempDir, "test_input.png")
	createTestImage(t, inputPath)

	// Convert to SVG
	outputPath := filepath.Join(tempDir, "test_output.svg")
	converter := NewSVGConverter()

	opts := ConversionOptions{
		Method:          MethodPrimitive,
		OutputPath:      outputPath,
		PrimitiveShapes: 50, // Use fewer shapes for faster test
	}

	result, err := converter.ConvertToSVG(inputPath, opts)
	if err != nil {
		t.Fatalf("ConvertToSVG failed: %v", err)
	}

	if result != outputPath {
		t.Errorf("Expected output path %s, got %s", outputPath, result)
	}

	// Check if output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output SVG file was not created")
	}

	// Check if output file has content
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Output SVG file is empty")
	}
}

func TestConvertWithPrimitiveDefault(t *testing.T) {
	// Create temp directory for test files
	tempDir := t.TempDir()

	// Create test image
	inputPath := filepath.Join(tempDir, "test_input.png")
	createTestImage(t, inputPath)

	// Convert to SVG using default function
	outputPath := filepath.Join(tempDir, "test_output.svg")

	result, err := ConvertWithPrimitiveDefault(inputPath, outputPath, 30)
	if err != nil {
		t.Fatalf("ConvertWithPrimitiveDefault failed: %v", err)
	}

	if result != outputPath {
		t.Errorf("Expected output path %s, got %s", outputPath, result)
	}

	// Check if output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output SVG file was not created")
	}
}

func TestConvertWithGotrace(t *testing.T) {
	t.Skip("Gotrace has CGO dependencies that require specific Go runtime patches - skipping for now")
	
	// Create temp directory for test files
	tempDir := t.TempDir()

	// Create test image
	inputPath := filepath.Join(tempDir, "test_input.png")
	createTestImage(t, inputPath)

	// Convert to SVG using gotrace
	outputPath := filepath.Join(tempDir, "test_output.svg")
	converter := NewSVGConverter()

	opts := ConversionOptions{
		Method:     MethodGotrace,
		OutputPath: outputPath,
	}

	result, err := converter.ConvertToSVG(inputPath, opts)
	if err != nil {
		t.Fatalf("ConvertToSVG with gotrace failed: %v", err)
	}

	if result != outputPath {
		t.Errorf("Expected output path %s, got %s", outputPath, result)
	}

	// Check if output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output SVG file was not created")
	}

	// Check if output file has content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Output SVG file is empty")
	}

	// Check if it looks like SVG
	contentStr := string(content)
	if !strings.Contains(contentStr, "<svg") && !strings.Contains(contentStr, "<?xml") {
		t.Error("Output file does not appear to be valid SVG")
	}
}

func TestConvertNonExistentFile(t *testing.T) {
	converter := NewSVGConverter()

	opts := ConversionOptions{
		Method:     MethodPrimitive,
		OutputPath: "/tmp/output.svg",
	}

	_, err := converter.ConvertToSVG("/nonexistent/file.png", opts)
	if err == nil {
		t.Error("Expected error for non-existent input file")
	}
}

func TestConvertWithDefaultOutputPath(t *testing.T) {
	// Create temp directory for test files
	tempDir := t.TempDir()

	// Create test image
	inputPath := filepath.Join(tempDir, "test_input.png")
	createTestImage(t, inputPath)

	// Convert without specifying output path
	converter := NewSVGConverter()

	opts := ConversionOptions{
		Method:          MethodPrimitive,
		PrimitiveShapes: 20,
	}

	result, err := converter.ConvertToSVG(inputPath, opts)
	if err != nil {
		t.Fatalf("ConvertToSVG failed: %v", err)
	}

	expectedPath := filepath.Join(tempDir, "test_input.svg")
	if result != expectedPath {
		t.Errorf("Expected output path %s, got %s", expectedPath, result)
	}

	// Check if output file exists
	if _, err := os.Stat(result); os.IsNotExist(err) {
		t.Error("Output SVG file was not created")
	}

	// Clean up
	os.Remove(result)
}
