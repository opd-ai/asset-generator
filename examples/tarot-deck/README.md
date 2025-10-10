# Complete Tarot Deck Generation Pipeline

This example demonstrates the full capabilities of `asset-generator` by creating a complete 78-card tarot deck with professional-quality assets, organized structure, and multiple output formats.

## Overview

This pipeline generates:
- **78 unique tarot cards** (22 Major Arcana + 56 Minor Arcana)
- **4 card back design variations**
- **Multiple output formats** (print-ready, web, mobile, SVG, thumbnails)
- **Organized directory structure** for easy management
- **Reproducible generation** using seed values

## Features Demonstrated

✅ **Batch Generation**: Automated creation of 78+ unique assets  
✅ **Custom Filename Templates**: Organized naming with prefixes and sorting  
✅ **Seed-Based Reproducibility**: Consistent regeneration of specific cards  
✅ **Directory Organization**: Multi-level folder structure  
✅ **Negative Prompts**: Quality control and style consistency  
✅ **Post-Processing Pipeline**: Downscaling, SVG conversion, thumbnails  
✅ **YAML-Based Specifications**: Structured asset definitions  
✅ **Shell Script Automation**: Complete end-to-end pipeline  

## Quick Start

### Prerequisites

```bash
# Install asset-generator CLI
# (See main README for installation instructions)

# Configure asset-generator
asset-generator config set api-url http://localhost:7801
```

### Generate Complete Deck

**Option 1: Using the Built-in Pipeline Command (Recommended)**

The native pipeline command eliminates the need for external scripts and yq:

```bash
cd examples/tarot-deck

# Generate all 78 cards
asset-generator pipeline --file tarot-spec.yaml --output-dir ./tarot-deck-output

# Preview before generating (dry run)
asset-generator pipeline --file tarot-spec.yaml --dry-run

# With custom parameters and postprocessing
asset-generator pipeline --file tarot-spec.yaml \
  --output-dir ./my-deck \
  --base-seed 42 \
  --steps 50 \
  --style-suffix "detailed illustration, ornate border, rich colors" \
  --auto-crop \
  --downscale-width 1024
```

**Option 2: Using the Shell Script Wrapper**

The `generate-tarot-deck.sh` script is now a convenient wrapper around the pipeline command:

```bash
cd examples/tarot-deck

# Make script executable
chmod +x generate-tarot-deck.sh

# Generate with default settings
./generate-tarot-deck.sh

# Custom output directory and seed
./generate-tarot-deck.sh ./my-deck 42

# Pass additional flags to pipeline command
./generate-tarot-deck.sh ./my-deck 42 --dry-run
```

**Benefits of the wrapper approach:**
- ✅ No external dependencies (yq no longer required)
- ✅ Native YAML parsing in Go
- ✅ Backward compatible interface
- ✅ 59% less code (101 vs 246 lines)
- ✅ Better error handling and progress tracking
- ✅ Consistent with other asset-generator commands

### Traditional Workflow

```bash
# Generate all 78 cards (takes 30-60 minutes depending on your hardware)
./generate-tarot-deck.sh

# Generate card back designs
./generate-card-back.sh

# Post-process for multiple formats
./post-process-deck.sh
```

## Pipeline Structure

```
examples/tarot-deck/
├── tarot-spec.yaml              # Complete deck specification (YAML)
├── generate-tarot-deck.sh       # Pipeline wrapper script
├── generate-card-back.sh        # Card back generator
├── post-process-deck.sh         # Multi-format conversion
├── quick-demo.sh                # Quick demo (5 sample cards)
├── README.md                    # This file
└── tarot-deck-output/           # Generated assets
    ├── major-arcana/            # 22 Major Arcana cards
    ├── minor-arcana/            # 56 Minor Arcana cards
    │   ├── wands/               # 14 cards
    │   ├── cups/                # 14 cards
    │   ├── swords/              # 14 cards
    │   └── pentacles/           # 14 cards
    └── card-backs/              # Card back designs
```

## Card Specifications

### Major Arcana (22 Cards)
Traditional Major Arcana from 0 (The Fool) to 21 (The World), each with detailed symbolic prompts based on classic tarot imagery.

**Examples:**
- `00-the_fool.png` - Young traveler at cliff edge with white rose
- `01-the_magician.png` - Figure with infinity symbol and four suit symbols
- `21-the_world.png` - Dancer in laurel wreath with four creatures

### Minor Arcana (56 Cards)
Four suits with 14 cards each (Ace through King):

**Suits:**
- **Wands** (Fire element) - Action, inspiration, ambition
- **Cups** (Water element) - Emotions, relationships, feelings
- **Swords** (Air element) - Thoughts, communication, conflict
- **Pentacles** (Earth element) - Material, practical, physical

