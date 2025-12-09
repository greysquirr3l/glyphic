# Modern Delve Debugging in Go with VS Code

## Overview

This guide covers setting up and using modern Delve debugging (dlv-dap mode) in VS Code for Go applications. The modern Delve implementation uses the Debug Adapter Protocol (DAP) natively, providing better performance and more features than the legacy debug adapter.

## Prerequisites

- **Go 1.18+** (recommended: Go 1.21+)
- **VS Code** with the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
- **Delve debugger** (dlv) - latest version

## Installation and Setup

### 1. Install/Update Delve

#### Automatic Installation (Recommended)
1. Open VS Code Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`)
2. Run `Go: Install/Update Tools`
3. Select `dlv` from the tool list
4. Click "OK" to install

#### Manual Installation
For Go 1.16+:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

For older Go versions:
```bash
cd $(mktemp -d)
GO111MODULE=on go get github.com/go-delve/delve/cmd/dlv@latest
```

#### Verify Installation
```bash
dlv version
# Should show version 1.20.0+ for modern DAP support
```

### 2. Configure Modern Debug Adapter

Add to your VS Code settings (`settings.json`):
```json
{
  "go.delveConfig": {
    "debugAdapter": "dlv-dap",
    "showGlobalVariables": false,
    "substitutePath": []
  },
  "go.toolsManagement.autoUpdate": true
}
```

## Debug Configuration

### 1. Basic Launch Configuration

Create `.vscode/launch.json` in your project root:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Package",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${fileDirname}",
      "debugAdapter": "dlv-dap"
    }
  ]
}
```

### 2. Advanced Configuration Options

#### Debug Main Package
```json
{
  "name": "Debug Main",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}/cmd/server",
  "args": ["-config", "config.yaml"],
  "env": {
    "GO_ENV": "development",
    "DB_HOST": "localhost"
  },
  "cwd": "${workspaceFolder}",
  "debugAdapter": "dlv-dap",
  "showLog": true,
  "trace": "verbose"
}
```

#### Debug Tests
```json
{
  "name": "Debug Test",
  "type": "go",
  "request": "launch",
  "mode": "test",
  "program": "${workspaceFolder}/internal/domain/user",
  "args": ["-test.v", "-test.run", "TestUserService"],
  "debugAdapter": "dlv-dap"
}
```

#### Debug Pre-built Binary
```json
{
  "name": "Debug Binary",
  "type": "go",
  "request": "launch",
  "mode": "exec",
  "program": "${workspaceFolder}/bin/server",
  "args": ["--port", "8080"],
  "debugAdapter": "dlv-dap"
}
```

#### Attach to Running Process
```json
{
  "name": "Attach to Process",
  "type": "go",
  "request": "attach",
  "mode": "local",
  "processId": "${command:pickGoProcess}",
  "debugAdapter": "dlv-dap"
}
```

### 3. Configuration Attributes Reference

| Attribute | Description | Example |
|-----------|-------------|---------|
| `name` | Configuration name in dropdown | `"Launch Server"` |
| `type` | Must be `"go"` | `"go"` |
| `request` | `"launch"` or `"attach"` | `"launch"` |
| `mode` | `"auto"`, `"debug"`, `"test"`, `"exec"` | `"debug"` |
| `program` | Path to main package or binary | `"${workspaceFolder}/cmd/server"` |
| `args` | Command line arguments | `["-port", "8080"]` |
| `env` | Environment variables | `{"DEBUG": "true"}` |
| `cwd` | Working directory | `"${workspaceFolder}"` |
| `debugAdapter` | Use `"dlv-dap"` for modern | `"dlv-dap"` |
| `stopOnEntry` | Pause at program start | `false` |
| `console` | Terminal type | `"integratedTerminal"` |
| `showLog` | Show Delve logs | `true` |
| `trace` | Log verbosity | `"verbose"` |

## Debugging Features

### 1. Breakpoints

#### Setting Breakpoints
- **Line Breakpoints**: Click editor gutter or press `F9`
- **Conditional Breakpoints**: Right-click gutter → "Add Conditional Breakpoint"
- **Function Breakpoints**: In BREAKPOINTS panel, click `+` and enter function name
- **Logpoints**: Right-click gutter → "Add Logpoint"

#### Conditional Breakpoint Examples
```go
// Expression condition
x > 10 && y < 5

// Hit count condition
% 5  // Every 5th hit

// Wait for another breakpoint
// Configure in UI - triggered when another breakpoint hits
```

