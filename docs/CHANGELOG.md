# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **Random seed by default for pipeline command**: The `--base-seed` flag now defaults to `-1` (random)
  instead of `42`. Both `0` and `-1` values trigger random seed generation. This provides more variety by 
  default while still allowing reproducibility by explicitly specifying a seed. The generated random seed 
  is displayed in the output for later reproduction.
- **Pipeline system refactored to be fully generic**: Major breaking change
  - Removed all tarot-specific code and data structures
  - Replaced with flexible, hierarchical asset group system
  - Metadata-driven prompt enhancement (metadata values appended to prompts)
  - Support for unlimited nesting via subgroups
  - Flexible directory structure configuration
  - Metadata cascading from parent to child groups
  - Generic seed offset and filename management
  - **Breaking**: Old tarot YAML format no longer supported
  - Migration guide provided in docs/PIPELINE_MIGRATION.md
  - Complete documentation in docs/GENERIC_PIPELINE.md
  - Example conversion in examples/tarot-deck-converted.yaml
  - New generic examples in examples/generic-pipeline.yaml
  - Demo script updated (demo-generic-pipeline.sh)

### Added
- **Server status command**: New `status` command to query SwarmUI server health and configuration
  - Real-time server connectivity verification
  - Response time measurement
  - Backend status and information display
  - Model availability and loading status
  - System information (GPU, memory, etc.)
  - Multiple output formats (table with ANSI colors, JSON, YAML)
  - Exit codes for automation (0=online, 1=offline)
  - Graceful handling of missing API endpoints
  - Comprehensive documentation in docs/STATUS_COMMAND.md
  - Quick reference guide in docs/STATUS_QUICKREF.md
  - Demo script (demo-status.sh) with usage examples
- **Skimmed CFG support**: Advanced sampling technique for improved generation quality and speed
  - `--skimmed-cfg` flag to enable Skimmed CFG (Distilled CFG)
  - `--skimmed-cfg-scale` for adjusting guidance scale (default: 3.0)
  - `--skimmed-cfg-start` for phase-specific application (0.0-1.0)
  - `--skimmed-cfg-end` for phase-specific application (0.0-1.0)
  - Available in both `generate image` and `pipeline` commands
  - Config file support for default Skimmed CFG settings
  - Comprehensive documentation in docs/SKIMMED_CFG.md
  - Quick reference guide in docs/SKIMMED_CFG_QUICKREF.md
- **Pipeline processing feature**: Native YAML pipeline file processing for batch generation
  - `pipeline` command for automated multi-asset generation workflows
  - Eliminates need for external shell scripts and yq dependencies
  - Supports tarot deck format with Major/Minor Arcana organization
  - Automatic directory structure creation for organized output
  - Reproducible generation using base seed + offsets
  - Progress tracking with detailed status updates
  - Dry-run mode for previewing pipeline without generating
  - Continue-on-error support for robust production pipelines
  - Style suffix and negative prompt application to all assets
  - Full postprocessing support (auto-crop, downscaling)
  - Cross-platform YAML parsing with gopkg.in/yaml.v3
  - Comprehensive documentation in docs/PIPELINE.md
- Initial CLI implementation with Cobra framework
- Image generation command with full parameter support
- Model listing and inspection commands
- Configuration management system with viper
- Support for multiple output formats (table, json, yaml)
- Comprehensive error handling and validation
- Signal handling for graceful shutdown
- Configuration precedence: flags > env vars > config file > defaults
- Detailed help documentation for all commands
- Unit tests with high coverage
- Makefile for build automation
- **Image download feature**: Automatically download generated images to local disk
- `--save-images` flag to enable image downloading
- `--output-dir` flag to specify download directory
- Comprehensive test suite for download functionality
- **SVG conversion feature**: Convert images to SVG format using two methods
  - Primitive method: Geometric shape approximation using fogleman/primitive
  - Gotrace method: Edge tracing vector conversion using potrace wrapper
- `convert svg` command with extensive customization options
- Support for 9 different shape modes (triangles, ellipses, beziers, etc.)
- Quality control via `--shapes`, `--mode`, and `--alpha` flags
- Comprehensive documentation and examples for SVG conversion
- Unit tests for converter package with 100% coverage
- **PNG Metadata Stripping**: Automatic removal of all PNG metadata for privacy and security
  - Mandatory, non-optional feature enforced on all PNG operations
  - Strips all ancillary chunks (tEXt, zTXt, iTXt, tIME, pHYs, gAMA, iCCP, etc.)
  - Preserves only critical chunks (IHDR, PLTE, IDAT, IEND)
  - Applied automatically during download, crop, and resize operations
  - Comprehensive test coverage for metadata stripping functionality
  - Full documentation in PNG_METADATA_STRIPPING.md

### Features
- Generate images using AI text-to-image APIs
- **Download and save generated images locally**
- **Automatic PNG metadata removal for all downloaded and processed images**
- **Convert images to SVG format with artistic and technical options**
- List and inspect available models
- Configure API endpoint and authentication
- Support for batch generation with automatic download
- Custom seeds for reproducible generation
- Negative prompts support
- Multiple sampling methods
- Verbose and quiet modes
- File output with timestamps
- Progress feedback for image downloads
- Graceful handling of partial download failures
- **Privacy protection through mandatory metadata stripping**

## [0.1.0] - 2025-10-07

### Added
- Initial release
- Basic CLI structure
- Core functionality for asset generation API integration
