# Asset Generation Pipeline for 2D Games

This guide demonstrates how to integrate `asset-generator` into a 2D game asset generation pipeline using the built-in `pipeline` command for batch processing.

## Overview

`asset-generator` provides a command-line interface to AI asset generation APIs, enabling automated generation of game assets like sprites, backgrounds, UI elements, and textures through text-to-image AI models.

The **`pipeline` command** allows you to define complete asset generation workflows in YAML files and process them automatically, eliminating the need for shell scripts and external dependencies.

## Prerequisites

- Asset generation service running (default: `http://localhost:7801`)
- `asset-generator` CLI installed and configured
- Stable Diffusion or compatible model loaded in the service

## Quick Setup

```bash
# Initialize configuration
asset-generator config init

# Set your asset generation service endpoint (if not default)
asset-generator config set api-url http://localhost:7801

# Verify connection
asset-generator models list
```

## Quick Reference: Pipeline Processing

The `pipeline` command is the recommended way to batch-generate assets. It processes YAML pipeline files automatically.

```bash
# Basic pipeline processing
asset-generator pipeline --file assets.yaml

# Preview before generating
asset-generator pipeline --file assets.yaml --dry-run

# Production generation with postprocessing
asset-generator pipeline --file assets.yaml \
  --base-seed 42 \
  --steps 50 \
  --auto-crop \
  --downscale-width 1024
```

**Key Benefits:**
- ✅ No shell scripts or external dependencies needed
- ✅ Cross-platform (Windows, Linux, macOS)
- ✅ Built-in progress tracking
- ✅ Automatic directory organization
- ✅ Reproducible with seed-based generation
- ✅ Integrated postprocessing

For complete pipeline documentation, see [docs/PIPELINE.md](docs/PIPELINE.md).

## Alternative: Manual Generation

For single assets or simple workflows, use the `generate image` command:

```bash
# Basic download to current directory
asset-generator generate image --prompt "sprite" --save-images

# Download to specific directory
asset-generator generate image --prompt "sprite" --save-images --output-dir "./assets"

# Download with custom filenames
asset-generator generate image \
  --prompt "sprite" \
  --save-images \
  --output-dir "./assets" \
  --filename-template "sprite-{i1}.png"

# Batch download with organized naming
asset-generator generate image \
  --prompt "sprite" \
  --batch 5 \
  --save-images \
  --output-dir "./assets" \
  --filename-template "{date}/sprite-{i1}.png"
```

**Key Flags:**
- `--save-images` - Enable image downloading
- `--output-dir <path>` - Specify download directory (default: current directory)
- `--filename-template <template>` - Customize filenames with placeholders

**Common Placeholders:**
- `{i1}` - One-based index (1, 2, 3, ...)
- `{index}` - Zero-padded index (000, 001, 002, ...)
- `{date}` - Date (2024-10-09)
- `{seed}` - Seed value
- `{model}` - Model name
- `{width}x{height}` - Dimensions
- `{prompt}` - Sanitized prompt (first 50 chars)

For complete placeholder documentation, see [docs/FILENAME_TEMPLATES.md](docs/FILENAME_TEMPLATES.md).

**Available Placeholders:**
- **Index:** `{index}` (zero-padded), `{i1}` (one-based)
- **Time:** `{date}`, `{time}`, `{datetime}`, `{timestamp}`
- **Parameters:** `{seed}`, `{model}`, `{width}`, `{height}`, `{prompt}`
- **Original:** `{original}`, `{ext}`

### Pipeline Integration Examples

#### Simple Sequential Downloads

```bash
# Download hero variants with sequential numbering
asset-generator generate image \
  --prompt "brave knight character" \
  --batch 3 \
  --save-images \
  --output-dir "./assets/characters" \
  --filename-template "hero-{i1}.png"
# Generates: hero-1.png, hero-2.png, hero-3.png
```

#### Organized by Category and Date

```bash
# Create organized directory structure
asset-generator generate image \
  --prompt "dungeon background" \
  --save-images \
  --output-dir "./assets" \
  --filename-template "{date}/backgrounds/dungeon-{index}.png"
# Generates: ./assets/2024-10-09/backgrounds/dungeon-000.png
```

#### Reproducible Asset Names

