# SwarmUI CLI Client

A powerful command-line interface for interacting with the SwarmUI API. Generate assets, manage models, and configure your SwarmUI instance with ease.

## Features

- 🎨 **Asset Generation**: Generate images using text-to-image models
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
swarmui config init
```

### 2. Set Your API Endpoint

```bash
swarmui config set api-url http://localhost:7801
```

### 3. Generate Your First Image

```bash
swarmui generate image --prompt "a beautiful sunset over mountains"
```

## Usage

### Asset Generation

Generate images with various parameters:

```bash
# Basic generation
swarmui generate image --prompt "a futuristic cityscape"

# Advanced generation with custom parameters
swarmui generate image \
  --prompt "astronaut riding a horse" \
  --width 1024 \
  --height 1024 \
  --steps 30 \
  --cfg-scale 7.5 \
  --sampler euler_a

# Save output to file
swarmui generate image \
  --prompt "cat wearing sunglasses" \
  --output result.json \
  --format json

# Batch generation
swarmui generate image \
  --prompt "beautiful landscape" \
  --batch 4
```

### Model Management

List and inspect available models:

```bash
# List all models
swarmui models list

# List models as JSON
swarmui models list --format json

# Get details about a specific model
swarmui models get stable-diffusion-xl
```

### Configuration

Manage your CLI configuration:

```bash
# View current configuration
swarmui config view

# Set configuration values
swarmui config set api-url https://api.swarm.example.com
swarmui config set api-key your-api-key-here

# Get a specific configuration value
swarmui config get api-url
```

## Configuration

Configuration can be provided through multiple sources with the following precedence:

1. **Command-line flags** (highest priority)
2. **Environment variables** (prefixed with `SWARMUI_`)
3. **Configuration file** (`~/.swarmui/config.yaml`)
4. **Default values** (lowest priority)

### Configuration File

Location: `~/.swarmui/config.yaml`

```yaml
api-url: http://localhost:7801
api-key: your-api-key
format: table
verbose: false
```

### Environment Variables

```bash
export SWARMUI_API_URL=http://localhost:7801
export SWARMUI_API_KEY=your-api-key
export SWARMUI_FORMAT=json
```

## Global Flags

Available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--api-url` | | SwarmUI API base URL | `http://localhost:7801` |
| `--api-key` | | SwarmUI API key | |
| `--format` | `-f` | Output format (table, json, yaml) | `table` |
| `--output` | `-o` | Write output to file | |
| `--quiet` | `-q` | Quiet mode (errors only) | `false` |
| `--verbose` | `-v` | Verbose output | `false` |
| `--config` | | Config file path | `~/.swarmui/config.yaml` |

## Generation Parameters

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--prompt` | `-p` | Generation prompt (required) | |
| `--model` | | Model to use | |
| `--width` | `-w` | Image width | `512` |
| `--height` | `-h` | Image height | `512` |
| `--steps` | | Inference steps | `20` |
| `--cfg-scale` | | CFG scale (guidance) | `7.5` |
| `--sampler` | | Sampling method | `euler_a` |
| `--batch` | `-b` | Number of images to generate | `1` |
| `--seed` | | Random seed (-1 for random) | `-1` |
| `--negative-prompt` | `-n` | Negative prompt | |

## Examples

### Example 1: Basic Image Generation

```bash
swarmui generate image --prompt "a serene lake at sunset"
```

### Example 2: High-Quality Generation

```bash
swarmui generate image \
  --prompt "professional portrait photo of a scientist" \
  --width 1024 \
  --height 1024 \
  --steps 50 \
  --cfg-scale 8.0
```

### Example 3: Batch Generation with Seed

```bash
swarmui generate image \
  --prompt "fantasy castle in the clouds" \
  --batch 4 \
  --seed 42 \
  --output results.json
```

### Example 4: Pipeline with JSON Output

```bash
# Generate and extract image paths
swarmui generate image \
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
│   ├── client/            # SwarmUI API client
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
- **pkg/client/**: SwarmUI API client (reusable library)
- **pkg/output/**: Output formatting utilities
- **internal/config/**: Configuration validation logic

### API Client

The SwarmUI API client is designed to be reusable and can be imported by other Go projects:

```go
import "github.com/opd-ai/asset-generator/pkg/client"

// Create client
config := &client.Config{
    BaseURL: "http://localhost:7801",
    APIKey:  "your-api-key",
}
swarmClient, err := client.NewSwarmClient(config)

// Generate image
req := &client.GenerationRequest{
    Prompt: "a beautiful landscape",
    Parameters: map[string]interface{}{
        "width": 1024,
        "height": 1024,
    },
}
result, err := swarmClient.GenerateImage(context.Background(), req)
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

If you're having trouble connecting to SwarmUI:

```bash
# Check your API URL
swarmui config get api-url

# Test with verbose output
swarmui models list -v

# Verify SwarmUI is running
curl http://localhost:7801/API/ListModels
```

### Configuration Issues

```bash
# View current configuration
swarmui config view

# Reinitialize configuration
rm -rf ~/.swarmui
swarmui config init
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
