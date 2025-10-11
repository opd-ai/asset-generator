# LoRA Support in Asset Generator

## Overview

The Asset Generator CLI now supports **LoRA (Low-Rank Adaptation)** models for enhanced image generation. LoRAs allow you to fine-tune and customize the style, content, and characteristics of generated images without changing the base model.

## What are LoRAs?

LoRAs are lightweight model adaptations that modify specific aspects of image generation:
- **Style modifications** (anime, photorealistic, artistic styles)
- **Character consistency** (specific characters, faces)
- **Content enhancement** (details, quality improvements)
- **Concept injection** (specific objects, themes, aesthetics)

LoRAs are typically much smaller than full models (5-200MB vs 2-7GB) and can be combined for powerful customization.

## Basic Usage

### Simple LoRA Application

Apply a single LoRA with default weight (1.0):

```bash
asset-generator generate image \
  --prompt "anime character in forest" \
  --lora "anime-style-v1"
```

### LoRA with Custom Weight

Specify weight inline using colon notation (format: `name:weight`):

```bash
asset-generator generate image \
  --prompt "anime character in forest" \
  --lora "anime-style-v1:0.8"
```

Weight values typically range from 0.0 to 2.0:
- **0.0**: No effect
- **0.5-0.8**: Subtle influence
- **1.0**: Standard strength (default)
- **1.2-2.0**: Strong influence

### Multiple LoRAs

Combine multiple LoRAs for complex styles:

```bash
asset-generator generate image \
  --prompt "cyberpunk cityscape at night" \
  --lora "cyberpunk-aesthetic:1.0" \
  --lora "neon-lights:0.7" \
  --lora "detailed-architecture:0.5"
```

## Advanced Usage

### Negative Weights

Use negative weights to **remove** or **reduce** unwanted styles:

```bash
asset-generator generate image \
  --prompt "portrait of a person" \
  --lora "realistic-faces:1.0" \
  --lora "cartoon-style:-0.5"  # Remove cartoon-like features
```

### Configuration File

Define default LoRAs in your config file (`~/.config/asset-generator/config.yaml`):

```yaml
generate:
  model: stable-diffusion-xl
  steps: 30
  cfg-scale: 7.5
  
  # Default LoRAs applied to all generations
  loras:
    - anime-style-v1:0.8
    - detailed-faces:0.6
  
  # Default weight for LoRAs without explicit weight
  lora-default-weight: 1.0
```

Override config LoRAs with command-line flags:

```bash
# This overrides config file LoRAs
asset-generator generate image \
  --prompt "portrait" \
  --lora "photorealistic:1.0"
```

### Alternative Syntax: Explicit Weights

Use separate flags for names and weights:

```bash
asset-generator generate image \
  --prompt "fantasy landscape" \
  --lora "fantasy-art" \
  --lora "detailed-background" \
  --lora-weight 0.9 \
  --lora-weight 0.6
```

Weights are applied in order to the corresponding LoRA.

### Custom Default Weight

Change the default weight for all LoRAs without explicit weights:

```bash
asset-generator generate image \
  --prompt "character design" \
  --lora "anime-style" \
  --lora "detailed-clothing:1.2" \
  --lora-default-weight 0.7  # "anime-style" uses 0.7
```

## Practical Examples

### Example 1: Anime Character Generation

```bash
asset-generator generate image \
  --prompt "anime girl with blue hair, school uniform" \
  --model "stable-diffusion-xl" \
  --lora "anime-style-v2:0.9" \
  --lora "detailed-eyes:0.7" \
  --batch 4 \
  --save-images
```

### Example 2: Photorealistic Portrait

```bash
asset-generator generate image \
  --prompt "portrait of a young woman, natural lighting" \
  --lora "realistic-faces:1.0" \
  --lora "professional-photography:0.8" \
  --width 768 --height 1024 \
  --save-images
```

### Example 3: Cyberpunk Scene

```bash
asset-generator generate image \
  --prompt "futuristic city street, rain, neon signs" \
  --lora "cyberpunk-aesthetic:1.1" \
  --lora "neon-lights:0.8" \
  --lora "cinematic-lighting:0.6" \
  --cfg-scale 8.0 \
  --steps 40
```

### Example 4: Mixed Style with Negative LoRA

```bash
asset-generator generate image \
  --prompt "fantasy warrior character" \
  --lora "fantasy-art:1.0" \
  --lora "detailed-armor:0.7" \
  --lora "cartoon-style:-0.3"  # Reduce cartoon influence
```

### Example 5: Batch Generation with LoRAs

```bash
asset-generator generate image \
  --prompt "fantasy landscape, magical forest" \
  --lora "fantasy-art:0.9" \
  --lora "detailed-nature:0.6" \
  --batch 8 \
  --save-images --output-dir ./fantasy-landscapes \
  --filename-template "forest-{index}-{seed}.png"
```

## LoRA Discovery

### List Available LoRAs

The current implementation of `models list` lists all available models including LoRAs. To filter for LoRAs specifically, use grep or other text filtering tools:

```bash
# List all models (including LoRAs)
asset-generator models list

# Filter for LoRA models using grep
asset-generator models list | grep -i lora

# Use JSON output for more precise filtering
asset-generator models list --format json | jq '.[] | select(.type | contains("LoRA"))'
```

**Note**: The SwarmUI API supports `ListModels` with a `subtype` parameter for filtering (e.g., "LoRA", "Wildcards"), but this is not yet exposed in the CLI. A future enhancement could add `--subtype` flag to the `models list` command.

