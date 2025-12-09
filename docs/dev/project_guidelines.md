# Project Guidelines: DDD, CQRS, and Clean Architecture in Go

<!-- REF: https://threedots.tech/tags/modern-go/ -->
<!-- REF: https://threedots.tech/post/clean-architecture-in-go/ -->

## Introduction

This document provides guidelines for developing our app using Domain-Driven Design
(DDD), Command Query Responsibility Segregation (CQRS), and Clean Architecture
principles in Go. These approaches combine to create maintainable, testable, and
business-focused code.

## Core Concepts

### Domain-Driven Design (DDD)

DDD focuses on creating a software model that reflects the business domain:

<!-- REF: https://threedots.tech/post/ddd-in-go/ -->

- **Ubiquitous Language**: Use consistent terminology between developers and domain experts
- **Bounded Contexts**: Divide the domain into distinct areas with clear boundaries
- **Aggregates**: Treat related entities as a single unit with a root entity
- **Domain Events**: Model significant occurrences within the domain
- **Value Objects**: Immutable objects defined by their attributes
- **Entities**: Objects with identity that persists across state changes
- **Repositories**: Abstractions for persisting and retrieving domain objects

#### DDD Lite Approach

Following the approach from the "Modern Business Software in Go" series, we adopt
a pragmatic version of DDD that focuses on the most valuable aspects without
excessive complexity:

<!-- REF: https://threedots.tech/post/ddd-lite-in-go-introduction/ -->
<!-- REF: https://threedots.tech/post/ddd-in-go-lite/ -->

**Core Principles of DDD Lite:**

- **Start with strategic design**: Focus on understanding the business domain and
  identifying bounded contexts before diving into tactical patterns
- **Identify what matters most**: Not every part of your system needs full DDD
  treatment—apply it where business complexity justifies it
- **Evolve gradually**: Begin simple and add complexity only when business needs require it
- **Pragmatic over pure**: Choose practical solutions that work for your team and domain size

**When to Apply DDD Lite:**

- **Medium complexity domains**: More than CRUD but not enterprise-scale complexity
- **Growing applications**: Starting simple but expecting business logic to grow
- **Learning teams**: Teams new to DDD who want to adopt patterns gradually
- **Brownfield projects**: Existing projects that need better domain modeling

**DDD Lite Implementation Strategy:**

- **Focus on bounded contexts first**: Clearly define boundaries between different
  parts of the system
- **Identify aggregates**: Group related entities, but keep aggregates small and
  focused
- **Use value objects**: Encapsulate validation and business rules in value objects
- **Apply tactical patterns selectively**: Use repositories, domain services, and
  entities where they provide clear value

### CQRS (Command Query Responsibility Segregation)

CQRS separates operations that read data from operations that write data:

<!-- REF: https://threedots.tech/post/cqrs-in-go/ -->
<!-- REF: https://threedots.tech/post/basic-cqrs-in-go/ -->

- **Commands**: Operations that change state but don't return data
- **Queries**: Operations that return data but don't change state
- **Command Handlers**: Process commands and update the domain model
- **Query Handlers**: Process queries and return data representations

#### Practical CQRS in Go

- **Keep it simple:** You don't need a framework—CQRS can be implemented with
  plain Go structs and interfaces.
- **Command/Query separation:** Define separate types and handlers for commands
  and queries, even if they are implemented in the same package.
- **Handler interfaces:** Use interfaces like `Handle(ctx, cmd)` for commands and
  `Handle(ctx, query)` for queries.
- **Decouple side effects:** Command handlers should encapsulate side effects
  (e.g., database writes, publishing events), while query handlers should be
  side-effect free.
- **Testing:** CQRS makes it easier to test business logic in isolation, as
  commands and queries are explicit and handlers are small.

**Example:**

```go
// Command
type CreateUserCommand struct {
    Email    string
    Password string
}

// Command Handler
type CreateUserHandler struct {
    repo UserRepository
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
    // ...validate, hash password, etc...
    user := NewUser(cmd.Email, cmd.Password)
    return h.repo.Save(ctx, user)
}

// Query
type GetUserByEmailQuery struct {
    Email string
}

// Query Handler
type GetUserByEmailHandler struct {
    repo UserRepository
}

func (h *GetUserByEmailHandler) Handle(ctx context.Context, q GetUserByEmailQuery) (*User, error) {
    return h.repo.FindByEmail(ctx, q.Email)
}
```