```bash
# Include seed for regeneration tracking
asset-generator generate image \
  --prompt "magic spell effect" \
  --seed 42 \
  --save-images \
  --output-dir "./assets/effects" \
  --filename-template "spell-seed{seed}-{i1}.png"
# Generates: spell-seed42-1.png
```

## Creating Asset Generation Pipelines

### Method 1: Pipeline Command (Recommended)

The `pipeline` command provides a native, cross-platform solution for batch asset generation.

#### Step 1: Create Pipeline YAML File

Create `game-assets.yaml`:

```yaml
# Game Asset Generation Pipeline
# Organized structure for sprite generation

major_arcana:
  - number: 0
    name: Hero Knight
    prompt: "pixel art character sprite, brave knight with sword and shield, 32x32, game asset, transparent background, centered, detailed"
  
  - number: 1
    name: Enemy Goblin
    prompt: "pixel art monster sprite, goblin warrior with club, 32x32, game asset, transparent background, centered, aggressive pose"
  
  - number: 2
    name: Enemy Orc
    prompt: "pixel art monster sprite, orc brute with axe, 32x32, game asset, transparent background, centered, menacing"
  
  - number: 3
    name: Forest Background
    prompt: "pixel art forest background, parallax layer, game background, vibrant green trees, depth, atmospheric"

minor_arcana:
  ui:
    suit_element: interface
    suit_color: blue
    cards:
      - rank: Health Bar
        prompt: "game UI health bar, fantasy style, red and gold gradient, transparent background, glossy effect"
      - rank: Button
        prompt: "game UI button, fantasy style, stone texture with glowing blue edges, transparent background"
      - rank: Icon Sword
        prompt: "game UI icon, sword weapon, pixel art style, 64x64, transparent background, centered"
      - rank: Icon Potion
        prompt: "game UI icon, health potion, red liquid, pixel art style, 64x64, transparent background, centered"
  
  items:
    suit_element: collectibles
    suit_color: gold
    cards:
      - rank: Gold Coin
        prompt: "pixel art gold coin sprite, shiny, 16x16, game asset, transparent background, centered, spinning frame"
      - rank: Chest
        prompt: "pixel art treasure chest, closed, wooden with metal bands, 32x32, game asset, transparent background"
      - rank: Key
        prompt: "pixel art golden key sprite, ornate design, 16x16, game asset, transparent background, centered"
```

#### Step 2: Generate Assets

```bash
# Preview the pipeline first
asset-generator pipeline --file game-assets.yaml --dry-run

# Generate all assets
asset-generator pipeline --file game-assets.yaml \
  --output-dir ./game-assets \
  --base-seed 42 \
  --steps 40 \
  --width 512 \
  --height 512 \
  --style-suffix "professional game art, crisp edges, vibrant colors" \
  --negative-prompt "blurry, distorted, low quality, watermark"
```

#### Step 3: Add Postprocessing

```bash
# Generate with auto-crop and downscaling
asset-generator pipeline --file game-assets.yaml \
  --output-dir ./game-assets \
  --base-seed 42 \
  --steps 40 \
  --width 1024 \
  --height 1024 \
  --auto-crop \
  --downscale-width 512 \
  --downscale-filter lanczos \
  --continue-on-error
```

**Output Structure:**
```
game-assets/
├── major-arcana/
│   ├── 00-hero_knight.png
│   ├── 01-enemy_goblin.png
│   ├── 02-enemy_orc.png
│   └── 03-forest_background.png
└── minor-arcana/
    ├── ui/
    │   ├── 01-health_bar.png
    │   ├── 02-button.png
    │   ├── 03-icon_sword.png
    │   └── 04-icon_potion.png
    └── items/
        ├── 01-gold_coin.png
        ├── 02-chest.png
        └── 03-key.png
```

**Pipeline Benefits:**
- ✅ Single command generates all assets
- ✅ Automatic directory organization
- ✅ Reproducible with seed-based generation
- ✅ Progress tracking built-in
- ✅ Continue on error support
- ✅ No shell scripts or yq needed
- ✅ Cross-platform compatible

For complete pipeline documentation, see [docs/PIPELINE.md](docs/PIPELINE.md).

### Method 2: Individual Generation (Simple Workflows)

For simple workflows or single assets, you can use individual `generate image` commands:

