# GitHub Branch Protection Configuration

This repository includes bulletproof branch protection rules designed to maintain
code quality, security, and collaboration standards.

## 🛡️ Protection Rules Overview

### Required Status Checks

- **Strict Mode**: Branches must be up-to-date before merging
- **Required Checks**:
  - `ci/build` - Build verification
  - `ci/test` - Full test suite execution
  - `ci/lint` - Code linting (golangci-lint)
  - `ci/security-scan` - Security vulnerability scanning
  - `ci/dependency-check` - Dependency vulnerability analysis
  - `ci/go-mod-tidy` - Go module cleanliness verification
  - `ci/go-vet` - Go static analysis
  - `ci/golangci-lint` - Comprehensive linting
  - `ci/coverage` - Code coverage requirements
  - `ci/integration-tests` - Integration test validation
  - `ci/vulnerability-scan` - Container/dependency scanning
  - `continuous-integration` - General CI status

### Pull Request Requirements

- **Minimum Reviews**: 2 required approving reviews
- **Stale Review Dismissal**: Automatically dismiss stale reviews on new pushes
- **Code Owner Reviews**: Required when CODEOWNERS file exists
- **Last Push Approval**: Require approval on the most recent push
- **No Bypass Allowances**: No users, teams, or apps can bypass PR requirements

### Security & Quality Controls

- **Admin Enforcement**: Even repository admins must follow these rules
- **Linear History**: Enforce linear Git history (no merge commits)
- **Force Push Protection**: Prevent destructive force pushes
- **Deletion Protection**: Prevent branch deletion
- **Conversation Resolution**: All PR conversations must be resolved before merging
- **Fork Syncing**: Allow fork syncing for easier contribution

## 🚀 Implementation

### `github-branch-protection.json` Configuration File

```json
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "ci/build",
      "ci/test",
      "ci/lint",
      "ci/security-scan",
      "ci/dependency-check",
      "ci/go-mod-tidy",
      "ci/go-vet",
      "ci/golangci-lint",
      "ci/coverage",
      "ci/integration-tests",
      "ci/vulnerability-scan",
      "continuous-integration"
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "required_approving_review_count": 2,
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": true,
    "require_last_push_approval": true,
    "bypass_pull_request_allowances": {
      "users": [],
      "teams": [],
      "apps": []
    }
  },
  "restrictions": {
    "users": [],
    "teams": [],
    "apps": []
  },
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "block_creations": false,
  "required_conversation_resolution": true,
  "lock_branch": false,
  "allow_fork_syncing": true
}
```

### GitHub CLI Method

```bash
# Apply to main branch
gh api repos/:owner/:repo/branches/main/protection \
  --method PUT \
  --input github-branch-protection.json

# Apply to develop branch
gh api repos/:owner/:repo/branches/develop/protection \
  --method PUT \
  --input github-branch-protection.json
```

### GitHub API Method

```bash
# Using curl with personal access token
curl -X PUT \
  -H "Authorization: token YOUR_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/OWNER/REPO/branches/BRANCH/protection \
  -d @github-branch-protection.json
```

### GitHub Web Interface

1. Navigate to **Settings** → **Branches**
2. Click **Add rule** or edit existing rule
3. Configure settings to match the JSON specification
4. Enable all protection features as specified

## 🔧 Customization

### Adding Status Checks

Add new required checks to the `contexts` array:

```json
"contexts": [
  "ci/build",
  "ci/test",
  "your-custom-check"
]
```

### Adjusting Review Requirements

Modify review settings as needed:

```json
"required_pull_request_reviews": {
  "required_approving_review_count": 1,
  "dismiss_stale_reviews": true,
  "require_code_owner_reviews": false
}
```

### Bypass Allowances

For emergency situations, specify users/teams that can bypass:

```json
"bypass_pull_request_allowances": {
  "users": ["emergency-user"],
  "teams": ["security-team"],
  "apps": ["dependabot"]
}
```

## 🎯 Benefits

### Security

- **Prevents Direct Pushes**: All changes must go through PR review
- **Enforces Quality Gates**: Comprehensive CI/CD validation
- **Audit Trail**: Complete history of all changes and approvals
- **Vulnerability Prevention**: Automated security scanning

### Code Quality

- **Peer Review**: Multiple developer oversight on all changes
- **Automated Testing**: Comprehensive test coverage verification
- **Linting Compliance**: Consistent code style and standards
- **Documentation**: Required conversation resolution ensures knowledge transfer

### Collaboration

- **Knowledge Sharing**: Required reviews spread domain knowledge
- **Mentorship**: Code review process facilitates learning
- **Accountability**: Clear approval trail for all changes
- **Conflict Prevention**: Linear history prevents merge conflicts

## 📋 Prerequisites

### CI/CD Pipeline

Ensure your GitHub Actions or CI system provides these status checks:

- Build and test automation
- Linting and code quality checks
- Security and vulnerability scanning
- Coverage reporting

### CODEOWNERS File

Create `.github/CODEOWNERS` for automatic reviewer assignment:

```plaintext
# Global ownership
* @greysquirr3l

# Go-specific files
*.go @greysquirr3l
go.mod @greysquirr3l
go.sum @greysquirr3l

# Documentation
*.md @greysquirr3l
docs/ @greysquirr3l

# CI/CD and GitHub configuration
.github/ @greysquirr3l
.github/workflows/ @greysquirr3l
.github/dependabot.yml @greysquirr3l

# Security files
SECURITY.md @greysquirr3l
.github/SECURITY.md @greysquirr3l

# Configuration files
.golangci.yml @greysquirr3l
.gitignore @greysquirr3l
```

### Team Structure

- **Repository Admins**: Should be limited to senior maintainers
- **Code Owners**: Domain experts for specific areas
- **Contributors**: Developers with write access following protection rules

## 🚨 Emergency Procedures

### Hotfix Process

1. Create hotfix branch from protected branch
2. Implement minimal fix with tests
3. Create PR with `[HOTFIX]` prefix
4. Expedite review process with senior maintainers
5. Merge immediately after required approvals

### Protection Bypass (Extreme Emergency)

1. Document the emergency in an issue
2. Temporarily modify protection rules via API
3. Make necessary direct changes
4. Immediately restore protection rules
5. Create follow-up PR to document changes

## 🔄 Maintenance

### Regular Review

- **Monthly**: Review protection effectiveness
- **Quarterly**: Update required status checks
- **Semi-annually**: Evaluate bypass allowances
- **Annually**: Comprehensive security audit

### Adaptation

- Add new status checks as CI/CD evolves
- Adjust review requirements based on team size
- Update bypass allowances for new tools/users
- Modify restrictions based on security needs

This configuration provides enterprise-grade protection while maintaining developer productivity and collaboration efficiency.
