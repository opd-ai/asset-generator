````markdown
# Asset Generator CLI - SwarmUI API Integration Reference

**Note**: This document describes how the Asset Generator CLI integrates with the SwarmUI API. For complete SwarmUI API documentation, refer to the [official SwarmUI repository](https://github.com/mcmonkeyprojects/SwarmUI).

## Overview

The Asset Generator CLI is a client application that communicates with SwarmUI (or compatible) image generation APIs. This document explains the API integration patterns used by the CLI.

### The Basics

SwarmUI provides a full-capability network API for image generation and model management. The Asset Generator CLI wraps these APIs to provide a convenient command-line interface.

The majority of API calls take the form of `POST` requests sent to `(your server)/API/(route)`, containing JSON formatted inputs and receiving JSON formatted outputs.

### WebSocket Support

The Asset Generator CLI supports WebSocket connections for real-time progress updates during image generation.

**WebSocket Feature Status**: ✅ **IMPLEMENTED**

Use the `--websocket` flag with the `generate image` command to enable real-time progress tracking:

```bash
asset-generator generate image --prompt "your prompt" --websocket
```

WebSocket connections provide:
- Real-time progress updates (actual progress, not simulated)
- Live generation status
- Detailed feedback during long-running generations (e.g., Flux models: 5-10 minutes)
- Automatic fallback to HTTP if WebSocket connection fails

### Authorization

All API routes, with the exception of `GetNewSession`, require a `session_id` input in the JSON. The Asset Generator CLI handles session management automatically.

If the SwarmUI instance is configured to require accounts, set your API key:

```bash
asset-generator config set api-key your-api-key-here
```

### Error Handling

The Asset Generator CLI handles SwarmUI API errors gracefully:
- `invalid_session_id`: Automatically obtains a new session and retries
- Other errors: Displays descriptive error messages with troubleshooting suggestions

## CLI Integration Examples

The Asset Generator CLI provides a simplified interface to SwarmUI's API:

### Example 1: Generate an Image (CLI)

```bash
# Using the Asset Generator CLI
asset-generator generate image --prompt "a cat"
```

This automatically:
1. Obtains a session ID from `/API/GetNewSession`
2. Submits generation request to `/API/GenerateText2Image`
3. Parses the response and displays results

### Example 2: Generate with WebSocket Progress (CLI)

```bash
# Use WebSocket for real-time progress
asset-generator generate image --prompt "a cat" --websocket
```

This connects to `/API/GenerateText2ImageWS` for live updates.

### Example 3: Download and Save Images (CLI)

```bash
# Download generated images to local disk
asset-generator generate image \
  --prompt "a cat" \
  --save-images \
  --output-dir ./my-images
```

## Implementation Details

### API Client Features

The Asset Generator CLI's API client (`pkg/client/client.go`) implements:

- **HTTP Generation**: `GenerateImage()` - Standard REST API calls
- **WebSocket Generation**: `GenerateImageWS()` - Real-time progress updates
- **Session Management**: Automatic session creation, caching, and renewal
- **Error Handling**: Automatic retry on session expiration
- **Context Support**: Cancellation for graceful shutdown
- **Progress Tracking**: Simulated progress (HTTP) or real-time progress (WebSocket)

### Session Management

Sessions are managed automatically:
- First API call triggers session creation via `/API/GetNewSession`
- Session ID is cached and reused for subsequent calls
- Expired sessions are automatically renewed
- Manual session management is not required

## Direct API Access (Advanced)
