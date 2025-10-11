# Auto-Crop Feature

The asset-generator includes automatic image cropping functionality that removes excess whitespace from image edges while optionally preserving the original aspect ratio.

## Overview

Auto-crop is a local post-processing step that:
- Detects whitespace by scanning from each edge inward
- Removes excess padding/borders from generated images
- Can preserve the original aspect ratio
- Works as part of the generation pipeline or as a standalone command
- Supports both PNG and JPEG formats

## How It Works

The auto-crop algorithm:

1. **Edge Detection**: Scans from each edge (left, right, top, bottom) inward to find the first non-whitespace pixels
2. **Content Bounds**: Determines the minimal rectangle that contains all content
3. **Aspect Ratio Adjustment** (optional): Expands the crop bounds to match the original image's aspect ratio
4. **Cropping**: Extracts the identified region and saves the result

### Whitespace Detection

Pixels are considered "whitespace" when all RGB values are within the threshold range:
- **Threshold** (0-255): Minimum RGB value to be considered whitespace (default: 250)
- **Tolerance** (0-255): Allowed variation from threshold (default: 10)

A pixel is whitespace if: `R >= (threshold - tolerance) AND G >= (threshold - tolerance) AND B >= (threshold - tolerance)`

Additionally, fully transparent pixels (alpha < 10) are always considered whitespace.

## Usage

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

Auto-crop can be enabled as part of the image generation pipeline:

```bash
# Generate and auto-crop
asset-generator generate image \
  --prompt "centered logo design" \
  --save-images --auto-crop

# Generate, crop, and downscale (in that order)
asset-generator generate image \
  --prompt "high resolution art" \
  --width 2048 --height 2048 \
  --save-images \
  --auto-crop \
  --downscale-width 1024

# Preserve aspect ratio when auto-cropping
asset-generator generate image \
  --prompt "product photo" \
  --save-images \
  --auto-crop --auto-crop-preserve-aspect
```

## Post-Processing Pipeline Order

When both auto-crop and downscaling are enabled, they are applied in this order:

1. **Download**: Image is downloaded from the generation API
2. **Auto-Crop**: Whitespace borders are removed
3. **Downscale**: Image is resized to target dimensions

This order ensures that:
- Unnecessary whitespace isn't included in the downscaled image
- Downscaling operates on the minimal content area
- Better quality and smaller file sizes

## Configuration Options

### Crop Command Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--threshold` | int | 250 | Whitespace detection threshold (0-255) |
| `--tolerance` | int | 10 | Tolerance for near-white colors (0-255) |
| `--preserve-aspect` | bool | false | Preserve original aspect ratio |
| `--quality` | int | 90 | JPEG quality (1-100) |
| `--output`, `-o` | string | - | Output file path (single file mode) |
| `--in-place`, `-i` | bool | false | Replace original file |

### Generate Command Auto-Crop Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--auto-crop` | bool | false | Enable auto-crop |
| `--auto-crop-threshold` | int | 250 | Whitespace detection threshold |
| `--auto-crop-tolerance` | int | 10 | Tolerance for near-white colors |
| `--auto-crop-preserve-aspect` | bool | false | Preserve aspect ratio |

## Examples

### Example 1: Centered Logo

Generate a logo and remove padding:

```bash
asset-generator generate image \
  --prompt "minimalist tech company logo, centered, white background" \
  --width 1024 --height 1024 \
  --save-images --auto-crop
```

Result: Logo is cropped to content bounds, removing all white padding.

### Example 2: Product Photo with Consistent Aspect Ratio

Generate product photos maintaining the original aspect ratio:

```bash
asset-generator generate image \
  --prompt "professional product photo, white background" \
  --width 1024 --height 768 \
  --batch 5 \
  --save-images \
  --auto-crop --auto-crop-preserve-aspect
```

Result: Products are cropped but maintain the 4:3 aspect ratio.

### Example 3: High-Res Generation with Aggressive Post-Processing

Generate at high resolution, crop, and downscale:

```bash
asset-generator generate image \
  --prompt "detailed digital art" \
  --width 2048 --height 2048 \
  --save-images \
  --auto-crop \
  --downscale-width 1024 \
  --downscale-filter lanczos
```

Result: 2048x2048 image is generated, cropped to content, then downscaled to max 1024px width.

### Example 4: Batch Processing Existing Images

Crop existing images in a directory:

```bash
# Crop all PNGs, preserving aspect ratio
asset-generator crop *.png --preserve-aspect --in-place

# Crop JPEGs with higher quality
asset-generator crop *.jpg --quality 95 --in-place
```

## Sensitivity Tuning

### Default Settings (Most Common)
- Threshold: 250
- Tolerance: 10
- Detects: Near-white backgrounds (RGB > 240)

### Light Gray Backgrounds
- Threshold: 230
- Tolerance: 20
- Detects: Light gray and white backgrounds

### Aggressive (Pure White Only)
- Threshold: 255
- Tolerance: 5
- Detects: Only pure white or very close

### Conservative (Darker Backgrounds)
- Threshold: 200
- Tolerance: 30
- Detects: Broader range of light colors

## Technical Details

### Algorithm Complexity
- Time: O(W × H) where W and H are image dimensions
- Space: O(1) for detection, O(W' × H') for output where W' and H' are cropped dimensions

### Edge Cases Handled
- **No whitespace found**: Image is returned unchanged
- **Entire image is whitespace**: Returns error
- **Aspect ratio preservation**: Expands bounds symmetrically to maintain ratio
- **Transparent images**: Treats fully transparent pixels as whitespace

### Performance Considerations
- Fast for images with large whitespace borders (early termination)
- Linear scan from each edge
- Minimal memory overhead (in-place processing)

## API Usage

For programmatic use in Go:

```go
import "github.com/opd-ai/asset-generator/pkg/processor"

// Auto-crop with defaults
opts := processor.CropOptions{
    Threshold:           250,
    Tolerance:           10,
    JPEGQuality:         90,
    PreserveAspectRatio: false,
}

// Crop to new file
err := processor.AutoCropImage("input.png", "output.png", opts)

// Crop in-place
err := processor.AutoCropInPlace("image.png", opts)
```

## Troubleshooting

### Problem: Too much content is cropped

**Solution**: Increase tolerance or decrease threshold
```bash
asset-generator crop image.png --threshold 240 --tolerance 5
```

### Problem: Not enough whitespace is removed

**Solution**: Decrease tolerance or increase threshold
```bash
asset-generator crop image.png --threshold 252 --tolerance 3
```

### Problem: Aspect ratio doesn't match original

**Solution**: Enable `--preserve-aspect`
```bash
asset-generator crop image.png --preserve-aspect
```

### Problem: Colored backgrounds not detected

**Note**: Auto-crop only detects light/white backgrounds. For colored backgrounds, consider using different threshold values or manual cropping tools.

## See Also

- [Downscaling Feature](DOWNSCALING_FEATURE.md) - Resize images after generation
- [Image Download](IMAGE_DOWNLOAD.md) - Save generated images locally
- [Filename Templates](FILENAME_TEMPLATES.md) - Template-based filename generation
