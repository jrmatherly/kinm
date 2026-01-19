# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**kinm** (pronounced "kim" - Kinm Is Not Mink) is a Go-based API server providing Kubernetes-like CRUD+Watch operations backed by PostgreSQL or SQLite. Unlike Kubernetes, it keeps all state in the database rather than memory, enabling efficient horizontal scaling.

**Key Philosophy:** Embrace SQL for state management. No in-memory caching - all state lives in the database for consistency and scalability.

- **Module**: `github.com/obot-platform/kinm`
- **Go Version**: 1.25.0
- **Primary Database**: PostgreSQL (SQLite for development/testing)
- **Key Dependencies**: k8s.io/apiserver, GORM, OpenTelemetry

## Development Commands

### Building

```bash
make build            # Build the project
go build ./...        # Direct build
```

### Testing

```bash
# Unit tests (SQLite)
make test             # Run all tests with race detector
make test-short       # Run tests in short mode
go test ./...         # Direct test

# Integration tests (PostgreSQL)
make test-integration # Requires PostgreSQL (KINM_TEST_DB=postgres)

# Coverage
make test-coverage    # Generate coverage.html
```

### Code Quality

```bash
make fmt              # Format code (go fmt + goimports)
make vet              # Run go vet
make lint             # Run golangci-lint
make validate         # Run all validation (lint + vet)
```

### CI Validation

```bash
make ci               # Full CI pipeline: deps-verify, validate, test
make validate-ci      # Validation only (used in CI)
make fmt-check        # Check formatting without modifying
```

### Dependencies

```bash
make deps             # Download dependencies
make deps-tidy        # Tidy go.mod
make deps-verify      # Verify dependency integrity
make deps-update      # Update all dependencies
```

## Architecture

### Core Design: Database-First API Server

```
Request → k8s.io/apiserver → Strategy → GORM → PostgreSQL/SQLite
                                ↓
                          Broadcast ← Watch Clients
```

**Key Architectural Decisions:**

1. **No Memory State**: All state in database (unlike Kine/Mink)
2. **Version Chaining**: Records linked via `previous_id` for Watch support
3. **Background Compaction**: Removes old versions to prevent unbounded growth
4. **Field Indexing**: Dynamic columns for efficient field selectors

### Package Structure

| Package | Purpose |
| --------- | --------- |
| `pkg/db` | **Core**: Database layer with PostgreSQL/SQLite support |
| `pkg/server` | HTTP server wrapping k8s genericapiserver |
| `pkg/stores` | Builder pattern for store construction |
| `pkg/strategy` | CRUD+Watch strategy implementations |
| `pkg/types` | Core types (Object, ObjectList, Fields) |
| `pkg/apigroup` | API group registration |
| `pkg/authn` | Authentication (static token) |
| `pkg/serializer` | Object serialization |
| `pkg/validator` | Name validation |
| `pkg/otel` | OpenTelemetry attributes |

### Key Components

**pkg/db/Factory**

- Creates database connections and strategies
- Supports PostgreSQL and SQLite
- Manages schema migrations and compaction

**pkg/db/Strategy**

- Implements k8s storage interfaces
- Handles Create/Get/Update/Delete/List/Watch
- Uses embedded SQL (13 files in `pkg/db/statements/`)

**pkg/stores/Builder**

- Fluent interface for store construction
- 15+ variants: `complete`, `createget`, `getlist`, `readwritewatch`, etc.
- Example: `NewBuilder(scheme, obj).WithCompleteCRUD().WithWatch().Build()`

**pkg/server/Server**

- Wraps k8s.io/apiserver/genericapiserver
- Provides `Run(ctx)` and `Handler(ctx)` methods

### Database Schema

```sql
-- Main table (per resource type)
CREATE TABLE {name} (
    id          INTEGER PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    namespace   VARCHAR(255) NOT NULL,
    previous_id INTEGER UNIQUE,      -- Version chain for Watch
    uid         VARCHAR(255),
    created     INTEGER,             -- 1=create, NULL=update
    deleted     INTEGER DEFAULT 0,
    value       TEXT NOT NULL        -- JSON payload
);

-- Compaction tracking
CREATE TABLE compaction (
    name VARCHAR(255) UNIQUE,
    id   INTEGER                    -- Last compacted ID
);
```

## Code Patterns & Conventions

### Interface-Based Design

Strategies implement k8s apiserver storage interfaces:

```go
type Base interface {
    rest.Storage
    rest.Scoper
    rest.TableConvertor
}

type CompleteStrategy interface {
    Base
    rest.Creater
    rest.Getter
    rest.Lister
    rest.Updater
    rest.GracefulDeleter
    rest.Watcher
}
```

### Builder Pattern for Stores

```go
store := stores.NewBuilder(scheme, &MyObject{}).
    WithCompleteCRUD().
    WithWatch().
    Build()
```

### Error Handling

- Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
- Use k8s.io/apimachinery errors for API responses
- Custom DB errors in `pkg/db/errors/`

### Naming Conventions

- **Exported**: PascalCase (`Factory`, `Strategy`, `Builder`)
- **Unexported**: camelCase (`broadcast`, `objTemplate`)
- **Interfaces**: Descriptive names ending in `-er` or capability (`Watcher`, `CompleteStrategy`)

### SQL Files

Embedded SQL in `pkg/db/statements/` (13 files):

- `create.sql`, `get.sql`, `list.sql`, `update.sql`, `delete.sql`
- `watch.sql`, `compact.sql`, `migrate.sql`, etc.

## Testing Strategy

### Unit Tests (SQLite)

- Default test mode uses SQLite
- Fast, no external dependencies
- Located in `pkg/db/db_test.go`, `pkg/db/strategy_test.go`

### Integration Tests (PostgreSQL)

```bash
KINM_TEST_DB=postgres go test ./pkg/db/...
```

### Test Patterns

- Use `testify` for assertions
- Table-driven tests for multiple scenarios
- Test both SQLite and PostgreSQL paths

## Important Constraints

### Database Compatibility

- Primary: PostgreSQL 15+
- Development/Testing: SQLite (with sqlite-vec if needed)
- All SQL must work on both databases

### Kubernetes API Compatibility

- NOT a goal to maintain full k8s compatibility
- Uses k8s libraries for convenience, may fork away
- Focus on efficient CRUD+Watch semantics

### Performance Considerations

- Watch uses long-polling with broadcast notifications
- Compaction prevents unbounded table growth
- Field indexing for efficient field selectors

## Workspace Integration

This project is part of the AI workspace. Additional resources:

- **Claude Code commands**: `AI/.claude/commands/` (expert-mode, etc.)
- **Shared agents**: `AI/.claude/agents/` (explore, security-audit, etc.)
- **SuperClaude skills**: `/sc:analyze`, `/sc:test`, `/sc:git`, etc.
- **Serena memories**: `AI/.serena/memories/` (task_completion_checklist, etc.)
- **GitHub Actions**: Workspace-level PR review and issue triage

For session initialization with full context, run `/expert-mode` from the workspace root.

## Documentation

- **PROJECT_INDEX.md** - Quick reference index (this file summarized)
- **docs/API.md** - Complete API reference with interfaces
- **docs/ARCHITECTURE.md** - System architecture and design decisions
- **README.md** - Project overview and quick start

## Common Tasks

### Adding a New Resource Type

1. Define type in `pkg/types/` implementing `Object` interface
2. Create store using `stores.NewBuilder()`
3. Register with API group using `apigroup.ForStores()`
4. Add migrations if needed

### Adding Database Operations

1. Add SQL file to `pkg/db/statements/`
2. Embed using `//go:embed` directive
3. Add method to `Strategy` struct
4. Write tests for both SQLite and PostgreSQL

### Debugging Watch Issues

1. Check `previous_id` chain integrity
2. Verify compaction isn't removing needed records
3. Check broadcast channel for notification delivery
4. Enable GORM logging for SQL debugging

## K8s v0.35.0+ Compatibility

### Initial-Events-End Bookmark (CRITICAL)

client-go v0.35.0 introduced the **watch-list pattern** which requires kinm to send a special `initial-events-end` bookmark after streaming all initial events. Without this bookmark, `WaitForCacheSync` blocks forever.

**Implementation** (`pkg/db/strategy.go`):

1. Capture original `ResourceVersion` before modifications
2. Detect initial sync via `isInitialEventsEndBookmarkRequired()`
3. Send bookmark with `metav1.InitialEventsAnnotationKey: "true"` annotation

**Detection logic** - send bookmark when:

- `SendInitialEvents=true` (explicit), OR
- `AllowWatchBookmarks=true` AND original `ResourceVersion=""` or `"0"` (inferred)

The inferred case handles controller-runtime/nah which doesn't propagate `SendInitialEvents`.

**Reference**: https://github.com/kubernetes/kubernetes/issues/120348
