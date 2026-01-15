# Kinm API Reference

## Overview

Kinm provides a Kubernetes-like API for CRUD operations with Watch support. This document covers the core APIs and their usage patterns.

## Core Interfaces

### Strategy Interfaces

Located in `pkg/strategy/`, these interfaces define the contract for storage operations.

#### CompleteStrategy
```go
// pkg/strategy/all.go
type CompleteStrategy interface {
    Creater
    Updater
    StatusUpdater
    Getter
    Lister
    Deleter
    Watcher
    Destroy()
    Scheme() *runtime.Scheme
}
```

Full-featured strategy implementing all CRUD+Watch operations.

#### Creater
```go
// pkg/strategy/create.go
type Creater interface {
    Create(ctx context.Context, object types.Object) (types.Object, error)
    New() types.Object
}
```

#### Getter
```go
// pkg/strategy/get.go
type Getter interface {
    Get(ctx context.Context, namespace, name string) (types.Object, error)
}
```

#### Lister
```go
// pkg/strategy/list.go
type Lister interface {
    List(ctx context.Context, namespace string, opts storage.ListOptions) (types.ObjectList, error)
    New() types.Object
    NewList() types.ObjectList
}
```

#### Updater
```go
// pkg/strategy/update.go
type Updater interface {
    Update(ctx context.Context, obj types.Object) (types.Object, error)
    Get(ctx context.Context, namespace, name string) (types.Object, error)
    New() types.Object
}
```

#### Deleter
```go
// pkg/strategy/delete.go
type Deleter interface {
    Delete(ctx context.Context, namespace, name string, obj types.Object) (types.Object, error)
}
```

#### Watcher
```go
// pkg/strategy/watch.go
type Watcher interface {
    Watch(ctx context.Context, namespace string, opts storage.ListOptions) (<-chan watch.Event, error)
    New() types.Object
}
```

## Database Layer

### Factory

```go
// pkg/db/factory.go

// NewFactory creates a database connection factory
func NewFactory(schema *runtime.Scheme, dsn string) (*Factory, error)

// Supported DSN formats:
// - sqlite://path/to/db.sqlite
// - postgres://user:pass@host:port/dbname
// - postgresql://user:pass@host:port/dbname
```

**Factory Methods:**
- `Scheme() *runtime.Scheme` - Returns the registered scheme
- `Name() string` - Returns "kinm"
- `Check() error` - Performs database health check
- `NewDBStrategy(ctx, gvk, tableName) (*Strategy, error)` - Creates a storage strategy

### Strategy (db.Strategy)

```go
// pkg/db/strategy.go

// New creates a storage strategy for a specific resource type
func New(
    ctx context.Context,
    sqlDB *sql.DB,
    gvk schema.GroupVersionKind,
    scheme *runtime.Scheme,
    tableName string,
) (*Strategy, error)
```

**Strategy implements:**
- `Create(ctx, obj) (types.Object, error)`
- `Get(ctx, namespace, name) (types.Object, error)`
- `List(ctx, namespace, opts) (types.ObjectList, error)`
- `Update(ctx, obj) (types.Object, error)`
- `Delete(ctx, namespace, name, obj) (types.Object, error)`
- `Watch(ctx, namespace, opts) (<-chan watch.Event, error)`
- `Destroy()` - Cleanup and stop compaction

## Store Builder

### Builder Pattern

```go
// pkg/stores/builder.go

// NewBuilder creates a store builder
func NewBuilder(scheme *runtime.Scheme, obj kclient.Object) Builder

// Fluent methods:
builder.WithGet(strategy)           // Add Get capability
builder.WithList(strategy)          // Add List capability
builder.WithCreate(strategy)        // Add Create capability
builder.WithUpdate(strategy)        // Add Update capability
builder.WithDelete(strategy)        // Add Delete capability
builder.WithWatch(strategy)         // Add Watch capability
builder.WithCompleteCRUD(strategy)  // Add all CRUD operations
builder.WithTableConverter(conv)    // Custom table conversion
builder.WithValidateCreate(v)       // Create validation
builder.WithValidateUpdate(v)       // Update validation
builder.WithValidateDelete(v)       // Delete validation
builder.WithPrepareCreate(p)        // Pre-create hook
builder.WithPrepareUpdate(p)        // Pre-update hook
builder.Build() rest.Storage        // Build final store
```

### Pre-built Store Types

