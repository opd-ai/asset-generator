# Pipeline Refactoring: Quick Reference

## Summary

The pipeline system has been refactored from tarot-specific to fully generic. All tarot-specific code has been removed and replaced with a flexible metadata-driven system.

## What You Need to Know

### 1. YAML Format Changed

**Old (Tarot-specific):**
```yaml
major_arcana:
  - number: 0
    name: The Fool
    prompt: "..."

minor_arcana:
  wands:
    suit_element: fire
    suit_color: red
    cards:
      - rank: Ace
        prompt: "..."
```

**New (Generic):**
```yaml
assets:
  - name: Major Arcana
    output_dir: major-arcana
    seed_offset: 0
    assets:
      - id: "00"
        name: The Fool
        prompt: "..."
  
  - name: Wands
    output_dir: minor-arcana/wands
    seed_offset: 100
    metadata:
      suit_element: fire
      suit_color: red
    assets:
      - id: ace
        name: Ace of Wands
        prompt: "..."
```

### 2. Metadata is Automatically Appended

Previously, fields like `suit_element` were hardcoded. Now they're metadata that gets appended to prompts:

```yaml
metadata:
  style: "fantasy art"
  quality: "detailed"
  element: "fire"

assets:
  - id: hero
    prompt: "warrior with sword"
    # Final prompt: "warrior with sword, fantasy art, detailed, fire"
```

### 3. Commands Stay the Same

```bash
# Still works exactly the same way
asset-generator pipeline --file my-assets.yaml --output-dir ./output
```

## Migration Steps

1. **Replace structure fields** with generic groups:
   - `major_arcana` → `assets[0]` with `output_dir: major-arcana`
   - `minor_arcana.wands` → subgroup with `output_dir: wands`

2. **Convert special fields to metadata**:
   - `suit_element: fire` → `metadata: { suit_element: fire }`
   - `suit_color: red` → `metadata: { suit_color: red }`

3. **Add ID fields** to each asset:
   - `number: 0` → `id: "00"`
   - `rank: Ace` → `id: ace_of_wands`

4. **Add output_dir** to each group to match desired structure

## Files to Check

- **Documentation:**
  - `docs/GENERIC_PIPELINE.md` - Complete guide
  - `docs/PIPELINE_MIGRATION.md` - Migration details
  - `docs/PIPELINE_REFACTOR_SUMMARY.md` - Technical summary

- **Examples:**
  - `examples/generic-pipeline.yaml` - New format example
  - `examples/tarot-deck-converted.yaml` - Tarot in new format

- **Demo:**
  - `demo-generic-pipeline.sh` - Interactive demonstration

## Key Benefits

✅ Works for any asset type (not just tarot)  
✅ Metadata automatically enhances prompts  
✅ Unlimited nesting with subgroups  
✅ Flexible directory structure  
✅ No special cases in code  
✅ Easier to maintain and extend  

## Quick Test

```bash
# See the new format in action
cd /home/user/go/src/github.com/opd-ai/asset-generator
./demo-generic-pipeline.sh

# Or try it directly
asset-generator pipeline --file examples/generic-pipeline.yaml --dry-run -v
```

## Need Help?

- Read `docs/GENERIC_PIPELINE.md` for full documentation
- See `examples/tarot-deck-converted.yaml` for conversion example
- Check `docs/PIPELINE_MIGRATION.md` for migration guide
