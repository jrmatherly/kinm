## Summary

<!-- Brief description of changes - what does this PR do? -->

## Type of Change

<!-- Mark relevant options with 'x' -->

- [ ] Bug fix (non-breaking change fixing an issue)
- [ ] New feature (non-breaking change adding functionality)
- [ ] Breaking change (fix or feature causing existing functionality to break)
- [ ] Documentation update
- [ ] Configuration/dependency update
- [ ] Performance improvement
- [ ] Code refactoring
- [ ] Test addition/modification
- [ ] CI/CD or infrastructure change

## Related Issues

<!-- Link related issues using keywords: Closes #123, Fixes #456, Relates to #789 -->

Closes #
Relates to #

## Changes Made

<!-- List specific changes in bullet points -->

-
-
-

## Testing

### Local Testing

- [ ] Tests pass (`go test ./...`)
- [ ] Linting passes (`golangci-lint run`)
- [ ] Build succeeds (`go build ./...`)
- [ ] New tests added for new functionality
- [ ] Manual testing performed

### Test Environment

- OS: <!-- e.g., macOS 14, Ubuntu 22.04 -->
- Go Version: <!-- e.g., 1.25 -->
- PostgreSQL Version: <!-- e.g., 17 (if integration testing) -->

### Test Results

<details>
<summary>Click to expand test output</summary>

```bash
# Paste relevant test output here
```

</details>

## Code Quality

### General

- [ ] Code follows project conventions
- [ ] Self-reviewed all changes
- [ ] Comments added for complex logic
- [ ] No debug statements left in code
- [ ] No sensitive data (tokens, passwords, secrets) in code or commits

### Go-Specific

- [ ] `gofmt` applied (enforced by CI)
- [ ] `golangci-lint` passes (enforced by CI)
- [ ] Exported functions/types have GoDoc comments
- [ ] Proper error handling (no ignored errors)
- [ ] No panics in library code (return errors instead)
- [ ] Thread-safe concurrent access (verified with `go test -race` if applicable)

## Database Changes

<!-- Complete if changes affect database operations -->

- [ ] No database changes
- [ ] SQL statements modified (complete below)

**If SQL changes:**

- [ ] Parameterized queries used (no string concatenation)
- [ ] SQL injection prevention verified
- [ ] PostgreSQL compatibility verified
- [ ] SQLite compatibility verified (for testing)
- [ ] Proper error handling for database operations
- [ ] Connection/resource cleanup verified

## Library Quality

### API Design

- [ ] Exported APIs have clear, descriptive names
- [ ] Public APIs documented with GoDoc comments
- [ ] API changes follow Go best practices
- [ ] No unnecessary exports (keep internal APIs private)

### Backward Compatibility

- [ ] Changes are backward compatible
- [ ] Breaking changes documented in release notes format below
- [ ] Migration guide provided (if breaking change)
- [ ] Deprecated APIs marked with deprecation comments

## Documentation

- [ ] README updated (if needed)
- [ ] API documentation updated (if applicable)
- [ ] Code comments added for public APIs
- [ ] Migration guide provided (if breaking change)

## Breaking Changes

- [ ] No breaking changes

**If breaking changes exist:**

**Description:**
<!-- Describe what breaks and why -->

**Migration Path:**
<!-- Step-by-step guide for users to migrate -->

1.
2.
3.

**Affected APIs:**
<!-- List affected public APIs -->

-
-

## Performance Impact

- [ ] No performance impact expected
- [ ] Performance improved (describe below)
- [ ] May impact performance (explain and justify)

**Details:**
<!-- If performance changed, provide benchmarks or analysis -->

```go
// Example: Benchmark results showing improvement
// BenchmarkOld-8   1000  1000000 ns/op
// BenchmarkNew-8   2000   500000 ns/op  (50% improvement)
```

## Security Considerations

- [ ] No security implications
- [ ] SQL injection prevention verified
- [ ] Input validation added where needed
- [ ] No panics that could be triggered by user input
- [ ] No resource leaks (goroutines, memory, database connections)
- [ ] Thread-safe concurrent access
- [ ] Dependencies don't introduce known vulnerabilities

## Upstream Compatibility

<!-- Important for maintaining the fork -->

- [ ] Changes are in fork-specific code only
- [ ] Changes touch shared upstream code (may conflict with future syncs)

**If shared code changed:**

- Conflict resolution strategy: <!-- describe how to handle future upstream conflicts -->

## Additional Context

<!-- Any other information reviewers should know -->

## Post-Merge Actions

<!-- Actions needed after merging (if any) -->

- [ ] No post-merge actions required
- [ ] Update examples in downstream projects (obot-entraid, nah)
- [ ] Announce breaking changes to users
- [ ] Update migration documentation

---

## Reviewer Checklist

<!-- For maintainers reviewing this PR -->

- [ ] Code quality meets standards
- [ ] Tests comprehensive and passing
- [ ] Documentation clear and complete
- [ ] Security considerations addressed
- [ ] Breaking changes properly handled
- [ ] Performance acceptable
- [ ] API design follows Go best practices
- [ ] Backward compatibility maintained (or properly documented)
- [ ] Database operations are safe and efficient

---

**Thank you for contributing to kinm!**

By submitting this PR, you agree to license your contribution under the Apache 2.0 license.
