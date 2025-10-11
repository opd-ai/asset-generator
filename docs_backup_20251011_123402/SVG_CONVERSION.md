````markdown
[🏠 Docs Home](README.md) | [📚 Quick Start](QUICKSTART.md) | [🔧 Postprocessing](POSTPROCESSING.md) | [🔗 Pipeline](PIPELINE.md)

---

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
- Text and logos requiring clean vector conversion
- Scanned documents

**Implementation:**
This tool uses a **pure-Go implementation** of the potrace algorithm:
- Uses `github.com/dennwc/gotrace` library - pure-Go implementation
- **No external dependencies required** - potrace binary is not needed
- **Cross-platform compatible** - works on all platforms Go supports
- **Native integration** - direct Go library calls for better performance

The implementation leverages `github.com/dennwc/gotrace` which provides the full potrace functionality as a Go library, eliminating the need for system-level potrace installation.

**Example:**
```bash
# Basic gotrace conversion
asset-generator convert svg sketch.png --method gotrace

# With custom output path
asset-generator convert svg lineart.png --method gotrace -o vectorized.svg
```

**Image Preparation for Best Results:**

For optimal gotrace results:
1. **Increase contrast**: Use image editing tools to enhance contrast
2. **Remove noise**: Clean up artifacts and speckles
3. **Use grayscale**: Color images are automatically converted to grayscale
4. **High resolution**: Higher resolution inputs produce better traces

**Comparison: Primitive vs Gotrace**

| Aspect | Primitive | Gotrace |
|--------|-----------|---------|
| **Best for** | Photos, artistic effects | Line art, sketches |
| **Processing** | Geometric approximation | Edge tracing |
| **Output style** | Geometric, stylized | Smooth curves, precise |
| **Speed** | Slower (many shapes) | Generally faster |
| **Dependencies** | None | None (pure-Go) |
| **File size** | Varies with shape count | Depends on edge complexity |
| **Quality control** | --shapes, --mode flags | Library defaults |

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

# Gotrace method (edge tracing)
asset-generator convert svg sketch.png --method gotrace -o vectorized.svg
```

## Comprehensive Examples

This section provides detailed examples organized by use case and technique.

### Quality Levels

#### Quick Preview (Fast)
```bash
# 30 shapes - good for previewing the result
asset-generator convert svg photo.jpg --shapes 30 -o preview.svg
```

#### Low Quality
```bash
# 50 shapes - stylized, abstract look
asset-generator convert svg photo.jpg --shapes 50
```

#### Medium Quality (Default)
```bash
# 100 shapes - balanced quality and file size
asset-generator convert svg photo.jpg --shapes 100
```

#### High Quality
```bash
# 300 shapes - detailed representation
asset-generator convert svg photo.jpg --shapes 300
```

#### Ultra High Quality
```bash
# 1000 shapes - very detailed, larger file
asset-generator convert svg photo.jpg --shapes 1000
```

### Shape Mode Examples

#### Mode 0: Combo (Mixed Shapes)
```bash
# Automatically selects best shape type for each iteration
asset-generator convert svg image.png --shapes 200 --mode 0
```

#### Mode 1: Triangles (Default)
```bash
# Clean, geometric look with good detail
asset-generator convert svg image.png --shapes 150 --mode 1
```

#### Mode 2: Rectangles
```bash
# Blocky, mosaic effect
asset-generator convert svg image.png --shapes 150 --mode 2
```

#### Mode 3: Ellipses
```bash
# Soft, organic appearance
asset-generator convert svg image.png --shapes 150 --mode 3
```

#### Mode 4: Circles
```bash
# Pointillist, dotted effect
asset-generator convert svg image.png --shapes 250 --mode 4
```

#### Mode 5: Rotated Rectangles
```bash
# Dynamic, angular style
asset-generator convert svg image.png --shapes 150 --mode 5
```

#### Mode 6: Bezier Curves
```bash
# Smooth, flowing lines
asset-generator convert svg image.png --shapes 100 --mode 6
```

#### Mode 7: Rotated Ellipses
```bash
# Flowing, organic style
asset-generator convert svg image.png --shapes 150 --mode 7
```

#### Mode 8: Polygons
```bash
# Variable-sided polygons
asset-generator convert svg image.png --shapes 120 --mode 8
```

### Transparency Effects Examples

#### High Transparency (Layered Look)
```bash
# Alpha 64 - very transparent, lots of visible overlap
asset-generator convert svg photo.jpg --shapes 200 --alpha 64
```

#### Medium Transparency (Default)
```bash
# Alpha 128 - balanced layering
asset-generator convert svg photo.jpg --shapes 200 --alpha 128
```

#### Low Transparency (Solid Look)
```bash
# Alpha 200 - mostly opaque shapes
asset-generator convert svg photo.jpg --shapes 200 --alpha 200
```

#### Full Opacity
```bash
# Alpha 255 - completely solid shapes
asset-generator convert svg photo.jpg --shapes 200 --alpha 255
```

### Method Comparison Examples

#### Primitive Method (Geometric)
```bash
# Good for photos, artistic effects
asset-generator convert svg photo.jpg --method primitive --shapes 200
```

#### Gotrace Method (Edge Tracing)
```bash
# Good for line art, no external dependencies needed
asset-generator convert svg sketch.png --method gotrace
```

#### Side-by-Side Comparison
```bash
# Generate both versions for comparison
asset-generator convert svg image.png --method primitive --shapes 150 -o primitive.svg
asset-generator convert svg image.png --method gotrace -o gotrace.svg
```

### Use Case Examples

#### Logo Vectorization
```bash
# Convert logo to clean SVG
asset-generator convert svg logo.png \
    --shapes 100 \
    --mode 2 \
    --alpha 255 \
    -o logo-vector.svg
