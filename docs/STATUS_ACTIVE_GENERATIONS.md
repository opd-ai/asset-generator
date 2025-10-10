# Status Command - Active Generation Tracking

## Overview

The `status` command now displays information about currently running generation sessions. This enhancement provides real-time visibility into ongoing image generation processes.

## What's New

### Active Generations Section

The status command now includes an "Active Generations" section that shows:

- **Session ID**: Unique identifier for each generation
- **Status**: Current state (generating, starting, pending)
- **Progress**: Percentage complete (0-100%)
- **Duration**: How long the generation has been running

### JSON Output

When using `--format json`, the output includes:

```json
{
  "active_generations": [
    {
      "session_id": "abc123",
      "status": "generating",
      "progress": 0.45,
      "start_time": "2025-10-10T12:34:56Z",
      "duration": "2.5m"
    }
  ],
  "generations_running": 1
}
```

## Example Output

### No Active Generations

```
SwarmUI Server Status
═══════════════════════════════════════════════

Server URL:      http://localhost:7801
Status:          online
Response Time:   145ms
Version:         1.0.0

Active Generations
───────────────────────────────────────────────
No generations currently running

Models
───────────────────────────────────────────────
Total Available: 15
Currently Loaded: 2
```

### With Active Generations

```
SwarmUI Server Status
═══════════════════════════════════════════════

Server URL:      http://localhost:7801
Status:          online
Response Time:   145ms
Version:         1.0.0

Active Generations
───────────────────────────────────────────────
  Generation 1:
    Session ID:    gen-abc123
    Status:        generating
    Progress:      45.0%
    Duration:      2.5m

  Generation 2:
    Session ID:    gen-def456
    Status:        starting
    Progress:      5.0%
    Duration:      12s

Models
───────────────────────────────────────────────
Total Available: 15
Currently Loaded: 2
```

## Status Colors

Generation statuses are color-coded in table output:
- 🟢 **Green**: `generating` - actively processing
- 🟡 **Yellow**: `starting`, `pending` - initializing
- 🔴 **Red**: (reserved for error states)

## Use Cases

### 1. Monitor Long-Running Generations

Check on Flux or other slow models that take several minutes:

```bash
# Start a generation in another terminal
asset-generator generate image --prompt "complex scene" --model flux

# In this terminal, check progress
asset-generator status
```

### 2. Automation Scripts

Query status to see if generation is still in progress:

```bash
#!/bin/bash
while true; do
  running=$(asset-generator status --format json | jq '.generations_running')
  if [ "$running" -eq 0 ]; then
    echo "All generations complete"
    break
  fi
  echo "Still running: $running generation(s)"
  sleep 10
done
```

### 3. Pre-flight Checks

Verify no other generations are running before starting a new one:

```bash
if [ "$(asset-generator status --format json | jq '.generations_running')" -eq 0 ]; then
  asset-generator generate image --prompt "..."
else
  echo "Another generation is in progress, waiting..."
fi
```

## How It Works

The CLI tracks generation sessions using a combination of in-memory and file-based state. When you run `status`, it:

1. Loads persisted state from `.asset-generator-state.json` in the current working directory
2. Queries the SwarmUI server for connectivity and backend info
3. Checks both file-based and in-memory session tracking for active generations
4. Displays any sessions with status: `generating`, `starting`, or `pending`
5. Shows progress and elapsed time for each active session

### File-Based State Sharing (NEW!)

Sessions are now persisted to a state file, enabling cross-process tracking:

**State File**: `.asset-generator-state.json` (in current working directory)

```json
{
  "sessions": {
    "gen-abc123": {
      "id": "gen-abc123",
      "status": "generating",
      "progress": 0.45,
      "start_time": "2025-10-10T14:30:00Z",
      "updated_at": "2025-10-10T14:32:15Z"
    }
  },
  "updated_at": "2025-10-10T14:32:15Z"
}
```

### Cross-Process Tracking Now Works! ✅

**Before**: Could only track within the same CLI process
**Now**: Can track across different terminals if in the same directory

```bash
# Terminal 1
cd ~/my-project
asset-generator generate image --prompt "artwork"

# Terminal 2 (different process!)
cd ~/my-project
asset-generator status
# ✅ Shows full generation details!
```

### Important: Directory Matters

The state file is directory-specific:
- ✅ Same directory = shared state
- ❌ Different directory = different state file

This provides project isolation - each directory has its own independent state.

### Alternative Detection

If state file is unavailable (different directory, deleted, etc.), the CLI can still detect generations by checking backend status:

```bash
asset-generator status --format json | jq '.backends[].status'
```

If a backend shows status `running`, it's likely processing a generation.

See [State File Sharing Documentation](STATE_FILE_SHARING.md) for complete details.

## Related Commands

- `asset-generator generate image` - Start a new generation
- `asset-generator cancel` - Cancel the current generation
- `asset-generator cancel --all` - Cancel all queued generations

## Implementation Details

### Added to pkg/client/client.go:
- `ActiveGeneration` struct to hold generation info
- `GetActiveGenerations()` method to retrieve active sessions
- `formatDuration()` helper for human-readable time display
- Extended `ServerStatus` struct with `ActiveGenerations` and `GenerationsRunning` fields

### Updated in cmd/status.go:
- Enhanced `formatStatusTable()` to display active generations
- Updated `colorizeStatus()` to include generation-related statuses
- Added "Active Generations" section to help text

### Tests:
- `TestFormatStatusTableWithActiveGenerations` validates the new output format
- Extended `TestColorizeStatus` to include `generating`, `starting`, and `pending` statuses

## Duration Formatting

Durations are displayed in human-readable format:
- Under 1 minute: `45s`
- 1-60 minutes: `2.5m`
- Over 1 hour: `1.2h`

## Limitations

- Only tracks generations from the current CLI process
- Sessions are stored in memory and cleared when the CLI exits
- WebSocket generations maintain more accurate progress than HTTP generations
- Completed or failed sessions are automatically cleaned up from the active list
