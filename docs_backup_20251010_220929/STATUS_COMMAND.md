# Status Command

The `status` command queries the SwarmUI server and displays comprehensive information about its current state.

## Overview

The status command provides real-time information about:
- Server connectivity and response time
- Available backends and their operational states
- Current session information
- Model availability and loading status
- System information (GPU, memory, etc.)

## Usage

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

## Output Fields

### Server Information
- **Server URL**: The SwarmUI API endpoint being queried
- **Status**: Connection status (online/offline)
- **Response Time**: Time taken to query the server
- **Version**: SwarmUI server version (if available)

### Session Information
- **Session ID**: Current active session identifier

### Backend Information
For each available backend:
- **ID**: Unique backend identifier
- **Type**: Backend type (e.g., ComfyUI, Automatic1111)
- **Status**: Operational status (running, idle, error)
- **Model Loaded**: Currently loaded model (if any)
- **GPU**: GPU device information

### Model Information
- **Total Available**: Count of all available models
- **Currently Loaded**: Number of models loaded in memory

### System Information
Additional system details like GPU memory, CPU count, etc. (availability depends on SwarmUI configuration)

## Output Formats

### Table Format (Default)

The table format provides a human-readable display:

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

Status values are color-coded:
- 🟢 Green: running, loaded, active, online, ready
- 🟡 Yellow: idle, unloaded
- 🔴 Red: error, failed, offline

### JSON Format

```bash
asset-generator status --format json
```

```json
{
  "server_url": "http://localhost:7801",
  "status": "online",
  "response_time": "145ms",
  "version": "1.0.0",
  "session_id": "abc123def456",
  "backends": [
    {
      "id": "backend-1",
      "type": "ComfyUI",
      "status": "running",
      "model_loaded": "stable-diffusion-xl",
      "gpu": "NVIDIA RTX 3090"
    }
  ],
  "models_count": 15,
  "models_loaded": 2,
  "system_info": {
    "gpu_memory": "24GB",
    "cpu_count": 16
  }
}
```

### YAML Format

```bash
asset-generator status --format yaml
```

```yaml
server_url: http://localhost:7801
status: online
response_time: 145ms
version: 1.0.0
session_id: abc123def456
backends:
  - id: backend-1
    type: ComfyUI
    status: running
    model_loaded: stable-diffusion-xl
    gpu: NVIDIA RTX 3090
models_count: 15
models_loaded: 2
system_info:
  gpu_memory: 24GB
  cpu_count: 16
```

## Use Cases

### Health Monitoring

Check if the SwarmUI server is responsive:

```bash
if asset-generator status > /dev/null 2>&1; then
    echo "Server is online"
else
    echo "Server is offline or unreachable"
    exit 1
fi
```

### Pre-flight Checks

Verify server status before starting a long generation job:

```bash
# Check status and verify we have models available
asset-generator status --format json | jq -e '.models_count > 0'
if [ $? -eq 0 ]; then
    asset-generator generate image --prompt "..."
fi
```

### Monitoring Scripts

Regularly poll server status:

```bash
#!/bin/bash
while true; do
    asset-generator status --format json > /tmp/swarm-status.json
    # Parse and log status
    status=$(jq -r '.status' /tmp/swarm-status.json)
    echo "$(date): SwarmUI status: $status"
    sleep 60
done
```

### Backend Availability

Check which backends are available and their status:

```bash
asset-generator status --format json | jq '.backends[] | {id, status, model_loaded}'
```

## Exit Codes

- `0`: Server is reachable and responding normally
- `1`: Server is unreachable or returned an error

## API Endpoints Used

The status command queries the following SwarmUI API endpoints:

1. `POST /API/GetNewSession` - Verify connectivity and get session
2. `POST /API/ListBackends` - Get backend status (if available)
3. `POST /API/ListModels` - Get model information

Note: Some endpoints may not be available in all SwarmUI versions. The command gracefully handles missing endpoints and displays available information.

## Troubleshooting

### Server Offline Error

```
Error: failed to get server status: server unreachable: ...
```

**Solutions:**
1. Verify SwarmUI is running: `curl http://localhost:7801`
2. Check the API URL in your config: `asset-generator config get api-url`
3. Update the API URL if needed: `asset-generator config set api-url http://localhost:7801`

### Backend Information Not Available

Some SwarmUI versions may not expose the `ListBackends` endpoint. This is normal - the status command will display available information.

### Slow Response Time

High response times (>5 seconds) may indicate:
- Server is under heavy load
- Network latency issues
- Server is processing other requests

Use `--verbose` flag to see detailed API request information.

## Examples

### Basic Health Check

```bash
asset-generator status
```

### Automation-Friendly Output

```bash
# Get just the status field
asset-generator status --format json | jq -r '.status'

# Count available models
asset-generator status --format json | jq '.models_count'

# List backend IDs
asset-generator status --format json | jq -r '.backends[].id'
```

### Monitoring with Different Servers

```bash
# Check production server
asset-generator status --api-url https://prod.example.com:7801

# Check local development server
asset-generator status --api-url http://localhost:7801
```

## Integration with CI/CD

The status command is useful for health checks in automated workflows:

```yaml
# Example GitHub Actions workflow
- name: Check SwarmUI Health
  run: |
    asset-generator status --api-url $SWARM_URL
  env:
    SWARM_URL: ${{ secrets.SWARM_API_URL }}
```

## Related Commands

- `asset-generator models list` - List available models in detail
- `asset-generator config get api-url` - View current API URL
- `asset-generator generate` - Generate assets (requires server to be online)
