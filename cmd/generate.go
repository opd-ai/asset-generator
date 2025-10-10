package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/opd-ai/asset-generator/pkg/client"
	"github.com/opd-ai/asset-generator/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	generatePrompt           string
	generateModel            string
	generateSteps            int
	generateWidth            int
	generateHeight           int
	generateSeed             int64
	generateBatchSize        int
	generateCfgScale         float64
	generateNegPrompt        string
	generateSampler          string
	generateUseWebSocket     bool   // Enable WebSocket for real-time progress
	generateSaveImages       bool   // Download and save images locally
	generateOutputDir        string // Directory to save downloaded images
	generateFilenameTemplate string // Template for custom filenames
	// SkimmedCFG (Distilled CFG) options
	generateSkimmedCFG      bool    // Enable Skimmed CFG for improved quality/speed
	generateSkimmedCFGScale float64 // Skimmed CFG scale value
	generateSkimmedCFGStart float64 // Start percentage for Skimmed CFG (0-1)
	generateSkimmedCFGEnd   float64 // End percentage for Skimmed CFG (0-1)
	// Auto-crop postprocessing options
	generateAutoCrop               bool // Enable auto-crop to remove whitespace
	generateAutoCropThreshold      int  // Whitespace detection threshold (0-255)
	generateAutoCropTolerance      int  // Tolerance for near-white colors (0-255)
	generateAutoCropPreserveAspect bool // Preserve aspect ratio when cropping
	// Downscale postprocessing options
	generateDownscaleWidth      int     // Target width for postprocessing downscale
	generateDownscaleHeight     int     // Target height for postprocessing downscale
	generateDownscalePercentage float64 // Scale by percentage
	generateDownscaleFilter     string  // Downscaling algorithm (lanczos, bilinear, nearest)
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate assets using AI models",
	Long: `Generate various types of assets using AI asset generation APIs.
Supports image generation with various models and parameters.`,
}

