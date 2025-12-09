# Building Highly Performant APIs in Go: Production Guide

> A comprehensive guide to implementing high-performance, scalable API services in Go
> with Clean Architecture, Domain-Driven Design, and security-first principles.

<!-- REF: https://dev.to/andrey_s/building-high-load-api-services-in-go-from-design-to-production-2626 -->
<!-- REF: https://go.dev/blog/ -->
<!-- REF: https://threedots.tech/post/clean-architecture-in-go/ -->

## Table of Contents

1. [Performance Requirements](#performance-requirements)
2. [Architectural Foundation](#architectural-foundation)
3. [High-Performance Implementation Patterns](#high-performance-implementation-patterns)
4. [Concurrency and Scalability](#concurrency-and-scalability)
5. [Database and Caching Optimization](#database-and-caching-optimization)
6. [Security-First Performance](#security-first-performance)
7. [Monitoring and Observability](#monitoring-and-observability)
8. [Production Optimization](#production-optimization)
9. [Testing High-Performance APIs](#testing-high-performance-apis)
10. [Best Practices and Anti-Patterns](#best-practices-and-anti-patterns)

---

## Performance Requirements

### Defining High-Performance APIs

High-performance APIs are characterized by their ability to handle substantial traffic
volumes while maintaining responsiveness and reliability under demanding conditions.

**Key Performance Metrics:**

- **Requests Per Second (RPS)**: Target 10K-100K+ RPS depending on use case
- **Latency**:
  - **P50**: Under 10ms for critical paths
  - **P95**: Under 50ms for most operations
  - **P99**: Under 100ms maximum
- **Uptime**: 99.99% (52 minutes/year) to 99.999% (5 minutes/year)
- **Error Rate**: < 0.1% for production systems

### Performance vs. Consistency Trade-offs

Following the CAP theorem, high-performance APIs must balance:

- **Consistency**: All clients see the same data simultaneously
- **Availability**: System responds to requests even under stress
- **Partition Tolerance**: Continues operating during network failures

**Performance-Optimized Approaches:**

```go
// CP (Consistency + Partition Tolerance) - Financial systems
type PaymentService struct {
    repo domain.PaymentRepository
    lock sync.RWMutex // Ensure consistency
}

func (s *PaymentService) ProcessPayment(ctx context.Context, payment domain.Payment) error {
    s.lock.Lock()
    defer s.lock.Unlock()
    
    // Strong consistency for financial operations
    return s.repo.SaveWithTransaction(ctx, payment)
}

// AP (Availability + Partition Tolerance) - Social media feeds  
type FeedService struct {
    cache   Cache
    repo    domain.FeedRepository
    eventBus EventBus
}

func (s *FeedService) GetFeed(ctx context.Context, userID string) (*domain.Feed, error) {
    // Serve from cache for availability, update asynchronously
    if feed, err := s.cache.Get(ctx, userID); err == nil {
        // Async update in background
        go s.updateFeedAsync(ctx, userID)
        return feed, nil
    }
    
    // Fallback to repository
    return s.repo.GetFeed(ctx, userID)
}
```

---

## Architectural Foundation

### Clean Architecture for Performance

Clean Architecture provides the foundation for high-performance APIs by separating
concerns and enabling targeted optimizations at each layer:

```
┌─────────────────────────────────────────────────────┐
│                Interface Layer                      │
│  - HTTP/gRPC handlers with connection pooling       │  
│  - Request validation and response optimization     │
│  - Rate limiting and middleware composition         │
└─────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────┐
│              Application Layer                      │
│  - Use case orchestration with async patterns       │
│  - Command/Query separation (CQRS)                  │
│  - Event handling and message queuing               │
└─────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────┐
│                Domain Layer                         │
│  - Pure business logic (no I/O)                     │
│  - Optimized domain models and value objects        │
│  - Domain events for decoupling                     │
└─────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────┐
│            Infrastructure Layer                     │
│  - Database connection pooling                      │
│  - Caching strategies (Redis, in-memory)            │
│  - External service clients with circuit breakers   │
└─────────────────────────────────────────────────────┘
```

### Domain-Driven Design for High Performance

DDD principles aligned with performance optimization:

```go
// Domain Layer - Optimized for performance
package domain

// Value objects are immutable and efficient
type UserID struct {
    value string // Pre-validated and cached
}

func NewUserID(id string) (UserID, error) {
    // Validation happens once during creation
    if len(id) == 0 || len(id) > 36 {
        return UserID{}, ErrInvalidUserID
    }
    return UserID{value: id}, nil
}

// Aggregate optimized for specific use cases
type User struct {
    id        UserID
    email     Email
    profile   Profile
    // Lazy loading for heavy objects
    preferences *UserPreferences // nil until accessed
}

// Repository with performance-first design
type UserRepository interface {
    // Specific methods avoid over-fetching
    GetByID(ctx context.Context, id UserID) (*User, error)
    GetEmailByID(ctx context.Context, id UserID) (Email, error)
    GetActiveUsersInRegion(ctx context.Context, region Region, limit int) ([]*User, error)
    
    // Batch operations for efficiency
    GetByIDs(ctx context.Context, ids []UserID) ([]*User, error)
    UpdateEmailsBatch(ctx context.Context, updates []EmailUpdate) error
}
```

### CQRS for Performance Optimization

Command Query Responsibility Segregation enables independent optimization of read and write operations:

```go
// Command side - optimized for write performance
type CreateUserCommand struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,max=100"`
    Region   string `json:"region" validate:"required"`
}

type CreateUserHandler struct {
    writeRepo  domain.UserWriteRepository
    eventBus   EventBus
    validator  Validator
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
    // Fast validation
    if err := h.validator.Validate(cmd); err != nil {
        return NewValidationError(err)
    }
    
    // Create domain object
    user, err := domain.NewUser(cmd.Email, cmd.Name, cmd.Region)
    if err != nil {
        return err
    }
    
    // Optimized write operation
    if err := h.writeRepo.Save(ctx, user); err != nil {
        return err
    }
    
    // Async event publishing
    event := domain.UserCreatedEvent{
        UserID: user.ID(),
        Email:  user.Email(),
        Region: user.Region(),
        Time:   time.Now(),
    }
    
    // Non-blocking event publishing
    go h.eventBus.Publish(ctx, event)
    
    return nil
}

// Query side - optimized for read performance
type GetUserQuery struct {
    UserID string `json:"user_id"`
}

type GetUserHandler struct {
    readRepo domain.UserReadRepository
    cache    Cache
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*UserView, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user:%s", query.UserID)
    if cached, err := h.cache.Get(ctx, cacheKey); err == nil {
        return cached.(*UserView), nil
    }
    
    // Read from optimized read store
    user, err := h.readRepo.GetUserView(ctx, query.UserID)
    if err != nil {
        return nil, err
    }
    
    // Cache for next request
    go h.cache.Set(ctx, cacheKey, user, 5*time.Minute)
    
    return user, nil
}

// Optimized read model
type UserView struct {
    ID          string    `json:"id"`
    Email       string    `json:"email"`
    Name        string    `json:"name"`
    Region      string    `json:"region"`
    LastActive  time.Time `json:"last_active"`
    // Pre-computed fields for performance
    DisplayName string    `json:"display_name"`
    AvatarURL   string    `json:"avatar_url"`
}
```

---

## High-Performance Implementation Patterns

### Optimized HTTP Server Configuration

```go
package main

import (
    "context"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// High-performance server configuration
func NewOptimizedServer(handlers *Handlers) *http.Server {
    // Production-optimized Gin setup
    gin.SetMode(gin.ReleaseMode)
    router := gin.New()
    
    // Minimal middleware for maximum performance
    router.Use(
        RequestIDMiddleware(),
        LoggingMiddleware(),
        RecoveryMiddleware(),
        CORSMiddleware(),
        RateLimitMiddleware(1000), // 1000 RPS per client
    )
    
    // Optimized server configuration
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
        
        // Connection optimization
        ReadTimeout:       5 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       60 * time.Second,
        ReadHeaderTimeout: 2 * time.Second,
        
        // Performance tuning
        MaxHeaderBytes: 1 << 20, // 1MB
    }
    
    return server
}

// Connection pooling for database
func NewOptimizedDB() *sql.DB {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    
    // Connection pool optimization
    db.SetMaxOpenConns(100)    // Maximum connections
    db.SetMaxIdleConns(25)     // Idle connections to retain
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(30 * time.Second)
    
    return db
}
```

### High-Performance Middleware

```go
// Request ID middleware with minimal allocation
func RequestIDMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID() // Fast UUID generation
        }
        c.Header("X-Request-ID", requestID)
        c.Set("request_id", requestID)
        c.Next()
    })
}

// High-performance rate limiting middleware
func RateLimitMiddleware(rps int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(rps), rps*2) // Allow bursts
    
    return gin.HandlerFunc(func(c *gin.Context) {
        if !limiter.Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
                "retry_after": "1s",
            })
            return
        }
        c.Next()
    })
}

// Zero-allocation logging middleware
func LoggingMiddleware() gin.HandlerFunc {
    logger := zap.L()
    
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()
        
        // Log only after request completes
        duration := time.Since(start)
        status := c.Writer.Status()
        
        logger.Info("request",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.Int("status", status),
            zap.Duration("duration", duration),
            zap.String("request_id", c.GetString("request_id")),
        )
    })
}
```

### Optimized JSON Handling

```go
// Pre-allocated JSON encoder pool
var jsonEncoderPool = sync.Pool{
    New: func() interface{} {
        return json.NewEncoder(nil)
    },
}

