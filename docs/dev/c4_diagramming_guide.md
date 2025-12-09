# C4 Architecture Diagramming in Go Projects

This guide provides practical guidance for generating and maintaining C4 architecture diagrams in Go
projects, specifically using the go-structurizr library and following Clean Architecture principles.

## Overview

C4 diagramming is a visual approach to describing software architecture through four different levels of abstraction:

- **Level 1: System Context** - Shows how your software system fits into the world around it
- **Level 2: Container** - Zooms into your system to show containers (applications, microservices, databases)
- **Level 3: Component** - Zooms into an individual container to show components
- **Level 4: Code** - Zooms into an individual component to show implementation details

This guide focuses primarily on Level 3 (Component) diagrams, which are most relevant for documenting
Clean Architecture implementations in Go.

## Why Auto-Generated C4 Diagrams?

### Advantages

- **Always Up-to-Date**: Diagrams reflect the current codebase state
- **Consistency**: Standardized representation across teams and projects
- **Reduced Maintenance**: No manual diagram updates required
- **Clean Architecture Alignment**: Visualizes dependency inversion and layering
- **Documentation as Code**: Diagrams live alongside implementation

### Use Cases

- Architecture reviews and discussions
- Onboarding new team members
- System documentation for stakeholders
- Refactoring planning and validation
- Dependency analysis and cleanup

## Getting Started with go-structurizr

### Installation

```bash
go get github.com/krzysztofreczek/go-structurizr
```

### Basic Implementation

1. **Define Component Information Interface**

Implement the `model.HasInfo` interface on your key architectural components:

```go
package domain

import "github.com/krzysztofreczek/go-structurizr/pkg/model"

type UserService struct {
    repository UserRepository
    eventBus   EventBus
}

func (s UserService) Info() model.Info {
    return model.ComponentInfo(
        "UserService",
        "Manages user business logic and workflows",
        "Go Service",
        "domain", "service",
    )
}
```

2. **Create Scraper Configuration**

```go
package main

import (
    "github.com/krzysztofreczek/go-structurizr/pkg/scraper"
    "github.com/krzysztofreczek/go-structurizr/pkg/view"
)

func createScraper() scraper.Scraper {
    config := scraper.NewConfiguration(
        "github.com/yourorg/yourproject", // Your project's module path
    )
    
    s := scraper.NewScraper(config)
    
    // Add rules for specific patterns
    rule, _ := scraper.NewRule().
        WithPkgRegexps("github.com/yourorg/yourproject/internal/.*").
        WithNameRegexp("^.*Repository$").
        WithApplyFunc(func(name string, _ ...string) model.Info {
            return model.ComponentInfo(
                name,
                "Data access layer",
                "Repository Pattern",
                "infrastructure", "repository",
            )
        }).
        Build()
    
    s.RegisterRule(rule)
    
    return s
}
```

3. **Generate Diagrams**

```go
func generateDiagram(app interface{}) error {
    scraper := createScraper()
    structure := scraper.Scrape(app)
    
    view := view.NewView().
        WithTitle("Application Architecture").
        WithComponentStyle(
            view.NewComponentStyle("domain").
                WithBackgroundColor(color.RGBA{135, 206, 235, 255}). // Sky blue
                WithFontColor(color.RGBA{0, 0, 0, 255}).
                WithBorderColor(color.RGBA{0, 0, 139, 255}).
                Build(),
        ).
        WithComponentStyle(
            view.NewComponentStyle("infrastructure").
                WithBackgroundColor(color.RGBA{255, 215, 0, 255}). // Gold
                WithFontColor(color.RGBA{0, 0, 0, 255}).
                WithBorderColor(color.RGBA{184, 134, 11, 255}).
                Build(),
        ).
        Build()
    
    outFile, err := os.Create("architecture.plantuml")
    if err != nil {
        return err
    }
    defer outFile.Close()
    
    return view.RenderStructureTo(structure, outFile)
}
```

## Clean Architecture Integration

### Layer-Based Component Tagging

Use consistent tags to represent Clean Architecture layers:

```go
// Domain Layer
func (s *UserService) Info() model.Info {
    return model.ComponentInfo(
        "UserService",
        "Core user business logic",
        "Go Service",
        "domain", "service",
    )
}

// Application Layer
func (h *CreateUserHandler) Info() model.Info {
    return model.ComponentInfo(
        "CreateUserHandler",
        "Handles user creation use case",
        "Command Handler",
        "application", "handler",
    )
}

// Interface Layer
func (c *UserController) Info() model.Info {
    return model.ComponentInfo(
        "UserController",
        "HTTP API for user operations",
        "REST Controller",
        "interface", "controller",
    )
}

// Infrastructure Layer
func (r *PostgreSQLUserRepository) Info() model.Info {
    return model.ComponentInfo(
        "PostgreSQLUserRepository",
        "PostgreSQL implementation of user repository",
        "Database Repository",
        "infrastructure", "repository",
    )
}
```

