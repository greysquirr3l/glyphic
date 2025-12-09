---
applyTo: '**/*_test.go'
description: Go testing standards and patterns for comprehensive test coverage
references:
  testing_guide: "docs/dev/testing_go.md"
  go_tutorial: "https://go.dev/doc/tutorial/add-a-test"
  table_driven: "https://threedots.tech/post/table-driven-tests-in-go/"
---

# Go Testing Standards

## Core Testing Philosophy

- **Test behavior, not implementation**: Focus on _what_ the code does, not _how_ it does it
- **Aim for high coverage, but prioritize critical paths**: Target 80%+ coverage for business logic
- **Integration tests complement unit tests** but should be kept separate
- **Test edge cases and error conditions thoroughly**
- **Keep tests readable and maintainable**: Use descriptive names and clear structure
- **Prefer table-driven tests** for clarity and comprehensive coverage

## Test Structure and Organization

### Project-Specific Testing Structure
```bash
# Place unit tests alongside the code, named *_test.go
internal/domain/user/user.go
internal/domain/user/user_test.go

# Place integration tests in a separate directory
test/integration/user_integration_test.go
test/e2e/api_test.go
```

### File Naming and Package Structure
```go
// ✅ Test files should end with _test.go
// File: user_service_test.go
package user_test // External testing package

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
    
    "your-project/internal/domain/user"
)

// ✅ Use descriptive test function names
func TestUserService_CreateUser_Success(t *testing.T) {
    // Test implementation
}

func TestUserService_CreateUser_UserAlreadyExists_ReturnsError(t *testing.T) {
    // Test implementation
}
```

## Table-Driven Test Patterns

### Comprehensive Test Cases
```go
// ✅ Table-driven tests for thorough coverage
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name        string
        email       string
        want        error
        description string
    }{
        {
            name:        "valid_email",
            email:       "user@example.com",
            want:        nil,
            description: "should accept valid email format",
        },
        {
            name:        "empty_email",
            email:       "",
            want:        ErrEmailRequired,
            description: "should reject empty email",
        },
        {
            name:        "invalid_format",
            email:       "notanemail",
            want:        ErrInvalidEmailFormat,
            description: "should reject invalid email format",
        },
        {
            name:        "missing_domain",
            email:       "user@",
            want:        ErrInvalidEmailFormat,
            description: "should reject email without domain",
        },
        {
            name:        "missing_at_symbol",
            email:       "userexample.com",
            want:        ErrInvalidEmailFormat,
            description: "should reject email without @ symbol",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange - test setup is in table
            
            // Act
            got := ValidateEmail(tt.email)
            
            // Assert
            if tt.want == nil {
                assert.NoError(t, got, tt.description)
            } else {
                assert.ErrorIs(t, got, tt.want, tt.description)
            }
        })
    }
}
```

### Complex Scenario Testing
```go
// ✅ Test complex business logic scenarios
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name           string
        request        CreateUserRequest
        setupMocks     func(*MockUserRepository)
        expectedResult *User
        expectedError  error
        description    string
    }{
        {
            name: "successful_user_creation",
            request: CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMocks: func(repo *MockUserRepository) {
                repo.On("GetUserByEmail", mock.Anything, "test@example.com").
                    Return(nil, ErrUserNotFound)
                repo.On("SaveUser", mock.Anything, mock.AnythingOfType("*User")).
                    Return(nil)
            },
            expectedResult: &User{
                Email: "test@example.com",
                Name:  "Test User",
            },
            expectedError: nil,
            description:   "should create user successfully when email is unique",
        },
        {
            name: "user_already_exists",
            request: CreateUserRequest{
                Email: "existing@example.com",
                Name:  "Existing User",
            },
            setupMocks: func(repo *MockUserRepository) {
                existingUser := &User{
                    ID:    "existing-id",
                    Email: "existing@example.com",
                    Name:  "Existing User",
                }
                repo.On("GetUserByEmail", mock.Anything, "existing@example.com").
                    Return(existingUser, nil)
            },
            expectedResult: nil,
            expectedError:  ErrUserAlreadyExists,
            description:    "should return error when user with email already exists",
        },
        {
            name: "database_error_on_check",
            request: CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMocks: func(repo *MockUserRepository) {
                repo.On("GetUserByEmail", mock.Anything, "test@example.com").
                    Return(nil, errors.New("database connection failed"))
            },
            expectedResult: nil,
            expectedError:  errors.New("database connection failed"),
            description:    "should handle database errors when checking existing user",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := new(MockUserRepository)
            tt.setupMocks(mockRepo)
            
            service := NewUserService(mockRepo, slog.Default())
            ctx := context.Background()
            
            // Act
            result, err := service.CreateUser(ctx, tt.request)
            
            // Assert
            if tt.expectedError != nil {
                require.Error(t, err, tt.description)
                assert.Nil(t, result)
                // For specific error types, use ErrorIs
                if errors.Is(tt.expectedError, ErrUserAlreadyExists) {
                    assert.ErrorIs(t, err, ErrUserAlreadyExists)
                }
            } else {
                require.NoError(t, err, tt.description)
                require.NotNil(t, result)
                assert.Equal(t, tt.expectedResult.Email, result.Email)
                assert.Equal(t, tt.expectedResult.Name, result.Name)
                assert.NotEmpty(t, result.ID)
                assert.WithinDuration(t, time.Now(), result.CreatedAt, time.Second)
            }
            
            // Verify all mock expectations were met
            mockRepo.AssertExpectations(t)
        })
    }
}
```

