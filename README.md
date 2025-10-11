# Asset Generator CLI

A powerful command-line interface for interacting with AI asset generation APIs. Generate assets, manage models, and configure your asset generation workflow with ease.

## Features

- 🎨 **Asset Generation**: Generate images using text-to-image models
- 🎯 **LoRA Support**: Apply Low-Rank Adaptation models for style and content customization
- 🎯 **Skimmed CFG**: Advanced sampling technique for improved quality and speed
- 🔄 **Pipeline Processing**: Automated batch generation from YAML pipeline files
- 💾 **Image Download**: Automatically download and save generated images locally
- 🔒 **Automatic Metadata Stripping**: All PNG images have metadata removed for privacy and security
- ✂️ **Auto-Crop**: Remove whitespace borders from images while preserving aspect ratio
- 🔽 **Image Postprocessing**: High-quality Lanczos downscaling after download
- 🎨 **SVG Conversion**: Convert images to SVG format using geometric shapes or edge tracing
- 🏥 **Server Status**: Check SwarmUI server health and backend information
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
  --sampler euler_a \
  --scheduler karras

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

# Use Skimmed CFG for improved quality and speed
asset-generator generate image \
  --prompt "detailed portrait photograph" \
  --skimmed-cfg \
  --skimmed-cfg-scale 3.0

# Apply LoRA models for style customization
asset-generator generate image \
  --prompt "anime character in magical forest" \
  --lora "anime-style-v2:0.9" \
  --lora "detailed-eyes:0.7"

# Combine LoRAs with other features
asset-generator generate image \
  --prompt "cyberpunk cityscape" \
  --lora "cyberpunk-aesthetic:1.0" \
  --lora "neon-lights:0.7" \
  --skimmed-cfg \
  --batch 4 \
  --save-images
```

### Pipeline Processing

Process YAML pipeline files for automated batch generation:

```bash
# Process a tarot deck pipeline (78 cards)
asset-generator pipeline --file tarot-spec.yaml --output-dir ./deck

# Preview what would be generated (dry run)
asset-generator pipeline --file tarot-spec.yaml --dry-run

# Custom generation parameters
asset-generator pipeline --file tarot-spec.yaml \
  --base-seed 42 \
  --steps 40 \
  --width 768 \
  --height 1344 \
  --style-suffix "detailed illustration, ornate border, rich colors"

# With postprocessing
asset-generator pipeline --file tarot-spec.yaml \
  --auto-crop \
  --downscale-width 1024 \
  --continue-on-error
```

See [docs/PIPELINE.md](docs/PIPELINE.md) for complete documentation and the [examples/tarot-deck/](examples/tarot-deck/) directory for a full working example.

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

### Server Status

Check the SwarmUI server status and health:

```bash
# Check server status
asset-generator status

# Get status as JSON for automation
asset-generator status --format json

# Verbose output with API details
asset-generator status -v
```

**NEW: Cross-Process Generation Tracking** 🎉

The status command now tracks active generations across different terminal windows using file-based state sharing:

```bash
# Terminal 1: Start a generation
cd ~/my-project
asset-generator generate image --prompt "artwork"

# Terminal 2: Check status (different process!)
cd ~/my-project
asset-generator status
# Shows: Generation details with progress!
```

State is persisted to `.asset-generator-state.json` in the working directory, enabling accurate cross-process tracking while maintaining project isolation.

See [docs/COMMANDS.md](docs/COMMANDS.md) and [docs/STATE_FILE_SHARING.md](docs/STATE_FILE_SHARING.md) for complete documentation.

### Cancel Generations

Stop ongoing or queued image generations:

```bash
# Cancel the current generation
asset-generator cancel

# Cancel all queued generations
asset-generator cancel --all

