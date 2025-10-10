# Pipeline Processing Guide

The `pipeline` command allows you to process YAML pipeline files for automated batch asset generation. This eliminates the need for external shell scripts and provides a native, cross-platform solution for complex generation workflows.

## Table of Contents

- [Overview](#overview)
- [Pipeline File Format](#pipeline-file-format)
- [Command Reference](#command-reference)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

The pipeline feature is designed for projects that require generating multiple related assets with:

- **Consistent styling** across all assets
- **Organized output structure** with subdirectories
- **Reproducible results** using seed-based generation
- **Progress tracking** with detailed status updates
- **Error handling** with continue-on-error support
- **Postprocessing** including auto-crop and downscaling

### Use Cases

- **Tarot decks**: 78 cards organized by Major/Minor Arcana and suits
- **Game sprites**: Character sets, enemy variants, item collections
- **Card games**: Deck of cards with multiple suits and ranks
- **Character sheets**: Multiple poses, expressions, or equipment variants
- **Icon sets**: Consistent icon families with variations

## Pipeline File Format

Pipeline files use YAML format with a specific structure:

### Basic Structure

```yaml
major_arcana:
  - number: 0
    name: The Fool
    prompt: "detailed description of the asset"
  
  - number: 1
    name: The Magician
    prompt: "another detailed description"

minor_arcana:
  wands:
    suit_element: fire
    suit_color: red
    cards:
      - rank: Ace
        prompt: "card description"
      - rank: Two
        prompt: "card description"
  
  cups:
    suit_element: water
    suit_color: blue
    cards:
      - rank: Ace
        prompt: "card description"
```

### Field Reference

#### Major Arcana Cards

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `number` | integer | Yes | Card number (used for seed calculation and sorting) |
| `name` | string | Yes | Card name (used in filename and progress output) |
| `prompt` | string | Yes | Generation prompt for this specific card |

#### Minor Arcana Structure

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `suit_element` | string | No | Metadata: element associated with suit |
| `suit_color` | string | No | Metadata: color theme for suit |
| `cards` | array | Yes | List of cards in this suit |

#### Card Objects

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `rank` | string | Yes | Card rank (Ace, Two, ..., King) |
| `prompt` | string | Yes | Generation prompt for this card |

## Command Reference

### Basic Usage

```bash
asset-generator pipeline --file pipeline.yaml
```

### All Flags

#### Required Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to pipeline YAML file |

#### Generation Parameters

| Flag | Default | Description |
|------|---------|-------------|
| `--output-dir` | `./pipeline-output` | Output directory for generated assets |
| `--base-seed` | `-1` (random) | Base seed for reproducible generation (0 or -1 for random) |
| `--model` | (none) | Model to use for all generations |
| `--steps` | `40` | Number of inference steps |
| `--width` | `768` | Image width in pixels |
| `--height` | `1344` | Image height in pixels |
| `--cfg-scale` | `7.5` | CFG scale (guidance strength) |
| `--sampler` | `euler_a` | Sampling method |

#### Prompt Enhancement

| Flag | Default | Description |
|------|---------|-------------|
| `--style-suffix` | (none) | Suffix appended to all prompts |
| `--negative-prompt` | (none) | Negative prompt for all generations |

#### Pipeline Control

| Flag | Default | Description |
|------|---------|-------------|
| `--dry-run` | `false` | Preview pipeline without generating |
| `--continue-on-error` | `false` | Continue if individual generations fail |

#### Postprocessing

| Flag | Default | Description |
|------|---------|-------------|
| `--auto-crop` | `false` | Automatically crop whitespace borders |
| `--auto-crop-threshold` | `250` | Whitespace detection threshold (0-255) |
| `--auto-crop-tolerance` | `10` | Tolerance for near-white colors (0-255) |
| `--auto-crop-preserve-aspect` | `false` | Preserve aspect ratio when cropping |
| `--downscale-width` | `0` | Downscale to this width (0=disabled) |
| `--downscale-height` | `0` | Downscale to this height (0=disabled) |
| `--downscale-percentage` | `0` | Downscale by percentage (0=disabled) |
| `--downscale-filter` | `lanczos` | Filter: lanczos, bilinear, nearest |

## Examples

### Example 1: Basic Pipeline

```bash
asset-generator pipeline --file tarot-spec.yaml --output-dir ./my-deck
```

Output structure:
```
my-deck/
├── major-arcana/
│   ├── 00-the_fool.png
│   ├── 01-the_magician.png
│   └── ...
└── minor-arcana/
    ├── wands/
    │   ├── 01-ace_of_wands.png
    │   └── ...
    ├── cups/
    ├── swords/
    └── pentacles/
```

### Example 2: Preview Before Generating

```bash
asset-generator pipeline --file tarot-spec.yaml --dry-run
```

Shows:
- All cards that would be generated
- Calculated seeds for each card
- Generation parameters
- Output structure

### Example 3: Custom Generation Parameters

```bash
asset-generator pipeline --file tarot-spec.yaml \
  --base-seed 1000 \
  --steps 50 \
  --width 1024 \
  --height 1792 \
  --cfg-scale 8.0
```

### Example 4: Style Enhancement

```bash
asset-generator pipeline --file tarot-spec.yaml \
  --style-suffix "detailed illustration, ornate border, rich colors, professional quality" \
  --negative-prompt "blurry, distorted, low quality, modern elements"
```

This appends the style suffix to every prompt in the pipeline.

### Example 5: High-Resolution with Downscaling

Generate at high resolution, then downscale for web use:

```bash
asset-generator pipeline --file tarot-spec.yaml \
  --width 1536 \
  --height 2688 \
  --steps 50 \
  --downscale-width 768 \
  --downscale-filter lanczos
```

### Example 6: Auto-Crop and Resize

Remove whitespace borders and resize to specific dimensions:

```bash
asset-generator pipeline --file tarot-spec.yaml \
  --auto-crop \
  --auto-crop-threshold 245 \
  --downscale-width 768 \
  --downscale-height 1344
```

### Example 7: Robust Production Pipeline

```bash
asset-generator pipeline --file tarot-spec.yaml \
  --output-dir ./production-deck \
  --base-seed 42 \
  --model "XE-_Pixel_Flux_-_0-1.safetensors" \
  --steps 50 \
  --width 1536 \
  --height 2688 \
  --cfg-scale 7.5 \
  --style-suffix "masterpiece, detailed, professional quality" \
  --negative-prompt "blurry, distorted, low quality" \
  --continue-on-error \
  --auto-crop \
  --downscale-width 768
```

### Example 8: Verbose Output for Debugging

```bash
asset-generator pipeline --file tarot-spec.yaml \
  --verbose \
  --dry-run
```

Shows detailed information including full prompts for each card.

## Best Practices

### 1. Always Preview First

Use `--dry-run` to verify your pipeline before generating:

```bash
asset-generator pipeline --file my-pipeline.yaml --dry-run
```

### 2. Use Consistent Seeds

For reproducible results, always use the same `--base-seed`:

```bash
asset-generator pipeline --file deck.yaml --base-seed 42
```

The pipeline automatically calculates unique seeds for each asset based on the base seed.

**Note:** By default (when `--base-seed` is not specified, or set to `0` or `-1`), a random seed is 
generated for each pipeline run. The generated seed is displayed in the output so you can reproduce 
the same results later by explicitly specifying that seed with `--base-seed`.

### 3. Style Consistency

Use `--style-suffix` to ensure consistent styling across all assets:

```bash
--style-suffix "detailed illustration, professional quality, rich colors"
```

This is cleaner than adding the style to every prompt in your YAML file.

### 4. Organize Output

Use descriptive output directories:

```bash
asset-generator pipeline --file cards.yaml --output-dir ./decks/v1-fantasy-theme
```

### 5. Error Handling

For large pipelines, use `--continue-on-error` to avoid stopping on individual failures:

```bash
asset-generator pipeline --file large-set.yaml --continue-on-error
```

You can review failures in the summary and regenerate specific assets later.

### 6. Postprocessing Pipeline

Combine auto-crop and downscaling for optimal results:

1. Generate at high resolution for quality
2. Auto-crop to remove borders
3. Downscale to target size

```bash
asset-generator pipeline --file deck.yaml \
  --width 2048 --height 3584 \
  --auto-crop \
  --downscale-width 1024 \
  --downscale-filter lanczos
```

### 7. Model Selection

Test with different models to find the best fit:

```bash
# Preview with model info
asset-generator models list

# Generate with specific model
asset-generator pipeline --file deck.yaml \
  --model "stable-diffusion-xl-base"
```

### 8. Version Your Pipeline Files

Keep your pipeline files in version control:

```
project/
├── pipelines/
│   ├── tarot-deck-v1.yaml
│   ├── tarot-deck-v2-refined.yaml
│   └── character-sprites.yaml
└── outputs/
    ├── deck-v1/
    └── deck-v2/
```

## Seed Calculation

Understanding how seeds are calculated helps ensure reproducibility:

### Major Arcana
```
seed = base_seed + card_number
```

Example with `--base-seed 42`:
- The Fool (0) → seed 42
- The Magician (1) → seed 43
- The World (21) → seed 63

### Minor Arcana
```
seed = base_seed + 100 + suit_offset + card_index
```

Suit offsets:
- Wands: 0
- Cups: 20
- Swords: 40
- Pentacles: 60

Example with `--base-seed 42`:
- Ace of Wands → seed 142
- King of Wands → seed 155
- Ace of Cups → seed 162
- King of Pentacles → seed 215

## Output Structure

The pipeline automatically creates this directory structure:

```
output-dir/
├── major-arcana/
│   ├── 00-the_fool.png
│   ├── 01-the_magician.png
│   ├── 02-the_high_priestess.png
│   └── ...
└── minor-arcana/
    ├── wands/
    │   ├── 01-ace_of_wands.png
    │   ├── 02-two_of_wands.png
    │   └── ...
    ├── cups/
    │   ├── 01-ace_of_cups.png
    │   └── ...
    ├── swords/
    │   ├── 01-ace_of_swords.png
    │   └── ...
    └── pentacles/
        ├── 01-ace_of_pentacles.png
        └── ...
```

### Filename Format

- Major Arcana: `{number}-{sanitized_name}.png`
- Minor Arcana: `{rank_number}-{sanitized_rank}_of_{suit}.png`

Examples:
- `00-the_fool.png`
- `01-ace_of_wands.png`
- `14-king_of_pentacles.png`

## Troubleshooting

### Pipeline Fails to Load

**Error:** `failed to load pipeline: failed to read file`

**Solution:** Check that the file path is correct and the file exists:

```bash
ls -l tarot-spec.yaml
asset-generator pipeline --file ./examples/tarot-deck/tarot-spec.yaml
```

### YAML Parse Error

**Error:** `failed to parse YAML`

**Solution:** Validate your YAML syntax:

```bash
# Install yamllint
pip install yamllint

# Validate file
yamllint tarot-spec.yaml
```

Common issues:
- Incorrect indentation (use 2 spaces)
- Missing colons after keys
- Unquoted strings with special characters

### Generation Failures

**Error:** `failed to generate {card_name}`

**Solutions:**

1. Check API connection:
```bash
asset-generator models list
```

2. Use `--continue-on-error` to see all failures:
```bash
asset-generator pipeline --file deck.yaml --continue-on-error
```

3. Enable verbose output:
```bash
asset-generator pipeline --file deck.yaml --verbose
```

### Model Not Found

**Error:** `model validation failed: model 'xyz' not found`

**Solution:** List available models and use the correct name:

```bash
asset-generator models list
asset-generator pipeline --file deck.yaml --model "correct-model-name"
```

### Out of Memory

For large pipelines generating many high-resolution images:

1. Reduce dimensions:
```bash
--width 1024 --height 1792
```

2. Process in batches by creating multiple smaller pipeline files

3. Use downscaling instead of generating at final resolution:
```bash
--width 2048 --height 3584 --downscale-width 1024
```

### Slow Generation

**Issue:** Pipeline taking too long

**Solutions:**

1. Reduce steps:
```bash
--steps 20
```

2. Use faster sampler:
```bash
--sampler euler_a
```

3. Generate lower resolution with upscaling later:
```bash
--width 512 --height 896
```

### Interrupted Pipeline

If the pipeline is interrupted (Ctrl+C), it will stop gracefully. To resume:

1. Check which assets were generated
2. Remove completed cards from your pipeline file
3. Re-run with the same seed settings

Or use `--continue-on-error` to skip already-generated files (you'll need to check manually).

## Advanced Usage

### Custom Pipeline Structures

While the tarot deck format is the default, you can adapt the structure for other use cases:

#### Game Sprites Example

```yaml
major_arcana:  # Use for main character sprites
  - number: 0
    name: Hero Idle
    prompt: "pixel art character, hero idle pose"
  - number: 1
    name: Hero Walk
    prompt: "pixel art character, hero walking"

minor_arcana:
  weapons:  # Use suits for categories
    cards:
      - rank: Sword
        prompt: "pixel art sword weapon icon"
      - rank: Axe
        prompt: "pixel art axe weapon icon"
```

### Integration with Shell Scripts

You can still integrate with scripts for additional processing:

```bash
#!/bin/bash
# Generate deck
asset-generator pipeline --file deck.yaml --output-dir ./raw

# Additional processing
for file in ./raw/major-arcana/*.png; do
  # Custom postprocessing
  convert "$file" -quality 95 "./processed/$(basename "$file")"
done
```

### Batch Multiple Pipelines

```bash
#!/bin/bash
PIPELINES=(
  "deck-light.yaml"
  "deck-dark.yaml"
  "deck-vintage.yaml"
)

for pipeline in "${PIPELINES[@]}"; do
  asset-generator pipeline --file "$pipeline" \
    --output-dir "./output/$(basename "$pipeline" .yaml)"
done
```

## See Also

- [QUICKSTART.md](../QUICKSTART.md) - Getting started guide
- [README.md](../README.md) - Main documentation
- [examples/tarot-deck/](../examples/tarot-deck/) - Complete tarot deck example
- [FILENAME_TEMPLATES.md](./FILENAME_TEMPLATES.md) - Filename customization
