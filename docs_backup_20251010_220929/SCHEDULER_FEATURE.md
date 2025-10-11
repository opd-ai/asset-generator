# Scheduler Selection Feature

## Overview

The scheduler (also known as noise schedule) controls how noise is added and removed during the diffusion process. Different schedulers can significantly impact generation quality, speed, and characteristics.

**Status**: ✅ **IMPLEMENTED**

The Asset Generator CLI now supports selecting different schedulers via:
- `--scheduler` flag for `generate image` command
- `--scheduler` flag for `pipeline` command
- `scheduler` config file setting
- Defaults to `simple` scheduler

## Available Schedulers

SwarmUI supports the following schedulers:

| Scheduler | Description | Best For |
|-----------|-------------|----------|
| `simple` | Basic linear schedule (default) | General purpose, fast |
| `normal` | Standard scheduler with balanced noise | Most use cases, reliable |
| `karras` | Karras noise schedule | High-quality, detailed images |
| `exponential` | Exponential decay schedule | Smooth transitions, artistic |
| `sgm_uniform` | Uniform noise distribution | Experimental, specialized models |

**Default**: `simple` - Provides reliable results with minimal overhead.

## Usage

### Command Line

#### Generate Image with Scheduler

```bash
# Use default scheduler (simple)
asset-generator generate image \
  --prompt "portrait of a wizard"

# Use Karras scheduler for high quality
asset-generator generate image \
  --prompt "detailed portrait of a wizard" \
  --scheduler karras \
  --steps 30

# Use normal scheduler (standard)
asset-generator generate image \
  --prompt "landscape photograph" \
  --scheduler normal

# Use exponential scheduler for artistic results
asset-generator generate image \
  --prompt "abstract art, cosmic energy" \
  --scheduler exponential
```

#### Pipeline with Scheduler

```bash
# Apply scheduler to all assets in pipeline
asset-generator pipeline \
  --file assets-spec.yaml \
  --output-dir ./assets \
  --scheduler karras \
  --steps 40

# Dry run to preview with scheduler
asset-generator pipeline \
  --file assets-spec.yaml \
  --scheduler normal \
  --dry-run
```

### Configuration File

Set default scheduler in `~/.asset-generator/config.yaml`:

```yaml
api-url: http://localhost:7801
format: table

generate:
  model: stable-diffusion-xl
  steps: 25
  width: 768
  height: 768
  cfg-scale: 7.5
  sampler: euler_a
  scheduler: karras  # Use Karras scheduler by default
```

Now all generations use Karras scheduler unless overridden:

```bash
# Uses karras from config
asset-generator generate image --prompt "wizard"

# Override with normal scheduler
asset-generator generate image --prompt "wizard" --scheduler normal
```

## Scheduler Selection Guide

### Simple (Default)

**When to use:**
- Quick iterations and testing
- When you need fast results
- General purpose generation
- Learning and experimentation

**Characteristics:**
- Fastest processing
- Consistent results
- Good for most use cases
- Minimal computational overhead

**Example:**
```bash
asset-generator generate image \
  --prompt "concept art, character design" \
  --scheduler simple \
  --steps 20
```

### Normal

**When to use:**
- Production-quality assets
- Balanced speed and quality
- Professional workflows
- Most stable results

**Characteristics:**
- Standard diffusion schedule
- Proven reliability
- Good detail preservation
- Widely compatible

**Example:**
```bash
asset-generator generate image \
  --prompt "product photography, professional lighting" \
  --scheduler normal \
  --steps 25
```

### Karras

**When to use:**
- High-quality final outputs
- Fine details are critical
- Portrait and character work
- When you need maximum quality

**Characteristics:**
- Best detail preservation
- Higher computational cost
- Superior edge definition
- Recommended for final renders

**Example:**
```bash
asset-generator generate image \
  --prompt "detailed portrait, intricate jewelry" \
  --scheduler karras \
  --steps 35 \
  --cfg-scale 8.0
```

