#!/bin/bash
# Demo script for scheduler selection feature
# This script demonstrates how to use different schedulers for various use cases

set -e  # Exit on error

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Scheduler Selection Feature Demo${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

if ! command -v asset-generator &> /dev/null; then
    echo -e "${RED}Error: asset-generator not found in PATH${NC}"
    echo "Please install it first: go install"
    exit 1
fi

# Check if SwarmUI is reachable
if ! asset-generator status > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to SwarmUI server${NC}"
    echo "Please configure the API URL:"
    echo "  asset-generator config set api-url http://your-server:7801"
    exit 1
fi

echo -e "${GREEN}✓ Prerequisites met${NC}"
echo

# Demo 1: Compare schedulers with same prompt and seed
echo "=========================================="
echo "Demo 1: Scheduler Comparison"
echo "=========================================="
echo
echo "This demo generates the same image with different schedulers"
echo "to show the quality differences."
echo
read -p "Press Enter to continue..."
echo

PROMPT="detailed portrait of a wizard, intricate details, magical atmosphere"
SEED=42
OUTPUT_DIR="./scheduler-demo"

mkdir -p "$OUTPUT_DIR"

echo -e "${BLUE}Prompt:${NC} $PROMPT"
echo -e "${BLUE}Seed:${NC} $SEED"
echo -e "${BLUE}Output:${NC} $OUTPUT_DIR"
echo

# Test each scheduler
for scheduler in simple normal karras; do
    echo -e "${YELLOW}Generating with $scheduler scheduler...${NC}"
    
    asset-generator generate image \
        --prompt "$PROMPT" \
        --seed $SEED \
        --scheduler $scheduler \
        --steps 25 \
        --width 512 \
        --height 512 \
        --save-images \
        --output-dir "$OUTPUT_DIR" \
        --filename-template "wizard-${scheduler}.png" \
        > /dev/null 2>&1
    
    echo -e "${GREEN}✓ Generated: wizard-${scheduler}.png${NC}"
done

echo
echo -e "${GREEN}✓ Demo 1 complete!${NC}"
echo "Compare the images in $OUTPUT_DIR to see the differences."
echo

# Demo 2: Fast iteration with simple scheduler
echo "=========================================="
echo "Demo 2: Fast Iteration (Simple Scheduler)"
echo "=========================================="
echo
echo "The simple scheduler is ideal for quick iteration and testing."
echo "This demo generates 3 variations quickly."
echo
read -p "Press Enter to continue..."
echo

for i in 1 2 3; do
    echo -e "${YELLOW}Generating concept $i/3...${NC}"
    
    asset-generator generate image \
        --prompt "fantasy character concept, variation $i" \
        --scheduler simple \
        --steps 15 \
        --width 512 \
        --height 512 \
        --save-images \
        --output-dir "$OUTPUT_DIR" \
        --filename-template "concept-${i}.png" \
        > /dev/null 2>&1
    
    echo -e "${GREEN}✓ Generated: concept-${i}.png${NC}"
done

echo
echo -e "${GREEN}✓ Demo 2 complete!${NC}"
echo "Fast iterations completed with simple scheduler."
echo

# Demo 3: High-quality final render with Karras
echo "=========================================="
echo "Demo 3: High-Quality Render (Karras)"
echo "=========================================="
echo
echo "The Karras scheduler is best for final, high-quality renders."
echo "This uses more steps for maximum detail."
echo
read -p "Press Enter to continue..."
echo

echo -e "${YELLOW}Generating high-quality portrait...${NC}"

asset-generator generate image \
    --prompt "photorealistic portrait, professional studio lighting, detailed eyes" \
    --scheduler karras \
    --steps 40 \
    --cfg-scale 8.0 \
    --width 768 \
    --height 768 \
    --save-images \
    --output-dir "$OUTPUT_DIR" \
    --filename-template "portrait-hq-karras.png" \
    > /dev/null 2>&1

echo -e "${GREEN}✓ Generated: portrait-hq-karras.png${NC}"
echo
echo -e "${GREEN}✓ Demo 3 complete!${NC}"
echo "High-quality render created with Karras scheduler."
echo

# Demo 4: Pipeline with scheduler
echo "=========================================="
echo "Demo 4: Pipeline with Scheduler"
echo "=========================================="
echo
echo "Schedulers can be applied to entire pipelines."
echo "This creates a mini-pipeline with the Karras scheduler."
echo
read -p "Press Enter to continue..."
echo

# Create a simple pipeline file
PIPELINE_FILE="$OUTPUT_DIR/demo-pipeline.yaml"
cat > "$PIPELINE_FILE" << 'EOF'
assets:
  - name: character_designs
    output_dir: characters
    assets:
      - id: hero
        name: Hero Character
        prompt: "brave hero character, fantasy armor, heroic pose"
      - id: villain
        name: Villain Character
        prompt: "dark villain character, mysterious cloak, menacing"
EOF

echo -e "${YELLOW}Processing pipeline with Karras scheduler...${NC}"

asset-generator pipeline \
    --file "$PIPELINE_FILE" \
    --output-dir "$OUTPUT_DIR/pipeline-output" \
    --scheduler karras \
    --steps 30 \
    --width 512 \
    --height 768 \
    > /dev/null 2>&1

echo -e "${GREEN}✓ Pipeline complete!${NC}"
echo "Characters generated in $OUTPUT_DIR/pipeline-output/characters/"
echo

# Demo 5: Configuration file usage
echo "=========================================="
echo "Demo 5: Scheduler in Config File"
echo "=========================================="
echo
echo "You can set a default scheduler in your config file."
echo
echo "Example config.yaml:"
echo -e "${BLUE}"
cat << 'EOF'
generate:
  scheduler: karras
  steps: 35
  cfg-scale: 8.0
  sampler: euler_a
EOF
echo -e "${NC}"
echo "With this config, all generations use Karras scheduler by default."
echo "Override with: --scheduler simple"
echo

# Summary
echo "=========================================="
echo "Demo Complete!"
echo "=========================================="
echo
echo "Summary of generated files in $OUTPUT_DIR:"
echo
echo "Scheduler Comparison:"
echo "  - wizard-simple.png  (fast, default)"
echo "  - wizard-normal.png  (balanced)"
echo "  - wizard-karras.png  (high quality)"
echo
echo "Fast Iteration:"
echo "  - concept-1.png"
echo "  - concept-2.png"
echo "  - concept-3.png"
echo
echo "High-Quality Render:"
echo "  - portrait-hq-karras.png"
echo
echo "Pipeline Output:"
echo "  - pipeline-output/characters/hero.png"
echo "  - pipeline-output/characters/villain.png"
echo
echo -e "${GREEN}All demos completed successfully!${NC}"
echo
echo "Quick Reference:"
echo "  --scheduler simple       Fast, default"
echo "  --scheduler normal       Balanced"
echo "  --scheduler karras       High quality"
echo "  --scheduler exponential  Artistic"
echo
echo "For more information, see: docs/SCHEDULER_FEATURE.md"
