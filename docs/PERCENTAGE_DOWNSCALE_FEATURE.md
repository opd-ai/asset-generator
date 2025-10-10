# Percentage-Based Downscaling Feature

## ✅ IMPLEMENTATION STATUS (Updated: October 10, 2025)

| Command | Status | Flag Name | Notes |
|---------|--------|-----------|-------|
| `downscale` | ✅ **FULLY IMPLEMENTED** | `--percentage` / `-p` | Works correctly |
| `pipeline` | ✅ **FULLY IMPLEMENTED** | `--downscale-percentage` | Works correctly |
| `generate image` | ✅ **FULLY IMPLEMENTED** | `--downscale-percentage` | **FIXED** - Flag now registered |

### Bug Fix - October 10, 2025

**Previous Issue**: The `--downscale-percentage` flag was not registered in the `generate image` command.

**Status**: ✅ **FIXED** - Flag has been added to `cmd/generate.go` and is now fully functional.

**Change Made**:
```go
// Added to cmd/generate.go init() function:
generateImageCmd.Flags().Float64Var(&generateDownscalePercentage, 
    "downscale-percentage", 0, "downscale by percentage (1-100, 0=disabled, overrides width/height)")
```

## Overview
Percentage-based downscaling support allows users to specify a percentage (e.g., 50%) instead of calculating exact pixel dimensions. This feature is now **fully working** in all commands: `downscale`, `pipeline`, and `generate image`.

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
- ✅ Added `--percentage` / `-p` flag **(WORKING)**
- Validation ensures percentage and dimensions are mutually exclusive
- Updated help text and examples
- Modified verbose output to show percentage or dimensions

**File: `cmd/pipeline.go`**
- ✅ Added `--downscale-percentage` flag **(WORKING)**
- Passes percentage through pipeline processing

**File: `cmd/generate.go`**
- ✅ Variable `generateDownscalePercentage` defined **(WORKING)**
- ✅ Flag `--downscale-percentage` registered **(FIXED - October 10, 2025)**
- ✅ Passes percentage through `DownloadOptions` to client **(WORKING)**

## Working Usage Examples

### Standalone Downscaling (✅ Fully Working)
```bash
# Downscale to 50% of original size
asset-generator downscale image.png --percentage 50

# Short form
asset-generator downscale image.png -p 50

# Batch downscale by 75%
asset-generator downscale *.jpg -p 75 --in-place
```

### Pipeline Processing (✅ Fully Working)
```bash
# Generate pipeline with 50% downscaling
asset-generator pipeline --file tarot-spec.yaml \
  --width 2048 --height 2048 \
  --downscale-percentage 50
```

### Generate Command (✅ NOW WORKING - FIXED!)
```bash
# Downscale by percentage - NOW WORKS!
asset-generator generate image \
  --prompt "high resolution artwork" \
  --width 2048 --height 2048 \
  --save-images \
  --downscale-percentage 50  # ✅ Results in 1024x1024

# Also works with explicit dimensions
asset-generator generate image \
  --prompt "artwork" \
  --width 2048 --height 2048 \
  --save-images \
  --downscale-width 1024  # ✅ Still works
```

## Command-Line Flags

### Downscale Command (✅ All Working)
- `--width` / `-w`: Target width in pixels
- `--height` / `-l`: Target height in pixels
- `--percentage` / `-p`: ✅ Scale by percentage (1-100, takes precedence over width/height)
- `--filter`: Resampling algorithm (lanczos, bilinear, nearest)
- `--quality`: JPEG quality (1-100)
- `--output-file`: Output path for single file
- `--in-place`: Replace original file(s)

### Pipeline Command (✅ All Working)
- `--downscale-percentage`: ✅ Scale by percentage (0=disabled)
- `--downscale-width`: Target width
- `--downscale-height`: Target height
- `--downscale-filter`: Resampling algorithm

### Generate Command (✅ Now Fully Working - FIXED!)
- `--downscale-width`: ✅ Target width for postprocessing
- `--downscale-height`: ✅ Target height for postprocessing  
- `--downscale-percentage`: ✅ **NOW AVAILABLE** (bug fixed October 10, 2025)
- `--downscale-filter`: ✅ Resampling algorithm

## Benefits
