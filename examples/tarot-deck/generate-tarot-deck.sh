#!/bin/sh
# Complete Tarot Deck Generator - Pipeline Wrapper
# Generates all 78 tarot cards using the native pipeline command
#
# This script is a convenience wrapper around the asset-generator pipeline command.
# It maintains backward compatibility with the original shell script interface
# while leveraging the new native pipeline functionality.
#
# Usage: ./generate-tarot-deck.sh [output_dir] [base_seed] [additional_flags...]
#
# Arguments:
#   output_dir  - Base directory for generated cards (default: ./tarot-deck-output)
#   base_seed   - Starting seed for reproducible generation (default: 42)
#   additional_flags - Any additional flags to pass to pipeline command
#
# Examples:
#   ./generate-tarot-deck.sh
#   ./generate-tarot-deck.sh ./my-deck 42
#   ./generate-tarot-deck.sh ./my-deck 42 --dry-run
#   ./generate-tarot-deck.sh ./my-deck 42 --continue-on-error

set -e

# Configuration
OUTPUT_DIR="${1:-./tarot-deck-output}"
BASE_SEED="${2:-42}"
SPEC_FILE="tarot-spec.yaml"

# Shift past the first two arguments so remaining args can be passed to pipeline
shift 2 2>/dev/null || shift $# 2>/dev/null || true

# Card dimensions (high resolution for printing)
WIDTH=768
HEIGHT=1344

# Generation parameters
STEPS=40
CFG_SCALE=7.5

# Styling keywords for consistency
STYLE_SUFFIX="detailed illustration, ornate decorative border, mystical symbols, rich colors, traditional tarot card design, professional quality"
NEGATIVE_PROMPT="blurry, unreadable, distorted, low quality"

# Color codes for terminal output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verify dependencies
command -v asset-generator >/dev/null 2>&1 || {
    echo "Error: asset-generator CLI not found. Please install it first." >&2
    echo "" >&2
    echo "Install from: https://github.com/opd-ai/asset-generator" >&2
    exit 1
}

# Check if spec file exists
if [ ! -f "$SPEC_FILE" ]; then
    echo "Error: Specification file '$SPEC_FILE' not found." >&2
    echo "Expected location: $(pwd)/$SPEC_FILE" >&2
    exit 1
fi

echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${BLUE}         Complete Tarot Deck Generation Pipeline          ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "Using native pipeline command (no external dependencies needed)"
echo ""

# Execute the pipeline command with all configured parameters
asset-generator pipeline \
    --file "$SPEC_FILE" \
    --output-dir "$OUTPUT_DIR" \
    --base-seed "$BASE_SEED" \
    --width "$WIDTH" \
    --height "$HEIGHT" \
    --steps "$STEPS" \
    --cfg-scale "$CFG_SCALE" \
    --style-suffix "$STYLE_SUFFIX" \
    --negative-prompt "$NEGATIVE_PROMPT" \
    "$@"

# Check exit status
if [ $? -eq 0 ]; then
    echo ""
    echo "${GREEN}✓ Pipeline completed successfully!${NC}"
    echo ""
    echo "${YELLOW}Next Steps:${NC}"
    echo "  1. Review generated cards in $OUTPUT_DIR"
    echo "  2. Post-process with: ./post-process-deck.sh"
    echo "  3. Create card backs with: ./generate-card-back.sh"
    echo "  4. Package for printing with: ./package-for-print.sh"
    echo ""
else
    echo ""
    echo "${YELLOW}⚠ Pipeline encountered errors. Check output above for details.${NC}"
    echo ""
    exit 1
fi

