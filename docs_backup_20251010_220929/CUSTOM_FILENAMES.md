# Custom Filename Template Feature - Implementation Summary

## Overview
Added the ability to customize downloaded image filenames using template placeholders. This enhancement allows users to organize their generated images with meaningful, structured filenames based on generation parameters and metadata.

## Changes Made

### 1. Core Library (`pkg/client/client.go`)

#### New Types
- **`DownloadOptions`**: New struct to hold download configuration including filename template and metadata
  ```go
  type DownloadOptions struct {
      OutputDir        string
      FilenameTemplate string
      Metadata         map[string]interface{}
  }
  ```

#### New Functions
- **`DownloadImagesWithOptions()`**: Enhanced download method that accepts `DownloadOptions` for customization
- **`generateFilename()`**: Template processing function that replaces placeholders with actual values
- **`sanitizeForFilename()`**: Helper function to ensure filenames are filesystem-safe

#### Modified Functions
- **`DownloadImages()`**: Now wraps `DownloadImagesWithOptions()` for backward compatibility

#### Supported Placeholders
- **Index**: `{index}`, `{i}` (zero-padded), `{index1}`, `{i1}` (one-based)
- **Time**: `{timestamp}`, `{ts}`, `{datetime}`, `{dt}`, `{date}`, `{time}`
- **Parameters**: `{seed}`, `{model}`, `{width}`, `{height}`, `{prompt}`
- **Original**: `{original}`, `{ext}`

### 2. CLI Interface (`cmd/generate.go`)

#### New Flags
- **`--filename-template`**: Optional flag to specify custom filename template
  - Default: empty (uses original server filename)
  - Example: `"image-{index}-seed{seed}.png"`

#### Updated Help Text
- Added comprehensive placeholder documentation to `generate image --help`
- Included practical examples demonstrating filename templates
- Added dedicated section explaining all available placeholders

#### Enhanced Download Logic
- Prepares metadata map from generation parameters
- Conditionally uses `DownloadImagesWithOptions()` when template is specified
- Falls back to standard `DownloadImages()` when no template provided

### 3. Testing (`pkg/client/download_test.go`)

#### New Test Functions
- **`TestDownloadImagesWithTemplate`**: Tests various template patterns
  - Index placeholders
  - Seed and metadata substitution
  - Extension handling
  - Prompt sanitization
  - Multiple placeholder combinations

- **`TestGenerateFilename`**: Unit tests for filename generation logic
  - Zero-padded and one-based indices
  - Timestamp formatting
  - Metadata substitution
  - Extension extraction and auto-append
  - Prompt sanitization

- **`TestSanitizeForFilename`**: Tests filename sanitization
  - Space to underscore conversion
  - Invalid character removal
  - Multiple space handling
  - Leading/trailing underscore trimming

#### Test Coverage
- All tests pass ✅
- Covers edge cases and error scenarios
- Validates filename safety and correctness

### 4. Documentation

#### Updated Files
- **`docs/IMAGE_DOWNLOAD.md`**: 
  - Added placeholder reference table
  - Updated examples with filename templates
  - Expanded usage scenarios

- **`README.md`**:
  - Added "Custom Filenames" section
  - Included practical examples
  - Quick reference to placeholders

#### New Files
- **`docs/FILENAME_TEMPLATES.md`**: Comprehensive guide covering:
  - Complete placeholder reference
  - Practical examples for common use cases
  - Special behaviors (auto-extension, sanitization)
  - Best practices and tips
  - Complex multi-placeholder examples

## Usage Examples

### Basic Template
```bash
asset-generator generate image \
  --prompt "fantasy landscape" \
  --save-images \
  --filename-template "image-{index}.png"
# Output: image-000.png, image-001.png, ...
```

### With Seed
```bash
asset-generator generate image \
  --prompt "portrait" \
  --seed 42 \
  --batch 5 \
  --save-images \
  --filename-template "portrait-seed{seed}-{i1}.png"
# Output: portrait-seed42-1.png, portrait-seed42-2.png, ...
```

### Complex Organization
```bash
asset-generator generate image \
  --prompt "cyberpunk street" \
  --model "flux-dev" \
  --width 1024 \
  --height 768 \
  --save-images \
  --filename-template "{date}/{model}-{width}x{height}-{index}.png"
# Output: 2024-10-08/flux-dev-1024x768-000.png
```

## Technical Details

### Filename Sanitization
- Spaces converted to underscores
- Invalid characters (`/`, `\`, `:`, `*`, `?`, `"`, `<`, `>`, `|`) removed
- Newlines and tabs stripped
- Multiple consecutive underscores collapsed
- Leading/trailing underscores trimmed
- Prompt text truncated to 50 characters

### Extension Handling
- If template lacks extension, original extension is auto-appended
- Use `{ext}` placeholder for explicit control
- Extension includes the dot (e.g., `.png`)

### Backward Compatibility
- Existing `DownloadImages()` function unchanged
- When `--filename-template` is not specified, original behavior is preserved
- No breaking changes to existing functionality

### Performance
- Template processing is lightweight (string replacement)
- No additional network overhead
- Minimal CPU usage for filename generation

## Testing Results

All tests passing:
```
✓ TestDownloadImages (existing tests - still pass)
✓ TestDownloadImagesWithTemplate (6 scenarios)
✓ TestGenerateFilename (9 test cases)
✓ TestSanitizeForFilename (7 test cases)
✓ Full integration test suite passes
```

## Benefits

1. **Better Organization**: Structure downloads by date, model, seed, etc.
2. **Reproducibility**: Include seed in filename for easy regeneration
3. **Searchability**: Meaningful filenames make finding images easier
4. **Automation**: Template-based naming enables scripting workflows
5. **Archival**: Timestamped filenames prevent collisions
6. **Flexibility**: Combine multiple placeholders for complex naming schemes

## Future Enhancements (Potential)

- Custom date/time format strings
- Counter with custom padding length
- Nested directory creation based on metadata
- Template validation before download
- Template presets (e.g., `--template-preset archive`)

## See Also

- [docs/FILENAME_TEMPLATES.md](../docs/FILENAME_TEMPLATES.md) - Complete template guide
- [docs/IMAGE_DOWNLOAD.md](../docs/IMAGE_DOWNLOAD.md) - General download documentation
- [README.md](../README.md) - Main documentation
