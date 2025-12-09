---
applyTo: '**/*.go'
description: Go-specific problem resolution and debugging protocols with zero-tolerance error handling
---

# Go Problem Resolution Instructions
**Version: 1.0.0 | Last Updated: 2025-11-18**

> **CRITICAL NOTICE**: These Go-specific instructions are NON-NEGOTIABLE for debugging and problem resolution. They MUST be followed EXACTLY as specified to prevent circular debugging and ensure systematic problem resolution in Go codebases.

---

## 🚨 MANDATORY PRE-FLIGHT CHECK FOR GO PROJECTS
**EXECUTE BEFORE ANY GO CODE CHANGE - NO EXCEPTIONS**

```yaml
BEFORE_ANY_GO_CODE_CHANGE:
  1. Run `get_errors` tool → Document ALL compilation/runtime errors
  2. Run `go mod tidy` → Verify dependencies are clean
  3. Run `go vet ./...` → Check for Go-specific issues
  4. Run `golangci-lint run` → Comprehensive static analysis
  5. Check `GO_ERROR_TRACKING.md` → Verify not repeating previous fixes
  6. Scan for Go anti-patterns → `panic()`, `recover()`, ignored errors
  7. Validate imports → Check for circular dependencies
  8. Test current build → `go build ./...`
```

---

## 🛑 GO-SPECIFIC ANTI-WHACK-A-MOLE PROTOCOL
**PATTERN DETECTION & PREVENTION FOR GO CODEBASES**

### GO ERROR TRACKING TEMPLATE
```yaml
MANDATORY_GO_TRACKING:
  file: GO_ERROR_TRACKING.md
  format: |
    ## Go Error Entry [TIMESTAMP]
    ### Symptom
    - Error Message: [EXACT GO ERROR]
    - Package: [PACKAGE PATH]
    - File:Line: [FILE:LINE]
    - Go Version: [GO VERSION]
    - Build Context: [GOOS/GOARCH]
    
    ### Go-Specific Context
    - Module: [GO MODULE PATH]
    - Dependencies: [RELATED GO MODULES]
    - Goroutines Involved: [IF CONCURRENCY ISSUE]
    - Interface/Struct: [TYPE INFORMATION]
    
    ### Root Cause Analysis
    - Surface Issue: [WHAT APPEARS BROKEN]
    - Actual Go Cause: [UNDERLYING GO PROBLEM]
    - Type Safety Issue: [IF TYPE-RELATED]
    - Memory/Goroutine Leak: [IF RESOURCE ISSUE]
    
    ### Fix Applied
    - Changed Code: [EXACT GO CODE CHANGES]
    - Types Modified: [INTERFACE/STRUCT CHANGES]
    - Dependencies Updated: [GO.MOD CHANGES]
    - Tests Added: [NEW TEST CASES]
    
    ### Go-Specific Validation
    - `go build`: [PASS/FAIL]
    - `go test`: [PASS/FAIL WITH COVERAGE]
    - `go vet`: [PASS/FAIL WITH DETAILS]
    - `golangci-lint`: [PASS/FAIL WITH ISSUES]
    - Race Detector: [`go test -race` RESULTS]
    - New Errors: [ANY NEW GO ISSUES]
```

### GO PATTERN RECOGNITION RULES
1. **STOP if seeing same import cycle 2+ times**
   - Review module structure immediately
   - Identify architectural issue, not syntax
   - Escalate: "GO CIRCULAR IMPORT DETECTED in [packages]"

2. **STOP if panic/recover pattern appears**
   - Go favors explicit error handling
   - Review error propagation strategy
   - Document: "GO PANIC PATTERN - Review error handling"

3. **STOP if interface{} usage increases**
   - Type safety regression indicator
   - Review generic types usage (Go 1.18+)
   - Consider: "Type safety compromise detected"

---

## 🔧 GO-SPECIFIC CONFIGURATION PROTOCOL
**MODULE AND DEPENDENCY MANAGEMENT**

