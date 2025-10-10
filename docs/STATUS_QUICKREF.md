# Status Command - Quick Reference

## Basic Usage

```bash
# Check server status
asset-generator status

# JSON output
asset-generator status --format json

# YAML output
asset-generator status --format yaml

# Verbose mode
asset-generator status -v
```

## What It Shows

✅ **Server connectivity** - Is the server reachable?
✅ **Response time** - How fast is the server responding?
✅ **Session info** - Current session ID
✅ **Backend status** - What backends are running and their state
✅ **Model info** - How many models are available and loaded
✅ **System info** - GPU, memory, and other system details (if available)

## Common Use Cases

### Health Check Before Generation

```bash
# Quick check before running a big job
if asset-generator status > /dev/null 2>&1; then
    asset-generator generate image --prompt "..."
else
    echo "Server is down!"
    exit 1
fi
```

### Monitoring Script

```bash
# Check status every minute
while true; do
    asset-generator status
    sleep 60
done
```

### Get Specific Info (JSON)

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

## Output Example

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

## Status Colors

- 🟢 **Green**: running, loaded, active, online, ready
- 🟡 **Yellow**: idle, unloaded
- 🔴 **Red**: error, failed, offline

## Exit Codes

- `0` - Server is online and responding
- `1` - Server is offline or error occurred

## Troubleshooting

**Server offline?**
```bash
# Check if SwarmUI is running
curl http://localhost:7801

# Update API URL if needed
asset-generator config set api-url http://localhost:7801
```

**Slow response?**
- Server may be under load
- Check with `--verbose` for details

## Related Commands

- `asset-generator models list` - Detailed model information
- `asset-generator config get api-url` - View API endpoint
- `asset-generator generate` - Generate images (requires online server)

## See Also

- [Full Documentation](STATUS_COMMAND.md)
- [API Documentation](API.md)
