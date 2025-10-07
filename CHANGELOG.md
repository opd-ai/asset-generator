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

### Features
- Generate images using SwarmUI text-to-image API
- List and inspect available models
- Configure API endpoint and authentication
- Support for batch generation
- Custom seeds for reproducible generation
- Negative prompts support
- Multiple sampling methods
- Verbose and quiet modes
- File output with timestamps

## [0.1.0] - 2025-10-07

### Added
- Initial release
- Basic CLI structure
- Core functionality for SwarmUI integration