// High-performance JSON response helper
func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    // Use pooled encoder to avoid allocations
    encoder := jsonEncoderPool.Get().(*json.Encoder)
    defer jsonEncoderPool.Put(encoder)
    
    encoder.Reset(w)
    return encoder.Encode(data)
}

// Optimized request parsing with validation
type FastValidator struct {
    validate *validator.Validate
}

func (fv *FastValidator) ValidateStruct(s interface{}) error {
    return fv.validate.Struct(s)
}

// Fast request binding with minimal allocations
func BindAndValidate(c *gin.Context, obj interface{}) error {
    if err := c.ShouldBindJSON(obj); err != nil {
        return NewValidationError("invalid JSON", err)
    }
    
    if err := validator.Validate(obj); err != nil {
        return NewValidationError("validation failed", err)
    }
    
    return nil
}
```

---

## Concurrency and Scalability

### Goroutine Patterns for High Performance

```go
// Worker pool pattern for high-throughput processing
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    resultChan chan Result
    quit       chan bool
    wg         sync.WaitGroup
}

type Job struct {
    ID      string
    Payload interface{}
    Context context.Context
}

type Result struct {
    JobID  string
    Data   interface{}
    Error  error
}

func NewWorkerPool(workers, queueSize int) *WorkerPool {
    return &WorkerPool{
        workers:    workers,
        jobQueue:   make(chan Job, queueSize),
        resultChan: make(chan Result, queueSize),
        quit:       make(chan bool),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for {
        select {
        case job := <-wp.jobQueue:
            result := wp.processJob(job)
            wp.resultChan <- result
            
        case <-wp.quit:
            return
        }
    }
}

func (wp *WorkerPool) processJob(job Job) Result {
    // Process job with context cancellation support
    ctx, cancel := context.WithTimeout(job.Context, 30*time.Second)
    defer cancel()
    
    // Your business logic here
    data, err := processBusinessLogic(ctx, job.Payload)
    
    return Result{
        JobID: job.ID,
        Data:  data,
        Error: err,
    }
}

// Submit job with backpressure handling
func (wp *WorkerPool) Submit(job Job) bool {
    select {
    case wp.jobQueue <- job:
        return true
    default:
        return false // Queue is full, handle backpressure
    }
}

// Graceful shutdown
func (wp *WorkerPool) Stop() {
    close(wp.quit)
    wp.wg.Wait()
    close(wp.jobQueue)
    close(wp.resultChan)
}
```

### Fan-Out/Fan-In Pattern for Parallel Processing

```go
// High-performance fan-out/fan-in for parallel API calls
func ProcessUsersBatch(ctx context.Context, userIDs []string, processor UserProcessor) ([]*UserResult, error) {
    const maxWorkers = 10
    const batchSize = 100
    
    // Channel for work distribution
    jobs := make(chan []string, maxWorkers)
    results := make(chan []*UserResult, maxWorkers)
    
    // Start worker goroutines
    var wg sync.WaitGroup
    for i := 0; i < maxWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for batch := range jobs {
                batchResults := make([]*UserResult, 0, len(batch))
                for _, userID := range batch {
                    result, err := processor.Process(ctx, userID)
                    if err != nil {
                        result = &UserResult{UserID: userID, Error: err}
                    }
                    batchResults = append(batchResults, result)
                }
                results <- batchResults
            }
        }()
    }
    
    // Distribute work in batches
    go func() {
        defer close(jobs)
        for i := 0; i < len(userIDs); i += batchSize {
            end := i + batchSize
            if end > len(userIDs) {
                end = len(userIDs)
            }
            jobs <- userIDs[i:end]
        }
    }()
    
    // Collect results
    go func() {
        wg.Wait()
        close(results)
    }()
    
    var allResults []*UserResult
    for batchResults := range results {
        allResults = append(allResults, batchResults...)
    }
    
    return allResults, nil
}
```

### Context-Based Cancellation and Timeouts

```go
// Context patterns for high-performance APIs
func (h *UserHandler) GetUserProfile(c *gin.Context) {
    // Create request-scoped context with timeout
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    userID := c.Param("user_id")
    
    // Parallel fetching with context
    var (
        user    *domain.User
        profile *domain.Profile
        userErr error
        profErr error
    )
    
    var wg sync.WaitGroup
    wg.Add(2)
    
    // Fetch user data
    go func() {
        defer wg.Done()
        user, userErr = h.userRepo.GetByID(ctx, userID)
    }()
    
    // Fetch profile data
    go func() {
        defer wg.Done()
        profile, profErr = h.profileRepo.GetByUserID(ctx, userID)
    }()
    
    wg.Wait()
    
    // Check for context cancellation
    if ctx.Err() != nil {
        c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
        return
    }
    
    // Handle errors
    if userErr != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    
    // Construct response
    response := UserProfileResponse{
        User:    user,
        Profile: profile,
    }
    
    c.JSON(http.StatusOK, response)
}
```

---

## Database and Caching Optimization

### High-Performance Database Patterns

```go
// Optimized repository implementation
type PostgresUserRepository struct {
    db    *sql.DB
    cache Cache
    
    // Pre-prepared statements for performance
    getUserStmt      *sql.Stmt
    getUsersStmt     *sql.Stmt
    updateEmailStmt  *sql.Stmt
}

