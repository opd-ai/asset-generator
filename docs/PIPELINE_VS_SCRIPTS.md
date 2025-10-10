# Pipeline Command vs Shell Scripts: Comparison

## Overview

The new `pipeline` command provides a native, cross-platform solution for batch asset generation that was previously only possible through shell scripts. This document compares the two approaches.

## Quick Comparison

| Feature | Shell Script | Pipeline Command |
|---------|-------------|------------------|
| **Dependencies** | bash, yq, grep, sed | None (built-in) |
| **Cross-platform** | ❌ Linux/macOS only | ✅ Windows/Linux/macOS |
| **YAML parsing** | External yq tool | ✅ Native Go parser |
| **Error handling** | Manual exit codes | ✅ Structured errors |
| **Progress tracking** | Custom echo statements | ✅ Built-in progress |
| **Dry-run preview** | Must implement manually | ✅ Built-in --dry-run |
| **Postprocessing** | Separate script calls | ✅ Integrated flags |
| **Signal handling** | Manual trap setup | ✅ Automatic handling |
| **Installation** | Multiple tools | ✅ Single binary |

## Old Approach: Shell Script

### generate-tarot-deck.sh (excerpt)

```bash
#!/bin/sh
set -e

# Dependencies check
command -v asset-generator >/dev/null 2>&1 || exit 1
command -v yq >/dev/null 2>&1 || exit 1

# Parse YAML manually
major_count=$(yq eval '.major_arcana | length' "$SPEC_FILE")
i=0

while [ "$i" -lt "$major_count" ]; do
    number=$(yq eval ".major_arcana[$i].number" "$SPEC_FILE")
    name=$(yq eval ".major_arcana[$i].name" "$SPEC_FILE")
    prompt=$(yq eval ".major_arcana[$i].prompt" "$SPEC_FILE")
    
    # Sanitize name for filename
    filename=$(echo "$name" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
    
    # Calculate seed
    card_seed=$((BASE_SEED + number))
    
    # Generate
    asset-generator generate image \
        --prompt "$prompt" \
        --seed "$card_seed" \
        --save-images \
        --output-dir "$(dirname "$output_path")" \
        --filename-template "$(basename "$output_path")" \
        > /dev/null 2>&1
    
    i=$((i + 1))
done
```

### Issues with Shell Script Approach

1. **Platform Dependency**: Requires POSIX shell (bash/sh)
   - Not native on Windows
   - Different behavior on macOS vs Linux

