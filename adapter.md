# Game Asset-Generator Pipeline Adapter

## Purpose

This document provides a comprehensive prompt for analyzing game codebases and automatically developing custom asset-generator pipelines that deliver all required visual assets for the game.

## The Adapter Prompt

You are an expert game asset analyst and pipeline developer. Your task is to analyze a game codebase and create a complete, production-ready asset-generator pipeline that produces all visual assets needed by the game.

### Analysis Framework

When given a game repository, systematically analyze the following:

#### 1. **Codebase Structure Analysis**
- Identify the game's technology stack (language, framework, engine)
- Locate asset directories and resource loading code
- Map file naming conventions and directory structures
- Identify asset types: sprites, tiles, UI elements, backgrounds, characters, items, effects
- Document image formats used (PNG, SVG, JPEG, etc.)
- Note any asset dimensions, size requirements, or constraints

#### 2. **Asset Requirements Discovery**
Search the codebase for:
- **Hardcoded asset references**: File paths, resource IDs, constants
- **Asset loading functions**: Image loaders, texture managers, sprite factories
- **Configuration files**: JSON, YAML, XML defining assets
- **Documentation**: README, wiki, design docs mentioning assets
- **Code comments**: TODOs, FIXMEs related to missing assets
- **Placeholder assets**: Existing temporary images or stubs

Look for patterns like:
```
# Common patterns to search for:
- LoadImage("path/to/asset.png")
- sprites["character_name"]
- texture_ids = { "background": 1, "player": 2 }
- <asset id="weapon_sword" file="sword.png"/>
- TODO: Need proper artwork for boss monster
```

#### 3. **Asset Categorization**
Organize discovered assets into logical groups:
- **Characters**: Player characters, NPCs, enemies, bosses
- **Items**: Weapons, consumables, equipment, collectibles
- **Environment**: Backgrounds, tiles, terrain, props
- **UI**: Buttons, icons, panels, menus, HUD elements
- **Effects**: Particles, animations, overlays, transitions
- **Text/Typography**: Logos, titles, decorative text

#### 4. **Visual Style Analysis**
Determine the game's artistic requirements:
- **Art style**: Pixel art, realistic, cartoon, anime, minimalist, retro
- **Color palette**: Dominant colors, mood, atmosphere
- **Perspective**: Top-down, side-view, isometric, first-person
- **Resolution**: Target dimensions for different asset types
- **Theme**: Fantasy, sci-fi, modern, historical, abstract

#### 5. **Technical Constraints**
Identify technical requirements:
- **File formats**: Required formats (PNG for transparency, JPG for photos)
- **Dimensions**: Specific sizes (32x32 tiles, 64x64 sprites, 1920x1080 backgrounds)
- **Aspect ratios**: Square, 16:9, 4:3, custom
- **Color depth**: 8-bit, 24-bit, with/without alpha channel
- **File size limits**: Memory constraints, mobile optimization
- **Naming conventions**: Prefixes, suffixes, numbering schemes

### Pipeline Generation Strategy

Once analysis is complete, generate a comprehensive YAML pipeline file with:

#### 1. **Hierarchical Structure**
```yaml
assets:
  - name: Characters
    output_dir: assets/characters
    seed_offset: 1000
    metadata:
      style: "[determined art style]"
      quality: "high detail, game asset"
    subgroups:
      - name: Player
        output_dir: player
        assets:
          - name: player-idle
            prompt: "[descriptive prompt]"
          - name: player-walk
            prompt: "[descriptive prompt]"
      - name: Enemies
        output_dir: enemies
        assets:
          # Enemy asset definitions
```

#### 2. **Metadata Cascading**
Use metadata inheritance to maintain consistency:
```yaml
metadata:
  style: "pixel art, 16-bit style, retro gaming"
  quality: "clean lines, vibrant colors"
  negative: "blurry, 3d render, photorealistic"
```

#### 3. **Seed Strategy**
Assign logical seed offsets for reproducibility:
- Major categories: 1000 apart (Characters: 0, Items: 1000, Environment: 2000)
- Subcategories: 100 apart (Heroes: 0, Enemies: 100, NPCs: 200)
- Individual assets: Sequential (enemy_1: 100, enemy_2: 101)

