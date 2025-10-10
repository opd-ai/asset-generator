# Percentage-Based Downscaling Feature

## Overview
Added percentage-based downscaling support to simplify image resizing operations. Users can now specify a percentage (e.g., 50%) instead of calculating exact pixel dimensions.

## Implementation Summary

### Core Changes

**File: `pkg/processor/resize.go`**
- Extended `DownscaleOptions` struct with `Percentage float64` field
- Modified `DownscaleImage()` to prioritize percentage over explicit dimensions
- Percentage calculation: `targetDimension = int(sourceDimension * (percentage / 100))`
- Validation: percentage must be between 0 and 100

**File: `pkg/processor/resize_test.go`**
- Added `TestDownscaleByPercentage()` with 5 test cases
- Tests 25%, 50%, and 75% scaling
- Tests error conditions (negative, over 100)
- Verifies aspect ratio preservation

**File: `pkg/client/client.go`**
- Extended `DownloadOptions` struct with `DownscalePercentage float64` field
- Updated `applyDownscale()` to pass percentage to processor
- Enhanced verbose logging to show percentage when used

**File: `cmd/downscale.go`**
- Added `--percentage` / `-p` flag
- Validation ensures percentage and dimensions are mutually exclusive
- Updated help text and examples
- Modified verbose output to show percentage or dimensions

**File: `cmd/generate.go`**
- Added `--downscale-percentage` flag for postprocessing
- Passes percentage through `DownloadOptions` to client

**Documentation Updates:**
- `README.md`: Added percentage examples, updated flags table
- `DOWNSCALING_FEATURE.md`: Added percentage documentation and examples
- Both docs include standalone and postprocessing usage

## Usage Examples

### Standalone Downscaling
```bash
# Downscale to 50% of original size
asset-generator downscale image.png -p 50

# Batch downscale by 75%
asset-generator downscale *.jpg -p 75 --in-place
```

### Postprocessing During Generation
```bash
# Generate at 2048x2048, save at 1024x1024 (50%)
asset-generator generate image \
  --prompt "high resolution artwork" \
  --width 2048 --height 2048 \
  --save-images \
  --downscale-percentage 50
```

## Command-Line Flags

### Downscale Command
- `--width` / `-w`: Target width in pixels
- `--height` / `-l`: Target height in pixels (length)
- `--percentage` / `-p`: Scale by percentage (1-100, takes precedence)
- `--filter`: Resampling algorithm (lanczos, bilinear, nearest)
- `--quality`: JPEG quality (1-100)
- `--output-file`: Output path for single file
- `--in-place`: Replace original file(s)

### Generate Command (Postprocessing)
- `--downscale-width`: Target width for postprocessing
- `--downscale-height`: Target height for postprocessing
- `--downscale-percentage`: Scale by percentage (takes precedence)
- `--downscale-filter`: Resampling algorithm

## Benefits

1. **Simplicity**: No need to calculate exact pixel dimensions
2. **Consistency**: Maintains aspect ratio automatically
3. **Intuitive**: "50%" is easier to understand than "1024x1024"
4. **Use Cases**:
   - Thumbnail generation (25% or 50%)
   - Web optimization (75%)
   - Quick downsizing without dimension math

## Testing

All tests pass (21 test cases total):
- ✅ Percentage calculations (50%, 75%, 25%)
- ✅ Error handling (negative, over 100)
- ✅ Aspect ratio preservation
- ✅ All existing dimension-based tests still pass
- ✅ Integration test with real images

## Validation Rules

- Percentage must be between 0 and 100 (exclusive of 0)
- Cannot specify both percentage and dimensions
- If percentage is specified, it takes precedence
- If percentage is 0, falls back to width/height logic
- Still prevents upscaling (target must be smaller than source)

## Notes

- The `-l` shorthand for height follows CLI convention (length)
- The `-w` shorthand for width is standard
- Percentage feature works with all three filters (lanczos, bilinear, nearest)
- JPEG quality setting still applies when using percentage
