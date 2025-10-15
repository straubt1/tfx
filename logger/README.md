# TFX Logger Package

A structured logging package for the TFX CLI built on Go's standard library `log/slog` package, providing colored, leveled logging for debugging and troubleshooting.

## Features

- **Built on Go's `log/slog`**: Uses the official standard library logging package (Go 1.21+)
- **Multiple Log Levels**: TRACE, DEBUG, INFO, WARN, ERROR, OFF
- **Colored Output**: Color-coded log levels for easy visual scanning
- **Structured Logging**: Support for key-value pairs
- **Thread-Safe**: Safe for concurrent use (guaranteed by slog)
- **Zero Impact When Disabled**: No performance overhead when logging is off
- **Simple API**: Easy to use from anywhere in the codebase
- **No External Dependencies**: Only uses stdlib + existing aurora for colors

## Usage

### Environment Variable

Set the `TFX_LOG` environment variable to control logging:

```bash
# No logging (default)
./tfx org list

# Info level - shows informational messages and above
TFX_LOG=info ./tfx org list

# Debug level - shows debug messages and above
TFX_LOG=debug ./tfx org list

# Trace level - shows everything
TFX_LOG=trace ./tfx org list

# Warn level - shows only warnings and errors
TFX_LOG=warn ./tfx org list

# Error level - shows only errors
TFX_LOG=error ./tfx org list
```

### In Code

```go
import "github.com/straubt1/tfx/logger"

// Simple messages
logger.Info("Fetching organizations")
logger.Debug("Starting operation")
logger.Warn("Deprecated feature used")
logger.Error("Failed to fetch data")
logger.Trace("Detailed trace information")

// Structured logging with key-value pairs
logger.Info("Organizations fetched successfully", "count", len(orgs))
logger.Debug("API response received", "status", resp.StatusCode, "duration", elapsed)
logger.Error("Failed to read organization", "organization", orgName, "error", err)
```

## Log Levels

| Level | Description | Color | Use Case |
|-------|-------------|-------|----------|
| TRACE | Most verbose, follow execution flow | Gray/Faint | Detailed debugging, following execution step-by-step |
| DEBUG | Detailed diagnostic information | Cyan | HTTP requests/responses, data structures, flow control |
| INFO | General informational messages | Green | Operations started/completed, high-level status |
| WARN | Warning messages | Yellow | Deprecated features, non-critical issues |
| ERROR | Error conditions | Red | Errors that don't stop execution |
| OFF | No logging | - | Production/normal use |

## Output Format

```
[HH:MM:SS] [LEVEL] message key=value key2=value2
```

### Examples

```
[16:57:21] [INFO ] Fetching organizations
[16:57:22] [INFO ] Organizations fetched successfully count=4
[16:57:28] [DEBUG] Building output table organizations=4
[16:57:38] [TRACE] Adding organization to table name=terraform-tom email=tstraub@hashicorp.com
[16:58:41] [ERROR] Failed to read organization organization=nonexistent-org error=resource not found
```

## Design Principles

1. **Separate from User Output**: Logging is for developers/debugging, not end-user output
2. **Disabled by Default**: No impact on normal CLI usage
3. **Written to stderr**: Keeps stdout clean for piping/parsing command results
4. **Contextual Information**: Use key-value pairs to provide context
5. **Standard Library**: Built on Go's official `log/slog` package for stability and compatibility

## Implementation

The logger is a thin wrapper around Go's `log/slog` package with a custom `ColorHandler` that:

- Formats output with timestamps and colors
- Uses `slog.Level` directly (no custom level type)
- Adds a TRACE level (slog.LevelDebug - 1) for most verbose logging
- Provides a simple, familiar API similar to popular logging libraries

### Log Level Constants

The package defines two additional level constants:

- `logger.LevelTrace` - More verbose than Debug (slog.LevelDebug - 1)
- `logger.LevelOff` - Disables all logging (slog.Level(100))

All other levels use the standard `slog.Level*` constants directly.

## Integration

The logger is automatically initialized during application startup in `cmd/root.go`. No manual initialization required.

## HTTP Logging

The logger package is separate from the existing HTTP logger (`TFX_LOG_PATH`):

- **`TFX_LOG`** - Controls structured logging to stderr
- **`TFX_LOG_PATH`** - Controls HTTP request/response logging to files

Both can be used simultaneously for maximum debugging capability.
