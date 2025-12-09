# Codecov Integration Guide

This document provides comprehensive guidance for the Codecov integration in the <project_name> project.

## 🎯 Overview

Codecov is integrated into our CI/CD pipeline to provide detailed code coverage analysis and reporting. The setup includes:

- **Automated coverage collection** during CI runs
- **Pull request coverage reporting** with diff analysis
- **Coverage trend tracking** over time
- **Quality gates** to maintain coverage standards

## 🔧 Configuration

### GitHub Repository Setup

1. **Codecov Token**: Add `CODECOV_TOKEN` to GitHub repository secrets
   - Go to Codecov.io and sign in with GitHub
   - Find your repository and copy the upload token
   - Add as `CODECOV_TOKEN` in GitHub repository secrets

2. **Repository Settings**:
   - Enable "Require status checks to pass before merging"
   - Add `codecov/project` and `codecov/patch` as required status checks

### CI/CD Integration

Our CI pipeline includes multiple coverage collection points:

#### Test Job Coverage

```yaml
- name: Run tests with race detection
  run: go test -v -race -timeout=30m -coverprofile=coverage.out -covermode=atomic ./...

- name: Upload coverage to Codecov
  if: matrix.os == 'ubuntu-latest'
  uses: codecov/codecov-action@v4
  with:
    file: ./coverage.out
    flags: unittests
    name: codecov-umbrella
    token: ${{ secrets.CODECOV_TOKEN }}
    fail_ci_if_error: false
    verbose: true
```

#### Dedicated Coverage Job

```yaml
- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v4
  with:
    file: ./coverage.out
    flags: coverage-check
    name: coverage-analysis
    token: ${{ secrets.CODECOV_TOKEN }}
    fail_ci_if_error: true
    verbose: true
```

## 📊 Coverage Standards

### Quality Gates

- **Project Coverage**: Minimum 80% overall coverage
- **Patch Coverage**: Minimum 85% for new code
- **Individual Files**: Minimum 70% per file
- **Critical Files**: Higher coverage requirements for core files

### Coverage Thresholds

```yaml
coverage:
  status:
    project:
      default:
        target: 80%
        threshold: 1%
    patch:
      default:
        target: 85%
        threshold: 5%
```

## 🎛️ Codecov Configuration

### `.codecov.yml` Features

Our configuration includes:

1. **Precision Settings**: 2 decimal places for accuracy
2. **Status Checks**: Project and patch coverage validation
3. **Ignore Patterns**: Exclude test files and generated code
4. **PR Comments**: Automated coverage reporting
5. **Critical Files**: Special handling for core components

### Ignored Files

```yaml
ignore:
  - "**/*_test.go"        # Test files
  - "**/test/**"          # Test directories
  - "**/example/**"       # Example code
  - "**/debug/**"         # Debug utilities
  - "**/*.pb.go"          # Generated files
```

## 📈 Coverage Reports

### Dashboard Features

- **Coverage Trends**: Historical coverage data
- **File Browser**: Line-by-line coverage analysis
- **Pull Request**: Diff coverage analysis
- **Commit**: Individual commit coverage
- **Sunburst**: Visual coverage representation

### Report Types

1. **Project Coverage**: Overall repository coverage
2. **Patch Coverage**: Coverage for changed lines in PRs
3. **Flag Coverage**: Coverage by test type or component
4. **File Coverage**: Individual file coverage details

## 🔍 Using Coverage Data

### Pull Request Workflow

1. **Automated Comments**: Coverage summary posted on PRs
2. **Status Checks**: Pass/fail based on coverage thresholds
3. **Diff Analysis**: Shows coverage for changed lines
4. **Trend Comparison**: Compares against base branch

### Local Development

```bash
# Generate coverage locally
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# View coverage report
go tool cover -html=coverage.out -o coverage.html
open coverage.html

# Check coverage percentage
go tool cover -func=coverage.out | grep total
```

## 🚀 Best Practices

### Writing Testable Code

1. **Small Functions**: Easier to test and cover
2. **Clear Interfaces**: Mockable dependencies
3. **Error Handling**: Test error paths
4. **Edge Cases**: Cover boundary conditions

### Test Organization

```go
func TestDialector_Migrate(t *testing.T) {
    tests := []struct {
        name    string
        setup   func() *gorm.DB
        want    error
        wantErr bool
    }{
        {
            name: "successful migration",
            setup: func() *gorm.DB {
                // Setup test database
            },
            want:    nil,
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Coverage Improvement

1. **Identify Gaps**: Use Codecov reports to find uncovered code
2. **Add Tests**: Write tests for uncovered functions
3. **Refactor**: Simplify complex functions
4. **Integration Tests**: Cover interaction scenarios

## 🛠️ Troubleshooting

### Common Issues

#### Coverage Not Uploading

```bash
# Check token configuration
echo $CODECOV_TOKEN

# Verify file generation
ls -la coverage.out

# Check upload logs
codecov -f coverage.out -v
```

#### Low Coverage Warnings

```bash
# Find uncovered code
go tool cover -func=coverage.out | grep -v "100.0%"

# Generate detailed report
go tool cover -html=coverage.out -o coverage.html
```

#### Status Check Failures

- Verify coverage thresholds in `.codecov.yml`
- Check for configuration syntax errors
- Ensure required coverage is achievable

### Configuration Validation

```bash
# Validate codecov.yml syntax
curl -X POST --data-binary @.codecov.yml https://codecov.io/validate
```

## 📋 Maintenance

### Regular Tasks

- **Monthly**: Review coverage trends and identify patterns
- **Per Release**: Ensure coverage meets or exceeds targets
- **PR Reviews**: Check that new code includes appropriate tests
- **Quarterly**: Update coverage targets and thresholds

### Coverage Goals

- **Short Term**: Maintain 80% project coverage
- **Medium Term**: Achieve 85% project coverage
- **Long Term**: 90%+ coverage for critical components

## 🔗 Resources

### Documentation

- [Codecov Documentation](https://docs.codecov.com/)
- [Go Testing Package](https://golang.org/pkg/testing/)
- [GORM Testing Guide](https://gorm.io/docs/testing.html)

### Tools

- [Codecov Browser Extension](https://github.com/codecov/browser-extension)
- [Go Coverage Tools](https://blog.golang.org/cover)
- [Coverage Visualization](https://github.com/alanshaw/go-coverage-badges)

### Best Practices

- [Effective Go Testing](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Test Coverage Anti-Patterns](https://testing.googleblog.com/2020/08/code-coverage-best-practices.html)
- [Go Test Patterns](https://blog.golang.org/subtests)

This comprehensive coverage setup ensures high code quality while providing detailed insights into test effectiveness
and code reliability.
