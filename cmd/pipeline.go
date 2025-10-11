package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opd-ai/asset-generator/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	pipelineFile          string
	pipelineOutputDir     string
	pipelineBaseSeed      int64
	pipelineModel         string
	pipelineSteps         int
	pipelineWidth         int
	pipelineHeight        int
	pipelineCfgScale      float64
	pipelineSampler       string
	pipelineScheduler     string
	pipelineStyleSuffix   string
	pipelineNegPrompt     string
	pipelineDryRun        bool
	pipelineContinueError bool
	// Postprocessing options
	pipelineAutoCrop               bool
	pipelineAutoCropThreshold      int
	pipelineAutoCropTolerance      int
	pipelineAutoCropPreserveAspect bool
	pipelineDownscaleWidth         int
	pipelineDownscaleHeight        int
	pipelineDownscalePercentage    float64
	pipelineDownscaleFilter        string
	// SkimmedCFG (Distilled CFG) options
	pipelineSkimmedCFG      bool
	pipelineSkimmedCFGScale float64
	pipelineSkimmedCFGStart float64
	pipelineSkimmedCFGEnd   float64
)

// PipelineSpec represents the structure of a generic pipeline YAML file
type PipelineSpec struct {
	Assets []AssetGroup `yaml:"assets"`
}

// AssetGroup represents a collection of related assets
type AssetGroup struct {
	Name       string                 `yaml:"name"`                // Group name (e.g., "characters", "backgrounds")
	OutputDir  string                 `yaml:"output_dir"`          // Subdirectory for this group
	SeedOffset int64                  `yaml:"seed_offset"`         // Offset to add to base seed
	Metadata   map[string]interface{} `yaml:"metadata,omitempty"`  // Group metadata (appended to prompts)
	Assets     []Asset                `yaml:"assets"`              // Individual assets in this group
	Subgroups  []AssetGroup           `yaml:"subgroups,omitempty"` // Nested groups
}

// Asset represents a single asset to generate
type Asset struct {
	ID       string                 `yaml:"id"`                 // Unique identifier for the asset
	Name     string                 `yaml:"name"`               // Display name
	Prompt   string                 `yaml:"prompt"`             // Generation prompt
	Filename string                 `yaml:"filename,omitempty"` // Custom filename (optional, defaults to sanitized ID)
	Metadata map[string]interface{} `yaml:"metadata,omitempty"` // Asset metadata (appended to prompt)
}

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Process asset generation pipeline files",
	Long: `Process YAML pipeline files for batch asset generation.

The pipeline command allows you to define complex asset generation workflows
in YAML files and process them automatically without external scripts.

Supports structured pipelines for any multi-asset project: game sprites,
character sheets, card decks, UI elements, and more.

Examples:
  # Process a generic asset pipeline
  asset-generator pipeline --file assets-spec.yaml --output-dir ./assets
  
  # Process a legacy tarot deck pipeline (backward compatible)
  asset-generator pipeline --file tarot-spec.yaml --output-dir ./deck
  
  # Preview what would be generated (dry run)
  asset-generator pipeline --file assets-spec.yaml --dry-run
  
  # Use custom generation parameters
  asset-generator pipeline --file assets-spec.yaml \
    --base-seed 42 --steps 40 --width 768 --height 1344
  
  # Add style suffix to all prompts
  asset-generator pipeline --file assets-spec.yaml \
    --style-suffix "detailed illustration, ornate border, rich colors"
  
  # Continue on error (don't stop if one asset fails)
  asset-generator pipeline --file assets-spec.yaml --continue-on-error
  
  # With postprocessing
  asset-generator pipeline --file assets-spec.yaml \
    --auto-crop --downscale-width 1024

Pipeline File Structure (Generic Format):
  assets:
    - name: Characters
      output_dir: characters
      seed_offset: 0
      assets:
        - id: hero_01
          name: Hero Character
          prompt: "heroic warrior, detailed armor..."
          filename: hero.png
        - id: villain_01
          name: Villain Character
          prompt: "dark sorcerer, mysterious robes..."
          
    - name: Backgrounds
      output_dir: backgrounds
      seed_offset: 100
      assets:
        - id: forest
          name: Forest Scene
          prompt: "mystical forest, sunbeams..."

Legacy Tarot Format (Backward Compatible):
  major_arcana:
    - number: 0
      name: The Fool
      prompt: "card description..."
  
  minor_arcana:
    wands:
      suit_element: fire
      suit_color: red
      cards:
        - rank: Ace
          prompt: "card description..."

Output Structure:
  Generated assets will be organized according to the structure
  defined in your pipeline file.`,
	RunE: runPipeline,
}

