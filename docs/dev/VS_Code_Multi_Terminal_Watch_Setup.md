# VS Code Multi-Terminal Watch Setup Guide

This guide explains how to configure VS Code to automatically start multiple watch
processes in parallel terminals when opening a workspace. This pattern is perfect
for complex development workflows that require multiple build processes, servers,
or monitoring tasks running simultaneously.

## How It Works

The system uses VS Code's task configuration (`.vscode/tasks.json`) to orchestrate multiple background processes:

1. **Auto-start on folder open** - Tasks launch automatically when you open the workspace
2. **Parallel execution** - Multiple tasks run simultaneously in separate terminals
3. **Grouped terminals** - All related terminals are organized together
4. **Background processes** - Tasks run continuously without blocking the UI
5. **No focus stealing** - Terminals don't interrupt your workflow

## Configuration Architecture

### 1. Main Watch Task (Entry Point)

```json
{
    "label": "watch",
    "dependsOn": [
        "ensure-deps",           // Optional: Install dependencies first
        "start-watch-tasks"      // Launch all parallel tasks
    ],
    "dependsOrder": "sequence",  // Run dependencies in order, then parallel tasks
    "group": {
        "kind": "build",
        "isDefault": true        // Makes this the default build task
    },
    "runOptions": {
        "runOn": "folderOpen"    // 🎯 Auto-start when workspace opens
    }
}
```

### 2. Parallel Task Orchestrator

```json
{
    "label": "start-watch-tasks",
    "dependsOn": [
        "watch:frontend",
        "watch:backend", 
        "watch:tests",
        "watch:docs"
    ],
    "dependsOrder": "parallel",  // 🎯 Run all tasks simultaneously
    "presentation": {
        "reveal": "never"        // Don't show this orchestrator task
    }
}
```

### 3. Individual Watch Tasks

```json
{
    "label": "watch:frontend",
    "type": "npm",
    "script": "watch:frontend",
    "isBackground": true,        // 🎯 Keep running indefinitely
    "problemMatcher": "$tsc-watch",
    "presentation": {
        "group": "watch",        // 🎯 Group terminals together
        "reveal": "never"        // 🎯 Don't steal focus on start
    }
}
```

## Complete Example: Full-Stack Development

Here's a complete `.vscode/tasks.json` for a typical full-stack project:

```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "watch",
            "dependsOn": [
                "ensure-deps",
                "start-watch-tasks"
            ],
            "dependsOrder": "sequence",
            "group": {
                "kind": "build",
                "isDefault": true
            },
            "runOptions": {
                "runOn": "folderOpen"
            },
            "presentation": {
                "reveal": "never"
            }
        },
        {
            "label": "ensure-deps",
            "type": "shell",
            "command": "test -d node_modules || npm ci",
            "group": "build",
            "isBackground": true,
            "presentation": {
                "reveal": "never",
                "close": true
            },
            "problemMatcher": []
        },
        {
            "label": "start-watch-tasks",
            "dependsOn": [
                "watch:typescript",
                "watch:webpack",
                "watch:api-server",
                "watch:tests"
            ],
            "dependsOrder": "parallel",
            "presentation": {
                "reveal": "never"
            }
        },
        {
            "label": "watch:typescript",
            "type": "npm",
            "script": "watch:tsc",
            "group": "build",
            "isBackground": true,
            "problemMatcher": "$tsc-watch",
            "presentation": {
                "group": "watch",
                "reveal": "never",
                "panel": "dedicated"
            }
        },
        {
            "label": "watch:webpack",
            "type": "npm", 
            "script": "watch:webpack",
            "group": "build",
            "isBackground": true,
            "problemMatcher": "$webpack-watch",
            "presentation": {
                "group": "watch",
                "reveal": "never",
                "panel": "dedicated"
            }
        },
        {
            "label": "watch:api-server",
            "type": "npm",
            "script": "dev:api",
            "group": "build", 
            "isBackground": true,
            "problemMatcher": {
                "pattern": {
                    "regexp": "^Server running on port (\\d+)$",
                    "file": 1
                },
                "background": {
                    "activeOnStart": true,
                    "beginsPattern": "^Starting server...$",
                    "endsPattern": "^Server running on port \\d+$"
                }
            },
            "presentation": {
                "group": "watch",
                "reveal": "never",
                "panel": "dedicated"
            }
        },
        {
            "label": "watch:tests",
            "type": "npm",
            "script": "test:watch",
            "group": "build",
            "isBackground": true,
            "presentation": {
                "group": "watch", 
                "reveal": "never",
                "panel": "dedicated"
            }
        }
    ]
}
```

## Common Patterns and Use Cases

### TypeScript + React + Express API

```json
"dependsOn": [
    "watch:tsc-backend",     // TypeScript compilation for backend
    "watch:react-dev",       // React development server  
    "watch:api-server",      // Express API server
    "watch:tailwind"         // Tailwind CSS compilation
]
```

### Microservices Development

```json
"dependsOn": [
    "watch:auth-service",    // Authentication microservice
    "watch:user-service",    // User management service
    "watch:payment-service", // Payment processing service
    "watch:frontend"         // Frontend application
]
```

### Documentation & Testing

```json
"dependsOn": [
    "watch:docs-server",     // Documentation site (Docusaurus, etc.)
    "watch:unit-tests",      // Unit test runner
    "watch:e2e-tests",       // End-to-end test suite
    "watch:storybook"        // Component documentation
]
```

### Monorepo Workspace

```json
"dependsOn": [
    "watch:shared-lib",      // Shared library compilation
    "watch:web-app",         // Web application
    "watch:mobile-app",      // Mobile app bundler
    "watch:admin-panel"      // Admin interface
]
```

