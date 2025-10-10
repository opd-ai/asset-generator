# File-Based State Sharing Implementation

## Overview

The Asset Generator CLI now persists generation session state to a file in the current working directory. This enables tracking of active generations across separate CLI process invocations.

## State File Details

### File Location
```
.asset-generator-state.json
```

The state file is created in the **current working directory** where you run the `asset-generator` command.

### Why Current Working Directory?

The state file is tied to the working directory for several important reasons:

1. **Project Isolation**: Each project/directory has its own independent state
2. **Multi-Project Support**: Work on multiple projects simultaneously without conflicts
3. **No Global State**: Avoids pollution of home directory or system temp
4. **Easy Cleanup**: Delete the directory, and the state goes with it
5. **Version Control**: Can be gitignored per-project

### File Format

```json
{
  "sessions": {
    "session-id-123": {
      "id": "session-id-123",
      "status": "generating",
      "progress": 0.45,
      "start_time": "2025-10-10T14:30:00Z",
      "updated_at": "2025-10-10T14:32:15Z"
    }
  },
  "updated_at": "2025-10-10T14:32:15Z"
}
```

## State Lifecycle

### 1. Session Creation
When a generation starts, a session is created and immediately persisted:

```bash
$ asset-generator generate image --prompt "a cat"
# Creates/updates .asset-generator-state.json with new session
```

### 2. Progress Updates
During generation, progress is periodically saved:
- **HTTP mode**: Updates on status changes
- **WebSocket mode**: Updates on every progress message from server

### 3. Session Completion
When generation completes or fails, the session is removed from the state file:

```bash
# Generation completes -> session removed from state file
# .asset-generator-state.json updated with empty sessions
```

### 4. State Cleanup
Old sessions are automatically cleaned up:
- **On client initialization**: Sessions older than 24 hours are removed
- **Incomplete sessions**: Sessions in "pending"/"starting" state are kept for cross-process tracking
- **Completed sessions**: Removed immediately after generation completes

## How It Works

### State Persistence Flow

```
┌─────────────────┐
│ Start Generation│
└────────┬────────┘
         │
         ▼
    ┌────────────────┐
    │ Create Session │
    │ in Memory      │
    └────────┬───────┘
             │
             ▼
      ┌──────────────┐
      │ Save to File │
      └──────┬───────┘
             │
             ▼
    ┌────────────────────┐
    │ Generation Running │
    │ (updates persisted)│
    └────────┬───────────┘
             │
             ▼
    ┌─────────────────┐
    │ Generation Done │
    └────────┬────────┘
             │
             ▼
    ┌──────────────────┐
    │ Remove from File │
    └──────────────────┘
```

### Cross-Process Detection

```
Terminal 1:                    Terminal 2:
┌─────────────────┐           ┌──────────────────┐
│ Start Generation│           │                  │
│                 │           │                  │
│ ↓ Save to file  │           │                  │
│                 │           │                  │
│ [Generating...] │───────────│ Check Status     │
│                 │  File     │                  │
│                 │  Shared   │ ↓ Read from file │
│                 │           │                  │
│                 │           │ Shows: Running!  │
└─────────────────┘           └──────────────────┘
```

## Usage Examples

### Example 1: Monitor from Different Terminal

```bash
# Terminal 1: Start a long-running generation
cd ~/my-project
asset-generator generate image --prompt "detailed artwork" --websocket

# Terminal 2: Check status (must be in same directory!)
cd ~/my-project
asset-generator status
```

**Output in Terminal 2:**
```
Active Generations
───────────────────────────────────────────────
  Generation 1:
    Session ID:    abc-123-def
    Status:        generating
    Progress:      45.0%
    Duration:      2.5m
```

### Example 2: Project Isolation

```bash
# Project A
cd ~/project-a
asset-generator generate image --prompt "logo design"

# Project B (different directory = different state!)
cd ~/project-b  
asset-generator status
# Shows: No generations currently running
# (because project-b has its own state file)
```

### Example 3: Automation with State File

```bash
#!/bin/bash
# Start generation in background
cd ~/my-project
asset-generator generate image --prompt "artwork" &

# Wait for generation to appear in state
sleep 2

# Monitor until complete
while [ -f .asset-generator-state.json ]; do
  count=$(jq '.sessions | length' .asset-generator-state.json 2>/dev/null || echo "0")
  if [ "$count" -eq 0 ]; then
    echo "Generation complete!"
    break
  fi
  echo "Still generating... ($count active)"
  sleep 5
done
```

### Example 4: Query State File Directly

```bash
# Check if any generations are active
jq '.sessions | length' .asset-generator-state.json

# Get session IDs
jq -r '.sessions | keys[]' .asset-generator-state.json

# Get progress of all sessions
jq -r '.sessions[] | "\(.id): \(.progress * 100)%"' .asset-generator-state.json
```