2. **External Dependencies**: 
   - Must install `yq` (correct version!)
   - Different yq versions (python-yq vs mikefarah's yq)
   - Installation complexity varies by platform

3. **Manual Implementation**:
   - Custom loop logic for iteration
   - Manual string sanitization
   - Custom error handling with exit codes
   - Progress output requires echo statements

4. **Limited Error Recovery**:
   - Script stops on first error (unless you add complex logic)
   - No built-in continue-on-error
   - Manual tracking of what succeeded/failed

5. **No Preview Mode**:
   - Must implement custom dry-run logic
   - Can't easily preview without generating

6. **Maintenance Burden**:
   - Shell script quirks and gotchas
   - Quoting and escaping complexity
   - Different behavior across shells

## New Approach: Pipeline Command

### Usage

```bash
asset-generator pipeline --file tarot-spec.yaml
```

### Advantages

#### 1. Zero Dependencies
```bash
# Only need asset-generator binary
asset-generator pipeline --file deck.yaml
```

No need to install or manage:
- yq
- Special shell versions
- Additional parsing tools

#### 2. Cross-Platform Native
```bash
# Works identically on all platforms
# Windows
asset-generator.exe pipeline --file deck.yaml

# Linux/macOS
asset-generator pipeline --file deck.yaml
```

#### 3. Built-in Preview
```bash
# Instant preview without generating
asset-generator pipeline --file deck.yaml --dry-run
```

Output:
```
Pipeline Preview:

Major Arcana:
  00 - The Fool                  (seed: 42)
  01 - The Magician              (seed: 43)
  ...

Generation Parameters:
  Dimensions: 768x1344
  Steps: 40
  CFG Scale: 7.5
```

#### 4. Robust Error Handling
```bash
# Continue generating even if some fail
asset-generator pipeline --file deck.yaml --continue-on-error
```

Output includes summary:
```
Pipeline Complete!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total cards generated: 76/78
Failed: 2
Output location: ./pipeline-output
```

#### 5. Integrated Postprocessing
```bash
# All in one command
asset-generator pipeline --file deck.yaml \
  --auto-crop \
  --downscale-width 1024 \
  --downscale-filter lanczos
```

No need for separate post-processing scripts.

#### 6. Style Management
```bash
# Apply style to all prompts without modifying YAML
asset-generator pipeline --file deck.yaml \
  --style-suffix "detailed illustration, ornate border" \
  --negative-prompt "blurry, low quality"
```

#### 7. Progress Tracking
Built-in progress with detailed status:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Processing Major Arcana (22 cards)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[1/78] Generating: 00 - The Fool
  ✓ Saved to: ./output/major-arcana/00-the_fool.png

[2/78] Generating: 01 - The Magician
  ✓ Saved to: ./output/major-arcana/01-the_magician.png
```

#### 8. Signal Handling
Automatic graceful shutdown:

```bash
# Press Ctrl+C during generation
^C
Received interrupt signal, cancelling...
pipeline cancelled: context canceled
```

Clean exit, no orphaned processes.

## Migration Guide

### From Shell Script to Pipeline Command

**Old:**
```bash
./generate-tarot-deck.sh ./output 42
```

**New:**
```bash
asset-generator pipeline \
  --file tarot-spec.yaml \
  --output-dir ./output \
  --base-seed 42
```

### Common Parameters

| Shell Script Variable | Pipeline Flag |
|----------------------|---------------|
| `OUTPUT_DIR` | `--output-dir` |
| `BASE_SEED` | `--base-seed` |
| `WIDTH` | `--width` |
| `HEIGHT` | `--height` |
| `STEPS` | `--steps` |
| `CFG_SCALE` | `--cfg-scale` |
| `STYLE_SUFFIX` | `--style-suffix` |
| `NEGATIVE_PROMPT` | `--negative-prompt` |

### Example Migration

**Before (generate-tarot-deck.sh):**
```bash
#!/bin/sh
set -e

OUTPUT_DIR="${1:-./tarot-deck-output}"
BASE_SEED="${2:-42}"
WIDTH=768
HEIGHT=1344
STEPS=40
CFG_SCALE=7.5
STYLE_SUFFIX="detailed illustration, ornate border"
NEGATIVE_PROMPT="blurry, distorted, low quality"

# ... 200+ lines of shell script logic ...
```

**After (single command):**
```bash
asset-generator pipeline \
  --file tarot-spec.yaml \
  --output-dir ./tarot-deck-output \
  --base-seed 42 \
  --width 768 \
  --height 1344 \
  --steps 40 \
  --cfg-scale 7.5 \
  --style-suffix "detailed illustration, ornate border" \
  --negative-prompt "blurry, distorted, low quality"
```

## Performance Comparison

| Metric | Shell Script | Pipeline Command |
|--------|-------------|------------------|
| **Startup Time** | ~200ms (fork yq, parse) | ~50ms (native) |
| **Memory Usage** | Multiple processes | Single process |
| **Error Recovery** | Stop or complex logic | Built-in continue |
| **Parallel Potential** | Complex | Future-ready |

## When to Use Each

### Use Pipeline Command When:
- ✅ You want cross-platform support
- ✅ You need a single tool without dependencies
- ✅ You want built-in error handling and progress
- ✅ You need dry-run preview
- ✅ You want integrated postprocessing
- ✅ You prefer structured configuration

### Keep Shell Scripts When:
- You need custom processing between generations
- You're integrating with complex shell-based workflows
- You need bash-specific features (arrays, functions, etc.)
- You're already deeply invested in shell infrastructure

## Recommendation

**For most users**: Use the `pipeline` command. It's simpler, more reliable, and cross-platform.

**For advanced workflows**: You can still use shell scripts that *call* the pipeline command, getting the best of both worlds:

```bash
#!/bin/bash
# Advanced workflow using pipeline command

# Generate deck
asset-generator pipeline --file deck-v1.yaml --output-dir ./v1

# Custom processing
for file in ./v1/major-arcana/*.png; do
  # Your custom logic here
  convert "$file" -custom-filter "./processed/$(basename "$file")"
done

# Generate another variant
asset-generator pipeline --file deck-v2.yaml --output-dir ./v2
```

## Future Enhancements

The pipeline command architecture supports future features:

- **Parallel generation**: Generate multiple cards concurrently
- **Resume capability**: Resume interrupted pipelines
- **Progress persistence**: Save/restore pipeline state
- **Custom hooks**: Run commands before/after generation
- **Template system**: Reusable pipeline templates
- **Validation**: Pre-flight checks before generation

These would be very difficult to implement reliably in shell scripts.

## Conclusion

The `pipeline` command brings professional-grade batch generation to asset-generator without sacrificing simplicity. It eliminates external dependencies, works across all platforms, and provides robust error handling and progress tracking out of the box.

**Bottom line**: What took 200+ lines of shell script now takes a single command.
