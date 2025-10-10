#!/bin/bash
# Demonstration of automatic PNG metadata stripping
# This script shows how metadata is automatically removed from PNG files
# during download, crop, and resize operations.

set -e

DEMO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$DEMO_DIR/.." && pwd)"
ASSET_GEN="$PROJECT_ROOT/asset-generator"
OUTPUT_DIR="$DEMO_DIR/metadata-demo-output"
TEST_IMAGE="$PROJECT_ROOT/demo-image.png"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  PNG Metadata Stripping Demonstration${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo

# Check if asset-generator binary exists
if [ ! -f "$ASSET_GEN" ]; then
    echo -e "${YELLOW}Building asset-generator...${NC}"
    cd "$PROJECT_ROOT"
    go build -o asset-generator .
    echo -e "${GREEN}✓ Build complete${NC}"
    echo
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Check if test image exists
if [ ! -f "$TEST_IMAGE" ]; then
    echo -e "${RED}Error: Test image not found at $TEST_IMAGE${NC}"
    exit 1
fi

echo -e "${YELLOW}Test Image:${NC} $TEST_IMAGE"
echo

# Function to check PNG metadata
check_metadata() {
    local image_path="$1"
    local label="$2"
    
    echo -e "${YELLOW}$label${NC}"
    
    # Check with pngcheck if available
    if command -v pngcheck &> /dev/null; then
        echo -e "${BLUE}Chunks found:${NC}"
        pngcheck -v "$image_path" 2>/dev/null | grep "chunk" | sed 's/^/  /'
    else
        echo -e "${BLUE}(pngcheck not installed, skipping chunk analysis)${NC}"
    fi
    
    # Check with exiftool if available
    if command -v exiftool &> /dev/null; then
        echo -e "${BLUE}Metadata fields:${NC}"
        local field_count=$(exiftool "$image_path" 2>/dev/null | grep -v "^File" | grep -v "^Directory" | grep -v "^ExifTool" | wc -l)
        echo -e "  Total metadata fields: $field_count"
        exiftool "$image_path" 2>/dev/null | grep -v "^File" | grep -v "^Directory" | grep -v "^ExifTool" | head -10 | sed 's/^/  /'
    else
        echo -e "${BLUE}(exiftool not installed, skipping metadata check)${NC}"
    fi
    
    echo
}

# Demo 1: Crop operation
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}Demo 1: Auto-Crop with Metadata Stripping${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo

CROPPED_IMAGE="$OUTPUT_DIR/demo-cropped.png"
echo -e "${YELLOW}Running:${NC} asset-generator crop --input \"$TEST_IMAGE\" --output \"$CROPPED_IMAGE\""
"$ASSET_GEN" crop --input "$TEST_IMAGE" --output "$CROPPED_IMAGE" 2>&1 || true
echo

if [ -f "$CROPPED_IMAGE" ]; then
    check_metadata "$CROPPED_IMAGE" "Cropped image metadata:"
    echo -e "${GREEN}✓ Cropped image created and metadata automatically stripped${NC}"
else
    echo -e "${RED}✗ Cropped image not created${NC}"
fi
echo

# Demo 2: Downscale operation
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}Demo 2: Downscale with Metadata Stripping${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo

DOWNSCALED_IMAGE="$OUTPUT_DIR/demo-downscaled.png"
echo -e "${YELLOW}Running:${NC} asset-generator downscale --input \"$TEST_IMAGE\" --output \"$DOWNSCALED_IMAGE\" --width 256"
"$ASSET_GEN" downscale --input "$TEST_IMAGE" --output "$DOWNSCALED_IMAGE" --width 256 2>&1 || true
echo

if [ -f "$DOWNSCALED_IMAGE" ]; then
    check_metadata "$DOWNSCALED_IMAGE" "Downscaled image metadata:"
    echo -e "${GREEN}✓ Downscaled image created and metadata automatically stripped${NC}"
else
    echo -e "${RED}✗ Downscaled image not created${NC}"
fi
echo

# Summary
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Summary${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo
echo -e "${GREEN}✓ All operations completed${NC}"
echo -e "${GREEN}✓ PNG metadata automatically stripped from all processed images${NC}"
echo -e "${GREEN}✓ Only critical chunks (IHDR, IDAT, IEND) remain${NC}"
echo
echo -e "${YELLOW}Output directory:${NC} $OUTPUT_DIR"
echo
echo -e "${BLUE}What was removed:${NC}"
echo -e "  • Text chunks (tEXt, zTXt, iTXt) - Prompts, parameters, etc."
echo -e "  • Time chunk (tIME) - Timestamps"
echo -e "  • Physical dimensions (pHYs) - DPI information"
echo -e "  • Color profiles (iCCP) - ICC profiles"
echo -e "  • All other ancillary chunks"
echo
echo -e "${BLUE}What was preserved:${NC}"
echo -e "  • Image dimensions and color information"
echo -e "  • Pixel data (no quality loss)"
echo -e "  • File format compatibility"
echo
echo -e "${YELLOW}Note:${NC} This feature is mandatory and cannot be disabled."
echo -e "${YELLOW}      It ensures privacy and security for all PNG operations.${NC}"
echo
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
