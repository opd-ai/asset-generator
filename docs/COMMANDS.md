# Commands Reference

This document provides comprehensive documentation for all Asset Generator CLI commands.

## Table of Contents

- [Cancel Command](#cancel-command)
  - [Overview](#cancel-overview)
  - [Usage](#cancel-usage)
  - [Examples](#cancel-examples)
  - [Flags](#cancel-flags)
  - [Quick Reference](#cancel-quick-reference)
- [Status Command](#status-command)
  - [Overview](#status-overview)
  - [Usage](#status-usage)
  - [Output Formats](#status-output-formats)
  - [Examples](#status-examples)
  - [Quick Reference](#status-quick-reference)

---

## Cancel Command {#cancel-command}

The `cancel` command allows you to interrupt ongoing or queued image generations on the SwarmUI server.

### Overview {#cancel-overview}

During long-running generation tasks (especially with models like Flux that can take 5-40 minutes), you may need to cancel the generation for various reasons:
- The generation is taking too long
- You made a mistake in your prompt
- You need to free up server resources
- You have queued multiple generations and want to clear the backlog

The `cancel` command provides two modes of operation:
1. **Single Cancel**: Cancel only the current generation in progress
2. **Cancel All**: Cancel all queued generations

### Usage {#cancel-usage}

#### Basic Usage

```bash
# Cancel the current generation
asset-generator cancel

# Cancel all queued generations
asset-generator cancel --all

# Cancel with verbose output to see API details
asset-generator cancel -v

# Cancel quietly (for scripts)
asset-generator cancel -q
```

### Examples {#cancel-examples}

#### Example 1: Cancel a Long-Running Generation

```bash
# Start a long-running generation
asset-generator generate image --prompt "detailed fantasy landscape" --model flux &

# Cancel it if you change your mind
asset-generator cancel
```

**Output:**
```
Cancelling current generation...
✓ Successfully cancelled current generation
```

#### Example 2: Clear All Queued Generations

If you've submitted multiple generations and want to cancel them all:

```bash
# Submit multiple generations
asset-generator generate image --prompt "prompt 1" --batch 5 &
asset-generator generate image --prompt "prompt 2" --batch 5 &
asset-generator generate image --prompt "prompt 3" --batch 5 &

# Cancel all of them
asset-generator cancel --all
```

**Output:**
```
Cancelling all queued generations...
✓ Successfully cancelled all queued generations
```

#### Example 3: Script Integration

For use in scripts where you don't want output:

```bash
#!/bin/bash
# Emergency stop script
asset-generator cancel --all -q || echo "Cancel failed"
```

#### Example 4: Interactive Development

When iterating on prompts:

```bash
# Oops, forgot to add negative prompt
asset-generator generate image --prompt "cat" &
asset-generator cancel
asset-generator generate image --prompt "cat" --negative-prompt "blurry, low quality"
```

#### Example 5: Resource Management

If the server is running slow and you need to free up resources:

```bash
# Check what's running first
asset-generator status
# Cancel everything to free up GPU memory
asset-generator cancel --all
```

### Flags {#cancel-flags}

| Flag | Short | Description |
|------|-------|-------------|
| `--all` | | Cancel all queued generations instead of just the current one |
| `--quiet` | `-q` | Suppress success messages (errors only) |
| `--verbose` | `-v` | Show detailed API communication |

#### Global Flags

All global flags from the root command are available:
- `--config`: Specify config file location
- `--api-url`: Override the configured API URL
- `--api-key`: Override the configured API key

### How It Works

The cancel command uses SwarmUI's interrupt API endpoints:

1. **Single Cancel** (`asset-generator cancel`):
   - Calls `/API/InterruptGeneration` endpoint
   - Interrupts the currently running generation
   - Queued generations continue processing

2. **Cancel All** (`asset-generator cancel --all`):
   - Calls `/API/InterruptAll` endpoint
   - Interrupts the current generation
   - Clears all queued generations
   - Server returns to idle state

### Error Handling

#### Common Errors

**No generations to cancel:**
```
failed to cancel generation: SwarmUI error: no active generation
```
This is informational - there was nothing to cancel.

**Server not responding:**
```
failed to cancel generation: request failed: context deadline exceeded
```
Check server status: `asset-generator status`

### Best Practices

1. **Use Single Cancel for Specific Generations**: Stop current generation, let queued ones continue
2. **Use Cancel All for Full Reset**: Clean slate when needed
3. **Combine with Status Checks**: `asset-generator status` before and after
4. **In Scripts, Use Quiet Mode**: `asset-generator cancel --all -q`
5. **Use Verbose Mode for Debugging**: `asset-generator cancel -v`

### Limitations

- No partial cancel (cannot cancel specific generations by ID)
- No undo (cancelled generations cannot be resumed)
- Server-side operation (affects all clients)
- Race conditions possible (generation might complete before cancel)

---

### Cancel Quick Reference {#cancel-quick-reference}

#### Synopsis
```bash
asset-generator cancel [--all] [-q] [-v]
```

#### Common Use Cases

| Task | Command |
|------|---------|
| Stop current generation | `asset-generator cancel` |
| Clear all pending work | `asset-generator cancel --all` |
| Check before cancelling | `asset-generator status && asset-generator cancel` |
| Script integration | `asset-generator cancel --all -q` |

#### Exit Codes

- `0` - Success
- `1` - Error (connection, API, etc.)

#### API Endpoints

| Operation | Endpoint | Description |
|-----------|----------|-------------|
| Single | `/API/InterruptGeneration` | Cancel current generation |
| All | `/API/InterruptAll` | Cancel all queued generations |

#### Error Messages

| Error | Meaning |
|-------|---------|
| "no active generation" | Nothing to cancel (informational) |
| "failed to get session" | Session/connection issue |
| "request failed" | Network/server problem |

---

## Status Command {#status-command}

The `status` command queries the SwarmUI server and displays comprehensive information about its current state.

### Overview {#status-overview}

The status command provides real-time information about:
- Server connectivity and response time
- Available backends and their operational states
- Current session information
- Model availability and loading status
- System information (GPU, memory, etc.)

### Usage {#status-usage}

```bash
# Basic status check
asset-generator status

# JSON output (for scripting/automation)
asset-generator status --format json

# YAML output
asset-generator status --format yaml

# Verbose output with additional details
asset-generator status -v

# Save status to file
asset-generator status --output status.txt
```

### Output Fields

#### Server Information
- **Server URL**: The SwarmUI API endpoint being queried
- **Status**: Connection status (online/offline)
- **Response Time**: Time taken to query the server
- **Version**: SwarmUI server version (if available)

#### Session Information
- **Session ID**: Current active session identifier

#### Backend Information
For each available backend:
- **ID**: Unique backend identifier
- **Type**: Backend type (e.g., ComfyUI, Automatic1111)
- **Status**: Operational status (running, idle, error)
- **Model Loaded**: Currently loaded model (if any)
- **GPU**: GPU device information

#### Model Information
- **Total Available**: Count of all available models
- **Currently Loaded**: Number of models loaded in memory

### Output Formats {#status-output-formats}

#### Table Format (Default)

Human-readable display with color coding:

```
SwarmUI Server Status
═══════════════════════════════════════════════

Server URL:      http://localhost:7801
Status:          online
Response Time:   145ms
Version:         1.0.0

Session
───────────────────────────────────────────────
Session ID:      abc123def456

Backends
───────────────────────────────────────────────
  • backend-1
    Type:          ComfyUI
    Status:        running
    Model Loaded:  stable-diffusion-xl
    GPU:           NVIDIA RTX 3090

Models
───────────────────────────────────────────────
Total Available: 15
Currently Loaded: 2
```

**Status Colors:**
- 🟢 Green: running, loaded, active, online, ready
- 🟡 Yellow: idle, unloaded
- 🔴 Red: error, failed, offline

#### JSON Format

```bash
asset-generator status --format json
```

```json
{
  "server_url": "http://localhost:7801",
  "status": "online",
  "response_time": "145ms",
  "session_id": "abc123def456",
  "backends": [
    {
      "id": "backend-1",
      "type": "ComfyUI",
      "status": "running"
    }
  ],
  "models_count": 15
}
```

#### YAML Format

```bash
asset-generator status --format yaml
```

### Examples {#status-examples}

#### Example 1: Health Check Before Generation

```bash
# Quick check before running a big job
if asset-generator status > /dev/null 2>&1; then
    asset-generator generate image --prompt "..."
else
    echo "Server is down!"
    exit 1
fi
```

#### Example 2: Monitoring Script

```bash
# Check status every minute
while true; do
    asset-generator status
    sleep 60
done
```

#### Example 3: Extract Specific Information (JSON)

```bash
# Get just the status
asset-generator status --format json | jq -r '.status'

# Count available models
asset-generator status --format json | jq '.models_count'

# List backends
asset-generator status --format json | jq '.backends[].id'

# Check response time
asset-generator status --format json | jq -r '.response_time'
```

#### Example 4: Monitoring Different Servers

```bash
# Check production server
asset-generator status --api-url https://prod.example.com:7801

# Check local development server
asset-generator status --api-url http://localhost:7801
```

#### Example 5: CI/CD Integration

```yaml
# Example GitHub Actions workflow
- name: Check SwarmUI Health
  run: |
    asset-generator status --api-url $SWARM_URL
  env:
    SWARM_URL: ${{ secrets.SWARM_API_URL }}
```

### Troubleshooting

#### Server Offline Error

```
Error: failed to get server status: server unreachable
```

**Solutions:**
1. Verify SwarmUI is running: `curl http://localhost:7801`
2. Check the API URL: `asset-generator config get api-url`
3. Update if needed: `asset-generator config set api-url http://localhost:7801`

#### Slow Response Time

High response times (>5 seconds) may indicate:
- Server under heavy load
- Network latency issues
- Server processing other requests

Use `--verbose` flag to see detailed API request information.

---

### Status Quick Reference {#status-quick-reference}

#### Synopsis
```bash
asset-generator status [--format FORMAT] [-v]
```

#### What It Shows

✅ Server connectivity and response time  
✅ Session information  
✅ Backend status and loaded models  
✅ Model counts  
✅ System information (if available)

#### Common Use Cases

| Task | Command |
|------|---------|
| Basic check | `asset-generator status` |
| JSON for scripting | `asset-generator status --format json` |
| Check before generation | `asset-generator status && asset-generator generate ...` |
| Get specific field | `asset-generator status --format json \| jq '.status'` |

#### Output Formats

| Format | Flag | Use Case |
|--------|------|----------|
| Table | (default) | Human-readable |
| JSON | `--format json` | Scripting, parsing |
| YAML | `--format yaml` | Configuration, readability |

#### Exit Codes

- `0` - Server online and responding
- `1` - Server offline or error

#### API Endpoints Used

- `/API/GetNewSession` - Connectivity check
- `/API/ListBackends` - Backend information
- `/API/ListModels` - Model counts

---

## See Also

- [Generation Features](GENERATION_FEATURES.md) - Scheduler, Skimmed CFG, and other generation parameters
- [Pipeline Processing](PIPELINE.md) - Batch generation workflows
- [Quick Start Guide](QUICKSTART.md) - Getting started with the CLI
- [Development Documentation](DEVELOPMENT.md) - API integration details
