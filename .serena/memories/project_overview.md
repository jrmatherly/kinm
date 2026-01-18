# Kinm Project Overview

## Purpose

Kinm (pronounced "kim") is a continuation of the Mink project. It provides an efficient and scalable API server that performs Kubernetes-like CRUD operations with Watch support, backed by a database (primarily PostgreSQL).

**Key difference from Kubernetes:** Compatibility with Kubernetes is not a goal - the focus is on efficiency and simplicity.

## Goals

- Efficient Postgres backend embracing SQL
- Keep all state in the database (not in memory like Kine/Mink)
- Kubernetes-like API semantics (CRUD + Watch) without the bloat

## Tech Stack

- **Language:** Go 1.24.0
- **Database:** PostgreSQL (primary), SQLite (for testing)
- **ORM:** GORM
- **API Framework:** k8s apiserver libraries
- **Logging:** Logrus
- **Testing:** testify (assert, require)
- **Tracing:** OpenTelemetry

## Repository

- **Module:** `github.com/obot-platform/kinm`
- **Main branch:** `main`

## Key Dependencies

- `k8s.io/apiserver` - API server framework
- `k8s.io/apimachinery` - Kubernetes types and utilities
- `gorm.io/gorm` and `gorm.io/driver/postgres` - ORM
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/sirupsen/logrus` - Logging
- `go.opentelemetry.io/otel` - Tracing