# Cancel quietly (for scripts)
asset-generator cancel -q
```

Useful for:
- Stopping long-running generations (e.g., Flux models: 5-40 minutes)
- Clearing a backlog of queued generations
- Recovering from stuck generations

See [docs/COMMANDS.md](docs/COMMANDS.md) for complete documentation.

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

Configuration file location priority (highest to lowest):

1. **Custom config file** (via `--config` flag) - When specified, bypasses all search paths and uses only this file
2. `./config/config.yaml` - Current directory
3. `~/.asset-generator/config.yaml` - User's home directory
4. `/etc/asset-generator/config.yaml` - System-wide configuration

When `--config` is specified, only that file is used. Otherwise, the first file found in locations 2-4 is used.

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
| `--config` | | Custom config file path (overrides default search) | (searches multiple locations) |

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
| `--scheduler` | | Scheduler/noise schedule (simple, normal, karras, exponential, sgm_uniform) | `simple` |
| `--batch` | `-b` | Number of images to generate | `1` |
| `--seed` | | Random seed (-1 for random) | `-1` |
| `--negative-prompt` | `-n` | Negative prompt | |
| `--websocket` | | Use WebSocket for real-time progress (falls back to HTTP if unavailable) | `false` |
| `--save-images` | | Download and save generated images to local disk | `false` |
| `--output-dir` | | Directory to save downloaded images | `.` (current directory) |
| `--filename-template` | | Custom filename pattern (see placeholders below) | |
| `--downscale-width` | | Downscale to this width after download (0=auto) | `0` |
| `--downscale-height` | | Downscale to this height after download (0=auto) | `0` |
| `--downscale-percentage` | | Downscale by percentage (1-100, takes precedence) | `0` |
| `--downscale-filter` | | Downscaling algorithm: lanczos, bilinear, nearest | `lanczos` |
| `--skimmed-cfg` | | Enable Skimmed CFG (Distilled CFG) for improved quality/speed | `false` |
| `--skimmed-cfg-scale` | | Skimmed CFG scale value (typically lower than standard CFG) | `3.0` |
| `--skimmed-cfg-start` | | Start percentage for Skimmed CFG application (0.0-1.0) | `0.0` |
| `--skimmed-cfg-end` | | End percentage for Skimmed CFG application (0.0-1.0) | `1.0` |

> **Note:** The `--length` flag is used for the vertical dimension (height) for API compatibility with SwarmUI. Both `--length` and `--height` are supported as aliases.

### About Skimmed CFG

Skimmed CFG (also known as Distilled CFG or Dynamic CFG) is an advanced sampling technique that can improve generation quality and speed with compatible models. It works by applying a more efficient guidance strategy during the denoising process.

**Key benefits:**
- Improved image quality and coherence
- Potentially faster generation times
- Better adherence to prompts with lower CFG scales

**Usage tips:**
- Use lower `--skimmed-cfg-scale` values (2.0-4.0) compared to standard CFG
- Adjust `--skimmed-cfg-start` and `--skimmed-cfg-end` to apply only during specific phases of generation
- Not all models support this feature - check your model documentation

**Learn more:** See the [Generation Features documentation](docs/GENERATION_FEATURES.md) for detailed usage examples and best practices.

### About Scheduler Selection

The scheduler (also called noise schedule) controls how noise is added and removed during the diffusion process. Different schedulers can significantly impact generation quality and speed.

**Available schedulers:**
- `simple` (default): Fast, reliable, good for general use
- `normal`: Standard schedule, balanced quality
- `karras`: High-quality, detailed images (recommended for final renders)
- `exponential`: Smooth transitions, good for artistic work
- `sgm_uniform`: Specialized for certain models

**Usage tips:**
- Use `simple` for quick iteration and testing
- Use `karras` for high-quality final outputs
- Different schedulers work better with different samplers
- Increase steps (30-50) when using `karras` for best results

**Learn more:** See the [Generation Features documentation](docs/GENERATION_FEATURES.md) for detailed comparisons and usage examples.

## LoRA Support

LoRA (Low-Rank Adaptation) models allow you to customize and fine-tune the style, content, and characteristics of generated images without changing the base model. Apply lightweight model adaptations for anime styles, photorealistic enhancements, specific characters, or artistic themes.

**Basic usage:**
```bash
# Single LoRA with default weight
asset-generator generate image --prompt "anime character" --lora "anime-style"

# Single LoRA with custom weight (0.8)
asset-generator generate image --prompt "portrait" --lora "realistic-faces:0.8"

# Multiple LoRAs for complex styles
asset-generator generate image \
  --prompt "cyberpunk cityscape" \
  --lora "cyberpunk-aesthetic:1.0" \
  --lora "neon-lights:0.7" \
  --lora "detailed-architecture:0.5"
```

**Key features:**
- Support for multiple LoRAs simultaneously
- Flexible weight control (inline or explicit)
- Negative weights to remove unwanted styles
- Config file support for default LoRAs
- Weight validation and helpful error messages

**Learn more:** 
- [LoRA Support Documentation](docs/LORA_SUPPORT.md) - Complete guide with examples

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

# Use custom filename template
asset-generator generate image \
  --prompt "fantasy landscape" \
  --batch 5 \
  --save-images \
  --filename-template "landscape-{index}-seed{seed}.png"

# Batch generation with download
asset-generator generate image \
  --prompt "abstract art" \
  --batch 5 \
  --save-images \
  --output-dir ./batch-output
```

### Custom Filenames

