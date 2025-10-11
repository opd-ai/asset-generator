# Skimmed CFG - Quick Reference

## What is it?
Advanced sampling technique for improved image quality and faster generation.

## Basic Usage

```bash
# Enable with defaults
--skimmed-cfg

# Custom scale (typically 2.0-4.0)
--skimmed-cfg --skimmed-cfg-scale 3.0

# Phase-specific (0.0-1.0)
--skimmed-cfg --skimmed-cfg-start 0.2 --skimmed-cfg-end 0.8
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--skimmed-cfg` | bool | false | Enable Skimmed CFG |
| `--skimmed-cfg-scale` | float | 3.0 | Scale value (lower than standard CFG) |
| `--skimmed-cfg-start` | float | 0.0 | Start phase (0.0 = beginning) |
| `--skimmed-cfg-end` | float | 1.0 | End phase (1.0 = end) |

## Recommended Settings

### Photorealistic
```bash
--skimmed-cfg --skimmed-cfg-scale 3.5
```

### Artistic/Stylized
```bash
--skimmed-cfg --skimmed-cfg-scale 2.5 --skimmed-cfg-start 0.2 --skimmed-cfg-end 0.8
```

### Fast Iteration
```bash
--skimmed-cfg --skimmed-cfg-scale 2.0 --steps 15
```

### High Detail
```bash
--skimmed-cfg --skimmed-cfg-scale 4.0 --steps 40
```

## Config File

```yaml
generate:
  skimmed-cfg: true
  skimmed-cfg-scale: 3.0
  skimmed-cfg-start: 0.0
  skimmed-cfg-end: 1.0
```

## Examples

### Basic Generation
```bash
asset-generator generate image \
  --prompt "detailed fantasy landscape" \
  --skimmed-cfg
```

### With Pipeline
```bash
asset-generator pipeline \
  --file spec.yaml \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.5
```

### Combined with Postprocessing
```bash
asset-generator generate image \
  --prompt "character portrait" \
  --skimmed-cfg \
  --save-images \
  --auto-crop \
  --downscale-width 1024
```

## Key Differences from Standard CFG

| Aspect | Standard CFG | Skimmed CFG |
|--------|--------------|-------------|
| Scale | 7.0-8.0 | 2.0-4.0 |
| Speed | Baseline | Potentially faster |
| Quality | Good | Equal or better |
| Phase Control | No | Yes (start/end) |

## Phase Guidelines

- **0.0-0.3**: Early (composition, structure)
- **0.3-0.7**: Middle (details develop)
- **0.7-1.0**: Late (refinement)

## Tips

✅ **Do:**
- Start with scale 3.0 and adjust
- Use lower scales than standard CFG
- Experiment with phase ranges
- Check model compatibility

❌ **Don't:**
- Use high scales (>5.0)
- Set start >= end
- Expect all models to support it

## Troubleshooting

**No visible effect?**
- Verify model supports Skimmed CFG
- Try more pronounced scale differences

**Errors?**
- Check SwarmUI compatibility
- Ensure start < end
- Use positive scale values

## See Also

- [Full Documentation](SKIMMED_CFG.md)
- [Generation Parameters](../README.md#generation-parameters)
- [Pipeline Guide](PIPELINE.md)
