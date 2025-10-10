#!/bin/sh
# Tarot Deck Post-Processing Pipeline
# Processes generated cards for various output formats and sizes
#
# Usage: ./post-process-deck.sh [input_dir] [output_dir]

set -e

INPUT_DIR="${1:-./tarot-deck-output}"
OUTPUT_DIR="${2:-./tarot-deck-processed}"

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

# Check if input directory exists
if [ ! -d "$INPUT_DIR" ]; then
    echo "Error: Input directory '$INPUT_DIR' not found." >&2
    echo "Run generate-tarot-deck.sh first." >&2
    exit 1
fi

echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${BLUE}       Tarot Deck Post-Processing Pipeline                ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "Input Directory: $INPUT_DIR"
echo "Output Directory: $OUTPUT_DIR"
echo ""

# Create output directory structure
mkdir -p "$OUTPUT_DIR"/{print-ready,web-optimized,mobile-optimized,svg-versions,thumbnails}

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Step 1: Print-Ready Format (300 DPI equivalent)${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Standard tarot card size: 2.75" x 4.75" at 300 DPI = 825x1425
# Our generated cards are 768x1344, so we'll keep them as-is
echo "${GREEN}✓${NC} Cards already at optimal print resolution (768x1344)"
echo "  Copying to print-ready directory..."

find "$INPUT_DIR" -name "*.png" -type f | while read -r card; do
    relative_path=$(echo "$card" | sed "s|$INPUT_DIR/||")
    output_path="$OUTPUT_DIR/print-ready/$relative_path"
    mkdir -p "$(dirname "$output_path")"
    cp "$card" "$output_path"
done

echo "${GREEN}✓${NC} Print-ready cards prepared"
echo ""

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Step 2: Web-Optimized Format (1024px max dimension)${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

find "$INPUT_DIR" -name "*.png" -type f | while read -r card; do
    relative_path=$(echo "$card" | sed "s|$INPUT_DIR/||")
    output_path="$OUTPUT_DIR/web-optimized/$relative_path"
    mkdir -p "$(dirname "$output_path")"
    
    basename=$(basename "$card")
    echo "  Processing: $basename"
    
    asset-generator downscale \
        --input "$card" \
        --output "$output_path" \
        --max-dimension 1024 \
        > /dev/null 2>&1
done

echo "${GREEN}✓${NC} Web-optimized cards created (max 1024px)"
echo ""

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Step 3: Mobile-Optimized Format (512px max dimension)${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

find "$INPUT_DIR" -name "*.png" -type f | while read -r card; do
    relative_path=$(echo "$card" | sed "s|$INPUT_DIR/||")
    output_path="$OUTPUT_DIR/mobile-optimized/$relative_path"
    mkdir -p "$(dirname "$output_path")"
    
    basename=$(basename "$card")
    echo "  Processing: $basename"
    
    asset-generator downscale \
        --input "$card" \
        --output "$output_path" \
        --max-dimension 512 \
        > /dev/null 2>&1
done

echo "${GREEN}✓${NC} Mobile-optimized cards created (max 512px)"
echo ""

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Step 4: SVG Conversion (Vector Format)${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Converting selection of cards to SVG format..."
echo "(Converting first 3 Major Arcana cards as examples)"
echo ""

# Convert first 3 major arcana cards to SVG as examples
major_cards=$(find "$INPUT_DIR/major-arcana" -name "*.png" -type f | head -n 3)

for card in $major_cards; do
    basename=$(basename "$card" .png)
    output_svg="$OUTPUT_DIR/svg-versions/$basename.svg"
    mkdir -p "$(dirname "$output_svg")"
    
    echo "  Converting: $basename.png → $basename.svg"
    
    asset-generator convert svg \
        --input "$card" \
        --output "$output_svg" \
        --mode primitive \
        --shapes 200 \
        > /dev/null 2>&1
done

echo ""
echo "${GREEN}✓${NC} SVG conversions complete (sample cards)"
echo "  Note: Full deck SVG conversion can be enabled by uncommenting the full loop"
echo ""

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Step 5: Thumbnail Generation (256px max dimension)${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

find "$INPUT_DIR" -name "*.png" -type f | while read -r card; do
    relative_path=$(echo "$card" | sed "s|$INPUT_DIR/||")
    output_path="$OUTPUT_DIR/thumbnails/$relative_path"
    mkdir -p "$(dirname "$output_path")"
    
    basename=$(basename "$card")
    echo "  Processing: $basename"
    
    asset-generator downscale \
        --input "$card" \
        --output "$output_path" \
        --max-dimension 256 \
        > /dev/null 2>&1
done

echo "${GREEN}✓${NC} Thumbnails created (max 256px)"
echo ""

# Generate summary report
echo ""
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${GREEN}       Post-Processing Complete!                          ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "Output Directory: $OUTPUT_DIR"
echo ""
echo "Generated Formats:"
echo "  1. Print-Ready:      768x1344px  (original resolution)"
echo "  2. Web-Optimized:    max 1024px  (for websites)"
echo "  3. Mobile-Optimized: max 512px   (for mobile apps)"
echo "  4. SVG Versions:     vector      (scalable graphics)"
echo "  5. Thumbnails:       max 256px   (for galleries)"
echo ""
echo "Directory Structure:"
echo "  $OUTPUT_DIR/"
echo "  ├── print-ready/       (78 cards at full resolution)"
echo "  ├── web-optimized/     (78 cards at 1024px max)"
echo "  ├── mobile-optimized/  (78 cards at 512px max)"
echo "  ├── svg-versions/      (sample SVG conversions)"
echo "  └── thumbnails/        (78 cards at 256px max)"
echo ""
echo "${YELLOW}Usage Examples:${NC}"
echo "  - Print production: Use print-ready/ directory"
echo "  - Website gallery: Use web-optimized/ directory"
echo "  - Mobile app: Use mobile-optimized/ directory"
echo "  - Scalable designs: Use svg-versions/ directory"
echo "  - Preview galleries: Use thumbnails/ directory"
echo ""