### YAML Configuration for Complex Projects

For larger projects, use YAML configuration files:

```yaml
# go-structurizr.yml
configuration:
  pkgs:
    - "github.com/yourorg/yourproject"

rules:
  # Domain Services
  - name_regexp: "^.*Service$"
    pkg_regexps:
      - "github.com/yourorg/yourproject/internal/domain/.*"
    component:
      description: "Domain service implementing business logic"
      technology: "Go Service"
      tags:
        - domain
        - service

  # Command Handlers
  - name_regexp: "^.*Handler$"
    pkg_regexps:
      - "github.com/yourorg/yourproject/internal/application/.*"
    component:
      description: "Application layer command/query handler"
      technology: "Handler"
      tags:
        - application
        - handler

  # Repositories
  - name_regexp: "^.*Repository$"
    pkg_regexps:
      - "github.com/yourorg/yourproject/internal/infrastructure/.*"
    component:
      description: "Infrastructure data access layer"
      technology: "Repository"
      tags:
        - infrastructure
        - repository

view:
  title: "Clean Architecture Component Diagram"
  line_color: "666666ff"
  styles:
    - id: domain
      background_color: "87ceebff"  # Sky blue
      font_color: "000000ff"
      border_color: "00008bff"      # Dark blue
      shape: rectangle

    - id: application
      background_color: "98fb98ff"  # Pale green
      font_color: "000000ff"
      border_color: "006400ff"      # Dark green
      shape: rectangle

    - id: interface
      background_color: "ffd700ff"  # Gold
      font_color: "000000ff"
      border_color: "b8860bff"      # Dark goldenrod
      shape: rectangle

    - id: infrastructure
      background_color: "ffa07aff"  # Light salmon
      font_color: "000000ff"
      border_color: "cd5c5cff"      # Indian red
      shape: rectangle

  component_tags:
    - domain
    - application
    - interface
    - infrastructure
```

## Best Practices

### 1. Component Naming and Description

- Use clear, descriptive component names that match your ubiquitous language
- Provide meaningful descriptions that explain the component's responsibility
- Include technology information to help understand implementation choices

```go
// Good
func (s *EmailNotificationService) Info() model.Info {
    return model.ComponentInfo(
        "EmailNotificationService",
        "Sends email notifications for user events using SMTP",
        "Go Service + SMTP",
        "infrastructure", "notification",
    )
}

// Avoid
func (s *EmailNotificationService) Info() model.Info {
    return model.ComponentInfo(
        "EmailService",
        "Sends emails",
        "Service",
        "service",
    )
}
```

### 2. Consistent Tagging Strategy

Establish a consistent tagging strategy across your project:

```go
// Layer tags
const (
    TagDomain         = "domain"
    TagApplication    = "application"
    TagInterface      = "interface"
    TagInfrastructure = "infrastructure"
)

// Pattern tags
const (
    TagService     = "service"
    TagRepository  = "repository"
    TagController  = "controller"
    TagHandler     = "handler"
    TagAdapter     = "adapter"
    TagGateway     = "gateway"
)

// Technology tags
const (
    TagHTTP       = "http"
    TagGRPC       = "grpc"
    TagDatabase   = "database"
    TagMessage    = "messaging"
    TagCache      = "cache"
)
```

### 3. Dependency Visualization

Structure your application context to clearly show dependencies:

```go
type ApplicationContext struct {
    // Domain Layer
    UserService    *domain.UserService
    OrderService   *domain.OrderService
    
    // Application Layer  
    CreateUserHandler *application.CreateUserHandler
    PlaceOrderHandler *application.PlaceOrderHandler
    
    // Interface Layer
    UserController  *interfaces.UserController
    OrderController *interfaces.OrderController
    
    // Infrastructure Layer
    UserRepository  repository.UserRepository
    OrderRepository repository.OrderRepository
    EmailService    *infrastructure.EmailService
}

func (ctx *ApplicationContext) Info() model.Info {
    return model.ComponentInfo(
        "ApplicationContext",
        "Main application composition root",
        "Dependency Container",
        "root",
    )
}
```

### 4. Filtering and Views

Create focused views for different audiences:

