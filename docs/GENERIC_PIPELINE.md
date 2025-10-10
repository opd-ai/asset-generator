# Generic Pipeline System

## Overview

The asset-generator now features a fully generic pipeline system that works with any multi-asset generation project. The system supports hierarchical organization, metadata cascading, and flexible prompt enhancement.

## Key Features

### 1. Hierarchical Asset Groups
Organize assets into groups and subgroups with unlimited nesting:

```yaml
assets:
  - name: Characters
    output_dir: characters
    assets: [...]
    subgroups:
      - name: Heroes
        output_dir: heroes
        assets: [...]
```

### 2. Metadata Cascading
Metadata automatically flows from parent to child and gets appended to prompts:

```yaml
assets:
  - name: Fantasy Characters
    metadata:
      style: "fantasy art"
      quality: "detailed"
    assets:
      - id: hero_01
        prompt: "warrior with sword"
        metadata:
          element: "fire"
        # Final prompt: "warrior with sword, fantasy art, detailed, fire"
```

### 3. Flexible Seed Management
Control reproducibility with seed offsets:

```yaml
assets:
  - name: Group A
    seed_offset: 0     # Seeds: base+0, base+1, base+2...
  - name: Group B
    seed_offset: 100   # Seeds: base+100, base+101, base+102...
```

### 4. Custom Filenames
Specify exact filenames or let the system generate them:

```yaml
assets:
  - id: hero_warrior
    filename: knight.png        # Custom filename
  - id: hero_mage
    # Auto-generates: hero_mage.png
```

## Basic Usage

### Create a Pipeline File

```yaml
assets:
  - name: Character Sprites
    output_dir: sprites/characters
    seed_offset: 0
    metadata:
      style: "pixel art"
      size: "32x32"
    assets:
      - id: hero_idle
        name: Hero Idle Animation
        prompt: "knight standing idle, animation frame"
      
      - id: hero_walk
        name: Hero Walk Animation
        prompt: "knight walking, animation frame"
```

### Run the Pipeline

```bash
# Preview without generating
asset-generator pipeline --file my-assets.yaml --dry-run

# Generate assets
asset-generator pipeline --file my-assets.yaml --output-dir ./output

# With custom parameters
asset-generator pipeline \
    --file my-assets.yaml \
    --base-seed 42 \
    --width 512 \
    --height 512 \
    --steps 30
```

## Complete Example

```yaml
assets:
  # Main character sprites
  - name: Player Characters
    output_dir: sprites/player
    seed_offset: 0
    metadata:
      style: "16-bit pixel art"
      perspective: "top-down"
    assets:
      - id: player_warrior
        name: Warrior Class
        prompt: "armored knight with sword"
        filename: warrior.png
      
      - id: player_mage
        name: Mage Class
        prompt: "wizard with staff and robes"
        filename: mage.png
        metadata:
          element: "arcane magic"

  # Environment tiles
  - name: Environment
    output_dir: tiles
    seed_offset: 100
    metadata:
      style: "tileable texture"
    
    subgroups:
      # Grass tiles
      - name: Grass
        output_dir: grass
        seed_offset: 0
        metadata:
          type: "natural ground"
        assets:
          - id: grass_01
            name: Light Grass
            prompt: "light green grass texture"
          
          - id: grass_02
            name: Dark Grass
            prompt: "dark green grass texture"
      
      # Stone tiles
      - name: Stone
        output_dir: stone
        seed_offset: 50
        metadata:
          type: "hard surface"
        assets:
          - id: stone_01
            name: Grey Stone
            prompt: "grey cobblestone texture"
          
          - id: stone_02
            name: Brick Wall
            prompt: "red brick wall texture"
```

## Metadata System

### How Metadata Works

1. **Defined at Group Level**: Applied to all assets in the group
2. **Inherited by Subgroups**: Child groups receive parent metadata
3. **Overridable at Any Level**: Child values override parent values
4. **Appended to Prompts**: All metadata values are added to the generation prompt

### Example Flow

```yaml
assets:
  - name: Fantasy Art
    metadata:
      style: "oil painting"
      mood: "epic"
    
    subgroups:
      - name: Battle Scenes
        metadata:
          mood: "intense"        # Overrides parent's "epic"
          action: "dynamic"      # New metadata
        
        assets:
          - id: battle_01
            prompt: "warriors clashing"
            metadata:
              lighting: "dramatic"  # Asset-specific

# Final prompt for battle_01:
# "warriors clashing, oil painting, intense, dynamic, dramatic"
```

### Best Practices

1. **Use Group Metadata for Common Traits**: Style, quality, aspect ratio
2. **Use Asset Metadata for Specific Details**: Element, color, mood
3. **Keep Metadata Descriptive**: Values are added directly to prompts
4. **Avoid Redundancy**: Don't repeat what's in the base prompt

## Advanced Features

### Seed Calculation

Seeds are calculated as: `base_seed + group_seed_offset + asset_index`

By default, if `--base-seed` is not specified, or set to `0` or `-1`, a random seed is generated 
for each pipeline run. The generated seed is displayed in the output so you can reproduce results later.

```yaml
# With --base-seed 1000 (or random if not specified or 0/-1)
assets:
  - name: Group A
    seed_offset: 0
    assets:
      - id: a1  # Seed: 1000 + 0 + 0 = 1000
      - id: a2  # Seed: 1000 + 0 + 1 = 1001
  
  - name: Group B
    seed_offset: 100
    assets:
      - id: b1  # Seed: 1000 + 100 + 0 = 1100
```

### Output Directory Structure

Output directories are joined hierarchically:

```yaml
assets:
  - name: Characters
    output_dir: game/characters
    subgroups:
      - name: Heroes
        output_dir: heroes
        # Final path: <base>/game/characters/heroes/
```

### Postprocessing

All pipeline generations support postprocessing:

```bash
asset-generator pipeline \
    --file assets.yaml \
    --auto-crop \
    --downscale-width 1024 \
    --skimmed-cfg
```

## Migration from Tarot Format

The old tarot-specific format has been replaced with this generic system. Key changes:

1. **No More Special Fields**: `major_arcana`, `minor_arcana`, `suit_element`, etc. are gone
2. **Use Metadata Instead**: Former special fields become metadata that enhances prompts
3. **Flexible Structure**: No longer limited to tarot card hierarchy

See `docs/PIPELINE_MIGRATION.md` for detailed migration guide.

## Examples

- `examples/generic-pipeline.yaml` - Multi-type asset generation
- Run `./demo-generic-pipeline.sh` to see it in action

## Command Reference

```bash
asset-generator pipeline [flags]

Flags:
  --file string               Pipeline YAML file (required)
  --output-dir string         Output directory (default: "./pipeline-output")
  --base-seed int            Base seed for reproducibility (0 or -1 for random, default: -1)
  --dry-run                  Preview without generating
  --continue-on-error        Don't stop if individual generations fail
  
  # Generation parameters
  --model string             Model to use
  --width int               Image width (default: 768)
  --height int              Image height (default: 1344)
  --steps int               Inference steps (default: 40)
  --cfg-scale float         CFG scale (default: 7.5)
  --sampler string          Sampling method (default: "euler_a")
  
  # Prompt enhancement
  --style-suffix string     Suffix to append to all prompts
  --negative-prompt string  Negative prompt for all generations
  
  # Postprocessing
  --auto-crop               Auto-crop whitespace
  --downscale-width int     Downscale to width
  --downscale-percentage    Downscale by percentage
  --skimmed-cfg             Enable Skimmed CFG
```
