# Centralized Error Handling in Go: Enterprise Patterns and Best Practices

> *Last updated: September 2025*

This document provides comprehensive guidance for implementing centralized error handling
in Go applications, incorporating patterns from frameworks like Laravel while maintaining
Go's idiomatic approaches and enterprise-grade reliability.

<!-- REF: https://dev.to/hrrydgls/centralized-error-handling-in-go-like-laravel-8je -->
<!-- REF: https://blog.golang.org/error-handling-and-go -->
<!-- REF: https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully -->
<!-- REF: https://pkg.go.dev/errors -->

---

## Table of Contents

1. [Core Principles](#core-principles)
2. [Error Handling Patterns](#error-handling-patterns)
3. [Centralized Response System](#centralized-response-system)
4. [Middleware-Based Error Recovery](#middleware-based-error-recovery)
5. [Domain-Driven Error Design](#domain-driven-error-design)
6. [Production-Ready Implementation](#production-ready-implementation)
7. [Testing Error Scenarios](#testing-error-scenarios)
8. [Performance Considerations](#performance-considerations)
9. [Integration with Clean Architecture](#integration-with-clean-architecture)
10. [Best Practices Summary](#best-practices-summary)

---

## Why Error Handling in Go Matters (And Why You Should Care)

If you've come from languages like Java or Python, you might think error handling in Go is a little... weird. In Go, error handling isn't this passive thing that gets shunted to a 'catch' block where you never look. It's an active part of the program's logic, built to ensure you can spot and fix problems at their source.

### The "Silent Panic" Reality

Consider this production scenario: You ignore an "edge case" error in a file-reading function. It doesn't seem like a big deal at first. Until it escalates into a fatal error in production that throws your system into silent panic mode. You have to halt everything and scramble to pinpoint the source. That one error derails weeks of work.

**Errors are part of Go's DNA for a reason** — they're Go's way of saying, "Hey, deal with this before it gets worse!"

### Errors as Values, Not Exceptions

In Go, **errors aren't exceptions; they're values**. This distinction is key. Unlike other languages where errors are thrown and caught, Go expects you to handle them explicitly:

```go
// Go forces you to deal with errors head-on
result, err := divide(10, 0)
if err != nil {
    // Handle the error explicitly - no hiding in catch blocks
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result)
```

**The upside?** You're always in control. **The downside?** If you don't check for errors, they'll quietly fester in your codebase, waiting to strike when you least expect it.

---

## Core Principles

### Goals of Centralized Error Handling

- **Consistency**: Uniform error responses across all API endpoints
- **Maintainability**: Single place to modify error handling logic
- **Observability**: Centralized logging and monitoring integration
- **Developer Experience**: Simple, Laravel-like error throwing
- **Production Safety**: Proper error sanitization and security
- **User Guidance**: Move beyond "Error: Something Broke" to meaningful messages
- **Debugging Intelligence**: Create breadcrumb trails for faster issue resolution

### Go-Idiomatic Approach

Unlike exceptions in other languages, Go's error handling should remain explicit
while providing convenient centralization patterns:

```go
// ❌ Avoid: Heavy reliance on panic/recover for normal error flow
if err != nil {
    panic(CustomError{Code: 400, Message: "Bad request"})
}

// ❌ Avoid: The "Blank Error Trap" - ignoring errors
_, err := riskyFunction()
// Don't just ignore err!

// ✅ Preferred: Explicit error handling with centralized formatting
if err != nil {
    return nil, errors.NewValidationError("invalid input", err)
}

// ✅ Preferred: Always check errors with context
result, err := fetchData()
if err != nil {
    return fmt.Errorf("failed to process user data: %w", err)
}
```

---

## Error Handling Patterns

### The Foundation: Explicit Error Checking

Before diving into advanced patterns, master the fundamentals. Every error is an opportunity to make your system more resilient:

```go
package main

import (
    "errors"
    "fmt"
)

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return // Handle the error, don't ignore it
    }
    fmt.Println("Result:", result)
}
```

**Key insight:** Go forces you to deal with errors head-on. This isn't a burden—it's a feature that keeps your system healthy.

### Pattern 1: Custom Error Types (Because "Error 500" Is Never Enough)

Go's default errors are functional, but they're not particularly helpful. Creating custom error messages helps guide users and debug your code with precision:

```go
package main

import "fmt"

// Custom error type with rich context
type MyError struct {
    When string
    What string
    Code int
}

func (e *MyError) Error() string {
    return fmt.Sprintf("at %s, %s (code: %d)", e.When, e.What, e.Code)
}

func doSomething() error {
    return &MyError{
        When: "Loading Data",
        What: "Connection Timeout", 
        Code: 504,
    }
}

func main() {
    err := doSomething()
    if err != nil {
        fmt.Println("Error occurred:", err)
        // Now you get: "at Loading Data, Connection Timeout (code: 504)"
        // Instead of: "something went wrong"
    }
}
```

### Pattern 2: Error Wrapping (The Sherlock Holmes of Debugging)

Error wrapping adds context at each level of your application, creating a breadcrumb trail that's a lifesaver in debugging:

```go
package main

import (
    "fmt"
    "errors"
)

func fetchData() error {
    return errors.New("failed to fetch data")
}

func processData() error {
    err := fetchData()
    if err != nil {
        return fmt.Errorf("processData: %w", err)
    }
    return nil
}

func handleRequest() error {
    err := processData()
    if err != nil {
        return fmt.Errorf("handleRequest failed: %w", err)
    }
    return nil
}

func main() {
    err := handleRequest()
    if err != nil {
        fmt.Println("Error:", err)
        // Output: "handleRequest failed: processData: failed to fetch data"
        // Perfect breadcrumb trail for debugging!
    }
}
```

When `fetchData` returns an error, `processData` wraps it, creating a detailed error trail. You'll spend less time hunting down issues and more time fixing them.

### Pattern 3: Specific Error Detection with errors.Is() and errors.As()

Avoid the "blank error trap" by using Go's error inspection functions:

```go
package main

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")

type ValidationError struct {
    Field string
    Issue string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s %s", e.Field, e.Issue)
}

func findSomething(id string) error {
    if id == "missing" {
        return ErrNotFound
    }
    if id == "invalid" {
        return ValidationError{Field: "id", Issue: "must be numeric"}
    }
    if id == "forbidden" {
        return fmt.Errorf("access denied: %w", ErrUnauthorized)
    }
    return nil
}

func main() {
    testCases := []string{"missing", "invalid", "forbidden"}
    
    for _, testCase := range testCases {
        err := findSomething(testCase)
        if err != nil {
            // Check for specific error types
            if errors.Is(err, ErrNotFound) {
                fmt.Println("Handle not found case:", err)
            } else if errors.Is(err, ErrUnauthorized) {
                fmt.Println("Handle unauthorized case:", err)
            } else {
                // Check for custom error types
                var valErr ValidationError
                if errors.As(err, &valErr) {
                    fmt.Printf("Handle validation error - Field: %s, Issue: %s\n", 
                        valErr.Field, valErr.Issue)
                } else {
                    fmt.Println("Handle unexpected error:", err)
                }
            }
        }
    }
}
```

By recognizing and handling specific errors, you save users from generic "uh-oh" messages and provide targeted guidance.

### Pattern 4: Error Result Pattern with Centralized Formatting

```go
// pkg/errors/types.go
package errors

import (
    "fmt"
    "net/http"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
    ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
    ErrorTypeNotFound     ErrorType = "NOT_FOUND"
    ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
    ErrorTypeForbidden    ErrorType = "FORBIDDEN"
    ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
    ErrorTypeConflict     ErrorType = "CONFLICT"
    ErrorTypeRateLimit    ErrorType = "RATE_LIMIT"
)

// AppError represents a structured application error
type AppError struct {
    Type        ErrorType `json:"type"`
    Code        int       `json:"code"`
    Message     string    `json:"message"`
    Details     string    `json:"details,omitempty"`
    Cause       error     `json:"-"` // Don't expose internal errors
    RequestID   string    `json:"request_id,omitempty"`
    Field       string    `json:"field,omitempty"` // For validation errors
    Retryable   bool      `json:"retryable,omitempty"`
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Cause
}

func (e *AppError) WithRequestID(requestID string) *AppError {
    e.RequestID = requestID
    return e
}

func (e *AppError) WithField(field string) *AppError {
    e.Field = field
    return e
}
```

### Pattern 2: Error Builders for Common Scenarios

```go
// pkg/errors/builders.go
package errors

import "net/http"

// Validation errors
func NewValidationError(message string, cause error) *AppError {
    return &AppError{
        Type:      ErrorTypeValidation,
        Code:      http.StatusUnprocessableEntity,
        Message:   message,
        Cause:     cause,
        Retryable: false,
    }
}

func NewFieldValidationError(field, message string, cause error) *AppError {
    return &AppError{
        Type:      ErrorTypeValidation,
        Code:      http.StatusUnprocessableEntity,
        Message:   message,
        Field:     field,
        Cause:     cause,
        Retryable: false,
    }
}

// Resource errors
func NewNotFoundError(resource string, cause error) *AppError {
    return &AppError{
        Type:      ErrorTypeNotFound,
        Code:      http.StatusNotFound,
        Message:   fmt.Sprintf("%s not found", resource),
        Cause:     cause,
        Retryable: false,
    }
}

func NewConflictError(message string, cause error) *AppError {
    return &AppError{
        Type:      ErrorTypeConflict,
        Code:      http.StatusConflict,
        Message:   message,
        Cause:     cause,
        Retryable: false,
    }
}

// Authentication/Authorization errors
func NewUnauthorizedError(message string) *AppError {
    return &AppError{
        Type:      ErrorTypeUnauthorized,
        Code:      http.StatusUnauthorized,
        Message:   message,
        Retryable: false,
    }
}

func NewForbiddenError(message string) *AppError {
    return &AppError{
        Type:      ErrorTypeForbidden,
        Code:      http.StatusForbidden,
        Message:   message,
        Retryable: false,
    }
}

// System errors
func NewInternalError(message string, cause error) *AppError {
    return &AppError{
        Type:      ErrorTypeInternal,
        Code:      http.StatusInternalServerError,
        Message:   "Internal server error", // Don't expose internal details
        Details:   message,                 // For logging only
        Cause:     cause,
        Retryable: true,
    }
}

func NewRateLimitError(message string) *AppError {
    return &AppError{
        Type:      ErrorTypeRateLimit,
        Code:      http.StatusTooManyRequests,
        Message:   message,
        Retryable: true,
    }
}
```

---

## Centralized Response System

### Response Types

```go
// pkg/responses/types.go
package responses

import (
    "encoding/json"
    "time"
)

// APIResponse represents the standard API response format
type APIResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Error     *ErrorInfo  `json:"error,omitempty"`
    Meta      *Meta       `json:"meta,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
    RequestID string      `json:"request_id,omitempty"`
}

// ErrorInfo represents error information in API responses
type ErrorInfo struct {
    Type      string      `json:"type"`
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Field     string      `json:"field,omitempty"`
    Details   interface{} `json:"details,omitempty"`
    Retryable bool        `json:"retryable,omitempty"`
}

// Meta contains pagination and additional metadata
type Meta struct {
    Page       int `json:"page,omitempty"`
    PerPage    int `json:"per_page,omitempty"`
    Total      int `json:"total,omitempty"`
    TotalPages int `json:"total_pages,omitempty"`
}

// ValidationErrorDetail represents detailed validation error information
type ValidationErrorDetail struct {
    Field   string `json:"field"`
    Value   string `json:"value,omitempty"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

// MultiFieldValidationError for multiple validation errors
type MultiFieldValidationError struct {
    Errors []ValidationErrorDetail `json:"errors"`
}
```

### Response Builder

```go
// pkg/responses/builder.go
package responses

import (
    "context"
    "net/http"
    "time"

    apperrors "yourapp/pkg/errors"
)

type ResponseBuilder struct {
    requestID string
}

func NewResponseBuilder(ctx context.Context) *ResponseBuilder {
    requestID, _ := ctx.Value("request_id").(string)
    return &ResponseBuilder{requestID: requestID}
}

// Success responses
func (rb *ResponseBuilder) Success(data interface{}) *APIResponse {
    return &APIResponse{
        Success:   true,
        Data:      data,
        Timestamp: time.Now(),
        RequestID: rb.requestID,
    }
}

func (rb *ResponseBuilder) SuccessWithMeta(data interface{}, meta *Meta) *APIResponse {
    return &APIResponse{
        Success:   true,
        Data:      data,
        Meta:      meta,
        Timestamp: time.Now(),
        RequestID: rb.requestID,
    }
}

// Error responses
func (rb *ResponseBuilder) Error(err error) *APIResponse {
    var appErr *apperrors.AppError
    var errorInfo *ErrorInfo

    if errors.As(err, &appErr) {
        errorInfo = &ErrorInfo{
            Type:      string(appErr.Type),
            Code:      appErr.Code,
            Message:   appErr.Message,
            Field:     appErr.Field,
            Retryable: appErr.Retryable,
        }

        // Add details for validation errors in development
        if appErr.Type == apperrors.ErrorTypeValidation && appErr.Details != "" {
            errorInfo.Details = appErr.Details
        }
    } else {
        // Fallback for unexpected errors
        errorInfo = &ErrorInfo{
            Type:      string(apperrors.ErrorTypeInternal),
            Code:      http.StatusInternalServerError,
            Message:   "Internal server error",
            Retryable: true,
        }
    }

    return &APIResponse{
        Success:   false,
        Error:     errorInfo,
        Timestamp: time.Now(),
        RequestID: rb.requestID,
    }
}

func (rb *ResponseBuilder) ValidationError(errors []ValidationErrorDetail) *APIResponse {
    return &APIResponse{
        Success: false,
        Error: &ErrorInfo{
            Type:    string(apperrors.ErrorTypeValidation),
            Code:    http.StatusUnprocessableEntity,
            Message: "Validation failed",
            Details: MultiFieldValidationError{Errors: errors},
        },
        Timestamp: time.Now(),
        RequestID: rb.requestID,
    }
}
```

---

## Middleware-Based Error Recovery

### HTTP Error Middleware

```go
// pkg/middleware/error_handling.go
package middleware

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"
    "runtime/debug"

    apperrors "yourapp/pkg/errors"
    "yourapp/pkg/responses"
)

// ErrorHandlingMiddleware provides centralized error handling
type ErrorHandlingMiddleware struct {
    logger   *slog.Logger
    isDev    bool
    notifier ErrorNotifier // For production error notifications
}

type ErrorNotifier interface {
    NotifyError(ctx context.Context, err error, requestInfo map[string]interface{})
}

func NewErrorHandlingMiddleware(logger *slog.Logger, isDev bool, notifier ErrorNotifier) *ErrorHandlingMiddleware {
    return &ErrorHandlingMiddleware{
        logger:   logger,
        isDev:    isDev,
        notifier: notifier,
    }
}

func (m *ErrorHandlingMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                m.handlePanic(w, r, rec)
            }
        }()

        // Create a custom response writer to capture errors
        cw := &captureWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(cw, r)
    })
}

// captureWriter captures response status for logging
type captureWriter struct {
    http.ResponseWriter
    statusCode int
}

func (cw *captureWriter) WriteHeader(code int) {
    cw.statusCode = code
    cw.ResponseWriter.WriteHeader(code)
}

func (m *ErrorHandlingMiddleware) handlePanic(w http.ResponseWriter, r *http.Request, rec interface{}) {
    ctx := r.Context()
    requestID := getRequestID(ctx)

    // Log the panic with stack trace
    stack := debug.Stack()
    m.logger.Error("Panic recovered",
        "error", rec,
        "request_id", requestID,
        "method", r.Method,
        "path", r.URL.Path,
        "stack", string(stack),
    )

    var err error
    var appErr *apperrors.AppError

    // Try to convert panic to AppError
    switch v := rec.(type) {
    case *apperrors.AppError:
        appErr = v
        err = v
    case apperrors.AppError:
        appErr = &v
        err = &v
    case error:
        err = v
        appErr = apperrors.NewInternalError("Unexpected error occurred", v)
    case string:
        err = fmt.Errorf("panic: %s", v)
        appErr = apperrors.NewInternalError("Unexpected error occurred", err)
    default:
        err = fmt.Errorf("panic: %v", v)
        appErr = apperrors.NewInternalError("Unexpected error occurred", err)
    }

    // Set request ID for error tracking
    if requestID != "" {
        appErr = appErr.WithRequestID(requestID)
    }

    // Notify error monitoring system for internal errors
    if appErr.Type == apperrors.ErrorTypeInternal && m.notifier != nil {
        requestInfo := map[string]interface{}{
            "method":     r.Method,
            "path":       r.URL.Path,
            "user_agent": r.UserAgent(),
            "ip":         getRealIP(r),
        }
        go m.notifier.NotifyError(ctx, err, requestInfo)
    }

    // Build and send error response
    responseBuilder := responses.NewResponseBuilder(ctx)
    response := responseBuilder.Error(appErr)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(appErr.Code)
    
    if err := json.NewEncoder(w).Encode(response); err != nil {
        m.logger.Error("Failed to encode error response", "error", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func getRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value("request_id").(string); ok {
        return requestID
    }
    return ""
}

func getRealIP(r *http.Request) string {
    if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
        return ip
    }
    if ip := r.Header.Get("X-Real-IP"); ip != "" {
        return ip
    }
    return r.RemoteAddr
}
```

### Result Pattern Alternative (Preferred for Go)

Instead of relying heavily on panic/recover, consider the Result pattern:

```go
// pkg/result/result.go
package result

import (
    apperrors "yourapp/pkg/errors"
)

// Result represents either success with data or failure with error
type Result[T any] struct {
    data  T
    error *apperrors.AppError
}

func Success[T any](data T) Result[T] {
    return Result[T]{data: data}
}

func Failure[T any](err *apperrors.AppError) Result[T] {
    var zero T
    return Result[T]{data: zero, error: err}
}

func (r Result[T]) IsSuccess() bool {
    return r.error == nil
}

func (r Result[T]) IsFailure() bool {
    return r.error != nil
}

func (r Result[T]) Data() T {
    return r.data
}

func (r Result[T]) Error() *apperrors.AppError {
    return r.error
}

// Unwrap returns data and error for traditional Go error handling
func (r Result[T]) Unwrap() (T, error) {
    if r.error != nil {
        return r.data, r.error
    }
    return r.data, nil
}

// Map transforms the data if result is successful
func Map[T, U any](r Result[T], fn func(T) U) Result[U] {
    if r.IsFailure() {
        return Failure[U](r.error)
    }
    return Success(fn(r.data))
}

// FlatMap chains operations that return Results
func FlatMap[T, U any](r Result[T], fn func(T) Result[U]) Result[U] {
    if r.IsFailure() {
        return Failure[U](r.error)
    }
    return fn(r.data)
}
```

---

## Domain-Driven Error Design

### Domain-Specific Errors

```go
// internal/domain/user/errors.go
package user

import (
    apperrors "yourapp/pkg/errors"
)

// Domain-specific error constructors
func ErrUserNotFound(userID string, cause error) *apperrors.AppError {
    return apperrors.NewNotFoundError("User", cause).
        WithField("user_id").
        WithDetails(map[string]interface{}{"user_id": userID})
}

func ErrInvalidEmail(email string, cause error) *apperrors.AppError {
    return apperrors.NewFieldValidationError("email", "Invalid email format", cause).
        WithDetails(map[string]interface{}{"provided_email": email})
}

func ErrEmailAlreadyExists(email string) *apperrors.AppError {
    return apperrors.NewConflictError("Email already registered", nil).
        WithField("email").
        WithDetails(map[string]interface{}{"email": email})
}

func ErrPasswordTooWeak(reason string) *apperrors.AppError {
    return apperrors.NewFieldValidationError("password", 
        "Password does not meet security requirements", nil).
        WithDetails(map[string]interface{}{"reason": reason})
}

func ErrInsufficientPermissions(action, resource string) *apperrors.AppError {
    return apperrors.NewForbiddenError("Insufficient permissions").
        WithDetails(map[string]interface{}{
            "action":   action,
            "resource": resource,
        })
}
```

### Application Service Error Handling

```go
// internal/application/user/service.go
package user

import (
    "context"

    "yourapp/internal/domain/user"
    "yourapp/pkg/result"
)

type UserService struct {
    userRepo user.Repository
    emailSvc EmailService
}

func (s *UserService) CreateUser(ctx context.Context, cmd CreateUserCommand) result.Result[*user.User] {
    // Validation
    if cmd.Email == "" {
        return result.Failure[*user.User](
            user.ErrInvalidEmail(cmd.Email, nil))
    }

    // Check if user exists
    existing, err := s.userRepo.FindByEmail(ctx, cmd.Email)
    if err != nil && !errors.Is(err, user.ErrNotFound) {
        return result.Failure[*user.User](
            apperrors.NewInternalError("Failed to check existing user", err))
    }
    
    if existing != nil {
        return result.Failure[*user.User](
            user.ErrEmailAlreadyExists(cmd.Email))
    }

    // Create domain object
    newUser, err := user.NewUser(cmd.Email, cmd.Password, cmd.Name)
    if err != nil {
        // Domain validation errors are already properly formatted
        var appErr *apperrors.AppError
        if errors.As(err, &appErr) {
            return result.Failure[*user.User](appErr)
        }
        return result.Failure[*user.User](
            apperrors.NewValidationError("Invalid user data", err))
    }

    // Save user
    if err := s.userRepo.Save(ctx, newUser); err != nil {
        return result.Failure[*user.User](
            apperrors.NewInternalError("Failed to save user", err))
    }

    // Send welcome email (non-blocking)
    go s.emailSvc.SendWelcomeEmail(ctx, newUser.Email(), newUser.Name())

    return result.Success(newUser)
}
```

---

## Production-Ready Implementation

### HTTP Handler with Centralized Error Handling

```go
// internal/adapters/primary/http/user_handler.go
package http

import (
    "encoding/json"
    "net/http"

    "yourapp/internal/application/user"
    apperrors "yourapp/pkg/errors"
    "yourapp/pkg/responses"
)

type UserHandler struct {
    userService *user.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    responseBuilder := responses.NewResponseBuilder(ctx)

    var cmd user.CreateUserCommand
    if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
        errorResponse := responseBuilder.Error(
            apperrors.NewValidationError("Invalid JSON payload", err))
        h.sendJSONResponse(w, http.StatusUnprocessableEntity, errorResponse)
        return
    }

    // Use result pattern
    result := h.userService.CreateUser(ctx, cmd)
    userData, err := result.Unwrap()
    
    if err != nil {
        errorResponse := responseBuilder.Error(err)
        var appErr *apperrors.AppError
        statusCode := http.StatusInternalServerError
        
        if errors.As(err, &appErr) {
            statusCode = appErr.Code
        }
        
        h.sendJSONResponse(w, statusCode, errorResponse)
        return
    }

    // Success response
    successResponse := responseBuilder.Success(map[string]interface{}{
        "user_id": userData.ID(),
        "email":   userData.Email(),
        "name":    userData.Name(),
    })
    
    h.sendJSONResponse(w, http.StatusCreated, successResponse)
}

func (h *UserHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, response *responses.APIResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    if err := json.NewEncoder(w).Encode(response); err != nil {
        // Log error and send basic error response
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

// Alternative: Using panic pattern (Laravel-style) - use sparingly
func (h *UserHandler) CreateUserWithPanic(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    var cmd user.CreateUserCommand
    if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
        panic(apperrors.NewValidationError("Invalid JSON payload", err))
    }

    result := h.userService.CreateUser(ctx, cmd)
    userData, err := result.Unwrap()
    
    if err != nil {
        panic(err) // Middleware will handle this
    }

    // Success response
    responseBuilder := responses.NewResponseBuilder(ctx)
    successResponse := responseBuilder.Success(map[string]interface{}{
        "user_id": userData.ID(),
        "email":   userData.Email(),
        "name":    userData.Name(),
    })
    
    h.sendJSONResponse(w, http.StatusCreated, successResponse)
}
```

### Error Monitoring Integration

```go
// pkg/monitoring/error_notifier.go
package monitoring

import (
    "context"
    "log/slog"
    
    "github.com/getsentry/sentry-go"
)

type SentryErrorNotifier struct {
    logger *slog.Logger
    hub    *sentry.Hub
}

func NewSentryErrorNotifier(logger *slog.Logger, hub *sentry.Hub) *SentryErrorNotifier {
    return &SentryErrorNotifier{
        logger: logger,
        hub:    hub,
    }
}

func (n *SentryErrorNotifier) NotifyError(ctx context.Context, err error, requestInfo map[string]interface{}) {
    // Log locally
    n.logger.Error("Application error occurred",
        "error", err.Error(),
        "request_info", requestInfo,
    )

    // Send to Sentry
    n.hub.WithScope(func(scope *sentry.Scope) {
        scope.SetTag("error_type", "application_error")
        
        for key, value := range requestInfo {
            scope.SetTag(key, fmt.Sprintf("%v", value))
        }
        
        n.hub.CaptureException(err)
    })
}
```

---

## Testing Error Scenarios

### Unit Tests for Error Handling

```go
// internal/application/user/service_test.go
package user_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "yourapp/internal/application/user"
    domainuser "yourapp/internal/domain/user"
    apperrors "yourapp/pkg/errors"
)

func TestUserService_CreateUser_ValidationErrors(t *testing.T) {
    tests := []struct {
        name        string
        command     user.CreateUserCommand
        expectedErr *apperrors.AppError
    }{
        {
            name: "empty email",
            command: user.CreateUserCommand{
                Email:    "",
                Password: "validpass123",
                Name:     "John Doe",
            },
            expectedErr: &apperrors.AppError{
                Type: apperrors.ErrorTypeValidation,
                Code: 422,
            },
        },
        {
            name: "invalid email format",
            command: user.CreateUserCommand{
                Email:    "invalid-email",
                Password: "validpass123",
                Name:     "John Doe",
            },
            expectedErr: &apperrors.AppError{
                Type:  apperrors.ErrorTypeValidation,
                Code:  422,
                Field: "email",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockRepo := &MockUserRepository{}
            service := user.NewUserService(mockRepo, nil)

            // Execute
            result := service.CreateUser(context.Background(), tt.command)

            // Assert
            assert.True(t, result.IsFailure())
            
            _, err := result.Unwrap()
            assert.Error(t, err)
            
            var appErr *apperrors.AppError
            assert.True(t, errors.As(err, &appErr))
            assert.Equal(t, tt.expectedErr.Type, appErr.Type)
            assert.Equal(t, tt.expectedErr.Code, appErr.Code)
            
            if tt.expectedErr.Field != "" {
                assert.Equal(t, tt.expectedErr.Field, appErr.Field)
            }
        })
    }
}

func TestUserService_CreateUser_ConflictError(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    existingUser := &domainuser.User{} // Mock existing user
    mockRepo.On("FindByEmail", mock.Anything, "existing@example.com").
        Return(existingUser, nil)
    
    service := user.NewUserService(mockRepo, nil)

    cmd := user.CreateUserCommand{
        Email:    "existing@example.com",
        Password: "validpass123",
        Name:     "John Doe",
    }

    // Execute
    result := service.CreateUser(context.Background(), cmd)

    // Assert
    assert.True(t, result.IsFailure())
    
    _, err := result.Unwrap()
    var appErr *apperrors.AppError
    assert.True(t, errors.As(err, &appErr))
    assert.Equal(t, apperrors.ErrorTypeConflict, appErr.Type)
    assert.Equal(t, 409, appErr.Code)
}
```

### HTTP Handler Tests

```go
// internal/adapters/primary/http/user_handler_test.go
package http_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"

    httphandlers "yourapp/internal/adapters/primary/http"
    "yourapp/pkg/responses"
)

func TestUserHandler_CreateUser_ValidationError(t *testing.T) {
    // Setup
    handler := setupUserHandler() // Your setup function
    
    invalidPayload := map[string]interface{}{
        "email": "invalid-email",
        "password": "short",
    }
    
    body, _ := json.Marshal(invalidPayload)
    req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    recorder := httptest.NewRecorder()

    // Execute
    handler.CreateUser(recorder, req)

    // Assert
    assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
    assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
    
    var response responses.APIResponse
    err := json.NewDecoder(recorder.Body).Decode(&response)
    assert.NoError(t, err)
    
    assert.False(t, response.Success)
    assert.NotNil(t, response.Error)
    assert.Equal(t, "VALIDATION_ERROR", response.Error.Type)
    assert.Equal(t, 422, response.Error.Code)
}

func TestUserHandler_CreateUser_InternalError(t *testing.T) {
    // Setup handler with mock service that returns internal error
    // ... test internal error handling
}
```

---

## Performance Considerations

### Efficient Error Handling

```go
// Avoid allocating errors in hot paths
var (
    ErrUserNotFound = apperrors.NewNotFoundError("User", nil)
    ErrInvalidEmail = apperrors.NewValidationError("Invalid email", nil)
)

// Use error wrapping instead of creating new errors
func (s *UserService) validateEmail(email string) error {
    if email == "" {
        return fmt.Errorf("email validation: %w", ErrInvalidEmail)
    }
    return nil
}
```

### Error Pooling for High-Traffic Applications

```go
// pkg/errors/pool.go
package errors

import "sync"

var errorPool = sync.Pool{
    New: func() interface{} {
        return &AppError{}
    },
}

func NewPooledError(errType ErrorType, code int, message string) *AppError {
    err := errorPool.Get().(*AppError)
    err.Type = errType
    err.Code = code
    err.Message = message
    err.Cause = nil
    err.Details = ""
    err.Field = ""
    err.RequestID = ""
    err.Retryable = false
    return err
}

func (e *AppError) Release() {
    errorPool.Put(e)
}
```

---

## Integration with Clean Architecture

### Repository Layer Error Handling

```go
// internal/adapters/secondary/postgres/user_repository.go
package postgres

import (
    "context"
    "database/sql"
    "errors"

    "github.com/lib/pq"
    
    domainuser "yourapp/internal/domain/user"
    apperrors "yourapp/pkg/errors"
)

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Save(ctx context.Context, user *domainuser.User) error {
    query := `
        INSERT INTO users (id, email, name, password_hash, created_at) 
        VALUES ($1, $2, $3, $4, $5)`
    
    _, err := r.db.ExecContext(ctx, query, 
        user.ID(), user.Email(), user.Name(), user.PasswordHash(), user.CreatedAt())
    
    if err != nil {
        // Convert database errors to domain errors
        var pqErr *pq.Error
        if errors.As(err, &pqErr) {
            switch pqErr.Code {
            case "23505": // unique_violation
                if pqErr.Constraint == "users_email_unique" {
                    return domainuser.ErrEmailAlreadyExists(user.Email().String())
                }
            case "23514": // check_violation
                return apperrors.NewValidationError("Data constraint violation", err)
            }
        }
        
        // Generic database error
        return apperrors.NewInternalError("Failed to save user", err)
    }
    
    return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
    query := `SELECT id, email, name, password_hash, created_at FROM users WHERE email = $1`
    
    row := r.db.QueryRowContext(ctx, query, email)
    
    var user domainuser.User
    err := row.Scan(&user.id, &user.email, &user.name, &user.passwordHash, &user.createdAt)
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domainuser.ErrUserNotFound(email, err)
        }
        return nil, apperrors.NewInternalError("Failed to find user", err)
    }
    
    return &user, nil
}
```

---

## Best Practices Summary

### Do's ✅

1. **Use explicit error handling** over panic/recover for normal business logic
2. **Create domain-specific errors** with meaningful messages
3. **Centralize error response formatting** in middleware
4. **Include request IDs** for error tracking
5. **Log errors with context** for debugging
6. **Sanitize errors** before sending to clients
7. **Use structured error types** for programmatic handling
8. **Test error scenarios** thoroughly

### Don'ts ❌

1. **Don't expose internal system details** in error messages
2. **Don't use panic/recover for control flow** (only for truly exceptional cases)
3. **Don't create errors in hot paths** without considering performance
4. **Don't forget to wrap errors** for better context
5. **Don't ignore errors** or fail silently
6. **Don't couple error handling** to specific HTTP frameworks
7. **Don't return both data and error** from functions
8. **Don't use magic status codes** without proper error types

## Know When to Fail Fast (And When to Recover Gracefully)

In Go, there's a balance between handling every error and knowing when to throw in the towel. Understanding this balance is crucial for building resilient systems.

### Fail Fast Scenarios

There are times when failing fast is necessary—when the system detects a fundamental issue:

```go
package main

import (
    "fmt"
    "log"
    "os"
)

func initializeApp() error {
    // Configuration errors should fail fast
    configFile := os.Getenv("CONFIG_FILE")
    if configFile == "" {
        return fmt.Errorf("CONFIG_FILE environment variable is required")
    }
    
    // Essential data missing should fail fast
    _, err := os.Stat(configFile)
    if err != nil {
        return fmt.Errorf("configuration file not found: %w", err)
    }
    
    return nil
}

func main() {
    if err := initializeApp(); err != nil {
        log.Fatal("Application initialization failed:", err)
        // Exits immediately - no point continuing without essential config
    }
    
    fmt.Println("Application started successfully")
}
```

### Graceful Recovery with recover()

For recoverable errors (like transient network issues), consider graceful recovery using Go's `recover()` function:

```go
package main

import (
    "fmt"
    "log"
    "time"
)

func processWithRecovery() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
            // Log the panic but continue serving other requests
            // Use recover() wisely - easy to misuse and ignore serious problems
        }
    }()
    
    // Potentially dangerous operation
    riskyOperation()
}

func riskyOperation() {
    // Simulate a panic that we can recover from
    panic("temporary processing error")
}

func handleNetworkError(err error) error {
    // For network errors, retry with exponential backoff
    retryable := []string{"timeout", "connection refused", "temporary failure"}
    
    for _, pattern := range retryable {
        if contains(err.Error(), pattern) {
            log.Printf("Retryable error detected: %v", err)
            return performRetry()
        }
    }
    
    // Non-retryable error
    return fmt.Errorf("permanent failure: %w", err)
}

func performRetry() error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        time.Sleep(time.Duration(i) * time.Second) // Simple backoff
        
        err := attemptOperation()
        if err == nil {
            return nil // Success!
        }
        
        log.Printf("Retry %d failed: %v", i+1, err)
    }
    
    return fmt.Errorf("operation failed after %d retries", maxRetries)
}

func attemptOperation() error {
    // Simulate operation that might fail
    return nil
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && s[:len(substr)] == substr
}

func main() {
    // Demonstrate graceful recovery
    processWithRecovery()
    fmt.Println("Application continues running after recovery")
}
```

### The Recovery Decision Matrix

| Error Type | Response Strategy | Rationale |
|------------|------------------|-----------|
| **Configuration Missing** | Fail Fast | App cannot function without config |
| **Database Connection** | Retry + Fail Fast | Retry briefly, then fail if persistent |
| **Network Timeout** | Retry + Graceful | Often transient, worth retrying |
| **Invalid User Input** | Handle + Continue | Expected scenario, inform user |
| **Memory Exhaustion** | Graceful Shutdown | Attempt cleanup, then fail fast |
| **Authentication Failure** | Handle + Continue | Security concern, log and deny |

---

## Turning Errors into Allies: A Production Philosophy

Error handling in Go isn't about avoiding trouble—it's about facing it head-on with a solid game plan. By handling errors explicitly, creating custom error messages, wrapping errors, and using error-checking functions, you're ensuring that your Go app is prepared to tackle whatever comes its way.

### The Error Handling Mindset Shift

**Stop thinking of errors as failures.** Start thinking of them as:
- **Information sources** about system state
- **User guidance opportunities** for better UX  
- **Debugging breadcrumbs** for faster resolution
- **System health indicators** for monitoring
- **Security checkpoints** for access control

### The "Debugging Insurance Policy"

Stop sweeping errors under the rug. Start making error handling a top priority. Think of it as a **debugging insurance policy**—one that saves you time, headache, and frustration down the line.

Whether you're a Go beginner or seasoned developer, mastering error handling is an investment that pays for itself every time something goes wrong. **Make error handling an ally, not an enemy.**

---

## References

1. [Go Blog: Error handling and Go](https://blog.golang.org/error-handling-and-go)
2. [Dave Cheney: Don't just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
3. [Centralized Error Handling in Go Like Laravel](https://dev.to/hrrydgls/centralized-error-handling-in-go-like-laravel-8je)
4. [Go errors package](https://pkg.go.dev/errors)
5. [Effective Go: Errors](https://golang.org/doc/effective_go.html#errors)
6. [Go Wiki: Error Handling](https://github.com/golang/go/wiki/ErrorHandling)
7. [Uber Go Style Guide: Error Handling](https://github.com/uber-go/guide/blob/master/style.md#error-handling)

---

*This document provides enterprise-grade error handling patterns for Go applications while maintaining idiomatic Go practices and supporting modern development workflows. Remember: Error handling isn't about avoiding trouble—it's about building systems that gracefully handle the unexpected and turn problems into opportunities for better user experiences and system reliability.*