```go
// Architecture overview - focus on main components
architectureView := view.NewView().
    WithTitle("Architecture Overview").
    WithComponentTag("domain", "application").
    WithRootComponentTag("root").
    Build()

// Infrastructure details - focus on data layer
infrastructureView := view.NewView().
    WithTitle("Infrastructure Layer").
    WithComponentTag("infrastructure").
    Build()

// API surface - focus on interface layer
apiView := view.NewView().
    WithTitle("API Surface").
    WithComponentTag("interface", "application").
    Build()
```

## Integration with Build Process

### Makefile Integration

```makefile
.PHONY: docs-generate
docs-generate:
	@echo "Generating architecture diagrams..."
	@go run cmd/diagram-generator/main.go
	@plantuml -tpng docs/architecture.plantuml
	@echo "Diagrams generated in docs/ directory"

.PHONY: docs-validate
docs-validate: docs-generate
	@echo "Validating documentation is up-to-date..."
	@git diff --exit-code docs/ || (echo "Documentation is out of date. Run 'make docs-generate'" && exit 1)
```

### GitHub Actions Integration

```yaml
name: Documentation

on:
  pull_request:
    paths:
      - '**/*.go'
      - 'docs/**'

jobs:
  validate-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Generate diagrams
        run: make docs-generate
      
      - name: Check if diagrams are up-to-date
        run: |
          if [ -n "$(git status --porcelain docs/)" ]; then
            echo "Documentation is out of date"
            git diff docs/
            exit 1
          fi
```

## Troubleshooting

### Common Issues

1. **Components Not Appearing**
   - Ensure components implement `model.HasInfo` interface
   - Check package prefixes in scraper configuration
   - Verify naming patterns match rules

2. **Circular Dependencies**
   - go-structurizr handles circular dependencies gracefully
   - Use interfaces to break dependency cycles in your architecture
   - Consider if circular dependencies indicate design issues

3. **Too Many Components**
   - Use component and root component tags to filter views
   - Create multiple focused diagrams instead of one large diagram
   - Consider grouping related components

4. **Missing Relationships**
   - Ensure dependencies are properly injected and visible to the scraper
   - Use public fields or getters for dependencies you want to visualize
   - Check that dependency types implement required interfaces

### Debug Mode

Enable debug logging to troubleshoot scraping issues:

```bash
LOG_LEVEL=DEBUG go run cmd/diagram-generator/main.go
```

## Advanced Patterns

### Custom Rules for External Dependencies

```go
// Rule for external HTTP clients
externalHTTPRule, _ := scraper.NewRule().
    WithPkgRegexps(".*").
    WithNameRegexp("^.*Client$").
    WithApplyFunc(func(name string, _ ...string) model.Info {
        return model.ComponentInfo(
            name,
            "External HTTP client dependency",
            "HTTP Client",
            "external", "client",
        )
    }).
    Build()

// Rule for message queue adapters
messageQueueRule, _ := scraper.NewRule().
    WithPkgRegexps("github.com/yourorg/yourproject/internal/infrastructure/.*").
    WithNameRegexp("^.*Publisher$|^.*Consumer$").
    WithApplyFunc(func(name string, _ ...string) model.Info {
        tech := "Message Queue"
        if strings.Contains(name, "Publisher") {
            tech = "Message Publisher"
        } else if strings.Contains(name, "Consumer") {
            tech = "Message Consumer"
        }
        
        return model.ComponentInfo(
            name,
            "Message queue integration component",
            tech,
            "infrastructure", "messaging",
        )
    }).
    Build()
```

### Dynamic Component Information

```go
type DatabaseRepository struct {
    db       *sql.DB
    tableName string
}

func (r *DatabaseRepository) Info() model.Info {
    return model.ComponentInfo(
        fmt.Sprintf("%sRepository", strings.Title(r.tableName)),
        fmt.Sprintf("Database repository for %s entities", r.tableName),
        "PostgreSQL Repository",
        "infrastructure", "repository", "database",
    )
}
```

## Conclusion

Auto-generated C4 diagrams using go-structurizr provide a powerful way to visualize and document Go
applications following Clean Architecture principles. By implementing consistent tagging strategies,
meaningful component descriptions, and integrating diagram generation into your build process, you can
maintain always up-to-date architectural documentation that serves both development teams and stakeholders.

Remember to:

- Start simple and evolve your diagramming approach
- Focus on the most important architectural relationships
- Use multiple focused views rather than one complex diagram
- Integrate diagram generation into your CI/CD pipeline
- Regularly review and refine your component descriptions and tags

For more detailed examples and advanced usage patterns, refer to the [go-structurizr documentation](https://pkg.go.dev/github.com/krzysztofreczek/go-structurizr)
and the [Three Dots Labs article](https://threedots.tech/post/auto-generated-c4-architecture-diagrams-in-go/).
