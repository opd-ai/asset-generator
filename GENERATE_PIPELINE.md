# Asset Generation Pipeline for 2D Games

This guide demonstrates how to integrate `asset-generator` into a 2D game asset generation pipeline.

## Overview

`asset-generator` provides a command-line interface to AI asset generation APIs, enabling automated generation of game assets like sprites, backgrounds, UI elements, and textures through text-to-image AI models.

## Prerequisites

- Asset generation service running (default: `http://localhost:7801`)
- `asset-generator` CLI installed and configured
- Stable Diffusion or compatible model loaded in the service
- POSIX-compliant shell (bash, dash, zsh, etc.)

## Quick Setup

```bash
# Initialize configuration
asset-generator config init

# Set your asset generation service endpoint (if not default)
asset-generator config set api-url http://localhost:7801

# Verify connection
asset-generator models list
```

## Quick Reference: Downloading Assets

By default, `asset-generator` outputs JSON metadata with server-side image paths. To download images directly:

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

## Downloading Assets in the Pipeline

The `asset-generator` CLI can automatically download generated images to disk during the generation process, eliminating the need to manually retrieve them from the server's output directory. This is essential for automated pipelines.

### Output Directory Option

Use `--save-images` and `--output-dir` to download images directly:

```bash
asset-generator generate image \
  --prompt "hero character sprite" \
  --save-images \
  --output-dir "./assets/characters"
```

**Without these flags:** The CLI only outputs JSON metadata with server-side image paths.  
**With these flags:** Images are downloaded to your specified directory.

### Filename Template Option

Use `--filename-template` to control how downloaded images are named:

```bash
asset-generator generate image \
  --prompt "forest background" \
  --save-images \
  --output-dir "./assets/backgrounds" \
  --filename-template "bg-forest-{index}.png"
```

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

## Basic Pipeline Example

### 1. Define Asset Requirements

Create an asset specification file `assets.yaml`:

```yaml
characters:
    - name: hero
        prompt: "pixel art character sprite, brave knight with sword and shield, 32x32, game asset, transparent background"
        width: 512
        height: 512
    - name: enemy
        prompt: "pixel art monster sprite, goblin warrior, 32x32, game asset, transparent background"
        width: 512
        height: 512

backgrounds:
    - name: forest
        prompt: "pixel art forest background, parallax layer, game background, vibrant colors"
        width: 1024
        height: 512
    - name: dungeon
        prompt: "pixel art dungeon interior, stone walls and torches, game background"
        width: 1024
        height: 512

ui:
    - name: health_bar
        prompt: "game UI health bar, fantasy style, red and gold, transparent background"
        width: 256
        height: 64
    - name: button
        prompt: "game UI button, fantasy style, stone texture, glowing edges"
        width: 256
        height: 128
```

### 2. Generation Script with Downloads

Create `generate_assets.sh`:

```bash
#!/bin/sh
set -e

OUTPUT_DIR="generated_assets"

# Function to generate and download asset
generate_asset() {
    prompt="$1"
    width="$2"
    height="$3"
    category="$4"
    name="$5"
    
    echo "Generating: $category/$name"
    asset-generator generate image \
        --prompt "$prompt" \
        --width "$width" \
        --height "$height" \
        --steps 30 \
        --cfg-scale 7.5 \
        --save-images \
        --output-dir "$OUTPUT_DIR/$category" \
        --filename-template "${name}.png"
}

# Generate character sprites
generate_asset \
    "pixel art character sprite, brave knight with sword and shield, 32x32, game asset, transparent background" \
    512 512 \
    "characters" "hero"

generate_asset \
    "pixel art monster sprite, goblin warrior, 32x32, game asset, transparent background" \
    512 512 \
    "characters" "enemy"

# Generate backgrounds
generate_asset \
    "pixel art forest background, parallax layer, game background, vibrant colors" \
    1024 512 \
    "backgrounds" "forest"

generate_asset \
    "pixel art dungeon interior, stone walls and torches, game background" \
    1024 512 \
    "backgrounds" "dungeon"

# Generate UI elements
generate_asset \
    "game UI health bar, fantasy style, red and gold, transparent background" \
    256 64 \
    "ui" "health_bar"

generate_asset \
    "game UI button, fantasy style, stone texture, glowing edges" \
    256 128 \
    "ui" "button"

echo "Asset generation complete!"
echo "Assets saved to: $OUTPUT_DIR/"
```

### 3. Execute Pipeline

```bash
chmod +x generate_assets.sh
./generate_assets.sh
```

## Advanced Pipeline Integration

### Shell Script with YAML Parsing

For more complex pipelines, use a POSIX shell script with `yq`:

