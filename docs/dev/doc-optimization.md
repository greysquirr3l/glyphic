# Documentation & Prompt Optimization Strategy

## Current State Analysis

### Problems Identified

1. **Duplication**: Similar content exists in `docs/dev/`, `.github/prompts/`, and `.github/copilot-instructions/`
2. **Inconsistent Access**: Some docs are manual-only, others are auto-loaded by Copilot
3. **Mixed Purposes**: Technical references mixed with actionable prompts
4. **Poor Discoverability**: No clear system for when to use which document
5. **Context Overload**: Too much auto-loaded content may dilute Copilot's focus

## Recommended Organization Structure

```text
.github/copilot/                    # Auto-loaded by GitHub Copilot
├── core-architecture.md           # Always-loaded: DDD+CQRS+Clean Architecture
├── security-first.md              # Always-loaded: Security principles
├── go-standards.md                # Always-loaded: Go coding standards
└── branch-naming.md               # Always-loaded: Git workflow standards

.github/prompts/                    # On-demand prompts (tagged)
├── validate/
│   ├── architecture-review.md     # #architecture #validation
│   ├── security-audit.md          # #security #audit
│   ├── cqrs-compliance.md         # #cqrs #validation
│   └── test-coverage.md           # #testing #coverage
├── generate/
│   ├── domain-entity.md           # #domain #entity #generation
│   ├── repository.md              # #repository #generation
│   ├── handler.md                 # #handler #generation
│   └── value-object.md            # #value-object #generation
├── analyze/
│   ├── performance.md             # #performance #analysis
│   ├── dependencies.md            # #dependencies #analysis
│   └── complexity.md              # #complexity #analysis
└── tools/
    ├── readme-generator.md        # #documentation #readme
    ├── changelog-generator.md     # #changelog #release
    └── migration-helper.md        # #database #migration

docs/dev/                          # Manual reference docs
├── architecture/
│   ├── ddd-cqrs-guide.md         # Comprehensive guide
│   ├── clean-architecture.md     # Layer definitions
│   └── patterns/
├── implementation/
│   ├── go-concurrency.md
│   ├── database-patterns.md
│   └── testing-strategies.md
├── infrastructure/
│   ├── deployment-guide.md
│   ├── monitoring-setup.md
│   └── security-hardening.md
└── tools/
    ├── development-setup.md
    ├── ci-cd-configuration.md
    └── debugging-guide.md
```

## Implementation Plan

### Phase 1: Core Copilot Files (Always Loaded)

Create essential files that should always influence Copilot's behavior:

1. **`.github/copilot/core-architecture.md`**
   - Extract core DDD+CQRS+Clean Architecture principles
   - Keep to 1-2 pages max
   - Focus on "what to do" not "how to do it"

2. **`.github/copilot/security-first.md`**
   - Extract security-first repository design principles
   - Authorization patterns
   - Input validation requirements

3. **`.github/copilot/go-standards.md`**
   - Extract Go coding standards and tooling versions
   - Import patterns and project structure
   - Testing requirements

### Phase 2: Tagged On-Demand Prompts

Transform existing content into actionable, tagged prompts:

#### Validation Prompts (`#validation`)

- `architecture-review.md` - Check Clean Architecture compliance
- `security-audit.md` - Audit repository security patterns
- `cqrs-compliance.md` - Validate command/query separation

#### Generation Prompts (`#generation`)

- `domain-entity.md` - Generate DDD entities with proper invariants
- `repository.md` - Generate secure repository implementations
- `handler.md` - Generate CQRS handlers

#### Analysis Prompts (`#analysis`)

- `performance.md` - Analyze performance bottlenecks
- `complexity.md` - Identify overly complex code

### Phase 3: Manual Reference Reorganization

Move detailed technical documentation to well-organized reference:

1. **Group by Purpose**: Architecture, Implementation, Infrastructure, Tools
2. **Cross-Reference**: Link related concepts
3. **Template Structure**: Standardize document format
4. **Index Creation**: Master index with tags and use cases