func init() {
	rootCmd.AddCommand(pipelineCmd)

	// Required flags
	pipelineCmd.Flags().StringVar(&pipelineFile, "file", "", "pipeline YAML file (required)")
	pipelineCmd.Flags().StringVar(&pipelineOutputDir, "output-dir", "./pipeline-output", "output directory for generated assets")

	// Generation parameters
	pipelineCmd.Flags().Int64Var(&pipelineBaseSeed, "base-seed", -1, "base seed for reproducible generation (0 or -1 for random)")
	pipelineCmd.Flags().StringVar(&pipelineModel, "model", "", "model to use for all generations")
	pipelineCmd.Flags().IntVar(&pipelineSteps, "steps", 40, "number of inference steps")
	pipelineCmd.Flags().IntVar(&pipelineWidth, "width", 768, "image width")
	pipelineCmd.Flags().IntVar(&pipelineHeight, "height", 1344, "image height")
	pipelineCmd.Flags().Float64Var(&pipelineCfgScale, "cfg-scale", 7.5, "CFG scale (guidance)")
	pipelineCmd.Flags().StringVar(&pipelineSampler, "sampler", "euler_a", "sampling method")
	pipelineCmd.Flags().StringVar(&pipelineScheduler, "scheduler", "simple", "scheduler/noise schedule (simple, normal, karras, exponential, sgm_uniform)")

	// Prompt enhancement
	pipelineCmd.Flags().StringVar(&pipelineStyleSuffix, "style-suffix", "", "suffix to append to all prompts")
	pipelineCmd.Flags().StringVar(&pipelineNegPrompt, "negative-prompt", "", "negative prompt for all generations")

	// Pipeline control
	pipelineCmd.Flags().BoolVar(&pipelineDryRun, "dry-run", false, "preview pipeline without generating")
	pipelineCmd.Flags().BoolVar(&pipelineContinueError, "continue-on-error", false, "continue processing if individual generations fail")

	// Postprocessing options
	pipelineCmd.Flags().BoolVar(&pipelineAutoCrop, "auto-crop", false, "automatically crop whitespace borders")
	pipelineCmd.Flags().IntVar(&pipelineAutoCropThreshold, "auto-crop-threshold", 250, "whitespace detection threshold (0-255)")
	pipelineCmd.Flags().IntVar(&pipelineAutoCropTolerance, "auto-crop-tolerance", 10, "tolerance for near-white colors")
	pipelineCmd.Flags().BoolVar(&pipelineAutoCropPreserveAspect, "auto-crop-preserve-aspect", false, "preserve aspect ratio when auto-cropping")
	pipelineCmd.Flags().IntVar(&pipelineDownscaleWidth, "downscale-width", 0, "downscale to this width (0=disabled)")
	pipelineCmd.Flags().IntVar(&pipelineDownscaleHeight, "downscale-height", 0, "downscale to this height (0=disabled)")
	pipelineCmd.Flags().Float64Var(&pipelineDownscalePercentage, "downscale-percentage", 0, "downscale by percentage (0=disabled)")
	pipelineCmd.Flags().StringVar(&pipelineDownscaleFilter, "downscale-filter", "lanczos", "downscaling filter (lanczos, bilinear, nearest)")
	// SkimmedCFG (Distilled CFG) options
	pipelineCmd.Flags().BoolVar(&pipelineSkimmedCFG, "skimmed-cfg", false, "enable Skimmed CFG for improved quality and speed")
	pipelineCmd.Flags().Float64Var(&pipelineSkimmedCFGScale, "skimmed-cfg-scale", 3.0, "Skimmed CFG scale value")
	pipelineCmd.Flags().Float64Var(&pipelineSkimmedCFGStart, "skimmed-cfg-start", 0.0, "start percentage for Skimmed CFG (0.0-1.0)")
	pipelineCmd.Flags().Float64Var(&pipelineSkimmedCFGEnd, "skimmed-cfg-end", 1.0, "end percentage for Skimmed CFG (0.0-1.0)")

	pipelineCmd.MarkFlagRequired("file")
}

