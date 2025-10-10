# Pipeline Refactoring Summary

## What Changed

The pipeline system has been completely refactored to be generic and flexible, removing all tarot-specific code.

### Before
- Hardcoded tarot structure (`major_arcana`, `minor_arcana`, suits)
- Special field parsing (`suit_element`, `suit_color`, `rank`)
- Tarot-specific filename generation
- Rigid directory structure

### After
- Generic hierarchical asset groups
- Metadata-driven prompt enhancement
- Flexible directory structure
- Unlimited nesting with subgroups

## Key Design Decisions

### 1. Metadata as Prompt Enhancement
Instead of hardcoding what fields like `suit_element: fire` mean, we treat all metadata as contextual information that gets appended to prompts:

```yaml
# Old approach: Special parsing
suit_element: fire  # Code knows this means "add fire to prompt"

# New approach: Generic metadata
metadata:
  suit_element: fire  # Automatically appended to prompt
  suit_color: red     # Also appended
  any_field: value    # Any metadata works!
```

**Benefit**: No special cases in code, works for any use case.

### 2. Hierarchical Groups
Replaced flat tarot structure with nestable groups:

```yaml
assets:
  - name: Parent Group
    subgroups:
      - name: Child Group
        subgroups:
          - name: Grandchild Group
```

**Benefit**: Supports any organizational structure, not just tarot.

### 3. Metadata Cascading
Parent metadata flows to children with override capability:

```yaml
assets:
  - name: Parent
    metadata:
      style: "oil painting"
    subgroups:
      - name: Child
        metadata:
          mood: "dark"  # Inherits style, adds mood
        assets:
          - id: asset_01
            metadata:
              style: "watercolor"  # Overrides parent's oil painting
```

**Benefit**: DRY principle - define once, applies to many.

## Backward Compatibility

The old tarot YAML format is **not directly compatible**, but:

1. **Same output can be recreated** by using the same directory names and seed offsets
2. **Migration is straightforward** - metadata replaces special fields
3. **Example conversions provided** in `docs/PIPELINE_MIGRATION.md`

## Code Simplification

### Lines of Code
- **Removed**: ~200 lines of tarot-specific logic
- **Added**: ~150 lines of generic group/metadata handling
- **Net**: Simpler, more maintainable code

### Functions Removed
- `createOutputDirectories()` - Hardcoded tarot directories
- `getRankNumber()` - Tarot rank mapping
- Tarot-specific structs and parsing

### Functions Added
- `processGroup()` - Generic recursive group processing
- `mergeMetadata()` - Metadata inheritance
- `buildEnhancedPrompt()` - Generic prompt enhancement
- `countAssets()` - Recursive asset counting
- `previewGroups()` - Generic preview

## Testing

All existing tests pass. The refactored code:
- ✅ Compiles without errors
- ✅ Passes all unit tests
- ✅ Works with dry-run preview
- ✅ Properly merges and displays metadata

## Examples Created

1. **examples/generic-pipeline.yaml** - Shows the new format
2. **docs/GENERIC_PIPELINE.md** - Complete documentation
3. **docs/PIPELINE_MIGRATION.md** - Migration guide
4. **demo-generic-pipeline.sh** - Interactive demo

## Usage Changes

### Old Command (tarot-specific)
```bash
asset-generator pipeline --file tarot-spec.yaml
# Only worked with tarot structure
```

### New Command (generic)
```bash
asset-generator pipeline --file any-assets.yaml
# Works with any asset structure
```

The command-line interface remains the same, but the YAML format is more flexible.

## Benefits

1. **Flexibility**: Works for game sprites, UI elements, backgrounds, cards, or any assets
2. **Maintainability**: Single code path for all use cases
3. **Extensibility**: Easy to add new features (metadata types, processing options)
4. **Simplicity**: No special cases or domain-specific logic
5. **Power**: Metadata system enables complex prompt engineering

## Future Enhancements Enabled

The generic architecture makes these additions trivial:

- Per-asset generation parameters (different sizes, steps, etc.)
- Conditional metadata (apply only when certain conditions met)
- Template-based filename generation
- Asset references and dependencies
- Batch variations (generate multiple versions of each asset)

All these would be much harder with the old tarot-specific code.

## Conclusion

This refactoring transforms the pipeline from a "tarot deck generator" into a true "generic asset pipeline system" while maintaining the same level of functionality and adding more flexibility.

The code is now:
- **More generic** - Works for any use case
- **More maintainable** - Fewer special cases
- **More powerful** - Metadata cascading and nesting
- **Better documented** - Clear examples and guides