func NewPostgresUserRepository(db *sql.DB, cache Cache) (*PostgresUserRepository, error) {
    repo := &PostgresUserRepository{
        db:    db,
        cache: cache,
    }
    
    // Prepare statements once at startup
    var err error
    repo.getUserStmt, err = db.Prepare(`
        SELECT id, email, name, created_at, updated_at 
        FROM users 
        WHERE id = $1`)
    if err != nil {
        return nil, err
    }
    
    repo.getUsersStmt, err = db.Prepare(`
        SELECT id, email, name, created_at, updated_at 
        FROM users 
        WHERE id = ANY($1)`)
    if err != nil {
        return nil, err
    }
    
    return repo, nil
}

// Cache-aside pattern with database fallback
func (r *PostgresUserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user:%s", userID)
    if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
        return cached.(*domain.User), nil
    }
    
    // Database fallback
    var user domain.User
    err := r.getUserStmt.QueryRowContext(ctx, userID).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrUserNotFound
        }
        return nil, err
    }
    
    // Cache for next request (async to avoid blocking)
    go r.cache.Set(ctx, cacheKey, &user, 10*time.Minute)
    
    return &user, nil
}

// Batch operations for efficiency
func (r *PostgresUserRepository) GetByIDs(ctx context.Context, userIDs []string) ([]*domain.User, error) {
    // Check cache for all IDs first
    users := make([]*domain.User, 0, len(userIDs))
    missingIDs := make([]string, 0)
    
    for _, id := range userIDs {
        cacheKey := fmt.Sprintf("user:%s", id)
        if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
            users = append(users, cached.(*domain.User))
        } else {
            missingIDs = append(missingIDs, id)
        }
    }
    
    // Fetch missing users in batch
    if len(missingIDs) > 0 {
        rows, err := r.getUsersStmt.QueryContext(ctx, pq.Array(missingIDs))
        if err != nil {
            return nil, err
        }
        defer rows.Close()
        
        var fetchedUsers []*domain.User
        for rows.Next() {
            var user domain.User
            err := rows.Scan(
                &user.ID,
                &user.Email,
                &user.Name,
                &user.CreatedAt,
                &user.UpdatedAt,
            )
            if err != nil {
                return nil, err
            }
            fetchedUsers = append(fetchedUsers, &user)
            users = append(users, &user)
        }
        
        // Cache fetched users asynchronously
        go func() {
            for _, user := range fetchedUsers {
                cacheKey := fmt.Sprintf("user:%s", user.ID)
                r.cache.Set(context.Background(), cacheKey, user, 10*time.Minute)
            }
        }()
    }
    
    return users, nil
}
```

### Multi-Layer Caching Strategy

```go
// Multi-layer cache implementation
type MultiLevelCache struct {
    l1 Cache // In-memory (fastest)
    l2 Cache // Redis (shared across instances)
    l3 Cache // Database query cache
}

