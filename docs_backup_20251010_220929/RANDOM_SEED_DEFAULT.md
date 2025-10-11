# Random Seed Behavior

## Overview

Both the `generate image` and `pipeline` commands now default to using random seeds when `--seed` or `--base-seed` is not explicitly specified. This provides more variety by default while maintaining the ability to reproduce results when needed.

## Generate Image Command

**Default:** `--seed -1` (random)

When generating images without specifying a seed:

```bash
# Uses random seed
asset-generator generate image --prompt "fantasy landscape" --save-images

# Specify seed for reproducibility
asset-generator generate image --prompt "fantasy landscape" --seed 42 --save-images
```

The random seed behavior is handled by the SwarmUI API. When seed is `-1` or omitted, the API generates a unique random seed for each image.

## Pipeline Command

**Default:** `--base-seed -1` (random)

When running pipelines without specifying a base seed, or when explicitly setting it to `0` or `-1`:

```bash
# Uses random base seed (displayed in output)
asset-generator pipeline --file assets.yaml

# Explicitly request random seed
asset-generator pipeline --file assets.yaml --base-seed 0
asset-generator pipeline --file assets.yaml --base-seed -1

# Output shows:
# Generated random base seed: 1760137005907907614

# Specify base seed for reproducibility
asset-generator pipeline --file assets.yaml --base-seed 42
```

### How It Works

1. If `--base-seed` is not specified, set to `0`, or set to `-1`, the pipeline command generates a random seed using `time.Now().UnixNano()`
2. The generated seed is displayed in the output so you can reproduce results later
3. Individual asset seeds are calculated as: `base_seed + group_seed_offset + asset_index`

### Example Output

```
Loading pipeline file: examples/generic-pipeline.yaml
Pipeline loaded: 11 total assets
  - Hero Characters: 3 assets
  - Enemy Characters: 2 assets
  - UI Elements: 2 assets

Generated random base seed: 1760137005907907614

Output directory: ./pipeline-output
Base seed: 1760137005907907614
Dimensions: 768x1344
Steps: 40, CFG Scale: 7.5
```

## Reproducing Results

### With Generate Command

The seed used for each image is included in the response metadata and can be tracked using filename templates:

```bash
# Save images with seed in filename
asset-generator generate image \
  --prompt "fantasy castle" \
  --save-images \
  --filename-template "castle-{seed}.png"

# Regenerate with specific seed
asset-generator generate image \
  --prompt "fantasy castle" \
  --seed 1760137005907907614 \
  --save-images
```

### With Pipeline Command

When the random seed is displayed, you can reproduce the exact same pipeline run:

```bash
# First run (random seed displayed)
asset-generator pipeline --file assets.yaml
# Output: Generated random base seed: 1760137005907907614

# Reproduce exact results
asset-generator pipeline --file assets.yaml --base-seed 1760137005907907614
```

## Benefits

### Random by Default
- **Variety:** Each run produces different results, great for exploration
- **Avoiding Repetition:** No need to manually change seeds between runs
- **Natural Workflow:** Matches expected behavior for creative tools

### Reproducible When Needed
- **Version Control:** Commit specific seed values for consistent results
- **Debugging:** Reproduce issues by using the same seed
- **Production:** Lock in seeds for final assets

## Configuration Files

You can set default seed values in your config file:

```yaml
# ~/.config/asset-generator/config.yaml
generate:
  seed: 42  # Override default for generate command

pipeline:
  base_seed: 42  # Override default for pipeline command
```

However, command-line flags always take precedence over config file values.

## Migration Notes

### Previous Behavior

Previously, the `pipeline` command defaulted to `--base-seed 42`, which meant:
- Every pipeline run produced identical results unless seed was explicitly changed
- Users had to manually specify different seeds for variation

### Current Behavior

Now, both commands default to random seeds:
- `generate image`: `--seed -1` (always was this way)
- `pipeline`: `--base-seed -1` (changed from 42 to -1)

### Updating Scripts

If you have scripts that depend on deterministic output:

```bash
# Old (implicit seed of 42)
asset-generator pipeline --file assets.yaml

# New (explicit seed for same behavior)
asset-generator pipeline --file assets.yaml --base-seed 42
```

## Best Practices

### For Exploration and Iteration
```bash
# Let it use random seeds
asset-generator pipeline --file game-sprites.yaml
```

### For Production and CI/CD
```bash
# Lock in a specific seed
asset-generator pipeline --file game-sprites.yaml --base-seed 12345
```

### For Version Control
Document the seed in your pipeline files or README:

```yaml
# game-sprites.yaml
# Production seed: 12345
# Use: asset-generator pipeline --file game-sprites.yaml --base-seed 12345
```

## See Also

- [Pipeline Documentation](PIPELINE.md)
- [Generic Pipeline System](GENERIC_PIPELINE.md)
- [Filename Templates](FILENAME_TEMPLATES.md)
