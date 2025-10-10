# Tarot Deck Pipeline - Quick Reference

## Commands

### Essential Pipeline
```bash
./quick-demo.sh                  # Test with 5 sample cards (~5 min)
./generate-tarot-deck.sh         # Generate all 78 cards (~90 min)
./generate-card-back.sh          # Generate 4 card back designs
./post-process-deck.sh           # Create multiple formats
./package-for-print.sh           # Package for distribution
```

### Custom Parameters
```bash
./generate-tarot-deck.sh OUTPUT_DIR BASE_SEED
./generate-card-back.sh OUTPUT_DIR SEED
./post-process-deck.sh INPUT_DIR OUTPUT_DIR
./package-for-print.sh INPUT_DIR OUTPUT_DIR
```

## File Structure

```
examples/tarot-deck/
├── tarot-spec.yaml              # 78 card definitions
├── generate-tarot-deck.sh       # Main generator
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
Edit `generate-tarot-deck.sh`:
```bash
STYLE_SUFFIX="your style here, artistic keywords"
```

### Change Card Size
Edit `generate-tarot-deck.sh`:
```bash
WIDTH=750    # Your width
HEIGHT=1050  # Your height
```

### Adjust Quality
Edit `generate-tarot-deck.sh`:
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
| Command not found | `chmod +x *.sh` |
| Asset-generator error | `asset-generator config view` |
| Slow generation | Reduce STEPS to 25-30 |
| Quality issues | Increase STEPS to 50 |
| Out of disk space | Check with `df -h` |
| yq not found | Install: see Prerequisites |

## Prerequisites Checklist

- [ ] asset-generator CLI installed
- [ ] Asset generation service running (http://localhost:7801)
- [ ] yq installed (mikefarah's Go version, NOT python-yq)
- [ ] Sufficient disk space (~500 MB minimum)
- [ ] GPU with adequate VRAM (recommended)
- [ ] Scripts executable (`chmod +x *.sh`)

**Install yq (correct version):**
```bash
wget -qO /tmp/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64
sudo mv /tmp/yq /usr/local/bin/yq
sudo chmod +x /usr/local/bin/yq
yq --version  # Should show mikefarah's version
```

## Quick Links

- **Full Documentation**: [README.md](README.md)
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
