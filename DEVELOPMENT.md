# SwarmUI CLI - Developer Documentation

## Architecture Overview

The SwarmUI CLI is built using clean architecture principles with clear separation of concerns:

### Directory Structure

```
asset-generator/
├── cmd/                    # Command definitions (CLI layer)
│   ├── root.go            # Root command with global flags
│   ├── generate.go        # Image generation commands
│   ├── models.go          # Model management commands
│   └── config.go          # Configuration management
├── pkg/                   # Public reusable packages
│   ├── client/            # SwarmUI API client
│   │   ├── client.go
│   │   └── client_test.go
│   └── output/            # Output formatting
│       ├── formatter.go
│       └── formatter_test.go
├── internal/              # Private application code
│   └── config/            # Configuration validation
│       ├── validate.go
│       └── validate_test.go
├── main.go               # Application entry point
├── go.mod                # Go module definition
├── Makefile              # Build automation
└── README.md             # User documentation
```

## Component Details

### 1. Command Layer (`cmd/`)

The command layer handles CLI interactions using Cobra:

#### `root.go`
- Sets up global flags (api-url, api-key, format, output, etc.)
- Initializes configuration from multiple sources
- Creates SwarmUI client instance
- Implements flag precedence: CLI flags > env vars > config file > defaults

#### `generate.go`
- Implements `asset-generator generate image` command
- Handles signal interrupts for graceful shutdown
- Tracks generation progress
- Supports batch operations

#### `models.go`
- Lists available models
- Retrieves model details
- Formats output in multiple formats

#### `config.go`
- Initializes config file
- Views current configuration
- Sets/gets configuration values
- Masks sensitive values in output

### 2. API Client (`pkg/client/`)

Reusable SwarmUI API client library:

#### Key Types

```go
type SwarmClient struct {
    config     *Config
    httpClient *http.Client
    wsConn     *websocket.Conn
    mu         sync.RWMutex
    sessions   map[string]*GenerationSession
}

type GenerationRequest struct {
    Prompt     string
    Model      string
    Parameters map[string]interface{}
    SessionID  string
}

type GenerationResult struct {
    SessionID  string
    ImagePaths []string
    Metadata   map[string]interface{}
    Status     string
    CreatedAt  time.Time
}
```

#### Features
- HTTP client for REST API calls
- WebSocket support (foundation for future real-time features)
- Session management for tracking generations
- Context support for cancellation
- Error handling with descriptive messages

### 3. Output Formatting (`pkg/output/`)

Flexible output formatting system:

#### Supported Formats
- **Table**: Tab-separated human-readable output
- **JSON**: Structured JSON with indentation
- **YAML**: YAML format with comments

#### Features
- Automatic format detection
- Timestamp injection for file outputs
- Type-aware formatting
- Error handling

### 4. Configuration (`internal/config/`)

Configuration validation and management:

#### Validation Rules
- URL must include scheme (http/https/ws/wss)
- Format must be table, json, or yaml
- All values sanitized before use

#### Configuration Sources (in order of precedence)
1. Command-line flags
2. Environment variables (SWARMUI_*)
3. Configuration file (~/.asset-generator/config.yaml)
4. Default values

## API Integration

### SwarmUI API Endpoints

The client currently implements:

1. **Generate Text2Image**: `POST /API/GenerateText2Image`
   - Generates images from text prompts
   - Supports various parameters (steps, size, seed, etc.)

2. **List Models**: `GET /API/ListModels`
   - Returns available models

3. **Get Model**: `GET /API/GetModel?name={name}`
   - Returns model details

### Adding New Endpoints

To add a new endpoint:

1. Add method to `SwarmClient` in `pkg/client/client.go`:
```go
func (c *SwarmClient) NewEndpoint(params) (*Result, error) {
    endpoint := fmt.Sprintf("%s/API/NewEndpoint", c.config.BaseURL)
    // Implementation
}
```

2. Create command in `cmd/`:
```go
var newCmd = &cobra.Command{
    Use: "new",
    RunE: runNewCommand,
}
```

3. Add tests in `pkg/client/client_test.go`

## Testing

### Test Coverage

Current coverage:
- `internal/config`: 95.0%
- `pkg/client`: 54.7%
- `pkg/output`: 60.6%

### Running Tests

```bash
# All tests
make test

# With coverage report
make coverage

# Specific package
go test -v ./pkg/client/...

# With race detector
go test -race ./...
```

### Test Structure

Tests use table-driven approach:

```go
tests := []struct {
    name    string
    input   InputType
    want    OutputType
    wantErr bool
}{
    // test cases
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test logic
    })
}
```

## Building and Distribution

### Building

```bash
# Local build
make build

# Install system-wide
make install

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o asset-generator-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o asset-generator-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o asset-generator-windows-amd64.exe
```

### Binary Size Optimization

The Makefile uses these flags:
```bash
-ldflags "-s -w"
```
- `-s`: Omit symbol table
- `-w`: Omit DWARF debug info

This reduces binary size by ~30%.

## Error Handling

### Guidelines

1. **User-Facing Errors**: Descriptive, actionable messages
   ```go
   return fmt.Errorf("failed to connect to SwarmUI at %s: %w (is the service running?)", url, err)
   ```

2. **Error Wrapping**: Use `%w` for error chains
   ```go
   return fmt.Errorf("failed to parse response: %w", err)
   ```

3. **Cobra Error Handling**: Use `cobra.CheckErr()` for fatal errors
   ```go
   cobra.CheckErr(err)
   ```

4. **Stderr vs Stdout**: Errors and progress to stderr, data to stdout
   ```go
   fmt.Fprintln(os.Stderr, "Processing...")
   fmt.Println(data) // Goes to stdout
   ```

## Configuration Management

### Config File Location

Default: `~/.asset-generator/config.yaml`

Can be overridden with `--config` flag.

### Config Structure

```yaml
# API Configuration
api-url: http://localhost:7801
api-key: your-api-key-here

# Output Configuration
format: table
verbose: false
quiet: false

# Default Generation Parameters
generate:
  model: stable-diffusion-xl
  steps: 20
  width: 512
  height: 512
  cfg-scale: 7.5
  sampler: euler_a
```

### Environment Variables

All config values can be set via environment:
```bash
SWARMUI_API_URL=http://localhost:7801
SWARMUI_API_KEY=secret
SWARMUI_FORMAT=json
```

## Best Practices

### Command Design

1. **Consistent Naming**: Use verb-noun pattern (`generate image`, `list models`)
2. **Short Flags**: Single letter for common options (`-p` for prompt, `-o` for output)
3. **Sensible Defaults**: Most operations work without flags
4. **Help Text**: Include examples in every command

### Code Style

1. **Error Messages**: Start with lowercase, no trailing punctuation
2. **Comments**: Godoc comments for all exported symbols
3. **Testing**: Minimum 60% coverage, critical paths 100%
4. **Formatting**: Run `make fmt` before committing

### Adding Features

1. Start with tests (TDD approach)
2. Implement in `pkg/` if reusable, `internal/` if private
3. Add command in `cmd/`
4. Update documentation
5. Add examples to help text

## Troubleshooting

### Common Issues

**Issue**: "connection refused"
- Check SwarmUI is running
- Verify API URL in config
- Check firewall settings

**Issue**: "invalid format"
- Format must be: table, json, or yaml
- Check config file syntax

**Issue**: "config file not found"
- Run `asset-generator config init` first
- Check `--config` flag path

### Debug Mode

Enable verbose output:
```bash
asset-generator --verbose generate image --prompt "test"
```

This shows:
- Config file location
- API requests
- Response codes
- Timing information

## Performance

### Benchmarks

- Local operations: <10ms
- API calls: Network dependent
- Large outputs: Uses streaming where possible

### Optimization Tips

1. Use `--quiet` for scripting
2. Batch operations when possible
3. Cache model lists locally
4. Use specific model names (avoid lookups)

## Future Enhancements

### Planned Features

1. **WebSocket Support**: Real-time progress updates
2. **Batch Files**: Read prompts from file
3. **Image Input**: img2img operations
4. **Plugins**: Extension system for custom commands
5. **Shell Completions**: Auto-generated completion scripts

### API Evolution

The client is designed for extension:
- Add methods to `SwarmClient`
- Implement in separate files for organization
- Keep backward compatibility

## Contributing

See main README.md for contribution guidelines.

### Code Review Checklist

- [ ] Tests pass (`make test`)
- [ ] Lint passes (`make lint`)
- [ ] Code formatted (`make fmt`)
- [ ] Documentation updated
- [ ] Examples added
- [ ] Error messages are clear
- [ ] Backward compatible

## Resources

- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [SwarmUI API Docs](https://github.com/mcmonkeyprojects/SwarmUI)
- [Go Best Practices](https://go.dev/doc/effective_go)
