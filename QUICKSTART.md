# SwarmUI CLI - Quick Start Guide

Welcome to the SwarmUI CLI! This guide will get you up and running in 5 minutes.

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/opd-ai/asset-generator.git
cd asset-generator

# Build the binary
make build

# (Optional) Install system-wide
sudo make install
```

### Option 2: Download Pre-built Binary

Download the latest release from [GitHub Releases](https://github.com/opd-ai/asset-generator/releases) and add it to your PATH.

## Initial Setup

### 1. Initialize Configuration

```bash
swarmui config init
```

This creates `~/.swarmui/config.yaml` with default settings.

### 2. Configure API Endpoint

If your SwarmUI is not running on localhost:7801, set the API URL:

```bash
swarmui config set api-url http://your-swarmui-host:7801
```

### 3. (Optional) Set API Key

If your SwarmUI requires authentication:

```bash
swarmui config set api-key YOUR_API_KEY
```

## Basic Usage

### Generate Your First Image

```bash
swarmui generate image --prompt "a beautiful sunset over mountains"
```

This will:
1. Connect to your SwarmUI instance
2. Generate an image with the prompt
3. Display the result in table format

### View Available Models

```bash
swarmui models list
```

### Generate with Specific Parameters

```bash
swarmui generate image \
  --prompt "cyberpunk city at night, neon lights, rainy" \
  --width 1024 \
  --height 1024 \
  --steps 30 \
  --cfg-scale 7.5 \
  --model stable-diffusion-xl
```

### Save Output to File

```bash
swarmui generate image \
  --prompt "a cute cat wearing sunglasses" \
  --format json \
  --output result.json
```

## Common Tasks

### Batch Generation

Generate multiple images at once:

```bash
swarmui generate image \
  --prompt "fantasy landscape with mountains" \
  --batch 4
```

### Reproducible Generation

Use a seed for reproducible results:

```bash
swarmui generate image \
  --prompt "portrait of a wizard" \
  --seed 42
```

### Use Negative Prompts

```bash
swarmui generate image \
  --prompt "beautiful portrait photo" \
  --negative-prompt "blurry, low quality, distorted"
```

### Different Output Formats

#### JSON Output
```bash
swarmui models list --format json
```

#### YAML Output
```bash
swarmui config view --format yaml
```

## Configuration Examples

### Using Environment Variables

```bash
export SWARMUI_API_URL=http://localhost:7801
export SWARMUI_API_KEY=your-key
export SWARMUI_FORMAT=json

swarmui models list  # Uses environment config
```

### Config File

Edit `~/.swarmui/config.yaml`:

```yaml
api-url: http://localhost:7801
api-key: your-api-key-here
format: table

# Default generation parameters
generate:
  model: stable-diffusion-xl
  steps: 25
  width: 768
  height: 768
  cfg-scale: 7.0
  sampler: euler_a
```

Now all generate commands will use these defaults!

## Scripting Examples

### Extract Image Paths

```bash
swarmui generate image \
  --prompt "mountain landscape" \
  --format json | jq -r '.image_paths[]'
```

### Loop Through Multiple Prompts

```bash
#!/bin/bash
prompts=(
  "sunset over ocean"
  "mountain landscape"
  "city skyline at night"
)

for prompt in "${prompts[@]}"; do
  echo "Generating: $prompt"
  swarmui generate image --prompt "$prompt" --quiet
done
```

### Batch Processing with Different Seeds

```bash
for seed in {1..5}; do
  swarmui generate image \
    --prompt "fantasy castle" \
    --seed $seed \
    --output "castle_${seed}.json" \
    --quiet
