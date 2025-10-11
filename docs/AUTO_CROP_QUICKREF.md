# Auto-Crop Quick Reference

## Common Use Cases

### 1. Remove Whitespace from Generated Images
```bash
asset-generator generate image \
  --prompt "logo design" \
  --save-images --auto-crop
```

### 2. Crop Existing Image
```bash
asset-generator crop image.png
```

### 3. Crop Preserving Aspect Ratio
```bash
asset-generator crop photo.jpg --preserve-aspect --in-place
```

### 4. Batch Crop Multiple Images
```bash
asset-generator crop *.png --in-place
```

### 5. Crop Then Downscale
```bash
asset-generator generate image \
  --prompt "artwork" \
  --width 2048 --height 2048 \
  --save-images \
  --auto-crop \
  --downscale-width 1024
```

## Sensitivity Tuning

### Very Light Backgrounds (Default)
```bash
--threshold 250 --tolerance 10
```
Detects: RGB > 240 (very light grays and whites)

### Pure White Only (Aggressive)
```bash
--threshold 255 --tolerance 5
```
Detects: RGB > 250 (almost pure white only)

### Light Gray Backgrounds
```bash
--threshold 230 --tolerance 20
```
Detects: RGB > 210 (light grays and whites)

### Darker Backgrounds (Conservative)
```bash
--threshold 200 --tolerance 30
```
Detects: RGB > 170 (broader range of light colors)

## Flags Quick Reference

### Generate Image Flags
| Flag | Description | Default |
|------|-------------|---------|
| `--auto-crop` | Enable cropping | `false` |
| `--auto-crop-threshold` | Detection threshold | `250` |
| `--auto-crop-tolerance` | Color tolerance | `10` |
| `--auto-crop-preserve-aspect` | Keep aspect ratio | `false` |

### Crop Command Flags
| Flag | Description | Default |
|------|-------------|---------|
| `--threshold` | Detection threshold | `250` |
| `--tolerance` | Color tolerance | `10` |
| `--preserve-aspect` | Keep aspect ratio | `false` |
| `--quality` | JPEG quality | `90` |
| `-o, --output` | Output path | - |
| `-i, --in-place` | Replace original | `false` |

## Examples by Scenario

### Logo with White Background
```bash
asset-generator generate image \
  --prompt "minimalist tech logo, white background" \
  --save-images --auto-crop
```

### Product Photos (Keep Proportions)
```bash
asset-generator generate image \
  --prompt "product photography" \
  --save-images \
  --auto-crop --auto-crop-preserve-aspect
```

### High-Res to Thumbnail Pipeline
```bash
asset-generator generate image \
  --prompt "detailed illustration" \
  --width 2048 --height 2048 \
  --save-images \
  --auto-crop \
  --downscale-width 512 \
  --downscale-filter lanczos
```

### Clean Up Downloaded Images
```bash
# Crop all images in current directory
asset-generator crop *.png --in-place

# Crop with custom settings
asset-generator crop *.jpg \
  --threshold 245 \
  --tolerance 15 \
  --quality 95 \
  --in-place
```

## Troubleshooting

### Too Much Cropped
**Problem**: Important content is being removed  
**Solution**: Lower threshold or increase tolerance
```bash
--threshold 245 --tolerance 15
```

### Not Enough Cropped
**Problem**: Whitespace remains  
**Solution**: Increase threshold or decrease tolerance
```bash
--threshold 252 --tolerance 3
```

### Aspect Ratio Changed
**Problem**: Image proportions look wrong  
**Solution**: Enable aspect ratio preservation
```bash
--preserve-aspect
```

### Colored Background Not Detected
**Note**: Auto-crop is designed for white/light backgrounds only. For colored backgrounds, use manual cropping tools.

## Performance Tips

1. **Batch Processing**: Use wildcards for multiple files
   ```bash
   asset-generator crop *.png --in-place
   ```

2. **Pipeline Order**: Crop before downscaling for best quality
   ```bash
   --auto-crop --downscale-width 1024
   ```

3. **In-Place Mode**: Save disk space when you don't need originals
   ```bash
   --in-place
   ```

## See Also

- [Full Documentation](AUTO_CROP_FEATURE.md)
- [Downscaling Feature](DOWNSCALING_FEATURE.md)
- [Image Download](IMAGE_DOWNLOAD.md)
