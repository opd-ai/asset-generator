# LoRA Support - Quick Reference

## Quick Command Examples

```bash
# Single LoRA with default weight
asset-generator generate image --prompt "anime character" --lora "anime-style"

# Single LoRA with custom weight
asset-generator generate image --prompt "portrait" --lora "realistic-faces:0.8"

# Multiple LoRAs
asset-generator generate image \
  --prompt "cyberpunk city" \
  --lora "cyberpunk:1.0" --lora "neon-lights:0.7"

# Negative weight (remove style)
asset-generator generate image \
  --prompt "character" \
  --lora "realistic:1.0" --lora "cartoon:-0.5"
```

## Flag Reference

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--lora` | string | - | LoRA model (format: `name:weight` or `name`) |
| `--lora-weight` | float | - | Explicit weight (alternative to inline) |
| `--lora-default-weight` | string | "1.0" | Default weight when not specified |

## Weight Guidelines

| Weight | Effect |
|--------|--------|
| `-0.5 to 0.0` | Removes/reduces style elements |
| `0.0` | No effect |
| `0.5-0.7` | Subtle influence |
| `0.8-1.0` | Standard strength |
| `1.1-1.5` | Strong influence |
| `1.6-2.0` | Very strong (may overpower) |

## Format Options

### Inline Weight (Recommended)
```bash
--lora "model-name:0.8"
```

### Name Only (Uses Default)
```bash
--lora "model-name" --lora-default-weight 0.8
```

### Explicit Weights
```bash
--lora "model-1" --lora "model-2" --lora-weight 0.8 --lora-weight 0.6
```

## Configuration File

```yaml
generate:
  loras:
    - anime-style:0.9
    - detailed-faces:0.6
  lora-default-weight: 1.0
```

## Common Use Cases

### Anime Style
```bash
--lora "anime-style-v2:0.9" --lora "detailed-eyes:0.7"
```

### Photorealistic
```bash
--lora "realistic-faces:1.0" --lora "professional-photo:0.8"
```

### Fantasy Art
```bash
--lora "fantasy-art:0.9" --lora "detailed-background:0.6"
```

### Cyberpunk
```bash
--lora "cyberpunk:1.1" --lora "neon-lights:0.8"
```

## Integration Examples

### With Skimmed CFG
```bash
asset-generator generate image \
  --prompt "detailed portrait" \
  --lora "realistic-faces:0.8" \
  --skimmed-cfg --skimmed-cfg-scale 3.0
```

### With Postprocessing
```bash
asset-generator generate image \
  --prompt "high-res art" \
  --lora "detailed-art:0.9" \
  --width 2048 --height 2048 \
  --save-images --auto-crop --downscale-width 1024
```

### Batch Generation
```bash
asset-generator generate image \
  --prompt "character designs" \
  --lora "anime-style:0.9" \
  --batch 8 --save-images \
  --filename-template "char-{index}.png"
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| LoRA not working | Verify name matches exactly (case-sensitive) |
| Too strong | Reduce weight: 1.0 → 0.7 |
| Too weak | Increase weight: 0.8 → 1.2 |
| Conflicts | Use fewer LoRAs (2-3 max) |
| Unexpected style | Check LoRA-model compatibility |

## List Available LoRAs

```bash
# List all LoRAs
asset-generator models list --subtype LoRA

# Search for specific LoRA
asset-generator models list --subtype LoRA | grep "anime"
```

## See Also

- [LORA_SUPPORT.md](LORA_SUPPORT.md) - Complete documentation
- [QUICKSTART.md](QUICKSTART.md) - Getting started guide
- [SKIMMED_CFG_QUICKREF.md](SKIMMED_CFG_QUICKREF.md) - Skimmed CFG reference
