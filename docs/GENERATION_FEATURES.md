# Generation Features

This document covers advanced generation-time parameters that control quality, speed, and output characteristics.

## Table of Contents

- [Scheduler Selection](#scheduler-selection)
  - [Overview](#scheduler-overview)
  - [Available Schedulers](#available-schedulers)
  - [Usage Examples](#scheduler-usage)
  - [Quick Reference](#scheduler-quick-reference)
- [Skimmed CFG](#skimmed-cfg)
  - [Overview](#skimmed-cfg-overview)
  - [How It Works](#skimmed-cfg-how-it-works)
  - [Usage Examples](#skimmed-cfg-usage)
  - [Quick Reference](#skimmed-cfg-quick-reference)

---

## Scheduler Selection {#scheduler-selection}

### Overview {#scheduler-overview}

The scheduler (also known as noise schedule) controls how noise is added and removed during the diffusion process. Different schedulers can significantly impact generation quality, speed, and characteristics.

**Status**: ✅ **IMPLEMENTED**

The Asset Generator CLI supports selecting different schedulers via:
- `--scheduler` flag for `generate image` command
- `--scheduler` flag for `pipeline` command  
- `scheduler` config file setting
- Defaults to `simple` scheduler

### Available Schedulers {#available-schedulers}

| Scheduler | Description | Best For |
|-----------|-------------|----------|
| `simple` | Basic linear schedule (default) | General purpose, fast |
| `normal` | Standard scheduler with balanced noise | Most use cases, reliable |
| `karras` | Karras noise schedule | High-quality, detailed images |
| `exponential` | Exponential decay schedule | Smooth transitions, artistic |
| `sgm_uniform` | Uniform noise distribution | Experimental, specialized models |

### Usage {#scheduler-usage}

#### Command Line Examples

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

# Exponential scheduler for artistic results
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

**When to use:**
- Quick iterations and testing
- Fast results needed
- General purpose generation
- Learning and experimentation

**Characteristics:**
- Fastest processing
- Consistent results
- Good for most use cases

**Example:**
```bash
asset-generator generate image \
  --prompt "concept art, character design" \
  --scheduler simple \
  --steps 20
```

#### Normal

**When to use:**
- Production-quality assets
- Balanced speed and quality
- Professional workflows

**Characteristics:**
- Standard diffusion schedule
- Proven reliability
- Good detail preservation

**Example:**
```bash
asset-generator generate image \
  --prompt "product photography, professional lighting" \
  --scheduler normal \
  --steps 25
```

#### Karras

**When to use:**
- High-quality final outputs
- Fine details are critical
- Portrait and character work
- Maximum quality needed

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

#### Exponential

**When to use:**
- Artistic and creative work
- Smooth gradients
- Abstract compositions

**Characteristics:**
- Smoother transitions
- More "painted" look
- Good for artistic styles

**Example:**
```bash
asset-generator generate image \
  --prompt "abstract art, flowing energy, cosmic" \
  --scheduler exponential \
  --steps 30
```

#### SGM Uniform

**When to use:**
- Specialized models
- Experimental workflows
- Research purposes

**Characteristics:**
- Uniform noise distribution
- Model-specific behavior
- Experimental results

### Combined Workflows

#### Testing Different Schedulers

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
done
```

#### Scheduler + Sampler Combinations

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
```

### Performance Considerations

| Scheduler | Relative Speed | Quality | Best Use |
|-----------|----------------|---------|----------|
| Simple | Fastest (1.0x) | Good | Iteration |
| Normal | Fast (1.1x) | Better | Production |
| Karras | Slower (1.3x) | Best | Final renders |
| Exponential | Medium (1.2x) | Artistic | Creative work |
| SGM Uniform | Variable | Experimental | Research |

### Recommended Step Counts

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

---

### Scheduler Quick Reference {#scheduler-quick-reference}

#### One-Line Summary

Control noise schedule with `--scheduler` flag - use `simple` for speed, `karras` for quality.

#### Available Options

```
simple       - Fast, reliable (default)
normal       - Balanced quality
karras       - High-quality details  
exponential  - Smooth, artistic
sgm_uniform  - Specialized/experimental
```

#### Quick Commands

```bash
# Fast iteration (default)
asset-generator generate image --prompt "test" --scheduler simple

# High quality
asset-generator generate image --prompt "portrait" --scheduler karras --steps 35

# Pipeline with quality scheduler
asset-generator pipeline --file assets.yaml --scheduler karras
```

#### When to Use Each

| Task | Scheduler | Steps |
|------|-----------|-------|
| Quick testing | `simple` | 15-20 |
| Production work | `normal` | 20-30 |
| Final renders | `karras` | 30-50 |
| Artistic/creative | `exponential` | 25-40 |

#### Best Combinations

```bash
# Speed: Simple + Euler_a
--sampler euler_a --scheduler simple --steps 20

# Quality: Karras + DPM++
--sampler dpm_2 --scheduler karras --steps 35

# Balanced: Normal + Heun
--sampler heun --scheduler normal --steps 25
```

---

## Skimmed CFG {#skimmed-cfg}

### Overview {#skimmed-cfg-overview}

**Status**: ✅ **IMPLEMENTED**

Skimmed CFG (Classifier-Free Guidance) is an advanced sampling technique that improves image generation quality and speed by applying a more efficient guidance strategy during the denoising process. Also known as "Distilled CFG" or "Dynamic CFG."

### How It Works {#skimmed-cfg-how-it-works}

Traditional CFG requires computing both conditional and unconditional predictions at every sampling step. Skimmed CFG optimizes this by:

1. **Smart Guidance Application**: Applies full CFG guidance only when most beneficial
2. **Reduced Overhead**: Skips unnecessary computations during certain phases
3. **Dynamic Scaling**: Adjusts guidance strength based on the generation phase

### Benefits

- **Improved Quality**: Better adherence to prompts with enhanced coherence
- **Faster Generation**: Reduced computational overhead in some cases
- **Lower Scale Values**: Achieves similar or better results with lower CFG scales (typically 2.0-4.0 vs standard 7.0-8.0)
- **Fine-Grained Control**: Ability to apply guidance only during specific generation phases

### Usage {#skimmed-cfg-usage}

#### Basic Usage

Enable Skimmed CFG with default settings:

```bash
asset-generator generate image \
  --prompt "detailed fantasy landscape" \
  --skimmed-cfg
```

This uses the default scale of 3.0 and applies throughout the entire generation process.

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

**Phase guidelines**:
- **Early phase (0.0-0.3)**: Composition and structure formation
- **Middle phase (0.3-0.7)**: Detail development
- **Late phase (0.7-1.0)**: Refinement and final details

#### Combined with Standard CFG

You can use both standard CFG and Skimmed CFG together:

```bash
asset-generator generate image \
  --prompt "anime character portrait" \
  --cfg-scale 7.5 \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.0
```

#### Pipeline Usage

Apply Skimmed CFG to all generations in a pipeline:

```bash
asset-generator pipeline \
  --file assets-spec.yaml \
  --output-dir ./output \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.5 \
  --skimmed-cfg-start 0.1 \
  --skimmed-cfg-end 0.9
```

#### Configuration File

Set default Skimmed CFG options:

```yaml
# ~/.asset-generator/config.yaml
generate:
  model: stable-diffusion-xl
  steps: 20
  cfg-scale: 7.5
  
  # Skimmed CFG settings
  skimmed-cfg: true
  skimmed-cfg-scale: 3.0
  skimmed-cfg-start: 0.0
  skimmed-cfg-end: 1.0
```

### Practical Examples

#### Example 1: High-Quality Portrait

```bash
asset-generator generate image \
  --prompt "professional portrait photograph, studio lighting, sharp focus" \
  --width 768 \
  --height 1024 \
  --steps 30 \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.5 \
  --save-images
```

#### Example 2: Fast Concept Art

```bash
asset-generator generate image \
  --prompt "fantasy sword weapon design, concept art" \
  --width 1024 \
  --height 1024 \
  --steps 20 \
  --skimmed-cfg \
  --skimmed-cfg-scale 2.5 \
  --batch 4
```

#### Example 3: Batch with Postprocessing

```bash
asset-generator generate image \
  --prompt "minimalist logo design" \
  --width 1024 \
  --height 1024 \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.0 \
  --batch 10 \
  --save-images \
  --auto-crop \
  --downscale-width 512
```

#### Example 4: Mid-Phase Guidance Only

Apply Skimmed CFG only during detail development:

```bash
asset-generator generate image \
  --prompt "intricate mechanical device, detailed gears" \
  --skimmed-cfg \
  --skimmed-cfg-start 0.3 \
  --skimmed-cfg-end 0.7 \
  --steps 40
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

#### Detailed/Complex Scenes
```bash
--skimmed-cfg \
--skimmed-cfg-scale 4.0 \
--skimmed-cfg-start 0.1 \
--skimmed-cfg-end 0.9 \
--steps 40
```

### Troubleshooting

#### Images Look Different Than Expected

**Solution**: 
- Start with a scale of 3.0 and adjust incrementally
- Try applying only during mid-phase (0.3-0.7) first
- Compare results with and without Skimmed CFG enabled

#### No Visible Effect

**Solution**:
- Verify your model supports Skimmed CFG
- Ensure you're using a compatible SwarmUI version
- Try more pronounced scale differences (e.g., 2.0 vs 4.0)

#### Generation Failures

**Solution**:
- Verify your SwarmUI instance supports the feature
- Check that start < end for phase ranges
- Ensure scale values are positive numbers

### Model Compatibility

Not all models support Skimmed CFG. Test with:

```bash
asset-generator generate image \
  --prompt "test image" \
  --skimmed-cfg \
  --model your-model-name
```

Models known to work well with Skimmed CFG:
- Stable Diffusion XL variants
- Flux models
- Most modern diffusion models with CFG support

### Integration with Other Features

```bash
# With auto-crop and downscaling
asset-generator generate image \
  --prompt "character design" \
  --skimmed-cfg \
  --save-images \
  --auto-crop \
  --downscale-percentage 50

# With scheduler selection
asset-generator generate image \
  --prompt "detailed portrait" \
  --scheduler karras \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.0 \
  --steps 30
```

---

### Skimmed CFG Quick Reference {#skimmed-cfg-quick-reference}

#### Common Commands

```bash
# Enable with defaults
asset-generator generate image --prompt "..." --skimmed-cfg

# Custom scale (typically 2.0-4.0)
asset-generator generate image --prompt "..." --skimmed-cfg --skimmed-cfg-scale 3.0

# Phase-specific (0.0-1.0)
asset-generator generate image --prompt "..." \
  --skimmed-cfg --skimmed-cfg-start 0.2 --skimmed-cfg-end 0.8
```

#### Flag Reference

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--skimmed-cfg` | bool | false | Enable Skimmed CFG |
| `--skimmed-cfg-scale` | float | 3.0 | Scale value (lower than standard CFG) |
| `--skimmed-cfg-start` | float | 0.0 | Start phase (0.0 = beginning) |
| `--skimmed-cfg-end` | float | 1.0 | End phase (1.0 = end) |

#### Recommended Settings

| Use Case | Scale | Start | End | Steps |
|----------|-------|-------|-----|-------|
| **Photorealistic** | 3.5 | 0.0 | 1.0 | 30 |
| **Artistic/Stylized** | 2.5 | 0.2 | 0.8 | 25 |
| **Fast Iteration** | 2.0 | 0.0 | 1.0 | 15 |
| **High Detail** | 4.0 | 0.1 | 0.9 | 40 |

#### Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| No visible effect | Verify model compatibility, try scale 2.0 vs 4.0 |
| Unexpected results | Start with scale 3.0, try mid-phase only (0.3-0.7) |
| Generation fails | Check SwarmUI version, ensure start < end |

---

## See Also

- [LoRA Support](LORA_SUPPORT.md) - Fine-tune generation with LoRA models
- [Pipeline Processing](PIPELINE.md) - Batch generation workflows
- [Postprocessing](POSTPROCESSING.md) - Auto-crop, downscaling, metadata stripping
- [Seed Behavior](SEED_BEHAVIOR.md) - Reproducibility and random seed generation
- [Quick Start Guide](QUICKSTART.md) - Getting started with the CLI
