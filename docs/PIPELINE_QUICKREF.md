# Pipeline Command - Quick Reference

## Basic Usage

```bash
asset-generator pipeline --file pipeline.yaml
```

## Common Commands

```bash
# Preview pipeline (dry run)
asset-generator pipeline --file deck.yaml --dry-run

# Generate with custom output directory
asset-generator pipeline --file deck.yaml --output-dir ./my-output

# High quality with postprocessing
asset-generator pipeline --file deck.yaml \
  --steps 50 \
  --auto-crop \
  --downscale-width 1024

# Continue on errors
asset-generator pipeline --file deck.yaml --continue-on-error

# Add style to all prompts
asset-generator pipeline --file deck.yaml \
  --style-suffix "detailed, professional quality, rich colors"
```

## Key Flags

### Required
- `--file` - Pipeline YAML file path

### Generation
- `--output-dir` - Output directory (default: `./pipeline-output`)
- `--base-seed` - Base seed for reproducibility (default: `-1` for random, 0 also triggers random)
- `--steps` - Inference steps (default: `40`)
- `--width` - Image width (default: `768`)
- `--height` - Image height (default: `1344`)
- `--cfg-scale` - CFG scale/guidance (default: `7.5`)
- `--model` - Model to use for generation

### Prompt Enhancement
- `--style-suffix` - Append to all prompts
- `--negative-prompt` - Negative prompt for all

### Control
- `--dry-run` - Preview without generating
- `--continue-on-error` - Don't stop on failures

### Postprocessing
- `--auto-crop` - Remove whitespace borders
- `--downscale-width` - Downscale to width
- `--downscale-height` - Downscale to height
- `--downscale-filter` - Filter: lanczos, bilinear, nearest

## Pipeline File Format

```yaml
major_arcana:
  - number: 0
    name: The Fool
    prompt: "detailed card description"
  
  - number: 1
    name: The Magician
    prompt: "another description"

minor_arcana:
  wands:
    suit_element: fire
    suit_color: red
    cards:
      - rank: Ace
        prompt: "card description"
      - rank: Two
        prompt: "card description"
  
  cups:
    suit_element: water
    cards:
      - rank: Ace
        prompt: "card description"
```

## Output Structure

```
output-dir/
├── major-arcana/
│   ├── 00-the_fool.png
│   ├── 01-the_magician.png
│   └── ...
└── minor-arcana/
    ├── wands/
    │   ├── 01-ace_of_wands.png
    │   └── ...
    ├── cups/
    ├── swords/
    └── pentacles/
```

## Seed Calculation

### Major Arcana
```
seed = base_seed + card_number
```

### Minor Arcana
```
seed = base_seed + 100 + suit_offset + card_index

Suit offsets:
  - wands: 0
  - cups: 20
  - swords: 40
  - pentacles: 60
```

## Examples

### Basic Generation
```bash
asset-generator pipeline --file tarot-spec.yaml
```

### Preview First
```bash
asset-generator pipeline --file tarot-spec.yaml --dry-run
```

### Production Quality
```bash
asset-generator pipeline --file tarot-spec.yaml \
  --base-seed 42 \
  --steps 50 \
  --width 1536 \
  --height 2688 \
  --style-suffix "masterpiece, detailed, professional" \
  --auto-crop \
  --downscale-width 768 \
  --continue-on-error
```

### Quick Test
```bash
asset-generator pipeline --file test.yaml \
  --steps 20 \
  --width 512 \
  --height 768
```

## Troubleshooting

### Preview First
Always preview with `--dry-run` before generating:
```bash
asset-generator pipeline --file deck.yaml --dry-run
```

### Check Configuration
```bash
asset-generator config view
```

### Test Connection
```bash
asset-generator models list
```

### Enable Verbose Output
```bash
asset-generator pipeline --file deck.yaml --verbose
```

### Continue on Errors
For large pipelines:
```bash
asset-generator pipeline --file deck.yaml --continue-on-error
```

## See Also

- [Full Documentation](PIPELINE.md)
- [Tarot Deck Example](../examples/tarot-deck/)
- [Main README](../README.md)
