# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |

## Reporting a Vulnerability

**Please do NOT report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in this project, please report it responsibly:

1. **GitHub Security Advisories**: Use the "Report a vulnerability" button in the [Security tab](https://github.com/jrmatherly/kinm/security)
2. **Email**: Contact the repository maintainers directly

### What to Include

When reporting a vulnerability, please provide:

- Type of vulnerability (e.g., SQL injection, command injection, resource leak)
- Full path(s) of affected source file(s)
- Location of affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact assessment and potential consequences for downstream users

### Response Timeline

- **Initial Response**: Within 48 hours (2 business days)
- **Status Update**: Within 7 days
- **Resolution Target**: Within 30 days for critical issues (complex issues may take longer)

## Disclosure Policy

We follow coordinated disclosure practices:

- We will work with you to validate and remediate the vulnerability
- After a fix or mitigation is available, we'll publish release notes
- Security researchers who wish to be acknowledged will be credited in release notes

## Scope

Security issues that impact the confidentiality, integrity, or availability of this library or downstream projects using it are in scope.

**In scope:**
- SQL injection vulnerabilities in database operations
- Panics or crashes in library code that can be triggered by user input
- Resource leaks (goroutine leaks, memory leaks, database connection leaks)
- API design flaws that could lead to misuse causing security issues
- Unsafe concurrent access patterns
- Vulnerabilities in dependencies that affect kinm's functionality

**Out of scope (non-exhaustive):**
- Vulnerabilities requiring privileged database access without a clear escalation path
- Issues in user-implemented application logic (not in kinm library itself)
- Vulnerabilities in third-party dependencies not owned by this project (please report upstream)
- Theoretical issues without practical exploitation path

## Safe Harbor

We will not pursue legal action against security researchers conducting good-faith research aligned with this policy.

Please avoid:
- Privacy violations
- Service degradation or denial of service
- Data destruction or corruption
- Testing against databases or data you do not own

Only test against your own databases and data.

## Receiving Security Fixes

Security fixes are shipped in patch releases. We recommend:
- Always use the latest patch version of the library
- Monitor release notes for security advisories
- Enable GitHub notifications for this repository

We may issue public advisories (GitHub Security Advisory/CVE) when appropriate.

## Security Best Practices

This project follows security best practices including:

- **Code scanning**: GitHub CodeQL for static analysis
- **Dependency scanning**: Renovate for automated vulnerability detection
- **Secure defaults**: Safe configurations out of the box
- **Input validation**: Robust validation of inputs
- **Parameterized queries**: All SQL statements use parameterized queries to prevent injection
- **Error handling**: No panics in production code paths

## Library-Specific Security Considerations

### Database Security

kinm is a database-backed API server. Implementations using kinm should follow these security principles:

#### SQL Injection Prevention
- All SQL statements use parameterized queries
- No string concatenation for SQL construction
- Input validation before database operations

#### Connection Security
- Support for TLS/SSL database connections
- Connection string secrets should be stored securely
- Connection pooling with proper cleanup

#### Resource Management
- Bounded connection pools
- Proper cleanup of database resources
- Leak-free query execution

### API Server Security

#### Authentication & Authorization
- Integration with Kubernetes API server authentication
- Proper RBAC support when deployed in-cluster
- No elevation of privileges without explicit configuration

#### Input Validation
- Validate all Kubernetes resource inputs
- Sanitize data before storage
- Use typed APIs where possible

### Breaking Changes and API Stability

Security through API stability:
- Breaking changes are documented in release notes
- Major version bumps for breaking changes
- Deprecation notices before removal
- Migration guides for breaking changes

This helps downstream users maintain secure configurations during upgrades.

## Dependency Security

kinm depends on:
- k8s.io/apiserver (Kubernetes API server framework)
- k8s.io/apimachinery (Kubernetes types and utilities)
- gorm.io/gorm (ORM for database operations)
- github.com/jackc/pgx (PostgreSQL driver)
- OpenTelemetry (distributed tracing)

We monitor these dependencies for vulnerabilities using:
- Renovate automated scanning
- GitHub Dependabot alerts
- Manual security review of dependency updates

## Credits

With permission from reporters, we will credit security researchers in release notes and acknowledge their contributions to improving the security of this project.

Thank you for helping keep kinm and its users safe!
