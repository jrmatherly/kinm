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

See [docs/API.md](docs/API.md) for comprehensive examples and usage patterns.

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