```

#### Artistic Photo Effect
```bash
# Create artistic interpretation of photo
asset-generator convert svg portrait.jpg \
    --shapes 500 \
    --mode 1 \
    --alpha 128 \
    -o portrait-artistic.svg
```

#### Icon Creation
```bash
# Simplified icon from image
asset-generator convert svg icon.png \
    --shapes 30 \
    --mode 4 \
    --alpha 200 \
    -o icon-simple.svg
```

#### Web Graphics
```bash
# Optimized for web with small file size
asset-generator convert svg graphic.png \
    --shapes 75 \
    --mode 3 \
    -o web-graphic.svg
```

#### Print-Quality Vector
```bash
# High quality for printing
asset-generator convert svg artwork.jpg \
    --shapes 1000 \
    --mode 6 \
    --repeat 1 \
    -o print-quality.svg
```

#### Line Art Conversion
```bash
# Precise vector from sketch
asset-generator convert svg sketch.png \
    --method gotrace \
    -o sketch-vector.svg
```

#### Stencil Creation
```bash
# High contrast for stencils
asset-generator convert svg image.png \
    --method gotrace \
    -o stencil.svg
```

#### Watermark Creation
```bash
# Transparent overlay
asset-generator convert svg watermark.png \
    --shapes 50 \
    --alpha 64 \
    -o watermark-overlay.svg
```

### Advanced Techniques

#### Quality vs Speed Trade-off
```bash
# Fast conversion for iteration
asset-generator convert svg test.png --shapes 50 -o quick.svg

# Final high-quality version
asset-generator convert svg test.png --shapes 500 --repeat 1 -o final.svg
```

#### Style Experimentation
```bash
#!/bin/bash
# Try different artistic styles
INPUT="photo.jpg"

# Geometric
asset-generator convert svg "$INPUT" --shapes 200 --mode 1 --alpha 200 -o geometric.svg

# Soft and organic
asset-generator convert svg "$INPUT" --shapes 200 --mode 3 --alpha 100 -o organic.svg

# Pointillist
asset-generator convert svg "$INPUT" --shapes 300 --mode 4 --alpha 150 -o pointillist.svg

