# Asset Generation Pipeline for 2D Games

This guide demonstrates how to integrate `asset-generator` (SwarmUI CLI) into a 2D game asset generation pipeline.

## Overview

`asset-generator` provides a command-line interface to SwarmUI, enabling automated generation of game assets like sprites, backgrounds, UI elements, and textures through text-to-image AI models.

## Prerequisites

- SwarmUI instance running (default: `http://localhost:7801`)
- `asset-generator` CLI installed and configured
- Stable Diffusion or compatible model loaded in SwarmUI
- POSIX-compliant shell (bash, dash, zsh, etc.)

## Quick Setup

```bash
# Initialize configuration
asset-generator config init

# Set your SwarmUI endpoint (if not default)
asset-generator config set api-url http://localhost:7801

# Verify connection
asset-generator models list
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

### 2. Generation Script

Create `generate_assets.sh`:

```bash
#!/bin/sh
set -e

OUTPUT_DIR="generated_assets"
mkdir -p "$OUTPUT_DIR/characters" "$OUTPUT_DIR/backgrounds" "$OUTPUT_DIR/ui"

# Function to generate and extract image path
generate_asset() {
    prompt="$1"
    width="$2"
    height="$3"
    output_file="$4"
    
    echo "Generating: $output_file"
    asset-generator generate image \
        --prompt "$prompt" \
        --width "$width" \
        --height "$height" \
        --steps 30 \
        --cfg-scale 7.5 \
        --format json \
        --output "$output_file"
}

# Generate character sprites
generate_asset \
    "pixel art character sprite, brave knight with sword and shield, 32x32, game asset, transparent background" \
    512 512 \
    "$OUTPUT_DIR/characters/hero.json"

generate_asset \
    "pixel art monster sprite, goblin warrior, 32x32, game asset, transparent background" \
    512 512 \
    "$OUTPUT_DIR/characters/enemy.json"

# Generate backgrounds
generate_asset \
    "pixel art forest background, parallax layer, game background, vibrant colors" \
    1024 512 \
    "$OUTPUT_DIR/backgrounds/forest.json"

generate_asset \
    "pixel art dungeon interior, stone walls and torches, game background" \
    1024 512 \
    "$OUTPUT_DIR/backgrounds/dungeon.json"

# Generate UI elements
generate_asset \
    "game UI health bar, fantasy style, red and gold, transparent background" \
    256 64 \
    "$OUTPUT_DIR/ui/health_bar.json"

generate_asset \
    "game UI button, fantasy style, stone texture, glowing edges" \
    256 128 \
    "$OUTPUT_DIR/ui/button.json"

echo "Asset generation complete!"
```

### 3. Execute Pipeline

```bash
chmod +x generate_assets.sh
./generate_assets.sh
```

## Advanced Pipeline Integration

### Shell Script with YAML Parsing

For more complex pipelines, use a POSIX shell script with `yq` or `jq`:

```sh
#!/bin/sh
# generate_from_spec.sh - Generate assets from YAML specification

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
    
    # Create output directory
    mkdir -p "$OUTPUT_DIR/$category"
    
    # Iterate through assets
    i=0
    while [ "$i" -lt "$count" ]; do
        name=$(yq eval ".${category}[$i].name" "$SPEC_FILE")
        prompt=$(yq eval ".${category}[$i].prompt" "$SPEC_FILE")
        width=$(yq eval ".${category}[$i].width // 512" "$SPEC_FILE")
        height=$(yq eval ".${category}[$i].height // 512" "$SPEC_FILE")
        
        output_file="$OUTPUT_DIR/$category/${name}.json"
        
        echo "Generating $category/$name..."
        asset-generator generate image \
            --prompt "$prompt" \
            --width "$width" \
            --height "$height" \
            --steps 30 \
            --cfg-scale 7.5 \
            --format json \
            --output "$output_file"
        
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
        --output "characters/${name}.json"
}

# Background template
generate_background() {
    name="$1"
    description="$2"
    asset-generator generate image \
        --prompt "pixel art background, ${description}, parallax layer, game background, vibrant colors" \
        --width 1024 --height 512 \
        --output "backgrounds/${name}.json"
}

