package cmd

import (
	"fmt"
	"os"

	"github.com/opd-ai/asset-generator/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// modelsCmd represents the models command
var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Manage AI models",
	Long: `List, inspect, and manage models available for asset generation.

Examples:
  # List all available models
  asset-generator models list
  
  # List models as JSON
  asset-generator models list --format json
  
  # Get details about a specific model
  asset-generator models get stable-diffusion-xl`,
}

// modelsListCmd lists all available models
var modelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available models",
	Long: `List all models available in the asset generation service.

Examples:
  asset-generator models list
  asset-generator models list --format json`,
	RunE: runModelsList,
}

// modelsGetCmd gets details about a specific model
var modelsGetCmd = &cobra.Command{
	Use:   "get [model-name]",
	Short: "Get details about a specific model",
	Long: `Get detailed information about a specific model.

Examples:
  asset-generator models get stable-diffusion-xl
  asset-generator models get sdxl-turbo --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runModelsGet,
}

func init() {
	rootCmd.AddCommand(modelsCmd)
	modelsCmd.AddCommand(modelsListCmd)
	modelsCmd.AddCommand(modelsGetCmd)
}

func runModelsList(cmd *cobra.Command, args []string) error {
	if !quiet {
		fmt.Fprintln(os.Stderr, "Fetching available models...")
	}

	models, err := assetClient.ListModels()
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	// Format and output result
	formatter := output.NewFormatter(viper.GetString("format"))
	outputData, err := formatter.Format(models)
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
		fmt.Fprintf(os.Stderr, "✓ Found %d models\n", len(models))
	}

	return nil
}

func runModelsGet(cmd *cobra.Command, args []string) error {
	modelName := args[0]

	if !quiet {
		fmt.Fprintf(os.Stderr, "Fetching model: %s\n", modelName)
	}

	model, err := assetClient.GetModel(modelName)
	if err != nil {
		return fmt.Errorf("failed to get model: %w", err)
	}

	// Format and output result
	formatter := output.NewFormatter(viper.GetString("format"))
	outputData, err := formatter.Format(model)
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

	return nil
}