# Smooth curves
asset-generator convert svg "$INPUT" --shapes 150 --mode 6 --alpha 128 -o smooth.svg
```

#### Optimization for File Size
```bash
# Minimize shapes while maintaining quality
asset-generator convert svg image.png \
    --shapes 100 \
    --mode 1 \
    --repeat 2 \
    -o optimized.svg
```

#### Create Multiple Quality Versions
```bash
#!/bin/bash
INPUT="photo.jpg"
for shapes in 50 100 200 500; do
    asset-generator convert svg "$INPUT" \
        -o "output_${shapes}shapes.svg" \
        --shapes $shapes
done
```

#### Try All Shape Modes
```bash
#!/bin/bash
INPUT="image.png"
MODES=("combo" "triangle" "rect" "ellipse" "circle" "rotatedrect" "beziers" "rotatedellipse" "polygon")
for i in {0..8}; do
    asset-generator convert svg "$INPUT" \
        -o "mode_${MODES[$i]}.svg" \
        --shapes 150 \
        --mode $i
done
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

# Convert all images in directory (robust)
for file in *.png *.jpg; do
    [ -f "$file" ] || continue
    asset-generator convert svg "$file" --shapes 150
done

# Parallel processing (requires GNU parallel)
find . -name "*.png" | parallel -j4 \
    "asset-generator convert svg {} --shapes 200"

# Memory-efficient batch processing
find . -name "*.png" -print0 | xargs -0 -n 1 -P 2 \
    sh -c 'asset-generator convert svg "$1" --shapes 150 -q' sh
```

## Integration with Image Generation

Combine with the image generation and download features:

### Generate and Convert
```bash
# Generate an AI image and convert to SVG
asset-generator generate image \
    --prompt "minimalist mountain landscape" \
    --save-images \
    --output-dir ./generated

asset-generator convert svg ./generated/image_001.png \
    --shapes 300 \
    --mode 1 \
    -o mountain-vector.svg
```

### Batch Generate and Convert
```bash
#!/bin/bash
# Generate multiple images and convert each
asset-generator generate image \
    --prompt "abstract geometric pattern" \
    --batch 5 \
    --save-images \
    --output-dir ./patterns

