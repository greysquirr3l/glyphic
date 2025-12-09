<!-- filepath: /Users/nickcampbell/Projects/go/documents/copilot-instructions.md -->

# 🤖 GitHub Copilot Instructions

This document provides guidelines for GitHub Copilot to follow when generating code for this project. These instructions represent our team's best practices, architectural preferences, and coding principles.

## 📋 Core Development Principles

### 1. Clean Architecture

Generate code that follows Clean Architecture principles:

- **Domain Layer**: Core business logic, entities, and rules
- **Application Layer**: Use cases, service orchestration
- **Adapters/Interface Layer**: Controllers, presenters, gateways
- **Infrastructure Layer**: Frameworks, tools, external services

💡 **Principle**: Dependencies should always point inward. Inner layers must not depend on outer layers.

### 2. Domain-Driven Design (DDD)

When suggesting domain models:

- Use **Ubiquitous Language** that matches our business domain
- Create **Value Objects** for concepts defined by their attributes (e.g., Email, Money)
- Model **Entities** with clear identity and lifecycle
- Define **Aggregates** with clear boundaries and access through aggregate roots
- Suggest appropriate **Domain Events** for significant state changes
- Keep business logic in the domain model, not in application services

#### DDD Lite Approach

Apply DDD pragmatically based on complexity:

- **Start with strategic design**: Focus on understanding bounded contexts before tactical patterns
- **Identify what matters most**: Apply full DDD treatment only where business complexity justifies it
- **Evolve gradually**: Begin simple and add complexity only when business needs require it
- **Pragmatic over pure**: Choose practical solutions that work for the team and domain size
- **Medium complexity domains**: More than CRUD but not enterprise-scale complexity
- **Keep aggregates small and focused**: Prefer multiple small aggregates over single large ones

### 3. Command Query Responsibility Segregation (CQRS)

Separate code that reads data from code that writes data:

- **Commands**: Operations that change state but don't return data
- **Queries**: Operations that return data but don't change state
- **Command Handlers**: Process commands and update the domain model
- **Query Handlers**: Process queries and return data representations

### 4. Idempotence and Predictability

All operations should be idempotent and predictable where appropriate:

- Executing an operation multiple times should produce the same outcome as executing it once
- Operations should be deterministic - same input always yields same output
- Side effects should be explicit and documented
- Avoid hidden state that affects the operation's outcome

### 5. Atomic Actions

Generate code that ensures related operations are atomic:

- State changes that must succeed or fail together should be in a single transaction
- No intermediate states should be visible to other components
- Use appropriate transaction mechanisms for the target platform

### 6. Testing Best Practices

Generated tests should follow:

- **Unit Tests**: Test business logic in isolation
- **Integration Tests**: Test infrastructure components with real dependencies
- **Test Clean Architecture**: Mock dependencies at architectural boundaries
- **Test Business Rules**: Focus on verifying business invariants
- **Test Readability**: Tests should document expected behavior

### 7. Security First Approach

All code must follow security best practices:

- **Input Validation**: Validate all inputs, especially user-supplied data
- **Authorization**: Include proper authorization checks
- **Secret Management**: No hardcoded credentials or secrets
- **SQL Injection Prevention**: Parameterized queries, not string concatenation
- **XSS Prevention**: Output encoding for user-supplied content
- **CSRF Protection**: For web applications

## 🧠 Software Engineering Principles

### 1. DRY (Don't Repeat Yourself)

Avoid code duplication but apply intelligently:

- **When DRY is beneficial**: True duplication, business rules, configuration
- **When to avoid DRY**: Different bounded contexts, accidental duplication, over-abstraction
- **Rule of three**: Don't abstract until you see the pattern three times
- **Bounded context boundaries**: Don't DRY across different domains
- **Test the abstraction**: If it's hard to test, it's probably wrong

#### DRY Anti-patterns to Avoid

- **Premature abstraction**: Don't abstract until you have at least 3 similar cases
- **Wrong abstraction**: Forcing unrelated code into shared functions
- **God functions**: Functions that do too many things to avoid duplication

### 2. YAGNI (You Aren't Gonna Need It)

Implement only what's needed now, not what might be needed later:

- Focus on current requirements
- Refactor when new needs arise
- Avoid "just in case" features

### 3. KISS (Keep It Simple, Stupid)

Prefer simple solutions over complex ones. Simplicity yields maintainability.

### 4. Single Responsibility Principle

Each component should have only one reason to change.

### 5. Favor Composition Over Inheritance

Build complex behaviors by combining simpler components.

