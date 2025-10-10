# Cancel Command

The `cancel` command allows you to interrupt ongoing or queued image generations on the SwarmUI server.

## Overview

During long-running generation tasks (especially with models like Flux that can take 5-40 minutes), you may need to cancel the generation for various reasons:
- The generation is taking too long
- You made a mistake in your prompt
- You need to free up server resources
- You have queued multiple generations and want to clear the backlog

The `cancel` command provides two modes of operation:
1. **Single Cancel**: Cancel only the current generation in progress
2. **Cancel All**: Cancel all queued generations

## Usage

### Basic Usage

```bash
# Cancel the current generation
asset-generator cancel

# Cancel all queued generations
asset-generator cancel --all

# Cancel with verbose output to see API details
asset-generator cancel -v
```

### Examples

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
asset-generator generate image --prompt "prompt 1" --images 5 &
asset-generator generate image --prompt "prompt 2" --images 5 &
asset-generator generate image --prompt "prompt 3" --images 5 &

# Cancel all of them
asset-generator cancel --all
```

**Output:**
```
Cancelling all queued generations...
✓ Successfully cancelled all queued generations
```

#### Example 3: Cancel in Quiet Mode

For use in scripts where you don't want output:

```bash
asset-generator cancel --all -q
```

No output is produced unless an error occurs.

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--all` | | Cancel all queued generations instead of just the current one |
| `--quiet` | `-q` | Suppress success messages (errors only) |
| `--verbose` | `-v` | Show detailed API communication |

### Global Flags

All global flags from the root command are available:
- `--config`: Specify config file location
- `--api-url`: Override the configured API URL
- `--api-key`: Override the configured API key

## How It Works

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

### Session Management

The cancel command automatically:
- Uses the current session ID (or creates one if needed)
- Handles session expiration and retries automatically
- Manages authentication if API keys are configured

## Use Cases

### 1. Interactive Development

When iterating on prompts, you might start a generation and realize you want to change the prompt:

```bash
# Oops, forgot to add negative prompt
asset-generator generate image --prompt "cat" &
asset-generator cancel
asset-generator generate image --prompt "cat" --negative-prompt "blurry, low quality"
```

### 2. Resource Management

If the server is running slow and you need to free up resources:

```bash
# Cancel everything to free up GPU memory
asset-generator cancel --all
```

### 3. Batch Processing Scripts

In automated workflows, you might want to cancel pending work:

```bash
#!/bin/bash
# Emergency stop script
asset-generator cancel --all -q
echo "All generations cancelled"
```

### 4. Testing and Development

During CLI development or testing, cancel is useful for cleaning up:

```bash
# Test script
for i in {1..10}; do
  asset-generator generate image --prompt "test $i" &
done

# Changed your mind?
asset-generator cancel --all
```

## Error Handling

### Common Errors

**No generations to cancel:**
```
failed to cancel generation: SwarmUI error: no active generation
```
This means there was nothing to cancel. This is informational, not an error.

**Session expired:**
The command automatically handles session expiration by obtaining a new session and retrying.

**Server not responding:**
```
failed to cancel generation: request failed: context deadline exceeded
```
The server may be overloaded or unreachable. Check server status with:
```bash
asset-generator status
```

### Verbose Mode

For debugging, use `-v` to see the API communication:

```bash
asset-generator cancel -v
```

**Output:**
```
Cancelling current generation...
Request: POST http://localhost:7801/API/InterruptGeneration
✓ Successfully cancelled current generation
```

## Integration with Other Commands

### With Generate Command

The cancel command complements the generate command:

```bash
# Start generation in background
asset-generator generate image --prompt "complex scene" --websocket &
GEN_PID=$!

# Monitor progress
watch -n 1 asset-generator status

# Cancel if taking too long
asset-generator cancel

# Kill background process
kill $GEN_PID
```

### With Status Command

Check what's running before cancelling:

```bash
# Check server status
asset-generator status

# If there are active generations, cancel them
asset-generator cancel --all
```

## Technical Details

### API Endpoints

The cancel command uses the following SwarmUI API endpoints:

**InterruptGeneration:**
- **Endpoint**: `/API/InterruptGeneration`
- **Method**: POST
- **Payload**: `{"session_id": "..."}`
- **Response**: `{"success": true}`

**InterruptAll:**
- **Endpoint**: `/API/InterruptAll`
- **Method**: POST
- **Payload**: `{"session_id": "..."}`
- **Response**: `{"success": true}`

### Implementation

The cancel functionality is implemented in:
- **Command**: `cmd/cancel.go` - CLI command definition
- **Client**: `pkg/client/client.go` - API client methods:
  - `Interrupt(ctx context.Context) error`
  - `InterruptAll(ctx context.Context) error`

### Context and Cancellation

The cancel command itself supports context cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := assetClient.Interrupt(ctx)
```

## Best Practices

### 1. Use Single Cancel for Specific Generations

If you only want to stop the current generation and let queued ones continue:

```bash
asset-generator cancel
```

### 2. Use Cancel All for Full Reset

When you need a clean slate:

```bash
asset-generator cancel --all
```

### 3. Combine with Status Checks

Always check status before and after cancelling:

```bash
asset-generator status
asset-generator cancel --all
asset-generator status  # Verify everything is cancelled
```

### 4. In Scripts, Use Quiet Mode

For automation, suppress output:

```bash
if ! asset-generator cancel --all -q; then
  echo "Warning: Failed to cancel generations"
fi
```

### 5. Use Verbose Mode for Debugging

When troubleshooting:

```bash
asset-generator cancel -v
```

## Limitations

1. **No Partial Cancel**: You cannot cancel specific generations by ID; only current or all
2. **No Undo**: Once cancelled, generations cannot be resumed
3. **Server-Side**: The cancel affects the server; other clients may also be affected
4. **Race Conditions**: A generation might complete before the cancel request is processed

## See Also

- [Status Command](STATUS_COMMAND.md) - Check server and generation status
- [Generate Command](GENERATE_PIPELINE.md) - Start image generation
- [API Documentation](API.md) - SwarmUI API details
- [Quick Reference](QUICKSTART.md) - Getting started guide

## Summary

The `cancel` command provides essential control over long-running generation tasks:

- **Simple**: Just run `asset-generator cancel`
- **Powerful**: Can cancel all queued work with `--all`
- **Reliable**: Automatic session management and error handling
- **Flexible**: Works in interactive and scripted scenarios

Use it whenever you need to interrupt generations and free up server resources.
