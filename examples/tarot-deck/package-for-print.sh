#!/bin/sh
# Package Tarot Deck for Print Production
# Creates organized archives ready for print shops
#
# Usage: ./package-for-print.sh [input_dir] [output_dir]

set -e

INPUT_DIR="${1:-./tarot-deck-processed/print-ready}"
OUTPUT_DIR="${2:-./tarot-deck-packages}"

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if input directory exists
if [ ! -d "$INPUT_DIR" ]; then
    echo "Error: Input directory '$INPUT_DIR' not found." >&2
    echo "Run post-process-deck.sh first." >&2
    exit 1
fi

mkdir -p "$OUTPUT_DIR"

echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${BLUE}       Package Tarot Deck for Print Production            ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "Input Directory: $INPUT_DIR"
echo "Output Directory: $OUTPUT_DIR"
echo ""

# Create timestamp for version tracking
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Creating Print-Ready Packages${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Package 1: Complete deck
echo "${GREEN}[1/5]${NC} Creating complete deck package..."
mkdir -p "$OUTPUT_DIR/complete-deck"
cp -r "$INPUT_DIR"/* "$OUTPUT_DIR/complete-deck/"
cd "$OUTPUT_DIR"
zip -q -r "tarot-deck-complete-${TIMESTAMP}.zip" complete-deck/
rm -rf complete-deck/
echo "  ${GREEN}✓${NC} tarot-deck-complete-${TIMESTAMP}.zip"
cd - > /dev/null
echo ""

# Package 2: Major Arcana only
echo "${GREEN}[2/5]${NC} Creating Major Arcana package..."
mkdir -p "$OUTPUT_DIR/major-arcana-only"
cp -r "$INPUT_DIR/major-arcana" "$OUTPUT_DIR/major-arcana-only/"
cd "$OUTPUT_DIR"
zip -q -r "tarot-major-arcana-${TIMESTAMP}.zip" major-arcana-only/
rm -rf major-arcana-only/
echo "  ${GREEN}✓${NC} tarot-major-arcana-${TIMESTAMP}.zip (22 cards)"
cd - > /dev/null
echo ""

# Package 3: Minor Arcana only
echo "${GREEN}[3/5]${NC} Creating Minor Arcana package..."
mkdir -p "$OUTPUT_DIR/minor-arcana-only"
cp -r "$INPUT_DIR/minor-arcana" "$OUTPUT_DIR/minor-arcana-only/"
cd "$OUTPUT_DIR"
zip -q -r "tarot-minor-arcana-${TIMESTAMP}.zip" minor-arcana-only/
rm -rf minor-arcana-only/
echo "  ${GREEN}✓${NC} tarot-minor-arcana-${TIMESTAMP}.zip (56 cards)"
cd - > /dev/null
echo ""

# Package 4: Individual suits
echo "${GREEN}[4/5]${NC} Creating individual suit packages..."
for suit in wands cups swords pentacles; do
    mkdir -p "$OUTPUT_DIR/suit-${suit}"
    cp -r "$INPUT_DIR/minor-arcana/${suit}" "$OUTPUT_DIR/suit-${suit}/"
    cd "$OUTPUT_DIR"
    zip -q -r "tarot-suit-${suit}-${TIMESTAMP}.zip" "suit-${suit}/"
    rm -rf "suit-${suit}/"
    echo "  ${GREEN}✓${NC} tarot-suit-${suit}-${TIMESTAMP}.zip (14 cards)"
    cd - > /dev/null
done
echo ""

# Package 5: Card backs (if available)
if [ -d "./tarot-deck-output/card-backs" ]; then
    echo "${GREEN}[5/5]${NC} Creating card backs package..."
    mkdir -p "$OUTPUT_DIR/card-backs"
    cp ./tarot-deck-output/card-backs/*.png "$OUTPUT_DIR/card-backs/" 2>/dev/null || true
    cd "$OUTPUT_DIR"
    zip -q -r "tarot-card-backs-${TIMESTAMP}.zip" card-backs/
    rm -rf card-backs/
    echo "  ${GREEN}✓${NC} tarot-card-backs-${TIMESTAMP}.zip"
    cd - > /dev/null
    echo ""
else
    echo "${YELLOW}[5/5]${NC} Card backs not found (run generate-card-back.sh first)"
    echo ""
fi

# Create print specification document
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Creating Print Specifications${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

cat > "$OUTPUT_DIR/PRINT_SPECIFICATIONS.txt" << 'EOF'
TAROT DECK PRINT SPECIFICATIONS
================================

CARD SPECIFICATIONS:
-------------------
Quantity: 78 cards total
  - 22 Major Arcana cards
  - 56 Minor Arcana cards (4 suits × 14 cards)

Dimensions: 768 × 1344 pixels
  - Equivalent to 2.56" × 4.48" at 300 DPI
  - Standard tarot proportion (approximately 4:7 ratio)

File Format: PNG
  - Full color (RGB)
  - No transparency
  - High resolution

RECOMMENDED PRINT SETTINGS:
--------------------------
Paper Stock: 300-350 GSM cardstock
Finish: Matte or linen finish recommended
Cutting: Standard tarot size (2.75" × 4.75")
Bleed: Add 3mm bleed on all sides during print prep

PRINT PREPARATION:
-----------------
1. Images are provided at optimal resolution
2. Add bleed area if required by printer (3mm recommended)
3. Consider card back design (multiple options provided)
4. Proof print recommended before full production
5. Color calibration may be needed for accurate reproduction

PACKAGE CONTENTS:
----------------
- tarot-deck-complete-YYYYMMDD-HHMMSS.zip
  Complete deck (78 cards)

- tarot-major-arcana-YYYYMMDD-HHMMSS.zip
  Major Arcana only (22 cards)

- tarot-minor-arcana-YYYYMMDD-HHMMSS.zip
  Minor Arcana only (56 cards)

- tarot-suit-wands-YYYYMMDD-HHMMSS.zip
  Wands suit (14 cards)

- tarot-suit-cups-YYYYMMDD-HHMMSS.zip
  Cups suit (14 cards)

- tarot-suit-swords-YYYYMMDD-HHMMSS.zip
  Swords suit (14 cards)

- tarot-suit-pentacles-YYYYMMDD-HHMMSS.zip
  Pentacles suit (14 cards)

- tarot-card-backs-YYYYMMDD-HHMMSS.zip
  Card back design variations (4 options)

FILE NAMING CONVENTION:
----------------------
Major Arcana: 00-the_fool.png through 21-the_world.png
Minor Arcana: 01-ace_of_wands.png through 14-king_of_pentacles.png

Files are named for easy identification and sorted by card order.

PRINT VENDOR CHECKLIST:
----------------------
□ Verify dimensions (768×1344px or equivalent print size)
□ Confirm bleed requirements (3mm recommended)
□ Check color profile (RGB provided, may need CMYK conversion)
□ Proof print single card before full run
□ Verify card stock weight (300-350 GSM recommended)
□ Confirm finish type (matte or linen recommended)
□ Test cutting registration
□ Verify card back alignment

ADDITIONAL NOTES:
----------------
- All cards generated with consistent style and quality
- Cards designed for traditional tarot deck layout
- Color palette optimized for mystical/esoteric aesthetic
- Borders included in design (no additional borders needed)

For questions or issues with files, refer to the README.md
in the examples/tarot-deck directory.
EOF

echo "${GREEN}✓${NC} Created PRINT_SPECIFICATIONS.txt"
echo ""

# Create manifest
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${YELLOW}Creating Package Manifest${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

cat > "$OUTPUT_DIR/MANIFEST.txt" << EOF
TAROT DECK PRINT PACKAGE MANIFEST
==================================

Package Date: $(date)
Generation Timestamp: ${TIMESTAMP}

PACKAGE CONTENTS:
----------------
EOF

# List all zip files with sizes
for file in "$OUTPUT_DIR"/*.zip; do
    if [ -f "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        basename=$(basename "$file")
        echo "  ✓ ${basename} (${size})" >> "$OUTPUT_DIR/MANIFEST.txt"
    fi
done

cat >> "$OUTPUT_DIR/MANIFEST.txt" << EOF

DOCUMENTATION:
-------------
  ✓ PRINT_SPECIFICATIONS.txt - Print vendor specifications
  ✓ MANIFEST.txt - This file
  ✓ README.md - Complete pipeline documentation

CARD COUNT VERIFICATION:
-----------------------
  Major Arcana: 22 cards (00-21)
  Minor Arcana: 56 cards
    - Wands: 14 cards
    - Cups: 14 cards
    - Swords: 14 cards
    - Pentacles: 14 cards
  Total: 78 cards

QUALITY ASSURANCE:
-----------------
  □ All 78 cards present
  □ File sizes consistent
  □ No corrupted images
  □ Correct aspect ratio
  □ Consistent quality
  □ Proper naming convention

NEXT STEPS:
----------
1. Extract complete deck package
2. Review all cards for quality
3. Send PRINT_SPECIFICATIONS.txt to print vendor
4. Request proof print before full production
5. Verify color accuracy on proof
6. Approve for full production run

For technical support or questions, refer to the
pipeline documentation in examples/tarot-deck/README.md
EOF

echo "${GREEN}✓${NC} Created MANIFEST.txt"
echo ""

# Summary
echo ""
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${GREEN}       Packaging Complete!                                ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "Created Print-Ready Packages:"
echo ""

# List all generated packages
for file in "$OUTPUT_DIR"/*.zip; do
    if [ -f "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        basename=$(basename "$file")
        echo "  ${GREEN}✓${NC} ${basename}"
        echo "    Size: ${size}"
    fi
done

echo ""
echo "Documentation:"
echo "  ${GREEN}✓${NC} PRINT_SPECIFICATIONS.txt"
echo "  ${GREEN}✓${NC} MANIFEST.txt"
echo ""
echo "Output Location: $OUTPUT_DIR"
echo ""
echo "${YELLOW}Ready for Print Production${NC}"
echo ""
echo "Next Steps:"
echo "  1. Review PRINT_SPECIFICATIONS.txt"
echo "  2. Extract and verify complete-deck package"
echo "  3. Send specifications and files to print vendor"
echo "  4. Request proof print before full production"
echo ""
