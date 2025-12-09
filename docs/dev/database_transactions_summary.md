# Database Transactions in Go: Summary

## Why Transactions Matter

- Transactions ensure that a set of related database operations are atomic: all succeed or none do.
- Without transactions, partial updates can leave your data inconsistent, especially in the face of errors or concurrent
  requests.

## SQL Transactions 101

- Use `BEGIN` and `COMMIT` to wrap related queries in a transaction.
- Use `SELECT ... FOR UPDATE` to lock rows and prevent race conditions in concurrent updates.
- Be aware: `FOR UPDATE` can cause performance bottlenecks under high contention.
- Test your code under load to catch issues early.

## Layered Architecture and Transactions

- Keep business logic and database code separate (e.g., using the Repository pattern).
- The challenge: how to handle transactions cleanly across layers.

## Anti-Patterns

- **Skipping transactions**: Never rely on luck; always use transactions for related updates.
- **Mixing transactions with logic**: Avoid passing transaction objects (`*sql.Tx`) through your business logic. It
    complicates code and testing.
- **One repository per table**: Instead, design repositories around aggregates—groups of data that must be consistent together.

## Recommended Patterns

### 1. Transactions Inside the Repository (Aggregate Repositories)

- Keep all data that must be consistent in a single aggregate and repository.
- The repository handles the transaction, so the application logic stays clean.
- Downside: Business logic may leak into the repository, and interfaces can grow large.

### 2. The UpdateFn Pattern (Preferred)

- The repository exposes an `UpdateByID` method that takes a closure (function) with the loaded aggregate.
- The closure contains the business logic; the repository manages the transaction and persistence.
- Keeps logic and data access cleanly separated.

```go
// Handler
err := userRepository.UpdateByID(ctx, userID, func(user *User) (bool, error) {
    err := user.UsePointsAsDiscount(points)
    if err != nil {
        return false, err
    }
    return true, nil
})
```

### 3. The Transaction Provider (For Edge Cases)

- For rare cases where multiple repositories must share a transaction (e.g., updating an audit log and a user), use a
  transaction provider.
- The provider manages the transaction and injects repositories into a closure.
- Use sparingly; can get out of hand if overused.

## General Advice

- Design repositories around aggregates, not tables.
- Keep business logic out of repositories when possible.
- Use the UpdateFn pattern for most transactional updates.
- Only use transaction providers for technical cross-cutting concerns.
- Always test under realistic concurrency and load.

## Related Topics

- For distributed transactions (across services), see the summary in `distributed_transactions_summary.md`.

---

This summary is based on the article: [Database Transactions in Go with Layered Architecture](https://threedots.tech/post/database-transactions-in-go/)
