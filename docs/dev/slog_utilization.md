# 📋 Structured Logging with slog in MCP YardGopher

This document explains how we utilize Go's `slog` package for structured logging
throughout the MCP YardGopher project. The implementation follows modern
logging best practices and provides a centralized, configurable logging system.

## 🎯 Overview

The MCP YardGopher project uses Go's built-in `log/slog` package
(introduced in Go 1.21) as the foundation for structured logging. Our
implementation provides:

- **Structured JSON and text output formats**
- **Context-aware logging** with request IDs, user IDs, and trace information
- **Component-specific loggers** for different parts of the application
- **Performance, audit, and security logging capabilities**
- **Configurable log levels and output destinations**
- **Type-safe logging** with strongly-typed attributes

## 🏗 Architecture

### Core Components

```text
internal/infrastructure/logger/
├── logger.go          # Main logger implementation
└── logger_test.go     # Comprehensive unit tests
```

### Logger Structure

```go
type Logger struct {
    *slog.Logger  // Embedded standard slog.Logger
}

type Config struct {
    Level      LogLevel   // debug, info, warn, error
    Format     LogFormat  // json, text  
    Output     string     // stdout, stderr, or file path
    AddSource  bool       // Include source code location
    TimeFormat string     // Custom time formatting
}
```

## 🚀 Usage Examples

### Basic Usage

```go
import "github.com/greysquirr3l/mcp-yardgopher/internal/infrastructure/logger"

// Initialize logger
appLogger := logger.Initialize(&logger.Config{
    Level:  logger.LevelInfo,
    Format: logger.FormatJSON,
    Output: "stdout",
})

// Basic logging
appLogger.Info("Server started", "port", 8080, "version", "1.0.0")
appLogger.Error("Database connection failed", "error", err)
```

### Component-Specific Logging

```go
// Create specialized loggers for different components
dbLogger := appLogger.Database()
httpLogger := appLogger.HTTP()
repoLogger := appLogger.Repository()
commandLogger := appLogger.Command()
queryLogger := appLogger.Query()
domainLogger := appLogger.Domain()
migrationLogger := appLogger.Migration()

// Each adds a "component" field automatically
dbLogger.Info("Connection established")
// Output: {"time":"2025-06-22T19:00:41Z","level":"INFO","msg":"Connection established","component":"database"}
```

### Context-Aware Logging

```go
import "context"

// Add context information
ctx := logger.WithRequestID(context.Background(), "req-123")
ctx = logger.WithUserID(ctx, "user-456")

// Create logger with context
contextLogger := appLogger.WithContext(ctx)
contextLogger.Info("Processing request")
// Output: {"time":"...","level":"INFO","msg":"Processing request","request_id":"req-123","user_id":"user-456"}
```

### Error Logging

```go
// Log with error information
err := errors.New("database timeout")
errorLogger := appLogger.WithError(err)
errorLogger.Error("Operation failed")
// Output: {"time":"...","level":"ERROR","msg":"Operation failed","error":"database timeout"}
```

### Custom Fields

```go
// Add custom structured fields
fields := map[string]interface{}{
    "operation": "create_memory",
    "duration":  "150ms",
    "count":     42,
}
fieldsLogger := appLogger.WithFields(fields)
fieldsLogger.Info("Operation completed")
```

### Performance Logging

```go
import "time"

// Log performance metrics
duration := 150 * time.Millisecond
appLogger.Performance("database_query", duration, 
    slog.String("query", "SELECT * FROM memories"),
    slog.Int("result_count", 25))
```

### Audit Logging

```go
// Log audit events
appLogger.Audit("CREATE", "memory", "user-123", 
    slog.String("memory_id", "mem-456"),
    slog.String("content_type", "knowledge"))
```

### Security Logging

```go
// Log security events
appLogger.Security("UNAUTHORIZED_ACCESS", "HIGH",
    slog.String("ip", "192.168.1.100"),
    slog.String("endpoint", "/admin/memories"))
```

## ⚙️ Configuration

### Environment Variables

```bash
# Log level (debug, info, warn, error)
export LOG_LEVEL=info

# Output format (json, text)
export LOG_FORMAT=json

# Output destination (stdout, stderr, /path/to/file)
export LOG_OUTPUT=stdout

# Include source code location (true, false)
export LOG_ADD_SOURCE=false

# Custom time format
export LOG_TIME_FORMAT="2006-01-02T15:04:05.000Z07:00"
```