## Test Suites and Setup/Teardown

### Using Testify Suites for Complex Setup
```go
// ✅ Use test suites for complex scenarios requiring setup/teardown
type UserServiceTestSuite struct {
    suite.Suite
    mockRepo   *MockUserRepository
    service    *UserService
    ctx        context.Context
    testUsers  []*User
}

func (suite *UserServiceTestSuite) SetupTest() {
    suite.mockRepo = new(MockUserRepository)
    suite.service = NewUserService(suite.mockRepo, slog.Default())
    suite.ctx = context.Background()
    
    // Create test data
    suite.testUsers = []*User{
        {
            ID:        "user-1",
            Email:     "user1@example.com",
            Name:      "User One",
            CreatedAt: time.Now(),
        },
        {
            ID:        "user-2", 
            Email:     "user2@example.com",
            Name:      "User Two",
            CreatedAt: time.Now(),
        },
    }
}

func (suite *UserServiceTestSuite) TearDownTest() {
    // Clean up any resources if needed
    suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestCreateUser_Success() {
    // Arrange
    request := CreateUserRequest{
        Email: "new@example.com",
        Name:  "New User",
    }
    
    suite.mockRepo.On("GetUserByEmail", suite.ctx, request.Email).
        Return(nil, ErrUserNotFound)
    suite.mockRepo.On("SaveUser", suite.ctx, mock.AnythingOfType("*User")).
        Return(nil)
    
    // Act
    user, err := suite.service.CreateUser(suite.ctx, request)
    
    // Assert
    suite.NoError(err)
    suite.Equal(request.Email, user.Email)
    suite.Equal(request.Name, user.Name)
    suite.NotEmpty(user.ID)
}

func TestUserServiceTestSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

## Mock Generation and Usage

### Generate Mocks with GoMock
```go
// ✅ Generate mocks for interfaces
//go:generate mockgen -source=user_repository.go -destination=mocks/user_repository_mock.go

// In test file:
func TestUserService_WithGoMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo, slog.Default())
    
    ctx := context.Background()
    userID := "test-user-id"
    expectedUser := &User{
        ID:    userID,
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    // Setup expectation
    mockRepo.EXPECT().
        GetUserByID(ctx, userID).
        Return(expectedUser, nil).
        Times(1)
    
    // Execute
    user, err := service.GetUserByID(ctx, userID)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, expectedUser, user)
}
```

## Integration Testing Patterns

### Database Integration Tests
```go
// ✅ Integration tests with real database
func TestUserRepository_Integration(t *testing.T) {
    // Skip integration tests in short mode
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewUserRepository(db)
    ctx := context.Background()
    
    t.Run("create_and_retrieve_user", func(t *testing.T) {
        // Arrange
        user := &User{
            ID:        generateTestID(),
            Email:     "integration@test.com",
            Name:      "Integration Test User",
            CreatedAt: time.Now(),
        }
        
        // Act - Save user
        err := repo.SaveUser(ctx, user)
        require.NoError(t, err)
        
        // Act - Retrieve user
        retrieved, err := repo.GetUserByID(ctx, user.ID)
        require.NoError(t, err)
        
        // Assert
        assert.Equal(t, user.ID, retrieved.ID)
        assert.Equal(t, user.Email, retrieved.Email)
        assert.Equal(t, user.Name, retrieved.Name)
        assert.WithinDuration(t, user.CreatedAt, retrieved.CreatedAt, time.Second)
    })
    
    t.Run("user_not_found", func(t *testing.T) {
        // Act
        user, err := repo.GetUserByID(ctx, "non-existent-id")
        
        // Assert
        assert.Nil(t, user)
        assert.ErrorIs(t, err, ErrUserNotFound)
    })
}