func (c *MultiLevelCache) Get(ctx context.Context, key string) (interface{}, error) {
    // Try L1 cache first (in-memory)
    if val, err := c.l1.Get(ctx, key); err == nil {
        return val, nil
    }
    
    // Try L2 cache (Redis)
    if val, err := c.l2.Get(ctx, key); err == nil {
        // Populate L1 cache asynchronously
        go c.l1.Set(context.Background(), key, val, 1*time.Minute)
        return val, nil
    }
    
    // Try L3 cache (database)
    return c.l3.Get(ctx, key)
}

func (c *MultiLevelCache) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
    // Set in all layers
    var wg sync.WaitGroup
    wg.Add(3)
    
    go func() {
        defer wg.Done()
        c.l1.Set(ctx, key, val, ttl)
    }()
    
    go func() {
        defer wg.Done()
        c.l2.Set(ctx, key, val, ttl)
    }()
    
    go func() {
        defer wg.Done()
        c.l3.Set(ctx, key, val, ttl)
    }()
    
    wg.Wait()
    return nil
}

// Cache warming strategies
func (c *MultiLevelCache) WarmCache(ctx context.Context, keys []string) error {
    const batchSize = 100
    for i := 0; i < len(keys); i += batchSize {
        end := i + batchSize
        if end > len(keys) {
            end = len(keys)
        }
        
        batch := keys[i:end]
        go c.warmBatch(ctx, batch)
    }
    return nil
}
```

---

## Security-First Performance

### Secure Repository Design with Performance

Following security-first repository design principles while maintaining high performance:

```go
// Security-first repository with performance optimizations
type SecureUserRepository struct {
    db           *sql.DB
    cache        Cache
    authz        AuthorizationService
    
    // Prepared statements for security and performance
    getUserByIDStmt          *sql.Stmt
    getUsersByTenantStmt     *sql.Stmt
    updateUserStmt           *sql.Stmt
}

// Repository method with built-in authorization and caching
func (r *SecureUserRepository) GetUser(ctx context.Context, userID string, requestingUser domain.User) (*domain.User, error) {
    // Security check built into repository interface
    if !r.authz.CanAccessUser(requestingUser, userID) {
        return nil, domain.ErrUnauthorized
    }
    
    // Performance optimization: use tenant-scoped cache key
    tenantID := requestingUser.TenantID()
    cacheKey := fmt.Sprintf("user:%s:%s", tenantID, userID)
    
    // Check cache with security context
    if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
        user := cached.(*domain.User)
        // Double-check authorization on cached data
        if user.TenantID() == tenantID {
            return user, nil
        }
        // Invalid cache entry, remove it
        go r.cache.Delete(context.Background(), cacheKey)
    }
    
    // Database query with tenant isolation
    var user domain.User
    err := r.getUserByIDStmt.QueryRowContext(ctx, userID, tenantID).Scan(
        &user.ID,
        &user.TenantID,
        &user.Email,
        &user.Name,
        &user.CreatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrUserNotFound
        }
        return nil, err
    }
    
    // Cache with tenant-scoped key
    go r.cache.Set(ctx, cacheKey, &user, 5*time.Minute)
    
    return &user, nil
}

