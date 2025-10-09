# Image Download Quick Reference

## What it does
The `--save-images` flag downloads generated images from the SwarmUI server directly to your local disk, with optional custom filename templates.

## Why use it?
- 💾 Preserve images locally
- 📂 Organize images in custom directories
- 🏷️ Use custom filenames with metadata
- 🔄 Work offline with generated images
- 🎨 Build local collections of generated art
- ⚡ Automatic download after generation completes

## Basic Usage

### Download to current directory
```bash
asset-generator generate image --prompt "your prompt" --save-images
```

### Download to specific directory
```bash
asset-generator generate image --prompt "your prompt" --save-images --output-dir ./my-images
```

### Download with custom filenames
```bash
asset-generator generate image --prompt "fantasy landscape" --batch 5 --save-images \
  --filename-template "landscape-{index}-seed{seed}.png"
```

### Batch download
```bash
asset-generator generate image --prompt "your prompt" --batch 5 --save-images --output-dir ./batch
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--save-images` | `false` | Enable image downloading |
| `--output-dir` | `.` | Directory where images will be saved |
| `--filename-template` | (empty) | Template for custom filenames (see placeholders below) |

## Filename Template Placeholders

When using `--filename-template`, you can use these placeholders:

| Placeholder | Description | Example Output |
|------------|-------------|----------------|
| `{index}` or `{i}` | Zero-padded index (001, 002, ...) | `001` |
| `{index1}` or `{i1}` | One-based index (1, 2, 3, ...) | `1` |
| `{timestamp}` or `{ts}` | Unix timestamp | `1696723200` |
| `{datetime}` or `{dt}` | Formatted datetime | `2024-10-08_14-30-45` |
| `{date}` | Date only | `2024-10-08` |
| `{time}` | Time only | `14-30-45` |
| `{seed}` | Seed value | `42` |
| `{model}` | Model name | `flux` |
| `{width}` | Image width | `1024` |
| `{height}` | Image height | `768` |
| `{prompt}` | First 50 chars of prompt (sanitized) | `a_beautiful_landscape` |
| `{original}` | Original filename from server | `image-abc123.png` |
| `{ext}` | Original file extension | `.png` |

**Note:** If no extension is in your template, it will be automatically appended from the original file.

## Output

When enabled, you'll see:
```
Generating image with prompt: your prompt
Downloading generated images...
  [1/3] Saved: ./my-images/image-123456.png
  [2/3] Saved: ./my-images/image-123457.png
  [3/3] Saved: ./my-images/image-123458.png
✓ Generation completed successfully (3 images)
```

## Notes

- Directory is created automatically if it doesn't exist
- Original filenames from the server are preserved
- Works with both single and batch generation
- Compatible with all generation parameters
- Local paths are added to JSON/YAML output metadata
- Partial failures are handled gracefully

## Examples

### Save landscape images with custom filenames
```bash
asset-generator generate image \
  --prompt "beautiful mountain landscape at sunset" \
  --width 1024 \
  --height 768 \
  --save-images \
  --output-dir ~/art/landscapes \
  --filename-template "landscape-{date}-{index}.png"
```

### Generate series with seed in filename
```bash
asset-generator generate image \
  --prompt "cyberpunk city street, neon lights" \
  --batch 10 \
  --save-images \
  --output-dir ./cyberpunk-series \
  --seed 42 \
  --filename-template "cyber-{seed}-{i1}.png"
```

### Organized by model and dimensions
```bash
asset-generator generate image \
  --prompt "portrait of a warrior" \
  --model "flux-dev" \
  --width 768 \
  --height 1024 \
  --batch 5 \
  --save-images \
  --filename-template "{model}-{width}x{height}-{index}.png"
```

### Include prompt in filename
```bash
asset-generator generate image \
  --prompt "red sports car in desert" \
  --batch 3 \
  --save-images \
  --filename-template "{prompt}-{timestamp}-{i1}.png"
# Results in: red_sports_car_in_desert-1696723200-1.png
```

### Download with metadata tracking
```bash
asset-generator generate image \
  --prompt "fantasy character portrait" \
  --save-images \
  --output-dir ./portraits \
  --output metadata.json \
  --format json
```

This saves both the images to `./portraits/` and the metadata (including local paths) to `metadata.json`.
