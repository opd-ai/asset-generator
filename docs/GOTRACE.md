# Gotrace Integration

## Overview

The `gotrace` conversion method provides edge-tracing vector conversion using [potrace](http://potrace.sourceforge.net/), a powerful tool for transforming bitmap images into smooth, scalable vector graphics.

## Installation

The gotrace method requires the `potrace` command-line tool to be installed on your system.

### Ubuntu/Debian

```bash
sudo apt-get update
sudo apt-get install potrace
```

### macOS

```bash
brew install potrace
```

### Fedora/RHEL/CentOS

```bash
sudo dnf install potrace
```

### Arch Linux

```bash
sudo pacman -S potrace
```

### From Source

```bash
# Download from http://potrace.sourceforge.net/
wget http://potrace.sourceforge.net/download/1.16/potrace-1.16.tar.gz
tar xzf potrace-1.16.tar.gz
cd potrace-1.16
./configure
make
sudo make install
```

## Verification

Verify potrace is installed:

```bash
potrace --version
```

You should see output like:
```
potrace 1.16. Copyright (C) 2001-2019 Peter Selinger.
```

## Usage

Once potrace is installed, you can use the gotrace method:

```bash
# Basic edge tracing conversion
asset-generator convert svg sketch.png --method gotrace

# With custom output path
asset-generator convert svg lineart.png --method gotrace -o vectorized.svg

# Pass additional potrace arguments
asset-generator convert svg image.png --method gotrace --gotrace-args="--turdsize,2,--alphamax,1"
```

## Potrace Arguments

You can pass additional arguments to potrace using the `--gotrace-args` flag. Common options:

### Quality Options

- `--turdsize <n>`: Suppress speckles of up to n pixels (default: 2)
- `--alphamax <n>`: Corner threshold (default: 1.0)
- `--opticurve`: Optimize Bezier curves
- `--opttolerance <n>`: Curve optimization tolerance (default: 0.2)

### Output Options

- `--flat`: Use only straight line segments
- `--tight`: Minimize margins

### Example with Custom Args

```bash
# High quality with speckle suppression
asset-generator convert svg drawing.png \
  --method gotrace \
  --gotrace-args="--turdsize,4,--opticurve,--alphamax,0.8"

# Simplified output with straight lines only
asset-generator convert svg sketch.png \
  --method gotrace \
  --gotrace-args="--flat,--turdsize,10"
```

## Best Practices

### When to Use Gotrace

- **Line art**: Sketches, technical drawings
- **High contrast images**: Black and white artwork
- **Text and logos**: Clean vector conversion
- **Scanned documents**: Converting paper scans to vectors

### When to Use Primitive Instead

- **Photographs**: Natural images, portraits
- **Gradients**: Smooth color transitions
- **Artistic effects**: Stylized, geometric look

### Image Preparation

For best gotrace results:

1. **Increase contrast**: Use image editing tools to enhance contrast
2. **Remove noise**: Clean up artifacts and speckles
3. **Use grayscale**: Color images are automatically converted to grayscale
4. **High resolution**: Higher resolution inputs produce better traces

### Optimization Tips

```bash
# For simple line drawings
asset-generator convert svg sketch.png \
  --method gotrace \
  --gotrace-args="--turdsize,2,--flat"

# For detailed artwork
asset-generator convert svg artwork.png \
  --method gotrace \
  --gotrace-args="--turdsize,1,--opticurve,--alphamax,0.5"

# For logos (clean and tight)
asset-generator convert svg logo.png \
  --method gotrace \
  --gotrace-args="--turdsize,5,--tight"
```

## Comparison: Primitive vs Gotrace

| Aspect | Primitive | Gotrace |
|--------|-----------|---------|
| **Best for** | Photos, artistic effects | Line art, sketches |
| **Processing** | Geometric approximation | Edge tracing |
| **Output style** | Geometric, stylized | Smooth curves, precise |
| **Speed** | Slower (many shapes) | Generally faster |
| **Dependencies** | None | Requires potrace |
| **File size** | Varies with shape count | Depends on edge complexity |
| **Quality control** | --shapes, --mode flags | Potrace arguments |

## Troubleshooting

### "potrace not found in PATH"

**Problem**: The potrace binary is not installed or not in your system PATH.

**Solution**: Install potrace using your package manager (see Installation section above).

### "failed to convert to PBM"

**Problem**: The input image format is not supported or is corrupted.

**Solution**: 
- Verify the image file is valid
- Try converting it to PNG first using image editing tools
- Check file permissions

### Poor trace quality

**Problem**: Output SVG doesn't match the input well.

**Solution**:
- Increase contrast of input image
- Adjust `--turdsize` to suppress noise
- Try different `--alphamax` values (0.0 to 1.333)
- Use `--opticurve` for smoother curves

### Inverted colors

**Problem**: Black and white are swapped in the output.

**Solution**: Potrace traces black regions. If needed, invert your input image first using an image editor.

## Advanced Examples

### Batch Conversion

Convert multiple images using shell scripting:

```bash
# Convert all line art images
for file in lineart_*.png; do
    asset-generator convert svg "$file" \
      --method gotrace \
      --gotrace-args="--turdsize,2,--opticurve"
done
```

### Pipeline with Image Download

```bash
# Generate and convert to SVG
asset-generator generate image \
  --prompt "simple line drawing of a cat" \
  --save-images \
  --output-dir ./generated

asset-generator convert svg ./generated/image_001.png \
  --method gotrace \
  -o cat-vector.svg
```

### Quality Comparison Script

```bash
#!/bin/bash
INPUT="sketch.png"

# Try different quality settings
for turdsize in 1 2 5 10; do
    asset-generator convert svg "$INPUT" \
      --method gotrace \
      --gotrace-args="--turdsize,$turdsize" \
      -o "output_turd${turdsize}.svg"
done
```

## Technical Details

### Conversion Process

1. **Input Loading**: PNG/JPEG image is loaded
2. **Grayscale Conversion**: Image converted to grayscale
3. **Threshold**: Converted to binary (black/white)
4. **PBM Generation**: Temporary PBM file created
5. **Potrace Execution**: Potrace traces the bitmap
6. **SVG Output**: Final SVG file is generated
7. **Cleanup**: Temporary files removed

### Format Support

The gotrace method internally converts images to PBM (Portable Bitmap) format before passing to potrace. This means any image format supported by Go's image libraries (PNG, JPEG, GIF, etc.) can be used as input.

## See Also

- [SVG Conversion Guide](SVG_CONVERSION.md)
- [Potrace Documentation](http://potrace.sourceforge.net/potracelib.pdf)
- [Primitive Method Documentation](SVG_CONVERSION.md#primitive-method-default)
