# SVG Conversion - Quick Reference

## Installation
No additional dependencies needed for primitive method (default). For gotrace method:
```bash
sudo apt-get install potrace  # Ubuntu/Debian
brew install potrace          # macOS
```

## Basic Commands

```bash
# Convert with defaults (100 triangles)
asset-generator convert svg image.png

# Custom output path
asset-generator convert svg image.png -o output.svg

# Specify quality
asset-generator convert svg image.png --shapes 300

# Use gotrace method
asset-generator convert svg sketch.png --method gotrace
```

## Shape Modes

| Mode | Type | Best For |
|------|------|----------|
| 0 | Combo | Automatic selection |
| 1 | Triangle | Default, balanced |
| 2 | Rectangle | Mosaic effect |
| 3 | Ellipse | Soft, organic |
| 4 | Circle | Pointillist |
| 5 | Rotated Rectangle | Angular, dynamic |
| 6 | Bezier | Smooth curves |
| 7 | Rotated Ellipse | Flowing style |
| 8 | Polygon | Variable sides |

## Quality Presets

| Shapes | Quality | Use Case | Processing Time |
|--------|---------|----------|-----------------|
| 30 | Preview | Quick test | ~1s |
| 50 | Low | Stylized | ~2s |
| 100 | Medium | Default | ~3s |
| 200 | Good | Web graphics | ~6s |
| 300 | High | Print-ready | ~10s |
| 500 | Very High | Detailed | ~15s |
| 1000 | Ultra | Maximum detail | ~30s |

## Common Use Cases

### Photos to Artistic SVG
```bash
asset-generator convert svg photo.jpg --shapes 200 --mode 1
```

### Logo Vectorization
```bash
asset-generator convert svg logo.png --shapes 100 --mode 2 --alpha 255
```

### Sketch to Vector
```bash
asset-generator convert svg sketch.png --method gotrace
```

### Soft/Organic Style
```bash
asset-generator convert svg image.png --shapes 200 --mode 3 --alpha 100
```

### Pointillist Effect
```bash
asset-generator convert svg photo.jpg --shapes 300 --mode 4
```

## Flags Reference

### General
- `-m, --method`: Conversion method (`primitive` or `gotrace`)
- `-o, --output`: Output file path
- `-q, --quiet`: Suppress output

### Primitive Method
- `--shapes`: Number of shapes (default: 100)
- `--mode`: Shape mode 0-8 (default: 1)
- `--alpha`: Transparency 0-255 (default: 128)
- `--repeat`: Optimization repeats (default: 0)

### Gotrace Method
- No additional flags (uses library defaults)

## File Size Guide

- 50 shapes: ~2-5 KB
- 100 shapes: ~4-10 KB
- 200 shapes: ~8-20 KB
- 500 shapes: ~20-50 KB
- 1000 shapes: ~50-100 KB

## Tips

1. Start with default settings, then adjust
2. Use primitive for photos, gotrace for line art
3. More shapes ≠ always better (diminishing returns)
4. Lower alpha = more layered/transparent effect
5. Different modes create different aesthetics
6. For web: 100-200 shapes is usually sufficient
7. For print: 300-500 shapes recommended

## Troubleshooting

**Slow processing?** Reduce `--shapes` value

**Large files?** Use fewer shapes or simpler modes

**Not enough detail?** Increase `--shapes` or try different `--mode`

**Need potrace?** `sudo apt-get install potrace`

## Full Documentation

- [Complete Guide](../SVG_CONVERSION.md)
- [Examples](SVG_EXAMPLES.md)
- [Gotrace Details](GOTRACE.md)
