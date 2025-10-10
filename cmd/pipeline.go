package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
)

// PipelineSpec represents the structure of a pipeline YAML file
type PipelineSpec struct {
	MajorArcana []MajorArcanaCard `yaml:"major_arcana"`
	MinorArcana MinorArcanaSpec   `yaml:"minor_arcana"`
}

// MajorArcanaCard represents a single Major Arcana card
type MajorArcanaCard struct {
	Number int    `yaml:"number"`
	Name   string `yaml:"name"`
	Prompt string `yaml:"prompt"`
}

// MinorArcanaSpec represents the Minor Arcana structure
type MinorArcanaSpec struct {
	Wands     SuitSpec `yaml:"wands"`
	Cups      SuitSpec `yaml:"cups"`
	Swords    SuitSpec `yaml:"swords"`
	Pentacles SuitSpec `yaml:"pentacles"`
}

// SuitSpec represents a suit with metadata and cards
type SuitSpec struct {
	SuitElement string     `yaml:"suit_element"`
	SuitColor   string     `yaml:"suit_color"`
	Cards       []RankCard `yaml:"cards"`
}

// RankCard represents a single card in a suit
type RankCard struct {
	Rank   string `yaml:"rank"`
	Prompt string `yaml:"prompt"`
}

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Process asset generation pipeline files",
	Long: `Process YAML pipeline files for batch asset generation.

The pipeline command allows you to define complex asset generation workflows
in YAML files and process them automatically without external scripts.

Supports structured pipelines like tarot decks, game sprite collections,
character sheets, and other multi-asset projects.

Examples:
  # Process a tarot deck pipeline
  asset-generator pipeline --file tarot-spec.yaml --output-dir ./deck
  
  # Preview what would be generated (dry run)
  asset-generator pipeline --file tarot-spec.yaml --dry-run
  
  # Use custom generation parameters
  asset-generator pipeline --file tarot-spec.yaml \
    --base-seed 42 --steps 40 --width 768 --height 1344
  
  # Add style suffix to all prompts
  asset-generator pipeline --file tarot-spec.yaml \
    --style-suffix "detailed illustration, ornate border, rich colors"
  
  # Continue on error (don't stop if one card fails)
  asset-generator pipeline --file tarot-spec.yaml --continue-on-error
  
  # With postprocessing
  asset-generator pipeline --file tarot-spec.yaml \
    --auto-crop --downscale-width 1024

Pipeline File Structure:
  The pipeline file should follow this structure:
  
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
  Generated assets will be organized in subdirectories:
  
  output-dir/
    ├── major-arcana/
    │   ├── 00-the_fool.png
    │   ├── 01-the_magician.png
    │   └── ...
    └── minor-arcana/
        ├── wands/
        ├── cups/
        ├── swords/
        └── pentacles/`,
	RunE: runPipeline,
}

