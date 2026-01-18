# Kinm Architecture

## Overview

Kinm is a database-backed API server providing Kubernetes-like CRUD+Watch semantics without the complexity of etcd. All state is persisted in PostgreSQL (or SQLite for development).

```
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP Clients                             │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       pkg/server                                 │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────────────┐  │
│  │ HTTP Server │  │ Middleware   │  │ k8s GenericAPIServer   │  │
│  └─────────────┘  └──────────────┘  └────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       pkg/stores                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              Builder (fluent configuration)              │    │
│  └─────────────────────────────────────────────────────────┘    │
│                              │                                   │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────────┐   │
│  │ Complete │ GetList  │ ReadOnly │ ListWatch│ ... (15+)    │   │
│  └──────────┴──────────┴──────────┴──────────┴──────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      pkg/strategy                                │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Adapters (Create, Get, List, Update, Delete, Watch)      │  │
│  │  - Translates k8s apiserver interfaces to Kinm interfaces │  │
│  │  - Validation, preparation, warnings                       │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                         pkg/db                                   │
│  ┌──────────────────┐      ┌──────────────────────────────────┐ │
│  │     Factory      │─────▶│           Strategy               │ │
│  │  - Connection    │      │  - Create/Get/List/Update/Delete │ │
│  │  - Schema        │      │  - Watch (long-poll + broadcast) │ │
│  │  - Migration     │      │  - Compaction (background)       │ │
│  └──────────────────┘      └──────────────────────────────────┘ │
│                                        │                         │
│                                        ▼                         │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                   pkg/db/statements                        │  │
│  │  - Embedded SQL files (.sql)                               │  │
│  │  - Parameterized queries                                   │  │
│  │  - PostgreSQL/SQLite compatible                            │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    PostgreSQL / SQLite                           │
└─────────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### 1. Database-First State Management

Unlike Kubernetes (etcd) or Kine/Mink (in-memory cache), Kinm keeps **all state in the database**.

**Benefits:**

- No memory pressure from large resource counts
- Stateless API servers (horizontal scaling)
- Database handles durability and consistency

**Trade-offs:**

- Higher database load
- Watch requires polling mechanism

### 2. Record Versioning via Linked List

Each record has a `previous_id` forming a version chain:

```
┌────────┐     ┌────────┐     ┌────────┐
│ id: 1  │◀────│ id: 2  │◀────│ id: 3  │  (current)
│ v1     │     │ v2     │     │ v3     │
│ created│     │        │     │        │
└────────┘     └────────┘     └────────┘
```

**Purpose:**

- Enables Watch by scanning `id > lastSeen`
- Provides optimistic concurrency via `previous_id`
- Supports conflict detection

### 3. Compaction

Background process removes old versions:

```sql
-- Removes records where:
-- 1. Not the latest (another record points to it as previous_id)
-- 2. Below compaction watermark
```

**Process:**

- Runs every 15 minutes per table
- Updates `compaction.id` watermark
- Clients with old `resourceVersion` get `410 Gone`

### 4. Watch Implementation

Long-polling with broadcast notification:

```go
// Producer (on write)
s.broadcastLock.Lock()
close(s.broadcast)
s.broadcast = make(chan struct{})
s.broadcastLock.Unlock()

// Consumer (watch)
for {
    select {
    case <-ctx.Done():
        return
    case <-s.broadcast:
        // Query for changes since lastID
    }
}
```

**Flow:**

1. Client connects with `resourceVersion`
2. Server queries for records with `id > resourceVersion`
3. Returns events or waits for broadcast
4. On any write, all waiters are notified

### 5. Field Indexing

Dynamic columns for field selectors:

```sql
-- Table with extra field columns
CREATE TABLE myresource (
    ...
    field_status VARCHAR(255),
    field_type VARCHAR(255)
);

-- Query with field selector
SELECT * FROM myresource
WHERE field_status = 'Active'
AND rn = 1 AND deleted = 0;
```

**Implementation:**

- Types implement `FieldNames()` interface
- Migration adds columns automatically
- Queries filter on indexed fields

## Package Responsibilities

### pkg/server

- HTTP server wrapping k8s genericapiserver
- TLS/authentication/authorization
- Middleware chain
- OpenAPI generation

### pkg/stores

- Pre-configured store combinations
- Builder pattern for custom stores
- Maps store variants to strategy adapters

### pkg/strategy

- Adapters between k8s apiserver and Kinm
- Validation hooks
- Namespace scoping
- Table conversion
- OpenTelemetry tracing

### pkg/db

- Database connection management
- SQL execution with GORM
- Strategy implementation
- Background compaction
- Watch broadcast coordination

### pkg/db/statements

- Embedded SQL files
- Template substitution (table names, field columns)
- PostgreSQL/SQLite compatibility

### pkg/types

- Base `Object` and `ObjectList` interfaces
- Field selection interfaces
- Attribute extraction for filtering

### pkg/apigroup

- API group registration helpers
- Scheme management

### pkg/authn

- Static token authentication
- Bearer token extraction

### pkg/otel

- OpenTelemetry attribute helpers
- Trace span attributes for operations

## Data Flow

### Create Operation

```
Client POST /api/v1/namespaces/default/myresources
    │
    ▼
Server (authentication, authorization)
    │
    ▼
CreateAdapter.Create()
    ├── FillObjectMetaSystemFields (UID, timestamps)
    ├── GenerateName (if generateName set)
    ├── BeforeCreate (validation)
    └── strategy.Create()
            │
            ▼
        db.Strategy.Create()
            ├── Serialize to JSON
            ├── INSERT with created=1
            └── Return with resourceVersion=id
```

### Watch Operation

```
Client GET /api/v1/myresources?watch=true&resourceVersion=100
    │
    ▼
WatchAdapter.Watch()
    │
    ▼
db.Strategy.Watch()
    ├── Query records with id > 100
    ├── Send initial events
    └── Loop:
          ├── Wait for broadcast OR timeout
          ├── Query new records
          └── Send events
```

### List with Field Selector

```
Client GET /api/v1/myresources?fieldSelector=status.phase=Running
    │
    ▼
ListAdapter.List()
    ├── Parse field selector
    ├── Build predicate
    └── strategy.List(opts)
            │
            ▼
        db.Strategy.List()
            ├── Build SQL with field filters
            ├── Window function for latest per name/namespace
            └── Return filtered list
```

## Scalability Considerations

### Horizontal Scaling

- API servers are stateless
- Database handles coordination
- Watch broadcast is per-process (no cross-server coordination)

### Database Performance

- Connection pooling (configurable)
- Index on (name, namespace, created)
- Compaction prevents unbounded growth

### Watch Efficiency

- Broadcast channel avoids polling
- Clients share same broadcast notification
- Query only changed records

## Security Model

### Authentication

- Pluggable authenticator interface
- Static token for simple deployments
- Anonymous fallback available

### Authorization

- Pluggable authorizer interface
- Per-request authorization
- Integrates with k8s RBAC concepts

### Input Validation

- Name validation (DNS subdomain by default)
- Namespace scoping enforcement
- Custom validators via hooks