**Ranks:** Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Page, Knight, Queen, King

## Generation Parameters

### Image Dimensions
- **Width:** 768px
- **Height:** 1344px
- **Aspect Ratio:** 4:7 (standard tarot proportion)
- **Print Size:** ~2.75" x 4.75" at 300 DPI equivalent

### Generation Settings
- **Steps:** 40 (high quality)
- **CFG Scale:** 7.5 (balanced prompt adherence)
- **Style:** Traditional tarot art with ornate borders
- **Negative Prompt:** Removes blurry, distorted, modern elements

### Seed Strategy
- **Major Arcana:** `BASE_SEED + card_number` (42-63)
- **Minor Arcana:** `BASE_SEED + 100 + offset` (142+)
- **Card Backs:** `BASE_SEED + 9999` (10041+)

## Advanced Usage

### Custom Base Seed

```bash
# Generate with different seed for variation
./generate-tarot-deck.sh ./custom-output 1234
```

### Regenerate Specific Card

```bash
# Regenerate just The Fool (card 0) with seed 42
asset-generator generate image \
  --prompt "tarot card art, The Fool, young traveler at cliff edge..." \
  --width 768 --height 1344 \
  --steps 40 --cfg-scale 7.5 \
  --seed 42 \
  --save-images \
  --output-dir ./tarot-deck-output/major-arcana \
  --filename-template "00-the_fool.png"
```

### Generate Single Suit

```bash
# Use pipeline command with a custom YAML containing just one suit
asset-generator pipeline --file wands-only.yaml --output-dir ./wands-output

# Or use generate command for manual single-card generation
asset-generator generate image \
  --prompt "Ace of Wands, flaming wand..." \
  --seed 142 \
  --output-dir ./single-cards
```

### Custom Styling

Modify styling via pipeline command flags or edit the wrapper script:

```bash
# Pass custom style via pipeline command
asset-generator pipeline --file tarot-spec.yaml \
  --style-suffix "watercolor painting, soft edges, artistic, hand-painted look"

# Or for modern minimalist style
asset-generator pipeline --file tarot-spec.yaml \
  --style-suffix "minimalist design, clean lines, modern interpretation, flat colors"

# Or edit generate-tarot-deck.sh to change default STYLE_SUFFIX variable:
STYLE_SUFFIX="vintage woodcut print, engraving style, black and white, high contrast"
```

## Post-Processing Options

The `post-process-deck.sh` script creates multiple output formats:

### 1. Print-Ready (Original Resolution)
- **Size:** 768x1344px
- **Use:** Professional printing, high-quality reproduction
- **Location:** `tarot-deck-processed/print-ready/`

### 2. Web-Optimized
- **Size:** Max 1024px (maintains aspect ratio)
- **Use:** Website galleries, online stores
- **Location:** `tarot-deck-processed/web-optimized/`

### 3. Mobile-Optimized
- **Size:** Max 512px (maintains aspect ratio)
- **Use:** Mobile apps, quick loading
- **Location:** `tarot-deck-processed/mobile-optimized/`

### 4. SVG Versions (Sample)
- **Format:** Scalable Vector Graphics
- **Use:** Infinite scaling, design work, laser cutting
- **Location:** `tarot-deck-processed/svg-versions/`
- **Note:** Sample conversion (3 cards) due to processing time

### 5. Thumbnails
- **Size:** Max 256px (maintains aspect ratio)
- **Use:** Preview galleries, navigation, selection UI
- **Location:** `tarot-deck-processed/thumbnails/`

## Quick Demo Mode

Test the pipeline without generating all 78 cards:

```bash
# Generate just 5 sample cards for testing
./quick-demo.sh
```

This generates:
- The Fool (Major Arcana 0)
- The Magician (Major Arcana 1)
- Ace of Wands
- Ace of Cups
- Ace of Swords

## Customization Examples

### Different Card Dimensions

For poker-sized cards (2.5" x 3.5"):

```bash
# Edit generate-tarot-deck.sh
WIDTH=750
HEIGHT=1050
```

### Different Art Styles

#### Pixel Art Tarot
```bash
STYLE_SUFFIX="pixel art style, 8-bit aesthetic, retro gaming, crisp pixels, limited color palette"
```

#### Watercolor Tarot
```bash
STYLE_SUFFIX="watercolor painting, soft gradients, artistic brushstrokes, flowing colors, hand-painted"
```

#### Minimalist Tarot
```bash
STYLE_SUFFIX="minimalist design, clean lines, simple shapes, limited colors, modern interpretation"
```

#### Dark Gothic Tarot
```bash
STYLE_SUFFIX="dark gothic art, dramatic shadows, intricate details, Victorian gothic, dark fantasy"
```