- **No need for a "bus":** For most Go projects, you can call handlers directly
  instead of introducing a command/query bus abstraction.
- **Keep models simple:** You can use the same domain models for both commands
  and queries unless you have specific performance or security needs.

> For more, see [Basic CQRS in Go](https://threedots.tech/post/basic-cqrs-in-go/)

### Clean Architecture

<!-- REF: https://threedots.tech/post/clean-architecture-in-go/ -->

Clean Architecture organizes code in concentric layers, with dependencies always
pointing inward. The innermost layers are the most stable and least likely to change,
while the outer layers handle technical details and frameworks.

- **Dependency Rule:** Code in outer layers can depend on inner layers, but never
  the other way around. The domain layer must not depend on any other layer.
- **Domain Layer:** Pure business logic, entities, value objects, and domain services.
  No dependencies on frameworks or infrastructure.
- **Application Layer:** Orchestrates use cases, coordinates domain objects, and
  defines interfaces (ports) for infrastructure. Contains business rules that are
  specific to application workflows.
- **Interface/Adapter Layer:** Implements ports, handles input/output (e.g., HTTP,
  gRPC, CLI), and translates between external and internal representations. Adapters
  implement interfaces defined in the application layer.
- **Infrastructure Layer:** Frameworks, databases, external services, and technical
  details. Implements interfaces (ports) defined in the application layer.

#### Ports and Adapters (Hexagonal Architecture)

<!-- REF: https://threedots.tech/post/ports-adapters-hexagonal-architecture-go/ -->

Use the Ports and Adapters (Hexagonal) pattern to decouple the core from the outside world:

- **Ports:** Interfaces defined by the application layer (e.g., repository interfaces, service interfaces).
- **Adapters:** Implementations of those interfaces (e.g., SQL repository, HTTP handler).

#### Dependency Injection

<!-- REF: https://threedots.tech/post/dependency-injection-in-go/ -->

Dependencies (e.g., repositories, services) should be injected into constructors,
not created inside the business logic. This makes the core easy to test and swap
implementations.

```go
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

#### Testing Benefits

Clean Architecture makes it easy to test business logic in isolation, as the domain
and application layers have no dependencies on frameworks or databases.

#### Practical Go Advice

<!-- REF: https://threedots.tech/post/modern-go-best-practices/ -->

- Start simple; only introduce more layers and abstractions as the project grows.
- Use Go interfaces to invert dependencies, but avoid unnecessary abstractions.

#### Example Directory Structure

Your project structure already reflects these principles. Map each folder to its
corresponding layer and keep dependencies flowing inward.

---

**Summary:**
Clean Architecture in Go means keeping your business logic at the center, protected
from changes in frameworks, databases, and delivery mechanisms. Use interfaces
(ports) and dependency injection to invert dependencies, and keep your domain pure
and testable. Only add complexity as your project demands it.

## Project Structure

<!-- REF: https://threedots.tech/post/modern-go-project-layout/ -->

```sh
/cmd                  # Application entry points
  /api                # API server
  /worker             # Background workers
/internal
  /domain             # Domain model, entities, value objects, domain services
    /{boundedcontext} # Specific bounded contexts
  /application        # Application services, commands, queries, handlers
    /command          # Command definitions and handlers
    /query            # Query definitions and handlers
  /ports              # Ports (interfaces) required by the application
  /adapters           # Adapters implementing the ports
    /primary          # Input adapters (REST API, gRPC, etc.)
    /secondary        # Output adapters (repositories, external services)
  /infrastructure     # Infrastructure concerns, frameworks, DB
/pkg                  # Public packages
/docs                 # Documentation
```

## Implementation Guidelines

<!-- REF: https://threedots.tech/post/ddd-in-go-lite/ -->

### Domain Layer

1. **Start with the domain model**:

   - Define entities and value objects
   - Implement domain logic within aggregates
   - Define domain events

   ```go
   // Example of a domain entity
   type User struct {
       ID       UserID
       Email    Email     // Value object
       Password Password  // Value object
       Role     UserRole  // Enum
       Active   bool
   }

   // Domain logic within entity
   func (u *User) ChangePassword(currentPassword, newPassword Password) error {
       if !u.Password.Matches(currentPassword) {
           return ErrInvalidPassword
       }
       u.Password = newPassword
       return nil
   }
   ```

2. **Use value objects for validation**:

   ```go
   type Email string

   func NewEmail(email string) (Email, error) {
       // Validation logic
       if !validEmail(email) {
           return "", ErrInvalidEmail
       }
       return Email(email), nil
   }
   ```

3. **Keep business logic in the domain**:

   ```go
   // Good: Domain logic in the domain model
   func (t *Training) CanBeAttendedBy(u *User) error {
       if !t.IsActive() {
           return ErrTrainingNotActive
       }
       if u.HasActiveSubscription() {
           return nil
       }
       return ErrUserHasNoActiveSubscription
   }

   // Bad: Domain logic in application services
   func (s *TrainingService) AttendTraining(userID, trainingID string) error {
       training, err := s.trainingRepo.ByID(trainingID)
       if err != nil {
           return err
       }

       if !training.IsActive() {
           return ErrTrainingNotActive
       }

       user, err := s.userRepo.ByID(userID)
       if err != nil {
           return err
       }

       if !user.HasActiveSubscription() {
           return ErrUserHasNoActiveSubscription
       }

       // Continue with attendance logic...
   }
   ```

### Application Layer

1. **Define commands and queries**:

   ```go
   // Command
   type CreateUser struct {
       Email    string
       Password string
   }

   // Query
   type GetUserByID struct {
       ID string
   }
   ```

2. **Implement command/query handlers**:

   ```go
   // Command handler
   type CreateUserHandler struct {
       userRepo UserRepository
   }

   func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUser) error {
       email, err := NewEmail(cmd.Email)
       if err != nil {
           return err
       }

       password, err := NewPassword(cmd.Password)
       if err != nil {
           return err
       }

       user := NewUser(email, password)
       return h.userRepo.Save(ctx, user)
   }
   ```

### Infrastructure Layer

1. **Implement repositories**:

   ```go
   type PostgresUserRepository struct {
       db *sql.DB
   }

   func (r *PostgresUserRepository) Save(ctx context.Context, user *User) error {
       // Implementation
   }
   ```

2. **Use dependency injection**:

   ```go
   func NewCreateUserHandler(repo UserRepository) *CreateUserHandler {
       return &CreateUserHandler{userRepo: repo}
   }
   ```

3. **Repository best practices**:

<!-- REF: https://threedots.tech/post/repository-pattern-in-go/ -->

```go
// Define repository interfaces in the domain layer
type TrainingRepository interface {
    GetByID(ctx context.Context, id string) (*Training, error)
    Save(ctx context.Context, training *Training) error
    FindAllAvailable(ctx context.Context) ([]*Training, error)
}