done
```

## Tips and Tricks

### 1. Use Quiet Mode for Scripts

```bash
swarmui --quiet generate image --prompt "test"
```

Suppresses progress messages, only shows errors.

### 2. Verbose Mode for Debugging

```bash
swarmui --verbose models list
```

Shows detailed information about API calls.

### 3. Check Current Configuration

```bash
swarmui config view
```

See all active settings and their sources.

### 4. Test Connection

```bash
swarmui models list
```

Quick way to verify your SwarmUI is accessible.

### 5. View Full Help

```bash
swarmui --help
swarmui generate --help
swarmui generate image --help
```

Every command has detailed help with examples.

## Troubleshooting

### Problem: Cannot connect to SwarmUI

**Solution:**
```bash
# Check config
swarmui config get api-url

# Test connection directly
curl http://localhost:7801/API/ListModels

# Update API URL if needed
swarmui config set api-url http://correct-host:7801
```

### Problem: "prompt is required" error

**Solution:**
The `--prompt` flag is mandatory for generation:
```bash
swarmui generate image --prompt "your prompt here"
```

### Problem: Output format not recognized

**Solution:**
Valid formats are: table, json, yaml
```bash
swarmui models list --format json
```

### Problem: Config file not found

**Solution:**
Initialize the config first:
```bash
swarmui config init
```

## Advanced Usage

### Custom Sampling Methods

```bash
swarmui generate image \
  --prompt "portrait" \
  --sampler euler_a \
  --steps 20
```

Available samplers: euler_a, euler, heun, dpm_2, dpm_2_ancestral, lms, dpm_fast, dpm_adaptive

### High-Quality Generation

```bash
swarmui generate image \
  --prompt "detailed portrait photograph" \
  --width 1024 \
  --height 1024 \
  --steps 50 \
  --cfg-scale 8.0
```

More steps = better quality but slower generation.

### Model-Specific Settings

Get model details first:
```bash
swarmui models get stable-diffusion-xl
```

Then use specific models:
```bash
swarmui generate image \
  --prompt "test" \
  --model stable-diffusion-xl
```

## Integration Examples

### Python Script

```python
#!/usr/bin/env python3
import subprocess
import json

def generate_image(prompt):
    result = subprocess.run(
        ['swarmui', 'generate', 'image',
         '--prompt', prompt,
         '--format', 'json'],
        capture_output=True,
        text=True
    )
    return json.loads(result.stdout)

# Usage
result = generate_image("a beautiful landscape")
print(f"Generated: {result['image_paths']}")
```

### Node.js Script

```javascript
const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

async function generateImage(prompt) {
  const { stdout } = await execPromise(
    `swarmui generate image --prompt "${prompt}" --format json`
  );
  return JSON.parse(stdout);
}

// Usage
generateImage('cyberpunk city')
  .then(result => console.log('Generated:', result.image_paths));
```

### Makefile Integration

```makefile
.PHONY: generate-assets

generate-assets:
	@echo "Generating assets..."
	@swarmui generate image --prompt "logo design" --output logo.json
	@swarmui generate image --prompt "banner image" --output banner.json
	@echo "Assets generated!"
```

## Next Steps

1. **Read the Full Documentation**: Check `README.md` for comprehensive docs
2. **Explore All Commands**: Run `swarmui --help` to see all available commands
3. **Customize Defaults**: Edit `~/.swarmui/config.yaml` for your workflow
4. **Contribute**: Found a bug or want a feature? Visit our GitHub!

## Quick Reference Card

```bash
# Configuration
swarmui config init                    # Initialize config
swarmui config set KEY VALUE           # Set config value
swarmui config view                    # View current config

# Generation
swarmui generate image -p "prompt"     # Basic generation
swarmui generate image -p "..." -b 4   # Batch of 4
swarmui generate image -p "..." --seed 42  # With seed

# Models
swarmui models list                    # List all models
swarmui models get MODEL_NAME          # Get model details

# Output Control
--format json                          # JSON output
--format yaml                          # YAML output
--output file.json                     # Save to file
--quiet                                # Suppress progress
--verbose                              # Show debug info
```

## Getting Help

- **Documentation**: See `README.md` and `DEVELOPMENT.md`
- **Issues**: [GitHub Issues](https://github.com/opd-ai/asset-generator/issues)
- **Examples**: Look at command `--help` output

---

Happy generating! 🎨