### Exponential

**When to use:**
- Artistic and creative work
- Smooth gradients
- Abstract compositions
- Experimental styles

**Characteristics:**
- Smoother transitions
- More "painted" look
- Good for artistic styles
- Unique aesthetic quality

**Example:**
```bash
asset-generator generate image \
  --prompt "abstract art, flowing energy, cosmic" \
  --scheduler exponential \
  --steps 30
```

### SGM Uniform

**When to use:**
- Specialized models
- Experimental workflows
- Custom trained models
- Research purposes

**Characteristics:**
- Uniform noise distribution
- Model-specific behavior
- Not all models support it
- Experimental results

**Example:**
```bash
asset-generator generate image \
  --prompt "specialized model test" \
  --scheduler sgm_uniform \
  --steps 25
```

## Combined Workflows

### Pipeline with Different Quality Levels

```yaml
# assets-spec.yaml
assets:
  - name: draft_concepts
    output_dir: drafts
    assets:
      - id: hero_draft
        name: Hero Concept
        prompt: "hero character, action pose"
        # CLI will use --scheduler simple for fast iteration

  - name: final_renders
    output_dir: finals
    assets:
      - id: hero_final
        name: Hero Final
        prompt: "hero character, action pose, detailed"
        # CLI will use --scheduler karras for quality
```

Process with different schedulers:

```bash
# Quick drafts with simple scheduler
asset-generator pipeline \
  --file assets-spec.yaml \
  --output-dir ./drafts \
  --scheduler simple \
  --steps 20

# High-quality finals with karras
asset-generator pipeline \
  --file assets-spec.yaml \
  --output-dir ./finals \
  --scheduler karras \
  --steps 40 \
  --cfg-scale 8.0
```

### Testing Different Schedulers

Script to compare schedulers:

```bash
#!/bin/bash

PROMPT="detailed fantasy castle, sunset"
SEED=42

for scheduler in simple normal karras exponential; do
  asset-generator generate image \
    --prompt "$PROMPT" \
    --seed $SEED \
    --scheduler $scheduler \
    --steps 30 \
    --save-images \
    --output-dir "./scheduler-test" \
    --filename-template "castle-${scheduler}.png"
  
  echo "✓ Generated with $scheduler scheduler"
done

echo "Compare results in ./scheduler-test/"
```

## Integration with Other Features

### Scheduler + Sampler Combinations

Different samplers work better with different schedulers:

```bash
# Euler_a + Simple: Fast, good for iteration
asset-generator generate image \
  --prompt "character design" \
  --sampler euler_a \
  --scheduler simple \
  --steps 20

# DPM++ + Karras: High quality, detailed
asset-generator generate image \
  --prompt "detailed portrait" \
  --sampler dpm_2 \
  --scheduler karras \
  --steps 35

# Heun + Normal: Balanced quality
asset-generator generate image \
  --prompt "landscape photo" \
  --sampler heun \
  --scheduler normal \
  --steps 25
```

### Scheduler + Skimmed CFG

Combine schedulers with Skimmed CFG for optimization:

```bash
# Fast, quality-focused generation
asset-generator generate image \
  --prompt "professional portrait" \
  --scheduler karras \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.0 \
  --steps 25 \
  --cfg-scale 7.5
```

### Scheduler + LoRA

Schedulers affect how LoRAs are applied:

```bash
# Karras scheduler helps preserve LoRA details
asset-generator generate image \
  --prompt "anime character portrait" \
  --lora "anime-style:0.8" \
  --scheduler karras \
  --steps 30
```

## Performance Considerations

### Speed vs Quality

| Scheduler | Relative Speed | Quality | Best Use |
|-----------|----------------|---------|----------|
| Simple | Fastest (1.0x) | Good | Iteration |
| Normal | Fast (1.1x) | Better | Production |
| Karras | Slower (1.3x) | Best | Final renders |
| Exponential | Medium (1.2x) | Artistic | Creative work |
| SGM Uniform | Variable | Experimental | Research |