func init() {
	rootCmd.AddCommand(pipelineCmd)

	// Required flags
	pipelineCmd.Flags().StringVar(&pipelineFile, "file", "", "pipeline YAML file (required)")
	pipelineCmd.Flags().StringVar(&pipelineOutputDir, "output-dir", "./pipeline-output", "output directory for generated assets")

	// Generation parameters
	pipelineCmd.Flags().Int64Var(&pipelineBaseSeed, "base-seed", 42, "base seed for reproducible generation")
	pipelineCmd.Flags().StringVar(&pipelineModel, "model", "", "model to use for all generations")
	pipelineCmd.Flags().IntVar(&pipelineSteps, "steps", 40, "number of inference steps")
	pipelineCmd.Flags().IntVar(&pipelineWidth, "width", 768, "image width")
	pipelineCmd.Flags().IntVar(&pipelineHeight, "height", 1344, "image height")
	pipelineCmd.Flags().Float64Var(&pipelineCfgScale, "cfg-scale", 7.5, "CFG scale (guidance)")
	pipelineCmd.Flags().StringVar(&pipelineSampler, "sampler", "euler_a", "sampling method")

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
	totalCards := len(spec.MajorArcana)
	totalCards += len(spec.MinorArcana.Wands.Cards)
	totalCards += len(spec.MinorArcana.Cups.Cards)
	totalCards += len(spec.MinorArcana.Swords.Cards)
	totalCards += len(spec.MinorArcana.Pentacles.Cards)

	if !quiet {
		fmt.Fprintf(os.Stderr, "Pipeline loaded: %d total assets\n", totalCards)
		fmt.Fprintf(os.Stderr, "  - Major Arcana: %d cards\n", len(spec.MajorArcana))
		fmt.Fprintf(os.Stderr, "  - Minor Arcana (Wands): %d cards\n", len(spec.MinorArcana.Wands.Cards))
		fmt.Fprintf(os.Stderr, "  - Minor Arcana (Cups): %d cards\n", len(spec.MinorArcana.Cups.Cards))
		fmt.Fprintf(os.Stderr, "  - Minor Arcana (Swords): %d cards\n", len(spec.MinorArcana.Swords.Cards))
		fmt.Fprintf(os.Stderr, "  - Minor Arcana (Pentacles): %d cards\n", len(spec.MinorArcana.Pentacles.Cards))
		fmt.Fprintf(os.Stderr, "\n")
	}

	if pipelineDryRun {
		fmt.Fprintf(os.Stderr, "DRY RUN - No assets will be generated\n\n")
		return previewPipeline(spec)
	}

	// Create output directory structure
	if err := createOutputDirectories(pipelineOutputDir); err != nil {
		return fmt.Errorf("failed to create output directories: %w", err)
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

	// Process Major Arcana
	if len(spec.MajorArcana) > 0 {
		if !quiet {
			fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Fprintf(os.Stderr, "Processing Major Arcana (%d cards)\n", len(spec.MajorArcana))
			fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		}

		for _, card := range spec.MajorArcana {
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("pipeline cancelled: %w", err)
			}

			seed := pipelineBaseSeed + int64(card.Number)
			paddedNum := fmt.Sprintf("%02d", card.Number)
			filename := fmt.Sprintf("%s-%s.png", paddedNum, sanitizeFilename(card.Name))
			outputPath := filepath.Join(pipelineOutputDir, "major-arcana", filename)

			if !quiet {
				fmt.Fprintf(os.Stderr, "[%d/%d] Generating: %s - %s\n", completed+1, totalCards, paddedNum, card.Name)
			}

			if err := generateCard(ctx, card.Prompt, card.Name, outputPath, seed); err != nil {
				failed++
				if pipelineContinueError {
					fmt.Fprintf(os.Stderr, "  ⚠ Warning: Failed to generate %s: %v\n", card.Name, err)
				} else {
					return fmt.Errorf("failed to generate %s: %w", card.Name, err)
				}
			} else {
				completed++
				if !quiet {
					fmt.Fprintf(os.Stderr, "  ✓ Saved to: %s\n", outputPath)
				}
			}
			fmt.Fprintln(os.Stderr)
		}
	}

	// Process Minor Arcana
	minorSeedOffset := 100
	suits := []struct {
		name   string
		spec   SuitSpec
		offset int
	}{
		{"wands", spec.MinorArcana.Wands, 0},
		{"cups", spec.MinorArcana.Cups, 20},
		{"swords", spec.MinorArcana.Swords, 40},
		{"pentacles", spec.MinorArcana.Pentacles, 60},
	}

	for _, suit := range suits {
		if len(suit.spec.Cards) == 0 {
			continue
		}

		if !quiet {
			fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Fprintf(os.Stderr, "Processing Minor Arcana - %s (%d cards)\n", strings.Title(suit.name), len(suit.spec.Cards))
			fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		}

		for i, card := range suit.spec.Cards {
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("pipeline cancelled: %w", err)
			}

			seed := pipelineBaseSeed + int64(minorSeedOffset+suit.offset+i)
			rankNum := getRankNumber(card.Rank)
			filename := fmt.Sprintf("%02d-%s_of_%s.png", rankNum, sanitizeFilename(card.Rank), suit.name)
			outputPath := filepath.Join(pipelineOutputDir, "minor-arcana", suit.name, filename)

			if !quiet {
				fmt.Fprintf(os.Stderr, "[%d/%d] Generating: %s of %s\n", completed+1, totalCards, card.Rank, strings.Title(suit.name))
			}

			cardName := fmt.Sprintf("%s of %s", card.Rank, strings.Title(suit.name))
			if err := generateCard(ctx, card.Prompt, cardName, outputPath, seed); err != nil {
				failed++
				if pipelineContinueError {
					fmt.Fprintf(os.Stderr, "  ⚠ Warning: Failed to generate %s: %v\n", cardName, err)
				} else {
					return fmt.Errorf("failed to generate %s: %w", cardName, err)
				}
			} else {
				completed++
				if !quiet {
					fmt.Fprintf(os.Stderr, "  ✓ Saved to: %s\n", outputPath)
				}
			}
			fmt.Fprintln(os.Stderr)
		}
	}

	// Summary
	if !quiet {
		fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Fprintf(os.Stderr, "Pipeline Complete!\n")
		fmt.Fprintf(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
		fmt.Fprintf(os.Stderr, "Total cards generated: %d/%d\n", completed, totalCards)
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

func createOutputDirectories(baseDir string) error {
	dirs := []string{
		filepath.Join(baseDir, "major-arcana"),
		filepath.Join(baseDir, "minor-arcana", "wands"),
		filepath.Join(baseDir, "minor-arcana", "cups"),
		filepath.Join(baseDir, "minor-arcana", "swords"),
		filepath.Join(baseDir, "minor-arcana", "pentacles"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func generateCard(ctx context.Context, prompt, name, outputPath string, seed int64) error {
	// Build full prompt with style suffix
	fullPrompt := prompt
	if pipelineStyleSuffix != "" {
		fullPrompt = fmt.Sprintf("%s, %s", prompt, pipelineStyleSuffix)
	}

	// Build generation request
	req := &client.GenerationRequest{
		Prompt: fullPrompt,
		Parameters: map[string]interface{}{
			"steps":    pipelineSteps,
			"width":    pipelineWidth,
			"height":   pipelineHeight,
			"cfgscale": pipelineCfgScale,
			"sampler":  pipelineSampler,
			"seed":     seed,
			"images":   1, // Always generate one at a time for pipelines
		},
	}

	if pipelineNegPrompt != "" {
		req.Parameters["negative_prompt"] = pipelineNegPrompt
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

	// Download with postprocessing options
	opts := &client.DownloadOptions{
		OutputDir:        filepath.Dir(outputPath),
		FilenameTemplate: filepath.Base(outputPath),
		Metadata: map[string]interface{}{
			"prompt": prompt,
			"name":   name,
			"seed":   seed,
		},
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

	fmt.Println("Major Arcana:")
	for _, card := range spec.MajorArcana {
		seed := pipelineBaseSeed + int64(card.Number)
		fmt.Printf("  %02d - %-25s (seed: %d)\n", card.Number, card.Name, seed)
		if verbose {
			fmt.Printf("       Prompt: %s\n", card.Prompt)
		}
	}

	fmt.Println()
	fmt.Println("Minor Arcana:")

	suits := []struct {
		name   string
		spec   SuitSpec
		offset int
	}{
		{"Wands", spec.MinorArcana.Wands, 0},
		{"Cups", spec.MinorArcana.Cups, 20},
		{"Swords", spec.MinorArcana.Swords, 40},
		{"Pentacles", spec.MinorArcana.Pentacles, 60},
	}

	for _, suit := range suits {
		fmt.Printf("\n  %s (%s):\n", suit.name, suit.spec.SuitElement)
		for i, card := range suit.spec.Cards {
			seed := pipelineBaseSeed + int64(100+suit.offset+i)
			fmt.Printf("    %02d - %-20s (seed: %d)\n", getRankNumber(card.Rank), card.Rank, seed)
			if verbose {
				fmt.Printf("         Prompt: %s\n", card.Prompt)
			}
		}
	}

	fmt.Println()
	fmt.Println("Generation Parameters:")
	fmt.Printf("  Dimensions: %dx%d\n", pipelineWidth, pipelineHeight)
	fmt.Printf("  Steps: %d\n", pipelineSteps)
	fmt.Printf("  CFG Scale: %.1f\n", pipelineCfgScale)
	fmt.Printf("  Sampler: %s\n", pipelineSampler)
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

func getRankNumber(rank string) int {
	rankMap := map[string]int{
		"Ace":    1,
		"Two":    2,
		"Three":  3,
		"Four":   4,
		"Five":   5,
		"Six":    6,
		"Seven":  7,
		"Eight":  8,
		"Nine":   9,
		"Ten":    10,
		"Page":   11,
		"Knight": 12,
		"Queen":  13,
		"King":   14,
	}

	if num, ok := rankMap[rank]; ok {
		return num
	}
	return 0
}
