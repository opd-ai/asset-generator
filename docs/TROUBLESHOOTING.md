[🏠 Docs Home](README.md) | [📚 Quick Start](QUICKSTART.md) | [🔧 Commands](COMMANDS.md) | [👤 User Guide](USER_GUIDE.md)

---

# Troubleshooting Guide

> **Complete guide to resolving common issues with Asset Generator CLI**

This guide covers the most common problems and their solutions, organized by category for quick resolution.

## Table of Contents

- [Connection Issues](#connection-issues)
- [Configuration Problems](#configuration-problems)
- [Generation Failures](#generation-failures)
- [Pipeline Issues](#pipeline-issues)
- [LoRA Problems](#lora-problems)
- [File and Permission Issues](#file-and-permission-issues)
- [Performance Problems](#performance-problems)
- [Output and Format Issues](#output-and-format-issues)
- [Feature-Specific Issues](#feature-specific-issues)
- [Error Message Reference](#error-message-reference)

---

## Connection Issues {#connection-issues}

### Problem: "Connection refused" or "Cannot connect to SwarmUI"

**Symptoms:**
- Error: `failed to connect to asset generation service`
- Error: `connection refused`
- Error: `no such host`

**Solutions:**

1. **Check if SwarmUI is running:**
   ```bash
   curl http://localhost:7801/API/ListModels
   ```

2. **Verify API URL configuration:**
   ```bash
   asset-generator config get api-url
   
   # If incorrect, update it:
   asset-generator config set api-url http://localhost:7801
   ```

3. **Test with verbose output:**
   ```bash
   asset-generator models list --verbose
   ```

4. **Check firewall settings:**
   - Ensure port 7801 is accessible
   - Verify no proxy blocking the connection

5. **Try alternative host:**
   ```bash
   asset-generator config set api-url http://127.0.0.1:7801
   ```

### Problem: API Key Authentication Failures

**Symptoms:**
- Error: `unauthorized`
- Error: `invalid API key`

**Solutions:**

1. **Set API key:**
   ```bash
   asset-generator config set api-key YOUR_API_KEY
   ```

2. **Use environment variable:**
   ```bash
   export ASSET_GENERATOR_API_KEY=your_key_here
   asset-generator models list
   ```

3. **Verify key in config:**
   ```bash
   asset-generator config view
   # Should show api-key: ******** (masked)
   ```

---

## Configuration Problems {#configuration-problems}

### Problem: "Config file not found"

**Symptoms:**
- Error: `config file not found`
- Commands fail with configuration errors

**Solutions:**

1. **Initialize configuration:**
   ```bash
   asset-generator config init
   ```

2. **Check config location:**
   ```bash
   asset-generator config view --verbose
   # Shows config file path
   ```

3. **Manually create config file:**
   ```bash
   mkdir -p ~/.asset-generator
   cat > ~/.asset-generator/config.yaml << EOF
   api-url: http://localhost:7801
   output-format: table
   EOF
   ```

### Problem: "Invalid format" error

**Symptoms:**
- Error: `invalid format`
- Output format not recognized

**Solutions:**

1. **Use valid formats only:**
   ```bash
   # Valid: table, json, yaml
   asset-generator models list --format json
   ```

2. **Check config file syntax:**
   ```bash
   asset-generator config view
   # Verify YAML syntax is correct
   ```

3. **Reset configuration:**
   ```bash
   rm -rf ~/.asset-generator
   asset-generator config init
   ```

---

## Generation Failures {#generation-failures}

### Problem: "Prompt is required" error

**Symptoms:**
- Error: `prompt is required`
- Generation command fails immediately

**Solutions:**

1. **Always include prompt:**
   ```bash
   asset-generator generate image --prompt "your prompt here"
   ```

2. **Check prompt length:**
   ```bash
   # Ensure prompt is not empty or just whitespace
   asset-generator generate image --prompt "detailed fantasy landscape"
   ```

### Problem: Model not found

**Symptoms:**
- Error: `model validation failed: model 'xyz' not found`
- Error: `model not available`

**Solutions:**

1. **List available models:**
   ```bash
   asset-generator models list
   ```

2. **Use exact model name:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --model "stable-diffusion-xl"  # Use exact name
   ```

3. **Check model loading status:**
   ```bash
   asset-generator models list --format json | jq '.[] | select(.loaded == true)'
   ```

### Problem: Generation timeout

**Symptoms:**
- Error: `generation timeout`
- Process hangs for extended periods

**Solutions:**

1. **Check server status:**
   ```bash
   asset-generator status
   ```

2. **Reduce complexity:**
   ```bash
   asset-generator generate image \
     --prompt "simple test" \
     --steps 20 \
     --width 512 \
     --height 512
   ```

3. **Use continue-on-error for pipelines:**
   ```bash
   asset-generator pipeline --file assets.yaml --continue-on-error
   ```

---

## Pipeline Issues {#pipeline-issues}

### Problem: YAML syntax errors

**Symptoms:**
- Error: `failed to parse YAML`
- Error: `mapping values are not allowed in this context`

**Solutions:**

1. **Validate YAML syntax:**
   ```bash
   # Install yamllint if available
   yamllint pipeline.yaml
   ```

2. **Check common issues:**
   - Use 2 spaces for indentation (not tabs)
   - Ensure colons after keys
   - Quote strings with special characters

3. **Test with minimal pipeline:**
   ```yaml
   assets:
     - name: Test
       assets:
         - name: test-image
           prompt: "simple test prompt"
   ```

4. **Use dry run to validate:**
   ```bash
   asset-generator pipeline --file pipeline.yaml --dry-run
   ```

### Problem: Pipeline file not found

**Symptoms:**
- Error: `failed to read pipeline file: no such file or directory`

**Solutions:**

1. **Use absolute paths:**
   ```bash
   asset-generator pipeline --file /full/path/to/pipeline.yaml
   ```

2. **Check current directory:**
   ```bash
   ls -la *.yaml
   pwd
   ```

3. **Verify file permissions:**
   ```bash
   ls -la pipeline.yaml
   # Should be readable (r-- at minimum)
   ```

### Problem: Interrupted pipeline

**Symptoms:**
- Pipeline stops mid-generation
- Partial results in output directory

**Solutions:**

1. **Check which assets were generated:**
   ```bash
   find output/ -name "*.png" | wc -l
   ```

2. **Resume with continue-on-error:**
   ```bash
   asset-generator pipeline --file assets.yaml --continue-on-error
   ```

3. **Remove completed items from pipeline file and re-run**

---

## LoRA Problems {#lora-problems}

### Problem: LoRA not taking effect

**Symptoms:**
- Generated images don't show expected style
- LoRA appears to have no impact

**Solutions:**

1. **Check LoRA name (case-sensitive):**
   ```bash
   asset-generator models list | grep -i lora
   ```

2. **Increase weight:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --lora "style-name:1.2"  # Try higher weight
   ```

3. **Verify model compatibility:**
   - SD 1.5 LoRAs work with SD 1.5 models
   - SDXL LoRAs work with SDXL models

4. **Test with single LoRA:**
   ```bash
   asset-generator generate image \
     --prompt "anime character" \
     --lora "anime-style:1.0"  # Test one at a time
   ```

### Problem: LoRA conflicts or unexpected results

**Symptoms:**
- Weird artifacts in generated images
- Multiple LoRAs producing strange combinations

**Solutions:**

1. **Reduce number of LoRAs:**
   ```bash
   # Use maximum 2-3 LoRAs
   asset-generator generate image \
     --prompt "portrait" \
     --lora "style1:0.8" \
     --lora "details:0.6"
   ```

2. **Adjust weights:**
   ```bash
   # Lower weights for subtle effects
   asset-generator generate image \
     --prompt "character" \
     --lora "strong-style:0.6"  # Reduce from 1.0 to 0.6
   ```

3. **Use negative weights to remove conflicts:**
   ```bash
   asset-generator generate image \
     --prompt "realistic portrait" \
     --lora "realism:1.0" \
     --lora "cartoon:-0.3"  # Remove cartoon elements
   ```

---

## File and Permission Issues {#file-and-permission-issues}

### Problem: Permission denied errors

**Symptoms:**
- Error: `permission denied`
- Cannot create output directories
- Cannot write configuration files

**Solutions:**

1. **Check directory permissions:**
   ```bash
   ls -la output/
   # Ensure write permissions (w)
   ```

2. **Use different output directory:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --save-images \
     --output-dir ~/my-images
   ```

3. **Fix permissions:**
   ```bash
   chmod 755 output/
   chmod 644 ~/.asset-generator/config.yaml
   ```

### Problem: Input file does not exist

**Symptoms:**
- Error: `input file does not exist` (for conversions)
- File path errors

**Solutions:**

1. **Check file path and extension:**
   ```bash
   ls -la input.png
   # Verify file exists and extension is correct
   ```

2. **Use absolute paths:**
   ```bash
   asset-generator convert svg /full/path/to/image.png
   ```

3. **Verify file permissions:**
   ```bash
   ls -la image.png
   # Should be readable (r--)
   ```

---

## Performance Problems {#performance-problems}

### Problem: Slow generation

**Symptoms:**
- Generation takes excessively long
- Pipeline processes very slowly

**Solutions:**

1. **Reduce complexity:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --steps 20 \        # Reduce from 50
     --width 512 \       # Smaller dimensions
     --height 512 \
     --scheduler simple  # Faster scheduler
   ```

2. **Use faster sampler:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --sampler euler_a
   ```

3. **Check server resources:**
   ```bash
   asset-generator status
   # Check GPU usage, available models
   ```

### Problem: Out of disk space

**Symptoms:**
- Error: `no space left on device`
- Generation failures due to disk space

**Solutions:**

1. **Check disk space:**
   ```bash
   df -h
   ```

2. **Clean up old images:**
   ```bash
   find output/ -name "*.png" -mtime +7 -delete
   ```

3. **Use smaller image dimensions:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --width 512 \
     --height 512
   ```

4. **Downscale automatically:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --save-images \
     --downscale-width 1024
   ```

---

## Output and Format Issues {#output-and-format-issues}

### Problem: Invalid filename template

**Symptoms:**
- Error with filename template
- Files not saving with expected names

**Solutions:**

1. **Test template with single image:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --batch 1 \
     --save-images \
     --filename-template "test-{index}.png"
   ```

2. **Use simple templates:**
   ```bash
   # Safe template
   --filename-template "{date}-{index}.png"
   ```

3. **Avoid special characters:**
   ```bash
   # The {prompt} placeholder automatically sanitizes
   --filename-template "{prompt}-{index}.png"
   ```

### Problem: No images downloaded

**Symptoms:**
- Generation succeeds but no local files
- Missing `--save-images` flag effects

**Solutions:**

1. **Ensure save-images is enabled:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --save-images  # This flag is required
   ```

2. **Check output directory:**
   ```bash
   ls -la ./  # Default output directory
   ls -la output/  # Check if using custom directory
   ```

3. **Verify generation completed:**
   ```bash
   asset-generator generate image \
     --prompt "test" \
     --save-images \
     --verbose  # Shows download progress
   ```

---

## Feature-Specific Issues {#feature-specific-issues}

### Skimmed CFG Issues

#### Problem: No visible effect

**Solutions:**
- Verify model supports Skimmed CFG
- Try more pronounced scale differences (2.0 vs 4.0)
- Check SwarmUI version compatibility

```bash
asset-generator generate image \
  --prompt "test" \
  --skimmed-cfg \
  --skimmed-cfg-scale 2.0  # Try different scales
```

#### Problem: Unexpected results

**Solutions:**
- Start with scale 3.0 and adjust incrementally
- Try mid-phase only (0.3-0.7) first
- Compare with and without Skimmed CFG

```bash
asset-generator generate image \
  --prompt "portrait" \
  --skimmed-cfg \
  --skimmed-cfg-start 0.3 \
  --skimmed-cfg-end 0.7
```

### SVG Conversion Issues

#### Problem: Output looks pixelated

**Solutions:**
- Increase number of shapes
- Try different conversion methods
- Ensure input has sufficient resolution

```bash
asset-generator convert svg input.png \
  --shapes 300 \  # Increase from default 100
  --mode 3        # Try ellipses for smoother look
```

#### Problem: Processing takes too long

**Solutions:**
- Reduce number of shapes
- Use simpler shape modes
- Set repeat to 0

```bash
asset-generator convert svg input.png \
  --shapes 50 \   # Reduce from 100
  --mode 1 \      # Triangles are faster
  --repeat 0      # No repetition
```

---

## Error Message Reference {#error-message-reference}

### Common Error Messages and Solutions

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `connection refused` | SwarmUI not running | Start SwarmUI service |
| `prompt is required` | Missing --prompt flag | Add --prompt "your text" |
| `model validation failed` | Invalid model name | Use `models list` to find correct name |
| `failed to parse YAML` | Invalid YAML syntax | Check indentation and colons |
| `permission denied` | File/directory permissions | Fix with chmod or use different path |
| `no such file or directory` | File path incorrect | Check path and file existence |
| `invalid format` | Bad output format | Use table, json, or yaml |
| `config file not found` | No configuration | Run `config init` |
| `generation timeout` | Server overloaded | Reduce complexity or retry |
| `no space left on device` | Disk full | Clean up files or use smaller images |

### Debug Mode

For any persistent issues, enable verbose output:

```bash
asset-generator --verbose [command] [flags]
```

This shows:
- Configuration file location
- API requests and responses
- Timing information
- Detailed error contexts

### Getting Help

1. **Check command help:**
   ```bash
   asset-generator --help
   asset-generator generate --help
   asset-generator generate image --help
   ```

2. **Test basic connectivity:**
   ```bash
   asset-generator status
   asset-generator models list
   ```

3. **Verify configuration:**
   ```bash
   asset-generator config view
   ```

4. **Report bugs:**
   - Include verbose output
   - Provide exact command used
   - Specify OS and asset-generator version
   - Include relevant config files (redact API keys)

---

## Quick Diagnostic Commands

Use these commands to quickly diagnose issues:

```bash
# Test basic functionality
asset-generator status

# Check configuration
asset-generator config view

# Test API connection
asset-generator models list

# Verify simple generation
asset-generator generate image --prompt "test" --steps 10

# Check disk space
df -h

# Test with minimal pipeline
echo 'assets:
  - name: Test
    assets:
      - name: test
        prompt: "simple test"' > test.yaml && \
asset-generator pipeline --file test.yaml --dry-run
```

## See Also

- [Quick Start Guide](QUICKSTART.md) - Basic setup and usage
- [Commands Reference](COMMANDS.md) - Complete command documentation
- [User Guide](USER_GUIDE.md) - Advanced features and generation options
- [Pipeline Processing](PIPELINE.md) - Batch generation workflows
- [Development Guide](DEVELOPMENT.md) - Architecture and debugging information