### Programmatic Configuration

```go
config := &logger.Config{
    Level:      logger.LevelInfo,
    Format:     logger.FormatJSON,
    Output:     "stdout",
    AddSource:  false,
    TimeFormat: time.RFC3339Nano,
}

appLogger := logger.Initialize(config)
```

### Integration with Database Config

```go
// Automatically configure logger from database config + environment
dbConfig := config.DefaultDatabaseConfig()
appLogger := logger.NewFromDatabase(&dbConfig)
```

## 📊 Output Formats

### JSON Format (Default)

```json
{
  "time": "2025-06-22T19:00:41.123456789Z",
  "level": "INFO",
  "msg": "Memory created successfully",
  "component": "domain",
  "request_id": "req-abc123",
  "user_id": "user-456",
  "memory_id": "mem-789",
  "type": "knowledge",
  "tags": ["go", "architecture"],
  "duration": "25ms"
}
```

### Text Format (Development)

```bash
time=2025-06-22T19:00:41.123Z level=INFO msg="Memory created successfully" component=domain request_id=req-abc123 user_id=user-456 memory_id=mem-789 type=knowledge
```

## 🎭 Context Integration

### HTTP Middleware Example

```go
func LoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestID := generateRequestID()
            ctx := logger.WithRequestID(r.Context(), requestID)
            
            httpLogger := logger.WithContext(ctx).HTTP()
            httpLogger.Info("Request started",
                "method", r.Method,
                "path", r.URL.Path)
            
            // Process request
            next.ServeHTTP(w, r.WithContext(ctx))
            
            httpLogger.Info("Request completed")
        })
    }
}
```

### Repository Pattern Integration

```go
type GormMemoryRepository struct {
    db     *gorm.DB
    logger *logger.Logger
}

func (r *GormMemoryRepository) Create(ctx context.Context, memory *memory.MemoryEntry) error {
    logger := r.logger.WithContext(ctx).Repository()
    
    start := time.Now()
    err := r.db.WithContext(ctx).Create(memory).Error
    duration := time.Since(start)
    
    if err != nil {
        logger.Error("Failed to create memory", "error", err)
        return err
    }
    
    logger.Performance("create_memory", duration,
        slog.String("memory_id", memory.ID.String()),
        slog.String("type", string(memory.Type)))
    
    return nil
}
```

### Command Handler Integration

```go
type CreateMemoryHandler struct {
    repository memory.Repository
    logger     *logger.Logger
}

func (h *CreateMemoryHandler) Handle(ctx context.Context, cmd CreateMemoryCommand) error {
    cmdLogger := h.logger.WithContext(ctx).Command()
    
    cmdLogger.Info("Processing create memory command",
        "type", cmd.Type,
        "project_path", cmd.ProjectPath)
    
    // ... business logic ...
    
    cmdLogger.Audit("CREATE", "memory", userID,
        slog.String("memory_id", memoryID))
    
    return nil
}
```

## 🛠 Best Practices

### 1. Use Appropriate Log Levels

```go
// DEBUG: Detailed information for debugging
logger.Debug("Query parameters validated", "params", params)

// INFO: General information about application flow
logger.Info("Memory created successfully", "id", memoryID)

// WARN: Something unexpected but recoverable
logger.Warn("Rate limit approaching", "current", 95, "limit", 100)

// ERROR: Error conditions that should be investigated
logger.Error("Database query failed", "error", err, "query", sqlQuery)
```

### 2. Structure Your Log Data

```go
// Good: Structured, searchable fields
logger.Info("User action",
    "action", "create_memory",
    "user_id", userID,
    "memory_type", memoryType,
    "duration_ms", duration.Milliseconds())

// Avoid: Unstructured string interpolation
logger.Info(fmt.Sprintf("User %s created memory of type %s", userID, memoryType))
```

### 3. Use Context Consistently

```go
// Pass context through the call chain
func (s *MemoryService) CreateMemory(ctx context.Context, cmd CreateMemoryCommand) error {
    logger := s.logger.WithContext(ctx).WithComponent("service")
    
    // Context is available for correlation
    logger.Info("Creating memory", "type", cmd.Type)
    
    return s.repository.Create(ctx, memory) // Pass context down
}
```

### 4. Component Separation