### GO MODULE STATE TRACKING
```yaml
GO_MODULE_CONFIG:
  primary_file: GO_MODULE_STATE.json
  structure:
    module:
      path: [GO MODULE PATH]
      version: [MODULE VERSION]
      go_version: [REQUIRED GO VERSION]
      
    dependencies:
      direct:
        - module: version
      indirect:
        - module: version
      replaced:
        - original: replacement
        
    build_info:
      goos: [TARGET OS]
      goarch: [TARGET ARCH]
      cgo_enabled: [true/false]
      tags: [BUILD TAGS]
      
    tools:
      golangci_lint: [VERSION]
      goimports: [VERSION]
      mockgen: [VERSION]
      
    last_tidy: [TIMESTAMP]
    last_vet: [TIMESTAMP]
    last_lint: [TIMESTAMP]
```

### GO DEPENDENCY CHANGE PROTOCOL
```yaml
BEFORE_ANY_GO_DEPENDENCY_CHANGE:
  1. Run `go mod download` → Verify current state
  2. Backup go.mod/go.sum → `cp go.{mod,sum} backup/`
  3. Document intent → "Adding [MODULE] for [PURPOSE]"
  4. Add/update dependency → `go get [MODULE]@[VERSION]`
  5. Run `go mod tidy` → Clean up dependencies
  6. Test integration → `go test ./...`
  7. Check for vulnerabilities → `go list -json -deps | nancy sleuth`
  8. Update GO_MODULE_STATE.json
  9. Commit or rollback → Within 5 minutes
```

---

## 📋 GO ERROR CLASSIFICATION SYSTEM
**MANDATORY GO ERROR CATEGORIZATION**

```yaml
GO_ERROR_CATEGORIES:
  COMPILATION_ERROR:
    code_prefix: "COMP"
    examples: [syntax, type mismatch, undefined]
    action: Fix code syntax/types
    tools: [go build, go vet]
    
  RUNTIME_ERROR:
    code_prefix: "RUN"
    examples: [panic, nil pointer, index out of range]
    action: Add nil checks, bounds checking
    tools: [go test -race, delve debugger]
    
  CONCURRENCY_ERROR:
    code_prefix: "CONC"
    examples: [race condition, deadlock, channel]
    action: Review goroutine patterns, add sync
    tools: [go test -race, go tool trace]
    
  IMPORT_ERROR:
    code_prefix: "IMP"
    examples: [circular import, missing package]
    action: Restructure packages/modules
    tools: [go mod graph, go list]
    
  TYPE_ERROR:
    code_prefix: "TYPE"
    examples: [interface compliance, assertion]
    action: Fix type definitions/usage
    tools: [go vet, type checker]
    
  DEPENDENCY_ERROR:
    code_prefix: "DEP"
    examples: [version conflict, missing module]
    action: Update go.mod, resolve versions
    tools: [go mod tidy, go mod why]
```

---

## ⚡ GO INCREMENTAL VALIDATION PROTOCOL
**TEST AFTER EVERY GO CHANGE - NO BATCHING**

```yaml
GO_VALIDATION_SEQUENCE:
  after_each_go_change:
    1. Save .go file
    2. Run `go build [package]` → Must pass
    3. Run `go vet [package]` → Must pass
    4. Run `go test [package]` → Must pass
    5. Run `go test -race [package]` → Check concurrency
    6. Run `golangci-lint run [package]` → Check style/bugs
    7. Check test coverage → `go test -cover`
    8. Update GO_ERROR_TRACKING.md → Before next change
    
  go_specific_checks:
    - gofmt compliance: Auto-format on save
    - import organization: Use goimports
    - unused variables: go vet catches these
    - nil pointer safety: Review all pointer usage
    - error handling: Check all error returns
    
  batch_threshold: 1  # NEVER batch Go changes
  rollback_on_fail: true
  continue_on_warning: false
```

---

## 🎯 GO SINGLE RESPONSIBILITY FIXES
**ONE GO PROBLEM = ONE FIX**

