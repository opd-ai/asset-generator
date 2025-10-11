[🏠 Docs Home](README.md) | [📚 Quick Start](QUICKSTART.md) | [🔧 Commands](COMMANDS.md) | [🔗 Pipeline](PIPELINE.md)

---

# Asset Generator User Guide

> **Complete guide to generation features, LoRA support, and image download with filename templates**

This comprehensive user guide covers all the advanced features for generating, customizing, and organizing your AI-generated assets.

## Table of Contents

- [Generation Features](#generation-features)
  - [Scheduler Selection](#scheduler-selection)
  - [Skimmed CFG](#skimmed-cfg)
- [LoRA Support](#lora-support)
  - [Basic Usage](#lora-basic-usage)
  - [Advanced LoRA Techniques](#advanced-lora-techniques)
  - [LoRA Best Practices](#lora-best-practices)
- [Image Download & Organization](#image-download-organization)
  - [Filename Templates](#filename-templates)
  - [Template Placeholders](#template-placeholders)
  - [Organization Examples](#organization-examples)

---

# Generation Features {#generation-features}

Advanced generation-time parameters that control quality, speed, and output characteristics.

## Scheduler Selection {#scheduler-selection}

### Overview

The scheduler (also known as noise schedule) controls how noise is added and removed during the diffusion process. Different schedulers can significantly impact generation quality, speed, and characteristics.

**Status**: ✅ **IMPLEMENTED**

The Asset Generator CLI supports selecting different schedulers via:
- `--scheduler` flag for `generate image` command
- `--scheduler` flag for `pipeline` command  
- `scheduler` config file setting
- Defaults to `simple` scheduler

### Available Schedulers

| Scheduler | Description | Best For |
|-----------|-------------|----------|
| `simple` | Basic linear schedule (default) | General purpose, fast |
| `normal` | Standard scheduler with balanced noise | Most use cases, reliable |
| `karras` | Karras noise schedule | High-quality, detailed images |
| `exponential` | Exponential decay schedule | Smooth transitions, artistic |
| `sgm_uniform` | Uniform noise distribution | Experimental, specialized models |

### Scheduler Usage Examples

#### Basic Usage

```bash
# Default scheduler (simple)
asset-generator generate image \
  --prompt "portrait of a wizard"

# Karras scheduler for high quality
asset-generator generate image \
  --prompt "detailed portrait of a wizard" \
  --scheduler karras \
  --steps 30

# Normal scheduler (standard)
asset-generator generate image \
  --prompt "landscape photograph" \
  --scheduler normal
```

#### Pipeline with Scheduler

```bash
# Apply scheduler to all assets in pipeline
asset-generator pipeline \
  --file assets-spec.yaml \
  --output-dir ./assets \
  --scheduler karras \
  --steps 40
```

#### Configuration File

Set default scheduler in `~/.asset-generator/config.yaml`:

```yaml
generate:
  model: stable-diffusion-xl
  steps: 25
  scheduler: karras  # Use Karras scheduler by default
```

### Scheduler Selection Guide

#### Simple (Default)
- **Best for**: Quick iterations, testing, general purpose
- **Speed**: Fastest (1.0x)
- **Steps**: 15-20

```bash
asset-generator generate image \
  --prompt "concept art, character design" \
  --scheduler simple --steps 20
```

#### Normal
- **Best for**: Production-quality assets, balanced workflows
- **Speed**: Fast (1.1x)
- **Steps**: 20-30

```bash
asset-generator generate image \
  --prompt "product photography, professional lighting" \
  --scheduler normal --steps 25
```

#### Karras
- **Best for**: High-quality final outputs, fine details, portraits
- **Speed**: Slower (1.3x)
- **Steps**: 30-50

```bash
asset-generator generate image \
  --prompt "detailed portrait, intricate jewelry" \
  --scheduler karras --steps 35 --cfg-scale 8.0
```

#### Exponential
- **Best for**: Artistic work, smooth gradients, abstract art
- **Speed**: Medium (1.2x)
- **Steps**: 25-40

```bash
asset-generator generate image \
  --prompt "abstract art, flowing energy, cosmic" \
  --scheduler exponential --steps 30
```

### Performance Considerations

| Scheduler | Relative Speed | Quality | Best Use |
|-----------|----------------|---------|----------|
| Simple | Fastest (1.0x) | Good | Iteration |
| Normal | Fast (1.1x) | Better | Production |
| Karras | Slower (1.3x) | Best | Final renders |
| Exponential | Medium (1.2x) | Artistic | Creative work |
| SGM Uniform | Variable | Experimental | Research |

---

## Skimmed CFG {#skimmed-cfg}

### Overview

**Status**: ✅ **IMPLEMENTED**

Skimmed CFG (Classifier-Free Guidance) is an advanced sampling technique that improves image generation quality and speed by applying a more efficient guidance strategy during the denoising process.

### How It Works

Traditional CFG requires computing both conditional and unconditional predictions at every sampling step. Skimmed CFG optimizes this by:

1. **Smart Guidance Application**: Applies full CFG guidance only when most beneficial
2. **Reduced Overhead**: Skips unnecessary computations during certain phases
3. **Dynamic Scaling**: Adjusts guidance strength based on the generation phase

### Benefits

- **Improved Quality**: Better adherence to prompts with enhanced coherence
- **Faster Generation**: Reduced computational overhead in some cases
- **Lower Scale Values**: Achieves similar results with lower CFG scales (2.0-4.0 vs standard 7.0-8.0)
- **Fine-Grained Control**: Apply guidance only during specific generation phases

### Skimmed CFG Usage

#### Basic Usage

Enable Skimmed CFG with default settings:

```bash
asset-generator generate image \
  --prompt "detailed fantasy landscape" \
  --skimmed-cfg
```

#### Custom Scale

Adjust the Skimmed CFG scale (typically lower than standard CFG):

```bash
asset-generator generate image \
  --prompt "portrait of a wizard" \
  --skimmed-cfg \
  --skimmed-cfg-scale 2.5
```

**Recommended scale values**: 2.0-4.0 (compared to standard CFG of 7.0-8.0)

#### Phase-Specific Application

Apply Skimmed CFG only during specific phases of generation:

```bash
# Apply only during middle phase (20%-80% of generation)
asset-generator generate image \
  --prompt "cyberpunk cityscape" \
  --skimmed-cfg \
  --skimmed-cfg-start 0.2 \
  --skimmed-cfg-end 0.8
```

### Recommended Settings by Use Case

#### Photorealistic Images
```bash
--skimmed-cfg \
--skimmed-cfg-scale 3.5 \
--skimmed-cfg-start 0.0 \
--skimmed-cfg-end 1.0
```

#### Artistic/Stylized Images
```bash
--skimmed-cfg \
--skimmed-cfg-scale 2.5 \
--skimmed-cfg-start 0.2 \
--skimmed-cfg-end 0.8
```

#### Fast Iteration/Concepts
```bash
--skimmed-cfg \
--skimmed-cfg-scale 2.0 \
--steps 15
```

### Skimmed CFG Quick Reference

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--skimmed-cfg` | bool | false | Enable Skimmed CFG |
| `--skimmed-cfg-scale` | float | 3.0 | Scale value (lower than standard CFG) |
| `--skimmed-cfg-start` | float | 0.0 | Start phase (0.0 = beginning) |
| `--skimmed-cfg-end` | float | 1.0 | End phase (1.0 = end) |

| Use Case | Scale | Start | End | Steps |
|----------|-------|-------|-----|-------|
| **Photorealistic** | 3.5 | 0.0 | 1.0 | 30 |
| **Artistic/Stylized** | 2.5 | 0.2 | 0.8 | 25 |
| **Fast Iteration** | 2.0 | 0.0 | 1.0 | 15 |
| **High Detail** | 4.0 | 0.1 | 0.9 | 40 |

---

# LoRA Support {#lora-support}

## Overview

The Asset Generator CLI supports **LoRA (Low-Rank Adaptation)** models for enhanced image generation. LoRAs allow you to fine-tune and customize the style, content, and characteristics of generated images without changing the base model.

## What are LoRAs?

LoRAs are lightweight model adaptations that modify specific aspects of image generation:
- **Style modifications** (anime, photorealistic, artistic styles)
- **Character consistency** (specific characters, faces)
- **Content enhancement** (details, quality improvements)
- **Concept injection** (specific objects, themes, aesthetics)

LoRAs are typically much smaller than full models (5-200MB vs 2-7GB) and can be combined for powerful customization.

## Basic LoRA Usage {#lora-basic-usage}

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

## Advanced LoRA Techniques {#advanced-lora-techniques}

### Negative Weights

Use negative weights to **remove** or **reduce** unwanted styles:

```bash
asset-generator generate image \
  --prompt "portrait of a person" \
  --lora "realistic-faces:1.0" \
  --lora "cartoon-style:-0.5"  # Remove cartoon-like features
```

### Configuration File

Define default LoRAs in your config file:

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

### Alternative Syntax

Use separate flags for names and weights:

```bash
asset-generator generate image \
  --prompt "fantasy landscape" \
  --lora "fantasy-art" \
  --lora "detailed-background" \
  --lora-weight 0.9 \
  --lora-weight 0.6
```

### Practical LoRA Examples

#### Anime Character Generation

```bash
asset-generator generate image \
  --prompt "anime girl with blue hair, school uniform" \
  --model "stable-diffusion-xl" \
  --lora "anime-style-v2:0.9" \
  --lora "detailed-eyes:0.7" \
  --batch 4 \
  --save-images
```

#### Photorealistic Portrait

```bash
asset-generator generate image \
  --prompt "portrait of a young woman, natural lighting" \
  --lora "realistic-faces:1.0" \
  --lora "professional-photography:0.8" \
  --width 768 --height 1024 \
  --save-images
```

#### Cyberpunk Scene

```bash
asset-generator generate image \
  --prompt "futuristic city street, rain, neon signs" \
  --lora "cyberpunk-aesthetic:1.1" \
  --lora "neon-lights:0.8" \
  --lora "cinematic-lighting:0.6" \
  --cfg-scale 8.0 \
  --steps 40
```

## LoRA Best Practices {#lora-best-practices}

1. **Start with one LoRA**: Test each LoRA individually before combining
2. **Use moderate weights**: Start with 0.7-0.9 for most LoRAs
3. **Limit combinations**: 2-3 LoRAs work best; more can cause conflicts
4. **Match prompts to LoRAs**: Align your prompt with the LoRA's training domain
5. **Test weight ranges**: Small weight adjustments (±0.1) can make big differences
6. **Document your setups**: Save working LoRA combinations in config files

### Weight Guidelines

| Weight | Effect |
|--------|--------|
| `-0.5 to 0.0` | Removes/reduces style elements |
| `0.0` | No effect |
| `0.5-0.7` | Subtle influence |
| `0.8-1.0` | Standard strength |
| `1.1-1.5` | Strong influence |
| `1.6-2.0` | Very strong (may overpower) |

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

### LoRA Troubleshooting

| Problem | Solution |
|---------|----------|
| LoRA not taking effect | Check name (case-sensitive), increase weight |
| Too strong/overpowering | Reduce weight (try 0.6-0.8) |
| Unexpected results | Use fewer LoRAs (2-3 max) |
| Conflicts between LoRAs | Reorder or reduce weights |

---

# Image Download & Organization {#image-download-organization}

## Overview

The `--save-images` flag downloads generated images from the SwarmUI server directly to your local disk. When combined with the `--filename-template` flag, you can customize filenames using various placeholders for maximum organization.

### Why Use Image Download?

- 💾 **Preserve images locally** - Keep permanent copies on your disk
- 📂 **Organize images** in custom directories
- 🏷️ **Use custom filenames** with metadata (seed, model, dimensions, etc.)
- 🔄 **Work offline** with generated images
- 🎨 **Build local collections** of generated art
- ⚡ **Automatic download** after generation completes
- 🎯 **Batch processing** support for multiple images

## Image Download Quick Start

### Enable Image Download

```bash
# Download to current directory
asset-generator generate image --prompt "your prompt" --save-images

# Download to specific directory
asset-generator generate image --prompt "your prompt" --save-images --output-dir ./my-images

# Download with custom filenames
asset-generator generate image --prompt "fantasy landscape" --batch 5 --save-images \
  --filename-template "landscape-{index}-seed{seed}.png"
```

## Filename Templates {#filename-templates}

By default, downloaded images keep their original filename from the server. However, you can customize filenames using the `--filename-template` flag with various placeholders that get replaced with actual values.

## Template Placeholders {#template-placeholders}

### Index Placeholders

| Placeholder | Description | Example |
|------------|-------------|---------|
| `{index}` or `{i}` | Zero-padded index (3 digits) | `000`, `001`, `002` |
| `{index1}` or `{i1}` | One-based index (no padding) | `1`, `2`, `3` |

### Time Placeholders

| Placeholder | Description | Example |
|------------|-------------|---------|
| `{timestamp}` or `{ts}` | Unix timestamp | `1696723200` |
| `{datetime}` or `{dt}` | Full datetime | `2024-10-08_14-30-45` |
| `{date}` | Date only | `2024-10-08` |
| `{time}` | Time only | `14-30-45` |

### Generation Parameter Placeholders

| Placeholder | Description | Example |
|------------|-------------|---------|
| `{seed}` | Seed value used | `42` |
| `{model}` | Model name | `flux-dev` |
| `{width}` | Image width | `1024` |
| `{height}` | Image height | `768` |
| `{prompt}` | First 50 chars of prompt (sanitized) | `a_beautiful_landscape` |

### Original Filename Placeholders

| Placeholder | Description | Example |
|------------|-------------|---------|
| `{original}` | Complete original filename | `image-abc123.png` |
| `{ext}` | Extension only (with dot) | `.png` |

## Organization Examples {#organization-examples}

### Sequential Numbering with Seed

Great for tracking which seed generated each image:

```bash
asset-generator generate image \
  --prompt "fantasy castle" \
  --seed 42 \
  --batch 10 \
  --save-images \
  --filename-template "castle-seed{seed}-{i1}.png"
```

**Output:**
```
castle-seed42-1.png
castle-seed42-2.png
...
castle-seed42-10.png
```

### Organized by Date and Model

Perfect for daily generation workflows:

```bash
asset-generator generate image \
  --prompt "abstract art" \
  --model "sdxl-turbo" \
  --batch 5 \
  --save-images \
  --filename-template "{date}/{model}-{index}.png"
```

**Output:**
```
2024-10-08/sdxl-turbo-000.png
2024-10-08/sdxl-turbo-001.png
...
```

### Complex Organization Template

Combine multiple placeholders for maximum organization:

```bash
asset-generator generate image \
  --prompt "cyberpunk street" \
  --model "flux-dev" \
  --width 1024 \
  --height 768 \
  --seed 12345 \
  --batch 3 \
  --save-images \
  --output-dir ./renders \
  --filename-template "{date}/{model}/{prompt}-{width}x{height}-seed{seed}-{i1}.png"
```

**Output:**
```
./renders/2024-10-08/flux-dev/cyberpunk_street-1024x768-seed12345-1.png
./renders/2024-10-08/flux-dev/cyberpunk_street-1024x768-seed12345-2.png
./renders/2024-10-08/flux-dev/cyberpunk_street-1024x768-seed12345-3.png
```

### Special Behaviors

#### Automatic Extension Handling

If your template doesn't include an extension, the original file extension is automatically appended:

```bash
--filename-template "image-{index}"
# If original is image.png: image-000.png
```

#### Filename Sanitization

The `{prompt}` placeholder is automatically sanitized for filesystem compatibility:
- Spaces become underscores: `"hello world"` → `hello_world`
- Invalid characters are removed: `"cat/dog"` → `catdog`
- Truncated to 50 characters maximum

#### Directory Creation

Templates can include directory separators (`/`) and directories are created automatically:

```bash
--filename-template "{date}/{model}/image-{index}.png"
```

### Organization Tips

1. **Use zero-padded indices** (`{index}`) for proper file sorting
2. **Include seed for reproducibility** when you need to regenerate images
3. **Add timestamps for archival** to avoid filename collisions
4. **Combine model and dimensions** when testing different configurations
5. **Keep templates short** to avoid exceeding filesystem path limits
6. **Test templates first** with `--batch 1` to verify output

---

## Combined Feature Examples

### Complete Workflow with All Features

```bash
asset-generator generate image \
  --prompt "detailed fantasy portrait of an elven mage" \
  --model "flux-dev" \
  --width 768 \
  --height 1024 \
  --scheduler karras \
  --steps 35 \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.0 \
  --lora "fantasy-art:0.9" \
  --lora "detailed-faces:0.7" \
  --batch 5 \
  --save-images \
  --output-dir ./fantasy-portraits \
  --filename-template "{model}-{scheduler}-{prompt}-{width}x{height}-{index}.png"
```

This command:
1. Uses high-quality Karras scheduler
2. Applies Skimmed CFG for improved quality
3. Combines two LoRAs for fantasy art style
4. Generates 5 variations
5. Saves with descriptive filenames including all parameters

### Pipeline with User Guide Features

```bash
asset-generator pipeline \
  --file character-designs.yaml \
  --scheduler karras \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.5 \
  --lora "character-design:0.9" \
  --save-images \
  --output-dir ./character-collection \
  --filename-template "{date}/{scheduler}/{index}-{prompt}.png"
```

## Quick Reference

### Essential Commands

```bash
# High-quality generation with LoRAs and custom naming
asset-generator generate image \
  --prompt "your prompt" \
  --scheduler karras \
  --lora "style-lora:0.8" \
  --save-images \
  --filename-template "{date}-{prompt}-{index}.png"

# Fast iteration with Skimmed CFG
asset-generator generate image \
  --prompt "test concept" \
  --skimmed-cfg \
  --skimmed-cfg-scale 2.0 \
  --steps 15

# Batch with complete organization
asset-generator generate image \
  --prompt "character design" \
  --batch 10 \
  --save-images \
  --filename-template "{model}/{prompt}-seed{seed}-{i1}.png"
```

### Flag Quick Reference

| Feature | Key Flags | Example Values |
|---------|-----------|----------------|
| **Scheduler** | `--scheduler` | `simple`, `normal`, `karras`, `exponential` |
| **Skimmed CFG** | `--skimmed-cfg`, `--skimmed-cfg-scale` | `true`, `2.0-4.0` |
| **LoRA** | `--lora` | `"style-name:0.8"` |
| **Download** | `--save-images`, `--output-dir` | `true`, `./output` |
| **Templates** | `--filename-template` | `"{date}-{prompt}-{index}.png"` |

## See Also

- [Commands Reference](COMMANDS.md) - Complete command documentation
- [Pipeline Processing](PIPELINE.md) - Batch generation workflows  
- [Postprocessing](POSTPROCESSING.md) - Auto-crop, downscaling, metadata stripping
- [Quick Start Guide](QUICKSTART.md) - Getting started with the CLI
- [Development Guide](DEVELOPMENT.md) - Architecture and API details