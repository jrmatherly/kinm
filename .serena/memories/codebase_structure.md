# Kinm Codebase Structure

## Directory Layout

```
pkg/
├── apigroup/          # API group registration and management
├── authn/             # Authentication (static token auth)
├── db/                # Database layer (core storage implementation)
│   ├── errors/        # Custom database errors
│   ├── glogrus/       # GORM-logrus integration
│   └── statements/    # SQL statements (embedded .sql files)
├── otel/              # OpenTelemetry tracing attributes
├── serializer/        # Object serialization
├── server/            # HTTP server configuration and setup
├── stores/            # Store interfaces and builders (various CRUD combinations)
├── strategy/          # Storage strategies (CRUD operations)
│   ├── remote/        # Remote strategy implementation
│   └── translation/   # Translation layer
├── types/             # Type definitions (Object, Fields, Attr)
└── validator/         # Input validation (name validation)
```

## Key Components

### pkg/server

- `Server` - Main API server wrapping k8s genericapiserver
- `Config` - Server configuration with auth, ports, middleware

### pkg/db

- Core database operations (get, list, insert, delete, compact)
- SQL statements loaded from embedded .sql files
- Supports PostgreSQL and SQLite

### pkg/strategy

- Implements k8s apiserver storage strategies
- Operations: create, get, list, update, delete, watch
- `Base` interface combines Storage, Scoper, TableConvertor, SingularNameProvider

### pkg/stores

- Pre-built store configurations with different capability combinations
- Examples: `createonly`, `getlist`, `complete`, `readwritewatch`
- Uses builder pattern for store construction
