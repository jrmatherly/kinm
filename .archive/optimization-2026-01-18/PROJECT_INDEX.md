# Project Index: Kinm

Generated: 2026-01-15

## Overview

Kinm (pronounced "kim") is a Go-based API server providing Kubernetes-like CRUD+Watch operations backed by PostgreSQL. Unlike k8s, it keeps all state in the database, not memory.

**Module:** `github.com/obot-platform/kinm`
**Go Version:** 1.24.0

## Project Structure

```
kinm/
├── go.mod, go.sum          # Dependencies
├── README.md               # Project documentation
└── pkg/
    ├── apigroup/           # API group registration
    ├── authn/              # Authentication (static token)
    ├── db/                 # Database layer (core)
    │   ├── errors/         # Custom DB errors
    │   ├── glogrus/        # GORM-logrus bridge
    │   └── statements/     # Embedded SQL (13 files)
    ├── otel/               # OpenTelemetry attributes
    ├── serializer/         # Object serialization
    ├── server/             # HTTP server setup
    ├── stores/             # Store interface builders
    ├── strategy/           # CRUD+Watch strategies
    │   ├── remote/         # Remote strategy
    │   └── translation/    # Field translation
    ├── types/              # Core types (Object, Fields)
    └── validator/          # Name validation
```

## Core Modules

### pkg/db - Database Layer

**Purpose:** PostgreSQL/SQLite storage with versioned records

- `Factory` - Creates DB connections and strategies
- `Strategy` - Implements Create/Get/Update/Delete/List/Watch
- `statements/` - 13 embedded SQL files for all operations

**Key Types:**

- `Factory{DB, SQLDB, schema, migrationTimeout}`
- `Strategy{db, objTemplate, broadcast...}`

### pkg/server - API Server

**Purpose:** HTTP server wrapping k8s genericapiserver

- `Server{Config, GenericAPIServer, Loopback}`
- `Config{Name, Authenticator, Authorization, HTTPListenPort...}`

**Entry Points:**

- `New(config *Config)` - Create server
- `(*Server).Run(ctx)` - Start HTTP server
- `(*Server).Handler(ctx)` - Get HTTP handler

### pkg/stores - Store Builders

**Purpose:** Pre-built store configurations via builder pattern

- `Builder` - Fluent interface for store construction
- 15+ store variants: `complete`, `createget`, `getlist`, `readwritewatch`, etc.

**Key Methods:**

- `NewBuilder(scheme, obj)` - Start building
- `WithCompleteCRUD()`, `WithWatch()`, `WithGet()`, `WithList()`...
- `Build()` - Produce final store

### pkg/strategy - CRUD Strategies

**Purpose:** Implements k8s apiserver storage interfaces

- `Base` - Interface combining Storage, Scoper, TableConvertor
- `CompleteStrategy` - Full CRUD+Watch interface
- `Watcher` - Watch interface for streaming changes

**Operations:** create, get, list, update, delete, watch, destroy

### pkg/types - Core Types

- `Object` - Base object with TypeMeta, ObjectMeta, Spec, Status
- `ObjectList` - List container
- `NamespaceScoper` - Interface for namespace scoping

### pkg/apigroup - API Registration

- `AddToScheme(scheme)` - Register types
- `ForStores(scheme, stores)` - Create API group from stores

### pkg/authn - Authentication

- `StaticToken` - Static bearer token authenticator
- `GetBearerToken(req)` - Extract token from request

## Database Schema

**Main Table (per resource type):**

```sql
CREATE TABLE {name} (
    id          INTEGER PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    namespace   VARCHAR(255) NOT NULL,
    previous_id INTEGER UNIQUE,      -- Version chain
    uid         VARCHAR(255),
    created     INTEGER,             -- 1=create, NULL=update
    deleted     INTEGER DEFAULT 0,
    value       TEXT NOT NULL,       -- JSON payload
    -- Dynamic field columns for indexing
);
```

**Compaction Table:**

```sql
CREATE TABLE compaction (
    name VARCHAR(255) UNIQUE,
    id   INTEGER                    -- Last compacted ID
);
```

## Key Dependencies

| Package | Purpose |
| --------- | --------- |
| `k8s.io/apiserver` | API server framework |
| `k8s.io/apimachinery` | K8s types, errors |
| `gorm.io/gorm` | ORM for database access |
| `gorm.io/driver/postgres` | PostgreSQL driver |
| `github.com/sirupsen/logrus` | Logging |
| `go.opentelemetry.io/otel` | Tracing |
| `github.com/stretchr/testify` | Testing |

## Test Coverage

| Package | Test Files |
| --------- | ------------ |
| pkg/db | `db_test.go`, `strategy_test.go` |

**Test Commands:**

```bash
go test ./...                              # All tests (SQLite)
KINM_TEST_DB=postgres go test ./pkg/db/... # PostgreSQL tests
```

## Quick Start

```bash
# Install dependencies
go mod download

# Build
go build ./...

# Run tests
go test ./...

# Format code
go fmt ./...
```

## File Statistics

| Category | Count |
| ---------- | ------- |
| Go source files | 57 |
| Test files | 2 |
| SQL files | 13 |
| Total packages | 12 |

## Documentation

| Document | Description |
| ---------- | ------------- |
| [docs/API.md](docs/API.md) | Complete API reference with interfaces and usage |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | System architecture and design decisions |

## Architecture Notes

1. **Versioning:** Records form a chain via `previous_id`, enabling watch
2. **Compaction:** Background process removes old versions
3. **Field Indexing:** Dynamic columns for field selectors
4. **Watch:** Long-polling with broadcast notifications
5. **No Memory State:** All state in database (unlike Kine/Mink)
