#!/bin/sh
# Tarot Card Back Generator
# Generates a decorative card back design with variations
#
# Usage: ./generate-card-back.sh [output_dir] [seed]

set -e

OUTPUT_DIR="${1:-./tarot-deck-output/card-backs}"
BASE_SEED="${2:-9999}"

# Card dimensions (matching front cards)
WIDTH=768
HEIGHT=1344

# Generation parameters
STEPS=50
CFG_SCALE=8.0

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Verify dependencies
command -v asset-generator >/dev/null 2>&1 || {
    echo "Error: asset-generator CLI not found." >&2
    exit 1
}

mkdir -p "$OUTPUT_DIR"

echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${BLUE}            Tarot Card Back Generator                      ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Card back design variations
designs=(
    "intricate mandala design, mystical symbols, ornate geometric patterns, celestial stars and moons, deep blue and gold, symmetrical, traditional tarot card back"
    "ornate victorian border, central mystical symbol, celestial sun and moon, decorative filigree, rich purple and gold, traditional tarot card back"
    "sacred geometry pattern, flower of life, merkaba, golden ratio spiral, cosmic energy, deep indigo and silver, mystical tarot card back"
    "art nouveau border, mystical eye in center, vines and flowers, crescent moons, stars, teal and gold, elegant tarot card back"
)

STYLE_SUFFIX="highly detailed, ornate decorative elements, traditional card back design, symmetrical composition, professional quality, no text"
NEGATIVE_PROMPT="text, words, letters, numbers, watermark, signature, blurry, distorted, asymmetrical, modern, photograph"

i=1
for prompt in "${designs[@]}"; do
    echo "${YELLOW}Design $i:${NC} Generating card back variation..."
    
    seed=$((BASE_SEED + i))
    full_prompt="${prompt}, ${STYLE_SUFFIX}"
    
    asset-generator generate image \
        --prompt "$full_prompt" \
        --negative-prompt "$NEGATIVE_PROMPT" \
        --width "$WIDTH" \
        --height "$HEIGHT" \
        --steps "$STEPS" \
        --cfg-scale "$CFG_SCALE" \
        --seed "$seed" \
        --save-images \
        --output-dir "$OUTPUT_DIR" \
        --filename-template "card-back-design-${i}.png" \
        > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        echo "  ${GREEN}✓${NC} Saved to: $OUTPUT_DIR/card-back-design-${i}.png"
    else
        echo "  ${YELLOW}⚠${NC} Failed to generate design $i"
    fi
    
    i=$((i + 1))
done

echo ""
echo "${GREEN}✓ Card back designs complete!${NC}"
echo ""
echo "Generated Designs:"
echo "  1. Mandala with celestial symbols (blue/gold)"
echo "  2. Victorian ornate border (purple/gold)"
echo "  3. Sacred geometry patterns (indigo/silver)"
echo "  4. Art nouveau with mystical eye (teal/gold)"
echo ""
echo "Output Location: $OUTPUT_DIR"
echo ""
echo "${YELLOW}Next Steps:${NC}"
echo "  1. Review designs and choose your favorite"
echo "  2. Rename chosen design to 'card-back-final.png'"
echo "  3. Post-process with the same pipeline as front cards"
echo ""
