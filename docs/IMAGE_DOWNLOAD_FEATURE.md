# Image Download Feature - Implementation Summary

## Overview
Added the ability to automatically download generated images from the SwarmUI server and save them locally to disk. This feature is particularly useful for preserving generated images, working with them offline, and organizing collections of generated art.

## Changes Made

### 1. Command-Line Interface (`cmd/generate.go`)
- Added `--save-images` flag to enable image downloading
- Added `--output-dir` flag to specify download directory (defaults to current directory)
- Updated command help text with new examples
- Integrated download logic into the generation workflow
- Added progress feedback for download operations

### 2. Client Library (`pkg/client/client.go`)
- Implemented `DownloadImages()` method for downloading multiple images
- Implemented `downloadFile()` helper for individual file downloads
- Implemented `ensureDir()` utility for directory creation
- Added proper error handling for partial failures
- Support for authentication via API key headers
- Added local file paths to result metadata

### 3. Testing (`pkg/client/download_test.go`)
- Created comprehensive test suite covering:
  - Single and multiple image downloads
  - Empty path handling
  - Invalid path error cases
  - Directory creation scenarios
  - Authentication with API keys
  - Partial failure handling
- All tests pass successfully

### 4. Documentation
Updated the following documentation files:

#### README.md
- Added image download to feature list
- Added examples for downloading images
- Added dedicated "Image Download Feature" section
- Updated usage examples with batch download scenarios
- Added flags documentation for `--save-images` and `--output-dir`

#### QUICKSTART.md
- Added "Download Generated Images" section
- Included basic and advanced examples
- Demonstrated batch download workflows

## Usage Examples

### Basic Download
```bash
asset-generator generate image \
  --prompt "beautiful landscape" \
  --save-images
```

### Download to Specific Directory
```bash
asset-generator generate image \
  --prompt "fantasy castle" \
  --save-images \
  --output-dir ./my-art
```

### Batch Download
```bash
asset-generator generate image \
  --prompt "abstract art" \
  --batch 5 \
  --save-images \
  --output-dir ./batch-output
```

### Downscale After Download (Postprocessing)
The tool supports local postprocessing to downscale images using high-quality Lanczos filtering:

```bash
# Generate large images and downscale to 1024px width
asset-generator generate image \
  --prompt "detailed artwork" \
  --width 2048 --height 2048 \
  --save-images \
  --downscale-width 1024

# Downscale with specific height (auto-calculates width)
asset-generator generate image \
  --prompt "portrait" \
  --width 1920 --height 1080 \
  --save-images \
  --downscale-height 720

# Use different downscaling filters
asset-generator generate image \
  --prompt "photo" \
  --save-images \
  --downscale-width 800 \
  --downscale-filter lanczos  # Options: lanczos (best), bilinear, nearest
```

**Downscaling Features:**
- Applied locally after downloading (no API bandwidth wasted)
- Uses Lanczos3 resampling by default for highest quality
- Maintains aspect ratio when only one dimension is specified
- Prevents accidental upscaling (will error if target > source)
- Replaces the original downloaded file with the downscaled version
- Supports multiple filter algorithms: `lanczos`, `bilinear`, `nearest`

## Technical Details

### Image Path Handling
- SwarmUI returns paths like: `View/local/raw/2024-05-19/filename.png`
- These are converted to full URLs: `http://localhost:7801/View/local/raw/2024-05-19/filename.png`
- Original filenames are preserved when saving locally
- Directory structure is created automatically if needed

### Error Handling
- Graceful handling of partial failures (some images succeed, some fail)
- Clear error messages with context
- Non-blocking: successful downloads are preserved even if some fail
- Proper cleanup and resource management

### Performance
- Uses the existing HTTP client with appropriate timeout settings
- Downloads are sequential to avoid overwhelming the server
- Progress feedback shows each downloaded file

### Authentication
- Respects API key configuration when downloading
- Uses Bearer token authentication if API key is set
- Compatible with both authenticated and unauthenticated servers

## Testing
All tests pass:
- ✅ Single image download
- ✅ Batch downloads
- ✅ Empty path handling
- ✅ Invalid path errors
- ✅ Directory creation
- ✅ Existing directory handling
- ✅ Path-is-file error case
- ✅ Nested directory creation
- ✅ Authentication with API keys
- ✅ Partial failure scenarios

## Compatibility
- Works with both HTTP and WebSocket generation modes
- Compatible with SwarmUI API standard
- No breaking changes to existing functionality
- Backward compatible (feature is opt-in via flag)

## Future Enhancements
Potential future improvements:
- Parallel downloads for batch operations
- Progress bars for large file downloads
- Custom filename patterns
- Automatic deduplication
- Resume capability for interrupted downloads
- Image format conversion options
