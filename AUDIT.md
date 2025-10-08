# SwarmUI CLI Functional Audit Report

**Date:** October 8, 2025  
**Auditor:** GitHub Copilot  
**Repository:** opd-ai/asset-generator  
**Branch:** main  

## Executive Summary

This comprehensive functional audit examined the Go codebase to identify discrepancies between documented functionality in README.md and actual implementation. The audit followed a systematic dependency-based analysis approach, examining files in order from Level 0 (no internal imports) through Level 2 (main.go).

## Audit Summary

```
Total Issues Found: 9
- CRITICAL BUG: 2 
- FUNCTIONAL MISMATCH: 3
- MISSING FEATURE: 3  
- EDGE CASE BUG: 1
- PERFORMANCE ISSUE: 0

Dependency Analysis: 
- Level 0 files: 2 (internal/config/validate.go, pkg/output/formatter.go)
- Level 1 files: 5 (pkg/client/client.go, cmd/*.go files)  
- Level 2 files: 1 (main.go)

Overall Assessment: Several critical bugs and missing features that prevent the CLI from working as documented.
```

## Methodology

The audit was conducted using the following systematic approach:

1. **Documentation Analysis**: Comprehensive review of README.md to extract all functional requirements
2. **Dependency Mapping**: Analysis of import dependencies to categorize files by dependency levels
3. **Level-based Code Review**: Systematic examination of files in dependency order (0→1→2)
4. **Feature Tracing**: Verification of each documented feature's implementation
5. **Edge Case Analysis**: Testing of boundary conditions and error scenarios

## Critical Issues

### 1. CRITICAL BUG: Incorrect SwarmUI API Endpoint Usage
**File:** `pkg/client/client.go:89-95`  
**Severity:** High  
**Impact:** Complete failure of image generation functionality

The client is using incorrect API endpoints that don't match the actual SwarmUI API specification. The code uses `/API/GenerateText2Image` and `/API/GetModel?name=` endpoints, but SwarmUI uses different endpoint patterns.

**Expected Behavior:** Should use proper SwarmUI API endpoints like `/API/GenerateText2ImageWS` for WebSocket generation or the correct REST endpoints

**Actual Behavior:** Makes HTTP requests to non-existent endpoints, causing all generation requests to fail

**Reproduction:** Run `swarmui generate image --prompt "test"` with any SwarmUI instance

**Code Reference:**
```go
endpoint := fmt.Sprintf("%s/API/GenerateText2Image", c.config.BaseURL)
// This endpoint doesn't exist in SwarmUI API
```

### 2. CRITICAL BUG: Hardcoded Field Name Mismatch in Generation Request
**File:** `pkg/client/client.go:97-105`  
**Severity:** High  
**Impact:** Even if endpoints were correct, generation would fail due to invalid request format

The generation request body uses incorrect field names. The code sets `body["images"] = req.Parameters["batch_size"]` but SwarmUI expects different parameter names.

**Expected Behavior:** Should use proper SwarmUI parameter names and structure according to API specification

**Actual Behavior:** Sends malformed requests that SwarmUI cannot process

**Reproduction:** Any generation request will fail at the API level

**Code Reference:**
```go
body := map[string]interface{}{
    "prompt": req.Prompt,
    "images": req.Parameters["batch_size"], // Wrong field name
}
```

## Functional Mismatches

### 3. FUNCTIONAL MISMATCH: WebSocket Support Not Implemented Despite Documentation
**File:** `pkg/client/client.go:26-28`  
**Severity:** High  
**Impact:** Users cannot see generation progress, contradicting documented features

The README and copilot instructions emphasize WebSocket integration for SwarmUI, and the client struct has a `wsConn` field, but no WebSocket functionality is actually implemented.

**Expected Behavior:** Should use WebSocket connections for real-time generation progress as documented in copilot instructions

**Actual Behavior:** Only implements HTTP requests, missing real-time progress capabilities

**Reproduction:** Check any generation - no progress feedback is provided

**Code Reference:**
```go
wsConn     *websocket.Conn  // Field exists but never used
// All generation uses HTTP instead of WebSocket
```

### 4. FUNCTIONAL MISMATCH: Output Format Table Implementation Incomplete
**File:** `pkg/output/formatter.go:77-95`  
**Severity:** Medium  
**Impact:** Poor user experience, output doesn't match professional CLI tool expectations

The table formatter has crude implementation that doesn't match the quality expected from the documentation screenshots and examples.

**Expected Behavior:** Should provide well-formatted, aligned tables with proper headers and spacing

**Actual Behavior:** Uses simple tab-separated values with basic formatting

**Reproduction:** Run any command with `--format table` (default) to see crude table output

**Code Reference:**
```go
// Write header
buf.WriteString(strings.Join(headers, "\t"))
// Uses tabs instead of proper column alignment
```

### 5. FUNCTIONAL MISMATCH: Error Messages Don't Match SwarmUI API Format
**File:** `pkg/client/client.go:142-145`  
**Severity:** Medium  
**Impact:** Users get unhelpful error messages when API calls fail

The error handling assumes generic HTTP error responses but SwarmUI likely returns structured error messages.

**Expected Behavior:** Should parse SwarmUI-specific error format and provide meaningful error messages