### Check LoRA Compatibility

Not all LoRAs work with all base models:
- **SD 1.5 LoRAs**: Work with Stable Diffusion 1.5 models
- **SDXL LoRAs**: Work with Stable Diffusion XL models
- **Pony Diffusion LoRAs**: Work with Pony-based models

Check the LoRA documentation or filename for compatibility information.

## Troubleshooting

### LoRA Not Taking Effect

1. **Verify LoRA name**: Ensure the LoRA name matches exactly (case-sensitive)
2. **Check compatibility**: Verify LoRA is compatible with your base model
3. **Adjust weight**: Try increasing the weight (0.8 → 1.2)
4. **Prompt alignment**: Ensure your prompt aligns with the LoRA's purpose

### Unexpected Results

1. **Lower weights**: Strong LoRAs may overpower at 1.0 - try 0.6-0.8
2. **Reduce LoRA count**: Too many LoRAs can conflict - use 2-3 maximum
3. **Check LoRA order**: Earlier LoRAs have slightly more influence

### Weight Out of Range Error

The CLI validates weights between -2.0 and 5.0:
- Normal range: 0.0 to 2.0
- Extended range: -2.0 to 5.0 (for advanced use cases)

If you need weights outside this range, the LoRA may be incompatible or incorrectly trained.

## Integration with Other Features

### LoRAs + Skimmed CFG

Combine LoRAs with Skimmed CFG for faster, high-quality generation:

```bash
asset-generator generate image \
  --prompt "detailed character portrait" \
  --lora "detailed-faces:0.8" \
  --skimmed-cfg --skimmed-cfg-scale 3.0
```

### LoRAs + Pipeline Processing

Use LoRAs with the complete postprocessing pipeline:

```bash
asset-generator generate image \
  --prompt "high-resolution artwork" \
  --lora "detailed-art:0.9" \
  --width 2048 --height 2048 \
  --save-images \
  --auto-crop \
  --downscale-width 1024
```

### LoRAs in Batch Scripts

Incorporate LoRAs in automated workflows:

```bash
#!/bin/bash
# Generate character variations with different LoRAs

for style in "anime:0.9" "realistic:1.0" "cartoon:0.8"; do
  asset-generator generate image \
    --prompt "character design, fantasy warrior" \
    --lora "$style" \
    --batch 4 \
    --save-images --output-dir "./output-${style%%:*}"
done
```

## API Implementation Details

### SwarmUI LoRA Format

LoRAs are sent to the SwarmUI API in the following format:

```json
{
  "prompt": "your prompt here",
  "loras": {
    "anime-style-v1": 0.8,
    "detailed-faces": 0.6
  }
}
```

The CLI automatically converts the `--lora` flags into this format.

### LoRA Parameter Precedence

1. **Inline weights** (highest priority): `--lora "name:0.8"`
2. **Explicit weights**: `--lora-weight 0.8`
3. **Default weight**: `--lora-default-weight 1.0`
4. **Built-in default**: 1.0

## Best Practices

1. **Start with one LoRA**: Test each LoRA individually before combining
2. **Use moderate weights**: Start with 0.7-0.9 for most LoRAs
3. **Limit combinations**: 2-3 LoRAs work best; more can cause conflicts
4. **Match prompts to LoRAs**: Align your prompt with the LoRA's training domain
5. **Test weight ranges**: Small weight adjustments (±0.1) can make big differences
6. **Document your setups**: Save working LoRA combinations in config files

---

## Quick Reference

### Common Commands

```bash
# Single LoRA with default weight
asset-generator generate image --prompt "anime character" --lora "anime-style"

# Single LoRA with custom weight
asset-generator generate image --prompt "portrait" --lora "realistic-faces:0.8"

# Multiple LoRAs
asset-generator generate image --prompt "cyberpunk city" \
  --lora "cyberpunk:1.0" --lora "neon-lights:0.7"

# Negative weight (remove style)
asset-generator generate image --prompt "character" \
  --lora "realistic:1.0" --lora "cartoon:-0.5"
```

### Weight Guidelines

| Weight | Effect |
|--------|--------|
| `-0.5 to 0.0` | Removes/reduces style elements |
| `0.0` | No effect |
| `0.5-0.7` | Subtle influence |
| `0.8-1.0` | Standard strength |
| `1.1-1.5` | Strong influence |
| `1.6-2.0` | Very strong (may overpower) |

### Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| LoRA not taking effect | Check name (case-sensitive), increase weight |
| Too strong/overpowering | Reduce weight (try 0.6-0.8) |
| Unexpected results | Use fewer LoRAs (2-3 max) |
| Conflicts between LoRAs | Reorder or reduce weights |

### Common Use Cases

**Anime Style:**
```bash
--lora "anime-style-v2:0.9" --lora "detailed-eyes:0.7"
```

**Photorealistic:**
```bash
--lora "realistic-faces:1.0" --lora "professional-photo:0.8"
```

**Fantasy Art:**
```bash
--lora "fantasy-art:0.9" --lora "detailed-background:0.6"
```

## See Also

- [SwarmUI Documentation](https://github.com/mcmonkeyprojects/SwarmUI) - SwarmUI API reference
- [QUICKSTART.md](QUICKSTART.md) - Basic usage guide
- [SKIMMED_CFG.md](SKIMMED_CFG.md) - Advanced sampling techniques
- [PIPELINE.md](PIPELINE.md) - Complete postprocessing pipeline
