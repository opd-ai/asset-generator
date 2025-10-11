# Pipeline Processing Guide

The `pipeline` command allows you to process YAML pipeline files for automated batch asset generation. This eliminates the need for external shell scripts and provides a native, cross-platform solution for complex generation workflows.

## Table of Contents

- [Overview](#overview)
- [Pipeline File Format](#pipeline-file-format)
- [Command Reference](#command-reference)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Legacy Format Support](#legacy-format-support)

## Overview

The pipeline feature is a fully generic system that works with any multi-asset generation project. It supports:

- **Hierarchical organization** with unlimited nesting of asset groups
- **Metadata cascading** that automatically enhances prompts
- **Consistent styling** across all assets
- **Organized output structure** with subdirectories
- **Reproducible results** using seed-based generation
- **Progress tracking** with detailed status updates
- **Error handling** with continue-on-error support
- **Postprocessing** including auto-crop and downscaling

### Use Cases

- **Tarot decks**: 78 cards organized by Major/Minor Arcana and suits
- **Card games**: Deck of cards with multiple suits and ranks
- **Game sprites**: Character sets, enemy variants, item collections
- **Character sheets**: Multiple poses, expressions, or equipment variants
- **Icon sets**: Consistent icon families with variations
- **UI elements**: Buttons, badges, and interface components
- **Environment assets**: Tiles, backgrounds, props
- **Any multi-asset project** with structured organization

## Pipeline File Format

Pipeline files use YAML format with a flexible, generic structure.

### Basic Structure

```yaml
assets:
  - name: Character Sprites
    output_dir: sprites/characters
    seed_offset: 0
    metadata:
      style: "pixel art"
      quality: "16-bit"
    assets:
      - id: hero_idle
        name: Hero Idle Animation
        prompt: "knight standing idle, animation frame"
      - id: hero_walk
        name: Hero Walk Animation
        prompt: "knight walking, animation frame"
```

### Hierarchical Structure with Subgroups

```yaml
major_arcana:
  - number: 0
    name: The Fool
    prompt: "detailed description of the asset"
  - number: 1
    name: The Magician
    prompt: "another detailed description"

minor_arcana:
  wands:
    suit_element: fire
    suit_color: red

    assets:    suit_color: red

      - id: hero_idle    cards:

        name: Hero Idle Animation      - rank: Ace

        prompt: "knight standing idle, animation frame"        prompt: "card description"

            - rank: Two

      - id: hero_walk        prompt: "card description"

        name: Hero Walk Animation  

        prompt: "knight walking, animation frame"  cups:

```    suit_element: water

    suit_color: blue

### Hierarchical Structure with Subgroups    cards:

      - rank: Ace

```yaml        prompt: "card description"

assets:```

  - name: Player Characters

    output_dir: sprites/player### Field Reference

    seed_offset: 0

    metadata:#### Major Arcana Cards

      style: "pixel art"

    | Field | Type | Required | Description |

    subgroups:|-------|------|----------|-------------|

      - name: Warrior Class| `number` | integer | Yes | Card number (used for seed calculation and sorting) |

        output_dir: warrior| `name` | string | Yes | Card name (used in filename and progress output) |

        seed_offset: 0| `prompt` | string | Yes | Generation prompt for this specific card |

        metadata:

          archetype: "melee fighter"#### Minor Arcana Structure

        assets:

          - id: warrior_idle| Field | Type | Required | Description |

            prompt: "armored knight standing"|-------|------|----------|-------------|

          | `suit_element` | string | No | Metadata: element associated with suit |

          - id: warrior_attack| `suit_color` | string | No | Metadata: color theme for suit |

            prompt: "armored knight swinging sword"| `cards` | array | Yes | List of cards in this suit |

      

      - name: Mage Class#### Card Objects

        output_dir: mage

        seed_offset: 50| Field | Type | Required | Description |

        metadata:|-------|------|----------|-------------|

          archetype: "spell caster"| `rank` | string | Yes | Card rank (Ace, Two, ..., King) |

        assets:| `prompt` | string | Yes | Generation prompt for this card |

          - id: mage_idle

            prompt: "wizard with staff standing"## Command Reference

          

          - id: mage_cast### Basic Usage

            prompt: "wizard casting spell, magical effects"

``````bash

asset-generator pipeline --file pipeline.yaml

### Field Reference```



#### Asset Group### All Flags



| Field | Type | Required | Description |#### Required Flags

|-------|------|----------|-------------|

| `name` | string | Yes | Group name for progress display || Flag | Description |

| `output_dir` | string | No | Subdirectory path for outputs ||------|-------------|

| `seed_offset` | integer | No | Offset added to base seed (default: 0) || `--file` | Path to pipeline YAML file |

| `metadata` | map | No | Key-value pairs appended to prompts |

| `assets` | array | No | Individual assets in this group |#### Generation Parameters

| `subgroups` | array | No | Nested asset groups |

| Flag | Default | Description |

#### Asset|------|---------|-------------|

| `--output-dir` | `./pipeline-output` | Output directory for generated assets |

| Field | Type | Required | Description || `--base-seed` | `-1` (random) | Base seed for reproducible generation (0 or -1 for random) |

|-------|------|----------|-------------|| `--model` | (none) | Model to use for all generations |

| `id` | string | Yes | Unique identifier (used for default filename) || `--steps` | `40` | Number of inference steps |

| `name` | string | No | Display name for progress output || `--width` | `768` | Image width in pixels |

| `prompt` | string | Yes | Generation prompt || `--height` | `1344` | Image height in pixels |

| `filename` | string | No | Custom filename (default: sanitized ID) || `--cfg-scale` | `7.5` | CFG scale (guidance strength) |

| `metadata` | map | No | Asset-specific metadata appended to prompt || `--sampler` | `euler_a` | Sampling method |



## Metadata System#### Prompt Enhancement



### How Metadata Works| Flag | Default | Description |

|------|---------|-------------|

Metadata provides a powerful way to enhance prompts automatically:| `--style-suffix` | (none) | Suffix appended to all prompts |

| `--negative-prompt` | (none) | Negative prompt for all generations |

1. **Defined at any level**: Group, subgroup, or individual asset

2. **Automatically cascaded**: Parent metadata flows to children#### Pipeline Control

3. **Appended to prompts**: All metadata values added to generation prompt

4. **Overridable**: Child values override parent values with same key| Flag | Default | Description |

|------|---------|-------------|

### Example| `--dry-run` | `false` | Preview pipeline without generating |

| `--continue-on-error` | `false` | Continue if individual generations fail |

```yaml

assets:#### Postprocessing

  - name: Fantasy Characters

    metadata:| Flag | Default | Description |

      style: "oil painting"|------|---------|-------------|

      mood: "epic"| `--auto-crop` | `false` | Automatically crop whitespace borders |

    | `--auto-crop-threshold` | `250` | Whitespace detection threshold (0-255) |

    subgroups:| `--auto-crop-tolerance` | `10` | Tolerance for near-white colors (0-255) |

      - name: Battle Scenes| `--auto-crop-preserve-aspect` | `false` | Preserve aspect ratio when cropping |

        metadata:| `--downscale-width` | `0` | Downscale to this width (0=disabled) |

          mood: "intense"        # Overrides parent's "epic"| `--downscale-height` | `0` | Downscale to this height (0=disabled) |

          action: "dynamic"      # New metadata| `--downscale-percentage` | `0` | Downscale by percentage (0=disabled) |

        | `--downscale-filter` | `lanczos` | Filter: lanczos, bilinear, nearest |

        assets:

          - id: battle_01## Examples

            prompt: "warriors clashing"

            metadata:### Example 1: Basic Pipeline

              lighting: "dramatic"  # Asset-specific

```bash

# Final prompt for battle_01:asset-generator pipeline --file tarot-spec.yaml --output-dir ./my-deck

# "warriors clashing, oil painting, intense, dynamic, dramatic"```

```

Output structure:

### Best Practices```

my-deck/

1. **Use group metadata for common attributes**: Style, quality, perspective├── major-arcana/

2. **Use subgroup metadata for variations**: Character class, environment type│   ├── 00-the_fool.png

3. **Use asset metadata for specifics**: Individual elements, lighting, mood│   ├── 01-the_magician.png

4. **Keep metadata concise**: Avoid redundant or conflicting descriptors│   └── ...

└── minor-arcana/

## Seed Management    ├── wands/

    │   ├── 01-ace_of_wands.png

Seeds control reproducibility and variation:    │   └── ...

    ├── cups/

```yaml    ├── swords/

assets:    └── pentacles/

  - name: Group A```

    seed_offset: 0      # Assets use: base+0, base+1, base+2...

  ### Example 2: Preview Before Generating

  - name: Group B

    seed_offset: 100    # Assets use: base+100, base+101, base+102...```bash

    asset-generator pipeline --file tarot-spec.yaml --dry-run

    subgroups:```

      - name: Subgroup B1

        seed_offset: 0  # Relative to parent: base+100+0Shows:

      - All cards that would be generated

      - name: Subgroup B2- Calculated seeds for each card

        seed_offset: 50 # Relative to parent: base+100+50- Generation parameters

```- Output structure



**Tips:**### Example 3: Custom Generation Parameters

- Leave enough space between offsets to avoid overlap

- Use consistent offset increments (100, 200, 300...)```bash

- Document your offset scheme for future referenceasset-generator pipeline --file tarot-spec.yaml \

  --base-seed 1000 \

## Command Reference  --steps 50 \

  --width 1024 \

### Basic Usage  --height 1792 \

  --cfg-scale 8.0

```bash```

asset-generator pipeline --file pipeline.yaml

```### Example 4: Style Enhancement



### All Flags```bash

asset-generator pipeline --file tarot-spec.yaml \

#### Required Flags  --style-suffix "detailed illustration, ornate border, rich colors, professional quality" \

  --negative-prompt "blurry, distorted, low quality, modern elements"

| Flag | Description |```

|------|-------------|

| `--file` | Path to pipeline YAML file |This appends the style suffix to every prompt in the pipeline.



#### Generation Parameters### Example 5: High-Resolution with Downscaling



| Flag | Default | Description |Generate at high resolution, then downscale for web use:

|------|---------|-------------|

| `--output-dir` | `./pipeline-output` | Output directory for generated assets |```bash

| `--base-seed` | `-1` (random) | Base seed for reproducible generation (0 or -1 for random) |asset-generator pipeline --file tarot-spec.yaml \

| `--model` | (none) | Model to use for all generations |  --width 1536 \

| `--steps` | `40` | Number of inference steps |  --height 2688 \

| `--width` | `768` | Image width in pixels |  --steps 50 \

| `--height` | `1344` | Image height in pixels |  --downscale-width 768 \

| `--cfg-scale` | `7.5` | CFG scale (guidance strength) |  --downscale-filter lanczos

| `--sampler` | `euler_a` | Sampling method |```

| `--scheduler` | `simple` | Scheduler (simple, normal, karras, exponential, sgm_uniform) |

### Example 6: Auto-Crop and Resize

#### Prompt Enhancement

Remove whitespace borders and resize to specific dimensions:

| Flag | Default | Description |

|------|---------|-------------|```bash

| `--style-suffix` | (none) | Suffix appended to all prompts (after metadata) |asset-generator pipeline --file tarot-spec.yaml \

| `--negative-prompt` | (none) | Negative prompt for all generations |  --auto-crop \

  --auto-crop-threshold 245 \

#### Pipeline Control  --downscale-width 768 \

  --downscale-height 1344

| Flag | Default | Description |```

|------|---------|-------------|

| `--dry-run` | `false` | Preview pipeline without generating |### Example 7: Robust Production Pipeline

| `--continue-on-error` | `false` | Continue if individual generations fail |

| `-v, --verbose` | `false` | Show detailed generation progress |```bash

asset-generator pipeline --file tarot-spec.yaml \

#### Postprocessing  --output-dir ./production-deck \

  --base-seed 42 \

| Flag | Default | Description |  --model "XE-_Pixel_Flux_-_0-1.safetensors" \

|------|---------|-------------|  --steps 50 \

| `--auto-crop` | `false` | Automatically crop whitespace borders |  --width 1536 \

| `--auto-crop-threshold` | `250` | Whitespace detection threshold (0-255) |  --height 2688 \

| `--auto-crop-tolerance` | `10` | Tolerance for near-white colors (0-255) |  --cfg-scale 7.5 \

| `--auto-crop-preserve-aspect` | `false` | Preserve aspect ratio when cropping |  --style-suffix "masterpiece, detailed, professional quality" \

| `--downscale-width` | `0` | Downscale to this width (0=disabled) |  --negative-prompt "blurry, distorted, low quality" \

| `--downscale-height` | `0` | Downscale to this height (0=disabled) |  --continue-on-error \

| `--downscale-percentage` | `0` | Downscale by percentage (0=disabled) |  --auto-crop \

| `--downscale-filter` | `lanczos` | Filter: lanczos, bilinear, nearest |  --downscale-width 768

```

## Examples

### Example 8: Verbose Output for Debugging

### Example 1: Game Character Sprites

```bash

```yamlasset-generator pipeline --file tarot-spec.yaml \

assets:  --verbose \

  - name: Player Characters  --dry-run

    output_dir: sprites/player```

    seed_offset: 0

    metadata:Shows detailed information including full prompts for each card.

      style: "pixel art, 32x32"

      perspective: "top-down view"## Best Practices

    assets:

      - id: warrior### 1. Always Preview First

        prompt: "armored knight with sword and shield"

        filename: player_warrior.pngUse `--dry-run` to verify your pipeline before generating:

      

      - id: mage```bash

        prompt: "wizard in robes with staff"asset-generator pipeline --file my-pipeline.yaml --dry-run

        filename: player_mage.png```

      

      - id: archer### 2. Use Consistent Seeds

        prompt: "ranger with bow and quiver"

        filename: player_archer.pngFor reproducible results, always use the same `--base-seed`:



  - name: Enemy Sprites```bash

    output_dir: sprites/enemiesasset-generator pipeline --file deck.yaml --base-seed 42

    seed_offset: 100```

    metadata:

      style: "pixel art, 32x32"The pipeline automatically calculates unique seeds for each asset based on the base seed.

      perspective: "top-down view"

    assets:**Note:** By default (when `--base-seed` is not specified, or set to `0` or `-1`), a random seed is 

      - id: goblingenerated for each pipeline run. The generated seed is displayed in the output so you can reproduce 

        prompt: "small green goblin with club"the same results later by explicitly specifying that seed with `--base-seed`.

      

      - id: skeleton### 3. Style Consistency

        prompt: "undead skeleton warrior"

      Use `--style-suffix` to ensure consistent styling across all assets:

      - id: dragon

        prompt: "fierce red dragon"```bash

```--style-suffix "detailed illustration, professional quality, rich colors"

```

**Generate:**

```bashThis is cleaner than adding the style to every prompt in your YAML file.

asset-generator pipeline --file game-sprites.yaml \

  --base-seed 42 \### 4. Organize Output

  --width 512 --height 512 \

  --steps 30Use descriptive output directories:

```

```bash

### Example 2: Icon Set with Variationsasset-generator pipeline --file cards.yaml --output-dir ./decks/v1-fantasy-theme

```

```yaml

assets:### 5. Error Handling

  - name: UI Icons

    output_dir: iconsFor large pipelines, use `--continue-on-error` to avoid stopping on individual failures:

    seed_offset: 0

    metadata:```bash

      style: "flat design, minimalist"asset-generator pipeline --file large-set.yaml --continue-on-error

      quality: "clean lines, simple shapes"```

    

    subgroups:You can review failures in the summary and regenerate specific assets later.

      - name: Actions

        output_dir: actions### 6. Postprocessing Pipeline

        metadata:

          category: "action icons"Combine auto-crop and downscaling for optimal results:

        assets:

          - id: icon_save1. Generate at high resolution for quality

            prompt: "floppy disk save icon"2. Auto-crop to remove borders

          3. Downscale to target size

          - id: icon_delete

            prompt: "trash bin delete icon"```bash

          asset-generator pipeline --file deck.yaml \

          - id: icon_edit  --width 2048 --height 3584 \

            prompt: "pencil edit icon"  --auto-crop \

        --downscale-width 1024 \

      - name: Status  --downscale-filter lanczos

        output_dir: status```

        metadata:

          category: "status indicators"### 7. Model Selection

        assets:

          - id: icon_successTest with different models to find the best fit:

            prompt: "green checkmark success icon"

          ```bash

          - id: icon_error# Preview with model info

            prompt: "red X error icon"asset-generator models list

          

          - id: icon_warning# Generate with specific model

            prompt: "yellow exclamation warning icon"asset-generator pipeline --file deck.yaml \

```  --model "stable-diffusion-xl-base"

```

**Generate with style enhancement:**

```bash### 8. Version Your Pipeline Files

asset-generator pipeline --file icons.yaml \

  --width 256 --height 256 \Keep your pipeline files in version control:

  --style-suffix "vector style, transparent background"

``````

project/

### Example 3: Environment Tileset├── pipelines/

│   ├── tarot-deck-v1.yaml

```yaml│   ├── tarot-deck-v2-refined.yaml

assets:│   └── character-sprites.yaml

  - name: Terrain Tiles└── outputs/

    output_dir: tiles/terrain    ├── deck-v1/

    seed_offset: 0    └── deck-v2/

    metadata:```

      style: "tileable texture, seamless"

      perspective: "top-down"## Seed Calculation

    

    subgroups:Understanding how seeds are calculated helps ensure reproducibility:

      - name: Grass

        output_dir: grass### Major Arcana

        metadata:```

          type: "natural ground"seed = base_seed + card_number

          biome: "temperate"```

        assets:

          - id: grass_lightExample with `--base-seed 42`:

            prompt: "light green grass texture"- The Fool (0) → seed 42

          - The Magician (1) → seed 43

          - id: grass_dark- The World (21) → seed 63

            prompt: "dark green grass texture"

          ### Minor Arcana

          - id: grass_flowers```

            prompt: "grass with small wildflowers"seed = base_seed + 100 + suit_offset + card_index

      ```

      - name: Stone

        output_dir: stoneSuit offsets:

        seed_offset: 50- Wands: 0

        metadata:- Cups: 20

          type: "hard surface"- Swords: 40

        assets:- Pentacles: 60

          - id: cobblestone

            prompt: "grey cobblestone pavement"Example with `--base-seed 42`:

          - Ace of Wands → seed 142

          - id: bricks- King of Wands → seed 155

            prompt: "red brick wall texture"- Ace of Cups → seed 162

          - King of Pentacles → seed 215

          - id: marble

            prompt: "white marble floor texture"## Output Structure

```

The pipeline automatically creates this directory structure:

**Generate with postprocessing:**

```bash```

asset-generator pipeline --file tileset.yaml \output-dir/

  --width 512 --height 512 \├── major-arcana/

  --auto-crop \│   ├── 00-the_fool.png

  --downscale-percentage 50│   ├── 01-the_magician.png

```│   ├── 02-the_high_priestess.png

│   └── ...

### Example 4: Dry Run (Preview)└── minor-arcana/

    ├── wands/

Preview what will be generated without actually generating:    │   ├── 01-ace_of_wands.png

    │   ├── 02-two_of_wands.png

```bash    │   └── ...

asset-generator pipeline --file my-assets.yaml --dry-run    ├── cups/

```    │   ├── 01-ace_of_cups.png

    │   └── ...

**Output:**    ├── swords/

```    │   ├── 01-ace_of_swords.png

Pipeline: my-assets.yaml    │   └── ...

═══════════════════════════════════════════════    └── pentacles/

        ├── 01-ace_of_pentacles.png

Configuration:        └── ...

  Base seed:      12345```

  Output dir:     ./pipeline-output

  Width:          768### Filename Format

  Height:         1344

  Steps:          40- Major Arcana: `{number}-{sanitized_name}.png`

  Model:          (default)- Minor Arcana: `{rank_number}-{sanitized_rank}_of_{suit}.png`



Assets to Generate:Examples:

═══════════════════════════════════════════════- `00-the_fool.png`

- `01-ace_of_wands.png`

Group: Player Characters (sprites/player)- `14-king_of_pentacles.png`

  [1] warrior → player_warrior.png

      Prompt: armored knight with sword and shield, pixel art, 32x32, top-down view## Troubleshooting

      Seed: 12345

### Pipeline Fails to Load

  [2] mage → player_mage.png

      Prompt: wizard in robes with staff, pixel art, 32x32, top-down view**Error:** `failed to load pipeline: failed to read file`

      Seed: 12346

**Solution:** Check that the file path is correct and the file exists:

Group: Enemy Sprites (sprites/enemies)

  [3] goblin → goblin.png```bash

      Prompt: small green goblin with club, pixel art, 32x32, top-down viewls -l tarot-spec.yaml

      Seed: 12445asset-generator pipeline --file ./examples/tarot-deck/tarot-spec.yaml

```

Total: 3 assets

### YAML Parse Error

DRY RUN MODE - No images will be generated

```**Error:** `failed to parse YAML`



## Best Practices**Solution:** Validate your YAML syntax:



### 1. Organize by Purpose```bash

# Install yamllint

Group related assets together:pip install yamllint

```yaml

assets:# Validate file

  - name: Charactersyamllint tarot-spec.yaml

  - name: Environment  ```

  - name: Items

  - name: UI ElementsCommon issues:

```- Incorrect indentation (use 2 spaces)

- Missing colons after keys

### 2. Use Meaningful IDs- Unquoted strings with special characters



IDs become default filenames:### Generation Failures

```yaml

# Good**Error:** `failed to generate {card_name}`

- id: hero_idle

- id: sword_fire_enchanted**Solutions:**



# Avoid1. Check API connection:

- id: img1```bash

- id: tempasset-generator models list

``````



### 3. Leverage Metadata Cascade2. Use `--continue-on-error` to see all failures:

```bash

Define common attributes once:asset-generator pipeline --file deck.yaml --continue-on-error

```yaml```

assets:

  - name: All Characters3. Enable verbose output:

    metadata:```bash

      style: "anime art style"asset-generator pipeline --file deck.yaml --verbose

      quality: "high detail"```

    subgroups:

      # All inherit style and quality### Model Not Found

      - name: Heroes

      - name: Villains**Error:** `model validation failed: model 'xyz' not found`

```

**Solution:** List available models and use the correct name:

### 4. Plan Seed Offsets

```bash

Leave room for expansion:asset-generator models list

```yamlasset-generator pipeline --file deck.yaml --model "correct-model-name"

assets:```

  - name: Group 1

    seed_offset: 0      # 0-99 range### Out of Memory

  

  - name: Group 2For large pipelines generating many high-resolution images:

    seed_offset: 100    # 100-199 range

  1. Reduce dimensions:

  - name: Group 3```bash

    seed_offset: 200    # 200-299 range--width 1024 --height 1792

``````



### 5. Use Continue-on-Error for Large Batches2. Process in batches by creating multiple smaller pipeline files



Don't let one failure stop everything:3. Use downscaling instead of generating at final resolution:

```bash```bash

asset-generator pipeline --file large-batch.yaml --continue-on-error--width 2048 --height 3584 --downscale-width 1024

``````



### 6. Test with Dry Run First### Slow Generation



Always preview before generating:**Issue:** Pipeline taking too long

```bash

asset-generator pipeline --file new-pipeline.yaml --dry-run**Solutions:**

```

1. Reduce steps:

### 7. Apply Consistent Postprocessing```bash

--steps 20

```bash```

asset-generator pipeline --file assets.yaml \

  --auto-crop \2. Use faster sampler:

  --downscale-width 1024```bash

```--sampler euler_a

```

## Troubleshooting

3. Generate lower resolution with upscaling later:

### Pipeline File Not Found```bash

--width 512 --height 896

``````

Error: failed to read pipeline file: open pipeline.yaml: no such file or directory

```### Interrupted Pipeline



**Solution:** Check the path and use absolute paths if needed:If the pipeline is interrupted (Ctrl+C), it will stop gracefully. To resume:

```bash

asset-generator pipeline --file /full/path/to/pipeline.yaml1. Check which assets were generated

```2. Remove completed cards from your pipeline file

3. Re-run with the same seed settings

### Invalid YAML Syntax

Or use `--continue-on-error` to skip already-generated files (you'll need to check manually).

```

Error: failed to parse pipeline: yaml: line 15: mapping values are not allowed in this context## Advanced Usage

```

### Custom Pipeline Structures

**Solution:** Validate your YAML syntax. Common issues:

- Incorrect indentation (use 2 spaces)While the tarot deck format is the default, you can adapt the structure for other use cases:

- Missing colons after keys

- Unquoted strings with special characters#### Game Sprites Example



Use a YAML validator or linter.```yaml

major_arcana:  # Use for main character sprites

### Generation Failures  - number: 0

    name: Hero Idle

```    prompt: "pixel art character, hero idle pose"

Error: failed to generate asset 'hero': generation timeout  - number: 1

```    name: Hero Walk

    prompt: "pixel art character, hero walking"

**Solution:** Use `--continue-on-error` to skip failures:

```bashminor_arcana:

asset-generator pipeline --file assets.yaml --continue-on-error  weapons:  # Use suits for categories

```    cards:

      - rank: Sword

### Output Directory Permissions        prompt: "pixel art sword weapon icon"

      - rank: Axe

```        prompt: "pixel art axe weapon icon"

Error: failed to create output directory: permission denied```

```

### Integration with Shell Scripts

**Solution:** Ensure write permissions or use a different directory:

```bashYou can still integrate with scripts for additional processing:

asset-generator pipeline --file assets.yaml --output-dir ~/my-assets

``````bash

#!/bin/bash

### Metadata Not Appearing in Prompts# Generate deck

asset-generator pipeline --file deck.yaml --output-dir ./raw

**Issue:** Metadata values not being added to prompts.

# Additional processing

**Solution:** Check YAML structure:for file in ./raw/major-arcana/*.png; do

```yaml  # Custom postprocessing

# Correct  convert "$file" -quality 95 "./processed/$(basename "$file")"

metadata:done

  style: "value"```



# Incorrect (typo)### Batch Multiple Pipelines

metdata:

  style: "value"```bash

```#!/bin/bash

PIPELINES=(

### Seed Reproducibility Issues  "deck-light.yaml"

  "deck-dark.yaml"

**Issue:** Different results with same seed.  "deck-vintage.yaml"

)

**Causes:**

- Different model versionsfor pipeline in "${PIPELINES[@]}"; do

- Different generation parameters  asset-generator pipeline --file "$pipeline" \

- Server-side model updates    --output-dir "./output/$(basename "$pipeline" .yaml)"

done

**Solution:** Document your exact setup:```

```yaml

# Add comments to pipeline file## See Also

# Model: stable-diffusion-xl-v1.0

# Generated: 2025-10-10- [QUICKSTART.md](../QUICKSTART.md) - Getting started guide

# Parameters: steps=40, cfg=7.5- [README.md](../README.md) - Main documentation

```- [examples/tarot-deck/](../examples/tarot-deck/) - Complete tarot deck example

- [FILENAME_TEMPLATES.md](./FILENAME_TEMPLATES.md) - Filename customization

## Legacy Format Support

The pipeline command maintains backward compatibility with the legacy tarot-specific format.

### Legacy Tarot Format

```yaml
major_arcana:
  - number: 0
    name: The Fool
    prompt: "young traveler at cliff edge"
  
  - number: 1
    name: The Magician
    prompt: "magician with tools of the trade"

minor_arcana:
  wands:
    suit_element: fire
    suit_color: red
    cards:
      - rank: Ace
        prompt: "single wand with flames"
      - rank: Two
        prompt: "two crossed wands"
  
  cups:
    suit_element: water
    suit_color: blue
    cards:
      - rank: Ace
        prompt: "ornate chalice with water"
```

### Migration to Generic Format

For new projects, use the generic format. The legacy tarot-specific format (`major_arcana`, `minor_arcana`) is still supported for backward compatibility but is deprecated.

**Key differences in generic format:**
- Uses `assets` instead of `major_arcana` and `minor_arcana`
- Supports unlimited nesting with `subgroups`
- Metadata is explicitly defined and cascaded
- More flexible file naming and organization

---

## Quick Reference

### Basic Usage

```bash
asset-generator pipeline --file pipeline.yaml
```

### Common Commands

```bash
# Preview pipeline (dry run)
asset-generator pipeline --file deck.yaml --dry-run

# Generate with custom output directory
asset-generator pipeline --file deck.yaml --output-dir ./my-output

# High quality with postprocessing
asset-generator pipeline --file deck.yaml \
  --steps 50 \
  --scheduler karras \
  --auto-crop \
  --downscale-width 1024

# Continue on errors
asset-generator pipeline --file deck.yaml --continue-on-error

# Add style to all prompts
asset-generator pipeline --file deck.yaml \
  --style-suffix "detailed, professional quality, rich colors"
```

### Key Flags

#### Required
- `--file` - Pipeline YAML file path

#### Generation
- `--output-dir` - Output directory (default: `./pipeline-output`)
- `--base-seed` - Base seed for reproducibility (default: `-1` for random)
- `--steps` - Inference steps (default: `40`)
- `--width` - Image width (default: `768`)
- `--height` - Image height (default: `1344`)
- `--cfg-scale` - CFG scale/guidance (default: `7.5`)
- `--model` - Model to use for generation
- `--scheduler` - Scheduler (simple, normal, karras, exponential, sgm_uniform)

#### Prompt Enhancement
- `--style-suffix` - Append to all prompts
- `--negative-prompt` - Negative prompt for all

#### Control
- `--dry-run` - Preview without generating
- `--continue-on-error` - Don't stop on failures
- `-v, --verbose` - Show detailed progress

#### Postprocessing
- `--auto-crop` - Remove whitespace borders
- `--downscale-width` - Downscale to width
- `--downscale-height` - Downscale to height
- `--downscale-filter` - Filter: lanczos, bilinear, nearest

### Troubleshooting

```bash
# Preview first
asset-generator pipeline --file deck.yaml --dry-run

# Check configuration
asset-generator config view

# Test connection
asset-generator models list

# Enable verbose output
asset-generator pipeline --file deck.yaml --verbose

# Continue on errors for large pipelines
asset-generator pipeline --file deck.yaml --continue-on-error
```

---

## See Also

- [Generation Features](GENERATION_FEATURES.md) - Scheduler, Skimmed CFG
- [Postprocessing](POSTPROCESSING.md) - Auto-crop, downscaling, metadata stripping
- [Seed Behavior](SEED_BEHAVIOR.md) - Random vs explicit seeds
- [Quick Start Guide](QUICKSTART.md) - Getting started
- [Examples](../examples/) - Sample pipeline files
