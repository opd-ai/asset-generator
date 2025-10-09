# Image Download Quick Reference

## What it does
The `--save-images` flag downloads generated images from the SwarmUI server directly to your local disk.

## Why use it?
- 💾 Preserve images locally
- 📂 Organize images in custom directories
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

### Batch download
```bash
asset-generator generate image --prompt "your prompt" --batch 5 --save-images --output-dir ./batch
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--save-images` | `false` | Enable image downloading |
| `--output-dir` | `.` | Directory where images will be saved |

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

### Save landscape images to a specific folder
```bash
mkdir -p ~/art/landscapes
asset-generator generate image \
  --prompt "beautiful mountain landscape at sunset" \
  --width 1024 \
  --height 768 \
  --save-images \
  --output-dir ~/art/landscapes
```

### Generate and download multiple variations
```bash
asset-generator generate image \
  --prompt "cyberpunk city street, neon lights" \
  --batch 10 \
  --save-images \
  --output-dir ./cyberpunk-series \
  --seed 42
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
