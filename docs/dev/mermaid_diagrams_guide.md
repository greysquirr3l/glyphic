# 📊 GitHub-Compatible Mermaid Diagrams Guide

> A comprehensive guide for creating effective Mermaid diagrams in GitHub markdown
> files, covering syntax, best practices, and advanced patterns for documenting
> Go applications.

<!-- REF: https://mermaid.js.org/syntax/examples.html -->
<!-- REF: https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/creating-diagrams -->
<!-- REF: https://github.blog/developer-skills/github/include-diagrams-markdown-files-mermaid/ -->

## 📋 Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Flowcharts](#flowcharts)
- [Sequence Diagrams](#sequence-diagrams)
- [Class Diagrams](#class-diagrams)
- [State Diagrams](#state-diagrams)
- [Entity Relationship Diagrams](#entity-relationship-diagrams)
- [Gantt Charts](#gantt-charts)
- [Architecture Patterns](#architecture-patterns)
- [Best Practices](#best-practices)
- [GitHub-Specific Considerations](#github-specific-considerations)
- [Troubleshooting](#troubleshooting)
- [Advanced Examples](#advanced-examples)

## Overview

Mermaid is a JavaScript-based diagramming tool that renders markdown-like text
definitions into diagrams. GitHub natively supports Mermaid diagrams in markdown
files, making it perfect for documenting software architecture, workflows, and
system designs directly in your repository.

### Supported Diagram Types

GitHub supports these Mermaid diagram types:

- **Flowcharts** - Process flows and decision trees
- **Sequence Diagrams** - Interactions between entities over time
- **Class Diagrams** - Object-oriented structure and relationships
- **State Diagrams** - State machines and transitions
- **Entity Relationship Diagrams** - Database schemas and relationships
- **User Journey** - User interaction flows
- **Gantt Charts** - Project timelines and schedules
- **Pie Charts** - Data distribution visualization
- **Git Graphs** - Git branching strategies

## Getting Started

### Basic Syntax

Mermaid diagrams in GitHub markdown are enclosed in code fences with `mermaid` as the language identifier:

````markdown
```mermaid
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action 1]
    B -->|No| D[Action 2]
```
````

```mermaid
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action 1]
    B -->|No| D[Action 2]
```

### Essential Rules

1. **Use triple backticks** with `mermaid` language identifier
2. **Start with diagram type declaration** (e.g., `graph`, `sequenceDiagram`, `classDiagram`)
3. **Follow consistent indentation** (spaces or tabs, but be consistent)
4. **Use meaningful node IDs** that describe the component
5. **Keep lines under 100 characters** for better readability

## Flowcharts

Flowcharts are perfect for documenting process flows, decision trees, and system architectures.

### Basic Flowchart Syntax

```mermaid
graph TD
    Start([Start]) --> Input[/User Input/]
    Input --> Validate{Valid?}
    Validate -->|Yes| Process[Process Data]
    Validate -->|No| Error[Show Error]
    Error --> Input
    Process --> Save[(Save to DB)]
    Save --> Success([Success])
```

### Node Shapes

```mermaid
graph LR
    A[Rectangle] --> B(Rounded Rectangle)
    B --> C([Stadium/Pill])
    C --> D[[Subroutine]]
    D --> E[(Database)]
    E --> F{Diamond/Decision}
    F --> G{{Hexagon}}
    G --> H[/Parallelogram/]
    H --> I[\Parallelogram Alt\]
    I --> J>/Flag/]
```

### Clean Architecture Example

```mermaid
graph TB
    subgraph "Interface Layer"
        HTTP[HTTP Controllers]
        CLI[CLI Commands]
        gRPC[gRPC Handlers]
    end
    
    subgraph "Application Layer"
        UseCase[Use Cases]
        CommandHandler[Command Handlers]
        QueryHandler[Query Handlers]
    end
    
    subgraph "Domain Layer"
        Entity[Entities]
        ValueObject[Value Objects]
        DomainService[Domain Services]
        Repository[Repository Interfaces]
    end
    
    subgraph "Infrastructure Layer"
        Database[(PostgreSQL)]
        FileSystem[File System]
        ExternalAPI[External APIs]
        RepoImpl[Repository Implementations]
    end
    
    HTTP --> UseCase
    CLI --> CommandHandler
    gRPC --> QueryHandler
    
    UseCase --> Entity
    CommandHandler --> DomainService
    QueryHandler --> Repository
    
    Repository --> RepoImpl
    RepoImpl --> Database
    RepoImpl --> FileSystem
    DomainService --> ExternalAPI
```

## Sequence Diagrams

Perfect for documenting API interactions, system communications, and process flows over time.

### Basic Sequence Diagram

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant Service
    participant DB
    
    Client->>API: POST /users
    activate API
    
    API->>Service: CreateUser(req)
    activate Service
    
    Service->>DB: INSERT user
    activate DB
    DB-->>Service: user_id
    deactivate DB
    
    Service-->>API: User created
    deactivate Service
    
    API-->>Client: 201 Created
    deactivate API
```

### CQRS Pattern Example

```mermaid
sequenceDiagram
    participant Client
    participant CommandAPI as Command API
    participant QueryAPI as Query API
    participant CommandHandler as Command Handler
    participant QueryHandler as Query Handler
    participant WriteDB as Write Database
    participant ReadDB as Read Database
    participant EventBus as Event Bus
    
    Note over Client,EventBus: Command Flow (Write)
    Client->>CommandAPI: POST /orders
    CommandAPI->>CommandHandler: CreateOrderCommand
    CommandHandler->>WriteDB: Save Order
    CommandHandler->>EventBus: OrderCreatedEvent
    CommandAPI-->>Client: 202 Accepted
    
    Note over EventBus,ReadDB: Event Processing
    EventBus->>QueryHandler: OrderCreatedEvent
    QueryHandler->>ReadDB: Update Order Projection
    
    Note over Client,ReadDB: Query Flow (Read)
    Client->>QueryAPI: GET /orders/123
    QueryAPI->>QueryHandler: GetOrderQuery
    QueryHandler->>ReadDB: SELECT order
    QueryHandler-->>QueryAPI: Order Data
    QueryAPI-->>Client: 200 OK
```

### Error Handling Flow

```mermaid
sequenceDiagram
    participant User
    participant API
    participant Validator
    participant Service
    participant DB
    
    User->>API: Invalid Request
    API->>Validator: Validate(request)
    Validator-->>API: ValidationError
    API-->>User: 400 Bad Request
    
    User->>API: Valid Request
    API->>Validator: Validate(request)
    Validator-->>API: OK
    API->>Service: ProcessRequest
    Service->>DB: Query
    DB-->>Service: Connection Error
    Service-->>API: DatabaseError
    API-->>User: 500 Internal Error
```

## Class Diagrams

Excellent for documenting Go struct relationships, interfaces, and object-oriented designs.

### Basic Class Diagram

```mermaid
classDiagram
    class User {
        +ID UserID
        +Email string
        +Name string
        +CreatedAt time.Time
        +Validate() error
        +ChangeEmail(email string) error
    }
    
    class UserRepository {
        <<interface>>
        +Save(user User) error
        +FindByID(id UserID) User
        +FindByEmail(email string) User
    }
    
    class PostgreSQLUserRepository {
        -db *sql.DB
        +Save(user User) error
        +FindByID(id UserID) User
        +FindByEmail(email string) User
    }
    
    class UserService {
        -repo UserRepository
        +CreateUser(email string) User
        +GetUser(id UserID) User
    }
    
    UserRepository <|-- PostgreSQLUserRepository : implements
    UserService --> UserRepository : uses
    UserService ..> User : creates
```

### Domain Model Example

```mermaid
classDiagram
    class Order {
        +ID OrderID
        +CustomerID CustomerID
        +Status OrderStatus
        +Items []OrderItem
        +Total Money
        +AddItem(item OrderItem)
        +RemoveItem(itemID ItemID)
        +Submit() error
        +Cancel() error
    }
    
    class OrderItem {
        +ID ItemID
        +ProductID ProductID
        +Quantity int
        +UnitPrice Money
        +Subtotal() Money
    }
    
    class Money {
        +Amount decimal.Decimal
        +Currency string
        +Add(other Money) Money
        +Multiply(factor decimal.Decimal) Money
    }
    
    class Customer {
        +ID CustomerID
        +Email Email
        +Name string
        +PlaceOrder(items []OrderItem) Order
    }
    
    Order "1" *-- "*" OrderItem : contains
    OrderItem --> Money : has unit price
    Order --> Money : has total
    Customer --> Order : places
```

## State Diagrams

Perfect for documenting state machines, workflow states, and business process flows.

### Order State Machine

```mermaid
stateDiagram-v2
    [*] --> Draft
    
    Draft --> Pending : submit()
    Draft --> Cancelled : cancel()
    
    Pending --> Confirmed : confirm()
    Pending --> Cancelled : cancel()
    
    Confirmed --> Processing : startProcessing()
    Confirmed --> Cancelled : cancel()
    
    Processing --> Shipped : ship()
    Processing --> Cancelled : cancel()
    
    Shipped --> Delivered : deliver()
    
    Delivered --> [*]
    Cancelled --> [*]
    
    state Processing {
        [*] --> Preparing
        Preparing --> Packaged : package()
        Packaged --> [*]
    }
```

### User Authentication States

```mermaid
stateDiagram-v2
    [*] --> Anonymous
    
    Anonymous --> Authenticating : login()
    Authenticating --> Authenticated : success
    Authenticating --> Anonymous : failure
    
    Authenticated --> Authenticating : refresh_token()
    Authenticated --> Anonymous : logout()
    Authenticated --> Anonymous : token_expired
    
    state Authenticated {
        [*] --> Active
        Active --> Idle : no_activity
        Idle --> Active : activity
        
        state Active {
            [*] --> Reading
            Reading --> Writing : edit()
            Writing --> Reading : save()
        }
    }
```

## Entity Relationship Diagrams

Great for documenting database schemas and data relationships.

### Database Schema Example

```mermaid
erDiagram
    users {
        uuid id PK
        string email UK
        string name
        timestamp created_at
        timestamp updated_at
    }
    
    orders {
        uuid id PK
        uuid user_id FK
        string status
        decimal total_amount
        string currency
        timestamp created_at
        timestamp updated_at
    }
    
    order_items {
        uuid id PK
        uuid order_id FK
        uuid product_id FK
        int quantity
        decimal unit_price
        string currency
    }
    
    products {
        uuid id PK
        string name
        string description
        decimal price
        string currency
        boolean active
        timestamp created_at
    }
    
    users ||--o{ orders : "places"
    orders ||--o{ order_items : "contains"
    products ||--o{ order_items : "referenced by"
```

### Microservices Data Flow

```mermaid
erDiagram
    USER_SERVICE {
        users table
        user_profiles table
    }
    
    ORDER_SERVICE {
        orders table
        order_items table
    }
    
    INVENTORY_SERVICE {
        products table
        stock_levels table
    }
    
    PAYMENT_SERVICE {
        payments table
        transactions table
    }
    
    USER_SERVICE ||--o{ ORDER_SERVICE : "user_id reference"
    ORDER_SERVICE ||--o{ INVENTORY_SERVICE : "product_id reference"
    ORDER_SERVICE ||--o{ PAYMENT_SERVICE : "order_id reference"
```

## Gantt Charts

Useful for project timelines, feature development schedules, and release planning.

### Project Timeline

```mermaid
gantt
    title Go Microservice Development Timeline
    dateFormat  YYYY-MM-DD
    section Planning
    Requirements Analysis    :req, 2024-01-01, 2024-01-15
    Architecture Design     :arch, after req, 10d
    
    section Development
    Domain Layer            :domain, 2024-01-20, 20d
    Application Layer       :app, after domain, 15d
    Infrastructure Layer    :infra, after app, 20d
    Interface Layer         :interface, after infra, 10d
    
    section Testing
    Unit Tests             :unit, after domain, 30d
    Integration Tests      :integration, after infra, 15d
    E2E Tests             :e2e, after interface, 10d
    
    section Deployment
    CI/CD Setup           :cicd, 2024-02-15, 10d
    Staging Deployment    :staging, after e2e, 5d
    Production Deployment :prod, after staging, 5d
```

## Architecture Patterns

### Microservices Architecture

```mermaid
graph TB
    subgraph "API Gateway"
        Gateway[Kong/Nginx]
    end
    
    subgraph "User Service"
        UserAPI[User API]
        UserDB[(User Database)]
        UserAPI --> UserDB
    end
    
    subgraph "Order Service"
        OrderAPI[Order API]
        OrderDB[(Order Database)]
        OrderAPI --> OrderDB
    end
    
    subgraph "Inventory Service"
        InventoryAPI[Inventory API]
        InventoryDB[(Inventory Database)]
        InventoryAPI --> InventoryDB
    end
    
    subgraph "Payment Service"
        PaymentAPI[Payment API]
        PaymentDB[(Payment Database)]
        PaymentAPI --> PaymentDB
    end
    
    subgraph "Message Bus"
        EventBus[Event Bus/Kafka]
    end
    
    Gateway --> UserAPI
    Gateway --> OrderAPI
    Gateway --> InventoryAPI
    Gateway --> PaymentAPI
    
    OrderAPI --> EventBus
    InventoryAPI --> EventBus
    PaymentAPI --> EventBus
    UserAPI --> EventBus
```

### Event-Driven Architecture

```mermaid
graph LR
    subgraph "Command Side"
        CommandAPI[Command API]
        CommandHandler[Command Handler]
        WriteDB[(Write Database)]
        
        CommandAPI --> CommandHandler
        CommandHandler --> WriteDB
    end
    
    subgraph "Event Store"
        EventStore[(Event Store)]
        EventBus[Event Bus]
        
        EventStore --> EventBus
    end
    
    subgraph "Query Side"
        ProjectionHandler[Projection Handler]
        ReadDB[(Read Database)]
        QueryAPI[Query API]
        
        ProjectionHandler --> ReadDB
        ReadDB --> QueryAPI
    end
    
    CommandHandler --> EventStore
    EventBus --> ProjectionHandler
```

## Best Practices

### 1. Naming Conventions

```mermaid
graph TD
    A[Use Clear, Descriptive Names] --> B[UserService not US]
    A --> C[CreateOrderCommand not COC]
    A --> D[PostgreSQLRepository not PSQLRepo]
```

### 2. Consistent Styling

```mermaid
graph TB
    subgraph "External Systems"
        ExtAPI[External API]
        ExtDB[(External Database)]
    end
    
    subgraph "Our System"
        API[Our API]
        Service[Business Service]
        DB[(Our Database)]
        
        API --> Service
        Service --> DB
    end
    
    Service --> ExtAPI
    Service --> ExtDB
    
    classDef external fill:#ffcccc,stroke:#ff0000
    classDef internal fill:#ccffcc,stroke:#00ff00
    
    class ExtAPI,ExtDB external
    class API,Service,DB internal
```

### 3. Layer Separation

```mermaid
graph TB
    subgraph "Interface Layer"
        direction TB
        HTTP[HTTP Controllers]
        gRPC[gRPC Handlers]
        CLI[CLI Commands]
    end
    
    subgraph "Application Layer"
        direction TB
        Commands[Command Handlers]
        Queries[Query Handlers]
        Services[Application Services]
    end
    
    subgraph "Domain Layer"
        direction TB
        Entities[Domain Entities]
        ValueObjects[Value Objects]
        DomainServices[Domain Services]
        Repositories[Repository Interfaces]
    end
    
    subgraph "Infrastructure Layer"
        direction TB
        Database[(Database)]
        FileSystem[File System]
        ExternalAPIs[External APIs]
        RepositoryImpl[Repository Implementations]
    end
    
    HTTP --> Commands
    gRPC --> Queries
    CLI --> Services
    
    Commands --> Entities
    Queries --> ValueObjects
    Services --> DomainServices
    
    DomainServices --> Repositories
    Repositories --> RepositoryImpl
    
    RepositoryImpl --> Database
    RepositoryImpl --> FileSystem
    RepositoryImpl --> ExternalAPIs
```

### 4. Error Flow Documentation

```mermaid
graph TD
    Start([Request]) --> Validate{Valid Input?}
    Validate -->|No| ValidationError[Validation Error]
    Validate -->|Yes| Authorize{Authorized?}
    
    Authorize -->|No| AuthError[Authorization Error]
    Authorize -->|Yes| Process[Process Request]
    
    Process --> BusinessLogic{Business Rules OK?}
    BusinessLogic -->|No| BusinessError[Business Logic Error]
    BusinessLogic -->|Yes| Persist[Persist Data]
    
    Persist --> DatabaseOp{Database Success?}
    DatabaseOp -->|No| DatabaseError[Database Error]
    DatabaseOp -->|Yes| Success[Success Response]
    
    ValidationError --> ErrorResponse[400 Bad Request]
    AuthError --> ErrorResponse2[401/403 Error]
    BusinessError --> ErrorResponse3[422 Unprocessable Entity]
    DatabaseError --> ErrorResponse4[500 Internal Server Error]
    
    Success --> SuccessResponse[200/201 Success]
```

## GitHub-Specific Considerations

### 1. Rendering Limitations

- **Maximum diagram size**: Keep diagrams reasonably sized for mobile viewing
- **Performance**: Very large diagrams may load slowly
- **Browser compatibility**: Test across different browsers

### 2. Accessibility

```mermaid
graph LR
    A[Clear Node Labels] --> B[Descriptive Text]
    B --> C[Good Color Contrast]
    C --> D[Logical Flow Direction]
```

### 3. Documentation Integration

```markdown
## System Architecture

The following diagram shows our Clean Architecture implementation:

```mermaid
graph TB
    %% This diagram shows the dependency flow in our Clean Architecture
    %% Dependencies point inward toward the domain layer
    subgraph "Infrastructure"
        DB[(Database)]
        API[External API]
    end
    
    subgraph "Interface"
        HTTP[HTTP Controllers]
    end
    
    subgraph "Application"
        UseCase[Use Cases]
    end
    
    subgraph "Domain"
        Entity[Entities]
        Repo[Repository Interface]
    end
    
    HTTP --> UseCase
    UseCase --> Entity
    UseCase --> Repo
    DB -.-> Repo
    API -.-> Repo
```

Key principles shown:

- Interface layer depends on Application layer
- Application layer depends on Domain layer
- Infrastructure implements Domain interfaces (dependency inversion)

### 4. File Organization

```consoletext
docs/
├── architecture/
│   ├── overview.md           # High-level system overview
│   ├── components.md         # Detailed component diagrams
│   └── data-flow.md         # Data flow diagrams
├── api/
│   ├── authentication.md    # Auth flow diagrams
│   └── endpoints.md         # API sequence diagrams
└── deployment/
    ├── infrastructure.md     # Infrastructure diagrams
    └── ci-cd.md             # CI/CD pipeline diagrams
```

## Troubleshooting

### Common Issues

1. **Diagram not rendering**

Check: -->
- Triple backticks with 'mermaid' language identifier
- Valid Mermaid syntax
- No extra spaces before diagram type

2. **Syntax errors**

   ```mermaid
   graph TD
       A[Good: Valid node ID] --> B[Another valid node]
       %% Bad: Special characters in node IDs
       %% A-1[Bad: Hyphen in ID] --> B+1[Bad: Plus in ID]
   ```

3. **Layout issues**

   ```mermaid
   graph TD
       %% Use subgraphs for better organization
       subgraph "Group 1"
           A --> B
       end
       
       subgraph "Group 2"
           C --> D
       end
       
       B --> C
   ```

### Debugging Tips

1. **Start simple** and add complexity gradually
2. **Use online editors** like [Mermaid Live Editor](https://mermaid.live/) for testing
3. **Check console errors** in browser dev tools
4. **Validate syntax** before committing

### Performance Optimization

```mermaid
graph LR
    A[Keep diagrams focused] --> B[Use subgraphs for organization]
    B --> C[Split large diagrams into multiple files]
    C --> D[Link related diagrams in documentation]
```

## Advanced Examples

### Complex Deployment Pipeline

```mermaid
graph TB
    subgraph "Development"
        Dev[Developer]
        PR[Pull Request]
        Dev --> PR
    end
    
    subgraph "CI/CD Pipeline"
        Build[Build & Test]
        Security[Security Scan]
        Deploy[Deploy to Staging]
        E2E[E2E Tests]
        Production[Deploy to Production]
        
        PR --> Build
        Build --> Security
        Security --> Deploy
        Deploy --> E2E
        E2E --> Production
    end
    
    subgraph "Monitoring"
        Metrics[Metrics Collection]
        Alerts[Alert System]
        Logs[Log Aggregation]
        
        Production --> Metrics
        Production --> Alerts
        Production --> Logs
    end
    
    subgraph "Environments"
        Staging[Staging Environment]
        Prod[Production Environment]
        
        Deploy --> Staging
        Production --> Prod
    end
```

### Multi-Service Communication

```mermaid
sequenceDiagram
    participant Client
    participant Gateway as API Gateway
    participant Auth as Auth Service
    participant User as User Service
    participant Order as Order Service
    participant Payment as Payment Service
    participant Notification as Notification Service
    
    Client->>Gateway: Create Order Request
    Gateway->>Auth: Validate Token
    Auth-->>Gateway: Token Valid
    
    Gateway->>User: Get User Details
    User-->>Gateway: User Data
    
    Gateway->>Order: Create Order
    activate Order
    Order->>Payment: Process Payment
    activate Payment
    Payment-->>Order: Payment Successful
    deactivate Payment
    
    Order->>Notification: Send Order Confirmation
    activate Notification
    Notification-->>Order: Notification Sent
    deactivate Notification
    
    Order-->>Gateway: Order Created
    deactivate Order
    Gateway-->>Client: 201 Created
```

### Event Sourcing Pattern

```mermaid
graph TB
    subgraph "Command Side"
        Command[Command]
        Aggregate[Aggregate]
        Events[Domain Events]
        EventStore[(Event Store)]
        
        Command --> Aggregate
        Aggregate --> Events
        Events --> EventStore
    end
    
    subgraph "Query Side"
        EventHandler[Event Handler]
        Projection[Projection]
        ReadModel[(Read Model)]
        
        EventStore --> EventHandler
        EventHandler --> Projection
        Projection --> ReadModel
    end
    
    subgraph "Snapshots"
        SnapshotStore[(Snapshot Store)]
        
        Aggregate --> SnapshotStore
        SnapshotStore --> Aggregate
    end
```

## Conclusion

Mermaid diagrams in GitHub provide a powerful way to document software architecture,
processes, and workflows directly in your repository. By following these patterns
and best practices, you can create maintainable, accessible, and informative diagrams
that enhance your project documentation.

### Key Takeaways

1. **Start simple** - Begin with basic diagrams and add complexity as needed
2. **Be consistent** - Use consistent naming and styling across all diagrams
3. **Document context** - Always provide explanation text around diagrams
4. **Keep it current** - Update diagrams when architecture changes
5. **Test rendering** - Verify diagrams display correctly on different devices

### Related Documentation

- [C4 Architecture Diagramming Guide](c4_diagramming_guide.md) - For auto-generated architecture diagrams
- [Project Guidelines](project_guidelines.md) - Clean Architecture implementation patterns
- [GitHub CI/CD](github_cicd.md) - Automated workflow documentation

---

*For more advanced diagramming needs, consider combining Mermaid with our C4*
*auto-generation guide for comprehensive architectural documentation.*
