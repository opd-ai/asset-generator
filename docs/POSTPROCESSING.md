# Postprocessing Features

This document covers all image postprocessing features that are applied after generation or as standalone operations.

## Table of Contents

- [Overview](#overview)
- [Auto-Crop](#auto-crop)
- [Downscaling](#downscaling)
- [PNG Metadata Stripping](#png-metadata-stripping)
- [Postprocessing Pipeline](#postprocessing-pipeline)

---

## Overview

The Asset Generator CLI provides powerful postprocessing capabilities:

1. **Auto-Crop** - Remove whitespace borders from images
2. **Downscaling** - High-quality image resizing with Lanczos filtering
3. **PNG Metadata Stripping** - Automatic removal of sensitive metadata

These features can be used:
- **Integrated**: During image generation (`generate` command)
- **Standalone**: On existing images (`crop`, `downscale` commands)
- **Pipeline**: Applied to all images in batch operations

### Processing Order

When multiple postprocessing steps are enabled, they execute in this order:

1. **Download** - Image retrieved from server
2. **Metadata Stripping** - PNG metadata removed (automatic)
3. **Auto-Crop** - Whitespace borders removed (if enabled)
4. **Downscaling** - Image resized (if enabled)

---

## Auto-Crop {#auto-crop}

Automatically detect and remove excess whitespace from image edges while optionally preserving the original aspect ratio.

### How It Works

The auto-crop algorithm:

1. **Edge Detection**: Scans from each edge (left, right, top, bottom) inward to find non-whitespace pixels
2. **Content Bounds**: Determines the minimal rectangle that contains all content
3. **Aspect Ratio Adjustment** (optional): Expands crop bounds to match original aspect ratio
4. **Cropping**: Extracts the identified region and saves the result

### Whitespace Detection

Pixels are considered "whitespace" when all RGB values are within the threshold range:
- **Threshold** (0-255): Minimum RGB value to be considered whitespace (default: 250)
- **Tolerance** (0-255): Allowed variation from threshold (default: 10)

A pixel is whitespace if: `R >= (threshold - tolerance) AND G >= (threshold - tolerance) AND B >= (threshold - tolerance)`

Additionally, fully transparent pixels (alpha < 10) are always considered whitespace.

### Standalone Crop Command

```bash
# Basic auto-crop
asset-generator crop image.png

# Auto-crop with output to new file
asset-generator crop input.png --output cropped.png

# Auto-crop in-place (replace original)
asset-generator crop image.png --in-place

# Preserve aspect ratio
asset-generator crop image.png --preserve-aspect --in-place

# Batch crop multiple images
asset-generator crop *.png --in-place

# Adjust sensitivity for darker backgrounds
asset-generator crop image.png --threshold 200 --tolerance 20
```

### Integrated with Generation

```bash
# Generate and auto-crop
asset-generator generate image \
  --prompt "centered logo design" \
  --save-images --auto-crop

# Generate, crop, and downscale (in that order)
asset-generator generate image \
  --prompt "high resolution art" \
  --save-images \
  --auto-crop \
  --downscale-width 1024

# Preserve aspect ratio when auto-cropping
asset-generator generate image \
  --prompt "portrait photo" \
  --save-images \
  --auto-crop \
  --auto-crop-preserve-aspect
```

### Pipeline Integration

```bash
# Apply auto-crop to all pipeline outputs
asset-generator pipeline \
  --file assets.yaml \
  --auto-crop \
  --auto-crop-threshold 245
```

### Crop Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--threshold` | 250 | Whitespace detection threshold (0-255) |
| `--tolerance` | 10 | Tolerance for near-white colors (0-255) |
| `--preserve-aspect` | false | Preserve original aspect ratio |
| `--quality` | 90 | JPEG quality (1-100) |
| `--output`, `-o` | (none) | Output file path (single file mode) |
| `--in-place`, `-i` | false | Replace original file(s) |

### Generate Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--auto-crop` | false | Enable auto-cropping |
| `--auto-crop-threshold` | 250 | Whitespace detection threshold (0-255) |
| `--auto-crop-tolerance` | 10 | Tolerance for near-white colors (0-255) |
| `--auto-crop-preserve-aspect` | false | Preserve original aspect ratio |

### Use Cases

#### Logo Design
```bash
asset-generator generate image \
  --prompt "minimalist tech company logo, centered, white background" \
  --save-images \
  --auto-crop \
  --auto-crop-preserve-aspect
```

#### Character Art
```bash
asset-generator generate image \
  --prompt "character portrait, full body" \
  --save-images \
  --auto-crop \
  --downscale-width 768
```

#### Screenshots Cleanup
```bash
asset-generator crop screenshot.png \
  --threshold 240 \
  --in-place
```

---

## Downscaling {#downscaling}

High-quality image resizing using advanced resampling algorithms, particularly Lanczos3 filtering for superior downscaling results.

### How It Works

The downscaling process:
1. **Dimension Calculation**: Determines target dimensions from width, height, or percentage
2. **Aspect Ratio**: Automatically maintains proportions if only one dimension specified
3. **Resampling**: Applies selected filter algorithm (Lanczos, Bilinear, or Nearest Neighbor)
4. **Format Preservation**: Maintains original format (PNG/JPEG)

### Standalone Downscale Command

```bash
# Downscale to specific width (maintains aspect ratio)
asset-generator downscale image.png --width 1024

# Downscale to 50% of original size
asset-generator downscale image.png --percentage 50

# Downscale to specific dimensions
asset-generator downscale photo.jpg --width 800 --height 600

# Downscale in-place (replaces original)
asset-generator downscale image.png --width 512 --in-place

# Downscale with specific output path
asset-generator downscale input.png --width 1024 --output-file resized.png

# Batch downscale multiple images by 75%
asset-generator downscale *.jpg --percentage 75 --in-place

# Use bilinear filter for speed
asset-generator downscale large.png --width 512 --filter bilinear
```

### Integrated with Generation

```bash
# Generate at high res, downscale for web
asset-generator generate image \
  --prompt "detailed artwork" \
  --width 2048 --height 2048 \
  --save-images \
  --downscale-width 1024

# Generate and downscale by percentage
asset-generator generate image \
  --prompt "character design" \
  --width 1536 --height 1536 \
  --save-images \
  --downscale-percentage 50

# Full postprocessing pipeline
asset-generator generate image \
  --prompt "logo design" \
  --width 2048 --height 2048 \
  --save-images \
  --auto-crop \
  --downscale-width 512 \
  --downscale-filter lanczos
```

### Pipeline Integration

```bash
# Downscale all pipeline outputs
asset-generator pipeline \
  --file assets.yaml \
  --width 1536 --height 2688 \
  --downscale-width 768 \
  --downscale-filter lanczos
```

### Filter Options

| Filter | Quality | Speed | Best For |
|--------|---------|-------|----------|
| `lanczos` | Highest | Slower | Photographs, final renders (default) |
| `bilinear` | Good | Faster | Web graphics, previews |
| `nearest` | Lowest | Fastest | Pixel art, icons |

### Downscale Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--width`, `-w` | 0 | Target width in pixels (0=auto from height) |
| `--height`, `-l` | 0 | Target height in pixels (0=auto from width) |
| `--percentage`, `-p` | 0 | Scale by percentage (1-100, 0=use width/height) |
| `--filter` | lanczos | Resampling filter: lanczos, bilinear, nearest |
| `--quality` | 90 | JPEG quality (1-100) |
| `--output-file` | (none) | Output file path (single file mode) |
| `--in-place` | false | Replace original file(s) |

### Generate Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--downscale-width` | 0 | Downscale to this width (0=disabled) |
| `--downscale-height` | 0 | Downscale to this height (0=disabled) |
| `--downscale-percentage` | 0 | Downscale by percentage (0=disabled) |
| `--downscale-filter` | lanczos | Filter: lanczos, bilinear, nearest |
| `--jpeg-quality` | 90 | JPEG quality (1-100) |

### Use Cases

#### Generate High, Save Small
```bash
# Generate at 2x resolution, downscale for delivery
asset-generator generate image \
  --prompt "product photography" \
  --width 2048 --height 2048 \
  --save-images \
  --downscale-width 1024
```

#### Batch Resize
```bash
# Resize all images in directory to 50%
asset-generator downscale *.png --percentage 50 --in-place
```

#### Web Optimization
```bash
asset-generator generate image \
  --prompt "website hero image" \
  --width 3840 --height 2160 \
  --save-images \
  --downscale-width 1920 \
  --downscale-filter bilinear
```

---

## PNG Metadata Stripping {#png-metadata-stripping}

**Mandatory, non-optional security and privacy feature** that ensures no sensitive information is accidentally embedded in output images.

### What Metadata is Removed

All PNG ancillary chunks are stripped, including:

- **Text chunks** (tEXt, zTXt, iTXt) - May contain generation prompts, parameters, or other textual data
- **Time chunk** (tIME) - Timestamp of when the image was created
- **Physical dimensions** (pHYs) - DPI and physical size information
- **Gamma information** (gAMA) - Display gamma correction data
- **Color profiles** (iCCP) - ICC color profile data
- **Other metadata** - Any other ancillary chunks

### What is Preserved

Only critical PNG chunks required for displaying the image:

- **IHDR** - Image header (dimensions, bit depth, color type)
- **PLTE** - Color palette (for indexed color images)
- **IDAT** - Actual image pixel data
- **IEND** - End of PNG file marker

### When Metadata is Stripped

Metadata stripping is automatically applied:

1. **During Image Download** - When images are downloaded from the generation API
2. **During Auto-Crop** - After cropping operations
3. **During Downscaling** - After resizing operations

**This is automatic and cannot be disabled** - it happens for all PNG files processed by the CLI.

### Why This is Important

#### Privacy Protection
Image generation APIs may embed sensitive information in PNG metadata:
- Full generation prompts (including potentially sensitive content)
- API keys or session tokens
- User identifiers
- Generation timestamps
- Server information

#### Security
Metadata can reveal:
- Which AI models were used
- Generation parameters that could be reverse-engineered
- Infrastructure details about your generation setup

#### Clean Deliverables
For professional or commercial use, clean images without metadata:
- No accidental information leakage
- Smaller file sizes
- Professional presentation
- Compliance with privacy requirements

### Technical Implementation

The metadata stripping is implemented in `pkg/processor/metadata.go` and uses pure Go PNG decoding/encoding:

```go
func StripPNGMetadata(inputPath, outputPath string) error {
    // Decode PNG
    img, err := png.Decode(file)
    
    // Re-encode with only critical chunks
    // (IHDR, PLTE, IDAT, IEND)
    err = png.Encode(outputFile, img)
}
```

### Verification

You can verify metadata has been removed using standard tools:

```bash
# Check metadata before
exiftool generated-image.png

# Process with asset-generator
asset-generator crop generated-image.png -o cleaned.png

# Verify metadata removed
exiftool cleaned.png  # Should show minimal information
```

### JPEG Handling

**Note**: This feature is specific to PNG files. JPEG files do not undergo metadata stripping in the current implementation, as JPEG metadata (EXIF) handling requires different tooling.

---

## Postprocessing Pipeline {#postprocessing-pipeline}

Combine multiple postprocessing steps for optimal results.

### Full Pipeline Example

```bash
# Complete postprocessing workflow
asset-generator generate image \
  --prompt "professional product photo" \
  --width 2048 --height 2048 \
  --save-images \
  --output-dir ./products \
  --auto-crop \
  --auto-crop-preserve-aspect \
  --downscale-width 1024 \
  --downscale-filter lanczos
```

**What happens**:
1. Image generated at 2048x2048
2. Downloaded to ./products/
3. PNG metadata automatically stripped
4. Whitespace borders removed (preserving aspect ratio)
5. Downscaled to 1024px width with Lanczos filtering

### Pipeline Batch Processing

```bash
# Apply to all assets in pipeline
asset-generator pipeline \
  --file tarot-deck.yaml \
  --width 1536 --height 2688 \
  --auto-crop \
  --downscale-width 768 \
  --downscale-percentage 0
```

### Standalone Batch Processing

```bash
# Process existing images
for img in *.png; do
  asset-generator crop "$img" --in-place --preserve-aspect
  asset-generator downscale "$img" --width 1024 --in-place
done
```

### Best Practices

#### 1. Generate High, Deliver Low

Always generate at higher resolution than needed, then downscale:

```bash
# Generate 2x, deliver 1x
asset-generator generate image \
  --width 2048 \
  --save-images \
  --downscale-width 1024
```

#### 2. Crop Before Downscale

Remove unnecessary borders before resizing for best results:

```bash
# Crop first, then downscale
--auto-crop \
--downscale-width 512
```

#### 3. Preserve Aspect Ratios

Use `--auto-crop-preserve-aspect` for consistent dimensions:

```bash
--auto-crop \
--auto-crop-preserve-aspect \
--downscale-width 768
```

#### 4. Choose Right Filter

- Production: `--downscale-filter lanczos`
- Preview: `--downscale-filter bilinear`
- Pixel art: `--downscale-filter nearest`

### Performance Considerations

- **Auto-crop**: Fast, minimal overhead
- **Downscaling (Lanczos)**: Slower but highest quality
- **Downscaling (Bilinear)**: Fast, good quality
- **Metadata stripping**: Fast, negligible overhead

### Troubleshooting

#### Auto-Crop Removes Too Much

```bash
# Reduce threshold to be less aggressive
--auto-crop-threshold 245 \
--auto-crop-tolerance 15
```

#### Auto-Crop Not Removing Enough

```bash
# Increase threshold for more aggressive cropping
--auto-crop-threshold 240 \
--auto-crop-tolerance 5
```

#### Downscaling Too Slow

```bash
# Use bilinear for speed
--downscale-filter bilinear
```

#### File Size Issues

```bash
# Adjust JPEG quality
--jpeg-quality 85
```

---

## See Also

- [Filename Templates](FILENAME_TEMPLATES.md) - Custom naming for downloaded images
- [Generation Features](GENERATION_FEATURES.md) - Scheduler, Skimmed CFG
- [Pipeline Processing](PIPELINE.md) - Batch generation workflows
- [SVG Conversion](SVG_CONVERSION.md) - Convert images to vector format
- [Quick Start Guide](QUICKSTART.md) - Getting started with the CLI
