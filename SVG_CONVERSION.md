# SVG Conversion Feature

## Overview

The asset-generator CLI now includes powerful image-to-SVG conversion capabilities using two complementary libraries:

- **[fogleman/primitive](https://github.com/fogleman/primitive)**: Converts images to SVG using geometric shape approximation
- **[dennwc/gotrace](https://github.com/dennwc/gotrace)**: Pure-Go implementation of potrace for edge-tracing vector conversion

## Quick Start

Convert an image to SVG with default settings (100 triangles):

```bash
asset-generator convert svg input.png
```

This creates `input.svg` with a geometric approximation of your image.

## Conversion Methods

### 1. Primitive Method (Default)

Uses geometric shapes to approximate the image, creating artistic, simplified versions.

**Best for:**
- Logos and branding materials
- Illustrations and artistic renders
- Creating stylized versions of photos
- Modern, minimalist graphics

**Example:**
```bash
# Basic conversion with 100 shapes
asset-generator convert svg photo.jpg

# High quality with 500 shapes
asset-generator convert svg photo.jpg --shapes 500

# Use ellipses for softer look
asset-generator convert svg photo.jpg --shapes 200 --mode 3
```

**Shape Modes:**
- `0` - Combo (mix of shapes)
- `1` - Triangle (default - good balance)
- `2` - Rectangle
- `3` - Ellipse (softer, organic)
- `4` - Circle
- `5` - Rotated Rectangle
- `6` - Beziers (smooth curves)
- `7` - Rotated Ellipse
- `8` - Polygon

### 2. Gotrace Method

Uses edge tracing to preserve fine details and create precise vector conversions.

**Best for:**
- Line art and sketches
- High-contrast images
- Technical drawings
- Preserving fine details

**Implementation:**
- Uses pure-Go `dennwc/gotrace` library
- No external dependencies required
- Cross-platform compatible

**Example:**
```bash
# Basic gotrace conversion
asset-generator convert svg sketch.png --method gotrace
```

## Command Reference

### Basic Usage

```bash
asset-generator convert svg <input-file> [flags]
```

### Flags

#### General Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--method` | `-m` | `primitive` | Conversion method: `primitive` or `gotrace` |
| `--output` | `-o` | `<input>.svg` | Output file path |
| `--quiet` | `-q` | `false` | Quiet mode (errors only) |

#### Primitive Method Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--shapes` | `100` | Number of shapes to use |
| `--mode` | `1` | Shape mode (0-8, see modes above) |
| `--alpha` | `128` | Alpha transparency (0-255) |
| `--repeat` | `0` | Optimization repeats per shape |

#### Gotrace Method Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--gotrace-args` | `[]` | Additional potrace arguments |

## Examples

### Basic Conversions

```bash
# Default conversion (100 triangles)
asset-generator convert svg image.png

# Specify output location
asset-generator convert svg image.png -o /path/to/output.svg

# Quiet mode (no progress output)
asset-generator convert svg image.png -q
```

### Quality Variations

```bash
# Low quality, fast (30 shapes)
asset-generator convert svg image.png --shapes 30

# Medium quality (100 shapes - default)
asset-generator convert svg image.png --shapes 100

# High quality (500 shapes)
asset-generator convert svg image.png --shapes 500

# Ultra high quality (1000 shapes)
asset-generator convert svg image.png --shapes 1000
```

### Artistic Styles

```bash
# Geometric triangles (default)
asset-generator convert svg photo.jpg --shapes 200 --mode 1

# Soft ellipses
asset-generator convert svg photo.jpg --shapes 200 --mode 3

# Circles for pointillist effect
asset-generator convert svg photo.jpg --shapes 300 --mode 4

# Smooth bezier curves
asset-generator convert svg photo.jpg --shapes 150 --mode 6

# Mixed shapes
asset-generator convert svg photo.jpg --shapes 200 --mode 0
```

### Transparency Effects

```bash
# High transparency (lighter, more layered)
asset-generator convert svg photo.jpg --shapes 200 --alpha 64

# Medium transparency (default)
asset-generator convert svg photo.jpg --shapes 200 --alpha 128

# Low transparency (more opaque)
asset-generator convert svg photo.jpg --shapes 200 --alpha 200
```

### Gotrace Conversions

```bash
# Basic trace conversion
asset-generator convert svg lineart.png --method gotrace

# High quality trace
asset-generator convert svg sketch.png --method gotrace --gotrace-args="--turdsize,2"

# Optimized for smoothness
asset-generator convert svg drawing.png --method gotrace --gotrace-args="--opticurve"
```

## Batch Processing

Convert multiple images using shell scripting:

```bash
# Convert all PNGs in a directory
for file in *.png; do
    asset-generator convert svg "$file" --shapes 150
done

# Convert with custom output directory
mkdir -p svg_output
for file in *.png; do
    asset-generator convert svg "$file" -o "svg_output/${file%.png}.svg"
done

# Try different quality levels
for shapes in 50 100 200 500; do
    asset-generator convert svg input.png -o "output_${shapes}.svg" --shapes $shapes
done
```

## Integration with Image Download

Combine with the image download feature to generate and convert images:

```bash
# Generate an image and convert to SVG
asset-generator generate image --prompt "mountain landscape" --save-images --output-dir ./generated

# Convert the downloaded image
asset-generator convert svg ./generated/image_001.png --shapes 200
```

## Performance Considerations

### Primitive Method

- **Shapes vs Quality**: More shapes = better quality but larger file size and longer processing time
- **Shape Complexity**: Beziers and polygons are more complex than triangles
- **Optimization**: Use `--repeat` flag for better shape placement (slower but better quality)

**Processing Time Examples** (approximate):
- 50 shapes: ~1-2 seconds
- 100 shapes: ~2-4 seconds  
- 200 shapes: ~4-8 seconds
- 500 shapes: ~10-20 seconds
- 1000 shapes: ~20-40 seconds

### Gotrace Method

- **Pre-processing**: Image is converted to bitmap format first
- **Detail Level**: High contrast images trace better
- **Speed**: Generally faster than primitive with many shapes
- **External Dependency**: Requires potrace installation

## File Size Considerations

SVG file sizes vary based on:

- **Number of shapes**: More shapes = larger files
- **Shape complexity**: Beziers > Ellipses > Triangles
- **Image complexity**: Complex images need more shapes

**Typical File Sizes:**
- 50 shapes: 2-5 KB
- 100 shapes: 4-10 KB
- 200 shapes: 8-20 KB
- 500 shapes: 20-50 KB
- 1000 shapes: 50-100 KB

## Tips and Best Practices

1. **Start with defaults**: Try default settings first, then adjust
2. **Match method to content**: Use primitive for photos, gotrace for line art
3. **Balance quality vs size**: More shapes isn't always better
4. **Experiment with modes**: Different shape modes create different aesthetics
5. **Use appropriate alpha**: Higher alpha for solid look, lower for layered effect
6. **Consider final use**: Web graphics can use fewer shapes than print

## Troubleshooting

### "input file does not exist"
- Check file path is correct
- Ensure file extension is included

### "potrace not found in PATH" (gotrace method)
- Install potrace: `sudo apt-get install potrace` (Ubuntu/Debian)
- Or use primitive method instead

### Output looks pixelated
- Increase `--shapes` value
- Try different `--mode` settings
- Consider using gotrace for line art

### Processing takes too long
- Reduce `--shapes` value
- Set `--repeat` to 0 (default)
- Use simpler shape modes (triangles, rectangles)

### Output file is too large
- Reduce `--shapes` value
- Use simpler shape modes
- Consider if SVG is the right format for your use case

## API Usage

Use the converter package directly in your Go code:

```go
package main

import (
    "github.com/opd-ai/asset-generator/pkg/converter"
)

func main() {
    // Create converter
    conv := converter.NewSVGConverter()
    
    // Primitive conversion
    opts := converter.ConversionOptions{
        Method:          converter.MethodPrimitive,
        OutputPath:      "output.svg",
        PrimitiveShapes: 150,
        PrimitiveMode:   1,
    }
    
    result, err := conv.ConvertToSVG("input.png", opts)
    if err != nil {
        panic(err)
    }
    
    println("Converted:", result)
}
```

### Convenience Functions

```go
// Quick primitive conversion
result, err := converter.ConvertWithPrimitiveDefault("input.png", "output.svg", 100)

// Quick gotrace conversion
result, err := converter.ConvertWithGotraceDefault("input.png", "output.svg")
```

## See Also

- [Image Download Feature](IMAGE_DOWNLOAD.md)
- [Custom Filenames](CUSTOM_FILENAMES.md)
- [API Documentation](API.md)

## License


This feature uses:
- [fogleman/primitive](https://github.com/fogleman/primitive) - MIT License
- [dennwc/gotrace](https://github.com/dennwc/gotrace) - BSD-2-Clause License

```
