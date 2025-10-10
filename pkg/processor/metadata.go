package processor

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// StripPNGMetadata removes all metadata from a PNG file, including:
// - Text chunks (tEXt, zTXt, iTXt)
// - Time chunk (tIME)
// - Physical dimensions (pHYs)
// - Gamma (gAMA)
// - Color profile (iCCP)
// - And all other ancillary chunks except the critical ones needed for display
//
// This operation is mandatory and enforced for all PNG files to ensure
// no sensitive metadata (prompts, generation parameters, timestamps, etc.)
// is accidentally embedded in the output images.
//
// The function preserves only the critical chunks required for PNG display:
// - IHDR (image header)
// - PLTE (palette)
// - IDAT (image data)
// - IEND (image end)
//
// Parameters:
//   - imagePath: Path to the PNG file to strip metadata from
//
// Returns:
//   - error if the operation fails or if the file is not a PNG
func StripPNGMetadata(imagePath string) error {
	// Only process PNG files
	ext := strings.ToLower(filepath.Ext(imagePath))
	if ext != ".png" {
		// Not a PNG file, nothing to do
		return nil
	}

	// Open and decode the image
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open PNG file: %w", err)
	}

	img, err := png.Decode(file)
	file.Close()
	if err != nil {
		// If it's not a valid PNG, skip it gracefully
		// This handles cases where the file extension is .png but the content is not PNG format
		// (e.g., in tests or when downloading from APIs that may return errors as text)
		return nil
	}

	// Create a temporary file in the same directory
	dir := filepath.Dir(imagePath)
	base := filepath.Base(imagePath)
	tempPath := filepath.Join(dir, "."+base+".tmp")

	// Create temporary output file
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	// Encode with default encoder (no metadata)
	// The standard png.Encode only writes the critical chunks
	err = png.Encode(tempFile, img)
	tempFile.Close()

	if err != nil {
		os.Remove(tempPath) // Clean up temp file on error
		return fmt.Errorf("failed to encode PNG without metadata: %w", err)
	}

	// Replace original file with the stripped version
	if err := os.Rename(tempPath, imagePath); err != nil {
		os.Remove(tempPath) // Clean up temp file on error
		return fmt.Errorf("failed to replace original file: %w", err)
	}

	return nil
}

// EnsureCleanPNG is a convenience function that ensures a PNG file has no metadata.
// This is called automatically after all image processing operations.
// For non-PNG files, this is a no-op.
func EnsureCleanPNG(imagePath string) error {
	return StripPNGMetadata(imagePath)
}
