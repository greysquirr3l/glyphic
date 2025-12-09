---
applyTo: '**/*.go'
description: Go development standards and patterns for enterprise applications
---

# Go Development Standards

## Code Structure and Organization

### Package Organization
- Use meaningful package names that describe functionality
- Keep packages focused and cohesive
- Avoid circular dependencies
- Place main packages in `cmd/` directory
- Internal packages in `internal/` are not importable from outside

### Naming Conventions
```go
// ✅ Exported types and functions use PascalCase
type UserRepository interface {
    GetUserByID(ctx context.Context, id string) (*User, error)
    SaveUser(ctx context.Context, user *User) error
}

// ✅ Unexported types and variables use camelCase
type userService struct {
    repo UserRepository
    logger *slog.Logger
}

// ✅ Constants use appropriate casing
const (
    DefaultTimeout = 30 * time.Second
    maxRetryAttempts = 3
)

// ✅ Error types end with 'Error'
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error in field %s: %s", e.Field, e.Message)
}
```

## Error Handling Best Practices

### Always Handle Errors Explicitly
```go
// ✅ Proper error handling
func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Validate input
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // Check if user exists
    existing, err := s.repo.GetUserByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, ErrUserNotFound) {
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    if existing != nil {
        return nil, ErrUserAlreadyExists
    }
    
    // Create user
    user := &User{
        ID:        generateID(),
        Email:     req.Email,
        Name:      req.Name,
        CreatedAt: time.Now(),
    }
    
    if err := s.repo.SaveUser(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    return user, nil
}

// ❌ Never ignore errors
func badExample() {
    data, _ := ioutil.ReadFile("config.json") // Don't do this
    // Handle the error!
}
```

### Use Sentinel Errors for Expected Conditions
```go
// ✅ Define sentinel errors
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
)

// ✅ Use errors.Is for comparison
if errors.Is(err, ErrUserNotFound) {
    // Handle user not found case
}
```

## Context Usage Patterns

### Always Use Context for Cancellation and Timeouts
```go
// ✅ Pass context through function calls
func (s *service) ProcessData(ctx context.Context, data []Item) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    for _, item := range data {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := s.processItem(ctx, item); err != nil {
                return fmt.Errorf("processing item %s: %w", item.ID, err)
            }
        }
    }
    return nil
}

// ✅ Use context in database operations
func (r *repository) GetUser(ctx context.Context, id string) (*User, error) {
    query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
    
    var user User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Name, &user.Email, &user.CreatedAt,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to query user: %w", err)
    }
    
    return &user, nil
}
```

## Interface Design Principles

### Keep Interfaces Small and Focused
```go
// ✅ Small, focused interfaces
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

// ✅ Compose interfaces when needed
type ReadWriter interface {
    Reader
    Writer
}

// ✅ Define interfaces where they're used, not where they're implemented
type UserService struct {
    repo UserRepository  // Interface defined in this package
}

type UserRepository interface {
    GetUser(ctx context.Context, id string) (*User, error)
    SaveUser(ctx context.Context, user *User) error
}
```

## Testing Patterns

### Use Table-Driven Tests
```go
// ✅ Table-driven test pattern
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {
            name:    "valid email",
            email:   "user@protonmail.com",
            wantErr: false,
        },
        {
            name:    "empty email",
            email:   "",
            wantErr: true,
        },
        {
            name:    "invalid format",
            email:   "notanemail",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Use Testify for Assertions and Mocks
```go
// ✅ Use testify for clean assertions
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func TestUserService_CreateUser(t *testing.T) {
    // Setup
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    
    ctx := context.Background()
    req := CreateUserRequest{
        Email: "test@protonmail.com",
        Name:  "Test User",
    }
    
    // Mock expectations
    mockRepo.On("GetUserByEmail", ctx, req.Email).Return(nil, ErrUserNotFound)
    mockRepo.On("SaveUser", ctx, mock.AnythingOfType("*User")).Return(nil)
    
    // Execute
    user, err := service.CreateUser(ctx, req)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, req.Email, user.Email)
    assert.Equal(t, req.Name, user.Name)
    assert.NotEmpty(t, user.ID)
    
    // Verify mocks
    mockRepo.AssertExpectations(t)
}
```

## Concurrency Patterns

### Use Channels for Communication
```go
// ✅ Producer-consumer pattern with channels
func ProcessFiles(ctx context.Context, filenames []string) error {
    const workers = 3
    jobs := make(chan string, len(filenames))
    results := make(chan error, len(filenames))
    
    // Start workers
    for i := 0; i < workers; i++ {
        go worker(ctx, jobs, results)
    }
    
    // Send jobs
    for _, filename := range filenames {
        jobs <- filename
    }
    close(jobs)
    
    // Collect results
    for i := 0; i < len(filenames); i++ {
        if err := <-results; err != nil {
            return fmt.Errorf("processing failed: %w", err)
        }
    }
    
    return nil
}

