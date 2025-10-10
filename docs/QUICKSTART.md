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
asset-generator config init
```

This creates `~/.asset-generator/config.yaml` with default settings.

### 2. Configure API Endpoint

If your SwarmUI is not running on localhost:7801, set the API URL:

```bash
asset-generator config set api-url http://your-asset-generator-host:7801
```

### 3. (Optional) Set API Key

If your SwarmUI requires authentication:

```bash
asset-generator config set api-key YOUR_API_KEY
```

## Basic Usage

### Generate Your First Image

```bash
asset-generator generate image --prompt "a beautiful sunset over mountains"
```

This will:
1. Connect to your SwarmUI instance
2. Generate an image with the prompt
3. Display the result in table format

### View Available Models

```bash
asset-generator models list
```

### Generate with Specific Parameters

```bash
asset-generator generate image \
  --prompt "cyberpunk city at night, neon lights, rainy" \
  --width 1024 \
  --height 1024 \
  --steps 30 \
  --cfg-scale 7.5 \
  --model stable-diffusion-xl
```

### Save Output to File

```bash
asset-generator generate image \
  --prompt "a cute cat wearing sunglasses" \
  --format json \
  --output result.json
```

### Download Generated Images

Download images directly to your local disk:

```bash
# Download to current directory
asset-generator generate image \
  --prompt "beautiful landscape" \
  --save-images

# Download to specific directory
asset-generator generate image \
  --prompt "fantasy castle" \
  --save-images \
  --output-dir ./my-art

# Batch download
asset-generator generate image \
  --prompt "abstract art" \
  --batch 5 \
  --save-images \
  --output-dir ./batch-output
```

The `--save-images` flag automatically downloads all generated images after completion.


## Common Tasks

### Batch Generation

Generate multiple images at once:

```bash
asset-generator generate image \
  --prompt "fantasy landscape with mountains" \
  --batch 4
```

### Reproducible Generation

Use a seed for reproducible results:

```bash
asset-generator generate image \
  --prompt "portrait of a wizard" \
  --seed 42
```

### Generate with Real-Time Progress (WebSocket)

Use WebSocket for authentic progress updates during generation:

```bash
asset-generator generate image \
  --prompt "detailed fantasy landscape" \
  --websocket
```

**Note**: WebSocket support requires SwarmUI server. Falls back to HTTP automatically if unavailable.

Benefits:
- Real-time progress updates (not simulated)
- Particularly useful for long-running generations (Flux models: 5-10 minutes)
- Live status feedback during generation

### Use Negative Prompts

```bash
asset-generator generate image \
  --prompt "beautiful portrait photo" \
  --negative-prompt "blurry, low quality, distorted"
```

### Different Output Formats

#### JSON Output
```bash
asset-generator models list --format json
```

#### YAML Output
```bash
asset-generator config view --format yaml
```

## Configuration Examples

### Using Environment Variables

```bash
export SWARMUI_API_URL=http://localhost:7801
export SWARMUI_API_KEY=your-key
export SWARMUI_FORMAT=json

asset-generator models list  # Uses environment config
```

### Config File

Edit `~/.asset-generator/config.yaml`:

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
asset-generator generate image \
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
  asset-generator generate image --prompt "$prompt" --quiet
done
```

### Batch Processing with Different Seeds

```bash
for seed in {1..5}; do
  asset-generator generate image \
    --prompt "fantasy castle" \
    --seed $seed \
    --output "castle_${seed}.json" \
    --quiet
done
```

## Tips and Tricks

### 1. Use Quiet Mode for Scripts

```bash
asset-generator --quiet generate image --prompt "test"
```

Suppresses progress messages, only shows errors.

### 2. Verbose Mode for Debugging

```bash
asset-generator --verbose models list
```

Shows detailed information about API calls.

### 3. Check Current Configuration

```bash
asset-generator config view
```

See all active settings and their sources.

### 4. Test Connection

```bash
asset-generator models list
```

Quick way to verify your SwarmUI is accessible.

### 5. View Full Help

```bash
asset-generator --help
asset-generator generate --help
asset-generator generate image --help
```

Every command has detailed help with examples.

## Troubleshooting

### Problem: Cannot connect to SwarmUI

**Solution:**
```bash
# Check config
asset-generator config get api-url

# Test connection directly
curl http://localhost:7801/API/ListModels

# Update API URL if needed
asset-generator config set api-url http://correct-host:7801
```

### Problem: "prompt is required" error

**Solution:**
The `--prompt` flag is mandatory for generation:
```bash
asset-generator generate image --prompt "your prompt here"
```

### Problem: Output format not recognized

**Solution:**
Valid formats are: table, json, yaml
```bash
asset-generator models list --format json
```

### Problem: Config file not found

**Solution:**
Initialize the config first:
```bash
asset-generator config init
```

## Advanced Usage

### Custom Sampling Methods

```bash
asset-generator generate image \
  --prompt "portrait" \
  --sampler euler_a \
  --steps 20
```

Available samplers: euler_a, euler, heun, dpm_2, dpm_2_ancestral, lms, dpm_fast, dpm_adaptive

### High-Quality Generation

```bash
asset-generator generate image \
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
asset-generator models get stable-diffusion-xl
```

Then use specific models:
```bash
asset-generator generate image \
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
        ['asset-generator', 'generate', 'image',
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
    `asset-generator generate image --prompt "${prompt}" --format json`
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
	@asset-generator generate image --prompt "logo design" --output logo.json
	@asset-generator generate image --prompt "banner image" --output banner.json
	@echo "Assets generated!"
```

## Next Steps

1. **Read the Full Documentation**: Check `README.md` for comprehensive docs
2. **Explore All Commands**: Run `asset-generator --help` to see all available commands
3. **Customize Defaults**: Edit `~/.asset-generator/config.yaml` for your workflow
4. **Contribute**: Found a bug or want a feature? Visit our GitHub!

## Quick Reference Card

```bash
# Configuration
asset-generator config init                    # Initialize config
asset-generator config set KEY VALUE           # Set config value
asset-generator config view                    # View current config

# Generation
asset-generator generate image -p "prompt"     # Basic generation
asset-generator generate image -p "..." -b 4   # Batch of 4
asset-generator generate image -p "..." --seed 42  # With seed

# Models
asset-generator models list                    # List all models
asset-generator models get MODEL_NAME          # Get model details

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
