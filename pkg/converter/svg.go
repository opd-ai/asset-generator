package converter

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dennwc/gotrace"
	"github.com/fogleman/primitive/primitive"
)

// SVGConverter provides methods to convert images to SVG format
type SVGConverter struct {
	// PrimitiveShapes is the number of shapes to use with primitive
	PrimitiveShapes int
	// PrimitiveMode is the shape mode (0=combo, 1=triangle, 2=rect, etc)
	PrimitiveMode int
	// PrimitiveAlpha is the alpha value (0-255)
	PrimitiveAlpha int
	// PrimitiveRepeat is the number of times to repeat the primitive process
	PrimitiveRepeat int
}

// SVGConversionMethod represents the conversion method to use
type SVGConversionMethod string

const (
	// MethodPrimitive uses fogleman/primitive for geometric shape approximation
	MethodPrimitive SVGConversionMethod = "primitive"
	// MethodGotrace uses dennwc/gotrace for edge tracing
	MethodGotrace SVGConversionMethod = "gotrace"
)

// ConversionOptions holds options for SVG conversion
type ConversionOptions struct {
	Method          SVGConversionMethod
	OutputPath      string
	PrimitiveShapes int // Number of shapes for primitive method
	PrimitiveMode   int // Shape mode for primitive
	PrimitiveAlpha  int // Alpha value for primitive
	PrimitiveRepeat int // Number of optimization repeats
}

// NewSVGConverter creates a new SVG converter with default settings
func NewSVGConverter() *SVGConverter {
	return &SVGConverter{
		PrimitiveShapes: 100,
		PrimitiveMode:   1, // triangles
		PrimitiveAlpha:  128,
		PrimitiveRepeat: 0,
	}
}

// ConvertToSVG converts an image file to SVG using the specified method
func (c *SVGConverter) ConvertToSVG(inputPath string, opts ConversionOptions) (string, error) {
	// Validate input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Determine output path if not specified
	outputPath := opts.OutputPath
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = strings.TrimSuffix(inputPath, ext) + ".svg"
	}

	// Convert based on method
	switch opts.Method {
	case MethodPrimitive:
		return c.convertWithPrimitive(inputPath, outputPath, opts)
	case MethodGotrace:
		return c.convertWithGotrace(inputPath, outputPath, opts)
	default:
		return "", fmt.Errorf("unknown conversion method: %s (supported: 'primitive', 'gotrace')", opts.Method)
	}
}

// convertWithPrimitive uses fogleman/primitive library for conversion
func (c *SVGConverter) convertWithPrimitive(inputPath, outputPath string, opts ConversionOptions) (string, error) {
	// Load the input image
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open input image: %w", err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode input image: %w", err)
	}

	// Set up primitive parameters
	shapes := opts.PrimitiveShapes
	if shapes == 0 {
		shapes = c.PrimitiveShapes
	}

	mode := opts.PrimitiveMode
	if mode == 0 {
		mode = c.PrimitiveMode
	}

	alpha := opts.PrimitiveAlpha
	if alpha == 0 {
		alpha = c.PrimitiveAlpha
	}

	repeat := opts.PrimitiveRepeat
	if repeat == 0 {
		repeat = c.PrimitiveRepeat
	}

	// Create background color from average image color
	bg := primitive.MakeColor(primitive.AverageImageColor(img))

	// Create model
	model := primitive.NewModel(img, bg, 1024, runtime.NumCPU())

	// Add shapes iteratively
	for i := 0; i < shapes; i++ {
		// Add shape
		model.Step(primitive.ShapeType(mode), alpha, repeat)
	}

	// Export to SVG
	svg := model.SVG()

	// Write SVG to file
	if err := os.WriteFile(outputPath, []byte(svg), 0644); err != nil {
		return "", fmt.Errorf("failed to write SVG file: %w", err)
	}

	return outputPath, nil
}

// convertWithGotrace uses dennwc/gotrace library for edge tracing conversion
func (c *SVGConverter) convertWithGotrace(inputPath, outputPath string, opts ConversionOptions) (string, error) {
	// Load the input image
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open input image: %w", err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode input image: %w", err)
	}

	// Convert image to bitmap for tracing using default alpha threshold
	bm := gotrace.NewBitmapFromImage(img, nil)

	// Trace the bitmap to get paths
	paths, err := gotrace.Trace(bm, nil) // nil uses default parameters
	if err != nil {
		return "", fmt.Errorf("failed to trace bitmap: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Write SVG
	if err := gotrace.WriteSvg(outFile, img.Bounds(), paths, ""); err != nil {
		return "", fmt.Errorf("failed to write SVG: %w", err)
	}

	return outputPath, nil
}

// ConvertWithPrimitiveDefault is a convenience function for quick primitive conversion
func ConvertWithPrimitiveDefault(inputPath, outputPath string, shapes int) (string, error) {
	converter := NewSVGConverter()
	opts := ConversionOptions{
		Method:          MethodPrimitive,
		OutputPath:      outputPath,
		PrimitiveShapes: shapes,
	}
	return converter.ConvertToSVG(inputPath, opts)
}