func worker(ctx context.Context, jobs <-chan string, results chan<- error) {
    for filename := range jobs {
        select {
        case <-ctx.Done():
            results <- ctx.Err()
            return
        default:
            results <- processFile(filename)
        }
    }
}
```

### Use sync Package for Synchronization
```go
// ✅ Use sync.WaitGroup for coordinating goroutines
func ProcessConcurrently(ctx context.Context, items []Item) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(items))
    
    for _, item := range items {
        wg.Add(1)
        go func(item Item) {
            defer wg.Done()
            if err := processItem(ctx, item); err != nil {
                errCh <- fmt.Errorf("processing item %s: %w", item.ID, err)
            }
        }(item)
    }
    
    // Wait for all goroutines to complete
    go func() {
        wg.Wait()
        close(errCh)
    }()
    
    // Check for errors
    for err := range errCh {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

## Performance Best Practices

### Minimize Allocations in Hot Paths
```go
// ✅ Reuse slices and maps when possible
type Processor struct {
    buffer []byte
    cache  map[string]interface{}
}

func (p *Processor) Process(data []byte) []byte {
    // Reuse buffer, growing if necessary
    if cap(p.buffer) < len(data)*2 {
        p.buffer = make([]byte, len(data)*2)
    }
    p.buffer = p.buffer[:0] // Reset length but keep capacity
    
    // Process data into buffer
    // ...
    
    return p.buffer
}

// ✅ Use string builder for string concatenation
func BuildMessage(parts []string) string {
    var builder strings.Builder
    builder.Grow(estimateSize(parts)) // Pre-allocate capacity
    
    for _, part := range parts {
        builder.WriteString(part)
    }
    
    return builder.String()
}
```

## Documentation Requirements

### Document Exported Functions and Types
```go
// ✅ Comprehensive package documentation
// Package user provides user management functionality including
// authentication, authorization, and user profile management.
//
// Example usage:
//
//     service := user.NewService(repo, logger)
//     user, err := service.CreateUser(ctx, createUserRequest)
//     if err != nil {
//         log.Printf("failed to create user: %v", err)
//     }
package user

// User represents a user in the system with authentication and profile information.
type User struct {
    // ID is the unique identifier for the user
    ID string `json:"id" db:"id"`
    
    // Email is the user's email address, used for authentication
    Email string `json:"email" db:"email"`
    
    // Name is the user's display name
    Name string `json:"name" db:"name"`
    
    // CreatedAt is when the user account was created
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateUser creates a new user with the provided information.
// It validates the input, checks for existing users with the same email,
// and returns an error if the user cannot be created.
//
// Returns ErrUserAlreadyExists if a user with the same email exists.
// Returns a ValidationError if the input is invalid.
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Implementation...
}
```

## Structured Logging with slog

### Basic Logging Patterns
#### slog Usage

```go
import "log/slog"

// ✅ Use structured logging with key-value pairs
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    slog.Info("Creating user", 
        slog.String("email", req.Email),
        slog.String("name", req.Name))
    
    user, err := s.repository.Create(ctx, req)
    if err != nil {
        slog.Error("Failed to create user", 
            slog.String("email", req.Email),
            slog.String("error", err.Error()))
        return nil, err
    }
    
    slog.Info("User created successfully", 
        slog.String("user_id", user.ID),
        slog.String("email", user.Email))
    
    return user, nil
}

// ✅ Context-aware logging with request tracking
func (s *Service) ProcessRequest(ctx context.Context, req Request) error {
    requestID := extractRequestID(ctx)
    
    logger := slog.With(
        slog.String("request_id", requestID),
        slog.String("operation", "process_request"))
    
    logger.Info("Processing request started")
    
    start := time.Now()
    err := s.process(ctx, req)
    duration := time.Since(start)
    
    if err != nil {
        logger.Error("Request processing failed", 
            slog.Duration("duration", duration),
            slog.String("error", err.Error()))
        return err
    }
    
    logger.Info("Request processing completed", 
        slog.Duration("duration", duration))
    
    return nil
}

// ✅ Component-specific logging
func (r *UserRepository) GetUser(ctx context.Context, id string) (*User, error) {
    logger := slog.With(slog.String("component", "repository"))
    
    logger.Debug("Fetching user from database", slog.String("user_id", id))
    
    // Database operation...
    
    logger.Info("User retrieved successfully", slog.String("user_id", id))
    return user, nil
}

// ✅ Performance and audit logging
func (s *Service) PerformSensitiveOperation(ctx context.Context, userID string) error {
    // Performance logging
    start := time.Now()
    defer func() {
        slog.Info("Performance metric",
            slog.String("operation", "sensitive_operation"),
            slog.Duration("duration", time.Since(start)))
    }()
    
    // Audit logging
    slog.Info("Audit event",
        slog.String("action", "SENSITIVE_OPERATION"),
        slog.String("user_id", userID),
        slog.String("timestamp", time.Now().Format(time.RFC3339)))
    
    // Operation implementation...
    
    return nil
}

// ✅ Error logging with context
func (s *Service) handleError(ctx context.Context, err error, operation string) error {
    slog.Error("Operation failed",
        slog.String("operation", operation),
        slog.String("error", err.Error()),
        slog.String("request_id", extractRequestID(ctx)))
    
    return fmt.Errorf("%s failed: %w", operation, err)
}
```
#### logrus Usage

```go
import "github.com/sirupsen/logrus"
// ✅ Use structured logging with key-value pairs
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    logrus.WithFields(logrus.Fields{
        "email": req.Email,
        "name":  req.Name,
    }).Info("Creating user")
}
```

## Security Best Practices

### Input Validation and Sanitization
```go
// ✅ Always validate inputs
func (r CreateUserRequest) Validate() error {
    var errors []string
    
    if r.Email == "" {
        errors = append(errors, "email is required")
    } else if !isValidEmail(r.Email) {
        errors = append(errors, "email format is invalid")
    }
    
    if r.Name == "" {
        errors = append(errors, "name is required")
    } else if len(r.Name) > 100 {
        errors = append(errors, "name must be less than 100 characters")
    }
    
    if len(errors) > 0 {
        return ValidationError{
            Field:   "request",
            Message: strings.Join(errors, "; "),
        }
    }
    
    return nil
}

// ✅ Use parameterized queries to prevent SQL injection
func (r *repository) UpdateUser(ctx context.Context, user *User) error {
    query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
    _, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
    return err
}
```

Always follow these patterns when generating Go code to ensure consistency, reliability, and maintainability across the codebase.