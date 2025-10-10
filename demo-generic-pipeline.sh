#!/bin/bash
# Demonstration of the generic pipeline system

set -e

echo "=== Generic Asset Pipeline Demo ==="
echo ""
echo "The pipeline system is now fully generic and works with any"
echo "multi-asset generation project, not just tarot decks."
echo ""

# Check if asset-generator is installed
if ! command -v asset-generator &> /dev/null; then
    echo "Installing asset-generator..."
    go install ./...
fi

echo "Step 1: Preview the pipeline structure"
echo "----------------------------------------"
echo ""
echo "The --dry-run flag shows what would be generated without"
echo "actually connecting to the generation API:"
echo ""

asset-generator pipeline \
    --file examples/generic-pipeline.yaml \
    --dry-run

echo ""
echo ""
echo "Step 2: Understanding the structure"
echo "------------------------------------"
echo ""
echo "The new generic format supports:"
echo "  • Hierarchical asset groups and subgroups"
echo "  • Custom output directories"
echo "  • Metadata that automatically enhances prompts"
echo "  • Flexible seed offsets for reproducibility"
echo ""
echo "Example structure:"
echo ""
cat << 'EOF'
assets:
  - name: Characters
    output_dir: characters
    seed_offset: 0
    metadata:
      style: "fantasy art"
    assets:
      - id: hero_01
        name: Hero
        prompt: "warrior with sword"
        metadata:
          element: "fire"
        # Final prompt: "warrior with sword, fantasy art, fire"
EOF

echo ""
echo ""
echo "Step 3: Metadata cascade"
echo "------------------------"
echo ""
echo "Metadata flows from parent to child:"
echo "  • Group metadata applies to all assets in the group"
echo "  • Subgroup metadata inherits and can override parent"
echo "  • Asset metadata has highest priority"
echo ""
echo "This makes it easy to add common styling to related assets"
echo "while still allowing per-asset customization."
echo ""

echo ""
echo "Step 4: Migration from tarot format"
echo "------------------------------------"
echo ""
echo "The old tarot-specific format has been replaced with a generic"
echo "structure that treats metadata (like suit_element) as prompt"
echo "enhancements rather than special fields."
echo ""
echo "Old format:"
echo "  minor_arcana:"
echo "    wands:"
echo "      suit_element: fire"
echo "      suit_color: red"
echo ""
echo "New format:"
echo "  assets:"
echo "    - name: Wands"
echo "      metadata:"
echo "        suit_element: fire"
echo "        suit_color: red"
echo ""
echo "The metadata values are automatically appended to prompts!"
echo ""

echo ""
echo "=== Demo complete ==="
echo ""
echo "To actually generate assets (requires SwarmUI running):"
echo "  asset-generator pipeline --file examples/generic-pipeline.yaml"
echo ""
echo "For more details, see:"
echo "  • docs/PIPELINE_MIGRATION.md - Migration guide"
echo "  • examples/generic-pipeline.yaml - Example file"
echo ""
