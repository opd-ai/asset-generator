package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/opd-ai/asset-generator/pkg/converter"
	"github.com/spf13/cobra"
)

var (
	convertMethod          string
	convertOutput          string
	convertPrimitiveShapes int
	convertPrimitiveMode   int
	convertPrimitiveAlpha  int
	convertPrimitiveRepeat int
	convertGotraceArgs     []string
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert images to different formats",
	Long: `Convert images to different formats, primarily focusing on SVG conversion.

Supports two conversion methods:
  - primitive: Uses geometric shapes to approximate the image (fogleman/primitive)
  - gotrace:   Uses edge tracing for vector conversion (requires potrace installed)

Examples:
  # Convert image to SVG using primitive method with 100 shapes
  asset-generator convert svg input.png --method primitive --shapes 100

  # Convert using gotrace method
  asset-generator convert svg input.png --method gotrace

  # Specify custom output path
  asset-generator convert svg input.png -o output.svg --shapes 200`,
}

// convertSvgCmd represents the svg conversion command
var convertSvgCmd = &cobra.Command{
	Use:   "svg <input-file>",
	Short: "Convert an image to SVG format",
	Long: `Convert an image to SVG (Scalable Vector Graphics) format.

Two conversion methods are available:

1. Primitive (default): Uses geometric shapes to approximate the image
   - Fast and produces clean, artistic results
   - Good for logos, illustrations, and simplified graphics
   - Control the quality with --shapes flag

2. Gotrace: Uses edge tracing to convert bitmap to vector
   - Requires 'potrace' to be installed on your system
   - Better for preserving fine details
   - Good for line art and high-contrast images

Examples:
  # Basic conversion with 100 shapes
  asset-generator convert svg input.png

  # High quality with 500 shapes
  asset-generator convert svg input.png --shapes 500

  # Use different shape modes (1=triangle, 2=rect, 3=ellipse, etc)
  asset-generator convert svg input.png --mode 3 --shapes 200

  # Use gotrace method
  asset-generator convert svg input.png --method gotrace

  # Specify output location
  asset-generator convert svg input.png -o /path/to/output.svg`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertSvg,
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.AddCommand(convertSvgCmd)

	// Conversion method flags
	convertSvgCmd.Flags().StringVarP(&convertMethod, "method", "m", "primitive",
		"Conversion method: 'primitive' or 'gotrace'")
	convertSvgCmd.Flags().StringVarP(&convertOutput, "output", "o", "",
		"Output file path (defaults to input filename with .svg extension)")

	// Primitive-specific flags
	convertSvgCmd.Flags().IntVar(&convertPrimitiveShapes, "shapes", 100,
		"Number of shapes to use (primitive method)")
	convertSvgCmd.Flags().IntVar(&convertPrimitiveMode, "mode", 1,
		"Shape mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect, 6=beziers, 7=rotatedellipse, 8=polygon")
	convertSvgCmd.Flags().IntVar(&convertPrimitiveAlpha, "alpha", 128,
		"Alpha value for shapes (0-255)")
	convertSvgCmd.Flags().IntVar(&convertPrimitiveRepeat, "repeat", 0,
		"Number of optimization repeats")

	// Gotrace-specific flags
	convertSvgCmd.Flags().StringSliceVar(&convertGotraceArgs, "gotrace-args", []string{},
		"Additional arguments to pass to potrace (gotrace method)")
}

func runConvertSvg(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Validate input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Determine output path
	outputPath := convertOutput
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = inputPath[:len(inputPath)-len(ext)] + ".svg"
	}

	// Validate conversion method
	var convMethod converter.SVGConversionMethod
	switch convertMethod {
	case "primitive":
		convMethod = converter.MethodPrimitive
	case "gotrace":
		convMethod = converter.MethodGotrace
	default:
		return fmt.Errorf("unknown conversion method: %s (use 'primitive' or 'gotrace')", convertMethod)
	}

	// Create converter
	svgConverter := converter.NewSVGConverter()

	// Set up conversion options
	opts := converter.ConversionOptions{
		Method:          convMethod,
		OutputPath:      outputPath,
		PrimitiveShapes: convertPrimitiveShapes,
		PrimitiveMode:   convertPrimitiveMode,
		PrimitiveAlpha:  convertPrimitiveAlpha,
		PrimitiveRepeat: convertPrimitiveRepeat,
		GotraceArgs:     convertGotraceArgs,
	}

	// Provide feedback
	if !quiet {
		fmt.Printf("Converting %s to SVG using %s method...\n", inputPath, convertMethod)
		if convMethod == converter.MethodPrimitive {
			fmt.Printf("  Shapes: %d\n", convertPrimitiveShapes)
			fmt.Printf("  Mode: %d\n", convertPrimitiveMode)
		}
	}

	// Perform conversion
	result, err := svgConverter.ConvertToSVG(inputPath, opts)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Success message
	if !quiet {
		fmt.Printf("✓ Successfully converted to: %s\n", result)

		// Show file size
		info, err := os.Stat(result)
		if err == nil {
			fmt.Printf("  Output size: %.2f KB\n", float64(info.Size())/1024)
		}
	}

	return nil
}
