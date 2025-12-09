# Software Engineering Principles I Ignored for Too Long

> _Inspired by Abhay Parashar, The Pythoneers, May 2025, Medium_

This document summarizes key software engineering principles that are often
overlooked, with practical Go examples and actionable advice.

---

## 1. DRY (Don't Repeat Yourself)

**Principle:**  
Avoid duplicating code or logic. Duplication increases maintenance cost and risk of inconsistencies.

**Bad Example:**

```go
func getUserName(user User) string {
    return user.FirstName + " " + user.LastName
}

func getAuthorName(author Author) string {
    return author.FirstName + " " + author.LastName
}
```

**Good Example:**

```go
type Nameable interface {
    GetFirstName() string
    GetLastName() string
}

func getFullName(n Nameable) string {
    return n.GetFirstName() + " " + n.GetLastName()
}
```

**Tip:**  
Extract common logic into functions or interfaces. Use composition over inheritance.

---

## 2. YAGNI (You Aren't Gonna Need It)

**Principle:**  
Don't implement features until they are actually needed. Premature generalization leads to wasted effort and complexity.

**Bad Example:**

```go
// Adding a feature "just in case"
func calculate(a, b int, operation string) int {
    if operation == "add" {
        return a + b
    }
    // Subtraction not needed yet, but added anyway
    if operation == "subtract" {
        return a - b
    }
    return 0
}
```

**Good Example:**

```go
// Only implement what's needed now
func add(a, b int) int {
    return a + b
}
```

**Tip:**  
Focus on current requirements. Refactor when new needs arise.

---

## 3. KISS (Keep It Simple, Stupid)

**Principle:**  
Prefer simple, clear solutions over clever or complex ones. Simplicity makes code easier to read, test, and maintain.

**Bad Example:**

```go
func isEven(n int) bool {
    if n%2 == 0 {
        return true
    }
    return false
}
```

**Good Example:**

```go
func isEven(n int) bool {
    return n%2 == 0
}
```

**Tip:**  
If you can't explain your code simply, it's probably too complex.

---

## 4. Single Responsibility Principle (SRP)

**Principle:**  
A module, class, or function should have only one reason to change. Each responsibility should be a separate concern.

**Bad Example:**

```go
type User struct {
    Name string
}

func (u *User) Save() error {
    // Save user to DB
}

func (u *User) Validate() bool {
    // Validate user fields
}
```

**Good Example:**

```go
type User struct {
    Name string
}

type UserRepository struct{}

func (r *UserRepository) Save(u *User) error {
    // Save user to DB
}

func (u *User) Validate() bool {
    // Validate user fields
}
```

**Tip:**  
Separate validation, persistence, and business logic.

---

## 5. Readability Over Cleverness

**Principle:**  
Write code for humans first, computers second. Favor clarity over brevity or clever tricks.

**Example:**

```go
// Less readable
res := make([]int, 0)
for i := 0; i < 10; i++ { if i%2==0 { res = append(res, i) } }

// More readable
var evens []int
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        evens = append(evens, i)
    }
}
```

**Tip:**  
Use meaningful names and consistent formatting.

---

## 6. Avoid Premature Optimization

**Principle:**  
First make it work, then make it right, then make it fast. Optimize only when necessary and with evidence.

**Example:**

```go
// Don't micro-optimize before profiling
```

**Tip:**  
Profile your code before optimizing. Maintain readability.

---

## 7. Test Early and Often

**Principle:**  
Write tests as you develop. Tests help catch bugs early and make refactoring safer.

**Example:**

```go
func Add(a, b int) int {
    return a + b
}

// In add_test.go
func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2,3) = %d; want %d", got, want)
    }
}
```

**Tip:**  
Aim for high test coverage, but prioritize meaningful tests.

---

## 8. Favor Composition Over Inheritance

**Principle:**  
Compose behavior using interfaces and embedding rather than deep inheritance trees.

**Example:**

```go
type Logger interface {
    Log(msg string)
}

type Service struct {
    Logger
}
```

**Tip:**  
Go encourages composition via interfaces and struct embedding.

---

## 9. Document Your Code

**Principle:**  
Good documentation helps others (and your future self) understand and use your code.

**Example:**

```go
// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}
```

**Tip:**  
Use Go doc comments and keep documentation up to date.

---

### 10. Embrace Code Reviews

**Principle:**  
Peer reviews catch issues you might miss and improve code quality.

**Tip:**  
Be open to feedback and review others' code constructively.

---

## Further Reading

- [Software Engineering Principles I Ignored for Too Long (Medium)](https://medium.com/@abhayparashar/software-engineering-principles-i-ignored-for-too-long-xxxxxx)
- [Go Proverbs](https://go-proverbs.github.io/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Clean Code by Robert C. Martin](https://www.oreilly.com/library/view/clean-code/9780136083238/)

---
