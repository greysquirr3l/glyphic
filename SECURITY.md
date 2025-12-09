# Security Policy

## Supported Versions

We actively support the following versions of glyphic with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.1   | :white_check_mark: |
| 0.1.0   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability in glyphic, please report it responsibly:

### How to Report

1. **DO NOT** open a public GitHub issue
2. Email security concerns to: [s0ma@protonmail.com]
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity
  - Critical: 1-7 days
  - High: 7-30 days
  - Medium: 30-90 days
  - Low: Next regular release

### Security Best Practices

When using glyphic:

1. **Keep Updated**: Always use the latest version
2. **Secure Environment**: Run in a trusted environment
3. **Memory Protection**: The tool uses memory locking where possible
4. **No Persistence**: Passwords are never written to disk
5. **Secure Randomness**: Only crypto/rand is used for randomness

## Security Features

### Cryptographic Security

- **PRNG**: Uses `crypto/rand` exclusively
- **Validation**: PRNG validated on startup
- **No Fallbacks**: Never falls back to insecure randomness

### Memory Security

- **Memory Locking**: Uses `unix.Mlock()` to prevent swapping
- **Secure Zeroing**: All password buffers zeroed after use
- **Runtime Protection**: Uses `runtime.KeepAlive()` for timing

### Data Security

- **No Persistence**: Passwords never written to disk or logs
- **No Network**: Only HTTPS for initial wordlist downloads
- **Checksum Verification**: SHA-256 verification for all wordlists
- **Constant-Time**: Sensitive comparisons use constant-time operations

### Supply Chain Security

- **Go Modules**: Dependencies locked with go.sum
- **Minimal Dependencies**: Only essential, well-audited packages
- **License Compliance**: All dependencies MIT/BSD/Apache 2.0

## Known Limitations

1. **Memory Locking**: Requires appropriate OS permissions
2. **Terminal Detection**: May not detect all terminal capabilities
3. **Wordlist Downloads**: Initial download requires internet connection

## Disclosure Policy

- We follow coordinated vulnerability disclosure
- Security patches released as soon as possible
- CVEs assigned for significant vulnerabilities
- Public disclosure after patch availability

## Security Audit History

- No formal security audits conducted yet
- Community review welcome
- Continuous security improvements

## Contact

For security-related questions (not vulnerabilities):
- GitHub Discussions: https://github.com/greysquirr3l/glyphic/discussions
- Issues: https://github.com/greysquirr3l/glyphic/issues (non-security)