#### 4. **Prompt Engineering**
Craft detailed, game-specific prompts:
- Include art style keywords matching the game
- Specify perspective and viewing angle
- Add technical requirements (transparent background, centered)
- Use negative prompts to avoid unwanted styles
- Reference specific game elements for consistency

#### 5. **Output Organization**
Mirror or enhance the game's directory structure:
```
game_assets/
├── characters/
│   ├── player/
│   │   ├── idle.png
│   │   ├── walk_01.png
│   │   └── jump.png
│   └── enemies/
├── items/
│   ├── weapons/
│   └── consumables/
└── environment/
```

### Game-Specific Templates

#### Template 1: Multiplayer Space Combat Game (à la Netrek)
For games like https://github.com/opd-ai/go-netrek:

**Asset Categories:**
1. **Spaceships** (multiple types, factions)
   - Fighter ships, cruisers, destroyers, scouts
   - Different angles/rotations (0°, 45°, 90°, etc.)
   - Faction variants (Federation, Romulan, Klingon, Orion)
   
2. **Weapons & Effects**
   - Torpedo sprites, phaser beams
   - Explosion animations (frame sequence)
   - Shield impact effects
   
3. **Celestial Objects**
   - Planets (various types)
   - Stars, asteroids
   - Space stations
   
4. **UI Elements**
   - Tactical displays, radar icons
   - Status indicators, health bars
   - Minimap symbols

**Pipeline Example:**
```yaml
assets:
  # Spaceship Assets
  - name: Spaceships
    output_dir: assets/ships
    seed_offset: 0
    metadata:
      style: "retro space game sprite, top-down view, clean pixel art"
      background: "transparent background"
      quality: "sharp edges, distinct silhouette"
    
    subgroups:
      - name: Federation
        output_dir: federation
        seed_offset: 0
        metadata:
          faction: "Federation"
          color_scheme: "white and blue"
        assets:
          - id: fed_scout
            name: federation-scout
            prompt: "small triangular scout ship, top-down view, white hull with blue markings, sleek design"
            filename: "scout.png"
          
          - id: fed_cruiser
            name: federation-cruiser
            prompt: "medium cruiser starship, top-down view, white hull with blue highlights, dual nacelles"
            filename: "cruiser.png"
          
          - id: fed_battleship
            name: federation-battleship
            prompt: "large battleship, top-down view, imposing white and blue vessel, heavy armor"
            filename: "battleship.png"
      
      - name: Klingon
        output_dir: klingon
        seed_offset: 50
        metadata:
          faction: "Klingon"
          color_scheme: "dark green and red"
        assets:
          - id: kling_bird
            name: klingon-bird-of-prey
            prompt: "bird of prey ship, top-down view, dark green hull with red accents, aggressive angular design"
            filename: "bird_of_prey.png"
          
          - id: kling_cruiser
            name: klingon-battle-cruiser
            prompt: "Klingon battle cruiser, top-down view, green and red warrior ship, menacing appearance"
            filename: "battle_cruiser.png"
  
  # Weapon Effects
  - name: Weapons
    output_dir: assets/weapons
    seed_offset: 1000
    metadata:
      style: "glowing energy effect, game sprite"
      background: "transparent background"
    assets:
      - name: torpedo-active
        prompt: "glowing blue energy torpedo, small orb, bright core with trailing glow"
        filename: "torpedo.png"
      
      - name: phaser-beam
        prompt: "red energy beam line, glowing laser effect, thin straight line"
        filename: "phaser.png"
      
      - name: explosion-frame-1
        prompt: "small explosion beginning, bright orange and yellow burst, frame 1 of animation"
        filename: "explosion_01.png"
      
      - name: explosion-frame-2
        prompt: "medium explosion expanding, orange yellow red, frame 2 of animation"
        filename: "explosion_02.png"
  
  # Celestial Objects
  - name: Space Objects
    output_dir: assets/space
    seed_offset: 2000
    metadata:
      style: "stylized space game graphics, vibrant colors"
    assets:
      - name: planet-earth
        prompt: "Earth-like planet, blue and green, clouds, space view"
        filename: "planet_earth.png"
      
      - name: planet-desert
        prompt: "desert planet, orange and tan colors, arid surface"
        filename: "planet_desert.png"
      
      - name: star-base
        prompt: "space station, metallic structure, docking ports, orbital facility"
        filename: "starbase.png"
  
  # UI Elements
  - name: UI
    output_dir: assets/ui
    seed_offset: 3000
    metadata:
      style: "clean sci-fi interface icons, minimal"
      background: "transparent background"
    assets:
      - name: shield-icon
        prompt: "shield indicator icon, hexagonal energy shield, blue glow"
        filename: "icon_shield.png"
      
      - name: target-reticle
        prompt: "targeting reticle, crosshair design, red lines, HUD element"
        filename: "target_reticle.png"
```