// Implement in the infrastructure layer
type SQLTrainingRepository struct {
    db *sql.DB
}

func (r *SQLTrainingRepository) GetByID(ctx context.Context, id string) (*Training, error) {
    // SQL implementation
}

// Return domain errors, not infrastructure errors
func (r *SQLTrainingRepository) GetByID(ctx context.Context, id string) (*Training, error) {
    training := &Training{}
    err := r.db.QueryRowContext(ctx, "SELECT * FROM trainings WHERE id = $1", id).Scan(
        &training.ID, &training.Name, &training.Status, // ...other fields
    )
    if err == sql.ErrNoRows {
        return nil, domain.ErrTrainingNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("error getting training: %w", err)
    }
    return training, nil
}
```

## Security-First Repository Design

<!-- REF: https://threedots.tech/post/repository-secure-by-design/ -->

Repository security should be built into the design, not added as an afterthought.
The key principle is to make it **impossible to use repositories incorrectly**.

### Core Security Principles

1. **Authorization at Repository Level**: Include user context in repository method signatures
2. **Domain-Driven Security**: Model security rules as domain logic, not technical details
3. **Explicit Over Implicit**: Make security requirements visible in interfaces
4. **Fail Secure**: Default to denying access rather than allowing it

### Secure Repository Interface Design

```go
// Good: Security built into interface
type TrainingRepository interface {
    GetTraining(ctx context.Context, trainingID string, user User) (*Training, error)
    UpdateTraining(ctx context.Context, trainingID string, user User, updateFn func(*Training) (*Training, error)) error
    
    // Internal operations clearly marked
    GetTrainingInternal(ctx context.Context, trainingID string) (*Training, error)
}