func runPipeline(cmd *cobra.Command, args []string) error {
	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handler
	setupSignalHandler(cancel)

	// Read and parse pipeline file
	if !quiet {
		fmt.Fprintf(os.Stderr, "Loading pipeline file: %s\n", pipelineFile)
	}

	spec, err := loadPipelineSpec(pipelineFile)
	if err != nil {
		return fmt.Errorf("failed to load pipeline: %w", err)
	}

	// Calculate total work
	totalAssets := countAssets(spec.Assets)

	if !quiet {
		fmt.Fprintf(os.Stderr, "Pipeline loaded: %d total assets\n", totalAssets)
		printGroupSummary(spec.Assets, "")
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Generate random seed if not specified (both -1 and 0 trigger random seed)
	if pipelineBaseSeed == -1 || pipelineBaseSeed == 0 {
		pipelineBaseSeed = time.Now().UnixNano()
		if !quiet {
			fmt.Fprintf(os.Stderr, "Generated random base seed: %d\n\n", pipelineBaseSeed)
		}
	}

	if pipelineDryRun {
		fmt.Fprintf(os.Stderr, "DRY RUN - No assets will be generated\n\n")
		return previewPipeline(spec)
	}

	// Create output directory
	if err := os.MkdirAll(pipelineOutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "Output directory: %s\n", pipelineOutputDir)
		fmt.Fprintf(os.Stderr, "Base seed: %d\n", pipelineBaseSeed)
		fmt.Fprintf(os.Stderr, "Dimensions: %dx%d\n", pipelineWidth, pipelineHeight)
		fmt.Fprintf(os.Stderr, "Steps: %d, CFG Scale: %.1f\n\n", pipelineSteps, pipelineCfgScale)
	}

	// Track progress
	completed := 0
	failed := 0

	// Process all groups
	for _, group := range spec.Assets {
		c, f, err := processGroup(ctx, group, pipelineOutputDir, &completed, totalAssets, nil)
		completed += c
		failed += f

		if err != nil && !pipelineContinueError {
			return err
		}
	}

	// Summary
	if !quiet {
		fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Fprintf(os.Stderr, "Pipeline Complete!\n")
		fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		fmt.Fprintf(os.Stderr, "Total assets generated: %d/%d\n", completed, totalAssets)
		if failed > 0 {
			fmt.Fprintf(os.Stderr, "Failed: %d\n", failed)
		}
		fmt.Fprintf(os.Stderr, "Output location: %s\n", pipelineOutputDir)
	}

	if failed > 0 && !pipelineContinueError {
		return fmt.Errorf("pipeline completed with %d failures", failed)
	}

	return nil
}

// countAssets recursively counts all assets in groups and subgroups
func countAssets(groups []AssetGroup) int {
	count := 0
	for _, group := range groups {
		count += len(group.Assets)
		count += countAssets(group.Subgroups)
	}
	return count
}

// printGroupSummary prints a summary of asset groups
func printGroupSummary(groups []AssetGroup, indent string) {
	for _, group := range groups {
		if len(group.Assets) > 0 {
			fmt.Fprintf(os.Stderr, "%s  - %s: %d assets\n", indent, group.Name, len(group.Assets))
		}
		if len(group.Subgroups) > 0 {
			printGroupSummary(group.Subgroups, indent+"  ")
		}
	}
}