You can customize the downloaded image filenames using the `--filename-template` flag with various placeholders:

```bash
# Include seed and index in filename
asset-generator generate image \
  --prompt "portrait" \
  --seed 42 \
  --batch 3 \
  --save-images \
  --filename-template "portrait-seed{seed}-{index}.png"
# Output: portrait-seed42-000.png, portrait-seed42-001.png, portrait-seed42-002.png

# Include timestamp and model
asset-generator generate image \
  --prompt "landscape" \
  --model "flux-dev" \
  --save-images \
  --filename-template "{model}-{datetime}.png"
# Output: flux-dev-2024-10-08_14-30-45.png
```

**Available placeholders:** `{index}`, `{i1}`, `{timestamp}`, `{datetime}`, `{date}`, `{time}`, `{seed}`, `{model}`, `{width}`, `{height}`, `{prompt}`, `{original}`, `{ext}`

See [Filename Templates documentation](docs/FILENAME_TEMPLATES.md) for complete placeholder reference.

### Behavior

- Images are downloaded immediately after generation completes
- By default, the original filename from the server is preserved
- Custom filenames can be specified with `--filename-template`
- Directory is created automatically if it doesn't exist
- Progress feedback shows each downloaded file
- Local file paths are added to the metadata output
- Partial failures are handled gracefully (some images may succeed even if others fail)

### Local Postprocessing

The CLI supports a multi-stage postprocessing pipeline applied after downloading:

1. **Auto-Crop** (optional): Remove whitespace borders
2. **Downscale** (optional): Resize with high-quality filtering

#### Auto-Crop

Automatically remove excess whitespace from image edges while optionally preserving aspect ratio:

```bash
# Generate and auto-crop
asset-generator generate image \
  --prompt "centered logo design" \
  --save-images --auto-crop

# Auto-crop preserving aspect ratio
asset-generator generate image \
  --prompt "product photo" \
  --save-images \
  --auto-crop --auto-crop-preserve-aspect

# Combine crop and downscale
asset-generator generate image \
  --prompt "high resolution art" \
  --width 2048 --height 2048 \
  --save-images \
  --auto-crop \
  --downscale-width 1024
```

**Auto-Crop Flags:**
| Flag | Description | Default |
|------|-------------|---------|
| `--auto-crop` | Enable automatic whitespace removal | `false` |
| `--auto-crop-threshold` | Whitespace detection threshold (0-255) | `250` |
| `--auto-crop-tolerance` | Tolerance for near-white colors (0-255) | `10` |
| `--auto-crop-preserve-aspect` | Preserve original aspect ratio | `false` |

See [Auto-Crop Documentation](AUTO_CROP_FEATURE.md) for detailed usage and sensitivity tuning.

#### Downscaling

Reduce image dimensions using high-quality filtering:

```bash
# Generate at high resolution, save downscaled version
asset-generator generate image \
  --prompt "detailed artwork" \
  --width 2048 --height 2048 \
  --save-images \
  --downscale-width 1024

# Downscale by percentage (simplest method)
asset-generator generate image \
  --prompt "high resolution photo" \
  --width 2048 --height 2048 \
  --save-images \
  --downscale-percentage 50  # Results in 1024x1024

# Downscale by height (width auto-calculated)
asset-generator generate image \
  --prompt "portrait" \
  --width 1920 --height 1080 \
  --save-images \
  --downscale-height 720

# Choose different downscaling algorithm
asset-generator generate image \
  --prompt "photo" \
  --save-images \
  --downscale-width 800 \
  --downscale-filter lanczos  # Options: lanczos, bilinear, nearest
```

**Downscaling Features:**
- Applied locally after download (saves API bandwidth)
- Uses Lanczos3 resampling by default for best quality
- Percentage-based scaling maintains aspect ratio automatically
- Auto-maintains aspect ratio when only one dimension specified
- Prevents accidental upscaling
- Three filter options: `lanczos` (highest quality), `bilinear` (balanced), `nearest` (fastest)

**Flags for `generate image` command:**
| Flag | Description | Default |
|------|-------------|---------|
| `--downscale-width` | Target width in pixels (0=auto from height) | `0` |
| `--downscale-height` | Target height in pixels (0=auto from width) | `0` |
| `--downscale-percentage` | Scale by percentage (1-100, overrides width/height) | `0` |
| `--downscale-filter` | Algorithm: `lanczos`, `bilinear`, `nearest` | `lanczos` |

### Output

When `--save-images` is enabled, you'll see output like:

```
Generating image with prompt: cyberpunk cityscape
Downloading generated images...
  [1/1] Saved: ./generated-art/cyberpunk cityscape-1234567890.png
✓ Generation completed successfully (1 image)
```

