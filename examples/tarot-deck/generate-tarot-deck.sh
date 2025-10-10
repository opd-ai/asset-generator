#!/bin/sh
# Complete Tarot Deck Generator
# Generates all 78 tarot cards with consistent styling and organized structure
#
# Usage: ./generate-tarot-deck.sh [output_dir] [base_seed]
#
# Arguments:
#   output_dir  - Base directory for generated cards (default: ./tarot-deck-output)
#   base_seed   - Starting seed for reproducible generation (default: 42)

set -e

# Configuration
OUTPUT_DIR="${1:-./tarot-deck-output}"
BASE_SEED="${2:-42}"
SPEC_FILE="tarot-spec.yaml"

# Card dimensions (high resolution for printing)
WIDTH=768
HEIGHT=1344

# Override dimensions for quicker generation:
WIDTH=512
HEIGHT=768

# Generation parameters
STEPS=40
CFG_SCALE=5.5

# Styling keywords for consistency
STYLE_SUFFIX="readable, detailed realistic illustration, ornate decorative border, mystical symbols, rich colors, traditional tarot card design, professional quality"
NEGATIVE_PROMPT="blurry, unreadable, distorted, low quality, text, watermark, signatures, modern elements, photograph"

# Color codes for terminal output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verify dependencies
command -v asset-generator >/dev/null 2>&1 || {
    echo "Error: asset-generator CLI not found. Please install it first." >&2
    exit 1
}

# Check for yq and verify it's the correct version (mikefarah's yq, not python-yq)
if ! command -v yq >/dev/null 2>&1; then
    echo "Error: yq is required but not installed." >&2
    echo "This pipeline requires mikefarah's yq (Go version), not python-yq" >&2
    echo "" >&2
    echo "Install with:" >&2
    echo "  wget -qO /tmp/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64" >&2
    echo "  sudo mv /tmp/yq /usr/local/bin/yq" >&2
    echo "  sudo chmod +x /usr/local/bin/yq" >&2
    exit 1
fi

# Verify it's the Go version (mikefarah's yq)
if ! yq --version 2>&1 | grep -q "mikefarah"; then
    echo "Error: Wrong version of yq detected." >&2
    echo "This pipeline requires mikefarah's yq (Go version), not python-yq" >&2
    echo "" >&2
    echo "Current yq version:" >&2
    yq --version >&2
    echo "" >&2
    echo "Please install the correct version:" >&2
    echo "  wget -qO /tmp/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64" >&2
    echo "  sudo mv /tmp/yq /usr/local/bin/yq" >&2
    echo "  sudo chmod +x /usr/local/bin/yq" >&2
    exit 1
fi

# Check if spec file exists
if [ ! -f "$SPEC_FILE" ]; then
    echo "Error: Specification file '$SPEC_FILE' not found." >&2
    exit 1
fi

# Create output directory structure
mkdir -p "$OUTPUT_DIR"/{major-arcana,minor-arcana/{wands,cups,swords,pentacles}}

echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${BLUE}         Complete Tarot Deck Generation Pipeline          ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "Output Directory: $OUTPUT_DIR"
echo "Base Seed: $BASE_SEED"
echo "Card Dimensions: ${WIDTH}x${HEIGHT}"
echo "Generation Steps: $STEPS"
echo ""

# Function to generate a single card
generate_card() {
    prompt="$1"
    output_path="$2"
    card_name="$3"
    seed="$4"
    
    echo "${GREEN}Generating:${NC} $card_name"
    
    # Combine prompt with style suffix
    full_prompt="${prompt}, ${STYLE_SUFFIX}"
    
    asset-generator generate image \
        --prompt "$full_prompt" \
        --negative-prompt "$NEGATIVE_PROMPT" \
        --width "$WIDTH" \
        --height "$HEIGHT" \
        --steps "$STEPS" \
        --cfg-scale "$CFG_SCALE" \
        --seed "$seed" \
        --model "XE-_Pixel_Flux_-_0-1.safetensors" \
        --save-images \
        --output-dir "$(dirname "$output_path")" \
        --filename-template "$(basename "$output_path")" \
        > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        echo "  ${GREEN}✓${NC} Saved to: $output_path"
    else
        echo "  ${YELLOW}⚠${NC} Failed to generate $card_name"
    fi
}

