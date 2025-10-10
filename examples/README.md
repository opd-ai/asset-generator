# Asset Generator Examples

This directory contains complete, production-ready examples demonstrating the full capabilities of the `asset-generator` CLI.

## Available Examples

### 🎴 [Tarot Deck Generation Pipeline](tarot-deck/)

A comprehensive example demonstrating professional asset generation workflow by creating a complete 78-card tarot deck.

**Features Demonstrated:**
- ✅ Batch generation of 78+ unique assets
- ✅ YAML-based specifications
- ✅ Organized multi-level directory structure
- ✅ Custom filename templates with sorting
- ✅ Seed-based reproducibility
- ✅ Negative prompts for quality control
- ✅ Consistent styling across large asset sets
- ✅ Post-processing pipeline (downscale, SVG, thumbnails)
- ✅ Print production packaging
- ✅ Shell script automation

**What's Included:**
- Complete 78-card tarot deck specification
- Automated generation scripts
- Card back design generator
- Multi-format post-processing
- Print production packaging
- Quick demo mode (5 sample cards)
- Comprehensive documentation

**Use Cases:**
- Game asset production (card games, collectibles)
- Print-on-demand products
- Digital collectibles and NFTs
- Mobile and web applications
- Commercial deck production

**Quick Start:**
```bash
cd tarot-deck/
./quick-demo.sh              # Generate 5 sample cards (~5 min)
./generate-tarot-deck.sh     # Generate complete deck (~90 min)
./post-process-deck.sh       # Create multiple formats
./package-for-print.sh       # Package for distribution
```

**Learn More:** [tarot-deck/README.md](tarot-deck/README.md)

---

## Example Categories

### 🎮 Game Asset Pipelines
- **Tarot Deck** - Complete card game asset generation

### 🖼️ Batch Generation
- **Tarot Deck** - 78+ unique assets with consistent styling

### 📐 Post-Processing
- **Tarot Deck** - Multi-format conversion (print, web, mobile, SVG)

### 🔧 Advanced Workflows
- **Tarot Deck** - YAML specifications, shell automation, packaging

---

## Coming Soon

Additional examples in development:

- **Pixel Art Sprite Sheets** - Character animations and variations
- **UI Element Sets** - Complete game UI asset library
- **Background Parallax Layers** - Multi-layer scrolling backgrounds
- **Icon Collections** - Themed icon sets with variations
- **Texture Atlases** - Seamless texture generation and packing

---

## Using These Examples

### As Learning Resources
Each example includes:
- Complete, working code
- Detailed documentation
- Step-by-step workflows
- Troubleshooting guides
- Customization examples

### As Templates
Examples are designed to be:
- Easily customizable for your needs
- Production-ready out of the box
- Well-documented for modification
- Extensible with your own features

### As Reference
Examples demonstrate:
- Best practices
- CLI feature usage
- Pipeline patterns
- Common workflows
- Performance optimization

---

## Requirements

All examples require:
- `asset-generator` CLI installed and configured
- Asset generation service running (default: http://localhost:7801)
- POSIX-compliant shell (bash, dash, zsh)
- Basic utilities (wget, zip, etc.)

**Note:** Previous examples may have required `yq` for YAML parsing, but the new native `pipeline` command eliminates this dependency. Shell scripts in examples are now simple wrappers around the pipeline command.

See individual example READMEs for specific requirements.

---

## Contributing Examples

Want to contribute an example? Great! Examples should:

1. **Be Complete** - Fully working, end-to-end pipeline
2. **Be Documented** - README with usage, customization, troubleshooting
3. **Demonstrate Features** - Showcase CLI capabilities effectively
4. **Be Practical** - Solve real-world use cases
5. **Be Maintainable** - Clean code, clear structure

Example structure:
```
examples/your-example/
├── README.md              # Complete documentation
├── specification.yaml     # Asset definitions (if applicable)
├── generate-*.sh          # Generation scripts
├── post-process-*.sh      # Post-processing scripts (if applicable)
└── WORKFLOW.md            # Detailed workflow documentation (optional)
```

---

## Support

For issues with examples:
1. Check the example's README troubleshooting section
2. Verify asset-generator configuration: `asset-generator config view`
3. Test with smaller sample sizes first
4. Review logs for specific error messages

For general CLI support, see the main [README.md](../README.md).

---

## Example Statistics

| Example | Assets | Time | Complexity | Best For |
|---------|--------|------|------------|----------|
| Tarot Deck | 78+ cards | 90 min | Advanced | Card games, print production, learning |

---

## Quick Navigation

- **[Main Documentation](../README.md)** - CLI reference
- **[Quick Start Guide](../QUICKSTART.md)** - Getting started
- **[Pipeline Guide](../GENERATE_PIPELINE.md)** - General pipeline patterns
- **[Development Guide](../DEVELOPMENT.md)** - Contributing to CLI

---

**Ready to create something amazing?** Start with the [Tarot Deck example](tarot-deck/) to see the full power of automated asset generation! 🎨✨
