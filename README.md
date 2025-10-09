# Asset Generator CLI

A powerful command-line interface for interacting with AI asset generation APIs. Generate assets, manage models, and configure your asset generation workflow with ease.

## Features

- 🎨 **Asset Generation**: Generate images using text-to-image models
- 💾 **Image Download**: Automatically download and save generated images locally
- 📦 **Model Management**: List and inspect available models
- ⚙️ **Configuration**: Easy configuration management with multiple sources
- 📊 **Multiple Output Formats**: Table, JSON, and YAML output support
- 🔧 **Flexible Parameters**: Configure via flags, environment variables, or config file
- 🚀 **Performance**: Fast, efficient CLI built with Go

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/opd-ai/asset-generator.git
cd asset-generator

# Build and install
make install
```

### Download Binary

Download the latest release from the [releases page](https://github.com/opd-ai/asset-generator/releases).

## Quick Start

### 1. Initialize Configuration

```bash
asset-generator config init
```

### 2. Set Your API Endpoint

```bash
asset-generator config set api-url http://localhost:7801
```

### 3. Generate Your First Image

```bash
asset-generator generate image --prompt "a beautiful sunset over mountains"
```

## Usage

### Asset Generation

Generate images with various parameters:

```bash
# Basic generation
asset-generator generate image --prompt "a futuristic cityscape"

# Advanced generation with custom parameters
asset-generator generate image \
  --prompt "astronaut riding a horse" \
  --width 1024 \
  --length 1024 \
  --steps 30 \
  --cfg-scale 7.5 \
  --sampler euler_a

# Save output to file
asset-generator generate image \
  --prompt "cat wearing sunglasses" \
  --output result.json \
  --format json

# Download and save generated images locally
asset-generator generate image \
  --prompt "beautiful mountain landscape" \
  --save-images \
  --output-dir ./my-images

# Batch generation with image download
asset-generator generate image \
  --prompt "beautiful landscape" \
  --batch 4 \
  --save-images \
  --output-dir ./landscapes
```

### Model Management

List and inspect available models:

```bash
# List all models
asset-generator models list

# List models as JSON
asset-generator models list --format json

# Get details about a specific model
asset-generator models get stable-diffusion-xl
```

### Configuration

Manage your CLI configuration:

```bash
# View current configuration
asset-generator config view

# Set configuration values
asset-generator config set api-url https://api.example.com
asset-generator config set api-key your-api-key-here

# Get a specific configuration value
asset-generator config get api-url
```

## Configuration

Configuration can be provided through multiple sources with the following precedence:

1. **Command-line flags** (highest priority)
2. **Environment variables** (prefixed with `ASSET_GENERATOR_`)
3. **Configuration file** (searches multiple locations)
4. **Default values** (lowest priority)

### Configuration File

The application searches for `config.yaml` in the following locations (in order of precedence):

1. `./config/config.yaml` - Current directory (highest precedence)
2. `~/.asset-generator/config.yaml` - User's home directory
3. `/etc/asset-generator/config.yaml` - System-wide configuration (lowest precedence)

You can also specify a custom config file location using the `--config` flag, which takes highest precedence among configuration files.

See [`config/example-config.yaml`](config/example-config.yaml) for a template configuration file.

```yaml
api-url: http://localhost:7801
api-key: your-api-key
format: table
verbose: false
```

### Environment Variables

```bash
export ASSET_GENERATOR_API_URL=http://localhost:7801
export ASSET_GENERATOR_API_KEY=your-api-key
export ASSET_GENERATOR_FORMAT=json
```

## Global Flags

Available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--api-url` | | Asset generation API base URL | `http://localhost:7801` |
| `--api-key` | | Asset generation API key | |
| `--format` | `-f` | Output format (table, json, yaml) | `table` |
| `--output` | `-o` | Write output to file | |
| `--quiet` | `-q` | Quiet mode (errors only) | `false` |
| `--verbose` | `-v` | Verbose output | `false` |
| `--config` | | Config file path | `~/.asset-generator/config.yaml` |

## Generation Parameters

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--prompt` | `-p` | Generation prompt (required) | |
| `--model` | | Model to use | |
| `--width` | `-w` | Image width | `512` |
| `--length` | `-l` | Image length (height) | `512` |
| `--height` | | Image height (alias for --length) | `512` |
| `--steps` | | Inference steps | `20` |
| `--cfg-scale` | | CFG scale (guidance) | `7.5` |
| `--sampler` | | Sampling method | `euler_a` |
| `--batch` | `-b` | Number of images to generate | `1` |
| `--seed` | | Random seed (-1 for random) | `-1` |
| `--negative-prompt` | `-n` | Negative prompt | |
| `--websocket` | | Use WebSocket for real-time progress updates | `false` |
| `--save-images` | | Download and save generated images to local disk | `false` |
| `--output-dir` | | Directory to save downloaded images | `.` (current directory) |

> **Note:** The `--length` flag is used for the vertical dimension (height) for API compatibility with SwarmUI. Both `--length` and `--height` are supported as aliases.

## Image Download Feature

The `--save-images` flag enables automatic downloading of generated images from the server to your local disk. This is particularly useful when:

- You want to work with generated images locally
- You're generating multiple images in batch mode
- You need to preserve images beyond the server's retention period
- You want to organize images in specific directories

### Usage

```bash
# Download a single generated image
asset-generator generate image \
  --prompt "cyberpunk cityscape" \
  --save-images

