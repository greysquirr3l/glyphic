# Serving JSON from Your Go API — Complete Guide 2025

> A comprehensive guide to implementing JSON responses in Go web APIs with modern
> patterns, performance optimizations, and best practices aligned with Clean
> Architecture principles.

## Table of Contents

1. [Why JSON in Go APIs?](#why-json-in-go-apis)
2. [Basic JSON Response Patterns](#basic-json-response-patterns)
3. [Advanced JSON Handling](#advanced-json-handling)
4. [Clean Architecture Integration](#clean-architecture-integration)
5. [Performance Optimizations](#performance-optimizations)
6. [Error Handling Patterns](#error-handling-patterns)
7. [Testing JSON Endpoints](#testing-json-endpoints)
8. [Production Considerations](#production-considerations)

---

## Why JSON in Go APIs?

JSON (JavaScript Object Notation) has become the de facto standard for REST APIs.
In Go, the `encoding/json` package provides excellent performance and ease of use:

**Benefits:**

- **Human-readable** and debugging-friendly
- **Lightweight** compared to XML (typically 30-40% smaller)
- **Native Go support** with struct tags for fine-grained control
- **Type safety** through Go's strong typing system
- **Performance** - Go's JSON encoder is highly optimized

### JSON Structure Recap

```json
{
  "id": 42,
  "title": "The Matrix",
  "genres": ["action", "sci-fi"],
  "released": true,
  "rating": 8.7,
  "metadata": null
}
```

**Key elements:**

- Objects: `{}` with key-value pairs
- Arrays: `[]` with ordered elements  
- Strings: `"quoted text"`
- Numbers: integers or floats
- Booleans: `true` or `false`
- Null: `null`

---

## Basic JSON Response Patterns

### 1. Static JSON Responses

For simple, unchanging responses like health checks:

```go
func (app *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
    // Static JSON for simple endpoints
    response := `{
        "status": "healthy",
        "timestamp": "` + time.Now().Format(time.RFC3339) + `",
        "version": "1.0.0"
    }`
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))
}
```

**When to use:** Health checks, API version info, simple status endpoints.

### 2. Dynamic JSON with json.Marshal

For data-driven responses:

```go
type StatusResponse struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Version   string            `json:"version"`
    Services  map[string]string `json:"services"`
}

func (app *Application) statusHandler(w http.ResponseWriter, r *http.Request) {
    status := StatusResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Version:   app.config.Version,
        Services: map[string]string{
            "database": "connected",
            "cache":    "connected", 
            "queue":    "connected",
        },
    }
    
    data, err := json.Marshal(status)
    if err != nil {
        app.logError(r, err)
        app.errorResponse(w, r, http.StatusInternalServerError, "failed to encode response")
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}
```

### 3. Reusable JSON Helper

Create a centralized helper to reduce boilerplate:

```go
// JSONResponse standardizes JSON responses across the application
func (app *Application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
    // Marshal data to JSON
    js, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    // Add trailing newline for better curl experience
    js = append(js, '\n')
    
    // Set any additional headers
    for key, value := range headers {
        w.Header()[key] = value
    }
    
    // Set content type and status
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    // Write response
    _, err = w.Write(js)
    return err
}

// Usage becomes much cleaner
func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
    user, err := app.models.Users.Get(getUserID(r))
    if err != nil {
        app.errorResponse(w, r, http.StatusNotFound, "user not found")
        return
    }
    
    err = app.writeJSON(w, http.StatusOK, user, nil)
    if err != nil {
        app.logError(r, err)
        app.errorResponse(w, r, http.StatusInternalServerError, "failed to write response")
    }
}
```

---

## Advanced JSON Handling

### 1. Struct Tags for Control

Go's struct tags provide fine-grained control over JSON serialization:

```go
type Movie struct {
    ID          int64     `json:"id"`
    CreatedAt   time.Time `json:"-"`                    // Never serialize
    UpdatedAt   time.Time `json:"updated_at,omitempty"` // Omit if zero
    Title       string    `json:"title"`
    Description *string   `json:"description,omitempty"` // Omit if nil
    Year        int       `json:"year,omitempty"`
    Duration    int       `json:"duration_minutes"`      // Custom field name
    Genres      []string  `json:"genres,omitempty"`
    Rating      float64   `json:"rating,omitempty"`
    
    // Computed fields
    DurationFormatted string `json:"duration_formatted"`
}

// Custom marshaling for computed fields
func (m Movie) MarshalJSON() ([]byte, error) {
    // Create an alias to avoid infinite recursion
    type MovieAlias Movie
    
    return json.Marshal(&struct {
        *MovieAlias
        DurationFormatted string `json:"duration_formatted"`
    }{
        MovieAlias:        (*MovieAlias)(&m),
        DurationFormatted: formatDuration(m.Duration),
    })
}

func formatDuration(minutes int) string {
    hours := minutes / 60
    mins := minutes % 60
    return fmt.Sprintf("%dh %dm", hours, mins)
}
```

### 2. JSON Streaming with json.Encoder

For large responses or when you want to stream directly:

```go
func (app *Application) exportMoviesHandler(w http.ResponseWriter, r *http.Request) {
    movies, err := app.models.Movies.GetAll()
    if err != nil {
        app.errorResponse(w, r, http.StatusInternalServerError, "failed to fetch movies")
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    // Stream JSON directly to response writer
    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "  ") // Pretty print for exports
    
    if err := encoder.Encode(movies); err != nil {
        // Note: Can't change status code after writing starts
        app.logError(r, fmt.Errorf("streaming encode error: %w", err))
    }
}
```

**Caveat:** With `json.Encoder`, if encoding fails partway through, you can't
change the HTTP status code since headers are already sent.

### 3. Custom JSON Types

Handle special formatting requirements:

```go
// Custom time format
type APITime time.Time

func (t APITime) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Time(t).Format("2006-01-02 15:04:05"))
}

func (t *APITime) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    parsed, err := time.Parse("2006-01-02 15:04:05", s)
    if err != nil {
        return err
    }
    *t = APITime(parsed)
    return nil
}

// Usage in structs
type Event struct {
    ID        int64   `json:"id"`
    Name      string  `json:"name"`
    StartTime APITime `json:"start_time"`
    EndTime   APITime `json:"end_time"`
}
```

---

## Clean Architecture Integration

### Domain Layer

Keep domain entities free of JSON concerns:

```go
// Domain entities remain pure
package domain

type Movie struct {
    ID          MovieID
    Title       string
    Description string
    Year        int
    Duration    time.Duration
    Genres      []Genre
    Rating      Rating
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type MovieRepository interface {
    Create(movie Movie) error
    GetByID(id MovieID) (Movie, error)
    List(filters MovieFilters) ([]Movie, error)
    Update(movie Movie) error
    Delete(id MovieID) error
}
```

### Application Layer

Use cases orchestrate domain logic:

```go
package application

type MovieUseCase struct {
    repo   domain.MovieRepository
    logger Logger
}

func (uc *MovieUseCase) GetMovie(ctx context.Context, id string) (domain.Movie, error) {
    movieID, err := domain.ParseMovieID(id)
    if err != nil {
        return domain.Movie{}, ErrInvalidMovieID
    }
    
    movie, err := uc.repo.GetByID(movieID)
    if err != nil {
        uc.logger.Error("failed to get movie", "id", id, "error", err)
        return domain.Movie{}, err
    }
    
    return movie, nil
}
```

### Infrastructure Layer

Handle JSON serialization in HTTP handlers:

```go
package http

// DTO for JSON serialization
type MovieResponse struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`
    Description string   `json:"description,omitempty"`
    Year        int      `json:"year"`
    Duration    string   `json:"duration"`
    Genres      []string `json:"genres"`
    Rating      float64  `json:"rating"`
    CreatedAt   string   `json:"created_at"`
    UpdatedAt   string   `json:"updated_at"`
}

// Convert domain model to response DTO
func toMovieResponse(movie domain.Movie) MovieResponse {
    return MovieResponse{
        ID:          movie.ID.String(),
        Title:       movie.Title,
        Description: movie.Description,
        Year:        movie.Year,
        Duration:    movie.Duration.String(),
        Genres:      genresToStrings(movie.Genres),
        Rating:      float64(movie.Rating),
        CreatedAt:   movie.CreatedAt.Format(time.RFC3339),
        UpdatedAt:   movie.UpdatedAt.Format(time.RFC3339),
    }
}

type MovieHandler struct {
    useCase *application.MovieUseCase
    app     *Application
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    movie, err := h.useCase.GetMovie(r.Context(), id)
    if err != nil {
        h.app.errorResponse(w, r, http.StatusNotFound, "movie not found")
        return
    }
    
    response := toMovieResponse(movie)
    err = h.app.writeJSON(w, http.StatusOK, response, nil)
    if err != nil {
        h.app.logError(r, err)
        h.app.errorResponse(w, r, http.StatusInternalServerError, "failed to write response")
    }
}
```

---

## Performance Optimizations

### 1. Response Pooling

Reuse response objects for high-traffic endpoints:

```go
var movieResponsePool = sync.Pool{
    New: func() interface{} {
        return &MovieResponse{}
    },
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
    // Get response object from pool
    response := movieResponsePool.Get().(*MovieResponse)
    defer movieResponsePool.Put(response)
    
    // Reset fields
    *response = MovieResponse{}
    
    // Populate response
    movie, err := h.useCase.GetMovie(r.Context(), chi.URLParam(r, "id"))
    if err != nil {
        h.app.errorResponse(w, r, http.StatusNotFound, "movie not found")
        return
    }
    
    populateMovieResponse(response, movie)
    
    err = h.app.writeJSON(w, http.StatusOK, response, nil)
    if err != nil {
        h.app.logError(r, err)
        h.app.errorResponse(w, r, http.StatusInternalServerError, "failed to write response")
    }
}
```

### 2. JSON Buffer Pooling

Pool JSON marshal buffers for reduced allocations:

```go
var jsonBufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 1024) // 1KB initial capacity
    },
}

func (app *Application) writeJSONPooled(w http.ResponseWriter, status int, data any, headers http.Header) error {
    buf := jsonBufferPool.Get().([]byte)
    defer jsonBufferPool.Put(buf[:0]) // Reset length, keep capacity
    
    // Use buffer for marshaling
    buf, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    buf = append(buf, '\n')
    
    for key, value := range headers {
        w.Header()[key] = value
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    _, err = w.Write(buf)
    return err
}
```

### 3. Conditional Responses

Implement ETag/Last-Modified for cacheable resources:

```go
func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    movie, err := h.useCase.GetMovie(r.Context(), id)
    if err != nil {
        h.app.errorResponse(w, r, http.StatusNotFound, "movie not found")
        return
    }
    
    // Generate ETag based on movie data
    etag := generateETag(movie)
    w.Header().Set("ETag", etag)
    w.Header().Set("Last-Modified", movie.UpdatedAt.Format(http.TimeFormat))
    
    // Check if client has current version
    if match := r.Header.Get("If-None-Match"); match == etag {
        w.WriteHeader(http.StatusNotModified)
        return
    }
    
    response := toMovieResponse(movie)
    err = h.app.writeJSON(w, http.StatusOK, response, nil)
    if err != nil {
        h.app.logError(r, err)
        h.app.errorResponse(w, r, http.StatusInternalServerError, "failed to write response")
    }
}

func generateETag(movie domain.Movie) string {
    h := sha256.New()
    h.Write([]byte(fmt.Sprintf("%d-%d", movie.ID, movie.UpdatedAt.Unix())))
    return fmt.Sprintf(`"%x"`, h.Sum(nil)[:8])
}
```

---

## Error Handling Patterns

### 1. Standardized Error Responses

```go
type ErrorResponse struct {
    Error struct {
        Code    string `json:"code"`
        Message string `json:"message"`
        Details any    `json:"details,omitempty"`
    } `json:"error"`
    RequestID string    `json:"request_id"`
    Timestamp time.Time `json:"timestamp"`
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
    response := ErrorResponse{
        RequestID: getRequestID(r),
        Timestamp: time.Now(),
    }
    
    response.Error.Code = http.StatusText(status)
    response.Error.Message = message
    
    err := app.writeJSON(w, status, response, nil)
    if err != nil {
        app.logError(r, err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

// Validation error with details
func (app *Application) validationErrorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
    response := ErrorResponse{
        RequestID: getRequestID(r),
        Timestamp: time.Now(),
    }
    
    response.Error.Code = "VALIDATION_FAILED"
    response.Error.Message = "The request contains invalid data"
    response.Error.Details = errors
    
    err := app.writeJSON(w, http.StatusUnprocessableEntity, response, nil)
    if err != nil {
        app.logError(r, err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
```

### 2. Recovery Middleware

Handle panics gracefully:

```go
func (app *Application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Log the panic
                app.logger.Error("panic recovered", 
                    "error", err,
                    "request_id", getRequestID(r),
                    "stack", string(debug.Stack()))
                
                // Return JSON error response
                app.errorResponse(w, r, http.StatusInternalServerError, "internal server error")
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}
```

---

## Testing JSON Endpoints

### 1. JSON Response Testing

```go
func TestMovieHandler_GetMovie(t *testing.T) {
    app := newTestApplication(t)
    
    tests := []struct {
        name           string
        movieID        string
        setupMock      func(*MockMovieUseCase)
        wantStatus     int
        wantResponse   MovieResponse
        wantError      bool
    }{
        {
            name:    "successful retrieval",
            movieID: "123",
            setupMock: func(m *MockMovieUseCase) {
                m.On("GetMovie", mock.Anything, "123").Return(domain.Movie{
                    ID:    domain.MovieID(123),
                    Title: "Test Movie",
                    Year:  2023,
                }, nil)
            },
            wantStatus: http.StatusOK,
            wantResponse: MovieResponse{
                ID:    "123",
                Title: "Test Movie",
                Year:  2023,
            },
        },
        {
            name:    "movie not found",
            movieID: "999",
            setupMock: func(m *MockMovieUseCase) {
                m.On("GetMovie", mock.Anything, "999").Return(domain.Movie{}, domain.ErrMovieNotFound)
            },
            wantStatus: http.StatusNotFound,
            wantError:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockUseCase := &MockMovieUseCase{}
            tt.setupMock(mockUseCase)
            
            handler := &MovieHandler{
                useCase: mockUseCase,
                app:     app,
            }
            
            // Create request
            req := httptest.NewRequest(http.MethodGet, "/movies/"+tt.movieID, nil)
            req = req.WithContext(chi.NewRouteContext(req.Context(), chi.NewRouter()))
            chi.RouteContext(req.Context()).URLParams.Add("id", tt.movieID)
            
            rr := httptest.NewRecorder()
            
            // Execute
            handler.GetMovie(rr, req)
            
            // Assert status
            assert.Equal(t, tt.wantStatus, rr.Code)
            
            // Assert content type
            if rr.Code < 400 || tt.wantError {
                assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
            }
            
            // Assert response body
            if !tt.wantError {
                var response MovieResponse
                err := json.Unmarshal(rr.Body.Bytes(), &response)
                require.NoError(t, err)
                assert.Equal(t, tt.wantResponse.ID, response.ID)
                assert.Equal(t, tt.wantResponse.Title, response.Title)
                assert.Equal(t, tt.wantResponse.Year, response.Year)
            } else {
                var errorResponse ErrorResponse
                err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
                require.NoError(t, err)
                assert.NotEmpty(t, errorResponse.Error.Message)
            }
            
            mockUseCase.AssertExpectations(t)
        })
    }
}
```

### 2. JSON Validation Testing

```go
func TestJSONValidation(t *testing.T) {
    tests := []struct {
        name       string
        input      string
        wantValid  bool
        wantErrors map[string]string
    }{
        {
            name: "valid movie data",
            input: `{
                "title": "Test Movie",
                "year": 2023,
                "duration": 120
            }`,
            wantValid: true,
        },
        {
            name: "missing required fields",
            input: `{
                "year": 2023
            }`,
            wantValid: false,
            wantErrors: map[string]string{
                "title": "This field is required",
            },
        },
        {
            name: "invalid year",
            input: `{
                "title": "Test Movie",
                "year": 1800,
                "duration": 120
            }`,
            wantValid: false,
            wantErrors: map[string]string{
                "year": "Year must be between 1888 and current year",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var movie CreateMovieRequest
            err := json.Unmarshal([]byte(tt.input), &movie)
            require.NoError(t, err)
            
            errors := validateCreateMovieRequest(movie)
            
            if tt.wantValid {
                assert.Empty(t, errors)
            } else {
                for field, expectedMsg := range tt.wantErrors {
                    assert.Contains(t, errors, field)
                    assert.Equal(t, expectedMsg, errors[field])
                }
            }
        })
    }
}
```

---

## Production Considerations

### 1. Content Compression

Enable gzip compression for JSON responses:

```go
func (app *Application) enableGzip(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }
        
        w.Header().Set("Content-Encoding", "gzip")
        
        gz := gzip.NewWriter(w)
        defer gz.Close()
        
        gzw := &gzipResponseWriter{
            ResponseWriter: w,
            Writer:         gz,
        }
        
        next.ServeHTTP(gzw, r)
    })
}

type gzipResponseWriter struct {
    http.ResponseWriter
    io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}
```

### 2. Rate Limiting

Protect endpoints from abuse:

```go
func (app *Application) rateLimit(next http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Every(time.Second), 100) // 100 requests per second
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            app.errorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded")
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### 3. Request Tracing

Add correlation IDs for request tracing:

```go
func (app *Application) requestTracing(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // Add to context
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        r = r.WithContext(ctx)
        
        // Add to response headers
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r)
    })
}

func getRequestID(r *http.Request) string {
    if id, ok := r.Context().Value("request_id").(string); ok {
        return id
    }
    return "unknown"
}
```

### 4. Security Headers

Add security headers to JSON responses:

```go
func (app *Application) securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        next.ServeHTTP(w, r)
    })
}
```

## Best Practices Summary

1. **Use struct tags effectively** - Control JSON field names, omit empty values, hide sensitive data
2. **Centralize JSON handling** - Create reusable helpers for consistent responses
3. **Separate concerns** - Keep domain models pure, handle JSON in infrastructure layer
4. **Handle errors gracefully** - Provide structured error responses with correlation IDs
5. **Optimize for performance** - Use object pooling, compression, and conditional responses
6. **Test thoroughly** - Test both success and error cases, validate JSON structure
7. **Monitor in production** - Add tracing, rate limiting, and security headers

## Conclusion

Go's `encoding/json` package provides powerful tools for building robust JSON APIs.
By following Clean Architecture principles and implementing proper error handling,
testing, and performance optimizations, you can create APIs that are both maintainable
and performant.

The key is to start simple with basic patterns and gradually add sophistication
as your requirements grow. Always prioritize clarity and maintainability over
premature optimization.

---

**Related Resources:**

- [Go JSON Package Documentation](https://pkg.go.dev/encoding/json)
- [Go HTTP Server Best Practices](https://go.dev/doc/articles/wiki/)
- [Clean Architecture in Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