## Key Configuration Properties

### Task Properties

| Property           | Purpose                              | Example                          |
| ------------------ | ------------------------------------ | -------------------------------- |
| `dependsOn`        | Tasks that must run before this task | `["ensure-deps", "start-watch"]` |
| `dependsOrder`     | How dependencies are executed        | `"parallel"` or `"sequence"`     |
| `isBackground`     | Keep task running indefinitely       | `true` for watch tasks           |
| `runOptions.runOn` | When to auto-start the task          | `"folderOpen"`                   |

### Presentation Properties

| Property | Purpose                       | Example                            |
| -------- | ----------------------------- | ---------------------------------- |
| `group`  | Group related terminals       | `"watch"`                          |
| `reveal` | When to show the terminal     | `"never"`, `"always"`, `"silent"`  |
| `panel`  | Terminal panel behavior       | `"shared"`, `"dedicated"`, `"new"` |
| `clear`  | Clear terminal before running | `true` or `false`                  |
| `focus`  | Focus terminal when shown     | `true` or `false`                  |

### Problem Matchers

Problem matchers help VS Code parse output and show errors in the Problems panel:

```json
"problemMatcher": "$tsc-watch"        // TypeScript watch mode
"problemMatcher": "$webpack-watch"    // Webpack watch mode  
"problemMatcher": "$esbuild-watch"    // ESBuild watch mode
"problemMatcher": []                  // No problem matching
```

## Corresponding package.json Scripts

Your `package.json` should have the corresponding scripts:

```json
{
  "scripts": {
    "watch:tsc": "tsc --watch",
    "watch:webpack": "webpack --watch --mode development", 
    "watch:tests": "jest --watch",
    "dev:api": "nodemon src/server.ts",
    "watch:tailwind": "tailwindcss -w -i ./src/input.css -o ./dist/output.css"
  }
}
```

## Advanced Features

### Custom Problem Matchers

For custom tools, define problem matchers:

```json
"problemMatcher": {
    "pattern": {
        "regexp": "^ERROR\\s+(.*?):(\\d+):(\\d+)\\s+(.*)$",
        "file": 1,
        "line": 2, 
        "column": 3,
        "message": 4
    },
    "background": {
        "activeOnStart": true,
        "beginsPattern": "^Starting build...$",
        "endsPattern": "^Build complete.$"
    }
}
```

### Conditional Commands

Platform-specific commands:

```json
{
    "label": "ensure-deps",
    "type": "shell",
    "windows": {
        "command": "if not exist node_modules npm ci"
    },
    "linux": {
        "command": "test -d node_modules || npm ci"
    },
    "osx": {
        "command": "test -d node_modules || npm ci"
    }
}
```

### Environment Variables

Set environment variables for tasks:

```json
{
    "label": "watch:api-dev",
    "type": "npm",
    "script": "dev:api",
    "options": {
        "env": {
            "NODE_ENV": "development",
            "PORT": "3001",
            "DEBUG": "api:*"
        }
    }
}
```

## Terminal Management

### Grouping Terminals

Use `presentation.group` to organize terminals:

```json
"presentation": {
    "group": "watch",        // All watch tasks together
    "group": "servers",      // All server tasks together  
    "group": "tests"         // All test tasks together
}
```

### Terminal Panel Options

```json
"presentation": {
    "reveal": "always",      // Always show terminal when task starts
    "reveal": "silent",      // Show terminal but don't focus it
    "reveal": "never",       // Never show terminal (run in background)
    "panel": "shared",       // Use shared terminal panel
    "panel": "dedicated",    // Use dedicated terminal for this task
    "panel": "new",          // Always create new terminal
    "clear": true,           // Clear terminal before running
    "close": true            // Close terminal when task completes
}
```

## Troubleshooting

### Common Issues

1. **Tasks don't auto-start**: Check `runOptions.runOn` is set to `"folderOpen"`
2. **Terminals steal focus**: Set `presentation.reveal` to `"never"` or `"silent"`
3. **Tasks don't run in parallel**: Ensure `dependsOrder` is `"parallel"`
4. **Problem matchers not working**: Verify the pattern matches your tool's output format

### Debugging Tasks

Add a simple debug task to test configuration:

```json
{
    "label": "debug-task",
    "type": "shell",
    "command": "echo 'Task running in terminal: ${workspaceFolder}'",
    "group": "build",
    "presentation": {
        "reveal": "always",
        "focus": true
    }
}
```

### Viewing Task Output

- **Terminal Panel**: View real-time output from each task
- **Problems Panel**: See parsed errors and warnings
- **Output Panel**: View task execution logs (select "Tasks" from dropdown)

## Benefits

✅ **Automatic Setup**: No manual terminal management  
✅ **Organized Workflow**: All development processes in one place  
✅ **Error Integration**: Problems show up in VS Code's Problems panel  
✅ **Team Consistency**: Same setup for all team members  
✅ **Focus Management**: Terminals don't interrupt your workflow  
✅ **Resource Efficiency**: Only start what you need when you need it

## Real-World Example: VS Code Copilot Chat Extension

This pattern is used in Microsoft's VS Code Copilot Chat extension with these 4 parallel tasks:

1. **`watch:tsc-extension`** - TypeScript compilation for main extension
2. **`watch:tsc-extension-web`** - TypeScript compilation for web version  
3. **`watch:tsc-simulation-workbench`** - TypeScript compilation for test workbench
4. **`watch:esbuild`** - ESBuild bundling for production assets

The result: A seamless development experience where all necessary build processes
start automatically and run continuously in the background, with organized terminals
and integrated error reporting.
