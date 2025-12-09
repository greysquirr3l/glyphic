# Security-First Repository Design in Go

This document provides comprehensive guidance for implementing security-firs
repository patterns in Go applications, following Clean Architecture and DDD
principles. It focuses on creating repositories that are secure by design, making
it impossible to use them incorrectly.

## Table of Contents

- [Core Security Principles](#core-security-principles)
- [Domain-Driven Security](#domain-driven-security)
- [Implementation Patterns](#implementation-patterns)
- [Collection Security](#collection-security)
- [Internal Operations](#internal-operations)
- [Anti-Patterns to Avoid](#anti-patterns-to-avoid)
- [Testing Security](#testing-security)
- [Advanced Patterns](#advanced-patterns)

## Core Security Principles

### Security by Design

The fundamental principle of secure repository design is to make it **impossible**
**to use the repository incorrectly**. Rather than relying on developers to remember
security checks, we build security directly into the repository interface.

**Key Concepts:**

- **Fail Secure**: Default to denying access rather than allowing it
- **Explicit Authorization**: Make authorization requirements visible in method signatures
- **Domain-Driven Security**: Security rules should reflect business domain concepts
- **Static Type Safety**: Use Go's type system to enforce security at compile time

### Authorization at the Repository Level

While some may debate whether authorization belongs in the repository, for
business applications it provides significant advantages:

- **Single Point of Control**: All data access goes through secured repositories
- **Cannot Be Bypassed**: No way to accidentally skip authorization checks
- **Domain Alignment**: Security rules match business concepts
- **Testable**: Easy to unit test authorization logic

## Domain-Driven Security

### Security as Domain Logic

Security rules are often core business logic that should be modeled in the domain
layer. For example, "only training owners and trainers can see training details"
is a business rule, not a technical detail.

```go
package domain

import (
    "fmt"
)

// User represents the authenticated user context
type User struct {
    uuid     string
    userType UserType
}

type UserType string

const (
    AttendeeUserType UserType = "attendee"
    TrainerUserType  UserType = "trainer"
    AdminUserType    UserType = "admin"
)

func (u User) UUID() string {
    return u.uuid
}

func (u User) Type() UserType {
    return u.userType
}

// Security error types - part of domain
type ForbiddenToSeeTrainingError struct {
    RequestingUserUUID string
    TrainingOwnerUUID  string
}

func (f ForbiddenToSeeTrainingError) Error() string {
    return fmt.Sprintf(
        "user '%s' cannot access training owned by '%s'",
        f.RequestingUserUUID,
        f.TrainingOwnerUUID,
    )
}

// Domain security rule
func CanUserSeeTraining(user User, training Training) error {
    // Trainers can see all trainings
    if user.Type() == TrainerUserType {
        return nil
    }
    
    // Users can see their own trainings
    if user.UUID() == training.UserUUID() {
        return nil
    }
    
    // Admins can see all trainings
    if user.Type() == AdminUserType {
        return nil
    }
    
    return ForbiddenToSeeTrainingError{
        RequestingUserUUID: user.UUID(),
        TrainingOwnerUUID:  training.UserUUID(),
    }
}

// Additional authorization functions
func CanUserModifyTraining(user User, training Training) error {
    // Only owners and trainers can modify trainings
    if user.Type() == TrainerUserType || user.UUID() == training.UserUUID() {
        return nil
    }
    
    return ForbiddenToModifyTrainingError{
        RequestingUserUUID: user.UUID(),
        TrainingOwnerUUID:  training.UserUUID(),
    }
}
```

### Modeling User Context

Create a domain `User` type that encapsulates authentication context rather than using primitive types:

```go
// Good: Type-safe user context
type User struct {
    uuid        string
    userType    UserType
    permissions []Permission
}

// Avoid: Primitive obsession
func GetTraining(ctx context.Context, trainingID string, userID string) (*Training, error)

// Better: Domain-driven user context  
func GetTraining(ctx context.Context, trainingID string, user User) (*Training, error)
```

## Implementation Patterns

### Secure Repository Interface

Design repository interfaces that require user context for all operations that need authorization:

```go
package domain

// TrainingRepository defines the interface in the domain layer
type TrainingRepository interface {
    // Operations requiring authorization
    GetTraining(ctx context.Context, trainingID string, user User) (*Training, error)
    UpdateTraining(ctx context.Context, trainingID string, user User, updateFn func(*Training) (*Training, error)) error
    DeleteTraining(ctx context.Context, trainingID string, user User) error
    
    // Collection operations with built-in filtering
    FindTrainingsForUser(ctx context.Context, userUUID string) ([]*Training, error)
    FindTrainingsForTrainer(ctx context.Context, trainerUUID string) ([]*Training, error)
    
    // Internal operations (clearly marked)
    GetTrainingInternal(ctx context.Context, trainingID string) (*Training, error)
    UpdateTrainingInternal(ctx context.Context, trainingID string, updateFn func(*Training) (*Training, error)) error
}
```

### Repository Implementation

Implement the repository with authorization checks integrated into every method:

```go
package infrastructure

import (
    "context"
    "myapp/internal/domain"
)

type PostgresTrainingRepository struct {
    db *sql.DB
}

func (r *PostgresTrainingRepository) GetTraining(
    ctx context.Context,
    trainingID string,
    user domain.User,
) (*domain.Training, error) {
    // First, retrieve the training
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

func (r *PostgresTrainingRepository) UpdateTraining(
    ctx context.Context,
    trainingID string,
    user domain.User,
    updateFn func(*domain.Training) (*domain.Training, error),
) error {
    return r.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
        // Get current training
        training, err := r.getTrainingByIDTx(ctx, tx, trainingID)
        if err != nil {
            return err
        }
        
        // Check authorization before allowing updates
        if err := domain.CanUserModifyTraining(user, *training); err != nil {
            return err
        }
        
        // Apply the update function
        updatedTraining, err := updateFn(training)
        if err != nil {
            return err
        }
        
        // Save the updated training
        return r.saveTrainingTx(ctx, tx, updatedTraining)
    })
}

// Private helper methods
func (r *PostgresTrainingRepository) getTrainingByID(ctx context.Context, id string) (*domain.Training, error) {
    // Implementation details...
    return nil, nil
}
```

### Application Layer Integration

The application layer uses the secure repository without additional security checks:

```go
package application

type ApproveTrainingHandler struct {
    trainingRepo domain.TrainingRepository
}

func (h *ApproveTrainingHandler) Handle(ctx context.Context, cmd ApproveTrainingCommand) error {
    // No additional security checks needed - repository handles it
    return h.trainingRepo.UpdateTraining(
        ctx,
        cmd.TrainingID,
        cmd.User, // User context passed through
        func(training *domain.Training) (*domain.Training, error) {
            // Apply business logic
            if err := training.Approve(cmd.User.Type()); err != nil {
                return nil, err
            }
            return training, nil
        },
    )
}
```

## Collection Security

### Query-Level Security

For collection operations, implement security at the query level rather than filtering results in application code:

```go
func (r *PostgresTrainingRepository) FindTrainingsForUser(
    ctx context.Context,
    userUUID string,
) ([]*domain.Training, error) {
    query := `
        SELECT t.id, t.title, t.user_uuid, t.trainer_uuid, t.scheduled_at
        FROM trainings t
        WHERE t.user_uuid = $1
           OR t.trainer_uuid = $1
        AND t.canceled = false
        AND t.scheduled_at >= NOW() - INTERVAL '24 hours'
        ORDER BY t.scheduled_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, userUUID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    return r.scanTrainings(rows)
}

// Separate method for trainers to see all trainings they're assigned to
func (r *PostgresTrainingRepository) FindTrainingsForTrainer(
    ctx context.Context,
    trainerUUID string,
) ([]*domain.Training, error) {
    query := `
        SELECT t.id, t.title, t.user_uuid, t.trainer_uuid, t.scheduled_at
        FROM trainings t
        WHERE t.trainer_uuid = $1
        AND t.canceled = false
        ORDER BY t.scheduled_at ASC
    `
    
    rows, err := r.db.QueryContext(ctx, query, trainerUUID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    return r.scanTrainings(rows)
}
```

### Role-Based Collection Access

Create specific methods for different user roles rather than generic "get all" methods:

```go
// Good: Role-specific methods
func (r *TrainingRepository) FindTrainingsForUser(ctx context.Context, userUUID string) ([]*Training, error)
func (r *TrainingRepository) FindTrainingsForTrainer(ctx context.Context, trainerUUID string) ([]*Training, error)
func (r *TrainingRepository) FindAllTrainingsForAdmin(ctx context.Context, filters AdminFilters) ([]*Training, error)

// Avoid: Generic method that requires external authorization
func (r *TrainingRepository) FindAllTrainings(ctx context.Context, user User) ([]*Training, error)
```

## Internal Operations

### Administrative Operations

For internal operations (migrations, admin functions, background jobs), create
explicitly named methods that make their purpose clear:

```go
// Clearly marked internal operations
type TrainingRepository interface {
    // Regular user operations
    GetTraining(ctx context.Context, trainingID string, user User) (*Training, error)
    
    // Internal operations - naming makes security implications clear
    GetTrainingInternal(ctx context.Context, trainingID string) (*Training, error)
    UpdateTrainingForMigration(ctx context.Context, trainingID string, updateFn func(*Training) (*Training, error)) error
    DeleteTrainingForCleanup(ctx context.Context, trainingID string) error
    
    // System operations with explicit naming
    GetTrainingForNotification(ctx context.Context, trainingID string) (*Training, error)
    UpdateTrainingBySystem(ctx context.Context, trainingID string, reason string, updateFn func(*Training) (*Training, error)) error
}
```

### Background Job Security

For background operations, create specific command types and handlers:

```go
type SystemUpdateTrainingCommand struct {
    TrainingID string
    Reason     string
    UpdateFn   func(*Training) (*Training, error)
}

type SystemUpdateTrainingHandler struct {
    repo TrainingRepository
}

func (h *SystemUpdateTrainingHandler) Handle(ctx context.Context, cmd SystemUpdateTrainingCommand) error {
    // Use internal repository method
    return h.repo.UpdateTrainingBySystem(ctx, cmd.TrainingID, cmd.Reason, cmd.UpdateFn)
}
```

## Anti-Patterns to Avoid

### 1. Context-Based Authentication

**Avoid** passing authentication details via `context.Context`:

```go
// Bad: Hidden authentication requirements
func GetTraining(ctx context.Context, trainingID string) (*Training, error) {
    user := ctx.Value("user").(User) // Loses type safety
    // ...
}

// Good: Explicit authentication requirements
func GetTraining(ctx context.Context, trainingID string, user User) (*Training, error) {
    // ...
}
```

### 2. External Authorization Checks

**Avoid** relying on callers to perform authorization:

```go
// Bad: Authorization responsibility placed on caller
if !CanUserSeeTraining(user, training) {
    return nil, ErrForbidden
}
training, err := repo.GetTraining(ctx, trainingID)

// Good: Repository handles authorization internally
training, err := repo.GetTraining(ctx, trainingID, user)
```

### 3. Fake Users and Backdoors

**Avoid** creating fake users or bypassing authorization:

```go
// Bad: Creates maintenance and audit issues
adminUser := User{UUID: "system", Type: AdminUserType}
training, err := repo.GetTraining(ctx, trainingID, adminUser)

// Good: Explicit internal methods
training, err := repo.GetTrainingInternal(ctx, trainingID)
```

### 4. Generic Repository Methods

**Avoid** overly generic methods that push authorization logic to callers:

```go
// Bad: Forces callers to implement authorization
func (r *Repository) FindAll(ctx context.Context, filters map[string]interface{}) ([]*Training, error)

// Good: Purpose-built methods with built-in security
func (r *Repository) FindTrainingsForUser(ctx context.Context, userUUID string) ([]*Training, error)
func (r *Repository) FindTrainingsForTrainer(ctx context.Context, trainerUUID string) ([]*Training, error)
```

## Testing Security

### Unit Testing Authorization Logic

Test domain authorization rules independently:

```go
func TestCanUserSeeTraining(t *testing.T) {
    tests := []struct {
        name        string
        user        domain.User
        training    domain.Training
        expectError bool
    }{
        {
            name:        "trainer can see any training",
            user:        domain.NewUser("trainer1", domain.TrainerUserType),
            training:    domain.NewTraining("user1", "trainer2"),
            expectError: false,
        },
        {
            name:        "user can see own training",
            user:        domain.NewUser("user1", domain.AttendeeUserType),
            training:    domain.NewTraining("user1", "trainer1"),
            expectError: false,
        },
        {
            name:        "user cannot see other's training",
            user:        domain.NewUser("user1", domain.AttendeeUserType),
            training:    domain.NewTraining("user2", "trainer1"),
            expectError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := domain.CanUserSeeTraining(tt.user, tt.training)
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Integration Testing Repository Security

Test that repositories properly enforce authorization:

```go
func TestTrainingRepository_GetTraining_Authorization(t *testing.T) {
    db := setupTestDB(t)
    repo := infrastructure.NewPostgresTrainingRepository(db)
    
    // Setup test data
    training := createTestTraining(t, db, "user1", "trainer1")
    
    t.Run("user can access own training", func(t *testing.T) {
        user := domain.NewUser("user1", domain.AttendeeUserType)
        result, err := repo.GetTraining(context.Background(), training.ID, user)
        
        assert.NoError(t, err)
        assert.Equal(t, training.ID, result.ID)
    })
    
    t.Run("user cannot access other's training", func(t *testing.T) {
        user := domain.NewUser("user2", domain.AttendeeUserType)
        _, err := repo.GetTraining(context.Background(), training.ID, user)
        
        assert.Error(t, err)
        assert.IsType(t, domain.ForbiddenToSeeTrainingError{}, err)
    })
    
    t.Run("trainer can access any training", func(t *testing.T) {
        trainer := domain.NewUser("trainer2", domain.TrainerUserType)
        result, err := repo.GetTraining(context.Background(), training.ID, trainer)
        
        assert.NoError(t, err)
        assert.Equal(t, training.ID, result.ID)
    })
}
```

### Security Regression Testing

Create tests that verify security cannot be bypassed:

```go
func TestRepositorySecurityRegression(t *testing.T) {
    t.Run("cannot access training without user context", func(t *testing.T) {
        // This test ensures we don't accidentally create methods that bypass authorization
        db := setupTestDB(t)
        repo := infrastructure.NewPostgresTrainingRepository(db)
        
        // Verify no public method exists that bypasses authorization
        trainingType := reflect.TypeOf(repo)
        for i := 0; i < trainingType.NumMethod(); i++ {
            method := trainingType.Method(i)
            
            // Any method that returns Training should require User parameter
            if strings.Contains(method.Type.String(), "Training") &&
               !strings.Contains(method.Type.String(), "User") &&
               !strings.Contains(method.Name, "Internal") {
                t.Errorf("Method %s might bypass authorization", method.Name)
            }
        }
    })
}
```

## Advanced Patterns

### Row-Level Security (RLS)

For PostgreSQL, consider implementing Row-Level Security policies:

```sql
-- Enable RLS on the trainings table
ALTER TABLE trainings ENABLE ROW LEVEL SECURITY;

-- Policy for users to see only their own trainings
CREATE POLICY user_trainings_policy ON trainings
    FOR ALL TO app_user
    USING (user_uuid = current_setting('app.user_uuid'));

-- Policy for trainers to see all trainings
CREATE POLICY trainer_trainings_policy ON trainings
    FOR ALL TO app_trainer
    USING (true);
```

Then set the user context in your application:

```go
func (r *PostgresTrainingRepository) GetTraining(
    ctx context.Context,
    trainingID string,
    user domain.User,
) (*domain.Training, error) {
    // Set user context for RLS
    if _, err := r.db.ExecContext(ctx, "SET app.user_uuid = $1", user.UUID()); err != nil {
        return nil, err
    }
    
    // Query will automatically filter based on RLS policies
    return r.getTrainingByID(ctx, trainingID)
}
```

### Audit Logging

Implement comprehensive audit logging for security-sensitive operations:

```go
type AuditLogger interface {
    LogAccess(ctx context.Context, user User, resource string, action string, result string)
}

func (r *PostgresTrainingRepository) GetTraining(
    ctx context.Context,
    trainingID string,
    user domain.User,
) (*domain.Training, error) {
    defer func() {
        r.auditLogger.LogAccess(ctx, user, "training:"+trainingID, "read", "success")
    }()
    
    training, err := r.getTrainingByID(ctx, trainingID)
    if err != nil {
        r.auditLogger.LogAccess(ctx, user, "training:"+trainingID, "read", "error")
        return nil, err
    }
    
    if err := domain.CanUserSeeTraining(user, *training); err != nil {
        r.auditLogger.LogAccess(ctx, user, "training:"+trainingID, "read", "forbidden")
        return nil, err
    }
    
    return training, nil
}
```

## Conclusion

Security-first repository design creates systems that are inherently secure and
impossible to use incorrectly. By following these patterns:

1. **Build authorization into repository interfaces** - Make it impossible to forget security checks
2. **Model security as domain logic** - Security rules should reflect business requirements
3. **Use explicit naming for internal operations** - Make security implications clear
4. **Avoid context-based authentication** - Maintain type safety and explicit dependencies
5. **Test security thoroughly** - Include both positive and negative security test cases

This approach results in repositories that are secure by design, making your
application more robust and your team more confident when making changes.

## References

- [Repository secure by design - Three Dots Labs](https://threedots.tech/post/repository-secure-by-design/)
- [DDD Lite in Go - Three Dots Labs](https://threedots.tech/post/ddd-lite-in-go-introduction/)
- [The Repository pattern in Go - Three Dots Labs](https://threedots.tech/post/repository-pattern-in-go/)
- [OSSF Security Baselines](https://github.com/ossf/security-baselines)