// Batch operations with security and performance
func (r *SecureUserRepository) GetUsersInTenant(ctx context.Context, tenantID string, requestingUser domain.User, limit int) ([]*domain.User, error) {
    // Authorization check
    if !r.authz.CanListUsersInTenant(requestingUser, tenantID) {
        return nil, domain.ErrUnauthorized
    }
    
    // Performance: limit check
    if limit > 1000 {
        limit = 1000 // Prevent excessive queries
    }
    
    cacheKey := fmt.Sprintf("tenant_users:%s:%d", tenantID, limit)
    
    // Check cache
    if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
        return cached.([]*domain.User), nil
    }
    
    // Query with tenant isolation built-in
    rows, err := r.getUsersByTenantStmt.QueryContext(ctx, tenantID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        var user domain.User
        err := rows.Scan(
            &user.ID,
            &user.TenantID,
            &user.Email,
            &user.Name,
            &user.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    
    // Cache with shorter TTL for larger datasets
    go r.cache.Set(ctx, cacheKey, users, 2*time.Minute)
    
    return users, nil
}
```

### Performance-Optimized Authentication & Authorization

```go
// Fast JWT validation with caching
type JWTAuthService struct {
    publicKey   *rsa.PublicKey
    cache       Cache
    tokenParser *jwt.Parser
}

func (a *JWTAuthService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
    // Check token cache first (for performance)
    cacheKey := fmt.Sprintf("token:%s", hashString(tokenString))
    if cached, err := a.cache.Get(ctx, cacheKey); err == nil {
        return cached.(*Claims), nil
    }
    
    // Parse and validate token
    token, err := a.tokenParser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("invalid signing method")
        }
        return a.publicKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }
    
    // Cache valid token (until expiry)
    ttl := time.Until(time.Unix(claims.ExpiresAt, 0))
    if ttl > 0 {
        go a.cache.Set(ctx, cacheKey, claims, ttl)
    }
    
    return claims, nil
}

// Authorization with caching and batch checks
type RBACAuthService struct {
    db    *sql.DB
    cache Cache
}

func (r *RBACAuthService) CanAccessResource(ctx context.Context, userID, resource, action string) (bool, error) {
    cacheKey := fmt.Sprintf("authz:%s:%s:%s", userID, resource, action)
    
    // Check cache
    if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
        return cached.(bool), nil
    }
    
    // Database check with optimized query
    var hasPermission bool
    query := `
        SELECT EXISTS(
            SELECT 1 FROM user_permissions up
            JOIN permissions p ON up.permission_id = p.id
            WHERE up.user_id = $1 
            AND p.resource = $2 
            AND p.action = $3
            AND up.expires_at > NOW()
        )`
    
    err := r.db.QueryRowContext(ctx, query, userID, resource, action).Scan(&hasPermission)
    if err != nil {
        return false, err
    }
    
    // Cache result
    go r.cache.Set(ctx, cacheKey, hasPermission, 5*time.Minute)
    
    return hasPermission, nil
}
```

---

## Monitoring and Observability

### High-Performance Metrics Collection

```go
// Optimized Prometheus metrics
var (
    // Use counters for cumulative values
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    // Use histograms for latency tracking
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "path"},
    )
    
    // Use gauges for current values
    activeConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )
)

// Performance-optimized middleware
func PrometheusMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        path := c.FullPath() // Use route pattern, not actual path
        method := c.Request.Method
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        // Update metrics (these are atomic operations)
        httpRequestsTotal.WithLabelValues(method, path, status).Inc()
        httpRequestDuration.WithLabelValues(method, path).Observe(duration)
    })
}
```

### Structured Logging for Performance

```go
// High-performance structured logging
type OptimizedLogger struct {
    *zap.Logger
    requestPool sync.Pool
}

func NewOptimizedLogger() *OptimizedLogger {
    logger, _ := zap.NewProduction(zap.Config{
        Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
        Encoding: "json",
        
        // Optimized encoder config
        EncoderConfig: zapcore.EncoderConfig{
            TimeKey:        "timestamp",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "", // Disable caller for performance
            MessageKey:     "message",
            StacktraceKey:  "", // Disable stacktrace for performance
            EncodeLevel:    zapcore.LowercaseLevelEncoder,
            EncodeTime:     zapcore.ISO8601TimeEncoder,
        },
        
        OutputPaths:      []string{"stdout"},
        ErrorOutputPaths: []string{"stderr"},
    })
    
    return &OptimizedLogger{
        Logger: logger,
        requestPool: sync.Pool{
            New: func() interface{} {
                return make([]zapcore.Field, 0, 10) // Pre-allocate capacity
            },
        },
    }
}