#### Template 2: Classic RPG Game (à la GoldBox)
For games like https://github.com/opd-ai/goldbox-rpg:

**Asset Categories:**
1. **Characters**
   - Multiple races (human, elf, dwarf, halfling)
   - Classes (fighter, mage, cleric, thief, ranger)
   - Equipment variations
   - Portrait views and battle sprites
   
2. **Monsters**
   - Classic D&D creatures
   - Boss monsters
   - Animal/beast variants
   
3. **Items & Equipment**
   - Weapons (swords, axes, bows, staves)
   - Armor sets (leather, chain, plate)
   - Magic items (potions, scrolls, rings)
   
4. **Environment**
   - Dungeon tiles (walls, floors, doors)
   - Outdoor terrain (grass, forest, mountains)
   - Town/city buildings
   
5. **UI & Interface**
   - Character portraits
   - Inventory icons
   - Combat interface elements

**Pipeline Example:**
```yaml
assets:
  # Character Portraits
  - name: Character Portraits
    output_dir: assets/portraits
    seed_offset: 0
    metadata:
      style: "fantasy RPG character portrait, detailed face, medieval fantasy art"
      quality: "high detail, expressive, heroic"
      framing: "head and shoulders portrait, neutral background"
      negative: "modern clothing, sci-fi, photorealistic"
    
    subgroups:
      - name: Fighters
        output_dir: fighters
        seed_offset: 0
        metadata:
          class: "fighter"
          equipment: "armor and weapons visible"
        assets:
          - id: fighter_human_male
            name: human-male-fighter
            prompt: "male human fighter, strong features, plate armor, determined expression, sword visible"
            filename: "fighter_human_m.png"
          
          - id: fighter_dwarf_male
            name: dwarf-male-fighter
            prompt: "male dwarf fighter, beard, sturdy build, heavy armor, axe bearer, gruff look"
            filename: "fighter_dwarf_m.png"
          
          - id: fighter_elf_female
            name: elf-female-fighter
            prompt: "female elf fighter, elegant features, light armor, longsword, noble bearing"
            filename: "fighter_elf_f.png"
      
      - name: Mages
        output_dir: mages
        seed_offset: 100
        metadata:
          class: "mage"
          equipment: "robes and magical items"
        assets:
          - id: mage_human_female
            name: human-female-mage
            prompt: "female human mage, wise eyes, flowing robes, holding staff, mystical aura"
            filename: "mage_human_f.png"
          
          - id: mage_elf_male
            name: elf-male-mage
            prompt: "male elf wizard, ancient appearance, ornate robes, spellbook visible, arcane focus"
            filename: "mage_elf_m.png"
      
      - name: Clerics
        output_dir: clerics
        seed_offset: 200
        metadata:
          class: "cleric"
          equipment: "holy symbols and armor"
        assets:
          - id: cleric_human_male
            name: human-male-cleric
            prompt: "male human cleric, kind face, chainmail with holy symbol, mace, faithful expression"
            filename: "cleric_human_m.png"
  
  # Monster Sprites
  - name: Monsters
    output_dir: assets/monsters
    seed_offset: 1000
    metadata:
      style: "fantasy RPG monster art, game sprite, menacing"
      quality: "detailed creature design, clear silhouette"
      view: "front view, battle-ready pose"
    
    subgroups:
      - name: Undead
        output_dir: undead
        seed_offset: 0
        metadata:
          creature_type: "undead"
          theme: "dark, necromantic"
        assets:
          - name: skeleton-warrior
            prompt: "skeleton warrior with sword and shield, animated bones, tattered armor"
            filename: "skeleton.png"
          
          - name: zombie
            prompt: "rotting zombie, shambling pose, decayed flesh, menacing"
            filename: "zombie.png"
          
          - name: vampire
            prompt: "vampire noble, dark cloak, pale skin, fangs, aristocratic and dangerous"
            filename: "vampire.png"
      
      - name: Dragons
        output_dir: dragons
        seed_offset: 100
        metadata:
          creature_type: "dragon"
          theme: "powerful, majestic"
        assets:
          - name: red-dragon
            prompt: "large red dragon, spread wings, breathing fire, scales and horns, fearsome"
            filename: "dragon_red.png"
          
          - name: black-dragon
            prompt: "black dragon, acid breath, dark scales, evil appearance, swamp dragon"
            filename: "dragon_black.png"
      
      - name: Goblins
        output_dir: goblins
        seed_offset: 200
        assets:
          - name: goblin-warrior
            prompt: "small green goblin with crude weapon, pointy ears, mean expression"
            filename: "goblin.png"
          
          - name: hobgoblin
            prompt: "large hobgoblin soldier, better armor than goblin, disciplined warrior"
            filename: "hobgoblin.png"
  
  # Equipment & Items
  - name: Items
    output_dir: assets/items
    seed_offset: 2000
    metadata:
      style: "fantasy RPG item icon, clean design, easily recognizable"
      view: "icon view, centered, transparent background"
      background: "transparent or simple backdrop"
    
    subgroups:
      - name: Weapons
        output_dir: weapons
        seed_offset: 0
        assets:
          - name: longsword
            prompt: "silver longsword, ornate crossguard, sharp blade, fantasy weapon icon"
            filename: "weapon_longsword.png"
          
          - name: battle-axe
            prompt: "double-bladed battle axe, wooden handle, sharp metal blades, dwarf weapon"
            filename: "weapon_battleaxe.png"
          
          - name: mage-staff
            prompt: "wooden staff with glowing crystal top, magical focus, wizard weapon"
            filename: "weapon_staff.png"
          
          - name: longbow
            prompt: "elegant wooden longbow with string, elven craftsmanship, ranged weapon"
            filename: "weapon_longbow.png"
      
      - name: Armor
        output_dir: armor
        seed_offset: 100
        assets:
          - name: leather-armor
            prompt: "brown leather armor set, studded protection, light armor"
            filename: "armor_leather.png"
          
          - name: chainmail
            prompt: "metal chainmail armor, linked rings, medium protection"
            filename: "armor_chain.png"
          
          - name: plate-armor
            prompt: "full plate armor, shining metal, heavy protection, knight armor"
            filename: "armor_plate.png"
      
      - name: Magic Items
        output_dir: magic
        seed_offset: 200
        assets:
          - name: health-potion
            prompt: "red potion bottle, glowing liquid, healing elixir, cork stopper"
            filename: "potion_health.png"
          
          - name: mana-potion
            prompt: "blue potion bottle, magical energy liquid, glowing blue"
            filename: "potion_mana.png"
          
          - name: scroll
            prompt: "ancient scroll, rolled parchment, magical runes visible, spell scroll"
            filename: "scroll_magic.png"
          
          - name: magic-ring
            prompt: "golden ring with gemstone, magical aura, enchanted jewelry"
            filename: "ring_magic.png"
  
  # Dungeon Tiles
  - name: Dungeon Tiles
    output_dir: assets/tiles
    seed_offset: 3000
    metadata:
      style: "top-down dungeon tile, tileable, stone texture"
      dimensions: "square tile, seamless edges"
      view: "top-down perspective"
    assets:
      - name: stone-floor
        prompt: "gray stone floor tile, weathered texture, dungeon floor, seamless tileable"
        filename: "floor_stone.png"
      
      - name: stone-wall
        prompt: "stone dungeon wall, vertical view, moss and cracks, medieval masonry"
        filename: "wall_stone.png"
      
      - name: wooden-door
        prompt: "wooden dungeon door, iron hinges, medieval style, closed door"
        filename: "door_wood.png"
      
      - name: treasure-chest
        prompt: "wooden treasure chest, iron bands, closed, dungeon loot container"
        filename: "chest.png"
  
  # UI Elements
  - name: UI Elements
    output_dir: assets/ui
    seed_offset: 4000
    metadata:
      style: "fantasy RPG UI element, medieval theme"
      background: "transparent background"
    assets:
      - name: button-normal
        prompt: "stone fantasy UI button, neutral state, ornate border"
        filename: "button_normal.png"
      
      - name: button-hover
        prompt: "stone fantasy UI button, glowing state, highlighted ornate border"
        filename: "button_hover.png"
      
      - name: health-bar
        prompt: "red health bar, medieval frame, ornate edges, game UI element"
        filename: "healthbar.png"
      
      - name: inventory-slot
        prompt: "square inventory slot frame, stone texture, medieval border, empty"
        filename: "inventory_slot.png"
```