## Image Cropping

Remove excess whitespace from image edges automatically.

### Standalone Crop Command

```bash
# Basic auto-crop
asset-generator crop image.png

# Crop in-place (replace original)
asset-generator crop image.png --in-place

# Preserve aspect ratio
asset-generator crop photo.jpg --preserve-aspect

# Batch crop multiple images
asset-generator crop *.png --in-place

# Adjust sensitivity
asset-generator crop image.png --threshold 240 --tolerance 5
```

**Crop Flags:**
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--threshold` | | Whitespace detection threshold (0-255) | `250` |
| `--tolerance` | | Tolerance for near-white colors (0-255) | `10` |
| `--preserve-aspect` | | Preserve original aspect ratio | `false` |
| `--quality` | | JPEG quality (1-100) | `90` |
| `--output` | `-o` | Output file path (single file mode) | |
| `--in-place` | `-i` | Replace original file(s) | `false` |

See [Auto-Crop Documentation](AUTO_CROP_FEATURE.md) for detailed usage and sensitivity tuning.

## Image Downscaling

The `downscale` command allows you to resize images independently of the generation workflow. Supports both absolute dimensions and percentage-based scaling.

### Usage

```bash
# Downscale to specific width (auto-calculates height)
asset-generator downscale image.png --width 1024

# Downscale by percentage (simplest method)
asset-generator downscale image.png --percentage 50

# Downscale to specific dimensions
asset-generator downscale photo.jpg --width 800 --height 600

# Downscale in-place (replaces original)
asset-generator downscale image.png --percentage 75 --in-place

# Batch downscale multiple images
asset-generator downscale *.jpg --percentage 50 --in-place

# Custom output location
asset-generator downscale input.png --width 512 --output resized.png

# Choose different filter for speed
asset-generator downscale large.png --width 512 --filter bilinear
```

### Downscale Flags
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--width` | `-w` | Target width in pixels (0=auto from height) | `0` |
| `--height` | `-l` | Target height in pixels (0=auto from width) | `0` |
| `--percentage` | `-p` | Scale by percentage (1-100, overrides width/height) | `0` |
| `--filter` | | Resampling filter: lanczos, bilinear, nearest | `lanczos` |
| `--quality` | | JPEG quality (1-100) | `90` |
| `--output-file` | | Output file path (single file mode) | |
| `--in-place` | | Replace original file(s) | `false` |

**Filter Options:**
- `lanczos` - Highest quality, best for photographs (default)
- `bilinear` - Good balance of speed and quality
- `nearest` - Fastest, best for pixel art or icons

**Features:**
- Percentage-based scaling automatically maintains aspect ratio
- Automatically calculates missing dimension to maintain aspect ratio
- Prevents accidental upscaling (will error if target > source)
- Preserves image format (PNG/JPEG)
- Batch processing support

## SVG Conversion

Convert images to SVG format using two powerful methods:

### Quick Start

```bash
# Convert any image to SVG (default: 100 geometric shapes)
asset-generator convert svg image.png

# High quality conversion with 500 shapes
asset-generator convert svg photo.jpg --shapes 500

# Use edge tracing for line art
asset-generator convert svg sketch.png --method gotrace
```

### Conversion Methods

**Primitive Method** (default): Creates artistic SVG using geometric shapes
- Good for: Photos, logos, illustrations
- Fast and produces clean results
- Adjustable quality via `--shapes` flag

**Gotrace Method**: Uses edge tracing for precise vector conversion
- Good for: Line art, sketches, high-contrast images
- Pure-Go implementation (no external dependencies required)
- Better detail preservation
- Uses default tracing parameters (no customization flags)

### Examples

```bash
# Basic conversion
asset-generator convert svg input.png

# Custom output location
asset-generator convert svg input.png -o artwork.svg

# Different shape types
asset-generator convert svg photo.jpg --shapes 200 --mode 3  # ellipses
asset-generator convert svg photo.jpg --shapes 150 --mode 6  # bezier curves

# Edge tracing conversion
asset-generator convert svg lineart.png --method gotrace
```

See [SVG Conversion Documentation](SVG_CONVERSION.md) for comprehensive guide and examples.

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
│   ├── converter/         # Image format converters (SVG)
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
- [PNG Metadata Stripping](PNG_METADATA_STRIPPING.md)
- [Auto-Crop Feature](AUTO_CROP_FEATURE.md)
- [Downscaling Feature](DOWNSCALING_FEATURE.md)
- [SVG Conversion](SVG_CONVERSION.md)

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
- Inspired by best practices from the Go community