// Helper functions for integration tests
func setupTestDB(t *testing.T) *sql.DB {
    // Setup test database connection
    // Run migrations
    // Return database connection
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
    // Clean up test data
    // Close connections
}
```

## Benchmark Testing

### Performance Testing Patterns
```go
// ✅ Benchmark critical performance paths
func BenchmarkValidateEmail(b *testing.B) {
    testEmail := "user@example.com"
    
    b.ResetTimer() // Reset timer after setup
    
    for i := 0; i < b.N; i++ {
        _ = ValidateEmail(testEmail)
    }
}

func BenchmarkProcessUsers_ParallelProcessing(b *testing.B) {
    users := generateTestUsers(1000) // Generate test data
    service := setupTestService()
    
    b.ResetTimer()
    b.SetBytes(int64(len(users))) // Set bytes processed per iteration
    
    for i := 0; i < b.N; i++ {
        err := service.ProcessUsers(context.Background(), users)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// ✅ Benchmark with different input sizes
func BenchmarkUserService_CreateUser(b *testing.B) {
    sizes := []int{1, 10, 100, 1000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("users_%d", size), func(b *testing.B) {
            service := setupBenchmarkService()
            requests := generateCreateUserRequests(size)
            
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                for _, req := range requests {
                    _, err := service.CreateUser(context.Background(), req)
                    if err != nil {
                        b.Fatal(err)
                    }
                }
            }
        })
    }
}
```

## Concurrency Testing Patterns

### Testing Concurrent Code Safely

```go
// ✅ Use channels and timeouts for deterministic concurrent tests
func TestConcurrentProcessor_Success(t *testing.T) {
    done := make(chan struct{})
    results := make(chan result, 10)
    
    // Start concurrent processing
    go func() {
        processor.ProcessConcurrently(context.Background(), jobs, results)
        close(done)
    }()
    
    // Wait for completion with timeout
    select {
    case <-done:
        // Success - verify results
        assert.Len(t, results, len(expectedResults))
    case <-time.After(5 * time.Second):
        t.Fatal("Test timed out - concurrent processing took too long")
    }
}

// ✅ Test worker pool behavior
func TestWorkerPool_ProcessesAllJobs(t *testing.T) {
    const numJobs = 100
    const numWorkers = 5
    
    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    
    // Add jobs
    for i := 0; i < numJobs; i++ {
        jobs <- Job{ID: i, Data: fmt.Sprintf("job-%d", i)}
    }
    close(jobs)
    
    // Start worker pool
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            processJobs(jobs, results)
        }()
    }
    
    // Close results when all workers finish
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect and verify results
    var processed []Result
    for result := range results {
        processed = append(processed, result)
    }
    
    assert.Len(t, processed, numJobs)
    // Additional assertions on results...
}

// ❌ AVOID: Timing-dependent tests that can be flaky
func TestConcurrent_Flaky(t *testing.T) {
    go doSomething()
    time.Sleep(100 * time.Millisecond) // This can fail randomly!
    // Assert something...
}

// ✅ Use race detector for concurrent tests
// Run with: go test -race ./...
func TestConcurrentAccess(t *testing.T) {
    counter := NewSafeCounter()
    
    var wg sync.WaitGroup
    const numGoroutines = 10
    const incrementsPerGoroutine = 100
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < incrementsPerGoroutine; j++ {
                counter.Increment()
            }
        }()
    }
    
    wg.Wait()
    
    expected := numGoroutines * incrementsPerGoroutine
    assert.Equal(t, expected, counter.Value())
}
```

## Test Data Management

### Test Data Factories
```go
// ✅ Create test data factories for consistent test data
type UserTestFactory struct {
    counter int
}