// Pool request context for reuse
func (l *OptimizedLogger) LogRequest(ctx context.Context, method, path string, status int, duration time.Duration) {
    fields := l.requestPool.Get().([]zapcore.Field)
    fields = fields[:0] // Reset length but keep capacity
    
    fields = append(fields,
        zap.String("method", method),
        zap.String("path", path),
        zap.Int("status", status),
        zap.Duration("duration", duration),
    )
    
    if requestID := ctx.Value("request_id"); requestID != nil {
        fields = append(fields, zap.String("request_id", requestID.(string)))
    }
    
    l.Info("request", fields...)
    
    // Return to pool
    l.requestPool.Put(fields)
}
```

### Distributed Tracing with OpenTelemetry

```go
// Optimized tracing setup
func NewOptimizedTracer(serviceName string) trace.Tracer {
    // Create resource
    resource := resource.NewWithAttributes(
        semconv.SchemaURL,
        semconv.ServiceName(serviceName),
        semconv.ServiceVersion("1.0.0"),
    )
    
    // Create OTLP exporter
    ctx := context.Background()
    exporter, err := otlptrace.New(ctx,
        otlptracegrpc.NewClient(
            otlptracegrpc.WithEndpoint("http://jaeger:14268"),
            otlptracegrpc.WithInsecure(),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create tracer provider with optimized settings
    tp := tracesdk.NewTracerProvider(
        tracesdk.WithBatcher(exporter,
            tracesdk.WithMaxExportBatchSize(512),
            tracesdk.WithBatchTimeout(5*time.Second),
            tracesdk.WithMaxExportBatchSize(512),
        ),
        tracesdk.WithResource(resource),
        
        // Sample traces for performance (adjust based on load)
        tracesdk.WithSampler(tracesdk.TraceIDRatioBased(0.1)), // 10% sampling
    )
    
    otel.SetTracerProvider(tp)
    
    return tp.Tracer(serviceName)
}

// Tracing middleware with minimal overhead
func TracingMiddleware(tracer trace.Tracer) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        ctx, span := tracer.Start(c.Request.Context(), 
            fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
            trace.WithSpanKind(trace.SpanKindServer),
        )
        defer span.End()
        
        // Add request attributes
        span.SetAttributes(
            semconv.HTTPMethod(c.Request.Method),
            semconv.HTTPRoute(c.FullPath()),
            semconv.HTTPScheme(c.Request.URL.Scheme),
        )
        
        c.Request = c.Request.WithContext(ctx)
        c.Next()
        
        // Add response attributes
        span.SetAttributes(
            semconv.HTTPStatusCode(c.Writer.Status()),
        )
        
        if c.Writer.Status() >= 400 {
            span.SetStatus(codes.Error, "HTTP error")
        }
    })
}
```

---

## Production Optimization

### Memory Optimization

```go
// Object pooling for frequently used objects
var responsePool = sync.Pool{
    New: func() interface{} {
        return &APIResponse{
            Data:   make(map[string]interface{}, 10),
            Errors: make([]string, 0, 5),
        }
    },
}