// Avoid: Security as external concern
type TrainingRepository interface {
    GetTraining(ctx context.Context, trainingID string) (*Training, error)
    UpdateTraining(ctx context.Context, trainingID string, updateFn func(*Training) (*Training, error)) error
}
```

### Domain Security Rules

```go
// Security logic belongs in the domain layer
func CanUserSeeTraining(user User, training Training) error {
    if user.Type() == TrainerUserType {
        return nil // Trainers can see all trainings
    }
    if user.UUID() == training.UserUUID() {
        return nil // Users can see their own trainings
    }
    
    return ForbiddenToSeeTrainingError{
        RequestingUserUUID: user.UUID(),
        TrainingOwnerUUID:  training.UserUUID(),
    }
}
```

### Implementation Example

```go
func (r *PostgresTrainingRepository) GetTraining(
    ctx context.Context,
    trainingID string,
    user domain.User,
) (*domain.Training, error) {
    // Get training from database
    training, err := r.getTrainingByID(ctx, trainingID)
    if err != nil {
        return nil, err
    }
    
    // Apply domain authorization rule
    if err := domain.CanUserSeeTraining(user, *training); err != nil {
        return nil, err
    }
    
    return training, nil
}
```

**For comprehensive security guidance, see [Security-First Repository Design](security_first_repository_design.md).**

## Repository Pattern Implementation

<!-- REF: https://threedots.tech/post/repository-pattern-in-go/ -->

1. **Interface definition best practices**:

   - Keep interfaces focused on domain needs
   - Use method names that reflect domain language

   ```go
   // Good: Domain-focused repository interface
   type UserRepository interface {
       GetByID(ctx context.Context, id string) (*User, error)
       FindByEmail(ctx context.Context, email string) (*User, error)
       Save(ctx context.Context, user *User) error
   }

   // Bad: Database-focused repository interface
   type UserRepository interface {
       Select(ctx context.Context, query string, args ...interface{}) (*User, error)
       Insert(ctx context.Context, user *User) error
       Update(ctx context.Context, user *User) error
   }
   ```

2. **Methods to abstract complex queries**:

   ```go
   type TrainingRepository interface {
       // ...existing methods...

       // Method for a specific business query
       FindAvailableForTrainer(
           ctx context.Context,
           trainerID string,
           from time.Time,
           to time.Time,
       ) ([]*Training, error)
   }
   ```

3. **Transaction handling**:

   ```go
   // Transaction interface
   type TxBeginner interface {
       BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
   }

   // Repository with transaction support
   type UserRepository interface {
       WithTx(tx *sql.Tx) UserRepository
       // ...other methods
   }

   // Implementation
   type SQLUserRepository struct {
       db DBTX // interface that both *sql.DB and *sql.Tx satisfy
   }

   func (r *SQLUserRepository) WithTx(tx *sql.Tx) UserRepository {
       return &SQLUserRepository{db: tx}
   }

   // Usage in application service
   func (s *UserService) TransferSubscription(fromID, toID string) error {
       tx, err := s.db.BeginTx(ctx, nil)
       if err != nil {
           return err
       }
       defer tx.Rollback()

       userRepo := s.userRepo.WithTx(tx)

       from, err := userRepo.GetByID(ctx, fromID)
       if err != nil {
           return err
       }

       to, err := userRepo.GetByID(ctx, toID)
       if err != nil {
           return err
       }

       sub := from.RemoveSubscription()
       to.AddSubscription(sub)

       if err := userRepo.Save(ctx, from); err != nil {
           return err
       }
       if err := userRepo.Save(ctx, to); err != nil {
           return err
       }

       return tx.Commit()
   }
   ```

## DDD Lite: A Pragmatic Approach

Following the approach from the "Modern Business Software in Go" series, we adopt
a pragmatic version of DDD that focuses on the most valuable aspects without
excessive complexity:

1. **Focus on bounded contexts first**:

   - Clearly define boundaries between different parts of the system
   - Use strategic DDD to identify subdomains before diving into implementation details

2. **Keep aggregates small and focused**:

   - Group related entities that should be consistent together
   - Prefer multiple small aggregates over a single large one
   - Ensure aggregates have clear root entities that control access

3. **Use value objects extensively**:

   - Encapsulate validation and business rules in value objects
   - Use value objects for concepts with no identity (like Email, Address, Money)
   - Make value objects immutable to prevent unexpected state changes

4. **Apply tactical patterns selectively**:
   - Use repositories, domain services, and entities where they provide clear value
   - Don't force patterns where they don't fit naturally

## Combining CQRS, DDD, and Clean Architecture

Based on the "Modern Business Software in Go" series, here's how to effectively combine these patterns:

1. **Separate read and write operations**:

   - Commands change state but don't return data
   - Queries return data but don't change state
   - Use separate models for reads and writes when beneficial

2. **Structure application layers**:

   - Domain Layer: Entities, value objects, domain services (independent of infrastructure)
   - Application Layer: Command/query handlers that orchestrate domain objects
   - Infrastructure Layer: Repositories, external services implementations, database access
   - Interface Layer: REST, gRPC, or GraphQL handlers that translate external requests

3. **Follow dependency rule**:
   - Inner layers should not depend on outer layers
   - Use dependency injection and ports/adapters to keep the domain pure

```go
// Application layer: Command handler
type CreateTrainingHandler struct {
    trainingRepo domain.TrainingRepository // dependency on interface, not implementation
    eventBus     domain.EventBus
}

