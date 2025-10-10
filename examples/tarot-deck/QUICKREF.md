# Tarot Deck Pipeline - Quick Reference

## Commands

### Native Pipeline Command (Recommended)
```bash
# Generate all 78 cards with pipeline command
asset-generator pipeline --file tarot-spec.yaml

# Preview before generating
asset-generator pipeline --file tarot-spec.yaml --dry-run

# Custom output directory and seed
asset-generator pipeline --file tarot-spec.yaml \
  --output-dir ./my-deck --base-seed 42

# With postprocessing
asset-generator pipeline --file tarot-spec.yaml \
  --auto-crop --downscale-width 1024
```

### Shell Script Wrappers (Legacy)
```bash
./quick-demo.sh                  # Test with 5 sample cards (~5 min)
./generate-tarot-deck.sh         # Wrapper: generates all 78 cards
./generate-card-back.sh          # Generate 4 card back designs
./post-process-deck.sh           # Create multiple formats
./package-for-print.sh           # Package for distribution
```

### Custom Parameters
```bash
# Pipeline command with custom parameters
asset-generator pipeline --file tarot-spec.yaml \
  --output-dir OUTPUT_DIR --base-seed BASE_SEED

# Wrapper script (backward compatible)
./generate-tarot-deck.sh OUTPUT_DIR BASE_SEED
./generate-tarot-deck.sh OUTPUT_DIR BASE_SEED --dry-run
```

## File Structure

```
examples/tarot-deck/
├── tarot-spec.yaml              # 78 card definitions (YAML)
├── generate-tarot-deck.sh       # Pipeline wrapper script
├── generate-card-back.sh        # Card back generator
├── post-process-deck.sh         # Format converter
├── package-for-print.sh         # Print packager
├── quick-demo.sh                # Quick test (5 cards)
├── README.md                    # Complete guide
├── WORKFLOW.md                  # Visual workflow
├── DEMONSTRATION.md             # Technical summary
└── QUICKREF.md                  # This file
```

## Output Structure

```
tarot-deck-output/               # Generated cards
├── major-arcana/                # 22 cards
├── minor-arcana/                # 56 cards
│   ├── wands/
│   ├── cups/
│   ├── swords/
│   └── pentacles/
└── card-backs/                  # 4 designs

tarot-deck-processed/            # Post-processed
├── print-ready/
├── web-optimized/
├── mobile-optimized/
├── svg-versions/
└── thumbnails/

tarot-deck-packages/             # Distribution
├── *.zip
├── PRINT_SPECIFICATIONS.txt
└── MANIFEST.txt
```

## Key Parameters

| Parameter | Value | Purpose |
|-----------|-------|---------|
| Width | 768px | Card width |
| Height | 1344px | Card height (4:7 ratio) |
| Steps | 40 | Quality iterations |
| CFG Scale | 7.5 | Prompt adherence |
| Base Seed | 42 | Reproducibility |

## Card Counts

| Category | Count | Notes |
|----------|-------|-------|
| Major Arcana | 22 | Cards 0-21 |
| Wands | 14 | Fire suit |
| Cups | 14 | Water suit |
| Swords | 14 | Air suit |
| Pentacles | 14 | Earth suit |
| **Total** | **78** | Complete deck |

## Timing Estimates

| Operation | Time | Notes |
|-----------|------|-------|
| Quick Demo | 5 min | 5 sample cards |
| Major Arcana | 25 min | 22 cards |
| Minor Arcana | 65 min | 56 cards |
| Card Backs | 5 min | 4 designs |
| Post-Process | 10 min | All formats |
| Packaging | <1 min | Archives |
| **Total** | **~105 min** | Complete pipeline |

## Customization Quick Tips

### Change Art Style
Use pipeline command flags:
```bash
asset-generator pipeline --file tarot-spec.yaml \
  --style-suffix "watercolor painting, soft edges, artistic"
```

Or edit wrapper script `generate-tarot-deck.sh`:
```bash
STYLE_SUFFIX="your style here, artistic keywords"
```

### Change Card Size
Use pipeline command flags:
```bash
asset-generator pipeline --file tarot-spec.yaml \
  --width 750 --height 1050
```

Or edit wrapper script `generate-tarot-deck.sh`:
```bash
WIDTH=750    # Your width
HEIGHT=1050  # Your height
```

### Adjust Quality
Use pipeline command flags:
```bash
asset-generator pipeline --file tarot-spec.yaml \
  --steps 50 --cfg-scale 8.0
```