### Deliverables

When analyzing a game, produce:

1. **Analysis Report** (`ASSET_ANALYSIS.md`)
   - Game overview and technology stack
   - Comprehensive asset inventory
   - Directory structure mapping
   - Visual style guide
   - Technical requirements summary

2. **Pipeline YAML** (`game-assets.yaml`)
   - Complete hierarchical structure
   - All discovered assets with prompts
   - Appropriate seed offsets
   - Metadata for consistency
   - Output directory organization

3. **Generation Scripts**
   - `generate-all.sh` - Main generation script
   - `generate-[category].sh` - Category-specific scripts
   - `post-process.sh` - Post-processing pipeline
   - `verify-assets.sh` - Asset validation script

4. **Integration Guide** (`ASSET_INTEGRATION.md`)
   - How to run the pipeline
   - Build system integration
   - Asset placement instructions
   - Regeneration procedures
   - Troubleshooting tips

5. **Makefile Targets**
```makefile
.PHONY: assets assets-clean assets-preview

assets:
	asset-generator pipeline --file game-assets.yaml \
		--output-dir ./assets \
		--auto-crop \
		--downscale-width 1024

assets-preview:
	asset-generator pipeline --file game-assets.yaml --dry-run

assets-clean:
	rm -rf ./assets/generated
```

