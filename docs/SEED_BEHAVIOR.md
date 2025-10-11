# Seed Behavior Reference

## Overview

This document describes how seed values are handled across different commands in `asset-generator`.

## Pipeline Command

### Random Seed Triggers

The following values will generate a random seed:

| Value | Behavior |
|-------|----------|
| *not specified* | Random seed generated (default) |
| `0` | Random seed generated |
| `-1` | Random seed generated |

### Explicit Seeds

Any positive integer (1, 2, 3, ... ) will be used as-is for reproducible generation.

### Examples

```bash
# All of these generate random seeds:
asset-generator pipeline --file assets.yaml
asset-generator pipeline --file assets.yaml --base-seed 0
asset-generator pipeline --file assets.yaml --base-seed -1

# These use explicit seeds:
asset-generator pipeline --file assets.yaml --base-seed 1
asset-generator pipeline --file assets.yaml --base-seed 42
asset-generator pipeline --file assets.yaml --base-seed 12345
```

### Output

When a random seed is generated, it's displayed in the output:

```
Loading pipeline file: assets.yaml
Pipeline loaded: 11 total assets
...

Generated random base seed: 1760137005907907614

Output directory: ./pipeline-output
Base seed: 1760137005907907614
...
```

## Generate Image Command

The `generate image` command follows the same pattern:

### Random Seed Triggers

| Value | Behavior |
|-------|----------|
| *not specified* | Random seed (default -1) |
| `-1` | Random seed |

**Note:** The generate command doesn't support `0` as a random trigger (only `-1`), but when omitted, 
it defaults to `-1` which triggers random seed generation by the API.

### Examples

```bash
# Random seed (default)
asset-generator generate image --prompt "fantasy landscape" --save-images

# Explicit seed for reproducibility
asset-generator generate image --prompt "fantasy landscape" --seed 42 --save-images
```

## Why Both 0 and -1?

Supporting both `0` and `-1` as random seed triggers provides flexibility:

- **`-1`**: Traditional convention in many AI/ML tools for "random"
- **`0`**: Natural choice for "no seed specified" or "zero seed"
- **Consistency**: Some users may expect `0` to mean random, others `-1`
- **Safety**: Prevents accidental use of seed `0` when random was intended

## Seed Calculation in Pipelines

When using explicit seeds in pipeline mode, individual asset seeds are calculated as:

```
asset_seed = base_seed + group_seed_offset + asset_index
```

Example with `--base-seed 1000`:

```yaml
assets:
  - name: Group A
    seed_offset: 0
    assets:
      - id: a1  # Seed: 1000 + 0 + 0 = 1000
      - id: a2  # Seed: 1000 + 0 + 1 = 1001
  
  - name: Group B
    seed_offset: 100
    assets:
      - id: b1  # Seed: 1000 + 100 + 0 = 1100
      - id: b2  # Seed: 1000 + 100 + 1 = 1101
```

## Best Practices

### For Exploration and Iteration

Use random seeds to explore different variations:

```bash
# Let the tool generate random seeds
asset-generator pipeline --file game-assets.yaml
```

### For Production and Reproducibility

Always specify explicit seeds:

```bash
# Lock in specific seed for consistent results
asset-generator pipeline --file game-assets.yaml --base-seed 12345
```

### For CI/CD Pipelines

Use explicit seeds in automated workflows:

```yaml
# .github/workflows/assets.yml
- name: Generate Assets
  run: |
    asset-generator pipeline \
      --file assets.yaml \
      --base-seed 42 \
      --output-dir ./generated
```

### For Version Control

Document seeds in your repository:

```markdown
# README.md
## Asset Generation

Generate production assets with:
```bash
asset-generator pipeline --file assets.yaml --base-seed 12345
```
```

## Troubleshooting

### Issue: Pipeline generates different results each time

**Cause:** Using default random seed

**Solution:** Specify an explicit seed:
```bash
asset-generator pipeline --file assets.yaml --base-seed 42
```

### Issue: Want to force random seed but seed is configured in config file

**Cause:** Config file has explicit seed set

**Solution:** Override with command-line flag:
```bash
asset-generator pipeline --file assets.yaml --base-seed 0
# or
asset-generator pipeline --file assets.yaml --base-seed -1
```

### Issue: Seed 0 produces unexpected results

This shouldn't happen! Seeds of `0` now trigger random generation. If you need 
to use seed `0` specifically, you'll need to use seed `1` or another value.

## Technical Details

### Random Seed Generation

When a random seed is needed, it's generated using:

```go
seed := time.Now().UnixNano()
```

This produces a 64-bit integer based on the current nanosecond timestamp, 
providing high-quality randomness for typical use cases.

### Seed Range

Valid seed values:
- **Random triggers:** `0` or `-1`
- **Explicit seeds:** Any positive 64-bit integer (1 to 9,223,372,036,854,775,807)

## See Also

- [Pipeline Documentation](PIPELINE.md)
- [Generate Command](../README.md#image-generation)
