# CORS (Cross-Origin Resource Sharing) in Go: Authoritative Guide

## Key Takeaways

### 1. CORS is a Browser Security Mechanism, Not a Server Security Feature

- CORS only protects browsers from malicious websites; server-to-server requests bypass CORS entirely.
- CORS relaxes the Same-Origin Policy by allowing servers to specify which origins can access their resources.
- Never rely on CORS alone for API security—always implement proper authentication and authorization.

### 2. Understanding Simple vs. Preflighted Requests

- **Simple requests** (GET, HEAD, POST with basic content types) are sent directly with CORS headers.
- **Preflighted requests** (PUT, DELETE, custom headers, non-simple content types) trigger an OPTIONS request first.
- The preflight OPTIONS request checks permissions before the actual request is made.
- Preflight responses can be cached by browsers using `Access-Control-Max-Age` header.

### 3. CORS Headers Hierarchy and Security

**Request Headers (Browser → Server):**
- `Origin`: Specifies the requesting domain (set by browser, cannot be spoofed in web contexts)

**Response Headers (Server → Browser):**
- `Access-Control-Allow-Origin`: Specifies allowed origins (use specific domains, not `*` with credentials)
- `Access-Control-Allow-Methods`: Lists allowed HTTP methods for preflight requests
- `Access-Control-Allow-Headers`: Lists allowed request headers for preflight requests
- `Access-Control-Allow-Credentials`: Enables cookies/auth headers (requires specific origin, not `*`)
- `Access-Control-Expose-Headers`: Lists response headers accessible to JavaScript
- `Access-Control-Max-Age`: Caches preflight response (in seconds)

### 4. Go Implementation Patterns

#### Echo Framework (labstack/echo)

**Basic CORS Setup:**
```go
import "github.com/labstack/echo/v4/middleware"

// Enable default CORS (allows all origins)
e.Use(middleware.CORS())
```

**Production-Ready Configuration:**
```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{
        "https://myapp.com", 
        "https://www.myapp.com",
        "https://admin.myapp.com",
    },
    AllowMethods: []string{
        echo.GET, echo.HEAD, echo.PUT, 
        echo.PATCH, echo.POST, echo.DELETE,
    },
    AllowHeaders: []string{
        echo.HeaderOrigin,
        echo.HeaderContentType,
        echo.HeaderAccept,
        echo.HeaderAuthorization,
        "X-Requested-With",
        "X-API-Key",
    },
    AllowCredentials: true,
    ExposeHeaders: []string{
        "X-Total-Count",
        "X-Page-Number", 
        "X-Rate-Limit-Remaining",
    },
    MaxAge: 300, // 5 minutes preflight cache
}))
```

**Dynamic Origin Validation:**
```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOriginFunc: func(origin string) (bool, error) {
        // Custom logic for dynamic origin validation
        if strings.HasSuffix(origin, ".mycompany.com") {
            return true, nil
        }
        if origin == "https://trusted-partner.com" {
            return true, nil
        }
        // Log suspicious origins for security monitoring
        log.Printf("CORS: Blocked origin %s", origin)
        return false, nil
    },
    AllowCredentials: true,
    AllowHeaders: []string{
        echo.HeaderOrigin, echo.HeaderContentType, 
        echo.HeaderAccept, echo.HeaderAuthorization,
    },
}))
```

#### Gin Framework (gin-gonic/gin)

```go
import "github.com/gin-contrib/cors"

// Production configuration
config := cors.Config{
    AllowOrigins:     []string{"https://myapp.com", "https://admin.myapp.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
    AllowCredentials: true,
    MaxAge:          12 * time.Hour,
}
router.Use(cors.New(config))
```

#### Standard net/http with Gorilla

```go
import "github.com/gorilla/handlers"

// Basic CORS with Gorilla handlers
handler := handlers.CORS(
    handlers.AllowedOrigins([]string{"https://myapp.com"}),
    handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
    handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    handlers.AllowCredentials(),
)(router)

http.ListenAndServe(":8080", handler)
```

#### Chi Router

```go
import "github.com/go-chi/cors"

r.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"https://myapp.com", "https://admin.myapp.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300,
}))
```

### 5. Security Best Practices

**Never Use Wildcard with Credentials:**
```go
// ❌ DANGER: This is a security vulnerability
AllowOrigins: []string{"*"},
AllowCredentials: true,  // This combination is forbidden by browsers

// ✅ CORRECT: Specify exact origins when using credentials
AllowOrigins: []string{"https://trusted-app.com"},
AllowCredentials: true,
```

**Validate Origins Dynamically for Multi-Tenant Apps:**
```go
AllowOriginFunc: func(origin string) (bool, error) {
    // Check against database of allowed tenant domains
    allowed, err := tenantService.IsOriginAllowed(origin)
    if err != nil {
        log.Printf("CORS origin validation error: %v", err)
        return false, nil // Fail closed on errors
    }
    return allowed, nil
}
```

