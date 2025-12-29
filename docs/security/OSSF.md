# OpenSSF Best Practices

This document outlines the OpenSSF (Open Source Security Foundation) security practices implemented in glyphic.

## OpenSSF Scorecard

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/greysquirr3l/glyphic/badge)](https://scorecard.dev/viewer/?uri=github.com/greysquirr3l/glyphic)

The project aims to maintain a high OpenSSF Scorecard score through:

## Implemented Security Measures

### Branch Protection ✅

- **Required status checks**: All tests and linting must pass
- **Required PR reviews**: At least 1 approval required
- **Code owner reviews**: Security-sensitive files require owner approval
- **Dismiss stale reviews**: Reviews dismissed when new commits pushed
- **Conversation resolution**: All discussions must be resolved
- **Linear history**: No merge commits, rebase only
- **Force push protection**: No force pushes allowed
- **Deletion protection**: Branch cannot be deleted

### Automated Security Scanning ✅

- **CodeQL Analysis**: Weekly automated code scanning for vulnerabilities
- **Dependency Review**: Automated review of dependency changes in PRs
- **gosec**: Go security scanner runs on every push/PR
- **govulncheck**: Go vulnerability database scanner
- **OpenSSF Scorecard**: Weekly security posture assessment

### Dependency Management ✅

- **Dependabot**: Automated dependency updates
  - Weekly Go module updates
  - Weekly GitHub Actions updates
  - Security updates prioritized
  - Grouped updates for related packages
- **Pinned dependencies**: All dependencies locked with go.sum
- **Minimal dependencies**: Only essential, well-audited packages

### Code Quality ✅

- **golangci-lint**: Comprehensive linting (45+ linters)
- **100% clean**: Zero linting issues
- **Test coverage**: 66.4% overall, 74%+ on security-critical code
- **Code owners**: CODEOWNERS file for sensitive areas

### Security Policy ✅

- **SECURITY.md**: Comprehensive security policy
- **Vulnerability reporting**: Clear disclosure process
- **Response timeline**: Documented SLAs for security issues

### Supply Chain Security ✅

- **SBOM generation**: Software Bill of Materials via go.mod
- **License compliance**: All dependencies MIT/BSD/Apache 2.0
- **Checksum verification**: SHA-256 for all wordlists
- **HTTPS only**: Secure downloads for external resources

### Cryptographic Security ✅

- **crypto/rand only**: No insecure randomness
- **Memory locking**: Prevents swapping sensitive data
- **Secure zeroing**: Immediate clearing of password buffers
- **Quantum resistant**: High-entropy passwords (77+ bits)

## Security Checklist

Following OpenSSF Best Practices criteria:

- [x] Branch protection enabled
- [x] Code review required for changes
- [x] Automated testing in CI/CD
- [x] Static analysis (golangci-lint, CodeQL)
- [x] Dependency vulnerability scanning
- [x] Security policy documented
- [x] Vulnerability disclosure process
- [x] SBOM available (go.mod/go.sum)
- [x] Pinned dependencies
- [x] Code owners defined
- [x] No secrets in repository
- [x] Signed releases (GPG)
- [x] Memory-safe language (Go)
- [ ] SLSA provenance (planned)
- [ ] Fuzzing (planned for critical paths)

## Reporting Security Issues

See [SECURITY.md](../SECURITY.md) for vulnerability reporting procedures.

## References

- [OpenSSF Best Practices](https://www.bestpractices.dev/)
- [OpenSSF Scorecard](https://scorecard.dev/)
- [SLSA Framework](https://slsa.dev/)
- [Supply-chain Levels for Software Artifacts](https://slsa.dev/spec/v1.0/)