```bash
# Generate a single sprite
asset-generator generate image \
  --prompt "pixel art knight sprite" \
  --width 512 \
  --height 512 \
  --save-images \
  --output-dir "./assets/characters" \
  --filename-template "knight.png"

# Generate UI element
asset-generator generate image \
  --prompt "game UI health bar" \
  --width 256 \
  --height 64 \
  --save-images \
  --output-dir "./assets/ui" \
  --filename-template "health_bar.png"
```

## Complete Example: Tarot Deck Generation

The [examples/tarot-deck/](examples/tarot-deck/) directory demonstrates a complete 78-card tarot deck generation pipeline.

### Using the Pipeline Command

```bash
cd examples/tarot-deck

# Preview the full deck before generating
asset-generator pipeline --file tarot-spec.yaml --dry-run

# Generate all 78 cards
asset-generator pipeline --file tarot-spec.yaml \
  --output-dir ./tarot-deck-output \
  --base-seed 42 \
  --steps 40 \
  --width 768 \
  --height 1344 \
  --style-suffix "detailed illustration, ornate border, rich colors" \
  --negative-prompt "blurry, distorted, low quality"

# High-quality production with postprocessing
asset-generator pipeline --file tarot-spec.yaml \
  --output-dir ./production-deck \
  --base-seed 42 \
  --steps 50 \
  --width 1536 \
  --height 2688 \
  --auto-crop \
  --downscale-width 768 \
  --continue-on-error
```

**Output Structure:**
```
tarot-deck-output/
├── major-arcana/
│   ├── 00-the_fool.png
│   ├── 01-the_magician.png
│   └── ... (22 cards total)
└── minor-arcana/
    ├── wands/ (14 cards)
    ├── cups/ (14 cards)
    ├── swords/ (14 cards)
    └── pentacles/ (14 cards)
```

**Key Features Demonstrated:**
- 78 unique cards with detailed prompts
- Organized Major/Minor Arcana structure
- Seed-based reproducible generation
- Automatic directory creation
- Progress tracking
- Error recovery with continue-on-error

See [examples/tarot-deck/README.md](examples/tarot-deck/README.md) for complete details.

## Advanced Integration

### Makefile Integration

Add to your game's `Makefile`:

```makefile
.PHONY: assets assets-clean assets-verify

ASSET_PIPELINE := game-assets.yaml
ASSET_DIR := assets/generated

assets: assets-clean
	@echo "Generating game assets..."
	@asset-generator pipeline --file $(ASSET_PIPELINE) \
		--output-dir $(ASSET_DIR) \
		--base-seed 42 \
		--continue-on-error
	@echo "Assets generated successfully!"

assets-clean:
	@rm -rf $(ASSET_DIR)

assets-verify:
	@asset-generator config view
	@asset-generator models list
	@echo "Asset generation system ready"

assets-preview:
	@asset-generator pipeline --file $(ASSET_PIPELINE) --dry-run
```

Then simply run:
```bash
make assets        # Generate all assets
make assets-preview  # Preview pipeline
```

### CI/CD Integration

Add to your `.github/workflows/assets.yml`:

```yaml
name: Generate Assets

on:
  push:
    paths:
      - 'assets-spec.yaml'
  workflow_dispatch:

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install asset-generator
        run: |
          wget https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-linux-amd64
          chmod +x asset-generator-linux-amd64
          sudo mv asset-generator-linux-amd64 /usr/local/bin/asset-generator
      
      - name: Configure
        run: |
          asset-generator config set api-url ${{ secrets.API_URL }}
          asset-generator config set api-key ${{ secrets.API_KEY }}
      
      - name: Generate Assets
        run: |
          asset-generator pipeline --file assets-spec.yaml \
            --output-dir ./generated-assets \
            --continue-on-error
      
      - name: Upload Artifacts
        uses: actions/upload-artifact@v3
        with:
          name: generated-assets
          path: generated-assets/
```

## Best Practices

### 1. Consistent Prompting

Use standardized prompt templates:

```sh
# Character sprites template
generate_character() {
    name="$1"
    description="$2"
    asset-generator generate image \
        --prompt "pixel art character sprite, ${description}, 32x32, game asset, transparent background, clean lines" \
        --width 512 --height 512 \
        --save-images \
        --output-dir "assets/characters" \
        --filename-template "${name}-{i1}.png"
}

# Background template
generate_background() {
    name="$1"
    description="$2"
    asset-generator generate image \
        --prompt "pixel art background, ${description}, parallax layer, game background, vibrant colors" \
        --width 1024 --height 512 \
        --save-images \
        --output-dir "assets/backgrounds" \
        --filename-template "${name}-{i1}.png"
}

# UI element template
generate_ui() {
    name="$1"
    element="$2"
    style="$3"
    asset-generator generate image \
        --prompt "game UI ${element}, ${style} style, clean design, transparent background" \
        --width 256 --height 128 \
        --save-images \
        --output-dir "assets/ui" \
        --filename-template "${name}.png"
}

# Usage:
generate_character "hero" "brave knight with sword and shield"
generate_background "forest" "dense forest with ancient trees"
generate_ui "health_bar" "health bar" "fantasy"
```

### 2. Reproducible Generations

Use seeds for consistent results and track them in filenames:

```bash
asset-generator generate image \
    --prompt "hero character sprite" \
    --seed 42 \
    --width 512 --height 512 \
    --save-images \
    --output-dir "assets/characters" \
    --filename-template "hero-seed{seed}.png"
# Generates: assets/characters/hero-seed42.png
```

### 3. Batch Generation

Generate variations efficiently with organized filenames:

```bash
# Generate 4 enemy variations with sequential naming
asset-generator generate image \
    --prompt "pixel art monster sprite, game asset" \
    --batch 4 \
    --width 512 --height 512 \
    --save-images \
    --output-dir "assets/enemies" \
    --filename-template "enemy-variant-{i1}.png"
# Generates: enemy-variant-1.png, enemy-variant-2.png, ..., enemy-variant-4.png
```

### 4. Model Selection

Choose appropriate models for your art style and include model name in filenames:

```bash
# List available models
asset-generator models list

# Use specific model for pixel art and track it in filename
asset-generator generate image \
    --prompt "pixel art sprite" \
    --model "pixel-art-diffusion" \
    --save-images \
    --output-dir "assets/sprites" \
    --filename-template "{model}-sprite-{i1}.png"
# Generates: pixel-art-diffusion-sprite-1.png
```

## Pipeline Stages

### Stage 1: Concept Generation

```bash
# Generate concept art variants with timestamps
asset-generator generate image \
    --prompt "fantasy knight character design, multiple poses" \
    --batch 5 \
    --steps 50 \
    --width 1024 --height 1024 \
    --save-images \
    --output-dir "concepts" \
    --filename-template "{date}-concept-{i1}.png"
# Generates: 2024-10-09-concept-1.png, ..., 2024-10-09-concept-5.png
```

### Stage 2: Asset Production

```bash
# Generate final game assets with descriptive names
asset-generator generate image \
    --prompt "pixel art knight sprite, 32x32, transparent" \
    --width 512 --height 512 \
    --steps 30 \
    --save-images \
    --output-dir "assets/final" \
    --filename-template "knight-{width}x{height}.png"
# Generates: assets/final/knight-512x512.png
```

### Stage 3: Variations

```sh
# Generate color variations with organized naming
for color in red blue green yellow; do
    asset-generator generate image \
        --prompt "knight sprite with ${color} armor" \
        --seed 42 \
        --width 512 --height 512 \
        --save-images \
        --output-dir "assets/knights" \
        --filename-template "knight-${color}-seed{seed}.png"
done
# Generates: knight-red-seed42.png, knight-blue-seed42.png, etc.
```

## Advanced Filename Template Patterns

### Pattern 1: Date-Based Organization

Organize assets by date for daily generation workflows:

```bash
asset-generator generate image \
    --prompt "daily asset generation" \
    --batch 10 \
    --save-images \
    --output-dir "assets" \
    --filename-template "{date}/asset-{i1}.png"
# Generates: assets/2024-10-09/asset-1.png, assets/2024-10-09/asset-2.png, ...
```

### Pattern 2: Multi-Level Directory Structure

Create complex hierarchies using slashes in templates:

```bash
asset-generator generate image \
    --prompt "fantasy castle background" \
    --model "sdxl-turbo" \
    --save-images \
    --output-dir "game_assets" \
    --filename-template "{date}/backgrounds/{model}/castle-{width}x{height}.png"
# Generates: game_assets/2024-10-09/backgrounds/sdxl-turbo/castle-1024x768.png
```

### Pattern 3: Prompt-Based Naming

Use the prompt in filenames for easy identification:

```bash
asset-generator generate image \
    --prompt "cyberpunk street scene" \
    --save-images \
    --output-dir "scenes" \
    --filename-template "{prompt}-{datetime}.png"
# Generates: cyberpunk_street_scene-2024-10-09_14-30-45.png
# Note: Prompt is sanitized (spaces → underscores, max 50 chars)
```

### Pattern 4: Complete Metadata Tracking

Include all relevant parameters for full traceability:

```bash
asset-generator generate image \
    --prompt "warrior character" \
    --model "flux-dev" \
    --width 1024 \
    --height 768 \
    --seed 12345 \
    --save-images \
    --output-dir "characters" \
    --filename-template "{model}-{width}x{height}-seed{seed}-{prompt}.png"
# Generates: flux-dev-1024x768-seed12345-warrior_character.png
```

### Pattern 5: Loop-Based Batch Processing

Combine shell loops with filename templates:

```sh
#!/bin/sh
# Generate multiple asset types with consistent naming

TYPES="hero enemy npc boss"
OUTPUT_BASE="assets/characters"

for type in $TYPES; do
    asset-generator generate image \
        --prompt "pixel art ${type} character sprite" \
        --batch 3 \
        --seed 42 \
        --save-images \
        --output-dir "$OUTPUT_BASE" \
        --filename-template "${type}-{i1}-seed{seed}.png"
done

# Generates:
#   assets/characters/hero-1-seed42.png
#   assets/characters/hero-2-seed42.png
#   assets/characters/hero-3-seed42.png
#   assets/characters/enemy-1-seed42.png
#   ... and so on
```

### Pattern 6: Environment-Specific Organization

Use environment variables for flexible pipeline configuration:

```sh
#!/bin/sh
# Production vs development asset generation

ENVIRONMENT="${DEPLOY_ENV:-dev}"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

asset-generator generate image \
    --prompt "game logo" \
    --save-images \
    --output-dir "assets/${ENVIRONMENT}" \
    --filename-template "logo-${ENVIRONMENT}-${TIMESTAMP}.png"

# Dev: assets/dev/logo-dev-20241009-143045.png
# Prod: assets/prod/logo-prod-20241009-143045.png
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Generate Game Assets

on:
    workflow_dispatch:
    push:
        paths:
            - 'assets/specs/*.yaml'

jobs:
    generate:
        runs-on: ubuntu-latest
        steps:
            - uses: actions/checkout@v3
            
            - name: Setup Asset Generator CLI
              run: |
                  wget https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-linux-amd64
                  chmod +x asset-generator-linux-amd64
                  sudo mv asset-generator-linux-amd64 /usr/local/bin/asset-generator
            
            - name: Configure CLI
              run: |
                  asset-generator config set api-url ${{ secrets.ASSET_GENERATOR_URL }}
                  asset-generator config set api-key ${{ secrets.ASSET_GENERATOR_KEY }}
            
            - name: Generate Assets with Download
              run: |
                  # Generate and download assets with organized filenames
                  asset-generator generate image \
                    --prompt "game character sprite" \
                    --batch 5 \
                    --save-images \
                    --output-dir generated_assets \
                    --filename-template "character-{date}-{i1}.png"
            
            - name: Upload Assets
              uses: actions/upload-artifact@v3
              with:
                  name: game-assets-${{ github.sha }}
                  path: generated_assets/
```

### GitLab CI Example

```yaml
generate_assets:
  stage: build
  image: alpine:latest
  before_script:
    - apk add --no-cache wget
    - wget -O /usr/local/bin/asset-generator https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-linux-amd64
    - chmod +x /usr/local/bin/asset-generator
    - asset-generator config set api-url $ASSET_GENERATOR_URL
  script:
    # Generate with organized directory structure
    - |
      asset-generator generate image \
        --prompt "background asset" \
        --batch 3 \
        --save-images \
        --output-dir ci_assets \
        --filename-template "{date}/bg-{i1}.png"
  artifacts:
    paths:
      - ci_assets/
    expire_in: 30 days
  only:
    - main
```

## Troubleshooting

### Connection Issues

```bash
# Verify asset generation service is accessible
asset-generator models list

# Check configuration
asset-generator config view

# Test with verbose output
asset-generator generate image --prompt "test" --verbose
```

### Quality Issues

```bash
# Increase steps for better quality
asset-generator generate image --prompt "detailed sprite" --steps 50

# Adjust CFG scale for prompt adherence
asset-generator generate image --prompt "sprite" --cfg-scale 8.0

# Use negative prompts
asset-generator generate image \
    --prompt "clean sprite" \
    --negative-prompt "blurry, distorted, low quality"
```