#### Function Breakpoint Examples
```
main.main
github.com/yourorg/ragnar/internal/domain/user.(*UserService).CreateUser
(*UserService).CreateUser:15
```

#### Logpoint Examples
```
User ID: {userID}, Status: {user.Status}
Processing {len(items)} items
Function called with args: {args}
```

### 2. Debug Actions

| Action | Shortcut | Description |
|--------|----------|-------------|
| Continue | `F5` | Resume execution |
| Step Over | `F10` | Execute next line |
| Step Into | `F11` | Enter function calls |
| Step Out | `Shift+F11` | Exit current function |
| Restart | `Ctrl+Shift+F5` | Restart debug session |
| Stop | `Shift+F5` | Stop debugging |

### 3. Data Inspection

#### Variables Panel
- View local variables and function arguments
- Expand complex types (structs, slices, maps)
- Modify simple values with right-click → "Set Value"
- Copy values or expressions

#### Watch Panel
Add expressions to monitor:
```go
user.ID
len(users)
config.Database.Host
time.Now().Format("15:04:05")
```

#### Debug Console
Evaluate expressions during debugging:
```go
// Variable inspection
p user
p user.Email
p len(tickets)

// Function calls (experimental)
call fmt.Printf("Debug: %v\n", user)

// Delve commands
dlv help
dlv config -list
```

### 4. Call Stack and Goroutines

- **Call Stack**: View function call hierarchy
- **Goroutines**: Switch between goroutines
- **Stack Frames**: Navigate through call stack
- **Runtime Frames**: View or hide Go runtime frames

## Advanced Debugging Scenarios

### 1. Remote Debugging

#### Server Setup
```bash
# Start headless Delve server
dlv debug ./cmd/server --headless --listen=:2345 --api-version=2
```

#### Client Configuration
```json
{
  "name": "Connect to Remote",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "remotePath": "/path/to/remote/workspace",
  "port": 2345,
  "host": "remote-server.com",
  "debugAdapter": "dlv-dap",
  "substitutePath": [
    {
      "from": "${workspaceFolder}",
      "to": "/path/to/remote/workspace"
    }
  ]
}
```

### 2. Container Debugging

#### Dockerfile for Debugging
```dockerfile
FROM golang:1.21-alpine

# Install delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app
COPY . .

# Build with debug symbols
RUN go build -gcflags="all=-N -l" -o server ./cmd/server

# Expose debug port
EXPOSE 40000

# Start with delve
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "./server"]
```

#### Docker Compose
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
      - "40000:40000"  # Debug port
    security_opt:
      - "apparmor:unconfined"
    cap_add:
      - SYS_PTRACE
