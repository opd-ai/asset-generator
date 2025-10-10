# Tarot Deck Pipeline Workflow

## Complete Pipeline Visualization

```
┌─────────────────────────────────────────────────────────────────┐
│                   TAROT DECK GENERATION PIPELINE                │
└─────────────────────────────────────────────────────────────────┘

STEP 1: SPECIFICATION
┌────────────────────┐
│  tarot-spec.yaml   │  ← Complete deck definition (78 cards)
│                    │    - 22 Major Arcana with detailed prompts
│  • Major Arcana    │    - 56 Minor Arcana (4 suits × 14 cards)
│  • Minor Arcana    │    - Symbolic descriptions
│    - Wands         │    - Traditional tarot imagery
│    - Cups          │
│    - Swords        │
│    - Pentacles     │
└────────────────────┘
         │
         ▼

STEP 2: GENERATION
┌────────────────────────────────────────────────────────────────┐
│  generate-tarot-deck.sh [output_dir] [base_seed]              │
├────────────────────────────────────────────────────────────────┤
│  Generates all 78 cards with consistent styling               │
│                                                                │
│  Parameters:                                                   │
│    • Dimensions: 768×1344px (tarot proportion)                │
│    • Steps: 40 (high quality)                                 │
│    • CFG Scale: 7.5 (balanced)                                │
│    • Seeds: BASE_SEED + card_number (reproducible)            │
│                                                                │
│  Output: tarot-deck-output/                                   │
│    ├── major-arcana/ (22 cards)                              │
│    └── minor-arcana/                                          │
│        ├── wands/ (14 cards)                                  │
│        ├── cups/ (14 cards)                                   │
│        ├── swords/ (14 cards)                                 │
│        └── pentacles/ (14 cards)                              │
│                                                                │
│  Time: ~60-90 minutes on typical GPU                          │
└────────────────────────────────────────────────────────────────┘
         │
         ▼

STEP 3: CARD BACKS
┌────────────────────────────────────────────────────────────────┐
│  generate-card-back.sh [output_dir] [seed]                    │
├────────────────────────────────────────────────────────────────┤
│  Generates 4 decorative card back designs                     │
│                                                                │
│  Designs:                                                      │
│    1. Mandala with celestial symbols (blue/gold)              │
│    2. Victorian ornate border (purple/gold)                   │
│    3. Sacred geometry (indigo/silver)                         │
│    4. Art nouveau mystical eye (teal/gold)                    │
│                                                                │
│  Output: tarot-deck-output/card-backs/                        │
│    ├── card-back-design-1.png                                 │
│    ├── card-back-design-2.png                                 │
│    ├── card-back-design-3.png                                 │
│    └── card-back-design-4.png                                 │
│                                                                │
│  Time: ~5 minutes                                             │
└────────────────────────────────────────────────────────────────┘
         │
         ▼

STEP 4: POST-PROCESSING
┌────────────────────────────────────────────────────────────────┐
│  post-process-deck.sh [input_dir] [output_dir]                │
├────────────────────────────────────────────────────────────────┤
│  Creates multiple output formats for different use cases      │
│                                                                │
│  Processes:                                                    │
│    1. Print-Ready      → 768×1344px (original)                │
│    2. Web-Optimized    → max 1024px                           │
│    3. Mobile-Optimized → max 512px                            │
│    4. SVG Versions     → vector format (sample)               │
│    5. Thumbnails       → max 256px                            │
│                                                                │
│  Output: tarot-deck-processed/                                │
│    ├── print-ready/                                           │
│    ├── web-optimized/                                         │
│    ├── mobile-optimized/                                      │
│    ├── svg-versions/                                          │
│    └── thumbnails/                                            │
│                                                                │
│  Time: ~5-10 minutes                                          │
└────────────────────────────────────────────────────────────────┘
         │
         ▼

STEP 5: PACKAGING
┌────────────────────────────────────────────────────────────────┐
│  package-for-print.sh [input_dir] [output_dir]                │
├────────────────────────────────────────────────────────────────┤
│  Creates organized archives for distribution                  │
│                                                                │
│  Packages:                                                     │
│    • Complete deck (78 cards)                                 │
│    • Major Arcana only (22 cards)                             │
│    • Minor Arcana only (56 cards)                             │
│    • Individual suits (4 × 14 cards)                          │
│    • Card backs (4 designs)                                   │
│                                                                │
│  Documentation:                                                │
│    • PRINT_SPECIFICATIONS.txt                                 │
│    • MANIFEST.txt                                             │
│                                                                │
│  Output: tarot-deck-packages/                                 │
│    ├── tarot-deck-complete-TIMESTAMP.zip                      │
│    ├── tarot-major-arcana-TIMESTAMP.zip                       │
│    ├── tarot-minor-arcana-TIMESTAMP.zip                       │
│    ├── tarot-suit-wands-TIMESTAMP.zip                         │
│    ├── tarot-suit-cups-TIMESTAMP.zip                          │
│    ├── tarot-suit-swords-TIMESTAMP.zip                        │
│    ├── tarot-suit-pentacles-TIMESTAMP.zip                     │
│    ├── tarot-card-backs-TIMESTAMP.zip                         │
│    ├── PRINT_SPECIFICATIONS.txt                               │
│    └── MANIFEST.txt                                           │
│                                                                │
│  Ready for: Print shops, web deployment, distribution         │
└────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════

QUICK START COMMANDS:

# Full pipeline (end-to-end)
./generate-tarot-deck.sh
./generate-card-back.sh
./post-process-deck.sh
./package-for-print.sh

# Quick demo (5 sample cards only)
./quick-demo.sh

# Custom output location
./generate-tarot-deck.sh ./my-custom-deck 1234

═══════════════════════════════════════════════════════════════════

TIMELINE ESTIMATE:

Stage                   Time Estimate    Progress
────────────────────────────────────────────────────
Specification           < 1 minute       ▓
Generation (78 cards)   60-90 minutes    ▓▓▓▓▓▓▓▓▓░
Card Backs (4 designs)  5 minutes        ▓
Post-Processing         5-10 minutes     ▓
Packaging               < 1 minute       ▓
────────────────────────────────────────────────────
Total                   ~70-105 minutes

═══════════════════════════════════════════════════════════════════

FEATURES DEMONSTRATED:

✓ Batch Generation         - 78+ unique assets automated
✓ YAML Specifications      - Structured, maintainable definitions
✓ Custom Filenames         - Organized, sortable naming
✓ Seed Reproducibility     - Regenerate specific cards
✓ Negative Prompts         - Quality control
✓ Style Consistency        - Unified aesthetic
✓ Directory Organization   - Clean structure
✓ Multi-Format Output      - Print, web, mobile, SVG
✓ Post-Processing Pipeline - Downscale, convert, optimize
✓ Print Production Ready   - Packaged with specifications

═══════════════════════════════════════════════════════════════════

CUSTOMIZATION POINTS:

1. Art Style (generate-tarot-deck.sh)
   • Modify STYLE_SUFFIX for different aesthetics
   • Examples: pixel art, watercolor, minimalist, gothic

2. Card Dimensions (generate-tarot-deck.sh)
   • Change WIDTH/HEIGHT for different card sizes
   • Example: 750×1050 for poker-sized cards

3. Quality Settings (generate-tarot-deck.sh)
   • Adjust STEPS (20-60) for speed vs quality
   • Adjust CFG_SCALE (5-12) for prompt adherence

4. Prompts (tarot-spec.yaml)
   • Customize individual card descriptions
   • Add custom cards or variations

5. Post-Processing (post-process-deck.sh)
   • Change max-dimension values
   • Enable/disable specific formats
   • Adjust SVG conversion parameters

═══════════════════════════════════════════════════════════════════
```

