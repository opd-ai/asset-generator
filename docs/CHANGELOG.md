# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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