```go
// Create component-specific loggers early
type MemoryRepository struct {
    db     *gorm.DB
    logger *logger.Logger
}

func NewMemoryRepository(db *gorm.DB, appLogger *logger.Logger) *MemoryRepository {
    return &MemoryRepository{
        db:     db,
        logger: appLogger.Repository(), // Component-specific logger
    }
}
```

### 5. Performance Logging

```go
func (r *Repository) expensiveOperation(ctx context.Context) error {
    start := time.Now()
    defer func() {
        r.logger.WithContext(ctx).Performance("expensive_operation", time.Since(start))
    }()
    
    // ... operation ...
    return nil
}
```

## 🔍 Observability Integration

### Correlation IDs

```go
// Generate and propagate correlation IDs
func generateRequestID() string {
    return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), randomString(8))
}

// Use in HTTP handlers
ctx = logger.WithRequestID(r.Context(), generateRequestID())
```

### Distributed Tracing

```go
// Integration with OpenTelemetry
import "go.opentelemetry.io/otel/trace"

func logWithTrace(ctx context.Context, logger *logger.Logger, msg string) {
    span := trace.SpanFromContext(ctx)
    if span.SpanContext().IsValid() {
        traceLogger := logger.WithFields(map[string]interface{}{
            "trace_id": span.SpanContext().TraceID().String(),
            "span_id":  span.SpanContext().SpanID().String(),
        })
        traceLogger.Info(msg)
    } else {
        logger.Info(msg)
    }
}
```

## 📈 Monitoring and Alerting

### Error Rate Monitoring

```go
// Log errors with consistent structure for monitoring
logger.Error("Critical operation failed",
    "operation", "create_memory",
    "error_type", "database_timeout",
    "severity", "high",
    "user_id", userID)
```

### Performance Monitoring

```go
// Log performance metrics for monitoring
logger.Performance("api_request", duration,
    slog.String("endpoint", "/memories"),
    slog.String("method", "POST"),
    slog.Int("status_code", 201))
```

## 🧪 Testing

### Test Logger Configuration

```go
func TestLoggerConfiguration(t *testing.T) {
    var buf bytes.Buffer
    
    config := &logger.Config{
        Level:  logger.LevelInfo,
        Format: logger.FormatJSON,
        Output: "stdout",
    }
    
    // Create logger that writes to buffer for testing
    opts := &slog.HandlerOptions{Level: slog.LevelInfo}
    handler := slog.NewJSONHandler(&buf, opts)
    testLogger := &logger.Logger{Logger: slog.New(handler)}
    
    testLogger.Info("test message", "key", "value")
    
    // Assert log output
    var logEntry map[string]interface{}
    err := json.Unmarshal(buf.Bytes(), &logEntry)
    assert.NoError(t, err)
    assert.Equal(t, "test message", logEntry["msg"])
    assert.Equal(t, "value", logEntry["key"])
}
```

## 🚨 Security Considerations

### Sensitive Data Filtering

```go
// Never log sensitive information
logger.Info("User authenticated",
    "user_id", userID,
    "email", maskEmail(email), // Mask PII
    // "password", password,    // NEVER log passwords
)

func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "***@***.***"
    }
    return fmt.Sprintf("%s***@%s", parts[0][:1], parts[1])
}
```

### Audit Trail Compliance

```go
// Maintain audit trails for compliance
logger.Audit("ACCESS", "sensitive_memory", userID,
    slog.String("resource_id", resourceID),
    slog.String("action", "READ"),
    slog.Time("accessed_at", time.Now()),
    slog.String("source_ip", clientIP))
```

## 📚 References

- [Official Go slog documentation](https://pkg.go.dev/log/slog)
- [Go blog: Structured Logging with slog](https://go.dev/blog/slog)
- [Comprehensive slog guide](https://betterstack.com/community/guides/logging/logging-in-go/)
- [Modern structured logging tutorial](https://dev.to/leapcell/gos-slog-modern-structured-logging-made-easy-4e04)
- [Detailed slog guide](https://last9.io/blog/logging-in-go-with-slog-a-detailed-guide/)

## 🔄 Migration from Standard log

If migrating from Go's standard `log` package:

```go
// Before (standard log)
log.Printf("User %s created memory %s", userID, memoryID)

// After (structured slog)
logger.Info("Memory created",
    "user_id", userID,
    "memory_id", memoryID)
```

This provides better searchability, structured queries, and integration with modern observability tools.

---

**Last Updated:** June 22, 2025  
**slog Version:** Go 1.21+  
**Implementation Status:** ✅ Complete