# Save to a specific directory
asset-generator generate image \
  --prompt "fantasy landscape" \
  --save-images \
  --output-dir ./generated-art

# Batch generation with download
asset-generator generate image \
  --prompt "abstract art" \
  --batch 5 \
  --save-images \
  --output-dir ./batch-output
```

### Behavior

- Images are downloaded immediately after generation completes
- The original filename from the server is preserved
- Directory is created automatically if it doesn't exist
- Progress feedback shows each downloaded file
- Local file paths are added to the metadata output
- Partial failures are handled gracefully (some images may succeed even if others fail)

### Output

When `--save-images` is enabled, you'll see output like:

```
Generating image with prompt: cyberpunk cityscape
Downloading generated images...
  [1/1] Saved: ./generated-art/cyberpunk cityscape-1234567890.png
✓ Generation completed successfully (1 image)
```

## Examples

### Example 1: Basic Image Generation

```bash
asset-generator generate image --prompt "a serene lake at sunset"
```

### Example 2: High-Quality Generation

```bash
asset-generator generate image \
  --prompt "professional portrait photo of a scientist" \
  --width 1024 \
  --length 1024 \
  --steps 50 \
  --cfg-scale 8.0
```

### Example 3: Batch Generation with Seed

```bash
asset-generator generate image \
  --prompt "fantasy castle in the clouds" \
  --batch 4 \
  --seed 42 \
  --output results.json
```

### Example 4: Pipeline with JSON Output

```bash
# Generate and extract image paths
asset-generator generate image \
  --prompt "cyberpunk street scene" \
  --format json | jq '.image_paths[]'
```

## Development

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile)

### Building

```bash
# Build binary
make build

# Run tests
make test

# Generate coverage report
make coverage

# Format code
make fmt

# Run linter
make lint
```

### Project Structure

```
asset-generator/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and global setup
│   ├── generate.go        # Generation commands
│   ├── models.go          # Model management commands
│   └── config.go          # Configuration commands
├── pkg/                   # Public packages
│   ├── client/            # Asset generation API client
│   └── output/            # Output formatters
├── internal/              # Private packages
│   └── config/            # Configuration validation
├── main.go               # Application entry point
├── go.mod                # Go module file
├── Makefile              # Build automation
└── README.md             # This file
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific test
go test -v ./pkg/client/...
```

## Architecture

### Clean Architecture Principles

The CLI follows clean architecture with clear separation of concerns:

- **cmd/**: CLI command definitions and user interaction
- **pkg/client/**: Asset generation API client (reusable library)
- **pkg/output/**: Output formatting utilities
- **internal/config/**: Configuration validation logic

### API Client

The asset generation API client is designed to be reusable and can be imported by other Go projects:

```go
import "github.com/opd-ai/asset-generator/pkg/client"

// Create client
config := &client.Config{
    BaseURL: "http://localhost:7801",
    APIKey:  "your-api-key",
}
assetClient, err := client.NewAssetClient(config)

// Generate image
req := &client.GenerationRequest{
    Prompt: "a beautiful landscape",
    Parameters: map[string]interface{}{
        "width":  1024,
        "height": 1024,
        "steps":  30,
        "images": 4,  // Batch size - number of images to generate
    },
}
result, err := assetClient.GenerateImage(context.Background(), req)

// Download generated images
savedPaths, err := assetClient.DownloadImages(context.Background(), result.ImagePaths, "./output")
if err != nil {
    log.Fatalf("Failed to download images: %v", err)
}
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards

- Follow Go best practices and idioms
- Write tests for new features
- Update documentation as needed
- Run `make fmt` and `make lint` before committing

## Troubleshooting

### Connection Issues

If you're having trouble connecting to the asset generation service:

```bash
# Check your API URL
asset-generator config get api-url

# Test with verbose output
asset-generator models list -v

# Verify the service is running
curl http://localhost:7801/API/ListModels
```

### Configuration Issues

```bash
# View current configuration
asset-generator config view

# Reinitialize configuration
rm -rf ~/.asset-generator
asset-generator config init
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Links

- [SwarmUI Documentation](https://github.com/mcmonkeyprojects/SwarmUI)
- [Report Issues](https://github.com/opd-ai/asset-generator/issues)
- [Changelog](CHANGELOG.md)

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
- Inspired by best practices from the Go community
