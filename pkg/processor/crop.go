package processor

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// CropOptions configures automatic cropping behavior
type CropOptions struct {
	// Threshold for detecting "whitespace" (0-255, default: 250)
	// Pixels with all RGB values >= this threshold are considered whitespace
	Threshold uint8
	// Tolerance for near-white colors (0-255, default: 10)
	// Allows RGB values within Threshold ± Tolerance to be considered whitespace
	Tolerance uint8
	// Quality for JPEG output (1-100, default: 90)
	JPEGQuality int
	// PreserveAspectRatio ensures the crop maintains the original aspect ratio
	PreserveAspectRatio bool
}

// AutoCropImage automatically detects and removes excess whitespace from image edges
// while optionally preserving the aspect ratio.
//
// Parameters:
//   - inputPath: Path to the source image file
//   - outputPath: Path where the cropped image will be saved (can be same as inputPath)
//   - opts: Cropping options (threshold, tolerance, quality)
//
// Returns:
//   - error if the operation fails
//
// The function:
//   - Detects whitespace by scanning from edges inward
//   - Preserves aspect ratio if requested
//   - Maintains the original image format (PNG/JPEG)
//   - Returns error if entire image would be cropped (no content detected)
func AutoCropImage(inputPath, outputPath string, opts CropOptions) error {
	// Set defaults
	if opts.Threshold == 0 {
		opts.Threshold = 250 // Default: very light colors are "whitespace"
	}
	if opts.Tolerance == 0 {
		opts.Tolerance = 10 // Default: allow some variation
	}
	if opts.JPEGQuality == 0 {
		opts.JPEGQuality = 90
	}

	// Open the source image
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input image: %w", err)
	}
	defer inputFile.Close()

	// Decode the image
	srcImg, format, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Detect content bounds
	bounds := detectContentBounds(srcImg, opts.Threshold, opts.Tolerance)

	// Validate that we found some content
	if bounds.Empty() {
		return fmt.Errorf("no content detected in image (entire image appears to be whitespace)")
	}

	// Check if crop would actually change the image
	srcBounds := srcImg.Bounds()
	if bounds.Eq(srcBounds) {
		// No cropping needed - image has no whitespace borders
		// If input and output are different, copy the file
		if inputPath != outputPath {
			return copyFile(inputPath, outputPath)
		}
		return nil
	}

	// Apply aspect ratio preservation if requested
	if opts.PreserveAspectRatio {
		bounds = preserveAspectRatio(srcBounds, bounds)
	}

	// Create cropped image
	croppedImg := cropImage(srcImg, bounds)

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Encode and save based on format
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outputFile, croppedImg, &jpeg.Options{Quality: opts.JPEGQuality})
	case "png":
		err = png.Encode(outputFile, croppedImg)
	default:
		// Default to PNG for unknown formats
		err = png.Encode(outputFile, croppedImg)
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

// detectContentBounds scans from each edge inward to find the first non-whitespace pixels
func detectContentBounds(img image.Image, threshold, tolerance uint8) image.Rectangle {
	bounds := img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	// Scan from left
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if !isWhitespace(img.At(x, y), threshold, tolerance) {
				minX = x
				goto foundLeft
			}
		}
	}
foundLeft:

	// Scan from right
	for x := bounds.Max.X - 1; x >= bounds.Min.X; x-- {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if !isWhitespace(img.At(x, y), threshold, tolerance) {
				maxX = x + 1 // exclusive bound
				goto foundRight
			}
		}
	}
foundRight:

	// Scan from top
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if !isWhitespace(img.At(x, y), threshold, tolerance) {
				minY = y
				goto foundTop
			}
		}
	}
foundTop:

	// Scan from bottom
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if !isWhitespace(img.At(x, y), threshold, tolerance) {
				maxY = y + 1 // exclusive bound
				goto foundBottom
			}
		}
	}
foundBottom:

	// Return the content bounds
	return image.Rect(minX, minY, maxX, maxY)
}

