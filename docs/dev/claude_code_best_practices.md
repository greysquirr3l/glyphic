# GitHub Copilot Best Practices for Go Development: Clean Architecture & DDD Integration

> A comprehensive guide for integrating GitHub Copilot into Go development workflows while maintaining Clean Architecture,
> Domain-Driven Design (DDD), and CQRS principles. Based on GitHub's engineering best practices and aligned with our
> project guidelines.

**📄 References:**

- [GitHub Copilot Best Practices](https://docs.github.com/en/copilot/using-github-copilot)
- [Project Guidelines: DDD, CQRS, and Clean Architecture](./project_guidelines.md)
- [Security-First Repository Design](./security_first_repository_design.md)
- [Software Engineering Principles](./software_principals.md)

---

## 📋 Table of Contents

1. [Setup Configuration](#-setup-configuration)
2. [Project Structure Integration](#️-project-structure-integration)
3. [Development Workflows](#-development-workflows)
4. [Clean Architecture Patterns](#-clean-architecture-patterns)
5. [Domain-Driven Design with GitHub Copilot](#-domain-driven-design-with-github-copilot)
6. [CQRS Implementation](#-cqrs-implementation)
7. [Testing Strategies](#-testing-strategies)
8. [Security Best Practices](#-security-best-practices)
9. [Advanced Workflows](#-advanced-workflows)
10. [Troubleshooting Optimization](#-troubleshooting-optimization)

---

## 🔧 Setup Configuration

### Suggested Project Structure

Following Clean Architecture and DDD principles, organize your Go project as follows:

```text
go-moss/                                    # Project root
├── .github/                               # GitHub configurations
│   ├── copilot-instructions.md           # GitHub Copilot project context
│   ├── prompts/                          # Claude analysis prompts
│   │   ├── README.md                    # Prompts collection guide
│   │   ├── validate-architecture.md      # Architecture validation
│   │   ├── audit-repositories.md         # Security audit
│   │   ├── domain-analysis.md           # DDD analysis
│   │   ├── cqrs-review.md              # CQRS implementation review
│   │   ├── test-analysis.md            # Test quality analysis
│   │   ├── performance-analysis.md     # Performance optimization
│   │   └── copilot-generation.md       # Code generation guide
│   └── workflows/                        # GitHub Actions
├── cmd/                                  # Application entry points
│   └── app/
│       └── main.go                      # Main application
├── internal/                            # Private application code
│   ├── domain/                         # Domain layer (Clean Architecture)
│   │   ├── entities/                   # Domain entities
│   │   │   ├── user.go
│   │   │   ├── training.go
│   │   │   └── session.go
│   │   ├── valueobjects/              # Value objects
│   │   │   ├── email.go
│   │   │   ├── money.go
│   │   │   └── address.go
│   │   ├── events/                    # Domain events
│   │   │   ├── user_created.go
│   │   │   └── training_completed.go
│   │   ├── services/                  # Domain services
│   │   │   └── pricing_service.go
│   │   ├── repositories/              # Repository interfaces
│   │   │   ├── user_repository.go
│   │   │   └── training_repository.go
│   │   ├── errors/                    # Domain-specific errors
│   │   │   └── domain_errors.go
│   │   └── COPILOT_INSTRUCTIONS.md   # Domain layer guidance
│   ├── application/                   # Application layer (CQRS)
│   │   ├── commands/                  # Command handlers
│   │   │   ├── create_user/
│   │   │   │   ├── command.go
│   │   │   │   ├── handler.go
│   │   │   │   └── handler_test.go
│   │   │   └── update_training/
│   │   ├── queries/                   # Query handlers
│   │   │   ├── get_user/
│   │   │   │   ├── query.go
│   │   │   │   ├── handler.go
│   │   │   │   └── handler_test.go
│   │   │   └── list_trainings/
│   │   ├── services/                  # Application services
│   │   │   └── user_service.go
│   │   └── COPILOT_INSTRUCTIONS.md   # Application layer guidance
│   ├── infrastructure/                # Infrastructure layer
│   │   ├── persistence/               # Database implementations
│   │   │   ├── postgres/
│   │   │   │   ├── user_repository.go
│   │   │   │   ├── training_repository.go
│   │   │   │   └── migrations/
│   │   │   └── memory/               # In-memory implementations
│   │   ├── messaging/                # Event bus implementations
│   │   │   └── redis_eventbus.go
│   │   ├── external/                 # External service clients
│   │   │   └── payment_gateway.go
│   │   └── config/                   # Configuration management
│   │       └── config.go
│   └── interface/                    # Interface layer
│       ├── http/                     # HTTP handlers
│       │   ├── handlers/
│       │   │   ├── user_handler.go
│       │   │   └── training_handler.go
│       │   ├── middleware/
│       │   │   ├── auth.go
│       │   │   └── logging.go
│       │   ├── routes/
│       │   │   └── routes.go
│       │   └── dto/                  # Data transfer objects
│       │       ├── user_dto.go
│       │       └── training_dto.go
│       ├── grpc/                     # gRPC servers (if applicable)
│       └── cli/                      # CLI commands (if applicable)
├── pkg/                              # Public library code
│   ├── errors/                       # Shared error types
│   ├── logger/                       # Logging utilities
│   └── validator/                    # Validation utilities
├── test/                             # Test files
│   ├── integration/                  # Integration tests
│   │   ├── user_integration_test.go
│   │   └── training_integration_test.go
│   ├── fixtures/                     # Test data
│   └── testutils/                    # Test utilities
├── scripts/                          # Build and deployment scripts
│   ├── install-go-tools.sh
│   ├── setup-pre-commit-hooks.sh
│   └── update-action-versions.sh
├── docs/                             # Documentation
│   ├── dev/                          # Development guides
│   │   ├── claude_code_best_practices.md
│   │   ├── project_guidelines.md
│   │   └── testing_go.md
│   └── api/                          # API documentation
├── deployments/                      # Deployment configurations
│   ├── docker/
│   └── kubernetes/
├── COPILOT_INSTRUCTIONS.md          # Root level project context
├── README.md                         # Project documentation
├── go.mod                           # Go module definition
├── go.sum                           # Go module checksums
├── Makefile                         # Build automation
└── .gitignore                       # Git ignore rules
```

### Create Project-Specific Context Files

Following GitHub Copilot best practices, create structured context files that align with our Go project standards:

#### Root Level: `COPILOT_INSTRUCTIONS.md`

```markdown
# Go-Moss Template: Clean Architecture & DDD

## Project Philosophy
- Clean Architecture with clear layer separation
- Domain-Driven Design (DDD) with ubiquitous language
- CQRS pattern for command/query separation
- Security-first repository design
- Test-driven development (TDD)

## Architecture Layers
- `internal/domain/` - Business logic, entities, value objects
- `internal/application/` - Use cases, command/query handlers
- `internal/infrastructure/` - Database, external services, repositories
- `internal/interface/` - HTTP handlers, gRPC servers, CLI

## Essential Commands
- `go run cmd/app/main.go` - Start the application
- `go test ./...` - Run all tests
- `go test -v -race ./internal/domain/...` - Run domain tests only
- `golangci-lint run` - Run linter
- `./scripts/install-go-tools.sh modern` - Install development tools
- `./scripts/setup-pre-commit-hooks.sh` - Setup pre-commit hooks
- `./scripts/update-action-versions.sh --dry-run` - Check GitHub Actions updates

## Code Standards
- Use `goimports` for import organization
- Follow Clean Architecture dependency rules (dependencies point inward)
- Use interfaces for dependency injection
- Implement domain models without external dependencies
- Use value objects for primitive obsession
- Repository interfaces defined in domain layer
- Error types should be domain-specific

## Testing Guidelines
- Unit tests alongside code (*_test.go)
- Integration tests in `/test/integration/`
- Use table-driven tests for multiple scenarios
- Mock external dependencies, not domain logic
- Test behavior, not implementation details

## Security Requirements
- All repository methods must include user context
- Authorization at repository level
- Input validation using domain value objects
- No secrets in environment variables or code
- Use prepared statements for SQL queries

## Domain Patterns
- Entities have identity and lifecycle
- Value objects are immutable and compared by value
- Aggregates maintain consistency boundaries
- Domain events for significant state changes
- Repository per aggregate root
- Use the UpdateFn pattern for transactional updates

## Git Workflow
- Use conventional commit format (feat, fix, docs, refactor, test)
- Create feature branches from main
- Include tests with all changes
- Update documentation when adding features
- Squash commits before merging to main

## Debug Mode
- Set `export GO_ENV=development` for debug logging
- Use `go run -race` to detect race conditions
- Enable SQL query logging with `?debug=true` connection parameter
```

#### Domain Layer: `internal/domain/COPILOT_INSTRUCTIONS.md`

```markdown
# Domain Layer Guidelines

## Core Principles
- No external dependencies (database, HTTP, etc.)
- Pure business logic only
- Express domain concepts using ubiquitous language
- Maintain aggregate consistency boundaries

## Patterns
- Entities: Identity + behavior (User, Order, Training)
- Value Objects: Immutable, compared by value (Email, Money, Address)
- Aggregates: Consistency boundaries with root entities
- Domain Services: Stateless business operations
- Domain Events: Capture significant state changes

## File Organization
- `/entities/` - Domain entities with identity
- `/valueobjects/` - Immutable value types
- `/events/` - Domain events for state changes
- `/services/` - Domain services for complex business rules
- `/repositories/` - Repository interfaces (implementation in infrastructure)

## Code Examples
```go
// Entity example
type User struct {
    id       UserID
    email    Email
    profile  UserProfile
    version  int // For optimistic locking
}

// Value Object example
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !isValidEmail(value) {
        return Email{}, ErrInvalidEmail
    }
    return Email{value: value}, nil
}

// Repository interface (in domain, implemented in infrastructure)
type UserRepository interface {
    Save(ctx context.Context, user *User, userCtx UserContext) error
    FindByID(ctx context.Context, id UserID, userCtx UserContext) (*User, error)
}
```

## Error Handling

- Domain-specific error types
- Errors should express business concepts
- Use structured errors with context

```go
// Example domain errors
type DomainError struct {
    Code    string
    Message string
    Context map[string]interface{}
}
```

## Application Layer: `internal/application/COPILOT_INSTRUCTIONS.md`

### Application Layer Guidelines

#### Responsibilities

- Orchestrate domain objects
- Handle use cases and business workflows
- Coordinate with infrastructure
- Manage transactions and persistence

#### CQRS Implementation

- Commands: State-changing operations
- Queries: Data retrieval operations
- Separate handlers for commands and queries
- Command handlers coordinate domain objects
- Query handlers return data representations

##### File Organization

- `/commands/` - Command definitions and handlers
- `/queries/` - Query definitions and handlers
- `/services/` - Application services for coordination

#### Command Handler Pattern

```go
type CreateUserCommand struct {
    Email string
    Name  string
}

type CreateUserHandler struct {
    userRepo domain.UserRepository
    eventBus domain.EventBus
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand, userCtx domain.UserContext) error {
    email, err := domain.NewEmail(cmd.Email)
    if err != nil {
        return err
    }
    
    user := domain.NewUser(email, cmd.Name)
    
    if err := h.userRepo.Save(ctx, user, userCtx); err != nil {
        return err
    }
    
    // Publish domain event
    event := domain.UserCreatedEvent{
        UserID: user.ID(),
        Email:  user.Email(),
    }
    
    return h.eventBus.Publish(ctx, event)
}
```

### Transaction Management

- Use the UpdateFn pattern for atomic updates
- Keep transactions short and focused
- Handle rollback scenarios gracefully

### Configure GitHub Copilot for Clean Architecture

Create a `.github/copilot-instructions.md` file for project-specific guidance:

```markdown
# GitHub Copilot Instructions for Go-Moss

## Code Generation Guidelines
- Always follow Clean Architecture dependency rules
- Generate domain models without external dependencies
- Use value objects for primitive obsession prevention
- Include proper error handling with domain-specific errors
- Generate comprehensive tests alongside code

## Security Requirements
- All repository methods must include UserContext parameter
- Use prepared statements for SQL queries
- Validate all inputs using domain value objects
- Include authorization checks in repository implementations

## Patterns to Follow
- Use CQRS pattern for commands and queries
- Implement repository interfaces in domain layer
- Use dependency injection for cross-layer communication
- Follow DDD patterns for aggregate design
```

### Claude Integration for Architecture Reviews

Use Claude 4 Sonnet alongside GitHub Copilot for comprehensive code reviews:

```bash
# Use Claude for architectural validation
# Example prompt for Claude:
"Review this GitHub Copilot-generated code for Clean Architecture compliance.
Check dependency directions, domain purity, and CQRS patterns."
```

---

## 🏗️ Project Structure Integration

### Use GitHub Copilot with Claude for Architecture Validation

Create a hybrid workflow combining GitHub Copilot's code generation with Claude's architectural analysis:

#### `.github/prompts/validate-architecture.md`

Please validate the Clean Architecture principles in this GitHub Copilot-generated Go code:

1. **Dependency Direction Check**:
   - Scan `internal/domain/` - should have NO imports from application, infrastructure, or interface layers
   - Scan `internal/application/` - should only import from domain layer
   - Scan `internal/infrastructure/` - can import from domain and application layers
   - Scan `internal/interface/` - can import from all layers

2. **Domain Purity Check**:
   - Domain entities should not import database drivers, HTTP libraries, or external APIs
   - Value objects should be immutable with proper validation
   - Repository interfaces should be defined in domain, implemented in infrastructure

3. **CQRS Compliance**:
   - Commands should change state, not return data
   - Queries should return data, not change state
   - Separate handlers for commands and queries

4. **Security Review**:
   - All repository methods should require user context
   - No hardcoded secrets or credentials
   - Input validation using value objects

Report any violations with specific file locations and recommendations for fixes.

### Repository Structure Validation

Use Claude to review GitHub Copilot-generated repository implementations for security compliance:

#### `.github/prompts/audit-repositories.md`

Audit all repository implementations for security compliance:

1. **User Context Requirements**:
   - Every repository method should include `UserContext` parameter
   - Methods should validate user permissions before data access
   - Internal methods should be clearly marked

2. **Authorization Patterns**:
   - Check for proper authorization at repository level
   - Verify domain-specific authorization rules
   - Ensure no bypass mechanisms exist

3. **SQL Security**:
   - All queries use prepared statements
   - No string concatenation for SQL building
   - Proper input sanitization

4. **Error Handling**:
   - Domain-specific error types
   - No sensitive information in error messages
   - Proper error propagation

Generate a security audit report with findings and recommendations.

---

## 🔄 Development Workflows

### 1. GitHub Copilot + Claude Domain-First Development Workflow

#### Step 1: Explore and Plan with Claude

```bash
# Start with domain exploration using Claude
```

**Claude Prompt:**

I need to implement a new feature: [describe feature].

First, help me understand the current domain model by reading the relevant files in `internal/domain/`.  
Don't write any code yet - just analyze the existing domain concepts and identify:

1. Which aggregates and entities are involved
2. What value objects might be needed
3. Any domain events that should be published
4. Repository interfaces that need to be defined

Think about this from a DDD perspective using ubiquitous language.

#### Step 2: Design Domain Model with Claude

**Claude Prompt:**

Based on your analysis, please create a plan for implementing this feature:

1. Domain layer changes (entities, value objects, events, repository interfaces)
2. Application layer changes (commands, queries, handlers)  
3. Infrastructure layer changes (repository implementations, database schema)
4. Interface layer changes (HTTP handlers, request/response models)

Make sure the plan follows Clean Architecture dependency rules. Create a GitHub issue with this plan before we start coding.

#### Step 3: Implement with GitHub Copilot + TDD

**GitHub Copilot Chat Prompt:**

Now let's implement using Test-Driven Development with GitHub Copilot:

1. Start with domain layer tests - write tests for the new domain behavior
2. Use Copilot to implement domain logic to make tests pass
3. Write application layer tests for command/query handlers
4. Use Copilot to implement application layer
5. Write integration tests for repositories
6. Use Copilot to implement infrastructure layer

Generate code that follows our Clean Architecture patterns and include comprehensive tests.

### 2. GitHub Copilot CQRS Command/Query Workflow

#### For Commands (State Changes)

**GitHub Copilot Chat Prompt:**

I need to implement a new command: [CommandName].

Follow the CQRS pattern and generate:

1. Command struct in internal/application/commands/
2. Command handler that:
   - Validates input using domain value objects
   - Loads required aggregates from repositories
   - Applies business logic through domain methods
   - Saves changes back to repositories
   - Publishes relevant domain events
3. Comprehensive tests for the handler
4. HTTP handler that uses the new command

Make sure the command doesn't return business data - only success/error status.
Follow our Clean Architecture patterns.

#### For Queries (Data Retrieval)

**GitHub Copilot Chat Prompt:**

I need to implement a new query: [QueryName].

Follow the CQRS pattern and generate:

1. Query struct in internal/application/queries/
2. Query handler that:
   - Validates query parameters
   - Retrieves data from read-optimized repositories
   - Returns structured data representations
   - Does NOT modify any state
3. Tests that verify correct data retrieval
4. HTTP handler that uses the new query

Optimize for read performance and include proper authorization checks.

### 3. Security-First Repository Implementation with GitHub Copilot

**GitHub Copilot Chat Prompt:**

I need to implement a new repository: [RepositoryName].

Follow our security-first design and generate:

1. Repository interface in internal/domain/ with:
   - All methods requiring UserContext parameter
   - Clear authorization semantics in method names
   - Internal methods clearly marked (e.g., GetTrainingInternal)

2. Implementation in internal/infrastructure/ with:
   - Authorization checks at the beginning of each method
   - Domain-specific error returns
   - Prepared statements for all SQL queries
   - Proper transaction handling

3. Comprehensive tests covering:
   - Authorization success and failure scenarios
   - Data access patterns
   - SQL injection prevention
   - Concurrent access scenarios

4. Documentation of the authorization rules in the repository interface comments.

Follow Clean Architecture patterns and our security-first repository design.

---

## 🎯 Clean Architecture Patterns

### Layer Interaction Validation with Claude

Use Claude to validate GitHub Copilot-generated architectural boundaries:

**Claude Prompt:**

Please analyze this GitHub Copilot-generated code change and verify it follows Clean Architecture:

1. Check dependency directions - inner layers should not depend on outer layers
2. Verify that domain layer remains pure (no external dependencies)
3. Ensure interfaces are defined in inner layers and implemented in outer layers
4. Confirm that all cross-layer communication goes through interfaces
5. Validate that infrastructure concerns don't leak into business logic

If any violations are found, suggest specific refactoring steps to fix them.

### Dependency Injection Pattern with GitHub Copilot

**GitHub Copilot Chat Prompt:**

I need to set up dependency injection for this application layer service.

Follow our Clean Architecture patterns and generate:

1. Interface definitions in the domain layer
2. Concrete type implementations in infrastructure layer  
3. Dependency wiring in the main function or DI container
4. Constructor injection pattern usage
5. Small, focused interfaces (Interface Segregation Principle)

Show me the complete dependency wiring including repository interfaces, application services, and HTTP handlers.

---

## 🎨 Domain-Driven Design with GitHub Copilot

### Ubiquitous Language Development

Help me develop the ubiquitous language for the [domain area] in our application.

1. Review existing domain models in internal/domain/
2. Identify the core concepts and their relationships
3. Suggest more expressive names that match business terminology
4. Recommend value objects to replace primitive obsession
5. Identify missing domain concepts that should be modeled explicitly

Work with me to refine the language until it clearly expresses the business domain.

### Aggregate Design

I need to design a new aggregate for [business concept].

Follow DDD principles:

1. Identify the aggregate root entity
2. Determine which entities and value objects belong in this aggregate
3. Define the consistency boundary - what must change together atomically
4. Design methods that enforce business invariants
5. Identify domain events that should be published when state changes
6. Keep the aggregate focused and not too large

Create the aggregate with proper encapsulation and business-focused methods.

### Value Object Implementation

I need to create value objects to eliminate primitive obsession in [area].

For each value object:

1. Identify the validation rules and business constraints
2. Make the value object immutable
3. Implement equality comparison by value
4. Add meaningful methods that express domain operations
5. Include proper error handling for invalid values
6. Add comprehensive tests covering edge cases

Replace primitive types throughout the domain model with these value objects.

---

## ⚡ CQRS Implementation

### Command Handler Generation

Generate a complete CQRS command handler for [operation]:

1. Command struct with proper validation tags
2. Handler implementation following our patterns:
   - Load aggregates using repositories
   - Apply business logic through domain methods
   - Handle transactions properly
   - Publish domain events
   - Return domain-specific errors
3. Comprehensive unit tests
4. Integration tests with real repositories
5. HTTP handler that uses the command

Ensure the command is idempotent and handles concurrent access safely.

### Query Optimization

Optimize this query handler for better performance:

1. Analyze the current query patterns and identify bottlenecks
2. Suggest database indexing strategies
3. Consider read model optimizations
4. Implement caching where appropriate
5. Add pagination for large result sets
6. Include performance tests to measure improvements

Maintain the separation between command and query sides while optimizing for read performance.

---

## 🧪 Testing Strategies

### Domain-Driven Testing

Create comprehensive tests for this domain model following our testing guidelines:

1. **Unit Tests**: Test business logic in isolation
   - Entity behavior and invariant enforcement
   - Value object validation and operations
   - Domain service business rules
   - Use table-driven tests for multiple scenarios

2. **Integration Tests**: Test infrastructure components
   - Repository implementations with real database
   - Event publishing and handling
   - External service integrations

3. **Architecture Tests**: Verify architectural compliance
   - Dependency direction validation
   - Layer isolation testing
   - Interface compliance verification

Include both positive and negative test cases, edge cases, and concurrent access scenarios.

### Test-Driven Development with Claude

Let's implement [feature] using Test-Driven Development:

1. **Red Phase**: Write failing tests that describe the desired behavior
   - Domain logic tests
   - Application layer tests  
   - Integration tests
   Run tests to confirm they fail for the right reasons

2. **Green Phase**: Write minimal code to make tests pass
   - Implement domain entities and value objects
   - Create application handlers
   - Add repository implementations
   - Don't over-engineer - just make tests pass

3. **Refactor Phase**: Improve code quality while keeping tests green
   - Extract common patterns
   - Improve naming and clarity
   - Optimize performance if needed
   - Maintain test coverage

Repeat this cycle for each piece of functionality.

---

## 🔒 Security Best Practices

### Security Audit Workflow

Perform a comprehensive security audit of this code:

1. **Input Validation**:
   - All user inputs validated using value objects
   - SQL injection prevention through prepared statements
   - Cross-site scripting (XSS) prevention in outputs

2. **Authorization**:
   - Every repository method includes user context
   - Authorization checks at appropriate boundaries
   - Proper error messages that don't leak sensitive information

3. **Data Protection**:
   - No hardcoded secrets or credentials
   - Sensitive data encrypted at rest and in transit
   - Proper session management

4. **Code Quality**:
   - No obvious security vulnerabilities
   - Proper error handling without information disclosure
   - Secure defaults throughout the codebase

Generate a detailed security report with findings and remediation steps.

### OSSF Security Baseline Compliance

Audit this Go project against OSSF Security Baselines:

1. **Code Security**:
   - Static analysis with gosec
   - Dependency vulnerability scanning
   - Code review requirements

2. **Supply Chain Security**:
   - Dependency pinning in go.mod
   - SBOM generation
   - Signed builds and releases

3. **Development Security**:
   - Pre-commit hooks for security checks
   - Secrets detection in repository
   - Secure development practices

Generate a compliance report and action items for any gaps found.

---

## 🚀 Advanced Workflows

### Multi-Tool Architecture Review: GitHub Copilot + Claude

For complex features spanning multiple architectural layers, use a hybrid approach:

#### Terminal 1: Claude Domain Expert

**Claude Prompt:**

You are a domain modeling expert. Focus only on the domain layer:

1. Review domain models for proper DDD implementation
2. Ensure entities, value objects, and aggregates are well-designed
3. Validate that business logic is properly encapsulated
4. Check that domain events capture important state changes
5. Verify repository interfaces express domain concepts

Do not implement anything - only provide domain modeling feedback for GitHub Copilot implementation.

#### Terminal 2: GitHub Copilot Application Layer

**GitHub Copilot Chat Prompt:**

You are implementing the application layer. Focus on CQRS and use cases:

1. Generate command and query handlers with proper separation
2. Ensure handlers coordinate domain objects correctly  
3. Implement transaction boundaries and error handling
4. Create application services without business logic
5. Generate proper event publishing and handling

Implement application layer components based on domain feedback from Claude.

#### Terminal 3: GitHub Copilot Infrastructure

**GitHub Copilot Chat Prompt:**

You are implementing the infrastructure layer. Focus on technical implementation:

1. Generate repository patterns with proper security
2. Create database schemas that support domain models
3. Configure dependency injection and application startup
4. Set up observability, logging, and monitoring
5. Ensure infrastructure doesn't leak into domain

Build infrastructure that supports the domain and application layers following Clean Architecture.

### Automated Code Quality Pipeline

Set up an automated code quality pipeline using GitHub Copilot and Claude integration:

1. **GitHub Copilot for Code Generation**:
   - Generate code following architectural patterns
   - Create comprehensive tests automatically
   - Implement security-first repository patterns
   - Follow CQRS and DDD principles

2. **Claude for Code Review and Validation**:
   - Architectural compliance validation
   - Security pattern verification
   - DDD principle enforcement
   - Performance analysis recommendations

3. **GitHub Actions Integration**:
   - Run GitHub Copilot-powered code generation
   - Use Claude for architecture compliance reports
   - Generate automated documentation updates
   - Trigger security and performance reviews

Make it fully automated while maintaining high code quality standards through the combined power of GitHub Copilot and Claude.

---

## 🔧 Troubleshooting Optimization

### Common Issues and Solutions

#### Issue: GitHub Copilot suggests breaking Clean Architecture rules

**Solution:**

**Claude Review Prompt:**

GitHub Copilot suggested importing a database package in the domain layer.
This violates our Clean Architecture principles.

Please provide guidance to revise the solution:

1. Keep the domain layer pure with no external dependencies
2. Define repository interfaces in the domain layer
3. Implement repositories in the infrastructure layer
4. Use dependency injection to wire them together

The domain should only depend on standard library packages and other domain concepts.

#### Issue: Overly complex domain models from Copilot

**Solution:**

**Claude Prompt:**

This GitHub Copilot-generated domain model seems overly complex. Let's simplify following DDD principles:

1. Focus on the core business concepts and relationships
2. Remove technical concerns from domain entities
3. Use value objects to encapsulate validation and business rules
4. Keep aggregates small and focused
5. Express business operations through domain methods

Can you provide guidance to refactor this to be more focused on the business domain?

#### Issue: Poor test coverage from Copilot generation

**Solution:**

**GitHub Copilot Chat + Claude Review:**

The GitHub Copilot-generated test coverage looks incomplete. Please enhance the tests to include:

1. **Business Logic Tests**: Test all domain entity methods and business rules
2. **Edge Cases**: Test boundary conditions and error scenarios  
3. **Integration Tests**: Test repository implementations with real data
4. **Concurrent Access**: Test thread safety where applicable
5. **Performance Tests**: Verify acceptable response times

Use table-driven tests for multiple scenarios and descriptive test names.

### Performance Optimization

**Combined GitHub Copilot + Claude Approach:**

Use GitHub Copilot to implement optimizations and Claude to validate architectural integrity:

**GitHub Copilot Chat Prompt:**
Optimize this Go application for better performance while maintaining architectural integrity:

1. **Database Optimization**:
   - Add appropriate indexes for query patterns
   - Optimize repository query implementations
   - Use connection pooling and prepared statements
   - Consider read replicas for query-heavy operations

2. **Application Performance**:
   - Profile memory and CPU usage
   - Optimize hot paths in domain logic
   - Add caching where appropriate (without breaking CQRS)
   - Use efficient data structures and algorithms

3. **Monitoring & Observability**:
   - Add structured logging with correlation IDs
   - Implement metrics collection and alerting
   - Set up distributed tracing for request flows
   - Monitor application and infrastructure health

**Claude Review Prompt:**
Review the GitHub Copilot-generated performance optimizations to ensure they maintain clean architecture principles.

---

## 📚 Integration with Existing Documentation

This guide integrates with and extends our existing documentation:

- **[Project Guidelines](./project_guidelines.md)**: Core DDD, CQRS, and Clean Architecture principles
- **[Security-First Repository Design](./security_first_repository_design.md)**: Repository security patterns
- **[Software Engineering Principles](./software_principals.md)**: Fundamental coding principles
- **[Testing Guidelines](./testing_go.md)**: Go-specific testing best practices
- **[Concurrency Patterns](./go_concurrency.md)**: Safe concurrent programming in Go
- **[Battle-Tested Libraries](./battle_tested_go_libraries.md)**: Recommended Go libraries for each layer

### Quick Reference Commands

Create these as aliases in your shell for GitHub Copilot + Claude workflows:

```bash
# Architecture validation with Claude
alias validate-arch="claude -p 'Review GitHub Copilot-generated code for Clean Architecture compliance'"

# Security audit with Claude
alias security-audit="claude -p 'Audit GitHub Copilot-generated repositories for security compliance'"

# Test coverage check with GitHub Copilot
alias check-tests="copilot chat 'Analyze test coverage and generate comprehensive tests for untested code'"

# Performance analysis with combined tools
alias perf-check="copilot chat 'Profile application performance and suggest optimizations' && claude -p 'Review optimizations for architectural compliance'"
```

---

## 🎯 Best Practices Summary

### Do's ✅

- **Start with domain modeling** - Use Claude to explore and design domain concepts, then implement with GitHub Copilot
- **Validate architecture** - Use Claude to regularly check Clean Architecture compliance of Copilot-generated code
- **Test-driven development** - Use GitHub Copilot to generate tests, validate with Claude
- **Security-first design** - Generate secure code with Copilot, validate authorization patterns with Claude
- **Iterative refinement** - Use GitHub Copilot for rapid iteration, Claude for architectural feedback loops
- **Multi-tool workflows** - Use GitHub Copilot for implementation, Claude for design and review
- **Document decisions** - Use GitHub issues and ADRs for architectural decisions

### Don'ts ❌

- **Don't skip domain analysis** - Always use Claude for domain understanding before Copilot implementation
- **Don't violate layer boundaries** - Use Claude to validate strict dependency direction rules
- **Don't bypass authorization** - Ensure Copilot-generated code includes proper user context
- **Don't ignore architectural feedback** - Always review Copilot output with Claude for compliance
- **Don't over-engineer** - Start simple with Copilot, add complexity only when needed
- **Don't mix concerns** - Use Claude to verify business logic separation from infrastructure
- **Don't hardcode secrets** - Configure Copilot to avoid secrets, validate with Claude

---

## 🔮 Future Enhancements

- **AI-Powered Architecture Reviews**: Automated Pull Request reviews using GitHub Copilot + Claude for architectural compliance
- **Domain Model Evolution**: Tracking and managing domain model changes over time with hybrid AI assistance
- **Performance Regression Detection**: Automated performance testing using GitHub Copilot in CI/CD with Claude validation
- **Security Vulnerability Scanning**: Continuous security monitoring using combined GitHub Copilot and Claude analysis
- **Documentation Generation**: Auto-generated API docs and architectural diagrams using GitHub Copilot with Claude oversight
- **Refactoring Assistance**: AI-guided code refactoring using GitHub Copilot implementation with Claude architectural validation

---

*This document is a living resource that should be updated as we discover new patterns and practices.  
Contribute improvements via pull requests and share your GitHub Copilot + Claude workflows with the team.*

**Last updated:** August 2025  
**Version:** 1.0  
**Authors:** Development Team  
**Reviewers:** Architecture Committee