### Adding Custom Cards

Edit `tarot-spec.yaml` to add custom cards:

```yaml
custom_cards:
  - name: The Phoenix
    prompt: "tarot card art, mythical phoenix rising from flames, rebirth and transformation, golden fire, detailed feathers"
  
  - name: The Dragon
    prompt: "tarot card art, majestic dragon coiled around tower, wisdom and power, scales and wings, mystical energy"
```

## Performance & Timing

**Estimated Generation Times** (on typical GPU):
- Single card: ~30-60 seconds
- Full Major Arcana (22 cards): ~20-30 minutes
- Full Minor Arcana (56 cards): ~45-60 minutes
- Complete deck (78 cards): ~60-90 minutes
- Post-processing: ~5-10 minutes

**Optimization Tips:**
- Reduce `--steps` to 25 for faster generation (slight quality loss)
- Use `--batch` flag for parallel generation (if API supports)
- Run on dedicated GPU for faster processing

## Troubleshooting

### Connection Issues

```bash
# Verify asset generation service is running
curl http://localhost:7801/api/health

# Check asset-generator configuration
asset-generator config view

# Test generation with single card
asset-generator generate image --prompt "test card" --verbose
```

### Quality Issues

**Problem:** Cards look blurry or low quality  
**Solution:** Increase steps to 50, adjust CFG scale to 8-9

**Problem:** Prompt not being followed accurately  
**Solution:** Adjust CFG scale up to 9-10, make prompt more specific

**Problem:** Unwanted text appears on cards  
**Solution:** Add to negative prompt: "text, words, letters, numbers, watermark"

### Script Errors

**Problem:** Permission denied on scripts  
**Solution:** `chmod +x *.sh`

**Problem:** `asset-generator: command not found`  
**Solution:** Install asset-generator and ensure it's in your PATH

**Problem:** Asset generator not found  
**Solution:** Ensure `asset-generator` is in PATH or use full path

## Integration Examples

### Game Development

```bash
# Generate and export for Unity
./generate-tarot-deck.sh ./Assets/Resources/Cards
./post-process-deck.sh ./Assets/Resources/Cards ./Assets/Resources/Cards-Optimized

# Use mobile-optimized version in game
cp -r tarot-deck-processed/mobile-optimized/* ./Assets/Resources/Cards-Game/
```

### Web Application

```bash
# Generate for web app
./generate-tarot-deck.sh
./post-process-deck.sh

# Copy web-optimized to public directory
cp -r tarot-deck-processed/web-optimized/* ./public/assets/cards/
cp -r tarot-deck-processed/thumbnails/* ./public/assets/cards/thumbs/
```

### Print Production

```bash
# Generate high-resolution for printing
./generate-tarot-deck.sh
./generate-card-back.sh

# Use print-ready directory
cd tarot-deck-processed/print-ready/

# Package for print shop
zip -r tarot-deck-printable.zip ./*
```

## File Naming Convention

Cards are named with consistent, sortable filenames:

**Major Arcana:**
```
00-the_fool.png
01-the_magician.png
02-the_high_priestess.png
...
21-the_world.png
```

**Minor Arcana:**
```
01-ace_of_wands.png
02-two_of_wands.png
...
14-king_of_wands.png
```

This ensures:
- ✅ Alphabetical sorting matches card order
- ✅ Easy to identify specific cards
- ✅ Compatible with most file systems
- ✅ Searchable and filterable

## License & Attribution

This example pipeline and specifications are provided as-is for demonstration purposes. Generated card artwork may be subject to the terms of your asset generation service and model licenses.

**Traditional tarot symbolism** is public domain, but specific artistic interpretations may have copyright protections depending on your generation model and service terms.

## Next Steps

1. **Generate your deck**: Run the complete pipeline
2. **Customize styling**: Modify prompts for your preferred aesthetic
3. **Create variations**: Use different seeds for alternate designs
4. **Print or publish**: Use post-processed formats for your medium
5. **Integrate**: Import into your game, app, or website

## Related Documentation

- [../GENERATE_PIPELINE.md](../../GENERATE_PIPELINE.md) - General pipeline documentation
- [../docs/FILENAME_TEMPLATES.md](../../docs/FILENAME_TEMPLATES.md) - Filename template guide
- [../docs/IMAGE_DOWNLOAD.md](../../docs/IMAGE_DOWNLOAD.md) - Image download features
- [../README.md](../../README.md) - Main CLI documentation

## Questions & Support

For issues with the pipeline:
1. Check troubleshooting section above
2. Verify asset-generator configuration
3. Test with quick-demo.sh first
4. Review generated assets for quality

---

**Happy deck creation!** 🎴✨
