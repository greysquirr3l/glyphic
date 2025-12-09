# Task Completion Checklist

When completing a task in the glyphic project, always perform these steps:

## 1. Code Quality

- [ ] Run `make fmt` to format code
- [ ] Run `make tidy` to tidy dependencies
- [ ] Run `make lint` to check for lint issues
- [ ] Verify no golangci-lint warnings

## 2. Testing

- [ ] Run `make test` to ensure all tests pass
- [ ] Check test coverage with `make test-coverage` (aim for >80%)
- [ ] Run race detector: `go test -race ./...`
- [ ] For generator changes, run `make fuzz`

## 3. Security

- [ ] Run `make security` to scan for security issues
- [ ] Verify only crypto/rand is used for randomness (never math/rand)
- [ ] Ensure sensitive data is zeroed with SecureZero()
- [ ] Check memory locking is used for password buffers

## 4. Documentation

- [ ] Update godoc comments for new/changed exported functions
- [ ] Document "why" not "what" in inline comments
- [ ] Update CHANGELOG.md if user-facing changes
- [ ] Update README.md if CLI flags or behavior changed

## 5. Build Verification

- [ ] Run `make build` to ensure clean build
- [ ] Test the binary: `./bin/glyphic --help`
- [ ] Verify version info: `./bin/glyphic --version`

## 6. Git

- [ ] Review changes with `git diff`
- [ ] Stage changes: `git add <files>`
- [ ] Write meaningful commit message following format:

  ```text
  feat: add new feature
  fix: resolve bug
  refactor: improve code structure
  test: add test coverage
  docs: update documentation
  ```

- [ ] Run `make check` before committing

## Quick Completion Command

```bash
# Run all checks at once
make check && git status
```

## Before Pull Request

- [ ] All tests passing
- [ ] No lint issues
- [ ] Security scan clean
- [ ] Code coverage maintained/improved
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] Build successful for all platforms: `make build-all`
