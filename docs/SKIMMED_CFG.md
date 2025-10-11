# Skimmed CFG Feature

**Feature Status**: ✅ **IMPLEMENTED**

## Overview

Skimmed CFG (Classifier-Free Guidance) is an advanced sampling technique that improves image generation quality and speed by applying a more efficient guidance strategy during the denoising process. Also known as "Distilled CFG" or "Dynamic CFG," this technique can produce higher quality results with lower computational overhead.

## How It Works

Traditional CFG requires computing both conditional and unconditional predictions at every sampling step. Skimmed CFG optimizes this by:

1. **Smart Guidance Application**: Applies full CFG guidance only when most beneficial
2. **Reduced Overhead**: Skips unnecessary computations during certain phases
3. **Dynamic Scaling**: Adjusts guidance strength based on the generation phase

## Benefits

- **Improved Quality**: Better adherence to prompts with enhanced coherence
- **Faster Generation**: Reduced computational overhead in some cases
- **Lower Scale Values**: Achieves similar or better results with lower CFG scales (typically 2.0-4.0 vs standard 7.0-8.0)
- **Fine-Grained Control**: Ability to apply guidance only during specific generation phases

## Usage

### Basic Usage

Enable Skimmed CFG with default settings:

```bash
asset-generator generate image \
  --prompt "detailed fantasy landscape" \
  --skimmed-cfg
```

This uses the default scale of 3.0 and applies throughout the entire generation process.

### Custom Scale

Adjust the Skimmed CFG scale (typically lower than standard CFG):

```bash
asset-generator generate image \
  --prompt "portrait of a wizard" \
  --skimmed-cfg \
  --skimmed-cfg-scale 2.5
```

**Recommended scale values**: 2.0-4.0 (compared to standard CFG of 7.0-8.0)

### Phase-Specific Application

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

### Combined with Standard CFG

You can use both standard CFG and Skimmed CFG together:

```bash
asset-generator generate image \
  --prompt "anime character portrait" \
  --cfg-scale 7.5 \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.0
```

The interaction between standard and Skimmed CFG depends on your model's implementation.

### Pipeline Usage

Apply Skimmed CFG to all generations in a pipeline:

```bash
asset-generator pipeline \
  --file tarot-spec.yaml \
  --output-dir ./tarot-output \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.5 \
  --skimmed-cfg-start 0.1 \
  --skimmed-cfg-end 0.9
```

## Configuration File

Set default Skimmed CFG options in your config file:

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

## API Parameters

The CLI translates flags to these SwarmUI API parameters:

| CLI Flag | API Parameter | Type | Description |
|----------|--------------|------|-------------|
| `--skimmed-cfg` | `skimmedcfg` | boolean | Enable/disable Skimmed CFG |
| `--skimmed-cfg-scale` | `skimmedcfgscale` | float | Scale value for guidance |
| `--skimmed-cfg-start` | `skimmedcfgstart` | float | Start percentage (0.0-1.0) |
| `--skimmed-cfg-end` | `skimmedcfgend` | float | End percentage (0.0-1.0) |

## Examples

### Example 1: High-Quality Portrait

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

### Example 2: Fast Concept Art

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

### Example 3: Batch with Postprocessing

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

### Example 4: Mid-Phase Guidance Only

Apply Skimmed CFG only during detail development:

```bash
asset-generator generate image \
  --prompt "intricate mechanical device, detailed gears" \
  --skimmed-cfg \
  --skimmed-cfg-start 0.3 \
  --skimmed-cfg-end 0.7 \
  --steps 40
```

## Recommended Settings by Use Case

### Photorealistic Images
```bash
--skimmed-cfg \
--skimmed-cfg-scale 3.5 \
--skimmed-cfg-start 0.0 \
--skimmed-cfg-end 1.0
```

### Artistic/Stylized Images
```bash
--skimmed-cfg \
--skimmed-cfg-scale 2.5 \
--skimmed-cfg-start 0.2 \
--skimmed-cfg-end 0.8
```

### Fast Iteration/Concepts
```bash
--skimmed-cfg \
--skimmed-cfg-scale 2.0 \
--steps 15
```

### Detailed/Complex Scenes
```bash
--skimmed-cfg \
--skimmed-cfg-scale 4.0 \
--skimmed-cfg-start 0.1 \
--skimmed-cfg-end 0.9 \
--steps 40
```

## Troubleshooting

### Images Look Different Than Expected

**Problem**: Results differ significantly from standard CFG generation.