## Integration Examples

### Game Development (Unity)
```bash
# Generate assets
./generate-tarot-deck.sh ./Assets/Resources/TarotCards
./post-process-deck.sh ./Assets/Resources/TarotCards ./Assets/Resources/TarotCards-Optimized

# Use mobile-optimized in game
cp -r tarot-deck-processed/mobile-optimized/* ./Assets/Resources/GameCards/
```

### Web Application (React/Vue/Angular)
```bash
# Generate and optimize for web
./generate-tarot-deck.sh
./post-process-deck.sh

# Deploy web-optimized version
cp -r tarot-deck-processed/web-optimized/* ./public/assets/tarot-cards/
cp -r tarot-deck-processed/thumbnails/* ./public/assets/tarot-thumbs/
```

### Mobile App (iOS/Android)
```bash
# Generate with mobile-first approach
./generate-tarot-deck.sh
./post-process-deck.sh

# Use mobile-optimized assets
cp -r tarot-deck-processed/mobile-optimized/* ./app/src/main/res/drawable-xhdpi/
```

### Print Production
```bash
# Full pipeline with packaging
./generate-tarot-deck.sh
./generate-card-back.sh
./post-process-deck.sh
./package-for-print.sh

# Send to print shop
# Use: tarot-deck-packages/tarot-deck-complete-*.zip
# Include: PRINT_SPECIFICATIONS.txt
```