type APIResponse struct {
    Data      map[string]interface{} `json:"data,omitempty"`
    Errors    []string               `json:"errors,omitempty"`
    Meta      *ResponseMeta          `json:"meta,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}

func (r *APIResponse) Reset() {
    // Reset maps and slices without reallocating
    for k := range r.Data {
        delete(r.Data, k)
    }
    r.Errors = r.Errors[:0]
    r.Meta = nil
    r.Timestamp = time.Time{}
}

// Use pooled responses to reduce GC pressure
func WriteAPIResponse(c *gin.Context, status int, data interface{}) {
    resp := responsePool.Get().(*APIResponse)
    defer func() {
        resp.Reset()
        responsePool.Put(resp)
    }()
    
    resp.Data["result"] = data
    resp.Timestamp = time.Now()
    
    c.JSON(status, resp)
}

// String interning for repeated strings
type StringInterner struct {
    mu    sync.RWMutex
    cache map[string]string
}

func NewStringInterner() *StringInterner {
    return &StringInterner{
        cache: make(map[string]string),
    }
}

func (si *StringInterner) Intern(s string) string {
    si.mu.RLock()
    if interned, exists := si.cache[s]; exists {
        si.mu.RUnlock()
        return interned
    }
    si.mu.RUnlock()
    
    si.mu.Lock()
    defer si.mu.Unlock()
    
    // Double-check after acquiring write lock
    if interned, exists := si.cache[s]; exists {
        return interned
    }
    
    // Intern the string
    interned := string([]byte(s)) // Force allocation
    si.cache[s] = interned
    return interned
}
```

### Go Runtime Optimization

```go
// Runtime optimization for high-performance APIs
func OptimizeRuntime() {
    // Set GOMAXPROCS to match container limits
    if maxProcs := os.Getenv("GOMAXPROCS"); maxProcs != "" {
        if n, err := strconv.Atoi(maxProcs); err == nil {
            runtime.GOMAXPROCS(n)
        }
    }
    
    // Tune GC for low latency
    debug.SetGCPercent(50) // More frequent GC for lower latency
    
    // Set memory limit if in container
    if memLimit := os.Getenv("MEMORY_LIMIT"); memLimit != "" {
        if limit, err := strconv.ParseInt(memLimit, 10, 64); err == nil {
            debug.SetMemoryLimit(limit)
        }
    }
}

// Performance monitoring
func startPerformanceMonitoring(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            log.Info("performance_stats",
                zap.Uint64("alloc_mb", m.Alloc/1024/1024),
                zap.Uint64("total_alloc_mb", m.TotalAlloc/1024/1024),
                zap.Uint64("sys_mb", m.Sys/1024/1024),
                zap.Uint32("num_gc", m.NumGC),
                zap.Float64("gc_cpu_fraction", m.GCCPUFraction),
                zap.Int("num_goroutines", runtime.NumGoroutine()),
            )
            
        case <-ctx.Done():
            return
        }
    }
}
```

### Circuit Breaker Pattern

```go
// High-performance circuit breaker
type CircuitBreaker struct {
    mu             sync.RWMutex
    state          State
    failures       int
    successes      int
    lastFailure    time.Time
    failureThresh  int
    successThresh  int
    timeout        time.Duration
}

type State int

const (
    Closed State = iota
    Open
    HalfOpen
)

func NewCircuitBreaker(failureThresh, successThresh int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:         Closed,
        failureThresh: failureThresh,
        successThresh: successThresh,
        timeout:       timeout,
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.canExecute() {
        return ErrCircuitBreakerOpen
    }
    
    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    
    switch cb.state {
    case Closed:
        return true
    case Open:
        return time.Since(cb.lastFailure) > cb.timeout
    case HalfOpen:
        return true
    default:
        return false
    }
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        
        if cb.state == Closed && cb.failures >= cb.failureThresh {
            cb.state = Open
        } else if cb.state == HalfOpen {
            cb.state = Open
        }
    } else {
        cb.successes++
        
        if cb.state == HalfOpen && cb.successes >= cb.successThresh {
            cb.state = Closed
            cb.failures = 0
            cb.successes = 0
        } else if cb.state == Open {
            cb.state = HalfOpen
            cb.failures = 0
            cb.successes = 0
        }
    }
}
```

---

## Testing High-Performance APIs

### Performance Testing Patterns

```go
// Benchmark tests for high-performance components
func BenchmarkUserServiceCreateUser(b *testing.B) {
    service := setupUserService()
    ctx := context.Background()
    
    // Pre-generate test data to avoid allocation overhead in benchmark
    users := make([]CreateUserRequest, b.N)
    for i := 0; i < b.N; i++ {
        users[i] = CreateUserRequest{
            Email: fmt.Sprintf("user%d@example.com", i),
            Name:  fmt.Sprintf("User %d", i),
        }
    }
    
    b.ResetTimer()
    b.ReportAllocs() // Report memory allocations
    
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            _, err := service.CreateUser(ctx, users[i%len(users)])
            if err != nil {
                b.Fatal(err)
            }
            i++
        }
    })
}

// Load testing with realistic concurrency patterns
func BenchmarkAPIEndpoint_Concurrent(b *testing.B) {
    router := setupRouter()
    
    b.RunParallel(func(pb *testing.PB) {
        client := &http.Client{
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 100,
                IdleConnTimeout:     90 * time.Second,
            },
        }
        
        for pb.Next() {
            req, _ := http.NewRequest("GET", "/api/users/123", nil)
            resp, err := client.Do(req)
            if err != nil {
                b.Fatal(err)
            }
            resp.Body.Close()
        }
    })
}

// Memory allocation testing
func TestMemoryUsage(t *testing.T) {
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    // Execute operation under test
    service := setupUserService()
    for i := 0; i < 1000; i++ {
        service.CreateUser(context.Background(), CreateUserRequest{
            Email: fmt.Sprintf("user%d@example.com", i),
            Name:  fmt.Sprintf("User %d", i),
        })
    }
    
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    allocPerOp := (m2.TotalAlloc - m1.TotalAlloc) / 1000
    t.Logf("Memory allocation per operation: %d bytes", allocPerOp)
    
    // Assert reasonable memory usage
    if allocPerOp > 1024 { // 1KB per operation max
        t.Errorf("Memory usage too high: %d bytes per operation", allocPerOp)
    }
}
```

### Integration Testing for Performance

```go
// High-performance integration tests
func TestAPIPerformance(t *testing.T) {
    // Setup test server
    server := httptest.NewServer(setupRouter())
    defer server.Close()
    
    client := &http.Client{
        Timeout: 5 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
        },
    }
    
    // Performance test with realistic load
    const numRequests = 1000
    const concurrency = 50
    
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, concurrency)
    
    start := time.Now()
    
    for i := 0; i < numRequests; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release
            
            url := fmt.Sprintf("%s/api/users/%d", server.URL, i%100)
            resp, err := client.Get(url)
            if err != nil {
                t.Errorf("Request failed: %v", err)
                return
            }
            defer resp.Body.Close()
            
            if resp.StatusCode != http.StatusOK {
                t.Errorf("Expected 200, got %d", resp.StatusCode)
            }
        }(i)
    }
    
    wg.Wait()
    duration := time.Since(start)
    
    rps := float64(numRequests) / duration.Seconds()
    t.Logf("Performance: %.2f RPS over %v", rps, duration)
    
    // Assert performance requirements
    if rps < 500 { // Minimum 500 RPS requirement
        t.Errorf("Performance too low: %.2f RPS", rps)
    }
}
```

---

## Best Practices and Anti-Patterns

### Performance Best Practices

```go
// ✅ DO: Use connection pooling and reuse
func setupHTTPClient() *http.Client {
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        
        // Enable HTTP/2
        ForceAttemptHTTP2: true,
        
        // Optimize dial settings
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
    }
    
    return &http.Client{
        Transport: transport,
        Timeout:   10 * time.Second,
    }
}

// ✅ DO: Use prepared statements for database operations
type OptimizedRepository struct {
    db        *sql.DB
    getStmt   *sql.Stmt
    saveStmt  *sql.Stmt
}

func (r *OptimizedRepository) init() error {
    var err error
    r.getStmt, err = r.db.Prepare("SELECT * FROM users WHERE id = $1")
    if err != nil {
        return err
    }
    
    r.saveStmt, err = r.db.Prepare("INSERT INTO users (id, email, name) VALUES ($1, $2, $3)")
    return err
}

// ✅ DO: Use context for cancellation and timeouts
func (s *Service) ProcessWithTimeout(ctx context.Context, data Data) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        done <- s.processData(ctx, data)
    }()
    
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}

// ✅ DO: Batch operations for efficiency
func (r *Repository) SaveUsersBatch(ctx context.Context, users []*User) error {
    const batchSize = 100
    
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        if err := r.saveBatch(ctx, batch); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Common Anti-Patterns to Avoid

```go
// ❌ DON'T: Create new database connections for each request
func BadHandler(w http.ResponseWriter, r *http.Request) {
    // WRONG: Creates new connection every time
    db, _ := sql.Open("postgres", dsn)
    defer db.Close()
    
    // Use connection...
}

// ❌ DON'T: Ignore context cancellation
func BadProcessor(ctx context.Context, data []Item) error {
    for _, item := range data {
        // WRONG: Ignores context cancellation
        processItem(item) // This can run forever
    }
    return nil
}

// ✅ DO: Respect context cancellation
func GoodProcessor(ctx context.Context, data []Item) error {
    for _, item := range data {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := processItem(item); err != nil {
                return err
            }
        }
    }
    return nil
}

