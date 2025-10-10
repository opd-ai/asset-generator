package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opd-ai/asset-generator/pkg/processor"
	"github.com/spf13/cobra"
)

var (
	downscaleWidth      int
	downscaleHeight     int
	downscalePercentage float64
	downscaleFilter     string
	downscaleQuality    int
	downscaleOutput     string
	downscaleInPlace    bool
)

// downscaleCmd represents the downscale command
var downscaleCmd = &cobra.Command{
	Use:   "downscale [image-file...]",
	Short: "Downscale images using high-quality filters",
	Long: `Downscale one or more images using high-quality resampling algorithms.

The downscale command uses advanced filtering to reduce image dimensions while
preserving quality. By default, it uses Lanczos3 resampling, which provides
the best quality for downscaling operations.

Examples:
  # Downscale a single image to 1024px width (maintains aspect ratio)
  asset-generator downscale image.png --width 1024

  # Downscale to 50% of original size
  asset-generator downscale image.png --percentage 50

  # Downscale to specific dimensions
  asset-generator downscale photo.jpg --width 800 --height 600

  # Downscale in-place (replaces original file)
  asset-generator downscale image.png --width 512 --in-place

  # Downscale with specific output path
  asset-generator downscale input.png --width 1024 --output resized.png

  # Batch downscale multiple images by 75%
  asset-generator downscale *.jpg --percentage 75 --in-place

  # Use different filter for speed
  asset-generator downscale large.png --width 512 --filter bilinear

Filter Options:
  lanczos   - Highest quality, best for photographs (default)
  bilinear  - Good balance of speed and quality
  nearest   - Fastest, best for pixel art or icons

The command will:
  - Automatically maintain aspect ratio if only one dimension is specified
  - Prevent accidental upscaling (will error if target > source)
  - Create output directory if needed (when using --output)
  - Preserve image format (PNG/JPEG)
  - Use quality setting for JPEG output (default: 90)`,
	Args: cobra.MinimumNArgs(1),
	RunE: runDownscale,
}

func init() {
	rootCmd.AddCommand(downscaleCmd)

	downscaleCmd.Flags().IntVarP(&downscaleWidth, "width", "w", 0, "target width in pixels (0=auto from height)")
	downscaleCmd.Flags().IntVarP(&downscaleHeight, "height", "l", 0, "target height in pixels (0=auto from width)")
	downscaleCmd.Flags().Float64VarP(&downscalePercentage, "percentage", "p", 0, "scale by percentage (1-100, 0=use width/height instead)")
	downscaleCmd.Flags().StringVar(&downscaleFilter, "filter", "lanczos", "resampling filter: lanczos, bilinear, nearest")
	downscaleCmd.Flags().IntVar(&downscaleQuality, "quality", 90, "JPEG quality (1-100)")
	downscaleCmd.Flags().StringVar(&downscaleOutput, "output-file", "", "output file path (single file mode only)")
	downscaleCmd.Flags().BoolVar(&downscaleInPlace, "in-place", false, "replace original file(s) with downscaled version")
}

func runDownscale(cmd *cobra.Command, args []string) error {
	// Validate that either percentage OR at least one dimension is specified
	if downscalePercentage == 0 && downscaleWidth == 0 && downscaleHeight == 0 {
		return fmt.Errorf("either --percentage or at least one dimension (--width or --height) must be specified")
	}

	// Validate that percentage and dimensions are not used together
	if downscalePercentage > 0 && (downscaleWidth > 0 || downscaleHeight > 0) {
		return fmt.Errorf("cannot specify both --percentage and explicit dimensions (--width/--height)")
	}

	// Validate percentage range
	if downscalePercentage < 0 || downscalePercentage > 100 {
		return fmt.Errorf("percentage must be between 0 and 100")
	}

	// Validate dimensions
	if downscaleWidth < 0 || downscaleHeight < 0 {
		return fmt.Errorf("dimensions cannot be negative")
	}

	// Validate quality
	if downscaleQuality < 1 || downscaleQuality > 100 {
		return fmt.Errorf("quality must be between 1 and 100")
	}

	// Validate filter
	validFilters := []string{"lanczos", "bilinear", "nearest"}
	filterValid := false
	filterLower := strings.ToLower(downscaleFilter)
	for _, f := range validFilters {
		if filterLower == f {
			filterValid = true
			break
		}
	}
	if !filterValid {
		return fmt.Errorf("invalid filter '%s' (valid options: %s)", downscaleFilter, strings.Join(validFilters, ", "))
	}

	// Validate output flag usage
	if downscaleOutput != "" && len(args) > 1 {
		return fmt.Errorf("--output can only be used with a single input file")
	}

	if downscaleOutput != "" && downscaleInPlace {
		return fmt.Errorf("cannot specify both --output and --in-place")
	}

	// Map filter name to processor type
	var filterType processor.ResizeFilter
	switch filterLower {
	case "lanczos":
		filterType = processor.FilterLanczos
	case "bilinear":
		filterType = processor.FilterBiLinear
	case "nearest":
		filterType = processor.FilterNearestNeighbor
	}

	// Build downscale options
	opts := processor.DownscaleOptions{
		Width:       downscaleWidth,
		Height:      downscaleHeight,
		Percentage:  downscalePercentage,
		Filter:      filterType,
		JPEGQuality: downscaleQuality,
	}

	// Process each input file
	processedCount := 0
	errorCount := 0

	for _, inputPath := range args {
		// Check if file exists
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			if !quiet {
				fmt.Fprintf(os.Stderr, "Error: File not found: %s\n", inputPath)
			}
			errorCount++
			continue
		}

		// Determine output path
		var outputPath string
		if downscaleInPlace {
			outputPath = inputPath
		} else if downscaleOutput != "" {
			outputPath = downscaleOutput
		} else {
			// Generate default output path: input_downscaled.ext
			ext := filepath.Ext(inputPath)
			nameWithoutExt := strings.TrimSuffix(inputPath, ext)
			outputPath = fmt.Sprintf("%s_downscaled%s", nameWithoutExt, ext)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Processing: %s -> %s\n", inputPath, outputPath)
			if downscalePercentage > 0 {
				fmt.Fprintf(os.Stderr, "  Scale: %.0f%% (filter: %s)\n", downscalePercentage, filterLower)
			} else {
				fmt.Fprintf(os.Stderr, "  Target dimensions: %dx%d (filter: %s)\n", 
					downscaleWidth, downscaleHeight, filterLower)
			}
		}

		// Perform downscaling
		var err error
		if downscaleInPlace {
			err = processor.DownscaleInPlace(inputPath, opts)
		} else {
			err = processor.DownscaleImage(inputPath, outputPath, opts)
		}

		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", inputPath, err)
			}
			errorCount++
			continue
		}

		processedCount++
		if !quiet {
			if downscaleInPlace {
				fmt.Fprintf(os.Stderr, "✓ Downscaled: %s\n", inputPath)
			} else {
				fmt.Fprintf(os.Stderr, "✓ Downscaled: %s -> %s\n", inputPath, outputPath)
			}
		}
	}

	// Summary
	if !quiet && len(args) > 1 {
		fmt.Fprintf(os.Stderr, "\nProcessed %d/%d images successfully\n", processedCount, len(args))
	}

	if errorCount > 0 {
		return fmt.Errorf("%d file(s) failed to process", errorCount)
	}

	return nil
}
