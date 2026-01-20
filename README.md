# Kinm

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)

**A Kubernetes-like API server backed by PostgreSQL/SQLite**

> 📖 **Documentation:** [API Reference](docs/API.md) • [Architecture](docs/ARCHITECTURE.md)

---

## What is Kinm?

Kinm (pronounced "kim", like the name) is a database-backed API server providing Kubernetes-like CRUD+Watch semantics without the complexity of etcd. It enables building scalable, resource-oriented APIs with familiar K8s patterns while keeping all state in PostgreSQL or SQLite.

### Origin Story

Kinm continues the learnings from **Mink** (Mink is not Kubernetes), a Kubernetes Aggregated API Server backed by a database. While Mink was archived when Acorn Labs pivoted away from Kubernetes products, the core idea remained valuable: provide K8s-like API capabilities without being tightly coupled to Kubernetes itself.

Kinm takes this further by:
- **Embracing SQL**: Full PostgreSQL support with efficient query patterns
- **Database-first design**: All state in the database, no in-memory caching
- **Independence from K8s**: Compatibility is not a goal (though it happens to work currently)
- **Modern architecture**: Removing K8s library bloat and version churn

---

## 🎯 Key Features

### 🔄 Kubernetes-like API

- **Complete CRUD Operations** - Create, Read, Update, Delete with K8s semantics
- **Watch Support** - Real-time change notifications via long-polling and broadcast
- **Field Selectors** - Efficient filtering using indexed fields
- **Resource Versioning** - Version chains via `previous_id` for consistency

### 🗄️ Database-Backed Storage

- **PostgreSQL** - Primary production database with advanced SQL features
- **SQLite** - Development and testing support
- **No etcd** - Pure SQL design without etcd dependencies
- **Stateless API Servers** - All state in database enables horizontal scaling

### ⚡ Performance & Scalability

- **Background Compaction** - Automatic cleanup of old resource versions
- **Efficient Watch** - Long-polling with broadcast notifications
- **Field Indexing** - Dynamic columns for fast field selector queries
- **Zero In-Memory State** - Horizontal scaling without sticky sessions

### 🔧 Developer Experience

- **Builder Pattern** - Fluent API for configuring stores (15+ pre-built variants)
- **Strategy Interfaces** - Clean separation of concerns (Create, Get, List, Update, Delete, Watch)
- **Embedded SQL** - Parameterized queries in `.sql` files
- **Testing Support** - SQLite for fast unit tests, PostgreSQL for integration tests

---

## 🚀 Quick Start

### Prerequisites

```bash
# Check requirements
go version       # 1.25+ required
git --version
make --version   # Optional, for development tasks

# Database (choose one or both)
# PostgreSQL 12+ (recommended for production)
psql --version

# SQLite 3.35+ (built-in for development/testing)
sqlite3 --version
```

**Minimum Requirements:**
- **Go:** 1.25 or later (for Go modules and generics support)
- **Database:** PostgreSQL 12+ (production) or SQLite 3.35+ (development)
- **Git:** Any recent version (for dependency resolution)
- **Make:** Optional, for running development tasks (`make build`, `make test`, etc.)

**Optional Tools:**
- **golangci-lint:** For code quality checks (`make lint`)
- **Docker/Podman:** For running PostgreSQL in containers
- **psql:** PostgreSQL client for database management

### Installation

**As a Library Dependency:**

```bash
# Add to your Go project
go get github.com/obot-platform/kinm

# Or add to go.mod
require github.com/obot-platform/kinm v0.0.0
```

**For Development:**

```bash
# Clone the repository
git clone https://github.com/obot-platform/kinm.git
cd kinm

# Install dependencies and build
go mod download
make build

# Run tests
make test
```

### Basic Usage

Here's a complete example showing how to create a Kubernetes-like API server with CRUD operations:

#### 1. Define Your Resource Type

```go
package main

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

// Widget is a simple custom resource
type Widget struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   WidgetSpec   `json:"spec,omitempty"`
    Status WidgetStatus `json:"status,omitempty"`
}

type WidgetSpec struct {
    Color string `json:"color,omitempty"`
    Size  int    `json:"size,omitempty"`
}

type WidgetStatus struct {
    Phase string `json:"phase,omitempty"`
}

// WidgetList contains a list of Widgets
type WidgetList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Widget `json:"items"`
}

// Required methods for runtime.Object
func (w *Widget) DeepCopyObject() runtime.Object {
    // Implementation omitted for brevity
    return w
}

func (wl *WidgetList) DeepCopyObject() runtime.Object {
    // Implementation omitted for brevity
    return wl
}
```

#### 2. Set Up the Server

```go
package main

import (
    "context"
    "log"

    "github.com/obot-platform/kinm/pkg/apigroup"
    "github.com/obot-platform/kinm/pkg/db"
    "github.com/obot-platform/kinm/pkg/server"
    "github.com/obot-platform/kinm/pkg/stores"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/apiserver/pkg/registry/rest"
)

func main() {
    ctx := context.Background()

    // 1. Create scheme and register types
    scheme := runtime.NewScheme()
    gvk := schema.GroupVersionKind{
        Group:   "example.com",
        Version: "v1",
        Kind:    "Widget",
    }
    scheme.AddKnownTypes(gvk.GroupVersion(), &Widget{}, &WidgetList{})

    // 2. Create database factory (SQLite for development)
    factory, err := db.NewFactory(scheme, "sqlite://widgets.db")
    if err != nil {
        log.Fatal(err)
    }

    // 3. Create storage strategy for Widgets
    strategy, err := factory.NewDBStrategy(ctx, gvk, "widgets")
    if err != nil {
        log.Fatal(err)
    }
    defer strategy.Destroy()

    // 4. Build REST storage with complete CRUD+Watch capabilities
    store := stores.NewBuilder(scheme, &Widget{}).
        WithCompleteCRUD(strategy).
        Build()

    // 5. Create API group with the store
    apiGroup := apigroup.ForStores(scheme, map[string]rest.Storage{
        "widgets": store,
    })

    // 6. Configure and start server
    srv, err := server.New(&server.Config{
        Name:           "widget-api",
        Version:        "v1",
        Scheme:         scheme,
        APIGroups:      []*server.APIGroupInfo{apiGroup},
        HTTPListenPort: 8080,
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Widget API server running on :8080")
    if err := srv.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
```

