# 🔄 GitHub Actions CI/CD Best Practices

> *Last updated: September 2025*

This document provides comprehensive guidance for implementing effective CI/CD workflows
with GitHub Actions, incorporating advanced patterns and security-first principles.

<!-- REF: https://docs.github.com/en/actions -->
<!-- REF: https:/## 🧪 Code Quality and Testing

- **Linting**: Check code quality in early stages
- **Automated tests**: Run unit, integration, and e2e tests
- **Coverage reports**: Track test coverage over time
- **Build artifacts**: Upload build artifacts for verification

```yaml
steps:
  - name: Run tests
    run: npm test

  - name: Upload coverage reports
    uses: actions/upload-artifact@v4
    with:
      name: coverage-report
      path: coverage/
```

## 🧪 Local Workflow Testing

### Testing with `act`

Instead of repeatedly pushing commits to test workflow changes, use [`act`](https://github.com/nektos/act) to run GitHub Actions locally:

```bash
# Install act (macOS)
brew install act

# Run workflow for push events
act push

# Run workflow for pull request events  
act pull_request

# Run specific job
act -j build

# Run with custom event payload
act -e event.json
```

#### Setting up Local Testing Environment

```yaml
# .actrc - Configuration file for act
-P ubuntu-latest=catthehacker/ubuntu:act-latest
-P ubuntu-22.04=catthehacker/ubuntu:act-22.04
--container-daemon-socket -
```

```bash
# Create mock event for testing
cat > event.json << EOF
{
  "ref": "refs/heads/feature-branch",
  "repository": {
    "name": "my-repo",
    "full_name": "myorg/my-repo"
  }
}
EOF

# Test workflow with custom event
act push -e event.json
```

#### Benefits of Local Testing

- **Fast feedback**: No need to commit/push for every workflow change
- **Cost savings**: Avoid consuming GitHub Actions minutes during development
- **Debugging**: Easier to debug workflow issues locally
- **Offline development**: Work on workflows without internet connectivity

> 💡 **Tip**: Use act during workflow development to reduce feedback cycles from minutes to seconds.pulse/mastering-cicd-best-practices-github-actions-tatiana-sava -->
<!-- REF: https://hackernoon.com/a-practical-guide-to-building-smarter-github-workflows -->

## 📋 Workflow Structure

### Understanding Workflows vs Actions

> 💡 **Important Distinction**: Use correct semantics - GitHub **Workflows** are the automated processes, GitHub **Actions** are reusable components within workflows.

**GitHub Workflows:**
- Configurable automated processes that run one or more jobs
- Defined by YAML files in `.github/workflows/` directory  
- Triggered by repository events, manually, or on schedule
- Comparable to Jenkins jobs but stored as code in your repository

**GitHub Actions:**
- Reusable components that perform specific tasks within workflows
- Pre-defined sets of jobs or code that reduce repetitive workflow code
- Available from GitHub Marketplace or custom-built
- Examples: `actions/checkout`, `actions/setup-node`

### Basic Workflow Anatomy

```yaml
name: CI Pipeline

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up environment
        uses: actions/setup-node@v4
        with:
          node-version: '18'
          cache: 'npm'  # Built-in caching feature
      - name: Install dependencies
        run: npm ci
      - name: Run tests
        run: npm test
```

### Key Components

- **Workflow file**: YAML files in `.github/workflows/`
- **Triggers** (`on`): Events that start the workflow
- **Jobs**: Groups of steps that run on the same runner
- **Steps**: Individual tasks that run commands or actions
- **Actions**: Reusable units of code (e.g., `actions/checkout@v3`)
- **Runners**: VMs that execute the jobs

## 🚀 Best Practices

### 1. Workflow Design

- **Keep workflows focused**: One workflow per logical CI/CD phase
- **Use descriptive names**: Clear workflow, job, and step names
- **Modularize with composite actions**: Create reusable components
- **Trigger precision**: Limit workflow runs to relevant files/branches
- **Prefer GitHub Actions over ad-hoc commands**: Leverage the ecosystem

```yaml
# Good practice: Focused triggers
on:
  push:
    branches: [ main ]
    paths:
      - 'src/**'
      - 'package.json'
  pull_request:
    branches: [ main ]
    paths-ignore:
      - '**/*.md'
```

### 2. 📚 Choosing the Right Actions

When selecting actions from GitHub Marketplace, apply the same diligence as choosing any dependency:

#### Action Selection Criteria

- **Vendor trust**: Prefer verified creators (blue checkmark)
- **License compatibility**: Check Apache v2 vs GPL vs proprietary
- **Source code access**: Review repository for quality and activity
- **Maintenance status**: Check commit frequency and issue resolution
- **Documentation quality**: Ensure comprehensive guides and examples
- **Community support**: Active discussions and contributions

#### Built-in Action Features

Many popular actions have built-in features that eliminate manual configuration:

```yaml
# Instead of manual caching setup
- uses: actions/cache@v4
  with:
    path: ~/.m2/repository
    key: ${{ runner.os }}-maven-${{ hashFiles('**/pom.xml') }}

# Use built-in caching in setup actions
- uses: actions/setup-java@v4
  with:
    distribution: 'temurin'
    java-version: '21'
    cache: 'maven'  # Built-in Maven caching

# Node.js with built-in npm caching
- uses: actions/setup-node@v4
  with:
    node-version: '18'
    cache: 'npm'  # Automatic package-lock.json based caching
```

> 💡 **Tip**: Thoroughly read action documentation. One hour of reading can save days of configuration work.

### 3. ⚡ Performance Optimization

- **Use dependency caching**: Speed up builds by caching dependencies
- **Conditional execution**: Skip unnecessary steps
- **Job parallelization**: Run independent jobs concurrently
- **Matrix builds**: Test across multiple configurations in parallel
- **Pin to specific versions**: Avoid `latest` tags for predictable builds

```yaml
# Example: Pinning to specific Ubuntu version
jobs:
  build:
    runs-on: ubuntu-22.04  # Instead of ubuntu-latest

# Example: Matrix builds
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-22.04, windows-2022, macos-12]
        node-version: [18.x, 20.x, 22.x]
```

### 4. 🔒 Security

- **Protect secrets**: Use GitHub Secrets for sensitive data
- **Restrict permissions**: Apply principle of least privilege
- **Pin action versions**: Use specific versions, ideally by SHA (SHA stapling) for maximum security
- **Scan for vulnerabilities**: Include security scanning in workflows
- **Review third-party actions**: Audit external actions for security risks

#### Pinning Actions by SHA (SHA Stapling)

Pinning actions by SHA (also called SHA stapling) ensures that your workflow always
uses the exact same version of an action, protecting against supply chain attacks
and unexpected changes. Avoid using only tags (like `@v3`), as tags can be moved
to point to different code. Instead, use the full commit SHA:

```yaml
steps:
  - uses: actions/checkout@a81bbbf8298c0fa03ea29cdc473d45769f953675 # v3.2.0
```

**Best practices for SHA stapling:**

- Always use the full 40-character commit SHA when referencing third-party actions.
- Document the original tag or version in a comment for clarity.
- Regularly review and update SHAs to pick up security patches and improvements.
- Use tools or scripts to automate checking for new SHAs of trusted actions.
- Never use untrusted or unknown actions, even if pinned by SHA.

For more, see [GitHub Actions Security Best Practices](https://blog.gitguardian.com/github-actions-security-best-practices/)
and [GitHub Docs: Security hardening for GitHub Actions](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions).

```yaml
# Good practice: Define specific permissions
permissions:
  contents: read
  pull-requests: write

# Good practice: Pin actions to specific SHAs
steps:
  - uses: actions/checkout@a81bbbf8298c0fa03ea29cdc473d45769f953675 # v3.2.0
```

### 5. 📝 Workflow Management

- **Document workflows**: Add comments explaining complex steps
- **Self-contained workflows**: Minimize external dependencies
- **Re-usable workflows**: Define common workflows in separate files
- **Timeout limits**: Add timeouts to prevent stuck jobs

```yaml
# Example: Job timeout
jobs:
  build:
    runs-on: ubuntu-22.04
    timeout-minutes: 15
    steps:
      # ...
```

### 6. 📊 Job Summaries and Enhanced Developer Experience

Use GitHub's job summary feature to highlight important information without diving into logs:

```yaml
# Example: Custom job summary with important metrics
- name: Display Test Results
  if: always()  # Run regardless of previous step outcomes
  run: |
    echo "## Test Results Summary 📊" >> $GITHUB_STEP_SUMMARY
    echo "- Tests run: $(cat test-results.json | jq '.total')" >> $GITHUB_STEP_SUMMARY  
    echo "- Failures: $(cat test-results.json | jq '.failures')" >> $GITHUB_STEP_SUMMARY
    echo "- Coverage: $(cat coverage-summary.json | jq '.coverage')%" >> $GITHUB_STEP_SUMMARY
    echo "" >> $GITHUB_STEP_SUMMARY
    echo "### Build Information" >> $GITHUB_STEP_SUMMARY
    echo "- Branch: ${{ github.ref_name }}" >> $GITHUB_STEP_SUMMARY
    echo "- Commit: ${{ github.sha }}" >> $GITHUB_STEP_SUMMARY
    echo "- Triggered by: ${{ github.event_name }}" >> $GITHUB_STEP_SUMMARY

# Example: Show workflow inputs for debugging
- name: Show Workflow Inputs
  run: |
    echo "## Workflow Configuration 🔧" >> $GITHUB_STEP_SUMMARY
    echo "- Environment: ${{ inputs.environment }}" >> $GITHUB_STEP_SUMMARY
    echo "- Deploy: ${{ inputs.deploy }}" >> $GITHUB_STEP_SUMMARY
    echo "- Debug mode: ${{ inputs.debug }}" >> $GITHUB_STEP_SUMMARY
```

> 💡 **Tip**: Job summaries support GitHub Flavored Markdown and make workflow results much more accessible to team members.

## 🔄 Understanding Workflow Lifecycle

### Sequential Execution and Failure Handling

Workflows execute steps sequentially within each job. By default, a failing step cancels all remaining steps and marks the workflow as failed. However, you can control this behavior using conditional execution.

#### Status Check Functions

GitHub provides several status check functions for conditional execution:

- `success()` - Default condition, step runs only if previous steps succeeded
- `always()` - Step runs regardless of previous step outcomes
- `failure()` - Step runs only if any previous step failed  
- `cancelled()` - Step runs only if workflow was cancelled

#### Advanced Conditional Execution

```yaml
jobs:
  test-and-report:
    runs-on: ubuntu-22.04
    steps:
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v4.1.1
      
      - name: Build application
        id: build
        run: |
          # This might fail
          ./build.sh
      
      - name: Run tests
        id: test  
        if: steps.build.outcome == 'success'
        run: |
          ./run-tests.sh
          
      - name: Generate test report
        if: always() && steps.test.conclusion != 'skipped'
        uses: test-summary/action@v2
        with:
          paths: "test-results.xml"
          
      - name: Upload failure logs  
        if: failure()
        uses: actions/upload-artifact@v4
        with:
          name: failure-logs
          path: logs/

      - name: Notify on failure
        if: failure() && github.ref == 'refs/heads/main'
        run: |
          curl -X POST "${{ secrets.SLACK_WEBHOOK_URL }}" \
            -H 'Content-type: application/json' \
            --data '{"text":"🚨 Main branch build failed!"}'
```

#### Step Outcome vs Conclusion

- **outcome**: The result before `continue-on-error` is applied (`success`, `failure`, `cancelled`, `skipped`)
- **conclusion**: The final result after `continue-on-error` is applied (`success`, `failure`, `cancelled`, `skipped`, `neutral`)

```yaml
- name: Optional step that might fail
  id: optional
  continue-on-error: true
  run: ./optional-command.sh

- name: Check optional step results
  run: |
    echo "Outcome: ${{ steps.optional.outcome }}"        # Could be 'failure'
    echo "Conclusion: ${{ steps.optional.conclusion }}"  # Will be 'success' due to continue-on-error
```

> 💡 **Tip**: Use workflow lifecycle controls to create robust pipelines that provide useful feedback even when things go wrong.

- **Linting**: Check code quality in early stages
- **Automated tests**: Run unit, integration, and e2e tests
- **Coverage reports**: Track test coverage over time
- **Build artifacts**: Upload build artifacts for verification

```yaml
steps:
  - name: Run tests
    run: npm test

  - name: Upload coverage reports
    uses: actions/upload-artifact@v3
    with:
      name: coverage-report
      path: coverage/
```

## 🚢 Deployment

### Environments and Approvals

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v3
      - name: Deploy
        run: ./deploy.sh
```

### Continuous Deployment Practices

- **Environment-specific workflows**: Separate workflows for dev/staging/prod
- **Approval gates**: Require approval for sensitive environments
- **Deployment verification**: Add post-deployment verification steps
- **Rollback capability**: Plan for failed deployments

## 🔄 Advanced Patterns

### 1. Smart Workflow Composition

#### Reusable Workflows with Input Validation

```yaml
# .github/workflows/reusable-build.yml
name: Reusable Build Workflow
on:
  workflow_call:
    inputs:
      environment:
        required: true
        type: string
        description: 'Deployment environment'
      node-version:
        required: false
        type: string
        default: '18'
        description: 'Node.js version'
      run-tests:
        required: false
        type: boolean
        default: true
        description: 'Whether to run tests'
    outputs:
      build-version:
        description: "The version that was built"
        value: ${{ jobs.build.outputs.version }}

jobs:
  validate-inputs:
    runs-on: ubuntu-22.04
    steps:
      - name: Validate environment
        run: |
          if [[ ! "${{ inputs.environment }}" =~ ^(dev|staging|prod)$ ]]; then
            echo "::error::Invalid environment: ${{ inputs.environment }}"
            exit 1
          fi

  build:
    needs: validate-inputs
    runs-on: ubuntu-22.04
    outputs:
      version: ${{ steps.version.outputs.version }}
    steps:
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v4.1.1
      - name: Setup Node.js
        uses: actions/setup-node@5e21ff4d9bc1a8cf6de233a3057d20ec6b3fb69d # v4.0.3
        with:
          node-version: ${{ inputs.node-version }}
          cache: 'npm'
          
      - name: Get version
        id: version
        run: echo "version=$(npm run version --silent)" >> $GITHUB_OUTPUT
        
      - name: Build
        run: npm run build:${{ inputs.environment }}
        
      - name: Test
        if: inputs.run-tests
        run: npm test

# Usage in main workflow
# .github/workflows/ci.yml  
jobs:
  build-dev:
    uses: ./.github/workflows/reusable-build.yml
    with:
      environment: 'dev'
      node-version: '20'
      run-tests: true
```

#### Dynamic Matrix Generation

```yaml
jobs:
  generate-matrix:
    runs-on: ubuntu-22.04
    outputs:
      matrix: ${{ steps.set-matrix.outputs.matrix }}
    steps:
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v4.1.1
      - name: Generate test matrix
        id: set-matrix
        run: |
          # Generate matrix based on changed files or other conditions
          if git diff --name-only HEAD~1 | grep -q "^frontend/"; then
            MATRIX='["frontend-unit", "frontend-e2e", "integration"]'
          else
            MATRIX='["backend-unit", "integration"]'
          fi
          echo "matrix=$MATRIX" >> $GITHUB_OUTPUT

  test:
    needs: generate-matrix
    strategy:
      matrix:
        test-suite: ${{ fromJson(needs.generate-matrix.outputs.matrix) }}
    runs-on: ubuntu-22.04
    steps:
      - name: Run ${{ matrix.test-suite }}
        run: npm run test:${{ matrix.test-suite }}
```

### 2. Intelligent Caching Strategies

#### Multi-layer Cache Strategy

```yaml
- name: Cache with fallback
  uses: actions/cache@88522ab9f39a2ea568f7027eddc7d8d8bc9d59c8 # v4.0.1
  with:
    path: |
      ~/.npm
      node_modules
      .next/cache
    key: ${{ runner.os }}-nextjs-${{ hashFiles('package-lock.json') }}-${{ hashFiles('**/*.js', '**/*.jsx', '**/*.ts', '**/*.tsx') }}
    restore-keys: |
      ${{ runner.os }}-nextjs-${{ hashFiles('package-lock.json') }}-
      ${{ runner.os }}-nextjs-
      ${{ runner.os }}-

# Cache invalidation based on conditions
- name: Clear cache on dependency changes
  if: steps.cache.outputs.cache-hit != 'true'
  run: |
    echo "Dependencies changed, clearing cache..."
    rm -rf node_modules/.cache
```

#### Conditional Dependency Installation

```yaml
- name: Get npm cache directory
  id: npm-cache-dir
  shell: bash
  run: echo "dir=$(npm config get cache)" >> $GITHUB_OUTPUT

- name: Cache npm dependencies
  uses: actions/cache@88522ab9f39a2ea568f7027eddc7d8d8bc9d59c8 # v4.0.1
  id: npm-cache
  with:
    path: ${{ steps.npm-cache-dir.outputs.dir }}
    key: ${{ runner.os }}-node-${{ hashFiles('**/package-lock.json') }}
    restore-keys: |
      ${{ runner.os }}-node-

- name: Install dependencies
  if: steps.npm-cache.outputs.cache-hit != 'true'
  run: npm ci --prefer-offline --no-audit
```

### 3. Custom Actions

Create custom composite actions for project-specific repeated tasks:

```yaml
# .github/actions/setup-project/action.yml
name: 'Setup Project Environment'
description: 'Sets up the complete project environment with caching'
inputs:
  node-version:
    description: 'Node.js version'
    required: false
    default: '18'
  install-deps:
    description: 'Install dependencies'
    required: false
    default: 'true'
outputs:
  cache-hit:
    description: 'Whether dependencies were cached'
    value: ${{ steps.cache.outputs.cache-hit }}

runs:
  using: 'composite'
  steps:
    - name: Setup Node.js
      uses: actions/setup-node@5e21ff4d9bc1a8cf6de233a3057d20ec6b3fb69d # v4.0.3
      with:
        node-version: ${{ inputs.node-version }}
        
    - name: Cache dependencies
      id: cache
      uses: actions/cache@88522ab9f39a2ea568f7027eddc7d8d8bc9d59c8 # v4.0.1
      with:
        path: |
          ~/.npm
          node_modules
        key: ${{ runner.os }}-node-${{ inputs.node-version }}-${{ hashFiles('**/package-lock.json') }}
        
    - name: Install dependencies
      if: inputs.install-deps == 'true' && steps.cache.outputs.cache-hit != 'true'
      shell: bash
      run: npm ci

# Usage in workflows
jobs:
  test:
    runs-on: ubuntu-22.04
    steps:
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v4.1.1
      - name: Setup project
        uses: ./.github/actions/setup-project
        with:
          node-version: '20'
      - name: Run tests  
        run: npm test
```

### 4. Advanced Workflow Patterns

#### Environment-based Workflow Selection

```yaml
name: Smart Deploy
on:
  push:
    branches: [main, develop, 'release/**']

jobs:
  determine-environment:
    runs-on: ubuntu-22.04
    outputs:
      environment: ${{ steps.env.outputs.environment }}
      should-deploy: ${{ steps.env.outputs.should-deploy }}
    steps:
      - name: Determine environment
        id: env
        run: |
          if [[ "${{ github.ref }}" == "refs/heads/main" ]]; then
            echo "environment=production" >> $GITHUB_OUTPUT
            echo "should-deploy=true" >> $GITHUB_OUTPUT
          elif [[ "${{ github.ref }}" == "refs/heads/develop" ]]; then
            echo "environment=staging" >> $GITHUB_OUTPUT  
            echo "should-deploy=true" >> $GITHUB_OUTPUT
          elif [[ "${{ github.ref }}" == refs/heads/release/* ]]; then
            echo "environment=preview" >> $GITHUB_OUTPUT
            echo "should-deploy=true" >> $GITHUB_OUTPUT
          else
            echo "environment=none" >> $GITHUB_OUTPUT
            echo "should-deploy=false" >> $GITHUB_OUTPUT
          fi

  deploy:
    needs: determine-environment
    if: needs.determine-environment.outputs.should-deploy == 'true'
    environment: ${{ needs.determine-environment.outputs.environment }}
    runs-on: ubuntu-22.04
    steps:
      - name: Deploy to ${{ needs.determine-environment.outputs.environment }}
        run: |
          echo "Deploying to ${{ needs.determine-environment.outputs.environment }}"
          # Deployment commands here
```

Monitor workflow performance, success rates, and durations.

## 🛡️ Security Hardening for GitHub Actions

GitHub Actions workflows can introduce significant security risks if not properly
hardened. This section provides comprehensive security guidance based on industry
best practices and official recommendations.

### Core Security Principles

1. **Defense in Depth**: Apply multiple security controls at different layers
2. **Least Privilege**: Grant minimal permissions required for each workflow
3. **Trust Verification**: Validate all external inputs and dependencies
4. **Secure by Default**: Start with secure patterns and override only when necessary

### Supply Chain Security

#### 🔍 Controlling Action Usage

- **Pin actions using SHA digests** (not just versions/tags which can be moved):

```yaml
# Vulnerable - using only a tag
- uses: actions/checkout@v3

# Better - using a specific version but still risky
- uses: actions/checkout@v3.1.0

# Best - pinning with full SHA
- uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v3.5.2
```

- **Verify action authors** and prefer actions from:
  - GitHub-owned actions (`actions/*`, `github/*`)
  - Verified creator actions (blue checkmark in Marketplace)
  - Actions from your own organization

- **Audit third-party actions** before use:
  - Review source code for malicious behavior
  - Check for excessive permissions
  - Verify maintenance status and security practices
  - Consider forking and hosting trusted versions internally

```yaml
# Use GitHub's secure action for downloading artifacts instead of 
# arbitrary scripts or custom actions
- name: Download artifact
  uses: actions/download-artifact@6b208ae046db98c579e8a3aa621ab581ff575935 # v4.1.1
  with:
    name: my-artifact
```

#### 🛑 Preventing Script Injection

- **Validate and sanitize all workflow inputs**, especially in:
  - `pull_request_target` events (which have write access)
  - Issue comments that trigger workflows
  - User-controlled inputs that flow to commands

```yaml
# Vulnerable to injection
- name: Run script with user input
  run: ./script.sh ${{ github.event.issue.title }}

# Better - validate input format
- name: Validate input
  run: |
    if ! [[ "${{ github.event.issue.title }}" =~ ^[a-zA-Z0-9_-]+$ ]]; then
      echo "Invalid input format"
      exit 1
    fi
    ./script.sh "${{ github.event.issue.title }}"
```

- **Use workflow expressions for input validation**:

```yaml
# Validate inputs with expressions
- name: Set safe branch name
  run: |
    SAFE_BRANCH="${{ 
      github.event.inputs.branch_name && 
      github.event.inputs.branch_name != '' && 
      matches(github.event.inputs.branch_name, '^[a-zA-Z0-9_/-]+$') && 
      github.event.inputs.branch_name || 
      'main' 
    }}"
    echo "Using branch: $SAFE_BRANCH"
```

### Permissions and Access Control

#### 🔑 Fine-Grained Permissions Model

- **Define workflow-level permissions** to minimize attack surface:

```yaml
name: Limited Scope Workflow

# Set default permissions to minimum required
permissions: {}  # No permissions by default

jobs:
  test:
    permissions:
      # Only grant permissions actually needed
      contents: read
      issues: write
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v3.5.2
      - name: Run tests
        run: ./run-tests.sh
```

- **Available permission scopes**:

| Scope | Description |
|-------|-------------|
| `actions` | Manage GitHub Actions |
| `checks` | Read/write check runs |
| `contents` | Read/write repository contents |
| `deployments` | Manage deployments |
| `id-token` | Request OpenID Connect tokens |
| `issues` | Manage issues |
| `packages` | Manage packages |
| `pull-requests` | Manage pull requests |
| `repository-projects` | Manage repository projects |
| `security-events` | Manage security events |
| `statuses` | Manage commit statuses |

#### 🔐 Secure Token and Secret Management

- **Use GitHub's built-in secrets** and never hardcode sensitive values
- **Never log secrets** or expose them in error messages
- **Limit secret access** to specific environments and branches
- **Implement secret rotation** practices

```yaml
# Good practice - Using repository secrets
steps:
  - name: Deploy to production
    run: ./deploy.sh
    env:
      API_TOKEN: ${{ secrets.API_TOKEN }}
```

- **Use environment-specific secrets** for sensitive deployments:

```yaml
name: Deploy

jobs:
  deploy-to-prod:
    environment: production  # Restricts access to production secrets
    runs-on: ubuntu-latest
    steps:
      - name: Deploy
        env:
          PROD_SECRET: ${{ secrets.PROD_SECRET }}
        run: ./deploy.sh
```

#### 🌐 Cloud Provider Authentication with OIDC

- **Use OpenID Connect (OIDC)** instead of storing long-lived cloud credentials:

```yaml
# AWS authentication with OIDC
jobs:
  deploy:
    permissions:
      id-token: write  # Required for OIDC authentication
      contents: read
    steps:
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@8c3f20df09ac63af7b3ae3d7c91f105f857d8497 # v4.0.0
        with:
          role-to-assume: arn:aws:iam::123456789012:role/my-github-actions-role
          aws-region: us-east-1
```

### Runner Security

#### 🖥️ Protecting Self-Hosted Runners

- **Use ephemeral runners** that are created and destroyed for each job
- **Isolate runners** by workload or sensitivity level
- **Implement defense-in-depth** for runner environments
- **Never use self-hosted runners** for public repositories with external contributors

```yaml
# Example of using labels to target specific runner groups
jobs:
  sensitive-job:
    runs-on: [self-hosted, isolated-secure-runner]
    steps:
      # ...job steps...
```

#### 🧹 Secure Checkout and Build Practices

- **Clean workspace** before and after jobs:

```yaml
steps:
  - name: Clean workspace
    run: |
      if [ -d "$GITHUB_WORKSPACE" ]; then
        rm -rf "$GITHUB_WORKSPACE"/*
      fi

  # ... job steps ...

  - name: Clean up sensitive files
    if: always()  # Run even if previous steps fail
    run: |
      find . -name "*.key" -delete
      find . -name "*.pem" -delete
```

- **Consider using hardened runner actions** like `step-security/harden-runner`:

```yaml
steps:
  - uses: step-security/harden-runner@17d0e2bd7d51742c71671bd71b76b365dc6242f4 # v2.7.0
    with:
      egress-policy: audit
      disable-sudo: true

  - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v3.5.2
```

### Workflow Security

#### ⚠️ Securing Event Triggers

- **Be cautious with `pull_request_target`** as it runs with repository token:

```yaml
# Safer approach for pull_request_target events
on:
  pull_request_target:
    types: [opened, synchronize]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      # Check out the base branch, NOT the PR code
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v3.5.2
        with:
          ref: ${{ github.base_ref }}
          
      # Only after validation, checkout PR code
      - name: Checkout PR
        uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v3.5.2
        with:
          ref: ${{ github.event.pull_request.head.sha }}
```

- **Secure `workflow_dispatch` inputs**:

```yaml
on:
  workflow_dispatch:
    inputs:
      environment:
        type: choice
        description: 'Deployment environment'
        required: true
        options:
          - development
          - staging
          - production
```

#### 📦 Preventing Dependency Confusion Attacks

- **Pin all dependencies** in your project and CI/CD pipeline
- **Use lockfiles** and checksums to verify dependency integrity
- **Configure private registry** ahead of public ones

```yaml
# Example of securing npm install
steps:
  - name: Verify and install dependencies
    run: |
      npm ci --package-lock-only
      npm audit --audit-level=high
      npm ci
```

### Implementing Advanced Security Controls

#### 🛡️ Branch Protection and Required Reviews

- **Enforce branch protection rules** for production code:
  - Require pull request reviews before merging
  - Require status checks to pass
  - Require signed commits
  - Do not allow bypassing the above settings

#### 🔍 Security Scanning Integration

- **Integrate security scanners** into your workflow:

```yaml
jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v3.5.2
      
      # CodeQL Analysis
      - name: Initialize CodeQL
        uses: github/codeql-action/init@cdcdbb579706841c47f7063dda365e292e5cad7a # v2.2.2
        with:
          languages: javascript, python
          
      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@cdcdbb579706841c47f7063dda365e292e5cad7a # v2.2.2
        
      # Dependency scanning
      - name: Dependency Review
        uses: actions/dependency-review-action@7bbfa034e752445ea40215fff1c3bf9597993d3f # v3.1.3
```

- **Run OSSF Scorecard** to evaluate security practices:

```yaml
name: OSSF Scorecard
on:
  branch_protection_rule:
  schedule:
    - cron: '0 0 * * 0'
  push:
    branches: [ main ]

jobs:
  scorecard:
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      id-token: write
    steps:
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v3.5.2
        with:
          persist-credentials: false
      
      - name: Run Scorecard
        uses: ossf/scorecard-action@dc50aa9dfc5b2356601546948edc5ca5f8a443f6 # v2.3.1
        with:
          results_file: results.sarif
          results_format: sarif
          publish_results: true
```

### 📊 Continuous Security Monitoring

- **Implement workflow auditing tools** to detect issues over time
- **Regularly review workflow changes** for security implications
- **Monitor and analyze workflow runs** for unusual patterns
- **Set up alerts** for failed security checks

#### Security Automation with Tools

```yaml
name: Security Automation

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 0 * * 0'  # Weekly scan

jobs:
  security-checks:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4 # v3.5.2
      
      # Find vulnerable dependencies
      - name: Run Dependabot
        uses: github/dependabot-action@e7d862d8466a5f005828579a94f880247.17.2
        
      # Scan for secrets in code
      - name: Detect leaked secrets
        uses: gitleaks/gitleaks-action@8a5b84c7baa459a01af412be2d3438dbd5d0b73d # v2.3.2
        
      # Container security scanning
      - name: Scan container image
        uses: anchore/scan-action@8b2b1e5912ab932b9e355276ee871ec043a10746 # v3.3.6
        with:
          image: "my-organization/my-app:latest"
```

## 🚨 Security Checklist for GitHub Actions

Use this checklist to evaluate the security posture of your GitHub Actions workflows:

- [ ] **Actions and Dependencies**
  - [ ] All actions are pinned using full SHA hashes
  - [ ] Third-party actions are reviewed and trusted
  - [ ] Dependencies are locked and verified
  - [ ] Avoid actions that require excessive permissions

- [ ] **Permissions and Access**
  - [ ] Workflows use minimal required permissions
  - [ ] OIDC is used for cloud provider authentication
  - [ ] Environment protection rules are enabled for production deployments
  - [ ] Branch protection rules are enforced

- [ ] **Input Validation and Execution**
  - [ ] User-controlled inputs are validated and sanitized
  - [ ] Scripts avoid command injection vulnerabilities
  - [ ] Workflow expressions are used for input validation
  - [ ] Cautious handling of `pull_request_target` events

- [ ] **Secrets Management**
  - [ ] Secrets are stored in GitHub Secrets, not in code
  - [ ] Secrets are scoped to specific environments
  - [ ] No secrets are exposed in logs or outputs
  - [ ] Secret rotation process is documented and followed

- [ ] **Runner Security**
  - [ ] Use GitHub-hosted runners when possible
  - [ ] Self-hosted runners are ephemeral and isolated
  - [ ] Runner environments are hardened and secured
  - [ ] Workspaces are cleaned before and after jobs

- [ ] **Security Scanning**
  - [ ] Code scanning is enabled (e.g., CodeQL)
  - [ ] Dependency scanning is integrated
  - [ ] OSSF Scorecard runs regularly
  - [ ] Container images are scanned for vulnerabilities

- [ ] **Continuous Monitoring**
  - [ ] Security checks run on a schedule
  - [ ] Failed security checks trigger alerts
  - [ ] Workflow changes undergo security review
  - [ ] Audit logs are monitored for suspicious activity

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Actions Marketplace](https://github.com/marketplace?type=actions)
- [GitHub Actions Security Hardening](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
- [GitHub Actions Workflow Syntax Reference](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)
- [Act - Run GitHub Actions Locally](https://github.com/nektos/act)
- [GitHub CLI](https://cli.github.com/)
- [OpenSSF Source Code Management Best Practices](https://best.openssf.org/SCM-BestPractices/github/)
- [OSSF Scorecard](https://github.com/ossf/scorecard)
- [GitHub Actions Security Cheat Sheet](https://blog.gitguardian.com/github-actions-security-cheat-sheet/)
- [GitHub Security Lab](https://securitylab.github.com/)
- [Security Hardening with OpenID Connect](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
- [StepSecurity Hardening Tools](https://www.stepsecurity.io/)

## 📝 Smart Workflow Development Summary

### Key Principles from Industry Experience

1. **Use Correct Semantics**
   - GitHub **Workflows** = Automated processes defined in YAML
   - GitHub **Actions** = Reusable components within workflows

2. **Action Selection Best Practices**
   - Prefer GitHub Actions over custom scripts when possible
   - Choose actions with same rigor as any dependency selection
   - Thoroughly read documentation - one hour reading saves days of work
   - Pin actions to specific commit SHAs for security

3. **Version Management**
   - Pin runner images to specific versions (e.g., `ubuntu-22.04` not `ubuntu-latest`)
   - Use commit SHAs for actions: `actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada58f8c4`
   - Leverage built-in caching features in popular actions

4. **Workflow Lifecycle Mastery**
   - Understand sequential execution within jobs
   - Use conditional execution (`if: always()`, `if: failure()`) strategically
   - Distinguish between step `outcome` and `conclusion`
   - Design workflows to provide useful feedback even on failures

5. **Developer Experience**
   - Use job summaries (`$GITHUB_STEP_SUMMARY`) for key information
   - Test workflows locally with `act` for faster feedback
   - Create custom composite actions for repeated patterns
   - Implement smart caching strategies with fallback keys

6. **Production Readiness**
   - Pin all versions for reproducible builds
   - Implement proper error handling and notifications
   - Use environment-specific deployment workflows
   - Monitor workflow performance and success rates

> 💡 **Remember**: Workflows are code - apply the same quality standards, testing practices, and security considerations you would to any other codebase.

## 📖 References and Further Reading

1. [GitHub Actions Documentation](https://docs.github.com/en/actions)
2. [Mastering CI/CD Best Practices](https://www.linkedin.com/pulse/mastering-cicd-best-practices-github-actions-tatiana-sava)
3. [A Practical Guide to Building Smarter GitHub Workflows](https://hackernoon.com/a-practical-guide-to-building-smarter-github-workflows)
4. [GitHub Actions Marketplace](https://github.com/marketplace?type=actions)
5. [GitHub Actions Security Best Practices](https://blog.gitguardian.com/github-actions-security-best-practices/)
6. [GitHub Actions Security Hardening](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
7. [OpenSSF Security Best Practices](https://openssf.org/blog/2023/09/14/openssf-releases-source-code-management-best-practices-guide/)
8. [Wiz Security Guide for GitHub Actions](https://www.wiz.io/blog/github-actions-security-guide)
9. [OpenSSF GitHub Security Best Practices](https://best.openssf.org/SCM-BestPractices/github/)

---

## References

1. [GitHub Actions Documentation](https://docs.github.com/en/actions)
2. [Mastering CI/CD Best Practices](https://www.linkedin.com/pulse/mastering-cicd-best-practices-github-actions-tatiana-sava)
3. [GitHub Actions Marketplace](https://github.com/marketplace?type=actions)
4. [GitHub Actions Security Best Practices](https://blog.gitguardian.com/github-actions-security-best-practices/)
5. [GitHub Actions Security Hardening](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
6. [OpenSSF Security Best Practices](https://openssf.org/blog/2023/09/14/openssf-releases-source-code-management-best-practices-guide/)
7. [Wiz Security Guide for GitHub Actions](https://www.wiz.io/blog/github-actions-security-guide)
8. [OpenSSF GitHub Security Best Practices](https://best.openssf.org/SCM-BestPractices/github/)

---

To share this document as a GitHub Gist for easy distribution:

1. Visit [gist.github.com](https://gist.github.com/)
2. Copy the entire content of this document
3. Paste into the gist editor
4. Set the filename to `github-actions-best-practices.md`
5. Add an optional description
6. Click "Create public gist" (or private if you prefer)
7. Share the resulting URL with your team

A gist can be updated later and will maintain version history of all changes.

> **Pro tip:** You can embed gists in other Markdown documents or web pages using
> GitHub's embed script, which is automatically generated for each gist.

---