// generateImageCmd represents the image generation command
var generateImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate an image",
	Long: `Generate an image using AI text-to-image capabilities.

Examples:
  # Basic generation
  asset-generator generate image --prompt "a beautiful landscape"
  
  # Advanced generation with parameters
  asset-generator generate image \
    --prompt "futuristic city at sunset" \
    --model "stable-diffusion-xl" \
    --width 1024 --length 1024 \
    --steps 30 --cfg-scale 7.5
  
  # Download and save generated images locally
  asset-generator generate image \
    --prompt "cat wearing sunglasses" \
    --save-images --output-dir ./my-images
  
  # Custom filename template
  asset-generator generate image \
    --prompt "fantasy landscape" \
    --batch 5 --save-images \
    --filename-template "landscape-{index}-{seed}.png"
  
  # Auto-crop whitespace borders from generated images
  asset-generator generate image \
    --prompt "centered logo design" \
    --save-images --auto-crop
  
  # Auto-crop with aspect ratio preservation
  asset-generator generate image \
    --prompt "product photo" \
    --save-images --auto-crop --auto-crop-preserve-aspect
  
  # Combined postprocessing: crop then downscale
  asset-generator generate image \
    --prompt "high resolution art" \
    --width 2048 --length 2048 \
    --save-images --auto-crop \
    --downscale-width 1024 --downscale-filter lanczos
  
  # Downscale images after download (postprocessing)
  asset-generator generate image \
    --prompt "high resolution art" \
    --width 2048 --length 2048 \
    --save-images --downscale-width 1024 \
    --downscale-filter lanczos
  
  # Use Skimmed CFG for improved quality and faster generation
  asset-generator generate image \
    --prompt "detailed portrait" \
    --skimmed-cfg --skimmed-cfg-scale 3.0
  
  # Skimmed CFG with custom range (apply only during middle of generation)
  asset-generator generate image \
    --prompt "landscape painting" \
    --skimmed-cfg --skimmed-cfg-start 0.2 --skimmed-cfg-end 0.8
  
  # Save metadata to specific file
  asset-generator generate image \
    --prompt "cat wearing sunglasses" \
    --output my-cat.json

Postprocessing Pipeline (applied in order):
  1. Auto-crop: Removes whitespace borders while optionally preserving aspect ratio
  2. Downscale: Reduces image dimensions using high-quality filtering

Filename Template Placeholders:
  {index}, {i}     - Zero-padded index (001, 002, ...)
  {index1}, {i1}   - One-based index (1, 2, 3, ...)
  {timestamp}, {ts} - Unix timestamp
  {datetime}, {dt} - Formatted datetime (YYYY-MM-DD_HH-MM-SS)
  {date}          - Date only (YYYY-MM-DD)
  {time}          - Time only (HH-MM-SS)
  {seed}          - Seed value used for generation
  {model}         - Model name
  {width}         - Image width
  {height}        - Image height
  {prompt}        - First 50 chars of prompt (sanitized)
  {original}      - Original filename from server
  {ext}           - Original file extension`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate that both --length and --height are not specified simultaneously
		// They are aliases for the same parameter but both being set creates ambiguity
		if cmd.Flags().Changed("length") && cmd.Flags().Changed("height") {
			return fmt.Errorf("cannot specify both --length and --height flags (they are aliases for the same parameter)")
		}
		return nil
	},
	RunE: runGenerateImage,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateImageCmd)

	// Image generation flags
	generateImageCmd.Flags().StringVarP(&generatePrompt, "prompt", "p", "", "generation prompt (required)")
	generateImageCmd.Flags().StringVar(&generateModel, "model", "", "model to use for generation")
	generateImageCmd.Flags().IntVar(&generateSteps, "steps", 20, "number of inference steps")
	generateImageCmd.Flags().IntVarP(&generateWidth, "width", "w", 512, "image width")
	generateImageCmd.Flags().IntVarP(&generateHeight, "length", "l", 512, "image length (height)")
	// Add --height as an alias for backward compatibility
	generateImageCmd.Flags().IntVar(&generateHeight, "height", 512, "image height (alias for --length)")
	generateImageCmd.Flags().Int64Var(&generateSeed, "seed", -1, "random seed (-1 for random)")
	generateImageCmd.Flags().IntVarP(&generateBatchSize, "batch", "b", 1, "number of images to generate")
	generateImageCmd.Flags().Float64Var(&generateCfgScale, "cfg-scale", 7.5, "CFG scale (guidance)")
	generateImageCmd.Flags().StringVarP(&generateNegPrompt, "negative-prompt", "n", "", "negative prompt")
	generateImageCmd.Flags().StringVar(&generateSampler, "sampler", "euler_a", "sampling method")
	generateImageCmd.Flags().BoolVar(&generateUseWebSocket, "websocket", false, "use WebSocket for real-time progress (requires SwarmUI)")
	generateImageCmd.Flags().BoolVar(&generateSaveImages, "save-images", false, "download and save generated images to local disk")
	generateImageCmd.Flags().StringVar(&generateOutputDir, "output-dir", ".", "directory to save downloaded images (default: current directory)")
	generateImageCmd.Flags().StringVar(&generateFilenameTemplate, "filename-template", "", "template for custom filenames (e.g., 'image-{index}-{seed}.png')")
	// Auto-crop postprocessing flags
	generateImageCmd.Flags().BoolVar(&generateAutoCrop, "auto-crop", false, "automatically crop whitespace borders from images")
	generateImageCmd.Flags().IntVar(&generateAutoCropThreshold, "auto-crop-threshold", 250, "whitespace detection threshold (0-255, higher = more aggressive)")
	generateImageCmd.Flags().IntVar(&generateAutoCropTolerance, "auto-crop-tolerance", 10, "tolerance for near-white colors (0-255)")
	generateImageCmd.Flags().BoolVar(&generateAutoCropPreserveAspect, "auto-crop-preserve-aspect", false, "preserve original aspect ratio when auto-cropping")
	// Downscale postprocessing flags
	generateImageCmd.Flags().IntVar(&generateDownscaleWidth, "downscale-width", 0, "downscale images to this width after download (0=auto from height)")
	generateImageCmd.Flags().IntVar(&generateDownscaleHeight, "downscale-height", 0, "downscale images to this height after download (0=auto from width)")
	generateImageCmd.Flags().Float64Var(&generateDownscalePercentage, "downscale-percentage", 0, "downscale by percentage (1-100, 0=disabled, overrides width/height)")
	generateImageCmd.Flags().StringVar(&generateDownscaleFilter, "downscale-filter", "lanczos", "downscaling algorithm: lanczos (best), bilinear, nearest")
	// SkimmedCFG (Distilled CFG) flags - advanced sampling technique for improved quality/speed
	generateImageCmd.Flags().BoolVar(&generateSkimmedCFG, "skimmed-cfg", false, "enable Skimmed CFG (Distilled CFG) for improved quality and speed")
	generateImageCmd.Flags().Float64Var(&generateSkimmedCFGScale, "skimmed-cfg-scale", 3.0, "Skimmed CFG scale value (typically lower than standard CFG)")
	generateImageCmd.Flags().Float64Var(&generateSkimmedCFGStart, "skimmed-cfg-start", 0.0, "start percentage for Skimmed CFG application (0.0-1.0)")
	generateImageCmd.Flags().Float64Var(&generateSkimmedCFGEnd, "skimmed-cfg-end", 1.0, "end percentage for Skimmed CFG application (0.0-1.0)")

	generateImageCmd.MarkFlagRequired("prompt")

	// Bind to viper
	viper.BindPFlag("generate.model", generateImageCmd.Flags().Lookup("model"))
	viper.BindPFlag("generate.steps", generateImageCmd.Flags().Lookup("steps"))
	viper.BindPFlag("generate.width", generateImageCmd.Flags().Lookup("width"))
	viper.BindPFlag("generate.length", generateImageCmd.Flags().Lookup("length"))
	viper.BindPFlag("generate.height", generateImageCmd.Flags().Lookup("height")) // Backward compatibility alias
	viper.BindPFlag("generate.cfg-scale", generateImageCmd.Flags().Lookup("cfg-scale"))
	viper.BindPFlag("generate.sampler", generateImageCmd.Flags().Lookup("sampler"))
	viper.BindPFlag("generate.skimmed-cfg", generateImageCmd.Flags().Lookup("skimmed-cfg"))
	viper.BindPFlag("generate.skimmed-cfg-scale", generateImageCmd.Flags().Lookup("skimmed-cfg-scale"))
	viper.BindPFlag("generate.skimmed-cfg-start", generateImageCmd.Flags().Lookup("skimmed-cfg-start"))
	viper.BindPFlag("generate.skimmed-cfg-end", generateImageCmd.Flags().Lookup("skimmed-cfg-end"))
}

