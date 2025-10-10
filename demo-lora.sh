#!/bin/bash
# LoRA Support Demonstration
# This script showcases various LoRA usage patterns with the asset-generator CLI

set -e

echo "=== LoRA Support Demonstration ==="
echo ""

# Configuration
OUTPUT_DIR="./lora-demo-output"
mkdir -p "$OUTPUT_DIR"

echo "Output directory: $OUTPUT_DIR"
echo ""

# Check if asset-generator is available
if ! command -v asset-generator &> /dev/null; then
    if [ -f "./asset-generator" ]; then
        ASSET_GEN="./asset-generator"
    else
        echo "Error: asset-generator not found in PATH or current directory"
        exit 1
    fi
else
    ASSET_GEN="asset-generator"
fi

echo "Using: $ASSET_GEN"
echo ""

# Demo 1: Single LoRA with default weight
echo "--- Demo 1: Single LoRA (default weight) ---"
echo "Command: generate image --prompt 'anime character' --lora 'anime-style'"
echo ""
# Uncomment to run (requires SwarmUI connection):
# $ASSET_GEN generate image \
#   --prompt "anime character in magical forest" \
#   --lora "anime-style" \
#   --save-images --output-dir "$OUTPUT_DIR/demo1"

# Demo 2: Single LoRA with custom weight
echo "--- Demo 2: Single LoRA with custom weight ---"
echo "Command: generate image --prompt 'portrait' --lora 'realistic-faces:0.8'"
echo ""
# $ASSET_GEN generate image \
#   --prompt "portrait of a young woman, natural lighting" \
#   --lora "realistic-faces:0.8" \
#   --save-images --output-dir "$OUTPUT_DIR/demo2"

# Demo 3: Multiple LoRAs with different weights
echo "--- Demo 3: Multiple LoRAs ---"
echo "Command: generate image --prompt 'cyberpunk city' --lora 'cyberpunk:1.0' --lora 'neon:0.7'"
echo ""
# $ASSET_GEN generate image \
#   --prompt "cyberpunk city street at night, rain, neon signs" \
#   --lora "cyberpunk-aesthetic:1.0" \
#   --lora "neon-lights:0.7" \
#   --lora "detailed-architecture:0.5" \
#   --save-images --output-dir "$OUTPUT_DIR/demo3"

# Demo 4: Negative weight to remove style
echo "--- Demo 4: Negative LoRA weight (style removal) ---"
echo "Command: --lora 'realistic:1.0' --lora 'cartoon:-0.5'"
echo ""
# $ASSET_GEN generate image \
#   --prompt "fantasy warrior character" \
#   --lora "realistic-rendering:1.0" \
#   --lora "cartoon-style:-0.5" \
#   --save-images --output-dir "$OUTPUT_DIR/demo4"

# Demo 5: LoRA with Skimmed CFG
echo "--- Demo 5: LoRA + Skimmed CFG ---"
echo "Command: --lora 'detailed:0.8' --skimmed-cfg --skimmed-cfg-scale 3.0"
echo ""
# $ASSET_GEN generate image \
#   --prompt "detailed fantasy landscape, magical atmosphere" \
#   --lora "fantasy-art:0.9" \
#   --lora "detailed-background:0.6" \
#   --skimmed-cfg --skimmed-cfg-scale 3.0 \
#   --save-images --output-dir "$OUTPUT_DIR/demo5"

# Demo 6: Batch generation with LoRAs
echo "--- Demo 6: Batch generation with LoRAs ---"
echo "Command: --lora 'anime:0.9' --batch 4 --filename-template 'char-{index}.png'"
echo ""
# $ASSET_GEN generate image \
#   --prompt "anime character, different poses" \
#   --lora "anime-style-v2:0.9" \
#   --lora "detailed-clothing:0.7" \
#   --batch 4 \
#   --save-images --output-dir "$OUTPUT_DIR/demo6" \
#   --filename-template "character-{index}-{seed}.png"

# Demo 7: LoRA with complete pipeline
echo "--- Demo 7: LoRA + Complete Pipeline ---"
echo "Command: --lora 'art:0.9' --auto-crop --downscale-width 1024"
echo ""
# $ASSET_GEN generate image \
#   --prompt "high resolution fantasy artwork" \
#   --lora "detailed-art:0.9" \
#   --width 2048 --height 2048 \
#   --save-images --output-dir "$OUTPUT_DIR/demo7" \
#   --auto-crop \
#   --downscale-width 1024 --downscale-filter lanczos

# Demo 8: Custom default weight
echo "--- Demo 8: Custom default weight ---"
echo "Command: --lora 'style1' --lora 'style2' --lora-default-weight 0.7"
echo ""
# $ASSET_GEN generate image \
#   --prompt "artistic portrait" \
#   --lora "artistic-style" \
#   --lora "detailed-features:1.0" \
#   --lora "soft-lighting" \
#   --lora-default-weight 0.7 \
#   --save-images --output-dir "$OUTPUT_DIR/demo8"

echo ""
echo "=== Demonstration Complete ==="
echo ""
echo "Note: All commands are commented out by default."
echo "To run the demos:"
echo "  1. Ensure SwarmUI is running and configured"
echo "  2. Uncomment the desired demo commands in this script"
echo "  3. Adjust LoRA names to match your available models"
echo "  4. Run: bash demo-lora.sh"
echo ""
echo "To list available LoRAs:"
echo "  asset-generator models list --subtype LoRA"
echo ""
echo "For more information, see:"
echo "  docs/LORA_SUPPORT.md - Complete documentation"
echo "  docs/LORA_QUICKREF.md - Quick reference guide"
