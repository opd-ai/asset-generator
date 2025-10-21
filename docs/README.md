# Asset Generator CLI - Documentation

> **Complete documentation for the Asset Generator command-line interface**

## 🚀 Getting Started

**New to Asset Generator?** Start here:

1. **[Quick Start Guide](QUICKSTART.md)** - Get up and running in 5 minutes
2. **[Command Reference](COMMANDS.md)** - All available commands
3. **[Main README](../README.md)** - Project overview and installation

## 📚 User Documentation

### Core Guides

- **[Quick Start Guide](QUICKSTART.md)** - Installation and basic usage
- **[Command Reference](COMMANDS.md)** - Comprehensive command documentation
- **[Integration Guide](INTEGRATE_PROMPT.md)** - Add asset-generator to your project
- **[Game Asset Pipeline Adapter](../adapter.md)** - Analyze game codebases and generate custom asset pipelines

### Feature Guides

- **[Pipeline Processing](PIPELINE.md)** - Automated batch generation workflows
- **[User Guide](USER_GUIDE.md)** - Complete feature guide with generation features, LoRA support, and filename templates
- **[Troubleshooting](TROUBLESHOOTING.md)** - Error resolution and common issues
- **[Postprocessing](POSTPROCESSING.md)** - Auto-crop, downscaling, and metadata stripping
- **[SVG Conversion](SVG_CONVERSION.md)** - Convert images to vector format
- **[Filename Templates](USER_GUIDE.md#filename-templates)** - Custom naming for downloaded images

### Reference Documentation

- **[Seed Behavior](SEED_BEHAVIOR.md)** - Understanding seed randomization
- **[Changelog](CHANGELOG.md)** - Version history and release notes

## 🔧 Developer Documentation

### Technical Guides

- **[Development Guide](DEVELOPMENT.md)** - Architecture, API integration, and contributing
- **[Project Summary](PROJECT_SUMMARY.md)** - High-level overview and statistics
- **[State File Sharing](STATE_FILE_SHARING.md)** - Cross-process session tracking

## 📖 Documentation by Task

### I want to...

#### Generate Images
- Basic generation → [Quick Start](QUICKSTART.md#generate-your-first-image)
- Advanced parameters → [User Guide](USER_GUIDE.md)
- Use LoRA models → [User Guide - LoRA Support](USER_GUIDE.md#lora-support)
- Batch processing → [Pipeline Processing](PIPELINE.md)

#### Process Images
- Download images → [User Guide - Filename Templates](USER_GUIDE.md#filename-templates)
- Remove borders → [Postprocessing: Auto-Crop](POSTPROCESSING.md#auto-crop)
- Resize images → [Postprocessing: Downscaling](POSTPROCESSING.md#downscaling)
- Convert to SVG → [SVG Conversion](SVG_CONVERSION.md)

#### Configure & Manage
- Initial setup → [Quick Start: Initial Setup](QUICKSTART.md#initial-setup)
- Configuration → [Quick Start: Configuration](QUICKSTART.md#configuration-examples)
- Check server status → [Commands: Status](COMMANDS.md#status-command)
- Cancel generation → [Commands: Cancel](COMMANDS.md#cancel-command)
- List models → [Quick Start: View Available Models](QUICKSTART.md#view-available-models)

#### Build Workflows
- Create pipelines → [Pipeline Processing](PIPELINE.md)
- Integrate into projects → [Integration Guide](INTEGRATE_PROMPT.md)
- Analyze games for assets → [Game Asset Adapter](../adapter.md)
- Reproducible generation → [Seed Behavior](SEED_BEHAVIOR.md)
- Automated scripts → [Quick Start: Scripting](QUICKSTART.md#scripting-examples)

#### Extend & Contribute
- Understand architecture → [Development Guide](DEVELOPMENT.md)
- API integration → [Development: SwarmUI API](DEVELOPMENT.md#swarmui-api-integration)
- Add features → [Development: Adding Features](DEVELOPMENT.md#adding-features)

## 🔍 Quick Reference

### Most Used Commands

```bash
# Generate an image
asset-generator generate image --prompt "your prompt"

# Generate with custom parameters
asset-generator generate image --prompt "detailed portrait" \
  --scheduler karras --steps 30 --width 768 --height 1024

# Process a pipeline
asset-generator pipeline --file pipeline.yaml

# Check server status
asset-generator status

# List available models
asset-generator models list
```

### Common Workflows

```bash
# High-quality generation with postprocessing
asset-generator generate image \
  --prompt "artwork" \
  --width 2048 --height 2048 \
  --scheduler karras --steps 40 \
  --save-images \
  --auto-crop \
  --downscale-width 1024

# Batch pipeline with LoRA
asset-generator pipeline \
  --file assets.yaml \
  --lora "style-lora:0.8" \
  --scheduler karras \
  --continue-on-error
```

## 📂 Documentation Structure

```
docs/
├── README.md                    ← You are here
├── QUICKSTART.md                Quick start guide
├── COMMANDS.md                  Command reference
├── INTEGRATE_PROMPT.md          Integration guide
├── CHANGELOG.md                 Version history
│
├── Feature Guides
│   ├── PIPELINE.md              Pipeline processing
│   ├── USER_GUIDE.md            Complete user features guide
│   ├── TROUBLESHOOTING.md       Error resolution guide
│   ├── POSTPROCESSING.md        Image processing
│   ├── SVG_CONVERSION.md        Vector conversion
│   └── COMMANDS.md              Complete command reference
│
├── Reference
│   └── SEED_BEHAVIOR.md         Seed handling
│
└── Developer Docs
    ├── DEVELOPMENT.md           Architecture guide
    ├── PROJECT_SUMMARY.md       Project overview
    └── STATE_FILE_SHARING.md    State management
```

## 🆘 Getting Help

### Troubleshooting

- **Connection issues** → [Quick Start: Troubleshooting](QUICKSTART.md#troubleshooting)
- **Generation failures** → [Pipeline: Troubleshooting](PIPELINE.md#troubleshooting)
- **Processing errors** → [Postprocessing: Troubleshooting](POSTPROCESSING.md#troubleshooting)
- **Command not found** → [Quick Start: Installation](QUICKSTART.md#installation)

### Additional Resources

- **GitHub Issues:** [Report bugs or request features](https://github.com/opd-ai/asset-generator/issues)
- **Examples:** See `examples/` directory in repository
- **Demo Scripts:** Check `demo-*.sh` files in repository root

## 📝 Documentation Standards

### For Users
- **Quick Start:** 5-minute setup and first generation
- **Guides:** Task-oriented with examples
- **Reference:** Comprehensive flag/option tables
- **Troubleshooting:** Common issues and solutions

### For Developers
- **Architecture:** Clean code principles and patterns
- **API:** SwarmUI integration details
- **Contributing:** Code style and submission process

## 🔗 External Links

- **Main Project:** [GitHub Repository](https://github.com/opd-ai/asset-generator)
- **SwarmUI:** [SwarmUI Documentation](https://github.com/mcmonkeyprojects/SwarmUI)
- **Releases:** [Version Downloads](https://github.com/opd-ai/asset-generator/releases)

---

**Last Updated:** October 10, 2025  
**Documentation Version:** Aligned with current `main` branch

**Navigation:** [🏠 Docs Home](README.md) | [📚 Quick Start](QUICKSTART.md) | [🔧 Commands](COMMANDS.md)
