#!/bin/bash
# Demo script for filename template feature
# This demonstrates the new --filename-template functionality

set -e

echo "=== Filename Template Demo ==="
echo ""
echo "This script demonstrates various filename template patterns."
echo "Note: Requires a running SwarmUI instance to actually download images."
echo ""

PROMPT="test image"
OUTPUT_DIR="./demo-output"

# Clean up previous demo output
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

echo "1. Testing basic index template..."
./asset-generator generate image \
  --prompt "$PROMPT" \
  --batch 3 \
  --save-images \
  --output-dir "$OUTPUT_DIR/basic" \
  --filename-template "image-{index}.png" \
  --dry-run 2>/dev/null || echo "  (Would create: image-000.png, image-001.png, image-002.png)"

echo ""
echo "2. Testing seed-based naming..."
./asset-generator generate image \
  --prompt "$PROMPT" \
  --seed 42 \
  --batch 2 \
  --save-images \
  --output-dir "$OUTPUT_DIR/seed" \
  --filename-template "img-seed{seed}-{i1}.png" \
  --dry-run 2>/dev/null || echo "  (Would create: img-seed42-1.png, img-seed42-2.png)"

echo ""
echo "3. Testing timestamp template..."
./asset-generator generate image \
  --prompt "$PROMPT" \
  --save-images \
  --output-dir "$OUTPUT_DIR/time" \
  --filename-template "{datetime}-{index}.png" \
  --dry-run 2>/dev/null || echo "  (Would create: 2024-10-08_14-30-45-000.png)"

echo ""
echo "4. Testing complex multi-placeholder template..."
./asset-generator generate image \
  --prompt "$PROMPT" \
  --model "flux-dev" \
  --width 1024 \
  --height 768 \
  --seed 12345 \
  --save-images \
  --output-dir "$OUTPUT_DIR/complex" \
  --filename-template "{model}-{width}x{height}-seed{seed}-{index}.png" \
  --dry-run 2>/dev/null || echo "  (Would create: flux-dev-1024x768-seed12345-000.png)"

echo ""
echo "5. Testing prompt in filename..."
./asset-generator generate image \
  --prompt "beautiful sunset" \
  --save-images \
  --output-dir "$OUTPUT_DIR/prompt" \
  --filename-template "{prompt}-{index}.png" \
  --dry-run 2>/dev/null || echo "  (Would create: beautiful_sunset-000.png)"

echo ""
echo "=== Demo Complete ==="
echo ""
echo "To test with actual image generation, remove --dry-run flag and ensure SwarmUI is running."