Or edit wrapper script `generate-tarot-deck.sh`:
```bash
STEPS=50         # Higher = better quality, slower
CFG_SCALE=8.0    # Higher = more prompt adherence
```

### Modify Prompts
Edit `tarot-spec.yaml`:
```yaml
- name: Your Card
  prompt: "your description here"
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Command not found | `chmod +x *.sh` or use pipeline directly |
| Asset-generator error | `asset-generator config view` |
| Slow generation | Use `--steps 25` flag |
| Quality issues | Use `--steps 50` flag |
| Out of disk space | Check with `df -h` |
| Pipeline fails | Try `--continue-on-error` flag |

## Prerequisites Checklist

- [ ] asset-generator CLI installed (`go install` or binary)
- [ ] Asset generation service running (http://localhost:7801)
- [ ] Sufficient disk space (~500 MB minimum)
- [ ] GPU with adequate VRAM (recommended)
- [ ] For shell scripts: `chmod +x *.sh`

**No external dependencies needed!** The native pipeline command uses built-in Go YAML parsing.

## Quick Links

- **Pipeline Documentation**: [../../docs/PIPELINE.md](../../docs/PIPELINE.md)
- **Pipeline Quick Reference**: [../../docs/PIPELINE_QUICKREF.md](../../docs/PIPELINE_QUICKREF.md)
- **Full Tarot Documentation**: [README.md](README.md)
- **Visual Workflow**: [WORKFLOW.md](WORKFLOW.md)
- **Technical Details**: [DEMONSTRATION.md](DEMONSTRATION.md)
- **Main CLI Docs**: [../../README.md](../../README.md)
- **Pipeline Guide**: [../../GENERATE_PIPELINE.md](../../GENERATE_PIPELINE.md)

## One-Liner Examples

```bash
# Complete pipeline in one go
./generate-tarot-deck.sh && ./generate-card-back.sh && ./post-process-deck.sh && ./package-for-print.sh

# Custom output and seed
./generate-tarot-deck.sh ./my-deck 1234

# Test quickly before full run
./quick-demo.sh && ls -lh tarot-deck-demo/

# Regenerate just Major Arcana (edit script first)
./generate-tarot-deck.sh ./major-only 42

# Process existing cards to different formats
./post-process-deck.sh ./tarot-deck-output ./processed
```

## File Format Guide

| Format | Size | Use Case |
|--------|------|----------|
| Print-Ready | 768×1344 | Professional printing |
| Web-Optimized | ≤1024px | Website galleries |
| Mobile-Optimized | ≤512px | Mobile apps |
| Thumbnails | ≤256px | Preview galleries |
| SVG | Vector | Scalable designs |

## Seed Strategy

| Card Type | Seed Formula | Example |
|-----------|--------------|---------|
| Major Arcana | BASE + number | 42 + 0 = 42 |
| Wands | BASE + 100 + n | 42 + 100 + 0 = 142 |
| Cups | BASE + 120 + n | 42 + 120 + 0 = 162 |
| Swords | BASE + 140 + n | 42 + 140 + 0 = 182 |
| Pentacles | BASE + 160 + n | 42 + 160 + 0 = 202 |
| Card Backs | BASE + 9999 | 42 + 9999 = 10041 |

## Common Asset Generator Commands

```bash
# Test connection
asset-generator models list

# View configuration
asset-generator config view

# Generate single card manually
asset-generator generate image \
  --prompt "tarot card art, description" \
  --width 768 --height 1344 \
  --steps 40 --cfg-scale 7.5 \
  --seed 42 \
  --save-images \
  --output-dir ./output \
  --filename-template "card-name.png"

# Downscale image
asset-generator downscale \
  --input ./input.png \
  --output ./output.png \
  --max-dimension 1024

# Convert to SVG
asset-generator convert svg \
  --input ./input.png \
  --output ./output.svg \
  --mode primitive \
  --shapes 200
```

## Success Indicators

✓ Scripts are executable (green in `ls -l`)  
✓ asset-generator responds to `models list`  
✓ Quick demo completes in ~5 minutes  
✓ Generated cards have consistent quality  
✓ File sizes are reasonable (200-500 KB each)  
✓ Directory structure is organized  
✓ Post-processing creates all formats  
✓ Packaging produces ZIP archives  

---

**Quick Start**: Run `./quick-demo.sh` to generate 5 sample cards and verify everything works!

**Full Pipeline**: Run `./generate-tarot-deck.sh` to generate the complete 78-card deck.

**Need Help?** See [README.md](README.md) for comprehensive documentation.