func (h *CreateTrainingHandler) Handle(ctx context.Context, cmd CreateTraining) error {
    // Validate and prepare domain objects
    training, err := domain.NewTraining(
        cmd.Title,
        domain.TrainingType(cmd.Type),
        cmd.Description,
    )
    if err != nil {
        return err
    }

    // Use domain services or entities for domain logic
    if err := training.ScheduleAt(cmd.Time, cmd.Duration); err != nil {
        return err
    }

    // Persist changes via repository
    if err := h.trainingRepo.Save(ctx, training); err != nil {
        return err
    }

    // Publish domain events
    h.eventBus.Publish(training.Events())

    return nil
}
```

## Repository Pattern Best Practices

<!-- REF: https://threedots.tech/post/repository-pattern-in-go/ -->

The repository pattern is crucial for separating domain logic from data access:

1. **Design repositories for domain needs**:

   - Focus on domain operations, not database operations
   - Abstract away data access technology details
   - Use domain language in repository method names

2. **Keep repositories secure by design**:

   - Use strongly typed parameters for queries to prevent SQL injection
   - Avoid using raw query strings in repositories
   - Consider using query builders for complex queries
   - **Enforce authorization at the repository level**: Repositories should require
     explicit authorization information as parameters, not via context, to prevent
     accidental privilege escalation or misuse.<!-- REF: https://threedots.tech/post/repository-secure-by-design/ -->
   - **Do not pass authorization via context**: Always pass authorization or ownership
     information explicitly to repository methods. <!-- REF: https://threedots.tech/post/repository-secure-by-design/ -->
   - **Design repositories to prevent misuse**: Avoid exposing methods that allow
     bulk updates or deletions without proper checks. Do not return collections
     that can be modified directly; always provide controlled update methods.
     <!-- REF: https://threedots.tech/post/repository-secure-by-design/ -->
   - **Handle internal updates securely**: Ensure that only allowed fields are
     updated, and internal state changes are protected from external manipulation.
     <!-- REF: https://threedots.tech/post/repository-secure-by-design/ -->

   ```go
   // Bad: SQL injection vulnerability
   func (r *SQLTrainerRepository) GetByName(name string) (*Trainer, error) {
       // NEVER do this - SQL injection vulnerability
       return r.db.QueryRow("SELECT * FROM trainers WHERE name = '" + name + "'")
   }

   // Good: Use parameterized queries
   func (r *SQLTrainerRepository) GetByName(ctx context.Context, name string) (*Trainer, error) {
       return r.db.QueryRowContext(ctx, "SELECT * FROM trainers WHERE name = $1", name)
   }
   ```

1. **Test repositories thoroughly**:
   - Use integration tests with real database connections
   - Consider using testcontainers for database testing
   - Test edge cases like empty results and constraint violations

## Things to Know About DRY

<!-- REF: https://threedots.tech/post/things-to-know-about-dry/ -->
<!-- REF: https://threedots.tech/post/when-not-to-dry/ -->

While DRY (Don't Repeat Yourself) is generally good practice, understanding when
and how to apply it is crucial for maintainable code:

### When DRY Is Beneficial

1. **True duplication**: When the same logic appears in multiple places and changes together
2. **Business rules**: Centralize domain logic to ensure consistency
3. **Configuration**: Single source of truth for settings and constants

```go
// Good: Centralized validation logic
func ValidateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return ErrInvalidEmail
    }
    return nil
}