```yaml
GO_FIX_ISOLATION_RULES:
  - Fix ONE compilation error at a time
  - Fix ONE type error at a time  
  - Fix ONE import issue at a time
  - Test after EACH Go fix
  - Document BEFORE moving to next
  - NEVER combine interface and implementation fixes
  - NEVER fix types while debugging logic
  
GO_PROHIBITED_ACTIONS:
  - Fixing multiple compilation errors in one commit
  - Refactoring while fixing bugs
  - Changing interfaces during error fixes
  - Updating Go modules during bug fixes
  - Ignoring `go vet` warnings
  - Committing with `//TODO: fix this` comments
```

---

## 🔄 GO CONCURRENCY DEBUGGING PROTOCOL
**GOROUTINE AND CHANNEL ISSUE RESOLUTION**

```yaml
CONCURRENCY_DEBUG_SEQUENCE:
  detection_tools:
    - `go test -race`: Detect race conditions
    - `go tool trace`: Analyze goroutine behavior  
    - `GODEBUG=schedtrace=1000`: Runtime scheduling
    - Deadlock detector: Third-party tools
    
  common_patterns:
    goroutine_leak:
      symptoms: [increasing memory, hanging]
      check: `go tool pprof goroutine`
      fix: Ensure goroutine termination
      
    race_condition:
      symptoms: [inconsistent results, crashes]
      check: `go test -race -count=100`
      fix: Add proper synchronization
      
    channel_deadlock:
      symptoms: [hanging, all goroutines asleep]
      check: Stack trace analysis
      fix: Review channel usage patterns
      
    select_starvation:
      symptoms: [some cases never execute]
      check: Add default case analysis
      fix: Reorder select cases, add timeouts

CONCURRENCY_RECOVERY_PROTOCOL:
  1. Isolate the concurrent code
  2. Add extensive logging with goroutine IDs
  3. Reduce concurrency to 1 goroutine
  4. Verify sequential correctness
  5. Gradually increase concurrency
  6. Add comprehensive tests with race detector
```

---

## 📊 GO PROGRESS TRACKING DASHBOARD
**MAINTAIN IN PROJECT ROOT**

```markdown
# GO_DEBUG_DASHBOARD.md

## Current Go Build Status
- [ ] `go build ./...` passing
- [ ] `go test ./...` passing
- [ ] `go vet ./...` clean
- [ ] `golangci-lint run` clean
- [ ] `go test -race ./...` clean
- [ ] All imports organized
- [ ] All errors handled explicitly
- [ ] No panic/recover usage

## Go Module Status  
- [ ] `go mod tidy` clean
- [ ] No replace directives in production
- [ ] All dependencies up to date
- [ ] No known vulnerabilities
- [ ] Build reproducible

## Error Count by Type
- Compilation: 0
- Runtime: 0  
- Concurrency: 0
- Type: 0
- Import: 0
- Dependency: 0

## Go Fix History
| Timestamp | Package | Error Type | Fix Applied | Tests Added |
|-----------|---------|------------|-------------|-------------|
| [TIME] | [PKG] | [TYPE] | [FIX] | [TESTS] |

## Go Anti-Patterns Detected
| Pattern | Count | Package | Resolution |
|---------|-------|---------|------------|
| Ignored errors | 0 | - | - |
| panic/recover | 0 | - | - |
| interface{} overuse | 0 | - | - |
| Goroutine leaks | 0 | - | - |
```

---

## 💀 GO DEATH SPIRAL PREVENTION
**EMERGENCY PROTOCOLS FOR GO PROJECTS**

```yaml
GO_EMERGENCY_TRIGGERS:
  - Same import cycle appearing 3+ times
  - Goroutine count increasing without bound
  - Test coverage dropping below 70%
  - Build time > 5x normal
  - Memory usage growing during tests
  - `go mod tidy` failing repeatedly
  
GO_EMERGENCY_ACTIONS:
  1. FULL STOP - No more Go code changes
  2. Run `git stash push -m "Emergency Go state"`
  3. Restore last working go.mod/go.sum
  4. Run full Go toolchain: build/test/vet/lint
  5. Document architectural failure in GO_ERROR_TRACKING.md
  6. Consider package restructuring
  7. Review Clean Architecture boundaries
  8. Seek architectural guidance