// processGroup processes a single asset group and its subgroups
func processGroup(ctx context.Context, group AssetGroup, baseOutputDir string, completed *int, totalAssets int, parentMetadata map[string]interface{}) (int, int, error) {
	groupCompleted := 0
	groupFailed := 0

	// Merge parent metadata with group metadata
	groupMetadata := mergeMetadata(parentMetadata, group.Metadata)

	// Create group output directory
	groupOutputDir := filepath.Join(baseOutputDir, group.OutputDir)
	if err := os.MkdirAll(groupOutputDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("failed to create group directory %s: %w", groupOutputDir, err)
	}

	if len(group.Assets) > 0 && !quiet {
		fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Fprintf(os.Stderr, "Processing: %s (%d assets)\n", group.Name, len(group.Assets))
		fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	}

	// Process assets in this group
	for i, asset := range group.Assets {
		if err := ctx.Err(); err != nil {
			return groupCompleted, groupFailed, fmt.Errorf("pipeline cancelled: %w", err)
		}

		// Merge group metadata with asset metadata
		assetMetadata := mergeMetadata(groupMetadata, asset.Metadata)

		// Calculate seed
		seed := pipelineBaseSeed + group.SeedOffset + int64(i)

		// Determine filename
		filename := asset.Filename
		if filename == "" {
			filename = sanitizeFilename(asset.ID) + ".png"
		}
		outputPath := filepath.Join(groupOutputDir, filename)

		if !quiet {
			fmt.Fprintf(os.Stderr, "[%d/%d] Generating: %s\n", *completed+1, totalAssets, asset.Name)
		}

		// Build enhanced prompt with metadata
		enhancedPrompt := buildEnhancedPrompt(asset.Prompt, assetMetadata)

		if err := generateAsset(ctx, enhancedPrompt, asset.Name, outputPath, seed, assetMetadata); err != nil {
			groupFailed++
			if pipelineContinueError {
				fmt.Fprintf(os.Stderr, "  ⚠ Warning: Failed to generate %s: %v\n", asset.Name, err)
			} else {
				return groupCompleted, groupFailed, fmt.Errorf("failed to generate %s: %w", asset.Name, err)
			}
		} else {
			groupCompleted++
			*completed++
			if !quiet {
				fmt.Fprintf(os.Stderr, "  ✓ Saved to: %s\n", outputPath)
			}
		}
		fmt.Fprintln(os.Stderr)
	}

	// Process subgroups
	for _, subgroup := range group.Subgroups {
		c, f, err := processGroup(ctx, subgroup, groupOutputDir, completed, totalAssets, groupMetadata)
		groupCompleted += c
		groupFailed += f
		if err != nil && !pipelineContinueError {
			return groupCompleted, groupFailed, err
		}
	}

	return groupCompleted, groupFailed, nil
}

// mergeMetadata merges parent metadata with child metadata (child takes precedence)
func mergeMetadata(parent, child map[string]interface{}) map[string]interface{} {
	if parent == nil && child == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Copy parent metadata
	for k, v := range parent {
		result[k] = v
	}

	// Override with child metadata
	for k, v := range child {
		result[k] = v
	}

	return result
}

// buildEnhancedPrompt builds a prompt with metadata appended
func buildEnhancedPrompt(basePrompt string, metadata map[string]interface{}) string {
	if len(metadata) == 0 {
		return basePrompt
	}

	// Collect metadata values as strings
	var metadataParts []string
	for _, v := range metadata {
		if str, ok := v.(string); ok && str != "" {
			metadataParts = append(metadataParts, str)
		}
	}

	if len(metadataParts) == 0 {
		return basePrompt
	}

	// Append metadata to prompt
	return fmt.Sprintf("%s, %s", basePrompt, strings.Join(metadataParts, ", "))
}

func loadPipelineSpec(filename string) (*PipelineSpec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var spec PipelineSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &spec, nil
}