```

#### VS Code Configuration
```json
{
  "name": "Debug Container",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "port": 40000,
  "host": "localhost",
  "debugAdapter": "dlv-dap"
}
```

### 3. Testing with Debugging

#### Debug Specific Test
```json
{
  "name": "Debug Current Test",
  "type": "go",
  "request": "launch",
  "mode": "test",
  "program": "${fileDirname}",
  "args": [
    "-test.run",
    "^${input:testName}$",
    "-test.v"
  ],
  "debugAdapter": "dlv-dap"
}
```

#### Debug Test with Coverage
```json
{
  "name": "Debug Test with Coverage",
  "type": "go",
  "request": "launch",
  "mode": "test",
  "program": "${fileDirname}",
  "args": [
    "-test.coverprofile=coverage.out",
    "-test.v"
  ],
  "debugAdapter": "dlv-dap"
}
```

### 4. Debugging as Root (Linux/macOS)

```json
{
  "name": "Debug as Root",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}/cmd/server",
  "asRoot": true,
  "console": "integratedTerminal",
  "debugAdapter": "dlv-dap"
}
```

### 5. Multi-Session Debugging

Debug multiple processes simultaneously:

```json
{
  "name": "Launch Client",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}/cmd/client",
  "debugAdapter": "dlv-dap"
},
{
  "name": "Launch Server", 
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}/cmd/server",
  "debugAdapter": "dlv-dap"
}
```

## Best Practices

### 1. Performance Optimization

```json
{
  "name": "Optimized Debug",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}",
  "debugAdapter": "dlv-dap",
  "trace": "error",  // Reduce logging
  "showLog": false,  // Disable logs
  "hideSystemGoroutines": true,  // Hide Go runtime
  "showGlobalVariables": false   // Hide globals
}
```

### 2. Environment-Specific Configs

```json
{
  "name": "Debug Development",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}/cmd/server",
  "envFile": "${workspaceFolder}/.env.development",
  "debugAdapter": "dlv-dap"
},
{
  "name": "Debug Production Simulation",
  "type": "go",
  "request": "launch",
  "mode": "debug", 
  "program": "${workspaceFolder}/cmd/server",
  "envFile": "${workspaceFolder}/.env.production",
  "debugAdapter": "dlv-dap"
}
```

### 3. Debugging Build Process

For debugging issues with build flags:

```json
{
  "name": "Debug with Build Flags",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}",
  "buildFlags": [
    "-tags=debug",
    "-ldflags=-X main.version=dev"
  ],
  "debugAdapter": "dlv-dap"
}
```

## Troubleshooting

### 1. Common Issues

#### Issue: Breakpoints not hitting
**Solution**: Check build flags and ensure debug symbols are included:
```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o myapp
```

#### Issue: "dlv" not found
**Solution**: Install Delve or specify path:
```json
{
  "go.alternateTools": {
    "dlv": "/path/to/your/dlv"
  }
}
```

#### Issue: Variables showing "optimized out"
**Solution**: Disable optimizations:
```json
{
  "buildFlags": ["-gcflags=all=-N -l"]
}
```

### 2. Logging and Diagnostics

Enable verbose logging for troubleshooting:
```json
{
  "name": "Debug with Logs",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}",
  "debugAdapter": "dlv-dap",
  "trace": "verbose",
  "showLog": true,
  "logOutput": "debugger,dap,rpc"
}
```

### 3. Performance Issues

If debugging is slow:
1. Reduce variable inspection depth
2. Hide system goroutines: `"hideSystemGoroutines": true`
3. Disable global variables: `"showGlobalVariables": false`
4. Use conditional breakpoints sparingly

## Integration with Ragnar Project

### 1. Ragnar-Specific Configuration

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Ragnar Server",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/server",
      "args": [
        "-config", "${workspaceFolder}/config/dev.yaml"
      ],
      "env": {
        "RAGNAR_ENV": "development",
        "RAGNAR_LOG_LEVEL": "debug"
      },
      "envFile": "${workspaceFolder}/.env.development",
      "cwd": "${workspaceFolder}",
      "debugAdapter": "dlv-dap",
      "console": "integratedTerminal"
    },
    {
      "name": "Debug Ragnar Worker",
      "type": "go", 
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/worker",
      "debugAdapter": "dlv-dap"
    },
    {
      "name": "Debug Migration",
      "type": "go",
      "request": "launch", 
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/migrate",
      "args": ["up"],
      "debugAdapter": "dlv-dap"
    },
    {
      "name": "Debug Unit Tests",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${fileDirname}",
      "args": ["-test.v"],
      "debugAdapter": "dlv-dap"
    },
    {
      "name": "Attach to Ragnar Process",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": "${command:pickGoProcess}",
      "debugAdapter": "dlv-dap"
    }
  ]
}
```

### 2. Environment Files

Create `.env.development`:
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ragnar_dev
DB_USER=ragnar
DB_PASSWORD=password

# Email
SMTP_HOST=localhost
SMTP_PORT=1025

# Security
JWT_SECRET=dev-secret-key

# Debugging
RAGNAR_DEBUG=true
RAGNAR_LOG_LEVEL=debug
```

### 3. VS Code Tasks Integration

Add to `.vscode/tasks.json`:
```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Build Debug Binary",
      "type": "shell",
      "command": "go",
      "args": [
        "build",
        "-gcflags=all=-N -l",
        "-o", "./bin/server-debug",
        "./cmd/server"
      ],
      "group": "build",
      "presentation": {
        "echo": true,
        "reveal": "silent",
        "focus": false,
        "panel": "shared"
      }
    }
  ]
}
```

## Additional Resources

- [Delve Documentation](https://github.com/go-delve/delve/tree/master/Documentation)
- [VS Code Go Extension Wiki](https://github.com/golang/vscode-go/wiki/debugging)
- [Debug Adapter Protocol](https://microsoft.github.io/debug-adapter-protocol/)
- [Go Debugging Best Practices](https://go.dev/doc/debugging)

## Summary

Modern Delve debugging with dlv-dap provides:
- ✅ Better performance than legacy adapter
- ✅ Rich debugging features (conditional breakpoints, logpoints, etc.)
- ✅ Excellent VS Code integration
- ✅ Support for complex debugging scenarios
- ✅ Active development and support

This setup enables efficient debugging of Go applications with modern tooling and best practices.