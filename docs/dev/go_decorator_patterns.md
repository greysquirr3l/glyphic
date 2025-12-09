# Go Decorator Pattern: Clean Architecture Implementation

This document provides comprehensive guidance for implementing the Decorator pattern
in Go applications, following Clean Architecture, Domain-Driven Design (DDD), and
CQRS principles. The Decorator pattern allows adding new behaviors to objects
dynamically by placing them inside special wrapper objects.

## Table of Contents

- [Pattern Overview](#pattern-overview)
- [Core Concepts](#core-concepts)
- [Implementation Patterns](#implementation-patterns)
- [Clean Architecture Integration](#clean-architecture-integration)
- [Domain-Driven Design Applications](#domain-driven-design-applications)
- [CQRS Enhancement Patterns](#cqrs-enhancement-patterns)
- [Advanced Decorator Techniques](#advanced-decorator-techniques)
- [Testing Decorated Components](#testing-decorated-components)
- [Performance Considerations](#performance-considerations)
- [Common Pitfalls and Solutions](#common-pitfalls-and-solutions)
- [Best Practices](#best-practices)

---

## Pattern Overview

### What is the Decorator Pattern?

The Decorator pattern is a structural design pattern that allows adding new behaviors
to objects dynamically by placing them inside special wrapper objects, called decorators.
Using decorators you can wrap objects countless number of times since both target objects
and decorators follow the same interface. The resulting object will get a stacking behavior
of all wrappers.

### Why Use Decorators?

- **Dynamic Behavior Addition**: Add functionality at runtime without modifying existing code
- **Composition Over Inheritance**: Avoid the "class explosion" problem of inheritance hierarchies
- **Single Responsibility**: Each decorator handles one specific concern
- **Flexible Combinations**: Stack decorators in any order to achieve different behaviors
- **Open/Closed Principle**: Open for extension, closed for modification

### Decorator vs. Alternatives

| Pattern | Interface | Composition | Use Case |
|---------|-----------|-------------|----------|
| **Decorator** | Same or enhanced | Recursive | Add behavior dynamically |
| **Adapter** | Different | Single wrap | Interface compatibility |
| **Proxy** | Same | Single wrap | Control access |
| **Chain of Responsibility** | Similar structure | Linear chain | Handle requests sequentially |

---

## Core Concepts

### Basic Decorator Structure

The Decorator pattern consists of four key components:

```go
// 1. Component Interface - defines operations that can be altered by decorators
type Component interface {
    Operation() string
}

// 2. Concrete Component - provides default implementation
type ConcreteComponent struct {
    data string
}

func (c *ConcreteComponent) Operation() string {
    return c.data
}

// 3. Base Decorator - maintains reference to component and delegates to it
type BaseDecorator struct {
    component Component
}

func (d *BaseDecorator) Operation() string {
    if d.component != nil {
        return d.component.Operation()
    }
    return ""
}

// 4. Concrete Decorator - extends behavior while maintaining interface
type ConcreteDecorator struct {
    BaseDecorator
    additionalState string
}

func (d *ConcreteDecorator) Operation() string {
    return "Decorated(" + d.BaseDecorator.Operation() + ")"
}
```

### Interface-First Design

In Go, interfaces are the foundation of the Decorator pattern:

```go
// Define behavior contract
type UserService interface {
    CreateUser(ctx context.Context, user domain.User) error
    GetUser(ctx context.Context, userID string) (*domain.User, error)
    UpdateUser(ctx context.Context, user domain.User) error
}

// Core implementation
type BasicUserService struct {
    repo domain.UserRepository
}

func (s *BasicUserService) CreateUser(ctx context.Context, user domain.User) error {
    return s.repo.Save(ctx, user)
}

func (s *BasicUserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
    return s.repo.GetByID(ctx, userID)
}

func (s *BasicUserService) UpdateUser(ctx context.Context, user domain.User) error {
    return s.repo.Update(ctx, user)
}
```

---

## Implementation Patterns

### 1. Functional Decorators

Leverage Go's first-class functions for simple decorations:

```go
// Function type definition
type HandlerFunc func(ctx context.Context, req Request) (Response, error)

// Decorator function type
type Decorator func(HandlerFunc) HandlerFunc

// Logging decorator
func WithLogging(logger *log.Logger) Decorator {
    return func(next HandlerFunc) HandlerFunc {
        return func(ctx context.Context, req Request) (Response, error) {
            start := time.Now()
            logger.Printf("Starting request: %+v", req)
            
            resp, err := next(ctx, req)
            
            duration := time.Since(start)
            if err != nil {
                logger.Printf("Request failed after %v: %v", duration, err)
            } else {
                logger.Printf("Request completed after %v", duration)
            }
            
            return resp, err
        }
    }
}

// Retry decorator
func WithRetry(maxAttempts int, backoff time.Duration) Decorator {
    return func(next HandlerFunc) HandlerFunc {
        return func(ctx context.Context, req Request) (Response, error) {
            var lastErr error
            
            for attempt := 1; attempt <= maxAttempts; attempt++ {
                resp, err := next(ctx, req)
                if err == nil {
                    return resp, nil
                }
                
                lastErr = err
                if attempt < maxAttempts {
                    select {
                    case <-time.After(backoff * time.Duration(attempt)):
                    case <-ctx.Done():
                        return Response{}, ctx.Err()
                    }
                }
            }
            
            return Response{}, fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
        }
    }
}

// Composition helper
func Compose(decorators ...Decorator) Decorator {
    return func(next HandlerFunc) HandlerFunc {
        for i := len(decorators) - 1; i >= 0; i-- {
            next = decorators[i](next)
        }
        return next
    }
}

// Usage
func main() {
    baseHandler := func(ctx context.Context, req Request) (Response, error) {
        // Core business logic
        return processRequest(req)
    }
    
    // Compose multiple decorators
    decoratedHandler := Compose(
        WithLogging(logger),
        WithRetry(3, 100*time.Millisecond),
        WithTimeout(5*time.Second),
    )(baseHandler)
    
    // Use the decorated handler
    resp, err := decoratedHandler(ctx, request)
}
```

### 2. Struct-Based Decorators

For more complex state management and behavior:

```go
// Service decorator with caching
type CachedUserService struct {
    next  UserService
    cache map[string]*domain.User
    mu    sync.RWMutex
    ttl   time.Duration
}

func NewCachedUserService(next UserService, ttl time.Duration) *CachedUserService {
    return &CachedUserService{
        next:  next,
        cache: make(map[string]*domain.User),
        ttl:   ttl,
    }
}

func (s *CachedUserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
    // Check cache first
    s.mu.RLock()
    if user, exists := s.cache[userID]; exists {
        s.mu.RUnlock()
        return user, nil
    }
    s.mu.RUnlock()
    
    // Cache miss - delegate to next service
    user, err := s.next.GetUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Update cache
    s.mu.Lock()
    s.cache[userID] = user
    s.mu.Unlock()
    
    return user, nil
}

func (s *CachedUserService) CreateUser(ctx context.Context, user domain.User) error {
    err := s.next.CreateUser(ctx, user)
    if err != nil {
        return err
    }
    
    // Update cache
    s.mu.Lock()
    s.cache[user.ID] = &user
    s.mu.Unlock()
    
    return nil
}

func (s *CachedUserService) UpdateUser(ctx context.Context, user domain.User) error {
    err := s.next.UpdateUser(ctx, user)
    if err != nil {
        return err
    }
    
    // Invalidate cache
    s.mu.Lock()
    delete(s.cache, user.ID)
    s.mu.Unlock()
    
    return nil
}
```

### 3. Builder Pattern for Decorator Composition

Create fluent APIs for decorator assembly:

```go
type ServiceBuilder struct {
    service UserService
}

func NewServiceBuilder(base UserService) *ServiceBuilder {
    return &ServiceBuilder{service: base}
}

func (b *ServiceBuilder) WithCache(ttl time.Duration) *ServiceBuilder {
    b.service = NewCachedUserService(b.service, ttl)
    return b
}

func (b *ServiceBuilder) WithRateLimit(limit int, window time.Duration) *ServiceBuilder {
    b.service = NewRateLimitedUserService(b.service, limit, window)
    return b
}

func (b *ServiceBuilder) WithMetrics(collector MetricsCollector) *ServiceBuilder {
    b.service = NewMetricsUserService(b.service, collector)
    return b
}

func (b *ServiceBuilder) WithLogging(logger *log.Logger) *ServiceBuilder {
    b.service = NewLoggingUserService(b.service, logger)
    return b
}

func (b *ServiceBuilder) Build() UserService {
    return b.service
}

// Usage
userService := NewServiceBuilder(basicUserService).
    WithLogging(logger).
    WithMetrics(metricsCollector).
    WithCache(5 * time.Minute).
    WithRateLimit(100, time.Minute).
    Build()
```

---

## Clean Architecture Integration

### Domain Layer - Pure Business Logic

The domain layer should remain decorator-agnostic:

```go
// Domain entities remain unchanged
type User struct {
    ID       UserID
    Email    Email
    Username Username
    Role     UserRole
    Active   bool
}

// Domain services focus on business rules
type UserDomainService struct{}

func (s *UserDomainService) CanUserAccessResource(user User, resource Resource) error {
    if !user.Active {
        return ErrUserInactive
    }
    
    if !user.Role.HasPermission(resource.RequiredPermission()) {
        return ErrInsufficientPermissions
    }
    
    return nil
}
```

### Application Layer - Orchestrates Decorators

Use cases coordinate decorated services:

```go
// Command handler uses decorated services
type CreateUserHandler struct {
    userService    UserService    // May be decorated
    emailService   EmailService   // May be decorated
    eventPublisher EventPublisher // May be decorated
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
    // Create user through potentially decorated service
    user, err := domain.NewUser(cmd.Email, cmd.Username)
    if err != nil {
        return err
    }
    
    // All cross-cutting concerns (logging, metrics, etc.) are handled by decorators
    if err := h.userService.CreateUser(ctx, user); err != nil {
        return err
    }
    
    // Send welcome email (also potentially decorated)
    if err := h.emailService.SendWelcomeEmail(ctx, user.Email); err != nil {
        // Log but don't fail the command
        log.Printf("Failed to send welcome email: %v", err)
    }
    
    // Publish event (with potential retry, dead letter queue decorators)
    event := domain.UserCreatedEvent{UserID: user.ID, Email: user.Email}
    return h.eventPublisher.Publish(ctx, event)
}
```

### Infrastructure Layer - Implements Decorated Adapters

Repository decorators for cross-cutting concerns:

```go
// Base repository implementation
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(ctx context.Context, user domain.User) error {
    query := `INSERT INTO users (id, email, username, role, active) 
              VALUES ($1, $2, $3, $4, $5)`
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.Role, user.Active)
    return err
}

// Metrics decorator for repository
type MetricsUserRepository struct {
    next    domain.UserRepository
    metrics MetricsCollector
}

func (r *MetricsUserRepository) Save(ctx context.Context, user domain.User) error {
    start := time.Now()
    err := r.next.Save(ctx, user)
    
    duration := time.Since(start)
    r.metrics.RecordRepositoryOperation("user_save", duration, err)
    
    return err
}

// Circuit breaker decorator for external services
type CircuitBreakerEmailService struct {
    next    EmailService
    breaker *CircuitBreaker
}

func (s *CircuitBreakerEmailService) SendWelcomeEmail(ctx context.Context, email string) error {
    return s.breaker.Execute(func() error {
        return s.next.SendWelcomeEmail(ctx, email)
    })
}
```

---

## Domain-Driven Design Applications

### Aggregate Decorators

Enhance aggregates with cross-cutting concerns:

```go
// Training aggregate with decorators
type Training struct {
    ID          TrainingID
    Title       string
    TrainerID   UserID
    Attendees   []UserID
    ScheduledAt time.Time
    Status      TrainingStatus
}

// Audit decorator for domain operations
type AuditedTraining struct {
    training    *Training
    auditLogger AuditLogger
    currentUser User
}

func NewAuditedTraining(training *Training, logger AuditLogger, user User) *AuditedTraining {
    return &AuditedTraining{
        training:    training,
        auditLogger: logger,
        currentUser: user,
    }
}

func (a *AuditedTraining) AddAttendee(userID UserID) error {
    if err := a.training.AddAttendee(userID); err != nil {
        return err
    }
    
    // Log the domain event
    a.auditLogger.LogDomainEvent(AuditEvent{
        AggregateID:   a.training.ID.String(),
        AggregateType: "Training",
        EventType:     "AttendeeAdded",
        UserID:        a.currentUser.ID,
        Timestamp:     time.Now(),
        Data:          map[string]interface{}{"attendeeID": userID.String()},
    })
    
    return nil
}

// Validation decorator for domain operations
type ValidatedTraining struct {
    training  *Training
    validator TrainingValidator
}

func (v *ValidatedTraining) ScheduleAt(scheduledAt time.Time) error {
    // Pre-validation
    if err := v.validator.ValidateScheduleTime(scheduledAt); err != nil {
        return err
    }
    
    // Delegate to core domain logic
    if err := v.training.ScheduleAt(scheduledAt); err != nil {
        return err
    }
    
    // Post-validation (if needed)
    return v.validator.ValidateScheduledTraining(v.training)
}
```

### Repository Security Decorators

Following our security-first repository design:

```go
// Authorization decorator for repositories
type AuthorizedUserRepository struct {
    next        domain.UserRepository
    authService AuthorizationService
}

func (r *AuthorizedUserRepository) GetByID(ctx context.Context, userID string, requestingUser domain.User) (*domain.User, error) {
    // Check authorization before delegating
    if err := r.authService.CanAccessUser(requestingUser, userID); err != nil {
        return nil, err
    }
    
    return r.next.GetByID(ctx, userID, requestingUser)
}

func (r *AuthorizedUserRepository) Update(ctx context.Context, user domain.User, requestingUser domain.User) error {
    // Check update permissions
    if err := r.authService.CanModifyUser(requestingUser, user.ID); err != nil {
        return err
    }
    
    return r.next.Update(ctx, user, requestingUser)
}

// Encryption decorator for sensitive operations
type EncryptedUserRepository struct {
    next      domain.UserRepository
    encryptor FieldEncryptor
}

func (r *EncryptedUserRepository) Save(ctx context.Context, user domain.User) error {
    // Encrypt sensitive fields before storage
    encryptedUser := user
    encryptedUser.Email = r.encryptor.Encrypt(user.Email)
    encryptedUser.PhoneNumber = r.encryptor.Encrypt(user.PhoneNumber)
    
    return r.next.Save(ctx, encryptedUser)
}

func (r *EncryptedUserRepository) GetByID(ctx context.Context, userID string, requestingUser domain.User) (*domain.User, error) {
    user, err := r.next.GetByID(ctx, userID, requestingUser)
    if err != nil {
        return nil, err
    }
    
    // Decrypt sensitive fields after retrieval
    user.Email = r.encryptor.Decrypt(user.Email)
    user.PhoneNumber = r.encryptor.Decrypt(user.PhoneNumber)
    
    return user, nil
}
```

---

## CQRS Enhancement Patterns

### Command Handler Decorators

Add cross-cutting concerns to command processing:

```go
// Base command handler
type CreateTrainingHandler struct {
    trainingRepo domain.TrainingRepository
    eventBus     EventBus
}

func (h *CreateTrainingHandler) Handle(ctx context.Context, cmd CreateTrainingCommand) error {
    training := domain.NewTraining(cmd.Title, cmd.TrainerID, cmd.ScheduledAt)
    
    if err := h.trainingRepo.Save(ctx, training); err != nil {
        return err
    }
    
    event := domain.TrainingCreatedEvent{TrainingID: training.ID}
    return h.eventBus.Publish(ctx, event)
}

// Transaction decorator
type TransactionalHandler struct {
    next      CommandHandler
    txManager TransactionManager
}

func (h *TransactionalHandler) Handle(ctx context.Context, cmd Command) error {
    return h.txManager.WithTransaction(ctx, func(ctx context.Context) error {
        return h.next.Handle(ctx, cmd)
    })
}

// Idempotency decorator
type IdempotentHandler struct {
    next           CommandHandler
    idempotencyKey string
    store          IdempotencyStore
}

func (h *IdempotentHandler) Handle(ctx context.Context, cmd Command) error {
    // Check if command was already processed
    if result, processed := h.store.GetResult(h.idempotencyKey); processed {
        return result.Error
    }
    
    // Process command
    err := h.next.Handle(ctx, cmd)
    
    // Store result
    h.store.StoreResult(h.idempotencyKey, IdempotencyResult{Error: err})
    
    return err
}

// Composition
func NewDecoratedCreateTrainingHandler(
    base *CreateTrainingHandler,
    txManager TransactionManager,
    idempotencyStore IdempotencyStore,
) CommandHandler {
    return &TransactionalHandler{
        next:      &IdempotentHandler{
            next:           base,
            store:          idempotencyStore,
        },
        txManager: txManager,
    }
}
```

### Query Handler Decorators

Enhance query processing with caching and monitoring:

```go
// Cached query handler
type CachedQueryHandler struct {
    next  QueryHandler
    cache Cache
    ttl   time.Duration
}

func (h *CachedQueryHandler) Handle(ctx context.Context, query Query) (interface{}, error) {
    // Generate cache key
    key := h.generateCacheKey(query)
    
    // Check cache
    if result, found := h.cache.Get(key); found {
        return result, nil
    }
    
    // Cache miss - execute query
    result, err := h.next.Handle(ctx, query)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    h.cache.Set(key, result, h.ttl)
    
    return result, nil
}

// Metrics decorator for query performance
type MetricsQueryHandler struct {
    next    QueryHandler
    metrics MetricsCollector
}

func (h *MetricsQueryHandler) Handle(ctx context.Context, query Query) (interface{}, error) {
    start := time.Now()
    queryType := fmt.Sprintf("%T", query)
    
    result, err := h.next.Handle(ctx, query)
    
    duration := time.Since(start)
    h.metrics.RecordQueryExecution(queryType, duration, err)
    
    return result, err
}
```

---

## Advanced Decorator Techniques

### 1. Conditional Decorators

Apply decorators based on runtime conditions:

```go
type ConditionalDecorator struct {
    condition func(context.Context, interface{}) bool
    decorator Decorator
    next      Handler
}

func (d *ConditionalDecorator) Handle(ctx context.Context, request interface{}) error {
    if d.condition(ctx, request) {
        decoratedHandler := d.decorator(d.next)
        return decoratedHandler.Handle(ctx, request)
    }
    
    return d.next.Handle(ctx, request)
}

// Usage: Only apply caching for specific query types
cachingCondition := func(ctx context.Context, query interface{}) bool {
    switch query.(type) {
    case GetUserQuery, GetTrainingQuery:
        return true
    default:
        return false
    }
}

conditionalHandler := &ConditionalDecorator{
    condition: cachingCondition,
    decorator: WithCaching(cache, 5*time.Minute),
    next:      baseHandler,
}
```

### 2. Async Decorators

Add asynchronous processing capabilities:

```go
type AsyncDecorator struct {
    next   Handler
    pool   *WorkerPool
    buffer chan AsyncTask
}

type AsyncTask struct {
    Context context.Context
    Request interface{}
    Result  chan AsyncResult
}

type AsyncResult struct {
    Response interface{}
    Error    error
}

func (d *AsyncDecorator) Handle(ctx context.Context, request interface{}) error {
    task := AsyncTask{
        Context: ctx,
        Request: request,
        Result:  make(chan AsyncResult, 1),
    }
    
    // Submit to worker pool
    d.buffer <- task
    
    // Wait for result or timeout
    select {
    case result := <-task.Result:
        return result.Error
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (d *AsyncDecorator) worker() {
    for task := range d.buffer {
        response, err := d.next.Handle(task.Context, task.Request)
        
        task.Result <- AsyncResult{
            Response: response,
            Error:    err,
        }
    }
}
```

### 3. Pipeline Decorators

Chain multiple processing stages:

```go
type PipelineDecorator struct {
    stages []ProcessingStage
}

type ProcessingStage interface {
    Process(ctx context.Context, data interface{}) (interface{}, error)
}

func (p *PipelineDecorator) Handle(ctx context.Context, input interface{}) (interface{}, error) {
    current := input
    
    for _, stage := range p.stages {
        var err error
        current, err = stage.Process(ctx, current)
        if err != nil {
            return nil, err
        }
    }
    
    return current, nil
}

// Example stages
type ValidationStage struct{}
func (s *ValidationStage) Process(ctx context.Context, data interface{}) (interface{}, error) {
    // Validation logic
    return data, nil
}

type TransformationStage struct{}
func (s *TransformationStage) Process(ctx context.Context, data interface{}) (interface{}, error) {
    // Transformation logic
    return transformData(data), nil
}

type EnrichmentStage struct{}
func (s *EnrichmentStage) Process(ctx context.Context, data interface{}) (interface{}, error) {
    // Enrichment logic
    return enrichData(data), nil
}
```

---

## Testing Decorated Components

### Unit Testing Individual Decorators

Test decorators in isolation:

```go
func TestCachedUserService_GetUser_CacheHit(t *testing.T) {
    // Arrange
    mockNext := &MockUserService{}
    cache := make(map[string]*domain.User)
    
    expectedUser := &domain.User{ID: "user123", Email: "test@example.com"}
    cache["user123"] = expectedUser
    
    service := &CachedUserService{
        next:  mockNext,
        cache: cache,
    }
    
    // Act
    user, err := service.GetUser(context.Background(), "user123")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    
    // Verify that next service was not called
    mockNext.AssertNotCalled(t, "GetUser")
}

func TestCachedUserService_GetUser_CacheMiss(t *testing.T) {
    // Arrange
    mockNext := &MockUserService{}
    expectedUser := &domain.User{ID: "user123", Email: "test@example.com"}
    
    mockNext.On("GetUser", mock.Anything, "user123").Return(expectedUser, nil)
    
    service := &CachedUserService{
        next:  mockNext,
        cache: make(map[string]*domain.User),
    }
    
    // Act
    user, err := service.GetUser(context.Background(), "user123")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    
    // Verify cache was populated
    cachedUser, exists := service.cache["user123"]
    assert.True(t, exists)
    assert.Equal(t, expectedUser, cachedUser)
    
    mockNext.AssertExpectations(t)
}
```

### Integration Testing Decorator Chains

Test the behavior of composed decorators:

```go
func TestDecoratedUserService_FullChain(t *testing.T) {
    // Arrange
    baseRepo := &InMemoryUserRepository{}
    metrics := &MockMetricsCollector{}
    logger := &MockLogger{}
    
    // Build decorated service
    service := NewServiceBuilder(baseRepo).
        WithLogging(logger).
        WithMetrics(metrics).
        WithCache(time.Minute).
        Build()
    
    user := domain.User{ID: "user123", Email: "test@example.com"}
    
    // Act & Assert
    t.Run("create user flows through all decorators", func(t *testing.T) {
        err := service.CreateUser(context.Background(), user)
        assert.NoError(t, err)
        
        // Verify metrics were recorded
        metrics.AssertCalled(t, "RecordOperation", "create_user", mock.Anything, nil)
        
        // Verify logging occurred
        logger.AssertCalled(t, "Log", mock.MatchedBy(func(msg string) bool {
            return strings.Contains(msg, "CreateUser")
        }))
    })
    
    t.Run("get user uses cache after creation", func(t *testing.T) {
        // First call should hit repository
        user1, err := service.GetUser(context.Background(), "user123")
        assert.NoError(t, err)
        assert.Equal(t, user.Email, user1.Email)
        
        // Reset mock to verify cache usage
        baseRepo.ResetCalls()
        
        // Second call should use cache
        user2, err := service.GetUser(context.Background(), "user123")
        assert.NoError(t, err)
        assert.Equal(t, user.Email, user2.Email)
        
        // Verify repository wasn't called again
        baseRepo.AssertNotCalled(t, "GetByID")
    })
}
```

### Testing Decorator Error Handling

Ensure decorators properly handle and propagate errors:

```go
func TestRetryDecorator_ErrorHandling(t *testing.T) {
    tests := []struct {
        name           string
        failures       int
        maxAttempts    int
        expectedCalls  int
        expectSuccess  bool
    }{
        {
            name:          "succeeds on first attempt",
            failures:      0,
            maxAttempts:   3,
            expectedCalls: 1,
            expectSuccess: true,
        },
        {
            name:          "succeeds after retry",
            failures:      2,
            maxAttempts:   3,
            expectedCalls: 3,
            expectSuccess: true,
        },
        {
            name:          "fails after max attempts",
            failures:      3,
            maxAttempts:   3,
            expectedCalls: 3,
            expectSuccess: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockNext := &MockHandler{}
            
            // Set up mock to fail specified number of times
            for i := 0; i < tt.failures; i++ {
                mockNext.On("Handle", mock.Anything, mock.Anything).Return(errors.New("temporary error")).Once()
            }
            if tt.expectSuccess {
                mockNext.On("Handle", mock.Anything, mock.Anything).Return(nil).Once()
            }
            
            decorator := &RetryDecorator{
                next:        mockNext,
                maxAttempts: tt.maxAttempts,
                backoff:     time.Millisecond,
            }
            
            err := decorator.Handle(context.Background(), "test")
            
            if tt.expectSuccess {
                assert.NoError(t, err)
            } else {
                assert.Error(t, err)
            }
            
            mockNext.AssertNumberOfCalls(t, "Handle", tt.expectedCalls)
        })
    }
}
```

---

## Performance Considerations

### Memory Management

Decorators can create deep object hierarchies. Consider memory implications:

```go
// Use object pooling for frequently created decorators
var decoratorPool = sync.Pool{
    New: func() interface{} {
        return &MetricsDecorator{}
    },
}

func NewPooledMetricsDecorator(next Handler, collector MetricsCollector) *MetricsDecorator {
    decorator := decoratorPool.Get().(*MetricsDecorator)
    decorator.next = next
    decorator.collector = collector
    return decorator
}

func (d *MetricsDecorator) Release() {
    d.next = nil
    d.collector = nil
    decoratorPool.Put(d)
}

// Use context for request-scoped decorators instead of creating new instances
type ContextualDecorator struct {
    next Handler
}

func (d *ContextualDecorator) Handle(ctx context.Context, request interface{}) error {
    // Get request-specific configuration from context
    if config := GetDecoratorConfig(ctx); config != nil {
        return d.handleWithConfig(ctx, request, config)
    }
    
    return d.next.Handle(ctx, request)
}
```

### Avoiding Deep Call Stacks

Prevent stack overflow with many decorators:

```go
// Flatten decorator chain to avoid deep recursion
type FlattenedDecorator struct {
    decorators []Decorator
    base       Handler
}

func (f *FlattenedDecorator) Handle(ctx context.Context, request interface{}) error {
    // Build execution chain without deep nesting
    handler := f.base
    
    // Apply decorators in reverse order
    for i := len(f.decorators) - 1; i >= 0; i-- {
        handler = f.decorators[i].Wrap(handler)
    }
    
    return handler.Handle(ctx, request)
}

// Decorator interface that returns wrapped handler instead of calling next
type Decorator interface {
    Wrap(Handler) Handler
}
```

### Benchmarking Decorator Performance

Measure the overhead of your decorator chains:

```go
func BenchmarkDecoratorChain(b *testing.B) {
    baseHandler := &BaseHandler{}
    
    tests := []struct {
        name    string
        handler Handler
    }{
        {"base", baseHandler},
        {"with_logging", NewLoggingDecorator(baseHandler, logger)},
        {"with_metrics", NewMetricsDecorator(baseHandler, metrics)},
        {"full_chain", buildFullDecoratorChain(baseHandler)},
    }
    
    for _, tt := range tests {
        b.Run(tt.name, func(b *testing.B) {
            ctx := context.Background()
            request := TestRequest{}
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                tt.handler.Handle(ctx, request)
            }
        })
    }
}

func BenchmarkDecoratorMemory(b *testing.B) {
    baseHandler := &BaseHandler{}
    
    b.Run("decorator_creation", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _ = buildFullDecoratorChain(baseHandler)
        }
    })
}
```

---

## Common Pitfalls and Solutions

### 1. Interface Pollution

**Problem**: Creating too many interfaces for decorator compatibility

```go
// BAD: Over-engineered interfaces
type UserServiceWithCaching interface {
    UserService
    ClearCache(userID string)
    GetCacheStats() CacheStats
}

type UserServiceWithMetrics interface {
    UserService
    GetMetrics() Metrics
    ResetCounters()
}
```

**Solution**: Keep interfaces focused on core behavior

```go
// GOOD: Focused interface
type UserService interface {
    CreateUser(ctx context.Context, user domain.User) error
    GetUser(ctx context.Context, userID string) (*domain.User, error)
    UpdateUser(ctx context.Context, user domain.User) error
}

// Separate management interfaces if needed
type CacheManager interface {
    ClearCache(key string)
    GetCacheStats() CacheStats
}

type MetricsProvider interface {
    GetMetrics() Metrics
}
```

### 2. Decorator Order Dependencies

**Problem**: Decorators that depend on specific ordering

```go
// BAD: Order-dependent decorators
func buildService() UserService {
    base := &BasicUserService{}
    
    // These decorators have hidden dependencies on order
    cached := NewCachedService(base)
    metriced := NewMetricsService(cached) // Metrics won't measure cache misses
    
    return metriced
}
```

**Solution**: Make dependencies explicit and order-independent

```go
// GOOD: Explicit dependencies and proper ordering
func buildService() UserService {
    base := &BasicUserService{}
    
    // Metrics should wrap the actual operations
    metriced := NewMetricsService(base)
    
    // Cache should be outermost to measure cache hits/misses
    cached := NewCachedService(metriced)
    
    return cached
}

// Or use a builder that enforces correct ordering
type ServiceBuilder struct {
    service UserService
    layers  []string
}

func (b *ServiceBuilder) WithMetrics(collector MetricsCollector) *ServiceBuilder {
    if contains(b.layers, "cache") {
        panic("metrics must be added before caching")
    }
    
    b.service = NewMetricsService(b.service, collector)
    b.layers = append(b.layers, "metrics")
    return b
}
```

### 3. Context Pollution

**Problem**: Passing too much configuration through context

```go
// BAD: Context becomes a grab bag
ctx = context.WithValue(ctx, "cache_ttl", 5*time.Minute)
ctx = context.WithValue(ctx, "retry_count", 3)
ctx = context.WithValue(ctx, "log_level", "debug")
ctx = context.WithValue(ctx, "metrics_tags", map[string]string{"service": "user"})
```

**Solution**: Use proper configuration injection

```go
// GOOD: Configuration through constructor
type ConfiguredDecorator struct {
    next   Handler
    config DecoratorConfig
}

type DecoratorConfig struct {
    CacheTTL     time.Duration
    RetryCount   int
    LogLevel     string
    MetricsTags  map[string]string
}

func NewConfiguredDecorator(next Handler, config DecoratorConfig) *ConfiguredDecorator {
    return &ConfiguredDecorator{
        next:   next,
        config: config,
    }
}
```

### 4. Error Wrapping Confusion

**Problem**: Multiple decorators wrapping the same error

```go
// BAD: Error gets wrapped multiple times
func (d *LoggingDecorator) Handle(ctx context.Context, req interface{}) error {
    err := d.next.Handle(ctx, req)
    if err != nil {
        return fmt.Errorf("logging decorator: %w", err)
    }
    return nil
}

func (d *MetricsDecorator) Handle(ctx context.Context, req interface{}) error {
    err := d.next.Handle(ctx, req)
    if err != nil {
        return fmt.Errorf("metrics decorator: %w", err)
    }
    return nil
}
```

**Solution**: Only wrap errors when adding meaningful context

```go
// GOOD: Preserve original error in most decorators
func (d *LoggingDecorator) Handle(ctx context.Context, req interface{}) error {
    err := d.next.Handle(ctx, req)
    if err != nil {
        d.logger.Error("operation failed", "error", err, "request", req)
        // Return original error, don't wrap
        return err
    }
    return nil
}

// Only wrap when adding meaningful context
func (d *RetryDecorator) Handle(ctx context.Context, req interface{}) error {
    var lastErr error
    
    for attempt := 1; attempt <= d.maxAttempts; attempt++ {
        lastErr = d.next.Handle(ctx, req)
        if lastErr == nil {
            return nil
        }
        
        // Wait before retry...
    }
    
    // Add meaningful context about retry attempts
    return fmt.Errorf("failed after %d attempts: %w", d.maxAttempts, lastErr)
}
```

### 5. Resource Leak in Decorators

**Problem**: Decorators that don't properly clean up resources

```go
// BAD: Background goroutines without cleanup
type AsyncDecorator struct {
    next   Handler
    buffer chan Task
}

func NewAsyncDecorator(next Handler) *AsyncDecorator {
    d := &AsyncDecorator{
        next:   next,
        buffer: make(chan Task, 100),
    }
    
    // Goroutine will run forever!
    go d.processBuffer()
    
    return d
}
```

**Solution**: Implement proper lifecycle management

```go
// GOOD: Proper resource management
type AsyncDecorator struct {
    next   Handler
    buffer chan Task
    done   chan struct{}
    wg     sync.WaitGroup
}

func NewAsyncDecorator(next Handler) *AsyncDecorator {
    d := &AsyncDecorator{
        next:   next,
        buffer: make(chan Task, 100),
        done:   make(chan struct{}),
    }
    
    d.wg.Add(1)
    go d.processBuffer()
    
    return d
}

func (d *AsyncDecorator) processBuffer() {
    defer d.wg.Done()
    
    for {
        select {
        case task := <-d.buffer:
            d.processTask(task)
        case <-d.done:
            return
        }
    }
}

func (d *AsyncDecorator) Close() error {
    close(d.done)
    d.wg.Wait()
    return nil
}
```

---

## Best Practices

### 1. Interface Design

- **Keep interfaces small and focused**: Follow the Interface Segregation Principle
- **Design for the client**: Interfaces should reflect what clients need, not what implementations provide
- **Use embedding for interface composition**: Combine smaller interfaces into larger ones when needed

```go
// Small, focused interfaces
type Reader interface {
    Read(ctx context.Context, id string) (interface{}, error)
}

type Writer interface {
    Write(ctx context.Context, data interface{}) error
}

// Compose when needed
type ReadWriter interface {
    Reader
    Writer
}
```

### 2. Decorator Composition

- **Use builders for complex compositions**: Make decorator assembly explicit and validated
- **Document decorator order requirements**: Some decorators must be applied in specific order
- **Provide pre-configured compositions**: Offer common decorator combinations as convenience functions

```go
// Convenience function for common patterns
func NewStandardUserService(base UserService, config ServiceConfig) UserService {
    return NewServiceBuilder(base).
        WithLogging(config.Logger).
        WithMetrics(config.MetricsCollector).
        WithRateLimit(config.RateLimit.RequestsPerMinute, time.Minute).
        WithCache(config.Cache.TTL).
        WithRetry(config.Retry.MaxAttempts, config.Retry.BackoffDuration).
        Build()
}
```

### 3. Error Handling

- **Preserve error context**: Don't lose important error information in decorator chains
- **Handle decorator-specific errors appropriately**: Some errors might require specific handling at decorator level
- **Use error wrapping judiciously**: Only wrap errors when adding meaningful context

### 4. Testing Strategy

- **Test decorators in isolation**: Unit test each decorator's specific behavior
- **Test decorator compositions**: Integration test common decorator chains
- **Use test doubles**: Mock the next handler to isolate decorator behavior
- **Test error propagation**: Ensure errors flow correctly through decorator chains

### 5. Performance Optimization

- **Profile decorator overhead**: Measure the performance impact of your decorator chains
- **Use object pooling for high-frequency decorators**: Reduce allocation overhead
- **Consider flattening deep chains**: Avoid excessive call stack depth
- **Cache decorator instances when possible**: Reuse decorators that don't hold request-specific state

### 6. Documentation

- **Document decorator purposes clearly**: Explain what each decorator does and why
- **Specify decorator order requirements**: Make ordering dependencies explicit
- **Provide usage examples**: Show how to compose decorators effectively
- **Document performance characteristics**: Explain the overhead and scaling behavior

---

## Conclusion

The Decorator pattern is a powerful tool for implementing cross-cutting concerns
in Go applications while maintaining clean architecture principles. When applied
correctly, decorators provide:

- **Flexible behavior composition**: Mix and match functionality as needed
- **Separation of concerns**: Keep business logic separate from infrastructure concerns
- **Testability**: Each decorator can be tested in isolation
- **Maintainability**: Add new behavior without modifying existing code

Remember to:

- Keep interfaces focused and minimal
- Be mindful of decorator order and dependencies
- Test thoroughly, both in isolation and composition
- Monitor performance impact of decorator chains
- Follow the principles outlined in your project guidelines

The examples in this document align with your Clean Architecture, DDD, and CQRS
principles, showing how decorators can enhance your application's architecture
without compromising its core design principles.

## References

- Decorator Pattern - Refactoring Guru
- Decorator Patterns in Go - DEV Community