# Generate Major Arcana (22 cards)
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${BLUE}Major Arcana (22 cards)${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

major_count=$(yq eval '.major_arcana | length' "$SPEC_FILE")
i=0

while [ "$i" -lt "$major_count" ]; do
    number=$(yq eval ".major_arcana[$i].number" "$SPEC_FILE")
    name=$(yq eval ".major_arcana[$i].name" "$SPEC_FILE")
    prompt=$(yq eval ".major_arcana[$i].prompt" "$SPEC_FILE")
    
    # Format number with leading zero
    padded_number=$(printf "%02d" "$number")
    
    # Sanitize name for filename
    filename=$(echo "$name" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
    
    # Calculate seed (increment from base)
    card_seed=$((BASE_SEED + number))
    
    output_file="$OUTPUT_DIR/major-arcana/${padded_number}-${filename}.png"
    
    generate_card "$prompt" "$output_file" "$padded_number - $name" "$card_seed"
    
    i=$((i + 1))
done

echo ""
echo "${GREEN}✓ Major Arcana complete!${NC}"
echo ""

# Generate Minor Arcana (56 cards)
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${BLUE}Minor Arcana (56 cards)${NC}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Seed offset for minor arcana
minor_seed_offset=100

# Process each suit
for suit in wands cups swords pentacles; do
    echo "${YELLOW}Suit: ${suit}${NC}"
    echo ""
    
    card_count=$(yq eval ".minor_arcana.${suit}.cards | length" "$SPEC_FILE")
    i=0
    
    while [ "$i" -lt "$card_count" ]; do
        rank=$(yq eval ".minor_arcana.${suit}.cards[$i].rank" "$SPEC_FILE")
        prompt=$(yq eval ".minor_arcana.${suit}.cards[$i].prompt" "$SPEC_FILE")
        
        # Sanitize rank for filename
        filename=$(echo "$rank" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
        
        # Calculate seed
        card_seed=$((BASE_SEED + minor_seed_offset + i))
        
        # Output filename with rank number prefix for sorting
        case "$rank" in
            Ace) rank_num="01" ;;
            Two) rank_num="02" ;;
            Three) rank_num="03" ;;
            Four) rank_num="04" ;;
            Five) rank_num="05" ;;
            Six) rank_num="06" ;;
            Seven) rank_num="07" ;;
            Eight) rank_num="08" ;;
            Nine) rank_num="09" ;;
            Ten) rank_num="10" ;;
            Page) rank_num="11" ;;
            Knight) rank_num="12" ;;
            Queen) rank_num="13" ;;
            King) rank_num="14" ;;
        esac
        
        output_file="$OUTPUT_DIR/minor-arcana/${suit}/${rank_num}-${filename}_of_${suit}.png"
        
        generate_card "$prompt" "$output_file" "$rank of ${suit}" "$card_seed"
        
        i=$((i + 1))
    done
    
    # Increment seed offset for next suit
    minor_seed_offset=$((minor_seed_offset + 20))
    
    echo ""
    echo "${GREEN}✓ Suit of ${suit} complete!${NC}"
    echo ""
done

# Generate summary
echo ""
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo "${GREEN}       Complete Tarot Deck Generation Finished!           ${NC}"
echo "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo "Total Cards Generated: 78"
echo "  - Major Arcana: 22 cards"
echo "  - Minor Arcana: 56 cards (4 suits × 14 cards)"
echo ""
echo "Output Location: $OUTPUT_DIR"
echo ""
echo "Directory Structure:"
echo "  $OUTPUT_DIR/"
echo "  ├── major-arcana/        (22 cards)"
echo "  └── minor-arcana/"
echo "      ├── wands/           (14 cards)"
echo "      ├── cups/            (14 cards)"
echo "      ├── swords/          (14 cards)"
echo "      └── pentacles/       (14 cards)"
echo ""
echo "${YELLOW}Next Steps:${NC}"
echo "  1. Review generated cards in $OUTPUT_DIR"
echo "  2. Post-process with: ./post-process-deck.sh"
echo "  3. Create card backs with: ./generate-card-back.sh"
echo "  4. Package for printing with: ./package-for-print.sh"
echo ""
