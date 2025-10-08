package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/opd-ai/asset-generator/pkg/client"
	"github.com/opd-ai/asset-generator/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	generatePrompt    string
	generateModel     string
	generateSteps     int
	generateWidth     int
	generateHeight    int
	generateSeed      int64
	generateBatchSize int
	generateCfgScale  float64
	generateNegPrompt string
	generateSampler   string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate assets using SwarmUI",
	Long: `Generate various types of assets using the SwarmUI API.
Supports image generation with various models and parameters.`,
}

// generateImageCmd represents the image generation command
var generateImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate an image",
	Long: `Generate an image using SwarmUI's text-to-image capabilities.

Examples:
  # Basic generation
  asset-generator generate image --prompt "a beautiful landscape"
  
  # Advanced generation with parameters
  asset-generator generate image \
    --prompt "futuristic city at sunset" \
    --model "stable-diffusion-xl" \
    --width 1024 --height 1024 \
    --steps 30 --cfg-scale 7.5
  
  # Save to specific file
  asset-generator generate image \
    --prompt "cat wearing sunglasses" \
    --output my-cat.json`,
	RunE: runGenerateImage,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateImageCmd)

	// Image generation flags
	generateImageCmd.Flags().StringVarP(&generatePrompt, "prompt", "p", "", "generation prompt (required)")
	generateImageCmd.Flags().StringVar(&generateModel, "model", "", "model to use for generation")
	generateImageCmd.Flags().IntVar(&generateSteps, "steps", 20, "number of inference steps")
	generateImageCmd.Flags().IntVar(&generateWidth, "width", 512, "image width")
	generateImageCmd.Flags().IntVar(&generateHeight, "height", 512, "image height")
	generateImageCmd.Flags().Int64Var(&generateSeed, "seed", -1, "random seed (-1 for random)")
	generateImageCmd.Flags().IntVarP(&generateBatchSize, "batch", "b", 1, "number of images to generate")
	generateImageCmd.Flags().Float64Var(&generateCfgScale, "cfg-scale", 7.5, "CFG scale (guidance)")
	generateImageCmd.Flags().StringVar(&generateNegPrompt, "negative-prompt", "", "negative prompt")
	generateImageCmd.Flags().StringVar(&generateSampler, "sampler", "euler_a", "sampling method")

	generateImageCmd.MarkFlagRequired("prompt")

	// Bind to viper
	viper.BindPFlag("generate.model", generateImageCmd.Flags().Lookup("model"))
	viper.BindPFlag("generate.steps", generateImageCmd.Flags().Lookup("steps"))
	viper.BindPFlag("generate.width", generateImageCmd.Flags().Lookup("width"))
	viper.BindPFlag("generate.height", generateImageCmd.Flags().Lookup("height"))
	viper.BindPFlag("generate.cfg-scale", generateImageCmd.Flags().Lookup("cfg-scale"))
	viper.BindPFlag("generate.sampler", generateImageCmd.Flags().Lookup("sampler"))
}

func runGenerateImage(cmd *cobra.Command, args []string) error {
	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handler
	setupSignalHandler(cancel)

	// Validate prompt
	if generatePrompt == "" {
		return fmt.Errorf("prompt is required")
	}

	// Build generation request
	req := &client.GenerationRequest{
		Prompt: generatePrompt,
		Parameters: map[string]interface{}{
			"steps":           generateSteps,
			"width":           generateWidth,
			"height":          generateHeight,
			"cfgscale":        generateCfgScale, // Match SwarmUI API parameter name
			"sampler":         generateSampler,
			"batch_size":      generateBatchSize,
			"negative_prompt": generateNegPrompt,
		},
	}

	// Set model if specified
	if generateModel != "" {
		req.Model = generateModel
	} else if viper.IsSet("generate.model") {
		req.Model = viper.GetString("generate.model")
	}

	// Validate model if specified
	if req.Model != "" {
		if err := validateModel(swarmClient, req.Model); err != nil {
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
	result, err := swarmClient.GenerateImage(ctx, req)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
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
func validateModel(swarmClient *client.SwarmClient, modelName string) error {
	models, err := swarmClient.ListModels()
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
		var suggestions []string
		for i, model := range models {
			if i < 5 { // Limit to first 5 suggestions
				suggestions = append(suggestions, model.Name)
			}
		}
		return fmt.Errorf("model '%s' not found\n\nAvailable models:\n  %s\n\nUse 'asset-generator models list' to see all available models",
			modelName, strings.Join(suggestions, "\n  "))
	}

	return fmt.Errorf("model '%s' not found (no models available from API)", modelName)
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