## Troubleshooting

### Pipeline fails during generation
- Check asset-generator connection: `asset-generator models list`
- Verify sufficient disk space: `df -h`
- Test with quick-demo.sh first

### Quality issues
- Increase STEPS to 50-60
- Adjust CFG_SCALE (try 8-9)
- Review prompts in tarot-spec.yaml

### Slow generation
- Reduce STEPS to 25-30 (faster, slight quality loss)
- Ensure GPU acceleration is working
- Run overnight for full deck

### Post-processing errors
- Verify input directory exists
- Check asset-generator has conversion tools installed
- Test individual operations manually

## Advanced Workflows

### Version Control
```bash
# Tag specific generation
./generate-tarot-deck.sh ./decks/v1.0 42
./generate-tarot-deck.sh ./decks/v1.1 43

# Compare versions
diff -r ./decks/v1.0 ./decks/v1.1
```

### Batch Variations
```bash
# Generate multiple complete decks with different styles
for style in "watercolor" "pixel_art" "gothic" "minimalist"; do
    # Modify STYLE_SUFFIX in script per style
    ./generate-tarot-deck.sh ./decks/${style} $((42 + RANDOM % 1000))
done
```

### Selective Regeneration
```bash
# Regenerate just Major Arcana
# Edit generate-tarot-deck.sh to skip Minor Arcana loop
./generate-tarot-deck.sh ./tarot-deck-output 42

# Regenerate specific card
asset-generator generate image \
  --prompt "..." \
  --seed 42 \
  --save-images \
  --output-dir ./tarot-deck-output/major-arcana \
  --filename-template "00-the_fool.png"
```

## Quality Assurance Checklist

Before distribution:

- [ ] All 78 cards generated successfully
- [ ] No corrupted or blank images
- [ ] Consistent style across all cards
- [ ] Proper aspect ratio maintained
- [ ] Filenames correct and sortable
- [ ] Card backs align with card fronts
- [ ] Print specifications documented
- [ ] Quality acceptable for intended use
- [ ] Color accuracy verified
- [ ] File sizes appropriate
- [ ] Packaging complete
- [ ] Documentation included

## Support & Resources

- Main documentation: ../../README.md
- Pipeline guide: ../../GENERATE_PIPELINE.md
- Filename templates: ../../docs/FILENAME_TEMPLATES.md
- Image download: ../../docs/IMAGE_DOWNLOAD.md
- Asset-generator CLI: ../../QUICKSTART.md
