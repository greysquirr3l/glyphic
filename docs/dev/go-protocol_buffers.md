# Protocol Buffers in Go: DDD + CQRS + Security First Guide

## Overview

This document provides comprehensive guidance on using Google's Protocol Buffers (protobuf-go) in Go applications following Domain-Driven Design (DDD), Command Query Responsibility Segregation (CQRS), and Security First methodologies.

**Repository**: https://github.com/protocolbuffers/protobuf-go

## Table of Contents

1. [Installation and Setup](#installation-and-setup)
2. [Core Concepts](#core-concepts)
3. [DDD Integration Patterns](#ddd-integration-patterns)
4. [CQRS Implementation](#cqrs-implementation)
5. [Security First Principles](#security-first-principles)
6. [Best Practices](#best-practices)
7. [Complete Example](#complete-example)

## Installation and Setup

### Prerequisites

Install the Protocol Buffer compiler and Go plugins:

```bash
# Install protoc compiler
# On macOS
brew install protobuf

# On Linux
apt-get install protobuf-compiler

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Project Structure

```
project/
├── domain/
│   ├── aggregates/
│   ├── entities/
│   ├── value_objects/
│   └── events/
├── application/
│   ├── commands/
│   └── queries/
├── infrastructure/
│   ├── proto/
│   │   ├── commands.proto
│   │   ├── queries.proto
│   │   └── events.proto
│   └── persistence/
└── api/
    └── grpc/
```

## Core Concepts

### Protocol Buffers Basics

Protocol Buffers define data structures in `.proto` files. The protoc compiler generates Go code from these definitions.

**Example proto file** (`user.proto`):

```protobuf
syntax = "proto3";

package domain.user.v1;

option go_package = "github.com/yourorg/project/gen/domain/user/v1;userv1";

message User {
  string user_id = 1;
  string email = 2;
  string hashed_password = 3;
  int64 created_at = 4;
  int64 updated_at = 5;
}
```

### Compilation

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    infrastructure/proto/*.proto
```

## DDD Integration Patterns

### Principle: Separate Domain Models from DTOs

**Critical Rule**: Never use generated protobuf structs as domain entities. Protobuf messages are Data Transfer Objects (DTOs), not domain models.

### Domain Layer (Pure Go)

```go
// domain/aggregates/user.go
package aggregates

import (
    "errors"
    "time"
    "github.com/yourorg/project/domain/value_objects"
)

// User is the domain aggregate root
type User struct {
    id        value_objects.UserID
    email     value_objects.Email
    password  value_objects.HashedPassword
    createdAt time.Time
    updatedAt time.Time
    version   int // for optimistic locking
}

// NewUser creates a new user with validation
func NewUser(email, password string) (*User, error) {
    emailVO, err := value_objects.NewEmail(email)
    if err != nil {
        return nil, err
    }

    hashedPwd, err := value_objects.HashPassword(password)
    if err != nil {
        return nil, err
    }

    return &User{
        id:        value_objects.NewUserID(),
        email:     emailVO,
        password:  hashedPwd,
        createdAt: time.Now(),
        updatedAt: time.Now(),
        version:   1,
    }, nil
}

// ChangeEmail is a domain behavior
func (u *User) ChangeEmail(newEmail string) error {
    emailVO, err := value_objects.NewEmail(newEmail)
    if err != nil {
        return err
    }

    u.email = emailVO
    u.updatedAt = time.Now()
    u.version++
    return nil
}
```

### Value Objects with Validation

```go
// domain/value_objects/email.go
package value_objects

import (
    "errors"
    "regexp"
)

type Email struct {
    value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(email string) (Email, error) {
    if !emailRegex.MatchString(email) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: email}, nil
}

func (e Email) String() string {
    return e.value
}
```

### Mapping Layer (Domain ↔ Protobuf)

```go
// infrastructure/mappers/user_mapper.go
package mappers

import (
    "github.com/yourorg/project/domain/aggregates"
    pb "github.com/yourorg/project/gen/domain/user/v1"
)

type UserMapper struct{}

// ToProto converts domain model to protobuf DTO
func (m *UserMapper) ToProto(user *aggregates.User) *pb.User {
    return &pb.User{
        UserId:         user.ID().String(),
        Email:          user.Email().String(),
        HashedPassword: user.Password().Hash(), // Never expose in API!
        CreatedAt:      user.CreatedAt().Unix(),
        UpdatedAt:      user.UpdatedAt().Unix(),
    }
}

// ToDomain converts protobuf DTO to domain model
func (m *UserMapper) ToDomain(pbUser *pb.User) (*aggregates.User, error) {
    // Reconstruction from persistence
    return aggregates.ReconstructUser(
        pbUser.UserId,
        pbUser.Email,
        pbUser.HashedPassword,
        time.Unix(pbUser.CreatedAt, 0),
        time.Unix(pbUser.UpdatedAt, 0),
    )
}
```

## CQRS Implementation

### Command Side (Write Model)

**Proto Definition** (`commands.proto`):

```protobuf
syntax = "proto3";

package application.commands.v1;

option go_package = "github.com/yourorg/project/gen/application/commands/v1;commandsv1";

message CreateUserCommand {
  string command_id = 1;  // idempotency key
  string email = 2;
  string password = 3;
  string ip_address = 4;  // for audit
  string user_agent = 5;  // for audit
}

message CreateUserResponse {
  string user_id = 1;
  bool success = 2;
  string error_message = 3;
}

message ChangeEmailCommand {
  string command_id = 1;
  string user_id = 2;
  string new_email = 3;
  int32 expected_version = 4; // optimistic locking
}

message ChangeEmailResponse {
  bool success = 1;
  string error_message = 2;
  int32 new_version = 3;
}
```

**Command Handler**:

```go
// application/commands/create_user_handler.go
package commands

import (
    "context"
    pb "github.com/yourorg/project/gen/application/commands/v1"
    "github.com/yourorg/project/domain/aggregates"
    "github.com/yourorg/project/domain/repositories"
)

type CreateUserHandler struct {
    userRepo repositories.UserRepository
    eventBus EventBus
    validator SecurityValidator
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd *pb.CreateUserCommand) (*pb.CreateUserResponse, error) {
    // Security: Validate command source
    if err := h.validator.ValidateCommand(ctx, cmd); err != nil {
        return &pb.CreateUserResponse{
            Success: false,
            ErrorMessage: "unauthorized",
        }, nil
    }

    // Security: Check for SQL injection attempts in email
    if containsSQLPatterns(cmd.Email) {
        return &pb.CreateUserResponse{
            Success: false,
            ErrorMessage: "invalid input detected",
        }, nil
    }

    // Create domain aggregate
    user, err := aggregates.NewUser(cmd.Email, cmd.Password)
    if err != nil {
        return &pb.CreateUserResponse{
            Success: false,
            ErrorMessage: err.Error(),
        }, nil
    }

    // Persist using repository
    if err := h.userRepo.Save(ctx, user); err != nil {
        return &pb.CreateUserResponse{
            Success: false,
            ErrorMessage: "failed to create user",
        }, err
    }

    // Publish domain event
    event := &UserCreatedEvent{
        UserID: user.ID().String(),
        Email: user.Email().String(),
        Timestamp: time.Now(),
    }
    h.eventBus.Publish(ctx, event)

    return &pb.CreateUserResponse{
        UserId: user.ID().String(),
        Success: true,
    }, nil
}
```

### Query Side (Read Model)

**Proto Definition** (`queries.proto`):

```protobuf
syntax = "proto3";

package application.queries.v1;

option go_package = "github.com/yourorg/project/gen/application/queries/v1;queriesv1";

message GetUserQuery {
  string user_id = 1;
  repeated string fields = 2; // field filtering for security
}

message UserDTO {
  string user_id = 1;
  string email = 2;
  // Note: NO password field!
  int64 created_at = 3;
  int64 updated_at = 4;
}

message ListUsersQuery {
  int32 page = 1;
  int32 page_size = 2;
  string sort_by = 3;
  repeated string filters = 4;
}

message ListUsersResponse {
  repeated UserDTO users = 1;
  int32 total_count = 2;
  int32 page = 3;
}
```

**Query Handler**:

```go
// application/queries/get_user_handler.go
package queries

import (
    "context"
    pb "github.com/yourorg/project/gen/application/queries/v1"
)

type GetUserHandler struct {
    readModel UserReadModel
    authz AuthorizationService
}

func (h *GetUserHandler) Handle(ctx context.Context, query *pb.GetUserQuery) (*pb.UserDTO, error) {
    // Security: Check authorization
    if !h.authz.CanReadUser(ctx, query.UserId) {
        return nil, errors.New("unauthorized")
    }

    // Security: Validate and sanitize user_id
    if !isValidUUID(query.UserId) {
        return nil, errors.New("invalid user id")
    }

    // Query read model (denormalized, optimized for reads)
    user, err := h.readModel.GetUser(ctx, query.UserId)
    if err != nil {
        return nil, err
    }

    // Security: Field-level filtering
    return h.filterFields(user, query.Fields), nil
}

func (h *GetUserHandler) filterFields(user *pb.UserDTO, fields []string) *pb.UserDTO {
    if len(fields) == 0 {
        return user
    }

    filtered := &pb.UserDTO{UserId: user.UserId}
    fieldSet := make(map[string]bool)
    for _, f := range fields {
        fieldSet[f] = true
    }

    if fieldSet["email"] {
        filtered.Email = user.Email
    }
    if fieldSet["created_at"] {
        filtered.CreatedAt = user.CreatedAt
    }
    if fieldSet["updated_at"] {
        filtered.UpdatedAt = user.UpdatedAt
    }

    return filtered
}
```

## Security First Principles

### 1. Input Validation at Protobuf Level

Use protobuf validation rules:

```protobuf
syntax = "proto3";

import "buf/validate/validate.proto";

message CreateUserCommand {
  string command_id = 1 [(buf.validate.field).string.uuid = true];
  
  string email = 2 [(buf.validate.field).string = {
    email: true,
    max_len: 255
  }];
  
  string password = 3 [(buf.validate.field).string = {
    min_len: 12,
    max_len: 128
  }];
}
```

### 2. Sensitive Data Handling

```go
// DO NOT include sensitive fields in query responses
message UserDTO {
  string user_id = 1;
  string email = 2;
  // NEVER include: password, password_hash, tokens, secrets
}

// Use separate secure channels for sensitive operations
message ChangePasswordCommand {
  string user_id = 1;
  string current_password = 2; // transmitted over TLS only
  string new_password = 3;
  string mfa_token = 4; // require MFA
}
```

### 3. Audit Logging with Protobuf

```protobuf
syntax = "proto3";

message AuditEvent {
  string event_id = 1;
  string event_type = 2;
  string user_id = 3;
  string ip_address = 4;
  string user_agent = 5;
  int64 timestamp = 6;
  bytes encrypted_payload = 7; // encrypt sensitive audit data
  string checksum = 8; // tamper detection
}
```

### 4. Rate Limiting Metadata

```go
// Embed rate limiting info in protobuf
message RateLimitMetadata {
  string client_id = 1;
  int32 requests_remaining = 2;
  int64 reset_at = 3;
}

message CreateUserResponse {
  string user_id = 1;
  bool success = 2;
  RateLimitMetadata rate_limit = 3;
}
```

### 5. Encryption for Sensitive Fields

```go
// infrastructure/crypto/field_encryption.go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
)

type FieldEncryptor struct {
    gcm cipher.AEAD
}

func NewFieldEncryptor(key []byte) (*FieldEncryptor, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    return &FieldEncryptor{gcm: gcm}, nil
}

func (e *FieldEncryptor) Encrypt(plaintext string) (string, error) {
    nonce := make([]byte, e.gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    ciphertext := e.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *FieldEncryptor) Decrypt(ciphertext string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    nonceSize := e.gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]

    plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
```

### 6. Zero Trust Field Access

```go
type FieldAccessControl struct {
    rules map[string][]string // field -> allowed roles
}

func (fac *FieldAccessControl) FilterMessage(msg proto.Message, userRoles []string) proto.Message {
    // Reflect over protobuf fields and remove unauthorized ones
    // This ensures zero-trust: explicitly allow, implicitly deny
    reflection := msg.ProtoReflect()
    descriptor := reflection.Descriptor()
    fields := descriptor.Fields()

    for i := 0; i < fields.Len(); i++ {
        field := fields.Get(i)
        if !fac.hasAccess(field.Name(), userRoles) {
            reflection.Clear(field)
        }
    }

    return msg
}
```

## Best Practices

### 1. Versioning Strategy

```
infrastructure/proto/
├── v1/
│   ├── commands.proto
│   ├── queries.proto
│   └── events.proto
└── v2/
    ├── commands.proto
    ├── queries.proto
    └── events.proto
```

Use package versioning:

```protobuf
package application.commands.v1;
option go_package = "github.com/yourorg/project/gen/application/commands/v1;commandsv1";
```

### 2. Backward Compatibility Rules

- Never change field numbers
- Never reuse field numbers
- Use `reserved` for deprecated fields
- Add new fields with new numbers
- Use `oneof` for evolving alternatives

```protobuf
message UserCommand {
  reserved 2; // old_field removed
  reserved "old_field";
  
  string user_id = 1;
  string new_field = 3; // safe to add
}
```

### 3. Domain Event Sourcing

```protobuf
syntax = "proto3";

package domain.events.v1;

message UserEvent {
  string event_id = 1;
  string aggregate_id = 2;
  int32 version = 3;
  int64 timestamp = 4;
  
  oneof event {
    UserCreated user_created = 10;
    EmailChanged email_changed = 11;
    PasswordChanged password_changed = 12;
    UserDeleted user_deleted = 13;
  }
}

message UserCreated {
  string email = 1;
  int64 created_at = 2;
}

message EmailChanged {
  string old_email = 1;
  string new_email = 2;
  int64 changed_at = 3;
}
```

### 4. Repository Pattern with Protobuf

```go
// domain/repositories/user_repository.go
package repositories

type UserRepository interface {
    Save(ctx context.Context, user *aggregates.User) error
    FindByID(ctx context.Context, id string) (*aggregates.User, error)
    FindByEmail(ctx context.Context, email string) (*aggregates.User, error)
}

// infrastructure/persistence/user_repository_impl.go
package persistence

type UserRepositoryImpl struct {
    store EventStore // stores protobuf events
    mapper *mappers.UserMapper
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user *aggregates.User) error {
    // Convert to protobuf for storage
    events := user.GetUncommittedEvents()
    
    for _, domainEvent := range events {
        pbEvent := r.mapper.EventToProto(domainEvent)
        if err := r.store.Append(ctx, pbEvent); err != nil {
            return err
        }
    }
    
    user.MarkEventsAsCommitted()
    return nil
}
```

### 5. Testing Strategy

```go
// Test with protobuf fixtures
func TestCreateUserCommand(t *testing.T) {
    cmd := &pb.CreateUserCommand{
        CommandId: "123e4567-e89b-12d3-a456-426614174000",
        Email: "test@example.com",
        Password: "SecurePass123!",
    }

    handler := NewCreateUserHandler(mockRepo, mockEventBus, mockValidator)
    resp, err := handler.Handle(context.Background(), cmd)

    assert.NoError(t, err)
    assert.True(t, resp.Success)
    assert.NotEmpty(t, resp.UserId)
}

// Security test
func TestCreateUserCommand_SQLInjection(t *testing.T) {
    cmd := &pb.CreateUserCommand{
        CommandId: "123e4567-e89b-12d3-a456-426614174000",
        Email: "test@example.com'; DROP TABLE users; --",
        Password: "SecurePass123!",
    }

    handler := NewCreateUserHandler(mockRepo, mockEventBus, mockValidator)
    resp, err := handler.Handle(context.Background(), cmd)

    assert.NoError(t, err)
    assert.False(t, resp.Success)
    assert.Contains(t, resp.ErrorMessage, "invalid")
}
```

## Complete Example

### Project Setup

```bash
# Initialize Go module
go mod init github.com/yourorg/project

# Install dependencies
go get google.golang.org/protobuf
go get google.golang.org/grpc

# Create directory structure
mkdir -p {domain/{aggregates,entities,value_objects,events,repositories},application/{commands,queries},infrastructure/{proto,persistence,mappers,crypto},api/grpc}
```

### Complete Working Example

This example demonstrates a user registration system with DDD, CQRS, and security:

**1. Proto Definitions** (infrastructure/proto/user_service.proto):

```protobuf
syntax = "proto3";

package userservice.v1;

option go_package = "github.com/yourorg/project/gen/userservice/v1;userservicev1";

import "buf/validate/validate.proto";

service UserService {
  rpc RegisterUser(RegisterUserRequest) returns (RegisterUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}

message RegisterUserRequest {
  string email = 1 [(buf.validate.field).string.email = true];
  string password = 2 [(buf.validate.field).string.min_len = 12];
}

message RegisterUserResponse {
  string user_id = 1;
  bool success = 2;
  string error = 3;
}

message GetUserRequest {
  string user_id = 1 [(buf.validate.field).string.uuid = true];
}

message GetUserResponse {
  string user_id = 1;
  string email = 2;
  int64 created_at = 3;
}
```

**2. Compile**:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    infrastructure/proto/user_service.proto
```

**3. Implementation** follows patterns shown above with proper separation of concerns, security validation, and domain logic isolation.

## Summary

This guide demonstrates how to effectively use Protocol Buffers in Go while maintaining:

- **DDD**: Domain models remain pure Go, protobuf is used only for serialization
- **CQRS**: Separate command and query models with different protobuf definitions
- **Security First**: Input validation, field-level access control, encryption, and audit logging

Key takeaways:
1. Protobuf messages are DTOs, not domain entities
2. Always validate and sanitize inputs
3. Never expose sensitive data in query responses
4. Use versioning for API evolution
5. Implement proper mapping layers between domain and infrastructure


