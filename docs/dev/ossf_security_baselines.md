# OSSF Security Baselines

This document provides guidance on implementing security practices aligned with
the Open Source Security Foundation (OSSF) security baselines for open source projects.

<!-- REF: https://github.com/ossf/security-baselines -->
<!-- REF: https://bestpractices.coreinfrastructure.org/en -->
<!-- REF: https://openssf.org/ -->
<!-- REF: https://github.com/ossf/scorecard -->
<!-- REF: https://github.com/ossf/allstar -->

## 🔐 Core Security Principles

The OSSF security baselines focus on these key areas:

1. **Code Security**
2. **Dependency Management**
3. **Build & Release Security**
4. **Vulnerability Disclosure**
5. **Security Testing**
6. **Documentation**

## 🛡️ AllStar Automated Policy Enforcement

[AllStar](https://github.com/ossf/allstar) is a GitHub App that enforces security
policies automatically across repositories.

To configure AllStar for your repository, follow these steps:

### 1. Install AllStar

- Visit the [AllStar GitHub App page](https://github.com/apps/allstar) and install it on your organization or repository.

### 2. Add AllStar Configuration

Create a `.allstar` directory at the root of your repository and add a `allstar.yaml` file to configure global settings.

You can also add individual policy files for fine-grained control.

#### Example: Minimal AllStar Configuration

```yaml
# .allstar/allstar.yaml
optConfig:
  optOutStrategy: false
```

#### Example: Branch Protection Policy

```yaml
# .allstar/branch_protection.yaml
enforcement: true
requireApprovalCount: 1
requireStatusChecks: true
requireUpToDateBranch: true
```

#### Example: Security Policy Enforcement

```yaml
# .allstar/security_policy.yaml
enforcement: true
```

#### Example: Outside Collaborators Policy

```yaml
# .allstar/outside_collaborators.yaml
enforcement: true
```

> **Tip:** For a full list of policies and configuration options, see the [AllStar Policy Reference](https://github.com/ossf/allstar/blob/main/docs/policies.md).

### 3. Commit and Push

Commit your `.allstar` configuration files to the default branch and push to GitHub.  
AllStar will automatically enforce the configured policies.

### 4. Review and Monitor

- Monitor AllStar's actions via GitHub Issues and repository checks.
- Adjust `.allstar/*.yaml` files as needed to refine enforcement.

---

## 📋 Security Baseline Checklist

### Code Security

- [ ] **Secure Coding Practices**
  - Follow language-specific secure coding guidelines
  - Enforce code quality standards through linters
  - Implement proper error handling
  - Validate all inputs, especially user inputs

- [ ] **Source Control Protection**
  - Protect default branches with required reviews
  - Enforce signed commits
  - Implement branch protection rules
  - Use tools to detect secrets and credentials in code

- [ ] **Authentication & Authorization**
  - Use strong authentication mechanisms
  - Implement proper authorization checks
  - Avoid hardcoded credentials
  - Practice least privilege principles

<!-- REF: https://github.com/ossf/secure-code-fundamentals -->

### Dependency Management

- [ ] **Dependency Verification**
  - Use a dependency scanning tool
  - Verify dependency integrity (checksums)
  - Implement Software Bill of Materials (SBOM)

- [ ] **Dependency Updates**
  - Regularly update dependencies
  - Automate dependency updates when possible
  - Monitor for vulnerabilities in dependencies

- [ ] **Dependency Minimization**
  - Minimize unnecessary dependencies
  - Document why each dependency is needed
  - Prefer well-maintained dependencies

<!-- REF: https://github.com/ossf/package-manager-best-practices -->

### Build & Release Security

- [ ] **Build Reproducibility**
  - Ensure builds are reproducible
  - Document build process thoroughly
  - Use automated builds

- [ ] **Artifact Signing**
  - Sign release artifacts
  - Verify signatures during installation/deployment
  - Document signature verification process

- [ ] **Supply Chain Protection**
  - Use trustworthy build environments
  - Implement CI/CD security controls
  - Consider SLSA (Supply-chain Levels for Software Artifacts) framework

<!-- REF: https://slsa.dev/ -->
<!-- REF: https://github.com/ossf/package-analysis -->

### Vulnerability Disclosure

- [ ] **Security Policy**
  - Maintain a clear SECURITY.md file
  - Define the vulnerability reporting process
  - Document supported versions

- [ ] **Vulnerability Management**
  - Track security issues appropriately
  - Provide timely fixes for security issues
  - Follow coordinated vulnerability disclosure practices

- [ ] **Security Advisories**
  - Publish security advisories for vulnerabilities
  - Use standard formats (e.g., CVE)
  - Communicate impact and mitigation clearly

<!-- REF: https://github.com/ossf/oss-vulnerability-guide -->

### Security Testing

- [ ] **Automated Testing**
  - Implement security-focused test cases
  - Use SAST (Static Application Security Testing) tools
  - Consider DAST (Dynamic Application Security Testing) if applicable
  - **Adopt modern Go testing practices** (see [Go testing best practices](https://threedots.tech/post/go-testing-practical-guide/))

- [ ] **Fuzz Testing**
  - Implement fuzzing for parsing or complex logic
  - Integrate fuzzing into CI pipeline
  - Have a process to triage fuzzing results

- [ ] **Penetration Testing**
  - Consider regular security reviews
  - Document security testing approach
  - Fix identified security issues promptly

<!-- REF: https://github.com/ossf/fuzz-introspector -->
<!-- REF: https://owasp.org/www-project-web-security-testing-guide/ -->
<!-- REF: https://threedots.tech/post/go-testing-practical-guide/ -->

### Security Documentation

- [ ] **User Documentation**
  - Document security features
  - Provide secure configuration guidance
  - Include threat model where appropriate

- [ ] **Developer Documentation**
  - Document security expectations for contributors
  - Provide security testing information
  - Include architecture security considerations

- [ ] **Security Risk Assessment**
  - Identify key security risks
  - Document trust boundaries
  - Maintain security assumptions

## 🛠️ Implementation Guidelines

### Starting Small

1. Begin with basic security hygiene:
   - Enable branch protection
   - Add a SECURITY.md file
   - Set up automated dependency scanning

2. Progress to intermediate measures:
   - Implement automated security testing
   - Sign releases
   - Create a vulnerability management process

3. Advanced security measures:
   - Generate and publish SBOMs
   - Implement fuzzing
   - Conduct regular security audits

### Tool Recommendations

#### General Purpose Tools

```bash
# Dependency scanning
$ dependency-check --project "Project Name" --scan /path/to/code

# Secret scanning
$ git-secrets --scan

# SAST tool example
$ semgrep --config=p/owasp-top-ten .

# SBOM generation
$ syft /path/to/project -o cyclonedx-json > sbom.json
```

#### Language-Specific Tools

Different tools are recommended based on programming language:

- **Go**: gosec, govulncheck, [see also: Go security testing tips](https://threedots.tech/post/go-testing-practical-guide/)
- **JavaScript/Node.js**: npm audit, eslint-plugin-security
- **Python**: bandit, safety
- **Java**: SpotBugs, OWASP Dependency Check
- **Ruby**: Brakeman, bundler-audit
- **Rust**: cargo-audit, cargo-deny

## 📊 Assessment & Improvement

### Measuring Security Maturity

The OSSF provides tools to assess your project's security posture:

- [OSSF Scorecard](https://securityscorecards.dev/): Automated checks for security best practices
- [OpenSSF Best Practices Badge Program](https://bestpractices.coreinfrastructure.org/): A way to show your project
  follows best practices

### Continuous Improvement

- Regularly review security posture
- Subscribe to security advisories for your dependencies
- Participate in security-focused communities
- Consider having periodic external security reviews

<!-- REF: https://github.com/ossf/scorecard -->
<!-- REF: https://github.com/ossf/allstar -->

## 📚 Additional Resources

- [Go testing practical guide (threedots.tech)](https://threedots.tech/post/go-testing-practical-guide/)
- [AllStar Documentation](https://github.com/ossf/allstar/blob/main/docs/)
- [AllStar Policy Reference](https://github.com/ossf/allstar/blob/main/docs/policies.md)
- [OSSF Security Insights](https://github.com/ossf/security-insights-spec): Standard format for security information
- [OSSF Security Tooling](https://github.com/ossf/wg-security-tooling): Working group on security tools
- [OSSF Best Practices](https://github.com/ossf/wg-best-practices-os-developers): Best practices for open source developers
- [OWASP Top 10](https://owasp.org/www-project-top-ten/): Common web application security risks
- [CII Best Practices](https://bestpractices.coreinfrastructure.org/): Security best practices for open source
- [OpenSSF](https://openssf.org/): Open Source Security Foundation homepage

## References

1. [OSSF Security Baselines](https://github.com/ossf/security-baselines)
2. [Open Source Security Foundation](https://openssf.org/)
3. [SLSA Framework](https://slsa.dev/)
4. [OSSF Scorecard](https://securityscorecards.dev/)
5. [OSSF Allstar](https://github.com/ossf/allstar)
6. [OWASP Top 10](https://owasp.org/www-project-top-ten/)
7. [CII Best Practices](https://bestpractices.coreinfrastructure.org/)
