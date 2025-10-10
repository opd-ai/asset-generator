#!/bin/sh
# Quick Tarot Deck Demo
# Generates a small sample of cards to demonstrate the pipeline
# without waiting for the full 78-card generation
#
# Usage: ./quick-demo.sh [output_dir]

set -e

OUTPUT_DIR="${1:-./tarot-deck-demo}"
BASE_SEED=42

# Card dimensions
WIDTH=768
HEIGHT=1344

# Generation parameters
STEPS=40
CFG_SCALE=7.5

# Styling
STYLE_SUFFIX="detailed illustration, ornate decorative border, mystical symbols, rich colors, traditional tarot card design, professional quality"
NEGATIVE_PROMPT="blurry, distorted, low quality, text, watermark, signatures, modern elements, photograph"

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Verify dependencies
command -v asset-generator >/dev/null 2>&1 || {
    echo "Error: asset-generator CLI not found." >&2
    echo "Please install asset-generator first." >&2
    exit 1
}

# Create output directory
mkdir -p "$OUTPUT_DIR"/{major-arcana,minor-arcana}

echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${BLUE}         Quick Tarot Deck Demo (5 Sample Cards)           ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "This demo generates 5 sample cards to showcase the pipeline:"
echo "  - 2 Major Arcana cards (The Fool, The Magician)"
echo "  - 3 Minor Arcana cards (Ace of Wands, Ace of Cups, Ace of Swords)"
echo ""
echo "Output Directory: $OUTPUT_DIR"
echo ""

# Function to generate a card
generate_card() {
    prompt="$1"
    output_path="$2"
    card_name="$3"
    seed="$4"
    
    echo "${GREEN}Generating:${NC} $card_name"
    
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
        --model "XE-_Pixel_Flux_-_0-1.safetensors" \
        --output-dir "$(dirname "$output_path")" \
        --filename-template "$(basename "$output_path")" \
        > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        echo "  ${GREEN}✓${NC} Saved to: $output_path"
    else
        echo "  ${YELLOW}⚠${NC} Failed to generate $card_name"
    fi
    echo ""
}

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Major Arcana Samples${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# The Fool
generate_card \
    "tarot card art, The Fool, young traveler at cliff edge, white rose, small dog, mountain peaks, mystical symbols, ornate border, art nouveau style, vibrant colors" \
    "$OUTPUT_DIR/major-arcana/00-the_fool.png" \
    "00 - The Fool" \
    $((BASE_SEED + 0))

# The Magician
generate_card \
    "tarot card art, The Magician, figure with infinity symbol above head, all four suit symbols on table, red robe white undergarment, pointing up and down, mystical atmosphere" \
    "$OUTPUT_DIR/major-arcana/01-the_magician.png" \
    "01 - The Magician" \
    $((BASE_SEED + 1))

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Minor Arcana Samples (Aces)${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Ace of Wands
generate_card \
    "tarot card art, Ace of Wands, hand emerging from cloud holding wooden staff, leaves sprouting, castle in distance, new beginnings and inspiration" \
    "$OUTPUT_DIR/minor-arcana/ace_of_wands.png" \
    "Ace of Wands" \
    $((BASE_SEED + 100))

# Ace of Cups
generate_card \
    "tarot card art, Ace of Cups, hand emerging from cloud holding chalice, dove descending, water overflowing, lotus flowers, love and new emotions" \
    "$OUTPUT_DIR/minor-arcana/ace_of_cups.png" \
    "Ace of Cups" \
    $((BASE_SEED + 120))

# Ace of Swords
generate_card \
    "tarot card art, Ace of Swords, hand emerging from cloud holding upright sword, crown with laurel, mountains and wind, mental clarity and truth" \
    "$OUTPUT_DIR/minor-arcana/ace_of_swords.png" \
    "Ace of Swords" \
    $((BASE_SEED + 140))

echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${GREEN}       Quick Demo Complete!                               ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "Generated 5 sample cards:"
echo "  ✓ 2 Major Arcana cards"
echo "  ✓ 3 Minor Arcana cards (Aces)"
echo ""
echo "Output Location: $OUTPUT_DIR"
echo ""
echo "Review the generated cards:"
echo "  cd $OUTPUT_DIR"
echo "  ls -lh major-arcana/"
echo "  ls -lh minor-arcana/"
echo ""
echo "${YELLOW}Next Steps:${NC}"
echo "  1. Review the sample cards to verify quality and style"
echo "  2. Adjust prompts in tarot-spec.yaml if needed"
echo "  3. Run full generation: ./generate-tarot-deck.sh"
echo "  4. Post-process for multiple formats: ./post-process-deck.sh"
echo ""
echo "${BLUE}Tip:${NC} The full deck generation creates 78 cards and takes 60-90 minutes."
echo ""
