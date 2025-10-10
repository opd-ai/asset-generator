# Pipeline Migration Guide

## Overview

The pipeline system has been refactored to be fully generic and work with any multi-asset generation project, not just tarot decks. The new format uses a flexible hierarchical structure with metadata support.

## Key Changes

### 1. Generic Structure
Instead of tarot-specific fields (`major_arcana`, `minor_arcana`, `wands`, etc.), the new format uses generic `assets` groups.

### 2. Metadata Enhancement
Metadata fields (like `suit_element`, `suit_color`) are now automatically appended to prompts, eliminating the need for specialized parsing logic.

### 3. Flexible Hierarchy
- Asset groups can contain subgroups for unlimited nesting
- Each group can have its own output directory and seed offset
- Metadata cascades from parent to child groups

## Migration Example

### Old Tarot Format
```yaml
major_arcana:
  - number: 0
    name: The Fool
    prompt: "young traveler at cliff edge"

minor_arcana:
  wands:
    suit_element: fire
    suit_color: red
    cards:
      - rank: Ace
        prompt: "single wand with flames"
```

### New Generic Format
```yaml
assets:
  - name: Major Arcana
    output_dir: major-arcana
    seed_offset: 0
    metadata:
      style: "tarot card art"
      quality: "ornate border"
    assets:
      - id: "00"
        name: The Fool
        prompt: "young traveler at cliff edge"
        filename: "00-the_fool.png"
  
  - name: Minor Arcana
    output_dir: minor-arcana
    seed_offset: 100
    metadata:
      style: "tarot card art"
    subgroups:
      - name: Wands
        output_dir: wands
        seed_offset: 0
        metadata:
          suit_element: fire
          suit_color: red
        assets:
          - id: ace_of_wands
            name: Ace of Wands
            prompt: "single wand with flames"
            filename: "01-ace_of_wands.png"
```

## How Metadata Works

Metadata values are automatically appended to prompts:

```yaml
# Group metadata applies to all assets in the group
metadata:
  style: "fantasy art"
  quality: "high detail"

# Asset with its own metadata
assets:
  - id: hero
    prompt: "warrior with sword"
    metadata:
      element: "fire"

# Final prompt becomes:
# "warrior with sword, fantasy art, high detail, fire"
```

### Metadata Inheritance

- Child metadata overrides parent metadata
- Subgroups inherit parent group metadata
- Asset metadata overrides group metadata

Example:
```yaml
assets:
  - name: Characters
    metadata:
      style: "anime"      # Applies to all
      quality: "detailed" # Applies to all
    subgroups:
      - name: Heroes
        metadata:
          style: "realistic"  # Overrides parent's "anime"
        assets:
          - id: hero_01
            prompt: "knight"
            metadata:
              quality: "ultra detailed"  # Overrides "detailed"
            # Final metadata: style="realistic", quality="ultra detailed"
```

## Benefits

1. **No Special Cases**: Works for any asset type (sprites, backgrounds, UI, cards, etc.)
2. **Flexible Organization**: Unlimited nesting and custom directory structure
3. **Metadata as Context**: Automatically enhance prompts with contextual information
4. **Backward Compatible Structure**: Can recreate the same output structure as before
5. **Easier Maintenance**: Single code path for all pipeline types

## Converting Existing Files

You can recreate the exact same output structure by:

1. Using the same output directory names (`major-arcana`, `minor-arcana/wands`, etc.)
2. Using the same seed offsets
3. Using the same filename patterns
4. Adding metadata fields (like `suit_element`) that were previously hardcoded

The generated images will be identical when using the same seeds and prompts.

## Examples

See:
- `examples/generic-pipeline.yaml` - Generic multi-asset pipeline
- `examples/tarot-deck-converted.yaml` - Tarot deck in new format