## State File Management

### Viewing State

```bash
# Pretty-print current state
cat .asset-generator-state.json | jq .

# Check last update time
jq -r '.updated_at' .asset-generator-state.json
```

### Manual Cleanup

```bash
# Remove state file (safe - will be recreated as needed)
rm .asset-generator-state.json

# Force cleanup of old sessions
asset-generator status  # This triggers cleanup on initialization
```

### Git Integration

Add to `.gitignore`:
```gitignore
# Asset Generator state (local session tracking)
.asset-generator-state.json
.asset-generator-state.json.tmp
```

## Technical Details

### Atomic Writes

State updates use atomic writes to prevent corruption:
1. Write to `.asset-generator-state.json.tmp`
2. Rename temp file to `.asset-generator-state.json`
3. Cleanup temp file on error

This ensures the state file is always valid JSON.

### Thread Safety

- In-memory state protected by `sync.RWMutex`
- File writes are non-blocking (don't block generation)
- Read-lock for queries, write-lock for updates

### Performance Considerations

**Write Frequency:**
- Initial session creation: 1 write
- Progress updates (WebSocket): ~1 write per update
- Progress updates (HTTP): ~1 write per status change
- Session completion: 1 write

**File Size:**
- Minimal: ~100-200 bytes per session
- JSON with indentation for human readability
- Old sessions automatically pruned

### Error Handling

**File I/O Errors:**
- Read errors on startup: Silently ignored (file may not exist)
- Write errors: Logged if verbose mode enabled
- Parse errors: State file ignored, new state created

**Fallback Behavior:**
If state file operations fail, the CLI continues to work using in-memory state only.

## Comparison: Before vs After

### Before (In-Memory Only)

```bash
# Terminal 1
asset-generator generate image --prompt "art"

# Terminal 2 (different process)
asset-generator status
# Shows: No generations running (❌ can't see Terminal 1's session)
```

### After (File-Based State)

```bash
# Terminal 1
cd ~/project
asset-generator generate image --prompt "art"

# Terminal 2 (same directory)
cd ~/project
asset-generator status
# Shows: Generation 1: generating, 45%, 2.5m (✅ sees shared state!)
```

## Troubleshooting

### State File Not Found

**Symptom:** Status shows "No generations running" but you know one is active

**Solutions:**
1. Check you're in the same directory where generation was started
2. Verify state file exists: `ls -la .asset-generator-state.json`
3. Check file permissions: `chmod 644 .asset-generator-state.json`

### Stale Sessions

**Symptom:** Old sessions showing as "running" but generation finished

**Solutions:**
1. Sessions older than 24h are auto-cleaned
2. Force cleanup: `asset-generator status` (triggers cleanup)
3. Manual removal: `rm .asset-generator-state.json`

### Permission Errors

**Symptom:** Can't write state file

**Solutions:**
1. Check directory permissions: `ls -ld .`
2. Ensure you own the directory: `stat .`
3. Fallback: CLI works in-memory-only mode if file write fails

## Security Considerations

### Sensitive Information

The state file contains:
- ✅ Session IDs (not sensitive - ephemeral identifiers)
- ✅ Status and progress (public information)
- ✅ Timestamps (not sensitive)
- ❌ NO prompts (not stored in state file)
- ❌ NO images (not stored in state file)
- ❌ NO API keys (not stored in state file)

### File Permissions

Default permissions: `0644` (owner read/write, others read)

For stricter security:
```bash
chmod 600 .asset-generator-state.json  # Owner only
```

## Future Enhancements

Potential improvements:

1. **Shared State Location**: Optional `--state-file` flag for custom location
2. **Network State**: Share state across machines (Redis, etc.)
3. **State History**: Keep log of completed sessions
4. **State Locking**: File locking for concurrent writes
5. **Compression**: Compress state for projects with many sessions

## Related Documentation

- [Generation Detection](STATUS_GENERATION_DETECTION.md) - How detection works
- [Status Command](STATUS_COMMAND.md) - Full status command documentation
- [Active Generations](STATUS_ACTIVE_GENERATIONS.md) - Feature overview

## Summary

File-based state sharing enables the Asset Generator CLI to track generation sessions across separate process invocations. By persisting state to `.asset-generator-state.json` in the working directory, you can now:

✅ Monitor generations from different terminals
✅ Track progress across CLI invocations  
✅ Build automation around persistent state
✅ Maintain project-isolated session tracking
✅ Query state directly via JSON

The implementation is transparent, atomic, and handles errors gracefully, providing a robust cross-process tracking solution.
