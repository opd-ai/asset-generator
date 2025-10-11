# Cancel Command - Quick Reference

## Synopsis

```bash
asset-generator cancel [--all] [-q] [-v]
```

## Quick Examples

```bash
# Cancel current generation
asset-generator cancel

# Cancel all queued generations  
asset-generator cancel --all

# Cancel quietly (for scripts)
asset-generator cancel -q

# Cancel with verbose output
asset-generator cancel -v
```

## Flags

| Flag | Description |
|------|-------------|
| `--all` | Cancel all queued generations (not just current) |
| `-q, --quiet` | Suppress output (errors only) |
| `-v, --verbose` | Show detailed API communication |

## Common Use Cases

### Stop Long-Running Generation
```bash
# Started a Flux generation that's taking too long
asset-generator cancel
```

### Clear All Pending Work
```bash
# Queued too many generations, clear them all
asset-generator cancel --all
```

### Check Before Cancelling
```bash
# See what's running first
asset-generator status
asset-generator cancel
```

### Script Integration
```bash
#!/bin/bash
# Cleanup script
asset-generator cancel --all -q || echo "Cancel failed"
```

## API Endpoints

| Operation | Endpoint | Description |
|-----------|----------|-------------|
| Single | `/API/InterruptGeneration` | Cancel current generation |
| All | `/API/InterruptAll` | Cancel all queued generations |

## Exit Codes

- `0` - Success
- `1` - Error (connection, API, etc.)

## Related Commands

- `asset-generator status` - Check server status
- `asset-generator generate image` - Start generation
- `asset-generator --help` - Full help

## Tips

- Use `--all` to clear the queue completely
- Combine with `status` to verify cancellation
- Use `-q` in scripts for clean automation
- Cancelled generations cannot be resumed

## Error Messages

| Error | Meaning |
|-------|---------|
| "no active generation" | Nothing to cancel (informational) |
| "failed to get session" | Session/connection issue |
| "request failed" | Network/server problem |

## See Full Documentation

See [CANCEL_COMMAND.md](CANCEL_COMMAND.md) for complete details.