for file in ./patterns/*.png; do
    basename="${file%.png}"
    asset-generator convert svg "$file" \
        --shapes 200 \
        --mode 0 \
        -o "${basename}.svg"
done
```

### Pipeline with Custom Filenames
```bash
#!/bin/bash
# Complete pipeline with organized output
PROMPT="fantasy landscape"
OUTPUT_DIR="./svg-art"
mkdir -p "$OUTPUT_DIR"

# Generate
asset-generator generate image \
    --prompt "$PROMPT" \
    --batch 3 \
    --save-images \
    --output-dir "$OUTPUT_DIR" \
    --filename-template "{index}-generated.png"

# Convert each with different styles
for i in 0 1 2; do
    asset-generator convert svg \
        "$OUTPUT_DIR/${i}-generated.png" \
        --shapes 250 \
        --mode 3 \
        -o "$OUTPUT_DIR/${i}-vector.svg"
done
```

## Error Handling

### Graceful Fallback
```bash
#!/bin/bash
# Always use primitive method (no dependencies)
asset-generator convert svg input.png --method primitive --shapes 200
```

### Validation and Retry
```bash
#!/bin/bash
convert_with_retry() {
    local input=$1
    local output=$2
    local max_attempts=3
    
    for ((i=1; i<=max_attempts; i++)); do
        if asset-generator convert svg "$input" -o "$output"; then
            echo "Success on attempt $i"
            return 0
        else
            echo "Attempt $i failed, retrying..."
            sleep 1
        fi
    done
    
    echo "Failed after $max_attempts attempts"
    return 1
}

convert_with_retry "input.png" "output.svg"
```

### Benchmark Different Settings
```bash
#!/bin/bash
INPUT="test.png"

echo "Benchmarking conversion speeds..."

time asset-generator convert svg "$INPUT" --shapes 50 -o bench50.svg -q
time asset-generator convert svg "$INPUT" --shapes 100 -o bench100.svg -q
time asset-generator convert svg "$INPUT" --shapes 200 -o bench200.svg -q
time asset-generator convert svg "$INPUT" --shapes 500 -o bench500.svg -q

ls -lh bench*.svg
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
**Problem**: File path is incorrect or file doesn't exist.

**Solution**: 
- Check file path is correct
- Ensure file extension is included
- Verify file permissions

### Output looks pixelated or lacks detail
**Problem**: Not enough shapes or wrong conversion method.

**Solution**:
- Increase `--shapes` value (try 300-500 for photos)
- Try different `--mode` settings
- Consider using gotrace for line art or high-contrast images
- Ensure input image has sufficient resolution

### Processing takes too long
**Problem**: Too many shapes or complex optimization settings.

**Solution**:
- Reduce `--shapes` value
- Set `--repeat` to 0 (default)
- Use simpler shape modes (triangles, rectangles)
- Start with 50-100 shapes for testing

### Output file is too large
**Problem**: Too many shapes or complex shape modes.

**Solution**:
- Reduce `--shapes` value
- Use simpler shape modes (triangles instead of beziers)
- Consider if SVG is the right format for your use case
- For web use: 100-200 shapes is usually sufficient

### Gotrace produces unexpected results
**Problem**: Input image not suitable for edge tracing.

**Solution**:
- Increase contrast of input image before conversion
- Remove noise and artifacts
- Consider using primitive method for photos/gradients
- Gotrace works best with line art, sketches, and high-contrast images

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

---

## Quick Reference

### Basic Commands

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

### Shape Modes

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

### Quality Presets

| Shapes | Quality | Use Case | Processing Time |
|--------|---------|----------|-----------------|
| 30 | Preview | Quick test | ~1s |
| 50 | Low | Stylized | ~2s |
| 100 | Medium | Default | ~3s |
| 200 | Good | Web graphics | ~6s |
| 300 | High | Print-ready | ~10s |
| 500 | Very High | Detailed | ~15s |
| 1000 | Ultra | Maximum detail | ~30s |

### Common Use Cases

```bash
# Photos to artistic SVG
asset-generator convert svg photo.jpg --shapes 200 --mode 1

# Logo vectorization
asset-generator convert svg logo.png --shapes 100 --mode 2 --alpha 255

# Sketch to vector
asset-generator convert svg sketch.png --method gotrace

# Soft/organic style
asset-generator convert svg image.png --shapes 200 --mode 3 --alpha 100

# Pointillist effect
asset-generator convert svg photo.jpg --shapes 300 --mode 4
```

### File Size Guide

- 50 shapes: ~2-5 KB
- 100 shapes: ~4-10 KB
- 200 shapes: ~8-20 KB
- 500 shapes: ~20-50 KB
- 1000 shapes: ~50-100 KB

### Tips

1. Start with default settings, then adjust
2. Use primitive for photos, gotrace for line art
3. More shapes ≠ always better (diminishing returns)
4. Lower alpha = more layered/transparent effect
5. Different modes create different aesthetics
6. For web: 100-200 shapes is usually sufficient
7. For print: 300-500 shapes recommended

---

## See Also

- [Postprocessing](POSTPROCESSING.md) - Auto-crop, downscaling, metadata stripping
- [Filename Templates](FILENAME_TEMPLATES.md) - Custom filename patterns for image downloads
- [Development Documentation](DEVELOPMENT.md) - Architecture and API integration
- [Project Summary](PROJECT_SUMMARY.md) - High-level project overview

## License

This feature uses:
- [fogleman/primitive](https://github.com/fogleman/primitive) - MIT License
- [dennwc/gotrace](https://github.com/dennwc/gotrace) - BSD-2-Clause License

```