**Actual Behavior:** Returns raw HTTP response body as error message

**Reproduction:** Try any API call with wrong credentials or invalid parameters

**Code Reference:**
```go
if resp.StatusCode != http.StatusOK {
    bodyBytes, _ := io.ReadAll(resp.Body)
    return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
}
```

## Missing Features

### 6. MISSING FEATURE: Config Init Command Not Implemented
**File:** `cmd/config.go:60-63`  
**Severity:** Medium  
**Impact:** Users cannot initialize configuration as documented in Quick Start guide

The README documents `swarmui config init` command for initializing configuration, but the implementation is missing the actual functionality.

**Expected Behavior:** Should create default config file at ~/.swarmui/config.yaml with default values

**Actual Behavior:** Function `runConfigInit` is declared but not implemented

**Reproduction:** Run `swarmui config init` - command exists but does nothing

**Code Reference:**
```go
var configInitCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize a new config file",
    Long:  `Create a new configuration file with default values.`,
    RunE:  runConfigInit, // Function not implemented
}
```

### 7. MISSING FEATURE: Batch Generation Not Properly Implemented
**File:** `cmd/generate.go:100-110`  
**Severity:** Medium  
**Impact:** Users cannot generate multiple images in one command as documented

The documentation shows batch generation support with `--batch` flag, but the implementation doesn't handle multiple image generation correctly.

**Expected Behavior:** Should generate multiple images when batch size > 1 and return array of results

**Actual Behavior:** Sends batch_size parameter but only handles single result structure

**Reproduction:** Run `swarmui generate image --prompt "test" --batch 4` - only one image would be attempted

**Code Reference:**
```go
req.Parameters["batch_size"] = generateBatchSize
// But response parsing only handles single image result
```

### 8. MISSING FEATURE: Model Validation Not Implemented
**File:** `cmd/generate.go:108-115`  
**Severity:** Medium  
**Impact:** Poor user experience - users get cryptic API errors instead of helpful validation messages

The code accepts any model name without validation against available models, despite having a ListModels function.

**Expected Behavior:** Should validate model names against available models before attempting generation

**Actual Behavior:** Accepts any model string and passes it through, leading to potential API errors

**Reproduction:** Run `swarmui generate image --prompt "test" --model "nonexistent-model"`

**Code Reference:**
```go
if generateModel != "" {
    req.Model = generateModel  // No validation performed
}
```

## Edge Cases and Minor Issues

### 9. EDGE CASE BUG: Race Condition in Session Management
**File:** `pkg/client/client.go:81-86, 163-166`  
**Severity:** Low  
**Impact:** Memory leak potential if sessions accumulate over time

The session management has potential race condition between session creation and cleanup, though currently mitigated by simple HTTP approach.

**Expected Behavior:** Should have proper synchronization for concurrent session access

**Actual Behavior:** Session map access is protected by mutex but session cleanup is not implemented

**Reproduction:** Run many generation commands in quick succession

**Code Reference:**
```go
c.mu.Lock()
c.sessions[sessionID] = session  // Added but never cleaned up
c.mu.Unlock()
```

## Architectural Assessment

### Positive Aspects
1. **Clean Architecture**: The project follows Go best practices with clear separation of concerns
2. **Proper Dependencies**: Uses established libraries (Cobra, Viper) as recommended
3. **Good Testing Structure**: Test files are present and follow Go conventions
4. **Configuration Management**: Proper hierarchy of configuration sources

### Areas for Improvement
1. **API Integration**: Complete mismatch with actual SwarmUI API specification
2. **Error Handling**: Generic error handling doesn't provide SwarmUI-specific context
3. **Progress Feedback**: Missing real-time progress despite WebSocket support claims
4. **Input Validation**: Lacks validation for critical parameters like model names

## Recommendations

### High Priority (Critical)
1. **Fix API Endpoints**: Research and implement correct SwarmUI API endpoints
2. **Implement WebSocket Support**: Add real-time progress feedback as documented
3. **Correct Request Format**: Use proper SwarmUI parameter names and structure

### Medium Priority
1. **Implement Config Init**: Complete the missing configuration initialization
2. **Add Model Validation**: Validate model names against available models
3. **Fix Batch Generation**: Properly handle multiple image generation
4. **Improve Error Messages**: Parse and display SwarmUI-specific errors

### Low Priority
1. **Enhance Table Formatting**: Implement proper column alignment and spacing
2. **Add Session Cleanup**: Implement proper session lifecycle management
3. **Add Integration Tests**: Test against actual SwarmUI instances

## Testing Recommendations

To verify fixes and prevent regressions:

1. **Integration Tests**: Test against real SwarmUI instances
2. **API Compatibility Tests**: Verify correct endpoint usage
3. **Progress Feedback Tests**: Validate WebSocket implementation
4. **Error Scenario Tests**: Test various failure modes
5. **Configuration Tests**: Verify config init and validation

## Conclusion

While the codebase demonstrates good Go programming practices and architectural decisions, it has significant functionality gaps that prevent it from working as documented. The most critical issues are the incorrect API integration and missing WebSocket implementation, which would prevent the core functionality from working with actual SwarmUI instances.

The issues identified represent a substantial amount of work to bring the implementation in line with the documentation, but the clean architecture provides a solid foundation for these improvements.