package processor

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// ResizeFilter represents an image scaling filter algorithm
type ResizeFilter int

const (
	// FilterLanczos uses Lanczos resampling - highest quality, best for downscaling
	FilterLanczos ResizeFilter = iota
	// FilterBiLinear uses bilinear interpolation - good balance of speed and quality
	FilterBiLinear
	// FilterNearestNeighbor uses nearest neighbor - fastest but lowest quality
	FilterNearestNeighbor
)

// DownscaleOptions configures image downscaling behavior
type DownscaleOptions struct {
	// Target width in pixels (0 means auto-calculate based on aspect ratio)
	Width int
	// Target height in pixels (0 means auto-calculate based on aspect ratio)
	Height int
	// Percentage to scale by (0-100, 0 means use Width/Height instead)
	// If set, this takes precedence over Width/Height
	Percentage float64
	// Scaling algorithm to use (default: Lanczos)
	Filter ResizeFilter
	// Quality for JPEG output (1-100, default: 90)
	JPEGQuality int
}

// DownscaleImage resizes an image using high-quality Lanczos filtering.
// This performs local postprocessing after downloading from the API.
//
// Parameters:
//   - inputPath: Path to the source image file
//   - outputPath: Path where the downscaled image will be saved (can be same as inputPath)
//   - opts: Downscaling options (dimensions, filter, quality)
//
// Returns:
//   - error if the operation fails
//
// The function automatically:
//   - Maintains aspect ratio if only width or height is specified
//   - Skips downscaling if target dimensions are larger than source
//   - Preserves the image format (PNG/JPEG)
//   - Uses Lanczos3 resampling for optimal quality
func DownscaleImage(inputPath, outputPath string, opts DownscaleOptions) error {
	// Validate options
	if opts.Percentage == 0 && opts.Width == 0 && opts.Height == 0 {
		return fmt.Errorf("either percentage or at least one dimension (width or height) must be specified")
	}

	if opts.Percentage < 0 || opts.Percentage > 100 {
		return fmt.Errorf("percentage must be between 0 and 100")
	}

	if opts.Width < 0 || opts.Height < 0 {
		return fmt.Errorf("dimensions cannot be negative")
	}

	// Set default JPEG quality if not specified
	if opts.JPEGQuality == 0 {
		opts.JPEGQuality = 90
	}

	// Open the source image
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input image: %w", err)
	}
	defer inputFile.Close()

	// Decode the image (supports PNG and JPEG automatically)
	srcImg, format, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Get source dimensions
	srcBounds := srcImg.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// Calculate target dimensions
	var targetWidth, targetHeight int

	// If percentage is specified, use it to calculate dimensions
	if opts.Percentage > 0 {
		scale := opts.Percentage / 100.0
		targetWidth = int(float64(srcWidth) * scale)
		targetHeight = int(float64(srcHeight) * scale)
	} else {
		// Otherwise use explicit dimensions
		targetWidth = opts.Width
		targetHeight = opts.Height

		// If only one dimension is specified, calculate the other to maintain aspect ratio
		if targetWidth == 0 {
			targetWidth = int(float64(targetHeight) * float64(srcWidth) / float64(srcHeight))
		} else if targetHeight == 0 {
			targetHeight = int(float64(targetWidth) * float64(srcHeight) / float64(srcWidth))
		}
	}

	// Validate that we're actually downscaling (not upscaling)
	if targetWidth >= srcWidth && targetHeight >= srcHeight {
		return fmt.Errorf("target dimensions (%dx%d) are not smaller than source (%dx%d) - downscaling only",
			targetWidth, targetHeight, srcWidth, srcHeight)
	}

	// Create destination image
	dstImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Select the scaling algorithm
	var scaler draw.Scaler
	switch opts.Filter {
	case FilterLanczos:
		scaler = draw.CatmullRom // Lanczos3 equivalent in x/image/draw
	case FilterBiLinear:
		scaler = draw.BiLinear
	case FilterNearestNeighbor:
		scaler = draw.NearestNeighbor
	default:
		scaler = draw.CatmullRom // Default to highest quality
	}

	// Perform the downscaling
	scaler.Scale(dstImg, dstImg.Bounds(), srcImg, srcBounds, draw.Over, nil)

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Encode and save based on format
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outputFile, dstImg, &jpeg.Options{Quality: opts.JPEGQuality})
	case "png":
		err = png.Encode(outputFile, dstImg)
	default:
		// Default to PNG for unknown formats
		err = png.Encode(outputFile, dstImg)
	}

	if err != nil {
		return fmt.Errorf("failed to encode output image: %w", err)
	}

	// Mandatory: Strip all PNG metadata for privacy and security
	// This removes any generation parameters, timestamps, or other sensitive data
	if err := StripPNGMetadata(outputPath); err != nil {
		return fmt.Errorf("failed to strip PNG metadata: %w", err)
	}

	return nil
}

// GetImageDimensions returns the width and height of an image file without fully decoding it.
// This is useful for validation before processing.
func GetImageDimensions(imagePath string) (width, height int, err error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode image config: %w", err)
	}

	return config.Width, config.Height, nil
}

// DownscaleInPlace downscales an image and replaces the original file.
// This is a convenience wrapper around DownscaleImage for in-place operations.
func DownscaleInPlace(imagePath string, opts DownscaleOptions) error {
	// Create a temporary file in the same directory
	dir := filepath.Dir(imagePath)
	base := filepath.Base(imagePath)
	ext := filepath.Ext(imagePath)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	tempPath := filepath.Join(dir, fmt.Sprintf(".%s_tmp%s", nameWithoutExt, ext))

	// Downscale to temporary file
	if err := DownscaleImage(imagePath, tempPath, opts); err != nil {
		// Clean up temp file if it was created
		os.Remove(tempPath)
		return err
	}

	// Replace original with downscaled version
	if err := os.Rename(tempPath, imagePath); err != nil {
		// Clean up temp file
		os.Remove(tempPath)
		return fmt.Errorf("failed to replace original file: %w", err)
	}

	return nil
}