## Tagging System

### Prompt Tags

```markdown
---
tags: [architecture, validation, ddd, cqrs]
mode: agent
purpose: Validate Clean Architecture compliance in Go code
trigger_keywords: [validate, architecture, clean, layers, dependencies]
---
```

### Usage Examples

```bash
# In Copilot Chat:
#architecture #validation - check my repository pattern
#security #audit - review this handler for security issues
#domain #generation - create a User aggregate
#performance #analysis - optimize this database query
```

## File Template Standards

### Copilot Core Files (Always Loaded)

```markdown
# [Topic] Core Principles

## Key Rules
1. [Rule 1 - actionable]
2. [Rule 2 - actionable]

## Common Patterns
- [Pattern with example]

## Anti-Patterns to Avoid
- [Anti-pattern with explanation]

Keep to 2 pages maximum for consistent loading
```

### On-Demand Prompts

```markdown
---
tags: [tag1, tag2, tag3]
mode: agent
purpose: [Clear purpose statement]
trigger_keywords: [word1, word2, word3]
---

# [Prompt Title]

[Clear instruction for AI]

## Context
[What the AI needs to know]

## Expected Output
[Format and structure expected]

## Examples
[Concrete examples]
```

### Reference Documentation

```markdown
# [Topic] Reference Guide

## Overview
[Brief description and scope]

## When to Use This Guide
[Specific use cases]

## Quick Reference
[Key points for scanning]

## Detailed Implementation
[Comprehensive details]

## Related Guides
[Cross-references]

## Examples
[Real-world examples]
```

## Migration Strategy

### Immediate Actions (Week 1)

1. Create core Copilot files from existing content
2. Audit current `.github/prompts/` for overlap
3. Tag existing prompts with new system

### Short-term (Month 1)

1. Reorganize `docs/dev/` by purpose
2. Create template-driven prompt files
3. Update README files with new structure

### Long-term (Quarter 1)

1. Implement prompt analytics to track usage
2. Create automated tools for prompt discovery
3. Integrate with development workflow

## Benefits of New Structure

### For GitHub Copilot

- **Focused Core Context**: Essential principles always loaded
- **Reduced Noise**: On-demand loading prevents information overload
- **Better Targeting**: Tagged prompts match specific needs
- **Consistent Quality**: Template-driven prompt structure

### For Developers

- **Clear Discovery**: Tags and structure make finding relevant prompts easy
- **Purposeful Reference**: Manual docs organized by actual use cases
- **Workflow Integration**: Prompts designed for specific development phases
- **Reduced Duplication**: Single source of truth per concept

### For AI Assistants

- **Context Clarity**: Clear separation between always-on rules and situational prompts
- **Task Specificity**: Prompts designed for specific AI capabilities
- **Output Consistency**: Standardized templates produce predictable results
- **Scalability**: System grows without degrading performance

## Measurement & Success Criteria

### Usage Metrics

- Prompt activation frequency by tag
- Developer adoption of new structure
- Reduction in documentation maintenance overhead

### Quality Metrics

- Code quality improvements in generated code
- Reduction in architectural violations
- Security issue prevention rate

### Developer Experience

- Time to find relevant information
- Satisfaction with AI-generated code
- Onboarding speed for new team members

---

## GitHub Copilot Technical Notes

### Key Limitation

- Custom prompts can **only reference files within** `.github/copilot/` directory
- References to files outside this folder are ignored/blocked for security reasons

### On-Demand Inclusion

Referenced files are only included when:

- You invoke a specific prompt that references them
- You manually reference them in chat (e.g., `@filename.md`)
- Files are **not** automatically included in every Copilot request

### Usage Examples

```bash
# In Copilot Chat:
@architecture-review.md - check my repository pattern
@security-audit.md - review this handler for security issues
#domain #generation - create a User aggregate
```