### Recommended Step Counts

Different schedulers work best with different step counts:

```bash
# Simple: 15-25 steps sufficient
asset-generator generate image --scheduler simple --steps 20

# Normal: 20-30 steps recommended
asset-generator generate image --scheduler normal --steps 25

# Karras: 30-50 steps for best results
asset-generator generate image --scheduler karras --steps 35

# Exponential: 25-40 steps
asset-generator generate image --scheduler exponential --steps 30
```

## API Integration

The scheduler parameter is passed directly to SwarmUI:

```go
// In cmd/generate.go and cmd/pipeline.go
req := &client.GenerationRequest{
    Prompt: prompt,
    Parameters: map[string]interface{}{
        "sampler":   "euler_a",
        "scheduler": "karras",  // Sent to SwarmUI API
        "steps":     30,
        // ... other parameters
    },
}
```

SwarmUI API expects the `scheduler` parameter in the generation request body.

## Troubleshooting

### Problem: Scheduler not affecting output

**Solution:**
Verify your SwarmUI version supports the scheduler parameter:
```bash
# Check if scheduler appears in generation info
asset-generator generate image \
  --prompt "test" \
  --scheduler karras \
  --format json \
  | jq '.metadata.scheduler'
```

### Problem: Unknown scheduler warning

**Solution:**
Use one of the supported schedulers: simple, normal, karras, exponential, sgm_uniform

```bash
# Correct usage
asset-generator generate image --scheduler karras

# Invalid scheduler (falls back to model default)
asset-generator generate image --scheduler invalid_name
```

### Problem: No visible difference between schedulers

**Solution:**
Increase step count and CFG scale to see scheduler impact:
```bash
asset-generator generate image \
  --prompt "detailed portrait" \
  --scheduler karras \
  --steps 40 \
  --cfg-scale 8.0
```

## Best Practices

### 1. Match Scheduler to Workflow Phase

```bash
# Phase 1: Concept exploration (simple)
asset-generator generate image --scheduler simple --steps 15

# Phase 2: Refinement (normal)
asset-generator generate image --scheduler normal --steps 25

# Phase 3: Final render (karras)
asset-generator generate image --scheduler karras --steps 40
```

### 2. Set Per-Project Defaults

Create project-specific configs:

```yaml
# project/.asset-generator-config.yaml
generate:
  scheduler: karras
  steps: 35
  cfg-scale: 8.0
```

### 3. Document Scheduler Choice

When sharing prompts or pipelines, include scheduler info:

```yaml
# Character renders - Karras scheduler for detail
assets:
  - name: characters
    metadata:
      notes: "Use --scheduler karras for best results"
```

### 4. Test Before Production

Always test scheduler selection on sample images before batch generation:

```bash
# Test single image first
asset-generator generate image \
  --prompt "test subject" \
  --scheduler karras \
  --seed 42 \
  --save-images

# Then run full pipeline
asset-generator pipeline \
  --file assets-spec.yaml \
  --scheduler karras
```

## Related Features

- **Samplers**: See [QUICKSTART.md](QUICKSTART.md) for sampler options
- **Skimmed CFG**: See [SKIMMED_CFG.md](SKIMMED_CFG.md) for quality optimization
- **Pipeline**: See [PIPELINE.md](PIPELINE.md) for batch workflows
- **Configuration**: See [QUICKSTART.md](QUICKSTART.md#configuration) for config file setup

## Summary

The scheduler selection feature provides:

- ✅ Full control over noise schedule
- ✅ Simple CLI interface (`--scheduler`)
- ✅ Config file support
- ✅ Pipeline integration
- ✅ Sensible default (simple)
- ✅ Five scheduler options
- ✅ Compatible with all other features

Choose the right scheduler for your workflow:
- **Quick iteration**: `simple`
- **Production work**: `normal`
- **Maximum quality**: `karras`
- **Artistic work**: `exponential`
- **Experimental**: `sgm_uniform`
