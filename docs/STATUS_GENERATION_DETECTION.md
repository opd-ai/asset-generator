# Generation Detection in Status Command

## Problem Statement

The `status` command needs to show if there is currently a generation in process. However, this presents a technical challenge due to how CLI tools work.

## The Challenge

### Process Isolation
Each time you run `asset-generator`, it creates a new process with its own memory:

```bash
# Terminal 1: Start generation (Process A)
asset-generator generate image --prompt "..."

# Terminal 2: Check status (Process B - different process!)
asset-generator status
```

Process B has no access to Process A's memory, so it can't see the generation session tracked in Process A.

### No SwarmUI API
SwarmUI doesn't currently provide an API endpoint to query all active generations globally. The only way to track generations is locally within each CLI process.

## The Solution

We implemented a **two-tier detection system**:

### Tier 1: Local Session Tracking
When a generation runs in the same process, we track it locally:
- Session ID
- Status (generating, starting, pending)
- Progress (0-100%)
- Duration

**Works for**: Checking status while generation is running in the same process

### Tier 2: Backend Status Inference
When local tracking unavailable, we infer from backend status:
- Check if any backend has status "running" or "generating"
- Count running backends as estimated active generations
- Display: "Estimated running: X (inferred from backend status)"

**Works for**: Detecting generations from other processes, terminals, or the web UI

## Output Examples

### Scenario 1: No Generations Running

```
Active Generations
───────────────────────────────────────────────
No generations currently running
```

### Scenario 2: Local Tracking Available (Same Process)

```
Active Generations
───────────────────────────────────────────────
  Generation 1:
    Session ID:    gen-abc123
    Status:        generating
    Progress:      45.0%
    Duration:      2.5m
```

### Scenario 3: Inferred from Backend (Cross-Process)

```
Backends
───────────────────────────────────────────────
  • backend-1
    Type:          ComfyUI
    Status:        running
    Model Loaded:  stable-diffusion-xl
    GPU:           NVIDIA RTX 3090

Active Generations
───────────────────────────────────────────────
Estimated running: 1 (inferred from backend status)

Note: Generation details unavailable - check backend status above
      for more information about active processing.
```

## Detection Logic

```go
// 1. Check local session tracking
activeGens := c.GetActiveGenerations()
status.GenerationsRunning = len(activeGens)

// 2. If no local sessions, infer from backends
if len(activeGens) == 0 && len(backends) > 0 {
    for _, backend := range backends {
        if backend.Status == "running" || backend.Status == "generating" {
            status.GenerationsRunning++
        }
    }
}
```

## Use Cases

### 1. Monitor Your Own Generation

Start a generation and check status in the same terminal:

```bash
# Start generation (will take a while)
asset-generator generate image --prompt "complex scene" &
PID=$!

# Check status while it runs
asset-generator status

# Wait for completion
wait $PID
```

### 2. Check If Server Is Busy

From any terminal, check if the server is processing:

```bash
asset-generator status

# Look for:
# - "Estimated running: X" in Active Generations
# - Backend status: "running"
```

### 3. Automation: Wait for Idle

```bash
#!/bin/bash
# Wait until no generations are running

while true; do
  running=$(asset-generator status --format json | jq '.generations_running')
  if [ "$running" -eq 0 ]; then
    echo "Server is idle"
    break
  fi
  echo "Server busy: $running generation(s)"
  sleep 10
done
```

### 4. JSON Output for Scripts

```bash
# Get generation count
asset-generator status --format json | jq '.generations_running'

# Check if any backends are running
asset-generator status --format json | jq '.backends[] | select(.status=="running")'
```

## Accuracy

### When Local Tracking Works
- ✅ **100% accurate** - session details, progress, duration
- ✅ Within same CLI process
- ✅ Real-time updates

### When Using Backend Inference
- ⚠️ **Estimated** - counts running backends
- ⚠️ Backends might be "running" for reasons other than generation
- ⚠️ Multiple generations on same backend show as one
- ✅ Works across processes and terminals
- ✅ Works for web UI generations too

## Limitations

1. **Process Isolation**: Can't share memory between separate CLI invocations
2. **No SwarmUI API**: No global generation query endpoint
3. **Backend Status Ambiguity**: "running" might not always mean generating
4. **No Progress for Inferred**: Can only estimate count, not see progress

## Future Improvements

Potential enhancements if SwarmUI adds APIs:

1. **Global Query API**: `GET /API/ListActiveGenerations`
   - Would return all active generations with progress
   - Enable accurate cross-process tracking

2. **Generation Events API**: WebSocket endpoint for generation events
   - Subscribe to generation start/stop events
   - Real-time notifications

3. **Session Persistence**: CLI could write to shared file
   - `/tmp/asset-generator-sessions.json`
   - Would enable cross-process tracking without server changes

## Related Documentation

- [Status Command](STATUS_COMMAND.md) - Full status command documentation
- [Active Generations Feature](STATUS_ACTIVE_GENERATIONS.md) - Feature overview
- [Cancel Command](CANCEL_COMMAND.md) - How to stop generations
- [WebSocket Support](API.md#websocket-support) - Real-time progress tracking

## Summary

The status command now intelligently detects generations using a two-tier approach:
1. **Local tracking** when possible (same process, full details)
2. **Backend inference** when needed (cross-process, estimated count)

This provides useful information in all scenarios while being honest about limitations when detection is based on inference rather than direct tracking.
