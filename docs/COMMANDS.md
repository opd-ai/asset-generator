[🏠 Docs Home](README.md) | [📚 Quick Start](QUICKSTART.md) | [🔗 Pipeline](PIPELINE.md) | [👤 User Guide](USER_GUIDE.md)

---

# Commands Reference

> **Complete reference for all Asset Generator CLI commands**

This document provides comprehensive documentation for all Asset Generator CLI commands, flags, and usage patterns.

## Table of Contents

- [Global Options](#global-options)
- [Generation Commands](#generation-commands)
  - [generate image](#generate-image)
- [Pipeline Commands](#pipeline-commands)
  - [pipeline](#pipeline)
- [Model Commands](#model-commands)
  - [models list](#models-list)
- [Configuration Commands](#configuration-commands)
  - [config init](#config-init)
  - [config get](#config-get)
  - [config set](#config-set)
  - [config view](#config-view)
- [Conversion Commands](#conversion-commands)
  - [convert svg](#convert-svg)
- [Postprocessing Commands](#postprocessing-commands)
  - [crop](#crop)
  - [downscale](#downscale)
- [Status Commands](#status-commands)
  - [status](#status)
  - [cancel](#cancel)

---

## Global Options {#global-options}

These flags are available for all commands:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--api-url` | string | `http://localhost:7801` | SwarmUI API endpoint |
| `--api-key` | string | (empty) | API key for authentication |
| `--config` | string | `~/.asset-generator/config.yaml` | Configuration file path |
| `--format` | string | `table` | Output format: table, json, yaml |
| `--output` | string | (stdout) | Output file path |
| `--quiet` | bool | false | Suppress progress output |
| `--verbose` | bool | false | Verbose logging |
| `--help` | bool | false | Show help |

### Examples

```bash
# Use different API endpoint
asset-generator --api-url http://remote-server:7801 models list

# Output to file in JSON format
asset-generator --format json --output models.json models list

# Quiet mode for scripting
asset-generator --quiet generate image --prompt "test"

# Verbose debugging
asset-generator --verbose status
```

---

## Generation Commands {#generation-commands}

### generate image {#generate-image}

Generate images using text-to-image models.

#### Synopsis

```bash
asset-generator generate image --prompt PROMPT [flags]
```

#### Required Flags

| Flag | Type | Description |
|------|------|-------------|
| `--prompt`, `-p` | string | Text prompt for image generation (required) |

#### Generation Parameters

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--model` | string | (from config) | Model name to use |
| `--width` | int | 1024 | Image width in pixels |
| `--height` | int | 1024 | Image height in pixels |
| `--steps` | int | 25 | Number of inference steps |
| `--cfg-scale` | float | 7.5 | Classifier-free guidance scale |
| `--seed` | int | -1 | Seed for reproducibility (-1 = random) |
| `--batch` | int | 1 | Number of images to generate |
| `--sampler` | string | (default) | Sampling method |
| `--scheduler` | string | `simple` | Noise scheduler |
| `--negative-prompt` | string | (empty) | Negative prompt |

#### LoRA Flags

| Flag | Type | Description |
|------|------|-------------|
| `--lora` | string | LoRA model with optional weight (name:weight) |
| `--lora-weight` | float | Weight for LoRA (when not specified inline) |
| `--lora-default-weight` | float | Default weight for LoRAs without weights |

#### Skimmed CFG Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--skimmed-cfg` | bool | false | Enable Skimmed CFG |
| `--skimmed-cfg-scale` | float | 3.0 | Skimmed CFG scale |
| `--skimmed-cfg-start` | float | 0.0 | Start phase (0.0-1.0) |
| `--skimmed-cfg-end` | float | 1.0 | End phase (0.0-1.0) |

#### Image Download Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--save-images` | bool | false | Download images to local disk |
| `--output-dir` | string | `.` | Directory to save images |
| `--filename-template` | string | (empty) | Custom filename template |

#### Postprocessing Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--auto-crop` | bool | false | Remove whitespace borders |
| `--crop-threshold` | int | 10 | Whitespace detection threshold |
| `--crop-preserve-aspect` | bool | false | Maintain original aspect ratio |
| `--downscale-width` | int | 0 | Target width for downscaling |
| `--downscale-height` | int | 0 | Target height for downscaling |
| `--downscale-percentage` | int | 0 | Scale by percentage |

#### Examples

##### Basic Generation

```bash
asset-generator generate image --prompt "beautiful sunset over mountains"
```

##### High-Quality Portrait

```bash
asset-generator generate image \
  --prompt "professional portrait, studio lighting" \
  --width 768 \
  --height 1024 \
  --scheduler karras \
  --steps 35 \
  --cfg-scale 8.0
```

##### Batch with LoRAs

```bash
asset-generator generate image \
  --prompt "anime character, detailed eyes" \
  --lora "anime-style:0.9" \
  --lora "detailed-faces:0.7" \
  --batch 5 \
  --save-images \
  --filename-template "character-{index}-seed{seed}.png"
```

##### Complete Workflow

```bash
asset-generator generate image \
  --prompt "fantasy landscape, mystical forest" \
  --model "stable-diffusion-xl" \
  --width 1024 \
  --height 768 \
  --scheduler karras \
  --steps 30 \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.0 \
  --lora "fantasy-art:0.8" \
  --batch 3 \
  --save-images \
  --output-dir ./fantasy-art \
  --filename-template "{date}/{model}-{prompt}-{index}.png" \
  --auto-crop \
  --downscale-width 512
```

---

## Pipeline Commands {#pipeline-commands}

### pipeline {#pipeline}

Process YAML pipeline files for automated batch generation.

#### Synopsis

```bash
asset-generator pipeline --file FILE [flags]
```

#### Required Flags

| Flag | Type | Description |
|------|------|-------------|
| `--file`, `-f` | string | Pipeline YAML file path (required) |

#### Pipeline Control Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--output-dir` | string | `./pipeline-output` | Base output directory |
| `--dry-run` | bool | false | Preview without generating |
| `--continue-on-error` | bool | false | Continue if individual assets fail |
| `--base-seed` | int | -1 | Base seed for reproducibility (-1 = random) |

#### Generation Override Flags

All generation flags from `generate image` can be used to override pipeline defaults:

```bash
asset-generator pipeline --file assets.yaml \
  --scheduler karras \
  --steps 30 \
  --save-images \
  --auto-crop
```

#### Examples

##### Basic Pipeline

```bash
asset-generator pipeline --file character-designs.yaml
```

##### Dry Run Preview

```bash
asset-generator pipeline --file assets.yaml --dry-run
```

##### Production Pipeline

```bash
asset-generator pipeline \
  --file production-assets.yaml \
  --output-dir ./final-assets \
  --scheduler karras \
  --steps 35 \
  --continue-on-error \
  --save-images \
  --auto-crop \
  --downscale-width 1024
```

---

## Model Commands {#model-commands}

### models list {#models-list}

List available models on the SwarmUI server.

#### Synopsis

```bash
asset-generator models list [flags]
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--loaded-only` | bool | false | Show only loaded models |

#### Examples

##### List All Models

```bash
asset-generator models list
```

##### JSON Output for Scripting

```bash
asset-generator models list --format json
```

##### Filter Loaded Models

```bash
asset-generator models list --loaded-only
```

##### Extract Specific Information

```bash
# Get model names only
asset-generator models list --format json | jq -r '.[].name'

# Filter for specific type
asset-generator models list --format json | jq '.[] | select(.type == "Stable-Diffusion")'

# Count total models
asset-generator models list --format json | jq 'length'
```

---

## Configuration Commands {#configuration-commands}

### config init {#config-init}

Initialize configuration file with default settings.

#### Synopsis

```bash
asset-generator config init [flags]
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force` | bool | false | Overwrite existing config |

#### Examples

```bash
# Create initial config
asset-generator config init

# Overwrite existing config
asset-generator config init --force
```

### config get {#config-get}

Get a configuration value.

#### Synopsis

```bash
asset-generator config get KEY
```

#### Examples

```bash
# Get API URL
asset-generator config get api-url

# Get output format
asset-generator config get output-format
```

### config set {#config-set}

Set a configuration value.

#### Synopsis

```bash
asset-generator config set KEY VALUE
```

#### Examples

```bash
# Set API URL
asset-generator config set api-url http://localhost:7801

# Set API key
asset-generator config set api-key your-api-key-here

# Set default output format
asset-generator config set output-format json
```

### config view {#config-view}

View current configuration.

#### Synopsis

```bash
asset-generator config view [flags]
```

#### Examples

```bash
# View current config
asset-generator config view

# View with file location
asset-generator config view --verbose
```

---

## Conversion Commands {#conversion-commands}

### convert svg {#convert-svg}

Convert images to SVG format using geometric shapes or edge tracing.

#### Synopsis

```bash
asset-generator convert svg INPUT [flags]
```

#### Required Arguments

| Argument | Type | Description |
|----------|------|-------------|
| `INPUT` | string | Input image file path |

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--output`, `-o` | string | (auto) | Output SVG file path |
| `--method` | string | `primitive` | Conversion method: primitive, gotrace |
| `--shapes` | int | 100 | Number of shapes (primitive method) |
| `--mode` | int | 1 | Shape mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, etc. |
| `--alpha` | int | 128 | Shape alpha transparency (0-255) |
| `--repeat` | int | 0 | Optimization iterations |

#### Examples

##### Basic Conversion

```bash
asset-generator convert svg photo.jpg
```

##### High Quality with Many Shapes

```bash
asset-generator convert svg artwork.png \
  --shapes 500 \
  --mode 3 \
  --output artwork-detailed.svg
```

##### Gotrace Method for Line Art

```bash
asset-generator convert svg lineart.png \
  --method gotrace \
  --output lineart-vector.svg
```

##### Batch Conversion

```bash
#!/bin/bash
for file in *.png; do
    asset-generator convert svg "$file" \
      --shapes 200 \
      --mode 1 \
      --output "${file%.png}.svg"
done
```

---

## Postprocessing Commands {#postprocessing-commands}

### crop {#crop}

Remove whitespace borders from images.

#### Synopsis

```bash
asset-generator crop INPUT [flags]
```

#### Required Arguments

| Argument | Type | Description |
|----------|------|-------------|
| `INPUT` | string | Input image file path |

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--output`, `-o` | string | (auto) | Output file path |
| `--threshold` | int | 10 | Whitespace detection threshold |
| `--preserve-aspect` | bool | false | Maintain original aspect ratio |

#### Examples

```bash
# Basic crop
asset-generator crop image-with-borders.png

# Custom threshold
asset-generator crop noisy-image.png --threshold 20

# Preserve aspect ratio
asset-generator crop logo.png --preserve-aspect --output logo-cropped.png
```

### downscale {#downscale}

Resize images with high-quality Lanczos filtering.

#### Synopsis

```bash
asset-generator downscale INPUT [flags]
```

#### Required Arguments

| Argument | Type | Description |
|----------|------|-------------|
| `INPUT` | string | Input image file path |

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--output`, `-o` | string | (auto) | Output file path |
| `--width` | int | 0 | Target width (0 = calculate from height) |
| `--height` | int | 0 | Target height (0 = calculate from width) |
| `--percentage` | int | 0 | Scale by percentage |

#### Examples

```bash
# Resize to specific width
asset-generator downscale large-image.png --width 1024

# Scale by percentage
asset-generator downscale huge-image.png --percentage 50

# Specific dimensions
asset-generator downscale wallpaper.png --width 1920 --height 1080
```

---

## Status Commands {#status-commands}

### status {#status}

Check SwarmUI server health and configuration.

#### Synopsis

```bash
asset-generator status [flags]
```

#### Examples

```bash
# Basic status check
asset-generator status

# JSON for scripting
asset-generator status --format json

# Check specific server
asset-generator --api-url http://remote:7801 status
```

#### Output Information

- Server connectivity and response time
- Session information
- Backend status and loaded models
- Model counts
- System information (if available)

### cancel {#cancel}

Cancel ongoing or queued image generations.

#### Synopsis

```bash
asset-generator cancel [flags]
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | Cancel all queued generations |

#### Examples

```bash
# Cancel current generation
asset-generator cancel

# Cancel all queued generations
asset-generator cancel --all

# Quiet cancellation for scripts
asset-generator --quiet cancel
```

---

## Flag Reference Quick Lookup

### Generation Flags

| Category | Key Flags | Values |
|----------|-----------|---------|
| **Basic** | `--prompt`, `--model` | Required, model name |
| **Dimensions** | `--width`, `--height` | 512-2048 pixels |
| **Quality** | `--steps`, `--cfg-scale` | 10-50 steps, 1.0-20.0 scale |
| **Reproducibility** | `--seed` | Integer or -1 for random |
| **Scheduler** | `--scheduler` | simple, normal, karras, exponential |
| **LoRA** | `--lora` | "name:weight" format |
| **Skimmed CFG** | `--skimmed-cfg`, `--skimmed-cfg-scale` | true/false, 1.0-5.0 |

### Output Flags

| Category | Key Flags | Values |
|----------|-----------|---------|
| **Download** | `--save-images`, `--output-dir` | true/false, directory path |
| **Naming** | `--filename-template` | Template with placeholders |
| **Processing** | `--auto-crop`, `--downscale-width` | true/false, pixels |

### Control Flags

| Category | Key Flags | Values |
|----------|-----------|---------|
| **Format** | `--format`, `--output` | table/json/yaml, file path |
| **Verbosity** | `--quiet`, `--verbose` | true/false |
| **API** | `--api-url`, `--api-key` | URL, key string |

---

## Common Command Patterns

### Quick Generation

```bash
# Basic
asset-generator generate image -p "prompt here"

# High quality
asset-generator generate image -p "prompt" --scheduler karras --steps 35

# With download
asset-generator generate image -p "prompt" --save-images --output-dir ./images
```

### Batch Processing

```bash
# Pipeline
asset-generator pipeline -f assets.yaml --save-images --auto-crop

# Multiple single generations
for prompt in "cat" "dog" "bird"; do
  asset-generator generate image -p "$prompt" --save-images --quiet
done
```

### Status and Configuration

```bash
# Health check
asset-generator status

# Setup
asset-generator config init
asset-generator config set api-url http://localhost:7801

# List models
asset-generator models list --format json | jq '.[].name'
```

---

## See Also

- [User Guide](USER_GUIDE.md) - Advanced features and generation options
- [Quick Start](QUICKSTART.md) - Getting started guide
- [Pipeline Processing](PIPELINE.md) - Batch generation workflows
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues and solutions


|-------|---------|
| "no active generation" | Nothing to cancel (informational) |
| "failed to get session" | Session/connection issue |
| "request failed" | Network/server problem |

---

## Status Command {#status-command}

The `status` command queries the SwarmUI server and displays comprehensive information about its current state.

### Overview {#status-overview}

The status command provides real-time information about:
- Server connectivity and response time
- Available backends and their operational states
- Current session information
- Model availability and loading status
- System information (GPU, memory, etc.)

### Usage {#status-usage}

```bash
# Basic status check
asset-generator status

# JSON output (for scripting/automation)
asset-generator status --format json

# YAML output
asset-generator status --format yaml

# Verbose output with additional details
asset-generator status -v

# Save status to file
asset-generator status --output status.txt
```