// isWhitespace determines if a pixel is considered whitespace based on threshold and tolerance
func isWhitespace(c color.Color, threshold, tolerance uint8) bool {
	r, g, b, a := c.RGBA()
	// Convert from 16-bit to 8-bit
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)
	a8 := uint8(a >> 8)

	// Check for full transparency (also considered whitespace for PNG with alpha)
	if a8 < 10 {
		return true
	}

	// Calculate the acceptable range
	minVal := threshold - tolerance
	if minVal > threshold { // Handle underflow
		minVal = 0
	}

	// Check if all RGB components are within the whitespace range
	return r8 >= minVal && g8 >= minVal && b8 >= minVal
}

// preserveAspectRatio adjusts crop bounds to maintain the original aspect ratio
func preserveAspectRatio(srcBounds, cropBounds image.Rectangle) image.Rectangle {
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()
	cropWidth := cropBounds.Dx()
	cropHeight := cropBounds.Dy()

	// Calculate aspect ratios
	srcAspect := float64(srcWidth) / float64(srcHeight)
	cropAspect := float64(cropWidth) / float64(cropHeight)

	// If aspect ratios are already similar (within 1%), no adjustment needed
	aspectDiff := cropAspect / srcAspect
	if aspectDiff > 0.99 && aspectDiff < 1.01 {
		return cropBounds
	}

	// Determine which dimension to adjust
	if cropAspect > srcAspect {
		// Crop is too wide, need to increase height
		targetHeight := int(float64(cropWidth) / srcAspect)
		heightDiff := targetHeight - cropHeight

		// Expand vertically, centered
		newMinY := cropBounds.Min.Y - heightDiff/2
		newMaxY := cropBounds.Max.Y + heightDiff/2 + heightDiff%2

		// Clamp to source bounds
		if newMinY < srcBounds.Min.Y {
			newMinY = srcBounds.Min.Y
			newMaxY = newMinY + targetHeight
		}
		if newMaxY > srcBounds.Max.Y {
			newMaxY = srcBounds.Max.Y
			newMinY = newMaxY - targetHeight
		}

		return image.Rect(cropBounds.Min.X, newMinY, cropBounds.Max.X, newMaxY)
	} else {
		// Crop is too tall, need to increase width
		targetWidth := int(float64(cropHeight) * srcAspect)
		widthDiff := targetWidth - cropWidth

		// Expand horizontally, centered
		newMinX := cropBounds.Min.X - widthDiff/2
		newMaxX := cropBounds.Max.X + widthDiff/2 + widthDiff%2

		// Clamp to source bounds
		if newMinX < srcBounds.Min.X {
			newMinX = srcBounds.Min.X
			newMaxX = newMinX + targetWidth
		}
		if newMaxX > srcBounds.Max.X {
			newMaxX = srcBounds.Max.X
			newMinX = newMaxX - targetWidth
		}

		return image.Rect(newMinX, cropBounds.Min.Y, newMaxX, cropBounds.Max.Y)
	}
}

// cropImage extracts a sub-image from the source using the specified bounds
func cropImage(src image.Image, bounds image.Rectangle) image.Image {
	// Check if the source image supports SubImage (most do)
	if subImager, ok := src.(interface {
		SubImage(r image.Rectangle) image.Image
	}); ok {
		return subImager.SubImage(bounds)
	}

	// Fallback: manually copy pixels
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x-bounds.Min.X, y-bounds.Min.Y, src.At(x, y))
		}
	}
	return dst
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}
	return nil
}

// AutoCropInPlace crops an image and replaces the original file.
// This is a convenience wrapper around AutoCropImage for in-place operations.
func AutoCropInPlace(imagePath string, opts CropOptions) error {
	// Create a temporary file in the same directory
	dir := filepath.Dir(imagePath)
	base := filepath.Base(imagePath)
	ext := filepath.Ext(imagePath)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	tempPath := filepath.Join(dir, fmt.Sprintf(".%s_crop_tmp%s", nameWithoutExt, ext))

	// Crop to temporary file
	if err := AutoCropImage(imagePath, tempPath, opts); err != nil {
		// Clean up temp file if it was created
		os.Remove(tempPath)
		return err
	}

	// Replace original with cropped version
	if err := os.Rename(tempPath, imagePath); err != nil {
		// Clean up temp file
		os.Remove(tempPath)
		return fmt.Errorf("failed to replace original file: %w", err)
	}

	return nil
}