### Best Practices for Game Asset Pipelines

1. **Maintain Consistency**
   - Use metadata cascading for shared style attributes
   - Apply consistent negative prompts across categories
   - Use logical seed offset patterns

2. **Organize Hierarchically**
   - Mirror or improve upon game's directory structure
   - Group related assets together
   - Use clear, descriptive names

3. **Optimize for Game Requirements**
   - Match dimensions to game's needs
   - Use appropriate file formats
   - Consider performance constraints

4. **Enable Reproducibility**
   - Document seed values
   - Commit pipeline YAML to version control
   - Include generation parameters in documentation

5. **Plan for Iteration**
   - Start with key assets
   - Use dry-run mode extensively
   - Test individual assets before batch generation
   - Build incrementally

6. **Post-Process Appropriately**
   - Auto-crop to remove excess whitespace
   - Downscale to target dimensions
   - Strip metadata for smaller file sizes
   - Convert formats as needed

### Quality Checklist

Before delivering a pipeline, verify:

- [ ] All referenced assets in code are included in pipeline
- [ ] Asset dimensions match game requirements
- [ ] Prompts are detailed and style-consistent
- [ ] Seed offsets are logical and documented
- [ ] Output directory structure matches game layout
- [ ] Metadata cascades correctly through hierarchy
- [ ] Generation scripts are executable and tested
- [ ] Documentation includes usage examples
- [ ] Integration instructions are clear and complete
- [ ] Dry-run mode produces expected asset list

### Example Usage

```bash
# 1. Analyze game repository
cd /path/to/game-repo
./analyze-game-assets.sh > ASSET_ANALYSIS.md

# 2. Generate pipeline
# (Use templates from this document)
vi game-assets.yaml

# 3. Preview generation
asset-generator pipeline --file game-assets.yaml --dry-run

# 4. Generate assets
asset-generator pipeline --file game-assets.yaml \
  --output-dir ./assets \
  --base-seed 42 \
  --auto-crop \
  --downscale-width 1024

# 5. Integrate into game
cp -r ./assets/* ./game/resources/
make build
./game
```

## Conclusion

This adapter framework enables systematic analysis of any game codebase to produce complete, production-ready asset generation pipelines. By following this methodology, you can rapidly develop all visual assets needed for games of any genre, from retro space combat to classic RPGs and beyond.

The key is thorough analysis, appropriate categorization, detailed prompts, and logical organization that matches the game's structure and style requirements.