**Solution**: 
- Start with a scale of 3.0 and adjust incrementally
- Try applying only during mid-phase (0.3-0.7) first
- Compare results with and without Skimmed CFG enabled

### No Visible Effect

**Problem**: Enabling Skimmed CFG doesn't seem to change results.

**Solution**:
- Verify your model supports Skimmed CFG (check model documentation)
- Ensure you're using a compatible SwarmUI version
- Try more pronounced scale differences (e.g., 2.0 vs 4.0)

### Generation Failures

**Problem**: Errors when using Skimmed CFG.

**Solution**:
- Verify your SwarmUI instance supports the feature
- Check that start < end for phase ranges
- Ensure scale values are positive numbers

## Model Compatibility

Not all models support Skimmed CFG. Check your model's documentation or test with:

```bash
# Test basic compatibility
asset-generator generate image \
  --prompt "test image" \
  --skimmed-cfg \
  --model your-model-name
```

Models known to work well with Skimmed CFG:
- Stable Diffusion XL variants
- Flux models
- Most modern diffusion models with CFG support

## Performance Considerations

- **Speed**: May be slightly faster or slower depending on implementation
- **Memory**: Similar memory usage to standard CFG
- **Quality**: Generally equal or better quality with appropriate settings

## Integration with Other Features

Skimmed CFG works seamlessly with other CLI features:

```bash
# With auto-crop and downscaling
asset-generator generate image \
  --prompt "character design" \
  --skimmed-cfg \
  --save-images \
  --auto-crop \
  --downscale-percentage 50

# With custom filenames
asset-generator generate image \
  --prompt "landscape scene" \
  --skimmed-cfg \
  --batch 5 \
  --save-images \
  --filename-template "landscape-{index}-cfg{skimmed-cfg-scale}.png"

# With WebSocket progress
asset-generator generate image \
  --prompt "detailed artwork" \
  --skimmed-cfg \
  --websocket
```

---

## Quick Reference

### Common Commands

```bash
# Enable with defaults
asset-generator generate image --prompt "..." --skimmed-cfg

# Custom scale (typically 2.0-4.0)
asset-generator generate image --prompt "..." --skimmed-cfg --skimmed-cfg-scale 3.0

# Phase-specific (0.0-1.0)
asset-generator generate image --prompt "..." \
  --skimmed-cfg --skimmed-cfg-start 0.2 --skimmed-cfg-end 0.8
```

### Flag Reference

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--skimmed-cfg` | bool | false | Enable Skimmed CFG |
| `--skimmed-cfg-scale` | float | 3.0 | Scale value (lower than standard CFG) |
| `--skimmed-cfg-start` | float | 0.0 | Start phase (0.0 = beginning) |
| `--skimmed-cfg-end` | float | 1.0 | End phase (1.0 = end) |

### Recommended Settings by Use Case

| Use Case | Scale | Start | End | Steps |
|----------|-------|-------|-----|-------|
| **Photorealistic** | 3.5 | 0.0 | 1.0 | 30 |
| **Artistic/Stylized** | 2.5 | 0.2 | 0.8 | 25 |
| **Fast Iteration** | 2.0 | 0.0 | 1.0 | 15 |
| **High Detail** | 4.0 | 0.1 | 0.9 | 40 |

### Config File Example

```yaml
generate:
  skimmed-cfg: true
  skimmed-cfg-scale: 3.0
  skimmed-cfg-start: 0.0
  skimmed-cfg-end: 1.0
```

### Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| No visible effect | Verify model compatibility, try scale 2.0 vs 4.0 |
| Unexpected results | Start with scale 3.0, try mid-phase only (0.3-0.7) |
| Generation fails | Check SwarmUI version, ensure start < end |

## See Also

- [Generation Parameters](../README.md#generation-parameters)
- [Pipeline Command](PIPELINE.md)
- [Configuration Guide](QUICKSTART.md#configuration)
- [SwarmUI API Documentation](API.md)

## Technical Notes

The feature is implemented by passing additional parameters to the SwarmUI API:

```go
// When --skimmed-cfg is enabled
req.Parameters["skimmedcfg"] = true
req.Parameters["skimmedcfgscale"] = generateSkimmedCFGScale

// Optional phase control
if generateSkimmedCFGStart != 0.0 {
    req.Parameters["skimmedcfgstart"] = generateSkimmedCFGStart
}
if generateSkimmedCFGEnd != 1.0 {
    req.Parameters["skimmedcfgend"] = generateSkimmedCFGEnd
}
```

The parameters are only included when Skimmed CFG is explicitly enabled to avoid unnecessary API payload overhead.