**Environment-Specific Configuration:**
```go
func getCORSConfig() middleware.CORSConfig {
    if os.Getenv("ENV") == "development" {
        return middleware.CORSConfig{
            AllowOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
            AllowCredentials: true,
            AllowHeaders: []string{"*"}, // More permissive for development
        }
    }
    
    return middleware.CORSConfig{
        AllowOrigins: []string{
            "https://app.production.com",
            "https://admin.production.com",
        },
        AllowCredentials: true,
        AllowHeaders: []string{
            "Origin", "Content-Type", "Accept", "Authorization", "X-API-Key",
        },
        MaxAge: 86400, // 24 hours for production
    }
}
```

### 6. Common CORS Pitfalls and Solutions

**Problem: Credentials Not Being Sent**
```go
// ✅ Server-side: Enable credentials
AllowCredentials: true,

// ✅ Client-side: Include credentials in fetch
fetch('https://api.example.com/data', {
    credentials: 'include',  // Send cookies/auth headers
    headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + token
    }
})
```

**Problem: Custom Headers Causing Preflight Failures**
```go
// ✅ Ensure custom headers are explicitly allowed
AllowHeaders: []string{
    "Origin", "Content-Type", "Accept", "Authorization",
    "X-Requested-With", "X-API-Key", "X-Client-Version", // Custom headers
}
```

**Problem: Preflight Cache Issues During Development**
```go
// ✅ Short cache time for development, longer for production
MaxAge: func() int {
    if os.Getenv("ENV") == "development" {
        return 10 // 10 seconds
    }
    return 86400 // 24 hours for production
}(),
```

### 7. Testing CORS Configuration

**Unit Testing with Echo:**
```go
func TestCORSConfiguration(t *testing.T) {
    e := echo.New()
    e.Use(middleware.CORSWithConfig(getCORSConfig()))
    
    // Test allowed origin
    req := httptest.NewRequest(http.MethodOptions, "/", nil)
    req.Header.Set("Origin", "https://myapp.com")
    req.Header.Set("Access-Control-Request-Method", "POST")
    rec := httptest.NewRecorder()
    
    e.ServeHTTP(rec, req)
    
    assert.Equal(t, "https://myapp.com", rec.Header().Get("Access-Control-Allow-Origin"))
    assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}
```

**Integration Testing:**
```bash
# Test preflight request
curl -X OPTIONS \
  -H "Origin: https://myapp.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,Authorization" \
  https://api.example.com/endpoint

# Test actual request
curl -X POST \
  -H "Origin: https://myapp.com" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token" \
  https://api.example.com/endpoint
```

### 8. Monitoring and Debugging

**CORS Error Logging:**
```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOriginFunc: func(origin string) (bool, error) {
        allowed := isOriginAllowed(origin)
        
        // Log CORS decisions for monitoring
        log.WithFields(log.Fields{
            "origin":  origin,
            "allowed": allowed,
            "action":  "cors_check",
        }).Info("CORS origin validation")
        
        if !allowed {
            // Alert on suspicious patterns
            metrics.Counter("cors.blocked_origins").Inc()
        }
        
        return allowed, nil
    },
}))
```

**Browser Developer Tools:**
- Check Network tab for preflight OPTIONS requests
- Look for CORS errors in Console tab
- Verify response headers match your configuration

### 9. Performance Considerations

**Preflight Cache Optimization:**
```go
// Longer cache times reduce preflight overhead
MaxAge: 86400, // 24 hours

// Group similar endpoints to share preflight cache
// Instead of: /api/users/123, /api/users/456
// Use: /api/users with route parameters
```

**Conditional CORS Middleware:**
```go
// Skip CORS for same-origin requests
Skipper: func(c echo.Context) bool {
    origin := c.Request().Header.Get("Origin")
    return origin == "" || origin == c.Scheme()+"://"+c.Request().Host
},
```

### 10. Advanced CORS Patterns

**Multi-Environment Origin Management:**
```go
type CORSManager struct {
    allowedOrigins map[string][]string
}

func (cm *CORSManager) GetConfig(env string) middleware.CORSConfig {
    origins := cm.allowedOrigins[env]
    if origins == nil {
        origins = []string{} // Empty for unknown environments
    }
    
    return middleware.CORSConfig{
        AllowOrigins:     origins,
        AllowCredentials: true,
        MaxAge:          getCacheTime(env),
    }
}
```

**Tenant-Specific CORS:**
```go
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        tenant := extractTenant(c)
        corsConfig := getTenantCORSConfig(tenant)
        
        return middleware.CORSWithConfig(corsConfig)(next)(c)
    }
})
```

---

**Summary for Future Projects:**

- Always specify exact origins in production; avoid wildcards with credentials
- Use framework-specific CORS middleware (Echo, Gin, Chi) rather than manual implementation
- Implement proper origin validation for multi-tenant or dynamic environments
- Test both simple and preflighted requests in your test suite
- Monitor CORS blocks and errors for security insights
- Configure appropriate preflight cache times for your use case
- Remember: CORS protects browsers, not APIs—implement proper authentication/authorization

## References

- [Authoritative Guide to CORS - Moesif](https://www.moesif.com/blog/technical/cors/Authoritative-Guide-to-CORS-Cross-Origin-Resource-Sharing-for-REST-APIs/)
- [Echo CORS Middleware Documentation](https://echo.labstack.com/docs/middleware/cors)
- [Gin CORS Middleware](https://github.com/gin-contrib/cors)
- [MDN CORS Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [OWASP CORS Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Origin_Resource_Sharing_Cheat_Sheet.html)
