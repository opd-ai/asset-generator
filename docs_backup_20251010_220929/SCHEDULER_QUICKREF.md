# Scheduler Selection - Quick Reference

## One-Line Summary

Control noise schedule with `--scheduler` flag - use `simple` for speed, `karras` for quality.

## Available Schedulers

```
simple       - Fast, reliable (default)
normal       - Balanced quality
karras       - High-quality details
exponential  - Smooth, artistic
sgm_uniform  - Specialized/experimental
```

## Quick Examples

```bash
# Fast iteration (default)
asset-generator generate image --prompt "test" --scheduler simple

# High quality
asset-generator generate image --prompt "portrait" --scheduler karras --steps 35

# Pipeline with quality scheduler
asset-generator pipeline --file assets.yaml --scheduler karras
```

## Config File

```yaml
generate:
  scheduler: karras
  steps: 35
```

## Flag Reference

### Generate Command
```bash
--scheduler string   # simple, normal, karras, exponential, sgm_uniform (default: "simple")
```

### Pipeline Command
```bash
--scheduler string   # Applies to all assets in pipeline (default: "simple")
```

## When to Use Each

| Task | Scheduler | Steps |
|------|-----------|-------|
| Quick testing | `simple` | 15-20 |
| Production work | `normal` | 20-30 |
| Final renders | `karras` | 30-50 |
| Artistic/creative | `exponential` | 25-40 |

## Best Combinations

```bash
# Speed: Simple + Euler_a
--sampler euler_a --scheduler simple --steps 20

# Quality: Karras + DPM++
--sampler dpm_2 --scheduler karras --steps 35

# Balanced: Normal + Heun
--sampler heun --scheduler normal --steps 25
```

## API Parameter

SwarmUI API receives:
```json
{
  "scheduler": "karras",
  "sampler": "euler_a",
  "steps": 35
}
```

## Viper Config Binding

```go
viper.BindPFlag("generate.scheduler", generateImageCmd.Flags().Lookup("scheduler"))
```

## Related Features

- Samplers: Control sampling method (`--sampler`)
- Skimmed CFG: Quality optimization (`--skimmed-cfg`)
- Steps: Number of inference iterations (`--steps`)

## Full Documentation

See [SCHEDULER_FEATURE.md](SCHEDULER_FEATURE.md) for:
- Detailed scheduler comparisons
- Performance benchmarks
- Advanced usage patterns
- Troubleshooting guide
