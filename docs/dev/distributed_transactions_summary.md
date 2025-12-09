# Distributed Transactions in Go: Summary

## Key Takeaways

### 1. Distributed Transactions Are (Almost Always) an Anti-Pattern

- Avoid transactions that span multiple services. They are hard to test, debug,
  and maintain, and they tightly couple your system.
- If you feel the need for distributed transactions, your service boundaries are
  likely wrong. Consider merging services or using a modular monolith.

### 2. The Distributed Monolith Problem

- Example: Deducting user points in one service and applying a discount in another
  can lead to inconsistencies if the second step fails.
- Synchronous HTTP calls between services do not guarantee atomicity.

### 3. Eventual Consistency as a Solution

- Instead of trying to make everything consistent immediately, use events and message
  queues (Pub/Sub) to achieve eventual consistency.
- Publish an event (e.g., `PointsUsedForDiscount`) after updating local state.
  Other services react to the event asynchronously.
- Most of the time, the system will be consistent within milliseconds. If a service
  is down, the event will be retried until processed.

### 4. Implementation Tips

- Use a message broker (e.g., Redis, Kafka, NATS) and a Go library like [Watermill](https://watermill.io/)
  for event-driven communication.
- Replace direct service calls with event publishing and handling.
- Handlers in other services subscribe to relevant events and process them as needed.

### 5. The Outbox Pattern

- To avoid losing events due to network issues, use the Outbox Pattern:
  - Store both the data change and the event in the same database transaction.
  - A separate process (forwarder) reads events from the database and publishes them to the message broker.
- This ensures that events are not lost even if the broker is temporarily unavailable.

### 6. Event Design and Coupling

- Events are contracts between services. Design them to state facts about your
  domain, not to expose internal workflows or intentions.
- Poorly designed events can create tight coupling between services.

### 7. Testing and Monitoring

- Use real message brokers in tests (e.g., via Docker).
- Test the system’s behavior via public APIs, not just event handlers.
- Use tools like `assert.Eventually` to wait for eventual consistency in tests.
- Monitor the queue for unprocessed messages—long delays usually indicate a problem.

### 8. When to Use Distributed Transactions

- Only use true distributed transactions or the saga pattern if there is absolutely
  no alternative and strong consistency is a must.
- In most business cases, eventual consistency is sufficient and much simpler to
  implement and maintain.

---

**Summary for Future Projects:**

- Prefer eventual consistency and event-driven architecture over distributed transactions.
- Use the Outbox Pattern to ensure reliable event delivery.
- Design events as domain facts, not as commands or workflow steps.
- Test and monitor your event-driven flows carefully.
- Rethink your service boundaries if you find yourself needing distributed transactions.

## Reference

- [Distributed Transactions in Go: Read Before You Try](https://threedots.tech/post/distributed-transactions-in-go/?utm_source=newsletter&utm_medium=email&utm_term=2025-04-29&utm_campaign=Distributed+Transactions+in+Go+Read+Before+You+Try)