#### 3. Perform CRUD Operations

Once the server is running, interact with it using kubectl or HTTP clients:

**Create a Widget:**

```bash
# Using curl
curl -X POST http://localhost:8080/apis/example.com/v1/namespaces/default/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion": "example.com/v1",
    "kind": "Widget",
    "metadata": {
      "name": "my-widget",
      "namespace": "default"
    },
    "spec": {
      "color": "blue",
      "size": 42
    }
  }'
```

**Get a Widget:**

```bash
curl http://localhost:8080/apis/example.com/v1/namespaces/default/widgets/my-widget
```

**List All Widgets:**

```bash
curl http://localhost:8080/apis/example.com/v1/namespaces/default/widgets
```

**Update a Widget:**

```bash
curl -X PUT http://localhost:8080/apis/example.com/v1/namespaces/default/widgets/my-widget \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion": "example.com/v1",
    "kind": "Widget",
    "metadata": {
      "name": "my-widget",
      "namespace": "default",
      "resourceVersion": "1"
    },
    "spec": {
      "color": "red",
      "size": 100
    }
  }'
```

**Delete a Widget:**

```bash
curl -X DELETE http://localhost:8080/apis/example.com/v1/namespaces/default/widgets/my-widget
```

**Watch for Changes:**

```bash
# Long-polling watch stream
curl http://localhost:8080/apis/example.com/v1/namespaces/default/widgets?watch=true
```

#### 4. Using with kubectl (Optional)

Kinm servers are compatible with kubectl (though compatibility is not a design goal):

```bash
# Configure kubectl
kubectl config set-cluster widget-api --server=http://localhost:8080
kubectl config set-context widget-api --cluster=widget-api
kubectl config use-context widget-api

# Use kubectl commands
kubectl get widgets -n default
kubectl describe widget my-widget -n default
kubectl delete widget my-widget -n default
```

**Next Steps:**

- See [docs/API.md](docs/API.md) for comprehensive API documentation
- Explore builder patterns in [docs/API.md#store-builder](docs/API.md#store-builder)
- Learn about field selectors, validation hooks, and custom table conversion
- Review [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for database schema and design decisions

---

## 📚 Documentation

| Document | Description |
| ---------- | ------------- |
| **[API Reference](docs/API.md)** | Complete API documentation with examples |
| **[Architecture](docs/ARCHITECTURE.md)** | Database schema, data flow, and design decisions |
| **[CLAUDE.md](CLAUDE.md)** | Development guide for Claude Code AI assistant |

---

## 🏗️ Architecture

Kinm uses a layered architecture:

```
HTTP Clients
    ↓
pkg/server (HTTP Server + k8s GenericAPIServer)
    ↓
pkg/stores (Builder + 15+ store variants)
    ↓
pkg/strategy (CRUD+Watch adapters)
    ↓
pkg/db (Factory + Strategy + SQL statements)
    ↓
PostgreSQL / SQLite
```

**Key architectural decisions:**

- **Database-first**: All state persisted in SQL, no in-memory caching
- **Version chains**: Resources linked via `previous_id` for Watch support
- **Background compaction**: Prevents unbounded growth of version history
- **Long-polling Watch**: Efficient change notifications without WebSockets

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed diagrams and explanations.

---

## 🔧 Development

### Building

```bash
make build            # Build the project
go build ./...        # Direct build
```

### Testing

```bash
# Unit tests (SQLite - fast)
make test             # Run all tests with race detector
make test-short       # Run tests in short mode

# Integration tests (PostgreSQL - comprehensive)
make test-integration # Requires PostgreSQL (KINM_TEST_DB=postgres)

# Coverage
make test-coverage    # Generate coverage.html
```

### Code Quality

```bash
make lint             # Run golangci-lint
make validate         # Run all validation (lint + tests)
```

---

## 🤝 Contributing

Contributions are welcome! Please follow the existing code patterns:

- **Interface-based design** - Use strategy interfaces for extensibility
- **Builder pattern** - Fluent configuration for stores
- **Error wrapping** - Use `fmt.Errorf("context: %w", err)` for error chains
- **SQL-first** - Embrace database capabilities, avoid in-memory caching

See [CLAUDE.md](CLAUDE.md) for detailed development guidelines and patterns.

---

## 📄 License

Apache 2.0 - See [LICENSE](LICENSE) for details

---

## 🔗 Related Projects

Part of the [AI/MCP Multi-Repo Workspace](https://github.com/obot-platform):

- **[nah](https://github.com/obot-platform/nah)** - Kubernetes controller framework (uses kinm concepts)
- **[obot-entraid](https://github.com/obot-platform/obot-entraid)** - MCP platform with custom auth

---

*Kinm is not Mink. Kinm embraces SQL and modern Go patterns while providing familiar Kubernetes-like API semantics.*
