# Branch Protection Rules

## Branch Naming Convention
```
<type>/<scope>/<short-description>
```

### Types
- `feat` - New feature
- `infra` - Infrastructure/platform
- `fix` - Bug fix
- `refactor` - Code refactoring
- `perf` - Performance improvement
- `test` - Tests
- `docs` - Documentation
- `chore` - Maintenance
- `security` - Security improvement

### Examples
```
feat/agents/tool-registry
infra/grpc/add-interceptors
fix/cache/redis-timeout
docs/platform/update-readme
```

## Protected Branches

### `main`
- Require pull request reviews (1 approver)
- Require status checks to pass
- Require branches to be up to date
- Enforce for administrators
- Require linear history

### `develop` (if using)
- Require pull request reviews (1 approver)
- Require status checks to pass

## Required Status Checks
- ✅ go-lint
- ✅ go-test
- ✅ python-lint
- ✅ python-test
- ✅ proto-lint