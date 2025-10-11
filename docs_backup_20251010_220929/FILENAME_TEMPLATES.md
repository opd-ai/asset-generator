# Filename Templates Guide

This guide explains how to use custom filename templates when downloading generated images.

## Overview

By default, downloaded images keep their original filename from the server. However, you can customize filenames using the `--filename-template` flag with various placeholders that get replaced with actual values.

## Basic Usage

```bash
asset-generator generate image \
  --prompt "your prompt" \
  --save-images \
  --filename-template "image-{index}.png"
```

## Placeholders Reference

### Index Placeholders

| Placeholder | Description | Example |
|------------|-------------|---------|
| `{index}` or `{i}` | Zero-padded index (3 digits) | `000`, `001`, `002` |
| `{index1}` or `{i1}` | One-based index (no padding) | `1`, `2`, `3` |

**Example:**
```bash
--filename-template "img-{index}.png"
# Results: img-000.png, img-001.png, img-002.png
```

### Time Placeholders

| Placeholder | Description | Example |
|------------|-------------|---------|
| `{timestamp}` or `{ts}` | Unix timestamp | `1696723200` |
| `{datetime}` or `{dt}` | Full datetime | `2024-10-08_14-30-45` |
| `{date}` | Date only | `2024-10-08` |
| `{time}` | Time only | `14-30-45` |

**Example:**
```bash
--filename-template "{date}-{time}-{i1}.png"
# Results: 2024-10-08-14-30-45-1.png
```

### Generation Parameter Placeholders

| Placeholder | Description | Example |
|------------|-------------|---------|
| `{seed}` | Seed value used | `42` |
| `{model}` | Model name | `flux-dev` |
| `{width}` | Image width | `1024` |
| `{height}` | Image height | `768` |
| `{prompt}` | First 50 chars of prompt (sanitized) | `a_beautiful_landscape` |

**Example:**
```bash
--filename-template "{model}-{width}x{height}-seed{seed}.png"
# Results: flux-dev-1024x768-seed42.png
```

### Original Filename Placeholders

| Placeholder | Description | Example |
|------------|-------------|---------|
| `{original}` | Complete original filename | `image-abc123.png` |
| `{ext}` | Extension only (with dot) | `.png` |

**Example:**
```bash
--filename-template "copy-{original}"
# Results: copy-image-abc123.png
```

## Practical Examples

### Example 1: Sequential Numbering with Seed

Great for tracking which seed generated each image:

```bash
asset-generator generate image \
  --prompt "fantasy castle" \
  --seed 42 \
  --batch 10 \
  --save-images \
  --filename-template "castle-seed{seed}-{i1}.png"
```

**Output:**
```
castle-seed42-1.png
castle-seed42-2.png
...
castle-seed42-10.png
```

### Example 2: Organized by Date and Model

Perfect for daily generation workflows:

```bash
asset-generator generate image \
  --prompt "abstract art" \
  --model "sdxl-turbo" \
  --batch 5 \
  --save-images \
  --filename-template "{date}/{model}-{index}.png"
```

**Output:**
```
2024-10-08/sdxl-turbo-000.png
2024-10-08/sdxl-turbo-001.png
...
```

### Example 3: Descriptive Names with Dimensions

Useful for organizing by aspect ratio:

```bash
asset-generator generate image \
  --prompt "portrait of warrior" \
  --width 768 \
  --height 1024 \
  --save-images \
  --filename-template "{prompt}-{width}x{height}.png"
```

**Output:**
```
portrait_of_warrior-768x1024.png
```

### Example 4: Timestamped with Prompt

Great for archiving and searching:

```bash
asset-generator generate image \
  --prompt "sunset over ocean" \
  --save-images \
  --filename-template "{datetime}-{prompt}.png"
```

**Output:**
```
2024-10-08_14-30-45-sunset_over_ocean.png
```

### Example 5: Complex Template

Combine multiple placeholders for maximum organization:

```bash
asset-generator generate image \
  --prompt "cyberpunk street" \
  --model "flux-dev" \
  --width 1024 \
  --height 768 \
  --seed 12345 \
  --batch 3 \
  --save-images \
  --output-dir ./renders \
  --filename-template "{date}/{model}/{prompt}-{width}x{height}-seed{seed}-{i1}.png"
```

**Output:**
```
./renders/2024-10-08/flux-dev/cyberpunk_street-1024x768-seed12345-1.png
./renders/2024-10-08/flux-dev/cyberpunk_street-1024x768-seed12345-2.png
./renders/2024-10-08/flux-dev/cyberpunk_street-1024x768-seed12345-3.png
```

## Special Behaviors

### Automatic Extension Handling

If your template doesn't include an extension, the original file extension is automatically appended:

```bash
--filename-template "image-{index}"
# If original is image.png: image-000.png
# If original is image.jpg: image-000.jpg
```

To override, include `{ext}` explicitly:
```bash
--filename-template "image-{index}{ext}"
```

### Filename Sanitization

The `{prompt}` placeholder is automatically sanitized for filesystem compatibility:
- Spaces become underscores: `"hello world"` → `hello_world`
- Invalid characters are removed: `"cat/dog"` → `catdog`
- Truncated to 50 characters maximum

### Directory Creation

Templates can include directory separators (`/`):

```bash
--filename-template "{date}/{model}/image-{index}.png"
```

Directories are created automatically.

## Tips and Best Practices

1. **Use zero-padded indices** (`{index}`) for proper file sorting
2. **Include seed for reproducibility** when you need to regenerate images
3. **Add timestamps for archival** to avoid filename collisions
4. **Combine model and dimensions** when testing different configurations
5. **Keep templates short** to avoid exceeding filesystem path limits
6. **Test templates first** with `--batch 1` to verify output

## Combining with Other Flags

Filename templates work with all other flags:

```bash
asset-generator generate image \
  --prompt "fantasy landscape" \
  --model "flux-dev" \
  --width 1024 \
  --height 768 \
  --seed 42 \
  --batch 5 \
  --steps 30 \
  --cfg-scale 7.5 \
  --save-images \
  --output-dir ./renders \
  --filename-template "{model}-{width}x{height}-s{seed}-{index}.png" \
  --output metadata.json \
  --format json
```

This generates images with custom names AND saves the metadata to a JSON file.

## See Also

- [Image Download Guide](IMAGE_DOWNLOAD.md) - General image download documentation
- [README](../README.md) - Main documentation
- [Quick Start](../QUICKSTART.md) - Getting started guide
