# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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

### Features
- Generate images using AI text-to-image APIs
- **Download and save generated images locally**
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

## [0.1.0] - 2025-10-07

### Added
- Initial release
- Basic CLI structure
- Core functionality for asset generation API integration
