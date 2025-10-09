package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/opd-ai/asset-generator/pkg/processor"
	"github.com/spf13/cobra"
)

var (
	cropThreshold      int
	cropTolerance      int
	cropPreserveAspect bool
	cropQuality        int
	cropOutput         string
	cropInPlace        bool
)

// cropCmd represents the crop command
var cropCmd = &cobra.Command{
	Use:   "crop [image-file...]",
	Short: "Auto-crop whitespace borders from images",
	Long: `Automatically detect and remove excess whitespace from image edges.

The crop command scans from each edge inward to find content boundaries and
removes whitespace borders while optionally preserving the original aspect ratio.

This is useful for:
  - Removing padding from generated images
  - Cleaning up screenshots with extra borders
  - Preparing images for further processing or display
  - Batch processing multiple images with consistent cropping

Examples:
  # Basic auto-crop of a single image
  asset-generator crop image.png

  # Auto-crop preserving original aspect ratio
  asset-generator crop photo.jpg --preserve-aspect

  # Auto-crop in-place (replaces original file)
  asset-generator crop image.png --in-place

  # Auto-crop with specific output path
  asset-generator crop input.png --output cropped.png

  # Batch crop multiple images
  asset-generator crop *.jpg --in-place

  # Adjust sensitivity for darker backgrounds
  asset-generator crop image.png --threshold 200 --tolerance 20

  # More aggressive whitespace detection
  asset-generator crop image.png --threshold 240 --tolerance 5

Whitespace Detection:
  The tool identifies pixels as "whitespace" when all RGB values are above
  (threshold - tolerance). Default settings (threshold=250, tolerance=10)
  detect very light colors as whitespace.

  - Higher threshold (250-255): Only pure white considered whitespace
  - Lower threshold (200-240): Light grays also considered whitespace
  - Higher tolerance (15-30): More forgiving, catches near-white colors
  - Lower tolerance (5-10): Stricter, only very close to threshold

Aspect Ratio Preservation:
  When --preserve-aspect is enabled, the crop bounds are expanded (if needed)
  to match the original image's aspect ratio. This ensures the cropped image
  maintains the same width:height proportions as the original.`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCrop,
}

func init() {
	rootCmd.AddCommand(cropCmd)

	cropCmd.Flags().IntVar(&cropThreshold, "threshold", 250, "whitespace detection threshold (0-255)")
	cropCmd.Flags().IntVar(&cropTolerance, "tolerance", 10, "tolerance for near-white colors (0-255)")
	cropCmd.Flags().BoolVar(&cropPreserveAspect, "preserve-aspect", false, "preserve original aspect ratio")
	cropCmd.Flags().IntVar(&cropQuality, "quality", 90, "JPEG quality (1-100)")
	cropCmd.Flags().StringVarP(&cropOutput, "output", "o", "", "output file path (single file mode only)")
	cropCmd.Flags().BoolVarP(&cropInPlace, "in-place", "i", false, "replace original file(s) with cropped version")
}

func runCrop(cmd *cobra.Command, args []string) error {
	// Validate threshold and tolerance
	if cropThreshold < 0 || cropThreshold > 255 {
		return fmt.Errorf("threshold must be between 0 and 255")
	}
	if cropTolerance < 0 || cropTolerance > 255 {
		return fmt.Errorf("tolerance must be between 0 and 255")
	}
	if cropQuality < 1 || cropQuality > 100 {
		return fmt.Errorf("quality must be between 1 and 100")
	}

	// Check for conflicting options
	if len(args) > 1 && cropOutput != "" {
		return fmt.Errorf("--output can only be used with a single input file")
	}

	if cropOutput != "" && cropInPlace {
		return fmt.Errorf("cannot use both --output and --in-place")
	}

	// Build crop options
	opts := processor.CropOptions{
		Threshold:           uint8(cropThreshold),
		Tolerance:           uint8(cropTolerance),
		JPEGQuality:         cropQuality,
		PreserveAspectRatio: cropPreserveAspect,
	}

	successCount := 0
	failCount := 0

	for _, imagePath := range args {
		// Verify file exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: file not found: %s\n", imagePath)
			failCount++
			continue
		}

		// Determine output path
		var outputPath string
		if cropInPlace || cropOutput == "" && len(args) == 1 {
			// In-place mode
			if verbose {
				fmt.Fprintf(os.Stderr, "Cropping: %s (in-place)\n", imagePath)
			}
			if err := processor.AutoCropInPlace(imagePath, opts); err != nil {
				fmt.Fprintf(os.Stderr, "Error cropping %s: %v\n", imagePath, err)
				failCount++
				continue
			}
			outputPath = imagePath
		} else {
			// Output to different file
			if cropOutput != "" {
				outputPath = cropOutput
			} else {
				// Generate output filename (add -cropped suffix)
				ext := filepath.Ext(imagePath)
				nameWithoutExt := imagePath[:len(imagePath)-len(ext)]
				outputPath = fmt.Sprintf("%s-cropped%s", nameWithoutExt, ext)
			}

			if verbose {
				fmt.Fprintf(os.Stderr, "Cropping: %s -> %s\n", imagePath, outputPath)
			}

			if err := processor.AutoCropImage(imagePath, outputPath, opts); err != nil {
				fmt.Fprintf(os.Stderr, "Error cropping %s: %v\n", imagePath, err)
				failCount++
				continue
			}
		}

		successCount++
		if !quiet {
			// Get dimensions to show the result
			width, height, err := processor.GetImageDimensions(outputPath)
			if err == nil {
				fmt.Fprintf(os.Stderr, "✓ Cropped: %s (%dx%d)\n", outputPath, width, height)
			} else {
				fmt.Fprintf(os.Stderr, "✓ Cropped: %s\n", outputPath)
			}
		}
	}

	// Summary
	if !quiet && len(args) > 1 {
		fmt.Fprintf(os.Stderr, "\nCrop complete: %d succeeded, %d failed\n", successCount, failCount)
	}

	if failCount > 0 {
		return fmt.Errorf("failed to crop %d image(s)", failCount)
	}

	return nil
}