func NewUserTestFactory() *UserTestFactory {
    return &UserTestFactory{}
}

func (f *UserTestFactory) CreateUser(overrides ...func(*User)) *User {
    f.counter++
    
    user := &User{
        ID:        fmt.Sprintf("test-user-%d", f.counter),
        Email:     fmt.Sprintf("user%d@example.com", f.counter),
        Name:      fmt.Sprintf("Test User %d", f.counter),
        CreatedAt: time.Now(),
    }
    
    // Apply overrides
    for _, override := range overrides {
        override(user)
    }
    
    return user
}

func (f *UserTestFactory) CreateUserWithEmail(email string) *User {
    return f.CreateUser(func(u *User) {
        u.Email = email
    })
}

// Usage in tests
func TestSomething(t *testing.T) {
    factory := NewUserTestFactory()
    
    user1 := factory.CreateUser()
    user2 := factory.CreateUserWithEmail("specific@example.com")
    user3 := factory.CreateUser(func(u *User) {
        u.Name = "Custom Name"
        u.CreatedAt = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    })
}
```

## Error Testing Patterns

### Comprehensive Error Scenarios
```go
// ✅ Test all error conditions thoroughly
func TestUserService_CreateUser_ErrorScenarios(t *testing.T) {
    tests := []struct {
        name          string
        request       CreateUserRequest
        setupMocks    func(*MockUserRepository)
        expectedError error
        errorCheck    func(t *testing.T, err error)
    }{
        {
            name: "validation_error",
            request: CreateUserRequest{
                Email: "", // Invalid empty email
                Name:  "Valid Name",
            },
            setupMocks: func(repo *MockUserRepository) {
                // No mock setup needed - validation happens first
            },
            expectedError: ErrValidation,
            errorCheck: func(t *testing.T, err error) {
                assert.ErrorIs(t, err, ErrValidation)
                assert.Contains(t, err.Error(), "email")
            },
        },
        {
            name: "database_connection_error",
            request: CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMocks: func(repo *MockUserRepository) {
                repo.On("GetUserByEmail", mock.Anything, "test@example.com").
                    Return(nil, errors.New("database connection lost"))
            },
            expectedError: nil, // We'll check the wrapped error
            errorCheck: func(t *testing.T, err error) {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), "failed to check existing user")
                assert.Contains(t, err.Error(), "database connection lost")
            },
        },
        {
            name: "timeout_error",
            request: CreateUserRequest{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMocks: func(repo *MockUserRepository) {
                repo.On("GetUserByEmail", mock.Anything, "test@example.com").
                    Return(nil, context.DeadlineExceeded)
            },
            expectedError: context.DeadlineExceeded,
            errorCheck: func(t *testing.T, err error) {
                assert.ErrorIs(t, err, context.DeadlineExceeded)
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := new(MockUserRepository)
            tt.setupMocks(mockRepo)
            
            service := NewUserService(mockRepo, slog.Default())
            ctx := context.Background()
            
            // Act
            user, err := service.CreateUser(ctx, tt.request)
            
            // Assert
            assert.Nil(t, user)
            require.Error(t, err)
            
            if tt.expectedError != nil {
                assert.ErrorIs(t, err, tt.expectedError)
            }
            
            if tt.errorCheck != nil {
                tt.errorCheck(t, err)
            }
            
            mockRepo.AssertExpectations(t)
        })
    }
}
```

## Test Coverage Requirements

### Coverage Guidelines
- **Minimum 80% coverage** for business logic packages
- **Minimum 60% coverage** for infrastructure packages  
- **100% coverage** for critical security functions
- **All exported functions** must have at least one test
- **All error paths** must be tested
- **Edge cases and boundary conditions** must be covered

### Running Tests with Coverage
```bash
# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Check coverage percentage
go tool cover -func=coverage.out | tail -1
```

Always strive for comprehensive test coverage that validates both the happy path and all possible error conditions to ensure robust, reliable code.