# UI element template
generate_ui() {
    name="$1"
    element="$2"
    style="$3"
    asset-generator generate image \
        --prompt "game UI ${element}, ${style} style, clean design, transparent background" \
        --width 256 --height 128 \
        --output "ui/${name}.json"
}
```

### 2. Reproducible Generations

Use seeds for consistent results:

```bash
asset-generator generate image \
    --prompt "hero character sprite" \
    --seed 42 \
    --width 512 --height 512
```

### 3. Batch Generation

Generate variations efficiently:

```bash
# Generate 4 enemy variations
asset-generator generate image \
    --prompt "pixel art monster sprite, game asset" \
    --batch 4 \
    --width 512 --height 512
```

### 4. Model Selection

Choose appropriate models for your art style:

```bash
# List available models
asset-generator models list

# Use specific model for pixel art
asset-generator generate image \
    --prompt "pixel art sprite" \
    --model "pixel-art-diffusion"
```

## Pipeline Stages

### Stage 1: Concept Generation

```bash
# Generate concept art variants
asset-generator generate image \
    --prompt "fantasy knight character design, multiple poses" \
    --batch 5 \
    --steps 50 \
    --width 1024 --height 1024
```

### Stage 2: Asset Production

```bash
# Generate final game assets
asset-generator generate image \
    --prompt "pixel art knight sprite, 32x32, transparent" \
    --width 512 --height 512 \
    --steps 30
```

### Stage 3: Variations

```sh
# Generate color variations
for color in red blue green yellow; do
    asset-generator generate image \
        --prompt "knight sprite with ${color} armor" \
        --seed 42 \
        --width 512 --height 512 \
        --output "knight_${color}.json"
done
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
            
            - name: Setup SwarmUI CLI
              run: |
                  wget https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-linux-amd64
                  chmod +x asset-generator-linux-amd64
                  sudo mv asset-generator-linux-amd64 /usr/local/bin/asset-generator
            
            - name: Configure CLI
              run: |
                  asset-generator config set api-url ${{ secrets.SWARMUI_URL }}
                  asset-generator config set api-key ${{ secrets.SWARMUI_KEY }}
            
            - name: Generate Assets
              run: |
                  ./scripts/generate_assets.sh
            
            - name: Upload Assets
              uses: actions/upload-artifact@v3
              with:
                  name: game-assets
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
    - asset-generator config set api-url $SWARMUI_URL
  script:
    - sh ./scripts/generate_assets.sh
  artifacts:
    paths:
      - generated_assets/
  only:
    - main
```

## Troubleshooting

### Connection Issues

```bash
# Verify SwarmUI is accessible
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

## Current Limitations & Future Enhancements

**Current Capabilities:**
- ✅ Text-to-image generation
- ✅ Batch generation
- ✅ JSON/YAML output
- ✅ Seed-based reproducibility
- ✅ Multiple output formats

**Missing Features (See PLAN.md):**
- ❌ Direct PNG file output (currently outputs JSON with image paths)
- ❌ Automatic post-processing (resize, crop, format conversion)
- ❌ Built-in batch processing from YAML specs
- ❌ Image-to-image generation (for iterations)
- ❌ Inpainting/outpainting support
- ❌ Template system for common asset types

For planned feature implementation, see [PLAN.md](PLAN.md).

## Example Output

After running the pipeline, you'll have:

```
generated_assets/
├── characters/
│   ├── hero.json
│   └── enemy.json
├── backgrounds/
│   ├── forest.json
│   └── dungeon.json
└── ui/
    ├── health_bar.json
    └── button.json
```

Each JSON file contains generation metadata and image paths:

```json
{
    "session_id": "abc123",
    "prompt": "pixel art character sprite...",
    "image_paths": ["/path/to/generated/image.png"],
    "parameters": {
        "width": 512,
        "height": 512,
        "steps": 30
    }
}
```

## Next Steps

1. Review generated assets in SwarmUI output directory
2. Integrate with your game engine's asset pipeline
3. Set up automated regeneration on specification changes
4. Implement post-processing for game-ready formats

For more information, see:
- [README.md](README.md) - Complete CLI documentation
- [QUICKSTART.md](QUICKSTART.md) - Getting started guide
- [DEVELOPMENT.md](DEVELOPMENT.md) - Developer documentation