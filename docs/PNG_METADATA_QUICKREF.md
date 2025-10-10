# PNG Metadata Stripping - Quick Reference

## What It Does

**Automatically removes all PNG metadata from downloaded and processed images.**

## Why It Matters

- 🔒 **Privacy**: Prevents leaking prompts, API keys, or generation parameters
- 🔐 **Security**: No sensitive data embedded in images
- 📦 **Clean Output**: Professional deliverables without metadata bloat
- 💾 **Smaller Files**: Reduced file size by removing unnecessary chunks

## When It Happens

Metadata is **automatically stripped** (no action required):

1. ✅ When downloading images via `generate` command
2. ✅ When cropping images via `crop` command  
3. ✅ When downscaling images via `downscale` command

## What Gets Removed

All PNG ancillary chunks:
- `tEXt`, `zTXt`, `iTXt` - Text metadata (prompts, parameters, etc.)
- `tIME` - Timestamps
- `pHYs` - Physical dimensions/DPI
- `gAMA` - Gamma correction
- `iCCP` - Color profiles
- All other non-critical chunks

## What Stays

Only critical display chunks:
- `IHDR` - Image header (dimensions, bit depth)
- `PLTE` - Color palette (if needed)
- `IDAT` - Image pixel data
- `IEND` - End marker

## Usage

**No flags or configuration needed!** It happens automatically.

```bash
# Download with automatic metadata stripping
asset-generator generate --prompt "..." --save-images

# Crop with automatic metadata stripping
asset-generator crop --input image.png --output cropped.png

# Downscale with automatic metadata stripping
asset-generator downscale --input image.png --width 512
```

## Verification

Check if metadata was removed:

```bash
# Using exiftool
exiftool image.png

# Using pngcheck (shows chunk structure)
pngcheck -v image.png

# Using ImageMagick
identify -verbose image.png
```

## Can I Disable It?

**No.** This is a mandatory security feature and cannot be disabled.

## Does It Affect Quality?

**No.** Only metadata is removed; pixel data is never re-compressed.

## Performance

- ⚡ Processing time: 5-50ms per image
- 📉 File size: Usually 1-5% smaller
- 💾 Memory: Brief duplication during processing

## Programmatic Usage

```go
import "github.com/opd-ai/asset-generator/pkg/processor"

// Strip metadata from any PNG file
err := processor.StripPNGMetadata("/path/to/image.png")

// Alternative name (same function)
err = processor.EnsureCleanPNG("/path/to/image.png")
```

## Related Commands

- `generate` - Generate and download images
- `crop` - Auto-crop whitespace borders
- `downscale` - Resize images with high quality

## Full Documentation

See [PNG_METADATA_STRIPPING.md](../PNG_METADATA_STRIPPING.md) for complete details.