```sh
#!/bin/sh
# generate_from_spec.sh - Generate and download assets from YAML specification

set -e

SPEC_FILE="${1:-assets.yaml}"
OUTPUT_DIR="${2:-generated_assets}"

# Check dependencies
command -v yq >/dev/null 2>&1 || { echo "Error: yq is required but not installed" >&2; exit 1; }
command -v asset-generator >/dev/null 2>&1 || { echo "Error: asset-generator CLI not found" >&2; exit 1; }

# Parse YAML and generate assets
parse_and_generate() {
    category="$1"
    
    # Get number of assets in category
    count=$(yq eval ".${category} | length" "$SPEC_FILE")
    
    if [ "$count" = "0" ] || [ "$count" = "null" ]; then
        return 0
    fi
    
    # Iterate through assets
    i=0
    while [ "$i" -lt "$count" ]; do
        name=$(yq eval ".${category}[$i].name" "$SPEC_FILE")
        prompt=$(yq eval ".${category}[$i].prompt" "$SPEC_FILE")
        width=$(yq eval ".${category}[$i].width // 512" "$SPEC_FILE")
        height=$(yq eval ".${category}[$i].height // 512" "$SPEC_FILE")
        
        echo "Generating $category/$name..."
        asset-generator generate image \
            --prompt "$prompt" \
            --width "$width" \
            --height "$height" \
            --steps 30 \
            --cfg-scale 7.5 \
            --save-images \
            --output-dir "$OUTPUT_DIR/$category" \
            --filename-template "${name}.png"
        
        i=$((i + 1))
    done
}

# Process all categories
for category in characters backgrounds ui; do
    echo "Processing category: $category"
    parse_and_generate "$category"
done

echo "All assets generated successfully!"
```

### Makefile Integration

Add to your game's `Makefile`:

```makefile
.PHONY: assets assets-clean assets-verify

ASSET_DIR := assets/generated

assets: assets-clean
    @echo "Generating game assets..."
    @./scripts/generate_assets.sh
    @echo "Assets generated successfully!"

assets-clean:
    @rm -rf $(ASSET_DIR)
    @mkdir -p $(ASSET_DIR)

assets-verify:
    @asset-generator config view
    @asset-generator models list
    @echo "Asset generation system ready"
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
- ✅ Text-to-image generation
- ✅ Batch generation
- ✅ Automatic image download with `--save-images`
- ✅ Custom output directories with `--output-dir`
- ✅ Filename templates with rich placeholders
- ✅ JSON/YAML output formats
- ✅ Seed-based reproducibility
- ✅ Image post-processing (downscale, crop, SVG conversion)

**Filename Template Placeholders:**
- Index: `{index}`, `{i1}`
- Time: `{date}`, `{time}`, `{datetime}`, `{timestamp}`
- Parameters: `{seed}`, `{model}`, `{width}`, `{height}`, `{prompt}`
- Original: `{original}`, `{ext}`

For complete filename template documentation, see [docs/FILENAME_TEMPLATES.md](docs/FILENAME_TEMPLATES.md).

**Missing Features (See PLAN.md):**
- ❌ Built-in batch processing from YAML specs
- ❌ Image-to-image generation (for iterations)
- ❌ Inpainting/outpainting support
- ❌ Built-in template system for common asset types

For planned feature implementation, see [PLAN.md](PLAN.md).

## Example Output

After running the pipeline with `--save-images`, you'll have actual image files:

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

Using `--filename-template "{date}/asset-{i1}.png"`:

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

### JSON Output (without --save-images)

When not using `--save-images`, output contains metadata with server image paths:

```json
{
    "session_id": "abc123",
    "prompt": "pixel art character sprite...",
    "image_paths": ["/path/to/server/output/image.png"],
    "parameters": {
        "width": 512,
        "height": 512,
        "steps": 30
    }
}
```

## Next Steps

1. Use `--save-images --output-dir` to download assets directly
2. Leverage `--filename-template` for organized asset naming
3. Integrate with your game engine's asset pipeline
4. Set up automated regeneration on specification changes
5. Implement post-processing for game-ready formats (see `asset-generator convert --help`)

For more information, see:
- [README.md](README.md) - Complete CLI documentation
- [QUICKSTART.md](QUICKSTART.md) - Getting started guide
- [docs/FILENAME_TEMPLATES.md](docs/FILENAME_TEMPLATES.md) - Complete filename template guide
- [docs/IMAGE_DOWNLOAD.md](docs/IMAGE_DOWNLOAD.md) - Image download feature documentation
- [DEVELOPMENT.md](DEVELOPMENT.md) - Developer documentation