# Gotrace Integration

## Overview

The `gotrace` conversion method provides edge-tracing vector conversion using the pure-Go [dennwc/gotrace](https://github.com/dennwc/gotrace) library, which implements the potrace algorithm natively in Go.

## Implementation

This tool uses a **pure-Go implementation** of the potrace algorithm:
- **No external dependencies required** - potrace binary is not needed
- **Cross-platform compatible** - works on all platforms Go supports
- **Native integration** - direct Go library calls for better performance

The implementation leverages `github.com/dennwc/gotrace` which provides the full potrace functionality as a Go library, eliminating the need for system-level potrace installation.

## Usage

```bash
# Basic edge tracing conversion
asset-generator convert svg sketch.png --method gotrace

# With custom output path
asset-generator convert svg lineart.png --method gotrace -o vectorized.svg
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

## Comparison: Primitive vs Gotrace

| Aspect | Primitive | Gotrace |
|--------|-----------|---------|
| **Best for** | Photos, artistic effects | Line art, sketches |
| **Processing** | Geometric approximation | Edge tracing |
| **Output style** | Geometric, stylized | Smooth curves, precise |
| **Speed** | Slower (many shapes) | Generally faster |
| **Dependencies** | None | None (pure-Go) |
| **File size** | Varies with shape count | Depends on edge complexity |
| **Quality control** | --shapes, --mode flags | Library defaults |

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

## Technical Details

### Conversion Process

1. **Input Loading**: PNG/JPEG image is loaded
2. **Image Decoding**: Image decoded using Go's image libraries
3. **Bitmap Conversion**: Image converted to bitmap for tracing
4. **Edge Tracing**: dennwc/gotrace library traces edges using potrace algorithm
5. **Path Generation**: Vector paths are generated from traced edges
6. **SVG Output**: Final SVG file is written with traced paths

### Format Support

The gotrace method uses Go's standard image libraries for input. Any image format supported by Go (PNG, JPEG, GIF, etc.) can be used as input. The pure-Go implementation eliminates the need for intermediate file format conversions.

## See Also

- [SVG Conversion Guide](SVG_CONVERSION.md)
- [dennwc/gotrace Library](https://github.com/dennwc/gotrace)
- [Primitive Method Documentation](SVG_CONVERSION.md#primitive-method-default)