// Use consistently across the domain
func (u *User) SetEmail(email string) error {
    if err := ValidateEmail(email); err != nil {
        return err
    }
    u.Email = email
    return nil
}
```

### When to Avoid DRY

1. **Different bounded contexts**: Code that looks similar but serves different domains

```go
// Don't DRY across bounded contexts
// User registration validation (Identity context)
func ValidateUserRegistration(user *User) error {
    if len(user.Password) < 8 {
        return ErrPasswordTooShort
    }
    return nil
}

// Training enrollment validation (Training context)
func ValidateTrainingEnrollment(enrollment *Enrollment) error {
    if len(enrollment.Notes) < 8 { // Different business rule, same check
        return ErrNotesTooShort
    }
    return nil
}
```

1. **Accidental duplication**: Similar code that may evolve differently
2. **Over-abstraction**: When the abstraction is more complex than the duplication

### DRY Anti-patterns to Avoid

1. **Premature abstraction**: Don't abstract until you have at least 3 similar cases
2. **Wrong abstraction**: Forcing unrelated code into shared functions
3. **God functions**: Functions that do too many things to avoid duplication

```go
// Bad: Over-abstracted function trying to handle everything
func ProcessEntity(entityType string, data interface{}, action string) error {
    switch entityType {
    case "user":
        // User-specific logic
    case "training":
        // Training-specific logic
    }
    // This becomes unmaintainable
}

// Good: Separate, focused functions
func ProcessUser(user *User, action UserAction) error {
    // User-specific logic
}

func ProcessTraining(training *Training, action TrainingAction) error {
    // Training-specific logic
}
```

### Practical Guidelines

- **Rule of three**: Don't abstract until you see the pattern three times
- **Bounded context boundaries**: Don't DRY across different domains
- **Test the abstraction**: If it's hard to test, it's probably wrong
- **Consider evolution**: Will these pieces of code change together or separately?

## Microservices Testing Strategy

<!-- REF: https://threedots.tech/post/microservices-test-architecture/ -->

Based on the "Microservices test architecture" article:

1. **Test pyramid approach**:

   - Unit tests for business logic (fast, isolated)
   - Integration tests for repositories and external services
   - Contract tests between services instead of full end-to-end tests
   - Limited end-to-end tests for critical paths only

2. **Make tests independent**:
   - Each test should create its own data and clean up after itself
   - Use database transactions to isolate test data
   - Use containers for integration tests

## Additional Considerations

<!-- REF: https://threedots.tech/post/event-sourcing-in-go/ -->

1. **Event Sourcing**:

   - Consider using event sourcing for complex domains with audit requirements
   - Store the sequence of events rather than just the current state

2. **Read Models**:

   - Optimize read operations with denormalized models
   - Update read models asynchronously via domain events

3. **Performance**:
   - Profile and optimize critical paths
   - Consider caching for frequently accessed data

## Conclusion

These guidelines aim to provide a balance between architecture purity and practical
implementation. The goal is to create maintainable, testable code that accurately
models our business domain while leveraging Go's strengths.

Remember that these patterns should serve our needs, not constrain us. Adapt them
as necessary for your specific context.
