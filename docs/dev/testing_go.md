# Go Testing Best Practices for Frank

## Core Testing Philosophy

- **Test behavior, not implementation**: Focus on _what_ the code does, not _how_ it does it
- **Aim for high coverage, but prioritize critical paths**
- **Integration tests complement unit tests** but should be kept separate
- **Test edge cases and error conditions thoroughly**
- **Keep tests readable and maintainable**: Use descriptive names and clear structure
- **Prefer table-driven tests** for clarity and coverage

<!-- REF: https://go.dev/doc/tutorial/add-a-test -->
<!-- REF: https://github.com/golang/go/wiki/TestComments -->
<!-- REF: https://threedots.tech/post/table-driven-tests-in-go/ -->

## Project-Specific Testing Structure

```bash
# Place unit tests alongside the code, named *_test.go
# Place integration tests in a separate /integration or /test directory
```

## Test Organization

### Unit Test Structure

```go
func TestFunctionName_Scenario_ExpectedBehavior(t *testing.T) {
    // ARRANGE: Set up test data and expectations
    input := "sample"
    expected := "result"

    // ACT: Call the function being tested
    actual := FunctionBeingTested(input)

    // ASSERT: Verify the results
    if actual != expected {
        t.Errorf("Expected %v but got %v", expected, actual)
    }
}
```

### Table-Driven Tests

Preferred for functions with multiple input/output scenarios:

```go
func TestIsValidEmail(t *testing.T) {
    tests := []struct {
        name  string
        email string
        want  bool
    }{
        {"Valid email", "user@example.com", true},
        {"Missing @", "userexample.com", false},
        {"Missing domain", "user@", false},
        {"Invalid TLD", "user@example.123", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := pst.IsValidEmail(tt.email)
            if got != tt.want {
                t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
            }
        })
    }
}
```

> **See also:** [Table-driven tests in Go](https://threedots.tech/post/table-driven-tests-in-go/)

### Subtests and Parallelism

Use subtests (`t.Run`) for logical grouping and `t.Parallel()` for concurrency:

```go
func TestSomething(t *testing.T) {
    cases := []struct {
        name string
        input int
        want int
    }{
        {"zero", 0, 0},
        {"positive", 1, 1},
    }

    for _, tc := range cases {
        tc := tc // capture range variable
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            got := SomeFunc(tc.input)
            if got != tc.want {
                t.Errorf("got %v, want %v", got, tc.want)
            }
        })
    }
}
```

<!-- REF: https://threedots.tech/post/advanced-testing-in-go-subtests-and-test-main/ -->

## Mocking and Dependency Injection

Use interfaces for testability:

```go
// Make dependencies explicit through interfaces
type PSTextractor interface {
    ExtractContacts(filePath string, extractAttachments bool, debug bool) ([]Contact, error)
}

// Mock implementation for testing
type MockPSTextractor struct {
    Contacts []Contact
    Err      error
}

func (m *MockPSTextractor) ExtractContacts(filePath string, extractAttachments bool, debug bool) ([]Contact, error) {
    return m.Contacts, m.Err
}

// Test using the mock
func TestProcessFile_WithValidFile(t *testing.T) {
    mockExtractor := &MockPSTextractor{
        Contacts: []Contact{
            {Email: "test@example.com", DisplayName: "Test User"},
        },
        Err: nil,
    }

    processor := NewProcessor(config, mockExtractor)
    // Test processor with the mock
}
```

> **See also:** [Dependency injection in Go](https://threedots.tech/post/dependency-injection-in-go/)

## Testing Concurrency

1. **Deterministic Tests**: Avoid timing-dependent tests

    ```go
        // AVOID
        func TestConcurrent_Flaky(t *testing.T) {
            go doSomething()
            time.Sleep(100 * time.Millisecond) // Unreliable!
            // Assert something
        }
    ```

    ```go
        func TestConcurrent_Reliable(t *testing.T) {
            done := make(chan struct{})
            go func() {
                doSomething()
                close(done)
            }()

            select {
            case <-done:
                // Success case
            case <-time.After(1 * time.Second):
                t.Fatal("Test timed out")
            }
        }
    ```

2. **Testing Worker Pools**: Verify behavior, not exact timing

    ```go
        func TestWorkerPool(t *testing.T) {
            // Arrange: Create sample jobs
            jobs := []string{"job1", "job2", "job3"}
            results := make(chan string, len(jobs))

            // Act: Run the worker pool
            processJobs(jobs, results)

            // Assert: Verify all jobs were processed
            var processed []string
            for i := 0; i < len(jobs); i++ {
                processed = append(processed, <-results)
            }

            // Use a helper like testify/assert for more readable assertions
            if len(processed) != len(jobs) {
                t.Errorf("Expected %d processed jobs, got %d", len(jobs), len(processed))
            }
        }
    ```

> **See also:** [Testing concurrent code in Go](https://threedots.tech/post/testing-concurrent-code-in-go/)

## Testing File Operations

Use `io/fs` testing utilities and temporary directories:

## Testing Resource Management

Focus on behavior change rather than absolute values:

## Recommended Testing Libraries

- **Standard library**: `testing` for most needs
- **Assertion helpers**: `github.com/stretchr/testify/assert` for readable assertions
- **HTTP testing**: `net/http/httptest` for API testing
- **Mock generation**: `github.com/golang/mock/gomock` for interface mocking
- **Golden files**: Use for complex output comparison ([Golden files in Go](https://threedots.tech/post/golden-files-in-go-tests/))

<!-- REF: https://pkg.go.dev/github.com/stretchr/testify/assert -->
<!-- REF: https://pkg.go.dev/net/http/httptest -->
<!-- REF: https://pkg.go.dev/github.com/golang/mock/gomock -->
<!-- REF: https://threedots.tech/post/golden-files-in-go-tests/ -->

## Common Testing Anti-patterns to Avoid

1. **Testing implementation details** rather than behavior
2. **Brittle tests** that break when internal implementation changes
3. **Slow tests** that access real resources unnecessarily
4. **Global state** that makes tests non-deterministic
5. **Partial assertions** that don't verify complete function behavior
6. **Exposing fields or methods** solely for testing purposes
7. **Ignoring error values**: Always check errors in tests
8. **Not cleaning up resources**: Use `t.Cleanup` or defer for cleanup

<!-- REF: https://dave.cheney.net/2016/08/20/solid-go-design -->
<!-- REF: https://threedots.tech/post/common-mistakes-in-go-tests/ -->

## Test Coverage

Use the built-in Go coverage tool:

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Check coverage percentage
go tool cover -func=coverage.out
```

Aim for high coverage (>80%) of critical components:

- Contact extraction logic
- Error handling paths
- Resource management
- Progress tracking

## Further Reading

- [Table-driven tests in Go](https://threedots.tech/post/table-driven-tests-in-go/)
- [Advanced testing in Go: subtests and TestMain](https://threedots.tech/post/advanced-testing-in-go-subtests-and-test-main/)
- [Dependency injection in Go](https://threedots.tech/post/dependency-injection-in-go/)
- [Testing concurrent code in Go](https://threedots.tech/post/testing-concurrent-code-in-go/)
- [Golden files in Go tests](https://threedots.tech/post/golden-files-in-go-tests/)
- [Common mistakes in Go tests](https://threedots.tech/post/common-mistakes-in-go-tests/)