```

---

## 🔍 GO-SPECIFIC DEBUGGING TOOLS & TECHNIQUES

### Essential Go Debugging Commands
```bash
# Error Detection & Analysis
go build ./...                    # Check compilation
go vet ./...                     # Static analysis
golangci-lint run                # Comprehensive linting
go test -v ./...                 # Run all tests
go test -race ./...              # Race condition detection
go test -cover ./...             # Test coverage

# Dependency Management
go mod tidy                      # Clean dependencies
go mod why [module]              # Why is module needed
go mod graph                     # Dependency graph
go list -m all                   # All modules
go list -json [package]          # Package details

# Performance & Profiling
go test -bench=.                 # Benchmarks
go tool pprof [profile]          # CPU/memory profiling
go tool trace [trace]            # Execution tracing
GODEBUG=gctrace=1 [program]      # GC tracing

# Type & Interface Analysis
go doc [package]                 # Package documentation
go list -json -deps [package]    # Dependencies
go tool objdump [binary]         # Assembly analysis

# Module Security
go list -json -deps | nancy sleuth  # Vulnerability scanning
go mod download -x               # Verbose download info
```

### Go Debugging Workflow
```yaml
SYSTEMATIC_GO_DEBUGGING:
  step_1_isolation:
    - Create minimal reproduction case
    - Remove all unnecessary code
    - Use single goroutine if concurrent
    - Add extensive logging with slog
    
  step_2_analysis:
    - Check all error handling paths
    - Verify nil pointer safety
    - Confirm interface compliance
    - Review channel operations
    
  step_3_validation:
    - Write failing test first
    - Fix until test passes
    - Run full test suite
    - Check with race detector
    
  step_4_integration:
    - Gradually add complexity back
    - Test at each integration step
    - Monitor memory/goroutine usage
    - Update documentation
```

---

## 🎓 GO LEARNING PROTOCOL
**CONTINUOUS IMPROVEMENT FOR GO DEVELOPMENT**

After EVERY Go debugging session:
1. Review GO_ERROR_TRACKING.md for patterns
2. Update these instructions with new Go patterns
3. Add Go-specific test cases for bugs encountered
4. Document Go module lessons learned
5. Share Go architectural findings with team
6. Review Go best practices alignment
7. Update package documentation

---

## 📝 GO COMMUNICATION PROTOCOL
**STATUS REPORTING FORMAT FOR GO ISSUES**

```yaml
GO_STATUS_UPDATE_TEMPLATE: |
  ## Go Status Update [TIMESTAMP]
  
  ### Go Build State
  - Compilation Errors: [COUNT]
  - Test Failures: [COUNT] 
  - Race Conditions: [COUNT]
  - Lint Issues: [COUNT]
  - Coverage: [PERCENTAGE]
  
  ### Go Module State
  - Dependencies Clean: [YES/NO]
  - Vulnerabilities: [COUNT]
  - Replace Directives: [COUNT]
  
  ### Current Go Issue
  - Package: [PACKAGE PATH]
  - Error Type: [CATEGORY]
  - Approach: [DEBUGGING STRATEGY]
  - Risk: [POTENTIAL GO IMPACTS]
  
  ### Go-Specific Blockers
  - [LIST GO BLOCKERS - imports, types, etc.]
  
  ### Estimate
  - [REALISTIC GO FIX ESTIMATE]
```

---

## 🚀 GO COMPLETION CRITERIA
**DEFINITION OF "DONE" FOR GO TASKS**

```yaml
GO_TASK_COMPLETE_ONLY_WHEN:
  - `go build ./...`: PASS
  - `go test ./...`: PASS  
  - `go test -race ./...`: PASS
  - `go vet ./...`: CLEAN
  - `golangci-lint run`: CLEAN
  - `go mod tidy`: CLEAN
  - Test coverage >= 80%: TRUE
  - All errors handled explicitly: TRUE
  - No panic/recover patterns: TRUE
  - No goroutine leaks: TRUE
  - Documentation updated: TRUE
  - GO_ERROR_TRACKING.md current: TRUE
  - GO_MODULE_STATE.json valid: TRUE
  - No anti-patterns detected: TRUE
