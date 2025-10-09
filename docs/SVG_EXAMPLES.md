# SVG Conversion Examples

This document provides comprehensive examples of the SVG conversion feature.

## Table of Contents
- [Basic Usage](#basic-usage)
- [Quality Levels](#quality-levels)
- [Shape Modes](#shape-modes)
- [Transparency Effects](#transparency-effects)
- [Method Comparison](#method-comparison)
- [Batch Processing](#batch-processing)
- [Integration Examples](#integration-examples)
- [Use Cases](#use-cases)

## Basic Usage

### Convert with Defaults
```bash
# Creates output.svg with 100 triangular shapes
asset-generator convert svg input.png
```

### Specify Output Path
```bash
asset-generator convert svg photo.jpg -o artwork.svg
```

### Quiet Mode
```bash
asset-generator convert svg image.png -q
```

## Quality Levels

### Quick Preview (Fast)
```bash
# 30 shapes - good for previewing the result
asset-generator convert svg photo.jpg --shapes 30 -o preview.svg
```

### Low Quality
```bash
# 50 shapes - stylized, abstract look
asset-generator convert svg photo.jpg --shapes 50
```

### Medium Quality (Default)
```bash
# 100 shapes - balanced quality and file size
asset-generator convert svg photo.jpg --shapes 100
```

### High Quality
```bash
# 300 shapes - detailed representation
asset-generator convert svg photo.jpg --shapes 300
```

### Ultra High Quality
```bash
# 1000 shapes - very detailed, larger file
asset-generator convert svg photo.jpg --shapes 1000
```

## Shape Modes

### Mode 0: Combo (Mixed Shapes)
```bash
# Automatically selects best shape type for each iteration
asset-generator convert svg image.png --shapes 200 --mode 0
```

### Mode 1: Triangles (Default)
```bash
# Clean, geometric look with good detail
asset-generator convert svg image.png --shapes 150 --mode 1
```

### Mode 2: Rectangles
```bash
# Blocky, mosaic effect
asset-generator convert svg image.png --shapes 150 --mode 2
```

### Mode 3: Ellipses
```bash
# Soft, organic appearance
asset-generator convert svg image.png --shapes 150 --mode 3
```

### Mode 4: Circles
```bash
# Pointillist, dotted effect
asset-generator convert svg image.png --shapes 250 --mode 4
```

### Mode 5: Rotated Rectangles
```bash
# Dynamic, angular style
asset-generator convert svg image.png --shapes 150 --mode 5
```

### Mode 6: Bezier Curves
```bash
# Smooth, flowing lines
asset-generator convert svg image.png --shapes 100 --mode 6
```

### Mode 7: Rotated Ellipses
```bash
# Flowing, organic style
asset-generator convert svg image.png --shapes 150 --mode 7
```

### Mode 8: Polygons
```bash
# Variable-sided polygons
asset-generator convert svg image.png --shapes 120 --mode 8
```

## Transparency Effects

### High Transparency (Layered Look)
```bash
# Alpha 64 - very transparent, lots of visible overlap
asset-generator convert svg photo.jpg --shapes 200 --alpha 64
```

### Medium Transparency (Default)
```bash
# Alpha 128 - balanced layering
asset-generator convert svg photo.jpg --shapes 200 --alpha 128
```

### Low Transparency (Solid Look)
```bash
# Alpha 200 - mostly opaque shapes
asset-generator convert svg photo.jpg --shapes 200 --alpha 200
```

### Full Opacity
```bash
# Alpha 255 - completely solid shapes
asset-generator convert svg photo.jpg --shapes 200 --alpha 255
```

## Method Comparison

### Primitive Method (Geometric)
```bash
# Good for photos, artistic effects
asset-generator convert svg photo.jpg --method primitive --shapes 200
```

### Gotrace Method (Edge Tracing)
```bash
# Good for line art, requires potrace installed
asset-generator convert svg sketch.png --method gotrace
```

### Side-by-Side Comparison
```bash
# Generate both versions for comparison
asset-generator convert svg image.png --method primitive --shapes 150 -o primitive.svg
asset-generator convert svg image.png --method gotrace -o gotrace.svg
```

## Batch Processing

### Convert All Images in Directory
```bash
#!/bin/bash
for file in *.png *.jpg; do
    [ -f "$file" ] || continue
    asset-generator convert svg "$file" --shapes 150
done
```

### Create Multiple Quality Versions
```bash
#!/bin/bash
INPUT="photo.jpg"
for shapes in 50 100 200 500; do
    asset-generator convert svg "$INPUT" \
        -o "output_${shapes}shapes.svg" \
        --shapes $shapes
done
```

### Try All Shape Modes
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

### Parallel Processing
```bash
#!/bin/bash
# Convert multiple images in parallel
find . -name "*.png" | parallel -j4 \
    "asset-generator convert svg {} --shapes 200"
```

## Integration Examples

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

## Use Cases

### Logo Vectorization
```bash
# Convert logo to clean SVG
asset-generator convert svg logo.png \
    --shapes 100 \
    --mode 2 \
    --alpha 255 \
    -o logo-vector.svg
```

### Artistic Photo Effect
```bash
# Create artistic interpretation of photo
asset-generator convert svg portrait.jpg \
    --shapes 500 \
    --mode 1 \
    --alpha 128 \
    -o portrait-artistic.svg
```

### Icon Creation
```bash
# Simplified icon from image
asset-generator convert svg icon.png \
    --shapes 30 \
    --mode 4 \
    --alpha 200 \
    -o icon-simple.svg
```

### Web Graphics
```bash
# Optimized for web with small file size
asset-generator convert svg graphic.png \
    --shapes 75 \
    --mode 3 \
    -o web-graphic.svg
```

### Print-Quality Vector
```bash
# High quality for printing
asset-generator convert svg artwork.jpg \
    --shapes 1000 \
    --mode 6 \
    --repeat 1 \
    -o print-quality.svg
```

### Line Art Conversion
```bash
# Precise vector from sketch
asset-generator convert svg sketch.png \
    --method gotrace \
    --gotrace-args="--turdsize,2,--opticurve" \
    -o sketch-vector.svg
```

### Stencil Creation
```bash
# High contrast for stencils
asset-generator convert svg image.png \
    --method gotrace \
    --gotrace-args="--flat,--turdsize,10" \
    -o stencil.svg
```

### Watermark Creation
```bash
# Transparent overlay
asset-generator convert svg watermark.png \
    --shapes 50 \
    --alpha 64 \
    -o watermark-overlay.svg
```

## Advanced Techniques

### Quality vs Speed Trade-off
```bash
# Fast conversion for iteration
asset-generator convert svg test.png --shapes 50 -o quick.svg

# Final high-quality version
asset-generator convert svg test.png --shapes 500 --repeat 1 -o final.svg
```

### Style Experimentation
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

### Optimization for File Size
```bash
# Minimize shapes while maintaining quality
asset-generator convert svg image.png \
    --shapes 100 \
    --mode 1 \
    --repeat 2 \
    -o optimized.svg
```

## Performance Examples

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

### Memory-Efficient Batch Processing
```bash
#!/bin/bash
# Process large batches without overwhelming memory
find . -name "*.png" -print0 | xargs -0 -n 1 -P 2 \
    sh -c 'asset-generator convert svg "$1" --shapes 150 -q' sh
```

## Error Handling Examples

### Graceful Fallback
```bash
#!/bin/bash
# Try gotrace, fall back to primitive if potrace not installed
if command -v potrace &> /dev/null; then
    asset-generator convert svg input.png --method gotrace
else
    echo "Potrace not found, using primitive method..."
    asset-generator convert svg input.png --method primitive --shapes 200
fi
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

## See Also

- [SVG Conversion Documentation](../SVG_CONVERSION.md)
- [Gotrace Integration](GOTRACE.md)
- [Image Download Feature](IMAGE_DOWNLOAD.md)
