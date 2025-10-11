[🏠 Docs Home](README.md) | [📚 Quick Start](QUICKSTART.md) | [🔧 Commands](COMMANDS.md) | [🔗 Pipeline](PIPELINE.md)

---

# Integration Guide: Adding Asset-Generator to Your Project

This guide shows you how to integrate `asset-generator` into your own project for automated asset creation. Whether you're building a game, card deck, icon set, or any other visual project, this guide provides copy-paste ready patterns for common workflows.

## Table of Contents

- [Quick Start (5 Minutes)](#quick-start-5-minutes)
- [Installation Methods](#installation-methods)
- [Project Structure](#project-structure)
- [Creating Your Pipeline File](#creating-your-pipeline-file)
- [Basic Workflow](#basic-workflow)
- [Build System Integration](#build-system-integration)
- [CI/CD Integration](#cicd-integration)
- [Configuration Management](#configuration-management)
- [Version Control](#version-control)
- [Real-World Examples](#real-world-examples)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Next Steps](#next-steps)

---

## Quick Start (5 Minutes)

Get your first assets generated in under 5 minutes:

```bash
# 1. Install asset-generator (choose one method)
curl -sSL https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-linux-amd64 -o asset-generator
chmod +x asset-generator
sudo mv asset-generator /usr/local/bin/

# 2. Initialize configuration
asset-generator config init
asset-generator config set api-url http://localhost:7801

# 3. Create a minimal pipeline file
cat > my-assets.yaml <<EOF
assets:
  - name: Test Assets
    output_dir: output/test
    assets:
      - name: Hero Character
        prompt: "heroic knight with sword and shield"
      - name: Enemy Monster
        prompt: "scary goblin creature"
EOF

# 4. Generate your assets
asset-generator pipeline --file my-assets.yaml

# Your assets are now in ./output/test/
ls -la output/test/
```

✅ **That's it!** You've just generated your first batch of assets.

---

## Installation Methods

### Method 1: Download Pre-built Binary (Recommended)

```bash
# Linux (amd64)
curl -sSL https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-linux-amd64 \
  -o asset-generator && chmod +x asset-generator && sudo mv asset-generator /usr/local/bin/

# macOS (amd64)
curl -sSL https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-darwin-amd64 \
  -o asset-generator && chmod +x asset-generator && sudo mv asset-generator /usr/local/bin/

# macOS (arm64 / Apple Silicon)
curl -sSL https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-darwin-arm64 \
  -o asset-generator && chmod +x asset-generator && sudo mv asset-generator /usr/local/bin/

# Verify installation
asset-generator --version
```

### Method 2: Build from Source

```bash
# Prerequisites: Go 1.21+
git clone https://github.com/opd-ai/asset-generator.git
cd asset-generator
make install

# Verify
asset-generator --version
```

### Method 3: Project-Local Installation

Keep the binary in your project directory instead of system-wide:

```bash
# Download to project directory
curl -sSL https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-linux-amd64 \
  -o ./tools/asset-generator && chmod +x ./tools/asset-generator

# Use relative path in scripts
./tools/asset-generator pipeline --file assets.yaml
```

> 💡 **Tip**: Project-local installation is useful for ensuring version consistency across team members and CI/CD environments.

---

## Project Structure

### Recommended Directory Layout

```
your-project/
├── assets/                      # Asset pipeline definitions
│   ├── characters.yaml          # Character generation pipeline
│   ├── items.yaml               # Item generation pipeline
│   └── ui-elements.yaml         # UI asset pipeline
├── output/                      # Generated assets (gitignored)
│   ├── characters/
│   ├── items/
│   └── ui/
├── config/                      # Configuration files
│   └── asset-generator.yaml     # Project-specific config
├── scripts/                     # Generation scripts
│   ├── generate-all.sh          # Master generation script
│   ├── generate-characters.sh   # Character-specific generation
│   └── post-process.sh          # Post-processing pipeline
├── .gitignore                   # Version control exclusions
├── Makefile                     # Build automation
└── README.md                    # Project documentation
```

### Create the Structure

```bash
# Create recommended directories
mkdir -p assets output config scripts

# Create basic .gitignore
cat > .gitignore <<EOF
# Generated assets
output/
*.png
*.jpg
*.jpeg
*.svg

# Except examples
!examples/**/*.png
!examples/**/*.svg

# State files (optional - see version control section)
*.state.json

# Local configuration overrides
config/local.yaml
.env
EOF
```

---

## Creating Your Pipeline File

### Basic Pipeline Structure

Create `assets/my-pipeline.yaml`:

```yaml
# Project: My Awesome Project
# Purpose: Generate game character assets
# Last updated: 2025-01-15

assets:
  # Top-level group
  - name: Player Characters
    output_dir: output/characters/players
    seed_offset: 0  # Base seed offset for this group
    
    # Metadata applies to all assets in this group
    metadata:
      style: "fantasy RPG character art, detailed"
      quality: "high quality, professional"
      negative: "blurry, low quality, deformed"
    
    # Individual assets
    assets:
      - id: char_warrior_001
        name: Warrior Class
        prompt: "armored knight with sword and shield"
        filename: warrior.png
      
      - id: char_mage_001
        name: Mage Class
        prompt: "robed wizard casting spell, magical effects"
        filename: mage.png
        # Asset-specific metadata overrides group metadata
        metadata:
          element: "fire"
      
      - id: char_archer_001
        name: Archer Class
        prompt: "elf ranger with bow, forest gear"
        filename: archer.png
```

### Hierarchical Organization

For complex projects with nested structure:

```yaml
assets:
  - name: Characters
    output_dir: output/characters
    seed_offset: 0
    metadata:
      style: "game character art"
    
    # Nested subgroups
    subgroups:
      - name: Heroes
        output_dir: heroes
        seed_offset: 0
        metadata:
          alignment: "good"
        assets:
          - name: Knight
            prompt: "heroic knight in shining armor"
      
      - name: Villains
        output_dir: villains
        seed_offset: 100
        metadata:
          alignment: "evil"
        assets:
          - name: Dark Sorcerer
            prompt: "evil wizard in dark robes"
```

> 📖 **Further Reading**: See [PIPELINE.md](PIPELINE.md) for complete pipeline file format documentation.

### Using Template Variables

Leverage filename templates for organized output:

```yaml
assets:
  - name: Enemy Variants
    output_dir: output/enemies
    # {name} and {id} are automatically available
    filename_template: "enemy_{name}_{id}.png"
    assets:
      - id: "001"
        name: goblin
        prompt: "small green goblin creature"
        # Generates: enemy_goblin_001.png
```

> 📖 **Further Reading**: See [FILENAME_TEMPLATES.md](FILENAME_TEMPLATES.md) for advanced template patterns.

---

## Basic Workflow

### 1. Generate Assets

```bash
# Basic generation
asset-generator pipeline --file assets/characters.yaml

# With custom output directory
asset-generator pipeline --file assets/characters.yaml --output-dir ./game-assets

# Preview before generating (dry run)
asset-generator pipeline --file assets/characters.yaml --dry-run

# Generate with custom seed for reproducibility
asset-generator pipeline --file assets/characters.yaml --base-seed 42
```

### 2. Add Post-Processing

```bash
# Auto-crop whitespace borders
asset-generator pipeline --file assets/characters.yaml --auto-crop

# Downscale to specific width (maintains aspect ratio)
asset-generator pipeline --file assets/characters.yaml --downscale-width 1024

# Combine multiple post-processing steps
asset-generator pipeline --file assets/characters.yaml \
  --auto-crop \
  --downscale-width 1024 \
  --strip-metadata
```

### 3. Convert to SVG (Optional)

```bash
# Convert existing PNG to SVG using shapes
asset-generator convert svg \
  --input output/characters/warrior.png \
  --output output/characters/warrior.svg \
  --mode shapes \
  --num-shapes 100

# Using edge tracing for line art
asset-generator convert svg \
  --input output/characters/archer.png \
  --output output/characters/archer.svg \
  --mode trace
```

> 📖 **Further Reading**: See [SVG_CONVERSION.md](SVG_CONVERSION.md) and [POSTPROCESSING.md](POSTPROCESSING.md) for details.

### 4. Regenerate Specific Assets

```bash
# Use asset IDs to regenerate specific items
asset-generator generate image \
  --prompt "armored knight with sword and shield, fantasy RPG art" \
  --seed 42 \
  --save-images \
  --output-dir output/characters \
  --filename warrior.png
```

> 💡 **Tip**: State files track generation parameters. See [STATE_FILE_SHARING.md](STATE_FILE_SHARING.md) for reproducibility patterns.

---

## Build System Integration

### Makefile Integration

Add asset generation to your project's `Makefile`:

```makefile
.PHONY: assets assets-characters assets-items assets-clean assets-preview

# Configuration
ASSET_GEN := asset-generator
ASSET_DIR := assets
OUTPUT_DIR := output

## Generate all assets
assets: assets-characters assets-items assets-ui
	@echo "✅ All assets generated"

## Generate character assets
assets-characters:
	@echo "🎨 Generating character assets..."
	$(ASSET_GEN) pipeline --file $(ASSET_DIR)/characters.yaml \
		--output-dir $(OUTPUT_DIR) \
		--auto-crop \
		--downscale-width 1024

## Generate item assets
assets-items:
	@echo "🎨 Generating item assets..."
	$(ASSET_GEN) pipeline --file $(ASSET_DIR)/items.yaml \
		--output-dir $(OUTPUT_DIR) \
		--auto-crop \
		--downscale-width 512

## Generate UI assets
assets-ui:
	@echo "🎨 Generating UI assets..."
	$(ASSET_GEN) pipeline --file $(ASSET_DIR)/ui-elements.yaml \
		--output-dir $(OUTPUT_DIR)

## Preview asset generation (dry run)
assets-preview:
	@echo "📋 Preview: Character Assets"
	$(ASSET_GEN) pipeline --file $(ASSET_DIR)/characters.yaml --dry-run
	@echo "\n📋 Preview: Item Assets"
	$(ASSET_GEN) pipeline --file $(ASSET_DIR)/items.yaml --dry-run

## Clean generated assets
assets-clean:
	@echo "🧹 Cleaning generated assets..."
	rm -rf $(OUTPUT_DIR)
	@echo "✅ Assets cleaned"

## Check asset-generator installation
check-asset-gen:
	@command -v $(ASSET_GEN) >/dev/null 2>&1 || \
		(echo "❌ asset-generator not found. Install from: https://github.com/opd-ai/asset-generator" && exit 1)
	@echo "✅ asset-generator found: $$($(ASSET_GEN) --version)"
```

Usage:

```bash
# Generate all assets
make assets

# Generate specific asset types
make assets-characters
make assets-items

# Preview what would be generated
make assets-preview

# Clean and regenerate
make assets-clean assets
```

### Shell Script Integration

Create `scripts/generate-all.sh`:

```bash
#!/bin/bash
set -e

# Configuration
ASSET_GEN="${ASSET_GEN:-asset-generator}"
BASE_SEED="${BASE_SEED:-42}"
OUTPUT_DIR="${OUTPUT_DIR:-./output}"
STEPS="${STEPS:-30}"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🎨 Asset Generation Pipeline${NC}"
echo "================================"
echo "Output Dir: $OUTPUT_DIR"
echo "Base Seed:  $BASE_SEED"
echo "Steps:      $STEPS"
echo ""

# Function to generate assets with error handling
generate_pipeline() {
    local pipeline_file=$1
    local pipeline_name=$2
    
    echo -e "${BLUE}Generating ${pipeline_name}...${NC}"
    
    if $ASSET_GEN pipeline \
        --file "$pipeline_file" \
        --output-dir "$OUTPUT_DIR" \
        --base-seed "$BASE_SEED" \
        --steps "$STEPS" \
        --auto-crop \
        --downscale-width 1024; then
        echo -e "${GREEN}✅ ${pipeline_name} complete${NC}\n"
    else
        echo -e "${RED}❌ ${pipeline_name} failed${NC}\n"
        return 1
    fi
}

# Generate all asset types
generate_pipeline "assets/characters.yaml" "Characters"
generate_pipeline "assets/items.yaml" "Items"
generate_pipeline "assets/environments.yaml" "Environments"

echo -e "${GREEN}🎉 All assets generated successfully!${NC}"
echo "Results in: $OUTPUT_DIR"
```

Make it executable and use:

```bash
chmod +x scripts/generate-all.sh

# Basic usage
./scripts/generate-all.sh

# With custom parameters
BASE_SEED=100 OUTPUT_DIR=./game-assets STEPS=50 ./scripts/generate-all.sh
```

### npm/package.json Integration

For JavaScript/Node.js projects, add to `package.json`:

```json
{
  "scripts": {
    "assets:gen": "asset-generator pipeline --file assets/game-assets.yaml",
    "assets:preview": "asset-generator pipeline --file assets/game-assets.yaml --dry-run",
    "assets:clean": "rm -rf output && mkdir -p output",
    "assets:rebuild": "npm run assets:clean && npm run assets:gen",
    "postinstall": "command -v asset-generator || echo 'Warning: asset-generator not installed'"
  }
}
```

Usage:

```bash
npm run assets:gen      # Generate assets
npm run assets:preview  # Preview generation
npm run assets:rebuild  # Clean and regenerate
```

---

## CI/CD Integration

### GitHub Actions

Create `.github/workflows/assets.yml`:

```yaml
name: Generate Assets

on:
  push:
    branches: [ main ]
    paths:
      - 'assets/**'
      - '.github/workflows/assets.yml'
  workflow_dispatch:  # Manual trigger

jobs:
  generate:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
      
      - name: Install asset-generator
        run: |
          curl -sSL https://github.com/opd-ai/asset-generator/releases/latest/download/asset-generator-linux-amd64 \
            -o asset-generator
          chmod +x asset-generator
          sudo mv asset-generator /usr/local/bin/
      
      - name: Verify installation
        run: asset-generator --version
      
      - name: Configure asset-generator
        run: |
          asset-generator config set api-url ${{ secrets.SWARMUI_API_URL }}
          # Optional: Set API key if required
          # asset-generator config set api-key ${{ secrets.SWARMUI_API_KEY }}
      
      - name: Generate assets
        run: |
          make assets
          # Or direct command:
          # asset-generator pipeline --file assets/game-assets.yaml --output-dir ./output
      
      - name: Upload generated assets
        uses: actions/upload-artifact@v3
        with:
          name: generated-assets
          path: output/
          retention-days: 30
      
      - name: Commit and push generated assets (optional)
        if: github.ref == 'refs/heads/main'
        run: |
          git config user.name "GitHub Actions Bot"
          git config user.email "actions@github.com"
          git add output/
          git diff --quiet && git diff --staged --quiet || (git commit -m "chore: regenerate assets [skip ci]" && git push)
```

### GitLab CI

Create `.gitlab-ci.yml`:

```yaml
stages:
  - generate
  - deploy

variables:
  ASSET_GEN_VERSION: "latest"
  OUTPUT_DIR: "output"

generate_assets:
  stage: generate
  image: ubuntu:22.04
  before_script:
    - apt-get update && apt-get install -y curl
    - curl -sSL https://github.com/opd-ai/asset-generator/releases/${ASSET_GEN_VERSION}/download/asset-generator-linux-amd64 -o /usr/local/bin/asset-generator
    - chmod +x /usr/local/bin/asset-generator
    - asset-generator config set api-url $SWARMUI_API_URL
  script:
    - asset-generator pipeline --file assets/game-assets.yaml --output-dir $OUTPUT_DIR
  artifacts:
    paths:
      - $OUTPUT_DIR/
    expire_in: 1 week
  only:
    - main
    - develop
```

### Jenkins Pipeline

Create `Jenkinsfile`:

```groovy
pipeline {
    agent any
    
    environment {
        ASSET_GEN_VERSION = 'latest'
        OUTPUT_DIR = 'output'
        SWARMUI_API_URL = credentials('swarmui-api-url')
    }
    
    stages {
        stage('Setup') {
            steps {
                sh '''
                    curl -sSL https://github.com/opd-ai/asset-generator/releases/${ASSET_GEN_VERSION}/download/asset-generator-linux-amd64 \
                        -o asset-generator
                    chmod +x asset-generator
                    sudo mv asset-generator /usr/local/bin/
                    asset-generator --version
                '''
            }
        }
        
        stage('Configure') {
            steps {
                sh '''
                    asset-generator config set api-url ${SWARMUI_API_URL}
                '''
            }
        }
        
        stage('Generate Assets') {
            steps {
                sh '''
                    make assets
                '''
            }
        }
        
        stage('Archive') {
            steps {
                archiveArtifacts artifacts: 'output/**/*', fingerprint: true
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
    }
}
```

> ⚠️ **Important**: Store API credentials as secrets, never commit them to version control.

---

## Configuration Management

### Project-Specific Configuration

Create `config/asset-generator.yaml`:

```yaml
# Project-specific asset-generator configuration
# Committed to version control for team consistency

# Default API endpoint (override with environment variable or user config)
api-url: http://localhost:7801

# Default output format
format: table

# Default generation parameters
generate:
  model: stable-diffusion-xl
  steps: 30
  width: 1024
  length: 1024
  cfg-scale: 7.5
  sampler: euler_a
  scheduler: karras
  
  # Project-specific negative prompt
  negative-prompt: "blurry, low quality, deformed, ugly, bad anatomy"

# Pipeline defaults
pipeline:
  auto-crop: true
  downscale-width: 1024
  strip-metadata: true
```

### User-Specific Configuration

Users can override project settings in `~/.asset-generator/config.yaml`:

```yaml
# User-specific overrides (not committed)
# Higher priority than project config

# Personal API endpoint
api-url: http://192.168.1.100:7801

# API authentication
api-key: your-personal-api-key

# Performance preferences
generate:
  steps: 50  # Override project default
```

### Configuration Hierarchy

Configuration is loaded in this order (later values override earlier):

1. **Project config**: `./config/asset-generator.yaml`
2. **User config**: `~/.asset-generator/config.yaml`
3. **Environment variables**: `ASSET_GEN_API_URL`, `ASSET_GEN_API_KEY`, etc.
4. **CLI flags**: `--api-url`, `--steps`, etc.

### Environment Variables

```bash
# Set in your shell profile or CI environment
export ASSET_GEN_API_URL=http://localhost:7801
export ASSET_GEN_API_KEY=your-api-key
export ASSET_GEN_FORMAT=json
export ASSET_GEN_VERBOSE=true

# Use in scripts
asset-generator pipeline --file assets/game.yaml
```

### Project Setup Script

Create `scripts/setup-asset-gen.sh`:

```bash
#!/bin/bash
# Setup script for new team members

set -e

echo "🔧 Setting up asset-generator for this project..."

# Check if asset-generator is installed
if ! command -v asset-generator &> /dev/null; then
    echo "❌ asset-generator not found"
    echo "📥 Install from: https://github.com/opd-ai/asset-generator/releases"
    exit 1
fi

# Initialize user config if it doesn't exist
if [ ! -f ~/.asset-generator/config.yaml ]; then
    echo "📝 Initializing user configuration..."
    asset-generator config init
fi

# Prompt for API URL
read -p "Enter SwarmUI API URL [http://localhost:7801]: " api_url
api_url=${api_url:-http://localhost:7801}
asset-generator config set api-url "$api_url"

# Optional API key
read -p "Enter API key (leave blank if none): " api_key
if [ -n "$api_key" ]; then
    asset-generator config set api-key "$api_key"
fi

echo "✅ Setup complete!"
echo "📋 Next steps:"
echo "   1. Review project pipeline files in ./assets/"
echo "   2. Run 'make assets-preview' to see what will be generated"
echo "   3. Run 'make assets' to generate assets"
```

---

## Version Control

### What to Commit

**DO commit** ✅:
- Pipeline YAML files (`assets/*.yaml`)
- Project configuration (`config/asset-generator.yaml`)
- Generation scripts (`scripts/*.sh`, `Makefile`)
- Example/reference assets (`examples/`, `references/`)
- Documentation (`README.md`, `ASSETS.md`)

**DON'T commit** ❌:
- Generated assets (`output/`, `*.png`, `*.jpg`)
- User-specific config (`~/.asset-generator/`)
- API keys or credentials
- Local configuration overrides

### .gitignore Template

```gitignore
# Generated Assets
output/
generated/
*.png
*.jpg
*.jpeg
*.webp

# Exception: Keep example/reference assets
!examples/**/*.png
!reference/**/*.png

# SVG outputs (optional - decide based on project needs)
output/**/*.svg

# State files (optional - see reproducibility section)
*.state.json
.asset-generator-state/

# User-specific configuration
config/local.yaml
config/user.yaml
.env
.env.local

# API credentials
.api-key
secrets.yaml

# Temporary files
*.tmp
*.swp
*~

# OS files
.DS_Store
Thumbs.db
```

### State File Management

State files track generation parameters for reproducibility. You have options:

**Option 1: Don't commit state files** (default)
```gitignore
*.state.json
```
- Pros: Clean repo, no clutter
- Cons: Team members can't reproduce exact assets

**Option 2: Commit state files for reproducibility**
```gitignore
# Remove or comment out:
# *.state.json
```
- Pros: Perfect reproducibility across team
- Cons: More files in repo

**Option 3: Selective state file commits**
```gitignore
# Ignore by default
*.state.json

# Commit specific important assets
!output/heroes/*.state.json
!output/key-art/*.state.json
```

> 📖 **Further Reading**: See [STATE_FILE_SHARING.md](STATE_FILE_SHARING.md) for reproducibility patterns.

### Git LFS for Large Assets

If you commit generated assets, use Git LFS:

```bash
# Install Git LFS
git lfs install

# Track large image files
git lfs track "*.png"
git lfs track "*.jpg"
git lfs track "*.psd"

# Commit .gitattributes
git add .gitattributes
git commit -m "chore: configure Git LFS for assets"
```

---

## Real-World Examples

### Example 1: Game Sprite Generation

**Project**: 2D platformer game  
**Goal**: Generate character sprites, enemies, and items

```yaml
# assets/game-sprites.yaml
assets:
  # Player character variants
  - name: Player Characters
    output_dir: output/characters/player
    seed_offset: 1000
    metadata:
      style: "pixel art style, 2D game sprite, clean lines"
      view: "side view"
      background: "transparent background"
    assets:
      - id: player_idle
        name: player-idle
        prompt: "character standing idle, neutral pose"
      
      - id: player_run
        name: player-running
        prompt: "character running animation frame, legs moving"
      
      - id: player_jump
        name: player-jumping
        prompt: "character jumping up, arms raised"

  # Enemy types
  - name: Enemies
    output_dir: output/characters/enemies
    seed_offset: 2000
    metadata:
      style: "pixel art style, 2D game sprite, menacing"
      view: "side view"
    assets:
      - name: slime-green
        prompt: "small green slime creature, blob shape"
      
      - name: skeleton-warrior
        prompt: "skeleton with sword, undead enemy"
      
      - name: flying-bat
        prompt: "bat with spread wings, flying enemy"

  # Collectible items
  - name: Items
    output_dir: output/items
    seed_offset: 3000
    metadata:
      style: "pixel art, game item icon, glowing"
      view: "top-down view"
    assets:
      - name: coin-gold
        prompt: "gold coin with shine effect"
      
      - name: health-potion
        prompt: "red health potion bottle"
      
      - name: key-golden
        prompt: "golden key with ornate design"
```

**Build integration**:

```makefile
# Makefile
sprites: sprites-characters sprites-enemies sprites-items

sprites-characters:
	asset-generator pipeline --file assets/game-sprites.yaml \
		--output-dir game/assets \
		--auto-crop \
		--downscale-width 512

sprites-post-process:
	# Convert to game engine format
	cd game/assets && for img in output/characters/**/*.png; do \
		convert $$img -resize 64x64 sprites/$$(basename $$img); \
	done
```

### Example 2: Card Game Deck

**Project**: Trading card game  
**Goal**: Generate 52 playing cards with consistent style

```yaml
# assets/playing-cards.yaml
assets:
  - name: Playing Cards
    output_dir: output/cards
    seed_offset: 5000
    metadata:
      style: "fantasy playing card art, ornate border, detailed illustration"
      quality: "high detail, professional quality"
      negative: "blurry, text, letters, words"
    
    subgroups:
      # Hearts suit
      - name: Hearts
        output_dir: hearts
        seed_offset: 0
        metadata:
          suit: "hearts"
          color: "red"
        assets:
          - id: hearts_ace
            name: ace
            prompt: "ace of hearts, single large heart symbol, elegant"
            filename: "01_ace_of_hearts.png"
          
          - id: hearts_2
            name: two
            prompt: "two of hearts, two heart symbols arranged"
            filename: "02_two_of_hearts.png"
          # ... continue for all ranks

      # Diamonds suit
      - name: Diamonds
        output_dir: diamonds
        seed_offset: 20
        metadata:
          suit: "diamonds"
          color: "red"
        assets:
          - id: diamonds_ace
            name: ace
            prompt: "ace of diamonds, single large diamond symbol"
            filename: "01_ace_of_diamonds.png"
          # ... continue

      # Clubs and Spades...
```

**Generation script**:

```bash
#!/bin/bash
# scripts/generate-deck.sh

OUTPUT_DIR="./deck-output"
PRINT_SIZE_WIDTH=2480  # 300 DPI at 2.5" x 3.5"

# Generate all cards
asset-generator pipeline \
    --file assets/playing-cards.yaml \
    --output-dir "$OUTPUT_DIR" \
    --base-seed 42 \
    --steps 50 \
    --auto-crop \
    --downscale-width $PRINT_SIZE_WIDTH

# Create print sheets (9 cards per sheet)
montage "$OUTPUT_DIR"/hearts/*.png \
    -tile 3x3 \
    -geometry +10+10 \
    "$OUTPUT_DIR"/print-sheet-hearts.png

echo "✅ Deck generation complete!"
echo "📋 52 cards in $OUTPUT_DIR"
```

### Example 3: Icon Set for Web App

**Project**: SaaS application UI  
**Goal**: Consistent icon set for navigation and actions

```yaml
# assets/ui-icons.yaml
assets:
  - name: Navigation Icons
    output_dir: output/icons/navigation
    seed_offset: 7000
    metadata:
      style: "minimal flat icon, simple geometric shapes, modern UI"
      background: "transparent background, white icon on transparent"
      size: "square icon, centered"
    assets:
      - name: home
        prompt: "house icon, simple roof and walls"
        filename: "icon-home.png"
      
      - name: search
        prompt: "magnifying glass icon, circular lens"
        filename: "icon-search.png"
      
      - name: settings
        prompt: "gear icon, mechanical cog wheel"
        filename: "icon-settings.png"
      
      - name: profile
        prompt: "user profile icon, person silhouette"
        filename: "icon-profile.png"

  - name: Action Icons
    output_dir: output/icons/actions
    seed_offset: 7100
    metadata:
      style: "minimal flat icon, simple geometric shapes, modern UI"
      background: "transparent background"
    assets:
      - name: add
        prompt: "plus sign icon, thick lines, centered"
      
      - name: delete
        prompt: "trash can icon, simple waste bin"
      
      - name: edit
        prompt: "pencil icon, diagonal angle"
      
      - name: save
        prompt: "floppy disk icon, retro save symbol"
```

**npm integration**:

```json
{
  "scripts": {
    "icons:generate": "asset-generator pipeline --file assets/ui-icons.yaml --downscale-width 256",
    "icons:optimize": "svgo -f output/icons -o public/icons",
    "icons:build": "npm run icons:generate && npm run icons:optimize"
  }
}
```

### Example 4: Marketing Asset Workflow

**Project**: Marketing campaign  
**Goal**: Generate social media assets in multiple sizes

```yaml
# assets/marketing-campaign.yaml
assets:
  - name: Social Media - Product Launch
    output_dir: output/marketing/product-launch
    seed_offset: 9000
    metadata:
      style: "professional marketing photo, vibrant colors, sharp focus"
      quality: "high resolution, commercial quality"
      branding: "modern tech aesthetic, clean composition"
    assets:
      - id: hero_desktop
        name: hero-desktop
        prompt: "product showcase on desk, laptop with software interface, professional workspace"
        filename: "hero-1920x1080.png"
        width: 1920
        length: 1080
      
      - id: social_square
        name: social-square
        prompt: "product logo centered, bold colors, minimal design"
        filename: "social-1080x1080.png"
        width: 1080
        length: 1080
      
      - id: story_vertical
        name: story-vertical
        prompt: "product in use, vertical composition, dynamic angle"
        filename: "story-1080x1920.png"
        width: 1080
        length: 1920
```

**Post-processing pipeline**:

```bash
#!/bin/bash
# scripts/marketing-post-process.sh

INPUT_DIR="output/marketing/product-launch"
OUTPUT_DIR="output/marketing/optimized"

mkdir -p "$OUTPUT_DIR"

# Process each size for web optimization
for img in "$INPUT_DIR"/*.png; do
    filename=$(basename "$img" .png)
    
    # Create WebP version (better compression)
    cwebp -q 90 "$img" -o "$OUTPUT_DIR/${filename}.webp"
    
    # Create JPEG version (wider compatibility)
    convert "$img" -quality 90 "$OUTPUT_DIR/${filename}.jpg"
    
    echo "✅ Optimized: $filename"
done

echo "🎉 Marketing assets ready for deployment!"
```

---

## Best Practices

### 1. Reproducibility with Seeds

Always use consistent seeds for reproducible results:

```yaml
# In pipeline files
assets:
  - name: Key Art
    seed_offset: 1000  # Fixed base seed
    assets:
      - name: hero-character
        seed: 1000  # Explicit seed for critical assets
```

```bash
# In scripts
SEED=42
asset-generator pipeline --file assets/critical-assets.yaml --base-seed $SEED
```

> 💡 **Tip**: Document seeds in your project README for reproducibility.

### 2. Organizing Outputs

Use consistent directory structures:

```
output/
├── characters/
│   ├── heroes/
│   │   ├── warrior.png
│   │   ├── warrior.state.json
│   │   └── warrior-metadata.json
│   └── enemies/
├── items/
│   ├── weapons/
│   └── consumables/
└── environments/
    ├── dungeons/
    └── towns/
```

Configure in pipeline:

```yaml
assets:
  - name: Characters
    output_dir: output/characters
    subgroups:
      - name: Heroes
        output_dir: heroes  # Relative to parent
      - name: Enemies
        output_dir: enemies
```

### 3. Postprocessing Workflows

Chain postprocessing steps efficiently:

```bash
# Single command with all postprocessing
asset-generator pipeline --file assets/game.yaml \
  --auto-crop \
  --downscale-width 1024 \
  --strip-metadata

# Or create dedicated postprocessing script
# scripts/post-process.sh
for img in output/**/*.png; do
    # Crop whitespace
    asset-generator crop --input "$img" --output "$img"
    
    # Downscale
    asset-generator downscale --input "$img" --width 1024 --output "$img"
    
    # Convert to SVG for web
    asset-generator convert svg --input "$img" \
        --output "${img%.png}.svg" \
        --mode shapes --num-shapes 100
done
```

### 4. Error Handling Strategies

Build robust generation scripts:

```bash
#!/bin/bash
set -e  # Exit on error

# Function with error handling
generate_assets() {
    local pipeline=$1
    local max_retries=3
    local retry_count=0
    
    while [ $retry_count -lt $max_retries ]; do
        if asset-generator pipeline --file "$pipeline"; then
            echo "✅ Success: $pipeline"
            return 0
        else
            retry_count=$((retry_count + 1))
            echo "⚠️  Retry $retry_count/$max_retries for $pipeline"
            sleep 5
        fi
    done
    
    echo "❌ Failed after $max_retries attempts: $pipeline"
    return 1
}

# Use the function
generate_assets "assets/characters.yaml" || exit 1
generate_assets "assets/items.yaml" || exit 1
```

### 5. Pipeline Validation

Preview before generating:

```bash
# Dry run shows what would be generated
asset-generator pipeline --file assets/game.yaml --dry-run

# Validate pipeline file structure
asset-generator pipeline --file assets/game.yaml --validate

# Generate only specific assets
asset-generator pipeline --file assets/game.yaml --filter "heroes/*"
```

### 6. Performance Optimization

```yaml
# Use appropriate image sizes
assets:
  - name: Thumbnails
    width: 512
    length: 512
    steps: 20  # Fewer steps for small assets
  
  - name: Hero Images
    width: 2048
    length: 2048
    steps: 50  # More steps for quality
```

```bash
# Parallel generation (if your SwarmUI supports it)
# Generate multiple pipelines concurrently
asset-generator pipeline --file assets/chars.yaml &
asset-generator pipeline --file assets/items.yaml &
wait
```

### 7. Documentation

Document your asset pipelines:

```yaml
# assets/game-assets.yaml
# Project: Fantasy RPG Game
# Purpose: Generate all game character assets
# Maintainer: team@example.com
# Last Updated: 2025-10-10
# 
# Usage:
#   make assets                    # Generate all
#   make assets-characters         # Characters only
#   
# Reproducibility:
#   Base seed: 42 (do not change)
#   Model: stable-diffusion-xl
#   
# Notes:
#   - Character heights normalized to 512px
#   - All assets auto-cropped to remove whitespace
#   - Metadata stripped for privacy

assets:
  # ... your pipeline
```

---

## Troubleshooting

### Issue: "Connection refused" or API errors

**Cause**: SwarmUI not running or wrong URL

**Solution**:
```bash
# Check if SwarmUI is running
curl http://localhost:7801/api/status

# Verify configuration
asset-generator config get api-url

# Test with explicit URL
asset-generator pipeline --file assets/test.yaml --api-url http://localhost:7801
```

### Issue: Generated assets don't match expectations

**Cause**: Prompt quality, model choice, or parameters

**Solution**:
```bash
# Use dry run to preview prompts
asset-generator pipeline --file assets/test.yaml --dry-run

# Experiment with single generation
asset-generator generate image \
  --prompt "your test prompt" \
  --steps 30 \
  --cfg-scale 7.5 \
  --save-images

# Check available models
asset-generator models list
```

### Issue: Pipeline file syntax errors

**Cause**: YAML formatting issues

**Solution**:
```bash
# Validate YAML syntax
yamllint assets/pipeline.yaml

# Check structure with dry run
asset-generator pipeline --file assets/pipeline.yaml --dry-run --verbose

# Start with minimal example
cat > test.yaml <<EOF
assets:
  - name: Test
    output_dir: output/test
    assets:
      - name: test-image
        prompt: "simple test prompt"
EOF
```

### Issue: Slow generation performance

**Cause**: Hardware limitations, too many steps, large images

**Solution**:
```yaml
# Optimize parameters
assets:
  - name: Fast Assets
    width: 512    # Smaller size
    length: 512
    steps: 20     # Fewer steps
    
    # Use faster scheduler
    scheduler: simple  # Instead of karras
```

### Issue: Out of disk space

**Cause**: Generated assets accumulating

**Solution**:
```bash
# Clean old outputs
make assets-clean

# Or selective cleanup
find output/ -type f -mtime +30 -delete  # Delete files older than 30 days

# Compress old assets
tar -czf archive-$(date +%Y%m%d).tar.gz output/
rm -rf output/
```

### Issue: Inconsistent results across team

**Cause**: Different configurations or model versions

**Solution**:
```bash
# Commit project configuration
git add config/asset-generator.yaml

# Document setup in README
cat >> README.md <<EOF
## Asset Generation Setup
1. Install asset-generator: https://github.com/opd-ai/asset-generator
2. Run: ./scripts/setup-asset-gen.sh
3. Generate: make assets
EOF

# Use explicit parameters
asset-generator pipeline \
  --file assets/game.yaml \
  --model stable-diffusion-xl \
  --base-seed 42 \
  --steps 30
```

### Issue: State file corruption

**Cause**: Interrupted generation or disk issues

**Solution**:
```bash
# Remove corrupt state files
find output/ -name "*.state.json" -size 0 -delete

# Regenerate with fresh state
rm output/**/*.state.json
asset-generator pipeline --file assets/game.yaml
```

---

## Next Steps

Now that you have asset-generator integrated into your project, explore advanced features:

### Essential Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Complete CLI usage guide
- **[PIPELINE.md](PIPELINE.md)** - Advanced pipeline patterns and metadata cascading
- **[COMMANDS.md](COMMANDS.md)** - Full command reference
- **[GENERATION_FEATURES.md](GENERATION_FEATURES.md)** - LoRA, Skimmed CFG, and advanced parameters

### Advanced Topics

- **[FILENAME_TEMPLATES.md](FILENAME_TEMPLATES.md)** - Dynamic naming with variables and functions
- **[POSTPROCESSING.md](POSTPROCESSING.md)** - Crop, downscale, and metadata management
- **[SVG_CONVERSION.md](SVG_CONVERSION.md)** - Convert raster to vector graphics
- **[STATE_FILE_SHARING.md](STATE_FILE_SHARING.md)** - Reproducibility and collaboration patterns
- **[LORA_SUPPORT.md](LORA_SUPPORT.md)** - Style adaptation with Low-Rank models

### Examples

- **[examples/generic-pipeline.yaml](../examples/generic-pipeline.yaml)** - Reference pipeline structure
- **[examples/tarot-deck/](../examples/tarot-deck/)** - Complete 78-card deck generation
- **[examples/tarot-deck/DEMONSTRATION.md](../examples/tarot-deck/DEMONSTRATION.md)** - Step-by-step walkthrough

### Community

- **GitHub Issues**: Report bugs or request features
- **GitHub Discussions**: Ask questions and share your projects
- **Contributing**: See [DEVELOPMENT.md](DEVELOPMENT.md) for contribution guidelines

---

## Summary

You now have everything needed to integrate asset-generator into your project:

✅ Installation methods (binary, source, project-local)  
✅ Project structure and configuration  
✅ Pipeline file creation  
✅ Build system integration (Makefile, npm, scripts)  
✅ CI/CD patterns (GitHub Actions, GitLab CI, Jenkins)  
✅ Version control best practices  
✅ Real-world examples (games, cards, icons, marketing)  
✅ Troubleshooting guide  

**Quick Reference Commands**:

```bash
# Generate assets
asset-generator pipeline --file assets/my-project.yaml

# Preview generation
asset-generator pipeline --file assets/my-project.yaml --dry-run

# With postprocessing
asset-generator pipeline --file assets/my-project.yaml \
  --auto-crop --downscale-width 1024

# Reproducible generation
asset-generator pipeline --file assets/my-project.yaml --base-seed 42
```

Happy generating! 🎨