func runGenerateImage(cmd *cobra.Command, args []string) error {
	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handler
	setupSignalHandler(cancel)

	// Validate that config file doesn't have both length and height set
	// Only check if neither flag was explicitly set on command line
	if !cmd.Flags().Changed("length") && !cmd.Flags().Changed("height") {
		if viper.IsSet("generate.length") && viper.IsSet("generate.height") {
			lengthVal := viper.GetInt("generate.length")
			heightVal := viper.GetInt("generate.height")
			if lengthVal != heightVal {
				return fmt.Errorf("config file has conflicting values for 'length' (%d) and 'height' (%d) - they are aliases, please use only one", lengthVal, heightVal)
			}
		}
	}

	// Validate prompt
	if generatePrompt == "" {
		return fmt.Errorf("prompt is required")
	}

	// Build generation request
	req := &client.GenerationRequest{
		Prompt: generatePrompt,
		Parameters: map[string]interface{}{
			"steps":    generateSteps,
			"width":    generateWidth,
			"height":   generateHeight,
			"cfgscale": generateCfgScale, // SwarmUI API parameter name
			"sampler":  generateSampler,
			"images":   generateBatchSize, // SwarmUI uses "images" for batch size
		},
	}

	// Only include negative prompt if non-empty to avoid unnecessary API payload
	if generateNegPrompt != "" {
		req.Parameters["negative_prompt"] = generateNegPrompt
	}

	// Add SkimmedCFG parameters if enabled
	if generateSkimmedCFG {
		req.Parameters["skimmedcfg"] = true
		req.Parameters["skimmedcfgscale"] = generateSkimmedCFGScale
		// Only include start/end if they differ from defaults
		if generateSkimmedCFGStart != 0.0 {
			req.Parameters["skimmedcfgstart"] = generateSkimmedCFGStart
		}
		if generateSkimmedCFGEnd != 1.0 {
			req.Parameters["skimmedcfgend"] = generateSkimmedCFGEnd
		}
	}

	// Set model if specified
	if generateModel != "" {
		req.Model = generateModel
	} else if viper.IsSet("generate.model") {
		req.Model = viper.GetString("generate.model")
	}

	// Validate model if specified
	if req.Model != "" {
		if err := validateModel(assetClient, req.Model); err != nil {
			return fmt.Errorf("model validation failed: %w", err)
		}
	}

	// Set seed if specified
	if generateSeed >= 0 {
		req.Parameters["seed"] = generateSeed
	}

	if !quiet {
		// Provide clear feedback about batch generation
		if generateBatchSize > 1 {
			fmt.Fprintf(os.Stderr, "Generating %d images with prompt: %s\n", generateBatchSize, generatePrompt)
		} else {
			fmt.Fprintf(os.Stderr, "Generating image with prompt: %s\n", generatePrompt)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Model: %s\n", req.Model)
			fmt.Fprintf(os.Stderr, "Steps: %d, Size: %dx%d, CFG: %.1f\n",
				generateSteps, generateWidth, generateHeight, generateCfgScale)
			if generateBatchSize > 1 {
				fmt.Fprintf(os.Stderr, "Batch size: %d\n", generateBatchSize)
			}
		}
	}

	// Execute generation with progress tracking
	// Use WebSocket if flag is enabled, otherwise use HTTP
	var result *client.GenerationResult
	var err error
	if generateUseWebSocket {
		if verbose {
			fmt.Fprintf(os.Stderr, "Using WebSocket for real-time progress updates\n")
		}
		result, err = assetClient.GenerateImageWS(ctx, req)
	} else {
		result, err = assetClient.GenerateImage(ctx, req)
	}

	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Download images if requested
	if generateSaveImages {
		if !quiet {
			fmt.Fprintf(os.Stderr, "Downloading generated images...\n")
		}

		// Prepare metadata for filename template
		templateMetadata := map[string]interface{}{
			"prompt": generatePrompt,
			"model":  req.Model,
			"width":  generateWidth,
			"height": generateHeight,
		}
		if generateSeed >= 0 {
			templateMetadata["seed"] = generateSeed
		}

		// Build download options with postprocessing
		opts := &client.DownloadOptions{
			OutputDir:        generateOutputDir,
			FilenameTemplate: generateFilenameTemplate,
			Metadata:         templateMetadata,
			// Auto-crop options
			AutoCrop:               generateAutoCrop,
			AutoCropThreshold:      uint8(generateAutoCropThreshold),
			AutoCropTolerance:      uint8(generateAutoCropTolerance),
			AutoCropPreserveAspect: generateAutoCropPreserveAspect,
			// Downscale options
			DownscaleWidth:      generateDownscaleWidth,
			DownscaleHeight:     generateDownscaleHeight,
			DownscalePercentage: generateDownscalePercentage,
			DownscaleFilter:     generateDownscaleFilter,
		}

		// Download images with options
		var savedPaths []string
		savedPaths, err = assetClient.DownloadImagesWithOptions(ctx, result.ImagePaths, opts)

		if err != nil {
			return fmt.Errorf("failed to download images: %w", err)
		}

		if !quiet {
			for i, path := range savedPaths {
				fmt.Fprintf(os.Stderr, "  [%d/%d] Saved: %s\n", i+1, len(savedPaths), path)
			}
		}

		// Update result with local paths for output formatting
		// Initialize metadata map if nil
		if result.Metadata == nil {
			result.Metadata = make(map[string]interface{})
		}
		result.Metadata["local_paths"] = savedPaths
	}

	// Format and output result
	formatter := output.NewFormatter(viper.GetString("format"))
	outputData, err := formatter.Format(result)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Write output
	outputFile := viper.GetString("output")
	if outputFile != "" {
		if err := output.WriteToFile(outputFile, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if !quiet {
			fmt.Fprintf(os.Stderr, "Output saved to: %s\n", outputFile)
		}
	} else {
		fmt.Println(outputData)
	}

	if !quiet {
		// Provide clear feedback about number of images generated
		imageCount := len(result.ImagePaths)
		if imageCount == 1 {
			fmt.Fprintf(os.Stderr, "✓ Generation completed successfully (1 image)\n")
		} else {
			fmt.Fprintf(os.Stderr, "✓ Generation completed successfully (%d images)\n", imageCount)
		}
	}

	return nil
}

// validateModel checks if the specified model exists in the available models list
func validateModel(assetClient *client.AssetClient, modelName string) error {
	models, err := assetClient.ListModels()
	if err != nil {
		// If we can't list models, allow the request to proceed
		// The SwarmUI API will handle the validation
		if verbose {
			fmt.Fprintf(os.Stderr, "Warning: Could not validate model (API unavailable): %v\n", err)
		}
		return nil
	}

	// Check if the model exists in the available models
	for _, model := range models {
		if model.Name == modelName {
			return nil // Model found
		}
	}

	// Model not found - provide helpful error with suggestions
	if len(models) > 0 {
		// Find most similar model names using fuzzy matching
		type modelScore struct {
			name  string
			score int
		}

		var scored []modelScore
		for _, model := range models {
			score := stringSimilarity(strings.ToLower(modelName), strings.ToLower(model.Name))
			scored = append(scored, modelScore{name: model.Name, score: score})
		}

		// Sort by similarity score (higher is more similar)
		sort.Slice(scored, func(i, j int) bool {
			return scored[i].score > scored[j].score
		})

		// Take top 5 most similar models
		var suggestions []string
		for i := 0; i < len(scored) && i < 5; i++ {
			suggestions = append(suggestions, scored[i].name)
		}

		return fmt.Errorf("model '%s' not found\n\nDid you mean one of these?\n  %s\n\nUse 'asset-generator models list' to see all available models",
			modelName, strings.Join(suggestions, "\n  "))
	}

	return fmt.Errorf("model '%s' not found (no models available from API)", modelName)
}

// stringSimilarity calculates a simple similarity score between two strings.
// Uses a combination of substring matching and common prefix length.
// Higher scores indicate more similar strings.
func stringSimilarity(s1, s2 string) int {
	score := 0

	// Exact match gets highest score
	if s1 == s2 {
		return 1000
	}

	// Check if one is a substring of the other
	if strings.Contains(s2, s1) {
		score += 500
	} else if strings.Contains(s1, s2) {
		score += 400
	}

	// Common prefix length
	prefixLen := 0
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}
	for i := 0; i < minLen; i++ {
		if s1[i] == s2[i] {
			prefixLen++
		} else {
			break
		}
	}
	score += prefixLen * 10

	// Count common characters (case-insensitive)
	commonChars := 0
	s1Chars := make(map[rune]int)
	for _, ch := range s1 {
		s1Chars[ch]++
	}
	for _, ch := range s2 {
		if s1Chars[ch] > 0 {
			commonChars++
			s1Chars[ch]--
		}
	}
	score += commonChars

	// Penalize length difference
	lenDiff := len(s1) - len(s2)
	if lenDiff < 0 {
		lenDiff = -lenDiff
	}
	score -= lenDiff

	return score
}

func setupSignalHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintln(os.Stderr, "\nReceived interrupt signal, cancelling...")
		cancel()
	}()
}