```

---

## ⚠️ GO-SPECIFIC FORBIDDEN ACTIONS
**NEVER DO THESE IN GO CODE - NO EXCEPTIONS**

1. **NEVER** ignore error returns (`_, err := func(); // ignore err`)
2. **NEVER** use panic/recover for normal error handling
3. **NEVER** commit code that doesn't pass `go vet`
4. **NEVER** use `interface{}` without strong justification
5. **NEVER** create goroutines without termination strategy
6. **NEVER** modify shared state without proper synchronization
7. **NEVER** ignore race condition warnings
8. **NEVER** commit with `//TODO` or `//HACK` comments
9. **NEVER** use `init()` functions without documentation
10. **NEVER** bypass type safety with unsafe package
11. **NEVER** commit failing tests "to save time"
12. **NEVER** use global variables for application state

---

## 📚 GO REFERENCE TOOLS & COMMANDS

```bash
# Build & Test
go build ./...                   # Build all packages
go test -v -race ./...          # Test with race detection  
go clean -cache                 # Clean build cache
go env                          # Go environment

# Static Analysis
go vet ./...                    # Go's built-in checker
golangci-lint run               # Comprehensive linting
go fmt ./...                    # Format code
goimports -w .                  # Organize imports

# Dependencies  
go mod init [module]            # Initialize module
go mod tidy                     # Clean dependencies
go get [package]@[version]      # Add/update dependency
go mod vendor                   # Vendor dependencies

# Debugging & Profiling
dlv debug                       # Delve debugger
go tool pprof                   # Profiling analysis
go tool trace                   # Execution tracing
go tool objdump                 # Disassembly

# Documentation
go doc [package]                # Package docs
go doc -all [package]           # Full package docs
godoc -http=:6060              # Local documentation server

# Tracking Files
touch GO_ERROR_TRACKING.md     # Initialize Go error log
touch GO_MODULE_STATE.json     # Initialize module tracking  
touch GO_DEBUG_DASHBOARD.md    # Initialize Go dashboard

# Emergency Recovery
git stash push -m "Go emergency: [REASON]"  # Emergency save
go clean -modcache                          # Clean module cache
go mod download                             # Re-download modules
```

---

## 🔴 GO CRITICAL REMINDERS

> **"Write Go code the Go way, or don't write Go at all"**

1. **Explicit Error Handling**: Every error must be handled - no exceptions
2. **Goroutine Lifecycle**: Every goroutine needs a termination strategy
3. **Interface Compliance**: Types should satisfy interfaces naturally
4. **Concurrency Safety**: Shared state requires synchronization
5. **Module Hygiene**: Keep dependencies clean and minimal
6. **Type Safety First**: Avoid interface{} unless absolutely necessary
7. **Test Everything**: Especially concurrent code with race detector

---

## 📌 GO QUICK REFERENCE CARD

```asciidoc
┌─────────────────────────────────────┐
│  BEFORE EVERY GO ACTION:            │
│  1. go build ./...                  │
│  2. go test ./...                   │
│  3. go vet ./...                    │
│  4. Read GO_ERROR_TRACKING.md       │
│                                     │
│  AFTER EVERY GO FIX:                │
│  1. go test -race                   │
│  2. golangci-lint run               │
│  3. Update tracking files           │
│  4. Check no new goroutine leaks    │
│                                     │
│  IF GO STUCK:                       │
│  1. Stop immediately                │
│  2. Check for import cycles         │
│  3. Review last 5 Go changes        │
│  4. Consider architectural issue    │
│  5. Use minimal reproduction case   │
└─────────────────────────────────────┘
```ass

---

**END OF GO INSTRUCTIONS - NO MODIFICATIONS WITHOUT VERSION INCREMENT**

*These Go-specific instructions eliminate circular debugging in Go codebases by leveraging Go's excellent tooling ecosystem. Follow them exactly for systematic, efficient Go problem resolution that respects Go idioms and best practices.*