// ❌ DON'T: Use strings for frequently compared values
type BadStatus string
const (
    StatusActive   BadStatus = "active"
    StatusInactive BadStatus = "inactive"
)

// ✅ DO: Use integers for better performance
type Status int
const (
    StatusActive Status = iota
    StatusInactive
)

// ❌ DON'T: Allocate unnecessarily in hot paths
func BadJSONResponse(data interface{}) []byte {
    // WRONG: Creates new buffer every time
    buffer := &bytes.Buffer{}
    encoder := json.NewEncoder(buffer)
    encoder.Encode(data)
    return buffer.Bytes()
}

// ✅ DO: Use object pools for frequent allocations
var bufferPool = sync.Pool{
    New: func() interface{} {
        return &bytes.Buffer{}
    },
}

func GoodJSONResponse(data interface{}) []byte {
    buffer := bufferPool.Get().(*bytes.Buffer)
    buffer.Reset()
    defer bufferPool.Put(buffer)
    
    encoder := json.NewEncoder(buffer)
    encoder.Encode(data)
    return buffer.Bytes()
}
```

### Performance Monitoring and Alerting

```go
// Performance SLI/SLO monitoring
type PerformanceMonitor struct {
    latencyHist   prometheus.HistogramVec
    errorCounter  prometheus.CounterVec
    sloViolations prometheus.CounterVec
}

func NewPerformanceMonitor() *PerformanceMonitor {
    return &PerformanceMonitor{
        latencyHist: *prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "api_request_duration_seconds",
                Help: "API request duration in seconds",
                Buckets: []float64{.001, .005, .01, .05, .1, .25, .5, 1, 2.5, 5},
            },
            []string{"method", "endpoint"},
        ),
        errorCounter: *prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "api_errors_total",
                Help: "Total API errors",
            },
            []string{"method", "endpoint", "error_type"},
        ),
        sloViolations: *prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "slo_violations_total",
                Help: "SLO violations",
            },
            []string{"slo_type"},
        ),
    }
}

func (pm *PerformanceMonitor) RecordRequest(method, endpoint string, duration time.Duration, err error) {
    // Record latency
    pm.latencyHist.WithLabelValues(method, endpoint).Observe(duration.Seconds())
    
    // Record errors
    if err != nil {
        pm.errorCounter.WithLabelValues(method, endpoint, "error").Inc()
    }
    
    // Check SLO violations
    if duration > 50*time.Millisecond {
        pm.sloViolations.WithLabelValues("latency_p95").Inc()
    }
}

// Grafana alerting query examples:
// Latency P95 > 50ms: histogram_quantile(0.95, rate(api_request_duration_seconds_bucket[5m])) > 0.05
// Error rate > 1%: rate(api_errors_total[5m]) / rate(api_requests_total[5m]) > 0.01
// Availability < 99.9%: (1 - (rate(api_errors_total[5m]) / rate(api_requests_total[5m]))) < 0.999
```

---

## Conclusion

Building highly performant APIs in Go requires a systematic approach that combines architectural best practices with performance-specific optimizations. Key takeaways:

1. **Architecture First**: Clean Architecture and DDD provide the foundation for scalable, maintainable high-performance systems
2. **Security Integration**: Security-first design doesn't compromise performance when implemented correctly
3. **Measurement-Driven**: Use profiling, monitoring, and benchmarks to guide optimization decisions
4. **Go's Strengths**: Leverage goroutines, channels, and the standard library for maximum performance
5. **Production Readiness**: Consider memory management, runtime tuning, and operational concerns from day one

### Performance Targets Summary

| Metric | Target | Monitoring |
|--------|--------|------------|
| **Latency P95** | < 50ms | Histogram percentiles |
| **Latency P99** | < 100ms | Histogram percentiles |
| **Throughput** | 10K+ RPS | Request counter rates |
| **Error Rate** | < 0.1% | Error counter / total requests |
| **Availability** | > 99.9% | Uptime monitoring |
| **Memory Usage** | < 1KB/req | Allocation profiling |

### Next Steps

1. **Profile First**: Use `go tool pprof` to identify actual bottlenecks
2. **Load Test**: Implement realistic load testing with tools like `wrk` or `k6`
3. **Monitor Continuously**: Set up comprehensive monitoring and alerting
4. **Iterate**: Performance optimization is an ongoing process

---

**References:**
- [Building High-Load API Services in Go](https://dev.to/andrey_s/building-high-load-api-services-in-go-from-design-to-production-2626)
- [Go Blog - Official Go Team](https://go.dev/blog/)
- [Clean Architecture in Go - Three Dots Labs](https://threedots.tech/post/clean-architecture-in-go/)
- [Go Performance Best Practices](https://github.com/golang/go/wiki/Performance)
