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

We implemented a **three-tier detection system**:

### Tier 1: File-Based State Sharing (NEW!)
Sessions are persisted to `.asset-generator-state.json` in the current working directory:
- Works across different CLI processes
- Tracks session ID, status, progress, and duration
- Automatically cleaned up after 24 hours
- Atomic writes prevent corruption

**Works for**: All scenarios where you run commands in the same directory

### Tier 2: Local Session Tracking
When a generation runs in the same process, we track it locally in memory:
- Session ID
- Status (generating, starting, pending)
- Progress (0-100%)
- Duration

**Works for**: Enhanced tracking within the same CLI process

### Tier 3: Backend Status Inference
When local tracking and file state are unavailable, we infer from backend status:
- Check if any backend has status "running" or "generating"
- Count running backends as estimated active generations
- Display: "Estimated running: X (inferred from backend status)"

**Works for**: Detecting generations started from web UI or other sources

## Output Examples

### Scenario 1: No Generations Running

```
Active Generations
───────────────────────────────────────────────
No generations currently running
```

### Scenario 2: File-Based Tracking (Cross-Process - NEW!)

```bash
# Terminal 1
cd ~/my-project
asset-generator generate image --prompt "artwork"

# Terminal 2 (different process, same directory!)
cd ~/my-project
asset-generator status
```

Output:
```
Active Generations
───────────────────────────────────────────────
  Generation 1:
    Session ID:    gen-abc123
    Status:        generating
    Progress:      45.0%
    Duration:      2.5m
```

### Scenario 3: Local Tracking (Same Process)

```
Active Generations
───────────────────────────────────────────────
  Generation 1:
    Session ID:    gen-abc123
    Status:        generating
    Progress:      45.0%
    Duration:      2.5m
```

### Scenario 4: Inferred from Backend (Web UI or Different Directory)

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
// 1. Load state from file (on client initialization)
client.loadStateFromFile()

// 2. Check local + file-based session tracking
activeGens := c.GetActiveGenerations()
status.GenerationsRunning = len(activeGens)

// 3. If no local/file sessions, infer from backends
if len(activeGens) == 0 && len(backends) > 0 {
    for _, backend := range backends {
        if backend.Status == "running" || backend.Status == "generating" {
            status.GenerationsRunning++
        }
    }
}

// 4. Persist state changes to file
c.saveStateToFile()
```

### State File: `.asset-generator-state.json`

Located in the current working directory:

```json
{
  "sessions": {
    "session-abc123": {
      "id": "session-abc123",
      "status": "generating",
      "progress": 0.45,
      "start_time": "2025-10-10T14:30:00Z",
      "updated_at": "2025-10-10T14:32:15Z"
    }
  },
  "updated_at": "2025-10-10T14:32:15Z"
}
```

See [State File Sharing Documentation](STATE_FILE_SHARING.md) for complete details.

## Use Cases

### 1. Monitor Your Own Generation (Cross-Process!)

Start a generation and check status from a different terminal:

```bash
# Terminal 1: Start generation
cd ~/my-project
asset-generator generate image --prompt "complex scene" &

# Terminal 2: Check status (same directory!)
cd ~/my-project  
asset-generator status
# Shows full generation details via state file!
```

### 2. Check If Server Is Busy

From any terminal, check if the server is processing:

```bash
cd ~/my-project
asset-generator status

# Look for:
# - Generation details from state file (if in same directory)
# - "Estimated running: X" from backend inference
# - Backend status: "running"
```

### 3. Automation: Wait for Idle

```bash
#!/bin/bash
# Wait until no generations are running

cd ~/my-project  # Important: be in the right directory!

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
cd ~/my-project

# Get generation count
asset-generator status --format json | jq '.generations_running'

# Get active session details
asset-generator status --format json | jq '.active_generations[]'

# Check if any backends are running
asset-generator status --format json | jq '.backends[] | select(.status=="running")'
```

## Accuracy

### When File-Based Tracking Works (NEW!)
- ✅ **100% accurate** - session details, progress, duration
- ✅ Works across different terminals/processes
- ✅ Persists even if CLI exits
- ✅ Real-time updates
- ⚠️ **Requires same working directory**

### When Local Tracking Works
- ✅ **100% accurate** - session details, progress, duration
- ✅ Within same CLI process
- ✅ Real-time updates
- ✅ Combined with file-based state

### When Using Backend Inference
- ⚠️ **Estimated** - counts running backends
- ⚠️ Backends might be "running" for reasons other than generation
- ⚠️ Multiple generations on same backend show as one
- ✅ Works across processes and terminals
- ✅ Works for web UI generations too
- ✅ No directory requirement

## Limitations

1. **Directory-Specific State**: File-based tracking requires same working directory
2. **No SwarmUI API**: No global generation query endpoint from server
3. **Backend Status Ambiguity**: "running" might not always mean generating
4. **No Progress for Inferred**: Can only estimate count, not see progress
5. **State File Management**: Need to be aware of `.asset-generator-state.json`

## Future Improvements

Potential enhancements:

1. **Global Query API**: `GET /API/ListActiveGenerations` from SwarmUI
   - Would return all active generations with progress
   - Enable accurate cross-directory tracking

2. **Generation Events API**: WebSocket endpoint for generation events
   - Subscribe to generation start/stop events
   - Real-time notifications

3. **Custom State Location**: `--state-file` flag for custom location
   - Share state across directories
   - Centralized state management

4. **Network State Sharing**: Redis/database backend
   - Share state across machines
   - Multi-user environments

## Related Documentation

- [State File Sharing](STATE_FILE_SHARING.md) - **NEW!** Complete guide to file-based state
- [Status Command](STATUS_COMMAND.md) - Full status command documentation
- [Active Generations Feature](STATUS_ACTIVE_GENERATIONS.md) - Feature overview
- [Cancel Command](CANCEL_COMMAND.md) - How to stop generations
- [WebSocket Support](API.md#websocket-support) - Real-time progress tracking

## Summary

The status command now intelligently detects generations using a three-tier approach:
1. **File-based state** (NEW!) - Cross-process tracking via `.asset-generator-state.json`
2. **Local tracking** - In-memory for same process (full details)
3. **Backend inference** - When state unavailable (estimated count)

This provides accurate cross-process generation tracking while being honest about limitations when detection is based on inference rather than direct tracking.