func generateAsset(ctx context.Context, prompt, name, outputPath string, seed int64, metadata map[string]interface{}) error {
	// Build full prompt with style suffix
	fullPrompt := prompt
	if pipelineStyleSuffix != "" {
		fullPrompt = fmt.Sprintf("%s, %s", prompt, pipelineStyleSuffix)
	}

	// Build generation request
	req := &client.GenerationRequest{
		Prompt: fullPrompt,
		Parameters: map[string]interface{}{
			"steps":     pipelineSteps,
			"width":     pipelineWidth,
			"height":    pipelineHeight,
			"cfgscale":  pipelineCfgScale,
			"sampler":   pipelineSampler,
			"scheduler": pipelineScheduler,
			"seed":      seed,
			"images":    1, // Always generate one at a time for pipelines
		},
	}

	if pipelineNegPrompt != "" {
		req.Parameters["negative_prompt"] = pipelineNegPrompt
	}

	// Add SkimmedCFG parameters if enabled
	if pipelineSkimmedCFG {
		req.Parameters["skimmedcfg"] = true
		req.Parameters["skimmedcfgscale"] = pipelineSkimmedCFGScale
		if pipelineSkimmedCFGStart != 0.0 {
			req.Parameters["skimmedcfgstart"] = pipelineSkimmedCFGStart
		}
		if pipelineSkimmedCFGEnd != 1.0 {
			req.Parameters["skimmedcfgend"] = pipelineSkimmedCFGEnd
		}
	}

	if pipelineModel != "" {
		req.Model = pipelineModel
	}

	// Generate image
	result, err := assetClient.GenerateImage(ctx, req)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	if len(result.ImagePaths) == 0 {
		return fmt.Errorf("no images generated")
	}

	// Merge metadata for download
	downloadMetadata := map[string]interface{}{
		"prompt": prompt,
		"name":   name,
		"seed":   seed,
	}
	for k, v := range metadata {
		downloadMetadata[k] = v
	}

	// Download with postprocessing options
	opts := &client.DownloadOptions{
		OutputDir:        filepath.Dir(outputPath),
		FilenameTemplate: filepath.Base(outputPath),
		Metadata:         downloadMetadata,
		// Auto-crop options
		AutoCrop:               pipelineAutoCrop,
		AutoCropThreshold:      uint8(pipelineAutoCropThreshold),
		AutoCropTolerance:      uint8(pipelineAutoCropTolerance),
		AutoCropPreserveAspect: pipelineAutoCropPreserveAspect,
		// Downscale options
		DownscaleWidth:      pipelineDownscaleWidth,
		DownscaleHeight:     pipelineDownscaleHeight,
		DownscalePercentage: pipelineDownscalePercentage,
		DownscaleFilter:     pipelineDownscaleFilter,
	}

	_, err = assetClient.DownloadImagesWithOptions(ctx, result.ImagePaths, opts)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	return nil
}

func previewPipeline(spec *PipelineSpec) error {
	fmt.Println("Pipeline Preview:")
	fmt.Println()

	previewGroups(spec.Assets, "", nil)

	fmt.Println()
	fmt.Println("Generation Parameters:")
	fmt.Printf("  Dimensions: %dx%d\n", pipelineWidth, pipelineHeight)
	fmt.Printf("  Steps: %d\n", pipelineSteps)
	fmt.Printf("  CFG Scale: %.1f\n", pipelineCfgScale)
	fmt.Printf("  Sampler: %s\n", pipelineSampler)
	fmt.Printf("  Scheduler: %s\n", pipelineScheduler)
	if pipelineModel != "" {
		fmt.Printf("  Model: %s\n", pipelineModel)
	}
	if pipelineStyleSuffix != "" {
		fmt.Printf("  Style Suffix: %s\n", pipelineStyleSuffix)
	}
	if pipelineNegPrompt != "" {
		fmt.Printf("  Negative Prompt: %s\n", pipelineNegPrompt)
	}

	return nil
}

func previewGroups(groups []AssetGroup, indent string, parentMetadata map[string]interface{}) {
	for _, group := range groups {
		fmt.Printf("%s%s (%s):\n", indent, group.Name, group.OutputDir)

		// Merge metadata
		groupMetadata := mergeMetadata(parentMetadata, group.Metadata)
		if len(groupMetadata) > 0 {
			fmt.Printf("%s  Metadata: %v\n", indent, groupMetadata)
		}

		// Preview assets
		for i, asset := range group.Assets {
			seed := pipelineBaseSeed + group.SeedOffset + int64(i)
			assetMetadata := mergeMetadata(groupMetadata, asset.Metadata)
			enhancedPrompt := buildEnhancedPrompt(asset.Prompt, assetMetadata)

			fmt.Printf("%s  [%s] %s (seed: %d)\n", indent, asset.ID, asset.Name, seed)
			if verbose {
				fmt.Printf("%s    Prompt: %s\n", indent, enhancedPrompt)
				if asset.Filename != "" {
					fmt.Printf("%s    Filename: %s\n", indent, asset.Filename)
				}
			}
		}

		// Preview subgroups
		if len(group.Subgroups) > 0 {
			fmt.Println()
			previewGroups(group.Subgroups, indent+"  ", groupMetadata)
		}

		fmt.Println()
	}
}

func sanitizeFilename(name string) string {
	// Convert to lowercase and replace spaces with underscores
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	// Remove special characters
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}

	return result.String()
}