### 6. Avoid Premature Optimization

Focus on correct functionality first, then optimize if needed based on measurements.

## 📁 Repository and Data Access

### 1. Repository Pattern Best Practices

When generating repository code:

- Focus on domain operations, not database operations
- Abstract away data access technology details
- Use domain language in repository method names
- Return domain objects, not data transfer objects
- Handle transactions appropriately for atomic operations
- Return domain errors, not infrastructure-specific errors

### 2. Idempotent Data Operations

For data modification operations:

- Implement conditional checks to prevent duplicate effects
- Use upsert patterns where appropriate
- Include idempotency keys for operations that may be retried

### 3. Secure Repository Design

When generating repository code, implement security by design:

- **Include user context in method signatures**: Make authorization requirements explicit
- **Implement authorization at the repository level**: Don't rely on external authorization checks
- **Use domain security rules**: Model authorization as business logic in the domain layer
- **Create purpose-built methods**: Avoid generic methods that require external authorization
- **Name internal methods explicitly**: Make security implications clear (e.g., `GetTrainingInternal`)
- **Fail secure by default**: Deny access unless explicitly authorized

```go
// Good: Secure by design
type TrainingRepository interface {
    GetTraining(ctx context.Context, trainingID string, user User) (*Training, error)
    FindTrainingsForUser(ctx context.Context, userUUID string) ([]*Training, error)
    GetTrainingInternal(ctx context.Context, trainingID string) (*Training, error) // Clearly marked
}

// Avoid: Security as external concern
type TrainingRepository interface {
    GetTraining(ctx context.Context, trainingID string) (*Training, error) // Who can access?
    FindAllTrainings(ctx context.Context) ([]*Training, error) // No filtering?
}
```

## 🌐 API Design

### 1. RESTful API Guidelines

- Use appropriate HTTP methods (GET, POST, PUT, DELETE)
- Return appropriate status codes
- Implement proper error handling and consistent error responses
- Design for idempotence (especially PUT and DELETE operations)
- Use consistent naming conventions

### 2. GraphQL Guidelines (if applicable)

- Model queries based on actual usage patterns, not database structure
- Implement proper authorization at the field level
- Consider query complexity and depth limitations
- Handle errors gracefully

## 🔄 Event-Driven Architecture

### 1. Event Design

- Events should describe facts that occurred, not commands
- Include all necessary context in the event payload
- Version events to allow for evolution
- Design events for idempotent handling

### 2. Event Handling

- Implement idempotent event consumers
- Handle out-of-order events gracefully
- Consider the Outbox Pattern for reliable event publishing
- Ensure event handlers can be retried safely

## 📦 Distributed Systems Considerations

### 1. Eventual Consistency

- Prefer eventual consistency over distributed transactions
- Design commands and events with eventual consistency in mind
- Implement appropriate compensating actions for failures

### 2. Resilience Patterns

- Implement circuit breakers for external service calls
- Use timeouts and retries with backoff strategies
- Design for partial failures

## 💻 Language-Specific Guidelines

When GitHub Copilot detects the programming language in use, it should adapt these principles to language-specific best practices.

### General Coding Standards

- Use meaningful variable, method, and class names
- Keep functions/methods small and focused
- Write code for humans first, computers second
- Add appropriate comments and documentation
- Follow the conventional style guide for the language in use

## 🧪 Example Patterns to Reference

When implementing functionality, reference these common patterns:

### Command Handler Pattern

```
function MyCommandHandler(dependencies) {
  return {
    handle: async (command) => {
      // 1. Validate command
      // 2. Load required entities/aggregates
      // 3. Execute domain logic
      // 4. Persist changes
      // 5. Publish events
    }
  };
}
```

### Repository Pattern

```
function MyRepository(dependencies) {
  return {
    findById: async (id) => {
      // Load entity and return domain object
    },
    save: async (entity) => {
      // Save entity with atomic guarantees
    }
  };
}
```

### Value Object Pattern

```
function createValueObject(data) {
  // Validate data
  // Create immutable object
  // Provide methods for business logic
  return Object.freeze({
    value: data,
    equals: (other) => data === other.value,
    // Domain methods...
  });
}
```

## 📚 Additional Considerations

- Apply [OSSF security baselines](https://github.com/ossf/security-baselines) when possible
- Consider accessibility requirements for UI code
- Optimize for performance only after functionality is correct
- Write for maintainability and readability over cleverness

## 🤝 Remember

These are guidelines, not rigid rules. Adapt them as necessary for the specific context and requirements of each task.