## Current Capabilities & Features

**Available Features:**
- ✅ **Pipeline Processing**: Native YAML batch generation with `pipeline` command
- ✅ Text-to-image generation
- ✅ Batch generation
- ✅ Automatic image download with `--save-images`
- ✅ Custom output directories with `--output-dir`
- ✅ Filename templates with rich placeholders
- ✅ JSON/YAML output formats
- ✅ Seed-based reproducibility
- ✅ Image post-processing (downscale, crop, SVG conversion)
- ✅ Progress tracking and error recovery
- ✅ Dry-run preview mode
- ✅ Cross-platform support (Windows/Linux/macOS)

**Pipeline Command Benefits:**
- ✅ No external dependencies (no yq, no shell scripts)
- ✅ Automatic directory organization
- ✅ Built-in progress tracking
- ✅ Continue-on-error support
- ✅ Integrated postprocessing
- ✅ Style suffix application
- ✅ Reproducible seed-based generation

**Filename Template Placeholders:**
- Index: `{index}`, `{i1}`
- Time: `{date}`, `{time}`, `{datetime}`, `{timestamp}`
- Parameters: `{seed}`, `{model}`, `{width}`, `{height}`, `{prompt}`
- Original: `{original}`, `{ext}`

For complete documentation:
- Pipeline: [docs/PIPELINE.md](docs/PIPELINE.md)
- Filename templates: [docs/FILENAME_TEMPLATES.md](docs/FILENAME_TEMPLATES.md)
- Image download: [docs/IMAGE_DOWNLOAD.md](docs/IMAGE_DOWNLOAD.md)

## Example Output

### Pipeline Command Output

After running `asset-generator pipeline --file game-assets.yaml`, you'll have:

```
game-assets/
├── major-arcana/
│   ├── 00-hero_knight.png
│   ├── 01-enemy_goblin.png
│   ├── 02-enemy_orc.png
│   └── 03-forest_background.png
└── minor-arcana/
    ├── ui/
    │   ├── 01-health_bar.png
    │   ├── 02-button.png
    │   ├── 03-icon_sword.png
    │   └── 04-icon_potion.png
    └── items/
        ├── 01-gold_coin.png
        ├── 02-chest.png
        └── 03-key.png
```

### Individual Generation Output

Using `--save-images` with individual commands:

```
generated_assets/
├── characters/
│   ├── hero.png
│   └── enemy.png
├── backgrounds/
│   ├── forest.png
│   └── dungeon.png
└── ui/
    ├── health_bar.png
    └── button.png
```

### With Filename Templates

Using individual generation with `--filename-template "{date}/asset-{i1}.png"`:

```
generated_assets/
└── 2024-10-09/
    ├── asset-1.png
    ├── asset-2.png
    └── asset-3.png
```

### With Complex Templates

Using `--filename-template "{model}-{width}x{height}-seed{seed}-{i1}.png"`:

```
generated_assets/
├── flux-dev-1024x768-seed42-1.png
├── flux-dev-1024x768-seed42-2.png
└── flux-dev-1024x768-seed42-3.png
```

## Next Steps

1. **Start with Pipeline Command**: Use `asset-generator pipeline` for batch workflows
2. **Preview First**: Always use `--dry-run` before generating
3. **Organize Assets**: Use the built-in directory structure
4. **Integrate with Build**: Add to Makefile or CI/CD pipeline
5. **Post-Process**: Leverage `--auto-crop` and `--downscale-*` flags
6. **Game Engine Integration**: Import generated assets into your game engine

For more information, see:
- [docs/PIPELINE.md](docs/PIPELINE.md) - Complete pipeline documentation
- [docs/PIPELINE_QUICKREF.md](docs/PIPELINE_QUICKREF.md) - Quick reference
- [examples/tarot-deck/](examples/tarot-deck/) - Complete working example
- [README.md](README.md) - Complete CLI documentation
- [QUICKSTART.md](QUICKSTART.md) - Getting started guide
- [docs/FILENAME_TEMPLATES.md](docs/FILENAME_TEMPLATES.md) - Filename template guide
- [docs/IMAGE_DOWNLOAD.md](docs/IMAGE_DOWNLOAD.md) - Image download documentation