| Store Type | Capabilities |
| ------------ | -------------- |
| `complete` | Create, Get, List, Update, Delete, Watch |
| `createget` | Create, Get |
| `createonly` | Create |
| `getlist` | Get, List |
| `getonly` | Get |
| `listonly` | List |
| `listwatch` | List, Watch |
| `getlistwatch` | Get, List, Watch |
| `getlistdelete` | Get, List, Delete |
| `getlistupdatedelete` | Get, List, Update, Delete |
| `getlistupdatedeletewatch` | Get, List, Update, Delete, Watch |
| `readdelete` | Get, List, Delete |
| `readwritewatch` | Get, List, Create, Update, Delete, Watch |
| `status` | Status subresource operations |

## Server Configuration

### Config

```go
// pkg/server/server.go

type Config struct {
    Name                  string                         // Server name
    Version               string                         // API version
    Authenticator         authenticator.Request          // Request authenticator
    Authorization         authorizer.Authorizer          // Request authorizer
    HTTPListenPort        int                            // HTTP port (default: 8080)
    HTTPSListenPort       int                            // HTTPS port (default: 8081)
    Listener              net.Listener                   // Custom listener
    LongRunningVerbs      []string                       // Long-running verbs
    LongRunningResources  []string                       // Long-running resources
    OpenAPIConfig         openapicommon.GetOpenAPIDefinitions
    Scheme                *runtime.Scheme                // Type scheme
    CodecFactory          *serializer.CodecFactory       // Codec factory
    APIGroups             []*server.APIGroupInfo         // API groups to install
    Middleware            []func(http.Handler) http.Handler
    PostStartFunc         server.PostStartHookFunc       // Post-start hook
    SupportAPIAggregation bool                           // K8s aggregation support
    DefaultOptions        *options.RecommendedOptions
    AuditConfig           *options.AuditOptions
    IgnoreStartFailure    bool                           // Continue on failure
    ReadinessCheckers     []healthz.HealthChecker
}
```

### Server Creation

```go
// Create server
server, err := server.New(&server.Config{
    Name:    "myapi",
    Version: "v1",
    Scheme:  myScheme,
    APIGroups: []*server.APIGroupInfo{...},
})

// Run server
err = server.Run(ctx)

// Or get handler for custom server
handler := server.Handler(ctx)
```

## Type System

### Object Interface

```go
// pkg/types/object.go

type Object interface {
    runtime.Object
    metav1.Object
}

type ObjectList interface {
    runtime.Object
    metav1.ListInterface
}

type NamespaceScoper interface {
    NamespaceScoped() bool
}
```

### Field Selection

```go
// pkg/types/fields.go

type Fields interface {
    FieldNames() []string
}

type FieldsIndexer interface {
    IndexFields() []string
}
```

## Authentication

### Static Token Auth

```go
// pkg/authn/statictoken.go

// Create authenticator
auth := authn.NewStaticToken(token, userName, groups)

// Extract token from request
token := authn.GetBearerToken(request)
```

## OpenTelemetry Integration

```go
// pkg/otel/attributes.go

// Convert list options to trace attributes
attrs := otel.ListOptionsToAttributes(opts)

// Convert object to trace attributes
attrs := otel.ObjectToAttributes(obj)
```

## Usage Example

```go
package main

import (
    "context"

    "github.com/obot-platform/kinm/pkg/db"
    "github.com/obot-platform/kinm/pkg/server"
    "github.com/obot-platform/kinm/pkg/stores"
    "github.com/obot-platform/kinm/pkg/apigroup"
)

func main() {
    ctx := context.Background()

    // 1. Create scheme and register types
    scheme := runtime.NewScheme()
    // ... register types

    // 2. Create database factory
    factory, err := db.NewFactory(scheme, "postgres://...")
    if err != nil {
        panic(err)
    }

    // 3. Create storage strategies
    myStrategy, err := factory.NewDBStrategy(ctx, myGVK, "mytable")
    if err != nil {
        panic(err)
    }

    // 4. Build stores
    store := stores.NewBuilder(scheme, &MyType{}).
        WithCompleteCRUD(myStrategy).
        Build()

    // 5. Create API group
    apiGroup := apigroup.ForStores(scheme, map[string]rest.Storage{
        "myresources": store,
    })

    // 6. Create and run server
    srv, err := server.New(&server.Config{
        Name:      "myapi",
        Scheme:    scheme,
        APIGroups: []*server.APIGroupInfo{apiGroup},
    })
    if err != nil {
        panic(err)
    }

    srv.Run(ctx)
}
```
